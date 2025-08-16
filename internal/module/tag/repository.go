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

func (repository *TagRepository) GetCategories(
	paginate *pagination.PaginationOptions,
	search *filter.SearchOptions,
) ([]models.ResponseTagCard, error) {
	// Handle pagination and order params
	orderByClause, limitClause, offsetClause, args := paginate.BuildPaginationClauses()

	// Handle search by name
	searchClause, searchArgs := filter.BuildSearchClause(search, []string{"name"})

	// Merge them and generate query
	query, allArgs := handlers.BuildModifiedSQL(
		QueryTagGetCategories,
		[]string{searchClause, orderByClause, limitClause, offsetClause},
		[][]any{searchArgs, args},
	)

	// Run query
	rows, err := repository.database.Query(query, allArgs...)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
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

// Post-Tag Relationship Methods
func (repository *TagRepository) GetPostTags(postId int64) ([]models.ResponseTagCard, error) {
	rows, err := repository.database.Query(QueryGetPostTags, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.ResponseTagCard{}
	for rows.Next() {
		var tag models.ResponseTagCard
		err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.Slug,
			&tag.BlogCount,
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

func (repository *TagRepository) AssignTagsToPost(postId int64, tagIds []int64) error {
	// Check if post exists
	if exists, err := repository.checkPostExists(postId); err != nil {
		return err
	} else if !exists {
		return apierrors.ErrNotFound
	}

	// Check if all tags exist
	for _, tagId := range tagIds {
		if exists, err := repository.checkTagExists(tagId); err != nil {
			return err
		} else if !exists {
			return apierrors.ErrNotFound
		}
	}

	// Assign tags to post
	for _, tagId := range tagIds {
		_, err := repository.database.Exec(QueryAssignTagToPost, postId, tagId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository *TagRepository) RemoveTagFromPost(postId int64, tagId int64) error {
	result, err := repository.database.Exec(QueryRemoveTagFromPost, postId, tagId)
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

func (repository *TagRepository) RemoveAllTagsFromPost(postId int64) error {
	_, err := repository.database.Exec(QueryRemoveAllTagsFromPost, postId)
	return err
}

func (repository *TagRepository) checkPostExists(postId int64) (bool, error) {
	row := repository.database.QueryRow(QueryCheckPostExists, postId)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *TagRepository) checkTagExists(tagId int64) (bool, error) {
	row := repository.database.QueryRow(QueryCheckTagExists, tagId)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
