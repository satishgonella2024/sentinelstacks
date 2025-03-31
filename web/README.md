# SentinelStacks Web UI

This directory contains a simple web-based user interface for the SentinelStacks API server. It provides a way to:

1. Manage agents (view, create, and interact with)
2. Browse available images
3. Chat with agents using WebSockets
4. Monitor agent events in real-time

## Getting Started

1. Start the SentinelStacks API server:
   ```
   sentinel api --cors
   ```

2. Open the `index.html` file in a web browser, or serve it using a simple HTTP server:
   ```
   python -m http.server
   ```

3. Navigate to http://localhost:8000 in your browser

## Features

### Authentication

The UI includes a login form that authenticates against the SentinelStacks API server. For the prototype, any credentials will work, but in a production environment, this would validate against actual user accounts.

### Agent Management

- View a list of running agents
- Create new agents from available images
- Delete agents

### Real-time Chat

The UI demonstrates the WebSocket chat capabilities:

- Connect to an agent's chat interface
- Send messages to the agent
- Receive streaming responses
- View events (thinking, processing, etc.)

### WebSocket Integration

The UI showcases both WebSocket endpoints:

- `/v1/agents/{id}/chat` - For chatting with an agent
- `/v1/agents/{id}/events` - For monitoring agent events

## Technical Details

This is a static HTML/CSS/JavaScript application that runs entirely in the browser. It uses:

- Bootstrap 5 for styling
- Native Fetch API for REST requests
- WebSockets for real-time communication
- Local Storage for token persistence

## Development

This is a prototype implementation. In a production environment, you would:

1. Use a modern JavaScript framework (React, Vue, etc.)
2. Implement proper error handling and retry logic
3. Add comprehensive test coverage
4. Enhance security features
5. Improve accessibility
6. Optimize for mobile devices

## Screenshots

(Add screenshots here once the application is finalized) 