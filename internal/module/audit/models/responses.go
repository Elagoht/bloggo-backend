package models

import "time"

type AuditLogResponse struct {
	ID         int64                   `json:"id"`
	UserID     *int64                  `json:"userId,omitempty"`
	UserName   *string                 `json:"userName,omitempty"`
	EntityType string                  `json:"entityType"`
	EntityID   int64                   `json:"entityId"`
	EntityName *string                 `json:"entityName,omitempty"`
	Action     string                  `json:"action"`
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time               `json:"createdAt"`
}

type AuditLogListResponse struct {
	Logs  []AuditLogResponse `json:"logs"`
	Total int                `json:"total"`
}