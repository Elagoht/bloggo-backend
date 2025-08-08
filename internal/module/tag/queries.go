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
)
