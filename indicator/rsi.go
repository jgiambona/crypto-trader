package indicator

// RSIOptions to configure the behaviour of RSI.
type RSIOptions struct {
	period        string
	periodLength  string
	minPeriods    int64
	rsiPeriods    int64
	oversoldRSI   int64
	overBoughtRSI int64
	rsiRecover    int64
	rsiDrop       int64
	rsiDivisor    int64
}

func rsi(s interface{}, key string, length int64) {
	//if len(s.lookback) >= length {
	//	avgGain := s.lookback[0]
	//	avgLoss := s.lookback[0]
	//}
}
