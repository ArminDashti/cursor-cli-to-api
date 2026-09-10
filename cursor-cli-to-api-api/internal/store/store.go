package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	for _, name := range files {
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(body)); err != nil {
			return fmt.Errorf("migrate %s: %w", name, err)
		}
	}
	return nil
}

func SeedDefaultUser(ctx context.Context, db *sql.DB, username, passwordHash string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (username) DO UPDATE SET password_hash = EXCLUDED.password_hash
	`, username, passwordHash)
	return err
}

func SaveResponse(ctx context.Context, db *sql.DB, id, model, status string, body any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `
		INSERT INTO responses (id, model, status, body)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET model = EXCLUDED.model, status = EXCLUDED.status, body = EXCLUDED.body
	`, id, model, status, raw)
	return err
}

func GetResponse(ctx context.Context, db *sql.DB, id string) (json.RawMessage, error) {
	var raw []byte
	err := db.QueryRowContext(ctx, `SELECT body FROM responses WHERE id = $1`, id).Scan(&raw)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}
