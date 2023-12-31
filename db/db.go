package db

import (
	"database/sql"
	"os"

	"github.com/morikuni/failure"
	"github.com/taro-28/saas-sample-api/fail"
)

func Connect() (*sql.DB, func(), error) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		return nil, nil, failure.Translate(err, fail.InternalServerError, failure.Message("failed to open db"))
	}

	return db, func() {
		db.Close()
	}, nil
}
