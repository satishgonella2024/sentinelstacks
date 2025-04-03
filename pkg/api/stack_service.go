// Package api provides a unified API for the Sentinel Stacks system
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/satishgonella2024/sentinelstacks/pkg/stack"
	"github.com/satishgonella2024/sentinelstacks/pkg/storage"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// StackServiceConfig contains configuration for the stack service
type StackServiceConfig struct {
	// StoragePath is where stack definitions are persisted
	StoragePath string

	// Verbose enables verbose logging
	Verbose bool
}

// StackService implements types.StackService
type StackService struct {
	config  StackServiceConfig
	stacks  map[string]*stackInfo
	storage *storage.Storage
	mu      sync.RWMutex
}

// stackInfo contains internal information about a stack
type stackInfo struct {
	id          string
	spec        types.StackSpec
	engine      *stack.Engine
	createdAt   time.Time
	lastUpdated time.Time
	status      types.StackStatus
}

// NewStackService creates a new stack service
func NewStackService(config StackServiceConfig) (*StackService, error) {
	var stor *storage.Storage
	var err error

	if config.StoragePath != "" {
		// Create storage instance
		stor, err = storage.NewStorage(config.StoragePath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize storage: %w", err)
		}
	}

	service := &StackService{
		config:  config,
		stacks:  make(map[string]*stackInfo),
		storage: stor,
	}

	// Load existing stacks if storage is enabled
	if stor != nil {
		if err := service.loadStacks(); err != nil {
			return nil, fmt.Errorf("failed to load stacks: %w", err)
		}
	}

	return service, nil
}

// CreateStack creates a new stack
func (s *StackService) CreateStack(ctx context.Context, spec types.StackSpec) (string, error) {
	// Generate a unique ID for the stack
	stackID := uuid.New().String()

	// Create a new stack engine
	engine, err := stack.NewEngine(spec, stack.WithVerbose(s.config.Verbose))
	if err != nil {
		return "", fmt.Errorf("failed to create stack engine: %w", err)
	}

	// Set creation time
	now := time.Now()

	// Store stack information
	info := &stackInfo{
		id:          stackID,
		spec:        spec,
		engine:      engine,
		createdAt:   now,
		lastUpdated: now,
		status:      types.StackStatusReady,
	}

	// Store in memory
	s.mu.Lock()
	s.stacks[stackID] = info
	s.mu.Unlock()

	// Persist stack definition if storage is enabled
	if s.storage != nil {
		storageInfo := storage.StackInfo{
			ID:          stackID,
			Spec:        spec,
			CreatedAt:   now,
			LastUpdated: now,
			Status:      types.StackStatusReady,
		}

		if err := s.storage.SaveStack(storageInfo); err != nil {
			// Log the error but continue with in-memory operation
			fmt.Printf("Warning: Failed to persist stack: %v\n", err)
		}
	}

	return stackID, nil
}

// ExecuteStack executes a stack with given inputs
func (s *StackService) ExecuteStack(ctx context.Context, stackID string, inputs map[string]interface{}) (map[string]interface{}, error) {
	// Get the stack
	s.mu.RLock()
	info, exists := s.stacks[stackID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("stack not found: %s", stackID)
	}

	// Create execution ID
	executionID := uuid.New().String()
	startTime := time.Now()

	// Create a channel to collect outputs
	resultChan := make(chan map[string]interface{}, 1)
	errChan := make(chan error, 1)

	// Execute the stack in a separate goroutine
	go func() {
		// Set up execution options
		options := []stack.ExecuteOption{
			stack.WithInput(inputs),
		}

		// Execute the stack
		err := info.engine.Execute(ctx, options...)
		if err != nil {
			errChan <- err
			return
		}

		// For now, we'll just return the inputs as outputs
		// In a real implementation, we would collect outputs from the final agent
		outputs := make(map[string]interface{})
		for k, v := range inputs {
			outputs[k] = v
		}
		outputs["status"] = "success"
		outputs["stack_id"] = stackID
		outputs["execution_time"] = time.Now().Format(time.RFC3339)

		resultChan <- outputs
	}()

	// Wait for result or error
	var result map[string]interface{}
	var execErr error

	select {
	case <-ctx.Done():
		execErr = ctx.Err()
	case err := <-errChan:
		execErr = err
	case res := <-resultChan:
		result = res
	}

	// Record execution summary if storage is enabled
	if s.storage != nil {
		endTime := time.Now()
		status := types.StackStatusSucceeded
		if execErr != nil {
			status = types.StackStatusFailed
		}

		// Create execution summary
		summary := storage.ExecutionSummary{
			ExecutionID:    executionID,
			StartTime:      startTime,
			EndTime:        endTime,
			Status:         status,
			CompletedCount: len(info.spec.Agents),
			FailedCount:    0,
			BlockedCount:   0,
			Inputs:         inputs,
			Outputs:        result,
		}

		// If execution failed, record the failure count
		if status == types.StackStatusFailed {
			summary.FailedCount = 1
			summary.CompletedCount = len(info.spec.Agents) - 1
		}

		// Update execution history
		if err := s.storage.UpdateStackExecution(stackID, summary); err != nil {
			// Log the error but continue with operation
			fmt.Printf("Warning: Failed to update stack execution: %v\n", err)
		}

		// Update stack status in memory
		s.mu.Lock()
		if info, ok := s.stacks[stackID]; ok {
			info.status = status
			info.lastUpdated = endTime
		}
		s.mu.Unlock()
	}

	if execErr != nil {
		return nil, execErr
	}

	return result, nil
}

