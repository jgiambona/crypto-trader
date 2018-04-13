package main

import (
	"sync"
	"time"

	"github.com/ffimnsr/trader/exchange"
)

// PollTicker fetches and updates the ticker for all exchanges.
func pollTicker() {

	var waitExchanges sync.WaitGroup
	for {
		waitExchanges.Add(len(bot.exchanges))
		for _, x := range bot.exchanges {
			//log.Printf("check updated prices")
			//log.Printf("calculate sell price")
			//log.Printf("calculate profit")
			//log.Printf("check if selling price is gt sell worth then begin selling")

			go func(c exchange.BotExchange) {
				defer waitExchanges.Done()
				// c.UpdateTicker()
			}(x)
		}

		waitExchanges.Wait()
		time.Sleep(10 * time.Second)
	}
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
