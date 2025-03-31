package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Registry is the interface for all registry implementations
type Registry interface {
	// List returns all images in the registry
	List() ([]Image, error)

	// Get returns an image by name and tag
	Get(name, tag string) (*Image, error)

	// Push pushes an image to the registry
	Push(image *Image) error

	// Pull pulls an image from the registry
	Pull(name, tag string) (*Image, error)

	// Delete removes an image from the registry
	Delete(name, tag string) error
}

// Image represents a Sentinel Image
type Image struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Tag        string                 `json:"tag"`
	Data       []byte                 `json:"data,omitempty"`
	Config     map[string]interface{} `json:"config"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	BaseModel  string                 `json:"baseModel"`
	CreatedAt  string                 `json:"createdAt"`
	Size       int64                  `json:"size"`
}

// LocalRegistry is a registry that stores images locally
type LocalRegistry struct {
	dataDir string
	mu      sync.RWMutex
}

var defaultRegistry *LocalRegistry
var once sync.Once

// GetLocalRegistry returns the default local registry
func GetLocalRegistry() (*LocalRegistry, error) {
	var initError error
	once.Do(func() {
		defaultRegistry, initError = NewLocalRegistry("")
	})

	return defaultRegistry, initError
}

// NewLocalRegistry creates a new local registry
func NewLocalRegistry(dataDir string) (*LocalRegistry, error) {
	// If data directory not specified, use default
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not get home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".sentinel", "registry")
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create registry directory: %w", err)
	}

	return &LocalRegistry{
		dataDir: dataDir,
	}, nil
}

// List returns all images in the registry
func (r *LocalRegistry) List() ([]Image, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// In a real implementation, this would list images from the registry
	// For now, return an empty list
	return []Image{}, nil
}

// Get returns an image by name and tag
func (r *LocalRegistry) Get(name, tag string) (*Image, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// In a real implementation, this would get the image from the registry
	// For now, return an error
	return nil, fmt.Errorf("image not found: %s:%s", name, tag)
}

// Push pushes an image to the registry
func (r *LocalRegistry) Push(image *Image) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// In a real implementation, this would push the image to the registry
	// For now, return nil
	return nil
}

// Pull pulls an image from the registry
func (r *LocalRegistry) Pull(name, tag string) (*Image, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// In a real implementation, this would pull the image from the registry
	// For now, return an error
	return nil, fmt.Errorf("image not found: %s:%s", name, tag)
}

// Delete removes an image from the registry
func (r *LocalRegistry) Delete(name, tag string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// In a real implementation, this would delete the image from the registry
	// For now, return nil
	return nil
}
