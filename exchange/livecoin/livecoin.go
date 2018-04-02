package livecoin

import (
	"github.com/ffimnsr/trader/exchange"
)

// Base API URL.
const (
	liveCoinAPIURL = "https://api.livecoin.net"
)

// The API error codes that are being returned.
const (
	UnknownError             = 1
	SystemError              = 2
	AuthenticationError      = 10
	AuthenticationIsRequired = 11
	AuthenticationFailed     = 12
	SignatureIncorrect       = 20
	AccessDenied             = 30
	APIDisabled              = 31
	APIRestrictedByIP        = 32
	IncorrectParameters      = 100
	IncorrectAPIKey          = 101
	IncorrectUserID          = 102
	IncorrectCurrency        = 103
	IncorrectAmount          = 104
	UnableToBlockFunds       = 150
)

type (
	// LiveCoin is an interface to LiveCoin rest API.
	LiveCoin interface {
	}

	liveCoin struct {
		exchange.Base
	}
)

// NewInstance creates an instance of LiveCoin Callable exchange API.
func NewInstance() LiveCoin {
	exchange := new(liveCoin)
	exchange.Name = "hello"
	exchange.Enabled = true
	return exchange
}
