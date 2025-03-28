package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	agentsDir = filepath.Join(os.Getenv("HOME"), ".sentinel", "agents")
)

func init() {
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		fmt.Printf("Error creating agents directory: %v\n", err)
	}
}

func Run(agentName, version string) error {
	agentDir := filepath.Join(agentsDir, agentName, version)

	// Check if agent exists
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("agent %s:%s not found", agentName, version)
	}

	// In a real implementation, this would:
	// 1. Load the agent configuration
	// 2. Set up the environment
	// 3. Execute the agent
	// 4. Handle the results

	fmt.Printf("Running agent %s:%s\n", agentName, version)
	return nil
}
