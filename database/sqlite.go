package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const dbPath = "./data/dacademy.db"
const migrationPath = "./database/migrations"

func Run() (*sqlx.DB, error) {
	// new dir
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create db dir: %w", err)
	}

	// new file
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		f, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("cannot create db file: %w", err)
		}
		f.Close()
	}

	// Connect
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable WAL mode
	_, _ = db.Exec("PRAGMA journal_mode=WAL;")

	// Run migrations
	m, err := migrate.New(
		"file://"+migrationPath,
		"sqlite3://"+dbPath,
	)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return db, nil
}
