// Package adapter provides adapters for type conversions
package adapter

import (
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// These type definitions mirror the internal ones but don't import them
// This allows us to convert without circular dependencies

// InternalStackAgentSpec mirrors the internal stack agent specification
type InternalStackAgentSpec struct {
	ID        string
	Uses      string
	InputFrom []string
	Depends   []string
	Params    map[string]interface{}
}

// InternalStackSpec mirrors the internal stack specification
type InternalStackSpec struct {
	Name        string
	Description string
	Version     string
	Agents      []InternalStackAgentSpec
}

// ToPublicAgentSpec converts an internal-like agent spec to a public one
func ToPublicAgentSpec(agentSpec InternalStackAgentSpec) types.StackAgentSpec {
	return types.StackAgentSpec{
		ID:        agentSpec.ID,
		Uses:      agentSpec.Uses,
		InputFrom: agentSpec.InputFrom,
		Depends:   agentSpec.Depends,
		With:      agentSpec.Params,
	}
}

// ToInternalAgentSpec converts a public agent spec to an internal-like one
func ToInternalAgentSpec(agentSpec types.StackAgentSpec) InternalStackAgentSpec {
	return InternalStackAgentSpec{
		ID:        agentSpec.ID,
		Uses:      agentSpec.Uses,
		InputFrom: agentSpec.InputFrom,
		Depends:   agentSpec.Depends,
		Params:    agentSpec.With,
	}
}

// ToPublicStackSpec converts an internal-like stack spec to a public one
func ToPublicStackSpec(stackSpec InternalStackSpec) types.StackSpec {
	agents := make([]types.StackAgentSpec, len(stackSpec.Agents))
	for i, agent := range stackSpec.Agents {
		agents[i] = ToPublicAgentSpec(agent)
	}

	return types.StackSpec{
		Name:        stackSpec.Name,
		Description: stackSpec.Description,
		Version:     stackSpec.Version,
		Agents:      agents,
	}
}

// ToInternalStackSpec converts a public stack spec to an internal-like one
func ToInternalStackSpec(stackSpec types.StackSpec) InternalStackSpec {
	agents := make([]InternalStackAgentSpec, len(stackSpec.Agents))
	for i, agent := range stackSpec.Agents {
		agents[i] = ToInternalAgentSpec(agent)
	}

	return InternalStackSpec{
		Name:        stackSpec.Name,
		Description: stackSpec.Description,
		Version:     stackSpec.Version,
		Agents:      agents,
	}
}
