package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/database"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/config"
)

func main() {
	config.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get database handle", "error", err)
		return
	}
	defer sqlDB.Close()

	row := sqlDB.QueryRowContext(context.Background(), "SELECT id, email, password_hash FROM users WHERE email = $1", "admin@example.com")
	var id int
	var email string
	var hash string
	if err := row.Scan(&id, &email, &hash); err != nil {
		slog.Error("Failed to scan row", "error", err)
		return
	}

	fmt.Printf("id=%d email=%s hash=%s\n", id, email, hash)
}
