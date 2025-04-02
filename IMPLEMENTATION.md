# SentinelStacks Implementation Plan

This document outlines the detailed implementation plan for SentinelStacks, providing specific technical tasks and architectural decisions for each phase of development.

## Phase 1: Developer-Complete Core System

### Week 1: Core Agent Flow & Stack Engine Foundations

#### Task 1.1: Rename "compose" to "stack" throughout codebase
- Create new `cmd/sentinel/stack` directory
- Migrate functionality from `cmd/sentinel/compose`
- Update CLI command structure and help text
- Maintain backward compatibility with temporary aliases
- Update all relevant documentation

#### Task 1.2: Implement Stack Specification Structure
```go
// internal/stack/types.go
type StackSpec struct {
    Name        string                 `json:"name" yaml:"name"`
    Description string                 `json:"description" yaml:"description"`
    Version     string                 `json:"version" yaml:"version"`
    Agents      []StackAgentSpec       `json:"agents" yaml:"agents"`
    Networks    []string               `json:"networks" yaml:"networks"`
    Volumes     []string               `json:"volumes" yaml:"volumes"`
    Metadata    map[string]interface{} `json:"metadata" yaml:"metadata"`
}

type StackAgentSpec struct {
    ID         string                 `json:"id" yaml:"id"`
    Uses       string                 `json:"uses" yaml:"uses"`
    InputFrom  []string               `json:"inputFrom" yaml:"inputFrom"`
    InputKey   string                 `json:"inputKey" yaml:"inputKey"`
    OutputKey  string                 `json:"outputKey" yaml:"outputKey"`
    Params     map[string]interface{} `json:"params" yaml:"params"`
    Depends    []string               `json:"depends" yaml:"depends"`
}
```

#### Task 1.3: Enhance NLP-to-Sentinelfile Parser
- Extend parser to recognize multi-agent patterns
- Implement natural language parsing for agent relationships
- Add capability extraction from descriptions
- Create parser test suite with example inputs

### Week 2: DAG Runner & Context Management

#### Task 2.1: Implement DAG Runner
```go
// internal/stack/engine.go
type StackEngine struct {
    spec      StackSpec
    agents    map[string]*AgentInstance
    ctx       context.Context
    stateManager StateManager
}

func NewStackEngine(spec StackSpec) *StackEngine {
    // Initialize stack engine with the given specification
}

func (e *StackEngine) BuildExecutionGraph() (*DAG, error) {
    // Convert StackSpec to a directed acyclic graph
}

func (e *StackEngine) Execute(ctx context.Context) error {
    // Execute the DAG in topological order
}
```

#### Task 2.2: Implement Context Propagation
- Create context management system for passing data between agents
- Implement input/output mapping between agents
- Add JSON path support for selecting specific data fields
- Create state persistence between agent executions

#### Task 2.3: CLI Command Implementation
- Implement `sentinel stack run` command
- Add options for verbose output and debug logging
- Support for providing inputs to the first agent in the stack

### Week 3: Multi-Agent Execution & Testing

#### Task 3.1: Complete Agent State Manager
```go
// internal/runtime/state.go
type StateManager interface {
    Get(agentID string, key string) (interface{}, error)
    Set(agentID string, key string, value interface{}) error
    GetAll(agentID string) (map[string]interface{}, error)
    Clear(agentID string) error
}
```

#### Task 3.2: REST API Extensions
- Add endpoints for stack management
- Implement WebSocket for real-time stack execution monitoring
- Create API documentation using OpenAPI/Swagger

#### Task 3.3: Testing
- Create integration tests for multi-agent stacks
- Test different DAG patterns (linear, branching, rejoining)
- Benchmark performance with different stack configurations

## Phase 2: Team Workflow & Collaboration Layer

### Week 4: Registry & Authentication

#### Task 4.1: Complete Registry API
```go
// internal/registry/api.go
type RegistryAPI interface {
    Login(credentials Credentials) (AuthToken, error)
    Logout() error
    Push(agentImage *agent.Image, tag string) error
    Pull(name string, tag string) (*agent.Image, error)
    Search(query string) ([]SearchResult, error)
    Tags(name string) ([]string, error)
}
```

#### Task 4.2: Authentication System
- Implement JWT-based authentication
- Add user management functionality
- Create role-based access control for organizations

