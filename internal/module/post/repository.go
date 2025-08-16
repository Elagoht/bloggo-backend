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
			&post.Author.Id,
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
		&post.Author.Id,
		&post.Author.Name,
		&post.Author.Avatar,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.CoverImage,
		&post.Description,
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
		&post.Author.Id,
		&post.Author.Name,
		&post.Author.Avatar,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.CoverImage,
		&post.Description,
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
		nil, // Initial (original) version has not duplicated from any
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

func (repository *PostRepository) GetPostVersionById(
	postId int64,
	versionId int64,
) (*models.ResponseVersionDetailsOfPost, error) {
	row := repository.database.QueryRow(QueryPostVersionGetById, postId, versionId)

	result := models.ResponseVersionDetailsOfPost{}
	if err := row.Scan(
		&result.VersionId,
		&result.DuplicatedFrom,
		&result.VersionAuthor.Id,
		&result.VersionAuthor.Name,
		&result.VersionAuthor.Avatar,
		&result.Title,
		&result.Slug,
		&result.Content,
		&result.CoverImage,
		&result.Description,
		&result.Spot,
		&result.Status,
		&result.StatusChangedAt,
		&result.StatusChangedBy,
		&result.StatusChangeNote,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Category.Id,
		&result.Category.Name,
		&result.Category.Slug,
	); err != nil {
		return nil, err
	}

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

func (repository *PostRepository) CreateVersionFromLatest(
	id int64,
	authorId int64,
) (int64, error) {
	transaction, err := repository.database.Begin()
	if err != nil {
		return 0, err
	}

	copyingRow := transaction.QueryRow(QueryGetPostVersionDuplicate, id)
	duplicate := models.QueryGetPostVersionDuplicateData{}
	if err := copyingRow.Scan(
		&duplicate.VersionId,
		&duplicate.PostId,
		&duplicate.Title,
		&duplicate.Slug,
		&duplicate.Content,
		&duplicate.CoverImage,
		&duplicate.Description,
		&duplicate.Spot,
		&duplicate.CategoryId,
		&duplicate.CreatedBy,
	); err != nil {
		transaction.Rollback()
		return 0, err
	}

	result, err := transaction.Exec(
		QueryPostVersionCreate,
		id,
		&duplicate.Title,
		&duplicate.Slug,
		&duplicate.Content,
		&duplicate.CoverImage,
		&duplicate.Description,
		&duplicate.Spot,
		&duplicate.CategoryId,
		authorId,
		&duplicate.VersionId,
	)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	createdId, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	transaction.Commit()

	return createdId, nil
}

func (repository *PostRepository) GetVersionCreatorAndStatus(
	id int64,
) (int64, int64, error) {
	row := repository.database.QueryRow(QueryGetVersionCreatorAndStatus, id)

	var creatorId, status int64
	if err := row.Scan(
		&creatorId,
		&status,
	); err != nil {
		return 0, 0, err
	}

	if creatorId == 0 {
		return 0, 0, apierrors.ErrNotFound
	}

	return creatorId, status, nil
}

func (repository *PostRepository) UpdateVersionById(
	postId int64,
	versionId int64,
	userId int64,
	model *models.RequestPostUpsert,
	filePath *string,
) error {
	_, err := repository.database.Exec(
		QueryPostVersionUpdate,
		model.Title,
		slugify.Slugify(model.Title),
		model.Content,
		filePath,
		model.Description,
		model.Spot,
		model.CategoryId,
		versionId,
	)
	return err
}

func (repository *PostRepository) UpdateVersionStatus(
	versionId int64,
	status int64,
	statusChangedBy int64,
) error {
	_, err := repository.database.Exec(
		QueryPostVersionUpdateStatus,
		status,
		statusChangedBy,
		versionId,
	)
	return err
}

func (repository *PostRepository) UpdateVersionStatusWithNote(
	versionId int64,
	status int64,
	statusChangedBy int64,
	note *string,
) error {
	_, err := repository.database.Exec(
		QueryPostVersionUpdateStatusWithNote,
		status,
		statusChangedBy,
		note,
		versionId,
	)
	return err
}

func (repository *PostRepository) GetVersionCoverImage(
	versionId int64,
) (*string, error) {
	row := repository.database.QueryRow(QueryGetVersionCoverImage, versionId)

	var coverImage sql.NullString
	if err := row.Scan(&coverImage); err != nil {
		if err == sql.ErrNoRows {
			return nil, apierrors.ErrNotFound
		}
		return nil, err
	}

	if coverImage.Valid {
		return &coverImage.String, nil
	}
	return nil, nil
}

func (repository *PostRepository) SoftDeleteVersionById(
	versionId int64,
) error {
	result, err := repository.database.Exec(QuerySoftDeleteVersion, versionId)
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

	return nil
}

func (repository *PostRepository) IsImageReferencedByOtherVersions(
	imagePath string,
	excludeVersionId int64,
) (bool, error) {
	row := repository.database.QueryRow(
		QueryCheckImageReferences,
		imagePath,
		excludeVersionId,
	)

	var count int64
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *PostRepository) IsVersionCurrentlyPublished(versionId int64) (bool, error) {
	row := repository.database.QueryRow(QueryCheckIfVersionIsCurrentlyPublished, versionId)
	
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (repository *PostRepository) SetPostCurrentVersionToNull(versionId int64) error {
	_, err := repository.database.Exec(QuerySetPostCurrentVersionToNull, versionId)
	return err
}
