package db

import (
	"database/sql"
)

var (
	Con *sql.DB
)

func Connect(f string) (*sql.DB, error) {
	var db *sql.DB

	db, err := sql.Open("sqlite3", f)
	if err != nil {
		return nil, err
	}

	return db, nil
}
