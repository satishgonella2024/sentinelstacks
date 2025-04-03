# Memory Management

SentinelStacks provides a flexible memory management system that enables agents to store, retrieve, and search through information. This guide explains how the memory system works and how to use it effectively.

## Overview

The memory system in SentinelStacks serves several critical purposes:

- **State Persistence**: Allows agents to maintain state across sessions
- **Knowledge Storage**: Provides a way to store and retrieve structured information
- **Semantic Search**: Enables finding relevant information using vector embeddings
- **Multi-Agent Sharing**: Facilitates knowledge sharing between agents in a stack

## Memory Store Types

SentinelStacks supports several types of memory stores:

### Key-Value Stores

| Store Type | Persistence | Description |
|------------|-------------|-------------|
| `local` | In-memory | Fast, non-persistent storage that exists only during runtime |
| `sqlite` | File-based | Persistent storage using SQLite database files |

### Vector Stores

| Store Type | Persistence | Description |
|------------|-------------|-------------|
| `local` | In-memory | Simple vector storage for development and testing |
| `sqlite` | File-based | Persistent vector storage using SQLite and vector extensions |
| `chroma` | External | Integration with the Chroma vector database for advanced retrieval |

## Using Memory in Sentinelfiles

To configure memory for an agent, add memory settings to your Sentinelfile:

```yaml
name: ResearchAssistant
description: An agent that remembers research topics
baseModel: claude-3-sonnet-20240229

memory:
  type: sqlite
  path: ./data/research_agent
  retention: 30d
  vectorStore: chroma
  
  settings:
    maxSize: 1048576
    embedModel: text-embedding-3-small
```

### Memory Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `type` | Memory store type (`local`, `sqlite`) | `local` |
| `path` | Storage location for persistent stores | `./data/memory` |
| `retention` | How long to retain memory entries | `7d` |
| `vectorStore` | Vector store type (`local`, `sqlite`, `chroma`) | `local` |
| `settings.maxSize` | Maximum size of stored values in bytes | `1048576` (1MB) |
| `settings.embedModel` | Model to use for generating embeddings | `text-embedding-3-small` |

## Memory Management Using API

For programmatic access, SentinelStacks provides a Memory service through its API:

```go
package main

import (
    "context"
    "fmt"
    "github.com/sentinelstacks/sentinel/pkg/api"
)

func main() {
    ctx := context.Background()
    
    // Initialize API
    sentinelAPI, err := api.NewAPI(api.Config{
        MemoryConfig: api.MemoryServiceConfig{
            StoragePath: "./data/memory",
        },
    })
    if err != nil {
        panic(err)
    }
    defer sentinelAPI.Close()
    
    // Get memory service
    memoryService := sentinelAPI.Memory()
    
    // Store a value
    err = memoryService.StoreValue(ctx, "agent123", "last_search", "quantum computing")
    if err != nil {
        panic(err)
    }
    
    // Retrieve a value
    value, err := memoryService.RetrieveValue(ctx, "agent123", "last_search")
    if err != nil {
        panic(err)
    }
    fmt.Println("Last search:", value)
}
```

## CLI Commands for Memory

SentinelStacks includes several CLI commands for working with agent and stack memory:

### Viewing Agent Memory

```bash
# View an agent's memory state
sentinel agent memory research-agent

# Export an agent's memory to a file
sentinel agent memory research-agent --export memory.json

# View specific memory keys
sentinel agent memory research-agent --keys last_search,user_preferences
```

### Managing Memory

```bash
# Clear an agent's memory
sentinel agent memory research-agent --clear

# Import memory from a file
sentinel agent memory research-agent --import memory.json

# Set retention period for an agent's memory
sentinel agent memory research-agent --retention 90d
```

### Working with Vector Memory

```bash
# Search agent's vector memory
sentinel agent vectorsearch research-agent "quantum computing"

# Add documents to vector memory
sentinel agent vectoradd research-agent --file research_paper.txt

# List vector memory entries
sentinel agent vectorlist research-agent
```

## Shell Integration

When using the interactive shell with an agent, you can access memory commands:

```bash
# Start a shell session with an agent
sentinel shell research-agent

# In the shell:
> memory                  # Display all memory entries
> memory last_search      # Display a specific memory entry
> memory clear            # Clear all memory
> memory set key value    # Set a memory value
```

## Memory Persistence

Memory persistence depends on the configured store type:

1. **Local Store**: Memory exists only during runtime and is lost when the agent stops
2. **SQLite Store**: Memory is persisted to disk in SQLite database files
3. **Chroma Store**: Vector embeddings are stored in a Chroma database

For persistent stores, data is saved in the location specified by the `path` configuration or in the default location (`./data/memory`).

## Vector Search

The vector memory system enables semantic search capabilities:

```go
// Store a document with automatic embedding
err := memoryService.StoreDocument(ctx, "agent123", "doc1", "Quantum computing uses qubits instead of classical bits", nil)

// Search for semantically similar content
results, err := memoryService.SearchSimilar(ctx, "agent123", "How do quantum computers differ from regular computers?", 3)
```

The vector search system supports:

- **Semantic Matching**: Find relevant information based on meaning, not just keywords
- **Metadata Filtering**: Filter search results by metadata attributes
- **Document Chunking**: Automatically split large documents into searchable chunks

## Stack Memory Sharing

Agents within a stack can share memory through the stack memory system:

```yaml
# In Stackfile.yaml
agents:
  - name: researcher
    # ... other configuration ...
    memory:
      type: sqlite
      shared: ["research_results"]
      
  - name: writer
    # ... other configuration ...
    memory:
      type: sqlite 
      access: ["research_results"]
```

The `shared` attribute specifies memory keys that other agents can access, while the `access` attribute defines which shared keys this agent can read from.

## Advanced Usage

### Memory Backends

You can configure different memory backends for different types of data:

```yaml
memory:
  keyValue:
    type: sqlite
    path: ./data/agent_kv
  
  vector:
    type: chroma
    endpoint: http://localhost:8000
    collection: agent_vectors
```

### Memory Plugins

SentinelStacks supports custom memory plugins for specialized storage needs:

```yaml
memory:
  type: custom
  plugin: redis-store
  settings:
    host: localhost
    port: 6379
    password: secret
```

To implement a custom memory plugin, see the [Developer Guide](/development/memory_plugins.md).

## Best Practices

1. **Choose the Right Store Type**: 
   - For development/testing, use `local` memory
   - For production, use `sqlite` or `chroma` for persistence

2. **Manage Memory Growth**: 
   - Set appropriate retention periods
   - Use `maxSize` to limit individual entry sizes
   - Periodically clean up obsolete data

3. **Optimize Vector Search**:
   - Use specific, focused queries
   - Include relevant metadata for filtering
   - Choose appropriate vector models for your content

4. **Security Considerations**:
   - Avoid storing sensitive information in memory
   - Use encryption for sensitive data if necessary
   - Implement access controls for shared memory

5. **Performance Optimization**:
   - Use the appropriate store type for your scale
   - Consider external vector databases for large collections
   - Implement caching for frequently accessed data 