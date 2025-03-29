package agentfile

// Agentfile represents the structure of an agent definition
type Agentfile struct {
	Name         string       `yaml:"name" json:"name"`
	Version      string       `yaml:"version" json:"version"`
	Description  string       `yaml:"description" json:"description"`
	Model        ModelConfig  `yaml:"model" json:"model"`
	Capabilities []string     `yaml:"capabilities" json:"capabilities"`
	Memory       MemoryConfig `yaml:"memory" json:"memory"`
	Tools        []ToolConfig `yaml:"tools,omitempty" json:"tools,omitempty"`
	Permissions  Permissions  `yaml:"permissions,omitempty" json:"permissions,omitempty"`
	Author       string       `yaml:"author,omitempty" json:"author,omitempty"`
	Tags         []string     `yaml:"tags,omitempty" json:"tags,omitempty"`
	Registry     RegistryInfo `yaml:"registry,omitempty" json:"registry,omitempty"`
}

// ModelConfig defines which AI model to use
type ModelConfig struct {
	Provider string                 `yaml:"provider" json:"provider"`
	Name     string                 `yaml:"name" json:"name"`
	Endpoint string                 `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Options  map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

// MemoryConfig defines how the agent stores state
type MemoryConfig struct {
	Type           string                 `yaml:"type" json:"type"`
	Persistence    bool                   `yaml:"persistence" json:"persistence"`
	MaxItems       int                    `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
	EmbeddingModel string                 `yaml:"embeddingModel,omitempty" json:"embeddingModel,omitempty"`
	VectorOptions  map[string]interface{} `yaml:"vectorOptions,omitempty" json:"vectorOptions,omitempty"`
}

// ToolConfig defines a tool the agent can use
type ToolConfig struct {
	ID      string `yaml:"id" json:"id"`
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
}

// Permissions defines what the agent is allowed to do
type Permissions struct {
	FileAccess []string `yaml:"file_access,omitempty" json:"file_access,omitempty"`
	Network    bool     `yaml:"network" json:"network"`
}

// RegistryInfo contains metadata about the agent in the registry
type RegistryInfo struct {
	Source     string `yaml:"source" json:"source"`
	Visibility string `yaml:"visibility" json:"visibility"`
	PulledAt   string `yaml:"pulled_at,omitempty" json:"pulled_at,omitempty"`
	PushedAt   string `yaml:"pushed_at,omitempty" json:"pushed_at,omitempty"`
}

// DefaultAgentfile creates a default agent configuration
func DefaultAgentfile(name string) Agentfile {
	return Agentfile{
		Name:        name,
		Version:     "0.1.0",
		Description: "A SentinelStacks agent",
		Model: ModelConfig{
			Provider: "ollama",
			Name:     "llama3",
			Options: map[string]interface{}{
				"temperature": 0.7,
			},
		},
		Capabilities: []string{"conversation"},
		Memory: MemoryConfig{
			Type:        "simple",
			Persistence: true,
			MaxItems:    1000,
		},
		Permissions: Permissions{
			FileAccess: []string{"read"},
			Network:    false,
		},
	}
}

// DefaultVectorAgentfile creates a default agent configuration with vector memory
func DefaultVectorAgentfile(name string) Agentfile {
	agent := DefaultAgentfile(name)
	agent.Memory = MemoryConfig{
		Type:           "vector",
		Persistence:    true,
		MaxItems:       1000,
		EmbeddingModel: "openai:text-embedding-3-small",
		VectorOptions: map[string]interface{}{
			"similarityType": "cosine",
		},
	}
	return agent
}
