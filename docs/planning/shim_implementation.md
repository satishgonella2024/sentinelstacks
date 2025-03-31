# Sentinel Shim Implementation Plan

This document details the implementation plan for the Sentinel Shim component, which provides abstraction between different LLM providers and the SentinelStacks runtime.

## Overview

The Sentinel Shim is a critical component that enables agent portability across different LLM providers. It abstracts away provider-specific details, standardizes interactions, and optimizes prompts for each provider.

## Architecture

### Component Structure

```
sentinel/
├── pkg/
│   ├── shim/                  # Shim package
│   │   ├── shim.go            # Core interface definitions
│   │   ├── provider.go        # Provider interface
│   │   ├── context.go         # Context management
│   │   ├── prompt.go          # Prompt engineering
│   │   ├── cache.go           # Response caching
│   │   ├── utils.go           # Utility functions
│   │   ├── config.go          # Configuration management
│   │   ├── providers/         # Provider implementations
│   │   │   ├── claude.go      # Claude provider
│   │   │   ├── openai.go      # OpenAI provider
│   │   │   ├── llama.go       # Llama provider
│   │   │   └── custom.go      # Custom provider
│   │   ├── formatters/        # Format converters
│   │   │   ├── json.go        # JSON formatter
│   │   │   ├── text.go        # Plain text formatter
│   │   │   └── markdown.go    # Markdown formatter
│   │   └── tools/             # Tool integrations
│   │       ├── registry.go    # Tool registry
│   │       ├── web.go         # Web tools
│   │       └── data.go        # Data tools
```

### Core Interfaces

```go
// Shim is the main interface for interacting with LLM providers
type Shim interface {
    // Initialize the shim with configuration
    Initialize(config Config) error
    
    // Generate a response from the LLM
    Generate(ctx context.Context, input GenerateInput) (*GenerateOutput, error)
    
    // Stream a response from the LLM
    Stream(ctx context.Context, input GenerateInput) (<-chan StreamChunk, error)
    
    // Get embeddings for text
    GetEmbeddings(ctx context.Context, input EmbeddingsInput) (*EmbeddingsOutput, error)
    
    // Close and clean up resources
    Close() error
}

// Provider is the interface implemented by specific LLM providers
type Provider interface {
    // Name of the provider
    Name() string
    
    // Available models from this provider
    AvailableModels() []string
    
    // Generate a response (provider-specific)
    GenerateResponse(ctx context.Context, prompt string, params ProviderParams) (string, error)
    
    // Stream a response (provider-specific)
    StreamResponse(ctx context.Context, prompt string, params ProviderParams) (<-chan string, error)
    
    // Get embeddings (provider-specific)
    GetEmbeddings(ctx context.Context, text []string) ([][]float32, error)
}
```

## Implementation Phases

### Phase 1: Core Framework (Weeks 1-2)

1. Define core interfaces and data structures
2. Implement basic configuration management
3. Create provider registration mechanism
4. Build prompt handling utilities
5. Set up error handling and logging

```go
// Example core shim implementation
type CoreShim struct {
    config      Config
    providers   map[string]Provider
    activeModel string
    activeProvider Provider
    cache       Cache
}

func NewShim() Shim {
    return &CoreShim{
        providers: make(map[string]Provider),
    }
}

func (s *CoreShim) Initialize(config Config) error {
    s.config = config
    
    // Register built-in providers
    if err := s.registerBuiltinProviders(); err != nil {
        return err
    }
    
    // Set active provider based on config
    provider, ok := s.providers[s.config.Provider]
    if !ok {
        return fmt.Errorf("provider not found: %s", s.config.Provider)
    }
    
    s.activeProvider = provider
    s.activeModel = s.config.Model
    
    // Initialize cache if enabled
    if s.config.CacheEnabled {
        s.cache = NewCache(s.config.CacheSize)
    }
    
    return nil
}

func (s *CoreShim) registerBuiltinProviders() error {
    // Register Claude provider
    claudeProvider := providers.NewClaudeProvider()
    s.providers[claudeProvider.Name()] = claudeProvider
    
    // Register OpenAI provider
    openaiProvider := providers.NewOpenAIProvider()
    s.providers[openaiProvider.Name()] = openaiProvider
    
    // Register additional providers
    // ...
    
    return nil
}
```

