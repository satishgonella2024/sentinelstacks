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

- WebSocket support for real-time agent interaction
- Rate limiting and quota management
- Enhanced authentication with OAuth2
- API versioning strategy
- Swagger/OpenAPI documentation 