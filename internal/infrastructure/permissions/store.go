package permissions

import (
	"database/sql"
)

type Store interface {
	Load(db *sql.DB) error
	HasPermission(role string, permission string) bool
}

type permissionCell = map[string]bool
type permissionStore = map[string]permissionCell
