package stack

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// checkAgentExists checks if an agent exists locally
func (r *RealAgentRuntime) checkAgentExists(agentRef string) (bool, error) {
	// Split reference into name and tag
	parts := strings.Split(agentRef, ":")
	agentName := parts[0]
	
	// Check in agent data directory
	agentsDir := filepath.Join(r.agentDataDir, agentName)
	if _, err := os.Stat(agentsDir); err == nil {
		return true, nil
	}
	
	// Check using sentinel images command
	cmd := exec.Command(r.sentinelPath, "images", "--quiet")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list images: %w", err)
	}
	
	// Check if agent is in the output
	return strings.Contains(string(output), agentRef), nil
}

// pullAgent pulls an agent from the registry
func (r *RealAgentRuntime) pullAgent(ctx context.Context, agentRef string) error {
	log.Printf("Pulling agent: %s", agentRef)
	
	// Execute sentinel pull command
	cmd := exec.CommandContext(ctx, r.sentinelPath, "pull", agentRef)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull agent: %w", err)
	}
	
	log.Printf("Successfully pulled agent: %s", agentRef)
	return nil
}

// createSentinelfile creates a Sentinelfile for a mock agent
func (r *RealAgentRuntime) createSentinelfile(agentType string) (string, error) {
	// Create a temporary directory for the Sentinelfile
	tmpDir, err := ioutil.TempDir("", "sentinelfile-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	
	// Create Sentinelfile based on agent type
	var content string
	switch {
	case strings.Contains(agentType, "extractor"):
		content = `
name: data-extractor
description: Extracts data from various sources
version: 1.0.0
base_model: claude-3-sonnet-20240229
input:
  - name: source
    type: string
    description: "Source to extract data from"
  - name: format
    type: string
    default: "json"
    description: "Output format"
output_format: "json"
system_prompt: |
  You are a specialized data extraction agent. Your task is to extract structured data from the input according to the specified format.
  Always return valid JSON.
  
  When extracting data:
  1. Identify key entities and their attributes
  2. Structure the data hierarchically
  3. Ensure all values have appropriate types
  4. Validate the extraction against the source

prompt_template: |
  Please extract data from the following source:
  
  {{input}}
  
  Format the output as {{format}}.
`
	case strings.Contains(agentType, "transformer"):
		content = `
name: data-transformer
description: Transforms and cleans data
version: 1.0.0
base_model: claude-3-sonnet-20240229
input:
  - name: data
    type: object
    description: "Input data to transform"
  - name: operations
    type: array
    description: "Transformation operations to apply"
output_format: "json"
system_prompt: |
  You are a specialized data transformation agent. Your task is to apply the specified operations to transform the input data.
  Always return valid JSON.
  
  Supported operations:
  - filter: Remove records that don't match criteria
  - normalize: Standardize values
  - enrich: Add computed or derived fields
  - clean: Remove nulls, fix inconsistencies

prompt_template: |
  Please transform the following data:
  
  {{data}}
  
  Apply these operations: {{operations}}
`
	case strings.Contains(agentType, "analyzer"):
		content = `
name: data-analyzer
description: Analyzes data to extract insights
version: 1.0.0
base_model: claude-3-sonnet-20240229
input:
  - name: data
    type: object
    description: "Input data to analyze"
  - name: analysis_type
    type: string
    description: "Type of analysis to perform"
output_format: "json"
system_prompt: |
  You are a specialized data analysis agent. Your task is to analyze the input data and extract meaningful insights.
  Always return valid JSON with numerical values where appropriate.
  
  Analysis types:
  - segmentation: Group data into segments based on common characteristics
  - correlation: Find relationships between variables
  - trend: Identify patterns over time or sequences
  - anomaly: Detect outliers or unusual patterns

prompt_template: |
  Please analyze the following data:
  
  {{data}}
  
  Perform {{analysis_type}} analysis.
`
	default:
		// Generic agent
		content = fmt.Sprintf(`
name: %s
description: Generic agent for testing
version: 1.0.0
base_model: claude-3-sonnet-20240229
input:
  - name: input
    type: any
    description: "Input data"
output_format: "json"
system_prompt: |
  You are a generic test agent. Echo back the input with minimal processing.

prompt_template: |
  Process this input:
  
  {{input}}
`, agentType)
	}
	
	// Write Sentinelfile
	sentinelfilePath := filepath.Join(tmpDir, "Sentinelfile")
	if err := ioutil.WriteFile(sentinelfilePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write Sentinelfile: %w", err)
	}
	
	return tmpDir, nil
}

// buildMockAgent builds a mock agent for testing
func (r *RealAgentRuntime) buildMockAgent(ctx context.Context, agentType string) (string, error) {
	// Create Sentinelfile
	tmpDir, err := r.createSentinelfile(agentType)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)
	
	// Build the agent
	agentName := fmt.Sprintf("%s:latest", agentType)
	cmd := exec.CommandContext(ctx, r.sentinelPath, "build", "-t", agentName, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build agent: %w", err)
	}
	
	return agentName, nil
}

// ensureAgentExists makes sure the agent exists, building it if necessary
func (r *RealAgentRuntime) ensureAgentExists(ctx context.Context, agentRef string) error {
	exists, err := r.checkAgentExists(agentRef)
	if err != nil {
		return err
	}
	
	if exists {
		return nil
	}
	
	// Try to pull the agent
	pullErr := r.pullAgent(ctx, agentRef)
	if pullErr == nil {
		return nil
	}
	
	log.Printf("Failed to pull agent %s: %v", agentRef, pullErr)
	log.Printf("Trying to build a mock agent instead")
	
	// If pull fails, build a mock agent
	agentName := strings.Split(agentRef, ":")[0]
	_, err = r.buildMockAgent(ctx, agentName)
	if err != nil {
		return fmt.Errorf("failed to build mock agent: %w", err)
	}
	
	return nil
}
