package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresDB(databaseURL string) (*sql.DB, error) {
	// Try to connect to the database with retries
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", databaseURL)
		if err == nil {
			// Try pinging the DB to check if it's ready
			if err = db.Ping(); err == nil {
				return db, nil
			}
		}
		time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
	}
	return nil, fmt.Errorf("unable to connect to database after retries: %w", err)
}
