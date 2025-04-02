package importcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/pkg/agent"
)

// NewImportCmd creates a new import command
func NewImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [path] [agent_name]",
		Short: "Import agent state from a file",
		Long:  `Import an agent's state, memory, and configuration from a file to restore or transfer it.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			
			// Optional agent name
			var agentName string
			if len(args) > 1 {
				agentName = args[1]
			}
			
			return runImport(path, agentName)
		},
	}

	return cmd
}

func runImport(path string, agentName string) error {
	// Import the agent state
	agentID, err := agent.ImportAgentState(path, agentName)
	if err != nil {
		return fmt.Errorf("failed to import agent state: %w", err)
	}

	fmt.Printf("Successfully imported agent from %s with ID %s\n", path, agentID)
	return nil
}
