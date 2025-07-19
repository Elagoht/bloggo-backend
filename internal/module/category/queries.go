package category

const (
	QueryCategoryGetBySlug = `
	SELECT c.id, c.name, c.slug, c.description, c.created_at, c.updated_at,
	(
    SELECT COUNT(*)
    FROM posts p
    WHERE p.category_id = c.id AND p.deleted_at IS NULL
  ) AS blogCount
	FROM categories c
	WHERE c.slug = ? AND c.deleted_at IS NULL;`
	QueryCategoryGetCategories = `
	SELECT c.id, c.name, c.slug,
	(
    SELECT COUNT(*)
    FROM posts p
    WHERE p.category_id = c.id AND p.deleted_at IS NULL
  ) AS blogCount
	FROM categories c
	WHERE c.deleted_at IS NULL;`
	QueryCategoryCreate = `
	INSERT INTO categories (
		name,
		slug,
		description
	) VALUES (?, ?, ?);`
	QueryCategoryPatch = `
	UPDATE categories
	SET
		name = COALESCE(?, name),
		description = COALESCE(?, description),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ? AND deleted_at IS NULL;`
	QueryCategorySoftDelete = `
	UPDATE categories
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ? AND deleted_at IS NULL;`
)
