package search

const (
	QuerySearchTags = `
	SELECT id, name as title, slug, NULL as avatar_url, NULL as cover_url, 'tag' as type
	FROM tags
	WHERE name LIKE ?
	ORDER BY name ASC
	LIMIT ?`

	QuerySearchCategories = `
	SELECT id, name as title, slug, NULL as avatar_url, cover_url, 'category' as type
	FROM categories
	WHERE name LIKE ?
	ORDER BY name ASC
	LIMIT ?`

	QuerySearchPosts = `
	SELECT p.id, pv.title, p.slug, NULL as avatar_url, pv.cover_url, 'post' as type
	FROM posts p
	JOIN post_versions pv ON p.current_version_id = pv.id
	WHERE pv.title LIKE ? AND p.published = 1
	ORDER BY pv.title ASC
	LIMIT ?`

	QuerySearchUsers = `
	SELECT id, name as title, NULL as slug, avatar_url, NULL as cover_url, 'user' as type
	FROM users
	WHERE name LIKE ?
	ORDER BY name ASC
	LIMIT ?`

	QueryCountSearchTags = `
	SELECT COUNT(*) FROM tags WHERE name LIKE ?`

	QueryCountSearchCategories = `
	SELECT COUNT(*) FROM categories WHERE name LIKE ?`

	QueryCountSearchPosts = `
	SELECT COUNT(*) FROM posts p
	JOIN post_versions pv ON p.current_version_id = pv.id
	WHERE pv.title LIKE ? AND p.published = 1`

	QueryCountSearchUsers = `
	SELECT COUNT(*) FROM users WHERE name LIKE ?`
)