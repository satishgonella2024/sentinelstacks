package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/satishgonella2024/sentinelstacks/internal/memory"
)

func setupTestDir() string {
	// Create a temp directory for testing
	tmpDir, _ := os.MkdirTemp("", "sentinel-test-")
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	return originalHome
}

func teardownTestDir(originalHome string) {
	// Restore original HOME
	os.Setenv("HOME", originalHome)
}

func TestSimpleMemory(t *testing.T) {
	originalHome := setupTestDir()
	defer teardownTestDir(originalHome)

	// Create a simple memory instance
	config := memory.MemoryConfig{
		Type:        memory.SimpleMemory,
		Persistence: true,
		MaxItems:    10,
	}

	mem, err := memory.NewMemory("test-agent", config)
	if err != nil {
		t.Fatalf("Failed to create memory: %v", err)
	}

	// Test adding entries
	id1, err := mem.Add("This is test entry 1", map[string]interface{}{
		"tag": "test",
	})
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	id2, err := mem.Add("This is test entry 2", map[string]interface{}{
		"tag": "important",
	})
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Test getting entries
	entry1, err := mem.Get(id1)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}
	if entry1.Content != "This is test entry 1" {
		t.Errorf("Expected content %q, got %q", "This is test entry 1", entry1.Content)
	}

	// Test searching
	results, err := mem.Search("test", 5)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Test listing
	allEntries, err := mem.List(100)
	if err != nil {
		t.Fatalf("Failed to list entries: %v", err)
	}
	if len(allEntries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(allEntries))
	}

	// Test persistence
	if err := mem.Save(); err != nil {
		t.Fatalf("Failed to save memory: %v", err)
	}

	// Create a new memory instance to test loading
	newMem, err := memory.NewMemory("test-agent", config)
	if err != nil {
		t.Fatalf("Failed to create new memory: %v", err)
	}

	// Verify the entries were loaded
	loadedEntry, err := newMem.Get(id1)
	if err != nil {
		t.Fatalf("Failed to get loaded entry: %v", err)
	}
	if loadedEntry.Content != "This is test entry 1" {
		t.Errorf("Expected content %q, got %q", "This is test entry 1", loadedEntry.Content)
	}

	// Test deleting
	if err := mem.Delete(id1); err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}

	// Verify deletion
	_, err = mem.Get(id1)
	if err == nil {
		t.Errorf("Expected error getting deleted entry, but got none")
	}

	// Test clearing
	if err := mem.Clear(); err != nil {
		t.Fatalf("Failed to clear memory: %v", err)
	}

	// Verify clearing
	allEntries, _ = mem.List(100)
	if len(allEntries) != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", len(allEntries))
	}
}

func TestVectorMemory(t *testing.T) {
	originalHome := setupTestDir()
	defer teardownTestDir(originalHome)

	// Create a vector memory instance
	config := memory.MemoryConfig{
		Type:        memory.VectorMemory,
		Persistence: true,
		MaxItems:    10,
	}

	mem, err := memory.NewMemory("test-agent-vector", config)
	if err != nil {
		t.Fatalf("Failed to create vector memory: %v", err)
	}

	// Test adding entries
	id1, err := mem.Add("The weather in London is rainy today", map[string]interface{}{
		"topic": "weather",
	})
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	id2, err := mem.Add("Paris has nice sunny weather this week", map[string]interface{}{
		"topic": "weather",
	})
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	id3, err := mem.Add("The stock market is trending upward", map[string]interface{}{
		"topic": "finance",
	})
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Test semantic search - should return weather-related entries
	results, err := mem.Search("What's the weather like in Europe?", 5)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	// Note: Since our dummy embedding isn't actually semantic,
	// we can't make strong assertions about search results,
	// but we can ensure the basic functionality works
	if len(results) == 0 {
		t.Errorf("Expected at least some search results, got none")
	}

	// Test deleting
	if err := mem.Delete(id3); err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}

	// Verify proper count after deletion
	allEntries, _ := mem.List(100)
	if len(allEntries) != 2 {
		t.Errorf("Expected 2 entries after deletion, got %d", len(allEntries))
	}

	// Test persistence
	if err := mem.Save(); err != nil {
		t.Fatalf("Failed to save memory: %v", err)
	}

	// Load into a new instance
	newMem, err := memory.NewMemory("test-agent-vector", config)
	if err != nil {
		t.Fatalf("Failed to create new memory: %v", err)
	}

	// Verify the entries were loaded
	loadedEntries, _ := newMem.List(100)
	if len(loadedEntries) != 2 {
		t.Errorf("Expected 2 loaded entries, got %d", len(loadedEntries))
	}

	// Test max items limit
	// Add more entries than the max limit
	for i := 0; i < 10; i++ {
		_, err := newMem.Add(fmt.Sprintf("Entry %d", i), nil)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}
	}

	// Verify that the memory pruned older entries
	allEntries, _ = newMem.List(100)
	if len(allEntries) > 10 {
		t.Errorf("Expected max 10 entries, got %d", len(allEntries))
	}
}

func TestMemoryFactory(t *testing.T) {
	originalHome := setupTestDir()
	defer teardownTestDir(originalHome)

	// Test creating different memory types via factory
	testCases := []struct {
		name        string
		memoryType  memory.MemoryType
		expectedErr bool
	}{
		{"SimpleMemory", memory.SimpleMemory, false},
		{"VectorMemory", memory.VectorMemory, false},
		{"InvalidType", "invalid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := memory.MemoryConfig{
				Type:        tc.memoryType,
				Persistence: false,
			}

			mem, err := memory.NewMemory("test-factory", config)
			if tc.expectedErr {
				if err == nil {
					t.Errorf("Expected error for memory type %q, but got none", tc.memoryType)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error creating memory type %q: %v", tc.memoryType, err)
				}
				if mem == nil {
					t.Errorf("Expected non-nil memory for type %q", tc.memoryType)
				}
			}
		})
	}
}
