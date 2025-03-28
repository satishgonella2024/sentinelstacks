# SentinelStacks Roadmap

This document outlines the development roadmap for SentinelStacks, showing our current progress and future plans.

## Current Status

SentinelStacks is currently in early development. The core CLI and Agentfile parser are being developed, with initial support for Ollama model integration.

## Development Phases

### Phase 1: Core CLI & Runtime (Current)

- [x] Project initialization
- [x] Basic CLI structure
- [ ] Agentfile parser (Natural Language → YAML)
- [ ] Simple agent runtime with Ollama adapter
- [ ] Local file-based registry

**Estimated Completion:** Q2 2025

### Phase 2: Desktop UI & Enhanced Features

- [ ] Basic Tauri desktop application
- [ ] Agent management screens
- [ ] Execution monitoring UI
- [ ] Settings management
- [ ] Enhanced state management
- [ ] Multi-model support (OpenAI, Claude)

**Estimated Completion:** Q3 2025

### Phase 3: Registry & Distribution

- [ ] Registry server basics
- [ ] Authentication and user management
- [ ] Publishing workflow
- [ ] Search and discovery features
- [ ] Versioning and dependency management
- [ ] Public registry hosting

**Estimated Completion:** Q4 2025

### Phase 4: Advanced Features

- [ ] Multi-agent orchestration
- [ ] Agent communication protocols
- [ ] Enhanced state management
- [ ] Analytics dashboard
- [ ] Plugin system for extensibility
- [ ] Enterprise features (SSO, audit logging)

**Estimated Completion:** Q1 2026

## Feature Backlog

These features are on our radar but not yet scheduled:

### Agent Capabilities
- [ ] Visual agent builder
- [ ] Conversation memory optimization
- [ ] Context window management
- [ ] Tool integration framework
- [ ] Custom prompt templates
- [ ] Agent performance metrics

### Runtime Enhancements
- [ ] Parallel agent execution
- [ ] Distributed execution
- [ ] GPU acceleration support
- [ ] Webhook triggers and automation
- [ ] Stream processing capabilities

### Developer Experience
- [ ] VS Code extension
- [ ] Agent debugging tools
- [ ] Performance profiling
- [ ] Testing framework for agents
- [ ] CI/CD integration

### Enterprise Features
- [ ] Role-based access control
- [ ] Compliance & audit tooling
- [ ] Private registry hosting
- [ ] Enterprise SSO
- [ ] Usage monitoring & quotas

## Release Schedule

| Version | Focus | Target Date |
|---------|-------|-------------|
| v0.1.0 | Basic CLI & Agentfile | May 2025 |
| v0.2.0 | Agent Runtime & Ollama | June 2025 |
| v0.3.0 | Registry Basics | July 2025 |
| v0.4.0 | Desktop UI Preview | August 2025 |
| v0.5.0 | Multi-Model Support | September 2025 |
| v1.0.0 | First Stable Release | December 2025 |

## How to Contribute

We welcome contributions to help us achieve these roadmap items! See our [contributing guide](developer-guide/contributing.md) for how to get involved.

## Feedback

Have suggestions for the roadmap? Please open an issue on GitHub with the tag 'roadmap' to share your thoughts.
