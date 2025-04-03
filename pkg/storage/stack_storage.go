package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// StackInfo represents a serializable stack information record
type StackInfo struct {
	ID          string              `json:"id"`
	Spec        types.StackSpec     `json:"spec"`
	CreatedAt   time.Time           `json:"created_at"`
	LastUpdated time.Time           `json:"last_updated"`
	Status      types.StackStatus   `json:"status"`
	Executions  []ExecutionSummary  `json:"executions,omitempty"`
}

// ExecutionSummary provides a summary of stack execution
type ExecutionSummary struct {
	ExecutionID    string                  `json:"execution_id"`
	StartTime      time.Time               `json:"start_time"`
	EndTime        time.Time               `json:"end_time"`
	Status         types.StackStatus       `json:"status"`
	CompletedCount int                     `json:"completed_count"`
	FailedCount    int                     `json:"failed_count"`
	BlockedCount   int                     `json:"blocked_count"`
	Inputs         map[string]interface{}  `json:"inputs,omitempty"`
	Outputs        map[string]interface{}  `json:"outputs,omitempty"`
}

// SaveStack persists a stack to storage
func (s *Storage) SaveStack(info StackInfo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create stacks directory if it doesn't exist
	stacksDir := filepath.Join(s.baseDir, "stacks")
	if err := os.MkdirAll(stacksDir, 0755); err != nil {
		return fmt.Errorf("failed to create stacks directory: %w", err)
	}

	// Update last updated time
	info.LastUpdated = time.Now()

	// Convert to JSON
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stack data: %w", err)
	}

	// Write to file
	filePath := filepath.Join(stacksDir, info.ID+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write stack data: %w", err)
	}

	return nil
}

// GetStack retrieves a stack from storage
func (s *Storage) GetStack(id string) (StackInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var info StackInfo

	// Read file
	filePath := filepath.Join(s.baseDir, "stacks", id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return info, fmt.Errorf("stack not found: %s", id)
		}
		return info, fmt.Errorf("failed to read stack data: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &info); err != nil {
		return info, fmt.Errorf("failed to unmarshal stack data: %w", err)
	}

	return info, nil
}

// ListStacks lists all stacks
func (s *Storage) ListStacks() ([]StackInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var stacks []StackInfo

	// Create stacks directory if it doesn't exist
	stacksDir := filepath.Join(s.baseDir, "stacks")
	if err := os.MkdirAll(stacksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create stacks directory: %w", err)
	}

	// Read directory
	files, err := os.ReadDir(stacksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return stacks, nil
		}
		return nil, fmt.Errorf("failed to read stacks directory: %w", err)
	}

	// Load each stack
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(stacksDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files that can't be read
		}

		var info StackInfo
		if err := json.Unmarshal(data, &info); err != nil {
			continue // Skip files that can't be parsed
		}

		stacks = append(stacks, info)
	}

	return stacks, nil
}

// DeleteStack removes a stack from storage
func (s *Storage) DeleteStack(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if stack exists
	filePath := filepath.Join(s.baseDir, "stacks", id+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("stack not found: %s", id)
	}

	// Remove file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete stack data: %w", err)
	}

	return nil
}

// UpdateStackExecution adds or updates an execution summary for a stack
func (s *Storage) UpdateStackExecution(id string, execution ExecutionSummary) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Read current stack info
	filePath := filepath.Join(s.baseDir, "stacks", id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("stack not found: %s", id)
		}
		return fmt.Errorf("failed to read stack data: %w", err)
	}

	var info StackInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return fmt.Errorf("failed to unmarshal stack data: %w", err)
	}

	// Update executions, limiting the history to the most recent 10 executions
	updated := false
	for i, exec := range info.Executions {
		if exec.ExecutionID == execution.ExecutionID {
			info.Executions[i] = execution
			updated = true
			break
		}
	}

	if !updated {
		// Add new execution to the beginning of the slice
		info.Executions = append([]ExecutionSummary{execution}, info.Executions...)
		
		// Limit executions history to 10 entries
		if len(info.Executions) > 10 {
			info.Executions = info.Executions[:10]
		}
	}

	// Update status based on the latest execution
	if len(info.Executions) > 0 {
		info.Status = info.Executions[0].Status
	}

	// Update last updated time
	info.LastUpdated = time.Now()

	// Write back to file
	data, err = json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stack data: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write stack data: %w", err)
	}

	return nil
}

// GetStackByName finds a stack by name
func (s *Storage) GetStackByName(name string) (StackInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// List all stacks
	stacks, err := s.ListStacks()
	if err != nil {
		return StackInfo{}, err
	}

	// Find stack by name
	for _, stack := range stacks {
		if stack.Spec.Name == name {
			return stack, nil
		}
	}

	return StackInfo{}, fmt.Errorf("stack not found: %s", name)
}
