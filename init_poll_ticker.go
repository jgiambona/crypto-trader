package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ffimnsr/trader/exchange"
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

	currentMovingAverage := float64(0)
	lenghtOfMA := 0
	simulation := false
	historicalData := []Period{}
	tradePlaced := false
	typeOfTrade := ""
	for {
		waitExchanges.Add(len(bot.exchanges))
		for _, x := range bot.exchanges {
			//log.Printf("check updated prices")
			//log.Printf("calculate sell price")
			//log.Printf("calculate profit")
			//log.Printf("check if selling price is gt sell worth then begin selling")

			go func(c exchange.BotExchange) {
				defer waitExchanges.Done()

				historicalData, err = getTestData()
				if err != nil {
					log.Fatal(err.Error())
				}

				var lastPairPrice float64
				if simulation && (len(historicalData) > 0) {
					var nextDataPoint Period
					// Pop and shift data.
					nextDataPoint, historicalData = historicalData[0], historicalData[1:]
					lastPairPrice = nextDataPoint.WeightedAverage
				} else if simulation && (len(historicalData) < 1) {
					log.Fatal("unable to check live data.")
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
							tradePlaced = true
							typeOfTrade = "short"
						} else if (lastPairPrice < currentMovingAverage) && (lastPairPrice > previousPrice) {
							log.Print("BUY ORDER")
							tradePlaced = true
							typeOfTrade = "long"
						}
					} else if typeOfTrade == "short" {
						if lastPairPrice < currentMovingAverage {
							log.Print("EXIT TRADE")
							tradePlaced = false
							typeOfTrade = ""
						}
					} else if typeOfTrade == "long" {
						if lastPairPrice > currentMovingAverage {
							log.Print("EXIT TRADE")
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

func getTestData() ([]Period, error) {
	result := []Period{}
	path := "https://poloniex.com/public?command=returnChartData&currencyPair=BTC_XMR&start=1405699200&end=9999999999&period=14400"
	return result, sendPayload("GET", path, nil, nil, &result)
}

func sendPayload(method, path string, headers map[string]string, body io.Reader, result interface{}) error {
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

func insertTransactions() {
	//bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
	//	Database:  "trader",
	//	Precision: "s",
	//})
	//if err != nil {
	//	log.Fatalf("%s", err.Error())
	//}

	//tags := map[string]string{
	//	"exchange": "livecoin",
	//	"pair":     "btc_usd",
	//	"type":     "buy",
	//}
	//fields := echo.Map{
	//	"price":  1.444,
	//	"amount": 1.444,
	//	"fee":    1.4444,
	//	"volume": 1.4444,
	//}

	//pt, err := influx.NewPoint("transactions", tags, fields, time.Now())
	//bp.AddPoint(pt)
	//err = bot.store.Write(bp)
	//if err != nil {
	//	log.Fatalf("%s", err.Error())
	//}
}
