# API Reference

This document provides detailed information about the SentinelStacks API endpoints, their parameters, and expected responses.

## Base URL

All API endpoints are relative to:
```
https://your-domain/api
```

## Authentication

Most API endpoints require authentication. Include your API key in the request header:
```
Authorization: Bearer your-api-key
```

## Endpoints

### Agents

#### List Agents

```http
GET /agents
```

Lists all available agents in the registry.

**Parameters:**
- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `filter` (optional): Filter by agent capabilities

**Response:**
```json
{
  "agents": [
    {
      "name": "terraform-agent",
      "version": "latest",
      "description": "Agent for managing Terraform infrastructure",
      "capabilities": ["terraform", "aws", "azure", "gcp"],
      "commands": [
        {
          "name": "plan",
          "description": "Generate Terraform plan",
          "args": [
            {
              "name": "path",
              "type": "string",
              "required": true,
              "description": "Path to Terraform configuration"
            }
          ]
        }
      ]
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

#### Get Agent Details

```http
GET /agents/{name}
```

Get detailed information about a specific agent.

**Parameters:**
- `name` (required): Name of the agent
- `version` (optional): Specific version (default: latest)

**Response:**
```json
{
  "name": "terraform-agent",
  "version": "latest",
  "description": "Agent for managing Terraform infrastructure",
  "capabilities": ["terraform", "aws", "azure", "gcp"],
  "commands": [
    {
      "name": "plan",
      "description": "Generate Terraform plan",
      "args": [
        {
          "name": "path",
          "type": "string",
          "required": true,
          "description": "Path to Terraform configuration"
        }
      ]
    }
  ],
  "metadata": {
    "created_at": "2024-03-28T10:00:00Z",
    "updated_at": "2024-03-28T10:00:00Z",
    "downloads": 100
  }
}
```

#### Push Agent

```http
POST /agents
```

Push a new agent to the registry.

**Request Body:**
```json
{
  "name": "my-agent",
  "version": "1.0.0",
  "description": "My custom agent",
  "capabilities": ["custom"],
  "commands": [],
  "files": {
    "agent.yaml": "base64-encoded-content",
    "README.md": "base64-encoded-content"
  }
}
```

**Response:**
```json
{
  "id": "agent-id",
  "status": "success",
  "message": "Agent pushed successfully"
}
```

### Registry Operations

#### Search Registry

```http
GET /registry/search
```

Search for agents in the registry.

**Parameters:**
- `q` (required): Search query
- `tags` (optional): Comma-separated list of tags
- `sort` (optional): Sort order (downloads, updated)
- `order` (optional): Sort direction (asc, desc)

**Response:**
```json
{
  "results": [
    {
      "name": "matching-agent",
      "version": "1.0.0",
      "description": "Agent matching search criteria",
      "score": 0.95
    }
  ],
  "total": 1
}
```

### Agent Execution

#### Run Agent Command

```http
POST /agents/{name}/run
```

Execute a command on a specific agent.

**Request Body:**
```json
{
  "command": "plan",
  "args": {
    "path": "/path/to/config"
  },
  "options": {
    "timeout": 300
  }
}
```

**Response:**
```json
{
  "id": "execution-id",
  "status": "running",
  "output": "Command output stream"
}
```

### Error Responses

All endpoints may return the following error responses:

#### 400 Bad Request
```json
{
  "error": "validation_error",
  "message": "Invalid request parameters",
  "details": {
    "field": "error description"
  }
}
```

#### 401 Unauthorized
```json
{
  "error": "unauthorized",
  "message": "Invalid or missing API key"
}
```

#### 404 Not Found
```json
{
  "error": "not_found",
  "message": "Resource not found"
}
```

#### 500 Internal Server Error
```json
{
  "error": "internal_error",
  "message": "An unexpected error occurred"
}
```

## Rate Limiting

API requests are rate limited to:
- 100 requests per minute for authenticated users
- 10 requests per minute for unauthenticated users

Rate limit headers are included in all responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1585395600
```

## Webhooks

SentinelStacks supports webhooks for event notifications. Configure webhooks in your account settings.

### Event Types

- `agent.pushed`: New agent version pushed
- `agent.pulled`: Agent pulled by a user
- `execution.started`: Agent execution started
- `execution.completed`: Agent execution completed
- `execution.failed`: Agent execution failed

### Webhook Payload
```json
{
  "event": "agent.pushed",
  "timestamp": "2024-03-28T10:00:00Z",
  "data": {
    "agent": "agent-name",
    "version": "1.0.0",
    "user": "user-id"
  }
}
```

## SDK Support

Official SDKs are available for:
- Go
- Python
- JavaScript/TypeScript
- Java

See the [SDK documentation](./03-sdk-guide.md) for detailed usage instructions. 