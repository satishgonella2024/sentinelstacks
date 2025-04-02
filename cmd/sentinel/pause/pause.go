package pause

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/pkg/agent"
)

// NewPauseCmd creates a new pause command
func NewPauseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause [agent_id...]",
		Short: "Pause one or more running agents",
		Long:  `Pause one or more running agents. This suspends processing but maintains state.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPause(args)
		},
	}

	return cmd
}

func runPause(agentIDs []string) error {
	for _, id := range agentIDs {
		// Get agent status
		status, err := agent.GetAgentStatus(id)
		if err != nil {
			return fmt.Errorf("failed to get agent status: %w", err)
		}

		// Check if agent is running
		if status.State != "running" {
			return fmt.Errorf("cannot pause agent %s (not running)", id)
		}

		// Pause the agent
		if err := agent.PauseAgent(id); err != nil {
			return fmt.Errorf("failed to pause agent %s: %w", id, err)
		}

		fmt.Printf("Paused agent: %s\n", id)
	}

	return nil
}