### Phase 2: Provider Implementations (Weeks 3-4)

1. Implement Claude provider
2. Implement OpenAI provider
3. Implement Llama provider
4. Create provider-specific parameter mapping
5. Add fallback mechanisms

```go
// Example Claude provider implementation
type ClaudeProvider struct {
    client      *anthropic.Client
    modelMap    map[string]string
}

func NewClaudeProvider() Provider {
    return &ClaudeProvider{
        modelMap: map[string]string{
            "claude-3.7-sonnet": "claude-3-sonnet-20240229",
            "claude-3-opus": "claude-3-opus-20240229",
            "claude-3.5-haiku": "claude-3-haiku-20240307",
            // Add other model mappings
        },
    }
}

func (p *ClaudeProvider) Name() string {
    return "claude"
}

func (p *ClaudeProvider) AvailableModels() []string {
    models := make([]string, 0, len(p.modelMap))
    for model := range p.modelMap {
        models = append(models, model)
    }
    return models
}

func (p *ClaudeProvider) GenerateResponse(ctx context.Context, prompt string, params ProviderParams) (string, error) {
    // Map generic model to provider-specific model
    modelName, ok := p.modelMap[params.Model]
    if !ok {
        return "", fmt.Errorf("unsupported model: %s", params.Model)
    }
    
    // Set up Claude-specific request
    message := anthropic.NewMessageRequest().
        WithModel(modelName).
        WithMaxTokens(params.MaxTokens).
        WithTemperature(params.Temperature).
        WithPrompt(prompt)
    
    // Send request to Claude API
    resp, err := p.client.Messages(ctx, message)
    if err != nil {
        return "", err
    }
    
    return resp.Content[0].Text, nil
}
```

### Phase 3: Context Management (Weeks 5-6)

1. Implement conversation history management
2. Build context window tracking
3. Create context pruning strategies
4. Develop memory management
5. Add context prioritization

```go
// Example context management
type Context struct {
    Messages        []Message
    TokenCount      int
    MaxTokens       int
    SystemPrompt    string
    ActiveTools     []Tool
}

func NewContext(maxTokens int, systemPrompt string) *Context {
    return &Context{
        Messages:     make([]Message, 0),
        MaxTokens:    maxTokens,
        SystemPrompt: systemPrompt,
    }
}

func (c *Context) AddMessage(role string, content string) {
    message := Message{
        Role:    role,
        Content: content,
    }
    
    // Estimate token count
    tokens := EstimateTokens(content)
    
    // Add to context
    c.Messages = append(c.Messages, message)
    c.TokenCount += tokens
    
    // Prune if necessary
    c.pruneIfNeeded()
}

func (c *Context) pruneIfNeeded() {
    // If we're under the token limit, no need to prune
    if c.TokenCount < c.MaxTokens {
        return
    }
    
    // Prune strategy: remove oldest messages first, but keep the latest exchange
    // Always preserve system prompt and the most recent user message + assistant response
    
    // Calculate how many tokens we need to remove
    toRemove := c.TokenCount - c.MaxTokens + 200 // Buffer of 200 tokens
    
    // Keep removing messages until we've freed up enough tokens
    removed := 0
    preserveCount := 2 // Preserve the latest exchange
    
    for i := 0; i < len(c.Messages) - preserveCount; i++ {
        messageTokens := EstimateTokens(c.Messages[i].Content)
        
        if removed + messageTokens >= toRemove {
            // We've removed enough tokens
            c.Messages = c.Messages[i:]
            c.TokenCount -= removed
            break
        }
        
        removed += messageTokens
    }
}
```

### Phase 4: Format Converters (Week 7)

1. Implement JSON formatter
2. Create Markdown formatter
3. Add plain text formatter
4. Build format detection
5. Develop structured output parsing

