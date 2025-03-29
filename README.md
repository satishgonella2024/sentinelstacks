# SentinelStacks

🤖 An AI-powered infrastructure management platform that helps organizations automate, secure, and manage their cloud resources using intelligent agents.

## Current Status

SentinelStacks is currently in active development. Here's what's working:

### ✅ Production Ready
- CLI tool with enhanced UI
- Model integrations (OpenAI, Claude, Ollama)
- Basic agent runtime
- Local file storage

### 🔄 In Development
- Memory system enhancements
- Desktop application
- Registry system
- Advanced agent features

## Quick Start

### Prerequisites
- Go 1.21 or later
- Node.js 18+ (for desktop UI)
- Rust (for Tauri desktop app)
- One of the supported LLM providers:
  - OpenAI API key
  - Anthropic API key
  - Ollama (local models)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/satishgonella2024/sentinelstacks.git
cd sentinelstacks
```

2. Build the CLI:
```bash
./scripts/build.sh
```

3. Set up your environment:
```bash
# For OpenAI
export OPENAI_API_KEY=your_key_here

# For Anthropic
export ANTHROPIC_API_KEY=your_key_here

# For local models
# Install Ollama from https://ollama.ai
```

### Basic Usage

1. List available agents:
```bash
./dist/darwin-arm64/sentinel registry list
```

2. Create a new agent:
```bash
./dist/darwin-arm64/sentinel agent create --name test-agent
```

3. Run an agent:
```bash
./dist/darwin-arm64/sentinel agent run --name test-agent --interactive
```

## Features

### Implemented ✅

1. **Enhanced CLI**
   - Animated progress indicators
   - Color-coded output
   - Interactive mode
   - Clear success/error states

2. **Model Integration**
   - OpenAI (GPT-3.5, GPT-4)
   - Anthropic Claude
   - Ollama (local models)
   - Configurable parameters

3. **Basic Memory System**
   - Key-value storage
   - Vector embeddings
   - Persistence
   - Simple search

### In Development 🔄

1. **Memory Enhancements**
   - Context window management
   - Advanced retrieval
   - Optimization

2. **Desktop Application**
   - Agent management
   - Monitoring
   - Settings

3. **Registry System**
   - Remote storage
   - Version control
   - Sharing

## Development

### Project Structure
```
sentinelstacks/
├── cmd/                    # Command-line tools
├── internal/              # Private application code
├── pkg/                   # Public libraries
├── desktop/              # Tauri desktop app
├── docs/                 # Documentation
└── scripts/              # Build and utility scripts
```

### Running Tests
```bash
go test ./...
```

### Building Desktop App
```bash
cd desktop
npm install
npm run tauri dev
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Documentation

- [Architecture](ARCHITECTURE.md)
- [Development Plan](DEVELOPMENT_PLAN.md)
- [API Reference](docs/api/README.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features and enhancements.

## Git Strategy

We follow a modified GitFlow workflow:

### Branch Structure
```
main
├── develop
│   ├── feature/*
│   ├── bugfix/*
│   └── docs/*
├── release/*
└── hotfix/*
```

### Branch Types
1. **Main Branches**
   - `main`: Production-ready code
   - `develop`: Integration branch for features

2. **Supporting Branches**
   - `feature/*`: New features
   - `bugfix/*`: Bug fixes
   - `docs/*`: Documentation updates
   - `release/*`: Release preparation
   - `hotfix/*`: Production fixes

### Branch Naming
- Features: `feature/[ticket-number]-description`
- Bugfixes: `bugfix/[ticket-number]-description`
- Docs: `docs/[topic]-update`
- Releases: `release/v[major].[minor].[patch]`
- Hotfixes: `hotfix/[ticket-number]-description`

### Workflow
1. Create feature branch from `develop`
2. Develop and test changes
3. Create PR to merge back to `develop`
4. After review and approval, merge using `--no-ff`
5. Delete feature branch after merge

### Commit Messages
We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `chore:` Maintenance tasks
- `test:` Adding/updating tests
- `refactor:` Code refactoring
