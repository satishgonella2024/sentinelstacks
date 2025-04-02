package file

import (
	"os"
	"path/filepath"

	"github.com/sentinelstacks/sentinel/internal/tools"
)

// RegisterFileTools registers all file tools
func RegisterFileTools() error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Set up allowed directories for file operations
	// By default, we only allow access to the .sentinel directory
	allowedDirs := []string{
		filepath.Join(homeDir, ".sentinel"),
	}

	// Get registry
	registry := tools.GetRegistry()

	// Register file tools
	registry.RegisterTool(NewReadFileTool(allowedDirs))
	registry.RegisterTool(NewWriteFileTool(allowedDirs))
	registry.RegisterTool(NewListDirectoryTool(allowedDirs))

	return nil
}
