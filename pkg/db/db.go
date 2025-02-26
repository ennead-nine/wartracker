package db

import (
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Connection        *sqlx.DB
	ErrDbErrorUnknown = errors.New("an unknown database error occured")
)

func Connect(f string) (*sqlx.DB, error) {
	var db *sqlx.DB

	db, err := sqlx.Open("sqlite3", f)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
