package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ffimnsr/trader/exchange"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

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

// SessionConfig stores the current session data in memory.
type SessionConfig struct {
	Exchange         string
	Strategy         string
	BTCWalletAddress string
	APIKey           string
	APISecret        string
	PPBuyPrice       string
	PPSellPrice      string
	BotStatus        string
}

// Bot is the singleton that holds all the data.
type Bot struct {
	config    *BotConfig
	exchanges []exchange.BotExchange
	store     influx.Client
	session   *SessionConfig
}

var bot Bot

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if len(os.Getenv("T_PROD")) > 0 {
		bot.store, _ = influx.NewHTTPClient(influx.HTTPConfig{
			Addr: "http://ec2-54-169-102-171.ap-southeast-1.compute.amazonaws.com:8086",
		})
	} else {
		bot.store, _ = influx.NewHTTPClient(influx.HTTPConfig{
			Addr: "http://localhost:8086",
		})
	}

	bot.session = &SessionConfig{
		Exchange:         "livecoin",
		Strategy:         "pp",
		BTCWalletAddress: "",
		APIKey:           "",
		APISecret:        "",
		PPBuyPrice:       "",
		PPSellPrice:      "",
	}

	loadConfig()
	loadExchanges()
	if len(bot.exchanges) == 0 {
		e.Logger.Fatal("no exchanges were loaded.")
	}

	go socketCheckBalance()
	go pollTicker()

	e.Add("GET", "/", index)
	e.Add("GET", "/bot/:name/info", getExchangeConfigInfo)
	e.Add("GET", "/bot/restart", restart)
	e.Add("GET", "/bot/suspend", suspend)
	e.Add("GET", "/ws", socket)

	e.Add("POST", "/bot/setup", setup)
	e.Add("POST", "/bot/setup/pp", setupPingpong)
	e.Add("POST", "/bot/setup/portfolio", setupPortfolio)

	if port, ok := os.LookupEnv("PORT"); ok {
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		e.Logger.Fatal(e.Start(":8000"))
	}
}

func index(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func setup(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}

func setupPingpong(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}

func setupPortfolio(c echo.Context) error {
	log.Printf("%+v", c.Request())
	log.Print(c.FormValue("btc"))
	log.Print(c.FormValue("key"))
	log.Print(c.FormValue("secret"))
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}

func restart(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func suspend(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func getExchangeConfigInfo(c echo.Context) error {
	flag := -1
	for i, x := range bot.config.Exchanges {
		param := c.Param("name")
		log.Print(param)
		if len(param) > 0 && param == strings.ToLower(x.Name) {
			flag = i
		}
	}

	if flag < 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "exchange is not been set",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"config":  bot.config.Exchanges[flag],
	})
}

func socket(c echo.Context) error {
	//websocket.Handler(func(ws *websocket.Conn) {
	//	defer ws.Close()
	//	for {
	//		err := websocket.Message.Send(ws, "Hello, Client!")
	//		if err != nil {
	//			c.Logger().Error(err)
	//		}

	//		msg := ""
	//		err = websocket.Message.Receive(ws, &msg)
	//		if err != nil {
	//			c.Logger().Error(err)
	//		}

	//		fmt.Printf("%s\n", msg)
	//	}
	//}).ServeHTTP(c.Response(), c.Request())
	return nil
}
