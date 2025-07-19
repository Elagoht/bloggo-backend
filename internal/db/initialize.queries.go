package db

const (
	// USERS
	CreateUsersTable = `CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL,
  avatar VARCHAR(100),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  last_login TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE,
  role_id INT NOT NULL,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);`
	CreateEmailUnique = `CREATE UNIQUE INDEX IF NOT EXISTS unique_active_email
  ON users(email)
  WHERE deleted_at IS NULL;`
	// ROLE BASED ACCESS CONTROL
	CreateRoles = `CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE
);`
	CreatePermission = `CREATE TABLE IF NOT EXISTS permissions (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE
);`
	CreateRolePermissions = `CREATE TABLE IF NOT EXISTS role_permissions (
  role_id INT NOT NULL,
  permission_id INT NOT NULL,
  PRIMARY KEY (role_id, permission_id),
  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (permission_id) REFERENCES permissions(id)
);`
	// CATEGORIES
	CreateCategories = `CREATE TABLE IF NOT EXISTS categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(120) NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);`
	CreateCategorySlugUnique = `CREATE UNIQUE INDEX IF NOT EXISTS unique_active_slug_category
  ON categories(slug)
  WHERE deleted_at IS NULL;`
	// POST AND VERSIONS
	CreatePosts = `CREATE TABLE IF NOT EXISTS posts (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  category_id INT NULL,
  slug VARCHAR(150) NOT NULL,
  published_version_id INT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);`
	CreatePostSlugUnique = `CREATE UNIQUE INDEX IF NOT EXISTS unique_active_slug_post
  ON posts(slug)
  WHERE deleted_at IS NULL;`
	CreatePostVersions = `CREATE TABLE IF NOT EXISTS post_versions (
  id SERIAL PRIMARY KEY,
  post_id INT NOT NULL,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  content TEXT NOT NULL,
  cover_image VARCHAR(255) NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'draft',
  created_by INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);`
	// TAGS
	CreateTags = `CREATE TABLE IF NOT EXISTS tags (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  slug VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);`
	CreatePostTags = `CREATE TABLE IF NOT EXISTS post_tags (
  post_id INT NOT NULL,
  tag_id INT NOT NULL,
  PRIMARY KEY (post_id, tag_id),
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);`
	// VIEWS
	CreateViews = `CREATE TABLE IF NOT EXISTS post_views (
  id SERIAL PRIMARY KEY,
  post_id INT NOT NULL,
  viewed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_agent TEXT NULL,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);`
	// SCHEDULES
	CreateSchedules = `CREATE TABLE IF NOT EXISTS schedules (
  id SERIAL PRIMARY KEY,
  post_id INT NOT NULL,
  scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  executed_at TIMESTAMP WITH TIME ZONE NULL,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);`
	// AUDIT LOGS
	CreateAuditLogs = `CREATE TABLE IF NOT EXISTS audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  entity TEXT NOT NULL,
  entity_id INTEGER NOT NULL,
  action TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);`
)

var InitializeQueries = []string{
	CreateUsersTable,
	CreateEmailUnique,
	CreateRoles,
	CreatePermission,
	CreateRolePermissions,
	CreateCategories,
	CreateCategorySlugUnique,
	CreatePosts,
	CreatePostSlugUnique,
	CreatePostVersions,
	CreateTags,
	CreatePostTags,
	CreateViews,
	CreateSchedules,
	CreateAuditLogs,
}
