package db

const (
	// USERS
	QueryCreateTableUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255) NOT NULL,
		avatar VARCHAR(100),
		passphrase_hash VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE,
		last_login TIMESTAMP WITH TIME ZONE,
		deleted_at TIMESTAMP WITH TIME ZONE,
		role_id INTEGER NOT NULL,
		FOREIGN KEY (role_id) REFERENCES roles(id)
		ON DELETE RESTRICT
	);
	CREATE UNIQUE INDEX IF NOT EXISTS unique_active_email
	ON users(email)
	WHERE deleted_at IS NULL;`
	// ROLE BASED ACCESS CONTROL
	QueryCreateTableRoles = `
	CREATE TABLE IF NOT EXISTS roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50) NOT NULL UNIQUE
	);`
	QueryCreateTablePermission = `
	CREATE TABLE IF NOT EXISTS permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50) NOT NULL UNIQUE
	);`
	QueryCreateTableRolePermissions = `
	CREATE TABLE IF NOT EXISTS role_permissions (
		role_id INTEGER NOT NULL,
		permission_id INTEGER NOT NULL,
		PRIMARY KEY (role_id, permission_id),
		FOREIGN KEY (role_id) REFERENCES roles(id),
		FOREIGN KEY (permission_id) REFERENCES permissions(id)
	);`
	// CATEGORIES
	QueryCreateTableCategories = `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(100) NOT NULL,
		slug VARCHAR(120) NOT NULL,
		spot VARCHAR(75) NOT NULL,
		description TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE,
		deleted_at TIMESTAMP WITH TIME ZONE
	);
	CREATE UNIQUE INDEX IF NOT EXISTS unique_active_slug_category
	ON categories(slug)
	WHERE deleted_at IS NULL;`
	// POSTS
	QueryCreateTablePosts = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_by INTEGER NOT NULL,
		current_version_id INTEGER NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP WITH TIME ZONE,
		FOREIGN KEY (created_by) REFERENCES users(id)
		ON DELETE RESTRICT,
		FOREIGN KEY (current_version_id) REFERENCES post_versions(id)
		ON DELETE SET NULL
	);
	CREATE INDEX IF NOT EXISTS idx_posts_deleted_at
	ON posts(deleted_at);`
	// POST VERSIONS
	QueryCreateTablePostVersions = `
	CREATE TABLE IF NOT EXISTS post_versions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
    duplicated_from INTEGER, -- If the initial version, must be null
		title VARCHAR(250),
		slug VARCHAR(250),
		content TEXT,
		cover_image VARCHAR(255),
		description VARCHAR(155),
		spot VARCHAR(75),
		category_id INTEGER,
		status INTEGER NOT NULL DEFAULT 0,
		created_by INTEGER NOT NULL,
		status_changed_at TIMESTAMP WITH TIME ZONE NULL,
		status_changed_by INTEGER NULL,
		status_change_note TEXT,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id)
		ON DELETE CASCADE,
		FOREIGN KEY (created_by) REFERENCES users(id)
		ON DELETE SET NULL,
		FOREIGN KEY (status_changed_by) REFERENCES users(id)
		ON DELETE SET NULL,
		FOREIGN KEY (category_id) REFERENCES categories(id)
		ON DELETE SET NULL
	);
	CREATE INDEX IF NOT EXISTS idx_post_versions_post_id_status
  ON post_versions(post_id, status);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_post_versions_slug_status
  ON post_versions(slug, status)
	WHERE status = 5; `
	// TAGS
	QueryCreateTableTags = `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50) NOT NULL UNIQUE,
		slug VARCHAR(100) NOT NULL UNIQUE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE,
		deleted_at TIMESTAMP WITH TIME ZONE
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_slug
	ON tags(slug) WHERE deleted_at IS NULL;`
	QueryCreateTablePostTags = `
	CREATE TABLE IF NOT EXISTS post_tags (
		post_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, tag_id),
		FOREIGN KEY (post_id) REFERENCES posts(id)
		ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id)
		ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_post_tags_tag_id
	ON post_tags(tag_id);
	CREATE INDEX IF NOT EXISTS idx_post_tags_post_id
	ON post_tags(post_id);`
	// VIEWS
	QueryCreateTableViews = `
	CREATE TABLE IF NOT EXISTS post_views (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		viewed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		user_agent TEXT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id)
		ON DELETE CASCADE
	);`
	// AUDIT LOGS
	QueryCreateTableAuditLogs = `
	CREATE TABLE IF NOT EXISTS audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NULL,
		entity TEXT NOT NULL,
		entity_id INTEGER NOT NULL,
		action TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
		ON DELETE SET NULL
	);`
)

var InitializeQueries = []string{
	QueryCreateTableUsersTable,
	QueryCreateTableRoles,
	QueryCreateTablePermission,
	QueryCreateTableRolePermissions,
	QueryCreateTableCategories,
	QueryCreateTablePosts,
	QueryCreateTablePostVersions,
	QueryCreateTableTags,
	QueryCreateTablePostTags,
	QueryCreateTableViews,
	QueryCreateTableAuditLogs,
}
