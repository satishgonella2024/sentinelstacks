# Development Guide

This guide provides detailed instructions for setting up your development environment and contributing to SentinelStacks.

## Development Environment Setup

### Prerequisites

1. **Go Installation**:
   ```bash
   # macOS
   brew install go

   # Linux
   wget https://go.dev/dl/go1.20.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz
   ```

2. **Docker and Docker Compose**:
   ```bash
   # macOS
   brew install docker docker-compose

   # Linux
   curl -fsSL https://get.docker.com | sh
   sudo apt-get install docker-compose
   ```

3. **Development Tools**:
   ```bash
   # Install development dependencies
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

### Repository Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/sentinelstacks.git
   cd sentinelstacks
   ```

2. **Install Dependencies**:
   ```bash
   go mod download
   go mod verify
   ```

3. **Set Up Pre-commit Hooks**:
   ```bash
   cp scripts/pre-commit .git/hooks/
   chmod +x .git/hooks/pre-commit
   ```

## Project Structure

```
.
├── cmd/                    # Command-line applications
│   ├── api/               # API server
│   └── sentinel/          # CLI tool
├── internal/              # Private application code
│   ├── api/              # API implementation
│   ├── agent/            # Agent management
│   └── registry/         # Registry implementation
├── pkg/                   # Public libraries
│   ├── client/           # API client library
│   └── models/           # Shared data models
├── scripts/              # Development scripts
├── examples/             # Example configurations
└── docs/                 # Documentation
```

## Development Workflow

### 1. Create a New Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

Follow these guidelines when making changes:

1. **Code Style**:
   - Follow Go style guidelines
   - Use meaningful variable names
   - Add comments for complex logic
   - Keep functions focused and small

2. **Testing**:
   - Write unit tests for new code
   - Update existing tests if needed
   - Ensure all tests pass
   ```bash
   go test ./...
   ```

3. **Documentation**:
   - Update relevant documentation
   - Add inline code comments
   - Update API documentation if needed

### 3. Local Testing

1. **Run Linters**:
   ```bash
   golangci-lint run
   ```

2. **Start Local Services**:
   ```bash
   docker-compose up -d
   ```

3. **Build and Run**:
   ```bash
   # Build CLI
   go build -o sentinel cmd/sentinel/main.go

   # Build API server
   go build -o api-server cmd/api/main.go

   # Run API server
   ./api-server
   ```

4. **Manual Testing**:
   - Test new features
   - Verify bug fixes
   - Check edge cases

### 4. Submit Changes

1. **Commit Changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

2. **Update Your Branch**:
   ```bash
   git fetch origin
   git rebase origin/main
   ```

3. **Push Changes**:
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Create Pull Request**:
   - Open a PR on GitHub
   - Fill out the PR template
   - Request reviews
   - Address feedback

## Testing Guide

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test ./... -tags=integration

# Clean up
docker-compose -f docker-compose.test.yml down
```

### Performance Tests

```bash
# Run benchmarks
go test -bench=. ./...

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=. ./...
```

## Debugging

### 1. API Server

```bash
# Run with debug logging
DEBUG=1 ./api-server

# Check logs
docker-compose logs -f api
```

### 2. Database

```bash
# Connect to database
docker-compose exec db psql -U postgres

# View logs
docker-compose logs -f db
```

### 3. Redis

```bash
# Connect to Redis CLI
docker-compose exec redis redis-cli

# Monitor commands
docker-compose exec redis redis-cli monitor
```

## Common Tasks

### Adding a New API Endpoint

1. Define the endpoint in `internal/api/routes.go`
2. Implement the handler in `internal/api/handlers/`
3. Add tests in `internal/api/handlers/handler_test.go`
4. Update API documentation
5. Update the OpenAPI specification

### Creating a New Agent

1. Create agent directory in `examples/agents/`
2. Define `agent.yaml` configuration
3. Implement required commands
4. Add documentation
5. Create example usage

### Adding a New Feature Flag

1. Add flag to `internal/config/features.go`
2. Update configuration loading
3. Add feature documentation
4. Update deployment configurations

## Troubleshooting

### Common Issues

1. **Database Connection Issues**:
   ```bash
   # Check database status
   docker-compose ps db
   # View logs
   docker-compose logs db
   ```

2. **API Server Errors**:
   ```bash
   # Enable debug logging
   export DEBUG=1
   # Check logs
   tail -f logs/api.log
   ```

3. **Redis Connection Issues**:
   ```bash
   # Check Redis status
   docker-compose exec redis redis-cli ping
   ```

### Getting Help

1. Check existing issues on GitHub
2. Join our Discord community
3. Contact the maintainers
4. Submit a new issue

## Release Process

1. **Prepare Release**:
   ```bash
   # Update version
   ./scripts/bump-version.sh

   # Update changelog
   ./scripts/update-changelog.sh
   ```

2. **Create Release Branch**:
   ```bash
   git checkout -b release/v1.0.0
   ```

3. **Run Tests**:
   ```bash
   # Run all tests
   go test ./...

   # Run integration tests
   go test -tags=integration ./...
   ```

4. **Build Release**:
   ```bash
   # Build binaries
   ./scripts/build-release.sh

   # Create release tag
   git tag v1.0.0
   ```

5. **Push Release**:
   ```bash
   git push origin v1.0.0
   ```

## Additional Resources

- [Go Style Guide](https://golang.org/doc/effective_go)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation) 