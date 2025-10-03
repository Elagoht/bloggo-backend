package db

const (
	SeedRolesSQL = `
	INSERT INTO roles (name)
  VALUES
    ('Author'),
    ('Editor'),
    ('Admin')
	ON CONFLICT(name) DO NOTHING;
`
	SeedPermissionsSQL = `
	INSERT INTO permissions (name)
  VALUES
    ('post:create'),
    ('post:delete'),
		('post:publish'),
    ('post:view'),
    ('post:list'),
    ('tag:list'),
    ('tag:view'),
    ('tag:create'),
    ('tag:update'),
    ('tag:delete'),
		('tag:assign'),
    ('category:list'),
    ('category:view'),
    ('category:create'),
    ('category:update'),
    ('category:delete'),
    ('user:list'),
    ('user:view'),
    ('user:register'),
    ('user:update'),
    ('user:delete'),
    ('user:change_passphrase'),
    ('user:assign_role'),
    ('statistics:view-total'),
    ('statistics:view-others'),
    ('statistics:view-self'),
    ('auditlog:view'),
    ('keyvalue:manage'),
    ('webhook:manage'),
    ('apidoc:view')
	ON CONFLICT(name)
  DO NOTHING;`
	InsertPermissionToRoleSQL = `
  INSERT INTO role_permissions (role_id, permission_id)
  VALUES (?, ?)
  ON CONFLICT(role_id, permission_id)
  DO NOTHING;`
	GetPermissionsSQL = `
  SELECT id, name
  FROM permissions;`
	GetRolesSQL = `
  SELECT id, name
  FROM roles;`
	// Default Admin user seed
	DefaultAdminName       = "Root User"
	DefaultAdminEmail      = "root@admin.dev"
	DefaultAdminPassphrase = "ChangeMeNow123!"
	InsertUserSQL          = `
	INSERT INTO users (name, email, passphrase_hash, role_id)
	VALUES (?, ?, ?, ?)
	ON CONFLICT DO NOTHING;`
)

var (
	RolePermissionsMatrix = map[string][]string{
		"Author": {
			"post:create", "post:view", "post:list",
			"tag:list", "tag:view",
			"category:list", "category:view",
			"statistics:view-self",
		},
		"Editor": {
			"post:create", "post:delete", "post:publish", "post:view", "post:list",
			"tag:list", "tag:view", "tag:create", "tag:update", "tag:delete", "tag:assign",
			"category:list", "category:view", "category:create", "category:update", "category:delete",
			"user:list", "user:view",
			"statistics:view-self", "statistics:view-others",
			"keyvalue:manage", "webhook:manage",
		},
		"Admin": {
			"post:create", "post:delete", "post:publish", "post:view", "post:list",
			"tag:list", "tag:view", "tag:create", "tag:update", "tag:delete", "tag:assign",
			"category:list", "category:view", "category:create", "category:update", "category:delete",
			"user:list", "user:view", "user:register", "user:update", "user:delete", "user:change_passphrase", "user:assign_role",
			"statistics:view-self", "statistics:view-others", "statistics:view-total",
			"keyvalue:manage", "auditlog:view", "webhook:manage", "apidoc:view",
		},
	}
	SeedQueries = []string{
		SeedRolesSQL,
		SeedPermissionsSQL,
	}
)
