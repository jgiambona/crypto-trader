package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
)

type (
	// Period holds the data returned from Poloniex historical data.
	Period struct {
		Date            int64   `json:"date"`
		High            float64 `json:"high"`
		Low             float64 `json:"low"`
		Open            float64 `json:"close"`
		Volume          float64 `json:"volume"`
		QuoteVolume     float64 `json:"quoteVolume"`
		WeightedAverage float64 `json:"weightedAverage"`
	}
)

func pollTicker() {
	var waitExchanges sync.WaitGroup

	bot.running = false
	insertBotStatus("OFF")

	currencyPair := "NOX/ETH"

	var placedOrder int64 = -1
	var botVolume float64 = 0

	for {
		waitExchanges.Add(1)

		if !bot.running {
			botVolume = 0
		}

		go func() {
			defer waitExchanges.Done()

		repeatCheckLowestBid:
			if bot.running {
				if len(bot.accountOne.APIKey) > 0 && len(bot.accountTwo.APIKey) > 0 {

					// In each loop reset this values.
					tradePlace := false
					tradePrice := -1.0
					tradeQuantity := -1.0

					if bot.ruleOne.Enabled {

						// Update ticker to get price
						p := updateTicker(currencyPair)
						lowest := p["ask"].(float64)
						volume := p["volume"].(float64)

						// Calculate the target price and quantity to trade
						v := bot.ruleOne.TransactionVolume * (bot.ruleOne.VarianceOfTransaction / 100.0)
						quantity := bot.ruleOne.TransactionVolume + getRandom(v)
						targetPrice := lowest - bot.ruleOne.BidPriceStepDown

						log.Printf("--- %.8f %0.8f", lowest, volume)
						log.Printf("--- %.8f %0.8f = %0.8f", targetPrice, quantity, targetPrice*quantity)

						// Switch seller roles if quantity is lower than can be traded.
						if err := switchAccountRolesSeller(quantity); err != nil {
							log.Print(err)
							return
						}

						// Switch buyer roles if total amount is lower than can be traded.
						if err := switchAccountRolesBuyer(targetPrice * quantity); err != nil {
							log.Print(err)
							return
						}

						// Check if the target price aligns and above minimum bid requirements
						if targetPrice >= bot.ruleOne.MinimumBid {

							// Check if the volume generated by bot is already exceeding the
							// maximum allowable volume per run set.
							if botVolume < bot.ruleOne.MaximumVolume {
								botVolume += quantity

								remarks := bot.accountOne.APIKey

								if !bot.simulate {
									o, err := sellLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, targetPrice, quantity)
									if err != nil {
										log.Print("error occurred in creating sell order ", err)
										remarks = fmt.Sprintf("Error on %s", bot.accountOne.APIKey)
										placedOrder = -1
									} else {
										placedOrder = o.OrderID
									}
								} else {
									placedOrder = 1
								}

								// Save target price and quantity
								tradePrice = targetPrice
								tradeQuantity = quantity

								// Record sell transaction
								insertTransaction("SELL", "nox_eth", tradePrice, tradeQuantity,
									strconv.FormatBool(bot.simulate), remarks)
							}

							{
								o, err := getOrderBook(currencyPair)
								if err != nil {
									log.Print(err)
								}

								qs := big.NewFloat(tradeQuantity).SetMode(big.AwayFromZero).Text('f', 7)
								qr := o.Asks[0][0]
								tp := big.NewFloat(tradePrice).SetMode(big.AwayFromZero).Text('f', 7)
								tc := o.Asks[0][1]

								log.Print("-- ", qs)
								log.Print("-- ", qr)
								log.Print("-- ", tp)
								log.Print("-- ", tc)
								if qr != qs || tc != tp {
									remarks := bot.accountOne.APIKey

									c, err := cancelLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, placedOrder)

									if err != nil {
										log.Print("error occurred in cancelling order")
										remarks = fmt.Sprintf("Error on %s", bot.accountOne.APIKey)
									}

									if !c.Success {
										log.Print("unable to cancel order")
									}

									// Record cancel order
									insertTransaction("CANCEL", "nox_eth", tradePrice, tradeQuantity,
										strconv.FormatBool(bot.simulate), remarks)
									placedOrder = -1

									// Repeat
									goto repeatCheckLowestBid
								}
							}

							if !bot.simulate && placedOrder > 1 {
								log.Print("-- ", placedOrder)
								c, err := getOrder(bot.accountOne.APIKey, bot.accountOne.APISecret,
									strconv.FormatInt(placedOrder, 10))
								if err != nil {
									log.Print("error occurred in cancelling order")
								}

								// Cancel if not aligned in target price and quantity
								qs := big.NewFloat(tradeQuantity).SetMode(big.AwayFromZero).Text('f', 7)
								qr := big.NewFloat(c.RemainingQuantity).SetMode(big.AwayFromZero).Text('f', 7)
								tp := big.NewFloat(tradePrice).SetMode(big.AwayFromZero).Text('f', 7)
								tc := big.NewFloat(c.Price).SetMode(big.AwayFromZero).Text('f', 7)

								log.Print("-- ", qs)
								log.Print("-- ", qr)
								log.Print("-- ", tp)
								log.Print("-- ", tc)
								if qr != qs || tc != tp {
									remarks := bot.accountOne.APIKey

									c, err := cancelLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, placedOrder)

									if err != nil {
										log.Print("error occurred in cancelling order")
										remarks = fmt.Sprintf("Error on %s", bot.accountOne.APIKey)
									}

									if !c.Success {
										log.Print("unable to cancel order")
									}

									// Record cancel order
									insertTransaction("CANCEL", "nox_eth", tradePrice, tradeQuantity,
										strconv.FormatBool(bot.simulate), remarks)
									placedOrder = -1

									// Repeat
									goto repeatCheckLowestBid
								}
							}

							// Check if the lowest is the trade price and check
							if lowest >= tradePrice && tradePrice > 0 && placedOrder > 0 {
								remarks := bot.accountTwo.APIKey

								if !bot.simulate {
									_, err := buyLimit(bot.accountTwo.APIKey, bot.accountTwo.APISecret,
										currencyPair, tradePrice, tradeQuantity)
									if err != nil {
										log.Print("error occurred in creating buy order ", err)
										remarks = fmt.Sprintf("Error on %s", bot.accountTwo.APIKey)
									}
								}

								// Record buy transaction
								insertTransaction("BUY", "nox_eth", tradePrice, tradeQuantity,
									strconv.FormatBool(bot.simulate), remarks)

								// Put a trade place flag and removed placed order ID.
								tradePlace = true
								placedOrder = -1
							}
						}
					}

					// Check if rule two and three is enabled and if there is already been
					// a trade place then skip.
					if bot.ruleTwo.Enabled && bot.ruleThree.Enabled && !tradePlace {

						// Update the price to get the latest update
						p := updateTicker(currencyPair)
						lowest := p["ask"].(float64)
						volume := p["volume"].(float64)

						// Calculate the trade price and the amount to be trade
						v := bot.ruleTwo.TransactionVolume * (bot.ruleTwo.VarianceOfTransaction / 100.0)
						quantity := bot.ruleTwo.TransactionVolume + getRandom(v)
						targetPrice := lowest - bot.ruleTwo.BidPriceStepDown

						// Log for easy debugging
						log.Printf("--- %.8f %0.8f", lowest, volume)
						log.Printf("--- %.8f %0.8f = %0.8f", targetPrice, quantity)

						// Switch roles if the seller has no amount specified.
						if err := switchAccountRolesSeller(quantity); err != nil {
							log.Print(err)
							return
						}

						// Check buyer if it can buy the specified amount.
						if err := switchAccountRolesBuyer(targetPrice * quantity); err != nil {
							log.Print(err)
							return
						}

						// Check if target price aligns the minimumBid and if no order has
						// been placed yet.
						if targetPrice >= bot.ruleTwo.MinimumBid && placedOrder < 0 {

							// Check if the bot volume generated is over the set maximum volume
							// otherwise stop.
							if botVolume < bot.ruleTwo.MaximumVolume {
								botVolume += quantity

								remarks := bot.accountOne.APIKey

								// Check if its in simulation mode otherwise mock placed order.
								if !bot.simulate {
									o, err := sellLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, targetPrice, quantity)
									if err != nil {
										log.Print("error occurred in creating sell order")
										remarks = fmt.Sprintf("Error on %s", bot.accountOne.APIKey)
										placedOrder = -1
									}
									placedOrder = o.OrderID
								} else {
									placedOrder = 1
								}

								// Save trade price and quantity
								tradePrice = targetPrice
								tradeQuantity = quantity

								// Record the sell transaction
								insertTransaction("SELL", "nox_eth", tradePrice, tradeQuantity,
									strconv.FormatBool(bot.simulate), remarks)
							}

							// Check if the lowest is not the trade price
							if lowest != tradePrice && tradePrice > 0 && placedOrder > 0 {
								remarks := bot.accountOne.APIKey

								// Check if bot is in simulation mode
								if !bot.simulate {

									c, err := cancelLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, placedOrder)

									if err != nil {
										log.Print("error occurred in cancelling order")
										remarks = fmt.Sprintf("Error on %s", bot.accountOne.APIKey)
									}

									if !c.Success {
										log.Print("unable to cancel order")
									}
								}

								// Record cancel order
								insertTransaction("CANCEL", "nox_eth", tradePrice, tradeQuantity,
									strconv.FormatBool(bot.simulate), remarks)
								placedOrder = -1

								// Repeat
								goto repeatCheckLowestBid
							}

							if !bot.simulate && placedOrder > 1 {
								log.Print("-- ", placedOrder)
								c, err := getOrder(bot.accountOne.APIKey, bot.accountOne.APISecret,
									strconv.FormatInt(placedOrder, 10))
								if err != nil {
									log.Print("error occurred in cancelling order")
								}

								// Cancel if not aligned in target price and quantity
								qs := big.NewFloat(tradeQuantity).SetMode(big.AwayFromZero).Text('f', 7)
								qr := big.NewFloat(c.RemainingQuantity).SetMode(big.AwayFromZero).Text('f', 7)
								tp := big.NewFloat(tradePrice).SetMode(big.AwayFromZero).Text('f', 7)
								tc := big.NewFloat(c.Price).SetMode(big.AwayFromZero).Text('f', 7)

								log.Print("-- ", qs)
								log.Print("-- ", qr)
								log.Print("-- ", tp)
								log.Print("-- ", tc)
								if qr != qs || tc != tp {
									remarks := bot.accountOne.APIKey

									c, err := cancelLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, placedOrder)

									if err != nil {
										log.Print("error occurred in cancelling order")
										remarks = fmt.Sprintf("Error on %s", bot.accountOne.APIKey)
									}

									if !c.Success {
										log.Print("unable to cancel order")
									}

									// Record cancel order
									insertTransaction("CANCEL", "nox_eth", tradePrice, tradeQuantity,
										strconv.FormatBool(bot.simulate), remarks)
									placedOrder = -1

									// Repeat
									goto repeatCheckLowestBid
								}
							}

							// Check if the lowest price comes from accounte one then
							// execute trade
							if lowest >= tradePrice && tradePrice > 0 && placedOrder > 0 {

								remarks := bot.accountTwo.APIKey

								// Check if bot is in simulation mode
								if !bot.simulate {
									_, err := buyLimit(bot.accountTwo.APIKey, bot.accountTwo.APISecret,
										currencyPair, tradePrice, tradeQuantity)
									if err != nil {
										log.Print("error occurred in creating buy order ", err)
										remarks = fmt.Sprintf("Error on %s", bot.accountTwo.APIKey)
									}
								}

								// Record transaction
								insertTransaction("BUY", "nox_eth", tradePrice, tradeQuantity,
									strconv.FormatBool(bot.simulate), remarks)

								// Put a trade place flag and removed placed order ID.
								tradePlace = true
								placedOrder = -1
							}
						}
					}
				}
			}
		}()

		waitExchanges.Wait()
		interval := time.Duration(bot.ruleOne.Interval)
		time.Sleep(interval)
	}
}

