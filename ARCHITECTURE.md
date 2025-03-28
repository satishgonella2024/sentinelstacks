# SentinelStacks Architecture

This document outlines the architecture and design decisions for SentinelStacks, a platform for creating, running, and sharing AI agents.

## System Components

### 1. CLI (`sentinel`)
- Command-line interface for interacting with all system components
- Implemented using Cobra for command structure
- Provides direct access to agent management, registry, and runtime

### 2. Agentfile Parser
- Converts natural language agent descriptions to structured YAML
- Uses LLMs for understanding and extraction
- Validates and normalizes the generated configuration

### 3. Agent Runtime
- Executes agents based on their Agentfile configuration
- Manages agent lifecycle (init, run, pause, terminate)
- Handles state persistence and retrieval

### 4. Model Adapters
- Provides unified interface to different LLM backends
- Handles model-specific quirks and capabilities
- Supports local models (Ollama) and API-based services (OpenAI, Claude)

### 5. Registry
- Stores and manages agent definitions
- Enables discovery and sharing
- Tracks versions and dependencies

### 6. Desktop Application
- Visual interface for SentinelStacks
- Agent monitoring and management
- Built with Tauri and React

## Data Flow

1. User defines agent in natural language
2. Parser converts to structured YAML
3. Agent runtime initializes based on configuration
4. Runtime connects to appropriate model through adapter
5. State is maintained and persisted as needed
6. Results are returned to user via CLI or Desktop UI

## Key Interfaces

### Agentfile Format

```yaml
name: agent-name
version: "1.0.0"
description: "Agent description"
model:
  provider: "ollama"
  name: "llama3"
capabilities:
  - capability1
  - capability2
memory:
  type: "simple"
  persistence: true
```

### Model Adapter Interface

```go
type ModelAdapter interface {
    Generate(prompt string, systemPrompt string, options Options) (string, error)
    GetCapabilities() ModelCapabilities
}
```

### Agent Runtime Interface

```go
type AgentRuntime interface {
    Initialize(agentfile string) error
    Run(input string) (string, error)
    GetState() map[string]interface{}
    SaveState() error
}
```

## Deployment Architecture

Phase 1: Local-first with file-based storage
Phase 2: Optional remote registry with authentication
Phase 3: Enterprise deployment with private registries
