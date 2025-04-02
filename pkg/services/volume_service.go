package services

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	"github.com/sentinelstacks/sentinel/pkg/models"
	"github.com/sentinelstacks/sentinel/pkg/repository"
)

// VolumeService provides volume management functionality
type VolumeService struct {
	repo repository.VolumeRepository
}

// NewVolumeService creates a new volume service
func NewVolumeService(repo repository.VolumeRepository) *VolumeService {
	return &VolumeService{repo: repo}
}

// CreateVolume creates a new persistent memory volume
func (s *VolumeService) CreateVolume(ctx context.Context, name, size string, encrypted bool) (*models.Volume, error) {
	volume := &models.Volume{
		ID:        uuid.New().String(),
		Name:      name,
		Size:      size,
		CreatedAt: time.Now(),
		Encrypted: encrypted,
		Used:      "0",
		Metadata:  make(map[string]string),
	}
	
	if err := s.repo.Create(ctx, volume); err != nil {
		return nil, err
	}
	
	return volume, nil
}

// GetVolume gets a volume by ID
func (s *VolumeService) GetVolume(ctx context.Context, id string) (*models.Volume, error) {
	return s.repo.Get(ctx, id)
}

// GetVolumeByName gets a volume by name
func (s *VolumeService) GetVolumeByName(ctx context.Context, name string) (*models.Volume, error) {
	return s.repo.GetByName(ctx, name)
}

// ListVolumes lists all volumes
func (s *VolumeService) ListVolumes(ctx context.Context) ([]*models.Volume, error) {
	return s.repo.List(ctx)
}

// DeleteVolume deletes a volume
func (s *VolumeService) DeleteVolume(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// MountVolume mounts a volume to an agent
func (s *VolumeService) MountVolume(ctx context.Context, volumeName, agentID, mountPath string) error {
	volume, err := s.repo.GetByName(ctx, volumeName)
	if err != nil {
		return err
	}
	
	return s.repo.Mount(ctx, volume.ID, agentID, mountPath)
}

// UnmountVolume unmounts a volume from an agent
func (s *VolumeService) UnmountVolume(ctx context.Context, volumeName, agentID string) error {
	volume, err := s.repo.GetByName(ctx, volumeName)
	if err != nil {
		return err
	}
	
	return s.repo.Unmount(ctx, volume.ID, agentID)
}

// InspectVolume returns detailed information about a volume
func (s *VolumeService) InspectVolume(ctx context.Context, name string) (*models.Volume, error) {
	return s.repo.GetByName(ctx, name)
}
