package shim

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Provider constants
const (
	ProviderClaude = "claude"
	ProviderOpenAI = "openai"
	ProviderOllama = "ollama"
	ProviderGoogle = "google"
	ProviderMock   = "mock"
)

// DefaultModels maps providers to their default models
var DefaultModels = map[string]string{
	ProviderClaude: "claude-3-5-sonnet-20240627",
	ProviderOpenAI: "gpt-4-turbo",
	ProviderOllama: "llama3",
	ProviderGoogle: "gemini-1.5-pro",
	ProviderMock:   "mock-model",
}

// DefaultEndpoints maps providers to their default API endpoints
var DefaultEndpoints = map[string]string{
	ProviderClaude: "https://api.anthropic.com/v1/messages",
	ProviderOpenAI: "https://api.openai.com/v1/chat/completions",
	ProviderOllama: "http://localhost:11434/api/generate",
	ProviderGoogle: "https://generativelanguage.googleapis.com/v1beta/models",
	ProviderMock:   "",
}

// GetProviderFromEnv gets the provider from environment variables or returns a default
func GetProviderFromEnv() string {
	provider := os.Getenv("SENTINEL_LLM_PROVIDER")
	if provider == "" {
		provider = ProviderClaude // Default provider
	}
	
	return provider
}

// GetModelFromEnv gets the model from environment variables or returns a default based on provider
func GetModelFromEnv(provider string) string {
	model := os.Getenv("SENTINEL_LLM_MODEL")
	if model == "" {
		model = DefaultModels[provider]
	}
	
	return model
}

// GetEndpointFromEnv gets the endpoint from environment variables or returns a default based on provider
func GetEndpointFromEnv(provider string) string {
	endpoint := os.Getenv("SENTINEL_LLM_ENDPOINT")
	if endpoint == "" {
		endpoint = DefaultEndpoints[provider]
	}
	
	return endpoint
}

// GetAPIKeyFromEnv gets the API key from environment variables
func GetAPIKeyFromEnv() string {
	// Check for provider-specific API key first
	provider := GetProviderFromEnv()
	switch provider {
	case ProviderGoogle:
		if key := os.Getenv("GOOGLE_API_KEY"); key != "" {
			return key
		}
	case ProviderOpenAI:
		if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			return key
		}
	case ProviderClaude:
		if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
			return key
		}
	}
	
	// Fall back to generic API key
	return os.Getenv("SENTINEL_API_KEY")
}

// CreateShimFromEnv creates a shim based on environment variables
func CreateShimFromEnv() (LLMShim, error) {
	provider := GetProviderFromEnv()
	model := GetModelFromEnv(provider)
	endpoint := GetEndpointFromEnv(provider)
	apiKey := GetAPIKeyFromEnv()
	
	// Create config
	config := Config{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
		Endpoint: endpoint,
	}
	
	// Create shim
	return ShimFactory(provider, endpoint, apiKey, model)
}

// IsMultimodalModel checks if a model supports multimodal inputs
func IsMultimodalModel(provider, model string) bool {
	switch strings.ToLower(provider) {
	case ProviderClaude:
		// Claude multimodal models
		claudeMultimodal := map[string]bool{
			"claude-3-opus-20240229":   true,
			"claude-3-sonnet-20240229": true,
			"claude-3-haiku-20240307":  true,
			"claude-3-5-sonnet-20240627": true,
			"claude-3-5-sonnet": true,
			"claude-3-opus": true,
			"claude-3-sonnet": true,
			"claude-3-haiku": true,
		}
		return claudeMultimodal[model]
		
	case ProviderOpenAI:
		// OpenAI multimodal models
		openaiMultimodal := map[string]bool{
			"gpt-4-vision-preview": true,
			"gpt-4-turbo-preview":  true,
			"gpt-4-turbo":          true,
			"gpt-4-1106-vision-preview": true,
			"gpt-4-1106-preview":   true,
			"gpt-4-vision":         true,
		}
		return openaiMultimodal[model]
		
	case ProviderOllama:
		// Ollama multimodal models
		ollamaMultimodal := map[string]bool{
			"llava": true,
			"llava:7b": true,
			"llava:13b": true,
			"llava:latest": true,
			"bakllava": true,
			"bakllava:latest": true,
			"moondream": true,
			"moondream:latest": true,
			"fuyu": true,
			"fuyu:latest": true,
			"yi-vl": true,
			"yi-vl:latest": true,
			"cogvlm": true,
			"cogvlm:latest": true,
		}
		
		// Check for exact match or prefix match
		if ollamaMultimodal[model] {
			return true
		}
		
		// Check for prefix match
		for modelName := range ollamaMultimodal {
			if strings.HasPrefix(model, modelName+":") {
				return true
			}
		}
		
		return false
		
	case ProviderGoogle:
		// Google multimodal models
		googleMultimodal := map[string]bool{
			"gemini-pro-vision":  true,
			"gemini-1.5-pro":     true,
			"gemini-1.5-flash":   true,
		}
		return googleMultimodal[model]
		
	case ProviderMock:
		// Mock always supports multimodal for testing
		return true
		
	default:
		return false
	}
}

