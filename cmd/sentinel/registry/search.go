package registry

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/auth"
	"github.com/satishgonella2024/sentinelstacks/internal/registry/client"
)

// NewSearchCmd creates a new search command
func NewSearchCmd() *cobra.Command {
	var (
		limit int
	)

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for packages in the registry",
		Long:  `Search for packages in the registry matching a query.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get query
			query := args[0]

			// Get registry URL from flag or env var
			registryURL, _ := cmd.Flags().GetString("registry")
			if registryURL == "" {
				registryURL = os.Getenv("SENTINEL_REGISTRY_URL")
			}
			if registryURL == "" {
				registryURL = "https://registry.sentinelstacks.io"
			}

			// Create auth provider
			authConfig := auth.AuthConfig{
				RegistryURL: registryURL,
			}
			authProvider, err := auth.NewFileTokenProvider(authConfig)
			if err != nil {
				return fmt.Errorf("failed to create auth provider: %w", err)
			}

			// Create registry client
			clientConfig := client.Config{
				BaseURL:      registryURL,
				AuthProvider: authProvider,
			}
			registryClient := client.NewClient(clientConfig)

			// Search for packages
			fmt.Printf("Searching for packages matching '%s' in registry %s...\n", query, registryURL)
			results, err := registryClient.Search(context.Background(), query, limit)
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}

			// Display results
			if len(results) == 0 {
				fmt.Println("No packages found matching the query.")
				return nil
			}

			// Initialize tabwriter
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION\tAUTHOR\tCREATED")
			for _, result := range results {
				// Format creation date
				createdStr := result.CreatedAt.Format("2006-01-02")

				// Format description (limit length and remove newlines)
				desc := result.Description
				if len(desc) > 50 {
					desc = desc[:47] + "..."
				}
				desc = strings.ReplaceAll(desc, "\n", " ")

				// Print row
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					result.Name,
					result.Version,
					desc,
					result.Author,
					createdStr,
				)
			}
			w.Flush()

			return nil
		},
	}

	// Add flags
	cmd.Flags().IntVarP(&limit, "limit", "n", 10, "Maximum number of results to return")
	cmd.Flags().String("registry", "", "Registry URL (default: https://registry.sentinelstacks.io)")

	return cmd
}
