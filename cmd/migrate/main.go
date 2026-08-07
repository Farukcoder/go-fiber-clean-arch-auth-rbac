package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

func main() {
	config.InitLogger()

	if len(os.Args) < 2 {
		slog.Error("Usage: go run ./cmd/migrate [up|down|fresh|version|force <version>]")
		os.Exit(1)
	}

	action := os.Args[1]
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	createTableStart := time.Now()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Failed to create migrate driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance("file://database/migrations", "pgx", driver)
	if err != nil {
		slog.Error("Failed to create migrate instance", "error", err)
		os.Exit(1)
	}

	switch action {
	case "up":
		if err := runMigrationsUp(m, "database/migrations", time.Since(createTableStart)); err != nil {
			slog.Error("Migration up failed", "error", err)
			os.Exit(1)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			slog.Error("Migration down failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Migration down completed")
	case "fresh":
		if err := m.Drop(); err != nil && err != migrate.ErrNoChange {
			slog.Error("Migration fresh drop failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Database dropped, recreating migrations")
		driver, err = postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			slog.Error("Failed to recreate migrate driver", "error", err)
			os.Exit(1)
		}
		m, err = migrate.NewWithDatabaseInstance("file://database/migrations", "pgx", driver)
		if err != nil {
			slog.Error("Failed to create migrate instance", "error", err)
			os.Exit(1)
		}
		if err := runMigrationsUp(m, "database/migrations", time.Since(createTableStart)); err != nil {
			slog.Error("Migration fresh failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Migration fresh completed")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			slog.Error("Migration version check failed", "error", err)
			os.Exit(1)
		}
		fmt.Printf("version=%d dirty=%v\n", version, dirty)
	case "force":
		if len(os.Args) < 3 {
			slog.Error("Usage: go run ./cmd/migrate force <version>")
			os.Exit(1)
		}
		var v int
		_, err := fmt.Sscanf(os.Args[2], "%d", &v)
		if err != nil {
			slog.Error("Invalid version", "error", err)
			os.Exit(1)
		}
		if err := m.Force(v); err != nil {
			slog.Error("Force migration failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Migration force completed", "version", v)
	default:
		slog.Error("Unsupported action", "action", action)
		os.Exit(1)
	}
}

type migrationFile struct {
	Version uint
	Name    string
}

func runMigrationsUp(m *migrate.Migrate, migrationDir string, createTableElapsed time.Duration) error {
	currentVersion, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion || isMissingSchemaMigrationsError(err) {
			currentVersion = 0
		} else {
			return err
		}
	}

	if dirty {
		return fmt.Errorf("database is dirty at version %d", currentVersion)
	}

	migrations, err := loadMigrationFiles(migrationDir)
	if err != nil {
		return err
	}

	pending := make([]migrationFile, 0, len(migrations))
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			pending = append(pending, migration)
		}
	}

	if len(pending) == 0 {
		fmt.Println("Nothing to migrate.")
		return nil
	}

	fmt.Println("INFO Preparing database.")
	fmt.Println(formatProgressLine("Creating migration table", createTableElapsed))
	fmt.Println("INFO Running migrations.")
	for _, migration := range pending {
		started := time.Now()
		if err := m.Steps(1); err != nil && err != migrate.ErrNoChange {
			return err
		}
		fmt.Println(formatProgressLine(migration.Name, time.Since(started)))
	}

	return nil
}

func loadMigrationFiles(migrationDir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	migrations := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		base := strings.TrimSuffix(name, ".up.sql")
		parts := strings.SplitN(base, "_", 2)
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid migration filename %q: %w", name, err)
		}

		migrations = append(migrations, migrationFile{
			Version: uint(version),
			Name:    base,
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func formatProgressLine(label string, elapsed time.Duration) string {
	duration := fmt.Sprintf("%dms", elapsed.Milliseconds())
	width := 60
	dots := width - len(label) - len(duration) - len(" DONE") - 2
	if dots < 1 {
		dots = 1
	}

	return fmt.Sprintf("%s %s %s DONE", label, strings.Repeat(".", dots), duration)
}

func isMissingSchemaMigrationsError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*pq.Error); ok {
		return e.Code.Name() == "undefined_table"
	}
	return strings.Contains(err.Error(), "relation \"public.schema_migrations\" does not exist") ||
		(strings.Contains(err.Error(), "schema_migrations") && strings.Contains(err.Error(), "does not exist"))
}
