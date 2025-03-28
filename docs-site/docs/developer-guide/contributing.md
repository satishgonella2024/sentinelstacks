# Contributing to SentinelStacks

This guide will help you set up your development environment and start contributing to SentinelStacks.

## Development Environment Setup

### Prerequisites

1. Install the following tools:
   - Go 1.20 or later
   - Docker and Docker Compose
   - Git
   - Make (optional, but recommended)
   - Your favorite code editor (VS Code recommended)

2. Install Go development tools:
```bash
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Repository Setup

1. Fork the repository on GitHub

2. Clone your fork:
```bash
git clone https://github.com/yourusername/sentinelstacks.git
cd sentinelstacks
```

3. Add the upstream remote:
```bash
git remote add upstream https://github.com/originalorg/sentinelstacks.git
```

4. Install dependencies:
```bash
go mod download
```

5. Install pre-commit hooks:
```bash
make install-hooks
```

## Project Structure

```
sentinelstacks/
├── cmd/                    # Command-line tools
│   ├── sentinel/          # Main CLI tool
│   └── registry/          # Registry service
├── internal/              # Private application code
│   ├── agent/            # Agent runtime
│   ├── auth/             # Authentication
│   ├── config/           # Configuration
│   └── registry/         # Registry implementation
├── pkg/                   # Public libraries
│   ├── agentfile/        # Agentfile parser
│   ├── models/           # Data models
│   └── api/              # API client
├── web/                  # Web interfaces
│   ├── landing/          # Landing page
│   └── registry-ui/      # Registry UI
├── scripts/              # Development scripts
├── examples/             # Example agents
└── docs/                 # Documentation
```

## Development Workflow

### Creating a New Feature

1. Create a feature branch:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes, following our coding standards

3. Run tests:
```bash
make test
```

4. Run linters:
```bash
make lint
```

5. Commit your changes:
```bash
git add .
git commit -m "feat: add your feature description"
```

6. Push to your fork:
```bash
git push origin feature/your-feature-name
```

7. Create a Pull Request on GitHub

### Running the Development Environment

1. Start all services:
```bash
docker-compose up -d
```

2. Run the API server:
```bash
make run-api
```

3. Run the registry:
```bash
make run-registry
```

### Testing

We use Go's built-in testing framework. Run tests with:

```bash
# Run all tests
make test

# Run specific tests
go test ./pkg/agentfile/...

# Run tests with coverage
make test-coverage
```

### Debugging

1. API Server logs:
```bash
docker-compose logs -f api
```

2. Database access:
```bash
docker-compose exec db psql -U postgres -d sentinelstacks
```

3. Redis CLI:
```bash
docker-compose exec redis redis-cli
```

## Common Tasks

### Adding a New API Endpoint

1. Define the endpoint in `internal/api/routes.go`
2. Create handler in `internal/api/handlers/`
3. Add tests in `internal/api/handlers/handler_test.go`
4. Update API documentation
5. Update OpenAPI specification

### Creating a New Agent

1. Create agent directory in `examples/`
2. Write Agentfile
3. Add documentation
4. Add tests
5. Submit to registry

### Adding a Feature Flag

1. Add flag to `internal/config/feature_flags.go`
2. Update configuration documentation
3. Add migration if needed
4. Update deployment scripts

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check PostgreSQL logs
   - Verify connection string
   - Ensure migrations are up to date

2. **Redis Connection Issues**
   - Check Redis logs
   - Verify Redis is running
   - Check connection settings

3. **API Server Issues**
   - Check API logs
   - Verify configuration
   - Check dependencies

### Getting Help

- Join our [Discord server](https://discord.gg/sentinelstacks)
- Check existing issues on GitHub
- Ask in our developer forum

## Release Process

1. Update version in `version.go`
2. Update CHANGELOG.md
3. Create release branch:
```bash
git checkout -b release/v1.2.3
```
4. Run full test suite:
```bash
make test-all
```
5. Create GitHub release
6. Push Docker images
7. Update documentation

## Additional Resources

- [Go Style Guide](https://golang.org/doc/effective_go.html)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)

## Code of Conduct

Please read our [Code of Conduct](../community/code-of-conduct.md) before contributing.
