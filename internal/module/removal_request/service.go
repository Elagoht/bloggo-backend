package removal_request

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/removal_request/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/audit"
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
	// Check if the version exists (but allow all statuses)
	_, err := service.repository.GetVersionOwnerAndStatus(postVersionId)
	if err != nil {
		return nil, err
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

	// Log the audit event
	audit.LogAction(&requestedBy, "removal_request", id, "requested")

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
	decisionNote *string,
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

	// Get the post ID from the version ID
	postId, err := service.repository.GetPostIdFromVersionId(request.PostVersionId)
	if err != nil {
		return err
	}

	// Get all versions for this post to collect cover images
	versions, err := service.repository.GetAllVersionsForPost(postId)
	if err != nil {
		return err
	}

	// Collect all unique cover images for later cleanup
	coverImages := make(map[string]bool)
	for _, version := range versions {
		if version.CoverImage != nil && *version.CoverImage != "" {
			coverImages[*version.CoverImage] = true
		}
	}

	// Soft delete all versions of the post
	if err := service.repository.SoftDeleteAllVersionsForPost(postId); err != nil {
		return err
	}

	// Soft delete the post itself
	if err := service.repository.SoftDeletePost(postId); err != nil {
		return err
	}

	// Clean up all cover images from storage
	for imagePath := range coverImages {
		// Delete the image from storage
		service.bucket.Delete(imagePath)
	}

	// Finally, approve the removal request
	err = service.repository.ApproveRemovalRequest(id, decidedBy, decisionNote)
	if err != nil {
		return err
	}

	// Auto-approve all other pending removal requests for the same post
	autoApprovalNote := "Automatically approved - post was already deleted"
	err = service.repository.AutoApproveOtherRemovalRequestsForPost(
		postId,
		id,
		decidedBy,
		&autoApprovalNote,
	)
	if err != nil {
		// Log error but don't fail the operation
		// The main request was already approved successfully
	}

	// Log the audit events
	audit.LogAction(&decidedBy, "removal_request", id, "approved")
	audit.LogAction(&decidedBy, "post", postId, "deleted")

	return nil
}

func (service *RemovalRequestService) RejectRemovalRequest(
	id int64,
	decidedBy int64,
	userRoleId int64,
	decisionNote *string,
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
	err = service.repository.RejectRemovalRequest(id, decidedBy, decisionNote)
	if err != nil {
		return err
	}

	// Log the audit event
	audit.LogAction(&decidedBy, "removal_request", id, "rejected")

	return nil
}