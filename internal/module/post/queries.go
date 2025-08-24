package post

const (
	QueryPostGetByCurrentVersionSlug = `
	SELECT
		p.id as post_id, pv.id as version_id,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		pv.title, pv.slug, pv.content, pv.cover_image, pv.description, pv.spot,
		pv.status, p.read_count, pv.created_at, pv.updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	JOIN post_versions pv
	ON pv.id = p.current_version_id
	LEFT JOIN categories c
	ON c.id = pv.category_id
	LEFT JOIN users u
	ON u.id = p.created_by
	WHERE pv.slug = ?
		AND p.deleted_at IS NULL
	LIMIT 1;`
	QueryPostGetById = `
	SELECT
		p.id as post_id, pv.id as version_id,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		pv.title, pv.slug, pv.content, pv.cover_image, pv.description, pv.spot,
		pv.status, p.read_count, pv.created_at, pv.updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	JOIN post_versions pv
	ON pv.id = p.current_version_id
	LEFT JOIN categories c
	ON c.id = pv.category_id
	LEFT JOIN users u
	ON u.id = p.created_by
	WHERE p.id = ?
		AND p.deleted_at IS NULL
	LIMIT 1;`
	QueryPostGetList = `
	SELECT
		p.id as post_id,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		pv.title, pv.slug, pv.cover_image, pv.spot,
		pv.status, p.read_count,
		pv.created_at, pv.updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	JOIN post_versions pv
	ON pv.id = p.current_version_id
	LEFT JOIN categories c
	ON c.id = pv.category_id
	LEFT JOIN users u
	ON u.id = p.created_by
	WHERE p.deleted_at IS NULL %s;`
	QueryPostVersionGetById = `
	SELECT
		pv.id, pv.duplicated_from,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		pv.title, pv.slug, pv.content, pv.cover_image, pv.description, pv.spot,
		pv.status, pv.status_changed_at, pv.status_changed_by, pv.status_change_note,
		pv.created_at, pv.updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	JOIN post_versions pv
	ON pv.post_id = p.id
	LEFT JOIN categories c
	ON c.id = pv.category_id
	LEFT JOIN users u
	ON u.id = pv.created_by
	WHERE p.id = ?
	AND pv.id = ?
	AND p.deleted_at IS NULL
	AND pv.deleted_at IS NULL;`
	QueryPostCreate = `
	INSERT INTO posts (created_by)
	VALUES (?);`
	QueryPostVersionCreate = `
	INSERT INTO post_versions (
		post_id, title, slug, content, cover_image,
		description, spot, category_id, read_time, created_by,
		duplicated_from
	) VALUES (
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?
	);`
	QueryPostSetCurrentVersion = `
	UPDATE posts
	SET current_version_id = ?
	WHERE id = ?
	AND deleted_at IS NULL;`
	QueryPostVersionsGetByPostId = `
	SELECT
		pv.id,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		pv.title, pv.cover_image, pv.status, pv.updated_at,
		c.id as category_id, c.name as category_name, c.slug as category_slug
	FROM post_versions pv
	LEFT JOIN users u ON pv.created_by = u.id
	LEFT JOIN categories c ON pv.category_id = c.id
	WHERE pv.post_id = ?
	AND pv.deleted_at IS NULL;`
	QueryPostDetailsForVersionsGetByPostId = `
	SELECT
	p.current_version_id, p.created_at,
	u.id as author_id, u.name  as author_name, u.avatar as author_avatar
	FROM posts p
	LEFT JOIN users u
	ON u.id = p.created_by
	WHERE p.id = ?
	AND p.deleted_at IS NULL;`
	QueryGetPostVersionDuplicate = `
	SELECT
		id, post_id, title, slug, content, cover_image,
		description, spot, category_id, read_time, created_by
	FROM post_versions
	WHERE post_id = ?
	AND deleted_at IS NULL;`
	QueryGetSpecificVersionForDuplicate = `
	SELECT
		id, post_id, title, slug, content, cover_image,
		description, spot, category_id, read_time, created_by
	FROM post_versions
	WHERE id = ?
	AND deleted_at IS NULL;`
	QueryPostSoftDelete = `
	UPDATE posts
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE id = ?
	AND deleted_at IS NULL;`
	QueryPostAllRelatedCovers = `
	SELECT cover_image
	FROM post_versions
	WHERE post_id = ?;`
	QueryGetVersionCreatorAndStatus = `
	SELECT created_by, status
	FROM post_versions
	WHERE id = ?
	AND deleted_at IS NULL;`
	QueryPostVersionUpdate = `
	UPDATE post_versions
	SET
		title = COALESCE(?, title),
		slug = COALESCE(?, slug),
		content = COALESCE(?, content),
		cover_image = COALESCE(?, cover_image),
		description = COALESCE(?, description),
		spot = COALESCE(?, spot),
		category_id = COALESCE(?, category_id),
		read_time = COALESCE(?, read_time),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ?;`
	QueryPostVersionUpdateStatus = `
	UPDATE post_versions
	SET
		status = ?,
		status_changed_by = ?,
		status_changed_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ?;`
	QueryPostVersionUpdateStatusWithNote = `
	UPDATE post_versions
	SET
		status = ?,
		status_changed_by = ?,
		status_change_note = ?,
		status_changed_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ?;`
	QueryGetVersionCoverImage = `
	SELECT cover_image
	FROM post_versions
	WHERE id = ?
	AND deleted_at IS NULL;`
	QuerySoftDeleteVersion = `
	UPDATE post_versions
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE id = ?
	AND deleted_at IS NULL;`
	QueryCheckImageReferences = `
	SELECT COUNT(*)
	FROM post_versions
	WHERE cover_image = ?
	AND id != ?
	AND deleted_at IS NULL;`
	QueryCheckIfVersionIsCurrentlyPublished = `
	SELECT COUNT(*)
	FROM posts
	WHERE current_version_id = ?
	AND deleted_at IS NULL;`
	QuerySetPostCurrentVersionToNull = `
	UPDATE posts
	SET current_version_id = NULL
	WHERE current_version_id = ?
	AND deleted_at IS NULL;`
	QueryIncrementReadCount = `
	UPDATE posts
	SET read_count = read_count + 1
	WHERE id = ?
	AND deleted_at IS NULL;`
	QueryInsertPostView = `
	INSERT INTO post_views (post_id, user_agent)
	VALUES (?, ?);`
)
