package main

var strategies = map[string]string{
	"uptrend":      "Uptrend",
	"bb":           "Bollinger Bands",
	"gain":         "Gain",
	"pp":           "Pingpong",
	"stepgain":     "Stepgain",
	"tssl":         "Trailing Stop / Stop Limit",
	"emotionless":  "Emotionless",
	"ichimoku":     "Ichimoku",
	"tsslbb":       "Trailing Stop / Stop Limit - Bollinger Bands",
	"tsslpp":       "Trailing Stop / Stop Limit - Pingpong",
	"tsslstepgain": "Trailing Stop / Stop Limit - Stepgain",
	"tsslgain":     "Trailing Stop / Stop Limit - Gain",
	"bbrsitssl":    "Bollinger Bands + RSI - Trailing Stop / Stop Limit",
	"pptssl":       "Pingpong - Trailing Stop / Stop Limit",
	"stepgaintssl": "Stepgain - Trailing Stop / Stop Limit",
	"gaintssl":     "Gain - Trailing Stop / Stop Limit",
	"bbtssl":       "Bollinger Bands - Trailing Stop / Stop Limit",
	"bbgain":       "Bollinger Bands - Gain",
	"gainbb":       "Gain - Bollinger Bands",
	"bbstepgain":   "Bollinger Bands - Stepgain",
	"stepgainbb":   "Stepgain - Bollinger Bands",
	"bbpp":         "Bollinger Bands - Pingpong",
	"ppbb":         "Pingpong - Bollinger Bands",
	"gainstepgain": "Gain - Stepgain",
	"stepgaingain": "Stepgain - Gain",
	"gainpp":       "Gain - Pingpong",
	"stepgainpp":   "Stepgain - Pingpong",
	"ppstepgain":   "Pingpong - Stepgain",
}

