package stack

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/stack"
)

// DependencyResolver handles resolving and managing stack dependencies
type DependencyResolver struct {
	RegistryClient *RegistryClient
	SentinelPath   string
}

// DependencyResult contains the results of dependency resolution
type DependencyResult struct {
	Required  []AgentReference
	Missing   []AgentReference
	Available []AgentReference
}

// NewDependencyResolver creates a new dependency resolver
func NewDependencyResolver(client *RegistryClient, sentinelPath string) *DependencyResolver {
	return &DependencyResolver{
		RegistryClient: client,
		SentinelPath:   sentinelPath,
	}
}

// ResolveStackDependencies identifies all agent dependencies for a stack
func (r *DependencyResolver) ResolveStackDependencies(spec stack.StackSpec) (*DependencyResult, error) {
	result := &DependencyResult{
		Required:  []AgentReference{},
		Missing:   []AgentReference{},
		Available: []AgentReference{},
	}

	// Extract all agent references
	for _, agent := range spec.Agents {
		ref := parseAgentReference(agent.Uses)
		result.Required = append(result.Required, ref)

		// Check if agent is available locally
		available, err := r.isAgentAvailable(ref)
		if err != nil {
			return nil, fmt.Errorf("failed to check