// GetStackState gets the current state of a stack
func (s *StackService) GetStackState(ctx context.Context, stackID string) (*types.StackExecutionSummary, error) {
	// Get the stack
	s.mu.RLock()
	info, exists := s.stacks[stackID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("stack not found: %s", stackID)
	}

	// Try to get execution history from storage
	if s.storage != nil {
		storedInfo, err := s.storage.GetStack(stackID)
		if err == nil && len(storedInfo.Executions) > 0 {
			// Get the most recent execution
			latestExec := storedInfo.Executions[0]

			// Create summary from the stored execution
			summary := &types.StackExecutionSummary{
				StackName:      info.spec.Name,
				ExecutionID:    latestExec.ExecutionID,
				StartTime:      latestExec.StartTime,
				EndTime:        latestExec.EndTime,
				TotalAgents:    len(info.spec.Agents),
				CompletedCount: latestExec.CompletedCount,
				FailedCount:    latestExec.FailedCount,
				BlockedCount:   latestExec.BlockedCount,
				AgentStates:    make(map[string]*types.AgentState),
			}

			return summary, nil
		}
	}

	// Fallback to a generic summary if no execution history is available
	summary := &types.StackExecutionSummary{
		StackName:      info.spec.Name,
		ExecutionID:    uuid.New().String(),
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalAgents:    len(info.spec.Agents),
		CompletedCount: len(info.spec.Agents),
		FailedCount:    0,
		BlockedCount:   0,
		AgentStates:    make(map[string]*types.AgentState),
	}

	return summary, nil
}

