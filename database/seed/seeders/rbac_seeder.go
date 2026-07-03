package seeders

import (
	"database/sql"
	"fmt"

	"go-fiber-clean-arch-auth-rbac/internal/config"
	"go-fiber-clean-arch-auth-rbac/internal/router"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type seedRole struct {
	name        string
	description string
}

var roles = []seedRole{
	{name: "super_admin", description: "Full access to all routes"},
	{name: "admin", description: "Administrative access to managed routes"},
	{name: "customer", description: "Default customer access"},
}

func SeedRBAC(cfg *config.Config) error {
	db, err := sql.Open("pgx", cfg.DBUrl)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	roleIDs := make(map[string]int64)
	for _, role := range roles {
		id, err := upsertRole(db, role)
		if err != nil {
			return err
		}
		roleIDs[role.name] = id
	}

	permissionIDs := make(map[string]int64)
	for _, permission := range router.PermissionDefinitions() {
		id, err := upsertPermission(db, permission)
		if err != nil {
			return err
		}
		permissionIDs[permission.Name] = id
	}

	if err := syncRolePermissions(db, roleIDs, permissionIDs); err != nil {
		return err
	}

	return nil
}

func upsertRole(db *sql.DB, role seedRole) (int64, error) {
	var id int64
	err := db.QueryRow(`
		INSERT INTO roles (name, description)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET
			description = EXCLUDED.description,
			updated_at = NOW()
		RETURNING id
	`, role.name, role.description).Scan(&id)
	return id, err
}

func upsertPermission(db *sql.DB, permission router.PermissionDefinition) (int64, error) {
	var id int64
	err := db.QueryRow(`
		INSERT INTO permissions (name, method, path, description)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO UPDATE SET
			method = EXCLUDED.method,
			path = EXCLUDED.path,
			description = EXCLUDED.description,
			updated_at = NOW()
		RETURNING id
	`, permission.Name, permission.Method, permission.Path, permission.Description).Scan(&id)
	return id, err
}

func syncRolePermissions(db *sql.DB, roleIDs map[string]int64, permissionIDs map[string]int64) error {
	adminPermissions := []string{
		"user:me",
		"log:list",
		"role:list",
		"role:create",
		"role:update",
		"role:delete",
		"permission:list",
		"permission:create",
		"permission:update",
		"permission:delete",
		"role_permission:assign",
		"role_permission:revoke",
		"user:assign_role",
	}
	customerPermissions := []string{
		"user:me",
	}

	if err := assignAll(db, roleIDs["super_admin"], permissionIDs); err != nil {
		return err
	}
	if err := assignSubset(db, roleIDs["admin"], permissionIDs, adminPermissions); err != nil {
		return err
	}
	if err := assignSubset(db, roleIDs["customer"], permissionIDs, customerPermissions); err != nil {
		return err
	}

	return nil
}

func assignAll(db *sql.DB, roleID int64, permissionIDs map[string]int64) error {
	for _, permissionID := range permissionIDs {
		if _, err := db.Exec(`
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES ($1, $2)
			ON CONFLICT (role_id, permission_id) DO NOTHING
		`, roleID, permissionID); err != nil {
			return err
		}
	}
	return nil
}

func assignSubset(db *sql.DB, roleID int64, permissionIDs map[string]int64, subset []string) error {
	for _, permissionName := range subset {
		permissionID, ok := permissionIDs[permissionName]
		if !ok {
			return fmt.Errorf("permission %s not found", permissionName)
		}
		if _, err := db.Exec(`
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES ($1, $2)
			ON CONFLICT (role_id, permission_id) DO NOTHING
		`, roleID, permissionID); err != nil {
			return err
		}
	}
	return nil
}
