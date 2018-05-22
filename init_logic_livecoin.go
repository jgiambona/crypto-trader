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

	OrderDetailResponse struct {
		ID                int64   `json:"id"`
		ClientID          int64   `json:"client_id"`
		Status            string  `json:"status"`
		Symbol            string  `json:"symbol"`
		Price             float64 `json:"price"`
		Quantity          float64 `json:"quantity"`
		RemainingQuantity float64 `json:"remaining_quantity"`
		Blocked           float64 `json:"blocked"`
		BlockedRemain     float64 `json:"blocked_remain"`
		CommissionRate    float64 `json:"commission_rate"`
		Trades            []int64 `json:"trades"`
	}

	OrderBookResponse struct {
		Timestamp int64      `json:"timestamp"`
		Asks      [][]string `json:"asks"`
		Bids      [][]string `json:"bids"`
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
	Fee     float64
}, error) {
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
	log.Print("-- ", path)
	return result, sendPayload("GET", path, nil, nil, &result)
}

// Get the order information by its ID.
func getOrder(apiKey, apiSecret, orderId string) (OrderDetailResponse, error) {
	construct := url.Values{}
	construct.Add("orderId", orderId)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key": apiKey,
		"Sign":    createSignature(message, apiSecret),
	}

	path := fmt.Sprintf("%s/exchange/order?%s", LiveCoinAPIURL, message)
	log.Print("-- ", path)
	data := OrderDetailResponse{}
	return data, sendPayload("GET", path, headers, nil, &data)
}

// Get the order book list.
func getOrderBook(currencyPair string) (OrderBookResponse, error) {
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("depth", "4")
	construct.Add("groupByPrice", "true")
	message := construct.Encode()

	path := fmt.Sprintf("%s/exchange/order_book?%s", LiveCoinAPIURL, message)
	log.Print("-- ", path)
	data := OrderBookResponse{}
	return data, sendPayload("GET", path, nil, nil, &data)
}

// Open a buy order (limit) for particular currency pair.
func buyLimit(apiKey, apiSecret, currencyPair string, price, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/buylimit", LiveCoinAPIURL)

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

	log.Print("-- ", path)
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

	log.Print("-- ", path)
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
		"API-Key": apiKey,
		"Sign":    createSignature(message, apiSecret),
	}

	path := fmt.Sprintf("%s/payment/balance?%s", LiveCoinAPIURL, message)
	log.Print("-- ", path)
	data := BalanceResponse{}
	return data, sendPayload("GET", path, headers, nil, &data)
}
