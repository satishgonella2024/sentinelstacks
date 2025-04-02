package resume

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/pkg/agent"
)

// NewResumeCmd creates a new resume command
func NewResumeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume [agent_id...]",
		Short: "Resume one or more paused agents",
		Long:  `Resume one or more paused agents, allowing them to continue processing.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResume(args)
		},
	}

	return cmd
}

func runResume(agentIDs []string) error {
	for _, id := range agentIDs {
		// Get agent status
		status, err := agent.GetAgentStatus(id)
		if err != nil {
			return fmt.Errorf("failed to get agent status: %w", err)
		}

		// Check if agent is paused
		if status.State != "paused" {
			return fmt.Errorf("cannot resume agent %s (not paused)", id)
		}

		// Resume the agent
		if err := agent.ResumeAgent(id); err != nil {
			return fmt.Errorf("failed to resume agent %s: %w", id, err)
		}

		fmt.Printf("Resumed agent: %s\n", id)
	}

	return nil
}
