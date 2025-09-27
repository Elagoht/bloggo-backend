package audit

import "bloggo/internal/module/audit/models"

// AuditLogger is the interface for audit logging
type AuditLogger interface {
	LogAction(entry *models.AuditLogEntry) error
	LogUserAction(userID, targetUserID *int64, action string, oldValues, newValues map[string]interface{}) error
	LogPostAction(userID *int64, postID int64, action string, oldValues, newValues map[string]interface{}) error
	LogVersionAction(userID *int64, versionID int64, action string, oldValues, newValues map[string]interface{}, metadata map[string]interface{}) error
	LogCategoryAction(userID *int64, categoryID int64, action string, oldValues, newValues map[string]interface{}) error
	LogTagAction(userID *int64, tagID int64, action string, oldValues, newValues map[string]interface{}) error
	LogAuthAction(userID *int64, action string, oldValues, newValues map[string]interface{}) error
}