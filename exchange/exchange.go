package exchange

import (
	"time"

	"github.com/labstack/echo"
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
		GetName() string
		IsEnabled() bool
		UpdateTicker() echo.Map
	}
)

// GetName returns the name of the exchange.
func (o *Base) GetName() string {
	return o.Name
}

// IsEnabled is a method that returns if the exchange
// is enabled.
func (o *Base) IsEnabled() bool {
	return o.Enabled
}
