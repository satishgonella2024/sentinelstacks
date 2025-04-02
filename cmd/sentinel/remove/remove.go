package remove

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/pkg/agent"
)

// NewRemoveCmd creates a new remove command
func NewRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [agent_id...]",
		Aliases: []string{"rm"},
		Short:   "Remove one or more agents",
		Long:    `Remove one or more agents. By default, you cannot remove a running agent (use -f to override).`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			return runRemove(args, force)
		},
	}

	// Add flags
	cmd.Flags().BoolP("force", "f", false, "Force the removal of a running agent")

	return cmd
}

func runRemove(agentIDs []string, force bool) error {
	for _, id := range agentIDs {
		// Get agent status
		status, err := agent.GetAgentStatus(id)
		if err != nil {
			return fmt.Errorf("failed to get agent status: %w", err)
		}

		// Check if agent is running and force is not enabled
		if status.State == "running" && !force {
			return fmt.Errorf("cannot remove running agent %s (use -f to force)", id)
		}

		// Stop the agent first if it's running
		if status.State == "running" {
			if err := agent.StopAgent(id); err != nil {
				return fmt.Errorf("failed to stop agent %s: %w", id, err)
			}
		}

		// Remove the agent
		if err := agent.RemoveAgent(id); err != nil {
			return fmt.Errorf("failed to remove agent %s: %w", id, err)
		}

		fmt.Printf("Removed agent: %s\n", id)
	}

	return nil
}
