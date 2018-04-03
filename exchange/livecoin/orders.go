package livecoin

import (
	"fmt"

	"github.com/labstack/echo"
)

// BuyLimit - Open a buy order (limit) for particular currency pair.
func (e *liveCoin) BuyLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// SellLimit - Open a sell order (limit) for a specified currency pair.
func (e *liveCoin) SellLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// BuyMarket - Open a buy order (market) of specified amount for
// particular currency pair.
func (e *liveCoin) BuyMarket(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// SellMarket - Open a sell order (market) for specified amount of
// selected currency pair.
func (e *liveCoin) SellMarket(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// CancelLimit - Cancel order.
func (e *liveCoin) CancelLimit(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/selllimit", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}
