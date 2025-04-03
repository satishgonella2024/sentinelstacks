# API Documentation with OpenAPI

This guide explains how to document, view, and use the SentinelStacks API using OpenAPI specification.

## Overview

SentinelStacks uses [OpenAPI 3.0](https://swagger.io/specification/) (formerly known as Swagger) to document its REST API interfaces. This provides:

1. **Interactive API Documentation**: Try out API calls directly from the browser
2. **Precise Schema Definitions**: Clear understanding of request/response formats
3. **Code Generation**: Generate client SDKs for multiple languages
4. **Automated Testing**: Validate API behavior against documentation

## API Documentation Interfaces

SentinelStacks provides three ways to interact with the API documentation:

### 1. Swagger UI

Swagger UI offers an interactive way to explore and test the API.

- **URL**: `http://localhost:8081/swagger/` (when running the API server)
- **Usage**: Click on endpoints to expand them, fill in parameters, and execute requests

### 2. ReDoc

ReDoc provides a more modern, responsive documentation view that's easier to read.

- **URL**: `http://localhost:8081/redoc/` (when running the API server)
- **Usage**: Browse the documentation with enhanced readability and search

### 3. Raw OpenAPI JSON

Access the raw OpenAPI specification:

- **URL**: `http://localhost:8081/swagger/doc.json`
- **Usage**: Import into API design tools or client generators

## Starting the API Server

To run the API server with documentation:

```bash
sentinel api
```

This starts:
- The main API server on port 8080
- The documentation server on port 8081

## API Server Options

```bash
# Run API server with custom port
sentinel api --port 9000

# Enable detailed request logging
sentinel api --log-requests

# Run API server on a specific network interface
sentinel api --host 0.0.0.0

# Run API server with TLS/HTTPS
sentinel api --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem
```

## Authenticating with the API

The API requires authentication for most endpoints:

1. **Get a token**:
   ```bash
   curl -X POST http://localhost:8080/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "user", "password": "password"}'
   ```

2. **Use the token in requests**:
   ```bash
   curl http://localhost:8080/v1/agents \
     -H "Authorization: Bearer your-token-here"
   ```

In the Swagger UI, you can click the "Authorize" button to enter your token.

## Core API Endpoints

The SentinelStacks API is organized into logical groups:

1. **Agents**: Create, manage, and interact with agents
2. **Images**: Work with agent images
3. **Stacks**: Manage multi-agent systems
4. **Registry**: Search, pull, and push agent packages
5. **Memory**: Access agent memory and vector storage

## For Developers: Documenting API Endpoints

When adding new API endpoints, use the Swaggo annotations in your handler functions:

```go
// @Summary Create a new agent
// @Description Create a new agent from an image
// @Tags agents
// @Accept json
// @Produce json
// @Param agent body AgentRequest true "Agent Request"
// @Success 201 {object} AgentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /agents [post]
func (s *Server) createAgentHandler(w http.ResponseWriter, r *http.Request) {
    // Implementation...
}
```

After adding annotations, generate the OpenAPI spec:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g internal/api/server.go -o docs/swagger
```

## Client Libraries

The OpenAPI specification can be used to generate client libraries in various languages:

```bash
# Generate TypeScript client
npx @openapitools/openapi-generator-cli generate \
  -i http://localhost:8081/swagger/doc.json \
  -g typescript-fetch \
  -o client/typescript

# Generate Python client
npx @openapitools/openapi-generator-cli generate \
  -i http://localhost:8081/swagger/doc.json \
  -g python \
  -o client/python
```

## API Versioning

The API is versioned with the `/v1` prefix. Future versions will use `/v2`, `/v3`, etc. 