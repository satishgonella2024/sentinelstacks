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
	exampleType := "comprehensive"
	if len(os.Args) > 1 {
		exampleType = os.Args[1]
	}

	// Run the appropriate example
	switch exampleType {
	case "comprehensive":
		runComprehensiveExample(ctx)
	case "stack":
		runStackExample(ctx)
	case "memory":
		runMemoryExample(ctx)
	case "registry":
		runRegistryExample(ctx)
	default:
		fmt.Printf("Unknown example type: %s\n", exampleType)
		fmt.Println("Available examples: comprehensive, stack, memory, registry")
		os.Exit(1)
	}
}

// runComprehensiveExample demonstrates using all API services together
func runComprehensiveExample(ctx context.Context) {
	fmt.Println("\n=== Comprehensive API Example ===")

	// Initialize the API
	api, err := initializeAPI()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}
	defer api.Close()

	// Create a stack
	stackID, err := createExampleStack(ctx, api)
	if err != nil {
		log.Fatalf("Failed to create stack: %v", err)
	}

	// Store some data in memory
	if err := storeExampleMemory(ctx, api); err != nil {
		log.Fatalf("Failed to store memory: %v", err)
	}

	// Search for packages
	if err := searchRegistryPackages(ctx, api); err != nil {
		log.Fatalf("Failed to search packages: %v", err)
	}

	// Execute the stack with data from memory
	result, err := executeStackWithMemory(ctx, api, stackID)
	if err != nil {
		log.Fatalf("Failed to execute stack: %v", err)
	}

	// Print the result
	fmt.Printf("\nStack Execution Result:\n")
	prettyPrint(result)

	// Get execution history
	history, err := api.Stack().GetStackExecutionHistory(ctx, stackID)
	if err != nil {
		log.Fatalf("Failed to get execution history: %v", err)
	}

	// Print execution history
	fmt.Printf("\nExecution History (%d records):\n", len(history))
	for i, exec := range history {
		fmt.Printf("%d. ID: %s, Status: %s, Completed: %d, Failed: %d\n",
			i+1, exec.ExecutionID, exec.Status, exec.CompletedCount, exec.FailedCount)
	}

	fmt.Println("\nComprehensive example completed successfully!")
}

// runStackExample demonstrates using the stack service
func runStackExample(ctx context.Context) {
	fmt.Println("\n=== Stack Service Example ===")

	// Initialize the API
	api, err := initializeAPI()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}
	defer api.Close()

	// Create a stack
	stackID, err := createExampleStack(ctx, api)
	if err != nil {
		log.Fatalf("Failed to create stack: %v", err)
	}

	// Execute the stack
	inputs := map[string]interface{}{
		"message": "Hello from the stack example!",
		"options": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	result, err := api.Stack().ExecuteStack(ctx, stackID, inputs)
	if err != nil {
		log.Fatalf("Failed to execute stack: %v", err)
	}

	// Print the result
	fmt.Printf("\nStack Execution Result:\n")
	prettyPrint(result)

	// Get stack state
	state, err := api.Stack().GetStackState(ctx, stackID)
	if err != nil {
		log.Fatalf("Failed to get stack state: %v", err)
	}

	fmt.Printf("\nStack State: %d/%d agents completed\n",
		state.CompletedCount, state.TotalAgents)

	// Export the stack
	exportPath := filepath.Join("data", "exports", stackID+".json")
	if err := os.MkdirAll(filepath.Dir(exportPath), 0755); err != nil {
		log.Fatalf("Failed to create export directory: %v", err)
	}

	if err := api.Stack().ExportStack(ctx, stackID, exportPath); err != nil {
		log.Fatalf("Failed to export stack: %v", err)
	}

	fmt.Printf("\nStack exported to: %s\n", exportPath)
}

