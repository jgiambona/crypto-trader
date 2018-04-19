package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func botRestart(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func botSuspend(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
