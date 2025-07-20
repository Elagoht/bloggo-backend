package db

import (
	"database/sql"
	"log"
)

func InitializeTables(database *sql.DB) {
	for _, query := range InitializeQueries {
		_, err := database.Exec(query)
		if err != nil {
			log.Fatal("Database cannot be initialized.")
		}
	}
}