// runMemoryExample demonstrates using the memory service
func runMemoryExample(ctx context.Context) {
	fmt.Println("\n=== Memory Service Example ===")

	// Initialize the API
	api, err := initializeAPI()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}
	defer api.Close()

	// Collection for the example
	collection := "example-collection"

	// Store a simple value
	key := "greeting"
	value := map[string]interface{}{
		"message":   "Hello from memory service!",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	fmt.Printf("Storing value with key '%s'\n", key)
	if err := api.Memory().StoreValue(ctx, collection, key, value); err != nil {
		log.Fatalf("Failed to store value: %v", err)
	}

	// Retrieve the value
	fmt.Printf("Retrieving value with key '%s'\n", key)
	retrieved, err := api.Memory().RetrieveValue(ctx, collection, key)
	if err != nil {
		log.Fatalf("Failed to retrieve value: %v", err)
	}

	fmt.Printf("\nRetrieved Value:\n")
	prettyPrint(retrieved)

	// Store embeddings
	docs := []struct {
		id   string
		text string
		meta map[string]interface{}
	}{
		{
			id:   "doc1",
			text: "The quick brown fox jumps over the lazy dog",
			meta: map[string]interface{}{"category": "animals", "priority": 1},
		},
		{
			id:   "doc2",
			text: "A journey of a thousand miles begins with a single step",
			meta: map[string]interface{}{"category": "wisdom", "priority": 2},
		},
		{
			id:   "doc3",
			text: "To be or not to be, that is the question",
			meta: map[string]interface{}{"category": "literature", "priority": 3},
		},
	}

	fmt.Println("\nStoring text embeddings:")
	for _, doc := range docs {
		fmt.Printf("- Storing document: %s\n", doc.id)
		if err := api.Memory().StoreEmbedding(ctx, collection, doc.id, doc.text, doc.meta); err != nil {
			log.Fatalf("Failed to store embedding for %s: %v", doc.id, err)
		}
	}

	// Search for similar embeddings
	query := "A fox quickly jumped"
	limit := 2

	fmt.Printf("\nSearching for documents similar to: '%s'\n", query)
	results, err := api.Memory().SearchSimilar(ctx, collection, query, limit)
	if err != nil {
		log.Fatalf("Failed to search similar documents: %v", err)
	}

	fmt.Printf("\nSearch Results (%d matches):\n", len(results))
	for i, match := range results {
		fmt.Printf("%d. ID: %s (Score: %.4f)\n", i+1, match.Key, match.Score)
		fmt.Printf("   Content: %s\n", match.Content)
		fmt.Printf("   Metadata: %v\n", match.Metadata)
	}
}

// runRegistryExample demonstrates using the registry service
func runRegistryExample(ctx context.Context) {
	fmt.Println("\n=== Registry Service Example ===")

	// Initialize the API
	api, err := initializeAPI()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}
	defer api.Close()

	// Search for packages
	query := "transform"
	limit := 5

	fmt.Printf("Searching for packages matching: '%s'\n", query)
	packages, err := api.Registry().SearchPackages(ctx, query, limit)
	if err != nil {
		log.Fatalf("Failed to search packages: %v", err)
	}

	fmt.Printf("\nSearch Results (%d packages):\n", len(packages))
	for i, pkg := range packages {
		fmt.Printf("%d. %s@%s\n", i+1, pkg.Name, pkg.Version)
		fmt.Printf("   Description: %s\n", pkg.Description)
		fmt.Printf("   Author: %s\n", pkg.Author)
		if len(pkg.Tags) > 0 {
			fmt.Printf("   Tags: %v\n", pkg.Tags)
		}
	}

	// Pull a package
	if len(packages) > 0 {
		pkg := packages[0]
		fmt.Printf("\nPulling package: %s@%s\n", pkg.Name, pkg.Version)

		path, err := api.Registry().PullPackage(ctx, pkg.Name, pkg.Version)
		if err != nil {
			log.Fatalf("Failed to pull package: %v", err)
		}

		fmt.Printf("Package downloaded to: %s\n", path)

		// Read the package content
		fmt.Println("\nPackage Content:")
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read package file: %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(content, &parsed); err != nil {
			fmt.Printf("Raw content: %s\n", content)
		} else {
			prettyPrint(parsed)
		}
	}
}

// initializeAPI creates and initializes the API
func initializeAPI() (types.API, error) {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Create API configuration
	config := api.APIConfig{
		StackConfig: api.StackServiceConfig{
			StoragePath: filepath.Join(dataDir, "stacks"),
			Verbose:     true,
		},
		MemoryConfig: api.MemoryServiceConfig{
			StoragePath:         filepath.Join(dataDir, "memory"),
			EmbeddingProvider:   "local",
			EmbeddingModel:      "local",
			EmbeddingDimensions: 1536,
		},
		RegistryConfig: api.RegistryServiceConfig{
			RegistryURL: "https://registry.sentinelstacks.example",
			CachePath:   filepath.Join(dataDir, "registry-cache"),
			Username:    os.Getenv("REGISTRY_USERNAME"),
			AccessToken: os.Getenv("REGISTRY_TOKEN"),
		},
	}

	// Create API instance
	return api.NewAPI(config)
}

