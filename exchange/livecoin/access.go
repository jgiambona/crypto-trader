package livecoin

import "fmt"

// GetAddress - Get deposit address for selected cryptocurrency.
func (e *LiveCoin) GetAddress(currency string) {
	path := fmt.Sprintf("%s/payment/get/address", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil)
}

// PayoutCoin - Submit a request to withdraw cryptocurrency.
func (e *LiveCoin) PayoutCoin(amount float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/coin", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil)
}

// PayoutPayeer - Submit a request to withdraw to Payeer account.
func (e *LiveCoin) PayoutPayeer(amount float64, currency, protect,
	protectCode string, protectPeriod int64, wallet string) {
	path := fmt.Sprintf("%s/payment/out/payeer", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil)
}

// PayoutCapitalist - Submit a request to withdraw to Capitalist account.
func (e *LiveCoin) PayoutCapitalist(amount float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/capitalist", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil)
}

// PayoutAdvcash - Submit a request to withdraw to Advcash account.
func (e *LiveCoin) PayoutAdvcash(amout float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/advcash", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil)
}

// PayoutCard - Submit a request to withdraw to a bank card.
func (e *LiveCoin) PayoutCard(amount float64, currency, cardNumber,
	expiryMonth, expiryYear string) {
	path := fmt.Sprintf("%s/payment/out/card", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil)
}

// PayoutOkpay - Submit a request to withdraw to Okpay account.
func (e *LiveCoin) PayoutOkpay(amount float64, currency, invoice, wallet string) {
	path := fmt.Sprintf("%s/payment/out/okpay", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil)
}

// PayoutPerfectMoney - Submit a request to withdraw to PerfectMoney account.
func (e *LiveCoin) PayoutPerfectMoney(amount float64, currency, protectCode,
	wallet string, protectPeriod int64) {
	path := fmt.Sprintf("%s/payment/out/perfectmoney", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil)
}
