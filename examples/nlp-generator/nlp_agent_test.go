package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator("test-provider", "test-model", "test-dir")

	if g.LLMProvider != "test-provider" {
		t.Errorf("Expected LLMProvider to be 'test-provider', got '%s'", g.LLMProvider)
	}

	if g.LLMModel != "test-model" {
		t.Errorf("Expected LLMModel to be 'test-model', got '%s'", g.LLMModel)
	}

	if g.OutputDir != "test-dir" {
		t.Errorf("Expected OutputDir to be 'test-dir', got '%s'", g.OutputDir)
	}

	if !g.GenerateFiles {
		t.Errorf("Expected GenerateFiles to be true by default")
	}
}

func TestExtractNameAndTag(t *testing.T) {
	testCases := []struct {
		input    string
		wantName string
		wantTag  string
	}{
		{
			input:    "I want an agent called helper-bot that does something",
			wantName: "helper-bot",
			wantTag:  "latest",
		},
		{
			input:    "Create an assistant named travel-agent for planning trips",
			wantName: "travel-agent",
			wantTag:  "latest",
		},
		{
			input:    "Just a regular agent without a specific name",
			wantName: "generated-agent",
			wantTag:  "latest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			g := NewGenerator("test", "test", "test")
			g.extractNameAndTag(tc.input)

			if g.AgentName != tc.wantName {
				t.Errorf("extractNameAndTag() got AgentName = %v, want %v", g.AgentName, tc.wantName)
			}

			if g.TagName != tc.wantTag {
				t.Errorf("extractNameAndTag() got TagName = %v, want %v", g.TagName, tc.wantTag)
			}
		})
	}
}

func TestBuildPrompt(t *testing.T) {
	g := NewGenerator("test", "test", "test")
	input := "Create a helpful assistant"
	prompt := g.buildPrompt(input)

	if !strings.Contains(prompt, input) {
		t.Errorf("Expected prompt to contain the input text")
	}

	if !strings.Contains(prompt, "Convert the following natural language description") {
		t.Errorf("Expected prompt to contain instructions for converting to YAML")
	}

	if !strings.Contains(prompt, "Sentinelfile YAML") {
		t.Errorf("Expected prompt to mention Sentinelfile YAML")
	}
}

func TestProcessNaturalLanguage(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "nlp-generator-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	g := NewGenerator("test", "test", tmpDir)

	input := "Create a chatbot assistant that can have casual conversations with users"
	response, err := g.ProcessNaturalLanguage(input)

	if err != nil {
		t.Fatalf("ProcessNaturalLanguage() error = %v", err)
	}

	if response == nil {
		t.Fatalf("ProcessNaturalLanguage() returned nil response")
	}

	if response.Sentinelfile == "" {
		t.Errorf("Expected non-empty Sentinelfile in response")
	}

	if response.Metadata == nil {
		t.Errorf("Expected non-nil Metadata in response")
	}

	// Check if files were created
	agentDir := filepath.Join(tmpDir, g.AgentName)
	sentinelfilePath := filepath.Join(agentDir, "Sentinelfile")
	metadataPath := filepath.Join(agentDir, "metadata.json")

	if _, err := os.Stat(sentinelfilePath); os.IsNotExist(err) {
		t.Errorf("Expected Sentinelfile to be created at %s", sentinelfilePath)
	}

	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		t.Errorf("Expected metadata.json to be created at %s", metadataPath)
	}
}