func switchAccountRolesSeller(quantity float64) (err error) {
	count := 0

switchAccount:
	swap := false

	b, err := getBalance(bot.accountOne.APIKey, bot.accountOne.APISecret, "NOX")
	if err != nil {
		log.Print("unable to get balance", err)
	}

	if err != nil || b.Type == "" {
		bot.running = false
		insertBotStatus("OFF")

		bot.simulate = true
		insertBotSimulateStatus("ON")

		return errors.New("-- error unable to get balance seller")
	}

	t1 := quantity
	if b.Value < t1 {
		log.Print("-- commence switch roles seller")
		swap = true

		// increment counter to check if both accounts is invalid
		count += 1
	}

	if swap {
		tmp := bot.accountOne
		bot.accountOne = bot.accountTwo
		bot.accountTwo = tmp

		if count == 1 {
			goto switchAccount
		} else {
			bot.running = false
			insertBotStatus("OFF")

			bot.simulate = true
			insertBotSimulateStatus("ON")

			return errors.New("-- both account can't sell NOX")
		}
	}

	time.Sleep(1 * time.Second)
	return nil
}

func switchAccountRolesBuyer(price float64) (err error) {
	count := 0

switchAccount:
	swap := false

	b, err := getBalance(bot.accountTwo.APIKey, bot.accountTwo.APISecret, "ETH")
	if err != nil {
		log.Print("unable to get balance", err)
	}

	if err != nil || b.Type == "" {
		bot.running = false
		insertBotStatus("OFF")

		bot.simulate = true
		insertBotSimulateStatus("ON")

		return errors.New("-- error unable to get balance buyer")
	}

	t1 := price
	if b.Value < t1 {
		log.Print("-- commence switch roles buyer")
		swap = true

		// increment counter to check if both accounts is invalid
		count += 1
	}

	if swap {
		tmp := bot.accountTwo
		bot.accountTwo = bot.accountOne
		bot.accountOne = tmp

		if count == 1 {
			goto switchAccount
		} else {
			bot.running = false
			insertBotStatus("OFF")

			bot.simulate = true
			insertBotSimulateStatus("ON")

			return errors.New("-- both account can't buy using ETH")
		}
	}
	return nil
}

