package main

import "github.com/labstack/echo"

type RuleConfiguration struct {
	ID                    int64
	Interval              int64
	MaximumVolume         int64
	TransactionVolume     int64
	VarianceOfTransaction float64
	BidPriceStepDown      float64
	Enabled               bool
}

func updateSettings(c echo.Context) error {
	//interval := c.FormValue("interval")
	//maxVolume := c.FormValue("maxVolume")
	//variance := c.FormValue("variance")
	//stepDownPrice := c.FormValue("stepDownPrice")

	return nil
}

