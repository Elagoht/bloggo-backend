package removal_request

import (
	"bloggo/internal/module/removal_request/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"bloggo/internal/utils/schemas/responses"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type RemovalRequestRepository struct {
	database *sql.DB
}

func NewRemovalRequestRepository(database *sql.DB) RemovalRequestRepository {
	return RemovalRequestRepository{
		database,
	}
}

func (repository *RemovalRequestRepository) CreateRemovalRequest(
	postVersionId int64,
	requestedBy int64,
	note *string,
) (int64, error) {
	result, err := repository.database.Exec(
		QueryCreateRemovalRequest,
		postVersionId,
		requestedBy,
		note,
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

func (repository *RemovalRequestRepository) GetRemovalRequestList(
	paginate *pagination.PaginationOptions,
	search *filter.SearchOptions,
	status *int,
) (*responses.PaginatedResponse[models.RemovalRequestCard], error) {
	// Build the base query
	baseQuery := `
	SELECT
		rr.id, rr.post_version_id, pv.title as post_title,
		u1.id as requested_by_id, u1.name as requested_by_name, u1.avatar as requested_by_avatar,
		rr.note, rr.status,
		u2.id as decided_by_id, u2.name as decided_by_name, u2.avatar as decided_by_avatar,
		rr.decision_note, rr.decided_at, rr.created_at
	FROM removal_requests rr
	JOIN post_versions pv ON rr.post_version_id = pv.id
	JOIN users u1 ON rr.requested_by = u1.id
	LEFT JOIN users u2 ON rr.decided_by = u2.id
	WHERE 1=1`

	var args []any
	var conditions []string

	// Add status filter
	if status != nil {
		conditions = append(conditions, "rr.status = ?")
		args = append(args, *status)
	}

	// Add search filter
	if search != nil && search.Q != nil {
		searchClause, searchArgs := filter.BuildSearchClause(search, []string{
			"pv.title", "u1.name", "rr.note",
		})
		if searchClause != "" {
			conditions = append(conditions, strings.TrimPrefix(searchClause, "AND "))
			args = append(args, searchArgs...)
		}
	}

	// Add conditions to query
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Add ordering
	orderClause := "ORDER BY rr.created_at DESC"
	if paginate != nil && paginate.OrderBy != nil {
		direction := "ASC"
		if paginate.Direction != nil && *paginate.Direction == "desc" {
			direction = "DESC"
		}

		switch *paginate.OrderBy {
		case "created_at":
			orderClause = "ORDER BY rr.created_at " + direction
		case "decided_at":
			orderClause = "ORDER BY rr.decided_at " + direction
		case "post_title":
			orderClause = "ORDER BY pv.title " + direction
		case "requested_by_name":
			orderClause = "ORDER BY u1.name " + direction
		case "status":
			orderClause = "ORDER BY rr.status " + direction
		}
	}

	baseQuery += " " + orderClause

	// Add pagination
	var limit, offset int = 10, 0
	if paginate != nil {
		if paginate.Take != nil {
			limit = *paginate.Take
		}
		if paginate.Page != nil {
			offset = (*paginate.Page - 1) * limit
		}
	}
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Execute the query
	rows, err := repository.database.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []models.RemovalRequestCard{}
	for rows.Next() {
		var request models.RemovalRequestCard
		var decidedById sql.NullInt64
		var decidedByName sql.NullString
		var decidedByAvatar sql.NullString
		var decisionNote sql.NullString
		var decidedAtStr sql.NullString
		var createdAtStr string

		err := rows.Scan(
			&request.Id,
			&request.PostVersionId,
			&request.PostTitle,
			&request.RequestedBy.Id,
			&request.RequestedBy.Name,
			&request.RequestedBy.Avatar,
			&request.Note,
			&request.Status,
			&decidedById,
			&decidedByName,
			&decidedByAvatar,
			&decisionNote,
			&decidedAtStr,
			&createdAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse created_at timestamp
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}
		request.CreatedAt = createdAt

		// Parse decided_at timestamp if present
		if decidedAtStr.Valid {
			decidedAt, err := time.Parse("2006-01-02 15:04:05", decidedAtStr.String)
			if err != nil {
				return nil, err
			}
			request.DecidedAt = &decidedAt
		}

		// Format requested_by avatar URL
		if request.RequestedBy.Avatar != nil && *request.RequestedBy.Avatar != "" {
			avatarPath := fmt.Sprintf("/uploads/avatar/%s", *request.RequestedBy.Avatar)
			request.RequestedBy.Avatar = &avatarPath
		}

		if decidedById.Valid {
			var decidedByAvatarPath *string
			if decidedByAvatar.Valid && decidedByAvatar.String != "" {
				avatarPath := fmt.Sprintf("/uploads/avatar/%s", decidedByAvatar.String)
				decidedByAvatarPath = &avatarPath
			}
			request.DecidedBy = &models.UserInfo{
				Id:     decidedById.Int64,
				Name:   decidedByName.String,
				Avatar: decidedByAvatarPath,
			}
		}

		// Set decision note
		if decisionNote.Valid {
			request.DecisionNote = &decisionNote.String
		}

		requests = append(requests, request)
	}

	// Get total count for pagination
	countQuery := `
	SELECT COUNT(*)
	FROM removal_requests rr
	JOIN post_versions pv ON rr.post_version_id = pv.id
	JOIN users u1 ON rr.requested_by = u1.id
	LEFT JOIN users u2 ON rr.decided_by = u2.id
	WHERE 1=1`

	if len(conditions) > 0 {
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}

	var total int64
	err = repository.database.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	currentPage := 1
	if paginate != nil && paginate.Page != nil {
		currentPage = *paginate.Page
	}

	return &responses.PaginatedResponse[models.RemovalRequestCard]{
		Data:  requests,
		Total: total,
		Page:  currentPage,
		Take:  limit,
	}, nil
}

func (repository *RemovalRequestRepository) GetRemovalRequestById(
	id int64,
) (*models.RemovalRequestDetails, error) {
	row := repository.database.QueryRow(QueryGetRemovalRequestById, id)

	var request models.RemovalRequestDetails
	var decidedById sql.NullInt64
	var decidedByName sql.NullString
	var decidedByAvatar sql.NullString
	var decisionNote sql.NullString
	var decidedAtStr sql.NullString
	var createdAtStr string
	var postCoverUrl sql.NullString
	var postCategory sql.NullString

	err := row.Scan(
		&request.Id,
		&request.PostVersionId,
		&request.PostTitle,
		&request.PostWriter.Id,
		&request.PostWriter.Name,
		&request.PostWriter.Avatar,
		&postCoverUrl,
		&postCategory,
		&request.RequestedBy.Id,
		&request.RequestedBy.Name,
		&request.RequestedBy.Avatar,
		&request.Note,
		&request.Status,
		&decidedById,
		&decidedByName,
		&decidedByAvatar,
		&decisionNote,
		&decidedAtStr,
		&createdAtStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierrors.ErrNotFound
		}
		return nil, err
	}

	// Parse created_at timestamp
	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, err
	}
	request.CreatedAt = createdAt

	// Parse decided_at timestamp if present
	if decidedAtStr.Valid {
		decidedAt, err := time.Parse("2006-01-02 15:04:05", decidedAtStr.String)
		if err != nil {
			return nil, err
		}
		request.DecidedAt = &decidedAt
	}

	// Format post writer avatar URL
	if request.PostWriter.Avatar != nil && *request.PostWriter.Avatar != "" {
		avatarPath := fmt.Sprintf("/uploads/avatar/%s", *request.PostWriter.Avatar)
		request.PostWriter.Avatar = &avatarPath
	}

	// Format post cover URL
	if postCoverUrl.Valid && postCoverUrl.String != "" {
		nameWithoutExt := strings.TrimSuffix(postCoverUrl.String, filepath.Ext(postCoverUrl.String))
		coverPath := fmt.Sprintf("/uploads/cover/%s", nameWithoutExt)
		request.PostCoverUrl = &coverPath
	}

	// Set post category
	if postCategory.Valid {
		request.PostCategory = &postCategory.String
	}

	// Format requested_by avatar URL
	if request.RequestedBy.Avatar != nil && *request.RequestedBy.Avatar != "" {
		avatarPath := fmt.Sprintf("/uploads/avatar/%s", *request.RequestedBy.Avatar)
		request.RequestedBy.Avatar = &avatarPath
	}

	if decidedById.Valid {
		var decidedByAvatarPath *string
		if decidedByAvatar.Valid && decidedByAvatar.String != "" {
			avatarPath := fmt.Sprintf("/uploads/avatar/%s", decidedByAvatar.String)
			decidedByAvatarPath = &avatarPath
		}
		request.DecidedBy = &models.UserInfo{
			Id:     decidedById.Int64,
			Name:   decidedByName.String,
			Avatar: decidedByAvatarPath,
		}
	}

	// Set decision note
	if decisionNote.Valid {
		request.DecisionNote = &decisionNote.String
	}

	return &request, nil
}

