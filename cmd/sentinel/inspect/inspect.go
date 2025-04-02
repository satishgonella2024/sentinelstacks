package inspect

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/sentinelstacks/sentinel/pkg/agent"
)

// NewInspectCmd creates a new inspect command
func NewInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect [agent_id]",
		Short: "Display detailed information about an agent",
		Long:  `Display detailed information about an agent, including configuration, state, and capabilities.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			return runInspect(args[0], format)
		},
	}

	// Add flags
	cmd.Flags().StringP("format", "f", "yaml", "Format the output (json or yaml)")

	return cmd
}

func runInspect(agentID string, format string) error {
	// Get agent details
	details, err := agent.GetAgentDetails(agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent details: %w", err)
	}

	// Format the output
	var outputBytes []byte
	switch strings.ToLower(format) {
	case "json":
		outputBytes, err = json.MarshalIndent(details, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal agent details to JSON: %w", err)
		}
	case "yaml":
		outputBytes, err = yaml.Marshal(details)
		if err != nil {
			return fmt.Errorf("failed to marshal agent details to YAML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s (use json or yaml)", format)
	}

	// Print the output
	fmt.Fprintln(os.Stdout, string(outputBytes))

	return nil
}
