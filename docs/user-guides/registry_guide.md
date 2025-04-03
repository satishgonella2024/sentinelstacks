# Registry Management

This guide covers everything you need to know about the SentinelStacks registry system, which allows you to share agents and stacks across machines and with other users.

## Overview

The SentinelStacks registry is a central repository for storing, sharing, and discovering agent and stack packages. Similar to Docker Hub or npm, it provides:

- Version management for agents and stacks
- Searchable catalogs of available packages
- Authentication and access controls
- Cryptographic verification of package integrity

## Registry Concepts

### Package Types

The registry supports two primary package types:

1. **Agent Packages** (`.agent.sntl`): Contain a single agent with its configuration and dependencies
2. **Stack Packages** (`.stack.sntl`): Contain a multi-agent system definition and orchestration rules

### Custom File Formats

SentinelStacks uses custom file formats with specific extensions:

| Extension | Description | MIME Type |
|-----------|-------------|-----------|
| `.agent.sntl` | Agent package | `application/x-sentinel-agent` |
| `.stack.sntl` | Stack package | `application/x-sentinel-stack` |
| `.agent.yaml` | Agent definition | `application/x-sentinel-agent-def+yaml` |
| `.stack.yaml` | Stack definition | `application/x-sentinel-stack-def+yaml` |
| `.sig.sntl` | Detached signature | `application/x-sentinel-signature` |

### Package Format

The `.sntl` package format has several key characteristics:

1. **Magic Headers**: Each package starts with a magic header identifying its type:
   - Agent packages: `SNTL-AGENT-PKG`
   - Stack packages: `SNTL-STACK-PKG`

2. **Version Field**: A 4-byte version field that identifies the package format version

3. **Compressed Archive**: The remainder of the file is a gzipped tar archive containing:
   - `sentinel.manifest.json`: Package metadata and file inventory
   - The main definition file (`.agent.yaml` or `.stack.yaml`)
   - Documentation and example files
   - Other supporting resources

4. **Signature Block**: Optional cryptographic signatures for verification

## Authentication

Before you can push packages to the registry, you need to authenticate. The registry supports:

1. **Username/Password Authentication**:
   ```bash
   sentinel login --username myuser --password mypassword
   ```

2. **Token Authentication**:
   ```bash
   sentinel login --token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```

If you don't provide credentials via flags, the command will interactively prompt you for them.

## Configuration

The registry connection can be configured in your SentinelStacks config file, typically located at `~/.sentinel/config.yaml`:

```yaml
registry:
  default: registry.example.com
  auth:
    registry_example_com:
      token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
      username: myuser
```

## Registry Commands

### Login and Authentication

```bash
# Log in to the default registry
sentinel login

# Log in to a specific registry
sentinel login registry.mycompany.com

# Log in with a username (password will be prompted)
sentinel login --username myuser

# Log out from the default registry
sentinel logout

# Log out from all registries
sentinel logout --all
```

### Pushing and Pulling Agents

```bash
# Push an agent with the latest tag
sentinel push myname/research-agent

# Push a specific version
sentinel push myname/research-agent:v1.2

# Push to a specific registry
sentinel push --registry registry.mycompany.com myname/research-agent

# Make an agent publicly accessible
sentinel push --public myname/research-agent:latest

# Pull an agent
sentinel pull myname/research-agent:v1.2

# Pull from a specific registry
sentinel pull --registry registry.mycompany.com myname/research-agent
```

### Working with Stacks

```bash
# Push a stack to the registry
sentinel stack push Stackfile.yaml --sign

# Pull a stack from the registry
sentinel stack pull username/text-analyzer:latest --extract-agents

# Search for stacks
sentinel stack search text-analyzer

# List available versions of a stack
sentinel stack versions text-analyzer
```

## Package Naming

Agent and stack packages follow a similar naming convention to Docker images:

```
[registry/][namespace/]name[:tag]
```

Where:
- `registry` is the optional registry hostname
- `namespace` is the optional user or organization namespace
- `name` is the required package name
- `tag` is the optional tag (defaults to `latest`)

## Security Features

### Package Signing

You can cryptographically sign packages when pushing them to the registry:

```bash
# Sign a package while pushing
sentinel push --sign myname/research-agent:v1.2

# Verify a package signature
sentinel verify myname/research-agent:v1.2
```

### Access Controls

The registry supports different access levels:

- **Read-only**: Can download public packages
- **Publisher**: Can upload packages under their own namespace
- **Admin**: Has full management access

## Registry API

The registry provides a REST API for programmatic integration:

- `GET /v1/packages/search` - Search for packages
- `GET /v1/packages/{name}/{version}` - Download a package
- `PUT /v1/packages/publish` - Upload a new package
- `GET /v1/packages/{name}/versions` - List available versions
- `DELETE /v1/packages/{name}/{version}` - Remove a package
- `GET /v1/packages/{name}/{version}/info` - Get package metadata

## Example Workflow

Here's a typical workflow for sharing agent stacks through the registry:

1. **Create** a stack locally with `sentinel stack init`
2. **Test** it to ensure it works as expected with `sentinel stack run`
3. **Package** it with `sentinel stack push --build` to create a `.stack.sntl` file
4. **Sign** it during packaging with `--sign` to add a cryptographic signature
5. **Push** it to the registry with `sentinel stack push`
6. **Share** the reference with colleagues
7. **Pull** it on another machine with `sentinel stack pull`

## Troubleshooting

### Connection Issues

If you have trouble connecting to the registry:

1. Check your internet connection
2. Verify the registry URL is correct
3. Ensure your authentication token is valid
4. Check if the registry service is running

### Authentication Problems

If you can't authenticate:

1. Try logging in again with `sentinel login`
2. Check if your account has the necessary permissions
3. Verify your API key hasn't expired

### Package Issues

If you can't push or pull packages:

1. Ensure the package exists
2. Check your namespace permissions
3. Verify the signature if required
4. Make sure you're using the correct package name and version 