// createExampleStack creates an example stack
func createExampleStack(ctx context.Context, api types.API) (string, error) {
	// Check if stack already exists
	stacks, err := api.Stack().ListStacks(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list stacks: %w", err)
	}

	// Return existing stack if any
	if len(stacks) > 0 {
		fmt.Printf("Using existing stack: %s (%s)\n", stacks[0].Name, stacks[0].ID)
		return stacks[0].ID, nil
	}

	// Create stack specification
	spec := types.StackSpec{
		Name:        "example-stack",
		Description: "A stack demonstrating the API",
		Version:     "1.0.0",
		Type:        types.StackTypeDefault,
		Agents: []types.StackAgentSpec{
			{
				ID:   "input-agent",
				Uses: "echo",
				With: map[string]interface{}{
					"message": "Processed input: ${message}",
				},
			},
			{
				ID:        "process-agent",
				Uses:      "transform",
				InputFrom: []string{"input-agent"},
				With: map[string]interface{}{
					"operation": "uppercase",
				},
			},
			{
				ID:        "output-agent",
				Uses:      "output",
				InputFrom: []string{"process-agent"},
				With: map[string]interface{}{
					"format": "json",
				},
			},
		},
	}

	// Create the stack
	fmt.Printf("Creating new stack: %s\n", spec.Name)
	stackID, err := api.Stack().CreateStack(ctx, spec)
	if err != nil {
		return "", fmt.Errorf("failed to create stack: %w", err)
	}

	fmt.Printf("Stack created with ID: %s\n", stackID)
	return stackID, nil
}

// storeExampleMemory stores example data in memory
func storeExampleMemory(ctx context.Context, api types.API) error {
	collection := "system-memory"

	// Store system info
	systemInfo := map[string]interface{}{
		"hostname":      "sentinel-host",
		"os":            "Linux",
		"uptime":        3600,
		"memory_total":  8589934592,
		"memory_used":   4294967296,
		"cpu_usage":     0.35,
		"timestamp":     time.Now().Format(time.RFC3339),
		"active_stacks": 1,
		"active_agents": 3,
	}

	fmt.Println("\nStoring system information in memory...")
	if err := api.Memory().StoreValue(ctx, collection, "system-info", systemInfo); err != nil {
		return fmt.Errorf("failed to store system info: %w", err)
	}

	// Store contextual data
	contextData := map[string]interface{}{
		"message": "Data from memory service",
		"options": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"source":    "memory-service",
			"priority":  "high",
		},
	}

	if err := api.Memory().StoreValue(ctx, collection, "stack-context", contextData); err != nil {
		return fmt.Errorf("failed to store context data: %w", err)
	}

	fmt.Println("Memory data stored successfully")
	return nil
}

// searchRegistryPackages searches for packages in the registry
func searchRegistryPackages(ctx context.Context, api types.API) error {
	query := "" // Empty query to get all packages
	limit := 5

	fmt.Println("\nSearching for available packages...")
	packages, err := api.Registry().SearchPackages(ctx, query, limit)
	if err != nil {
		return fmt.Errorf("failed to search packages: %w", err)
	}

	fmt.Printf("Found %d packages:\n", len(packages))
	for i, pkg := range packages {
		fmt.Printf("%d. %s@%s - %s\n",
			i+1, pkg.Name, pkg.Version, pkg.Description)
	}

	return nil
}

// executeStackWithMemory executes a stack with data from memory
func executeStackWithMemory(ctx context.Context, api types.API, stackID string) (map[string]interface{}, error) {
	fmt.Println("\nRetrieving context data from memory...")
	contextData, err := api.Memory().RetrieveValue(ctx, "system-memory", "stack-context")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve context data: %w", err)
	}

	// Convert to map
	inputs, ok := contextData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid context data format")
	}

	fmt.Printf("Executing stack with data from memory...\n")
	return api.Stack().ExecuteStack(ctx, stackID, inputs)
}

// prettyPrint prints a value as formatted JSON
func prettyPrint(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("%+v\n", v)
		return
	}
	fmt.Println(string(data))
}