func sum(c []float64) float64 {
	sum := float64(0)
	for _, x := range c {
		sum += x
	}
	return sum
}

func sendPayload(method, path string, headers map[string]string, body io.Reader,
	result interface{}) error {
	method = strings.ToUpper(method)

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Printf("%s %d", res.Status, res.StatusCode)
	log.Printf("%s", string(contents))

	err = json.Unmarshal(contents, &result)
	return err
}

func insertTransaction(t, pair string, price, quantity float64, simulate, remarks string) {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "trader",
		Precision: "s",
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("-- %s %s", t, pair)

	tags := map[string]string{
		"exchange": "livecoin",
		"pair":     pair,
		"type":     t,
		"simulate": simulate,
	}

	fields := echo.Map{
		"price":    price,
		"quantity": quantity,
		"remarks":  remarks,
	}

	pt, err := influx.NewPoint("transactions", tags, fields, time.Now())
	bp.AddPoint(pt)
	if err != nil {
		log.Println(err)
	}

	err = bot.store.Write(bp)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(1 * time.Second)
}

func updateTicker(pair string) echo.Map {
	p, err := getTicker(pair)
	if err != nil {
		log.Println(err.Error())
	}

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "trader",
		Precision: "s",
	})
	if err != nil {
		log.Println(err.Error())
	}

	tags := map[string]string{
		"type":     "ticker",
		"pair":     "nox_eth",
		"exchange": "livecoin",
	}

	fields := echo.Map{
		"symbol":        p.Currency,
		"high":          p.High,
		"low":           p.Low,
		"volume":        p.Volume,
		"ask":           p.BestAsk,
		"askVolume":     -1,
		"bid":           p.BestBid,
		"bidVolume":     -1,
		"vwap":          p.Vwap,
		"open":          -1,
		"close":         p.Last,
		"previousClose": -1,
		"change":        -1,
		"percentage":    -1,
		"average":       -1,
		"baseVolume":    p.Volume,
		"quoteVolume":   p.Volume * p.Vwap,
	}

	pt, err := influx.NewPoint("stream", tags, fields, time.Now())
	bp.AddPoint(pt)
	err = bot.store.Write(bp)
	if err != nil {
		log.Println(err.Error())
	}

	return fields
}

