package keyvalues

const (
	QueryAPIGetAllKeyValues = `
	SELECT key, value
	FROM key_value_store
	ORDER BY key ASC;`

	QueryAPIGetKeyValueByKey = `
	SELECT key, value
	FROM key_value_store
	WHERE key = ?;`

	QueryAPIGetKeyValuesStartingWith = `
	SELECT key, value
	FROM key_value_store
	WHERE key LIKE ?
	ORDER BY key ASC;`
)
