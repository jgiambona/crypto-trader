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

	"github.com/ffimnsr/trader/exchange"
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

	// OrderResponse receives the order status.
	OrderResponse struct {
		Success bool  `json:"success"`
		Added   bool  `json:"added"`
		OrderID int64 `json:"orderId"`
	}
)

// PollTicker fetches and updates the ticker for all exchanges.
func pollTicker() {

	var waitExchanges sync.WaitGroup

	//pair := "btc_usd"
	//quantity := 0.001
	//simulationStart := false
	//historicalData := []Period{}
	//tradePlaced := false
	//typeOfTrade := ""

	// accountOne := true
	// accountTwo := true

	ruleOne := true
	ruleTwo := true

	for {
		waitExchanges.Add(len(bot.exchanges))
		for _, x := range bot.exchanges {
			go func(c exchange.BotExchange) {
				defer waitExchanges.Done()

				// switchAccountRoles

				tradeStop := false

				// repeatCheckLowestBid:
				// lowest := getLowestBidInQueue()
				// targetPrice := lowest - stepDownPrice
				// if targetPrice >= rangePriceAndAmount
				//   accountOnePlaceBidLowerThanInQueue

				// lowest := getLowesetBidInQueue()
				if ruleOne {
					// if lowest == fromAccountOne
					//   accountTwoBuyBid
					//   tradeStop = true
				}

				if ruleTwo && !tradeStop {
					// if lowest != fromAccountOne
					//   accountOneCancelBid
					//   goto repeatCheckLowestBid
					//
					//
					// if lowest == fromAccountOne
					//   accountTwoBuyBid
				}

				//var lastPairPrice float64
				//values := c.UpdateTicker()
				//lastPairPrice = values["close"].(float64)
			}(x)
		}

		waitExchanges.Wait()
		time.Sleep(7 * time.Second)
	}
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

