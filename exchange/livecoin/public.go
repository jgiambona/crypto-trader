package livecoin

import (
	"fmt"
)

// GetTicker - Get information on specified currency pair for the last 24 hours.
func (e *liveCoin) GetTicker(currencyPair string) {
	path := fmt.Sprintf("%s/exchange/ticker", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetLastTrades - Get a detailed review on the latest transactions for
// requested currency pair. You may receive the update for the last hour or
// for the last minute.
func (e *liveCoin) GetLastTrades(currencyPair, minutesOrHour, tradeType string) {
	path := fmt.Sprintf("%s/exchange/last_trades", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetOrderBook - Get the orderbook for specified currency pair (you may
// enable the feature of grouping orders by price).
func (e *liveCoin) GetOrderBook(currencyPair, groupByPrice string, depth int64) {
	path := fmt.Sprintf("%s/exchange/order_book", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetAllOrderBook - Returns orderbook for every currency pair.
func (e *liveCoin) GetAllOrderBook(groupByPrice string, depth int64) {
	path := fmt.Sprintf("%s/exchange/all/order_book", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetMaxBidMinAsk - Returns maximum bid and minimum ask in the current orderbook.
func (e *liveCoin) GetMaxBidMinAsk(currencyPair string) {
	path := fmt.Sprintf("%s/exchange/maxbid_minask", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetRestrictions - Returns the limit for minimum amount to open order, for each
// pair. Also returns maximum number of digits after the decimal point in price value.
func (e *liveCoin) GetRestrictions() {
	path := fmt.Sprintf("%s/exchange/restrictions", liveCoinAPIURL)
	fmt.Println(path)
}

// GetCoinInfo - Returns public data for currencies.
func (e *liveCoin) GetCoinInfo() {
	path := fmt.Sprintf("%s/info/coinInfo", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}
