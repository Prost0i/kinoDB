package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func ConnectDB() error {
	connectString := `host=localhost port=5432 user=postgres dbname=kinodb sslmode=disable`
	var err error
	db, err = sqlx.Open("postgres", connectString)
	if err != nil {
		return err
	}

	return nil
}