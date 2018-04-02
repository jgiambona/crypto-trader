package livecoin

// GetAddress - Get deposit address for selected cryptocurrency.
func GetAddress(currency string) {

}

// PayoutCoin - Submit a request to withdraw cryptocurrency.
func PayoutCoin(amount float64, currency, wallet string) {

}

// PayoutPayeer - Submit a request to withdraw to Payeer account.
func PayoutPayeer(amount float64, currency, protect,
	protectCode string, protectPeriod int64, wallet string) {

}

// PayoutCapitalist - Submit a request to withdraw to Capitalist account.
func PayoutCapitalist(amount float64, currency, wallet string) {

}

// PayoutAdvcash - Submit a request to withdraw to Advcash account.
func PayoutAdvcash(amout float64, currency, wallet string) {

}

// PayoutCard - Submit a request to withdraw to a bank card.
func PayoutCard(amount float64, currency, cardNumber,
	expiryMonth, expiryYear string) {

}

// PayoutOkpay - Submit a request to withdraw to Okpay account.
func PayoutOkpay(amount float64, currency, invoice, wallet string) {

}

// PayoutPerfectMoney - Submit a request to withdraw to PerfectMoney account.
func PayoutPerfectMoney(amount float64, currency, protectCode,
	wallet string, protectPeriod int64) {

}
