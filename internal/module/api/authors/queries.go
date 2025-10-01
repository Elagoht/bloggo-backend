package authors

const (
	// Authors queries
	QueryAPIGetAllAuthors = `
	SELECT
		u.id,
		u.name,
		u.avatar,
		u.created_at,
		COUNT(DISTINCT p.id) as published_post_count
	FROM users u
	LEFT JOIN posts p ON p.created_by = u.id
		AND p.deleted_at IS NULL
	LEFT JOIN post_versions pv ON pv.id = p.current_version_id
		AND pv.deleted_at IS NULL
		AND pv.status = 5
	WHERE u.deleted_at IS NULL
	GROUP BY u.id, u.name, u.avatar, u.created_at
	HAVING published_post_count > 0
	ORDER BY u.name ASC;`

	QueryAPIGetAuthorById = `
	SELECT
		u.id,
		u.name,
		u.avatar,
		u.created_at,
		COUNT(DISTINCT p.id) as published_post_count
	FROM users u
	LEFT JOIN posts p ON p.created_by = u.id
		AND p.deleted_at IS NULL
	LEFT JOIN post_versions pv ON pv.id = p.current_version_id
		AND pv.deleted_at IS NULL
		AND pv.status = 5
	WHERE u.deleted_at IS NULL
		AND u.id = ?
	GROUP BY u.id, u.name, u.avatar, u.created_at
	LIMIT 1;`
)
