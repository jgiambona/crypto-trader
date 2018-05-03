package main

import (
	"log"
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
		MaximumVolume         float64
		TransactionVolume     float64
		VarianceOfTransaction float64
		BidPriceStepDown      float64
		MinimumBid            float64
		Enabled               bool
	}
)

func botControls(c echo.Context) error {
	power, err := strconv.ParseInt(c.FormValue("power"), 10, 64)
	if err != nil {
		jsonBadRequest(c, err)
	}

	if power == 1 {
		bot.running = true
		insertBotStatus("ON")
	} else {
		bot.running = false
		insertBotStatus("OFF")
	}
	return jsonSuccess(c, echo.Map{})
}

func updateAccounts(c echo.Context) error {
	id, _ := strconv.ParseInt(c.FormValue("id"), 10, 64)
	if id == 1 {
		bot.accountOne.APIKey = c.FormValue("key")
		bot.accountOne.APISecret = c.FormValue("secret")
	} else if id == 2 {
		bot.accountTwo.APIKey = c.FormValue("key")
		bot.accountTwo.APISecret = c.FormValue("secret")
	} else {
		return jsonBadRequest(c, "error no such account.")
	}

	log.Printf(`
	----- One
	--- %s
	--- %s
	----- Two
	--- %s
	--- %s
	----- End`,
		bot.accountOne.APIKey, bot.accountOne.APISecret,
		bot.accountTwo.APIKey, bot.accountTwo.APISecret)

	return jsonSuccess(c, echo.Map{
		"account": id,
		"key":     c.FormValue("key"),
	})
}

func updateSettings(c echo.Context) error {
	id, _ := strconv.ParseInt(c.FormValue("id"), 10, 64)
	interval, _ := strconv.ParseInt(c.FormValue("interval"), 10, 64)
	maximumVolume, _ := strconv.ParseFloat(c.FormValue("maximumVolume"), 64)
	transactionVolume, _ := strconv.ParseFloat(c.FormValue("transactionVolume"), 64)
	variance, _ := strconv.ParseFloat(c.FormValue("variance"), 64)
	stepDownPrice, _ := strconv.ParseFloat(c.FormValue("stepDownPrice"), 64)
	minimumBid, _ := strconv.ParseFloat(c.FormValue("minimumBid"), 64)

	if id == 1 {
		bot.ruleOne.Enabled = true
		bot.ruleOne.Interval = interval
		bot.ruleOne.MaximumVolume = maximumVolume
		bot.ruleOne.TransactionVolume = transactionVolume
		bot.ruleOne.VarianceOfTransaction = variance
		bot.ruleOne.BidPriceStepDown = stepDownPrice
		bot.ruleOne.MinimumBid = minimumBid
	} else if id == 2 {
		bot.ruleTwo.Enabled = true
		bot.ruleTwo.Interval = interval
		bot.ruleTwo.MaximumVolume = maximumVolume
		bot.ruleTwo.TransactionVolume = transactionVolume
		bot.ruleTwo.VarianceOfTransaction = variance
		bot.ruleTwo.BidPriceStepDown = stepDownPrice
		bot.ruleTwo.MinimumBid = minimumBid
	} else {
		return jsonBadRequest(c, "error no such account.")
	}

	log.Printf(`
	----- One
	--- %d
	--- %.8f
	--- %.8f
	--- %.8f
	--- %.8f
	--- %.8f
	----- Two
	--- %d
	--- %.8f
	--- %.8f
	--- %.8f
	--- %.8f
	--- %.8f
	----- End`,
		bot.ruleOne.Interval, bot.ruleOne.MaximumVolume, bot.ruleOne.TransactionVolume,
		bot.ruleOne.VarianceOfTransaction, bot.ruleOne.BidPriceStepDown,
		bot.ruleOne.MinimumBid,
		bot.ruleTwo.Interval, bot.ruleTwo.MaximumVolume, bot.ruleTwo.TransactionVolume,
		bot.ruleTwo.VarianceOfTransaction, bot.ruleTwo.BidPriceStepDown,
		bot.ruleTwo.MinimumBid)

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
