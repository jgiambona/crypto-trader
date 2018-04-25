package main

import (
	"database/sql"
	"fmt"
	"log"
)

type (
	Account struct {
		ID          int64
		Key, Secret string
	}
)

func repoCreateDB() error {
	query := "CREATE TABLE IF NOT EXISTS accounts (id SERIAL PRIMARY KEY, apiKey TEXT, apiSecret TEXT)"

	_, err := bot.db.Exec(query)
	if err != nil {
		return err
	}

	query = "CREATE TABLE IF NOT EXISTS logs (id SERIAL PRIMARY KEY, message TEXT)"

	_, err = bot.db.Exec(query)
	if err != nil {
		return err
	}

	err = repoInsertNewLog("Sync database.")
	if err != nil {
		return err
	}

	return nil
}

func repoRowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT EXISTS (%s)", query)
	err := bot.db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	return exists
}

func repoListAccounts() ([]Account, error) {
	var list []Account

	rows, err := bot.db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Account
		err = rows.Scan(&r.ID, &r.Key, &r.Secret)
		if err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}

func repoInsertNewLog(message string) error {
	_, err := bot.db.Query(`INSERT INTO logs(message) VALUES($1)`, message)
	if err != nil {
		return err
	}

	log.Println("inserted new log")
	return nil
}

func repoInsertNewAccount(apiKey, apiSecret string) (int64, error) {
	var  accountID int64
	err := bot.db.QueryRow(
		`INSERT INTO accounts(apiKey, apiSecret) VALUES($1, $2) RETURNING id`,
		apiKey, apiSecret).Scan(&accountID)
	if err != nil {
		return -1, err
	}

	bot.nextID += 1
	return accountID, nil
}

func repoGetLastAccountID() (int64, error) {
	var id sql.NullInt64

	err := bot.db.QueryRow("select max(id) from accounts").Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	if id.Int64 < 1 {
		return 0, nil
	}
	return id.Int64, nil
}

func repoAccountExists(apiKey string) bool {
	return repoRowExists("select id from accounts where apiKey = $1", apiKey)
}