// ListStacks lists all available stacks
func (s *StackService) ListStacks(ctx context.Context) ([]types.StackInfo, error) {
	// If storage is enabled, try to get the list from storage first
	if s.storage != nil {
		storedStacks, err := s.storage.ListStacks()
		if err == nil {
			result := make([]types.StackInfo, 0, len(storedStacks))
			for _, stored := range storedStacks {
				result = append(result, types.StackInfo{
					ID:          stored.ID,
					Name:        stored.Spec.Name,
					Description: stored.Spec.Description,
					Version:     stored.Spec.Version,
					Type:        stored.Spec.Type,
					CreatedAt:   stored.CreatedAt.Format(time.RFC3339),
				})
			}
			return result, nil
		}
	}

	// Fallback to in-memory data if storage failed or is not enabled
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create list of stack info
	result := make([]types.StackInfo, 0, len(s.stacks))

	for id, info := range s.stacks {
		result = append(result, types.StackInfo{
			ID:          id,
			Name:        info.spec.Name,
			Description: info.spec.Description,
			Version:     info.spec.Version,
			Type:        info.spec.Type,
			CreatedAt:   info.createdAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

// UpdateStack updates an existing stack
func (s *StackService) UpdateStack(ctx context.Context, stackID string, spec types.StackSpec) error {
	// Check if stack exists
	s.mu.RLock()
	_, exists := s.stacks[stackID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("stack not found: %s", stackID)
	}

	// Create a new stack engine
	engine, err := stack.NewEngine(spec, stack.WithVerbose(s.config.Verbose))
	if err != nil {
		return fmt.Errorf("failed to create stack engine: %w", err)
	}

	// Update in memory
	now := time.Now()
	s.mu.Lock()
	s.stacks[stackID] = &stackInfo{
		id:          stackID,
		spec:        spec,
		engine:      engine,
		createdAt:   s.stacks[stackID].createdAt,
		lastUpdated: now,
		status:      types.StackStatusReady,
	}
	s.mu.Unlock()

	// Update in storage if enabled
	if s.storage != nil {
		// Get the current stored stack to preserve execution history
		storedStack, err := s.storage.GetStack(stackID)
		if err == nil {
			// Update fields but preserve execution history
			storedStack.Spec = spec
			storedStack.LastUpdated = now
			storedStack.Status = types.StackStatusReady

			if err := s.storage.SaveStack(storedStack); err != nil {
				// Log the error but continue with in-memory operation
				fmt.Printf("Warning: Failed to update stack in storage: %v\n", err)
			}
		} else {
			// If getting the stored stack failed, save a new one
			storageInfo := storage.StackInfo{
				ID:          stackID,
				Spec:        spec,
				CreatedAt:   s.stacks[stackID].createdAt,
				LastUpdated: now,
				Status:      types.StackStatusReady,
			}

			if err := s.storage.SaveStack(storageInfo); err != nil {
				// Log the error but continue with in-memory operation
				fmt.Printf("Warning: Failed to update stack in storage: %v\n", err)
			}
		}
	}

	return nil
}

// DeleteStack removes a stack
func (s *StackService) DeleteStack(ctx context.Context, stackID string) error {
	// Check if stack exists
	s.mu.RLock()
	_, exists := s.stacks[stackID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("stack not found: %s", stackID)
	}

	// Remove from memory
	s.mu.Lock()
	delete(s.stacks, stackID)
	s.mu.Unlock()

	// Remove from storage if enabled
	if s.storage != nil {
		if err := s.storage.DeleteStack(stackID); err != nil {
			// Log the error but continue with in-memory operation
			fmt.Printf("Warning: Failed to delete stack from storage: %v\n", err)
		}
	}

	return nil
}

// ImportStack imports a stack from a file
func (s *StackService) ImportStack(ctx context.Context, filePath string) (string, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Parse stack definition
	var spec types.StackSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return "", fmt.Errorf("failed to parse stack definition: %w", err)
	}

	// Create stack
	return s.CreateStack(ctx, spec)
}

// ExportStack exports a stack to a file
func (s *StackService) ExportStack(ctx context.Context, stackID, filePath string) error {
	// Get stack
	s.mu.RLock()
	info, exists := s.stacks[stackID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("stack not found: %s", stackID)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Convert to JSON
	data, err := json.MarshalIndent(info.spec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stack definition: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write stack definition: %w", err)
	}

	return nil
}

// loadStacks loads all stacks from storage
func (s *StackService) loadStacks() error {
	if s.storage == nil {
		return nil
	}

	// Get stacks from storage
	storedStacks, err := s.storage.ListStacks()
	if err != nil {
		return fmt.Errorf("failed to list stacks from storage: %w", err)
	}

	// Load each stack
	for _, stored := range storedStacks {
		// Create a new stack engine
		engine, err := stack.NewEngine(stored.Spec, stack.WithVerbose(s.config.Verbose))
		if err != nil {
			fmt.Printf("Warning: Failed to create engine for stack %s: %v\n", stored.ID, err)
			continue
		}

		// Store in memory
		s.stacks[stored.ID] = &stackInfo{
			id:          stored.ID,
			spec:        stored.Spec,
			engine:      engine,
			createdAt:   stored.CreatedAt,
			lastUpdated: stored.LastUpdated,
			status:      stored.Status,
		}
	}

	return nil
}

// GetStackExecutionHistory gets the execution history for a stack
func (s *StackService) GetStackExecutionHistory(ctx context.Context, stackID string) ([]types.ExecutionSummary, error) {
	// Check if stack exists
	s.mu.RLock()
	_, exists := s.stacks[stackID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("stack not found: %s", stackID)
	}

	// Get execution history from storage
	if s.storage != nil {
		storedInfo, err := s.storage.GetStack(stackID)
		if err != nil {
			return nil, fmt.Errorf("failed to get stack from storage: %w", err)
		}

		// Convert to API types
		result := make([]types.ExecutionSummary, 0, len(storedInfo.Executions))
		for _, exec := range storedInfo.Executions {
			result = append(result, types.ExecutionSummary{
				ExecutionID:    exec.ExecutionID,
				StartTime:      exec.StartTime,
				EndTime:        exec.EndTime,
				Status:         exec.Status,
				CompletedCount: exec.CompletedCount,
				FailedCount:    exec.FailedCount,
				BlockedCount:   exec.BlockedCount,
			})
		}

		return result, nil
	}

	// If storage is not enabled, return empty history
	return []types.ExecutionSummary{}, nil
}
