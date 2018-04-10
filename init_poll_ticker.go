package main

import (
	"log"
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
			log.Printf("check updated prices")
			log.Printf("calculate sell price")
			log.Printf("calculate profit")
			log.Printf("check if selling price is gt sell worth then begin selling")

			go func(c exchange.BotExchange) {
				defer waitExchanges.Done()
				c.UpdateTicker()
			}(x)
		}

		waitExchanges.Wait()
		time.Sleep(10 * time.Second)
	}
}
