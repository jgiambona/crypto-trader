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
	APIKey    = "gZaXWxSudtwt1AR3cW6Fdh5UY3BgVG4r"
	APISecret = "FYw6X3gzcJ4F5JvqmYBqAMwdMexzAay7"
)

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
		e.Logger.Fatal("No exchanges were loaded.")
	}

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
	return c.Render(http.StatusOK, "index.tmpl", nil)
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
