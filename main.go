package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

// Contains forward declaration of API secrets.
const (
	APIKey       = "gZaXWxSudtwt1AR3cW6Fdh5UY3BgVG4r"
	APISecretKey = "FYw6X3gzcJ4F5JvqmYBqAMwdMexzAay7"
)

// Bot is the singleton that holds all the data.
type Bot struct {
	exchanges []int64
}

func main() {
	e := echo.New()
	e.Renderer = loadTemplates(e)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	loadExchanges()

	go templateWatch(e)
	go pollTicker()

	e.Static("/public", "public")
	e.Add("GET", "/", index)
	e.Add("GET", "/ws", socket)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}

func templateWatch(e *echo.Echo) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer watcher.Close()

	mu := &sync.Mutex{}
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					e.Logger.Print("modified file: ", event.Name)
					mu.Lock()
					e.Renderer = loadTemplates(e)
					mu.Unlock()
				}
			case err := <-watcher.Errors:
				e.Logger.Fatal(err)
			}
		}
	}()

	err = watcher.Add("view/templates/index.tmpl")
	if err != nil {
		e.Logger.Fatal(err)
	}
	<-done
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
