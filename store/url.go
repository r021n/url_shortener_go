package store

import (
	"database/sql"
	"errors"
	"time"
)

type URL struct {
	ID          int64   `json:"id"`
	Code        string  `json:"code"`
	OriginalURL string  `json:"original_url"`
	CreatedAt   string  `json:"created_at"`
	ExpiresAt   *string `json:"expires_at"`
	Clicks      int64   `json:"clicks"`
	IsActive    bool    `json:"is_active"`
}

func InsertURL(code string, originalURL string, expiresAt *string) (*URL, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	result, err := DB.Exec(
		`INSERT INTO urls (code, original_url, created_at, expires_at)
		VALUES (?, ?, ?, ?)`,
		code, originalURL, now, expiresAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetURLByID(id)
}

func GetURLByCode(code string) (*URL, error) {
	var url URL
	var isActive int

	err := DB.QueryRow(
		`SELECT id, code, original_url, created_at, expires_at, clicks, is_active
		FROM urls WHERE code = ?`, code,
	).Scan(&url.ID, &url.Code, &url.OriginalURL, &url.CreatedAt, &url.ExpiresAt, &url.Clicks, &url.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	url.IsActive = isActive == 1
	return &url, nil
}
