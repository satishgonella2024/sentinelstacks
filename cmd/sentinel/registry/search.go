package registry

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/pkg/registry"
)

// NewSearchCmd creates a new search command
func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for agent images in a registry",
		Long:  `Search for agent images in a registry based on a query string.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			registryURL, _ := cmd.Flags().GetString("registry")
			limit, _ := cmd.Flags().GetInt("limit")
			return runSearch(args[0], registryURL, limit)
		},
	}

	// Add flags
	cmd.Flags().StringP("registry", "r", registry.DefaultRegistryURL, "Registry URL to search")
	cmd.Flags().IntP("limit", "n", 25, "Maximum number of results to return")

	return cmd
}

func runSearch(query, registryURL string, limit int) error {
	// Search for images
	results, err := registry.SearchImages(registryURL, query, limit)
	if err != nil {
		return fmt.Errorf("failed to search registry: %w", err)
	}

	// Display results
	if len(results) == 0 {
		fmt.Println("No results found for query:", query)
		return nil
	}

	// Create a tabwriter for formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	// Print headers
	fmt.Fprintln(w, "NAME\tDESCRIPTION\tSTARS\tOFFICIAL\tAUTOMATED")

	// Print each result
	for _, result := range results {
		official := ""
		if result.IsOfficial {
			official = "[OK]"
		}

		automated := ""
		if result.IsAutomated {
			automated = "[OK]"
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
			result.Name,
			result.Description,
			result.Stars,
			official,
			automated,
		)
	}

	return nil
}
