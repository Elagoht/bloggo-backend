package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetInstance() *sql.DB {
	once.Do(func() {
		db, err := sql.Open("sqlite3", "bloggo.sqlite")
		if err != nil {
			log.Fatal("Cannot open the database.")
		}

		if err = db.Ping(); err != nil {
			log.Fatal("Cannot connect to database.")
		}
	})
	return db
}
