package livecoin

import (
	"fmt"
)

// MakeVoucher - Creates a new voucher.
func MakeVoucher(amount float64, currency, description string) {
	path := fmt.Sprintf("%s/payment/voucher/make", liveCoinAPIURL)
	fmt.Println(path)
}

// CheckVoucherAmount - Returns a voucher amount upon its code.
func CheckVoucherAmount(voucherCode string) {
	path := fmt.Sprintf("%s/payment/voucher/amount", liveCoinAPIURL)
	fmt.Println(path)
}

// RedeemVoucher - Redeem a voucher.
func RedeemVoucher(voucherCode string) {
	path := fmt.Sprintf("%s/payment/voucher/redeem", liveCoinAPIURL)
	fmt.Println(path)

}
