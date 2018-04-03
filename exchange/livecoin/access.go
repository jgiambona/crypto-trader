package livecoin

import "fmt"

// GetAddress - Get deposit address for selected cryptocurrency.
func (e *liveCoin) GetAddress(currency string) {
	path := fmt.Sprintf("%s/payment/get/address", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// PayoutCoin - Submit a request to withdraw cryptocurrency.
func (e *liveCoin) PayoutCoin(amount float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/coin", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// PayoutPayeer - Submit a request to withdraw to Payeer account.
func (e *liveCoin) PayoutPayeer(amount float64, currency, protect,
	protectCode string, protectPeriod int64, wallet string) {
	path := fmt.Sprintf("%s/payment/out/payeer", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// PayoutCapitalist - Submit a request to withdraw to Capitalist account.
func (e *liveCoin) PayoutCapitalist(amount float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/capitalist", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// PayoutAdvcash - Submit a request to withdraw to Advcash account.
func (e *liveCoin) PayoutAdvcash(amout float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/advcash", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// PayoutCard - Submit a request to withdraw to a bank card.
func (e *liveCoin) PayoutCard(amount float64, currency, cardNumber,
	expiryMonth, expiryYear string) {
	path := fmt.Sprintf("%s/payment/out/card", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// PayoutOkpay - Submit a request to withdraw to Okpay account.
func (e *liveCoin) PayoutOkpay(amount float64, currency, invoice, wallet string) {
	path := fmt.Sprintf("%s/payment/out/okpay", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// PayoutPerfectMoney - Submit a request to withdraw to PerfectMoney account.
func (e *liveCoin) PayoutPerfectMoney(amount float64, currency, protectCode,
	wallet string, protectPeriod int64) {
	path := fmt.Sprintf("%s/payment/out/perfectmoney", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}
