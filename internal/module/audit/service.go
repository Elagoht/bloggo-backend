package audit

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/audit/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/pagination"
	"bloggo/internal/utils/schemas/responses"
)

type AuditService struct {
	repository  AuditRepository
	permissions permissions.Store
}

func NewAuditService(repository AuditRepository, permissions permissions.Store) AuditService {
	return AuditService{
		repository,
		permissions,
	}
}

func (service *AuditService) GetAuditLogs(
	pagination *pagination.PaginationOptions,
	userRoleId int64,
) (*responses.PaginatedResponse[models.ResponseAuditLog], error) {
	// Check if user has permission to view audit logs
	hasPermission := service.permissions.HasPermission(userRoleId, "auditlog:view")
	if !hasPermission {
		return nil, apierrors.ErrForbidden
	}

	// Get the audit logs data
	auditLogs, err := service.repository.GetAuditLogs(pagination)
	if err != nil {
		return nil, err
	}

	// Get the total count
	total, err := service.repository.GetAuditLogsCount()
	if err != nil {
		return nil, err
	}

	// Set default values for page and take if they're nil
	page := 1
	if pagination.Page != nil {
		page = *pagination.Page
	}

	take := 20 // default take value
	if pagination.Take != nil {
		take = *pagination.Take
	}

	return &responses.PaginatedResponse[models.ResponseAuditLog]{
		Data:  auditLogs,
		Page:  page,
		Take:  take,
		Total: int64(total),
	}, nil
}