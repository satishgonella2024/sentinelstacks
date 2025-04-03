package runtime

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// SimpleAgentRuntime is a basic implementation for testing
type SimpleAgentRuntime struct {
	logToConsole bool
}

// NewSimpleAgentRuntime creates a new simple agent runtime
func NewSimpleAgentRuntime(logToConsole bool) (*SimpleAgentRuntime, error) {
	return &SimpleAgentRuntime{
		logToConsole: logToConsole,
	}, nil
}

// Execute simulates agent execution
func (r *SimpleAgentRuntime) Execute(ctx context.Context, agentSpec types.StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error) {
	// Log execution if enabled
	if r.logToConsole {
		log.Printf("Executing agent %s (uses: %s)", agentSpec.ID, agentSpec.Uses)
		log.Printf("Inputs: %+v", inputs)
	}

	// Simulate processing time
	select {
	case <-time.After(500 * time.Millisecond):
		// Continue after delay
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Create outputs based on agent type
	outputs := make(map[string]interface{})
	outputs["_agent_id"] = agentSpec.ID
	outputs["_agent_type"] = agentSpec.Uses
	outputs["_processed_at"] = time.Now().Format(time.RFC3339)
	
	// Process inputs based on agent ID
	switch agentSpec.ID {
	case "processor":
		// Process data
		outputs["processed_data"] = "This is processed data"
		outputs["status"] = "completed"
	
	case "analyzer":
		// Analyze data
		outputs["analysis"] = map[string]interface{}{
			"sentiment": "positive",
			"entities": []string{"entity1", "entity2"},
			"confidence": 0.87,
		}
		outputs["status"] = "completed"
	
	case "summarizer":
		// Generate summary
		outputs["summary"] = "This is a summary of the analyzed data"
		outputs["bullet_points"] = []string{
			"Point 1: The analysis shows positive sentiment",
			"Point 2: Multiple entities were detected",
			"Point 3: Confidence level is high (87%)",
		}
		outputs["status"] = "completed"
	
	default:
		// Default behavior
		outputs["result"] = "Generic agent execution completed"
		outputs["status"] = "completed"
	}

	if r.logToConsole {
		log.Printf("Agent execution completed")
		log.Printf("Outputs: %+v", outputs)
	}
	
	return outputs, nil
}

// Cleanup performs any necessary cleanup
func (r *SimpleAgentRuntime) Cleanup() error {
	// No cleanup needed for simple runtime
	return nil
}
