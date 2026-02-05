package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/price-comparison/server/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run cmd/migrate/main.go [up|down|version|drop|force <version>]")
	}

	cfg := config.Load()
	command := os.Args[1]

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = filepath.Clean("../../packages/database/migrations")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations applied")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Rolled back one migration")
	case "version":
		version, dirty, err := m.Version()
		if err == migrate.ErrNilVersion {
			log.Println("No migrations applied")
			return
		}
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		log.Printf("Version: %d (dirty=%v)", version, dirty)
	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Drop failed: %v", err)
		}
		log.Println("Database dropped")
	case "force":
		if len(os.Args) < 3 {
			log.Fatalf("Usage: go run cmd/migrate/main.go force <version>")
		}
		var version int
		if _, err := fmt.Sscanf(os.Args[2], "%d", &version); err != nil {
			log.Fatalf("Invalid version: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("Force failed: %v", err)
		}
		log.Printf("Forced version to %d", version)
	default:
		log.Fatalf("Unknown command: %s", command)
	}

	_, _ = m.Close()
}
