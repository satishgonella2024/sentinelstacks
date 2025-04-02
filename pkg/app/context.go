package app

// Import repository implementations directly to avoid import cycles
import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// ServiceRegistry holds all application services
type ServiceRegistry struct {
	networkService NetworkService
	volumeService  VolumeService 
	composeService ComposeService
}

// Initialize the service registry
func NewServiceRegistry(dataDir string) *ServiceRegistry {
	// Ensure data directory exists
	os.MkdirAll(dataDir, 0755)
	
	// Create basic implementation for demonstration
	registry := &ServiceRegistry{
		networkService: &BasicNetworkService{basePath: filepath.Join(dataDir, "networks")},
		volumeService:  &BasicVolumeService{basePath: filepath.Join(dataDir, "volumes")},
		composeService: &BasicComposeService{basePath: filepath.Join(dataDir, "systems")},
	}
	
	// Create data directories
	os.MkdirAll(filepath.Join(dataDir, "networks"), 0755)
	os.MkdirAll(filepath.Join(dataDir, "volumes"), 0755)
	os.MkdirAll(filepath.Join(dataDir, "systems"), 0755)
	
	return registry
}

// Context keys
type contextKey int

const serviceRegistryKey contextKey = 0

// WithRegistry adds the service registry to the context
func WithRegistry(ctx context.Context, registry *ServiceRegistry) context.Context {
	return context.WithValue(ctx, serviceRegistryKey, registry)
}

// FromContext retrieves the service registry from the context
func FromContext(ctx context.Context) *ServiceRegistry {
	registry, ok := ctx.Value(serviceRegistryKey).(*ServiceRegistry)
	if !ok {
		// If no registry is found, create a default one
		fmt.Println("Warning: No service registry found in context, creating a default one")
		registry = NewServiceRegistry(filepath.Join(os.TempDir(), "sentinel"))
	}
	return registry
}

// Service accessor methods
func (r *ServiceRegistry) NetworkService() NetworkService {
	return r.networkService
}

func (r *ServiceRegistry) VolumeService() VolumeService {
	return r.volumeService
}

func (r *ServiceRegistry) ComposeService() ComposeService {
	return r.composeService
}

// Service interfaces 

// NetworkService defines the interface for network management
type NetworkService interface {
	CreateNetwork(ctx context.Context, name, driver string) (interface{}, error)
	GetNetwork(ctx context.Context, id string) (interface{}, error)
	GetNetworkByName(ctx context.Context, name string) (interface{}, error)
	ListNetworks(ctx context.Context) ([]interface{}, error)
	DeleteNetwork(ctx context.Context, id string) error
	ConnectAgent(ctx context.Context, networkName, agentID string) error
	DisconnectAgent(ctx context.Context, networkName, agentID string) error
	InspectNetwork(ctx context.Context, name string) (interface{}, error)
}

// VolumeService defines the interface for volume management
type VolumeService interface {
	CreateVolume(ctx context.Context, name, size string, encrypted bool) (interface{}, error)
	GetVolume(ctx context.Context, id string) (interface{}, error) 
	GetVolumeByName(ctx context.Context, name string) (interface{}, error)
	ListVolumes(ctx context.Context) ([]interface{}, error)
	DeleteVolume(ctx context.Context, id string) error
	MountVolume(ctx context.Context, volumeName, agentID, mountPath string) error
	UnmountVolume(ctx context.Context, volumeName, agentID string) error
	InspectVolume(ctx context.Context, name string) (interface{}, error)
}

// ComposeService defines the interface for multi-agent system management
type ComposeService interface {
	CreateSystem(ctx context.Context, name string, agents map[string]interface{}) (interface{}, error)
	GetSystem(ctx context.Context, id string) (interface{}, error)
	GetSystemByName(ctx context.Context, name string) (interface{}, error)
	ListSystems(ctx context.Context) ([]interface{}, error)
	DeleteSystem(ctx context.Context, id string) error
	UpdateSystemStatus(ctx context.Context, name, status string) error
	PauseSystem(ctx context.Context, name string) error
	ResumeSystem(ctx context.Context, name string) error
	StopSystem(ctx context.Context, name string) error
}

// Basic implementations for demonstration purposes

// BasicNetworkService provides a simple implementation of NetworkService
type BasicNetworkService struct {
	basePath string
}

func (s *BasicNetworkService) CreateNetwork(ctx context.Context, name, driver string) (interface{}, error) {
	return map[string]interface{}{
		"id": "net_" + name,
		"name": name,
		"driver": driver,
		"status": "active",
		"agents": []string{},
	}, nil
}

func (s *BasicNetworkService) GetNetwork(ctx context.Context, id string) (interface{}, error) {
	return map[string]interface{}{
		"id": id,
		"name": "network-" + id,
		"driver": "default",
		"status": "active",
		"agents": []string{},
	}, nil
}

