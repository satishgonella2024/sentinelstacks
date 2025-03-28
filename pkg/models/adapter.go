package models

// ModelAdapter defines the interface for all model adapters
type ModelAdapter interface {
	// Generate sends a prompt to the model and returns the response
	Generate(prompt string, systemPrompt string, options Options) (string, error)
	
	// GetCapabilities returns the capabilities of the model
	GetCapabilities() ModelCapabilities
}

// ModelCapabilities defines the features supported by a model
type ModelCapabilities struct {
	Streaming       bool
	FunctionCalling bool
	MaxTokens       int
	Multimodal      bool
}
