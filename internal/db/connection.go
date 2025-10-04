package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db      *sql.DB
	once    sync.Once
	initErr error
)

// MustConnect initializes the database connection and returns an error if it fails
func MustConnect() error {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", "bloggo.sqlite")
		if err != nil {
			initErr = fmt.Errorf("cannot open the database: %w", err)
			return
		}

		if err = db.Ping(); err != nil {
			initErr = fmt.Errorf("cannot connect to database: %w", err)
			return
		}

		// If the database recently created, create the tables
		if err = InitializeTables(db); err != nil {
			initErr = fmt.Errorf("cannot initialize tables: %w", err)
			return
		}

		if err = SeedDatabase(db); err != nil {
			initErr = fmt.Errorf("cannot seed database: %w", err)
			return
		}
	})

	return initErr
}

// Get returns the singleton database instance
func Get() *sql.DB {
	return db
}
