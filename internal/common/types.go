package common

// CompletionRequest represents a request to an LLM provider
type CompletionRequest struct {
	Model        string
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
	Temperature  float64
}

// CompletionResponse represents a response from an LLM provider
type CompletionResponse struct {
	Text       string
	FinishReason string
	Usage       struct {
		PromptTokens     int
		CompletionTokens int
		TotalTokens      int
	}
}

// AgentDefinition represents the core definition of an agent
type AgentDefinition struct {
	Name          string
	Description   string
	Version       string
	BaseModel     string
	SystemPrompt  string
	PromptTemplate string
	MaxTokens     int
	Temperature   float64
	Tools         []string
	OutputFormat  string
}
