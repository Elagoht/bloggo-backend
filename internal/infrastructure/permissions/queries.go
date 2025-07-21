package permissions

const (
	QueryGetPermissionRoles = `
	SELECT r.name, p.name
	FROM roles r
	JOIN role_permissions rp
	ON r.id = rp.role_id
	JOIN permissions p
	ON p.id = rp.permission_id;`
)
