// Package memory provides memory store implementations
package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore is a SQLite-backed implementation of MemoryStore
type SQLiteStore struct {
	db   *sql.DB
	path string
}

// NewSQLiteStore creates a new SQLite-backed memory store
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS memory (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			metadata TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteStore{
		db:   db,
		path: path,
	}, nil
}

// Save stores a value by key
func (s *SQLiteStore) Save(ctx context.Context, key string, value interface{}) error {
	// Convert value to JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Store value in database
	now := time.Now().UTC()
	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO memory (key, value, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at;`,
		key, string(valueJSON), "{}", now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to insert/update value: %w", err)
	}

	return nil
}

// Load retrieves a value by key
func (s *SQLiteStore) Load(ctx context.Context, key string) (interface{}, error) {
	var valueJSON string
	err := s.db.QueryRowContext(
		ctx,
		`SELECT value FROM memory WHERE key = ?;`,
		key,
	).Scan(&valueJSON)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query value: %w", err)
	}

	// Parse JSON value
	var value interface{}
	err = json.Unmarshal([]byte(valueJSON), &value)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return value, nil
}

// Delete removes a value by key
func (s *SQLiteStore) Delete(ctx context.Context, key string) error {
	_, err := s.db.ExecContext(
		ctx,
		`DELETE FROM memory WHERE key = ?;`,
		key,
	)
	if err != nil {
		return fmt.Errorf("failed to delete value: %w", err)
	}

	return nil
}

// List returns all keys
func (s *SQLiteStore) List(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT key FROM memory;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query keys: %w", err)
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("failed to scan key: %w", err)
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return keys, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
