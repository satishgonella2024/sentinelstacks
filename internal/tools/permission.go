package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Permission represents a permission level
type Permission string

const (
	// PermissionNone represents no permissions
	PermissionNone Permission = "none"
	
	// PermissionFile represents file access permissions
	PermissionFile Permission = "file"
	
	// PermissionNetwork represents network access permissions
	PermissionNetwork Permission = "network"
	
	// PermissionShell represents shell access permissions
	PermissionShell Permission = "shell"
	
	// PermissionAPI represents API access permissions
	PermissionAPI Permission = "api"
	
	// PermissionAll represents all permissions
	PermissionAll Permission = "all"
)

// AgentPermissions stores permissions for an agent
type AgentPermissions struct {
	AgentID     string       `json:"agent_id"`
	Permissions []Permission `json:"permissions"`
}

// PermissionManager manages tool permissions
type PermissionManager struct {
	permissions map[string][]Permission
	dataDir     string
	mu          sync.RWMutex
}

// NewPermissionManager creates a new permission manager
func NewPermissionManager(dataDir string) (*PermissionManager, error) {
	// If data directory not specified, use default
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not get home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".sentinel")
	}

	// Create full path for permissions data
	permissionsDir := filepath.Join(dataDir, "permissions")
	if err := os.MkdirAll(permissionsDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create permissions directory: %w", err)
	}

	manager := &PermissionManager{
		permissions: make(map[string][]Permission),
		dataDir:     permissionsDir,
	}

	// Load existing permissions
	if err := manager.load(); err != nil {
		return nil, fmt.Errorf("could not load permissions: %w", err)
	}

	return manager, nil
}

// load loads permissions from the data directory
func (p *PermissionManager) load() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Clear existing permissions
	p.permissions = make(map[string][]Permission)

	// Get all files in the permissions directory
	entries, err := os.ReadDir(p.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			// If directory doesn't exist, return without error
			return nil
		}
		return fmt.Errorf("failed to read permissions directory: %w", err)
	}

	// Iterate through each file
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Read the file
		filePath := filepath.Join(p.dataDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read permissions file: %w", err)
		}

		// Parse the file
		var agentPerms AgentPermissions
		if err := json.Unmarshal(data, &agentPerms); err != nil {
			return fmt.Errorf("failed to parse permissions: %w", err)
		}

		// Add to map
		p.permissions[agentPerms.AgentID] = agentPerms.Permissions
	}

	return nil
}

// save saves permissions to the data directory
func (p *PermissionManager) save() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Save each agent's permissions
	for agentID, perms := range p.permissions {
		// Create agent permissions object
		agentPerms := AgentPermissions{
			AgentID:     agentID,
			Permissions: perms,
		}

		// Marshal to JSON
		data, err := json.MarshalIndent(agentPerms, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal permissions: %w", err)
		}

		// Write to file
		filePath := filepath.Join(p.dataDir, fmt.Sprintf("%s.json", agentID))
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write permissions file: %w", err)
		}
	}

	return nil
}

// Grant grants a permission to an agent
func (p *PermissionManager) Grant(agentID string, perm Permission) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If agent doesn't exist, create entry
	if _, exists := p.permissions[agentID]; !exists {
		p.permissions[agentID] = []Permission{}
	}

	// Special case for "all" permission
	if perm == PermissionAll {
		p.permissions[agentID] = []Permission{PermissionAll}
		return p.save()
	}

	// Check if already has permission
	for _, existingPerm := range p.permissions[agentID] {
		if existingPerm == perm {
			return nil // Already has permission
		}
		
		// If already has "all" permission, no need to add specific ones
		if existingPerm == PermissionAll {
			return nil
		}
	}

	// Add permission
	p.permissions[agentID] = append(p.permissions[agentID], perm)

	// Save changes
	return p.save()
}

// Revoke revokes a permission from an agent
func (p *PermissionManager) Revoke(agentID string, perm Permission) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If agent doesn't exist, nothing to revoke
	perms, exists := p.permissions[agentID]
	if !exists {
		return nil
	}

	// Special case for "all" permission
	if perm == PermissionAll {
		p.permissions[agentID] = []Permission{}
		return p.save()
	}

	// Create new permissions list without the revoked permission
	newPerms := []Permission{}
	for _, existingPerm := range perms {
		if existingPerm != perm {
			// If we're keeping "all" permission, don't add specific ones
			if existingPerm == PermissionAll {
				return p.save() // Just keep "all" permission
			}
			newPerms = append(newPerms, existingPerm)
		}
	}

	// Update permissions
	p.permissions[agentID] = newPerms

	// Save changes
	return p.save()
}

// HasPermission checks if an agent has a permission
func (p *PermissionManager) HasPermission(agentID string, perm Permission) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Get agent permissions
	perms, exists := p.permissions[agentID]
	if !exists {
		return false
	}

	// Check for permission
	for _, existingPerm := range perms {
		if existingPerm == perm || existingPerm == PermissionAll {
			return true
		}
	}

	return false
}

// GetPermissions returns all permissions for an agent
func (p *PermissionManager) GetPermissions(agentID string) []Permission {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Get agent permissions
	perms, exists := p.permissions[agentID]
	if !exists {
		return []Permission{}
	}

	// Return a copy of the permissions
	result := make([]Permission, len(perms))
	copy(result, perms)
	return result
}

// RemoveAgent removes all permissions for an agent
func (p *PermissionManager) RemoveAgent(agentID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Remove from memory
	delete(p.permissions, agentID)

	// Remove file if it exists
	filePath := filepath.Join(p.dataDir, fmt.Sprintf("%s.json", agentID))
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove permissions file: %w", err)
	}

	return nil
}

// GetAgentIDs returns all agent IDs with permissions
func (p *PermissionManager) GetAgentIDs() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Get all agent IDs
	agentIDs := make([]string, 0, len(p.permissions))
	for agentID := range p.permissions {
		agentIDs = append(agentIDs, agentID)
	}

	return agentIDs
}