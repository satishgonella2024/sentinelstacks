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
	"github.com/satishgonella2024/sentinelstacks/pkg/models"
	"github.com/satishgonella2024/sentinelstacks/pkg/repository"
)

// FSNetworkRepository implements NetworkRepository using the filesystem
type FSNetworkRepository struct {
	dataDir string
	mutex   sync.RWMutex
}

// NewFSNetworkRepository creates a new filesystem-based network repository
func NewFSNetworkRepository(dataDir string) repository.NetworkRepository {
	networksDir := filepath.Join(dataDir, "networks")
	os.MkdirAll(networksDir, 0755)
	return &FSNetworkRepository{dataDir: networksDir}
}

// getFilePath returns the file path for a network
func (r *FSNetworkRepository) getFilePath(id string) string {
	return filepath.Join(r.dataDir, id+".json")
}

// Create stores a new network
func (r *FSNetworkRepository) Create(ctx context.Context, network *models.Network) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Generate ID if not provided
	if network.ID == "" {
		network.ID = uuid.New().String()
	}
	
	// Set creation time if not set
	if network.CreatedAt.IsZero() {
		network.CreatedAt = time.Now()
	}
	
	// Initialize metadata if nil
	if network.Metadata == nil {
		network.Metadata = make(map[string]string)
	}
	
	// Initialize agents slice if nil
	if network.Agents == nil {
		network.Agents = []string{}
	}
	
	// Check for existing network with same name
	existing, err := r.GetByName(ctx, network.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("network with name '%s' already exists", network.Name)
	}
	
	// Serialize to JSON
	data, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal network data: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(r.getFilePath(network.ID), data, 0644); err != nil {
		return fmt.Errorf("failed to write network data: %w", err)
	}
	
	return nil
}

// Get retrieves a network by ID
func (r *FSNetworkRepository) Get(ctx context.Context, id string) (*models.Network, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	data, err := os.ReadFile(r.getFilePath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("network not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read network data: %w", err)
	}
	
	var network models.Network
	if err := json.Unmarshal(data, &network); err != nil {
		return nil, fmt.Errorf("failed to unmarshal network data: %w", err)
	}
	
	return &network, nil
}

// GetByName retrieves a network by name
func (r *FSNetworkRepository) GetByName(ctx context.Context, name string) (*models.Network, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	networks, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	
	for _, network := range networks {
		if network.Name == name {
			return network, nil
		}
	}
	
	return nil, fmt.Errorf("network not found: %s", name)
}

// List returns all networks
func (r *FSNetworkRepository) List(ctx context.Context) ([]*models.Network, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	files, err := os.ReadDir(r.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read networks directory: %w", err)
	}
	
	var networks []*models.Network
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		data, err := os.ReadFile(filepath.Join(r.dataDir, file.Name()))
		if err != nil {
			continue
		}
		
		var network models.Network
		if err := json.Unmarshal(data, &network); err != nil {
			continue
		}
		
		networks = append(networks, &network)
	}
	
	return networks, nil
}

// Update updates an existing network
func (r *FSNetworkRepository) Update(ctx context.Context, network *models.Network) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if network exists
	existing, err := r.Get(ctx, network.ID)
	if err != nil {
		return err
	}
	
	// Preserve creation time
	network.CreatedAt = existing.CreatedAt
	
	// Serialize to JSON
	data, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal network data: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(r.getFilePath(network.ID), data, 0644); err != nil {
		return fmt.Errorf("failed to write network data: %w", err)
	}
	
	return nil
}

// Delete removes a network
func (r *FSNetworkRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if network exists
	if _, err := r.Get(ctx, id); err != nil {
		return err
	}
	
	// Remove file
	if err := os.Remove(r.getFilePath(id)); err != nil {
		return fmt.Errorf("failed to delete network data: %w", err)
	}
	
	return nil
}

// ConnectAgent connects an agent to a network
func (r *FSNetworkRepository) ConnectAgent(ctx context.Context, networkID, agentID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Get network
	network, err := r.Get(ctx, networkID)
	if err != nil {
		return err
	}
	
	// Check if agent is already connected
	for _, id := range network.Agents {
		if id == agentID {
			return nil // Already connected
		}
	}
	
	// Add agent to network
	network.Agents = append(network.Agents, agentID)
	
	// Update network
	return r.Update(ctx, network)
}

// DisconnectAgent disconnects an agent from a network
func (r *FSNetworkRepository) DisconnectAgent(ctx context.Context, networkID, agentID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Get network
	network, err := r.Get(ctx, networkID)
	if err != nil {
		return err
	}
	
	// Remove agent from network
	var updatedAgents []string
	for _, id := range network.Agents {
		if id != agentID {
			updatedAgents = append(updatedAgents, id)
		}
	}
	network.Agents = updatedAgents
	
	// Update network
	return r.Update(ctx, network)
}
