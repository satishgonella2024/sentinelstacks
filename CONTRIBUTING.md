# Contributing to SentinelStacks

Thank you for your interest in contributing to SentinelStacks! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Pull Request Process](#pull-request-process)
5. [Coding Standards](#coding-standards)
6. [Testing](#testing)
7. [Documentation](#documentation)
8. [Agent Development](#agent-development)
9. [Tools Development](#tools-development)
10. [Community](#community)

## Code of Conduct

By participating in this project, you agree to uphold our Code of Conduct:

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Go 1.20 or later
- Docker and Docker Compose (for running services)
- Node.js and npm (for UI development)
- Ollama (for local model execution)

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your forked repository:
   ```bash
   git clone https://github.com/yourusername/sentinelstacks.git
   cd sentinelstacks
   ```

3. Add the original repository as a remote:
   ```bash
   git remote add upstream https://github.com/satishgonella2024/sentinelstacks.git
   ```

4. Install dependencies:
   ```bash
   go mod tidy
   ```

5. Build the project:
   ```bash
   ./scripts/build.sh
   ```

## Development Workflow

1. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes, following our coding standards

3. Run tests to ensure your changes don't break existing functionality:
   ```bash
   go test ./...
   ```

4. Commit your changes with a descriptive commit message:
   ```bash
   git commit -m "Add detailed description of your changes"
   ```

5. Keep your branch updated with the upstream:
   ```bash
   git pull upstream main
   ```

6. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a Pull Request (PR) from your branch to the upstream main branch

## Pull Request Process

1. Ensure your PR includes only relevant changes
2. Update the README.md or documentation with details of changes, if applicable
3. Verify that your code passes all tests
4. Include new tests for new functionality
5. The PR must receive approval from at least one maintainer
6. Once approved, a maintainer will merge your PR

## Coding Standards

We follow standard Go coding conventions:

- Use `gofmt` to format your code
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Document all exported functions, types, and packages
- Write meaningful error messages
- Use descriptive variable and function names

## Testing

- All new features should include appropriate tests
- Aim for at least 80% test coverage for new code
- Test both normal operation and error conditions
- Use the standard Go testing package for unit tests
- Place integration tests in a separate directory

Example test:

```go
func TestAgentCreation(t *testing.T) {
    config := agent.AgentConfig{
        Name:        "test-agent",
        Version:     "1.0.0",
        Description: "Test agent",
        // ... other fields
    }
    
    ag, err := agent.NewAgent(config)
    if err != nil {
        t.Fatalf("Failed to create agent: %v", err)
    }
    
    if ag.Name != "test-agent" {
        t.Errorf("Expected name %q, got %q", "test-agent", ag.Name)
    }
    
    // ... other assertions
}
```

## Documentation

- Document all public APIs
- Keep README.md up to date
- Update ARCHITECTURE.md when changing the system design
- Add new commands to the help text
- Document new features in the docs directory

## Agent Development

### Creating a New Agent

1. Design your agent's purpose and capabilities
2. Create a directory in `examples/agents/` or `agents/`
3. Create an Agentfile (YAML) with the agent configuration
4. Implement the agent's logic
5. Test your agent thoroughly
6. Document how to use your agent

### Agent Structure

```
agent-name/
├── Agentfile           # Agent configuration
├── agent.py            # Agent implementation (Python)
├── agent.go            # Agent implementation (Go)
├── README.md           # Agent documentation
└── examples/           # Example usages
```

## Tools Development

### Creating a New Tool

1. Identify a useful capability for agents
2. Design a clear interface for the tool
3. Create a new package in `pkg/tools/`
4. Implement the tool functionality
5. Add comprehensive tests
6. Document the tool API and usage
7. Create examples demonstrating the tool

### Tool Interface

Tools should implement the `Tool` interface:

```go
type Tool interface {
    // Execute runs the tool with the given parameters and returns the result
    Execute(params map[string]interface{}) (interface{}, error)
    
    // GetCapabilities returns information about what the tool can do
    GetCapabilities() ToolCapabilities
}
```

## Community

- Join our [Discord server](https://discord.gg/sentinelstacks) to chat with other contributors
- Subscribe to our [mailing list](https://sentinelstacks.io/mailing-list) for announcements
- Follow us on [Twitter](https://twitter.com/sentinelstacks) for updates

## Recognition

Contributors who make significant contributions will be recognized in our README.md and on our website. We appreciate and value all contributions, big and small!

Thank you for helping make SentinelStacks better!