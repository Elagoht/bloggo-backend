package tag

import (
	"bloggo/internal/module/tag/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/handlers"
	"bloggo/internal/utils/pagination"
	"database/sql"
)

type TagRepository struct {
	database *sql.DB
}

func NewTagRepository(database *sql.DB) TagRepository {
	return TagRepository{
		database,
	}
}

func (repository *TagRepository) TagCreate(
	model *models.QueryParamsTagCreate,
) (int64, error) {
	statement, err := repository.database.Prepare(QueryTagCreate)
	if err != nil {
		return 0, err
	}

	result, err := statement.Exec(
		model.Name,
		model.Slug,
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

func (repository *TagRepository) GetTagBySlug(
	slug string,
) (*models.ResponseTagDetails, error) {
	row := repository.database.QueryRow(QueryTagGetBySlug, slug)

	var category models.ResponseTagDetails
	err := row.Scan(
		&category.Id,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.BlogCount,
	)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (repository *TagRepository) GetTags(
	paginate *pagination.PaginationOptions,
	search *filter.SearchOptions,
) ([]models.ResponseTagCard, int64, error) {
	// Handle pagination and order params
	orderByClause, limitClause, offsetClause, args := paginate.BuildPaginationClauses()

	// Handle search by name
	searchClause, searchArgs := filter.BuildSearchClause(search, []string{"name"})

	// Get total count
	total, err := repository.getTagsCount(search)
	if err != nil {
		return nil, 0, err
	}

	// Merge them and generate query
	query, allArgs := handlers.BuildModifiedSQL(
		QueryTagGetCategories,
		[]string{searchClause, orderByClause, limitClause, offsetClause},
		[][]any{searchArgs, args},
	)

	// Run query
	rows, err := repository.database.Query(query, allArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	categories := []models.ResponseTagCard{}
	for rows.Next() {
		var category models.ResponseTagCard
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Slug,
			&category.BlogCount,
		)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (repository *TagRepository) TagUpdate(
	slug string,
	model *models.QueryParamsTagUpdate,
) error {
	statement, err := repository.database.Prepare(QueryTagPatch)
	if err != nil {
		return err
	}

	result, err := statement.Exec(
		model.Name,
		model.Slug,

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

func (repository *TagRepository) TagDelete(
	slug string,
) error {
	statement, err := repository.database.Prepare(QueryTagSoftDelete)
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


func (repository *TagRepository) GetTagList() ([]models.ResponseTagListItem, error) {
	rows, err := repository.database.Query(QueryTagList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.ResponseTagListItem{}
	for rows.Next() {
		var tag models.ResponseTagListItem
		err := rows.Scan(
			&tag.Id,
			&tag.Name,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (repository *TagRepository) getTagsCount(search *filter.SearchOptions) (int64, error) {
	searchClause, searchArgs := filter.BuildSearchClause(search, []string{"name"})
	
	countQuery, countArgs := handlers.BuildModifiedSQL(
		QueryTagGetCategoriesCount,
		[]string{searchClause},
		[][]any{searchArgs},
	)
	
	var total int64
	err := repository.database.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return 0, err
	}
	
	return total, nil
}
