# Core DAG Implementation

This directory contains the core Directed Acyclic Graph (DAG) implementation for the SentinelStacks project.

## Files

- `dag.go`: Implementation of the DAG with topological sorting and cycle detection
- `types.go`: Core type definitions for the stack engine

## How It Works

The DAG implementation is used to determine the execution order of agents in a stack. It:

1. Creates nodes for each agent in the stack
2. Establishes dependencies between nodes based on data flow
3. Detects cycles to ensure the graph is acyclic
4. Performs topological sorting to determine execution order

This ensures that agents are executed in the correct order, with dependencies executed before the agents that depend on them.

## Testing

The DAG implementation has been thoroughly tested using:

1. Simple test cases with linear dependencies
2. Complex test cases with branching and merging flows
3. Edge cases including cycle detection

All tests pass, confirming that the core implementation works correctly.
