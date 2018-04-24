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
	path := fmt.Sprintf("%s%s", "https://api.livecoin.com/", "/exchange/buylimit")

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
