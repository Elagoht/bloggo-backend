package db

import (
	"database/sql"
)

func InitializeTables(database *sql.DB) {
	for _, query := range InitializeQueries {
		_, err := database.Exec(query)
		if err != nil {
			panic("Database cannot be initialized.")
		}
	}
}
