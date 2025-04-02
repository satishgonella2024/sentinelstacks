package services

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	"github.com/sentinelstacks/sentinel/pkg/models"
	"github.com/sentinelstacks/sentinel/pkg/repository"
)

// ComposeService provides multi-agent system management functionality
type ComposeService struct {
	repo repository.MultiAgentSystemRepository
}

// NewComposeService creates a new compose service
func NewComposeService(repo repository.MultiAgentSystemRepository) *ComposeService {
	return &ComposeService{repo: repo}
}

// CreateSystem creates a new multi-agent system
func (s *ComposeService) CreateSystem(ctx context.Context, name string, agents map[string]models.AgentConfig) (*models.MultiAgentSystem, error) {
	system := &models.MultiAgentSystem{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		Status:    "running",
		Agents:    agents,
		Networks:  []string{},
		Volumes:   []string{},
		Metadata:  make(map[string]string),
	}
	
	if err := s.repo.Create(ctx, system); err != nil {
		return nil, err
	}
	
	return system, nil
}

// GetSystem gets a system by ID
func (s *ComposeService) GetSystem(ctx context.Context, id string) (*models.MultiAgentSystem, error) {
	return s.repo.Get(ctx, id)
}

// GetSystemByName gets a system by name
func (s *ComposeService) GetSystemByName(ctx context.Context, name string) (*models.MultiAgentSystem, error) {
	return s.repo.GetByName(ctx, name)
}

// ListSystems lists all systems
func (s *ComposeService) ListSystems(ctx context.Context) ([]*models.MultiAgentSystem, error) {
	return s.repo.List(ctx)
}

// DeleteSystem deletes a system
func (s *ComposeService) DeleteSystem(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// UpdateSystemStatus updates the status of a system
func (s *ComposeService) UpdateSystemStatus(ctx context.Context, name, status string) error {
	system, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return err
	}
	
	system.Status = status
	
	return s.repo.Update(ctx, system)
}

// PauseSystem pauses a running system
func (s *ComposeService) PauseSystem(ctx context.Context, name string) error {
	return s.UpdateSystemStatus(ctx, name, "paused")
}

// ResumeSystem resumes a paused system
func (s *ComposeService) ResumeSystem(ctx context.Context, name string) error {
	return s.UpdateSystemStatus(ctx, name, "running")
}

// StopSystem stops a running or paused system
func (s *ComposeService) StopSystem(ctx context.Context, name string) error {
	return s.UpdateSystemStatus(ctx, name, "stopped")
}

// AddAgentToSystem adds an agent to a system
func (s *ComposeService) AddAgentToSystem(ctx context.Context, systemName string, agentName string, config models.AgentConfig) error {
	system, err := s.repo.GetByName(ctx, systemName)
	if err != nil {
		return err
	}
	
	if system.Agents == nil {
		system.Agents = make(map[string]models.AgentConfig)
	}
	
	system.Agents[agentName] = config
	
	return s.repo.Update(ctx, system)
}

// RemoveAgentFromSystem removes an agent from a system
func (s *ComposeService) RemoveAgentFromSystem(ctx context.Context, systemName string, agentName string) error {
	system, err := s.repo.GetByName(ctx, systemName)
	if err != nil {
		return err
	}
	
	if system.Agents == nil {
		return nil
	}
	
	delete(system.Agents, agentName)
	
	return s.repo.Update(ctx, system)
}
