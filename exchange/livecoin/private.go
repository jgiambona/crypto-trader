package livecoin

import (
	"fmt"

	"github.com/labstack/echo"
)

// GetTrades - Get information on your latest transactions. The return may be
// limited by the parameters below.
func GetTrades() {
	path := fmt.Sprintf("%s/exchange/trades", liveCoinAPIURL)
	fmt.Println(path)
}

// GetClientOrders - Get a detailed review of your orders for requested currency
// pair. You can optionally limit return of orders of a certain type (return
// only open or only closed orders).
func GetClientOrders(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/client_orders", liveCoinAPIURL)
	fmt.Println(path)

}

// GetOrder - Get the order information by its ID.
func GetOrder(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/order", liveCoinAPIURL)
	fmt.Println(path)

}

// GetBalances - Returns an array of your balances. There are four types of
// balances for every currency: totaltotal, funds available for trading,
// funds in open orders , funds available for withdrawal.
func GetBalances(currency string) {
	path := fmt.Sprintf("%s/payment/balances", liveCoinAPIURL)
	fmt.Println(path)
}

// GetBalance - Returns available balance for selected currency.
func GetBalance(currency string) {
	path := fmt.Sprintf("%s/payment/balance", liveCoinAPIURL)
	fmt.Println(path)
}

// GetHistoryTransactions - Returns a list of your transactions.
func GetHistoryTransactions(params echo.Map) {
	path := fmt.Sprintf("%s/payment/history/transactions", liveCoinAPIURL)
	fmt.Println(path)
}

// GetHistorySize - Returns the number of transactions with pre-defined
// parameters.
func GetHistorySize(params echo.Map) {
	path := fmt.Sprintf("%s/payment/history/size", liveCoinAPIURL)
	fmt.Println(path)
}

// GetCommission - Returns actual trading fee for customer.
func GetCommission(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/commission", liveCoinAPIURL)
	fmt.Println(path)

}

// GetCommissionCommonInfo - Returns actual trading fee and volume for the
// last 30 days in USD.
func GetCommissionCommonInfo(params echo.Map) {
	path := fmt.Sprintf("%s/exchange/commissionCommonInfo", liveCoinAPIURL)
	fmt.Println(path)

}
