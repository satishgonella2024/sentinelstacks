package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sentinelstacks/sentinel/pkg/services"
	"github.com/sentinelstacks/sentinel/pkg/storage"
)

// App represents the main application
type App struct {
	NetworkService *services.NetworkService
	// Add other services as needed
}

// NewApp creates a new application instance
func NewApp(dataDir string) (*App, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Initialize storage
	storage, err := storage.NewStorage(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize services
	networkService := services.NewNetworkService(storage)

	return &App{
		NetworkService: networkService,
	}, nil
}

// DefaultDataDir returns the default data directory
func DefaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "sentinelstacks")
	}
	return filepath.Join(home, ".sentinel", "data")
}
