package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/ffimnsr/trader/exchange/livecoin"
)

// Contains forward declaration of API secrets.
const (
	APIKey       = "gZaXWxSudtwt1AR3cW6Fdh5UY3BgVG4r"
	APISecretKey = "FYw6X3gzcJ4F5JvqmYBqAMwdMexzAay7"
)

func main() {
	e := echo.New()
	e.Renderer = LoadTemplates(e)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", index)
	e.Logger.Fatal(e.Start(":4000"))
}

func index(c echo.Context) error {
	hello := echo.Map{}
	livecoin.GetTicker(hello)
	return c.Render(http.StatusOK, "index.tmpl", nil)
}
