package main

import (
	"log"

	"github.com/labstack/echo"
)

func portfolioList(c echo.Context) error {
	list, err := repoListAccounts()
	if err != nil {
		return jsonServerError(c, err)
	}

	return jsonSuccess(c, echo.Map{
		"accounts": list,
	})
}

func portfolioAddNewAccount(c echo.Context) error {
	key := c.FormValue("key")
	secret := c.FormValue("secret")

	if exists := repoAccountExists(key); exists {
		return jsonBadRequest(c, "account already exists")
	}

	log.Print(key)
	id, err := repoInsertNewAccount(key, secret)
	if err != nil {
		return jsonServerError(c, err)
	}
	log.Print("Hello")

	return jsonSuccess(c, echo.Map{
		"id": id,
	})
}
