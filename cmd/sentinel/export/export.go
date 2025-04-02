package export

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/pkg/agent"
)

// NewExportCmd creates a new export command
func NewExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export [agent_id] [path]",
		Short: "Export agent state to a file",
		Long:  `Export an agent's state, memory, and configuration to a file for backup or transfer.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]
			
			// Default export path is [agent_id]-[timestamp].state.json in current directory
			exportPath := fmt.Sprintf("%s-%d.state.json", agentID, time.Now().Unix())
			if len(args) > 1 {
				exportPath = args[1]
			}
			
			return runExport(agentID, exportPath)
		},
	}

	return cmd
}

func runExport(agentID, exportPath string) error {
	// Ensure the directory exists
	dir := filepath.Dir(exportPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	// Export the agent state
	if err := agent.ExportAgentState(agentID, exportPath); err != nil {
		return fmt.Errorf("failed to export agent state: %w", err)
	}

	fmt.Printf("Successfully exported agent %s state to %s\n", agentID, exportPath)
	return nil
}