// ValidateProviderConfig validates a provider configuration
func ValidateProviderConfig(config Config) error {
	// Check provider
	switch strings.ToLower(config.Provider) {
	case ProviderClaude, ProviderOpenAI, ProviderOllama, ProviderGoogle, ProviderMock:
		// Valid provider
	default:
		return fmt.Errorf("unsupported provider: %s", config.Provider)
	}
	
	// Check API key for providers that require it
	if (config.Provider == ProviderClaude || 
		config.Provider == ProviderOpenAI || 
		config.Provider == ProviderGoogle) && config.APIKey == "" {
		return fmt.Errorf("API key is required for %s provider", config.Provider)
	}
	
	// Check endpoint
	if config.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	
	return nil
}

// GetDefaultSystemPrompt returns the default system prompt for a provider
func GetDefaultSystemPrompt(provider string, purpose string) string {
	switch purpose {
	case "agent":
		return "You are a helpful AI assistant working as part of an agent system."
	case "parser":
		switch provider {
		case ProviderClaude:
			return `You are a specialized AI that parses natural language Sentinelfiles into structured JSON. 
A Sentinelfile defines an AI agent's capabilities, behavior, and requirements. 
Your task is to extract information from the Sentinelfile and output a valid JSON structure with the following fields:
- name: The name of the agent (string, required)
- description: A short description of the agent (string, required)
- baseModel: The LLM model to use (string, required)
- capabilities: List of capabilities the agent should have (array of strings, optional)
- tools: List of tools the agent should have access to (array of strings, optional)
- stateSchema: Description of the state the agent should maintain (object, optional)
- parameters: Configuration parameters for the agent (object, optional)
- lifecycle: Object containing initialization and termination behaviors (object, optional)

Output ONLY the JSON object, no additional text or explanation.`
		case ProviderOpenAI:
			return `You are a specialized AI that parses natural language Sentinelfiles into structured JSON. 
Parse the provided Sentinelfile text into a JSON object with these fields:
- name: The agent name (string, required)
- description: Agent description (string, required)
- baseModel: LLM model to use (string, required)
- capabilities: List of agent capabilities (array of strings, optional)
- tools: List of required tools (array of strings, optional)
- stateSchema: State schema definition (object, optional)
- parameters: Configuration parameters (object, optional)
- lifecycle: Initialization and termination behaviors (object, optional)

Output ONLY valid JSON with no additional text.`
		case ProviderOllama:
			return `You are an AI assistant that specializes in parsing Sentinelfiles. Extract the name, description, baseModel, capabilities, tools, and other configuration from the Sentinelfile and return them as a valid JSON object.`
		case ProviderGoogle:
			return `You are a specialized AI that parses natural language Sentinelfiles into structured JSON. 
A Sentinelfile defines an AI agent's capabilities, behavior, and requirements. 
Extract information from the Sentinelfile into a JSON object with these fields:
- name: The agent name (string, required)
- description: Agent description (string, required)
- baseModel: LLM model to use (string, required)
- capabilities: List of agent capabilities (array of strings, optional)
- tools: List of required tools (array of strings, optional)
- stateSchema: State schema definition (object, optional)
- parameters: Configuration parameters (object, optional)
- lifecycle: Initialization and termination behaviors (object, optional)

Output ONLY valid JSON with no additional text.`
		default:
			return "You are a helpful AI assistant."
		}
	default:
		return "You are a helpful AI assistant."
	}
}

// CreateShimWithDefaultConfig creates a shim with default configuration
func CreateShimWithDefaultConfig(provider string) (LLMShim, error) {
	model := DefaultModels[provider]
	endpoint := DefaultEndpoints[provider]
	apiKey := GetAPIKeyFromEnv()
	
	// Set default timeout
	timeout := 60 * time.Second
	
	// Create config
	config := Config{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
		Endpoint: endpoint,
		Timeout:  timeout,
	}
	
	// Create shim
	shim, err := ShimFactory(provider, endpoint, apiKey, model)
	if err != nil {
		return nil, err
	}
	
	// Set default system prompt
	shim.SetSystemPrompt(GetDefaultSystemPrompt(provider, "agent"))
	
	return shim, nil
}
