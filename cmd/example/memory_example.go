// Package main provides examples of using the Sentinel Stacks system
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/api"
)

// runMemoryExample demonstrates using the MemoryService
func runMemoryExample() {
	fmt.Println("\n=== Memory Service Example ===")

	// Create API configuration
	config := api.APIConfig{
		MemoryConfig: api.MemoryServiceConfig{
			StoragePath:         "data/memory",
			EmbeddingProvider:   "local", // Use local in-memory implementation
			EmbeddingModel:      "mock",
			EmbeddingDimensions: 1536,
		},
		StackConfig: api.StackServiceConfig{
			StoragePath: "data/stacks",
			Verbose:     true,
		},
		RegistryConfig: api.RegistryServiceConfig{
			RegistryURL: "https://registry.example.com",
			CachePath:   "data/cache",
			Username:    os.Getenv("REGISTRY_USERNAME"),
			AccessToken: os.Getenv("REGISTRY_TOKEN"),
		},
	}

	// Create the API
	sentinel, err := api.NewAPI(config)
	if err != nil {
		log.Fatalf("Failed to create API: %v", err)
	}

	// Get the memory service
	memoryService := sentinel.Memory()

	// Create context
	ctx := context.Background()

	// Store and retrieve values
	fmt.Println("Testing key-value storage:")

	// Store a value
	err = memoryService.StoreValue(ctx, "test-collection", "key1", "Hello, Memory!")
	if err != nil {
		log.Fatalf("Failed to store value: %v", err)
	}
	fmt.Println("Stored value for key1")

	// Retrieve the value
	value, err := memoryService.RetrieveValue(ctx, "test-collection", "key1")
	if err != nil {
		log.Fatalf("Failed to retrieve value: %v", err)
	}
	fmt.Printf("Retrieved value for key1: %v\n", value)

	// Store a structured value
	person := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
		"address": map[string]string{
			"city":  "San Francisco",
			"state": "CA",
		},
	}

	err = memoryService.StoreValue(ctx, "test-collection", "person1", person)
	if err != nil {
		log.Fatalf("Failed to store person: %v", err)
	}
	fmt.Println("Stored person data")

	// Retrieve the structured value
	retrievedPerson, err := memoryService.RetrieveValue(ctx, "test-collection", "person1")
	if err != nil {
		log.Fatalf("Failed to retrieve person: %v", err)
	}
	fmt.Printf("Retrieved person data: %v\n", retrievedPerson)

	// Test vector storage and similarity search
	fmt.Println("\nTesting vector storage and similarity search:")

	// Store some texts with embeddings
	texts := []struct {
		key  string
		text string
		meta map[string]interface{}
	}{
		{
			key:  "doc1",
			text: "The quick brown fox jumps over the lazy dog",
			meta: map[string]interface{}{"category": "animals", "priority": 1},
		},
		{
			key:  "doc2",
			text: "The lazy dog slept all day long",
			meta: map[string]interface{}{"category": "animals", "priority": 2},
		},
		{
			key:  "doc3",
			text: "The quick rabbit runs through the forest",
			meta: map[string]interface{}{"category": "animals", "priority": 3},
		},
		{
			key:  "doc4",
			text: "The brown bear climbs the tall tree",
			meta: map[string]interface{}{"category": "animals", "priority": 4},
		},
		{
			key:  "doc5",
			text: "The red car drives down the highway",
			meta: map[string]interface{}{"category": "vehicles", "priority": 5},
		},
	}

	// Store all texts
	for _, item := range texts {
		err = memoryService.StoreEmbedding(ctx, "test-vectors", item.key, item.text, item.meta)
		if err != nil {
			log.Fatalf("Failed to store embedding for %s: %v", item.key, err)
		}
		fmt.Printf("Stored embedding for: %s\n", item.key)

		// Small delay to avoid overwhelming the system
		time.Sleep(50 * time.Millisecond)
	}

	// Search for similar texts
	fmt.Println("\nSearching for texts similar to 'dog sleeping in the sun':")
	results, err := memoryService.SearchSimilar(ctx, "test-vectors", "dog sleeping in the sun", 3)
	if err != nil {
		log.Fatalf("Failed to search similar texts: %v", err)
	}

	// Display results
	for i, match := range results {
		fmt.Printf("%d. Score: %.4f, Text: %s\n", i+1, match.Score, match.Content)
		fmt.Printf("   Metadata: %v\n", match.Metadata)
	}

	// Filter by metadata
	fmt.Println("\nSearching for texts with category=vehicles:")
	results, err = memoryService.SearchSimilar(ctx, "test-vectors", "car driving", 3)
	if err != nil {
		log.Fatalf("Failed to search similar texts: %v", err)
	}

	// Display results
	for i, match := range results {
		if match.Metadata["category"] == "vehicles" {
			fmt.Printf("%d. Score: %.4f, Text: %s\n", i+1, match.Score, match.Content)
			fmt.Printf("   Metadata: %v\n", match.Metadata)
		}
	}

	fmt.Println("\nMemory service example completed successfully!")
}
