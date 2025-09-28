package audit

import (
	"database/sql"
	"log"
)

type AuditLogger struct {
	db *sql.DB
}

func NewAuditLogger(db *sql.DB) *AuditLogger {
	return &AuditLogger{
		db: db,
	}
}

func (a *AuditLogger) LogAction(userID *int64, entity string, entityID int64, action string) {
	query := `INSERT INTO audit_logs (user_id, entity_type, entity_id, action) VALUES (?, ?, ?, ?)`

	_, err := a.db.Exec(query, userID, entity, entityID, action)
	if err != nil {
		log.Printf("Failed to log audit action: %v", err)
	}
}

var GlobalAuditLogger *AuditLogger

func InitializeAuditLogger(db *sql.DB) {
	GlobalAuditLogger = NewAuditLogger(db)
}

func LogAction(userID *int64, entity string, entityID int64, action string) {
	if GlobalAuditLogger != nil {
		GlobalAuditLogger.LogAction(userID, entity, entityID, action)
	}
}
