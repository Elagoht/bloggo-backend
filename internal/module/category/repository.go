package category

import (
	"bloggo/internal/module/category/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"database/sql"
	"fmt"
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
		&category.BlogCount,
	)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Remove CategoryQueryOptions, use handlers.PaginationOptions instead

func (repository *CategoryRepository) GetCategories(
	paginate *pagination.PaginationOptions,
	search *filter.SearchOptions,
) ([]models.ResponseCategoryCard, error) {
	orderByClause, limitClause, offsetClause, args := paginate.BuildPaginationClauses()
	searchClause, searchArgs := "", []any{}
	if search != nil && search.Q != nil {
		searchClause, searchArgs = search.BuildSearchClause([]string{"name"})
		if searchClause != "" {
			searchClause = " AND " + searchClause
		}
	}
	query := fmt.Sprintf(QueryCategoryGetCategories, searchClause, orderByClause, limitClause, offsetClause)
	allArgs := append(searchArgs, args...)

	rows, err := repository.database.Query(query, allArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.ResponseCategoryCard{}
	for rows.Next() {
		var category models.ResponseCategoryCard
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Slug,
			&category.BlogCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (repository *CategoryRepository) CategoryUpdate(
	slug string,
	model *models.QueryParamsCategoryUpdate,
) error {
	statement, err := repository.database.Prepare(QueryCategoryPatch)
	if err != nil {
		return err
	}

	result, err := statement.Exec(
		model.Name,
		model.Slug,
		model.Description,
		slug,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return apierrors.ErrNotFound
	}

	return nil
}

func (repository *CategoryRepository) CategoryDelete(
	slug string,
) error {
	statement, err := repository.database.Prepare(QueryCategorySoftDelete)
	if err != nil {
		return err
	}

	result, err := statement.Exec(slug)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return apierrors.ErrNotFound
	}

	return nil
}
