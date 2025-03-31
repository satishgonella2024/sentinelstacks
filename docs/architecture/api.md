# SentinelStacks API Reference

This document provides a reference for the SentinelStacks API, allowing developers to integrate with and extend the SentinelStacks ecosystem programmatically.

## API Basics

- **Base URL**: `https://api.sentinelstacks.com/v1`
- **Authentication**: Bearer token in the Authorization header
- **Content Type**: `application/json`
- **Rate Limiting**: 100 requests per minute per authenticated user

### Authentication

To authenticate with the API, obtain an access token:

```bash
# Using the CLI
sentinel login

# The token is stored in ~/.sentinel/config.json
```

For API requests, include the token in the Authorization header:

```
Authorization: Bearer your-access-token
```

## Core API Endpoints

### Agent Management

#### List Images

```
GET /images
```

Response:

```json
{
  "images": [
    {
      "id": "sha256:a1b2c3...",
      "name": "username/agent-name",
      "tag": "latest",
      "created_at": "2025-03-30T12:00:00Z",
      "size": 1024,
      "llm": "claude-3.7-sonnet"
    },
    ...
  ]
}
```

#### Get Image Details

```
GET /images/{image_id}
```

Response:

```json
{
  "id": "sha256:a1b2c3...",
  "name": "username/agent-name",
  "tag": "latest",
  "created_at": "2025-03-30T12:00:00Z",
  "size": 1024,
  "llm": "claude-3.7-sonnet",
  "capabilities": ["web_search", "document_analysis"],
  "parameters": {
    "memory_retention": "7d",
    "search_depth": 10
  },
  "metadata": {
    "description": "An agent that helps with research tasks",
    "author": "username",
    "version": "1.0.0"
  }
}
```

#### List Running Agents

```
GET /agents
```

Response:

```json
{
  "agents": [
    {
      "id": "agent_abc123",
      "image": "username/agent-name:latest",
      "status": "running",
      "created_at": "2025-03-30T14:00:00Z",
      "memory_usage": 128,
      "api_calls": 42
    },
    ...
  ]
}
```

#### Start Agent

```
POST /agents
```

Request:

```json
{
  "image": "username/agent-name:latest",
  "parameters": {
    "memory_retention": "30d",
    "search_depth": 15
  },
  "environment": {
    "API_KEY": "your-service-api-key"
  }
}
```

Response:

```json
{
  "id": "agent_def456",
  "image": "username/agent-name:latest",
  "status": "starting",
  "created_at": "2025-03-31T09:00:00Z",
  "endpoints": {
    "chat": "wss://api.sentinelstacks.com/v1/agents/agent_def456/chat",
    "events": "wss://api.sentinelstacks.com/v1/agents/agent_def456/events"
  }
}
```

#### Stop Agent

```
DELETE /agents/{agent_id}
```

Response:

```json
{
  "id": "agent_def456",
  "status": "stopping"
}
```

#### Get Agent Logs

```
GET /agents/{agent_id}/logs
```

Parameters:
- `limit` (optional): Number of log entries to return (default: 100)
- `since` (optional): Return logs since timestamp

Response:

```json
{
  "logs": [
    {
      "timestamp": "2025-03-31T09:01:00Z",
      "level": "info",
      "message": "Agent initialized successfully"
    },
    {
      "timestamp": "2025-03-31T09:01:05Z",
      "level": "debug",
      "message": "Connecting to LLM provider"
    },
    ...
  ]
}
```

### Registry Operations

#### Search Registry

```
GET /registry/search
```

Parameters:
- `q` (required): Search query
- `limit` (optional): Number of results to return (default: 20)
- `page` (optional): Page number for pagination (default: 1)
- `category` (optional): Filter by category
- `capabilities` (optional): Filter by capabilities (comma-separated)

Response:

```json
{
  "results": [
    {
      "name": "username/agent-name",
      "description": "An agent that helps with research tasks",
      "stars": 42,
      "downloads": 1024,
      "created_at": "2025-02-15T10:00:00Z",
      "tags": ["research", "productivity"],
      "verified": true
    },
    ...
  ],
  "total": 156,
  "page": 1,
  "limit": 20
}
```

#### Push Image

```
POST /registry/push
```

Request:

```json
{
  "image_id": "sha256:a1b2c3...",
  "name": "username/agent-name",
  "tag": "latest",
  "description": "An agent that helps with research tasks",
  "tags": ["research", "productivity"],
  "public": true
}
```

Response:

