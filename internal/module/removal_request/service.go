package removal_request

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/permissions"
	postmodels "bloggo/internal/module/post/models"
	"bloggo/internal/module/removal_request/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"bloggo/internal/utils/schemas/responses"
)

type RemovalRequestService struct {
	repository  RemovalRequestRepository
	permissions permissions.Store
	bucket      bucket.Bucket
}

func NewRemovalRequestService(
	repository RemovalRequestRepository,
	permissions permissions.Store,
	bucket bucket.Bucket,
) RemovalRequestService {
	return RemovalRequestService{
		repository,
		permissions,
		bucket,
	}
}

func (service *RemovalRequestService) CreateRemovalRequest(
	postVersionId int64,
	requestedBy int64,
	note *string,
) (*responses.ResponseCreated, error) {
	// Check if the version exists and get its status
	versionInfo, err := service.repository.GetVersionOwnerAndStatus(postVersionId)
	if err != nil {
		return nil, err
	}

	// Check if the version is published or approved (can only request removal for these)
	if versionInfo.Status != postmodels.STATUS_PUBLISHED && versionInfo.Status != postmodels.STATUS_APPROVED {
		return nil, apierrors.ErrPreconditionFailed
	}

	// Check if there's already a pending removal request for this version by this user
	exists, err := service.repository.CheckExistingRemovalRequest(postVersionId, requestedBy)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apierrors.ErrConflict
	}

	// Create the removal request
	id, err := service.repository.CreateRemovalRequest(postVersionId, requestedBy, note)
	if err != nil {
		return nil, err
	}

	return &responses.ResponseCreated{
		Id: id,
	}, nil
}

func (service *RemovalRequestService) GetRemovalRequestList(
	userRoleId int64,
	paginate *pagination.PaginationOptions,
	search *filter.SearchOptions,
	status *int,
) (*responses.PaginatedResponse[models.RemovalRequestCard], error) {
	// Check if user has permission to view all removal requests (editors/admins)
	hasPermission := service.permissions.HasPermission(userRoleId, "post:delete")
	if !hasPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetRemovalRequestList(paginate, search, status)
}

func (service *RemovalRequestService) GetRemovalRequestById(
	id int64,
	userId int64,
	userRoleId int64,
) (*models.RemovalRequestDetails, error) {
	request, err := service.repository.GetRemovalRequestById(id)
	if err != nil {
		return nil, err
	}

	// Check if user has permission to view this request
	hasPermission := service.permissions.HasPermission(userRoleId, "post:delete")
	isOwner := request.RequestedBy.Id == userId

	if !hasPermission && !isOwner {
		return nil, apierrors.ErrForbidden
	}

	return request, nil
}

func (service *RemovalRequestService) GetUserRemovalRequests(
	userId int64,
) ([]models.RemovalRequestCard, error) {
	return service.repository.GetUserRemovalRequests(userId)
}

func (service *RemovalRequestService) ApproveRemovalRequest(
	id int64,
	decidedBy int64,
	userRoleId int64,
) error {
	// Check if user has permission to approve removal requests (editors/admins)
	hasPermission := service.permissions.HasPermission(userRoleId, "post:delete")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	// Get the removal request to verify it exists and is pending
	request, err := service.repository.GetRemovalRequestById(id)
	if err != nil {
		return err
	}

	// Check if request is still pending
	if request.Status != models.STATUS_PENDING {
		return apierrors.ErrPreconditionFailed
	}

	// Get the cover image before deletion for cleanup
	coverImagePath, err := service.repository.GetVersionCoverImage(request.PostVersionId)
	if err != nil {
		return err
	}

	// Check if this version is currently published
	isCurrentlyPublished, err := service.repository.IsVersionCurrentlyPublished(request.PostVersionId)
	if err != nil {
		return err
	}

	// If it's currently published, set the post's current_version_id to NULL
	if isCurrentlyPublished {
		if err := service.repository.SetPostCurrentVersionToNull(request.PostVersionId); err != nil {
			return err
		}
	}

	// Soft delete the version
	if err := service.repository.SoftDeleteVersion(request.PostVersionId); err != nil {
		return err
	}

	// Clean up cover image if not referenced by other versions
	if coverImagePath != nil {
		isImageStillReferenced, err := service.repository.IsImageReferencedByOtherVersions(*coverImagePath, request.PostVersionId)
		if err != nil {
			// Log error but don't fail the operation
		} else if !isImageStillReferenced {
			// Delete the image from storage
			service.bucket.Delete(*coverImagePath)
		}
	}

	// Finally, approve the removal request
	return service.repository.ApproveRemovalRequest(id, decidedBy)
}

func (service *RemovalRequestService) RejectRemovalRequest(
	id int64,
	decidedBy int64,
	userRoleId int64,
) error {
	// Check if user has permission to reject removal requests (editors/admins)
	hasPermission := service.permissions.HasPermission(userRoleId, "post:delete")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	// Get the removal request to verify it exists and is pending
	request, err := service.repository.GetRemovalRequestById(id)
	if err != nil {
		return err
	}

	// Check if request is still pending
	if request.Status != models.STATUS_PENDING {
		return apierrors.ErrPreconditionFailed
	}

	// Reject the request
	return service.repository.RejectRemovalRequest(id, decidedBy)
}