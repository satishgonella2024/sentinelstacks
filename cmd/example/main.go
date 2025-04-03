// Package main provides an example of using the Sentinel Stacks API
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/api"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

func main() {
	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		fmt.Println("Received termination signal. Shutting down...")
		cancel()
		os.Exit(0)
	}()

	// Check for command line flags
	exampleType := "stack"
	if len(os.Args) > 1 {
		exampleType = os.Args[1]
	}

	// Run the appropriate example
	switch exampleType {
	case "stack":
		runStackExample(ctx)
	case "memory":
		runMemoryExample()
	default:
		fmt.Printf("Unknown example type: %s\n", exampleType)
		fmt.Println("Available examples: stack, memory")
		os.Exit(1)
	}
}

// runStackExample demonstrates using the StackService
func runStackExample(ctx context.Context) {
	fmt.Println("\n=== Stack Service Example ===")

	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Create API configuration with persistence enabled
	config := api.APIConfig{
		StackConfig: api.StackServiceConfig{
			StoragePath: filepath.Join(dataDir, "stacks"),
			Verbose:     true,
		},
		MemoryConfig: api.MemoryServiceConfig{
			StoragePath:         filepath.Join(dataDir, "memory"),
			EmbeddingProvider:   "local",
			EmbeddingModel:      "mock",
			EmbeddingDimensions: 1536,
		},
		RegistryConfig: api.RegistryServiceConfig{
			RegistryURL: "https://registry.example.com",
			CachePath:   filepath.Join(dataDir, "cache"),
			Username:    os.Getenv("REGISTRY_USERNAME"),
			AccessToken: os.Getenv("REGISTRY_TOKEN"),
		},
	}

	// Create the API
	sentinel, err := api.NewAPI(config)
	if err != nil {
		log.Fatalf("Failed to create API: %v", err)
	}

	// Get the stack service
	stackService := sentinel.Stack()

	// List existing stacks
	existingStacks, err := stackService.ListStacks(ctx)
	if err != nil {
		log.Fatalf("Failed to list stacks: %v", err)
	}

	var stackID string
	if len(existingStacks) > 0 {
		// Use an existing stack
		fmt.Println("Found existing stacks:")
		for i, stack := range existingStacks {
			fmt.Printf("%d. %s: %s (%s)\n", i+1, stack.ID, stack.Name, stack.Description)
		}

		stackID = existingStacks[0].ID
		fmt.Printf("\nUsing existing stack: %s\n", stackID)
	} else {
		// Create a new stack
		spec := createSampleStack()

		fmt.Printf("Creating new stack: %s\n", spec.Name)
		fmt.Printf("Description: %s\n", spec.Description)
		fmt.Printf("Agents: %d\n", len(spec.Agents))

		stackID, err = stackService.CreateStack(ctx, spec)
		if err != nil {
			log.Fatalf("Failed to create stack: %v", err)
		}

		fmt.Printf("Stack created with ID: %s\n", stackID)
	}

	// Execute the stack with sample inputs
	inputs := map[string]interface{}{
		"message": "Hello, Sentinel Stack!",
		"options": map[string]interface{}{
			"verbose":   true,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	fmt.Printf("Executing stack with inputs: %v\n", inputs)

	// Execute the stack
	results, err := stackService.ExecuteStack(ctx, stackID, inputs)
	if err != nil {
		log.Fatalf("Failed to execute stack: %v", err)
	}

	fmt.Printf("Stack execution complete!\n")
	fmt.Printf("Results: %v\n", results)

	// Get stack state
	state, err := stackService.GetStackState(ctx, stackID)
	if err != nil {
		log.Fatalf("Failed to get stack state: %v", err)
	}

	fmt.Printf("Stack state: %d/%d agents completed\n", state.CompletedCount, state.TotalAgents)

	// Get execution history
	history, err := stackService.GetStackExecutionHistory(ctx, stackID)
	if err != nil {
		log.Fatalf("Failed to get execution history: %v", err)
	}

	fmt.Printf("\nExecution History (%d records):\n", len(history))
	for i, execution := range history {
		fmt.Printf("%d. Execution ID: %s\n", i+1, execution.ExecutionID)
		fmt.Printf("   Start: %s\n", execution.StartTime.Format(time.RFC3339))
		fmt.Printf("   End: %s\n", execution.EndTime.Format(time.RFC3339))
		fmt.Printf("   Status: %s\n", execution.Status)
		fmt.Printf("   Agents: %d completed, %d failed, %d blocked\n",
			execution.CompletedCount, execution.FailedCount, execution.BlockedCount)
	}

	// Export stack to file
	exportPath := filepath.Join(dataDir, "exports", stackID+".json")
	if err := os.MkdirAll(filepath.Dir(exportPath), 0755); err != nil {
		log.Fatalf("Failed to create exports directory: %v", err)
	}

	if err := stackService.ExportStack(ctx, stackID, exportPath); err != nil {
		log.Fatalf("Failed to export stack: %v", err)
	}

	fmt.Printf("\nStack exported to: %s\n", exportPath)

	// Update stack if it's a new one
	if len(existingStacks) == 0 {
		// Get the exported stack definition
		data, err := os.ReadFile(exportPath)
		if err != nil {
			log.Fatalf("Failed to read exported stack: %v", err)
		}

		// Parse stack definition
		var spec types.StackSpec
		if err := json.Unmarshal(data, &spec); err != nil {
			log.Fatalf("Failed to parse stack definition: %v", err)
		}

		// Update the description
		spec.Description += " (Updated)"

		fmt.Printf("\nUpdating stack description to: %s\n", spec.Description)

		// Update the stack
		if err := stackService.UpdateStack(ctx, stackID, spec); err != nil {
			log.Fatalf("Failed to update stack: %v", err)
		}

		fmt.Printf("Stack updated successfully!\n")
	}

	// List all stacks again to show the update
	stacks, err := stackService.ListStacks(ctx)
	if err != nil {
		log.Fatalf("Failed to list stacks: %v", err)
	}

	fmt.Printf("\nAll stacks (%d):\n", len(stacks))
	for _, stack := range stacks {
		fmt.Printf("- %s: %s (%s)\n", stack.ID, stack.Name, stack.Description)
	}
}

// createSampleStack creates a sample stack for testing
func createSampleStack() types.StackSpec {
	return types.StackSpec{
		Name:        "example-stack",
		Description: "A sample stack for demonstration",
		Version:     "1.0.0",
		Type:        types.StackTypeDefault,
		Agents: []types.StackAgentSpec{
			{
				ID:   "agent1",
				Uses: "echo",
				With: map[string]interface{}{
					"message": "Hello from Agent 1",
				},
			},
			{
				ID:        "agent2",
				Uses:      "transform",
				InputFrom: []string{"agent1"},
				With: map[string]interface{}{
					"operation": "uppercase",
				},
			},
			{
				ID:        "agent3",
				Uses:      "output",
				InputFrom: []string{"agent2"},
				With: map[string]interface{}{
					"format": "json",
				},
			},
		},
	}
}
