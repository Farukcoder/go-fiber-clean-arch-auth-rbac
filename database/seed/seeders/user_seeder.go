package seeders

import (
	"database/sql"
	"fmt"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	name     string
	email    string
	password string
	roleName string
}

var Users = []seedUser{
	{name: "Super Admin", email: "superadmin@gmail.com", password: "Password123!", roleName: "super_admin"},
	{name: "Admin User", email: "admin@example.com", password: "Password123!", roleName: "admin"},
	{name: "Demo User", email: "demo@example.com", password: "DemoPassword1", roleName: "customer"},
}

func SeedUsers(cfg *config.Config) error {
	db, err := sql.Open("pgx", cfg.DBUrl)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	for _, u := range Users {
		if err := seedUserRow(db, u); err != nil {
			return fmt.Errorf("failed to seed user %s: %w", u.email, err)
		}
	}

	return nil
}

func seedUserRow(db *sql.DB, user seedUser) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	roleName := user.roleName
	if roleName == "" {
		roleName = "customer"
	}

	var roleID int64
	if err := db.QueryRow(`SELECT id FROM roles WHERE name = $1`, roleName).Scan(&roleID); err != nil {
		return fmt.Errorf("role %s not found: %w", roleName, err)
	}

	var userID int64
	err = db.QueryRow(`
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET
		name = EXCLUDED.name,
		password_hash = EXCLUDED.password_hash
		RETURNING id
	`, user.name, user.email, string(hash)).Scan(&userID)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET
		role_id = EXCLUDED.role_id,
		updated_at = NOW()
	`, userID, roleID)
	return err
}
