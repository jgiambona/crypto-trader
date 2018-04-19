package main

import (
	"errors"
	"log"
	"strings"

	"github.com/ffimnsr/trader/exchange"
	"github.com/ffimnsr/trader/exchange/livecoin"
)

// LoadExchanges initialize all exchange that are available
// on the platform.
func loadExchanges() error {
	for _, x := range bot.config.Exchanges {
		initializeExchange(x.Name)
	}

	return nil
}

func initializeExchange(name string) error {
	var exch exchange.BotExchange

	log.Print("--- Exchanges:")
	switch strings.ToLower(name) {
	case "livecoin":
		log.Printf("--- + -- [active] %s", name)
		exch = livecoin.NewInstance()
	default:
		return errors.New("exchange not found")
	}

	if exch == nil {
		return errors.New("exchange failed to load")
	}

	bot.exchanges = append(bot.exchanges, exch)
	return nil
}
