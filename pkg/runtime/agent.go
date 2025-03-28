package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/satishgonella2024/sentinelstacks/pkg/agentfile"
	"github.com/satishgonella2024/sentinelstacks/pkg/models"
	"github.com/satishgonella2024/sentinelstacks/pkg/tools"
)

// AgentRuntime manages the execution of an agent
type AgentRuntime struct {
	Agentfile     agentfile.Agentfile
	Adapter       models.ModelAdapter
	State         map[string]interface{}
	StatePath     string
	ModelEndpoint string             // For overriding the default endpoint
	ToolManager   *tools.ToolManager // Tool manager for agent tools
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

	// Set up model options
	if r.Agentfile.Model.Options == nil {
		r.Agentfile.Model.Options = make(map[string]interface{})
	}

	// Add endpoint to options if specified in agentfile
	if r.Agentfile.Model.Endpoint != "" {
		r.Agentfile.Model.Options["endpoint"] = r.Agentfile.Model.Endpoint
	}

	// Override the endpoint if specified in runtime
	if r.ModelEndpoint != "" {
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

	// Set up tools if specified in the agentfile
	if len(r.Agentfile.Tools) > 0 {
		// Get tool IDs
		toolIDs := make([]string, 0, len(r.Agentfile.Tools))
		for _, tool := range r.Agentfile.Tools {
			toolIDs = append(toolIDs, tool.ID)
		}

		// Create tool manager
		toolRegistry := tools.GetToolRegistry()
		toolManager, err := toolRegistry.CreateToolManager(toolIDs)
		if err != nil {
			return fmt.Errorf("failed to create tool manager: %w", err)
		}

		r.ToolManager = toolManager
	}

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

	// Process the response for any tool calls
	processedResponse, err := r.processToolCalls(response)
	if err != nil {
		fmt.Printf("Warning: Error processing tool calls: %v\n", err)
		// Continue with the original response
		processedResponse = response
	}

	// Update state with response
	conversation = append(conversation, map[string]string{
		"role":    "assistant",
		"content": processedResponse,
	})
	r.State["conversation"] = conversation

	// Update additional state metrics
	r.updateStateMetrics(input, processedResponse)

	// Save state if persistence is enabled
	if r.Agentfile.Memory.Persistence {
		if err := r.SaveState(); err != nil {
			fmt.Printf("Warning: Failed to save state: %v\n", err)
		}
	}

	return processedResponse, nil
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

	// Add tool information if available
	if r.ToolManager != nil {
		prompt += "\nYou have access to the following tools:\n"
		tools := r.ToolManager.ListTools()
		for _, tool := range tools {
			prompt += fmt.Sprintf("- %s: %s\n", tool.Name(), tool.Description())
		}

		// Add tool usage instructions
		prompt += "\nTo use a tool, respond with the following format:\n"
		prompt += "{{tool:tool_name,param1:value1,param2:value2}}\n"
		prompt += "For example: {{tool:calculator,operation:add,a:5,b:3}}\n"
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

// processToolCalls looks for tool invocations in the response and executes them
func (r *AgentRuntime) processToolCalls(response string) (string, error) {
	// If no tool manager, return the original response
	if r.ToolManager == nil {
		return response, nil
	}

	// Pattern for tool calls: {{tool:tool_name,param1:value1,param2:value2}}
	toolPattern := regexp.MustCompile(`\{\{tool:([^,}]+)(?:,([^}]+))?\}\}`)
	matches := toolPattern.FindAllStringSubmatch(response, -1)

	// If no matches, return the original response
	if len(matches) == 0 {
		return response, nil
	}

	// Process each tool call
	processedResponse := response
	for _, match := range matches {
		fullMatch := match[0]                   // The entire match
		toolName := strings.TrimSpace(match[1]) // The tool name

		// Parse parameters
		params := make(map[string]interface{})
		if len(match) > 2 && match[2] != "" {
			paramStr := match[2]
			paramPairs := strings.Split(paramStr, ",")

			for _, pair := range paramPairs {
				kv := strings.SplitN(pair, ":", 2)
				if len(kv) == 2 {
					key := strings.TrimSpace(kv[0])
					value := strings.TrimSpace(kv[1])

					// Try to convert to appropriate types
					if numVal, err := strconv.ParseFloat(value, 64); err == nil {
						params[key] = numVal
					} else if value == "true" {
						params[key] = true
					} else if value == "false" {
						params[key] = false
					} else {
						params[key] = value
					}
				}
			}
		}

		// Execute the tool
		result, err := r.ToolManager.ExecuteTool(toolName, params)
		if err != nil {
			// Replace the tool call with an error message
			errMessage := fmt.Sprintf("[Error executing tool '%s': %v]", toolName, err)
			processedResponse = strings.Replace(processedResponse, fullMatch, errMessage, 1)
		} else {
			// Replace the tool call with the result
			resultStr := fmt.Sprintf("%v", result)
			processedResponse = strings.Replace(processedResponse, fullMatch, resultStr, 1)
		}
	}

	return processedResponse, nil
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
