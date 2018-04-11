package livecoin

import (
	"fmt"
)

// MakeVoucher - Creates a new voucher.
func (e *LiveCoin) MakeVoucher(amount float64, currency, description string) {
	path := fmt.Sprintf("%s/payment/voucher/make", LiveCoinAPIURL)
	e.SendPayload("POST", path, nil, nil, nil)
}

// CheckVoucherAmount - Returns a voucher amount upon its code.
func (e *LiveCoin) CheckVoucherAmount(voucherCode string) {
	path := fmt.Sprintf("%s/payment/voucher/amount", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}

// RedeemVoucher - Redeem a voucher.
func (e *LiveCoin) RedeemVoucher(voucherCode string) {
	path := fmt.Sprintf("%s/payment/voucher/redeem", LiveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil, nil)
}
