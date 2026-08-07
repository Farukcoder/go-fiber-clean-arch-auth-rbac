package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName          string
	AppEnv           string
	Port             string
	DBConnection     string
	DBHost           string
	DBPort           string
	DBDatabase       string
	DBUsername       string
	DBPassword       string
	DBSSLMode        string
	DBUrl            string
	JwtSecret        string
	JwtRefreshSecret string
	AllowedOrigins   string
	MailMailer       string
	MailHost         string
	MailPort         string
	MailUsername     string
	MailPassword     string
	MailEncryption   string
	MailFromAddress  string
	MailFromName     string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	appName := getEnv("APP_NAME", "GO")

	cfg := &Config{
		AppName:          appName,
		AppEnv:           getEnv("APP_ENV", "local"),
		Port:             getEnv("PORT", "8080"),
		DBConnection:     getEnv("DB_CONNECTION", "mysql"),
		DBHost:           getEnv("DB_HOST", "127.0.0.1"),
		DBPort:           getEnv("DB_PORT", "3306"),
		DBDatabase:       getEnv("DB_DATABASE", ""),
		DBUsername:       getEnv("DB_USERNAME", ""),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
		JwtRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		AllowedOrigins:   getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173,http://localhost:5174"),
		MailMailer:       getEnv("MAIL_MAILER", "smtp"),
		MailHost:         getEnv("MAIL_HOST", "smtp.gmail.com"),
		MailPort:         getEnv("MAIL_PORT", "587"),
		MailUsername:     getEnv("MAIL_USERNAME", ""),
		MailPassword:     getEnv("MAIL_PASSWORD", ""),
		MailEncryption:   getEnv("MAIL_ENCRYPTION", "tls"),
		MailFromAddress:  getEnv("MAIL_FROM_ADDRESS", ""),
		MailFromName:     getEnv("MAIL_FROM_NAME", appName),
	}

	// Build DSN from individual components
	cfg.DBUrl = buildDSN(cfg)

	if cfg.DBUrl == "" {
		return nil, errors.New("missing database configuration")
	}

	if cfg.JwtSecret == "" {
		return nil, errors.New("missing JWT_SECRET environment variable")
	}

	if cfg.JwtRefreshSecret == "" {
		return nil, errors.New("missing JWT_REFRESH_SECRET environment variable")
	}

	return cfg, nil
}

func buildDSN(cfg *Config) string {
	if cfg.DBConnection == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUsername,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBDatabase,
		)
	} else if cfg.DBConnection == "postgres" || cfg.DBConnection == "pgsql" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUsername,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBDatabase,
			cfg.DBSSLMode,
		)
	}
	return ""
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
