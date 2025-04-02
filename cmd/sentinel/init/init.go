package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Available templates
var templates = map[string]string{
	"default": `# Sentinelfile for MyAgent

Create an agent that [describe your agent's purpose].

The agent should be able to:
- [Capability 1]
- [Capability 2]
- [Capability 3]

The agent should use claude-3.7-sonnet as its base model.

It should maintain state about [state description].

When the conversation starts, the agent should [initialization behavior].

Allow the agent to access the following tools:
- [Tool 1]
- [Tool 2]

Set [parameter name] to [value].
`,
	"chatbot": `# Sentinelfile for BasicChatbot

Create an agent that provides helpful, friendly chat responses with a unique personality.

The agent should be able to:
- Engage in casual conversation
- Answer general knowledge questions
- Remember context from earlier in the conversation
- Provide thoughtful and nuanced responses
- Maintain a consistent personality
- Gracefully handle inappropriate requests

The agent should use llama3 as its base model.

It should maintain state about the conversation history and user preferences.

When the conversation starts, the agent should introduce itself as a friendly assistant and ask how it can help the user today.

When the conversation ends, the agent should thank the user and offer assistance for next time.

Allow the agent to access the following tools:
- Web search (for factual information)
- Calculator (for mathematical operations)
- Date/time (for time-related queries)

Set personality to friendly.
Set response_length to medium.
Set memory_depth to 10.
`,
	"assistant": `# Sentinelfile for ResearchAssistant

Create an agent that helps users conduct research, find information, and synthesize knowledge.

The agent should be able to:
- Search for relevant information on a topic
- Summarize findings in a clear, concise manner
- Cite sources accurately
- Answer follow-up questions about the research
- Organize information into logical categories
- Identify gaps in research or understanding

The agent should use claude-3-opus-20240229 as its base model.

It should maintain state about the research queries, results, and user preferences.

When the conversation starts, the agent should introduce itself as a research assistant and ask what topic the user would like to explore.

When the conversation ends, the agent should summarize the key findings and offer to continue the research in the future.

Allow the agent to access the following tools:
- Web search (for finding information)
- Document parser (for processing PDFs, docs, etc.)
- Citation generator (for creating properly formatted citations)
- Knowledge graph (for connecting related concepts)

Set detail_level to high.
Set citation_style to APA.
Set search_depth to comprehensive.
`,
	"analyzer": `# Sentinelfile for DataAnalyzer

Create an agent that analyzes data, identifies patterns, and presents insights.

The agent should be able to:
- Process and clean structured data
- Generate statistical analyses
- Create visualizations to represent data
- Identify trends and anomalies
- Provide explanations of analytical findings
- Make recommendations based on data insights

The agent should use claude-3-opus-20240229 as its base model.

It should maintain state about the datasets, analyses performed, and visualization preferences.

When the conversation starts, the agent should introduce itself as a data analyzer and ask what data the user would like to explore.

When the conversation ends, the agent should summarize the key insights and offer to save the analysis for future reference.

Allow the agent to access the following tools:
- Data processor (for importing and cleaning data)
- Statistical engine (for advanced statistical tests)
- Visualization generator (for creating charts and graphs)
- Export tool (for saving results in various formats)

Set analysis_depth to comprehensive.
Set explanation_style to detailed.
Set visualization_quality to high.
`}

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	var (
		name     string
		template string
	)

	initCmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize a new Sentinelfile",
		Long:  `Create a new Sentinelfile in the current directory or specified directory`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine the directory
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			// Override with --name if provided
			if name != "" {
				dir = name
			}

			// Create the directory if it doesn't exist
			if dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
			}

			// Determine which template to use
			templateContent, ok := templates["default"]
			if template != "" {
				templateContent, ok = templates[strings.ToLower(template)]
				if !ok {
					// List available templates
					var availableTemplates []string
					for k := range templates {
						if k != "default" {
							availableTemplates = append(availableTemplates, k)
						}
					}
					return fmt.Errorf("template '%s' not found. Available templates: %s",
						template, strings.Join(availableTemplates, ", "))
				}
			}

			// Create the Sentinelfile
			filename := filepath.Join(dir, "Sentinelfile")
			if _, err := os.Stat(filename); err == nil {
				return fmt.Errorf("file %s already exists", filename)
			}

			if err := os.WriteFile(filename, []byte(templateContent), 0644); err != nil {
				return fmt.Errorf("failed to write Sentinelfile: %w", err)
			}

			fmt.Printf("Created Sentinelfile at %s\n", filename)
			fmt.Printf("To build your agent, run: sentinel build -t myname/%s:latest\n",
				filepath.Base(dir))

			return nil
		},
	}

	initCmd.Flags().StringVar(&name, "name", "", "Name for the agent directory")
	initCmd.Flags().StringVar(&template, "template", "", "Template to use (chatbot, assistant, analyzer)")

	return initCmd
}
