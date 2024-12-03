package db

import (
	"github.com/jmoiron/sqlx"
)

var (
	Connection *sqlx.DB
)

func Connect(f string) (*sqlx.DB, error) {
	var db *sqlx.DB

	db, err := sqlx.Open("sqlite3", f)
	if err != nil {
		return nil, err
	}

	return db, nil
}
