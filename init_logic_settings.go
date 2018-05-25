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
		MinInterval           time.Duration
		MaxInterval           time.Duration
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

	minIntervalRune := []rune(c.FormValue("minInterval"))
	log.Print(c.FormValue("minInterval"))

	maxIntervalRune := []rune(c.FormValue("maxInterval"))
	log.Print(c.FormValue("maxInterval"))

	if id == 1 {
		bot.ruleOne.Enabled = true
		bot.ruleOne.MinInterval = convertInterval(minIntervalRune)
		bot.ruleOne.MaxInterval = convertInterval(maxIntervalRune)
		bot.ruleOne.MaximumVolume = maximumVolume
		bot.ruleOne.TransactionVolume = transactionVolume
		bot.ruleOne.VarianceOfTransaction = variance
		bot.ruleOne.BidPriceStepDown = stepDownPrice
		bot.ruleOne.MinimumBid = minimumBid
	} else {
		return jsonBadRequest(c, "error no such account.")
	}

	log.Printf(`
	----- One
	--- %.8f
	--- %.8f
	--- %.8f
	--- %.8f
	--- %.8f
	----- End`,
		bot.ruleOne.MaximumVolume, bot.ruleOne.TransactionVolume,
		bot.ruleOne.VarianceOfTransaction, bot.ruleOne.BidPriceStepDown,
		bot.ruleOne.MinimumBid)

	return jsonSuccess(c, echo.Map{
		"account": id,
	})
}

func convertInterval(intervalRune []rune) time.Duration {
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
	return interval
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
