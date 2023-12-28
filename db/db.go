package db

import (
	"database/sql"
	"os"
)

func Connect() (*sql.DB, func(), error) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		db.Close()
	}, nil
}
