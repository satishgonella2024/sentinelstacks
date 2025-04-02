package ps

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/runtime"
)

// RuntimeInterface defines the minimum set of runtime methods needed by the ps command
type RuntimeInterface interface {
	GetRunningAgents() ([]runtime.AgentInfo, error)
}

// AgentInfo is a simplified version of runtime.AgentInfo for the ps command
type AgentInfo struct {
	ID        string    // Unique identifier for the agent
	Name      string    // Name of the agent
	Image     string    // Image used to create the agent
	Status    string    // Current status of the agent
	CreatedAt time.Time // When the agent was created
	Model     string    // LLM model being used
}

// NewPsCmd creates a new ps command
func NewPsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: "List running agents",
		Long:  `List all agents that are currently running on the local system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			quiet, _ := cmd.Flags().GetBool("quiet")
			format, _ := cmd.Flags().GetString("format")

			// Get the runtime
			rt, err := runtime.GetRuntime()
			if err != nil {
				return fmt.Errorf("failed to get runtime: %w", err)
			}

			return runPs(rt, all, quiet, format)
		},
	}

	// Add flags
	cmd.Flags().BoolP("all", "a", false, "Show all agents (default shows just running)")
	cmd.Flags().BoolP("quiet", "q", false, "Only display agent IDs")
	cmd.Flags().StringP("format", "f", "", "Format the output using a custom format")

	return cmd
}

// runPs executes the ps command
func runPs(rt RuntimeInterface, all, quiet bool, format string) error {
	// Get agents
	runtimeAgents, err := rt.GetRunningAgents()
	if err != nil {
		return fmt.Errorf("failed to get agents: %w", err)
	}

	// Convert runtime.AgentInfo to our local AgentInfo
	var agents []AgentInfo
	for _, agent := range runtimeAgents {
		agents = append(agents, AgentInfo{
			ID:        agent.ID,
			Name:      agent.Name,
			Image:     agent.Image,
			Status:    agent.Status,
			CreatedAt: agent.CreatedAt,
			Model:     agent.Model,
		})
	}

	// Filter agents if not showing all
	if !all {
		var runningAgents []AgentInfo
		for _, agent := range agents {
			if agent.Status == "running" {
				runningAgents = append(runningAgents, agent)
			}
		}
		agents = runningAgents
	}

	// Handle quiet mode (only show IDs)
	if quiet {
		for _, agent := range agents {
			fmt.Println(agent.ID)
		}
		return nil
	}

	// Handle custom format
	if format != "" {
		return formatOutput(agents, format)
	}

	// Default output format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "AGENT ID\tNAME\tIMAGE\tSTATUS\tCREATED\tMODEL")

	for _, agent := range agents {
		createdTime := formatTime(agent.CreatedAt)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			agent.ID[:12],
			agent.Name,
			agent.Image,
			agent.Status,
			createdTime,
			agent.Model,
		)
	}

	return w.Flush()
}

// formatOutput applies a custom format to the output
func formatOutput(agents []AgentInfo, format string) error {
	for _, agent := range agents {
		line := format
		// Replace placeholders with actual values
		line = strings.ReplaceAll(line, "{{.ID}}", agent.ID)
		line = strings.ReplaceAll(line, "{{.Name}}", agent.Name)
		line = strings.ReplaceAll(line, "{{.Image}}", agent.Image)
		line = strings.ReplaceAll(line, "{{.Status}}", agent.Status)
		line = strings.ReplaceAll(line, "{{.CreatedAt}}", agent.CreatedAt.Format(time.RFC3339))
		line = strings.ReplaceAll(line, "{{.Model}}", agent.Model)

		fmt.Println(line)
	}

	return nil
}

// formatTime formats the time relative to now
func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "Less than a minute ago"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, plural(minutes))
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, plural(hours))
	} else if diff < 48*time.Hour {
		return "Yesterday"
	} else {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, plural(days))
	}
}

// plural returns "s" if the number is not 1
func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
