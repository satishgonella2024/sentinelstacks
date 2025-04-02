package stack

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// AgentRuntime manages the execution of agents within a stack
type AgentRuntime interface {
	// Execute runs an agent with the given inputs and returns its outputs
	Execute(ctx context.Context, agentSpec StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error)
}

// RealAgentRuntime executes actual agents using the sentinel CLI
type RealAgentRuntime struct {
	workDir      string
	sentinelPath string
	agentDataDir string
}

// NewRealAgentRuntime creates a new agent runtime
func NewRealAgentRuntime() (*RealAgentRuntime, error) {
	// Create temporary working directory
	workDir, err := ioutil.TempDir("", "sentinel-stack-")
	if err != nil {
		return nil, fmt.Errorf("failed to create working directory: %w", err)
	}

	// Find sentinel executable path
	sentinelPath, err := findSentinelPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find sentinel executable: %w", err)
	}
	
	// Get agent data directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	
	agentDataDir := filepath.Join(home, ".sentinel", "agents")
	// Ensure it exists
	if err := os.MkdirAll(agentDataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create agent data directory: %w", err)
	}

	return &RealAgentRuntime{
		workDir:      workDir,
		sentinelPath: sentinelPath,
		agentDataDir: agentDataDir,
	}, nil
}

// Execute runs an agent and returns its outputs
func (r *RealAgentRuntime) Execute(ctx context.Context, agentSpec StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error) {
	// Create a unique directory for this agent run
	agentDir := filepath.Join(r.workDir, agentSpec.ID)
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create agent directory: %w", err)
	}

	// Write inputs to a file
	inputsFile := filepath.Join(agentDir, "inputs.json")
	inputsJSON, err := json.Marshal(inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inputs: %w", err)
	}
	if err := ioutil.WriteFile(inputsFile, inputsJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write inputs file: %w", err)
	}

	// Ensure the agent exists
	if err := r.ensureAgentExists(ctx, agentSpec.Uses); err != nil {
		return nil, fmt.Errorf("failed to ensure agent exists: %w", err)
	}

	// Generate a unique container name
	containerName := fmt.Sprintf("stack-%s-%d", agentSpec.ID, time.Now().Unix())

	// Build command to run the agent
	args := []string{
		"run",
		"--name", containerName,
		"--detach=false",  // Run in foreground
	}

	// Add agent reference
	args = append(args, agentSpec.Uses)

	// Add inputs file as bind mount
	args = append(args, "--mount", fmt.Sprintf("type=bind,source=%s,target=/sentinel/inputs.json", inputsFile))
	
	// Add output file
	outputsFile := filepath.Join(agentDir, "outputs.json")
	args = append(args, "--mount", fmt.Sprintf("type=bind,source=%s,target=/sentinel/outputs.json", outputsFile))

	// Execute the agent using the sentinel CLI
	cmd := exec.CommandContext(ctx, r.sentinelPath, args...)
	cmd.Dir = agentDir
	
	// Capture stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	log.Printf("Running agent: %s %s", r.sentinelPath, strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("agent execution failed: %w", err)
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
		// No outputs file, create a default output
		outputs = map[string]interface{}{
			"status": "completed",
			"message": fmt.Sprintf("Agent %s executed successfully but produced no output file", agentSpec.ID),
		}
	}

	return outputs, nil
}

// Cleanup removes temporary files
func (r *RealAgentRuntime) Cleanup() error {
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

// MockAgentRuntime is a simulated runtime for testing
type MockAgentRuntime struct {}

// NewMockAgentRuntime creates a new mock runtime
func NewMockAgentRuntime() *MockAgentRuntime {
	return &MockAgentRuntime{}
}

// Execute simulates agent execution
func (r *MockAgentRuntime) Execute(ctx context.Context, agentSpec StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error) {
	// Simulate processing time
	select {
	case <-time.After(500 * time.Millisecond):
		// Continue after delay
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Create outputs based on agent type
	outputs := make(map[string]interface{})
	outputs["_agent_id"] = agentSpec.ID
	outputs["_agent_type"] = agentSpec.Uses
	outputs["_processed_at"] = time.Now().Format(time.RFC3339)

	// Process inputs based on agent type
	switch {
	case strings.Contains(agentSpec.Uses, "processor"):
		// Process data
		if data, ok := inputs["data"].(string); ok {
			outputs["processed_data"] = fmt.Sprintf("Processed: %s", data)
		} else {
			outputs["processed_data"] = "Processed unknown data"
		}
	
	case strings.Contains(agentSpec.Uses, "analyzer"):
		// Analyze data
		outputs["analysis"] = map[string]interface{}{
			"sentiment": "positive",
			"entities": []string{"entity1", "entity2"},
			"confidence": 0.87,
		}
	
	case strings.Contains(agentSpec.Uses, "generator"):
		// Generate content
		outputs["generated_text"] = "This is generated content based on the inputs"
		outputs["generation_parameters"] = map[string]interface{}{
			"temperature": 0.7,
			"max_tokens": 100,
		}
	
	case strings.Contains(agentSpec.Uses, "summarizer"):
		// Summarize text
		if text, ok := inputs["text"].(string); ok {
			wordCount := len(strings.Fields(text))
			outputs["summary"] = fmt.Sprintf("Summary of %d words: The text discusses important topics.", wordCount)
		} else {
			outputs["summary"] = "Summary: The input contained relevant information."
		}
	
	default:
		// Default behavior - echo inputs
		for k, v := range inputs {
			outputs[k] = v
		}
		outputs["note"] = "Processed with default handler"
	}

	return outputs, nil
}
