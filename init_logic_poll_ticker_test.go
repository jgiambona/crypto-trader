package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPollTicker(t *testing.T) {
	assert.Equal(t, true, true)
}

func TestBuy(t *testing.T) {
	path := fmt.Sprintf("%s%s", "https://api.livecoin.net", "/exchange/buylimit")

	construct := url.Values{}
	construct.Add("currencyPair", "BTC/USD")
	construct.Add("price", "60")
	construct.Add("quantity", "0.001")
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        "vGMWuD7nw6WDQsRSA3SxW4mnqe4CnG4x",
		"Sign":           createSignature(message, "dVvBd2HCCmEsJCUHSYCxUgutskkUJEbe"),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	err := sendPayload("POST", path, headers, strings.NewReader(message), &data)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Printf("Hello World")
}

func TestGetOrder(t *testing.T) {
	construct := url.Values{}
	construct.Add("orderId", string(7392549951))
	message := construct.Encode()

	headers := map[string]string{
		"API-Key": "bVFxGcKEawAWeRR6je5jpVTF8jP3Qm47",
		"Sign":    createSignature(message, "iwqhegnHSq3Y14VVPVY1AFMDfGgGYKEhF"),
	}

	path := fmt.Sprintf("%s/exchange/order?%s", LiveCoinAPIURL, message)
	log.Print("-- ", path)
	data := OrderDetailResponse{}
	err := sendPayload("GET", path, headers, nil, &data)
	t.Logf("+%v", data)
	assert.NotNil(t, err)
}

func TestGetOrderBook(t *testing.T) {
	construct := url.Values{}
	construct.Add("currencyPair", "NOX/ETH")
	construct.Add("depth", "4")
	message := construct.Encode()

	path := fmt.Sprintf("%s/exchange/order_book?%s", LiveCoinAPIURL, message)
	log.Print("-- ", path)
	data := OrderBookResponse{}
	err := sendPayload("GET", path, nil, nil, &data)
	c := len(data.Asks[0][0])
	t.Logf("%#+v", c)
	t.Logf("%#+v", data.Asks[0][0])
	t.Logf("%#+v", data.Asks[0][0][0:c])
	assert.NotNil(t, err)
}
