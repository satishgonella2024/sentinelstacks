package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/sentinelstacks/sentinel/internal/shim"
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

	// Read the config file
	data, err := os.ReadFile(r.configFile)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	// Unmarshal the data
	var agentInfos map[string]AgentInfo
	if err := json.Unmarshal(data, &agentInfos); err != nil {
		return fmt.Errorf("could not unmarshal agent data: %w", err)
	}

	// Convert AgentInfo to Agent
	for id, info := range agentInfos {
		// Create state directory if it doesn't exist
		stateDir := filepath.Join(r.dataDir, "agents", id)
		if err := os.MkdirAll(stateDir, 0755); err != nil {
			return fmt.Errorf("could not create agent state directory: %w", err)
		}

		// Create agent object
		agent := &Agent{
			ID:        id,
			Name:      info.Name,
			Image:     info.Image,
			Status:    AgentStatus(info.Status), // Convert string to AgentStatus
			CreatedAt: info.CreatedAt,
			Model:     info.Model,
			Memory:    info.Memory,
			APIUsage:  info.APIUsage,
			StateDir:  stateDir,
		}

		// Add agent to map
		r.agents[id] = agent
	}

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

	// Save agent configuration - method is called within lock context
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

	// Build command to run the agent
	cmd := exec.Command(os.Args[0], "run", "--non-interactive", agent.Image)

	// Set environment variables for the agent
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("SENTINEL_AGENT_ID=%s", agent.ID),
		fmt.Sprintf("SENTINEL_AGENT_NAME=%s", agent.Name),
		fmt.Sprintf("SENTINEL_AGENT_MODEL=%s", agent.Model),
	)

	// Create log file for the agent
	logFile, err := os.Create(filepath.Join(agent.StateDir, "agent.log"))
	if err != nil {
		return fmt.Errorf("could not create log file: %w", err)
	}

	// Set command output to log file
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Start the process
	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("could not start agent process: %w", err)
	}

	// Store process
	agent.Process = cmd.Process
	agent.Status = StatusRunning

	// Save agent configuration - called within lock context
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

	// Check if process exists
	if agent.Process == nil {
		// Just update status if process doesn't exist
		agent.Status = StatusStopped
		// Called within lock context
		return r.saveAgents()
	}

	// First try to send a SIGTERM signal
	if err := agent.Process.Signal(syscall.SIGTERM); err != nil {
		// If SIGTERM fails, try SIGKILL
		if err := agent.Process.Kill(); err != nil {
			return fmt.Errorf("could not kill agent process: %w", err)
		}
	}

	// Wait for a short time to let the process terminate gracefully
	done := make(chan error, 1)
	go func() {
		state, err := agent.Process.Wait()
		if err != nil {
			done <- err
			return
		}
		done <- nil
		fmt.Printf("Process exited with code: %d\n", state.ExitCode())
	}()

	// Wait for process to exit or timeout
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("error waiting for process: %w", err)
		}
	case <-time.After(5 * time.Second):
		// Process didn't exit in time, force kill
		if err := agent.Process.Kill(); err != nil {
			return fmt.Errorf("could not force kill agent process: %w", err)
		}
	}

	// Update agent state
	agent.Process = nil
	agent.Status = StatusStopped

	// Save agent configuration - called within lock context
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

	// Save agent configuration - called within lock context
	return r.saveAgents()
}

// saveAgents saves all agents to the config file
func (r *Runtime) saveAgents() error {
	// This method should be called within a lock context
	// to avoid deadlocks

	// Convert agents to AgentInfo for serialization
	agentInfos := make(map[string]AgentInfo)
	for id, agent := range r.agents {
		agentInfos[id] = AgentInfo{
			ID:        agent.ID,
			Name:      agent.Name,
			Image:     agent.Image,
			Status:    string(agent.Status),
			CreatedAt: agent.CreatedAt,
			Model:     agent.Model,
			Memory:    agent.Memory,
			APIUsage:  agent.APIUsage,
		}
	}

	// Marshal the data to JSON
	data, err := json.Marshal(agentInfos)
	if err != nil {
		return fmt.Errorf("could not marshal agent data: %w", err)
	}

	// Write the data to the config file
	if err := os.WriteFile(r.configFile, data, 0644); err != nil {
		return fmt.Errorf("could not write agent data to file: %w", err)
	}

	return nil
}

// CreateMultimodalAgent creates a new multimodal agent
func (r *Runtime) CreateMultimodalAgent(name, image, model, provider, apiKey, endpoint string) (*MultimodalAgent, error) {
	// Create regular agent first
	agent, err := r.CreateAgent(name, image, model)
	if err != nil {
		return nil, fmt.Errorf("failed to create base agent: %w", err)
	}

	// Create shim configuration
	shimConfig := shim.Config{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
		Endpoint: endpoint,
	}

	// Create multimodal agent
	return NewMultimodalAgent(agent, shimConfig)
}

// GetAgentLogs returns the logs for a specific agent
func (r *Runtime) GetAgentLogs(id string, tail int) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[id]
	if !exists {
		return "", fmt.Errorf("agent not found: %s", id)
	}

	// Get log file path
	logFile := filepath.Join(agent.StateDir, "agent.log")

	// Check if log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return "", fmt.Errorf("no logs found for agent: %s", id)
	}

	// Read log file
	data, err := os.ReadFile(logFile)
	if err != nil {
		return "", fmt.Errorf("could not read log file: %w", err)
	}

	// If tail is specified, return only the last N lines
	if tail > 0 {
		lines := strings.Split(string(data), "\n")
		if len(lines) > tail {
			lines = lines[len(lines)-tail:]
		}
		return strings.Join(lines, "\n"), nil
	}

	return string(data), nil
}

// UpdateAgentMetrics updates the metrics for a specific agent
func (r *Runtime) UpdateAgentMetrics(id string, apiCalls int, memoryUsage int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[id]
	if !exists {
		return fmt.Errorf("agent not found: %s", id)
	}

	// Update metrics
	agent.APIUsage += apiCalls
	agent.Memory = memoryUsage

	// Save agent configuration - called within lock context
	return r.saveAgents()
}

// GetAgentMetrics returns metrics for an agent
func (r *Runtime) GetAgentMetrics(id string) (map[string]interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", id)
	}

	// Basic metrics
	metrics := map[string]interface{}{
		"apiCalls":    agent.APIUsage,
		"memoryUsage": agent.Memory,
		"uptime":      time.Since(agent.CreatedAt).Seconds(),
		"status":      string(agent.Status),
	}

	// If the agent is running and has a process, try to get CPU usage
	if agent.Status == StatusRunning && agent.Process != nil {
		// This would be platform-specific in a real implementation
		// For now, just add a placeholder
		metrics["cpuUsage"] = 0.0
	}

	return metrics, nil
}
