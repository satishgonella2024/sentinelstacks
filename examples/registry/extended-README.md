# SentinelStacks Registry Integration

This directory contains examples and documentation for the SentinelStacks registry system, which enables sharing stacks and agents between users.

## Custom File Formats

SentinelStacks uses custom file formats with specific extensions to clearly identify different types of content:

| Extension | Description | MIME Type |
|-----------|-------------|-----------|
| `.agent.sntl` | Agent package | `application/x-sentinel-agent` |
| `.stack.sntl` | Stack package | `application/x-sentinel-stack` |
| `.agent.yaml` | Agent definition | `application/x-sentinel-agent-def+yaml` |
| `.stack.yaml` | Stack definition | `application/x-sentinel-stack-def+yaml` |
| `.sig.sntl` | Detached signature | `application/x-sentinel-signature` |

This consistent naming scheme makes it easier to identify the purpose of each file and ensures proper handling by the system.

## Package Format

The `.sntl` package format has several unique characteristics:

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

## Registry API

The registry system provides a complete REST API for managing stacks and agents:

- `GET /v1/packages/search` - Search for packages
- `GET /v1/packages/{name}/{version}` - Download a package
- `PUT /v1/packages/publish` - Upload a new package
- `GET /v1/packages/{name}/versions` - List available versions
- `DELETE /v1/packages/{name}/{version}` - Remove a package
- `GET /v1/packages/{name}/{version}/info` - Get package metadata

## Authentication

The registry uses OAuth2 and JWT tokens for authentication, with support for different access levels:

- **Read-only**: Can download public packages
- **Publisher**: Can upload packages under their own namespace
- **Admin**: Has full management access

## Example Workflow

Here's how a typical workflow with the registry works:

1. **Create** a stack locally with the `sentinel stack init` command
2. **Test** it to ensure it works as expected with `sentinel stack run`
3. **Package** it with `sentinel stack push --build` to create a `.stack.sntl` file
4. **Sign** it during packaging with `--sign` to add a cryptographic signature
5. **Push** it to the registry with `sentinel stack push`
6. **Share** the reference with colleagues
7. **Pull** it on another machine with `sentinel stack pull`

## Configuration

The registry connection can be configured in your SentinelStacks config:

```yaml
registry:
  url: https://registry.sentinelstacks.io
  auth_token: your-auth-token
  cache_dir: ~/.sentinel/cache
  verify_ssl: true
  timeout: 30
```

## Example Files

This directory includes:

- `text-analyzer.stack.yaml`: An example stack definition
- `sentinel.manifest.json`: An example manifest file from a packaged stack
- `README.md`: Basic usage instructions

## Security Features

The registry system includes several security features:

1. **Package Signing**: Cryptographic signatures to verify authorship
2. **Integrity Checks**: SHA-256 hashes to verify file contents
3. **Access Control**: Granular permissions for package management
4. **Audit Logs**: Track all package operations
5. **Version Immutability**: Published versions cannot be modified

## Offline Mode

Even without internet access, you can use locally cached packages:

```bash
# Enable offline mode
sentinel config set registry.offline true

# Use cached packages
sentinel stack run mystack.stack.yaml
```
