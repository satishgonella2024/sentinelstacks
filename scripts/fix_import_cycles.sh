#!/bin/bash

# Script to fix import cycles by creating common interfaces
echo "Fixing import cycles..."

# 1. Create common directory if it doesn't exist
mkdir -p internal/common

# 2. Create interfaces for components involved in import cycles
cat > internal/common/interfaces.go << 'EOF'
package common

import "context"

// Provider defines the interface for LLM providers
type Provider interface {
	// Complete generates a completion given a prompt
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
	
	// GetName returns the name of the provider
	GetName() string
	
	// GetDefaultModel returns the default model for this provider
	GetDefaultModel() string
}

// CompletionRequest represents a request to an LLM provider
type CompletionRequest struct {
	Model        string
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
	Temperature  float64
}

// CompletionResponse represents a response from an LLM provider
type CompletionResponse struct {
	Text         string
	FinishReason string
	Usage        struct {
		PromptTokens     int
		CompletionTokens int
		TotalTokens      int
	}
}

// AgentDefinition represents the definition of an agent
type AgentDefinition struct {
	Name           string
	Description    string
	Version        string
	BaseModel      string
	SystemPrompt   string
	PromptTemplate string
	MaxTokens      int
	Temperature    float64
	OutputFormat   string
}

// RuntimeInterface defines the interface for agent runtimes
type RuntimeInterface interface {
	// Execute executes an agent with inputs and returns outputs
	Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)
	
	// Stop stops the execution of an agent
	Stop() error
}
EOF

# 3. Create a file that identifies problematic imports
cat > scripts/problematic_imports.txt << 'EOF'
# This file contains a list of problematic import cycles

# 1. Shim <-> Runtime cycle
shim -> runtime -> shim

# 2. Stack -> Memory -> Stack cycle
stack -> memory -> stack

# Resolve these by:
# - Use internal/common for shared types
# - Create simplified implementations for testing
# - Break direct dependencies between packages
EOF

echo "Common interfaces created in internal/common/interfaces.go"
echo "See scripts/problematic_imports.txt for documentation on import cycles"
echo ""
echo "Next step: Update the imports in affected packages to use common interfaces"