func (repository *RemovalRequestRepository) GetUserRemovalRequests(
	userId int64,
) ([]models.RemovalRequestCard, error) {
	rows, err := repository.database.Query(QueryGetUserRemovalRequests, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []models.RemovalRequestCard{}
	for rows.Next() {
		var request models.RemovalRequestCard
		var decidedById sql.NullInt64
		var decidedByName sql.NullString
		var decidedByAvatar sql.NullString
		var decisionNote sql.NullString
		var decidedAtStr sql.NullString
		var createdAtStr string

		err := rows.Scan(
			&request.Id,
			&request.PostVersionId,
			&request.PostTitle,
			&request.RequestedBy.Id,
			&request.RequestedBy.Name,
			&request.RequestedBy.Avatar,
			&request.Note,
			&request.Status,
			&decidedById,
			&decidedByName,
			&decidedByAvatar,
			&decisionNote,
			&decidedAtStr,
			&createdAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse created_at timestamp
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}
		request.CreatedAt = createdAt

		// Parse decided_at timestamp if present
		if decidedAtStr.Valid {
			decidedAt, err := time.Parse("2006-01-02 15:04:05", decidedAtStr.String)
			if err != nil {
				return nil, err
			}
			request.DecidedAt = &decidedAt
		}

		// Format requested_by avatar URL
		if request.RequestedBy.Avatar != nil && *request.RequestedBy.Avatar != "" {
			avatarPath := fmt.Sprintf("/uploads/avatar/%s", *request.RequestedBy.Avatar)
			request.RequestedBy.Avatar = &avatarPath
		}

		if decidedById.Valid {
			var decidedByAvatarPath *string
			if decidedByAvatar.Valid && decidedByAvatar.String != "" {
				avatarPath := fmt.Sprintf("/uploads/avatar/%s", decidedByAvatar.String)
				decidedByAvatarPath = &avatarPath
			}
			request.DecidedBy = &models.UserInfo{
				Id:     decidedById.Int64,
				Name:   decidedByName.String,
				Avatar: decidedByAvatarPath,
			}
		}

		// Set decision note
		if decisionNote.Valid {
			request.DecisionNote = &decisionNote.String
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (repository *RemovalRequestRepository) GetVersionOwnerAndStatus(
	postVersionId int64,
) (*models.QueryGetVersionOwnerAndStatus, error) {
	row := repository.database.QueryRow(QueryGetVersionOwnerAndStatus, postVersionId)

	var result models.QueryGetVersionOwnerAndStatus
	err := row.Scan(&result.CreatedBy, &result.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierrors.ErrNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (repository *RemovalRequestRepository) CheckExistingRemovalRequest(
	postVersionId int64,
	requestedBy int64,
) (bool, error) {
	row := repository.database.QueryRow(
		QueryCheckExistingRemovalRequest,
		postVersionId,
		requestedBy,
	)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *RemovalRequestRepository) ApproveRemovalRequest(
	id int64,
	decidedBy int64,
	decisionNote *string,
) error {
	result, err := repository.database.Exec(QueryApproveRemovalRequest, decidedBy, decisionNote, id)
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

func (repository *RemovalRequestRepository) RejectRemovalRequest(
	id int64,
	decidedBy int64,
	decisionNote *string,
) error {
	result, err := repository.database.Exec(QueryRejectRemovalRequest, decidedBy, decisionNote, id)
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

func (repository *RemovalRequestRepository) AutoApproveOtherRemovalRequests(
	postVersionId int64,
	excludeRequestId int64,
	decidedBy int64,
	decisionNote *string,
) error {
	_, err := repository.database.Exec(
		QueryAutoApproveOtherRemovalRequests,
		decidedBy,
		decisionNote,
		postVersionId,
		excludeRequestId,
	)
	return err
}

func (repository *RemovalRequestRepository) IsVersionCurrentlyPublished(
	versionId int64,
) (bool, error) {
	row := repository.database.QueryRow(QueryCheckIfVersionIsCurrentlyPublished, versionId)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *RemovalRequestRepository) SetPostCurrentVersionToNull(
	versionId int64,
) error {
	_, err := repository.database.Exec(QuerySetPostCurrentVersionToNull, versionId)
	return err
}

func (repository *RemovalRequestRepository) SoftDeleteVersion(
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

func (repository *RemovalRequestRepository) GetVersionCoverImage(
	versionId int64,
) (*string, error) {
	row := repository.database.QueryRow(QueryGetVersionCoverImage, versionId)

	var coverImage sql.NullString
	err := row.Scan(&coverImage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if coverImage.Valid && coverImage.String != "" {
		return &coverImage.String, nil
	}
	return nil, nil
}

func (repository *RemovalRequestRepository) IsImageReferencedByOtherVersions(
	imagePath string,
	excludeVersionId int64,
) (bool, error) {
	row := repository.database.QueryRow(QueryCheckImageReferences, imagePath, excludeVersionId)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPostIdFromVersionId gets the post ID from a version ID
func (repository *RemovalRequestRepository) GetPostIdFromVersionId(
	versionId int64,
) (int64, error) {
	row := repository.database.QueryRow(QueryGetPostIdFromVersionId, versionId)

	var postId int64
	err := row.Scan(&postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, apierrors.ErrNotFound
		}
		return 0, err
	}

	return postId, nil
}

// VersionInfo holds minimal version information for image cleanup
type VersionInfo struct {
	CoverImage *string
}

// GetAllVersionsForPost gets all versions for a post
func (repository *RemovalRequestRepository) GetAllVersionsForPost(
	postId int64,
) ([]VersionInfo, error) {
	rows, err := repository.database.Query(QueryGetAllVersionsForPost, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := []VersionInfo{}
	for rows.Next() {
		var version VersionInfo
		var coverImage sql.NullString

		err := rows.Scan(&coverImage)
		if err != nil {
			return nil, err
		}

		if coverImage.Valid && coverImage.String != "" {
			version.CoverImage = &coverImage.String
		}

		versions = append(versions, version)
	}

	return versions, nil
}

// SoftDeleteAllVersionsForPost soft deletes all versions of a post
func (repository *RemovalRequestRepository) SoftDeleteAllVersionsForPost(
	postId int64,
) error {
	_, err := repository.database.Exec(QuerySoftDeleteAllVersionsForPost, postId)
	return err
}

// SoftDeletePost soft deletes a post
func (repository *RemovalRequestRepository) SoftDeletePost(
	postId int64,
) error {
	result, err := repository.database.Exec(QuerySoftDeletePost, postId)
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

// AutoApproveOtherRemovalRequestsForPost auto-approves all other pending removal requests for the same post
func (repository *RemovalRequestRepository) AutoApproveOtherRemovalRequestsForPost(
	postId int64,
	excludeRequestId int64,
	decidedBy int64,
	decisionNote *string,
) error {
	_, err := repository.database.Exec(
		QueryAutoApproveOtherRemovalRequestsForPost,
		decidedBy,
		decisionNote,
		postId,
		excludeRequestId,
	)
	return err
}
