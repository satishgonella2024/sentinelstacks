package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/satishgonella2024/sentinelstacks/pkg/storage"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// NetworkService provides operations for managing networks
type NetworkService struct {
	storage *storage.Storage
}

// NewNetworkService creates a new network service
func NewNetworkService(storage *storage.Storage) *NetworkService {
	return &NetworkService{
		storage: storage,
	}
}

// CreateNetwork creates a new network
func (s *NetworkService) CreateNetwork(name, driver string) (types.Network, error) {
	// Check if network with same name already exists
	_, err := s.storage.GetNetworkByName(name)
	if err == nil {
		return types.Network{}, fmt.Errorf("network with name '%s' already exists", name)
	}

	// Create new network
	network := types.Network{
		ID:        uuid.New().String(),
		Name:      name,
		Driver:    driver,
		CreatedAt: time.Now(),
		Status:    "active",
		Agents:    []string{},
	}

	// Save network
	if err := s.storage.SaveNetwork(network); err != nil {
		return types.Network{}, err
	}

	return network, nil
}

// ListNetworks returns all networks
func (s *NetworkService) ListNetworks() ([]types.Network, error) {
	return s.storage.ListNetworks()
}

// GetNetwork returns a network by ID
func (s *NetworkService) GetNetwork(id string) (types.Network, error) {
	return s.storage.GetNetwork(id)
}

// GetNetworkByName returns a network by name
func (s *NetworkService) GetNetworkByName(name string) (types.Network, error) {
	return s.storage.GetNetworkByName(name)
}

// DeleteNetwork removes a network
func (s *NetworkService) DeleteNetwork(id string) error {
	// Check if the network has connected agents
	network, err := s.storage.GetNetwork(id)
	if err != nil {
		return err
	}

	if len(network.Agents) > 0 {
		return fmt.Errorf("network has %d connected agents, cannot remove", len(network.Agents))
	}

	return s.storage.DeleteNetwork(id)
}

// DeleteNetworkForce forcibly removes a network
func (s *NetworkService) DeleteNetworkForce(id string) error {
	return s.storage.DeleteNetwork(id)
}

// ConnectAgent connects an agent to a network
func (s *NetworkService) ConnectAgent(networkName, agentID string) error {
	// Get network
	network, err := s.storage.GetNetworkByName(networkName)
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
	return s.storage.SaveNetwork(network)
}

// DisconnectAgent disconnects an agent from a network
func (s *NetworkService) DisconnectAgent(networkName, agentID string) error {
	// Get network
	network, err := s.storage.GetNetworkByName(networkName)
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
	return s.storage.SaveNetwork(network)
}
