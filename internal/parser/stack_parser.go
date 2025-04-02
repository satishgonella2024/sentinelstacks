package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"regexp"
	
	"gopkg.in/yaml.v3"
	
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/stack"
)

// StackParser is responsible for parsing natural language descriptions into structured stack specifications
type StackParser struct {
	// Configuration options could be added here
}

// NewStackParser creates a new stack parser
func NewStackParser() *StackParser {
	return &StackParser{}
}

// ParseFromYAML parses a YAML string into a StackSpec
func (p *StackParser) ParseFromYAML(yamlContent string) (stack.StackSpec, error) {
	var spec stack.StackSpec
	err := yaml.Unmarshal([]byte(yamlContent), &spec)
	if err != nil {
		return stack.StackSpec{}, fmt.Errorf("failed to parse YAML: %w", err)
	}
	
	// Validate the parsed stack
	if err := validateStackSpec(spec); err != nil {
		return stack.StackSpec{}, err
	}
	
	return spec, nil
}

// ParseFromJSON parses a JSON string into a StackSpec
func (p *StackParser) ParseFromJSON(jsonContent string) (stack.StackSpec, error) {
	var spec stack.StackSpec
	err := json.Unmarshal([]byte(jsonContent), &spec)
	if err != nil {
		return stack.StackSpec{}, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// Validate the parsed stack
	if err := validateStackSpec(spec); err != nil {
		return stack.StackSpec{}, err
	}
	
	return spec, nil
}

// ParseFromNaturalLanguage parses a natural language description into a StackSpec
func (p *StackParser) ParseFromNaturalLanguage(description string) (stack.StackSpec, error) {
	// Extract stack name
	stackName := extractStackName(description)
	
	// Extract agents and their relationships
	agents, err := extractAgents(description)
	if err != nil {
		return stack.StackSpec{}, err
	}
	
	// Build the stack specification
	spec := stack.StackSpec{
		Name:        stackName,
		Description: description,
		Version:     "1.0.0",
		Agents:      agents,
	}
	
	// Validate the parsed stack
	if err := validateStackSpec(spec); err != nil {
		return stack.StackSpec{}, err
	}
	
	return spec, nil
}

// extractStackName attempts to extract a stack name from the description
func extractStackName(description string) string {
	// Try to find phrases like "Create a stack called X" or "Build a X stack"
	nameRegexes := []string{
		`(?i)stack (?:called|named) ["']?([a-zA-Z0-9_-]+)["']?`,
		`(?i)["']?([a-zA-Z0-9_-]+)["']? stack`,
		`(?i)create (?:a|an) ["']?([a-zA-Z0-9_-]+)["']?`,
	}
	
	for _, regex := range nameRegexes {
		re := regexp.MustCompile(regex)
		matches := re.FindStringSubmatch(description)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	
	// Default name if none found
	return "generated-stack"
}

// extractAgents parses the description to identify agents and their dependencies
func extractAgents(description string) ([]stack.StackAgentSpec, error) {
	// Split the description into sentences for analysis
	sentences := splitIntoSentences(description)
	
	// Map to store agents by ID
	agentMap := make(map[string]stack.StackAgentSpec)
	
	// First pass: identify agents
	for _, sentence := range sentences {
		// Look for agent definitions
		if strings.Contains(strings.ToLower(sentence), "agent") || 
		   strings.Contains(strings.ToLower(sentence), "step") {
			
			// Extract agent ID and purpose
			agentID, agentType := extractAgentInfo(sentence)
			if agentID != "" {
				// Create agent spec
				agentMap[agentID] = stack.StackAgentSpec{
					ID:       agentID,
					Uses:     agentType,
					InputFrom: []string{},
					Params:   map[string]interface{}{},
				}
			}
		}
	}
	
	// If no agents found, try a different approach - look for steps or numbered items
	if len(agentMap) == 0 {
		agentMap = extractAgentsFromSteps(description)
	}
	
	// Second pass: identify relationships
	for _, sentence := range sentences {
		// Look for relationship descriptions
		if strings.Contains(strings.ToLower(sentence), "output") || 
		   strings.Contains(strings.ToLower(sentence), "input") ||
		   strings.Contains(strings.ToLower(sentence), "then") ||
		   strings.Contains(strings.ToLower(sentence), "after") {
			
			// Extract relationship
			sourceID, targetID := extractRelationship(sentence)
			if sourceID != "" && targetID != "" {
				// Check if agents exist
				if spec, exists := agentMap[targetID]; exists {
					// Add dependency
					spec.InputFrom = append(spec.InputFrom, sourceID)
					agentMap[targetID] = spec
				}
			}
		}
	}
	
	// Convert map to slice
	agents := make([]stack.StackAgentSpec, 0, len(agentMap))
	for _, spec := range agentMap {
		agents = append(agents, spec)
	}
	
	// If we found no agents, return an error
	if len(agents) == 0 {
		return nil, errors.New("no agents identified in the description")
	}
	
	return agents, nil
}

// extractAgentsFromSteps attempts to identify agents in a step-by-step description
func extractAgentsFromSteps(description string) map[string]stack.StackAgentSpec {
	agentMap := make(map[string]stack.StackAgentSpec)
	
	// Look for numbered steps or bulleted lists
	stepRegex := regexp.MustCompile(`(?i)(?:\d+\.\s*|\*\s*|\-\s*)([^.!?]+)`)
	steps := stepRegex.FindAllStringSubmatch(description, -1)
	
	for i, step := range steps {
		if len(step) > 1 {
			stepText := step[1]
			agentID := fmt.Sprintf("agent-%d", i+1)
			
			// Try to extract a better ID from the step text
			customID := extractCustomAgentID(stepText)
			if customID != "" {
				agentID = customID
			}
			
			// Extract the agent type or purpose
			agentType := determineAgentType(stepText)
			
			// Create agent spec
			agentMap[agentID] = stack.StackAgentSpec{
				ID:       agentID,
				Uses:     agentType,
				InputFrom: []string{},
				Params:   map[string]interface{}{
					"description": stepText,
				},
			}
			
			// Add dependency on previous step if applicable
			if i > 0 {
				prevAgentID := fmt.Sprintf("agent-%d", i)
				if customPrevID := extractCustomAgentID(steps[i-1][1]); customPrevID != "" {
					prevAgentID = customPrevID
				}
				
				spec := agentMap[agentID]
				spec.InputFrom = append(spec.InputFrom, prevAgentID)
				agentMap[agentID] = spec
			}
		}
	}
	
	return agentMap
}

// extractCustomAgentID attempts to find a custom ID for an agent from its description
func extractCustomAgentID(text string) string {
	// Look for phrases like "a data processor agent" or "the summarizer step"
	idRegex := regexp.MustCompile(`(?i)(?:a|an|the)\s+([a-zA-Z0-9_-]+)(?:\s+agent|\s+step)`)
	matches := idRegex.FindStringSubmatch(text)
	
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	
	return ""
}

// determineAgentType tries to determine the type of agent based on its description
func determineAgentType(text string) string {
	text = strings.ToLower(text)
	
	// Map common tasks to agent types
	typeMap := map[string]string{
		"summar":      "summarizer",
		"extract":     "extractor",
		"analyz":      "analyzer",
		"process":     "processor",
		"transform":   "transformer",
		"classify":    "classifier",
		"categoriz":   "categorizer",
		"translat":    "translator",
		"generat":     "generator",
		"retriev":     "retriever",
		"search":      "searcher",
		"collect":     "collector",
		"filter":      "filter",
		"sort":        "sorter",
		"rank":        "ranker",
		"recommend":   "recommender",
		"predict":     "predictor",
		"detect":      "detector",
		"identify":    "identifier",
		"validate":    "validator",
		"verify":      "verifier",
		"authenticate":"authenticator",
		"authorize":   "authorizer",
	}
	
	for keyword, agentType := range typeMap {
		if strings.Contains(text, keyword) {
			return agentType
		}
	}
	
	// Default to generic agent
	return "generic-agent"
}

// extractAgentInfo tries to extract agent ID and type from a sentence
func extractAgentInfo(sentence string) (string, string) {
	// Look for explicit agent definitions
	agentRegex := regexp.MustCompile(`(?i)(?:a|an|the)\s+([a-zA-Z0-9_-]+)(?:\s+agent|\s+step)(?:[^a-zA-Z0-9_-]|$)`)
	matches := agentRegex.FindStringSubmatch(sentence)
	
	if len(matches) > 1 {
		agentID := strings.ToLower(matches[1])
		agentType := determineAgentType(sentence)
		return agentID, agentType
	}
	
	return "", ""
}

// extractRelationship tries to identify relationships between agents
func extractRelationship(sentence string) (string, string) {
	// Look for phrases like "Agent A outputs to Agent B" or "Agent B uses input from Agent A"
	relRegex := regexp.MustCompile(`(?i)(?:the\s+)?([a-zA-Z0-9_-]+)(?:\s+agent|\s+step)?\s+(?:outputs?|sends|provides|feeds)(?:[^a-zA-Z0-9_-]+)(?:to|into|for)(?:[^a-zA-Z0-9_-]+)(?:the\s+)?([a-zA-Z0-9_-]+)(?:\s+agent|\s+step)?`)
	matches := relRegex.FindStringSubmatch(sentence)
	
	if len(matches) > 2 {
		sourceID := strings.ToLower(matches[1])
		targetID := strings.ToLower(matches[2])
		return sourceID, targetID
	}
	
	// Try reverse pattern: "Agent B takes input from Agent A"
	revRegex := regexp.MustCompile(`(?i)(?:the\s+)?([a-zA-Z0-9_-]+)(?:\s+agent|\s+step)?\s+(?:takes|uses|receives|needs)(?:[^a-zA-Z0-9_-]+)(?:input|data|results)(?:[^a-zA-Z0-9_-]+)(?:from|of)(?:[^a-zA-Z0-9_-]+)(?:the\s+)?([a-zA-Z0-9_-]+)(?:\s+agent|\s+step)?`)
	matches = revRegex.FindStringSubmatch(sentence)
	
	if len(matches) > 2 {
		targetID := strings.ToLower(matches[1])
		sourceID := strings.ToLower(matches[2])
		return sourceID, targetID
	}
	
	// Try sequence pattern: "After Agent A, Agent B runs"
	seqRegex := regexp.MustCompile(`(?i)(?:after|following|once)(?:[^a-zA-Z0-9_-]+)(?:the\s+)?([a-zA-Z0-9_-]+)(?:\s+agent|\s+step)?[^a-zA-Z0-9_-]+(?:the\s+)?([a-zA-Z0-9_-]+)(?:\s+agent|\s+step)?`)
	matches = seqRegex.FindStringSubmatch(sentence)
	
	if len(matches) > 2 {
		sourceID := strings.ToLower(matches[1])
		targetID := strings.ToLower(matches[2])
		return sourceID, targetID
	}
	
	return "", ""
}

// splitIntoSentences breaks a text into sentences
func splitIntoSentences(text string) []string {
	// Basic sentence splitting - could be enhanced with NLP libraries
	re := regexp.MustCompile(`[.!?]+`)
	sentences := re.Split(text, -1)
	
	var result []string
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	
	return result
}

// validateStackSpec checks if a stack specification is valid
func validateStackSpec(spec stack.StackSpec) error {
	if spec.Name == "" {
		return errors.New("stack name cannot be empty")
	}
	
	if len(spec.Agents) == 0 {
		return errors.New("stack must contain at least one agent")
	}
	
	// Check for duplicate agent IDs
	agentIDs := make(map[string]bool)
	for _, agent := range spec.Agents {
		if agent.ID == "" {
			return errors.New("agent ID cannot be empty")
		}
		
		if agentIDs[agent.ID] {
			return fmt.Errorf("duplicate agent ID detected: %s", agent.ID)
		}
		
		agentIDs[agent.ID] = true
	}
	
	// Check for references to non-existent agents
	for _, agent := range spec.Agents {
		for _, inputFrom := range agent.InputFrom {
			if !agentIDs[inputFrom] {
				return fmt.Errorf("agent %s references non-existent agent %s", agent.ID, inputFrom)
			}
		}
	}
	
	return nil
}

// GenerateYAML converts a StackSpec to YAML
func (p *StackParser) GenerateYAML(spec stack.StackSpec) (string, error) {
	yamlBytes, err := yaml.Marshal(spec)
	if err != nil {
		return "", fmt.Errorf("failed to generate YAML: %w", err)
	}
	
	return string(yamlBytes), nil
}

// GenerateJSON converts a StackSpec to JSON
func (p *StackParser) GenerateJSON(spec stack.StackSpec) (string, error) {
	jsonBytes, err := json.Marshal(spec)
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON: %w", err)
	}
	
	return string(jsonBytes), nil
}
