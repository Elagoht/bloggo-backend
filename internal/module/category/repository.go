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

func (repository *CategoryRepository) GetCategoryBySlug(
	slug string,
) (*models.ResponseCategoryDetails, error) {
	row := repository.database.QueryRow(QueryCategoryGetBySlug, slug)

	var category models.ResponseCategoryDetails
	err := row.Scan(
		&category.Id,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