```go
// Example formatter implementation
type Formatter interface {
    Format(input interface{}) (string, error)
    Parse(input string, into interface{}) error
}

type JSONFormatter struct{}

func NewJSONFormatter() Formatter {
    return &JSONFormatter{}
}

func (f *JSONFormatter) Format(input interface{}) (string, error) {
    data, err := json.Marshal(input)
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func (f *JSONFormatter) Parse(input string, into interface{}) error {
    return json.Unmarshal([]byte(input), into)
}

// Format converter registry
type FormatRegistry struct {
    formatters map[string]Formatter
}

func NewFormatRegistry() *FormatRegistry {
    registry := &FormatRegistry{
        formatters: make(map[string]Formatter),
    }
    
    // Register default formatters
    registry.Register("json", NewJSONFormatter())
    registry.Register("markdown", NewMarkdownFormatter())
    registry.Register("text", NewTextFormatter())
    
    return registry
}

func (r *FormatRegistry) Register(name string, formatter Formatter) {
    r.formatters[name] = formatter
}

func (r *FormatRegistry) Get(name string) (Formatter, error) {
    formatter, ok := r.formatters[name]
    if !ok {
        return nil, fmt.Errorf("formatter not found: %s", name)
    }
    return formatter, nil
}
```

### Phase 5: Caching & Optimization (Week 8)

1. Implement response caching
2. Add prompt optimization
3. Create batch processing
4. Develop rate limiting
5. Build request deduplication

```go
// Example cache implementation
type Cache interface {
    Get(key string) (string, bool)
    Set(key string, value string)
    Clear()
}

type InMemoryCache struct {
    cache map[string]string
    mu    sync.RWMutex
    size  int
}

func NewCache(size int) Cache {
    return &InMemoryCache{
        cache: make(map[string]string),
        size:  size,
    }
}

func (c *InMemoryCache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    value, ok := c.cache[key]
    return value, ok
}

func (c *InMemoryCache) Set(key string, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Check if we need to evict entries
    if len(c.cache) >= c.size {
        // Simple eviction strategy: remove a random entry
        for k := range c.cache {
            delete(c.cache, k)
            break
        }
    }
    
    c.cache[key] = value
}

func (c *InMemoryCache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache = make(map[string]string)
}

// Response caching in the shim
func (s *CoreShim) Generate(ctx context.Context, input GenerateInput) (*GenerateOutput, error) {
    // Generate cache key
    cacheKey := generateCacheKey(input)
    
    // Check cache if enabled
    if s.cache != nil {
        if cachedResponse, found := s.cache.Get(cacheKey); found {
            return &GenerateOutput{
                Text:      cachedResponse,
                FromCache: true,
            }, nil
        }
    }
    
    // Prepare prompt for the provider
    prompt, err := s.preparePrompt(input)
    if err != nil {
        return nil, err
    }
    
    // Call provider
    response, err := s.activeProvider.GenerateResponse(ctx, prompt, ProviderParams{
        Model:       s.activeModel,
        MaxTokens:   input.MaxTokens,
        Temperature: input.Temperature,
    })
    if err != nil {
        return nil, err
    }
    
    // Cache response if enabled
    if s.cache != nil {
        s.cache.Set(cacheKey, response)
    }
    
    return &GenerateOutput{
        Text:      response,
        FromCache: false,
    }, nil
}
```

### Phase 6: Tool Integration (Weeks 9-10)

1. Implement tool registry
2. Add basic tool set
3. Create tool response handling
4. Build tool error handling
5. Develop tool authorization

