package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
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

	currencyPair := "NOX/ETH"
	//quantity := 0.001
	//placedOrder := 0.001

	for {
		waitExchanges.Add(1)
		go func() {
			defer waitExchanges.Done()
			updateTicker(currencyPair)

			//tradePlace := false
			switchAccountRoles()

//		repeatCheckLowestBid:
//			lowest := getLowestBidInQueue()
//
//			if bot.ruleOne.Enabled {
//				targetPrice := lowest - bot.ruleOne.BidPriceStepDown
//				if targetPrice >= rangePriceAndAmount {
//					insertTransaction("SELL", "nox_eth", targetPrice, quantity)
//					placedOrder = 0.001
//				}
//
//				if lowest == fromAccountOne {
//					insertTransaction("BUY", "nox_eth", targetPrice, quantity)
//					buyLimit(bot.accountTwo.APIKey, bot.accountTwo.APISecret,
//						currencyPair, targetPrice, quantity)
//					tradePlace = true
//				}
//			}
//
//			if bot.ruleTwo.Enabled && !tradePlace {
//				targetPrice := lowest - bot.ruleTwo.BidPriceStepDown
//				if targetPrice >= rangePriceAndAmount {
//					insertTransaction("SELL", "nox_eth", targetPrice, quantity)
//					sellLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
//						currencyPair, targetPrice, quantity)
//					placedOrder = 0.001
//				}
//
//				if lowest != fromAccountOne {
//					insertTransaction("CANCEL", "nox_eth", targetPrice, quantity)
//					cancelLimit(bot.accountOne.APIKey, bot.accountOne.APISecret,
//						currencyPair, placedOrder)
//					goto repeatCheckLowestBid
//				}
//
//				if lowest == fromAccountOne {
//					insertTransaction("BUY", "nox_eth", targetPrice, quantity)
//					buyLimit(bot.accountTwo.APIKey, bot.accountTwo.APISecret,
//						currencyPair, targetPrice, quantity)
//				}
//			}
		}()

		waitExchanges.Wait()
		time.Sleep(7 * time.Second)
	}
}

func switchAccountRoles() (err error) {

	swap := false
	if swap {
		tmp := bot.accountOne
		bot.accountOne = bot.accountTwo
		bot.accountTwo = tmp
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

func createSignature(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	d := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(d)
}
