package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/auth"
	"github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/config"
	httpserver "github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/http"
	"github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/store"
)

func main() {
	// Load .env from cwd or next to the binary's module root when run via go run.
	config.LoadDotEnv(".env")
	if _, err := os.Stat(".env"); err != nil {
		if wd, err2 := os.Getwd(); err2 == nil {
			config.LoadDotEnv(filepath.Join(wd, "cursor-cli-to-api-api", ".env"))
		}
	}

	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := store.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	migDir := cfg.MigrationsDir
	if !filepath.IsAbs(migDir) {
		if _, err := os.Stat(migDir); err != nil {
			alt := filepath.Join("cursor-cli-to-api-api", migDir)
			if _, err2 := os.Stat(alt); err2 == nil {
				migDir = alt
			}
		}
	}
	if err := store.Migrate(db, migDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	hash, err := auth.HashPassword(cfg.AuthPassword)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	if err := store.SeedDefaultUser(ctx, db, cfg.AuthUsername, hash); err != nil {
		log.Fatalf("seed user: %v", err)
	}

	srv := httpserver.New(cfg, db)
	log.Printf("cursor-cli-to-api-api listening on %s", cfg.Addr)
	if err := srv.Router().Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
