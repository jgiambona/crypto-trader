package main

import "github.com/labstack/echo"

func loadRoutes(e *echo.Echo) {
	e.Add("GET", "/", index)
	e.Add("GET", "/bot/restart", botRestart)
	e.Add("GET", "/bot/suspend", botSuspend)
	e.Add("GET", "/bot/exchange/:name/info", exchangeGetConfigInfo)
	e.Add("GET", "/bot/portfolios", portfolioList)

	e.Add("POST", "/bot/setup", setup)
	e.Add("POST", "/bot/setup/pp", setupPingpong)
	e.Add("POST", "/bot/portfolio", portfolioAddNewAccount)
}
