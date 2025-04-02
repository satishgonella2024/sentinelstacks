package models

import (
	"time"
)

// AgentConfig represents an agent configuration within a multi-agent system
type AgentConfig struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Networks    []string          `json:"networks,omitempty"`
	Volumes     []string          `json:"volumes,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Resources   Resources         `json:"resources,omitempty"`
}

// Resources defines computational resources for an agent
type Resources struct {
	Memory     string `json:"memory,omitempty"`
	CPULimit   string `json:"cpu_limit,omitempty"`
	GPUEnabled bool   `json:"gpu_enabled,omitempty"`
}

// MultiAgentSystem represents a composition of multiple agents
type MultiAgentSystem struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	Status    string                 `json:"status"` // "running", "paused", "stopped"
	Agents    map[string]AgentConfig `json:"agents"`
	Networks  []string               `json:"networks"`
	Volumes   []string               `json:"volumes"`
	Metadata  map[string]string      `json:"metadata"`
}
