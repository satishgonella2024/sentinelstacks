package stack

import (
	"testing"
)

// TestDAGCreation tests the creation of a DAG from a stack specification
func TestDAGCreation(t *testing.T) {
	// Create a sample stack spec
	spec := StackSpec{
		Name:        "test-stack",
		Description: "Test stack for DAG",
		Version:     "1.0.0",
		Agents: []StackAgentSpec{
			{
				ID:   "agent1",
				Uses: "processor",
			},
			{
				ID:        "agent2",
				Uses:      "analyzer",
				InputFrom: []string{"agent1"},
			},
			{
				ID:        "agent3",
				Uses:      "summarizer",
				InputFrom: []string{"agent2"},
			},
		},
	}

	// Create DAG
	dag, err := NewDAG(spec)
	if err != nil {
		t.Fatalf("Error creating DAG: %v", err)
	}

	// Check that all nodes are present
	if len(dag.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(dag.Nodes))
	}

	// Check that the start node is correct
	if len(dag.StartNodes) != 1 || dag.StartNodes[0].ID != "agent1" {
		t.Errorf("Incorrect start node(s)")
	}

	// Get execution order
	order, err := dag.TopologicalSort()
	if err != nil {
		t.Fatalf("Error sorting DAG: %v", err)
	}

	// Check execution order
	expectedOrder := []string{"agent1", "agent2", "agent3"}
	if len(order) != len(expectedOrder) {
		t.Fatalf("Expected %d items in execution order, got %d", len(expectedOrder), len(order))
	}

	for i, id := range order {
		if id != expectedOrder[i] {
			t.Errorf("Expected %s at position %d, got %s", expectedOrder[i], i, id)
		}
	}
}

// TestCycleDetection tests the cycle detection in the DAG
func TestCycleDetection(t *testing.T) {
	// Create a stack spec with a cycle
	spec := StackSpec{
		Name:        "cyclic-stack",
		Description: "Stack with a cycle for testing",
		Version:     "1.0.0",
		Agents: []StackAgentSpec{
			{
				ID:   "agent1",
				Uses: "processor",
			},
			{
				ID:        "agent2",
				Uses:      "analyzer",
				InputFrom: []string{"agent1"},
			},
			{
				ID:        "agent3",
				Uses:      "summarizer",
				InputFrom: []string{"agent2"},
			},
			{
				ID:        "agent1", // Duplicate ID creates a cycle
				Uses:      "duplicator",
				InputFrom: []string{"agent3"},
			},
		},
	}

	// Create DAG - should fail due to duplicate ID
	_, err := NewDAG(spec)
	if err == nil {
		t.Fatalf("Expected error due to duplicate ID, but got none")
	}

	// Create a stack spec with an explicit cycle
	spec = StackSpec{
		Name:        "cyclic-stack",
		Description: "Stack with a cycle for testing",
		Version:     "1.0.0",
		Agents: []StackAgentSpec{
			{
				ID:   "agent1",
				Uses: "processor",
			},
			{
				ID:        "agent2",
				Uses:      "analyzer",
				InputFrom: []string{"agent1", "agent3"}, // Cyclic dependency
			},
			{
				ID:        "agent3",
				Uses:      "summarizer",
				InputFrom: []string{"agent2"},
			},
		},
	}

	// Create DAG - should fail due to cycle
	_, err = NewDAG(spec)
	if err == nil {
		t.Fatalf("Expected error due to cycle, but got none")
	}
}
