package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/satishgonella2024/sentinelstacks/pkg/agentfile"
	"github.com/satishgonella2024/sentinelstacks/pkg/models"
)

// AgentRuntime manages the execution of an agent
type AgentRuntime struct {
	Agentfile agentfile.Agentfile
	Adapter   models.ModelAdapter
	State     map[string]interface{}
	StatePath string
}

// NewAgentRuntime creates a new agent runtime
func NewAgentRuntime() *AgentRuntime {
	return &AgentRuntime{
		State: make(map[string]interface{}),
	}
}

// LoadAgentfile loads an agent configuration from a file
func (r *AgentRuntime) LoadAgentfile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read agentfile: %w", err)
	}

	var af agentfile.Agentfile
	err = yaml.Unmarshal(data, &af)
	if err != nil {
		return fmt.Errorf("failed to parse agentfile: %w", err)
	}

	r.Agentfile = af
	
	// Set up state path
	dir := filepath.Dir(path)
	baseFileName := filepath.Base(path)
	ext := filepath.Ext(baseFileName)
	baseName := baseFileName[:len(baseFileName)-len(ext)]
	r.StatePath = filepath.Join(dir, baseName+".state.json")

	return nil
}

// Initialize sets up the agent runtime
func (r *AgentRuntime) Initialize() error {
	// Set up model adapter based on Agentfile
	switch r.Agentfile.Model.Provider {
	case "ollama":
		r.Adapter = models.NewOllamaAdapter("http://model.gonella.co.uk", r.Agentfile.Model.Name)
	default:
		return fmt.Errorf("unsupported model provider: %s", r.Agentfile.Model.Provider)
	}

	// Load state if it exists
	if r.Agentfile.Memory.Persistence {
		if _, err := os.Stat(r.StatePath); err == nil {
			data, err := ioutil.ReadFile(r.StatePath)
			if err == nil {
				err = yaml.Unmarshal(data, &r.State)
				if err != nil {
					return fmt.Errorf("failed to parse state: %w", err)
				}
			}
		}
	}

	return nil
}

// Run processes user input and returns the agent's response
func (r *AgentRuntime) Run(input string) (string, error) {
	// Create system prompt based on agent capabilities
	systemPrompt := fmt.Sprintf("You are %s, an AI assistant with these capabilities: %v", 
		r.Agentfile.Name, r.Agentfile.Capabilities)

	// Set model options
	temperature := 0.7
	if temp, ok := r.Agentfile.Model.Options["temperature"].(float64); ok {
		temperature = temp
	}
	
	options := models.Options{
		Temperature: temperature,
	}

	// Track conversation in state
	if _, ok := r.State["conversation"]; !ok {
		r.State["conversation"] = []map[string]string{}
	}
	
	conversation := r.State["conversation"].([]map[string]string)
	conversation = append(conversation, map[string]string{"role": "user", "content": input})
	r.State["conversation"] = conversation

	// Generate response
	response, err := r.Adapter.Generate(input, systemPrompt, options)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	// Update state with response
	conversation = append(conversation.([]map[string]string), map[string]string{"role": "assistant", "content": response})
	r.State["conversation"] = conversation

	// Save state if persistence is enabled
	if r.Agentfile.Memory.Persistence {
		r.SaveState()
	}

	return response, nil
}

// SaveState persists the agent's state to disk
func (r *AgentRuntime) SaveState() error {
	data, err := yaml.Marshal(r.State)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	err = ioutil.WriteFile(r.StatePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	return nil
}
