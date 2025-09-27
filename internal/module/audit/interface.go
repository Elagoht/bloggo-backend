package audit

import "bloggo/internal/module/audit/models"

// AuditLogger is the interface for audit logging
type AuditLogger interface {
	LogAction(entry *models.AuditLogEntry) error
	LogUserAction(userID, targetUserID *int64, action string) error
	LogPostAction(userID *int64, postID int64, action string) error
	LogVersionAction(userID *int64, versionID int64, action string, metadata map[string]interface{}) error
	LogCategoryAction(userID *int64, categoryID int64, action string) error
	LogTagAction(userID *int64, tagID int64, action string) error
	LogAuthAction(userID *int64, action string) error
}