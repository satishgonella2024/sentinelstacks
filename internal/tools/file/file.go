package file

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/satishgonella2024/sentinelstacks/internal/tools"
)

// ReadFileTool implements a tool for reading files
type ReadFileTool struct {
	tools.BaseTool
	allowedDirs []string
}

// NewReadFileTool creates a new read file tool
func NewReadFileTool(allowedDirs []string) *ReadFileTool {
	return &ReadFileTool{
		BaseTool: tools.BaseTool{
			Name:        "file/read",
			Description: "Read the contents of a file from the file system",
			Parameters: []tools.Parameter{
				{
					Name:        "path",
					Type:        "string",
					Description: "Path to the file to read",
					Required:    true,
				},
				{
					Name:        "encoding",
					Type:        "string",
					Description: "Encoding to use (text, base64)",
					Required:    false,
					Default:     "text",
				},
			},
			Permission: tools.PermissionFile,
		},
		allowedDirs: allowedDirs,
	}
}

// Execute reads a file
func (t *ReadFileTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Get parameters
	path, _ := params["path"].(string)
	encoding, _ := params["encoding"].(string)
	if encoding == "" {
		encoding = "text"
	}
	
	// Check if path is allowed
	if !t.isPathAllowed(path) {
		return nil, fmt.Errorf("access to path not allowed: %s", path)
	}
	
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Return based on encoding
	if encoding == "base64" {
		return base64.StdEncoding.EncodeToString(data), nil
	}
	
	// Default to text
	return string(data), nil
}

// isPathAllowed checks if a path is within allowed directories
func (t *ReadFileTool) isPathAllowed(path string) bool {
	// If no allowed dirs specified, deny all
	if len(t.allowedDirs) == 0 {
		return false
	}
	
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	
	// Check against allowed directories
	for _, allowedDir := range t.allowedDirs {
		// Get absolute path for allowed dir
		absAllowedDir, err := filepath.Abs(allowedDir)
		if err != nil {
			continue
		}
		
		// Check if path is within allowed dir
		if strings.HasPrefix(absPath, absAllowedDir) {
			return true
		}
	}
	
	return false
}

// WriteFileTool implements a tool for writing files
type WriteFileTool struct {
	tools.BaseTool
	allowedDirs []string
}

// NewWriteFileTool creates a new write file tool
func NewWriteFileTool(allowedDirs []string) *WriteFileTool {
	return &WriteFileTool{
		BaseTool: tools.BaseTool{
			Name:        "file/write",
			Description: "Write content to a file in the file system",
			Parameters: []tools.Parameter{
				{
					Name:        "path",
					Type:        "string",
					Description: "Path to the file to write",
					Required:    true,
				},
				{
					Name:        "content",
					Type:        "string",
					Description: "Content to write to the file",
					Required:    true,
				},
				{
					Name:        "encoding",
					Type:        "string",
					Description: "Encoding of the content (text, base64)",
					Required:    false,
					Default:     "text",
				},
				{
					Name:        "append",
					Type:        "boolean",
					Description: "Whether to append to the file",
					Required:    false,
					Default:     false,
				},
			},
			Permission: tools.PermissionFile,
		},
		allowedDirs: allowedDirs,
	}
}

// Execute writes to a file
func (t *WriteFileTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Get parameters
	path, _ := params["path"].(string)
	content, _ := params["content"].(string)
	encoding, _ := params["encoding"].(string)
	append, _ := params["append"].(bool)
	
	if encoding == "" {
		encoding = "text"
	}
	
	// Check if path is allowed
	if !t.isPathAllowed(path) {
		return nil, fmt.Errorf("access to path not allowed: %s", path)
	}
	
	// Decode content if necessary
	var data []byte
	var err error
	if encoding == "base64" {
		data, err = base64.StdEncoding.DecodeString(content)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 content: %w", err)
		}
	} else {
		// Default to text
		data = []byte(content)
	}
	
	// Create parent directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write or append to file
	if append {
		// Open file for appending
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file for appending: %w", err)
		}
		defer file.Close()
		
		// Append to file
		if _, err := file.Write(data); err != nil {
			return nil, fmt.Errorf("failed to append to file: %w", err)
		}
	} else {
		// Write to file
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to write to file: %w", err)
		}
	}
	
	return map[string]interface{}{
		"success": true,
		"path":    path,
		"size":    len(data),
	}, nil
}

