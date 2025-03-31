package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generator represents the NLP-to-Agent generator
type Generator struct {
	LLMProvider   string
	LLMModel      string
	OutputDir     string
	AgentName     string
	TagName       string
	GenerateFiles bool
}

// NLPRequest represents user's natural language request
type NLPRequest struct {
	Input string `json:"input"`
}

// SentinelfileResponse represents the LLM response with generated Sentinelfile
type SentinelfileResponse struct {
	Sentinelfile string            `json:"sentinelfile"`
	Metadata     map[string]string `json:"metadata"`
}

// NewGenerator creates a new NLP-to-Agent generator
func NewGenerator(llmProvider, llmModel, outputDir string) *Generator {
	return &Generator{
		LLMProvider:   llmProvider,
		LLMModel:      llmModel,
		OutputDir:     outputDir,
		GenerateFiles: true,
	}
}

// ProcessNaturalLanguage converts NL to Sentinelfile YAML
func (g *Generator) ProcessNaturalLanguage(input string) (*SentinelfileResponse, error) {
	// Step 1: Extract agent name and tag from input if possible
	g.extractNameAndTag(input)

	// Step 2: Prepare prompt for LLM
	prompt := g.buildPrompt(input)

	// Step 3: Send to LLM (simulated here)
	fmt.Println("Sending to LLM:", g.LLMProvider, g.LLMModel)
	fmt.Println("Processing natural language request...")

	// This would be replaced with actual LLM API call
	yamlContent, metadata, err := g.callLLM(prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM processing error: %v", err)
	}

	response := &SentinelfileResponse{
		Sentinelfile: yamlContent,
		Metadata:     metadata,
	}

	// Step 4: Generate files if enabled
	if g.GenerateFiles {
		err = g.generateFiles(response)
		if err != nil {
			return response, fmt.Errorf("file generation error: %v", err)
		}
	}

	return response, nil
}

// extractNameAndTag attempts to extract agent name and tag from input
func (g *Generator) extractNameAndTag(input string) {
	// Simple heuristic for extracting name - would be enhanced by LLM
	input = strings.ToLower(input)
	words := strings.Fields(input)

	// Default name if we can't extract one
	g.AgentName = "generated-agent"
	g.TagName = "latest"

	// Look for phrases like "create a [name] agent" or similar patterns
	for i, word := range words {
		if (word == "called" || word == "named") && i < len(words)-1 {
			candidateName := words[i+1]
			// Clean the name
			candidateName = strings.Trim(candidateName, ".,\"'!?;:")
			if candidateName != "" {
				g.AgentName = candidateName
				break
			}
		}
	}
}

// buildPrompt creates the prompt for the LLM
func (g *Generator) buildPrompt(input string) string {
	return fmt.Sprintf(`
You are an expert at creating SentinelStacks agent definitions.
Convert the following natural language description into a valid Sentinelfile YAML.

The Sentinelfile should follow this structure:
- name: A short name for the agent
- description: A concise description of what the agent does
- capabilities: List of what the agent can do
- model: Configuration for the LLM
- state: State variables to track
- tools: Tools that the agent can use
- initialization: How the agent introduces itself
- termination: How the agent ends a session

Include any other appropriate sections based on the user's requirements.

User's description: "%s"

Return ONLY the YAML content in valid format.
`, input)
}

// callLLM simulates calling an LLM API
func (g *Generator) callLLM(prompt string) (string, map[string]string, error) {
	// This is a simulation - in a real implementation, we'd call the actual LLM API
	// For this example, we'll just generate a sample Sentinelfile

	// This would be a response from the LLM
	yamlContent := fmt.Sprintf(`name: %s
description: Agent generated from natural language description
capabilities:
  - Understand user requests
  - Provide helpful responses
  - Maintain conversation context
model:
  base: %s
  parameters:
    temperature: 0.7
    top_p: 0.9
state:
  - conversation_history
  - user_preferences
initialization:
  introduction: "Hello! I'm a custom agent created from your natural language description."
termination:
  farewell: "Thank you for the conversation. Goodbye!"
tools:
  - web_search:
      purpose: For looking up information online
`, g.AgentName, strings.ToLower(strings.Replace(g.LLMModel, "-", "", -1)))

	metadata := map[string]string{
		"source":      "nlp_generator",
		"confidence":  "medium",
		"agent_name":  g.AgentName,
		"description": "Agent generated from natural language",
	}

	return yamlContent, metadata, nil
}

// generateFiles creates the necessary files for the agent
func (g *Generator) generateFiles(response *SentinelfileResponse) error {
	// Create output directory if needed
	agentDir := filepath.Join(g.OutputDir, g.AgentName)
	err := os.MkdirAll(agentDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write Sentinelfile
	sentinelfilePath := filepath.Join(agentDir, "Sentinelfile")
	err = os.WriteFile(sentinelfilePath, []byte(response.Sentinelfile), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Sentinelfile: %v", err)
	}

	// Write metadata
	metadataPath := filepath.Join(agentDir, "metadata.json")
	metadataBytes, err := json.MarshalIndent(response.Metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %v", err)
	}

	err = os.WriteFile(metadataPath, metadataBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %v", err)
	}

	return nil
}

// BuildAgent builds the agent using sentinel CLI
func (g *Generator) BuildAgent() error {
	agentDir := filepath.Join(g.OutputDir, g.AgentName)
	sentinelfilePath := filepath.Join(agentDir, "Sentinelfile")

	fmt.Printf("Building agent %s:%s from %s\n", g.AgentName, g.TagName, sentinelfilePath)

	// In a real implementation, this would call the sentinel CLI
	// For example:
	// cmd := exec.Command("sentinel", "build", "-t", g.AgentName+":"+g.TagName,
	//                     "-f", sentinelfilePath, "--llm", g.LLMProvider, "--llm-model", g.LLMModel)
	// return cmd.Run()

	fmt.Println("Agent built successfully (simulated)")
	return nil
}

func main() {
	// Simple flag parsing for this example
	demoMode := flag.String("mode", "interactive", "Demo mode: interactive, cli, package, or both")
	flag.Parse()

	fmt.Println("SentinelStacks NLP-to-Agent Generator")
	fmt.Println("====================================")

	switch *demoMode {
	case "interactive":
		// Original interactive mode
		fmt.Println("Describe the agent you want to create in natural language.")
		fmt.Println("Type your description and press Enter, then Ctrl+D (Unix) or Ctrl+Z (Windows) to finish:")

		// Read multiline input
		var input strings.Builder
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			input.WriteString(line)
			input.WriteString("\n")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		// Process the input
		generator := NewGenerator("anthropic", "claude-3-sonnet", "generated")
		response, err := generator.ProcessNaturalLanguage(input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Build the agent
		err = generator.BuildAgent()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building agent: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nGenerated Sentinelfile:")
		fmt.Println("----------------------")
		fmt.Println(response.Sentinelfile)

	case "cli":
		fmt.Println("Running CLI integration demo...")
		RunCLI()
	case "package":
		fmt.Println("Running package usage demo...")
		DemoAsPackage()
	case "both":
		fmt.Println("Running CLI integration demo...")
		RunCLI()
		fmt.Println("\n----------")
		fmt.Println("Running package usage demo...")
		DemoAsPackage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown demo mode: %s\n", *demoMode)
		fmt.Fprintf(os.Stderr, "Valid modes are: interactive, cli, package, both\n")
		os.Exit(1)
	}
}
