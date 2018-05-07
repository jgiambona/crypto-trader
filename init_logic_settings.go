package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
)

type (
	Account struct {
		APIKey    string
		APISecret string
	}

	RuleConfiguration struct {
		ID                    int64
		Interval              time.Duration
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
	maximumVolume, _ := strconv.ParseFloat(c.FormValue("maximumVolume"), 64)
	transactionVolume, _ := strconv.ParseFloat(c.FormValue("transactionVolume"), 64)
	variance, _ := strconv.ParseFloat(c.FormValue("variance"), 64)
	stepDownPrice, _ := strconv.ParseFloat(c.FormValue("stepDownPrice"), 64)
	minimumBid, _ := strconv.ParseFloat(c.FormValue("minimumBid"), 64)

	intervalRune := []rune(c.FormValue("interval"))
	log.Print(c.FormValue("interval"))

	in, _ := strconv.ParseInt(string(intervalRune[0:len(intervalRune)-1]), 10, 64)
	interval := time.Duration(in)
	switch intervalRune[len(intervalRune)-1] {
	case 'd':
		interval *= (time.Hour * 24)
	case 'h':
		interval *= time.Hour
	case 'm':
		interval *= time.Minute
	case 's':
		interval *= time.Second
	default:
		interval *= time.Second
	}
	log.Print(interval)

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

func botExported(c echo.Context) error {
	q := influx.Query{
		Command:  "select * from transactions",
		Database: "trader",
	}

	resp, err := bot.store.Query(q)
	if err != nil {
		log.Print("error ", err)
	}

	var b bytes.Buffer

	writer := csv.NewWriter(&b)

	for i := 0; i < len(resp.Results[0].Series[0].Values); i++ {
		err = writer.Write([]string{
			resp.Results[0].Series[0].Values[i][0].(string),
			resp.Results[0].Series[0].Values[i][1].(string),
			string(resp.Results[0].Series[0].Values[i][2].(string)),
			string(resp.Results[0].Series[0].Values[i][3].(json.Number)),
			string(resp.Results[0].Series[0].Values[i][4].(json.Number)),
		})
		if err != nil {
			log.Print("error ", err)
		}
	}
	writer.Flush()

	return c.Blob(http.StatusOK, "text/csv", b.Bytes())
}

func updateSimulate(c echo.Context) error {
	power, err := strconv.ParseInt(c.FormValue("power"), 10, 64)
	if err != nil {
		jsonBadRequest(c, err)
	}

	if power == 1 {
		bot.simulate = true
		insertBotSimulateStatus("ON")
	} else {
		bot.simulate = false
		insertBotSimulateStatus("OFF")
	}

	return jsonSuccess(c, echo.Map{})
}
