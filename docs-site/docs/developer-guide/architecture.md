# Architecture

This document outlines the architecture and design decisions for SentinelStacks.

## System Components

SentinelStacks is composed of several key components that work together to create a flexible and powerful AI agent platform.

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    CLI / GUI    │     │     Registry    │     │  Model Adapters │
│                 │     │                 │     │                 │
│ User Interfaces │────▶│ Agent Discovery │────▶│ LLM Providers   │
│                 │     │ & Sharing       │     │                 │
└────────┬────────┘     └─────────────────┘     └────────┬────────┘
         │                                               │
         ▼                                               ▼
┌─────────────────┐                           ┌─────────────────┐
│   Agentfile     │                           │ Agent Runtime   │
│                 │                           │                 │
│ Natural Language│◀──────────────────────────│ Execution       │
│ Parser & Schema │                           │ Environment     │
└─────────────────┘                           └─────────────────┘
```

### 1. CLI & GUI Interfaces

- **CLI** (`sentinel`): Command-line interface for developers
- **Desktop App**: Visual interface for monitoring and management
- **Web Interface**: Registry browsing and agent discovery

### 2. Agentfile System

- **Natural Language Parser**: Converts plain English to structured YAML
- **Schema Validation**: Ensures agent definitions follow the standard format
- **Versioning**: Tracks changes to agent definitions

### 3. Agent Runtime

- **Lifecycle Management**: Initializes, runs, and terminates agents
- **State Management**: Persists and retrieves agent memory
- **Resource Management**: Controls memory usage and execution limits
- **Event System**: Provides hooks for monitoring and integration
- **Sandboxing**: Isolates agent execution for security

### 4. Model Adapters

- **Interface Abstraction**: Common API across different model providers
- **Capability Negotiation**: Handles differences in model capabilities
- **Connection Management**: Handles authentication and session management
- **Error Handling**: Provides robust error recovery and fallbacks

### 5. Registry

- **Discovery**: Search and browse available agents
- **Distribution**: Sharing and downloading agent definitions
- **Versioning**: Tracks agent versions and compatibility
- **Access Control**: Manages permissions and visibility

## Data Flow

The diagram below illustrates how data flows through the system when running an agent:

```
┌──────────┐    ┌──────────┐    ┌────────────┐    ┌────────────┐
│  User    │    │ Agentfile│    │  Agent     │    │  Model     │
│ Input    │───▶│ Config   │───▶│  Runtime   │───▶│  Adapter   │
└──────────┘    └──────────┘    └────────────┘    └─────┬──────┘
                                      ▲                  │
                                      │                  ▼
                                      │              ┌────────────┐
                                      │              │   LLM      │
                                      └──────────────│  Provider  │
                                                     └────────────┘
```

1. User provides input through CLI or GUI
2. Agent runtime loads the Agentfile configuration
3. Runtime prepares context (including state and user input)
4. Model adapter translates this into provider-specific format
5. LLM provider processes the request and returns a response
6. Response is processed, state is updated, and results are returned to user

## Key Interfaces

### ModelAdapter Interface

```go
type ModelAdapter interface {
    // Generate sends a prompt to the model and returns the response
    Generate(prompt string, systemPrompt string, options Options) (string, error)
    
    // GetCapabilities returns the capabilities of the model
    GetCapabilities() ModelCapabilities
}
```

### Agent Runtime Interface

```go
type AgentRuntime interface {
    // Initialize sets up the agent with an Agentfile
    Initialize(agentfile string) error
    
    // Run processes user input and returns the agent's response
    Run(input string) (string, error)
    
    // GetState returns the current agent state
    GetState() map[string]interface{}
    
    // SaveState persists the current state
    SaveState() error
}
```

## Design Decisions

### Natural Language First

We've chosen to make natural language the primary way to define agents for several reasons:

1. **Accessibility**: Domain experts can create agents without learning YAML
2. **Expressiveness**: Natural language can capture nuanced requirements
3. **Standardization**: The parser ensures consistent YAML output

### Local-First Execution

SentinelStacks prioritizes local execution where possible:

1. **Privacy**: Sensitive data stays on the user's machine
2. **Cost Control**: Reduces API usage for cloud-based models
3. **Reliability**: Functions without internet access

### Registry Structure

The registry follows a Git-inspired model:

1. **Namespaces**: Username/organization scoping for agents
2. **Semantic Versioning**: Clear compatibility expectations
3. **Discovery Metadata**: Tags, capabilities, and model requirements

### Secure by Default

Security considerations are built in:

1. **Permission System**: Explicit access control for files and network
2. **Sandboxed Execution**: Isolated environments for agent execution
3. **Verification**: Signature checking for registry agents

## Deployment Architecture

SentinelStacks supports multiple deployment models:

### Local Development

```
┌─────────────────┐
│    Developer    │
│                 │
│  sentinel CLI   │◀───┐
└────────┬────────┘    │
         │             │
         ▼             │
┌─────────────────┐    │
│     Ollama      │    │
│                 │    │
│  Local Models   │    │
└─────────────────┘    │
                       │
┌─────────────────┐    │
│   Local File    │    │
│                 │    │
│     Registry    │────┘
└─────────────────┘
```

### Team Collaboration

```
┌─────────────────┐    ┌─────────────────┐
│   Developer A   │    │   Developer B   │
│                 │    │                 │
│  sentinel CLI   │    │  sentinel CLI   │
└────────┬────────┘    └────────┬────────┘
         │                      │
         └──────────┬───────────┘
                    │
                    ▼
          ┌─────────────────┐
          │  Shared Registry│
          │                 │
          │  Git-backed or  │
          │   Database      │
          └─────────────────┘
```

### Enterprise Deployment

```
┌─────────────────┐    ┌─────────────────┐
│    Business     │    │     Data        │
│     Users       │    │    Scientists   │
│                 │    │                 │
│  Desktop GUI    │    │  sentinel CLI   │
└────────┬────────┘    └────────┬────────┘
         │                      │
         └──────────┬───────────┘
                    │
                    ▼
          ┌─────────────────┐
          │  Enterprise     │
          │                 │
          │  Registry       │
          └───────┬─────────┘
                  │
         ┌────────┴─────────┐
         │                  │
         ▼                  ▼
┌─────────────────┐ ┌─────────────────┐
│  Self-hosted    │ │   Cloud-based   │
│                 │ │                 │
│  Models (Ollama)│ │   API Models    │
└─────────────────┘ └─────────────────┘
```

## Technology Stack

- **Backend**: Go for performance and cross-platform support
- **CLI**: Cobra framework for command structure
- **Desktop**: Tauri (Rust + Web Technologies)
- **Data Format**: YAML for configuration, JSON for state
- **Registry Backend**: Git-based (v1), PostgreSQL (future)
- **Documentation**: MkDocs with Material theme
