package livecoin

import (
	"fmt"

	"github.com/labstack/echo"
)

// GetTrades - Get information on your latest transactions. The return may be
// limited by the parameters below.
func (e *liveCoin) GetTrades() {
	path := fmt.Sprintf("%s/exchange/trades", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetClientOrders - Get a detailed review of your orders for requested currency
// pair. You can optionally limit return of orders of a certain type (return
// only open or only closed orders).
func (e *liveCoin) GetClientOrders(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/client_orders", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)

}

// GetOrder - Get the order information by its ID.
func (e *liveCoin) GetOrder(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/order", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetBalances - Returns an array of your balances. There are four types of
// balances for every currency: totaltotal, funds available for trading,
// funds in open orders , funds available for withdrawal.
func (e *liveCoin) GetBalances(currency string) {
	path := fmt.Sprintf("%s/payment/balances", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetBalance - Returns available balance for selected currency.
func (e *liveCoin) GetBalance(currency string) {
	path := fmt.Sprintf("%s/payment/balance", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetHistoryTransactions - Returns a list of your transactions.
func (e *liveCoin) GetHistoryTransactions(params echo.Map) {
	path := fmt.Sprintf("%s/payment/history/transactions", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetHistorySize - Returns the number of transactions with pre-defined
// parameters.
func (e *liveCoin) GetHistorySize(params echo.Map) {
	path := fmt.Sprintf("%s/payment/history/size", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// GetCommission - Returns actual trading fee for customer.
func (e *liveCoin) GetCommission(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/commission", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)

}

// GetCommissionCommonInfo - Returns actual trading fee and volume for the
// last 30 days in USD.
func (e *liveCoin) GetCommissionCommonInfo(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/commissionCommonInfo", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)

}
