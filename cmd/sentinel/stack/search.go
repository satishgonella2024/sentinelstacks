package stack

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/client"
	packages "github.com/satishgonella2024/sentinelstacks/internal/registry/package"
)

// NewSearchCommand creates a 'stack search' command
func NewSearchCommand() *cobra.Command {
	var (
		format string
		limit  int
	)

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for stacks in the registry",
		Long:  `Search for stacks in the registry with optional filtering by name, description, or labels`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get query string
			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			// Get registry URL from config
			registryURL := viper.GetString("registry.url")
			if registryURL == "" {
				registryURL = "https://registry.sentinelstacks.io"
			}

			// Get auth token from config
			authToken := viper.GetString("registry.auth_token")

			// Create registry client
			registryClient := client.NewRegistryClient(registryURL, authToken)

			fmt.Printf("Searching for stacks in registry...\n")

			// Search for stacks
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := registryClient.SearchPackages(ctx, query, packages.PackageTypeStack, limit)
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}

			// Display results
			if len(result.Items) == 0 {
				fmt.Printf("No stacks found matching '%s'\n", query)
				return nil
			}

			fmt.Printf("Found %d stacks matching '%s'\n", result.TotalCount, query)

			switch format {
			case "wide":
				printWideSearch(result)
			case "json":
				printJSONSearch(result)
			default:
				printDefaultSearch(result)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&format, "format", "f", "default", "Output format (default, wide, json)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Maximum number of results to display")

	return cmd
}

// printDefaultSearch prints search results in a concise format
func printDefaultSearch(result *client.SearchResult) {
	fmt.Printf("%-30s %-15s %-10s %-5s\n", "NAME", "VERSION", "AUTHOR", "DOWNLOADS")
	fmt.Println(strings.Repeat("-", 70))

	for _, item := range result.Items {
		author := item.Author
		if len(author) > 15 {
			author = author[:12] + "..."
		}

		fmt.Printf("%-30s %-15s %-10s %-5d\n",
			truncateString(item.Name, 30),
			truncateString(item.Version, 15),
			truncateString(author, 10),
			item.Downloads)
	}
}

// printWideSearch prints search results in a detailed format
func printWideSearch(result *client.SearchResult) {
	fmt.Printf("%-25s %-12s %-15s %-5s %-40s\n", "NAME", "VERSION", "AUTHOR", "DOWNLOADS", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 100))

	for _, item := range result.Items {
		fmt.Printf("%-25s %-12s %-15s %-5d %-40s\n",
			truncateString(item.Name, 25),
			truncateString(item.Version, 12),
			truncateString(item.Author, 15),
			item.Downloads,
			truncateString(item.Description, 40))
	}

	// Print each item's dependencies
	fmt.Println("\nDEPENDENCIES:")
	for _, item := range result.Items {
		if len(item.Dependencies) > 0 {
			fmt.Printf("\n%s:%s dependencies:\n", item.Name, item.Version)
			for _, dep := range item.Dependencies {
				fmt.Printf("  - %s:%s (%s)%s\n", 
					dep.Name, 
					dep.Version, 
					string(dep.Type),
					requiredString(dep.Required))
			}
		}
	}
}

// printJSONSearch prints search results in JSON format
func printJSONSearch(result *client.SearchResult) {
	// In a real implementation, this would use json.Marshal
	fmt.Println("{")
	fmt.Printf("  \"totalCount\": %d,\n", result.TotalCount)
	fmt.Println("  \"items\": [")

	for i, item := range result.Items {
		fmt.Printf("    {\n")
		fmt.Printf("      \"name\": \"%s\",\n", item.Name)
		fmt.Printf("      \"version\": \"%s\",\n", item.Version)
		fmt.Printf("      \"type\": \"%s\",\n", item.Type)
		fmt.Printf("      \"description\": \"%s\",\n", escapeString(item.Description))
		fmt.Printf("      \"author\": \"%s\",\n", escapeString(item.Author))
		fmt.Printf("      \"createdAt\": \"%s\",\n", item.CreatedAt.Format(time.RFC3339))
		fmt.Printf("      \"downloads\": %d,\n", item.Downloads)
		fmt.Printf("      \"verified\": %t", item.Verified)

		if len(item.Dependencies) > 0 {
			fmt.Printf(",\n      \"dependencies\": [\n")
			for j, dep := range item.Dependencies {
				fmt.Printf("        {\n")
				fmt.Printf("          \"name\": \"%s\",\n", dep.Name)
				fmt.Printf("          \"version\": \"%s\",\n", dep.Version)
				fmt.Printf("          \"type\": \"%s\",\n", dep.Type)
				fmt.Printf("          \"required\": %t\n", dep.Required)
				if j < len(item.Dependencies)-1 {
					fmt.Printf("        },\n")
				} else {
					fmt.Printf("        }\n")
				}
			}
			fmt.Printf("      ]\n")
		} else {
			fmt.Printf("\n")
		}

		if i < len(result.Items)-1 {
			fmt.Printf("    },\n")
		} else {
			fmt.Printf("    }\n")
		}
	}

	fmt.Println("  ]")
	fmt.Println("}")
}

// truncateString truncates a string to max length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// escapeString escapes special characters in a string for JSON
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// requiredString returns a string indicating if a dependency is required
func requiredString(required bool) string {
	if required {
		return " [required]"
	}
	return " [optional]"
}
