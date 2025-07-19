package db

import (
	"database/sql"
	"fmt"
)

func InitializeTables(database *sql.DB) {
	for _, query := range InitializeQueries {
		_, err := database.Exec(query)
		if err != nil {
			panic(fmt.Errorf("Database cannot be initialized: %w", err))
		}
	}
}
