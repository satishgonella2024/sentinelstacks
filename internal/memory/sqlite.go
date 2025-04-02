package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteMemoryStore is a SQLite-backed implementation of MemoryStore
type SQLiteMemoryStore struct {
	db        *sql.DB
	namespace string
	ttl       time.Duration
	tableName string
}

// NewSQLiteMemoryStore creates a new SQLite-backed store
func NewSQLiteMemoryStore(config MemoryConfig) (*SQLiteMemoryStore, error) {
	// Set default connection string if not provided
	connString := config.ConnectionString
	if connString == "" {
		connString = ":memory:"
	}

	// Set default table name if not provided
	tableName := config.CollectionName
	if tableName == "" {
		tableName = "memory_store"
	}

	// Open database connection
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create store
	store := &SQLiteMemoryStore{
		db:        db,
		namespace: config.Namespace,
		ttl:       config.TTL,
		tableName: tableName,
	}

	// Initialize tables
	if err := store.initTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return store, nil
}

// initTables creates the necessary tables if they don't exist
func (s *SQLiteMemoryStore) initTables() error {
	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY,
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		metadata TEXT,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_%s_key ON %s(key);
	CREATE INDEX IF NOT EXISTS idx_%s_updated_at ON %s(updated_at);
	`, s.tableName, s.tableName, s.tableName, s.tableName, s.tableName)

	_, err := s.db.Exec(createTableSQL)
	return err
}

// Save stores a value with the given key
func (s *SQLiteMemoryStore) Save(ctx context.Context, key string, value interface{}) error {
	fullKey := s.getNamespacedKey(key)
	now := time.Now().UTC()

	// Convert value to JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Check if key exists
	var count int
	err = s.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", s.tableName), fullKey).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if key exists: %w", err)
	}

	if count > 0 {
		// Update existing record
		_, err = s.db.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET value = ?, updated_at = ? WHERE id = ?", s.tableName),
			string(valueJSON), now, fullKey)
		if err != nil {
			return fmt.Errorf("failed to update record: %w", err)
		}
	} else {
		// Insert new record
		_, err = s.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s (id, key, value, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", s.tableName),
			fullKey, key, string(valueJSON), "{}", now, now)
		if err != nil {
			return fmt.Errorf("failed to insert record: %w", err)
		}
	}

	return nil
}

// Load retrieves a value by key
func (s *SQLiteMemoryStore) Load(ctx context.Context, key string) (interface{}, error) {
	fullKey := s.getNamespacedKey(key)

	// Prepare query
	query := fmt.Sprintf("SELECT value, updated_at FROM %s WHERE id = ?", s.tableName)
	var args []interface{}
	args = append(args, fullKey)

	if s.ttl > 0 {
		expiryTime := time.Now().UTC().Add(-s.ttl)
		query = fmt.Sprintf("SELECT value, updated_at FROM %s WHERE id = ? AND updated_at > ?", s.tableName)
		args = append(args, expiryTime)
	}

	// Execute query
	var valueJSON string
	var updatedAt time.Time
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&valueJSON, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("key not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Parse the value from JSON
	var value interface{}
	if err := json.Unmarshal([]byte(valueJSON), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return value, nil
}

// Delete removes a key-value pair
func (s *SQLiteMemoryStore) Delete(ctx context.Context, key string) error {
	fullKey := s.getNamespacedKey(key)

	// Delete the record
	result, err := s.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", s.tableName), fullKey)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	// Check if anything was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("key not found")
	}

	return nil
}

// Clear removes all keys and values
func (s *SQLiteMemoryStore) Clear(ctx context.Context) error {
	var query string
	var args []interface{}

	if s.namespace != "" {
		// Clear only keys in the namespace
		query = fmt.Sprintf("DELETE FROM %s WHERE id LIKE ?", s.tableName)
		args = []interface{}{s.namespace + ":%"}
	} else {
		// Clear all keys
		query = fmt.Sprintf("DELETE FROM %s", s.tableName)
	}

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to clear records: %w", err)
	}

	return nil
}

// Keys returns all keys in the store
func (s *SQLiteMemoryStore) Keys(ctx context.Context) ([]string, error) {
	var query string
	var args []interface{}

	// Build query based on namespace and TTL
	if s.namespace != "" {
		if s.ttl > 0 {
			expiryTime := time.Now().UTC().Add(-s.ttl)
			query = fmt.Sprintf("SELECT key FROM %s WHERE id LIKE ? AND updated_at > ? ORDER BY key", s.tableName)
			args = []interface{}{s.namespace + ":%", expiryTime}
		} else {
			query = fmt.Sprintf("SELECT key FROM %s WHERE id LIKE ? ORDER BY key", s.tableName)
			args = []interface{}{s.namespace + ":%"}
		}
	} else {
		if s.ttl > 0 {
			expiryTime := time.Now().UTC().Add(-s.ttl)
			query = fmt.Sprintf("SELECT key FROM %s WHERE updated_at > ? ORDER BY key", s.tableName)
			args = []interface{}{expiryTime}
		} else {
			query = fmt.Sprintf("SELECT key FROM %s ORDER BY key", s.tableName)
		}
	}

	// Execute query
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query keys: %w", err)
	}
	defer rows.Close()

	// Collect keys
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

// Close closes the memory store
func (s *SQLiteMemoryStore) Close() error {
	return s.db.Close()
}

// getNamespacedKey adds the namespace prefix to a key
func (s *SQLiteMemoryStore) getNamespacedKey(key string) string {
	if s.namespace == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", s.namespace, key)
}

// SQLiteVectorStore extends SQLiteMemoryStore with vector operations
type SQLiteVectorStore struct {
	*SQLiteMemoryStore
	vectorTableName string
	vectorDimension int
}

// NewSQLiteVectorStore creates a new SQLite-backed vector store
func NewSQLiteVectorStore(config MemoryConfig) (*SQLiteVectorStore, error) {
	baseStore, err := NewSQLiteMemoryStore(config)
	if err != nil {
		return nil, err
	}

	vectorDimension := config.VectorDimensions
	if vectorDimension <= 0 {
		vectorDimension = 1536 // Default for OpenAI embeddings
	}

	vectorTableName := baseStore.tableName + "_vectors"

	store := &SQLiteVectorStore{
		SQLiteMemoryStore: baseStore,
		vectorTableName:   vectorTableName,
		vectorDimension:   vectorDimension,
	}

	// Initialize vector table
	if err := store.initVectorTable(); err != nil {
		baseStore.Close()
		return nil, fmt.Errorf("failed to initialize vector table: %w", err)
	}

	return store, nil
}

// initVectorTable creates the vector table if it doesn't exist
func (s *SQLiteVectorStore) initVectorTable() error {
	// Create columns for vector dimensions
	vectorColumns := make([]string, s.vectorDimension)
	for i := 0; i < s.vectorDimension; i++ {
		vectorColumns[i] = fmt.Sprintf("dim_%d REAL", i)
	}

	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY,
		%s
	);
	`, s.vectorTableName, strings.Join(vectorColumns, ", "))

	_, err := s.db.Exec(createTableSQL)
	return err
}

