package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	
	"github.com/google/uuid"
	"github.com/sentinelstacks/sentinel/pkg/models"
	"github.com/sentinelstacks/sentinel/pkg/repository"
)

// FSMultiAgentSystemRepository implements MultiAgentSystemRepository using the filesystem
type FSMultiAgentSystemRepository struct {
	dataDir string
	mutex   sync.RWMutex
}

// NewFSMultiAgentSystemRepository creates a new filesystem-based multi-agent system repository
func NewFSMultiAgentSystemRepository(dataDir string) repository.MultiAgentSystemRepository {
	systemsDir := filepath.Join(dataDir, "systems")
	os.MkdirAll(systemsDir, 0755)
	return &FSMultiAgentSystemRepository{dataDir: systemsDir}
}

// getFilePath returns the file path for a system
func (r *FSMultiAgentSystemRepository) getFilePath(id string) string {
	return filepath.Join(r.dataDir, id+".json")
}

// Create stores a new multi-agent system
func (r *FSMultiAgentSystemRepository) Create(ctx context.Context, system *models.MultiAgentSystem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Generate ID if not provided
	if system.ID == "" {
		system.ID = uuid.New().String()
	}
	
	// Set creation time if not set
	if system.CreatedAt.IsZero() {
		system.CreatedAt = time.Now()
	}
	
	// Initialize metadata if nil
	if system.Metadata == nil {
		system.Metadata = make(map[string]string)
	}
	
	// Initialize networks and volumes if nil
	if system.Networks == nil {
		system.Networks = []string{}
	}
	
	if system.Volumes == nil {
		system.Volumes = []string{}
	}
	
	// Check for existing system with same name
	existing, err := r.GetByName(ctx, system.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("multi-agent system with name '%s' already exists", system.Name)
	}
	
	// Serialize to JSON
	data, err := json.MarshalIndent(system, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal multi-agent system data: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(r.getFilePath(system.ID), data, 0644); err != nil {
		return fmt.Errorf("failed to write multi-agent system data: %w", err)
	}
	
	return nil
}

// Get retrieves a multi-agent system by ID
func (r *FSMultiAgentSystemRepository) Get(ctx context.Context, id string) (*models.MultiAgentSystem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	data, err := os.ReadFile(r.getFilePath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("multi-agent system not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read multi-agent system data: %w", err)
	}
	
	var system models.MultiAgentSystem
	if err := json.Unmarshal(data, &system); err != nil {
		return nil, fmt.Errorf("failed to unmarshal multi-agent system data: %w", err)
	}
	
	return &system, nil
}

// GetByName retrieves a multi-agent system by name
func (r *FSMultiAgentSystemRepository) GetByName(ctx context.Context, name string) (*models.MultiAgentSystem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	systems, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	
	for _, system := range systems {
		if system.Name == name {
			return system, nil
		}
	}
	
	return nil, fmt.Errorf("multi-agent system not found: %s", name)
}

// List returns all multi-agent systems
func (r *FSMultiAgentSystemRepository) List(ctx context.Context) ([]*models.MultiAgentSystem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	files, err := os.ReadDir(r.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read multi-agent systems directory: %w", err)
	}
	
	var systems []*models.MultiAgentSystem
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		data, err := os.ReadFile(filepath.Join(r.dataDir, file.Name()))
		if err != nil {
			continue
		}
		
		var system models.MultiAgentSystem
		if err := json.Unmarshal(data, &system); err != nil {
			continue
		}
		
		systems = append(systems, &system)
	}
	
	return systems, nil
}

// Update updates an existing multi-agent system
func (r *FSMultiAgentSystemRepository) Update(ctx context.Context, system *models.MultiAgentSystem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if system exists
	existing, err := r.Get(ctx, system.ID)
	if err != nil {
		return err
	}
	
	// Preserve creation time
	system.CreatedAt = existing.CreatedAt
	
	// Serialize to JSON
	data, err := json.MarshalIndent(system, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal multi-agent system data: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(r.getFilePath(system.ID), data, 0644); err != nil {
		return fmt.Errorf("failed to write multi-agent system data: %w", err)
	}
	
	return nil
}

// Delete removes a multi-agent system
func (r *FSMultiAgentSystemRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if system exists
	if _, err := r.Get(ctx, id); err != nil {
		return err
	}
	
	// Remove file
	if err := os.Remove(r.getFilePath(id)); err != nil {
		return fmt.Errorf("failed to delete multi-agent system data: %w", err)
	}
	
	return nil
}
