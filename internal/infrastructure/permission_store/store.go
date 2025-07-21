package permissionstore

import (
	"database/sql"
)

type PermissionStore interface {
	Load(db *sql.DB) error
	HasPermission(role string, permission string) bool
}

type permissionCell = map[string]bool
type permissionStore = map[string]permissionCell
