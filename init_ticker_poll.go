package main

import (
	"fmt"
	"time"
)

// PollTicker fetches and updates the ticker for all exchanges.
func PollTicker() {
	for {
		fmt.Println("Poll Ticker")

		time.Sleep(8 * time.Second)
	}
}
