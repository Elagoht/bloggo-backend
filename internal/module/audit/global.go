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
func LogAction(userID *int64, entityType string, entityID int64, action string) error {
	service := GetGlobalAuditService()
	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
	}
	return service.LogAction(entry)
}

func LogUserAction(userID, targetUserID *int64, action string) error {
	service := GetGlobalAuditService()
	return service.LogUserAction(userID, targetUserID, action)
}

func LogPostAction(userID *int64, postID int64, action string) error {
	service := GetGlobalAuditService()
	return service.LogPostAction(userID, postID, action)
}

func LogVersionAction(userID *int64, versionID int64, action string, metadata map[string]interface{}) error {
	service := GetGlobalAuditService()
	return service.LogVersionAction(userID, versionID, action, metadata)
}

func LogCategoryAction(userID *int64, categoryID int64, action string) error {
	service := GetGlobalAuditService()
	return service.LogCategoryAction(userID, categoryID, action)
}

func LogTagAction(userID *int64, tagID int64, action string) error {
	service := GetGlobalAuditService()
	return service.LogTagAction(userID, tagID, action)
}

func LogAuthAction(userID *int64, action string) error {
	service := GetGlobalAuditService()
	return service.LogAuthAction(userID, action)
}

func LogWebhookAction(userID *int64, action string) error {
	service := GetGlobalAuditService()
	entry := &models.AuditLogEntry{
		UserID:     userID,
		EntityType: "webhook",
		EntityID:   1, // Webhook config is always ID 1
		Action:     action,
	}
	return service.LogAction(entry)
}