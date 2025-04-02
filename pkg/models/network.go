package models

import (
	"time"
)

// Network represents an agent communication network
type Network struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Driver    string            `json:"driver"`
	CreatedAt time.Time         `json:"created_at"`
	Status    string            `json:"status"` // "active", "idle"
	Agents    []string          `json:"agents"` // IDs of connected agents
	Metadata  map[string]string `json:"metadata"`
}
