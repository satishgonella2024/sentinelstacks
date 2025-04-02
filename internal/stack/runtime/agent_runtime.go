package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/stack"
)

// AgentRuntime manages the execution of agents within a stack
type AgentRuntime interface {
	// Execute runs an agent with the given inputs and returns its outputs
	Execute(ctx context.Context, agentSpec stack.StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error)
	
	// Cleanup performs any necessary cleanup after execution
	Cleanup() error
}

// DirectAgentRuntime runs agents directly using the LLM provider
type DirectAgentRuntime struct {
	workDir      string
	agentManager *agent.Manager
	logToConsole bool
}

// NewDirectAgentRuntime creates a new direct agent runtime
func NewDirectAgentRuntime(logToConsole bool) (*DirectAgentRuntime, error) {
	// Create temporary working directory
	workDir, err := ioutil.TempDir("", "sentinel-stack-")
	if err != nil {
		return nil, fmt.Errorf("failed to create working directory: %w", err)
	}
	
	// Create agent manager
	agentManager, err := agent.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create agent manager: %w", err)
	}

	return &DirectAgentRuntime{
		workDir:      workDir,
		agentManager: agentManager,
		logToConsole: logToConsole,
	}, nil
}

// Execute runs an agent directly and returns its outputs
func (r *DirectAgentRuntime) Execute(ctx context.Context, agentSpec stack.StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error) {
	// Create a unique directory for this agent run
	agentDir := filepath.Join(r.workDir, agentSpec.ID)
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create agent directory: %w", err)
	}
	
	// Log execution start if enabled
	if r.logToConsole {
		log.Printf("Starting agent execution: %s (uses: %s)", agentSpec.ID, agentSpec.Uses)
		log.Printf("Inputs: %+v", inputs)
	}

	// 1. Resolve the agent reference (could be an image name, local path, etc.)
	agentImage, err := r.agentManager.ResolveAgent(agentSpec.Uses)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve agent: %w", err)
	}
	
	// 2. Load the agent definition
	agentDef := agentImage.Definition
	
	// 3. Configure the LLM provider based on agent definition
	provider, err := shim.GetProvider(agentDef.BaseModel)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider for model %s: %w", agentDef.BaseModel, err)
	}
	
	// 4. Prepare agent inputs by combining stack inputs with agent params
	agentInputs := make(map[string]interface{})
	
	// Add agent parameters
	if agentSpec.Params != nil {
		for k, v := range agentSpec.Params {
			agentInputs[k] = v
		}
	}
	
	// Add stack inputs
	for k, v := range inputs {
		// Don't override existing params with stack inputs if keys conflict
		if _, exists := agentInputs[k]; !exists {
			agentInputs[k] = v
		}
	}
	
	// 5. Format the prompt using the agent's prompt template
	formattedPrompt, err := formatPrompt(agentDef.PromptTemplate, agentInputs)
	if err != nil {
		return nil, fmt.Errorf("failed to format prompt: %w", err)
	}
	
	// 6. Prepare the system prompt
	systemPrompt := agentDef.SystemPrompt
	
	// Log prompts if enabled
	if r.logToConsole {
		log.Printf("System prompt: %s", systemPrompt)
		log.Printf("Formatted prompt: %s", formattedPrompt)
	}
	
	// 7. Execute the LLM call
	startTime := time.Now()
	response, err := provider.Complete(ctx, shim.CompletionRequest{
		Model:        agentDef.BaseModel,
		SystemPrompt: systemPrompt,
		UserPrompt:   formattedPrompt,
		MaxTokens:    agentDef.MaxTokens,
		Temperature:  agentDef.Temperature,
	})
	
	duration := time.Since(startTime)
	
	if err != nil {
		return nil, fmt.Errorf("LLM completion failed: %w", err)
	}
	
	// 8. Process the response
	if r.logToConsole {
		log.Printf("Agent execution completed in %v", duration)
		log.Printf("Response: %s", response.Text)
	}
	
	// 9. Parse the response based on expected output format
	outputs, err := parseResponse(response.Text, agentDef.OutputFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Add metadata to outputs
	outputs["_agent_id"] = agentSpec.ID
	outputs["_agent_name"] = agentDef.Name
	outputs["_execution_time"] = duration.String()
	outputs["_executed_at"] = time.Now().Format(time.RFC3339)
	outputs["_model"] = agentDef.BaseModel
	
	// Write outputs to a file for debugging/audit
	outputsFile := filepath.Join(agentDir, "outputs.json")
	outputsJSON, err := json.MarshalIndent(outputs, "", "  ")
	if err == nil {
		if err := ioutil.WriteFile(outputsFile, outputsJSON, 0644); err != nil && r.logToConsole {
			log.Printf("Warning: Failed to write outputs file: %v", err)
		}
	}
	
	return outputs, nil
}

