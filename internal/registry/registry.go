package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Agent struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var (
	configDir = filepath.Join(os.Getenv("HOME"), ".sentinel")
	authFile  = filepath.Join(configDir, "auth.json")
)

type AuthConfig struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func init() {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
	}
}

func Login(server, username, password string) error {
	// In a real implementation, this would make an API call to get a token
	token := "dummy-token" // This should come from the server

	auth := AuthConfig{
		Server:   server,
		Username: username,
		Token:    token,
	}

	data, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("error marshaling auth config: %v", err)
	}

	if err := os.WriteFile(authFile, data, 0600); err != nil {
		return fmt.Errorf("error writing auth config: %v", err)
	}

	return nil
}

func Push(agentName, version string) error {
	// Read auth config
	auth, err := readAuthConfig()
	if err != nil {
		return fmt.Errorf("not logged in: %v", err)
	}

	// In a real implementation, this would:
	// 1. Package the agent
	// 2. Upload to the registry
	// 3. Update metadata

	fmt.Printf("Pushing agent %s:%s to %s\n", agentName, version, auth.Server)
	return nil
}

func Pull(agentName, version string) error {
	// Read auth config
	auth, err := readAuthConfig()
	if err != nil {
		return fmt.Errorf("not logged in: %v", err)
	}

	// In a real implementation, this would:
	// 1. Download from the registry
	// 2. Verify the package
	// 3. Install locally

	fmt.Printf("Pulling agent %s:%s from %s\n", agentName, version, auth.Server)
	return nil
}

func List() ([]Agent, error) {
	agentsDir := filepath.Join(os.Getenv("HOME"), ".sentinel", "agents")

	// Read the agents directory
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Agent{}, nil
		}
		return nil, fmt.Errorf("error reading agents directory: %v", err)
	}

	var agents []Agent
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Read versions directory
		versionsDir := filepath.Join(agentsDir, entry.Name())
		versions, err := os.ReadDir(versionsDir)
		if err != nil {
			continue
		}

		// Add each version as a separate agent
		for _, version := range versions {
			if version.IsDir() {
				agents = append(agents, Agent{
					Name:    entry.Name(),
					Version: version.Name(),
				})
			}
		}
	}

	return agents, nil
}

func readAuthConfig() (*AuthConfig, error) {
	data, err := os.ReadFile(authFile)
	if err != nil {
		return nil, fmt.Errorf("error reading auth config: %v", err)
	}

	var auth AuthConfig
	if err := json.Unmarshal(data, &auth); err != nil {
		return nil, fmt.Errorf("error unmarshaling auth config: %v", err)
	}

	return &auth, nil
}
