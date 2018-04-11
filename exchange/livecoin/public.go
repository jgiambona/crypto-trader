package livecoin

import (
	"fmt"
	"log"
)

// GetTicker - Get information on specified currency pair for the last 24 hours.
func (e *LiveCoin) GetTicker(currencyPair string) Price {
	path := fmt.Sprintf("%s/exchange/ticker?currencyPair=%s", LiveCoinAPIURL, currencyPair)

	result := Price{}
	if err := e.SendPayload("GET", path, nil, nil, &result); err != nil {
		log.Fatalf("%s", err.Error())
	}
	return result
}

// GetLastTrades - Get a detailed review on the latest transactions for
// requested currency pair. You may receive the update for the last hour or
// for the last minute.
func (e *LiveCoin) GetLastTrades(currencyPair, minutesOrHour, tradeType string) {
	path := fmt.Sprintf("%s/exchange/last_trades", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetOrderBook - Get the orderbook for specified currency pair (you may
// enable the feature of grouping orders by price).
func (e *LiveCoin) GetOrderBook(currencyPair, groupByPrice string, depth int64) {
	path := fmt.Sprintf("%s/exchange/order_book", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetAllOrderBook - Returns orderbook for every currency pair.
func (e *LiveCoin) GetAllOrderBook(groupByPrice string, depth int64) {
	path := fmt.Sprintf("%s/exchange/all/order_book", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetMaxBidMinAsk - Returns maximum bid and minimum ask in the current orderbook.
func (e *LiveCoin) GetMaxBidMinAsk(currencyPair string) {
	path := fmt.Sprintf("%s/exchange/maxbid_minask", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetRestrictions - Returns the limit for minimum amount to open order, for each
// pair. Also returns maximum number of digits after the decimal point in price value.
func (e *LiveCoin) GetRestrictions() {
	path := fmt.Sprintf("%s/exchange/restrictions", LiveCoinAPIURL)
	fmt.Println(path)
}

// GetCoinInfo - Returns public data for currencies.
func (e *LiveCoin) GetCoinInfo() {
	path := fmt.Sprintf("%s/info/coinInfo", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}
