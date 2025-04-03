// Package types defines common type definitions used across packages
package types

import (
	"context"
)

// StackService defines operations for managing stacks
type StackService interface {
	// CreateStack creates a new stack
	CreateStack(ctx context.Context, spec StackSpec) (string, error)

	// ExecuteStack executes a stack with given inputs
	ExecuteStack(ctx context.Context, stackID string, inputs map[string]interface{}) (map[string]interface{}, error)

	// GetStackState gets the current state of a stack
	GetStackState(ctx context.Context, stackID string) (*StackExecutionSummary, error)

	// ListStacks lists all available stacks
	ListStacks(ctx context.Context) ([]StackInfo, error)

	// UpdateStack updates an existing stack
	UpdateStack(ctx context.Context, stackID string, spec StackSpec) error

	// DeleteStack removes a stack
	DeleteStack(ctx context.Context, stackID string) error

	// ImportStack imports a stack from a file
	ImportStack(ctx context.Context, filePath string) (string, error)

	// ExportStack exports a stack to a file
	ExportStack(ctx context.Context, stackID, filePath string) error

	// GetStackExecutionHistory gets the execution history for a stack
	GetStackExecutionHistory(ctx context.Context, stackID string) ([]ExecutionSummary, error)
}

// MemoryService defines operations for managing memory
type MemoryService interface {
	// StoreValue stores a value in memory
	StoreValue(ctx context.Context, collection string, key string, value interface{}) error

	// RetrieveValue retrieves a value from memory
	RetrieveValue(ctx context.Context, collection string, key string) (interface{}, error)

	// StoreEmbedding stores text with vector embedding
	StoreEmbedding(ctx context.Context, collection string, key string, text string, metadata map[string]interface{}) error

	// SearchSimilar finds similar texts using vector similarity
	SearchSimilar(ctx context.Context, collection string, text string, limit int) ([]MemoryMatch, error)
}

// RegistryService defines operations for registry interaction
type RegistryService interface {
	// PushPackage pushes a package to the registry
	PushPackage(ctx context.Context, path string) error

	// PullPackage pulls a package from the registry
	PullPackage(ctx context.Context, name string, version string) (string, error)

	// SearchPackages searches for packages in the registry
	SearchPackages(ctx context.Context, query string, limit int) ([]PackageInfo, error)
}

// API provides access to all services
type API interface {
	// Stack returns the stack service
	Stack() StackService

	// Memory returns the memory service
	Memory() MemoryService

	// Registry returns the registry service
	Registry() RegistryService

	// Close releases any resources held by the API services
	Close() error
}

// StackInfo contains basic information about a stack
type StackInfo struct {
	ID          string
	Name        string
	Description string
	Version     string
	Type        StackType
	CreatedAt   string
}
