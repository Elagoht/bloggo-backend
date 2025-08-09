package post

const (
	QueryPostGetByCurrentVersionSlug = `
	SELECT
		p.id as post_id, pv.id as version_id,
		u.name as author, u.email as author_email, u.avatar as author_avatar,
		pv.title, pv.slug, pv.content, pv.cover_image, pv.description, pv.spot,
		pv.status, pv.status_changed_at, pv.status_changed_by, pv.status_change_note,
		pv.created_by, pv.created_at, pv.updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	JOIN post_versions pv ON pv.id = p.current_version_id
	LEFT JOIN categories c ON c.id = pv.category_id
	LEFT JOIN users u ON u.id = p.created_by
	WHERE pv.slug = ?
		AND p.deleted_at IS NULL
	LIMIT 1;`
	QueryPostGetById = `
	SELECT
		p.id as post_id, pv.id as version_id,
		u.name as author, u.email as author_email, u.avatar as author_avatar,
		pv.title, pv.slug, pv.content, pv.cover_image, pv.description, pv.spot,
		pv.status, pv.status_changed_at, pv.status_changed_by, pv.status_change_note,
		pv.created_by, pv.created_at, pv.updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	JOIN post_versions pv ON pv.id = p.current_version_id
	LEFT JOIN categories c ON c.id = pv.category_id
	LEFT JOIN users u ON u.id = p.created_by
	WHERE p.id = ?
		AND p.deleted_at IS NULL
	LIMIT 1;`
	QueryPostGetList = `
	SELECT
		p.id as post_id,
		u.name as author, u.avatar as author_avatar,
		pv.title, pv.slug, pv.cover_image, pv.spot,
		pv.status,
		pv.created_at, pv.updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	JOIN post_versions pv ON pv.id = p.current_version_id
	LEFT JOIN categories c ON c.id = pv.category_id
	LEFT JOIN users u ON u.id = p.created_by
	WHERE p.deleted_at IS NULL;`
	QueryPostCreate = `
	INSERT INTO posts (created_by)
	VALUES (?);`
	QueryPostVersionCreate = `
	INSERT INTO post_versions (
		post_id, title, slug, content, cover_image,
		description, spot, category_id, created_by
	) VALUES (
		?, ?, ?, ?, ?,
		?, ?, ?, ?
	);`
	QueryPostSetCurrentVersion = `
	UPDATE posts
	SET current_version_id = ?
	WHERE id = ? AND deleted_at IS NULL;`
	QueryPostVersionsGetByPostId = `
	SELECT
		pv.id,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		pv.title, pv.status, pv.updated_at
	FROM post_versions pv
	LEFT JOIN users u
	ON pv.created_by = u.id
	WHERE pv.post_id = ? AND p.deleted_at IS NULL;;`
	QueryPostDetailsForVersionsGetByPostId = `
	SELECT
	p.current_version_id, p.created_at,
	u.id as author_id, u.name  as author_name, u.avatar as author_avatar
	FROM posts p
	LEFT JOIN users u
	ON u.id = p.created_by
	WHERE p.id = ? AND p.deleted_at IS NULL;;`
	QueryPostPatch      = ``
	QueryPostSoftDelete = `
	UPDATE posts
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE id = ? AND deleted_at IS NULL;`
	QueryPostAllRelatedCovers = `
	SELECT cover_image
	FROM post_versions
	WHERE post_id = ?;`
)