// Cleanup removes temporary files
func (r *DirectAgentRuntime) Cleanup() error {
	if r.workDir != "" {
		return os.RemoveAll(r.workDir)
	}
	return nil
}

// formatPrompt replaces placeholders in the prompt template with values from inputs
func formatPrompt(template string, inputs map[string]interface{}) (string, error) {
	result := template
	
	// Replace placeholders like {{key}} with values from inputs
	for key, value := range inputs {
		placeholder := fmt.Sprintf("{{%s}}", key)
		
		// Convert value to string based on its type
		var stringValue string
		switch v := value.(type) {
		case string:
			stringValue = v
		case []byte:
			stringValue = string(v)
		default:
			// For other types, use JSON marshaling
			bytes, err := json.Marshal(value)
			if err != nil {
				return "", fmt.Errorf("failed to marshal value for placeholder %s: %w", key, err)
			}
			stringValue = string(bytes)
		}
		
		result = strings.ReplaceAll(result, placeholder, stringValue)
	}
	
	return result, nil
}

// parseResponse converts the LLM response to structured outputs based on the format
func parseResponse(responseText, outputFormat string) (map[string]interface{}, error) {
	outputs := make(map[string]interface{})
	
	switch strings.ToLower(outputFormat) {
	case "json":
		// Try to parse response as JSON
		if err := json.Unmarshal([]byte(responseText), &outputs); err != nil {
			// If parsing fails, try to extract JSON from markdown code blocks
			jsonText := extractJsonFromMarkdown(responseText)
			if jsonText != "" {
				if err := json.Unmarshal([]byte(jsonText), &outputs); err != nil {
					return nil, fmt.Errorf("failed to parse JSON response: %w", err)
				}
			} else {
				return nil, fmt.Errorf("response is not valid JSON and no JSON block found: %w", err)
			}
		}
		
	case "text":
		// Store as plain text
		outputs["text"] = responseText
		
	case "key_value":
		// Parse as key-value pairs (one per line)
		lines := strings.Split(responseText, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				outputs[key] = value
			}
		}
		
	default:
		// Default to storing as text
		outputs["text"] = responseText
	}
	
	return outputs, nil
}

// extractJsonFromMarkdown tries to extract JSON from markdown code blocks
func extractJsonFromMarkdown(text string) string {
	// Look for ```json ... ``` blocks
	jsonBlockRegex := "```json\\s*\\n([\\s\\S]*?)\\n\\s*```"
	matches := regexp.MustCompile(jsonBlockRegex).FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	
	// Look for generic code blocks that might contain JSON
	codeBlockRegex := "```\\s*\\n([\\s\\S]*?)\\n\\s*```"
	matches = regexp.MustCompile(codeBlockRegex).FindStringSubmatch(text)
	if len(matches) > 1 {
		// Check if content is valid JSON
		jsonContent := matches[1]
		var testJSON map[string]interface{}
		if json.Unmarshal([]byte(jsonContent), &testJSON) == nil {
			return jsonContent
		}
	}
	
	return ""
}

// CliAgentRuntime executes agents using the sentinel CLI
type CliAgentRuntime struct {
	workDir      string
	sentinelPath string
	logToConsole bool
}

