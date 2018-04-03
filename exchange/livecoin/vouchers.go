package livecoin

import (
	"fmt"
)

// MakeVoucher - Creates a new voucher.
func (e *liveCoin) MakeVoucher(amount float64, currency, description string) {
	path := fmt.Sprintf("%s/payment/voucher/make", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// CheckVoucherAmount - Returns a voucher amount upon its code.
func (e *liveCoin) CheckVoucherAmount(voucherCode string) {
	path := fmt.Sprintf("%s/payment/voucher/amount", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}

// RedeemVoucher - Redeem a voucher.
func (e *liveCoin) RedeemVoucher(voucherCode string) {
	path := fmt.Sprintf("%s/payment/voucher/redeem", liveCoinAPIURL)
	e.SendPayload("GET", path, nil, nil)
}
