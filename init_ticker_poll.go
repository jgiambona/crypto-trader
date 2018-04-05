package main

import (
	"time"
)

// PollTicker fetches and updates the ticker for all exchanges.
func pollTicker() {
	for {
		time.Sleep(8 * time.Second)
	}
}
