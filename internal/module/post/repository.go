package post

import (
	"bloggo/internal/module/post/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/schemas/responses"
	"bloggo/internal/utils/slugify"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

type PostRepository struct {
	database *sql.DB
}

// formatCoverImagePath converts database cover image filename to API path format
// Input: "abc123.webp" -> Output: "/uploads/cover/abc123"
func formatCoverImagePath(filename *string) *string {
	if filename == nil || *filename == "" {
		return nil
	}

	// Remove extension and format as API path
	nameWithoutExt := strings.TrimSuffix(*filename, filepath.Ext(*filename))
	formatted := "/uploads/cover/" + nameWithoutExt
	return &formatted
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
		var rawCoverImage *string

		err := rows.Scan(
			&post.PostId,
			&post.Author.Id,
			&post.Author.Name,
			&post.Author.Avatar,
			&post.Title,
			&post.Slug,
			&rawCoverImage,
			&post.Spot,
			&post.Status,
			&post.ReadCount,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Category.Slug,
			&post.Category.Id,
			&post.Category.Name,
		)
		if err != nil {
			return nil, err
		}

		// Format cover image path
		post.CoverImage = formatCoverImagePath(rawCoverImage)

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository *PostRepository) GetPostListPaginated(
	filters *models.RequestPostFilters,
) (*responses.PaginatedResponse[models.ResponsePostCard], error) {
	var whereClauses []string
	var args []any

	// Add search filter
	if filters.Q != nil && strings.TrimSpace(*filters.Q) != "" {
		whereClauses = append(whereClauses, "(COALESCE(current_pv.title, best_pv.title) LIKE ? OR COALESCE(current_pv.spot, best_pv.spot) LIKE ? OR COALESCE(current_pv.content, best_pv.content) LIKE ?)")
		searchTerm := "%" + *filters.Q + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	// Add status filter
	if filters.Status != nil {
		whereClauses = append(whereClauses, "COALESCE(current_pv.status, best_pv.status) = ?")
		args = append(args, *filters.Status)
	}

	// Add category filter
	if filters.CategoryId != nil {
		whereClauses = append(whereClauses, "COALESCE(current_pv.category_id, best_pv.category_id) = ?")
		args = append(args, *filters.CategoryId)
	}

	// Add author filter
	if filters.AuthorId != nil {
		whereClauses = append(whereClauses, "p.created_by = ?")
		args = append(args, *filters.AuthorId)
	}

	// Build WHERE clause
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " AND " + strings.Join(whereClauses, " AND ")
	}

	// Build ORDER BY clause
	orderClause := ""
	if filters.Order != nil {
		direction := "ASC"
		if filters.Dir != nil && strings.ToUpper(*filters.Dir) == "DESC" {
			direction = "DESC"
		}

		// Map order fields to their proper table prefixes
		var orderField string
		switch *filters.Order {
		case "title":
			orderField = "title"
		case "created_at":
			orderField = "created_at"
		case "updated_at":
			orderField = "updated_at"
		case "read_count":
			orderField = "p.read_count"
		default:
			orderField = "updated_at" // default fallback
		}

		orderClause = fmt.Sprintf(" ORDER BY %s %s", orderField, direction)
	} else {
		orderClause = " ORDER BY updated_at DESC"
	}

	// Set pagination defaults
	page := 1
	take := 12
	if filters.Page != nil {
		page = *filters.Page
	}
	if filters.Take != nil {
		take = *filters.Take
	}

	// Calculate offset
	offset := (page - 1) * take

	// Build final query using the existing QueryPostGetList
	finalQuery := fmt.Sprintf(QueryPostGetList, whereClause+orderClause+" LIMIT ? OFFSET ?")

	// Build count query using the declared query
	countQuery := fmt.Sprintf(QueryPostGetListCount, whereClause)

	// Add pagination args
	queryArgs := append(args, take, offset)

	// Get total count
	var total int64
	err := repository.database.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get data
	rows, err := repository.database.Query(finalQuery, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.ResponsePostCard{}
	for rows.Next() {
		var post models.ResponsePostCard
		var rawCoverImage *string

		err := rows.Scan(
			&post.PostId,
			&post.Author.Id,
			&post.Author.Name,
			&post.Author.Avatar,
			&post.Title,
			&post.Slug,
			&rawCoverImage,
			&post.Spot,
			&post.Status,
			&post.ReadCount,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Category.Slug,
			&post.Category.Id,
			&post.Category.Name,
		)
		if err != nil {
			return nil, err
		}

		// Format cover image path
		post.CoverImage = formatCoverImagePath(rawCoverImage)

		posts = append(posts, post)
	}

	return &responses.PaginatedResponse[models.ResponsePostCard]{
		Data:  posts,
		Page:  page,
		Take:  take,
		Total: total,
	}, nil
}

func (repository *PostRepository) GetPostById(
	id int64,
) (*models.ResponsePostDetails, error) {
	row := repository.database.QueryRow(QueryPostGetById, id)

	var post models.ResponsePostDetails
	var rawCoverImage *string
	err := row.Scan(
		&post.PostId,
		&post.VersionId,
		&post.Author.Id,
		&post.Author.Name,
		&post.Author.Avatar,
		&post.Title,
		&post.Slug,
		&post.Content,
		&rawCoverImage,
		&post.Description,
		&post.Spot,
		&post.Status,
		&post.ReadCount,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Category.Slug,
		&post.Category.Id,
		&post.Category.Name,
	)
	if err != nil {
		return nil, err
	}

	// Format cover image path
	post.CoverImage = formatCoverImagePath(rawCoverImage)

	return &post, nil
}

func (repository *PostRepository) GetPostGetByCurrentVersionSlug(
	slug string,
) (*models.ResponsePostDetails, error) {
	row := repository.database.QueryRow(QueryPostGetByCurrentVersionSlug, slug, slug)

	var post models.ResponsePostDetails
	var rawCoverImage *string
	err := row.Scan(
		&post.PostId,
		&post.VersionId,
		&post.Author.Id,
		&post.Author.Name,
		&post.Author.Avatar,
		&post.Title,
		&post.Slug,
		&post.Content,
		&rawCoverImage,
		&post.Description,
		&post.Spot,
		&post.Status,
		&post.ReadCount,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Category.Slug,
		&post.Category.Id,
		&post.Category.Name,
	)
	if err != nil {
		return nil, err
	}

	// Format cover image path
	post.CoverImage = formatCoverImagePath(rawCoverImage)

	return &post, nil
}

func (repository *PostRepository) CreatePost(
	model *models.RequestPostUpsert,
	coverPath string,
	readTime int,
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
		readTime,
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
		var rawCoverImage *string
		if err := rows.Scan(
			&version.VersionId,
			&version.VersionAuthor.Id,
			&version.VersionAuthor.Name,
			&version.VersionAuthor.Avatar,
			&version.Title,
			&rawCoverImage,
			&version.Status,
			&version.UpdatedAt,
			&version.Category.Id,
			&version.Category.Name,
			&version.Category.Slug,
		); err != nil {
			return nil, err
		}

		// Process cover image
		version.CoverImage = formatCoverImagePath(rawCoverImage)

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
	var rawCoverImage *string
	var statusChangedById *int64
	var statusChangedByName *string
	var statusChangedByAvatar *string
	if err := row.Scan(
		&result.VersionId,
		&result.DuplicatedFrom,
		&result.VersionAuthor.Id,
		&result.VersionAuthor.Name,
		&result.VersionAuthor.Avatar,
		&result.Title,
		&result.Slug,
		&result.Content,
		&rawCoverImage,
		&result.Description,
		&result.Spot,
		&result.Status,
		&result.StatusChangedAt,
		&result.StatusChangeNote,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Category.Id,
		&result.Category.Name,
		&result.Category.Slug,
		&statusChangedById,
		&statusChangedByName,
		&statusChangedByAvatar,
	); err != nil {
		return nil, err
	}

	// Build StatusChangedBy user object if data exists
	if statusChangedById != nil && statusChangedByName != nil {
		result.StatusChangedBy = &struct {
			Id     int64   `json:"id"`
			Name   string  `json:"name"`
			Avatar *string `json:"avatar"`
		}{
			Id:     *statusChangedById,
			Name:   *statusChangedByName,
			Avatar: statusChangedByAvatar,
		}
	}

	// Format cover image path
	result.CoverImage = formatCoverImagePath(rawCoverImage)

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
		&duplicate.ReadTime,
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
		&duplicate.ReadTime,
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

func (repository *PostRepository) CreateVersionFromSpecificVersion(
	versionId int64,
	authorId int64,
) (int64, error) {
	transaction, err := repository.database.Begin()
	if err != nil {
		return 0, err
	}

	copyingRow := transaction.QueryRow(QueryGetSpecificVersionForDuplicate, versionId)
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
		&duplicate.ReadTime,
		&duplicate.CreatedBy,
	); err != nil {
		transaction.Rollback()
		return 0, err
	}

	result, err := transaction.Exec(
		QueryPostVersionCreate,
		duplicate.PostId,
		&duplicate.Title,
		&duplicate.Slug,
		&duplicate.Content,
		&duplicate.CoverImage,
		&duplicate.Description,
		&duplicate.Spot,
		&duplicate.CategoryId,
		&duplicate.ReadTime,
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
	readTime *int,
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
		readTime,
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

func (repository *PostRepository) IncrementReadCount(postId int64) error {
	_, err := repository.database.Exec(QueryIncrementReadCount, postId)
	return err
}

func (repository *PostRepository) SetCurrentVersionForPost(postId int64, versionId int64) error {
	_, err := repository.database.Exec(QueryPostSetCurrentVersion, versionId, postId)
	return err
}

func (repository *PostRepository) GetPublishedVersionBySlug(slug string) (*struct {
	Id     int64
	PostId int64
}, error) {
	row := repository.database.QueryRow(QueryGetPublishedVersionBySlug, slug)

	var result struct {
		Id     int64
		PostId int64
	}
	err := row.Scan(&result.Id, &result.PostId)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repository *PostRepository) UnpublishVersionBySlug(slug string) error {
	_, err := repository.database.Exec(QueryUnpublishVersionBySlug, slug)
	return err
}

func (repository *PostRepository) GetVersionSlug(versionId int64) (string, error) {
	row := repository.database.QueryRow(QueryGetVersionSlug, versionId)

	var slug string
	err := row.Scan(&slug)
	if err != nil {
		return "", err
	}

	return slug, nil
}

func (repository *PostRepository) TrackView(postId int64, userAgent string) error {
	transaction, err := repository.database.Begin()
	if err != nil {
		return err
	}

	// Insert view record
	_, err = transaction.Exec(QueryInsertPostView, postId, userAgent)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Increment read count
	_, err = transaction.Exec(QueryIncrementReadCount, postId)
	if err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit()
}

func (repository *PostRepository) AssignTagsToPost(postId int64, tagIds []int64) error {
	// Check if all tags exist
	for _, tagId := range tagIds {
		if exists, err := repository.checkTagExists(tagId); err != nil {
			return err
		} else if !exists {
			return apierrors.ErrNotFound
		}
	}

	// Get current tags for the post
	currentTagIds, err := repository.getCurrentPostTagIds(postId)
	if err != nil {
		return err
	}

	// Convert slices to maps for easier comparison
	currentTagMap := make(map[int64]bool)
	for _, tagId := range currentTagIds {
		currentTagMap[tagId] = true
	}

	newTagMap := make(map[int64]bool)
	for _, tagId := range tagIds {
		newTagMap[tagId] = true
	}

	// Find tags to remove (in current but not in new)
	var tagsToRemove []int64
	for tagId := range currentTagMap {
		if !newTagMap[tagId] {
			tagsToRemove = append(tagsToRemove, tagId)
		}
	}

	// Find tags to add (in new but not in current)
	var tagsToAdd []int64
	for tagId := range newTagMap {
		if !currentTagMap[tagId] {
			tagsToAdd = append(tagsToAdd, tagId)
		}
	}

	// If no changes needed, return early
	if len(tagsToRemove) == 0 && len(tagsToAdd) == 0 {
		return nil
	}

	// Start transaction
	transaction, err := repository.database.Begin()
	if err != nil {
		return err
	}

	// Remove tags that are no longer needed (batch operation)
	if len(tagsToRemove) > 0 {
		err := repository.removeTagsFromPostBatch(transaction, postId, tagsToRemove)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	// Add new tags (batch operation)
	if len(tagsToAdd) > 0 {
		err := repository.addTagsToPostBatch(transaction, postId, tagsToAdd)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	return transaction.Commit()
}

func (repository *PostRepository) getCurrentPostTagIds(postId int64) ([]int64, error) {
	rows, err := repository.database.Query(QueryGetCurrentPostTagIds, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tagIds []int64
	for rows.Next() {
		var tagId int64
		err := rows.Scan(&tagId)
		if err != nil {
			return nil, err
		}
		tagIds = append(tagIds, tagId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tagIds, nil
}

func (repository *PostRepository) checkTagExists(tagId int64) (bool, error) {
	row := repository.database.QueryRow(QueryCheckTagExists, tagId)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *PostRepository) removeTagsFromPostBatch(transaction *sql.Tx, postId int64, tagIds []int64) error {
	if len(tagIds) == 0 {
		return nil
	}
	
	// Build placeholders for IN clause
	placeholders := make([]string, len(tagIds))
	args := make([]any, len(tagIds)+1)
	args[0] = postId
	
	for i, tagId := range tagIds {
		placeholders[i] = "?"
		args[i+1] = tagId
	}
	
	query := fmt.Sprintf(QueryRemoveTagsFromPost, strings.Join(placeholders, ","))
	_, err := transaction.Exec(query, args...)
	return err
}

func (repository *PostRepository) addTagsToPostBatch(transaction *sql.Tx, postId int64, tagIds []int64) error {
	if len(tagIds) == 0 {
		return nil
	}
	
	// Build VALUES clause
	valueParts := make([]string, len(tagIds))
	args := make([]any, len(tagIds)*2)
	
	for i, tagId := range tagIds {
		valueParts[i] = "(?, ?)"
		args[i*2] = postId
		args[i*2+1] = tagId
	}
	
	query := fmt.Sprintf(QueryAssignTagsToPost, strings.Join(valueParts, ","))
	_, err := transaction.Exec(query, args...)
	return err
}

func (repository *PostRepository) GetPostTags(postId int64) ([]models.TagCard, error) {
	rows, err := repository.database.Query(QueryGetPostTags, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.TagCard{}
	for rows.Next() {
		var tag models.TagCard
		err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.Slug,
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
