package category

const (
	QueryCategoryGetBySlug = `
	SELECT c.id, c.name, c.slug, c.spot, c.description, c.created_at, c.updated_at,
	(
    SELECT COUNT(*)
    FROM posts p
    WHERE p.category_id = c.id AND p.deleted_at IS NULL
  ) AS blogCount
	FROM categories c
	WHERE c.slug = ? AND c.deleted_at IS NULL;`
	QueryCategoryGetCategories = `
	SELECT c.id, c.name, c.slug, c.spot,
	(
    SELECT COUNT(*)
    FROM posts p
    WHERE p.category_id = c.id AND p.deleted_at IS NULL
  ) AS blogCount
	FROM categories c
	WHERE c.deleted_at IS NULL%s;`
	QueryCategoryCreate = `
	INSERT INTO categories (
		name,
		slug,
		spot,
		description
	) VALUES (?, ?, ?, ?);`
	QueryCategoryPatch = `
	UPDATE categories
	SET
		name = COALESCE(?, name),
		slug = COALESCE(?, slug),
		spot = COALESCE(?, spot),
		description = COALESCE(?, description),
		updated_at = CURRENT_TIMESTAMP
	WHERE slug = ? AND deleted_at IS NULL;`
	QueryCategorySoftDelete = `
	UPDATE categories
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE slug = ? AND deleted_at IS NULL;`
)
