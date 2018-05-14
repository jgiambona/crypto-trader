package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type (
	// OrderResponse parses order creation status.
	OrderResponse struct {
		Success bool  `json:"success"`
		Added   bool  `json:"added"`
		OrderID int64 `json:"orderId"`
	}

	// CancelOrderResponse stores status of cancellation
	// order.
	CancelOrderResponse struct {
		Success       bool    `json:"success"`
		Cancelled     bool    `json:"cancelled"`
		Message       string  `json:"message"`
		Quantity      float64 `json:"quantity"`
		TradeQuantity float64 `json:"tradeQuantity"`
	}

	// BalanceResponse stores available balance status.
	BalanceResponse struct {
		Type     string  `json:"type"`
		Currency string  `json:"currency"`
		Value    float64 `json:"value"`
	}

	// TickerResponse stores the pricing information.
	TickerResponse struct {
		Currency     string  `json:"cur"`
		CurrencyPair string  `json:"symbol"`
		Last         float64 `json:"last"`
		High         float64 `json:"high"`
		Low          float64 `json:"low"`
		Volume       float64 `json:"volume"`
		Vwap         float64 `json:"vwap"`
		MaxBid       float64 `json:"max_bid"`
		MinAsk       float64 `json:"min_ask"`
		BestBid      float64 `json:"best_bid"`
		BestAsk      float64 `json:"best_ask"`
	}
)

// Returns the current fee for the exchange.
func getFee(maker bool) float64 {
	if maker {
		return 0.18 / 100
	}
	return 0.18 / 100
}

// Returns actual trading fee.
func getCommission(apiKey, apiSecret string) (struct {
	Success bool
	Fee float64
	}, error){
	path := fmt.Sprintf("%s/exchange/commission", LiveCoinAPIURL)

	construct := url.Values{}
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        apiKey,
		"Sign":           createSignature(message, apiSecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := struct {
		Success bool
		Fee     float64
	}{}
	return data, sendPayload("GET", path, headers, strings.NewReader(message), &data)
}

// Get information on all currency pair for the last 24 hours.
func getTickerAll() (TickerResponse, error) {
	result := TickerResponse{}
	path := fmt.Sprintf("%s/exchange/ticker", LiveCoinAPIURL)
	return result, sendPayload("GET", path, nil, nil, &result)
}

// Get information on specified currency pair for the last 24 hours.
func getTicker(currencyPair string) (TickerResponse, error) {
	result := TickerResponse{}
	path := fmt.Sprintf("%s/exchange/ticker?currencyPair=%s", LiveCoinAPIURL, currencyPair)
	return result, sendPayload("GET", path, nil, nil, &result)
}

// Get a detailed review on the latest transactions for
// requested currency pair. You may receive the update for the last hour or
// for the last minute.
func getLastTrades(currencyPair, minutesOrHour, tradeType string) {
	path := fmt.Sprintf("%s/exchange/last_trades", LiveCoinAPIURL)
	if err := sendPayload("GET", path, nil, nil, nil); err != nil {
		log.Printf("%s", err.Error())
	}
}

// Get the orderbook for specified currency pair (you may
// enable the feature of grouping orders by price).
func getOrderBook(currencyPair, groupByPrice string, depth int64) {
	path := fmt.Sprintf("%s/exchange/order_book", LiveCoinAPIURL)
	sendPayload("GET", path, nil, nil, nil)
}

// Returns orderbook for every currency pair.
func getAllOrderBook(groupByPrice string, depth int64) {
	path := fmt.Sprintf("%s/exchange/all/order_book", LiveCoinAPIURL)
	sendPayload("GET", path, nil, nil, nil)
}

// Returns maximum bid and minimum ask in the current orderbook.
func getMaxBidMinAsk(currencyPair string) {
	path := fmt.Sprintf("%s/exchange/maxbid_minask", LiveCoinAPIURL)
	sendPayload("GET", path, nil, nil, nil)
}

// Returns the limit for minimum amount to open order, for each
// pair. Also returns maximum number of digits after the decimal point in price value.
func getRestrictions() {
	path := fmt.Sprintf("%s/exchange/restrictions", LiveCoinAPIURL)
	fmt.Println(path)
}

// Returns public data for currencies.
func getCoinInfo() {
	path := fmt.Sprintf("%s/info/coinInfo", LiveCoinAPIURL)
	sendPayload("GET", path, nil, nil, nil)
}

// Open a buy order (limit) for particular currency pair.
func buyLimit(apiKey, apiSecret, currencyPair string, price, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/buylimit", LiveCoinAPIURL)

	p := strconv.FormatFloat(price, 'f', 5, 64)
	q := strconv.FormatFloat(quantity, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("price", p)
	construct.Add("quantity", q)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        apiKey,
		"Sign":           createSignature(message, apiSecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	return data, sendPayload("POST", path, headers, strings.NewReader(message), &data)
}

// Open a sell order (limit) for a specified currency pair.
func sellLimit(apiKey, apiSecret, currencyPair string, price, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/selllimit", LiveCoinAPIURL)

	p := strconv.FormatFloat(price, 'f', 8, 64)
	q := strconv.FormatFloat(quantity, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("price", p)
	construct.Add("quantity", q)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        apiKey,
		"Sign":           createSignature(message, apiSecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	return data, sendPayload("POST", path, headers, strings.NewReader(message), &data)
}

// Open a buy order (market) of specified amount for
// particular currency pair.
func buyMarket(apiKey, apiSecret, currencyPair string, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/buymarket", LiveCoinAPIURL)

	q := strconv.FormatFloat(quantity, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("quantity", q)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        apiKey,
		"Sign":           createSignature(message, apiSecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	return data, sendPayload("POST", path, headers, strings.NewReader(message), &data)
}

// Open a sell order (market) for specified amount of
// selected currency pair.
func sellMarket(apiKey, apiSecret, currencyPair string, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/sellmarket", LiveCoinAPIURL)

	q := strconv.FormatFloat(quantity, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("quantity", q)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        apiKey,
		"Sign":           createSignature(message, apiSecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	return data, sendPayload("POST", path, headers, strings.NewReader(message), &data)
}

// Cancel order.
func cancelLimit(apiKey, apiSecret, currencyPair string, orderID int64) (CancelOrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/cancellimit", LiveCoinAPIURL)

	o := strconv.FormatInt(orderID, 10)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("orderId", o)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        apiKey,
		"Sign":           createSignature(message, apiSecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := CancelOrderResponse{}
	return data, sendPayload("POST", path, headers, strings.NewReader(message), &data)
}

// Returns available balance for selected currency
func getBalance(apiKey, apiSecret, currency string) (BalanceResponse, error) {
	construct := url.Values{}
	construct.Add("currency", currency)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        apiKey,
		"Sign":           createSignature(message, apiSecret),
	}

	path := fmt.Sprintf("%s/payment/balance?%s", LiveCoinAPIURL, message)
	log.Print("-- ", apiKey)
	log.Print("-- ", apiSecret)
	log.Print("-- ", createSignature)
	log.Print("-- ", headers)
	log.Print("-- ", path)

	data := BalanceResponse{}
	return data, sendPayload("GET", path, headers, nil, &data)
}

