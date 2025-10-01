package posts

const (
	// Posts queries - Only published posts (status = 5)
	QueryAPIGetPublishedPosts = `
	SELECT
		pv.slug,
		pv.title,
		pv.description,
		pv.spot,
		pv.cover_image,
		p.read_count,
		pv.read_time,
		pv.updated_at as published_at,
		u.id as author_id,
		u.name as author_name,
		u.avatar as author_avatar,
		c.slug as category_slug,
		c.name as category_name
	FROM posts p
	JOIN post_versions pv ON pv.id = p.current_version_id
	JOIN users u ON u.id = p.created_by
	JOIN categories c ON c.id = pv.category_id
	WHERE p.deleted_at IS NULL
		AND pv.deleted_at IS NULL
		AND pv.status = 5
		AND c.deleted_at IS NULL
		%s
	%s;`

	QueryAPICountPublishedPosts = `
	SELECT COUNT(*)
	FROM posts p
	JOIN post_versions pv ON pv.id = p.current_version_id
	JOIN users u ON u.id = p.created_by
	JOIN categories c ON c.id = pv.category_id
	WHERE p.deleted_at IS NULL
		AND pv.deleted_at IS NULL
		AND pv.status = 5
		AND c.deleted_at IS NULL
		%s;`

	QueryAPIGetPublishedPostBySlug = `
	SELECT
		pv.slug,
		pv.title,
		pv.content,
		pv.description,
		pv.spot,
		pv.cover_image,
		p.read_count,
		pv.read_time,
		pv.updated_at as published_at,
		pv.updated_at,
		p.id as post_id,
		u.id as author_id,
		u.name as author_name,
		u.avatar as author_avatar,
		c.slug as category_slug,
		c.name as category_name,
		c.description as category_description
	FROM posts p
	JOIN post_versions pv ON pv.id = p.current_version_id
	JOIN users u ON u.id = p.created_by
	JOIN categories c ON c.id = pv.category_id
	WHERE p.deleted_at IS NULL
		AND pv.deleted_at IS NULL
		AND pv.status = 5
		AND c.deleted_at IS NULL
		AND pv.slug = ?
	LIMIT 1;`

	QueryAPIGetPostTags = `
	SELECT t.slug, t.name
	FROM tags t
	JOIN post_tags pt ON pt.tag_id = t.id
	WHERE pt.post_id = ?
		AND t.deleted_at IS NULL;`
)
