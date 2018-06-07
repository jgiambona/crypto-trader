package main

import (
	"os"
	"strings"
	"time"

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
		running                bool
		simulate               bool
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

	bot.availableCurrencyPairs = strings.Fields(CurrencyPairAllowed)
	bot.baseCurrencies = strings.Fields(CurrencyAllowed)

	bot.accountOne.APIKey = "bot"

	bot.accountTwo.APIKey = "bot"

	bot.ruleOne.Enabled = true
	bot.ruleOne.MinInterval = time.Duration(7 * time.Second)
	bot.ruleOne.MaxInterval = time.Duration(10 * time.Second)
	bot.ruleOne.CheckOrderDelay = time.Duration(1 * time.Second)
	bot.ruleOne.MaximumVolume = 500000.0
	bot.ruleOne.TransactionVolume = 300.0
	bot.ruleOne.VarianceOfTransaction = 10.0
	bot.ruleOne.MinBidPriceStepDown = 0.00000001
	bot.ruleOne.MaxBidPriceStepDown = 0.00000010
	bot.ruleOne.FloorPriceGap = 0.00020000
	bot.ruleOne.MinimumBid = 0.00011500

	loadRoutes(e)

	go pollTicker()

	if port, ok := os.LookupEnv("PORT"); ok {
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		if len(os.Getenv("T_PROD")) > 0 {
			e.Logger.Fatal(e.StartTLS("0.0.0.0:8000", "ssl/fullchain.pem", "ssl/privkey.pem"))
		} else {
			e.Logger.Fatal(e.Start("localhost:8000"))
		}
	}
}

func loadRoutes(e *echo.Echo) {
	e.Add("GET", "/", index)
	e.Add("GET", "/bot/exported", botExported)
	e.Add("POST", "/bot/controls", botControls)
	e.Add("POST", "/bot/settings", updateSettings)
	e.Add("POST", "/bot/accounts", updateAccounts)
	e.Add("POST", "/bot/simulate", updateSimulate)
}