```json
{
  "upload_url": "https://storage.sentinelstacks.com/upload/token123",
  "token": "upload_token_xyz",
  "expires_at": "2025-03-31T10:00:00Z"
}
```

#### Pull Image

```
POST /registry/pull
```

Request:

```json
{
  "name": "username/agent-name",
  "tag": "latest"
}
```

Response:

```json
{
  "download_url": "https://storage.sentinelstacks.com/download/token456",
  "token": "download_token_abc",
  "expires_at": "2025-03-31T10:00:00Z",
  "image": {
    "id": "sha256:d4e5f6...",
    "size": 1024,
    "checksum": "sha256:abcdef123456..."
  }
}
```

### Chat Interaction

#### Connect to Agent Chat

Establish a WebSocket connection:

```
WebSocket: wss://api.sentinelstacks.com/v1/agents/{agent_id}/chat
```

Send messages:

```json
{
  "type": "message",
  "content": "Hello, I need help with my research on climate change.",
  "message_id": "msg_123"
}
```

Receive responses:

```json
{
  "type": "response",
  "content": "I'd be happy to help with your climate change research. What specific aspects are you interested in?",
  "message_id": "msg_123",
  "response_id": "resp_456",
  "tools_used": []
}
```

#### Tool Usage

When the agent uses tools:

```json
{
  "type": "tool_request",
  "tool": "web_search",
  "parameters": {
    "query": "latest climate change research 2025",
    "limit": 5
  },
  "request_id": "tool_req_789"
}
```

Tool response (from client to server):

```json
{
  "type": "tool_response",
  "request_id": "tool_req_789",
  "status": "success",
  "data": {
    "results": [
      {
        "title": "Climate Change Report 2025",
        "url": "https://example.com/climate-report-2025",
        "snippet": "The latest findings suggest..."
      },
      ...
    ]
  }
}
```

## Webhook Events

### Configure Webhooks

```
POST /webhooks
```

Request:

```json
{
  "url": "https://your-server.com/webhook",
  "events": ["agent.started", "agent.stopped", "agent.error"],
  "secret": "your_webhook_secret"
}
```

Response:

```json
{
  "id": "webhook_123",
  "url": "https://your-server.com/webhook",
  "events": ["agent.started", "agent.stopped", "agent.error"],
  "created_at": "2025-03-31T12:00:00Z"
}
```

### Webhook Payload Example

```json
{
  "event": "agent.started",
  "timestamp": "2025-03-31T14:00:00Z",
  "data": {
    "agent_id": "agent_ghi789",
    "image": "username/agent-name:latest",
    "status": "running"
  },
  "signature": "sha256:123abc..."
}
```

## Error Handling

All API errors follow this format:

```json
{
  "error": {
    "code": "not_found",
    "message": "Agent not found",
    "details": {
      "agent_id": "agent_xyz789"
    }
  },
  "request_id": "req_456abc"
}
```

Common error codes:
- `unauthorized`: Authentication required or failed
- `forbidden`: Insufficient permissions
- `not_found`: Resource not found
- `bad_request`: Invalid parameters
- `rate_limited`: Too many requests
- `internal_error`: Server error

## SDK Examples

### Go SDK

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sentinelstacks/sdk-go/client"
)

func main() {
	// Create client
	client := client.NewClient("your-access-token")

	// Start an agent
	agent, err := client.StartAgent(context.Background(), &client.StartAgentRequest{
		Image: "username/research-assistant:latest",
		Parameters: map[string]interface{}{
			"memory_retention": "14d",
		},
	})
	if err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	fmt.Printf("Agent started with ID: %s\n", agent.ID)
	fmt.Printf("Chat at: %s\n", agent.Endpoints.Chat)
}
```

### Python SDK

```python
from sentinelstacks import Client

# Create client
client = Client(token="your-access-token")

# Start an agent
agent = client.start_agent(
    image="username/research-assistant:latest",
    parameters={
        "memory_retention": "14d"
    }
)

print(f"Agent started with ID: {agent.id}")
print(f"Chat at: {agent.endpoints.chat}")
```

## Rate Limiting

- Headers include rate limit information:
  - `X-RateLimit-Limit`: Requests allowed per minute
  - `X-RateLimit-Remaining`: Requests remaining in the current window
  - `X-RateLimit-Reset`: Time when the rate limit resets (Unix timestamp)

- When rate limited, the API returns a 429:

```json
{
  "error": {
    "code": "rate_limited",
    "message": "Rate limit exceeded",
    "details": {
      "limit": 100,
      "reset_at": 1717144800
    }
  },
  "request_id": "req_789def"
}
```
