package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/satishgonella2024/sentinelstacks/internal/memory"
	"github.com/satishgonella2024/sentinelstacks/pkg/models"
)

var (
	agentsDir = filepath.Join(os.Getenv("HOME"), ".sentinel", "agents")
)

// Agent represents an AI agent instance
type Agent struct {
	Name         string
	Version      string
	Description  string
	Config       AgentConfig
	Memory       memory.Memory
	ModelAdapter models.ModelAdapter
}

// AgentConfig defines the configuration for an agent
type AgentConfig struct {
	Name         string              `yaml:"name"`
	Version      string              `yaml:"version"`
	Description  string              `yaml:"description"`
	Model        ModelConfig         `yaml:"model"`
	Capabilities []string            `yaml:"capabilities"`
	Memory       memory.MemoryConfig `yaml:"memory"`
	Tools        []ToolConfig        `yaml:"tools,omitempty"`
	Permissions  PermissionsConfig   `yaml:"permissions,omitempty"`
}

// ModelConfig defines the model settings
type ModelConfig struct {
	Provider string                 `yaml:"provider"`
	Name     string                 `yaml:"name"`
	Options  map[string]interface{} `yaml:"options,omitempty"`
}

// ToolConfig defines a tool that the agent can use
type ToolConfig struct {
	ID      string `yaml:"id"`
	Version string `yaml:"version"`
}

// PermissionsConfig defines what the agent is allowed to do
type PermissionsConfig struct {
	FileAccess []string `yaml:"file_access,omitempty"`
	Network    bool     `yaml:"network"`
}

func init() {
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		fmt.Printf("Error creating agents directory: %v\n", err)
	}
}

// NewAgent creates a new agent from a configuration
func NewAgent(config AgentConfig) (*Agent, error) {
	// Create memory system
	mem, err := memory.NewMemory(config.Name, config.Memory)
	if err != nil {
		return nil, fmt.Errorf("error creating memory: %w", err)
	}

	// Create model adapter
	factory := models.NewModelAdapterFactory()
	adapter, err := factory.CreateAdapter(config.Model.Provider, config.Model.Name, config.Model.Options)
	if err != nil {
		return nil, fmt.Errorf("error creating model adapter: %w", err)
	}

	agent := &Agent{
		Name:         config.Name,
		Version:      config.Version,
		Description:  config.Description,
		Config:       config,
		Memory:       mem,
		ModelAdapter: adapter,
	}

	return agent, nil
}

// LoadAgent loads an agent from the filesystem
func LoadAgent(name, version string) (*Agent, error) {
	if version == "" || version == "latest" {
		version = "latest"
	}

	agentDir := filepath.Join(agentsDir, name, version)

	// Check if agent exists
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("agent %s:%s not found", name, version)
	}

	// Load configuration
	configPath := filepath.Join(agentDir, "Agentfile")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try alternative names
		altNames := []string{"agent.yaml", "agent.yml", "config.yaml", "config.yml"}
		found := false
		for _, altName := range altNames {
			configPath = filepath.Join(agentDir, altName)
			if _, err := os.Stat(configPath); err == nil {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("no agent configuration file found in %s", agentDir)
		}
	}

	// Read configuration
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var config AgentConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	// Set defaults if not provided
	if config.Memory.Type == "" {
		config.Memory = memory.DefaultConfig()
	}

	return NewAgent(config)
}

// Run executes the agent
func Run(agentName, version string) error {
	// Load agent
	a, err := LoadAgent(agentName, version)
	if err != nil {
		return fmt.Errorf("error loading agent: %w", err)
	}

	fmt.Printf("Running agent %s:%s\n", a.Name, a.Version)

	// Here you would:
	// 1. Set up the environment
	// 2. Connect to the model
	// 3. Execute the agent's main loop
	// 4. Handle the results

	return nil
}

// Execute runs the agent with specific input
func (a *Agent) Execute(input string) (string, error) {
	// In a real implementation, this would:
	// 1. Retrieve relevant memories
	// 2. Format system prompt
	// 3. Send to the model
	// 4. Process the response
	// 5. Store in memory if appropriate

	// For now, we'll just make a simple call to the model
	systemPrompt := fmt.Sprintf("You are %s, a helpful AI assistant. %s", a.Name, a.Description)
	response, err := a.ModelAdapter.Generate(input, systemPrompt, models.Options{
		Temperature: 0.7,
	})

	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	// Store interaction in memory
	_, err = a.Memory.Add(response, map[string]interface{}{
		"input": input,
		"type":  "response",
	})
	if err != nil {
		// Log but don't fail
		fmt.Printf("Warning: Failed to store in memory: %v\n", err)
	}

	return response, nil
}
