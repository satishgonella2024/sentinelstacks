# Registry Guide

The SentinelStacks Registry allows you to share, discover, and reuse agents created by the community. This guide explains how to use the registry features.

## Registry Overview

The registry works similarly to package registries like npm, Docker Hub, or the Terraform Registry. It provides:

- A central repository for agents
- Discovery through search and browsing
- Version management
- Access control (public/private agents)

## Pushing an Agent

To share your agent with others, use the `registry push` command:

```bash
sentinel registry push my-agent
```

By default, agents are pushed as public, meaning anyone can discover and pull them. To push a private agent:

```bash
sentinel registry push my-agent --visibility private
```

### What Gets Pushed

When you push an agent, the following files are included:
- `agentfile.yaml`: The structured configuration
- `agentfile.natural.txt`: The natural language description
- Any additional files in the agent directory (except state files)

State files (`.state.json`) are not pushed to the registry to avoid sharing user-specific conversation history.

## Pulling an Agent

To use an agent someone else has created:

```bash
sentinel registry pull username/agent-name
```

This downloads the agent to your local machine with all its files. You can optionally specify a version:

```bash
sentinel registry pull username/agent-name@1.0.0
```

If no version is specified, the latest version is pulled.

## Searching the Registry

To find agents in the registry:

```bash
sentinel registry search "keyword"
```

You can also filter by tags:

```bash
sentinel registry search --tags research,summarization
```

## Listing Agents

To see all available agents in the registry:

```bash
sentinel registry list
```

This shows all public agents and your private agents.

## Registry Architecture

SentinelStacks uses a hybrid registry approach:

1. **Local Registry**: File-based storage in `~/.sentinelstacks/registry`
2. **Remote Registry**: (Planned) Cloud-based registry for wider sharing

The current implementation focuses on the local registry, with remote registry capabilities planned for future releases.

## Agent Naming Convention

Agents in the registry follow this naming convention:

```
username/agent-name[@version]
```

For example:
- `janesmith/research-assistant`: Latest version of janesmith's research assistant
- `janesmith/research-assistant@1.2.0`: Specific version of the agent

## Registry Metadata

Each agent in the registry includes metadata:
- Author name
- Version history
- Download count
- Tags and capabilities
- Compatible models

This metadata helps users discover and choose the right agents for their needs.

## Future Registry Features

Planned enhancements to the registry system include:
- Remote registry with user accounts
- Rating and review system
- Enhanced version management
- Dependencies between agents
- Team-based access control
