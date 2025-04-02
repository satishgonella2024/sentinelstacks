package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sentinelstacks/sentinel/pkg/agent"
)

// LocalRegistry represents a local registry for Sentinel Images
type LocalRegistry struct {
	imagesDir string
}

// GetLocalRegistry returns a LocalRegistry instance
func GetLocalRegistry() (*LocalRegistry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	imagesDir := filepath.Join(homeDir, ".sentinel/images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create images directory: %w", err)
	}

	return &LocalRegistry{
		imagesDir: imagesDir,
	}, nil
}

// ImageInfo contains information about a Sentinel Image
type ImageInfo struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Tag        string                 `json:"tag"`
	CreatedAt  time.Time              `json:"createdAt"`
	Size       int64                  `json:"size"`
	BaseModel  string                 `json:"baseModel"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ListImageInfo returns all images in the local registry
func (r *LocalRegistry) ListImageInfo() ([]ImageInfo, error) {
	var images []ImageInfo

	// Get all files in the images directory
	entries, err := os.ReadDir(r.imagesDir)
	if err != nil {
		if os.IsNotExist(err) {
			// If the directory doesn't exist, return an empty list
			return images, nil
		}
		return nil, fmt.Errorf("failed to read images directory: %w", err)
	}

	// Iterate through each file
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// Read the file
		filePath := filepath.Join(r.imagesDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: Failed to read %s: %v\n", filePath, err)
			continue
		}

		// Parse the file as a Sentinel Image
		var image agent.Image
		if err := json.Unmarshal(data, &image); err != nil {
			fmt.Printf("Warning: Failed to parse %s: %v\n", filePath, err)
			continue
		}

		// Get file info for size and modification time
		fileInfo, err := entry.Info()
		if err != nil {
			fmt.Printf("Warning: Failed to get file info for %s: %v\n", filePath, err)
			continue
		}

		// Create ImageInfo from the parsed image
		imageInfo := ImageInfo{
			ID:        image.ID,
			Name:      image.Name,
			Tag:       image.Tag,
			CreatedAt: time.Unix(image.CreatedAt, 0),
			Size:      fileInfo.Size(),
			BaseModel: image.Definition.BaseModel,
		}

		// Add parameters if they exist
		if image.Definition.Parameters != nil {
			imageInfo.Parameters = image.Definition.Parameters
		}

		images = append(images, imageInfo)
	}

	return images, nil
}

// Get returns an image by name and tag
func (r *LocalRegistry) Get(name, tag string) (*Image, error) {
	// If tag is not specified, use "latest"
	if tag == "" {
		tag = "latest"
	}

	// Create the file path
	filename := fmt.Sprintf("%s_%s.json", strings.ReplaceAll(name, "/", "_"), tag)
	filePath := filepath.Join(r.imagesDir, filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image %s:%s not found", name, tag)
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	// Parse the file as a Sentinel Image
	var agentImage agent.Image
	if err := json.Unmarshal(data, &agentImage); err != nil {
		return nil, fmt.Errorf("failed to parse image: %w", err)
	}

	// Convert to registry Image
	image := ConvertFromAgentImage(&agentImage)

	return image, nil
}

// Save stores an image in the local registry
func (r *LocalRegistry) Save(image *Image) error {
	// Ensure the images directory exists
	if err := os.MkdirAll(r.imagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create images directory: %w", err)
	}

	// Create the file path
	filename := fmt.Sprintf("%s_%s.json", strings.ReplaceAll(image.Name, "/", "_"), image.Tag)
	filePath := filepath.Join(r.imagesDir, filename)

	// Convert to agent Image
	agentImage := ConvertToAgentImage(image)

	// Marshal the image to JSON
	data, err := json.MarshalIndent(agentImage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal image: %w", err)
	}

	// Write the file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write image file: %w", err)
	}

	return nil
}

// Delete removes an image from the local registry
func (r *LocalRegistry) Delete(name, tag string) error {
	// If tag is not specified, use "latest"
	if tag == "" {
		tag = "latest"
	}

	// Create the file path
	filename := fmt.Sprintf("%s_%s.json", strings.ReplaceAll(name, "/", "_"), tag)
	filePath := filepath.Join(r.imagesDir, filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("image %s:%s not found", name, tag)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}
