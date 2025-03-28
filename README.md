# SentinelStacks

SentinelStacks is an AI-powered infrastructure management platform that helps you automate, secure, and manage your cloud resources using intelligent agents.

## Features

- 🤖 **AI-Powered Automation**: Intelligent agents that understand your infrastructure and automate complex tasks
- 🔒 **Security First**: Built-in security features and compliance checks
- ☁️ **Multi-Cloud Support**: Manage resources across multiple cloud providers
- 🔌 **Extensible Platform**: Create and share custom agents

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.20 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/sentinelstacks.git
cd sentinelstacks
```

2. Start the services:
```bash
docker-compose up -d
```

3. Access the web interface:
- Landing page: https://localhost
- Registry UI: https://localhost/registry

### Using the CLI

Install agents:
```bash
./sentinel registry pull -name terraform-agent -version latest
./sentinel registry pull -name kubernetes-agent -version latest
```

Run an agent:
```bash
./sentinel agent run -name terraform-agent -version latest
```

## Project Structure

```
.
├── cmd/                    # Command-line tools
│   ├── api/               # API server
│   └── sentinel/          # CLI tool
├── internal/              # Internal packages
│   ├── api/              # API implementation
│   ├── agent/            # Agent management
│   └── registry/         # Registry implementation
├── landing/              # Landing page
├── registry-ui/          # Registry web interface
├── nginx/                # Nginx configuration
└── docker-compose.yml    # Docker services configuration
```

## Development

### Building from Source

```bash
# Build the CLI tool
go build -o sentinel cmd/sentinel/main.go

# Build the API server
go build -o api-server cmd/api/main.go
```

### Running Tests

```bash
go test ./...
```

### Adding a New Agent

1. Create a new directory in `examples/agents/`
2. Add your agent configuration in `agent.yaml`
3. Implement the required commands
4. Build and push to the registry:
```bash
./sentinel registry push -path examples/agents/your-agent
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all contributors who have helped shape SentinelStacks
- Built with Go, Docker, and ❤️
