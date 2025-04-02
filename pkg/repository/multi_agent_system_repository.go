package repository

import (
	"context"
	
	"github.com/sentinelstacks/sentinel/pkg/models"
)

// MultiAgentSystemRepository defines the interface for multi-agent system data operations
type MultiAgentSystemRepository interface {
	Create(ctx context.Context, system *models.MultiAgentSystem) error
	Get(ctx context.Context, id string) (*models.MultiAgentSystem, error)
	GetByName(ctx context.Context, name string) (*models.MultiAgentSystem, error)
	List(ctx context.Context) ([]*models.MultiAgentSystem, error)
	Update(ctx context.Context, system *models.MultiAgentSystem) error
	Delete(ctx context.Context, id string) error
}
