package main

import "github.com/ffimnsr/trader/exchange/livecoin"

// LoadExchanges initialize all exchange that are available
// on the platform.
func LoadExchanges() {
	lc = livecoin.NewInstance()
}
