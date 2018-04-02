package exchange

type (
	// Base holds common features that identifies
	// the exchange platform.
	Base struct {
		Name    string
		Enabled bool
	}
)

// IsEnabled is a method that returns if the exchange
// is enabled.
func (o *Base) IsEnabled() bool {
	return o.Enabled
}
