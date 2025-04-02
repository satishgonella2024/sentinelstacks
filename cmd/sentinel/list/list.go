package list

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/pkg/agent"
)

// NewListCmd creates a new list command
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List agents",
		Long:    `List all available agent images and running agents.`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			showAll, _ := cmd.Flags().GetBool("all")
			quiet, _ := cmd.Flags().GetBool("quiet")
			return runList(showAll, quiet)
		},
	}

	// Add flags
	cmd.Flags().BoolP("all", "a", false, "Show all agents (default shows just running)")
	cmd.Flags().BoolP("quiet", "q", false, "Only display agent IDs")

	return cmd
}

func runList(showAll, quiet bool) error {
	// Get running agents
	agents, err := agent.ListAgents()
	if err != nil {
		return fmt.Errorf("failed to list agents: %w", err)
	}

	// Filter if not showing all
	if !showAll {
		var runningAgents []agent.AgentStatus
		for _, a := range agents {
			if a.State == "running" {
				runningAgents = append(runningAgents, a)
			}
		}
		agents = runningAgents
	}

	// Handle quiet mode
	if quiet {
		for _, a := range agents {
			fmt.Println(a.ID)
		}
		return nil
	}

	// Create a tabwriter for nicely formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	// Print headers
	fmt.Fprintln(w, "AGENT ID\tNAME\tIMAGE\tSTATUS\tCREATED\tPORT")

	// Print each agent's information
	for _, a := range agents {
		created := time.Unix(a.CreatedAt, 0).Format("2006-01-02 15:04:05")
		
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\n",
			a.ID[:12],     // Show shortened ID for readability
			a.Name,
			a.Image,
			a.State,
			created,
			a.Port,
		)
	}

	return nil
}
