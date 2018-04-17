package livecoin

import (
	"fmt"
	"log"
)

// GetAddress - Get deposit address for selected cryptocurrency.
// Not Use
func (e *LiveCoin) GetAddress(currency string) {
	path := fmt.Sprintf("%s/payment/get/address", LiveCoinAPIURL)
	log.Print(path)
}

// PayoutCoin - Submit a request to withdraw cryptocurrency.
// Not Use
func (e *LiveCoin) PayoutCoin(amount float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/coin", LiveCoinAPIURL)
	log.Print(path)
}

// PayoutPayeer - Submit a request to withdraw to Payeer account.
// Not Use
func (e *LiveCoin) PayoutPayeer(amount float64, currency, protect,
	protectCode string, protectPeriod int64, wallet string) {
	path := fmt.Sprintf("%s/payment/out/payeer", LiveCoinAPIURL)
	log.Print(path)
}

// PayoutCapitalist - Submit a request to withdraw to Capitalist account.
// Not Use
func (e *LiveCoin) PayoutCapitalist(amount float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/capitalist", LiveCoinAPIURL)
	log.Print(path)
}

// PayoutAdvcash - Submit a request to withdraw to Advcash account.
// Not Use
func (e *LiveCoin) PayoutAdvcash(amout float64, currency, wallet string) {
	path := fmt.Sprintf("%s/payment/out/advcash", LiveCoinAPIURL)
	log.Print(path)
}

// PayoutCard - Submit a request to withdraw to a bank card.
// Not Use
func (e *LiveCoin) PayoutCard(amount float64, currency, cardNumber,
	expiryMonth, expiryYear string) {
	path := fmt.Sprintf("%s/payment/out/card", LiveCoinAPIURL)
	log.Print(path)
}

// PayoutOkpay - Submit a request to withdraw to Okpay account.
// Not Use
func (e *LiveCoin) PayoutOkpay(amount float64, currency, invoice, wallet string) {
	path := fmt.Sprintf("%s/payment/out/okpay", LiveCoinAPIURL)
	log.Print(path)
}

// PayoutPerfectMoney - Submit a request to withdraw to PerfectMoney account.
// Not Use
func (e *LiveCoin) PayoutPerfectMoney(amount float64, currency, protectCode,
	wallet string, protectPeriod int64) {
	path := fmt.Sprintf("%s/payment/out/perfectmoney", LiveCoinAPIURL)
	log.Print(path)
}
