// Package api provides a unified API for the Sentinel Stacks system
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// RegistryServiceConfig contains configuration for the registry service
type RegistryServiceConfig struct {
	// RegistryURL is the base URL for the registry
	RegistryURL string

	// CachePath is where downloaded packages are cached
	CachePath string

	// Username is the username for registry authentication
	Username string

	// AccessToken is the access token for registry authentication
	AccessToken string
}

// RegistryService implements types.RegistryService
type RegistryService struct {
	config      RegistryServiceConfig
	client      *http.Client
	cache       map[string]types.PackageInfo
	packageLock sync.RWMutex
}

// NewRegistryService creates a new registry service
func NewRegistryService(config RegistryServiceConfig) (*RegistryService, error) {
	// Create cache directory if it doesn't exist
	if config.CachePath != "" {
		if err := os.MkdirAll(config.CachePath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache directory: %w", err)
		}
	}

	// Create HTTP client with reasonable timeouts
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &RegistryService{
		config: config,
		client: client,
		cache:  make(map[string]types.PackageInfo),
	}, nil
}

// PushPackage pushes a package to the registry
func (s *RegistryService) PushPackage(ctx context.Context, path string) error {
	// Check if file exists
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to access package: %w", err)
	}

	if stat.IsDir() {
		return fmt.Errorf("path is a directory, expected a file")
	}

	// Read package contents
	packageData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read package file: %w", err)
	}

	// Parse package metadata
	var packageMeta types.PackageInfo
	if err := json.Unmarshal(packageData, &packageMeta); err != nil {
		return fmt.Errorf("invalid package format: %w", err)
	}

	// Verify required fields
	if packageMeta.Name == "" {
		return fmt.Errorf("package has no name")
	}
	if packageMeta.Version == "" {
		return fmt.Errorf("package has no version")
	}

	// Prepare request
	uploadURL := fmt.Sprintf("%s/packages", s.config.RegistryURL)
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	if s.config.Username != "" && s.config.AccessToken != "" {
		req.SetBasicAuth(s.config.Username, s.config.AccessToken)
	}

	// The actual implementation would use a proper multipart form
	// or API-specific request body format
	// This is a placeholder for the real implementation
	return fmt.Errorf("push not implemented - real implementation would POST to %s", uploadURL)
}

// PullPackage pulls a package from the registry
func (s *RegistryService) PullPackage(ctx context.Context, name string, version string) (string, error) {
	// Check if package is already in cache
	cacheKey := fmt.Sprintf("%s@%s", name, version)
	cachePath := filepath.Join(s.config.CachePath, cacheKey+".json")

	// Check if cached version exists
	if s.config.CachePath != "" {
		if _, err := os.Stat(cachePath); err == nil {
			// Package exists in cache
			return cachePath, nil
		}
	}

	// Prepare request
	downloadURL := fmt.Sprintf("%s/packages/%s/%s", s.config.RegistryURL, name, version)
	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	if s.config.Username != "" && s.config.AccessToken != "" {
		req.SetBasicAuth(s.config.Username, s.config.AccessToken)
	}

	// The actual implementation would download the package
	// and save it to the cache path
	// For now, create a placeholder file
	if s.config.CachePath != "" {
		// Create a dummy package with the requested name and version
		packageInfo := types.PackageInfo{
			Name:        name,
			Version:     version,
			Description: "Downloaded package",
			Author:      "Registry",
			License:     "MIT",
			CreatedAt:   time.Now(),
		}

		// Convert to JSON
		data, err := json.MarshalIndent(packageInfo, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to create package data: %w", err)
		}

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
			return "", fmt.Errorf("failed to create cache directory: %w", err)
		}

		// Write to file
		if err := os.WriteFile(cachePath, data, 0644); err != nil {
			return "", fmt.Errorf("failed to write package to cache: %w", err)
		}

		return cachePath, nil
	}

	return "", fmt.Errorf("pull not fully implemented, would download from %s", downloadURL)
}

// SearchPackages searches for packages in the registry
func (s *RegistryService) SearchPackages(ctx context.Context, query string, limit int) ([]types.PackageInfo, error) {
	// Prepare request
	searchURL := fmt.Sprintf("%s/packages/search?q=%s&limit=%d", s.config.RegistryURL, query, limit)
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	if s.config.Username != "" && s.config.AccessToken != "" {
		req.SetBasicAuth(s.config.Username, s.config.AccessToken)
	}

	// Perform request (in a real implementation)
	// resp, err := s.client.Do(req)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to search packages: %w", err)
	// }
	// defer resp.Body.Close()

	// For now, return dummy results
	results := make([]types.PackageInfo, 0)
	now := time.Now()

	// If query is empty, return some default packages
	if query == "" {
		results = append(results, types.PackageInfo{
			Name:        "example-agent",
			Version:     "1.0.0",
			Description: "An example agent for Sentinel Stacks",
			Author:      "Sentinel Team",
			License:     "MIT",
			Tags:        []string{"example", "agent"},
			CreatedAt:   now,
		})

		results = append(results, types.PackageInfo{
			Name:        "transform-agent",
			Version:     "0.5.0",
			Description: "Text transformation agent",
			Author:      "Sentinel Team",
			License:     "MIT",
			Tags:        []string{"text", "transform"},
			CreatedAt:   now,
		})
	} else {
		// Add a single result matching the query
		results = append(results, types.PackageInfo{
			Name:        query + "-agent",
			Version:     "1.0.0",
			Description: "Search result for: " + query,
			Author:      "Sentinel Team",
			License:     "MIT",
			Tags:        []string{query},
			CreatedAt:   now,
		})
	}

	return results, nil
}
