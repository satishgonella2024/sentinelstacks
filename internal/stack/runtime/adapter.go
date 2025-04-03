// Package runtime provides runtime adapters for internal stack package
package runtime

import (
	"context"

	"github.com/satishgonella2024/sentinelstacks/internal/stack"
	"github.com/satishgonella2024/sentinelstacks/pkg/adapter"
	pkgRuntime "github.com/satishgonella2024/sentinelstacks/pkg/runtime"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// InternalAgentRuntime defines the interface for internal agent execution
type InternalAgentRuntime interface {
	// Execute runs an agent with the provided inputs and returns its outputs
	Execute(ctx context.Context, agentSpec stack.StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error)

	// Cleanup releases resources after execution
	Cleanup() error
}

// InternalRuntimeAdapter adapts pkg runtime implementations to internal agent runtime
type InternalRuntimeAdapter struct {
	runtime types.AgentRuntime
}

// NewInternalRuntimeAdapter creates a new adapter for public runtimes
func NewInternalRuntimeAdapter(runtime types.AgentRuntime) *InternalRuntimeAdapter {
	return &InternalRuntimeAdapter{
		runtime: runtime,
	}
}

// Execute runs an agent and returns its outputs
func (a *InternalRuntimeAdapter) Execute(ctx context.Context, agentSpec stack.StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error) {
	// Convert stack agent spec to adapter's internal type
	internalSpec := adapter.InternalStackAgentSpec{
		ID:        agentSpec.ID,
		Uses:      agentSpec.Uses,
		InputFrom: agentSpec.InputFrom,
		Depends:   agentSpec.Depends,
		Params:    agentSpec.Params,
	}

	// Convert to public type using the adapter
	publicSpec := adapter.ToPublicAgentSpec(internalSpec)

	// Execute using the public runtime
	return a.runtime.Execute(ctx, publicSpec, inputs)
}

// Cleanup performs cleanup
func (a *InternalRuntimeAdapter) Cleanup() error {
	return a.runtime.Cleanup()
}

// GetDirectRuntime returns an adapter for the direct runtime
func GetDirectRuntime(verbose bool) (InternalAgentRuntime, error) {
	// Create a public direct runtime
	publicRuntime, err := pkgRuntime.NewDirectRuntime(verbose)
	if err != nil {
		return nil, err
	}

	// Wrap with adapter
	return NewInternalRuntimeAdapter(publicRuntime), nil
}

// GetCliRuntime returns an adapter for the CLI runtime
func GetCliRuntime(verbose bool) (InternalAgentRuntime, error) {
	// Create a public CLI runtime
	publicRuntime, err := pkgRuntime.NewCliRuntime(verbose)
	if err != nil {
		return nil, err
	}

	// Wrap with adapter
	return NewInternalRuntimeAdapter(publicRuntime), nil
}

// GetDefaultRuntime returns the default runtime
func GetDefaultRuntime(verbose bool) (InternalAgentRuntime, error) {
	// Create a public factory
	factory := pkgRuntime.NewRuntimeFactory(verbose)

	// Get default runtime
	runtime, err := factory.DefaultRuntime()
	if err != nil {
		return nil, err
	}

	// Wrap with adapter
	return NewInternalRuntimeAdapter(runtime), nil
}
