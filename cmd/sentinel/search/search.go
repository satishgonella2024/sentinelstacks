package search

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSearchCmd creates the search command
func NewSearchCmd() *cobra.Command {
	var (
		registry string
		limit    int
		filter   string
	)

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for agent images in a registry",
		Long:  `Search for agent images matching a query in a registry`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			
			// Determine registry URL
			registryURL := registry
			if registryURL == "" {
				registryURL = viper.GetString("registry.default")
				if registryURL == "" {
					registryURL = "sentinel.registry.ai" // Default registry
				}
			}
			
			fmt.Printf("Searching for '%s' in registry '%s'\n", query, registryURL)
			
			// TODO: Implement actual registry search
			// For now, simulate search results
			
			// Simulate some search results based on the query
			results := simulateSearchResults(query, filter, limit)
			
			// Display results
			if len(results) == 0 {
				fmt.Println("No results found")
				return nil
			}
			
			fmt.Printf("Found %d results:\n\n", len(results))
			fmt.Println("NAME                  DESCRIPTION                    STARS  PULLS   TIMESTAMP")
			fmt.Println("---------------------------------------------------------------------")
			
			for _, result := range results {
				fmt.Printf("%-20s %-30s %-6d %-7d %s\n",
					result.name,
					truncateString(result.description, 28),
					result.stars,
					result.pulls,
					result.updated)
			}
			
			return nil
		},
	}

	cmd.Flags().StringVar(&registry, "registry", "", "Registry URL to search in")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to return")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter results (official, stars, downloads)")

	return cmd
}

// searchResult represents a search result
type searchResult struct {
	name        string
	description string
	stars       int
	pulls       int
	updated     string
}

// simulateSearchResults simulates search results based on the query
func simulateSearchResults(query, filter string, limit int) []searchResult {
	results := []searchResult{
		{
			name:        "research-assistant",
			description: "Agent for academic research",
			stars:       245,
			pulls:       1820,
			updated:     "2023-10-15",
		},
		{
			name:        "code-reviewer",
			description: "Code review and analysis agent",
			stars:       187,
			pulls:       1243,
			updated:     "2023-10-25",
		},
		{
			name:        "customer-support",
			description: "Customer service agent",
			stars:       92,
			pulls:       780,
			updated:     "2023-09-30",
		},
		{
			name:        "data-analyst",
			description: "Data analysis and visualization",
			stars:       156,
			pulls:       980,
			updated:     "2023-10-20",
		},
		{
			name:        "content-writer",
			description: "Content creation agent",
			stars:       134,
			pulls:       870,
			updated:     "2023-10-12",
		},
	}
	
	// Filter results based on the query
	filtered := []searchResult{}
	for _, result := range results {
		if strings.Contains(result.name, query) || strings.Contains(result.description, query) {
			filtered = append(filtered, result)
		}
	}
	
	// Apply additional filters
	if filter != "" {
		switch filter {
		case "stars":
			// Sort by stars (simplified)
			for i := 0; i < len(filtered)-1; i++ {
				for j := i + 1; j < len(filtered); j++ {
					if filtered[i].stars < filtered[j].stars {
						filtered[i], filtered[j] = filtered[j], filtered[i]
					}
				}
			}
		case "downloads":
			// Sort by pulls (simplified)
			for i := 0; i < len(filtered)-1; i++ {
				for j := i + 1; j < len(filtered); j++ {
					if filtered[i].pulls < filtered[j].pulls {
						filtered[i], filtered[j] = filtered[j], filtered[i]
					}
				}
			}
		}
	}
	
	// Limit the number of results
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	
	return filtered
}

// truncateString truncates a string to the specified length and adds ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
