package category

import (
	"bloggo/internal/module/category/models"
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

func (repository *CategoryRepository) CategoryCreate(
	model *models.QueryParamsCategoryCreate,
) (int64, error) {
	statement, err := repository.database.Prepare(QueryCategoryCreate)
	if err != nil {
		return 0, err
	}

	result, err := statement.Exec(
		model.Name,
		model.Slug,
		model.Description,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
