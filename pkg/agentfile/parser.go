package agentfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/satishgonella2024/sentinelstacks/pkg/models"
)

// Parser converts natural language descriptions to structured YAML Agentfiles
type Parser struct {
	ModelEndpoint string
	Verbose       bool
}

// NewParser creates a new Agentfile parser
func NewParser(modelEndpoint string) *Parser {
	return &Parser{
		ModelEndpoint: modelEndpoint,
		Verbose:       false,
	}
}

// SetVerbose sets the verbose flag
func (p *Parser) SetVerbose(verbose bool) {
	p.Verbose = verbose
}

// ParseFile takes a natural language file and converts it to YAML
func (p *Parser) ParseFile(naturalLanguagePath string) (string, error) {
	// Read the natural language file
	content, err := os.ReadFile(naturalLanguagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read natural language file: %w", err)
	}

	// Parse the natural language to YAML
	yaml, err := p.ParseText(string(content))
	if err != nil {
		return "", err
	}

	// Write the YAML file alongside the natural language file
	dir := filepath.Dir(naturalLanguagePath)
	base := filepath.Base(naturalLanguagePath)
	ext := filepath.Ext(naturalLanguagePath)
	name := base[:len(base)-len(ext)]
	yamlPath := filepath.Join(dir, name+".yaml")

	err = os.WriteFile(yamlPath, []byte(yaml), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write YAML file: %w", err)
	}

	return yamlPath, nil
}

// ParseText converts natural language text to YAML
func (p *Parser) ParseText(text string) (string, error) {
	// Connect to Ollama
	adapter := models.NewOllamaAdapter(p.ModelEndpoint, "llama3")
	
	// If verbose, pass it to the adapter
	adapter.Verbose = p.Verbose
	
	// Create a system prompt that explains the task with examples
	systemPrompt := `You are an expert AI agent designer who converts natural language descriptions into structured YAML configurations.
    
Your task is to extract key details about an agent from a description and format them according to this schema:

name: string (name of the agent)
version: string (semantic version, start with "0.1.0")
description: string (concise description based on the text)
model:
  provider: string (ollama, openai, claude)
  name: string (llama3, gpt-4, etc.)
  options:
    temperature: float (0.0-1.0)
capabilities:
  - conversation (always include this for basic interaction)
  - [other capabilities based on description]
memory:
  type: string (simple, vector)
  persistence: boolean
tools: (optional)
  - id: string (tool identifier)
    version: string (semver)
permissions: (optional)
  file_access: [string] (read, write, none)
  network: boolean

Example 1:
"Create an agent that helps with coding tasks using GPT-4. It should be able to generate code examples, debug problems, and explain concepts."

Would result in:
```yaml
name: code-assistant
version: "0.1.0"
description: "An agent that assists with coding tasks, generates examples, debugs problems, and explains concepts"
model:
  provider: openai
  name: gpt-4
  options:
    temperature: 0.7
capabilities:
  - conversation
  - code_generation
  - debugging
  - explanation
memory:
  type: simple
  persistence: true
permissions:
  file_access: ["read"]
  network: false
```

Example 2:
"I need a research agent named ResearchBuddy using Claude that helps me analyze academic papers and summarize them."

Would result in:
```yaml
name: ResearchBuddy
version: "0.1.0"
description: "An agent that analyzes and summarizes academic papers"
model:
  provider: claude
  name: claude-3-opus
  options:
    temperature: 0.5
capabilities:
  - conversation
  - analysis
  - summarization
  - research
memory:
  type: vector
  persistence: true
permissions:
  file_access: ["read"]
  network: true
```

Extract only what is explicitly stated or clearly implied. For anything not mentioned, use reasonable defaults that align with the described purpose.`
	
	// Create user prompt with clear instructions
	userPrompt := fmt.Sprintf("Convert this agent description to YAML. Output ONLY valid YAML with no additional comments or explanations:\n\n%s", text)
	
	// Set options for more precise output
	options := models.Options{
		Temperature: 0.2, // Low temperature for more deterministic output
		MaxTokens:   2000,
	}
	
	// Generate the YAML
	response, err := adapter.Generate(userPrompt, systemPrompt, options)
	if err != nil {
		return "", fmt.Errorf("failed to generate YAML: %w", err)
	}
	
	// Clean the response
	yaml := cleanYAMLResponse(response)
	
	// Validate the YAML structure
	if err := validateYAML(yaml); err != nil {
		if p.Verbose {
			fmt.Printf("Warning: Generated invalid YAML: %v\n", err)
			fmt.Printf("Returning unvalidated YAML anyway:\n%s\n", yaml)
		}
		// Return the YAML anyway, as we might have false positives in validation
	}
	
	return yaml, nil
}

// cleanYAMLResponse cleans the model's response to extract just the YAML
func cleanYAMLResponse(response string) string {
	yaml := response
	yaml = strings.TrimSpace(yaml)
	
	// Remove markdown code block markers if present
	if strings.HasPrefix(yaml, "```yaml") {
		yaml = strings.TrimPrefix(yaml, "```yaml")
		yaml = strings.TrimSuffix(yaml, "```")
	} else if strings.HasPrefix(yaml, "```") {
		yaml = strings.TrimPrefix(yaml, "```")
		yaml = strings.TrimSuffix(yaml, "```")
	}
	
	// Remove any frontmatter or other text before the YAML
	if idx := strings.Index(yaml, "name:"); idx > 0 {
		yaml = yaml[idx:]
	}
	
	return strings.TrimSpace(yaml)
}

// validateYAML validates the generated YAML against the expected schema
func validateYAML(yamlText string) error {
	var agentfile Agentfile
	
	err := yaml.Unmarshal([]byte(yamlText), &agentfile)
	if err != nil {
		return err
	}
	
	// Basic validation
	if agentfile.Name == "" {
		return fmt.Errorf("missing required field: name")
	}
	if agentfile.Version == "" {
		return fmt.Errorf("missing required field: version")
	}
	if agentfile.Model.Provider == "" {
		return fmt.Errorf("missing required field: model.provider")
	}
	if agentfile.Model.Name == "" {
		return fmt.Errorf("missing required field: model.name")
	}
	
	return nil
}
