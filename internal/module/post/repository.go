package post

import (
	"bloggo/internal/module/post/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/slugify"
	"database/sql"
)

type PostRepository struct {
	database *sql.DB
}

func NewPostRepository(database *sql.DB) PostRepository {
	return PostRepository{
		database,
	}
}

func (repository *PostRepository) GetPostList() (
	[]models.ResponsePostCard,
	error,
) {
	rows, err := repository.database.Query(QueryPostGetList)
	if err != nil {
		return nil, err
	}

	posts := []models.ResponsePostCard{}
	for rows.Next() {
		var post models.ResponsePostCard

		err := rows.Scan(
			&post.PostId,
			&post.Author.Name,
			&post.Author.Avatar,
			&post.Title,
			&post.Slug,
			&post.CoverImage,
			&post.Spot,
			&post.Status,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Category.Slug,
			&post.Category.Id,
			&post.Category.Name,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repository *PostRepository) GetPostById(
	id int64,
) (*models.ResponsePostDetails, error) {
	row := repository.database.QueryRow(QueryPostGetById, id)

	var post models.ResponsePostDetails
	err := row.Scan(
		&post.PostId,
		&post.VersionId,
		&post.Author.Name,
		&post.Author.Email,
		&post.Author.Avatar,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.CoverImage,
		&post.Description,
		&post.Spot,
		&post.Status,
		&post.StatusChangedAt,
		&post.StatusChangedBy,
		&post.StatusChangeNote,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Category.Slug,
		&post.Category.Id,
		&post.Category.Name,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (repository *PostRepository) GetPostGetByCurrentVersionSlug(
	slug string,
) (*models.ResponsePostDetails, error) {
	row := repository.database.QueryRow(QueryPostGetByCurrentVersionSlug, slug)

	var post models.ResponsePostDetails
	err := row.Scan(
		&post.PostId,
		&post.VersionId,
		&post.Author.Name,
		&post.Author.Email,
		&post.Author.Avatar,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.CoverImage,
		&post.Description,
		&post.Spot,
		&post.Status,
		&post.StatusChangedAt,
		&post.StatusChangedBy,
		&post.StatusChangeNote,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Category.Slug,
		&post.Category.Id,
		&post.Category.Name,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (repository *PostRepository) CreatePost(
	model *models.RequestPostUpsert,
	coverPath string,
	authorId int64,
) (int64, error) {
	transaction, err := repository.database.Begin()
	if err != nil {
		return 0, err
	}

	// Create Post
	createdPost, err := transaction.Exec(QueryPostCreate, authorId)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	createdPostId, err := createdPost.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	// Create first version automatically
	createdVersion, err := transaction.Exec(
		QueryPostVersionCreate,
		createdPostId,
		model.Title,
		slugify.Slugify(model.Title),
		model.Content,
		coverPath,
		model.Description,
		model.Spot,
		model.CategoryId,
		authorId,
	)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	createdVersionId, err := createdVersion.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	// Assign its first version to created post
	_, err = transaction.Exec(
		QueryPostSetCurrentVersion,
		createdVersionId,
		createdPostId,
	)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	// Run transaction
	err = transaction.Commit()
	if err != nil {
		return 0, err
	}

	return createdPostId, nil
}

func (repository *PostRepository) ListPostVersionsGetByPostId(
	id int64,
) (*models.ResponseVersionsOfPost, error) {
	rows, err := repository.database.Query(QueryPostVersionsGetByPostId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := models.ResponseVersionsOfPost{}

	details := repository.database.QueryRow(
		QueryPostDetailsForVersionsGetByPostId, id,
	)
	if err := details.Scan(
		&result.CurrentVersionId,
		&result.CreatedAt,
		&result.OriginalAuthor.Id,
		&result.OriginalAuthor.Name,
		&result.OriginalAuthor.Avatar,
	); err != nil {
		return nil, err
	}

	versions := []models.PostVersionsCard{}
	for rows.Next() {
		version := models.PostVersionsCard{}
		if err := rows.Scan(
			&version.VersionId,
			&version.VersionAuthor.Id,
			&version.VersionAuthor.Name,
			&version.VersionAuthor.Avatar,
			&version.Title,
			&version.Status,
			&version.UpdatedAt,
		); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	result.Versions = versions

	return &result, nil
}

func (repository *PostRepository) GetAllRelatedCovers(
	id int64,
) ([]string, error) {
	rows, err := repository.database.Query(QueryPostAllRelatedCovers, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	coverPaths := []string{}
	for rows.Next() {
		var coverPath string
		if err := rows.Scan(&coverPath); err != nil {
			return nil, err
		}
		coverPaths = append(coverPaths, coverPath)
	}

	return coverPaths, nil
}

func (repository *PostRepository) SoftDeletePostById(
	id int64,
) error {
	result, err := repository.database.Exec(QueryPostSoftDelete, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected < 1 {
		return apierrors.ErrNotFound
	}

	return err
}
