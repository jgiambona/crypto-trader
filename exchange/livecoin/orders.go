package livecoin

import (
	"fmt"

	"github.com/labstack/echo"
)

// BuyLimit - Open a buy order (limit) for particular currency pair.
func BuyLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	fmt.Println(path)
}

// SellLimit - Open a sell order (limit) for a specified currency pair.
func SellLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	fmt.Println(path)
}

// BuyMarket - Open a buy order (market) of specified amount for
// particular currency pair.
func BuyMarket(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	fmt.Println(path)
}

// SellMarket - Open a sell order (market) for specified amount of
// selected currency pair.
func SellMarket(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	fmt.Println(path)
}

// CancelLimit - Cancel order.
func CancelLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	fmt.Println(path)
}
