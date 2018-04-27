package main

import (
	"os"

	"github.com/ffimnsr/trader/exchange"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

// Bot is the singleton that holds all the data.
type Bot struct {
	exchanges []exchange.BotExchange
	store     influx.Client
}

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

	loadRoutes(e)

	go pollTicker()

	if port, ok := os.LookupEnv("PORT"); ok {
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		e.Logger.Fatal(e.Start("localhost:8000"))
	}
}

func loadRoutes(e *echo.Echo) {
	e.Add("GET", "/", index)

	e.Add("POST", "/bot/settings", updateSettings)
	e.Add("POST", "/bot/accounts", updateSettings)
}
