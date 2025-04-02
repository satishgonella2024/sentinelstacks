package stop

import (
	"fmt"

	"github.com/satishgonella2024/sentinelstacks/internal/runtime"
	"github.com/spf13/cobra"
)

// NewStopCmd creates a new stop command
func NewStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [agent_id]",
		Short: "Stop a running agent",
		Long:  `Stop a running agent. The agent will be gracefully terminated.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]
			force, _ := cmd.Flags().GetBool("force")
			timeout, _ := cmd.Flags().GetInt("timeout")

			return runStop(agentID, force, timeout)
		},
	}

	// Add flags
	cmd.Flags().BoolP("force", "f", false, "Force stop the agent without graceful shutdown")
	cmd.Flags().IntP("timeout", "t", 30, "Timeout in seconds for graceful shutdown")

	return cmd
}

// runStop executes the stop command
func runStop(agentID string, force bool, timeout int) error {
	// Get the runtime
	runtime, err := runtime.GetRuntime()
	if err != nil {
		return fmt.Errorf("failed to get runtime: %w", err)
	}

	// Get agent information first
	agent, err := runtime.GetAgent(agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %s", err)
	}

	// Print stopping message
	fmt.Printf("Stopping agent %s (%s)...\n", agent.Name, agentID)

	// Stop the agent
	if err := runtime.StopAgent(agentID); err != nil {
		return fmt.Errorf("failed to stop agent: %w", err)
	}

	fmt.Printf("Agent %s stopped successfully\n", agentID)
	return nil
}