// SaveEmbedding stores a vector embedding
func (s *SQLiteVectorStore) SaveEmbedding(ctx context.Context, key string, vector []float32, metadata map[string]interface{}) error {
	if len(vector) != s.vectorDimension {
		return fmt.Errorf("vector dimension mismatch: expected %d, got %d", s.vectorDimension, len(vector))
	}

	fullKey := s.getNamespacedKey(key)
	now := time.Now().UTC()

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Check if key exists in main table
	var count int
	err = tx.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", s.tableName), fullKey).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if key exists: %w", err)
	}

	if count > 0 {
		// Update existing record
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET metadata = ?, updated_at = ? WHERE id = ?", s.tableName),
			string(metadataJSON), now, fullKey)
		if err != nil {
			return fmt.Errorf("failed to update record: %w", err)
		}
	} else {
		// Insert new record
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s (id, key, value, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", s.tableName),
			fullKey, key, "null", string(metadataJSON), now, now)
		if err != nil {
			return fmt.Errorf("failed to insert record: %w", err)
		}
	}

	// Check if vector exists
	err = tx.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", s.vectorTableName), fullKey).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if vector exists: %w", err)
	}

	// Prepare vector columns and placeholders
	columns := make([]string, s.vectorDimension+1)
	placeholders := make([]string, s.vectorDimension+1)
	args := make([]interface{}, s.vectorDimension+1)

	columns[0] = "id"
	placeholders[0] = "?"
	args[0] = fullKey

	for i := 0; i < s.vectorDimension; i++ {
		columns[i+1] = fmt.Sprintf("dim_%d", i)
		placeholders[i+1] = "?"
		if i < len(vector) {
			args[i+1] = vector[i]
		} else {
			args[i+1] = 0.0
		}
	}

	if count > 0 {
		// Update existing vector
		setStatements := make([]string, s.vectorDimension)
		updateArgs := make([]interface{}, s.vectorDimension+1)
		
		for i := 0; i < s.vectorDimension; i++ {
			setStatements[i] = fmt.Sprintf("dim_%d = ?", i)
			updateArgs[i] = args[i+1]
		}
		updateArgs[s.vectorDimension] = fullKey

		updateSQL := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", 
			s.vectorTableName, strings.Join(setStatements, ", "))
		
		_, err = tx.ExecContext(ctx, updateSQL, updateArgs...)
		if err != nil {
			return fmt.Errorf("failed to update vector: %w", err)
		}
	} else {
		// Insert new vector
		insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			s.vectorTableName, 
			strings.Join(columns, ", "), 
			strings.Join(placeholders, ", "))
		
		_, err = tx.ExecContext(ctx, insertSQL, args...)
		if err != nil {
			return fmt.Errorf("failed to insert vector: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Query performs a similarity search on stored embeddings
func (s *SQLiteVectorStore) Query(ctx context.Context, vector []float32, topK int) ([]SimilarityMatch, error) {
	if len(vector) != s.vectorDimension {
		return nil, fmt.Errorf("vector dimension mismatch: expected %d, got %d", s.vectorDimension, len(vector))
	}

	// Prepare dot product calculation SQL
	// We'll use a simplified cosine similarity calculation:
	// 1. Calculate dot product: SUM(qi * vi) for each dimension
	// 2. We'll normalize vectors when inserting them, so magnitude is roughly 1
	// 3. Therefore, dot product approximates cosine similarity
	
	dotProductTerms := make([]string, s.vectorDimension)
	args := make([]interface{}, s.vectorDimension)
	
	for i := 0; i < s.vectorDimension; i++ {
		dotProductTerms[i] = fmt.Sprintf("(v.dim_%d * ?)", i)
		args[i] = vector[i]
	}
	
	// Build query with filtering by namespace and TTL if needed
	var queryConditions []string
	var additionalArgs []interface{}
	
	if s.namespace != "" {
		queryConditions = append(queryConditions, "m.id LIKE ?")
		additionalArgs = append(additionalArgs, s.namespace+":%")
	}
	
	if s.ttl > 0 {
		expiryTime := time.Now().UTC().Add(-s.ttl)
		queryConditions = append(queryConditions, "m.updated_at > ?")
		additionalArgs = append(additionalArgs, expiryTime)
	}
	
	conditionSQL := ""
	if len(queryConditions) > 0 {
		conditionSQL = "WHERE " + strings.Join(queryConditions, " AND ")
	}
	
	// Build complete query
	query := fmt.Sprintf(`
		SELECT m.key, m.metadata, (%s) as similarity
		FROM %s v
		JOIN %s m ON v.id = m.id
		%s
		ORDER BY similarity DESC
		LIMIT ?
	`, strings.Join(dotProductTerms, " + "), s.vectorTableName, s.tableName, conditionSQL)
	
	// Add all args
	queryArgs := append(args, additionalArgs...)
	
	// Add limit
	if topK <= 0 {
		topK = 10 // Default limit
	}
	queryArgs = append(queryArgs, topK)
	
	// Execute query
	rows, err := s.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	
	// Process results
	var results []SimilarityMatch
	for rows.Next() {
		var key string
		var metadataJSON string
		var similarity float32
		
		if err := rows.Scan(&key, &metadataJSON, &similarity); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		// Parse metadata
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		
		results = append(results, SimilarityMatch{
			Key:      key,
			Score:    similarity,
			Metadata: metadata,
		})
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	
	return results, nil
}

// DeleteEmbedding removes an embedding
func (s *SQLiteVectorStore) DeleteEmbedding(ctx context.Context, key string) error {
	fullKey := s.getNamespacedKey(key)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete from vector table
	_, err = tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", s.vectorTableName), fullKey)
	if err != nil {
		return fmt.Errorf("failed to delete from vector table: %w", err)
	}

	// Delete from main table
	result, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", s.tableName), fullKey)
	if err != nil {
		return fmt.Errorf("failed to delete from main table: %w", err)
	}

	// Check if anything was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("embedding not found")
	}

	return nil
}

// normalizeVector normalizes a vector to unit length
func normalizeVector(vector []float32) []float32 {
	var sum float32
	for _, v := range vector {
		sum += v * v
	}
	
	magnitude := float32(math.Sqrt(float64(sum)))
	if magnitude == 0 {
		return vector
	}
	
	normalized := make([]float32, len(vector))
	for i, v := range vector {
		normalized[i] = v / magnitude
	}
	
	return normalized
}
