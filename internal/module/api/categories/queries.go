package categories

const (
	// Categories queries
	QueryAPIGetAllCategories = `
	SELECT
		c.slug,
		c.name,
		c.description,
		c.spot,
		COUNT(DISTINCT p.id) as post_count
	FROM categories c
	LEFT JOIN post_versions pv ON pv.category_id = c.id
		AND pv.deleted_at IS NULL
		AND pv.status = 5
	LEFT JOIN posts p ON p.current_version_id = pv.id
		AND p.deleted_at IS NULL
	WHERE c.deleted_at IS NULL
	GROUP BY c.id, c.slug, c.name, c.description, c.spot
	ORDER BY c.name ASC;`

	QueryAPIGetCategoryBySlug = `
	SELECT
		c.slug,
		c.name,
		c.description,
		c.spot,
		COUNT(DISTINCT p.id) as post_count
	FROM categories c
	LEFT JOIN post_versions pv ON pv.category_id = c.id
		AND pv.deleted_at IS NULL
		AND pv.status = 5
	LEFT JOIN posts p ON p.current_version_id = pv.id
		AND p.deleted_at IS NULL
	WHERE c.deleted_at IS NULL
		AND c.slug = ?
	GROUP BY c.id, c.slug, c.name, c.description, c.spot
	LIMIT 1;`
)
