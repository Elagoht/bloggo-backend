package webhook

const (
	// Config queries
	QueryGetConfig = `
	SELECT id, url, updated_at
	FROM webhook_config
	WHERE id = 1;`

	QueryUpsertConfig = `
	INSERT INTO webhook_config (id, url, updated_at)
	VALUES (1, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(id) DO UPDATE SET
		url = excluded.url,
		updated_at = CURRENT_TIMESTAMP;`

	// Header queries
	QueryGetAllHeaders = `
	SELECT id, key, value, created_at, updated_at
	FROM webhook_headers
	ORDER BY key ASC;`

	QueryUpsertHeader = `
	INSERT INTO webhook_headers (key, value, created_at, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT(key) DO UPDATE SET
		value = excluded.value,
		updated_at = CURRENT_TIMESTAMP;`

	QueryDeleteAllHeaders = `
	DELETE FROM webhook_headers;`

	// Request queries
	QueryInsertRequest = `
	INSERT INTO webhook_requests (
		event, entity, entity_id, slug, request_body,
		response_status, response_body, attempt_count, error_message, webhook_url, webhook_headers, created_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);`

	QueryUpdateRequest = `
	UPDATE webhook_requests
	SET response_status = ?, response_body = ?, attempt_count = ?, error_message = ?
	WHERE id = ?;`

	QueryGetAllRequests = `
	SELECT id, event, entity, entity_id, slug, request_body,
		response_status, response_body, attempt_count, error_message, webhook_url, webhook_headers, created_at
	FROM webhook_requests
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?;`

	QueryGetRequestsBySearch = `
	SELECT id, event, entity, entity_id, slug, request_body,
		response_status, response_body, attempt_count, error_message, webhook_url, webhook_headers, created_at
	FROM webhook_requests
	WHERE event LIKE ? OR entity LIKE ?
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?;`

	QueryGetRequestByID = `
	SELECT id, event, entity, entity_id, slug, request_body,
		response_status, response_body, attempt_count, error_message, webhook_url, webhook_headers, created_at
	FROM webhook_requests
	WHERE id = ?;`

	QueryCountRequests = `
	SELECT COUNT(*) FROM webhook_requests;`

	QueryCountRequestsBySearch = `
	SELECT COUNT(*) FROM webhook_requests
	WHERE event LIKE ? OR entity LIKE ?;`
)
