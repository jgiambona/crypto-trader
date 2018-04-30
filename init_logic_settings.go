package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type (
	Account struct {
		APIKey    string
		APISecret string
	}

	RuleConfiguration struct {
		ID                    int64
		Interval              int64
		MaximumVolume         int64
		TransactionVolume     int64
		VarianceOfTransaction float64
		BidPriceStepDown      float64
		MinimumBid            float64
		Enabled               bool
	}
)

func updateAccounts(c echo.Context) error {
	id, _ := strconv.ParseInt(c.FormValue("id"), 10, 64)
	if id == 1 {
		bot.accountOne.A.APIKey = c.FormValue("key")
		bot.accountOne.A.APISecret = c.FormValue("secret")
	} else if id == 2 {
		bot.accountTwo.A.APIKey = c.FormValue("key")
		bot.accountTwo.A.APISecret = c.FormValue("secret")
	} else {
		return jsonBadRequest(c, "error no such account.")
	}
	return jsonSuccess(c, echo.Map{
		"account": id,
		"key":     c.FormValue("key"),
	})
}

func updateSettings(c echo.Context) error {
	id, _ := strconv.ParseInt(c.FormValue("id"), 10, 64)
	interval := c.FormValue("interval")
	maximumVolume := c.FormValue("maximumVolume")
	transactionVolume := c.FormValue("transactionVolume")
	variance := c.FormValue("variance")
	stepDownPrice := c.FormValue("stepDownPrice")
	minimumBid := c.FormValue("minimumBid")

	if id == 1 {
		bot.accountOne.R.Enabled = true
		bot.accountOne.R.Interval = interval
		bot.accountOne.R.MaximumVolume = maximumVolume
		bot.accountOne.R.TransactionVolume = transactionVolume
		bot.accountOne.R.VarianceOfTransaction = variance
		bot.accountOne.R.BidPriceStepDown = stepDownPrice
		bot.accountOne.R.MinimumBid = minimumBid
	} else if id == 2 {
		bot.accountTwo.R.Enabled = true
		bot.accountTwo.R.Interval = interval
		bot.accountTwo.R.MaximumVolume = maximumVolume
		bot.accountTwo.R.TransactionVolume = transactionVolume
		bot.accountTwo.R.VarianceOfTransaction = variance
		bot.accountTwo.R.BidPriceStepDown = stepDownPrice
		bot.accountTwo.R.MinimumBid = minimumBid
	} else {
		return jsonBadRequest(c, "error no such account.")
	}

	return jsonSuccess(c, echo.Map{
		"account": id,
	})
}

func jsonSuccess(c echo.Context, o echo.Map) error {
	o["success"] = true
	return c.JSON(http.StatusOK, o)
}

func jsonBadRequest(c echo.Context, i interface{}) error {
	return c.JSON(http.StatusBadRequest, echo.Map{
		"success": false,
		"message": i,
	})
}

func jsonServerError(c echo.Context, i interface{}) error {
	return c.JSON(http.StatusInternalServerError, echo.Map{
		"success": false,
		"message": i,
	})
}
