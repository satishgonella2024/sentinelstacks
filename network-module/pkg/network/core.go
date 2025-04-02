/*
Package network implements Docker-inspired network management commands for SentinelStacks.
*/
package network

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Network represents a communication network for agents
type Network struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Driver          string                 `json:"driver"`
	CreatedAt       time.Time              `json:"created_at"`
	Status          string                 `json:"status"`
	Agents          []string               `json:"agents"`
	SupportedFormats []string              `json:"supported_formats,omitempty"`
	Config          map[string]interface{} `json:"config,omitempty"`
}

// NetworkManager handles network operations
type NetworkManager struct {
	DataDir string
	mutex   sync.RWMutex
}

// NewNetworkManager creates a new network manager
func NewNetworkManager(dataDir string) (*NetworkManager, error) {
	// Create networks directory if it doesn't exist
	networksDir := filepath.Join(dataDir, "networks")
	if err := os.MkdirAll(networksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create networks directory: %w", err)
	}

	return &NetworkManager{
		DataDir: dataDir,
	}, nil
}

// CreateNetwork creates a new network
func (m *NetworkManager) CreateNetwork(name, driver string, config map[string]interface{}) (Network, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if network with same name already exists
	exists, err := m.networkExists(name)
	if err != nil {
		return Network{}, fmt.Errorf("error checking if network exists: %w", err)
	}
	
	if exists {
		return Network{}, fmt.Errorf("network with name '%s' already exists", name)
	}

	// Set default supported formats if not specified in config
	formats := []string{"text"}
	if configFormats, ok := config["supported_formats"].([]string); ok && len(configFormats) > 0 {
		formats = configFormats
	}

	// Create new network
	network := Network{
		ID:              uuid.New().String(),
		Name:            name,
		Driver:          driver,
		CreatedAt:       time.Now(),
		Status:          "active",
		Agents:          []string{},
		SupportedFormats: formats,
		Config:          config,
	}

	// Save network
	if err := m.saveNetwork(network); err != nil {
		return Network{}, err
	}

	return network, nil
}

// networkExists checks if a network with the given name exists
func (m *NetworkManager) networkExists(name string) (bool, error) {
	networksDir := filepath.Join(m.DataDir, "networks")
	files, err := os.ReadDir(networksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(networksDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files that can't be read
		}

		var network Network
		if err := json.Unmarshal(data, &network); err != nil {
			continue // Skip files that can't be parsed
		}

		if network.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// ListNetworks returns all networks
func (m *NetworkManager) ListNetworks() ([]Network, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var networks []Network

	// Read networks directory
	networksDir := filepath.Join(m.DataDir, "networks")
	files, err := os.ReadDir(networksDir)
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

		filePath := filepath.Join(networksDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files that can't be read
		}

		var network Network
		if err := json.Unmarshal(data, &network); err != nil {
			continue // Skip files that can't be parsed
		}

		networks = append(networks, network)
	}

	return networks, nil
}

// GetNetwork returns a network by ID
func (m *NetworkManager) GetNetwork(id string) (Network, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	filePath := filepath.Join(m.DataDir, "networks", id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Network{}, fmt.Errorf("network not found: %s", id)
		}
		return Network{}, fmt.Errorf("failed to read network data: %w", err)
	}

	var network Network
	if err := json.Unmarshal(data, &network); err != nil {
		return Network{}, fmt.Errorf("failed to parse network data: %w", err)
	}

	return network, nil
}

// GetNetworkByName returns a network by name
func (m *NetworkManager) GetNetworkByName(name string) (Network, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Read networks directory
	networksDir := filepath.Join(m.DataDir, "networks")
	files, err := os.ReadDir(networksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return Network{}, fmt.Errorf("network not found: %s", name)
		}
		return Network{}, fmt.Errorf("failed to read networks directory: %w", err)
	}

	// Find network with matching name
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(networksDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files that can't be read
		}

		var network Network
		if err := json.Unmarshal(data, &network); err != nil {
			continue // Skip files that can't be parsed
		}

		if network.Name == name {
			return network, nil
		}
	}

	return Network{}, fmt.Errorf("network not found: %s", name)
}

// DeleteNetwork removes a network
func (m *NetworkManager) DeleteNetwork(id string, force bool) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Get network
	filePath := filepath.Join(m.DataDir, "networks", id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("network not found: %s", id)
		}
		return fmt.Errorf("failed to read network data: %w", err)
	}

	var network Network
	if err := json.Unmarshal(data, &network); err != nil {
		return fmt.Errorf("failed to parse network data: %w", err)
	}

	// Check if network has connected agents
	if len(network.Agents) > 0 && !force {
		return fmt.Errorf("network has %d connected agents, use --force to remove", len(network.Agents))
	}

	// Remove network file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove network file: %w", err)
	}

	return nil
}

// ConnectAgent connects an agent to a network
func (m *NetworkManager) ConnectAgent(networkName, agentID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Get network
	network, err := m.getNetworkByName(networkName)
	if err != nil {
		return err
	}

	// Check if agent is already connected
	for _, id := range network.Agents {
		if id == agentID {
			return fmt.Errorf("agent '%s' is already connected to network '%s'", agentID, networkName)
		}
	}

	// Add agent to network
	network.Agents = append(network.Agents, agentID)

	// Save network
	return m.saveNetwork(network)
}

// DisconnectAgent disconnects an agent from a network
func (m *NetworkManager) DisconnectAgent(networkName, agentID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Get network
	network, err := m.getNetworkByName(networkName)
	if err != nil {
		return err
	}

	// Check if agent is connected
	found := false
	var newAgents []string
	for _, id := range network.Agents {
		if id == agentID {
			found = true
		} else {
			newAgents = append(newAgents, id)
		}
	}

	if !found {
		return fmt.Errorf("agent '%s' is not connected to network '%s'", agentID, networkName)
	}

	// Update agents list
	network.Agents = newAgents

	// Save network
	return m.saveNetwork(network)
}

// UpdateNetwork updates network configuration
func (m *NetworkManager) UpdateNetwork(name string, config map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Get network
	network, err := m.getNetworkByName(name)
	if err != nil {
		return err
	}

	// Update config
	if network.Config == nil {
		network.Config = make(map[string]interface{})
	}
	
	for k, v := range config {
		network.Config[k] = v
	}

	// Update supported formats if provided
	if formats, ok := config["supported_formats"].([]string); ok {
		network.SupportedFormats = formats
	}

	// Save network
	return m.saveNetwork(network)
}

// getNetworkByName is an internal helper that doesn't need an external lock
// since it's called from methods that already have a lock
func (m *NetworkManager) getNetworkByName(name string) (Network, error) {
	// Read networks directory
	networksDir := filepath.Join(m.DataDir, "networks")
	files, err := os.ReadDir(networksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return Network{}, fmt.Errorf("network not found: %s", name)
		}
		return Network{}, fmt.Errorf("failed to read networks directory: %w", err)
	}

	// Find network with matching name
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(networksDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files that can't be read
		}

		var network Network
		if err := json.Unmarshal(data, &network); err != nil {
			continue // Skip files that can't be parsed
		}

		if network.Name == name {
			return network, nil
		}
	}

	return Network{}, fmt.Errorf("network not found: %s", name)
}

// saveNetwork saves a network to disk
func (m *NetworkManager) saveNetwork(network Network) error {
	// Convert to JSON
	data, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal network data: %w", err)
	}

	// Write to file
	filePath := filepath.Join(m.DataDir, "networks", network.ID+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write network data: %w", err)
	}

	return nil
}
