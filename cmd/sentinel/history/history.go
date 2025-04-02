package history

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// NewHistoryCmd creates the history command
func NewHistoryCmd() *cobra.Command {
	var (
		limit  int
		agent  string
		all    bool
		format string
	)

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show agent execution history",
		Long:  `Display the history of agent executions and interactions`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Showing agent execution history:")
			
			// TODO: Implement actual history retrieval
			// For now, simulate history results
			
			// Simulate history entries
			history := simulateHistory(agent, limit)
			
			if len(history) == 0 {
				fmt.Println("No history found")
				return nil
			}
			
			// Display history based on format
			switch format {
			case "table":
				displayTableFormat(history, all)
			case "detailed":
				displayDetailedFormat(history, all)
			default:
				displayTableFormat(history, all)
			}
			
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of history entries to show")
	cmd.Flags().StringVar(&agent, "agent", "", "Filter history by agent ID")
	cmd.Flags().BoolVar(&all, "all", false, "Show all history entries (including system events)")
	cmd.Flags().StringVar(&format, "format", "table", "Output format (table, detailed)")

	return cmd
}

// historyEntry represents an entry in the execution history
type historyEntry struct {
	timestamp time.Time
	agentID   string
	agentName string
	action    string
	status    string
	details   string
	duration  time.Duration
	isSystem  bool
}

// simulateHistory simulates history entries
func simulateHistory(agentFilter string, limit int) []historyEntry {
	now := time.Now()
	history := []historyEntry{
		{
			timestamp: now.Add(-30 * time.Minute),
			agentID:   "agent123",
			agentName: "research-assistant",
			action:    "run",
			status:    "completed",
			details:   "Executed research task on AI ethics",
			duration:  15 * time.Minute,
			isSystem:  false,
		},
		{
			timestamp: now.Add(-45 * time.Minute),
			agentID:   "agent456",
			agentName: "code-reviewer",
			action:    "run",
			status:    "completed",
			details:   "Reviewed pull request #123",
			duration:  5 * time.Minute,
			isSystem:  false,
		},
		{
			timestamp: now.Add(-1 * time.Hour),
			agentID:   "agent123",
			agentName: "research-assistant",
			action:    "build",
			status:    "completed",
			details:   "Built from Sentinelfile",
			duration:  30 * time.Second,
			isSystem:  false,
		},
		{
			timestamp: now.Add(-1*time.Hour - 15*time.Minute),
			agentID:   "system",
			agentName: "registry",
			action:    "pull",
			status:    "completed",
			details:   "Pulled base image researcher:v1",
			duration:  10 * time.Second,
			isSystem:  true,
		},
		{
			timestamp: now.Add(-2 * time.Hour),
			agentID:   "agent789",
			agentName: "customer-support",
			action:    "run",
			status:    "failed",
			details:   "Error accessing knowledge base",
			duration:  2 * time.Minute,
			isSystem:  false,
		},
	}
	
	// Filter by agent ID if specified
	if agentFilter != "" {
		filtered := []historyEntry{}
		for _, entry := range history {
			if entry.agentID == agentFilter {
				filtered = append(filtered, entry)
			}
		}
		history = filtered
	}
	
	// Filter system events if not showing all
	if !all {
		filtered := []historyEntry{}
		for _, entry := range history {
			if !entry.isSystem {
				filtered = append(filtered, entry)
			}
		}
		history = filtered
	}
	
	// Limit the number of entries
	if len(history) > limit {
		history = history[:limit]
	}
	
	return history
}

// displayTableFormat displays history in table format
func displayTableFormat(history []historyEntry, showAll bool) {
	fmt.Println("TIME                  AGENT               ACTION     STATUS     DURATION")
	fmt.Println("----------------------------------------------------------------------")
	
	for _, entry := range history {
		// Format the timestamp
		timeStr := entry.timestamp.Format("2006-01-02 15:04:05")
		
		// Format agent name (show system tag if applicable)
		agentStr := entry.agentName
		if entry.isSystem {
			agentStr += " [SYSTEM]"
		}
		
		// Format duration
		durationStr := formatDuration(entry.duration)
		
		// Print the row
		fmt.Printf("%-20s %-20s %-10s %-10s %-10s\n",
			timeStr,
			truncateString(agentStr, 18),
			entry.action,
			entry.status,
			durationStr)
	}
}

// displayDetailedFormat displays history in detailed format
func displayDetailedFormat(history []historyEntry, showAll bool) {
	for i, entry := range history {
		fmt.Printf("Entry %d:\n", i+1)
		fmt.Printf("  Time:     %s\n", entry.timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Agent ID: %s\n", entry.agentID)
		fmt.Printf("  Name:     %s", entry.agentName)
		if entry.isSystem {
			fmt.Print(" [SYSTEM]")
		}
		fmt.Println()
		fmt.Printf("  Action:   %s\n", entry.action)
		fmt.Printf("  Status:   %s\n", entry.status)
		fmt.Printf("  Details:  %s\n", entry.details)
		fmt.Printf("  Duration: %s\n", formatDuration(entry.duration))
		fmt.Println()
	}
}

// formatDuration formats a duration in a human-readable format
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	} else if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%dm%.0fs", int(d.Minutes()), d.Seconds()-float64(int(d.Minutes()))*60)
	} else {
		return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
	}
}

// truncateString truncates a string to the specified length and adds ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
