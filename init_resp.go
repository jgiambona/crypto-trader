package main

import (
	"net/http"

	"github.com/labstack/echo"
)

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