func (s *BasicNetworkService) GetNetworkByName(ctx context.Context, name string) (interface{}, error) {
	return map[string]interface{}{
		"id": "net_" + name,
		"name": name,
		"driver": "default",
		"status": "active",
		"agents": []string{},
	}, nil
}

func (s *BasicNetworkService) ListNetworks(ctx context.Context) ([]interface{}, error) {
	return []interface{}{
		map[string]interface{}{
			"id": "net_default",
			"name": "default",
			"driver": "default",
			"status": "active",
			"agents": []string{},
		},
	}, nil
}

func (s *BasicNetworkService) DeleteNetwork(ctx context.Context, id string) error {
	return nil
}

func (s *BasicNetworkService) ConnectAgent(ctx context.Context, networkName, agentID string) error {
	return nil
}

func (s *BasicNetworkService) DisconnectAgent(ctx context.Context, networkName, agentID string) error {
	return nil
}

func (s *BasicNetworkService) InspectNetwork(ctx context.Context, name string) (interface{}, error) {
	return map[string]interface{}{
		"id": "net_" + name,
		"name": name,
		"driver": "default",
		"status": "active",
		"agents": []string{},
	}, nil
}

// BasicVolumeService provides a simple implementation of VolumeService
type BasicVolumeService struct {
	basePath string
}

func (s *BasicVolumeService) CreateVolume(ctx context.Context, name, size string, encrypted bool) (interface{}, error) {
	return map[string]interface{}{
		"id": "vol_" + name,
		"name": name,
		"size": size,
		"encrypted": encrypted,
		"used": "0",
	}, nil
}

func (s *BasicVolumeService) GetVolume(ctx context.Context, id string) (interface{}, error) {
	return map[string]interface{}{
		"id": id,
		"name": "volume-" + id,
		"size": "1GB",
		"encrypted": false,
		"used": "0",
	}, nil
}

func (s *BasicVolumeService) GetVolumeByName(ctx context.Context, name string) (interface{}, error) {
	return map[string]interface{}{
		"id": "vol_" + name,
		"name": name,
		"size": "1GB",
		"encrypted": false,
		"used": "0",
	}, nil
}

func (s *BasicVolumeService) ListVolumes(ctx context.Context) ([]interface{}, error) {
	return []interface{}{
		map[string]interface{}{
			"id": "vol_default",
			"name": "default",
			"size": "1GB",
			"encrypted": false,
			"used": "0",
		},
	}, nil
}

func (s *BasicVolumeService) DeleteVolume(ctx context.Context, id string) error {
	return nil
}

func (s *BasicVolumeService) MountVolume(ctx context.Context, volumeName, agentID, mountPath string) error {
	return nil
}

func (s *BasicVolumeService) UnmountVolume(ctx context.Context, volumeName, agentID string) error {
	return nil
}

func (s *BasicVolumeService) InspectVolume(ctx context.Context, name string) (interface{}, error) {
	return map[string]interface{}{
		"id": "vol_" + name,
		"name": name,
		"size": "1GB",
		"encrypted": false,
		"used": "0",
		"mountedBy": "",
		"mountPath": "",
	}, nil
}

// BasicComposeService provides a simple implementation of ComposeService
type BasicComposeService struct {
	basePath string
}

func (s *BasicComposeService) CreateSystem(ctx context.Context, name string, agents map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"id": "sys_" + name,
		"name": name,
		"status": "running",
		"agents": agents,
	}, nil
}

func (s *BasicComposeService) GetSystem(ctx context.Context, id string) (interface{}, error) {
	return map[string]interface{}{
		"id": id,
		"name": "system-" + id,
		"status": "running",
		"agents": map[string]interface{}{},
	}, nil
}

func (s *BasicComposeService) GetSystemByName(ctx context.Context, name string) (interface{}, error) {
	return map[string]interface{}{
		"id": "sys_" + name,
		"name": name,
		"status": "running",
		"agents": map[string]interface{}{},
	}, nil
}

func (s *BasicComposeService) ListSystems(ctx context.Context) ([]interface{}, error) {
	return []interface{}{
		map[string]interface{}{
			"id": "sys_example",
			"name": "example",
			"status": "running",
			"agents": map[string]interface{}{},
		},
	}, nil
}

func (s *BasicComposeService) DeleteSystem(ctx context.Context, id string) error {
	return nil
}

func (s *BasicComposeService) UpdateSystemStatus(ctx context.Context, name, status string) error {
	return nil
}

func (s *BasicComposeService) PauseSystem(ctx context.Context, name string) error {
	return nil
}

func (s *BasicComposeService) ResumeSystem(ctx context.Context, name string) error {
	return nil
}

func (s *BasicComposeService) StopSystem(ctx context.Context, name string) error {
	return nil
}
