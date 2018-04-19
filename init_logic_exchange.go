package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

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
