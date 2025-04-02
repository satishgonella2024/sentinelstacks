package main

import (
	"fmt"
	"time"
)

// Simple representation of a DAG node
type Node struct {
	ID          string
	Dependents  []string
	Dependencies []string
}

// Simple representation of a DAG
type DAG struct {
	Nodes map[string]*Node
}

// Simple test function
func testDAG() {
	// Create a simple DAG
	dag := &DAG{
		Nodes: make(map[string]*Node),
	}
	
	// Add nodes
	dag.Nodes["A"] = &Node{ID: "A", Dependents: []string{"B", "C"}, Dependencies: []string{}}
	dag.Nodes["B"] = &Node{ID: "B", Dependents: []string{"D"}, Dependencies: []string{"A"}}
	dag.Nodes["C"] = &Node{ID: "C", Dependents: []string{"D"}, Dependencies: []string{"A"}}
	dag.Nodes["D"] = &Node{ID: "D", Dependents: []string{}, Dependencies: []string{"B", "C"}}
	
	// Find start nodes
	var startNodes []string
	for id, node := range dag.Nodes {
		if len(node.Dependencies) == 0 {
			startNodes = append(startNodes, id)
		}
	}
	
	fmt.Printf("Start nodes: %v\n", startNodes)
	
	// Execute in topological order
	executionOrder, err := topologicalSort(dag)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Execution order: %v\n", executionOrder)
	
	// Simulate execution
	for _, id := range executionOrder {
		fmt.Printf("Executing node %s...\n", id)
		// Simulate processing time
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Node %s completed\n", id)
	}
}

// Simple topological sort implementation
func topologicalSort(dag *DAG) ([]string, error) {
	// Create a map of in-degrees
	inDegree := make(map[string]int)
	for id, node := range dag.Nodes {
		inDegree[id] = len(node.Dependencies)
	}
	
	// Find all nodes with in-degree 0
	var queue []string
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}
	
	// Process the queue
	var result []string
	for len(queue) > 0 {
		// Remove a node from the queue
		node := queue[0]
		queue = queue[1:]
		
		// Add it to the result
		result = append(result, node)
		
		// Decrease in-degree of all its dependents
		for _, depID := range dag.Nodes[node].Dependents {
			inDegree[depID]--
			
			// If in-degree becomes 0, add to queue
			if inDegree[depID] == 0 {
				queue = append(queue, depID)
			}
		}
	}
	
	// Check if we processed all nodes
	if len(result) != len(dag.Nodes) {
		return nil, fmt.Errorf("cycle detected in graph")
	}
	
	return result, nil
}

func main() {
	fmt.Println("Running simple DAG test...")
	testDAG()
	fmt.Println("Test completed")
}
