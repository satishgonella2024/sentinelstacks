package models

import (
	"time"
)

// Volume represents a persistent memory volume
type Volume struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Size      string            `json:"size"` // "500MB", "1GB", etc.
	CreatedAt time.Time         `json:"created_at"`
	Encrypted bool              `json:"encrypted"`
	MountPath string            `json:"mount_path,omitempty"`
	MountedBy string            `json:"mounted_by,omitempty"` // Agent ID if mounted
	Used      string            `json:"used,omitempty"`       // Current usage
	Metadata  map[string]string `json:"metadata"`
}
