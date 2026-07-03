package main

import (
	"context"
	"fmt"
	"log/slog"

	"go-fiber-clean-arch-auth-rbac/database"
	"go-fiber-clean-arch-auth-rbac/internal/config"
	"go-fiber-clean-arch-auth-rbac/internal/dto"
	"go-fiber-clean-arch-auth-rbac/internal/repository"
	"go-fiber-clean-arch-auth-rbac/internal/service"
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

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	rbacRepo := repository.NewRBACRepository(db)
	authService := service.NewAuthService(userRepo, rbacRepo, refreshTokenRepo, cfg.JwtSecret, cfg.JwtRefreshSecret)

	resp, err := authService.Login(context.Background(), dto.LoginInput{
		EmailOrPhone: "admin@example.com",
		Password:     "Password123!",
	})
	if err != nil {
		slog.Error("Login failed", "error", err)
		return
	}

	fmt.Printf("login succeeded; access_token=%s refresh_token=%s\n", resp.AccessToken, resp.RefreshToken)
}