func insertTickerUpdate(p *Period) echo.Map {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "trader",
		Precision: "s",
	})
	if err != nil {
		log.Println(err)
	}

	tags := map[string]string{
		"type":     "ticker",
		"pair":     "btc_usd",
		"exchange": "poloniex",
	}
	fields := echo.Map{
		"symbol":        "BTC",
		"high":          p.High,
		"low":           p.Low,
		"volume":        p.Volume,
		"ask":           -1.0,
		"askVolume":     -1.0,
		"bid":           -1.0,
		"bidVolume":     -1.0,
		"vwap":          -1.0,
		"open":          p.Open,
		"close":         p.WeightedAverage,
		"previousClose": -1.0,
		"change":        -1.0,
		"percentage":    -1.0,
		"average":       -1.0,
		"baseVolume":    p.Volume,
		"quoteVolume":   p.QuoteVolume,
	}

	pt, err := influx.NewPoint("stream", tags, fields, time.Now())
	bp.AddPoint(pt)
	err = bot.store.Write(bp)
	if err != nil {
		log.Println(err)
	}

	return fields
}

func insertBotStatus(status string) echo.Map {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "trader",
		Precision: "s",
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("-- Bot power %s", status)

	tags := map[string]string{
		"set":  "bot",
		"type": "power",
	}
	fields := echo.Map{
		"Status": status,
	}

	pt, err := influx.NewPoint("bot", tags, fields, time.Now())
	bp.AddPoint(pt)
	err = bot.store.Write(bp)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(1 * time.Second)
	return fields
}

func insertBotSimulateStatus(status string) echo.Map {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "trader",
		Precision: "s",
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("-- Bot simulate %s", status)

	tags := map[string]string{
		"set":  "bot",
		"type": "simulate",
	}
	fields := echo.Map{
		"Status": status,
	}

	pt, err := influx.NewPoint("bot", tags, fields, time.Now())
	bp.AddPoint(pt)
	err = bot.store.Write(bp)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(1 * time.Second)
	return fields
}

func createSignature(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	d := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(d)
}

type IntRange struct {
	min, max int64
}

func (ir *IntRange) NextRandom(r *rand.Rand) int64 {
	return r.Int63n(ir.max-ir.min) + ir.min
}

func getRandom(v float64) float64 {
	value := int64(v)
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	ir := IntRange{-value, value}
	return float64(ir.NextRandom(r)) + r.Float64()
}
