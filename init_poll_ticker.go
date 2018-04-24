package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
)

// PollTicker fetches and updates the ticker for all exchanges.
func pollTicker() {

	var waitExchanges sync.WaitGroup
	var prices []float64
	var err error

	pair := "btc_usd"
	quantity := 0.001
	currentMovingAverage := float64(0)
	lenghtOfMA := 0
	simulation := bot.simulation
	historicalData := []Period{}
	tradePlaced := false
	typeOfTrade := ""
	if simulation {
		log.Print("getting historical data")

		// Try to get fresh historical data from poloniex server if none
		// check the data folder for old save data.
		historicalData, err = getHistoricalData()
		if err != nil {
			historicalData, err = getLocalHistoricalData()
		}
	}

	for {
		waitExchanges.Add(len(bot.exchanges))
		for _, x := range bot.exchanges {
			//log.Printf("check updated prices")
			//log.Printf("calculate sell price")
			//log.Printf("calculate profit")
			//log.Printf("check if selling price is gt sell worth then begin selling")

			go func(c exchange.BotExchange) {
				defer waitExchanges.Done()

				var lastPairPrice float64
				if simulation && (len(historicalData) > 0) {
					var nextDataPoint Period

					// Pop and shift data.
					nextDataPoint, historicalData = historicalData[0], historicalData[1:]
					lastPairPrice = nextDataPoint.WeightedAverage
				} else if simulation && (len(historicalData) < 1) {
					log.Print("finish running historical data")
					os.Exit(0)
				} else {
					values := c.UpdateTicker()
					lastPairPrice = values["close"].(float64)
				}

				if len(prices) > 0 {
					currentMovingAverage = sum(prices) / float64(len(prices))
					previousPrice := prices[len(prices)-1]
					if !tradePlaced {
						if (lastPairPrice > currentMovingAverage) && (lastPairPrice < previousPrice) {
							log.Print("SELL ORDER")
							insertTransaction("sell", pair, lastPairPrice, quantity)
							if !simulation {

							}
							tradePlaced = true
							typeOfTrade = "short"
						} else if (lastPairPrice < currentMovingAverage) && (lastPairPrice > previousPrice) {
							log.Print("BUY ORDER")
							insertTransaction("buy", pair, lastPairPrice, quantity)
							if !simulation {

							}
							tradePlaced = true
							typeOfTrade = "long"
						}
					} else if typeOfTrade == "short" {
						if lastPairPrice < currentMovingAverage {
							log.Print("EXIT TRADE")
							insertTransaction("cancel", pair, lastPairPrice, quantity)
							if !simulation {

							}
							tradePlaced = false
							typeOfTrade = ""
						}
					} else if typeOfTrade == "long" {
						if lastPairPrice > currentMovingAverage {
							log.Print("EXIT TRADE")
							insertTransaction("cancel", pair, lastPairPrice, quantity)
							if !simulation {

							}
							tradePlaced = false
							typeOfTrade = ""
						}
					}
				}
				prices = append(prices, lastPairPrice)
				prices = prices[-lenghtOfMA:]

				log.Printf("last_price: %0.6f moving avg.: %0.6f", lastPairPrice, currentMovingAverage)
			}(x)
		}

		waitExchanges.Wait()
		if !simulation {
			time.Sleep(5 * time.Second)
		}
	}
}

func sum(c []float64) float64 {
	sum := float64(0)
	for _, x := range c {
		sum += x
	}
	return sum
}

func getHistoricalData() ([]Period, error) {
	result := []Period{}
	path := "https://poloniex.com/public?command=returnChartData&currencyPair=BTC_XMR&start=1405699200&end=9999999999&period=14400"
	return result, sendPayload("GET", path, nil, nil, &result)
}

func getLocalHistoricalData() ([]Period, error) {
	result := []Period{}
	path := "./data/historical-data-poloniex.json"
	
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	return result, json.Unmarshal(raw, &result)
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
		log.Fatal(err)
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
		log.Fatal(err)
	}

	err = bot.store.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
}
