package livecoin

import (
	"fmt"
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
)

// BuyLimit - Open a buy order (limit) for particular currency pair.
func (e *LiveCoin) BuyLimit(currencyPair string, price, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/buylimit", LiveCoinAPIURL)

	p := strconv.FormatFloat(price, 'f', 5, 64)
	q := strconv.FormatFloat(quantity, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("price", p)
	construct.Add("quantity", q)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        e.APIKey,
		"Sign":           createSignature(message, e.APISecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	return data, e.SendPayload("GET", path, headers, strings.NewReader(message), &data)
}

// SellLimit - Open a sell order (limit) for a specified currency pair.
func (e *LiveCoin) SellLimit(currencyPair string, price, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/selllimit", LiveCoinAPIURL)

	p := strconv.FormatFloat(price, 'f', 5, 64)
	q := strconv.FormatFloat(quantity, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("price", p)
	construct.Add("quantity", q)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        e.APIKey,
		"Sign":           createSignature(message, e.APISecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	return data, e.SendPayload("GET", path, headers, strings.NewReader(message), &data)
}

// BuyMarket - Open a buy order (market) of specified amount for
// particular currency pair.
func (e *LiveCoin) BuyMarket(currencyPair string, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/buymarket", LiveCoinAPIURL)

	q := strconv.FormatFloat(quantity, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("quantity", q)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        e.APIKey,
		"Sign":           createSignature(message, e.APISecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	return data, e.SendPayload("GET", path, headers, strings.NewReader(message), &data)
}

// SellMarket - Open a sell order (market) for specified amount of
// selected currency pair.
func (e *LiveCoin) SellMarket(currencyPair string, quantity float64) (OrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/sellmarket", LiveCoinAPIURL)

	q := strconv.FormatFloat(quantity, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("quantity", q)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        e.APIKey,
		"Sign":           createSignature(message, e.APISecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := OrderResponse{}
	return data, e.SendPayload("GET", path, headers, strings.NewReader(message), &data)
}

// CancelLimit - Cancel order.
func (e *LiveCoin) CancelLimit(currencyPair string, orderID float64) (CancelOrderResponse, error) {
	path := fmt.Sprintf("%s/exchange/cancellimit", LiveCoinAPIURL)

	o := strconv.FormatFloat(orderID, 'f', 8, 64)
	construct := url.Values{}
	construct.Add("currencyPair", currencyPair)
	construct.Add("orderId", o)
	message := construct.Encode()

	headers := map[string]string{
		"API-Key":        e.APIKey,
		"Sign":           createSignature(message, e.APISecret),
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(message)),
	}

	data := CancelOrderResponse{}
	return data, e.SendPayload("GET", path, headers, strings.NewReader(message), &data)
}
