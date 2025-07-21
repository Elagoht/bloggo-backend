package permissionstore

import (
	"database/sql"
	"sync"
)

type memoryStore struct {
	permissions permissionStore
	lock        sync.RWMutex
}

var (
	once     sync.Once
	instance PermissionStore
)

func GetStore() PermissionStore {
	once.Do(func() {
		instance = newMemoryStore()
	})
	return instance
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		permissions: make(permissionStore),
	}
}

// Loads the role-permission mapping from the database into memory.
func (store *memoryStore) Load(db *sql.DB) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	rows, err := db.Query(QueryGetPermissionRoles)
	if err != nil {
		return err
	}
	defer rows.Close()

	store.permissions = make(permissionStore)

	for rows.Next() {
		var role, permission string
		if err := rows.Scan(&role, &permission); err != nil {
			return err
		}
		if _, ok := store.permissions[role]; !ok {
			store.permissions[role] = make(permissionCell)
		}
		store.permissions[role][permission] = true
	}
	return nil
}

// Checks if the given role has the specified permission.
func (store *memoryStore) HasPermission(
	role string,
	permission string,
) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()
	perms, ok := store.permissions[role]
	if !ok {
		return false
	}
	return perms[permission]
}
