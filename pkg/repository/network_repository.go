package repository

import (
	"context"
	
	"github.com/satishgonella2024/sentinelstacks/pkg/models"
)

// NetworkRepository defines the interface for network data operations
type NetworkRepository interface {
	Create(ctx context.Context, network *models.Network) error
	Get(ctx context.Context, id string) (*models.Network, error)
	GetByName(ctx context.Context, name string) (*models.Network, error)
	List(ctx context.Context) ([]*models.Network, error)
	Update(ctx context.Context, network *models.Network) error
	Delete(ctx context.Context, id string) error
	ConnectAgent(ctx context.Context, networkID, agentID string) error
	DisconnectAgent(ctx context.Context, networkID, agentID string) error
}
