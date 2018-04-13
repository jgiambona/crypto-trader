package strategy

func vwap(s, key string, length, maxPeriod int64, sourceKey string) {
	if len(sourceKey) < 0 {
		sourceKey = "close"
	}

	// if s.lookbac[k.length >= length {
	// 	if !s.vwap {
	// 		s.vwap = 0
	// 		s.vwapMultiplier = 0
	// 		s.vwapDivider = 0
	// 		s.vwapCount = 0
	// 	}

	// 	if maxPeriod && s.vwapCount > maxPeriod {
	// 		s.vwap = 0
	// 		s.vwapMultiplier = 0
	// 		s.vwapDivider = 0
	// 		s.vwapCount = 0
	// 	}

	// 	s.vwapMultiplier = s.vwapMultiplier + parseFloat(s.period[sourceKey]) * parseFloat(s.period["volume"])
	// 	s.vwapDivider = s.vwapDivider + parseFloat(s.period["volume"])
	// 	s.period[key] = s.vwap = s.vwapMultiplier / s.vwapDivider
	// 	s.vwapCount += 1
	// }
}
