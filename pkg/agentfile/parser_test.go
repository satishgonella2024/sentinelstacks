package agentfile

import (
	"strings"
	"testing"
)

func TestParseText(t *testing.T) {
	// Create a test parser
	parser := NewParser("http://test-endpoint")
	
	// Skip tests that require a real LLM endpoint
	t.Skip("Skipping tests that require LLM endpoint")
	
	// Test cases
	testCases := []struct {
		name        string
		description string
		expectKeys  []string
	}{
		{
			name:        "Simple agent",
			description: "Create a coding assistant using GPT-4",
			expectKeys:  []string{"name", "version", "description", "model", "capabilities"},
		},
		{
			name:        "Research agent",
			description: "Create a research agent named ResearchBuddy using Claude that helps analyze papers",
			expectKeys:  []string{"name", "version", "description", "model", "capabilities", "memory"},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			yaml, err := parser.ParseText(tc.description)
			if err != nil {
				t.Fatalf("Error parsing text: %v", err)
			}
			
			// Verify YAML contains expected keys
			for _, key := range tc.expectKeys {
				if !strings.Contains(yaml, key+":") {
					t.Errorf("Expected key '%s' not found in YAML output", key)
				}
			}
		})
	}
}

func TestCleanYAMLResponse(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Markdown code block",
			input:    "```yaml\nname: test\nversion: 0.1.0\n```",
			expected: "name: test\nversion: 0.1.0",
		},
		{
			name:     "Plain YAML",
			input:    "name: test\nversion: 0.1.0",
			expected: "name: test\nversion: 0.1.0",
		},
		{
			name:     "YAML with extra text",
			input:    "Here's the YAML:\n\nname: test\nversion: 0.1.0",
			expected: "name: test\nversion: 0.1.0",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cleanYAMLResponse(tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestValidateYAML(t *testing.T) {
	testCases := []struct {
		name      string
		yaml      string
		expectErr bool
	}{
		{
			name: "Valid YAML",
			yaml: `name: test
version: 0.1.0
description: Test agent
model:
  provider: ollama
  name: llama3
capabilities:
  - conversation
memory:
  type: simple
  persistence: true`,
			expectErr: false,
		},
		{
			name: "Missing name",
			yaml: `version: 0.1.0
description: Test agent
model:
  provider: ollama
  name: llama3`,
			expectErr: true,
		},
		{
			name: "Missing model provider",
			yaml: `name: test
version: 0.1.0
description: Test agent
model:
  name: llama3`,
			expectErr: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateYAML(tc.yaml)
			if tc.expectErr && err == nil {
				t.Error("Expected error but got nil")
			} else if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
