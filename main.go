package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ffimnsr/trader/exchange"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

// Contains forward declaration of API secrets.
const (
	apiKey    = "gZaXWxSudtwt1AR3cW6Fdh5UY3BgVG4r"
	apiSecret = "FYw6X3gzcJ4F5JvqmYBqAMwdMexzAay7"
)

var strategies = map[string]string{
	"uptrend":      "Uptrend",
	"bb":           "Bollinger Bands",
	"gain":         "Gain",
	"pp":           "Pingpong",
	"stepgain":     "Stepgain",
	"tssl":         "Trailing Stop / Stop Limit",
	"emotionless":  "Emotionless",
	"ichimoku":     "Ichimoku",
	"tsslbb":       "Trailing Stop / Stop Limit - Bollinger Bands",
	"tsslpp":       "Trailing Stop / Stop Limit - Pingpong",
	"tsslstepgain": "Trailing Stop / Stop Limit - Stepgain",
	"tsslgain":     "Trailing Stop / Stop Limit - Gain",
	"bbrsitssl":    "Bollinger Bands + RSI - Trailing Stop / Stop Limit",
	"pptssl":       "Pingpong - Trailing Stop / Stop Limit",
	"stepgaintssl": "Stepgain - Trailing Stop / Stop Limit",
	"gaintssl":     "Gain - Trailing Stop / Stop Limit",
	"bbtssl":       "Bollinger Bands - Trailing Stop / Stop Limit",
	"bbgain":       "Bollinger Bands - Gain",
	"gainbb":       "Gain - Bollinger Bands",
	"bbstepgain":   "Bollinger Bands - Stepgain",
	"stepgainbb":   "Stepgain - Bollinger Bands",
	"bbpp":         "Bollinger Bands - Pingpong",
	"ppbb":         "Pingpong - Bollinger Bands",
	"gainstepgain": "Gain - Stepgain",
	"stepgaingain": "Stepgain - Gain",
	"gainpp":       "Gain - Pingpong",
	"stepgainpp":   "Stepgain - Pingpong",
	"ppstepgain":   "Pingpong - Stepgain",
}

// Bot is the singleton that holds all the data.
type Bot struct {
	config    *BotConfig
	exchanges []exchange.BotExchange
}

var bot Bot

func main() {
	e := echo.New()
	e.Renderer = loadTemplates(e)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	loadConfig()
	loadExchanges()
	if len(bot.exchanges) == 0 {
		e.Logger.Fatal("no exchanges were loaded.")
	}

	//c, err := influx.NewHTTPClient(influx.HTTPConfig{
	//	Addr: "http://localhost:8086",
	//})
	//if err != nil {
	//	e.Logger.Fatal(err)
	//}

	//influx.

	go templateWatch(e)
	go pollTicker()

	e.Static("/public", "public")
	e.Add("GET", "/", index)
	e.Add("GET", "/bot/restart", index)
	e.Add("GET", "/bot/suspend", index)
	e.Add("GET", "/ws", socket)

	if port, ok := os.LookupEnv("PORT"); ok {
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		e.Logger.Fatal(e.Start(":8000"))
	}
}

func index(c echo.Context) error {
	context := map[string]interface{}{
		"strategies": strategies,
	}
	return c.Render(http.StatusOK, "index.tmpl", context)
}

func socket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			err := websocket.Message.Send(ws, "Hello, Client!")
			if err != nil {
				c.Logger().Error(err)
			}

			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}

			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
