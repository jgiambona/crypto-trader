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