// isPathAllowed checks if a path is within allowed directories
func (t *WriteFileTool) isPathAllowed(path string) bool {
	// If no allowed dirs specified, deny all
	if len(t.allowedDirs) == 0 {
		return false
	}
	
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	
	// Check against allowed directories
	for _, allowedDir := range t.allowedDirs {
		// Get absolute path for allowed dir
		absAllowedDir, err := filepath.Abs(allowedDir)
		if err != nil {
			continue
		}
		
		// Check if path is within allowed dir
		if strings.HasPrefix(absPath, absAllowedDir) {
			return true
		}
	}
	
	return false
}

// ListDirectoryTool implements a tool for listing directory contents
type ListDirectoryTool struct {
	tools.BaseTool
	allowedDirs []string
}

// NewListDirectoryTool creates a new list directory tool
func NewListDirectoryTool(allowedDirs []string) *ListDirectoryTool {
	return &ListDirectoryTool{
		BaseTool: tools.BaseTool{
			Name:        "file/list",
			Description: "List the contents of a directory in the file system",
			Parameters: []tools.Parameter{
				{
					Name:        "path",
					Type:        "string",
					Description: "Path to the directory to list",
					Required:    true,
				},
				{
					Name:        "recursive",
					Type:        "boolean",
					Description: "Whether to list contents recursively",
					Required:    false,
					Default:     false,
				},
			},
			Permission: tools.PermissionFile,
		},
		allowedDirs: allowedDirs,
	}
}

// Execute lists a directory
func (t *ListDirectoryTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Get parameters
	path, _ := params["path"].(string)
	recursive, _ := params["recursive"].(bool)
	
	// Check if path is allowed
	if !t.isPathAllowed(path) {
		return nil, fmt.Errorf("access to path not allowed: %s", path)
	}
	
	// Get directory stats
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory: %w", err)
	}
	
	// Ensure path is a directory
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", path)
	}
	
	// List directory
	if recursive {
		return t.listRecursive(path)
	}
	
	return t.listDirectory(path)
}

// listDirectory lists the contents of a directory (non-recursive)
func (t *ListDirectoryTool) listDirectory(path string) (interface{}, error) {
	// Read directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	
	// Format entries
	result := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		result = append(result, map[string]interface{}{
			"name":      entry.Name(),
			"path":      filepath.Join(path, entry.Name()),
			"size":      info.Size(),
			"is_dir":    entry.IsDir(),
			"modified":  info.ModTime().Format("2006-01-02 15:04:05"),
			"mode":      info.Mode().String(),
		})
	}
	
	return result, nil
}

// listRecursive lists the contents of a directory recursively
func (t *ListDirectoryTool) listRecursive(path string) (interface{}, error) {
	var result []map[string]interface{}
	
	// Walk directory
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		
		// Skip the root path itself
		if filePath == path {
			return nil
		}
		
		// Add to result
		result = append(result, map[string]interface{}{
			"name":      filepath.Base(filePath),
			"path":      filePath,
			"size":      info.Size(),
			"is_dir":    info.IsDir(),
			"modified":  info.ModTime().Format("2006-01-02 15:04:05"),
			"mode":      info.Mode().String(),
		})
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	
	return result, nil
}

// isPathAllowed checks if a path is within allowed directories
func (t *ListDirectoryTool) isPathAllowed(path string) bool {
	// If no allowed dirs specified, deny all
	if len(t.allowedDirs) == 0 {
		return false
	}
	
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	
	// Check against allowed directories
	for _, allowedDir := range t.allowedDirs {
		// Get absolute path for allowed dir
		absAllowedDir, err := filepath.Abs(allowedDir)
		if err != nil {
			continue
		}
		
		// Check if path is within allowed dir
		if strings.HasPrefix(absPath, absAllowedDir) {
			return true
		}
	}
	
	return false
}