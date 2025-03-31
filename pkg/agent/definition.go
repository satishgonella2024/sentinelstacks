package agent

// Definition represents a structured agent definition parsed from a Sentinelfile
type Definition struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	BaseModel    string                 `json:"baseModel"`
	Capabilities []string               `json:"capabilities,omitempty"`
	Tools        []string               `json:"tools,omitempty"`
	StateSchema  map[string]StateField  `json:"stateSchema,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	Lifecycle    Lifecycle              `json:"lifecycle,omitempty"`
}

// StateField represents a field in the agent's state schema
type StateField struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// Lifecycle represents the agent's lifecycle behaviors
type Lifecycle struct {
	Initialization string `json:"initialization,omitempty"`
	Termination    string `json:"termination,omitempty"`
}

// Image represents a built Sentinel Image
type Image struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Tag          string     `json:"tag"`
	CreatedAt    int64      `json:"createdAt"`
	Definition   Definition `json:"definition"`
	Dependencies []string   `json:"dependencies,omitempty"`
}
