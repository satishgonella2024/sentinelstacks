package stack

import (
	"errors"
	"fmt"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// Node represents a node in the DAG
type Node struct {
	ID           string
	Dependencies []*Node
	Dependents   []*Node
	State        string
}

// DAG represents a directed acyclic graph for stack execution
type DAG struct {
	Nodes      map[string]*Node
	StartNodes []*Node
}

// NewDAG creates a new DAG from a stack specification
func NewDAG(spec types.StackSpec) (*DAG, error) {
	if len(spec.Agents) == 0 {
		return nil, fmt.Errorf("stack must contain at least one agent")
	}

	dag := &DAG{
		Nodes: make(map[string]*Node),
	}

	// Create nodes for all agents
	for _, agentSpec := range spec.Agents {
		if agentSpec.ID == "" {
			return nil, fmt.Errorf("agent ID cannot be empty")
		}

		// Check if node with this ID already exists
		if _, exists := dag.Nodes[agentSpec.ID]; exists {
			return nil, fmt.Errorf("duplicate agent ID detected: %s", agentSpec.ID)
		}

		// Create node
		node := &Node{
			ID:    agentSpec.ID,
			State: "pending",
		}

		dag.Nodes[agentSpec.ID] = node
	}

	// Connect nodes based on dependencies
	for _, agentSpec := range spec.Agents {
		node := dag.Nodes[agentSpec.ID]

		// Add dependencies based on explicit 'depends' field
		for _, depID := range agentSpec.Depends {
			depNode, exists := dag.Nodes[depID]
			if !exists {
				return nil, fmt.Errorf("agent %s depends on non-existent agent %s", agentSpec.ID, depID)
			}
			node.Dependencies = append(node.Dependencies, depNode)
			depNode.Dependents = append(depNode.Dependents, node)
		}

		// Add dependencies based on inputFrom field
		for _, inputFrom := range agentSpec.InputFrom {
			if inputFrom == "" {
				continue
			}

			inputNode, exists := dag.Nodes[inputFrom]
			if !exists {
				return nil, fmt.Errorf("agent %s takes input from non-existent agent %s", agentSpec.ID, inputFrom)
			}

			// Only add if not already in dependencies
			alreadyExists := false
			for _, dep := range node.Dependencies {
				if dep.ID == inputFrom {
					alreadyExists = true
					break
				}
			}

			if !alreadyExists {
				node.Dependencies = append(node.Dependencies, inputNode)
				inputNode.Dependents = append(inputNode.Dependents, node)
			}
		}
	}

	// Detect cycles in the graph
	if hasCycle(dag) {
		return nil, errors.New("cycle detected in agent dependencies - DAG must be acyclic")
	}

	// Find start nodes (nodes with no dependencies)
	for _, node := range dag.Nodes {
		if len(node.Dependencies) == 0 {
			dag.StartNodes = append(dag.StartNodes, node)
		}
	}

	if len(dag.StartNodes) == 0 {
		return nil, errors.New("no start nodes found - every agent has dependencies")
	}

	return dag, nil
}

// hasCycle checks if the DAG contains a cycle
func hasCycle(dag *DAG) bool {
	// Create a map to track visited nodes
	visited := make(map[string]bool)
	inProgress := make(map[string]bool)

	// Check each node
	for id := range dag.Nodes {
		if !visited[id] {
			if dfs(dag.Nodes[id], visited, inProgress) {
				return true
			}
		}
	}

	return false
}

// dfs performs depth-first search to detect cycles
func dfs(node *Node, visited, inProgress map[string]bool) bool {
	visited[node.ID] = true
	inProgress[node.ID] = true

	// Check all dependents
	for _, dep := range node.Dependents {
		if !visited[dep.ID] {
			if dfs(dep, visited, inProgress) {
				return true
			}
		} else if inProgress[dep.ID] {
			// If the dependent is already in progress, we have a cycle
			return true
		}
	}

	inProgress[node.ID] = false
	return false
}

// TopologicalSort returns a valid execution order for the DAG
func (d *DAG) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	tempMark := make(map[string]bool)
	order := []string{}

	var visit func(node *Node) error
	visit = func(node *Node) error {
		if tempMark[node.ID] {
			return errors.New("cycle detected in DAG")
		}

		if !visited[node.ID] {
			tempMark[node.ID] = true

			// Visit all dependents
			for _, dep := range node.Dependencies {
				if err := visit(dep); err != nil {
					return err
				}
			}

			visited[node.ID] = true
			tempMark[node.ID] = false
			order = append(order, node.ID) // Add to order
		}

		return nil
	}

	// Visit all nodes
	for _, node := range d.StartNodes {
		if err := visit(node); err != nil {
			return nil, err
		}
	}

	// Visit any remaining nodes (shouldn't happen if DAG is connected)
	for id, node := range d.Nodes {
		if !visited[id] {
			if err := visit(node); err != nil {
				return nil, err
			}
		}
	}

	return order, nil
}
