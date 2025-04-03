// Package memory provides memory store implementations
package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// SQLiteVectorStore is a SQLite-backed vector store
type SQLiteVectorStore struct {
	db         *sql.DB
	path       string
	dimensions int
}

// NewSQLiteVectorStore creates a new SQLite-backed vector store
func NewSQLiteVectorStore(path string, dimensions int) (*SQLiteVectorStore, error) {
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

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS vectors (
			key TEXT PRIMARY KEY,
			text TEXT NOT NULL,
			metadata TEXT,
			timestamp TIMESTAMP NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS embeddings (
			key TEXT PRIMARY KEY,
			embedding BLOB NOT NULL,
			dimensions INTEGER NOT NULL,
			FOREIGN KEY (key) REFERENCES vectors(key) ON DELETE CASCADE
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &SQLiteVectorStore{
		db:         db,
		path:       path,
		dimensions: dimensions,
	}, nil
}

// StoreVector stores a text with its vector embedding
func (s *SQLiteVectorStore) StoreVector(ctx context.Context, key string, text string, metadata map[string]interface{}) error {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Marshal metadata
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Insert or update vector record
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO vectors (key, text, metadata, timestamp)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			text = excluded.text,
			metadata = excluded.metadata,
			timestamp = excluded.timestamp;`,
		key, text, string(metadataJSON), time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert/update vector: %w", err)
	}

	// In a real implementation, we would generate embeddings here
	// For now, we'll just create a mock vector
	vector := make([]float32, s.dimensions)

	// Generate a simple mock vector based on the text content
	// This is just for demonstration purposes
	for i := 0; i < s.dimensions && i < len(text); i++ {
		if i < len(text) {
			vector[i] = float32(text[i]) / 255.0
		}
	}

	// Serialize the vector
	vectorBytes := make([]byte, len(vector)*4)
	for i, v := range vector {
		copy(vectorBytes[i*4:], float32ToBytes(v))
	}

	// Insert or update embedding record
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO embeddings (key, embedding, dimensions)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			embedding = excluded.embedding,
			dimensions = excluded.dimensions;`,
		key, vectorBytes, s.dimensions,
	)
	if err != nil {
		return fmt.Errorf("failed to insert/update embedding: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SearchVector finds similar vectors using cosine similarity
func (s *SQLiteVectorStore) SearchVector(ctx context.Context, text string, limit int, filter map[string]interface{}) ([]types.MemoryMatch, error) {
	if limit <= 0 {
		limit = 10
	}

	// In a real implementation, we would generate embeddings for the query text
	// For now, we'll create a mock query vector
	queryVector := make([]float32, s.dimensions)
	for i := 0; i < s.dimensions && i < len(text); i++ {
		if i < len(text) {
			queryVector[i] = float32(text[i]) / 255.0
		}
	}

	// Get all vectors from the database
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT v.key, v.text, v.metadata, e.embedding, v.timestamp
		FROM vectors v
		JOIN embeddings e ON v.key = e.key;`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query vectors: %w", err)
	}
	defer rows.Close()

	// Calculate similarity for all vectors
	type matchScore struct {
		key       string
		text      string
		metadata  map[string]interface{}
		score     float64
		timestamp time.Time
	}

	var matches []matchScore
	for rows.Next() {
		var key, text, metadataJSON string
		var embeddingBytes []byte
		var timestamp time.Time

		if err := rows.Scan(&key, &text, &metadataJSON, &embeddingBytes, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Parse metadata
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		// Check filter
		if !matchesFilter(metadata, filter) {
			continue
		}

		// Convert embedding bytes to vector
		if len(embeddingBytes) != s.dimensions*4 {
			return nil, fmt.Errorf("invalid embedding size: %d", len(embeddingBytes))
		}
		vector := make([]float32, s.dimensions)
		for i := 0; i < s.dimensions; i++ {
			vector[i] = bytesToFloat32(embeddingBytes[i*4 : (i+1)*4])
		}

		// Calculate cosine similarity
		similarity := cosineSimilarity(queryVector, vector)

		matches = append(matches, matchScore{
			key:       key,
			text:      text,
			metadata:  metadata,
			score:     float64(similarity),
			timestamp: timestamp,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Sort by similarity (descending)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	// Limit results
	if len(matches) > limit {
		matches = matches[:limit]
	}

	// Convert to MemoryMatch objects
	result := make([]types.MemoryMatch, len(matches))
	for i, match := range matches {
		result[i] = types.MemoryMatch{
			Key:       match.key,
			Content:   match.text,
			Metadata:  match.metadata,
			Score:     match.score,
			Distance:  1.0 - match.score,
			Timestamp: match.timestamp,
		}
	}

	return result, nil
}

// GetVector retrieves a vector by key
func (s *SQLiteVectorStore) GetVector(ctx context.Context, key string) (*types.MemoryMatch, error) {
	var text, metadataJSON string
	var timestamp time.Time

	err := s.db.QueryRowContext(
		ctx,
		`SELECT text, metadata, timestamp FROM vectors WHERE key = ?;`,
		key,
	).Scan(&text, &metadataJSON, &timestamp)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vector not found: %s", key)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query vector: %w", err)
	}

	// Parse metadata
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &types.MemoryMatch{
		Key:       key,
		Content:   text,
		Metadata:  metadata,
		Score:     1.0, // Perfect match for direct retrieval
		Distance:  0.0,
		Timestamp: timestamp,
	}, nil
}

// DeleteVector removes a vector by key
func (s *SQLiteVectorStore) DeleteVector(ctx context.Context, key string) error {
	_, err := s.db.ExecContext(
		ctx,
		`DELETE FROM vectors WHERE key = ?;`,
		key,
	)
	if err != nil {
		return fmt.Errorf("failed to delete vector: %w", err)
	}
	return nil
}

// ListVectors returns all keys in the store
func (s *SQLiteVectorStore) ListVectors(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT key FROM vectors;`)
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
func (s *SQLiteVectorStore) Close() error {
	return s.db.Close()
}

// float32ToBytes converts a float32 to bytes
func float32ToBytes(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 4)
	bytes[0] = byte(bits)
	bytes[1] = byte(bits >> 8)
	bytes[2] = byte(bits >> 16)
	bytes[3] = byte(bits >> 24)
	return bytes
}

// bytesToFloat32 converts bytes to a float32
func bytesToFloat32(bytes []byte) float32 {
	bits := uint32(bytes[0]) |
		uint32(bytes[1])<<8 |
		uint32(bytes[2])<<16 |
		uint32(bytes[3])<<24
	return math.Float32frombits(bits)
}
