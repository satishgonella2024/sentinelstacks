# SentinelStacks Registry

The SentinelStacks Registry is a centralized repository for sharing and discovering agents and stacks. This document explains how to use the registry features to collaborate with others.

## Overview

The registry system provides:

- **Stack Sharing**: Share your stack configurations with others
- **Agent Distribution**: Publish agents for others to use in their stacks
- **Versioning**: Maintain multiple versions of your stacks and agents
- **Security**: Sign and verify packages to ensure authenticity
- **Discovery**: Search and browse available stacks and agents

## File Formats

SentinelStacks uses standardized file formats with unique extensions:

- **Agent Packages**: `.agent.sntl` - Packaged agent ready for distribution
- **Stack Packages**: `.stack.sntl` - Packaged stack ready for distribution
- **Agent Definitions**: `.agent.yaml` - Agent definition file (similar to Dockerfile)
- **Stack Definitions**: `.stack.yaml` - Stack definition file (formerly Stackfile.yaml)
- **Signatures**: `.sig.sntl` - Detached signature for verification

## Package Signing and Verification

All packages in the registry can be cryptographically signed to verify their authenticity and integrity.

### Generating Signing Keys

```bash
# Generate a new key pair
sentinel key generate --id developer-key

# List your keys
sentinel key list
```

### Signing Packages

```bash
# Sign a package during push
sentinel stack push my-stack.stack.yaml --sign --key developer-key

# Create a detached signature
sentinel sign my-stack.stack.sntl --key developer-key
```

### Verifying Packages

```bash
# Verify a package during pull (default)
sentinel stack pull my-stack:1.0.0 --verify

# Verify a package explicitly
sentinel verify my-stack-1.0.0.stack.sntl
```

## Using the Registry

### Authentication

```bash
# Log in to the registry
sentinel login

# Log out
sentinel logout

# View current authentication status
sentinel whoami
```

### Pushing Stacks

```bash
# Push a stack to the registry
sentinel stack push my-stack.stack.yaml

# Push with a specific author
sentinel stack push my-stack.stack.yaml --author "Your Name <email@example.com>"

# Push and build a package file
sentinel stack push my-stack.stack.yaml --build --output my-stack.stack.sntl
```

### Pulling Stacks

```bash
# Pull the latest version of a stack
sentinel stack pull my-stack

# Pull a specific version
sentinel stack pull my-stack:1.2.3

# Pull to a specific directory
sentinel stack pull my-stack --output ./my-project

# Pull with dependency resolution
sentinel stack pull my-stack --extract-agents
```

### Searching for Stacks

```bash
# Search for stacks
sentinel stack search keywords

# Limit results
sentinel stack search keywords --limit 5

# Format results
sentinel stack search keywords --format wide
```

## Stack Dependencies

When you push a stack to the registry, the system automatically analyzes it for dependencies on specific agents. When others pull your stack, they can choose to automatically pull the required agents as well.

```bash
# Check dependencies before pulling
sentinel stack inspect-remote my-stack:1.0.0

# Pull with dependencies
sentinel stack pull my-stack:1.0.0 --extract-agents
```

## Private Registries

You can configure SentinelStacks to use private registries for your organization:

```bash
# Configure a private registry
sentinel config set registry.url https://registry.your-company.com

# Log in to the private registry
sentinel login --registry https://registry.your-company.com
```

## Custom Metadata and Labels

You can add custom metadata and labels to your packages for better organization:

```bash
# Add labels during push
sentinel stack push my-stack.stack.yaml --label category=NLP --label team=DataScience

# Search by label
sentinel stack search --label team=DataScience
```

## Registry API

The registry provides a REST API that can be used programmatically. Documentation for the API can be found at `https://registry.sentinelstacks.io/api/docs`.

## Configuration

The registry client configuration is stored in your SentinelStacks config file:

```yaml
registry:
  url: https://registry.sentinelstacks.io
  auth_token: your-auth-token
  cache_dir: ~/.sentinel/cache

security:
  keys_dir: ~/.sentinel/keys
  default_key: developer-key
```

You can modify these settings using the `sentinel config` command.

## Troubleshooting

### Common Issues

1. **Authentication Failures**: Make sure you are logged in with `sentinel login`
2. **Push Failures**: Ensure your stack file is valid and all dependencies are available
3. **Signature Verification**: Import the required public keys with `sentinel key import`

### Registry Status

Check the registry status at `https://status.sentinelstacks.io` for any service disruptions.
