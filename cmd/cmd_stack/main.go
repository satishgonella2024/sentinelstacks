package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/stack"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

func main() {
	// Create a sample stack specification
	spec := types.StackSpec{
		Name:        "test-stack",
		Description: "A test stack for trying out the engine",
		Version:     "1.0.0",
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
				Uses:      "echo",
				InputFrom: []string{"agent1"},
				With: map[string]interface{}{
					"prefix": "Agent 2 says: ",
				},
			},
			{
				ID:        "agent3",
				Uses:      "echo",
				InputFrom: []string{"agent2"},
				With: map[string]interface{}{
					"suffix": " (from Agent 3)",
				},
			},
		},
	}

	// Create a new stack engine
	engine, err := stack.NewEngine(spec, stack.WithVerbose(true))
	if err != nil {
		fmt.Printf("Error creating engine: %v\n", err)
		os.Exit(1)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute the stack
	fmt.Println("Starting stack execution...")
	err = engine.Execute(ctx,
		stack.WithTimeout(10),
		stack.WithInput(map[string]interface{}{
			"global": "This is a global input",
		}),
	)
	if err != nil {
		fmt.Printf("Error executing stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Stack execution completed successfully!")
}
