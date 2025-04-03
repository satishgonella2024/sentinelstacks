# Agent Packaging and Security

SentinelStacks uses a custom packaging format for agents and stacks that ensures portability, security, and integrity. This guide explains how agent packaging works and how to protect your agents from tampering.

## Packaging Format Overview

### File Extensions

SentinelStacks uses distinct file extensions to identify different types of content:

| Extension | Description | MIME Type |
|-----------|-------------|-----------|
| `.agent.sntl` | Agent package | `application/x-sentinel-agent` |
| `.stack.sntl` | Stack package | `application/x-sentinel-stack` |
| `.agent.yaml` | Agent definition | `application/x-sentinel-agent-def+yaml` |
| `.stack.yaml` | Stack definition | `application/x-sentinel-stack-def+yaml` |
| `.sig.sntl` | Detached signature | `application/x-sentinel-signature` |

### Package Structure

Each `.sntl` package follows a specific structure:

1. **Magic Header**: Identifies the package type
   - Agent packages: `SNTL-AGENT-PKG`
   - Stack packages: `SNTL-STACK-PKG`

2. **Version Field**: A 4-byte version field identifying the package format version

3. **Compressed Archive**: A gzipped tar archive containing:
   - `sentinel.manifest.json`: Package metadata and file inventory
   - The main definition file (`.agent.yaml` or `.stack.yaml`)
   - Prompt templates and system prompts
   - Tool configurations
   - Documentation
   - Resources (like images or data files)

4. **Signature Block**: Optional cryptographic signatures

## Creating Packages

### Building an Agent Package

To create an agent package from a Sentinelfile:

```bash
# Build an agent package from a Sentinelfile
sentinel build ResearchAgent.yaml --output research-agent.agent.sntl
```

You can include additional resources:

```bash
# Include specific resources
sentinel build ResearchAgent.yaml --include-resources ./prompts,./data
```

### Building a Stack Package

For multi-agent stacks:

```bash
# Build a stack package from a Stackfile
sentinel stack build CollaborationStack.yaml --output collab-system.stack.sntl
```

## Package Security

SentinelStacks provides several features to secure agent packages and protect them from tampering.

### Package Signing

You can cryptographically sign packages to verify their authenticity:

```bash
# Sign a package during build
sentinel build ResearchAgent.yaml --sign --key-id my-signing-key

# Sign an existing package
sentinel sign research-agent.agent.sntl --key-id my-signing-key
```

### Signature Verification

To verify a signed package:

```bash
# Verify a package signature
sentinel verify research-agent.agent.sntl

# Verify against a specific key
sentinel verify research-agent.agent.sntl --key-id trusted-key
```

### Key Management

SentinelStacks includes a key management system for package signing:

```bash
# Generate a new signing key
sentinel key generate --name my-signing-key

# List available keys
sentinel key list

# Export a public key for sharing
sentinel key export --name my-signing-key --public-only

# Import a trusted key
sentinel key import colleague-key.pub --name colleague
```

### Tamper Protection

Packages include several tamper protection mechanisms:

1. **Content Verification**: Package manifests include content hashes for all files
2. **Signature Verification**: Cryptographic signatures ensure authenticity
3. **Runtime Verification**: Runtime environment verifies package integrity before execution

## Working with Packages

### Inspecting Packages

To inspect the contents of a package:

```bash
# Show package metadata
sentinel inspect research-agent.agent.sntl

# List files in a package
sentinel inspect research-agent.agent.sntl --list

# Extract a specific file
sentinel inspect research-agent.agent.sntl --extract prompts/main.txt
```

### Extracting Packages

You can extract packages to view or modify their contents:

```bash
# Extract an entire package
sentinel extract research-agent.agent.sntl --output ./extracted

# Extract with verification
sentinel extract research-agent.agent.sntl --verify
```

### Modifying Packages

To modify a package, extract it, make changes, and rebuild:

```bash
# Extract, modify, and rebuild (with new signature)
sentinel extract research-agent.agent.sntl --output ./temp
# ... make modifications ...
sentinel build ./temp/ResearchAgent.yaml --sign --key-id my-key
```

## Package Distribution

### Registry Storage

When stored in the registry, packages maintain their integrity through:

1. Package-level signatures
2. Transport-level encryption (TLS)
3. Registry-added checksums and metadata

### Secure Distribution

Best practices for secure package distribution:

1. **Always sign packages** before distribution
2. **Share public keys** with users who need to verify packages
3. **Use trusted registries** for distribution
4. **Verify packages** before running them in production

## Advanced Features

### Version Locking

You can lock a package to a specific LLM version:

```yaml
# In Sentinelfile
name: CriticalAgent
version: 1.0.0
baseModel: claude-3-opus-20240229
lockToModel: true  # Prevents running on different models
```

### Environment Requirements

Specify required environment features:

```yaml
# In Sentinelfile
name: SecurityAgent
environment:
  requiredFeatures:
    - secure-execution
    - isolated-network
  minVersion: "1.2.0"
```

### Content Versioning

Implement content versioning with package metadata:

```yaml
# In Sentinelfile
metadata:
  contentVersion: "2023-05-15"
  trainingCutoff: "2023-03-01"
  requiredData: ["2023Q1-security-db", "compliance-rules-v3"]
```

## Compatibility and Portability

### Cross-Platform Compatibility

SentinelStacks packages work across platforms:

- Linux, macOS, and Windows
- Server and desktop environments
- Container-based deployments

### Runtime Compatibility

Packages specify runtime compatibility:

```yaml
# In Sentinelfile
compatibility:
  sentinel: ">=1.0.0 <2.0.0"
  providers: ["claude", "openai"]
  features: ["streaming", "tools"]
```

## Best Practices

1. **Always sign packages** for distribution
2. **Include documentation** within your packages
3. **Version your packages** with semantic versioning
4. **Test packages** in isolated environments before deployment
5. **Maintain key security** for your signing keys
6. **Audit package contents** regularly
7. **Implement renewal** procedures for long-lived agents 