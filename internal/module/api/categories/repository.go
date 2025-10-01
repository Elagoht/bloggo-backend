package categories

import (
	"bloggo/internal/module/api/categories/models"
	"bloggo/internal/utils/apierrors"
	"database/sql"
)

type CategoriesAPIRepository struct {
	database *sql.DB
}

func NewCategoriesAPIRepository(database *sql.DB) CategoriesAPIRepository {
	return CategoriesAPIRepository{database}
}

func (r *CategoriesAPIRepository) GetAllCategories() (*models.APICategoriesResponse, error) {
	rows, err := r.database.Query(QueryAPIGetAllCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.APICategoryDetails{}
	for rows.Next() {
		var cat models.APICategoryDetails
		err := rows.Scan(
			&cat.Slug,
			&cat.Name,
			&cat.Description,
			&cat.Spot,
			&cat.PostCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return &models.APICategoriesResponse{
		Categories: categories,
	}, nil
}

func (r *CategoriesAPIRepository) GetCategoryBySlug(slug string) (*models.APICategoryDetails, error) {
	row := r.database.QueryRow(QueryAPIGetCategoryBySlug, slug)

	var cat models.APICategoryDetails
	err := row.Scan(
		&cat.Slug,
		&cat.Name,
		&cat.Description,
		&cat.Spot,
		&cat.PostCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierrors.ErrNotFound
		}
		return nil, err
	}

	return &cat, nil
}
