package search

const (
	QuerySearchTags = `
	SELECT id, name as title, slug, NULL as avatar_url, NULL as cover_url, 'tag' as type
	FROM tags
	WHERE name LIKE ? AND deleted_at IS NULL
	ORDER BY name ASC
	LIMIT ?`

	QuerySearchCategories = `
	SELECT id, name as title, slug, NULL as avatar_url, NULL as cover_url, 'category' as type
	FROM categories
	WHERE name LIKE ? AND deleted_at IS NULL
	ORDER BY name ASC
	LIMIT ?`

	QuerySearchPosts = `
	SELECT DISTINCT pv.post_id as id, pv.title, pv.slug, NULL as avatar_url, pv.cover_image as cover_url, 'post' as type
	FROM post_versions pv
	WHERE pv.title LIKE ? AND pv.deleted_at IS NULL
	ORDER BY pv.title ASC
	LIMIT ?`

	QuerySearchUsers = `
	SELECT id, name as title, NULL as slug, avatar as avatar_url, NULL as cover_url, 'user' as type
	FROM users
	WHERE name LIKE ? AND deleted_at IS NULL
	ORDER BY name ASC
	LIMIT ?`

	QueryCountSearchTags = `
	SELECT COUNT(*) FROM tags WHERE name LIKE ? AND deleted_at IS NULL`

	QueryCountSearchCategories = `
	SELECT COUNT(*) FROM categories WHERE name LIKE ? AND deleted_at IS NULL`

	QueryCountSearchPosts = `
	SELECT COUNT(DISTINCT pv.post_id) FROM post_versions pv
	WHERE pv.title LIKE ? AND pv.deleted_at IS NULL`

	QueryCountSearchUsers = `
	SELECT COUNT(*) FROM users WHERE name LIKE ? AND deleted_at IS NULL`
)