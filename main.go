package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ffimnsr/trader/exchange"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

// Bot is the singleton that holds all the data.
type Bot struct {
	config     *BotConfig
	exchanges  []exchange.BotExchange
	store      influx.Client
	simulation bool
	db         *sql.DB
	nextID     int64
}

var bot Bot

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Enable paper trading.
	bot.simulation = true

	// Create and open SQLite3 database file.
	var err error
	connStr := "postgres://trader:trader@localhost/trader?sslmode=disable"
	bot.db, err = sql.Open("postgres", connStr)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer bot.db.Close()

	if err := bot.db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}

	if err := repoCreateDB(); err != nil {
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

	id, err := repoGetLastAccountID()
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

	loadRoutes(e)

	//go socketCheckBalance()
	//go pollTicker()

	if port, ok := os.LookupEnv("PORT"); ok {
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		e.Logger.Fatal(e.Start("localhost:8000"))
	}
}

