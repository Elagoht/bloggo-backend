package permissions

import (
	"database/sql"
)

type Store interface {
	Load(db *sql.DB) error
	HasPermission(role int64, permission string) bool
}

type permissionCell = map[string]bool
type permissionStore = map[int64]permissionCell
