package main

import (
	"os"
	"strings"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type (
	// Bot is the singleton that holds all the data.
	Bot struct {
		store                  influx.Client
		accountOne             Account
		accountTwo             Account
		ruleOne                RuleConfiguration
		ruleTwo                RuleConfiguration
		running                bool
		availableCurrencyPairs []string
		baseCurrencies         []string
	}
)

var bot Bot

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if len(os.Getenv("T_PROD")) > 0 {
		bot.store, _ = influx.NewHTTPClient(influx.HTTPConfig{
			Addr: "http://ec2-54-169-102-171.ap-southeast-1.compute.amazonaws.com:8086",
		})
	} else {
		bot.store, _ = influx.NewHTTPClient(influx.HTTPConfig{
			Addr: "http://localhost:8086",
		})
	}

	bot.running = false
	bot.availableCurrencyPairs = strings.Fields(CurrencyPairAllowed)
	bot.baseCurrencies = strings.Fields(CurrencyAllowed)

	bot.accountOne.APIKey = "bot"
	bot.accountTwo.APIKey = "bot"
	bot.ruleOne.Enabled = true
	bot.ruleOne.Interval = 7
	bot.ruleOne.MaximumVolume = 50000.0
	bot.ruleOne.TransactionVolume = 3000.0
	bot.ruleOne.VarianceOfTransaction = 10.0
	bot.ruleOne.BidPriceStepDown = 0.0001
	bot.ruleOne.MinimumBid = 0.00001

	bot.ruleTwo.Enabled = true
	bot.ruleTwo.Interval = 7
	bot.ruleTwo.MaximumVolume = 50000.0
	bot.ruleTwo.TransactionVolume = 3000.0
	bot.ruleTwo.VarianceOfTransaction = 10.0
	bot.ruleTwo.BidPriceStepDown = 0.0001
	bot.ruleTwo.MinimumBid = 0.00001

	loadRoutes(e)

	go pollTicker()

	if port, ok := os.LookupEnv("PORT"); ok {
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		if len(os.Getenv("T_PROD")) > 0 {
			e.Logger.Fatal(e.Start("0.0.0.0:8000"))
		} else {
			e.Logger.Fatal(e.Start("localhost:8000"))
		}
	}
}

func loadRoutes(e *echo.Echo) {
	e.Add("GET", "/", index)

	e.Add("POST", "/bot/controls", botControls)
	e.Add("POST", "/bot/settings", updateSettings)
	e.Add("POST", "/bot/accounts", updateAccounts)
}
