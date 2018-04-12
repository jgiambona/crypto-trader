package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

func socketCheckBalance() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: "ws.blockchain.info", Path: "/inv"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer c.Close()

	data, err := json.Marshal(echo.Map{
		"op": "ping",
	})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write: ", err.Error())
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read: ", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write: ", err)
				return
			}
			<-done
			return
		}
	}
}
