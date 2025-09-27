package audit

import (
	"bloggo/internal/module/audit/models"
)

type AuditService struct {
	repository AuditRepository
}

func NewAuditService(repository AuditRepository) AuditService {
	return AuditService{
		repository,
	}
}

// LogAction logs an audit event
func (service *AuditService) LogAction(entry *models.AuditLogEntry) error {
	return service.repository.LogAction(entry)
}

// GetAuditLogs retrieves audit logs with pagination
func (service *AuditService) GetAuditLogs(limit, offset int) (models.AuditLogListResponse, error) {
	logs, err := service.repository.GetAuditLogs(limit, offset)
	if err != nil {
		return models.AuditLogListResponse{}, err
	}

	// Ensure logs is never nil
	if logs == nil {
		logs = []models.AuditLogResponse{}
	}

	total, err := service.repository.CountAuditLogs()
	if err != nil {
		return models.AuditLogListResponse{}, err
	}

	return models.AuditLogListResponse{
		Logs:  logs,
		Total: total,
	}, nil
}

// GetAuditLogsWithFilters retrieves audit logs with filters and pagination
func (service *AuditService) GetAuditLogsWithFilters(limit, offset int, userIDs []int64, entityTypes, actions []string) (models.AuditLogListResponse, error) {
	logs, err := service.repository.GetAuditLogsWithFilters(limit, offset, userIDs, entityTypes, actions)
	if err != nil {
		return models.AuditLogListResponse{}, err
	}

	// Ensure logs is never nil
	if logs == nil {
		logs = []models.AuditLogResponse{}
	}

	total, err := service.repository.CountAuditLogsWithFilters(userIDs, entityTypes, actions)
	if err != nil {
		return models.AuditLogListResponse{}, err
	}

	return models.AuditLogListResponse{
		Logs:  logs,
		Total: total,
	}, nil
}

// GetAuditLogsByEntity retrieves audit logs for a specific entity
func (service *AuditService) GetAuditLogsByEntity(entityType string, entityID int64, limit, offset int) (models.AuditLogListResponse, error) {
	logs, err := service.repository.GetAuditLogsByEntity(entityType, entityID, limit, offset)
	if err != nil {
		return models.AuditLogListResponse{}, err
	}

	// Ensure logs is never nil
	if logs == nil {
		logs = []models.AuditLogResponse{}
	}

	total, err := service.repository.CountAuditLogsByEntity(entityType, entityID)
	if err != nil {
		return models.AuditLogListResponse{}, err
	}

	return models.AuditLogListResponse{
		Logs:  logs,
		Total: total,
	}, nil
}

// GetAuditLogsByUser retrieves audit logs for a specific user
func (service *AuditService) GetAuditLogsByUser(userID int64, limit, offset int) (models.AuditLogListResponse, error) {
	logs, err := service.repository.GetAuditLogsByUser(userID, limit, offset)
	if err != nil {
		return models.AuditLogListResponse{}, err
	}

	// Ensure logs is never nil
	if logs == nil {
		logs = []models.AuditLogResponse{}
	}

	total, err := service.repository.CountAuditLogsByUser(userID)
	if err != nil {
		return models.AuditLogListResponse{}, err
	}

	return models.AuditLogListResponse{
		Logs:  logs,
		Total: total,
	}, nil
}

// Helper functions to create audit log entries

// LogUserAction logs a user-related action
func (service *AuditService) LogUserAction(userID, targetUserID *int64, action string) error {
	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: models.EntityUser,
		EntityID:   *targetUserID,
		Action:     action,
	}

	return service.LogAction(entry)
}

// LogPostAction logs a post-related action
func (service *AuditService) LogPostAction(userID *int64, postID int64, action string) error {
	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: models.EntityPost,
		EntityID:   postID,
		Action:     action,
	}

	return service.LogAction(entry)
}

// LogVersionAction logs a post version-related action
func (service *AuditService) LogVersionAction(userID *int64, versionID int64, action string, metadata map[string]interface{}) error {
	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: models.EntityPostVersion,
		EntityID:   versionID,
		Action:     action,
		Metadata:   metadata,
	}

	return service.LogAction(entry)
}

// LogCategoryAction logs a category-related action
func (service *AuditService) LogCategoryAction(userID *int64, categoryID int64, action string) error {
	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: models.EntityCategory,
		EntityID:   categoryID,
		Action:     action,
	}

	return service.LogAction(entry)
}

// LogTagAction logs a tag-related action
func (service *AuditService) LogTagAction(userID *int64, tagID int64, action string) error {
	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: models.EntityTag,
		EntityID:   tagID,
		Action:     action,
	}

	return service.LogAction(entry)
}

// LogAuthAction logs an authentication-related action
func (service *AuditService) LogAuthAction(userID *int64, action string) error {
	// For auth actions, we use the userID as the entityID since auth is tied to a specific user
	// If userID is nil (system action), we'll use 0 as placeholder
	entityID := int64(0)
	if userID != nil {
		entityID = *userID
	}

	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: models.EntityAuth,
		EntityID:   entityID,
		Action:     action,
	}

	return service.LogAction(entry)
}

