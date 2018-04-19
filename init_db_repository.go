package main

import (
	"database/sql"
	"log"
)

func createLocalDB() error {
	query := "create table if not exists accounts (id integer not null primary key, apiKey text, apiSecret text);"

	_, err := bot.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func insertNewAccount(apiKey, apiSecret string) (int64, error) {
	tx, err := bot.db.Begin()
	if err != nil {
		return -1, err
	}

	stmt, err := tx.Prepare("insert into accounts(id, apiKey, apiSecret) values(?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var res sql.Result
	res, err = stmt.Exec(bot.nextID, apiKey, apiSecret)
	if err != nil {
		return -1, err
	}
	tx.Commit()

	bot.nextID += 1
	log.Println("insert new account")
	return res.LastInsertId()
}

func getLastAccountID() (int64, error) {
	var id sql.NullInt64

	row := bot.db.QueryRow("select max(id) from accounts")
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return 1, nil
	case nil:
		if id.Int64 < 1 {
			return 1, nil
		}
		return id.Int64, nil
	default:
		return -1, err
	}
}
