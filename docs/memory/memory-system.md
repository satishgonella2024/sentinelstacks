# SentinelStacks Memory System

The memory system enables agents to remember previous interactions, store information, and retrieve relevant context during conversations. SentinelStacks provides two types of memory implementations: Simple Memory and Vector Memory.

## Memory Types

### Simple Memory

Simple memory is a basic key-value store that associates unique IDs with text content. It's efficient for small-scale agents and straightforward use cases.

**Features:**
- Fast text-based lookup
- Low resource requirements
- Simple string matching for searches
- Persistence between sessions

**Use cases:**
- Basic conversational agents
- Agents with limited needs for context recall
- When resources are constrained

### Vector Memory

Vector memory uses embeddings to represent text semantically, enabling more powerful retrieval based on meaning rather than just keywords.

**Features:**
- Semantic search capabilities
- Better handling of natural language variations
- Finds related information even with different wording
- More sophisticated relevance ranking

**Use cases:**
- Advanced conversational agents
- Knowledge-intensive applications
- When semantic understanding is important

## Configuration

Memory is configured in the agent's Agentfile:

```yaml
memory:
  type: "vector"  # Options: "simple" or "vector"
  persistence: true
  maxItems: 1000  # Maximum number of items to store
```

## Persistence

By default, memory is persisted to disk, allowing agents to maintain context across multiple sessions. Memory files are stored in:

- Simple Memory: `~/.sentinel/memory/<agent-name>.json`
- Vector Memory: `~/.sentinel/vectors/<agent-name>` and `~/.sentinel/memory/<agent-name>.json`

## Using Memory in Your Agent

### Adding to Memory

```go
// Store something in memory
memoryID, err := agent.Memory.Add("This is important information", map[string]interface{}{
    "source": "user",
    "timestamp": time.Now(),
})
```

### Retrieving from Memory

```go
// Get by ID
entry, err := agent.Memory.Get(memoryID)

// Search by content
results, err := agent.Memory.Search("important information", 5)
```

### Listing Memory

```go
// Get recent memories
entries, err := agent.Memory.List(10)
```

## Example: Creating an Agent with Vector Memory

From the command line:

```bash
sentinel agent create --name knowledgebot --description "A bot that remembers facts" --model llama3 --memory vector
```

Or in an Agentfile:

```yaml
name: knowledgebot
version: "1.0.0"
description: "A bot that remembers facts"
model:
  provider: "ollama"
  name: "llama3"
capabilities:
  - conversation
memory:
  type: "vector"
  persistence: true
  maxItems: 2000
```

## Performance Considerations

- **Simple Memory** uses less resources but has limited search capabilities
- **Vector Memory** provides better semantic search but requires more computational resources
- For large knowledge bases, consider setting appropriate `maxItems` limits

## Implementation Details

The memory system follows a clean interface design that allows for future expansion with additional memory types:

```go
type Memory interface {
    Add(content string, metadata map[string]interface{}) (string, error)
    Get(id string) (*MemoryEntry, error)
    Search(query string, limit int) ([]MemoryEntry, error)
    List(limit int) ([]MemoryEntry, error)
    Delete(id string) error
    Clear() error
    Save() error
    Load() error
}
```

This interface ensures that all memory implementations provide consistent functionality while allowing for specialized features.

## Future Enhancements

Planned improvements to the memory system include:

1. Hierarchical memory with short-term and long-term storage
2. Automatic summarization of memories
3. Importance-based retention
4. Multi-agent shared memory
5. Improved embedding models for vector memory

For more information on how to optimize memory usage for your specific agent, see the [Agent Development Guide](../agent/agent-development.md).