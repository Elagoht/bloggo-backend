package db

import (
	"database/sql"
	"fmt"
)

func InitializeTables(database *sql.DB) error {
	for _, query := range InitializeQueries {
		_, err := database.Exec(query)
		if err != nil {
			return fmt.Errorf("database cannot be initialized: %w", err)
		}
	}
	return nil
}
