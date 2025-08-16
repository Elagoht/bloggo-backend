package user

const (
	QueryUserGetUserCards = `
	SELECT users.id, users.name, users.email, users.avatar, users.role_id, roles.name AS role_name,
  (SELECT COUNT(*) FROM posts WHERE created_by = users.id AND deleted_at IS NULL) AS writtenPostCount,
  (SELECT COUNT(*) FROM posts WHERE created_by = users.id AND deleted_at IS NULL AND current_version_id IS NOT NULL) AS publishedPostCount
	FROM users
	LEFT JOIN roles ON users.role_id = roles.id
	WHERE users.deleted_at IS NULL%s;`
	QueryUserGetById = `SELECT users.id, users.name, users.email, users.avatar, users.created_at, users.last_login, users.role_id, roles.name AS role_name,
  (SELECT COUNT(*) FROM posts WHERE created_by = users.id AND deleted_at IS NULL) AS writtenPostCount,
  (SELECT COUNT(*) FROM posts WHERE created_by = users.id AND deleted_at IS NULL AND current_version_id IS NOT NULL) AS publishedPostCount
	FROM users
	LEFT JOIN roles ON users.role_id = roles.id
	WHERE users.id = ? AND users.deleted_at IS NULL;`
	QueryUserCreate = `
	INSERT INTO users (name, email, avatar, passphrase_hash, role_id)
	VALUES (?, ?, ?, ?, ?);`
	QueryUserUpdateById = `
	UPDATE users
	SET 
		name = COALESCE(?, name),
		email = COALESCE(?, email),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ? AND deleted_at IS NULL;`
	QueryUserUpdateAvatarById = `
	UPDATE users
	SET avatar = ?
	WHERE id = ? AND deleted_at IS NULL;`
	QueryUserAssignRole = `
	UPDATE users
	SET 
		role_id = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ? AND deleted_at IS NULL;`
	QueryUserDelete = `
	UPDATE users
	SET 
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ? AND deleted_at IS NULL;`
)
