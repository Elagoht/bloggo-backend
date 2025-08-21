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
    ('post:edit'),
    ('post:edit_own'),
    ('post:delete'),
    ('post:delete_own'),
    ('post:publish'),
    ('post:view'),
    ('post:schedule'),
    ('tag:manage'),
    ('category:manage'),
    ('user:view'),
    ('user:create'),
    ('user:update'),
    ('user:delete'),
    ('user:change_passphrase'),
    ('statistics:view-all'),
    ('statistics:view-self'),
    ('role:assign'),
    ('auditlog:view'),
    ('schedule:create'),
    ('schedule:update'),
    ('schedule:delete'),
    ('schedule:view')
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
		"Admin": {
			"post:create", "post:edit", "post:edit_own", "post:delete", "post:delete_own", "post:publish", "post:view", "post:schedule", "tag:manage", "category:manage", "user:view", "user:create", "user:update", "user:delete", "user:change_passphrase", "statistics:view-all", "statistics:view-self", "role:assign", "auditlog:view", "schedule:create", "schedule:update", "schedule:delete", "schedule:view",
		},
		"Editor": {
			"post:create", "post:edit", "post:edit_own", "post:delete", "post:delete_own", "post:publish", "post:view", "post:schedule", "statistics:view-all", "statistics:view-self", "schedule:create", "schedule:update", "schedule:delete", "schedule:view",
		},
		"Author": {
			"post:create", "post:edit_own", "post:delete_own", "post:view", "statistics:view-all", "statistics:view-self",
		},
	}
	SeedQueries = []string{
		SeedRolesSQL,
		SeedPermissionsSQL,
	}
)
