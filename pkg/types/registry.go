// Package types contains common interfaces and types used across packages
package types

import (
	"context"
	"time"
)

// PackageType defines the type of content in a package
type PackageType string

const (
	// PackageTypeAgent represents an agent package
	PackageTypeAgent PackageType = "agent"

	// PackageTypeStack represents a stack package
	PackageTypeStack PackageType = "stack"
)

// Dependency represents a dependency on another package
type Dependency struct {
	Name     string      `json:"name"`
	Version  string      `json:"version"`
	Type     PackageType `json:"type"`
	Required bool        `json:"required"`
}

// PackageInfo represents metadata about a package
type PackageInfo struct {
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Type         PackageType            `json:"type"`
	Description  string                 `json:"description"`
	Author       string                 `json:"author"`
	License      string                 `json:"license,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty"`
	Size         int64                  `json:"size,omitempty"`
	Downloads    int                    `json:"downloads,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Dependencies []Dependency           `json:"dependencies,omitempty"`
	Verified     bool                   `json:"verified,omitempty"`
}

// PackageSearchResult represents a search result from the registry
type PackageSearchResult struct {
	TotalCount int           `json:"totalCount"`
	Items      []PackageInfo `json:"items"`
}

// AuthProvider defines the interface for registry authentication
type AuthProvider interface {
	// GetToken returns an authentication token
	GetToken(ctx context.Context) (string, error)

	// Login performs authentication and returns a token
	Login(ctx context.Context, username, password string) (string, error)

	// Logout invalidates the current token
	Logout(ctx context.Context) error

	// IsAuthenticated checks if the client is authenticated
	IsAuthenticated(ctx context.Context) bool
}

// RegistryClient defines the interface for interacting with package registries
type RegistryClient interface {
	// Push pushes a package to the registry
	Push(ctx context.Context, packagePath string) error

	// Pull pulls a package from the registry and returns the local path
	Pull(ctx context.Context, name, version string) (string, error)

	// Search searches for packages in the registry
	Search(ctx context.Context, query string, limit int) ([]PackageInfo, error)

	// GetPackageInfo gets information about a package
	GetPackageInfo(ctx context.Context, name, version string) (*PackageInfo, error)

	// ListTags lists all tags (versions) for a package
	ListTags(ctx context.Context, name string) ([]string, error)
}

// RegistryPackage represents a package in the registry
type RegistryPackage struct {
	ID          string
	Name        string
	Version     string
	Description string
	Tags        []string
	Metadata    map[string]interface{}
	CreatedAt   string
	UpdatedAt   string
}

// RegistrySearchRequest represents a search request to the registry
type RegistrySearchRequest struct {
	Query  string
	Tags   []string
	Limit  int
	Offset int
}

// RegistrySearchResponse represents a search response from the registry
type RegistrySearchResponse struct {
	Results []RegistryPackage
	Total   int
	Offset  int
	Limit   int
}
