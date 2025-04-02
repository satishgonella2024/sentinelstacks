package repository

import (
	"context"
	
	"github.com/sentinelstacks/sentinel/pkg/models"
)

// VolumeRepository defines the interface for volume data operations
type VolumeRepository interface {
	Create(ctx context.Context, volume *models.Volume) error
	Get(ctx context.Context, id string) (*models.Volume, error)
	GetByName(ctx context.Context, name string) (*models.Volume, error)
	List(ctx context.Context) ([]*models.Volume, error)
	Update(ctx context.Context, volume *models.Volume) error
	Delete(ctx context.Context, id string) error
	Mount(ctx context.Context, volumeID, agentID, mountPath string) error
	Unmount(ctx context.Context, volumeID, agentID string) error
}
