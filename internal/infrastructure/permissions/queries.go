package permissions

const (
	QueryGetPermissionRoles = `
	SELECT rp.role_id, p.name AS permission_name
	FROM role_permissions rp
	JOIN permissions p ON rp.permission_id = p.id
	JOIN roles r ON rp.role_id = r.id;`
)
