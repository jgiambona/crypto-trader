package livecoin

import (
	"fmt"

	"github.com/labstack/echo"
)

// BuyLimit - Open a buy order (limit) for particular currency pair.
func (e *LiveCoin) BuyLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/buylimit", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// SellLimit - Open a sell order (limit) for a specified currency pair.
func (e *LiveCoin) SellLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// BuyMarket - Open a buy order (market) of specified amount for
// particular currency pair.
func (e *LiveCoin) BuyMarket(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// SellMarket - Open a sell order (market) for specified amount of
// selected currency pair.
func (e *LiveCoin) SellMarket(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// CancelLimit - Cancel order.
func (e *LiveCoin) CancelLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}
