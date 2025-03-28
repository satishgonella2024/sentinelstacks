package agentfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Parser converts natural language descriptions to structured YAML Agentfiles
type Parser struct {
	ModelEndpoint string
}

// NewParser creates a new Agentfile parser
func NewParser(modelEndpoint string) *Parser {
	return &Parser{
		ModelEndpoint: modelEndpoint,
	}
}

// ParseFile takes a natural language file and converts it to YAML
func (p *Parser) ParseFile(naturalLanguagePath string) (string, error) {
	// Read the natural language file
	content, err := ioutil.ReadFile(naturalLanguagePath)
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

	err = ioutil.WriteFile(yamlPath, []byte(yaml), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write YAML file: %w", err)
	}

	return yamlPath, nil
}

// ParseText converts natural language text to YAML
func (p *Parser) ParseText(text string) (string, error) {
	// TODO: Implement actual LLM-based parsing
	// For now, we'll just return a hardcoded example

	yaml := `name: my-agent
version: "1.0.0"
description: "An agent created from natural language"
model:
  provider: "ollama"
  name: "llama3"
capabilities:
  - text_processing
  - conversation
memory:
  type: "simple"
  persistence: true
`
	return yaml, nil
}
