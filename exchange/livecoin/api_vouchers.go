package livecoin

import (
	"fmt"
	"log"
)

// MakeVoucher - Creates a new voucher.
// Not Use
func (e *LiveCoin) MakeVoucher(amount float64, currency, description string) {
	path := fmt.Sprintf("%s/payment/voucher/make", LiveCoinAPIURL)
	log.Print(path)
}

// CheckVoucherAmount - Returns a voucher amount upon its code.
// Not Use
func (e *LiveCoin) CheckVoucherAmount(voucherCode string) {
	path := fmt.Sprintf("%s/payment/voucher/amount", LiveCoinAPIURL)
	log.Print(path)
}

// RedeemVoucher - Redeem a voucher.
// Not Use
func (e *LiveCoin) RedeemVoucher(voucherCode string) {
	path := fmt.Sprintf("%s/payment/voucher/redeem", LiveCoinAPIURL)
	log.Print(path)
}
