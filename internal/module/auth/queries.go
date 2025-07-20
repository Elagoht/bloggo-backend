package auth

const (
	QueryUserLoginDataByEmail = `
	SELECT id, role_id, passphrase_hash
	FROM users
	WHERE email = ? AND deleted_at IS NULL;`
)
