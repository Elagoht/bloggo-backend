package tag

const (
	QueryTagGetBySlug = `
	SELECT t.id, t.name, t.slug, t.created_at, t.updated_at,
  (
    SELECT COUNT(*)
    FROM post_tags pt
    JOIN posts p ON p.id = pt.post_id
    WHERE pt.tag_id = t.id AND p.deleted_at IS NULL
  ) AS blogCount
	FROM tags t
	WHERE t.slug = ? AND t.deleted_at IS NULL;`
	QueryTagGetCategories = `
	SELECT t.id, t.name, t.slug,
	(
    SELECT COUNT(*)
    FROM post_tags pt
    JOIN posts p ON p.id = pt.post_id
    WHERE pt.tag_id = t.id AND p.deleted_at IS NULL
  ) AS blogCount
	FROM tags t
	WHERE t.deleted_at IS NULL%s;`
	QueryTagCreate = `
	INSERT INTO tags (
		name,
		slug
	) VALUES (?, ?);`
	QueryTagPatch = `
	UPDATE tags
	SET
		name = COALESCE(?, name),
		slug = COALESCE(?, slug),
		updated_at = CURRENT_TIMESTAMP
	WHERE slug = ? AND deleted_at IS NULL;`
	QueryTagSoftDelete = `
	UPDATE tags
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE slug = ? AND deleted_at IS NULL;`
	// Post-Tag Relationships
	QueryGetPostTags = `
	SELECT t.id, t.name, t.slug,
	(
		SELECT COUNT(*)
		FROM post_tags pt2
		JOIN posts p2 ON p2.id = pt2.post_id
		WHERE pt2.tag_id = t.id AND p2.deleted_at IS NULL
	) AS blogCount
	FROM tags t
	JOIN post_tags pt ON pt.tag_id = t.id
	WHERE pt.post_id = ? AND t.deleted_at IS NULL;`
	QueryAssignTagToPost = `
	INSERT INTO post_tags (post_id, tag_id)
	VALUES (?, ?)
	ON CONFLICT (post_id, tag_id) DO NOTHING;`
	QueryRemoveTagFromPost = `
	DELETE FROM post_tags
	WHERE post_id = ? AND tag_id = ?;`
	QueryRemoveAllTagsFromPost = `
	DELETE FROM post_tags
	WHERE post_id = ?;`
	QueryCheckPostExists = `
	SELECT COUNT(*)
	FROM posts
	WHERE id = ? AND deleted_at IS NULL;`
	QueryCheckTagExists = `
	SELECT COUNT(*)
	FROM tags
	WHERE id = ? AND deleted_at IS NULL;`
)
