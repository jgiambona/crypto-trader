package main

import (
	"fmt"

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
		repoInsertNewLog(fmt.Sprintf("Account %s alreadly exists.", key))
		return jsonBadRequest(c, "account already exists")
	}

	id, err := repoInsertNewAccount(key, secret)
	if err != nil {
		repoInsertNewLog(fmt.Sprintf("Can't create account %s.", key))
		return jsonServerError(c, err)
	}

	repoInsertNewLog(fmt.Sprintf("Account %s created.", key))
	return jsonSuccess(c, echo.Map{
		"id": id,
	})
}
