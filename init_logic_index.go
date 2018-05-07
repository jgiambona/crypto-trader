package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func index(c echo.Context) error {
	msg := `
 _______  ______    _______  ______   _______  ______
|       ||    _ |  |   _   ||      | |       ||    _ |
|_     _||   | ||  |  |_|  ||  _    ||    ___||   | ||
  |   |  |   |_||_ |       || | |   ||   |___ |   |_||_
  |   |  |    __  ||       || |_|   ||    ___||    __  |
  |   |  |   |  | ||   _   ||       ||   |___ |   |  | |
  |___|  |___|  |_||__| |__||______| |_______||___|  |_|
	`
	return c.String(http.StatusOK, msg)
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

