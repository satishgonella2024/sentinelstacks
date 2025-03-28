package runtime

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/satishgonella2024/sentinelstacks/pkg/agentfile"
	"github.com/satishgonella2024/sentinelstacks/pkg/models"
)

// AgentRuntime manages the execution of an agent
type AgentRuntime struct {
	Agentfile    agentfile.Agentfile
	Adapter      models.ModelAdapter
	State        map[string]interface{}
	StatePath    string
	ModelEndpoint string // For overriding the default endpoint
}

// NewAgentRuntime creates a new agent runtime
func NewAgentRuntime() *AgentRuntime {
	return &AgentRuntime{
		State: make(map[string]interface{}),
	}
}

// LoadAgentfile loads an agent configuration from a file
func (r *AgentRuntime) LoadAgentfile(path string) error {
	data, err := os.ReadFile(path)
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
	// Create model adapter factory
	factory := models.NewModelAdapterFactory()
	
	// Override the endpoint if specified
	if r.ModelEndpoint != "" {
		if r.Agentfile.Model.Options == nil {
			r.Agentfile.Model.Options = make(map[string]interface{})
		}
		r.Agentfile.Model.Options["endpoint"] = r.ModelEndpoint
	}
	
	// Create model adapter
	adapter, err := factory.CreateAdapter(
		r.Agentfile.Model.Provider,
		r.Agentfile.Model.Name,
		r.Agentfile.Model.Options,
	)
	if err != nil {
		return fmt.Errorf("failed to create model adapter: %w", err)
	}
	
	r.Adapter = adapter

	// Load state if it exists
	if r.Agentfile.Memory.Persistence {
		if _, err := os.Stat(r.StatePath); err == nil {
			data, err := os.ReadFile(r.StatePath)
			if err == nil {
				err = yaml.Unmarshal(data, &r.State)
				if err != nil {
					return fmt.Errorf("failed to parse state: %w", err)
				}
			}
		}
	}
	
	// Initialize state metrics
	if r.State == nil {
		r.State = make(map[string]interface{})
	}
	
	// Add initialization timestamp
	r.State["initialized_at"] = time.Now().Format(time.RFC3339)

	return nil
}

// Run processes user input and returns the agent's response
func (r *AgentRuntime) Run(input string) (string, error) {
	// Create a more detailed system prompt based on agent definition
	systemPrompt := r.buildSystemPrompt()
	
	// Update conversation history
	conversation := r.getConversationHistory()
	conversation = append(conversation, map[string]string{"role": "user", "content": input})
	
	// If using history, format the entire conversation for the model
	var prompt string
	if r.useConversationHistory() {
		prompt = r.formatConversationHistory(conversation)
	} else {
		prompt = input
	}
	
	// Set model options from agent configuration
	options := r.buildModelOptions()
	
	// Generate response
	response, err := r.Adapter.Generate(prompt, systemPrompt, options)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}
	
	// Update state with response
	conversation = append(conversation, map[string]string{
		"role": "assistant", 
		"content": response,
	})
	r.State["conversation"] = conversation
	
	// Update additional state metrics
	r.updateStateMetrics(input, response)
	
	// Save state if persistence is enabled
	if r.Agentfile.Memory.Persistence {
		if err := r.SaveState(); err != nil {
			fmt.Printf("Warning: Failed to save state: %v\n", err)
		}
	}
	
	return response, nil
}

// buildSystemPrompt creates a detailed system prompt based on the agent definition
func (r *AgentRuntime) buildSystemPrompt() string {
	// Start with the basic identity
	prompt := fmt.Sprintf("You are %s, an AI assistant. %s\n\n", 
		r.Agentfile.Name, r.Agentfile.Description)
	
	// Add capabilities information
	prompt += "Your capabilities include:\n"
	for _, capability := range r.Agentfile.Capabilities {
		prompt += fmt.Sprintf("- %s\n", capability)
	}
	
	// Add any constraints based on permissions
	prompt += "\nConstraints:\n"
	if len(r.Agentfile.Permissions.FileAccess) > 0 {
		prompt += fmt.Sprintf("- File access: %v\n", r.Agentfile.Permissions.FileAccess)
	} else {
		prompt += "- No file access\n"
	}
	if r.Agentfile.Permissions.Network {
		prompt += "- Network access allowed\n"
	} else {
		prompt += "- No network access\n"
	}
	
	return prompt
}

// getConversationHistory retrieves the conversation history from state
func (r *AgentRuntime) getConversationHistory() []map[string]string {
	if _, ok := r.State["conversation"]; !ok {
		r.State["conversation"] = []map[string]string{}
	}
	
	conversationVal := r.State["conversation"]
	conversation, ok := conversationVal.([]map[string]string)
	if !ok {
		return []map[string]string{}
	}
	
	return conversation
}

// useConversationHistory determines if conversation history should be used
func (r *AgentRuntime) useConversationHistory() bool {
	// Check if conversation capability is enabled
	for _, capability := range r.Agentfile.Capabilities {
		if capability == "conversation" {
			return true
		}
	}
	return false
}

// formatConversationHistory formats the conversation history for the model
func (r *AgentRuntime) formatConversationHistory(conversation []map[string]string) string {
	var formattedHistory strings.Builder
	
	// Use the most recent X messages to avoid context length issues
	maxHistoryLength := 10
	startIdx := 0
	if len(conversation) > maxHistoryLength {
		startIdx = len(conversation) - maxHistoryLength
	}
	
	for i := startIdx; i < len(conversation); i++ {
		msg := conversation[i]
		role := msg["role"]
		content := msg["content"]
		
		if role == "user" {
			formattedHistory.WriteString("User: " + content + "\n\n")
		} else if role == "assistant" {
			formattedHistory.WriteString("Assistant: " + content + "\n\n")
		}
	}
	
	return formattedHistory.String()
}

// buildModelOptions creates model options based on the agent configuration
func (r *AgentRuntime) buildModelOptions() models.Options {
	options := models.Options{}
	
	// Set temperature
	if temp, ok := r.Agentfile.Model.Options["temperature"].(float64); ok {
		options.Temperature = temp
	} else {
		options.Temperature = 0.7 // Default
	}
	
	// Set other options if defined
	if topP, ok := r.Agentfile.Model.Options["top_p"].(float64); ok {
		options.TopP = topP
	}
	
	if maxTokens, ok := r.Agentfile.Model.Options["max_tokens"].(float64); ok {
		options.MaxTokens = int(maxTokens)
	}
	
	return options
}

// updateStateMetrics updates state with metrics about the conversation
func (r *AgentRuntime) updateStateMetrics(input, response string) {
	// Update message count
	if count, ok := r.State["message_count"].(int); ok {
		r.State["message_count"] = count + 1
	} else {
		r.State["message_count"] = 1
	}
	
	// Update last active timestamp
	r.State["last_active"] = time.Now().Format(time.RFC3339)
	
	// Calculate and update response length stats
	if _, ok := r.State["response_lengths"]; !ok {
		r.State["response_lengths"] = []int{}
	}
	
	lengthsVal := r.State["response_lengths"]
	if lengths, ok := lengthsVal.([]int); ok {
		lengths = append(lengths, len(response))
		r.State["response_lengths"] = lengths
		
		// Calculate average
		total := 0
		for _, length := range lengths {
			total += length
		}
		r.State["avg_response_length"] = total / len(lengths)
	}
}

// SaveState persists the agent's state to disk
func (r *AgentRuntime) SaveState() error {
	data, err := yaml.Marshal(r.State)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	err = os.WriteFile(r.StatePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	return nil
}
