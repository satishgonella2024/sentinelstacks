package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AgentStatus represents the status of an agent
type AgentStatus string

const (
	// StatusCreating indicates the agent is being created
	StatusCreating AgentStatus = "creating"
	// StatusRunning indicates the agent is running
	StatusRunning AgentStatus = "running"
	// StatusStopped indicates the agent has stopped
	StatusStopped AgentStatus = "stopped"
	// StatusError indicates the agent is in an error state
	StatusError AgentStatus = "error"
	// StatusPaused indicates the agent is paused
	StatusPaused AgentStatus = "paused"
)

// AgentInfo contains information about a running agent
type AgentInfo struct {
	ID        string    `json:"id"`        // Unique identifier for the agent
	Name      string    `json:"name"`      // Name of the agent
	Image     string    `json:"image"`     // Image used to create the agent
	Status    string    `json:"status"`    // Current status of the agent
	CreatedAt time.Time `json:"createdAt"` // When the agent was created
	Model     string    `json:"model"`     // LLM model being used
	Memory    int64     `json:"memory"`    // Memory usage in bytes
	APIUsage  int       `json:"apiUsage"`  // Number of API calls made
}

// Runtime manages agent execution
type Runtime struct {
	agents     map[string]*Agent
	dataDir    string
	configFile string
	mu         sync.RWMutex
}

// Agent represents a running agent instance
type Agent struct {
	ID        string
	Name      string
	Image     string
	Status    AgentStatus
	CreatedAt time.Time
	Model     string
	Memory    int64
	APIUsage  int
	Process   *os.Process
	StateDir  string
}

// defaultRuntime is the singleton runtime instance
var defaultRuntime *Runtime
var once sync.Once

// GetRuntime returns the default runtime instance
func GetRuntime() (*Runtime, error) {
	var initError error
	once.Do(func() {
		defaultRuntime, initError = NewRuntime("")
	})

	return defaultRuntime, initError
}

// NewRuntime creates a new runtime instance
func NewRuntime(dataDir string) (*Runtime, error) {
	// If data directory not specified, use default
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not get home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".sentinel")
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create data directory: %w", err)
	}

	// Create agents directory if it doesn't exist
	agentsDir := filepath.Join(dataDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create agents directory: %w", err)
	}

	configFile := filepath.Join(dataDir, "agents.json")

	runtime := &Runtime{
		agents:     make(map[string]*Agent),
		dataDir:    dataDir,
		configFile: configFile,
	}

	// Load existing agents
	if err := runtime.loadAgents(); err != nil {
		return nil, fmt.Errorf("could not load agents: %w", err)
	}

	return runtime, nil
}

// loadAgents loads agent information from the config file
func (r *Runtime) loadAgents() error {
	// If config file doesn't exist, return without error
	if _, err := os.Stat(r.configFile); os.IsNotExist(err) {
		return nil
	}

	// For now, simply return without loading agents
	// In a real implementation, this would load agent data from the config file
	return nil
}

// CreateAgent creates a new agent from an image
func (r *Runtime) CreateAgent(name, image, model string) (*Agent, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a new agent ID
	id := uuid.New().String()

	// Create agent state directory
	stateDir := filepath.Join(r.dataDir, "agents", id)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create agent state directory: %w", err)
	}

	// Create agent object
	agent := &Agent{
		ID:        id,
		Name:      name,
		Image:     image,
		Status:    StatusCreating,
		CreatedAt: time.Now(),
		Model:     model,
		StateDir:  stateDir,
	}

	// Save agent to runtime
	r.agents[id] = agent

	// Save agent configuration
	if err := r.saveAgents(); err != nil {
		return nil, fmt.Errorf("could not save agent: %w", err)
	}

	return agent, nil
}

// StartAgent starts a previously created agent
func (r *Runtime) StartAgent(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[id]
	if !exists {
		return fmt.Errorf("agent not found: %s", id)
	}

	if agent.Status == StatusRunning {
		return fmt.Errorf("agent already running: %s", id)
	}

	// In a real implementation, this would start the agent process
	// For now, just update the status
	agent.Status = StatusRunning

	// Save agent configuration
	return r.saveAgents()
}

// StopAgent stops a running agent
func (r *Runtime) StopAgent(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[id]
	if !exists {
		return fmt.Errorf("agent not found: %s", id)
	}

	if agent.Status != StatusRunning {
		return fmt.Errorf("agent not running: %s", id)
	}

	// In a real implementation, this would stop the agent process
	// For now, just update the status
	agent.Status = StatusStopped

	// Save agent configuration
	return r.saveAgents()
}

// GetAgent returns information about a specific agent
func (r *Runtime) GetAgent(id string) (AgentInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[id]
	if !exists {
		return AgentInfo{}, fmt.Errorf("agent not found: %s", id)
	}

	return AgentInfo{
		ID:        agent.ID,
		Name:      agent.Name,
		Image:     agent.Image,
		Status:    string(agent.Status),
		CreatedAt: agent.CreatedAt,
		Model:     agent.Model,
		Memory:    agent.Memory,
		APIUsage:  agent.APIUsage,
	}, nil
}

// GetRunningAgents returns information about all running agents
func (r *Runtime) GetRunningAgents() ([]AgentInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var agents []AgentInfo
	for _, agent := range r.agents {
		agents = append(agents, AgentInfo{
			ID:        agent.ID,
			Name:      agent.Name,
			Image:     agent.Image,
			Status:    string(agent.Status),
			CreatedAt: agent.CreatedAt,
			Model:     agent.Model,
			Memory:    agent.Memory,
			APIUsage:  agent.APIUsage,
		})
	}

	return agents, nil
}

// DeleteAgent removes an agent
func (r *Runtime) DeleteAgent(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[id]
	if !exists {
		return fmt.Errorf("agent not found: %s", id)
	}

	if agent.Status == StatusRunning {
		return fmt.Errorf("cannot delete running agent: %s", id)
	}

	// Remove agent from map
	delete(r.agents, id)

	// Remove agent state directory
	if err := os.RemoveAll(agent.StateDir); err != nil {
		return fmt.Errorf("could not remove agent state directory: %w", err)
	}

	// Save agent configuration
	return r.saveAgents()
}

// saveAgents saves all agents to the config file
func (r *Runtime) saveAgents() error {
	// In a real implementation, this would serialize and save agent data
	// For now, simply return success
	return nil
}
