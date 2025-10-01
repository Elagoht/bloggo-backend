package tags

const (
	// Tags queries
	QueryAPIGetAllTags = `
	SELECT
		t.slug,
		t.name,
		COUNT(DISTINCT pt.post_id) as post_count
	FROM tags t
	LEFT JOIN post_tags pt ON pt.tag_id = t.id
	LEFT JOIN posts p ON p.id = pt.post_id
		AND p.deleted_at IS NULL
	LEFT JOIN post_versions pv ON pv.id = p.current_version_id
		AND pv.deleted_at IS NULL
		AND pv.status = 5
	WHERE t.deleted_at IS NULL
	GROUP BY t.id, t.slug, t.name
	ORDER BY t.name ASC;`

	QueryAPIGetTagBySlug = `
	SELECT
		t.slug,
		t.name,
		COUNT(DISTINCT pt.post_id) as post_count
	FROM tags t
	LEFT JOIN post_tags pt ON pt.tag_id = t.id
	LEFT JOIN posts p ON p.id = pt.post_id
		AND p.deleted_at IS NULL
	LEFT JOIN post_versions pv ON pv.id = p.current_version_id
		AND pv.deleted_at IS NULL
		AND pv.status = 5
	WHERE t.deleted_at IS NULL
		AND t.slug = ?
	GROUP BY t.id, t.slug, t.name
	LIMIT 1;`
)
