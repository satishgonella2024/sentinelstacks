package parser

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/satishgonella2024/sentinelstacks/internal/shim"
	"github.com/satishgonella2024/sentinelstacks/pkg/agent"
)

// SentinelfileParser handles parsing Sentinelfiles
type SentinelfileParser struct {
	llmProvider string
	llmEndpoint string
	llmAPIKey   string
	llmModel    string
}

// NewSentinelfileParser creates a new SentinelfileParser
func NewSentinelfileParser(llmProvider string) *SentinelfileParser {
	if llmProvider == "" {
		llmProvider = "claude" // Default provider
	}

	// Get API key from environment
	apiKey := os.Getenv("SENTINEL_API_KEY")

	return &SentinelfileParser{
		llmProvider: llmProvider,
		llmEndpoint: os.Getenv("SENTINEL_LLM_ENDPOINT"),
		llmAPIKey:   apiKey,
		llmModel:    os.Getenv("SENTINEL_LLM_MODEL"),
	}
}

// ParseFile reads a Sentinelfile and returns an agent definition
func (p *SentinelfileParser) ParseFile(filePath string) (*agent.Definition, error) {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Sentinelfile: %w", err)
	}

	return p.Parse(string(content))
}

// Parse parses a Sentinelfile content and returns an agent definition
func (p *SentinelfileParser) Parse(content string) (*agent.Definition, error) {
	// For sophisticated parsing, use the LLM
	if p.llmAPIKey != "" && !strings.Contains(content, "# DEBUG_SIMPLE_PARSE") {
		return p.parseLLM(content)
	}

	// Fall back to simple parsing for development or when API key is not available
	return p.parseSimple(content)
}

// parseLLM uses an LLM to parse the Sentinelfile content
func (p *SentinelfileParser) parseLLM(content string) (*agent.Definition, error) {
	// Create an LLM shim
	llm, err := shim.ShimFactory(p.llmProvider, p.llmEndpoint, p.llmAPIKey, p.llmModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM shim: %w", err)
	}

	// Parse the Sentinelfile with the LLM
	result, err := llm.ParseSentinelfile(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Sentinelfile with LLM: %w", err)
	}

	// Convert the result to an agent definition
	def := &agent.Definition{}

	// Extract the basic fields
	if name, ok := result["name"].(string); ok {
		def.Name = name
	} else {
		def.Name = "unnamed-agent"
	}

	if description, ok := result["description"].(string); ok {
		def.Description = description
	} else {
		def.Description = "No description provided"
	}

	if baseModel, ok := result["baseModel"].(string); ok {
		def.BaseModel = baseModel
	} else {
		def.BaseModel = "claude-3.7-sonnet"
	}

	// Extract the capabilities
	if capabilities, ok := result["capabilities"].([]interface{}); ok {
		for _, cap := range capabilities {
			if capStr, ok := cap.(string); ok {
				def.Capabilities = append(def.Capabilities, capStr)
			}
		}
	}

	// Extract the tools
	if tools, ok := result["tools"].([]interface{}); ok {
		for _, tool := range tools {
			if toolStr, ok := tool.(string); ok {
				def.Tools = append(def.Tools, toolStr)
			}
		}
	}

	// Extract the parameters
	if parameters, ok := result["parameters"].(map[string]interface{}); ok {
		def.Parameters = parameters
	} else {
		def.Parameters = make(map[string]interface{})
	}

	// Extract the lifecycle
	if lifecycle, ok := result["lifecycle"].(map[string]interface{}); ok {
		if init, ok := lifecycle["initialization"].(string); ok {
			if def.Lifecycle.Initialization == "" {
				def.Lifecycle = agent.Lifecycle{
					Initialization: init,
				}
			} else {
				def.Lifecycle.Initialization = init
			}
		}

		if term, ok := lifecycle["termination"].(string); ok {
			def.Lifecycle.Termination = term
		}
	}

	// Extract the state schema
	if stateSchema, ok := result["stateSchema"].(map[string]interface{}); ok {
		def.StateSchema = make(map[string]agent.StateField)

		for k, v := range stateSchema {
			if field, ok := v.(map[string]interface{}); ok {
				stateField := agent.StateField{}

				if typ, ok := field["type"].(string); ok {
					stateField.Type = typ
				} else {
					stateField.Type = "string"
				}

				if desc, ok := field["description"].(string); ok {
					stateField.Description = desc
				}

				if def, ok := field["default"]; ok {
					stateField.Default = def
				}

				def.StateSchema[k] = stateField
			} else if typeStr, ok := v.(string); ok {
				// Simple case where just the type is provided
				def.StateSchema[k] = agent.StateField{
					Type: typeStr,
				}
			}
		}
	}

	return def, nil
}

