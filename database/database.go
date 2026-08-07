package database

import (
	"fmt"
	"time"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DBConnection {
	case "mysql":
		dialector = mysql.Open(cfg.DBUrl)
	case "postgres", "pgsql", "postgresql":
		dialector = postgres.Open(cfg.DBUrl)
	default:
		return nil, fmt.Errorf("unsupported database connection type: %s", cfg.DBConnection)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
