package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
)

// Ticker for websocket GDAX.
type Ticker struct {
	Type        string `json:"type"`
	Sequence    int64  `json:"sequence"`
	ProductID   string `json:"product_id"`
	Price       string `json:"price"`
	Open        string `json:"open_24h"`
	Volume      string `json:"volume_24h"`
	Low         string `json:"low_24h"`
	High        string `json:"high_24h"`
	VolumeMonth string `json:"volume_30d"`
	BestBid     string `json:"best_bid"`
	BestAsk     string `json:"best_ask"`
	Side        string `json:"side"`
	Timestamp   string `json:"time"`
	TradeID     int64  `json:"trade_id"`
	LastSize    string `json:"last_size"`
}

func socketCheckBalance() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: "ws-feed.gdax.com", Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer c.Close()

	data, err := json.Marshal(echo.Map{
		"type": "subscribe",
		"channels": []echo.Map{
			echo.Map{
				"name":        "ticker",
				"product_ids": []string{"BTC-USD"},
			},
		},
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

			go func(x []byte) {
				writeTicker(x)
			}(message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
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

func writeTicker(message []byte) {
	ticker := Ticker{}
	if err := json.Unmarshal(message, &ticker); err != nil {
		log.Fatal(err.Error())
	}

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "trader",
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	tags := map[string]string{
		"type": "ticker",
	}
	fields := map[string]interface{}{
		"open":     ticker.Open,
		"price":    ticker.Price,
		"high":     ticker.High,
		"low":      ticker.Low,
		"volume":   ticker.Volume,
		"best_ask": ticker.BestAsk,
		"best_bid": ticker.BestBid,
	}

	pt, err := influx.NewPoint("ticker", tags, fields, time.Now())
	bp.AddPoint(pt)
	err = bot.store.Write(bp)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