// parseSimple does a simple parsing of the Sentinelfile content without using an LLM
func (p *SentinelfileParser) parseSimple(content string) (*agent.Definition, error) {
	// Extract name from the first line if it starts with a comment
	name := "unnamed-agent"
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# Sentinelfile for ") {
			name = strings.TrimPrefix(line, "# Sentinelfile for ")
			name = strings.ToLower(name)
			name = strings.ReplaceAll(name, " ", "-")
			break
		}
	}

	// Extract capabilities
	capabilities := []string{}
	inCapabilities := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "should be able to:") {
			inCapabilities = true
			continue
		} else if inCapabilities && strings.HasPrefix(line, "-") {
			capability := strings.TrimPrefix(line, "-")
			capability = strings.TrimSpace(capability)
			if capability != "" && capability != "[Capability 1]" && capability != "[Capability 2]" && capability != "[Capability 3]" {
				capabilities = append(capabilities, capability)
			}
		} else if inCapabilities && line == "" {
			// Empty line ends the capabilities section
			inCapabilities = false
		}
	}

	// Extract base model
	baseModel := "claude-3.7-sonnet" // Default
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "should use") && strings.Contains(line, "as its base model") {
			parts := strings.Split(line, "should use")
			if len(parts) > 1 {
				modelPart := parts[1]
				modelPart = strings.Split(modelPart, "as its base model")[0]
				baseModel = strings.TrimSpace(modelPart)
			}
			break
		}
	}

	// Extract tools
	tools := []string{}
	inTools := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "access the following tools:") {
			inTools = true
			continue
		} else if inTools && strings.HasPrefix(line, "-") {
			tool := strings.TrimPrefix(line, "-")
			tool = strings.TrimSpace(tool)
			if tool != "" && tool != "[Tool 1]" && tool != "[Tool 2]" {
				tools = append(tools, tool)
			}
		} else if inTools && line == "" {
			// Empty line ends the tools section
			inTools = false
		}
	}

	// Extract parameters
	parameters := map[string]interface{}{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Set ") && strings.Contains(line, " to ") {
			line = strings.TrimPrefix(line, "Set ")
			parts := strings.SplitN(line, " to ", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove trailing period if present
				if strings.HasSuffix(value, ".") {
					value = value[:len(value)-1]
				}

				// Try to convert value to appropriate type
				if value == "true" {
					parameters[key] = true
				} else if value == "false" {
					parameters[key] = false
				} else if intVal, err := parseInt(value); err == nil {
					parameters[key] = intVal
				} else if floatVal, err := parseFloat(value); err == nil {
					parameters[key] = floatVal
				} else {
					parameters[key] = value
				}
			}
		}
	}

	// Create the agent definition
	def := &agent.Definition{
		Name:         name,
		BaseModel:    baseModel,
		Capabilities: capabilities,
		Tools:        tools,
		Parameters:   parameters,
	}

	// Extract description (first non-empty, non-comment line)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			def.Description = line
			break
		}
	}

	// Ensure we have at least basic information
	if def.Description == "" {
		def.Description = "No description provided"
	}

	return def, nil
}

// parseInt tries to parse a string as an integer
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// parseFloat tries to parse a string as a float
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// ValidateDefinition checks if the agent definition is valid
func ValidateDefinition(def *agent.Definition) error {
	if def.Name == "" {
		return errors.New("agent name is required")
	}
	if def.BaseModel == "" {
		return errors.New("base model is required")
	}
	return nil
}