// NewCliAgentRuntime creates a new CLI-based agent runtime
func NewCliAgentRuntime(logToConsole bool) (*CliAgentRuntime, error) {
	// Create temporary working directory
	workDir, err := ioutil.TempDir("", "sentinel-stack-cli-")
	if err != nil {
		return nil, fmt.Errorf("failed to create working directory: %w", err)
	}

	// Find sentinel executable path
	sentinelPath, err := findSentinelPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find sentinel executable: %w", err)
	}

	return &CliAgentRuntime{
		workDir:      workDir,
		sentinelPath: sentinelPath,
		logToConsole: logToConsole,
	}, nil
}

// Execute runs an agent using the sentinel CLI
func (r *CliAgentRuntime) Execute(ctx context.Context, agentSpec stack.StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error) {
	// Create a unique directory for this agent run
	agentDir := filepath.Join(r.workDir, agentSpec.ID)
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create agent directory: %w", err)
	}

	// Write inputs to a file
	inputsFile := filepath.Join(agentDir, "inputs.json")
	inputsJSON, err := json.MarshalIndent(inputs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inputs: %w", err)
	}
	if err := ioutil.WriteFile(inputsFile, inputsJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write inputs file: %w", err)
	}

	// Generate a unique container name
	containerName := fmt.Sprintf("stack-%s-%d", agentSpec.ID, time.Now().Unix())

	// Prepare the output file path
	outputsFile := filepath.Join(agentDir, "outputs.json")

	// Build command to run the agent
	args := []string{
		"run",
		"--name", containerName,
		"--input-file", inputsFile,
		"--output-file", outputsFile,
		agentSpec.Uses,
	}

	// Log command if enabled
	if r.logToConsole {
		log.Printf("Running agent with command: %s %s", r.sentinelPath, strings.Join(args, " "))
	}

	// Execute the agent using the sentinel CLI
	cmd := exec.CommandContext(ctx, r.sentinelPath, args...)
	cmd.Dir = agentDir
	
	// Capture output
	var stdout, stderr bytes.Buffer
	if r.logToConsole {
		// If logging enabled, use both console and capture
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		// Otherwise just capture output
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	// Run the command
	startTime := time.Now()
	err = cmd.Run()
	duration := time.Since(startTime)
	
	if err != nil {
		cmdOutput := stderr.String()
		if cmdOutput == "" {
			cmdOutput = stdout.String()
		}
		return nil, fmt.Errorf("agent execution failed: %w\nOutput: %s", err, cmdOutput)
	}

	// Read outputs file if it exists
	var outputs map[string]interface{}
	if _, err := os.Stat(outputsFile); err == nil {
		outputData, err := ioutil.ReadFile(outputsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read outputs file: %w", err)
		}
		
		if err := json.Unmarshal(outputData, &outputs); err != nil {
			return nil, fmt.Errorf("failed to parse outputs JSON: %w", err)
		}
	} else {
		// No outputs file, create a default output with stdout as content
		outputs = map[string]interface{}{
			"status": "completed",
			"stdout": stdout.String(),
			"message": fmt.Sprintf("Agent %s executed successfully but produced no output file", agentSpec.ID),
		}
	}
	
	// Add execution metadata
	outputs["_agent_id"] = agentSpec.ID
	outputs["_execution_time"] = duration.String()
	outputs["_executed_at"] = time.Now().Format(time.RFC3339)
	
	return outputs, nil
}

// Cleanup removes temporary files
func (r *CliAgentRuntime) Cleanup() error {
	if r.workDir != "" {
		return os.RemoveAll(r.workDir)
	}
	return nil
}

// findSentinelPath locates the sentinel executable
func findSentinelPath() (string, error) {
	// First check if we're running as a sentinel command
	if execPath, err := os.Executable(); err == nil {
		if filepath.Base(execPath) == "sentinel" {
			return execPath, nil
		}
	}

	// Check in PATH
	if path, err := exec.LookPath("sentinel"); err == nil {
		return path, nil
	}

	// Check common installation locations
	commonLocations := []string{
		"/usr/local/bin/sentinel",
		"/usr/bin/sentinel",
		"/opt/sentinel/bin/sentinel",
	}

	for _, location := range commonLocations {
		if _, err := os.Stat(location); err == nil {
			return location, nil
		}
	}

	return "", fmt.Errorf("sentinel executable not found")
}
