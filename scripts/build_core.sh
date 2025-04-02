#!/bin/bash

# This script builds just the core stack components for testing
set -e  # Exit on error

echo "Building core stack components..."

# Create a temporary main file that only uses the stack components
TMP_DIR=$(mktemp -d)
TMP_MAIN="$TMP_DIR/main.go"

cat > "$TMP_MAIN" << 'EOF'
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/stack"
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/memory"
)

func main() {
	fmt.Println("SentinelStacks Core Test")
	
	// Initialize memory
	memory.Initialize()
	
	// Create a simple stack spec for testing
	spec := stack.StackSpec{
		Name:        "test-stack",
		Description: "A test stack",
		Version:     "1.0.0",
		Agents: []stack.StackAgentSpec{
			{
				ID:   "agent1",
				Uses: "test-agent",
				Params: map[string]interface{}{
					"param1": "value1",
				},
			},
			{
				ID:        "agent2",
				Uses:      "test-agent",
				InputFrom: []string{"agent1"},
				Params: map[string]interface{}{
					"param2": "value2",
				},
			},
		},
	}
	
	// Create a stack engine
	eng, err := stack.NewStackEngine(spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stack engine: %v\n", err)
		os.Exit(1)
	}
	
	// Build execution graph
	dag, err := eng.BuildExecutionGraph()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building execution graph: %v\n", err)
		os.Exit(1)
	}
	
	// Get execution order
	order, err := dag.TopologicalSort()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error determining execution order: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Execution order: %v\n", order)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Execute stack
	fmt.Println("Executing stack...")
	err = eng.Execute(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing stack: %v\n", err)
		os.Exit(1)
	}
	
	// Get execution summary
	summary := eng.GetState()
	fmt.Printf("\nExecution summary:\n")
	fmt.Printf("Total agents: %d\n", summary.TotalAgents)
	fmt.Printf("Completed: %d\n", summary.CompletedCount)
	fmt.Printf("Failed: %d\n", summary.FailedCount)
	
	fmt.Println("Core test completed successfully.")
}
EOF

# Try to build the test program
echo "Building test program..."
go build -o stack-test "$TMP_MAIN"

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Running test program..."
    ./stack-test
    rm ./stack-test
else
    echo "Build failed."
    exit 1
fi

# Clean up
rm -rf "$TMP_DIR"

echo "Core stack component test complete."
