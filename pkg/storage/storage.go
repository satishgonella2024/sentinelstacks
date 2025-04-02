package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// Storage handles data persistence
type Storage struct {
	baseDir string
	mutex   sync.RWMutex
}

// NewStorage creates a new storage instance
func NewStorage(baseDir string) (*Storage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Create subdirectories
	for _, dir := range []string{"networks", "volumes", "systems"} {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			return nil, fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
	}

	return &Storage{
		baseDir: baseDir,
	}, nil
}

// SaveNetwork persists a network to storage
func (s *Storage) SaveNetwork(network types.Network) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Convert to JSON
	data, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal network data: %w", err)
	}

	// Write to file
	filePath := filepath.Join(s.baseDir, "networks", network.ID+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write network data: %w", err)
	}

	return nil
}

// GetNetwork retrieves a network from storage
func (s *Storage) GetNetwork(id string) (types.Network, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var network types.Network

	// Read file
	filePath := filepath.Join(s.baseDir, "networks", id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return network, fmt.Errorf("network not found: %s", id)
		}
		return network, fmt.Errorf("failed to read network data: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &network); err != nil {
		return network, fmt.Errorf("failed to unmarshal network data: %w", err)
	}

	return network, nil
}

// ListNetworks lists all networks
func (s *Storage) ListNetworks() ([]types.Network, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var networks []types.Network

	// Read directory
	networkDir := filepath.Join(s.baseDir, "networks")
	files, err := os.ReadDir(networkDir)
	if err != nil {
		if os.IsNotExist(err) {
			return networks, nil
		}
		return nil, fmt.Errorf("failed to read networks directory: %w", err)
	}

	// Load each network
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(networkDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files that can't be read
		}

		var network types.Network
		if err := json.Unmarshal(data, &network); err != nil {
			continue // Skip files that can't be parsed
		}

		networks = append(networks, network)
	}

	return networks, nil
}

// DeleteNetwork removes a network from storage
func (s *Storage) DeleteNetwork(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if network exists
	filePath := filepath.Join(s.baseDir, "networks", id+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("network not found: %s", id)
	}

	// Remove file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete network data: %w", err)
	}

	return nil
}

// GetNetworkByName finds a network by name
func (s *Storage) GetNetworkByName(name string) (types.Network, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// List all networks
	networks, err := s.ListNetworks()
	if err != nil {
		return types.Network{}, err
	}

	// Find network by name
	for _, network := range networks {
		if network.Name == name {
			return network, nil
		}
	}

	return types.Network{}, fmt.Errorf("network not found: %s", name)
}
