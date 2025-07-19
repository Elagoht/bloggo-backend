package category

import (
	"database/sql"
)

type CategoryRepository struct {
	database *sql.DB
}

func NewCategoryRepository(database *sql.DB) CategoryRepository {
	return CategoryRepository{
		database,
	}
}
