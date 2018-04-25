package main

import (
	"github.com/labstack/echo"
)

func setup(c echo.Context) error {
	return jsonSuccess(c, echo.Map{})
}

func setupPingpong(c echo.Context) error {
	return jsonSuccess(c, echo.Map{})
}


