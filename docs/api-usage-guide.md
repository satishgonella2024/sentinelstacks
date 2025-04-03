# SentinelStacks API Usage Guide

This guide provides a practical overview of how to interact with the SentinelStacks API, with a focus on the memory management capabilities.

## Authentication

Most API endpoints require authentication. To authenticate:

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "your_username", "password": "your_password"}'
```

Store the returned JWT token for use in subsequent requests:

```json
{
  "token": "eyJhbGciOiJ...",
  "username": "your_username",
  "expires": "2023-04-20T12:00:00Z"
}
```

## Memory Management API

SentinelStacks provides a memory system for agents to store and retrieve information across sessions. The API exposes endpoints for working with both key-value memory and vector embeddings.

### Storing Values in Memory

Store a key-value pair in an agent's memory:

```bash
curl -X POST http://localhost:8080/v1/memory/store \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent_abc123",
    "key": "user_preference",
    "value": "dark_mode",
    "metadata": {
      "source": "user_input",
      "timestamp": "2023-04-01T12:00:00Z"
    }
  }'
```

Response:

```json
{
  "success": true,
  "key": "user_preference"
}
```

### Retrieving Values from Memory

Retrieve a value by key:

```bash
curl -X GET "http://localhost:8080/v1/memory/retrieve?agent_id=agent_abc123&key=user_preference" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:

```json
{
  "key": "user_preference",
  "value": "dark_mode",
  "metadata": {
    "source": "user_input",
    "timestamp": "2023-04-01T12:00:00Z"
  }
}
```

### Semantic Search in Memory

Search for semantically similar entries in an agent's memory:

```bash
curl -X GET "http://localhost:8080/v1/memory/search?agent_id=agent_abc123&query=user%20settings&limit=5&threshold=0.7" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:

```json
{
  "results": [
    {
      "key": "user_preference",
      "value": "dark_mode",
      "score": 0.92,
      "metadata": {
        "source": "user_input",
        "timestamp": "2023-04-01T12:00:00Z"
      }
    },
    {
      "key": "color_scheme",
      "value": "blue",
      "score": 0.85,
      "metadata": {
        "source": "api",
        "timestamp": "2023-04-02T15:30:00Z"
      }
    }
  ]
}
```

### Deleting Memory Entries

Delete a specific memory entry by key:

```bash
curl -X DELETE "http://localhost:8080/v1/memory/delete?agent_id=agent_abc123&key=user_preference" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:

```json
{
  "success": true,
  "key": "user_preference"
}
```

## Using Memory in Agent Configuration

When creating a new agent, you can configure its memory settings:

```bash
curl -X POST http://localhost:8080/v1/agents \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "username/research-agent:latest",
    "parameters": {
      "memory_retention": "30d",
      "memory_type": "persistent",
      "embedding_model": "text-embedding-ada-002"
    }
  }'
```

## Memory Types and Persistence

SentinelStacks supports different types of memory stores with varying persistence characteristics:

1. **In-memory**: Fast but volatile. Data is lost when the agent stops.
2. **File-backed**: Persists data to disk. Suitable for most use cases.
3. **Database-backed**: Uses SQLite for storage. Best for production use with large memory volumes.

The memory type can be configured when creating an agent or in the `Sentinelfile` configuration.

## Best Practices

1. **Use meaningful keys**: Structure your memory keys hierarchically (e.g., `user/preferences/theme`) for easier organization.
2. **Include metadata**: Add metadata like timestamps, sources, and tags to make memory entries more useful.
3. **Set appropriate retention**: Use the `memory_retention` parameter to control how long memory entries are kept.
4. **Vector search optimization**: For semantic search, provide clear, descriptive text that captures the essence of what you want to store and retrieve.
5. **Memory cleanup**: Periodically delete obsolete entries to maintain performance and minimize storage usage.

## Error Handling

The API returns standard HTTP status codes:
- `200/201`: Success
- `400`: Bad request (invalid parameters)
- `401`: Unauthorized (authentication failed)
- `404`: Resource not found
- `500`: Server error

Error responses include an `error` field with a description:

```json
{
  "error": "Agent not found"
}
```

## Rate Limiting

The API enforces rate limits to protect system resources. If you exceed the limits, you'll receive a `429 Too Many Requests` response.

For higher limits, contact the system administrator or consider running a self-hosted instance.

## WebSocket Connections

For real-time interaction with agents, WebSocket connections are available:

```javascript
// Chat with an agent
const chatSocket = new WebSocket('ws://localhost:8080/v1/agents/agent_abc123/chat');

// Listen for agent events
const eventsSocket = new WebSocket('ws://localhost:8080/v1/agents/agent_abc123/events');
```

See the [WebSocket API documentation](./websocket-api.md) for details on the message formats and event types. 