package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"

	"github.com/ffimnsr/trader/exchange/livecoin"
)

// Contains forward declaration of API secrets.
const (
	APIKey       = "gZaXWxSudtwt1AR3cW6Fdh5UY3BgVG4r"
	APISecretKey = "FYw6X3gzcJ4F5JvqmYBqAMwdMexzAay7"
)

func main() {
	e := echo.New()
	e.Renderer = LoadTemplates(e)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Add("GET", "/", index)
	e.Add("GET", "/ws", socket)
	e.Logger.Fatal(e.Start(":4000"))
}

func index(c echo.Context) error {
	o := livecoin.NewInstance()
	o.GetTicker("BTC")
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
