package livecoin

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetTicker - Get information on specified currency pair for the last 24 hours.
func GetTicker(currencyPair string) {
	path := fmt.Sprintf("%s/exchange/ticker", liveCoinAPIURL)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println("response body:", string(body))
}

// GetLastTrades - Get a detailed review on the latest transactions for
// requested currency pair. You may receive the update for the last hour or
// for the last minute.
func GetLastTrades(currencyPair, minutesOrHour, tradeType string) {
	path := fmt.Sprintf("%s/exchange/last_trades", liveCoinAPIURL)
	fmt.Println(path)

}

// GetOrderBook - Get the orderbook for specified currency pair (you may
// enable the feature of grouping orders by price).
func GetOrderBook(currencyPair, groupByPrice string, depth int64) {
	path := fmt.Sprintf("%s/exchange/order_book", liveCoinAPIURL)
	fmt.Println(path)

}

// GetAllOrderBook - Returns orderbook for every currency pair.
func GetAllOrderBook(groupByPrice string, depth int64) {
	path := fmt.Sprintf("%s/exchange/all/order_book", liveCoinAPIURL)
	fmt.Println(path)

}

// GetMaxBidMinAsk - Returns maximum bid and minimum ask in the current orderbook.
func GetMaxBidMinAsk(currencyPair string) {
	path := fmt.Sprintf("%s/exchange/maxbid_minask", liveCoinAPIURL)
	fmt.Println(path)

}

// GetRestrictions - Returns the limit for minimum amount to open order, for each
// pair. Also returns maximum number of digits after the decimal point in price value.
func GetRestrictions() {
	path := fmt.Sprintf("%s/exchange/restrictions", liveCoinAPIURL)
	fmt.Println(path)

}

// GetCoinInfo - Returns public data for currencies.
func GetCoinInfo() {
	path := fmt.Sprintf("%s/info/coinInfo", liveCoinAPIURL)
	fmt.Println(path)

}
