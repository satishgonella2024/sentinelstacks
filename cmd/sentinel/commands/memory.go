package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/satishgonella2024/sentinelstacks/internal/memory"
	"github.com/spf13/cobra"
)

// memoryCmd represents the memory command
var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage agent memory",
	Long: `Manage and inspect agent memory.
	
This command allows you to view, search, and manipulate the memory of an agent.`,
}

// memoryListCmd represents the memory list command
var memoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List memory entries",
	Long:  `List all memory entries for an agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		agentName, _ := cmd.Flags().GetString("agent")
		limit, _ := cmd.Flags().GetInt("limit")
		format, _ := cmd.Flags().GetString("format")

		// Load agent from registry
		agent, err := loadAgent(agentName)
		if err != nil {
			color.Red("Error loading agent: %v", err)
			return
		}

		// Get memory
		memoryInstance, err := memory.NewMemory(agent.Name, convertToInternalMemoryConfig(agent.Memory))
		if err != nil {
			color.Red("Error initializing memory: %v", err)
			return
		}

		// List entries
		entries, err := memoryInstance.List(limit)
		if err != nil {
			color.Red("Error listing memory entries: %v", err)
			return
		}

		if len(entries) == 0 {
			color.Yellow("No memory entries found.")
			return
		}

		// Sort entries by timestamp (newest first)
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.After(entries[j].Timestamp)
		})

		// Display results
		if format == "json" {
			displayMemoryEntriesJson(entries)
		} else {
			displayMemoryEntriesTable(entries)
		}
	},
}

// memorySearchCmd represents the memory search command
var memorySearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search memory entries",
	Long:  `Search memory entries using semantic or keyword search.`,
	Run: func(cmd *cobra.Command, args []string) {
		agentName, _ := cmd.Flags().GetString("agent")
		query, _ := cmd.Flags().GetString("query")
		limit, _ := cmd.Flags().GetInt("limit")
		format, _ := cmd.Flags().GetString("format")

		if query == "" {
			color.Red("Error: query is required")
			return
		}

		// Load agent from registry
		agent, err := loadAgent(agentName)
		if err != nil {
			color.Red("Error loading agent: %v", err)
			return
		}

		// Get memory
		memoryInstance, err := memory.NewMemory(agent.Name, convertToInternalMemoryConfig(agent.Memory))
		if err != nil {
			color.Red("Error initializing memory: %v", err)
			return
		}

		// Search entries
		entries, err := memoryInstance.Search(query, limit)
		if err != nil {
			color.Red("Error searching memory entries: %v", err)
			return
		}

		if len(entries) == 0 {
			color.Yellow("No matching entries found.")
			return
		}

		// Display results
		if format == "json" {
			displayMemoryEntriesJson(entries)
		} else {
			displayMemoryEntriesTable(entries)
		}
	},
}

// memoryClearCmd represents the memory clear command
var memoryClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear agent memory",
	Long:  `Clear all memory entries for an agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		agentName, _ := cmd.Flags().GetString("agent")
		confirm, _ := cmd.Flags().GetBool("confirm")

		if !confirm {
			fmt.Printf("Are you sure you want to clear all memory for agent '%s'? [y/N] ", agentName)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				color.Yellow("Operation cancelled.")
				return
			}
		}

		// Load agent from registry
		agent, err := loadAgent(agentName)
		if err != nil {
			color.Red("Error loading agent: %v", err)
			return
		}

		// Get memory
		memoryInstance, err := memory.NewMemory(agent.Name, convertToInternalMemoryConfig(agent.Memory))
		if err != nil {
			color.Red("Error initializing memory: %v", err)
			return
		}

		// Clear memory
		if err := memoryInstance.Clear(); err != nil {
			color.Red("Error clearing memory: %v", err)
			return
		}

		color.Green("Memory cleared successfully.")
	},
}

// Helper function to convert between memory config formats
func convertToInternalMemoryConfig(config MemoryConfig) memory.MemoryConfig {
	var memType memory.MemoryType
	switch config.Type {
	case "simple":
		memType = memory.SimpleMemory
	case "vector":
		memType = memory.VectorMemory
	default:
		memType = memory.SimpleMemory
	}

	return memory.MemoryConfig{
		Type:           memType,
		Persistence:    config.Persistence,
		MaxItems:       config.MaxItems,
		EmbeddingModel: config.EmbeddingModel,
	}
}

// Helper function to display memory entries in a table
func displayMemoryEntriesTable(entries []memory.MemoryEntry) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Content", "Timestamp", "Metadata"})
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	for _, entry := range entries {
		// Truncate content if too long
		content := entry.Content
		if len(content) > 50 {
			content = content[:47] + "..."
		}

		// Format timestamp
		timestamp := entry.Timestamp.Format(time.RFC3339)

		// Format metadata
		metadata := ""
		if entry.Metadata != nil {
			// Pretty print metadata, excluding some internal fields
			metadataItems := []string{}
			for k, v := range entry.Metadata {
				if k != "vector_id" {
					metadataItems = append(metadataItems, fmt.Sprintf("%s: %v", k, v))
				}
			}
			metadata = strings.Join(metadataItems, ", ")
		}

		// Add row to table
		table.Append([]string{entry.ID, content, timestamp, metadata})
	}

	table.Render()
}

// Helper function to display memory entries as JSON
func displayMemoryEntriesJson(entries []memory.MemoryEntry) {
	outputJson(entries)
}

func init() {
	rootCmd.AddCommand(memoryCmd)
	memoryCmd.AddCommand(memoryListCmd)
	memoryCmd.AddCommand(memorySearchCmd)
	memoryCmd.AddCommand(memoryClearCmd)

	// Flags for all memory commands
	memoryCmd.PersistentFlags().String("agent", "", "Agent name")
	memoryCmd.MarkPersistentFlagRequired("agent")

	// Flags for list command
	memoryListCmd.Flags().Int("limit", 10, "Maximum number of entries to return")
	memoryListCmd.Flags().String("format", "table", "Output format (table or json)")

	// Flags for search command
	memorySearchCmd.Flags().String("query", "", "Search query")
	memorySearchCmd.Flags().Int("limit", 10, "Maximum number of entries to return")
	memorySearchCmd.Flags().String("format", "table", "Output format (table or json)")
	memorySearchCmd.MarkFlagRequired("query")

	// Flags for clear command
	memoryClearCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")
}
