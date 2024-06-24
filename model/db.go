package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func ConnectDB(connectString string) error {
	var err error
	db, err = sqlx.Open("postgres", connectString)
	if err != nil {
		return err
	}

	return nil
}
