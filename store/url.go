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
	).Scan(&url.ID, &url.Code, &url.OriginalURL, &url.CreatedAt, &url.ExpiresAt, &url.Clicks, &isActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	url.IsActive = isActive == 1
	return &url, nil
}

func GetURLByID(id int64) (*URL, error) {
	var url URL
	var isActive int

	err := DB.QueryRow(`SELECT id, code, original_url, created_at, expires_at, clicks, is_active FROM urls WHERE id = ?`, id).Scan(&url.ID, &url.Code, &url.OriginalURL, &url.CreatedAt, &url.ExpiresAt, &url.Clicks, &isActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	url.IsActive = isActive == 1
	return &url, nil
}

func IncrementClicks(code string) error {
	_, err := DB.Exec(
		`UPDATE urls SET clicks = clicks + 1 WHERE code = ?`, code,
	)
	return err
}

func DeactivateURL(code string) error {
	_, err := DB.Exec(
		`UPDATE urls SET is_active = 0 WHERE code = ?`, code,
	)
	return err
}

func ListURLs(limit int, offset int) ([]URL, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := DB.Query(
		`SELECT id, code, original_url, created_at, expires_at, clicks, is_active
		FROM urls ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		var isActive int
		if err := rows.Scan(&url.ID, &url.Code, &url.OriginalURL, &url.CreatedAt, &url.ExpiresAt, &url.Clicks, &isActive); err != nil {
			return nil, err
		}
		url.IsActive = isActive == 1
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func CountURLs() (int64, error) {
	var count int64
	err := DB.QueryRow(`SELECT COUNT(*) FROM urls`).Scan(&count)
	return count, err
}

func CodeExists(code string) (bool, error) {
	var exists int
	err := DB.QueryRow(
		`SELECT COUNT(*) FROM urls WHERE code = ?`, code,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
