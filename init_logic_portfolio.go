package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func listPortfolio(c echo.Context) error {
	return nil
}

func addNewPortfolio(c echo.Context) error {
	key := c.FormValue("key")
	secret := c.FormValue("secret")
	id, err := insertNewAccount(key, secret)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"id":      id,
	})
}
