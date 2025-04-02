package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	
	"github.com/google/uuid"
	"github.com/sentinelstacks/sentinel/pkg/models"
	"github.com/sentinelstacks/sentinel/pkg/repository"
)

// FSVolumeRepository implements VolumeRepository using the filesystem
type FSVolumeRepository struct {
	dataDir string
	mutex   sync.RWMutex
}

// NewFSVolumeRepository creates a new filesystem-based volume repository
func NewFSVolumeRepository(dataDir string) repository.VolumeRepository {
	volumesDir := filepath.Join(dataDir, "volumes")
	os.MkdirAll(volumesDir, 0755)
	return &FSVolumeRepository{dataDir: volumesDir}
}

// getFilePath returns the file path for a volume
func (r *FSVolumeRepository) getFilePath(id string) string {
	return filepath.Join(r.dataDir, id+".json")
}

// Create stores a new volume
func (r *FSVolumeRepository) Create(ctx context.Context, volume *models.Volume) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Generate ID if not provided
	if volume.ID == "" {
		volume.ID = uuid.New().String()
	}
	
	// Set creation time if not set
	if volume.CreatedAt.IsZero() {
		volume.CreatedAt = time.Now()
	}
	
	// Initialize metadata if nil
	if volume.Metadata == nil {
		volume.Metadata = make(map[string]string)
	}
	
	// Check for existing volume with same name
	existing, err := r.GetByName(ctx, volume.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("volume with name '%s' already exists", volume.Name)
	}
	
	// Serialize to JSON
	data, err := json.MarshalIndent(volume, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal volume data: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(r.getFilePath(volume.ID), data, 0644); err != nil {
		return fmt.Errorf("failed to write volume data: %w", err)
	}
	
	return nil
}

// Get retrieves a volume by ID
func (r *FSVolumeRepository) Get(ctx context.Context, id string) (*models.Volume, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	data, err := os.ReadFile(r.getFilePath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("volume not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read volume data: %w", err)
	}
	
	var volume models.Volume
	if err := json.Unmarshal(data, &volume); err != nil {
		return nil, fmt.Errorf("failed to unmarshal volume data: %w", err)
	}
	
	return &volume, nil
}

// GetByName retrieves a volume by name
func (r *FSVolumeRepository) GetByName(ctx context.Context, name string) (*models.Volume, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	volumes, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	
	for _, volume := range volumes {
		if volume.Name == name {
			return volume, nil
		}
	}
	
	return nil, fmt.Errorf("volume not found: %s", name)
}

// List returns all volumes
func (r *FSVolumeRepository) List(ctx context.Context) ([]*models.Volume, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	files, err := os.ReadDir(r.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read volumes directory: %w", err)
	}
	
	var volumes []*models.Volume
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		data, err := os.ReadFile(filepath.Join(r.dataDir, file.Name()))
		if err != nil {
			continue
		}
		
		var volume models.Volume
		if err := json.Unmarshal(data, &volume); err != nil {
			continue
		}
		
		volumes = append(volumes, &volume)
	}
	
	return volumes, nil
}

// Update updates an existing volume
func (r *FSVolumeRepository) Update(ctx context.Context, volume *models.Volume) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if volume exists
	existing, err := r.Get(ctx, volume.ID)
	if err != nil {
		return err
	}
	
	// Preserve creation time
	volume.CreatedAt = existing.CreatedAt
	
	// Serialize to JSON
	data, err := json.MarshalIndent(volume, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal volume data: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(r.getFilePath(volume.ID), data, 0644); err != nil {
		return fmt.Errorf("failed to write volume data: %w", err)
	}
	
	return nil
}

// Delete removes a volume
func (r *FSVolumeRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if volume exists
	volume, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	
	// Check if volume is mounted
	if volume.MountedBy != "" {
		return fmt.Errorf("volume is mounted by agent '%s' and cannot be removed", volume.MountedBy)
	}
	
	// Remove file
	if err := os.Remove(r.getFilePath(id)); err != nil {
		return fmt.Errorf("failed to delete volume data: %w", err)
	}
	
	return nil
}

// Mount mounts a volume to an agent
func (r *FSVolumeRepository) Mount(ctx context.Context, volumeID, agentID, mountPath string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Get volume
	volume, err := r.Get(ctx, volumeID)
	if err != nil {
		return err
	}
	
	// Check if volume is already mounted
	if volume.MountedBy != "" {
		return fmt.Errorf("volume is already mounted by agent '%s'", volume.MountedBy)
	}
	
	// Mount volume
	volume.MountedBy = agentID
	volume.MountPath = mountPath
	
	// Update volume
	return r.Update(ctx, volume)
}

// Unmount unmounts a volume from an agent
func (r *FSVolumeRepository) Unmount(ctx context.Context, volumeID, agentID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Get volume
	volume, err := r.Get(ctx, volumeID)
	if err != nil {
		return err
	}
	
	// Check if volume is mounted by the specified agent
	if volume.MountedBy != agentID {
		return fmt.Errorf("volume is not mounted by agent '%s'", agentID)
	}
	
	// Unmount volume
	volume.MountedBy = ""
	volume.MountPath = ""
	
	// Update volume
	return r.Update(ctx, volume)
}
