package db

import (
	"bloggo/internal/utils/cryptography"
	"database/sql"
	"fmt"
	"log"
)

func SeedDatabase(database *sql.DB) error {
	// -- Roles and Permissions Seeding -- //
	for _, query := range SeedQueries {
		_, err := database.Exec(query)
		if err != nil {
			return fmt.Errorf("database cannot be seeded: %w", err)
		}
	}

	// -- Role-Permission Mapping Seeding -- //
	// This process requires to get generated id fields

	// Query all roles and permissions into maps
	roleIDs := map[string]int64{}
	permissionIDs := map[string]int64{}

	rows, err := database.Query(GetRolesSQL)
	if err != nil {
		return fmt.Errorf("cannot query roles: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err == nil {
			roleIDs[name] = id
		}
	}

	rows, err = database.Query(GetPermissionsSQL)
	if err != nil {
		return fmt.Errorf("cannot query permissions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err == nil {
			permissionIDs[name] = id
		}
	}

	// Insert role-permission assignments
	for role, permissions := range RolePermissionsMatrix {
		roleID, ok := roleIDs[role]
		if !ok {
			log.Printf("Role not found: %s", role)
			continue
		}
		for _, permission := range permissions {
			permissionID, ok := permissionIDs[permission]
			if !ok {
				log.Printf("Permission not found: %s", permission)
				continue
			}
			_, err := database.Exec(
				InsertPermissionToRoleSQL,
				roleID,
				permissionID,
			)
			if err != nil {
				log.Printf("Failed to insert role_permission: %s-%s", role, permission)
			}
		}
	}

	// -- Admin User Seeding -- //
	adminRoleID, ok := roleIDs["Admin"]
	if !ok {
		log.Printf("Admin role not found, cannot create default Admin user.")
	} else {
		hashedPass, err := cryptography.HashPassphrase(DefaultAdminPassphrase)
		if err != nil {
			log.Printf("Failed to hash Admin passphrase: %v", err)
		} else {
			_, err := database.Exec(
				InsertUserSQL,
				DefaultAdminName,
				DefaultAdminEmail,
				hashedPass,
				adminRoleID,
			)
			if err != nil {
				log.Printf("Failed to insert default Admin user: %v", err)
			}
		}
	}

	return nil
}
