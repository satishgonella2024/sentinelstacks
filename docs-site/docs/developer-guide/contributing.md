# Contributing to SentinelStacks

Thank you for your interest in contributing to SentinelStacks! This document provides guidelines and instructions for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Git
- Basic knowledge of LLMs and AI agents

### Setting Up the Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/sentinelstacks.git
   cd sentinelstacks
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/satishgonella2024/sentinelstacks.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

## Development Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and write tests as needed

3. Run the tests:
   ```bash
   go test ./...
   ```

4. Build the project:
   ```bash
   go build -o sentinel ./cmd/sentinel
   ```

5. Commit your changes with a descriptive message:
   ```bash
   git commit -m "Feature: Add new capability to the agent runtime"
   ```

6. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a pull request on GitHub

## Project Structure

SentinelStacks follows standard Go project layout:

- `cmd/sentinel`: CLI application entry point
- `pkg/`: Core libraries
  - `agentfile/`: Agentfile parser and schema
  - `models/`: Model adapters
  - `runtime/`: Agent execution runtime
  - `registry/`: Registry client and server
- `docs/`: Documentation
- `examples/`: Example agent definitions

## Coding Guidelines

### Go Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Add comments to exported functions, types, and constants
- Use meaningful variable names

### Documentation

- Update documentation when changing functionality
- Document new features, including examples
- Keep the API documentation up-to-date

### Testing

- Write unit tests for new functionality
- Ensure existing tests pass before submitting a PR
- Include integration tests for complex features

## Pull Request Process

1. Update the documentation with details of your changes
2. Add or update tests as needed
3. Make sure all tests pass
4. Update the README.md if appropriate
5. The PR will be merged once it's reviewed and approved

## Communication

- GitHub Issues: For bug reports and feature requests
- GitHub Discussions: For general questions and discussions
- Discord: For real-time communication (link to be added)

## Project Governance

SentinelStacks is currently maintained by Satish Gonella. The project follows a benevolent dictator model, with input from community contributors.

### Roles

- **Maintainer**: Has commit access and is responsible for reviewing PRs
- **Contributor**: Anyone who contributes code, documentation, or other artifacts
- **User**: Anyone who uses SentinelStacks

## Roadmap and Feature Requests

See the [ROADMAP.md](https://github.com/satishgonella2024/sentinelstacks/blob/main/ROADMAP.md) file for the current development roadmap.

To suggest new features, create an issue using the feature request template.

## License

By contributing to SentinelStacks, you agree that your contributions will be licensed under the project's [MIT License](https://github.com/satishgonella2024/sentinelstacks/blob/main/LICENSE).
