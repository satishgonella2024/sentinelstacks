package memory

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Initialize sets up the memory subsystem
func Initialize() error {
	// Create default memory directories
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	// Create directories
	dirs := []string{
		filepath.Join(home, ".sentinel", "memory"),
		filepath.Join(home, ".sentinel", "memory", "local"),
		filepath.Join(home, ".sentinel", "memory", "sqlite"),
		filepath.Join(home, ".sentinel", "memory", "chroma"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	// Log initialization
	log.Printf("Memory subsystem initialized")
	
	return nil
}

// Shutdown performs cleanup operations for the memory subsystem
func Shutdown() error {
	// Currently no cleanup needed
	return nil
}

// DefaultMemoryPath returns the default path for memory storage
func DefaultMemoryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	return filepath.Join(home, ".sentinel", "memory"), nil
}
