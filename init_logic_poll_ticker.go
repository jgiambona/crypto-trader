package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
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

	bot.simulate = true
	insertBotSimulateStatus("ON")

	currencyPair := "NOX/ETH"

	var placedOrder int64
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

					switchAccountRolesSeller()
					switchAccountRolesBuyer()

					tradePlace := false
					fromAccountOne := -1.0

					p := updateTicker(currencyPair)
					lowest := p["ask"].(float64)
					volume := p["volume"].(float64)

					log.Printf("--- %.8f %0.8f", lowest, volume)

					if bot.ruleOne.Enabled {
						v := bot.ruleOne.TransactionVolume * 0.10
						quantity := bot.ruleOne.TransactionVolume + getRandom(v)
						targetPrice := lowest - bot.ruleOne.BidPriceStepDown
						if targetPrice >= bot.ruleOne.MinimumBid {
							if botVolume < bot.ruleOne.MaximumVolume {
								botVolume += quantity
								insertTransaction("SELL", "nox_eth", targetPrice, quantity)
								if !bot.simulate {
									o, err := sellLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, targetPrice, quantity)
									if err != nil {
										log.Print("error occurred in creating sell order ", err)
									}
									placedOrder = o.OrderID
								}
								fromAccountOne = targetPrice
							}

							if lowest >= fromAccountOne && fromAccountOne > -1 &&
								placedOrder > 0 {
								insertTransaction("BUY", "nox_eth", targetPrice, quantity)
								if !bot.simulate {
									buyLimit(bot.accountTwo.APIKey, bot.accountTwo.APISecret,
										currencyPair, targetPrice, quantity)
								}
								tradePlace = true
							}
						}
					}

					if bot.ruleTwo.Enabled && !tradePlace {
						v := bot.ruleOne.TransactionVolume * 0.10
						quantity := bot.ruleOne.TransactionVolume + getRandom(v)
						targetPrice := lowest - bot.ruleTwo.BidPriceStepDown
						if targetPrice >= bot.ruleTwo.MinimumBid {
							if volume < bot.ruleTwo.MaximumVolume {
								botVolume += quantity
								insertTransaction("SELL", "nox_eth", targetPrice, quantity)
								if bot.simulate {
									o, err := sellLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, targetPrice, quantity)
									if err != nil {
										log.Print("error occurred in creating sell order")
									}
									placedOrder = o.OrderID
								}
								fromAccountOne = targetPrice
								tradePlace = true
							}

							if lowest != fromAccountOne && fromAccountOne > -1 &&
								placedOrder > 0 {
								insertTransaction("CANCEL", "nox_eth", targetPrice, quantity)
								if bot.simulate {
									c, err := cancelLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
										currencyPair, placedOrder)

									if err != nil {
										log.Print("error occurred in cancelling order")
									}

									if !c.Success {
										log.Print("unable to cancel order")
									}
								}
								goto repeatCheckLowestBid
							}

							if lowest >= fromAccountOne && fromAccountOne > -1 &&
								placedOrder > 0 {
								insertTransaction("BUY", "nox_eth", targetPrice, quantity)
								if bot.simulate {
									buyLimit(bot.accountTwo.APIKey, bot.accountTwo.APISecret,
										currencyPair, targetPrice, quantity)
								}
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

func switchAccountRolesSeller() (err error) {
	swap := false

	b, err := getBalance(bot.accountOne.APIKey, bot.accountOne.APISecret, "NOX")
	if err != nil {
		log.Print("unable to get balance", err)
	}

	t1 := bot.ruleOne.TransactionVolume + (bot.ruleOne.TransactionVolume * bot.ruleOne.VarianceOfTransaction)
	t2 := bot.ruleTwo.TransactionVolume + (bot.ruleTwo.TransactionVolume * bot.ruleTwo.VarianceOfTransaction)
	if b.Value < t1 || b.Value < t2 {
		swap = true
	}

	if swap {
		tmp := bot.accountOne
		bot.accountOne = bot.accountTwo
		bot.accountTwo = tmp
	}

	time.Sleep(1 * time.Second)
	return nil
}

func switchAccountRolesBuyer() (err error) {
	swap := false

	b, err := getBalance(bot.accountTwo.APIKey, bot.accountTwo.APISecret, "ETH")
	if err != nil {
		log.Print("unable to get balance", err)
	}

	t1 := bot.ruleOne.TransactionVolume + (bot.ruleOne.TransactionVolume * bot.ruleOne.VarianceOfTransaction)
	t2 := bot.ruleTwo.TransactionVolume + (bot.ruleTwo.TransactionVolume * bot.ruleTwo.VarianceOfTransaction)
	if b.Value < t1 || b.Value < t2 {
		swap = true
	}

	if swap {
		tmp := bot.accountTwo
		bot.accountTwo = bot.accountOne
		bot.accountOne = tmp
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
	log.Printf("%+v", contents)

	err = json.Unmarshal(contents, &result)
	return err
}

func insertTransaction(t, pair string, price, quantity float64) {
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
	}

	fields := echo.Map{
		"price":    price,
		"quantity": quantity,
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
