package audit

const (
	QueryInsertAuditLog = `
	INSERT INTO audit_logs (
		user_id, entity_type, entity_id, action,
		old_values, new_values, metadata
	) VALUES (?, ?, ?, ?, ?, ?, ?);`

	QueryGetAuditLogs = `
	SELECT
		al.id, al.user_id, u.name as user_name,
		al.entity_type, al.entity_id, al.action,
		al.old_values, al.new_values, al.metadata,
		al.created_at
	FROM audit_logs al
	LEFT JOIN users u ON al.user_id = u.id
	ORDER BY al.created_at DESC
	LIMIT ? OFFSET ?;`

	QueryGetAuditLogsByEntity = `
	SELECT
		al.id, al.user_id, u.name as user_name,
		al.entity_type, al.entity_id, al.action,
		al.old_values, al.new_values, al.metadata,
		al.created_at
	FROM audit_logs al
	LEFT JOIN users u ON al.user_id = u.id
	WHERE al.entity_type = ? AND al.entity_id = ?
	ORDER BY al.created_at DESC
	LIMIT ? OFFSET ?;`

	QueryGetAuditLogsByUser = `
	SELECT
		al.id, al.user_id, u.name as user_name,
		al.entity_type, al.entity_id, al.action,
		al.old_values, al.new_values, al.metadata,
		al.created_at
	FROM audit_logs al
	LEFT JOIN users u ON al.user_id = u.id
	WHERE al.user_id = ?
	ORDER BY al.created_at DESC
	LIMIT ? OFFSET ?;`

	QueryCountAuditLogs = `
	SELECT COUNT(*) FROM audit_logs;`

	QueryCountAuditLogsByEntity = `
	SELECT COUNT(*) FROM audit_logs
	WHERE entity_type = ? AND entity_id = ?;`

	QueryCountAuditLogsByUser = `
	SELECT COUNT(*) FROM audit_logs
	WHERE user_id = ?;`
)