```go
// Example tool implementation
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

type WebSearchTool struct {
    client *search.Client
}

func NewWebSearchTool(client *search.Client) Tool {
    return &WebSearchTool{
        client: client,
    }
}

func (t *WebSearchTool) Name() string {
    return "web_search"
}

func (t *WebSearchTool) Description() string {
    return "Search the web for information"
}

func (t *WebSearchTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Extract parameters
    query, ok := params["query"].(string)
    if !ok {
        return nil, fmt.Errorf("missing 'query' parameter")
    }
    
    limit := 5 // Default limit
    if limitParam, ok := params["limit"].(float64); ok {
        limit = int(limitParam)
    }
    
    // Execute search
    results, err := t.client.Search(ctx, query, limit)
    if err != nil {
        return nil, err
    }
    
    return results, nil
}

// Tool registry
type ToolRegistry struct {
    tools map[string]Tool
}

func NewToolRegistry() *ToolRegistry {
    return &ToolRegistry{
        tools: make(map[string]Tool),
    }
}

func (r *ToolRegistry) Register(tool Tool) {
    r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Get(name string) (Tool, error) {
    tool, ok := r.tools[name]
    if !ok {
        return nil, fmt.Errorf("tool not found: %s", name)
    }
    return tool, nil
}

func (r *ToolRegistry) List() []string {
    names := make([]string, 0, len(r.tools))
    for name := range r.tools {
        names = append(names, name)
    }
    return names
}
```

## Testing Strategy

### Unit Tests

- Test each provider implementation
- Verify correct context management
- Ensure proper format conversion
- Validate caching behavior
- Test tool execution

### Integration Tests

- End-to-end tests with mock LLM providers
- Verify cross-provider compatibility
- Test fallback mechanisms
- Validate context pruning strategies

### Performance Tests

- Benchmark response times
- Measure memory usage
- Test caching effectiveness
- Verify rate limiting behavior

## Deployment Considerations

### Configuration

The Shim will be configured through environment variables or configuration files:

```yaml
shim:
  provider: claude  # Default provider
  model: claude-3.7-sonnet  # Default model
  fallback_providers: [openai, llama]  # Fallback order
  cache:
    enabled: true
    size: 1000  # Number of entries
  context:
    max_tokens: 15000  # Max context window
    system_prompt: "You are a helpful assistant."
```

### Monitoring

The Shim will expose metrics for monitoring:

- Request counts by provider
- Error rates
- Response times
- Cache hit rates
- Token usage

### Logging

Detailed logging will be implemented:

- Request/response pairs (sanitized)
- Provider switches
- Cache events
- Context pruning events
- Tool executions

## Documentation

### API Documentation

```go
// Generate a response from the LLM
func (s *Shim) Generate(ctx context.Context, input GenerateInput) (*GenerateOutput, error)

// GenerateInput defines the input parameters for generation
type GenerateInput struct {
    Prompt      string            // The prompt to send to the LLM
    MaxTokens   int               // Maximum tokens to generate
    Temperature float32           // Temperature (0.0-1.0)
    Tools       []string          // Tools to make available
    Context     *Context          // Conversation context (optional)
    Format      string            // Desired output format (json, text, markdown)
}

// GenerateOutput contains the generated response
type GenerateOutput struct {
    Text      string        // Generated text
    FromCache bool          // Whether the response was from cache
    TokensUsed int          // Number of tokens used
    Tools     []ToolUsage   // Tools used in generation
}
```

### Usage Examples

```go
// Initialize the shim
shim := shim.NewShim()
err := shim.Initialize(shim.Config{
    Provider: "claude",
    Model: "claude-3.7-sonnet",
    CacheEnabled: true,
})
if err != nil {
    log.Fatalf("Failed to initialize shim: %v", err)
}

// Generate a response
output, err := shim.Generate(ctx, shim.GenerateInput{
    Prompt: "What is the capital of France?",
    MaxTokens: 100,
    Temperature: 0.7,
})
if err != nil {
    log.Fatalf("Failed to generate response: %v", err)
}

fmt.Println(output.Text)
```

## Implementation Schedule

| Week | Focus | Deliverables |
|------|-------|-------------|
| 1 | Core Framework | Interface definitions, config management |
| 2 | Core Framework | Provider registration, error handling |
| 3 | Provider Implementations | Claude and OpenAI providers |
| 4 | Provider Implementations | Llama provider, fallback mechanisms |
| 5 | Context Management | Conversation history, context tracking |
| 6 | Context Management | Context pruning, memory management |
| 7 | Format Converters | JSON, Markdown, text formatters |
| 8 | Caching & Optimization | Response caching, prompt optimization |
| 9 | Tool Integration | Tool registry, basic tools |
| 10 | Tool Integration | Tool response handling, authorization |
