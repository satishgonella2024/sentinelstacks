package types

import "time"

// Network represents a communication network for agents
type Network struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Driver    string    `json:"driver"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Agents    []string  `json:"agents"`
}
