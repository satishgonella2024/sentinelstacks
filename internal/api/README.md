# SentinelStacks API Server

This package implements the REST API server for SentinelStacks, providing endpoints for managing agents, images, and registry operations.

## Overview

The API server provides a RESTful interface for interacting with SentinelStacks programmatically. It enables developers to:

- Manage agents (create, list, delete, get logs)
- Query images (list, details)
- Interact with the registry (search, push, pull)
- Authenticate using JWT tokens

## Architecture

The API server follows a standard HTTP server architecture with middleware for common functionality:

- **Server**: Central component that manages routes and server lifecycle
- **Middleware**: Components that process requests (logging, CORS, authentication, recovery)
- **Handlers**: Request handlers for specific endpoints
- **Models**: Data structures shared between components

## API Endpoints

### Authentication

- `POST /v1/auth/login` - Authenticate and get a JWT token

### Agents

- `GET /v1/agents` - List all running agents
- `POST /v1/agents` - Create a new agent
- `GET /v1/agents/{id}` - Get details about a specific agent
- `DELETE /v1/agents/{id}` - Stop and delete an agent
- `GET /v1/agents/{id}/logs` - Get logs for an agent

### Images

- `GET /v1/images` - List all available images
- `GET /v1/images/{id}` - Get details about a specific image

### Registry

- `GET /v1/registry/search` - Search for images in the registry
- `POST /v1/registry/push` - Push an image to the registry
- `POST /v1/registry/pull` - Pull an image from the registry

## WebSocket Support

The API server includes WebSocket support for real-time communication with agents:

### Chat WebSocket (`/v1/agents/{id}/chat`)

Enables real-time chat with an agent. The connection follows this protocol:

1. Client connects to the WebSocket endpoint
2. Server sends a welcome message
3. Client sends messages to the agent
4. Server streams responses back to the client

#### Message Types (Client to Server)

- `message`: Regular message to the agent
  ```json
  {
    "type": "message",
    "content": "Hello, how can you help me?",
    "message_id": "optional-client-generated-id"
  }
  ```

- `tool_request`: Request to use a tool
  ```json
  {
    "type": "tool_request",
    "tool": "web_search",
    "parameters": {
      "query": "latest AI developments"
    },
    "request_id": "optional-client-generated-id"
  }
  ```

#### Message Types (Server to Client)

- `response`: Complete response from the agent
- `stream_start`: Indicates the start of a streaming response
- `stream_chunk`: A chunk of a streaming response
- `stream_end`: Indicates the end of a streaming response
- `event`: Notification of an event (thinking, processing, etc.)
- `tool_result`: Result from a tool execution
- `error`: Error message

### Events WebSocket (`/v1/agents/{id}/events`)

Provides a stream of events related to an agent:

- Agent status changes
- Message processing events
- Tool usage events
- Error events

Example event:
```json
{
  "type": "event",
  "event_type": "agent_status",
  "timestamp": "2024-04-10T12:34:56Z",
  "content": "Agent demo-agent is running",
  "data": {
    "agent_id": "agent123",
    "status": "running",
    "memory": "256MB",
    "api_usage": {
      "requests": 10,
      "tokens": 1500
    }
  }
}
```

## Running the API Server

The API server can be started using the CLI command:

```
sentinel api [flags]
```

### Available Flags

- `--port, -p` - Port to listen on (default: 8080)
- `--host` - Host address to listen on (default: localhost)
- `--tls-cert` - TLS certificate file path
- `--tls-key` - TLS key file path
- `--token-auth-secret` - Secret for JWT authentication
- `--cors` - Enable CORS (default: true)
- `--log-requests` - Log API requests (default: true)

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. To obtain a token, send a POST request to `/v1/auth/login` with your credentials:

```json
{
  "username": "user",
  "password": "password"
}
```

The response will include a token:

```json
{
  "token": "eyJhbGciOiJ...",
  "username": "user",
  "expires": "2024-04-20T12:00:00Z"
}
```

Include this token in the Authorization header for subsequent requests:

```
Authorization: Bearer eyJhbGciOiJ...
```

## Implementation Details

### Middleware

- **Logging**: Logs all API requests
- **CORS**: Handles Cross-Origin Resource Sharing
- **Recovery**: Recovers from panics and returns appropriate error responses
- **Authentication**: Validates JWT tokens and adds user info to request context

### Error Handling

Errors are returned in a consistent format:

```json
{
  "error": "Error message describing what went wrong"
}
```

HTTP status codes are used appropriately to indicate the type of error.

## Future Enhancements

- Enhanced WebSocket authentication
- Rate limiting and quota management
- OAuth2 authentication
- API versioning strategy
- Swagger/OpenAPI documentation 