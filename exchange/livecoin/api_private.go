package livecoin

import (
	"fmt"

	"github.com/labstack/echo"
)

// GetTrades - Get information on your latest transactions. The return may be
// limited by the parameters below.
func (e *LiveCoin) GetTrades() {
	path := fmt.Sprintf("%s/exchange/trades", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetClientOrders - Get a detailed review of your orders for requested currency
// pair. You can optionally limit return of orders of a certain type (return
// only open or only closed orders).
func (e *LiveCoin) GetClientOrders(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/client_orders", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetOrder - Get the order information by its ID.
func (e *LiveCoin) GetOrder(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/order", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetBalances - Returns an array of your balances. There are four types of
// balances for every currency: total, funds available for trading,
// funds in open orders, funds available for withdrawal.
func (e *LiveCoin) GetBalances(currency string) {
	path := fmt.Sprintf("%s/payment/balances", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetBalance - Returns available balance for selected currency.
func (e *LiveCoin) GetBalance(currency string) {
	path := fmt.Sprintf("%s/payment/balance", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetHistoryTransactions - Returns a list of your transactions.
func (e *LiveCoin) GetHistoryTransactions(params echo.Map) {
	path := fmt.Sprintf("%s/payment/history/transactions", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetHistorySize - Returns the number of transactions with pre-defined
// parameters.
func (e *LiveCoin) GetHistorySize(params echo.Map) {
	path := fmt.Sprintf("%s/payment/history/size", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetCommission - Returns actual trading fee for customer.
func (e *LiveCoin) GetCommission(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/commission", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// GetCommissionCommonInfo - Returns actual trading fee and volume for the
// last 30 days in USD.
func (e *LiveCoin) GetCommissionCommonInfo(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/commissionCommonInfo", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)

}