#### Task 4.3: Agent Signature
- Implement cryptographic signing of agent definitions
- Add verification of signatures during pull operations
- Create key management utilities

### Week 5: Memory & Observability

#### Task 5.1: Memory Plugin Interface
```go
// internal/memory/plugin.go
type MemoryStore interface {
    Save(key string, data interface{}) error
    Load(key string) (interface{}, error)
    Query(embedding []float32, topK int) ([]MemoryMatch, error)
    Delete(key string) error
    Clear() error
}
```

#### Task 5.2: Memory Implementations
- Create in-memory storage implementation
- Add SQLite persistence layer
- Implement Chroma vector store integration

#### Task 5.3: Observability
- Add run ID for stack executions
- Implement structured logging
- Create metrics collection for agent performance
- Add trace context propagation

## Phase 3: UX & Product Readiness

### Week 6: UI Core Components

#### Task 6.1: DAG Visualizer
- Create React component for displaying agent graphs
- Implement interactive visualization
- Add ability to edit connections visually

#### Task 6.2: Stack Launcher
- Create form-based UI for configuring stack runs
- Add parameter validation
- Implement file upload for inputs

#### Task 6.3: Documentation Structure
- Update README with new features
- Create comprehensive setup guide
- Begin writing tutorials

### Week 7: UI Enhancements & Real-time Features

#### Task 7.1: WebSocket Implementation
- Create WebSocket server for real-time updates
- Implement client-side WebSocket consumer
- Add reconnection handling

#### Task 7.2: Model Selector
- Create UI for selecting LLM models
- Add model parameter configuration
- Implement preview functionality

#### Task 7.3: User Onboarding
- Create first-time user experience
- Add interactive tutorials
- Implement sample agent library

### Week 8: Final Polish & Quality Assurance

#### Task 8.1: Template Gallery
- Create template system for common agents
- Add categorization and search
- Implement one-click creation

#### Task 8.2: Documentation Completion
- Finalize all documentation
- Create video tutorials
- Complete API reference

#### Task 8.3: Final Testing
- End-to-end testing of all workflows
- Performance benchmarking
- Security review
- User acceptance testing

## Technical Architecture Details

### Stack Engine Architecture

The Stack Engine will be the core component for executing multi-agent workflows:

```
internal/stack/
├── engine.go       # Main execution engine
├── dag.go          # Directed acyclic graph implementation
├── types.go        # Core data structures
├── parser.go       # Stack file parser
├── validator.go    # Validation logic
└── context.go      # Context propagation
```

### Memory Management Architecture

The memory subsystem will provide persistence and vector search capabilities:

```
internal/memory/
├── plugin.go       # Interface definition
├── local.go        # In-memory implementation
├── sqlite.go       # SQLite persistence
├── chroma.go       # Chroma vector store integration
└── factory.go      # Factory for creating memory instances
```

### Registry System Architecture

The registry system will enable sharing and collaboration:

```
internal/registry/
├── api.go          # Registry interface
├── client.go       # HTTP client implementation
├── auth.go         # Authentication
├── models.go       # Data models
└── cache.go        # Local cache
```

## Integration Tests

To ensure the system works as expected, we'll create integration tests for key workflows:

1. **Simple Chain** - Linear execution of 3 agents
2. **Branching Workflow** - One agent feeding into multiple parallel agents
3. **Converging Workflow** - Multiple agents feeding into a single agent
4. **Complex DAG** - Combination of branching and converging patterns
5. **Error Handling** - Testing recovery from agent failures

## Deployment Strategy

For both self-hosted and managed deployments:

1. **Binary Distribution** - Single binary for CLI usage
2. **Docker Images** - Containerized deployment for server components
3. **Helm Chart** - Kubernetes deployment option

## Post-Release Planning

After completing Phase 3:

1. **Community Engagement** - Open source community building
2. **Extension Ecosystem** - Plugin/extension marketplace
3. **Enterprise Features** - Compliance, audit logging, SSO
4. **Managed Service** - Hosted offering

## Technical Debt Management

Throughout development, we'll actively manage technical debt:

1. Weekly code review sessions
2. Maintain test coverage above 70%
3. Regular dependency updates
4. Architectural decision records (ADRs) for major decisions
5. Refactoring sessions for problem areas
