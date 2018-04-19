package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ffimnsr/trader/exchange"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// Bot is the singleton that holds all the data.
type Bot struct {
	config    *BotConfig
	exchanges []exchange.BotExchange
	store     influx.Client
	db        *sql.DB
	nextID    int64
}

var bot Bot

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Create and open SQLite3 database file.
	var err error
	bot.db, err = sql.Open("sqlite3", "./data/trader.db")
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer bot.db.Close()

	if err := createLocalDB(); err != nil {
		e.Logger.Fatal(err)
	}

	if len(os.Getenv("T_PROD")) > 0 {
		bot.store, _ = influx.NewHTTPClient(influx.HTTPConfig{
			Addr: "http://ec2-54-169-102-171.ap-southeast-1.compute.amazonaws.com:8086",
		})
	} else {
		bot.store, _ = influx.NewHTTPClient(influx.HTTPConfig{
			Addr: "http://localhost:8086",
		})
	}

	id, err := getLastAccountID()
	if err != nil {
		e.Logger.Fatal(err)
	}
	log.Printf("--- Number of Accounts: %d", id)
	bot.nextID = id + 1

	loadConfig()
	loadExchanges()
	if len(bot.exchanges) == 0 {
		e.Logger.Fatal("no exchanges were loaded.")
	}

	//go socketCheckBalance()
	//go pollTicker()

	e.Add("GET", "/", index)
	e.Add("GET", "/bot/:name/info", getExchangeConfigInfo)
	e.Add("GET", "/bot/restart", restart)
	e.Add("GET", "/bot/suspend", suspend)

	e.Add("POST", "/bot/setup", setup)
	e.Add("POST", "/bot/setup/pp", setupPingpong)
	e.Add("POST", "/bot/setup/portfolio", addNewPortfolio)
	e.Add("POST", "/bot/setup/portfolio/create", addNewPortfolio)

	if port, ok := os.LookupEnv("PORT"); ok {
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		e.Logger.Fatal(e.Start(":8000"))
	}
}

func setup(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}

func setupPingpong(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}

func restart(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func suspend(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
