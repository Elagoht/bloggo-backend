package category

const (
	QueryCategoryGetBySlug = `
	SELECT id, name, slug, description, created_at, updated_at
	FROM categories
	WHERE slug = ? AND deleted_at IS NULL;`
	QueryCategoryGetBlogsCards = `
	SELECT `
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
