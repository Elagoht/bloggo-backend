package db

const (
	SeedRolesSQL = `
	INSERT INTO roles (name) VALUES
  ('author'),
  ('editor'),
  ('admin')
	ON CONFLICT(name) DO NOTHING;
`
	SeedPermissionsSQL = `
	INSERT INTO permissions (name) VALUES
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
  ('user:change_password'),
  ('stats:view'),
  ('role:assign'),
  ('auditlog:view'),
  ('schedule:create'),
  ('schedule:update'),
  ('schedule:delete'),
  ('schedule:view')
	ON CONFLICT(name) DO NOTHING;`
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
)

var (
	RolePermissionsMatrix = map[string][]string{
		"admin": {
			"post:create", "post:edit", "post:edit_own", "post:delete", "post:delete_own", "post:publish", "post:view", "post:schedule", "tag:manage", "category:manage", "user:view", "user:create", "user:update", "user:delete", "user:change_password", "stats:view", "role:assign", "auditlog:view", "schedule:create", "schedule:update", "schedule:delete", "schedule:view",
		},
		"editor": {
			"post:create", "post:edit", "post:edit_own", "post:delete", "post:delete_own", "post:publish", "post:view", "post:schedule", "stats:view", "schedule:create", "schedule:update", "schedule:delete", "schedule:view",
		},
		"author": {
			"post:create", "post:edit_own", "post:delete_own", "post:view",
		},
	}
	SeedQueries = []string{
		SeedRolesSQL,
		SeedPermissionsSQL,
	}
)
