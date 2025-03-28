package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/agentfile"
)

// RegistryClient handles interaction with agent registries
type RegistryClient struct {
	// Base path for the local registry cache
	LocalPath string
	// Remote registry URL (for future use)
	RemoteURL string
}

// AgentMetadata stores metadata about an agent in the registry
type AgentMetadata struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Author      string            `json:"author"`
	Description string            `json:"description"`
	Tags        []string          `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Downloads   int               `json:"downloads"`
	Models      []string          `json:"models"`
	Visibility  string            `json:"visibility"` // public or private
	Capabilities []string         `json:"capabilities"`
	Extra       map[string]string `json:"extra,omitempty"`
}

// NewRegistryClient creates a new registry client
func NewRegistryClient() (*RegistryClient, error) {
	// Create local registry directory if it doesn't exist
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	localPath := filepath.Join(home, ".sentinelstacks", "registry")
	err = os.MkdirAll(localPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create local registry directory: %w", err)
	}

	return &RegistryClient{
		LocalPath: localPath,
		RemoteURL: "https://registry.sentinelstacks.io", // This is a placeholder for future use
	}, nil
}

// PushAgent pushes an agent to the registry
func (c *RegistryClient) PushAgent(agentPath string, visibility string) error {
	// Validate agent path
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		return fmt.Errorf("agent path does not exist: %s", agentPath)
	}

	// Read the agentfile
	agentfilePath := filepath.Join(agentPath, "agentfile.yaml")
	data, err := os.ReadFile(agentfilePath)
	if err != nil {
		return fmt.Errorf("failed to read agentfile: %w", err)
	}

	// Parse the agentfile
	var af agentfile.Agentfile
	err = json.Unmarshal(data, &af)
	if err != nil {
		return fmt.Errorf("failed to parse agentfile: %w", err)
	}

	// Create metadata
	username := getUsername()
	metadata := AgentMetadata{
		Name:        af.Name,
		Version:     af.Version,
		Author:      username,
		Description: af.Description,
		Tags:        []string{}, // Extract from capabilities or extra fields
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Downloads:   0,
		Models:      []string{af.Model.Provider + "/" + af.Model.Name},
		Visibility:  visibility,
		Capabilities: af.Capabilities,
	}

	// Create registry directory for this agent
	registryAgentPath := filepath.Join(c.LocalPath, username, af.Name, af.Version)
	err = os.MkdirAll(registryAgentPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create registry directory: %w", err)
	}

	// Write metadata
	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = os.WriteFile(filepath.Join(registryAgentPath, "metadata.json"), metadataBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// Copy agent files
	err = copyDir(agentPath, registryAgentPath)
	if err != nil {
		return fmt.Errorf("failed to copy agent files: %w", err)
	}

	fmt.Printf("Successfully pushed %s/%s@%s to registry\n", username, af.Name, af.Version)
	return nil
}

// PullAgent pulls an agent from the registry
func (c *RegistryClient) PullAgent(agentRef string) (string, error) {
	// Parse agent reference (username/name[@version])
	parts := strings.Split(agentRef, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid agent reference: %s (should be username/name[@version])", agentRef)
	}

	username := parts[0]
	nameVersion := parts[1]
	name := nameVersion
	version := "latest"

	// Check if version is specified
	if strings.Contains(nameVersion, "@") {
		nvParts := strings.Split(nameVersion, "@")
		name = nvParts[0]
		version = nvParts[1]
	}

	// Find the agent in the registry
	var registryAgentPath string
	if version == "latest" {
		// Find the latest version
		versionsPath := filepath.Join(c.LocalPath, username, name)
		versions, err := os.ReadDir(versionsPath)
		if err != nil {
			return "", fmt.Errorf("failed to read versions directory: %w", err)
		}

		if len(versions) == 0 {
			return "", fmt.Errorf("no versions found for agent: %s/%s", username, name)
		}

		// Sort versions (simple lexicographical sort for now)
		// TODO: Use proper semver sorting
		latestVersion := versions[0].Name()
		for _, v := range versions {
			if v.Name() > latestVersion {
				latestVersion = v.Name()
			}
		}

		version = latestVersion
		registryAgentPath = filepath.Join(versionsPath, version)
	} else {
		registryAgentPath = filepath.Join(c.LocalPath, username, name, version)
	}

	if _, err := os.Stat(registryAgentPath); os.IsNotExist(err) {
		return "", fmt.Errorf("agent not found in registry: %s/%s@%s", username, name, version)
	}

	// Create destination directory
	destPath := filepath.Join(".", name)
	if _, err := os.Stat(destPath); err == nil {
		return "", fmt.Errorf("destination directory already exists: %s", destPath)
	}

	err := os.MkdirAll(destPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy agent files
	err = copyDir(registryAgentPath, destPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy agent files: %w", err)
	}

	// Update download count
	metadata, err := readMetadata(filepath.Join(registryAgentPath, "metadata.json"))
	if err != nil {
		fmt.Printf("Warning: Failed to update download count: %v\n", err)
	} else {
		metadata.Downloads++
		metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			fmt.Printf("Warning: Failed to marshal updated metadata: %v\n", err)
		} else {
			err = os.WriteFile(filepath.Join(registryAgentPath, "metadata.json"), metadataBytes, 0644)
			if err != nil {
				fmt.Printf("Warning: Failed to write updated metadata: %v\n", err)
			}
		}
	}

	fmt.Printf("Successfully pulled %s/%s@%s to %s\n", username, name, version, destPath)
	return destPath, nil
}

// SearchAgents searches for agents in the registry
func (c *RegistryClient) SearchAgents(query string, tags []string) ([]AgentMetadata, error) {
	var results []AgentMetadata

	// Walk the registry directory
	err := filepath.Walk(c.LocalPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process metadata.json files
		if !info.IsDir() && filepath.Base(path) == "metadata.json" {
			// Read metadata
			metadata, err := readMetadata(path)
			if err != nil {
				fmt.Printf("Warning: Failed to read metadata at %s: %v\n", path, err)
				return nil
			}

			// Check if agent matches query
			queryLower := strings.ToLower(query)
			nameLower := strings.ToLower(metadata.Name)
			descLower := strings.ToLower(metadata.Description)
			authorLower := strings.ToLower(metadata.Author)

			if (query == "" || 
				strings.Contains(nameLower, queryLower) || 
				strings.Contains(descLower, queryLower) || 
				strings.Contains(authorLower, queryLower)) {
				
				// Check tags if specified
				if len(tags) > 0 {
					tagMatch := false
					for _, tag := range tags {
						for _, mtag := range metadata.Tags {
							if strings.EqualFold(tag, mtag) {
								tagMatch = true
								break
							}
						}
						if tagMatch {
							break
						}
					}
					if !tagMatch {
						return nil
					}
				}

				results = append(results, metadata)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search registry: %w", err)
	}

	return results, nil
}

// ListAgents lists all agents in the registry
func (c *RegistryClient) ListAgents() ([]AgentMetadata, error) {
	return c.SearchAgents("", nil)
}

// Helper functions

// copyDir copies a directory and its contents
func copyDir(src, dst string) error {
	// Get file info
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if info.IsDir() {
		err = os.MkdirAll(dst, info.Mode())
		if err != nil {
			return err
		}

		// Read source directory
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}

		// Copy each entry
		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			dstPath := filepath.Join(dst, entry.Name())

			// Skip state files
			if strings.HasSuffix(entry.Name(), ".state.json") {
				continue
			}

			if entry.IsDir() {
				err = copyDir(srcPath, dstPath)
				if err != nil {
					return err
				}
			} else {
				err = copyFile(srcPath, dstPath)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return copyFile(src, dst)
	}

	return nil
}

// copyFile copies a file
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get file info
	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy file contents
	_, err = dstFile.ReadFrom(srcFile)
	return err
}

// readMetadata reads and parses a metadata file
func readMetadata(path string) (AgentMetadata, error) {
	var metadata AgentMetadata

	data, err := os.ReadFile(path)
	if err != nil {
		return metadata, err
	}

	err = json.Unmarshal(data, &metadata)
	return metadata, err
}

// getUsername returns the current username
func getUsername() string {
	// Try to get username from environment variables
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	if username == "" {
		username = "anonymous"
	}
	return username
}
