package postgres

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// New creates a new PostgreSQL database connection.
func New(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}