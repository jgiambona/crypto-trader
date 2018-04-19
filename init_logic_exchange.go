package main

import (
	"log"
	"strings"

	"github.com/labstack/echo"
)

func exchangeGetConfigInfo(c echo.Context) error {
	flag := -1
	for i, x := range bot.config.Exchanges {
		param := c.Param("name")
		log.Print(param)
		if len(param) > 0 && param == strings.ToLower(x.Name) {
			flag = i
		}
	}

	if flag < 0 {
		return jsonBadRequest(c, "exchange is not been set")
	}

	return jsonSuccess(c, echo.Map{
		"config": bot.config.Exchanges[flag],
	})
}
