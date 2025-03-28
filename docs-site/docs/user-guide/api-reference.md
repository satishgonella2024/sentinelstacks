# API Reference

This document provides detailed information about the SentinelStacks API endpoints, authentication, and usage.

## Base URL

All API endpoints are relative to:
```
https://your-domain/api
```

## Authentication

Most endpoints require authentication using an API key. Include it in the request header:

```http
Authorization: Bearer your-api-key
```

## Endpoints

### Agents

#### List Agents

```http
GET /agents
```

Query Parameters:
- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 20)
- `status` (optional): Filter by agent status (active, inactive)

Response:
```json
{
  "agents": [
    {
      "id": "terraform-agent",
      "version": "1.0.0",
      "status": "active",
      "description": "Manages Terraform infrastructure",
      "created_at": "2024-03-20T10:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100
  }
}
```

#### Get Agent Details

```http
GET /agents/{name}
```

Parameters:
- `name`: Agent name (required)
- `version` (optional): Specific version (defaults to latest)

Response:
```json
{
  "id": "terraform-agent",
  "version": "1.0.0",
  "status": "active",
  "description": "Manages Terraform infrastructure",
  "created_at": "2024-03-20T10:00:00Z",
  "config": {
    "model": "gpt-4",
    "memory_size": 10000,
    "capabilities": ["infrastructure", "terraform"]
  }
}
```

#### Push Agent

```http
POST /agents
```

Request Body:
```json
{
  "name": "terraform-agent",
  "version": "1.0.0",
  "description": "Manages Terraform infrastructure",
  "config": {
    "model": "gpt-4",
    "memory_size": 10000,
    "capabilities": ["infrastructure", "terraform"]
  }
}
```

### Registry Operations

#### Search Registry

```http
GET /registry/search
```

Query Parameters:
- `q`: Search query
- `sort` (optional): Sort field (name, created_at, downloads)
- `order` (optional): Sort order (asc, desc)

Response:
```json
{
  "results": [
    {
      "id": "terraform-agent",
      "version": "1.0.0",
      "description": "Manages Terraform infrastructure",
      "downloads": 1500
    }
  ],
  "total": 25
}
```

### Agent Execution

#### Run Agent Command

```http
POST /agents/{name}/run
```

Request Body:
```json
{
  "command": "plan infrastructure changes",
  "context": {
    "workspace": "production",
    "dry_run": true
  }
}
```

Response:
```json
{
  "execution_id": "exec-123",
  "status": "running",
  "started_at": "2024-03-20T10:05:00Z"
}
```

## Error Responses

### Common Error Codes

- `400 Bad Request`: Invalid parameters
```json
{
  "error": "Invalid request",
  "message": "Missing required field: name",
  "code": "INVALID_REQUEST"
}
```

- `401 Unauthorized`: Missing or invalid API key
```json
{
  "error": "Unauthorized",
  "message": "Invalid API key",
  "code": "INVALID_AUTH"
}
```

- `404 Not Found`: Resource not found
```json
{
  "error": "Not found",
  "message": "Agent 'unknown-agent' not found",
  "code": "NOT_FOUND"
}
```

- `500 Internal Server Error`: Server error
```json
{
  "error": "Internal error",
  "message": "An unexpected error occurred",
  "code": "INTERNAL_ERROR"
}
```

## Rate Limiting

API requests are rate limited based on your authentication status:

- Authenticated users: 1000 requests per hour
- Unauthenticated users: 60 requests per hour

Rate limit headers are included in all responses:
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1616201400
```

## Webhooks

SentinelStacks can send webhook notifications for various events:

- Agent status changes
- Execution completion
- Error notifications

Webhook payload example:
```json
{
  "event": "agent.execution.completed",
  "agent": "terraform-agent",
  "execution_id": "exec-123",
  "status": "success",
  "timestamp": "2024-03-20T10:10:00Z",
  "data": {
    "duration": 300,
    "output": "Infrastructure changes applied successfully"
  }
}
```

## SDK Support

Official SDKs are available for:
- Python: [sentinelstacks-python](https://github.com/yourusername/sentinelstacks-python)
- Go: [sentinelstacks-go](https://github.com/yourusername/sentinelstacks-go)
- JavaScript: [sentinelstacks-js](https://github.com/yourusername/sentinelstacks-js)

For more information about using the SDKs, see the [SDK Documentation](../developer-guide/sdks.md). 