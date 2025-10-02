package keyvalue

const (
	QueryGetAll = `
	SELECT key, value, created_at, updated_at
	FROM key_value_store
	ORDER BY key ASC;`

	QueryUpsert = `
	INSERT INTO key_value_store (key, value, created_at, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT(key) DO UPDATE SET
		value = excluded.value,
		updated_at = CURRENT_TIMESTAMP;`

	QueryDelete = `
	DELETE FROM key_value_store WHERE key = ?;`

	QueryDeleteAll = `
	DELETE FROM key_value_store;`
)
