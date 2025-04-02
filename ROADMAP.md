# SentinelStacks Project Roadmap

This document outlines the strategic direction, implementation plan, and progress tracking for the SentinelStacks project.

## 🔭 Strategic Direction: Phase Prioritization

SentinelStacks development is organized into three focused phases that build upon each other to deliver a complete AI agent management system.

### ✅ Phase 1: Developer-Complete Core System (Weeks 1–3)
"Ship a complete developer experience for solo/local usage with full CLI parity and NLP agent generation."

**Must-Have Goals:**
- Finalize single-agent flow (init → run → log → stop → push/pull)
- Expose all runtime capabilities via REST/gRPC
- NLP-to-Stackfile support
- `stack run` to support multi-agent execution
- Agent state context manager

**Core Enhancements:**
- Rename compose → stack throughout the codebase
- Implement StackEngine with DAG runner (Go)
- Add `sentinel stack run` CLI entrypoint
- NLP Parser updates to emit multi-agent Stackfile (with dependent Sentinelfiles)
- Runtime context propagation logic (data passing)

### ✅ Phase 2: Team Workflow + Collaboration Layer (Weeks 4–5)
"Unlock collaboration with registry, agent versioning, and shared memory."

**Must-Have Goals:**
- Complete Registry API (auth, push/pull, tags, search, versioning)
- Add CLI commands for registry (stack push, stack pull, stack login)
- Enable teams to import existing agents into new stacks
- Secure stack/agent signature + verification
- Agent run ID tracing for observability

**Suggested Add-ons:**
- Memory plugins (Chroma, SQLite, in-memory)
- Stack metadata view (stack inspect, stack history)

### ✅ Phase 3: UX & Product Readiness (Weeks 6–8)
"Complete UI and prepare for early launch or OSS release."

**UI Enhancements:**
- Stack visualizer: DAG viewer
- Stack launcher from UI: form-based input → invoke CLI or backend
- Agent logs view per run
- Model selector (Claude, OpenAI, Ollama)
- Real-time WebSocket logs
- User onboarding + template agents

**Docs + DX:**
- Finalize README, SETUP, and mkdocs structure
- CLI man pages
- Examples showcase
- Add test coverage to runtime/, shim/, stack/, and api/

## 🧱 Structural Design Enhancements

### 🧩 1. Stack Execution Engine (Replace Compose)

Move `cmd/sentinel/compose/` → `cmd/sentinel/stack/`
- Rename compose.go → run.go
- Create internal/stack/engine.go:

```go
type StackSpec struct {
  Agents []StackAgentSpec
}

type StackAgentSpec struct {
  ID         string
  Uses       string
  InputFrom  string
  InputKey   string
  OutputKey  string
  Params     map[string]interface{}
}
```

Use topological sort + DAG traversal to run stack.

### 🔌 2. Vector Memory Plugin Interface

```go
type MemoryStore interface {
  Save(key string, data interface{}) error
  Load(key string) (interface{}, error)
  Query(embedding []float32, topK int) ([]MemoryMatch, error)
}
```

**Backends:**
- memory_local.go: in-memory
- memory_chroma.go: REST to Chroma
- memory_pg.go: Postgres + pgvector

### 🔧 3. Enhanced NLP → Stackfile Generator

Add internal/parser/stack_parser.go:
- Analyze user prompt (multi-step intent)
- Generate:
  - Stackfile with agent DAG
  - Sentinelfile for each node
  - Optional stub test files

## 📂 Folder Refactor

```
internal/
├── runtime/           # Agent execution
├── stack/             # Stack DAG runner, parser, engine
├── memory/            # Vector store plugins
├── parser/            # NLP & YAML parser
├── shim/              # LLM providers
├── registry/          # Registry service
├── api/               # REST / gRPC APIs
```

## 📊 Visual Roadmap

![SentinelStacks Phase Roadmap](docs/visualizations/phase-roadmap.png)

## 🏗️ Enhanced Architecture

![SentinelStacks Enhanced Architecture](docs/visualizations/enhanced-architecture.png)

## 🔄 Stack Engine Architecture

![Stack Engine Architecture](docs/visualizations/stack-engine-detail.png)

## 📋 Progress Tracker

The latest progress for each feature can be found in [PROGRESS.md](./PROGRESS.md).

## 📝 Detailed Implementation Plan

For technical details and specific task breakdowns, refer to [IMPLEMENTATION.md](./IMPLEMENTATION.md).

## 🚀 Getting Started

To contribute to this roadmap:

1. Review the overall plan and progress tracker
2. Choose a task from the current phase
3. Create a new branch with the format `feature/phase{1,2,3}-{task-name}`
4. Submit a PR with implementation
5. Update the progress tracker with your changes

## 📞 Weekly Sync

We'll hold weekly sync meetings to review progress and adjust priorities:
- **When**: Fridays at 10:00 AM PT
- **Where**: Project Discord #roadmap-sync channel
- **Format**: 15-minute status review, 15-minute blockers discussion, 15-minute next steps

## 🤝 Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on contributing to this project.
