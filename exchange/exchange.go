package exchange

import (
	"time"
)

type (
	// Base holds common features that identifies
	// the exchange platform.
	Base struct {
		Name         string
		APIBaseURL   string
		APIKey       string
		APISecret    string
		PollingDelay time.Duration
		Enabled      bool
	}

	// BotExchange TODO
	BotExchange interface {
		IsEnabled() bool
		// GetTickerPrice()
		// UpdateTicker()
	}
)

// NewExchange TODO
func NewExchange(name string, apiURL, apiKey, apiSecret string,
	pollingDelay time.Duration, enabled bool) BotExchange {
	return &Base{
		Name:         name,
		APIBaseURL:   apiURL,
		APIKey:       apiKey,
		APISecret:    apiSecret,
		PollingDelay: pollingDelay,
		Enabled:      enabled,
	}
}

// IsEnabled is a method that returns if the exchange
// is enabled.
func (o *Base) IsEnabled() bool {
	return o.Enabled
}
