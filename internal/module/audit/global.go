package audit

import (
	"bloggo/internal/db"
	"bloggo/internal/module/audit/models"
	"sync"
)

var (
	globalAuditService AuditService
	once               sync.Once
)

// GetGlobalAuditService returns the global audit service instance
func GetGlobalAuditService() AuditService {
	once.Do(func() {
		database := db.Get()
		repository := NewAuditRepository(database)
		globalAuditService = NewAuditService(repository)
	})
	return globalAuditService
}

// Helper functions for easy access to audit logging
func LogAction(userID *int64, entityType string, entityID int64, action string, oldValues, newValues map[string]interface{}) error {
	service := GetGlobalAuditService()
	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		OldValues:  oldValues,
		NewValues:  newValues,
	}
	return service.LogAction(entry)
}

func LogUserAction(userID, targetUserID *int64, action string, oldValues, newValues map[string]interface{}) error {
	service := GetGlobalAuditService()
	return service.LogUserAction(userID, targetUserID, action, oldValues, newValues)
}

func LogPostAction(userID *int64, postID int64, action string, oldValues, newValues map[string]interface{}) error {
	service := GetGlobalAuditService()
	return service.LogPostAction(userID, postID, action, oldValues, newValues)
}

func LogVersionAction(userID *int64, versionID int64, action string, oldValues, newValues map[string]interface{}, metadata map[string]interface{}) error {
	service := GetGlobalAuditService()
	return service.LogVersionAction(userID, versionID, action, oldValues, newValues, metadata)
}

func LogCategoryAction(userID *int64, categoryID int64, action string, oldValues, newValues map[string]interface{}) error {
	service := GetGlobalAuditService()
	return service.LogCategoryAction(userID, categoryID, action, oldValues, newValues)
}

func LogTagAction(userID *int64, tagID int64, action string, oldValues, newValues map[string]interface{}) error {
	service := GetGlobalAuditService()
	return service.LogTagAction(userID, tagID, action, oldValues, newValues)
}

func LogAuthAction(userID *int64, action string, oldValues, newValues map[string]interface{}) error {
	service := GetGlobalAuditService()
	return service.LogAuthAction(userID, action, oldValues, newValues)
}