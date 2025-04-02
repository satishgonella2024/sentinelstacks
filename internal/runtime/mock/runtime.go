package mock

import (
	"fmt"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/runtime"
)

// Runtime is a mock implementation of the runtime.Runtime interface
type Runtime struct {
	agents map[string]runtime.AgentInfo
}

// NewRuntime creates a new mock runtime with sample data
func NewRuntime() *Runtime {
	mock := &Runtime{
		agents: make(map[string]runtime.AgentInfo),
	}

	// Add some sample agents
	mock.agents["12345678abcd"] = runtime.AgentInfo{
		ID:        "12345678abcdef0123456789",
		Name:      "assistant",
		Image:     "sentinelstacks/assistant:latest",
		Status:    "running",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		Model:     "claude-3-sonnet",
		Memory:    1024 * 1024 * 10,
		APIUsage:  42,
	}

	mock.agents["98765432fedc"] = runtime.AgentInfo{
		ID:        "98765432fedcba0987654321",
		Name:      "researcher",
		Image:     "sentinelstacks/researcher:latest",
		Status:    "stopped",
		CreatedAt: time.Now().Add(-48 * time.Hour),
		Model:     "gpt-4-turbo",
		Memory:    1024 * 1024 * 20,
		APIUsage:  128,
	}

	mock.agents["abcdef123456"] = runtime.AgentInfo{
		ID:        "abcdef123456abcdef123456",
		Name:      "coder",
		Image:     "sentinelstacks/coder:latest",
		Status:    "running",
		CreatedAt: time.Now().Add(-12 * time.Hour),
		Model:     "llama3:8b",
		Memory:    1024 * 1024 * 5,
		APIUsage:  17,
	}

	return mock
}

// GetRunningAgents returns all running agents
func (m *Runtime) GetRunningAgents() ([]runtime.AgentInfo, error) {
	agents := make([]runtime.AgentInfo, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}
	return agents, nil
}

// GetAgent returns information about a specific agent
func (m *Runtime) GetAgent(id string) (runtime.AgentInfo, error) {
	for agentID, agent := range m.agents {
		if agentID == id {
			return agent, nil
		}
	}
	return runtime.AgentInfo{}, fmt.Errorf("agent not found: %s", id)
}

// CreateAgent creates a new agent
func (m *Runtime) CreateAgent(name, image, model string) (*runtime.Agent, error) {
	// This is a stub implementation
	return &runtime.Agent{
		ID:        "new-agent-id",
		Name:      name,
		Image:     image,
		Status:    runtime.StatusCreating,
		CreatedAt: time.Now(),
		Model:     model,
	}, nil
}

// StartAgent starts a previously created agent
func (m *Runtime) StartAgent(id string) error {
	if agent, exists := m.agents[id]; exists {
		agent.Status = "running"
		m.agents[id] = agent
		return nil
	}
	return fmt.Errorf("agent not found: %s", id)
}

// StopAgent stops a running agent
func (m *Runtime) StopAgent(id string) error {
	if agent, exists := m.agents[id]; exists {
		agent.Status = "stopped"
		m.agents[id] = agent
		return nil
	}
	return fmt.Errorf("agent not found: %s", id)
}

// DeleteAgent deletes an agent
func (m *Runtime) DeleteAgent(id string) error {
	if _, exists := m.agents[id]; exists {
		delete(m.agents, id)
		return nil
	}
	return fmt.Errorf("agent not found: %s", id)
}
