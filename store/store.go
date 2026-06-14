package store

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func Init(dsn string) error {
	var err error
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)
	DB.SetConnMaxLifetime(0)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}

func runMigrations() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
        id          INTEGER PRIMARY KEY AUTOINCREMENT,
        code        TEXT    NOT NULL UNIQUE,
        original_url TEXT   NOT NULL,
        created_at  TEXT    NOT NULL,
        expires_at  TEXT,
        clicks      INTEGER NOT NULL DEFAULT 0,
        is_active   INTEGER NOT NULL DEFAULT 1
    );

    CREATE INDEX IF NOT EXISTS idx_urls_code ON urls(code);
    CREATE INDEX IF NOT EXISTS idx_urls_is_active ON urls(is_active);
	`

	_, err := DB.Exec(query)
	return err
}
