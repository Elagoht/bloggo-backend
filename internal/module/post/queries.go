package post

const (
	QueryPostGetByCurrentVersionSlug = `
	SELECT
		p.id as post_id,
		COALESCE(current_pv.id, best_pv.id) as version_id,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		COALESCE(current_pv.title, best_pv.title) as title,
		COALESCE(current_pv.slug, best_pv.slug) as slug,
		COALESCE(current_pv.content, best_pv.content) as content,
		COALESCE(current_pv.cover_image, best_pv.cover_image) as cover_image,
		COALESCE(current_pv.description, best_pv.description) as description,
		COALESCE(current_pv.spot, best_pv.spot) as spot,
		COALESCE(current_pv.status, best_pv.status) as status,
		p.read_count,
		COALESCE(current_pv.created_at, best_pv.created_at) as created_at,
		COALESCE(current_pv.updated_at, best_pv.updated_at) as updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	LEFT JOIN post_versions current_pv ON current_pv.id = p.current_version_id AND current_pv.deleted_at IS NULL
	LEFT JOIN post_versions best_pv ON best_pv.post_id = p.id
		AND best_pv.deleted_at IS NULL
		AND best_pv.id = (
			SELECT pv2.id FROM post_versions pv2
			WHERE pv2.post_id = p.id AND pv2.deleted_at IS NULL
			ORDER BY
			CASE pv2.status
					WHEN 5 THEN 1  -- Published
					WHEN 2 THEN 2  -- Approved
					WHEN 1 THEN 3  -- Pending
					WHEN 0 THEN 4  -- Draft
					WHEN 3 THEN 5  -- Rejected
					ELSE 6
				END,
				pv2.updated_at DESC
			LIMIT 1
		)
	LEFT JOIN categories c ON c.id = COALESCE(current_pv.category_id, best_pv.category_id)
	LEFT JOIN users u ON u.id = p.created_by
	WHERE (current_pv.slug = ? OR best_pv.slug = ?)
		AND p.deleted_at IS NULL
		AND (current_pv.id IS NOT NULL OR best_pv.id IS NOT NULL)
	LIMIT 1;`
	QueryPostGetById = `
	SELECT
		p.id as post_id,
		COALESCE(current_pv.id, best_pv.id) as version_id,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		COALESCE(current_pv.title, best_pv.title) as title,
		COALESCE(current_pv.slug, best_pv.slug) as slug,
		COALESCE(current_pv.content, best_pv.content) as content,
		COALESCE(current_pv.cover_image, best_pv.cover_image) as cover_image,
		COALESCE(current_pv.description, best_pv.description) as description,
		COALESCE(current_pv.spot, best_pv.spot) as spot,
		COALESCE(current_pv.status, best_pv.status) as status,
		p.read_count,
		COALESCE(current_pv.created_at, best_pv.created_at) as created_at,
		COALESCE(current_pv.updated_at, best_pv.updated_at) as updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	LEFT JOIN post_versions current_pv ON current_pv.id = p.current_version_id AND current_pv.deleted_at IS NULL
	LEFT JOIN post_versions best_pv ON best_pv.post_id = p.id
		AND best_pv.deleted_at IS NULL
		AND best_pv.id = (
			SELECT pv2.id FROM post_versions pv2
			WHERE pv2.post_id = p.id AND pv2.deleted_at IS NULL
			ORDER BY
			CASE pv2.status
					WHEN 5 THEN 1  -- Published
					WHEN 2 THEN 2  -- Approved
					WHEN 1 THEN 3  -- Pending
					WHEN 0 THEN 4  -- Draft
					WHEN 3 THEN 5  -- Rejected
					ELSE 6
				END,
				pv2.updated_at DESC
			LIMIT 1
		)
	LEFT JOIN categories c ON c.id = COALESCE(current_pv.category_id, best_pv.category_id)
	LEFT JOIN users u ON u.id = p.created_by
	WHERE p.id = ?
		AND p.deleted_at IS NULL
		AND (current_pv.id IS NOT NULL OR best_pv.id IS NOT NULL)
	LIMIT 1;`
	QueryPostGetList = `
	SELECT
		p.id as post_id,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		COALESCE(current_pv.title, best_pv.title) as title,
		COALESCE(current_pv.slug, best_pv.slug) as slug,
		COALESCE(current_pv.cover_image, best_pv.cover_image) as cover_image,
		COALESCE(current_pv.spot, best_pv.spot) as spot,
		COALESCE(current_pv.status, best_pv.status) as status,
		p.read_count,
		COALESCE(current_pv.created_at, best_pv.created_at) as created_at,
		COALESCE(current_pv.updated_at, best_pv.updated_at) as updated_at,
		c.slug AS category_slug, c.id AS category_id, c.name AS category_name
	FROM posts p
	LEFT JOIN post_versions current_pv ON current_pv.id = p.current_version_id AND current_pv.deleted_at IS NULL
	LEFT JOIN post_versions best_pv ON best_pv.post_id = p.id
		AND best_pv.deleted_at IS NULL
		AND best_pv.id = (
			SELECT pv2.id FROM post_versions pv2
			WHERE pv2.post_id = p.id AND pv2.deleted_at IS NULL
			ORDER BY
			CASE pv2.status
					WHEN 5 THEN 1  -- Published
					WHEN 2 THEN 2  -- Approved
					WHEN 1 THEN 3  -- Pending
					WHEN 0 THEN 4  -- Draft
					WHEN 3 THEN 5  -- Rejected
					ELSE 6
				END,
				pv2.updated_at DESC
			LIMIT 1
		)
	LEFT JOIN categories c ON c.id = COALESCE(current_pv.category_id, best_pv.category_id)
	LEFT JOIN users u ON u.id = p.created_by
	WHERE p.deleted_at IS NULL
	AND (current_pv.id IS NOT NULL OR best_pv.id IS NOT NULL) %s;`
	QueryPostVersionGetById = `
	SELECT
		pv.id, pv.duplicated_from,
		u.id as author_id, u.name as author_name, u.avatar as author_avatar,
		pv.title, pv.slug, pv.content, pv.cover_image, pv.description, pv.spot,
		pv.status, pv.status_changed_at, pv.status_change_note,
		pv.created_at, pv.updated_at,
		c.id AS category_id, c.name AS category_name, c.slug AS category_slug,
		scb.id as status_changed_by_id, scb.name as status_changed_by_name, scb.avatar as status_changed_by_avatar
	FROM posts p
	JOIN post_versions pv
	ON pv.post_id = p.id
	LEFT JOIN categories c
	ON c.id = pv.category_id
	LEFT JOIN users u
	ON u.id = pv.created_by
	LEFT JOIN users scb
	ON scb.id = pv.status_changed_by
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
	QueryGetPublishedVersionBySlug = `
	SELECT id, post_id
	FROM post_versions
	WHERE slug = ? AND status = 5 AND deleted_at IS NULL`
	QueryUnpublishVersionBySlug = `
	UPDATE post_versions
	SET status = 2, updated_at = CURRENT_TIMESTAMP
	WHERE slug = ? AND status = 5 AND deleted_at IS NULL`
	QueryGetVersionSlug = `
	SELECT slug FROM post_versions WHERE id = ? AND deleted_at IS NULL`
	QueryIncrementReadCount = `
	UPDATE posts
	SET read_count = read_count + 1
	WHERE id = ?
	AND deleted_at IS NULL;`
	QueryInsertPostView = `
	INSERT INTO post_views (post_id, user_agent)
	VALUES (?, ?);`
	QueryPostGetListCount = `
	SELECT COUNT(*)
	FROM posts p
	LEFT JOIN post_versions current_pv ON current_pv.id = p.current_version_id AND current_pv.deleted_at IS NULL
	LEFT JOIN post_versions best_pv ON best_pv.post_id = p.id
		AND best_pv.deleted_at IS NULL
		AND best_pv.id = (
			SELECT pv2.id FROM post_versions pv2
			WHERE pv2.post_id = p.id AND pv2.deleted_at IS NULL
			ORDER BY
			CASE pv2.status
					WHEN 5 THEN 1  -- Published
					WHEN 2 THEN 2  -- Approved
					WHEN 1 THEN 3  -- Pending
					WHEN 0 THEN 4  -- Draft
					WHEN 3 THEN 5  -- Rejected
					ELSE 6
				END,
				pv2.updated_at DESC
			LIMIT 1
		)
	LEFT JOIN categories c ON c.id = COALESCE(current_pv.category_id, best_pv.category_id)
	LEFT JOIN users u ON u.id = p.created_by
	WHERE p.deleted_at IS NULL
	AND (current_pv.id IS NOT NULL OR best_pv.id IS NOT NULL)%s;`
	// Post-Tag Relationships
	QueryGetPostTags = `
	SELECT t.id, t.name, t.slug
	FROM tags t
	JOIN post_tags pt ON pt.tag_id = t.id
	WHERE pt.post_id = ? AND t.deleted_at IS NULL;`
	QueryAssignTagsToPost = `
	INSERT INTO post_tags (post_id, tag_id)
	VALUES %s
	ON CONFLICT (post_id, tag_id) DO NOTHING;`
	QueryRemoveTagsFromPost = `
	DELETE FROM post_tags
	WHERE post_id = ? AND tag_id IN (%s);`
	QueryGetCurrentPostTagIds = `
	SELECT tag_id
	FROM post_tags
	WHERE post_id = ?;`
	QueryCheckTagExists = `
	SELECT COUNT(*)
	FROM tags
	WHERE id = ? AND deleted_at IS NULL;`
)
