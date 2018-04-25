package strategy

import talib "github.com/markcheno/go-talib"

// Buy when (EMA - last(EMA) > 0) and sell when (EMA - last(EMA) < 0). Optional
// buy on low RSI.
func TrendEMA() {
}
