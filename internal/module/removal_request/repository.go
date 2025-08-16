package removal_request

import (
	"bloggo/internal/module/removal_request/models"
	"bloggo/internal/utils/apierrors"
	"database/sql"
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

func (repository *RemovalRequestRepository) GetRemovalRequestList() (
	[]models.RemovalRequestCard,
	error,
) {
	rows, err := repository.database.Query(QueryGetRemovalRequestList)
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
			&request.DecidedAt,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if decidedById.Valid {
			request.DecidedBy = &models.UserInfo{
				Id:     decidedById.Int64,
				Name:   decidedByName.String,
				Avatar: &decidedByAvatar.String,
			}
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (repository *RemovalRequestRepository) GetRemovalRequestById(
	id int64,
) (*models.RemovalRequestDetails, error) {
	row := repository.database.QueryRow(QueryGetRemovalRequestById, id)

	var request models.RemovalRequestDetails
	var decidedById sql.NullInt64
	var decidedByName sql.NullString
	var decidedByAvatar sql.NullString

	err := row.Scan(
		&request.Id,
		&request.PostVersionId,
		&request.PostTitle,
		&request.PostContent,
		&request.RequestedBy.Id,
		&request.RequestedBy.Name,
		&request.RequestedBy.Avatar,
		&request.Note,
		&request.Status,
		&decidedById,
		&decidedByName,
		&decidedByAvatar,
		&request.DecidedAt,
		&request.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierrors.ErrNotFound
		}
		return nil, err
	}

	if decidedById.Valid {
		request.DecidedBy = &models.UserInfo{
			Id:     decidedById.Int64,
			Name:   decidedByName.String,
			Avatar: &decidedByAvatar.String,
		}
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
			&request.DecidedAt,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if decidedById.Valid {
			request.DecidedBy = &models.UserInfo{
				Id:     decidedById.Int64,
				Name:   decidedByName.String,
				Avatar: &decidedByAvatar.String,
			}
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
) error {
	result, err := repository.database.Exec(QueryApproveRemovalRequest, decidedBy, id)
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
) error {
	result, err := repository.database.Exec(QueryRejectRemovalRequest, decidedBy, id)
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