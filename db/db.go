package db

import (
	"database/sql"
	"log"
	"os"
)

func Get() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return db
}
