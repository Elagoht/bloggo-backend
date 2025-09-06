package audit

const (
	QueryGetAuditLogs = `
	SELECT id, user_id, entity, entity_id, action, created_at
	FROM audit_logs
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?`

	QueryGetAuditLogsCount = `
	SELECT COUNT(*)
	FROM audit_logs`
)