package audit

const (
	QueryInsertAuditLog = `
	INSERT INTO audit_logs (
		user_id, entity_type, entity_id, action,
		metadata
	) VALUES (?, ?, ?, ?, ?);`

	QueryGetAuditLogs = `
	SELECT
		al.id, al.user_id, u.name as user_name,
		al.entity_type, al.entity_id, al.action,
		al.metadata,
		al.created_at,
		CASE
			WHEN al.entity_type = 'user' THEN (SELECT name FROM users WHERE id = al.entity_id)
			WHEN al.entity_type = 'category' THEN (SELECT name FROM categories WHERE id = al.entity_id)
			WHEN al.entity_type = 'tag' THEN (SELECT name FROM tags WHERE id = al.entity_id)
			WHEN al.entity_type = 'post' THEN (SELECT pv.title FROM post_versions pv JOIN posts p ON p.current_version_id = pv.id WHERE p.id = al.entity_id)
			WHEN al.entity_type = 'post_version' THEN (SELECT title FROM post_versions WHERE id = al.entity_id)
			WHEN al.entity_type = 'removal_request' THEN (SELECT pv.title FROM post_versions pv JOIN removal_requests rr ON rr.post_version_id = pv.id WHERE rr.id = al.entity_id)
			ELSE NULL
		END as entity_name
	FROM audit_logs al
	LEFT JOIN users u ON al.user_id = u.id
	ORDER BY al.created_at DESC
	LIMIT ? OFFSET ?;`

	QueryGetAuditLogsByEntity = `
	SELECT
		al.id, al.user_id, u.name as user_name,
		al.entity_type, al.entity_id, al.action,
		al.metadata,
		al.created_at,
		CASE
			WHEN al.entity_type = 'user' THEN (SELECT name FROM users WHERE id = al.entity_id)
			WHEN al.entity_type = 'category' THEN (SELECT name FROM categories WHERE id = al.entity_id)
			WHEN al.entity_type = 'tag' THEN (SELECT name FROM tags WHERE id = al.entity_id)
			WHEN al.entity_type = 'post' THEN (SELECT pv.title FROM post_versions pv JOIN posts p ON p.current_version_id = pv.id WHERE p.id = al.entity_id)
			WHEN al.entity_type = 'post_version' THEN (SELECT title FROM post_versions WHERE id = al.entity_id)
			WHEN al.entity_type = 'removal_request' THEN (SELECT pv.title FROM post_versions pv JOIN removal_requests rr ON rr.post_version_id = pv.id WHERE rr.id = al.entity_id)
			ELSE NULL
		END as entity_name
	FROM audit_logs al
	LEFT JOIN users u ON al.user_id = u.id
	WHERE al.entity_type = ? AND al.entity_id = ?
	ORDER BY al.created_at DESC
	LIMIT ? OFFSET ?;`

	QueryGetAuditLogsByUser = `
	SELECT
		al.id, al.user_id, u.name as user_name,
		al.entity_type, al.entity_id, al.action,
		al.metadata,
		al.created_at,
		CASE
			WHEN al.entity_type = 'user' THEN (SELECT name FROM users WHERE id = al.entity_id)
			WHEN al.entity_type = 'category' THEN (SELECT name FROM categories WHERE id = al.entity_id)
			WHEN al.entity_type = 'tag' THEN (SELECT name FROM tags WHERE id = al.entity_id)
			WHEN al.entity_type = 'post' THEN (SELECT pv.title FROM post_versions pv JOIN posts p ON p.current_version_id = pv.id WHERE p.id = al.entity_id)
			WHEN al.entity_type = 'post_version' THEN (SELECT title FROM post_versions WHERE id = al.entity_id)
			WHEN al.entity_type = 'removal_request' THEN (SELECT pv.title FROM post_versions pv JOIN removal_requests rr ON rr.post_version_id = pv.id WHERE rr.id = al.entity_id)
			ELSE NULL
		END as entity_name
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

	QueryGetAuditLogsWithFiltersBase = `
	SELECT
		al.id, al.user_id, u.name as user_name,
		al.entity_type, al.entity_id, al.action,
		al.metadata,
		al.created_at,
		CASE
			WHEN al.entity_type = 'user' THEN (SELECT name FROM users WHERE id = al.entity_id)
			WHEN al.entity_type = 'category' THEN (SELECT name FROM categories WHERE id = al.entity_id)
			WHEN al.entity_type = 'tag' THEN (SELECT name FROM tags WHERE id = al.entity_id)
			WHEN al.entity_type = 'post' THEN (SELECT pv.title FROM post_versions pv JOIN posts p ON p.current_version_id = pv.id WHERE p.id = al.entity_id)
			WHEN al.entity_type = 'post_version' THEN (SELECT title FROM post_versions WHERE id = al.entity_id)
			WHEN al.entity_type = 'removal_request' THEN (SELECT pv.title FROM post_versions pv JOIN removal_requests rr ON rr.post_version_id = pv.id WHERE rr.id = al.entity_id)
			ELSE NULL
		END as entity_name
	FROM audit_logs al
	LEFT JOIN users u ON al.user_id = u.id
	WHERE 1=1`

	QueryCountAuditLogsWithFiltersBase = `
	SELECT COUNT(*) FROM audit_logs WHERE 1=1`
)