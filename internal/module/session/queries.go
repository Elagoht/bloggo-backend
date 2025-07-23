package session

const (
	QuerySessionCreateDataByEmail = `
	SELECT u.id, u.name, u.role_id, r.name AS role_name, u.passphrase_hash
	FROM users u
	JOIN roles r ON r.id = u.role_id
	WHERE u.email = ? AND u.deleted_at IS NULL;`
	QuerySessionCreateDataById = `
	SELECT u.id, u.name, u.role_id, r.name AS role_name, u.passphrase_hash
	FROM users u
	JOIN roles r ON r.id = u.role_id
	WHERE u.id = ? AND u.deleted_at IS NULL;`
	QueryUserPermissionsById = `
	SELECT p.name
	FROM users u
	JOIN roles r ON u.role_id = r.id
	JOIN role_permissions rp ON r.id = rp.role_id
	JOIN permissions p ON rp.permission_id = p.id
	WHERE u.id = ? AND u.deleted_at IS NULL;`
)
