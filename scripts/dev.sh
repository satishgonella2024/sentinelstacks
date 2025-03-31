#!/usr/bin/env bash

# Exit on error for most commands
set -e

# Function to handle cleanup on exit
cleanup() {
  echo "Shutting down services..."
  if [ ! -z "$BACKEND_PID" ]; then
    kill $BACKEND_PID 2>/dev/null || true
  fi
  if [ ! -z "$FRONTEND_PID" ]; then
    kill $FRONTEND_PID 2>/dev/null || true
  fi
  exit 0
}

# Register the cleanup function for when script exits
trap cleanup SIGINT SIGTERM EXIT

cd "$(dirname "$0")/.."
REPO_ROOT=$(pwd)

# Flag to track if we're using mock mode
USE_MOCK=false

# Try to build the backend command
echo "Building the sentinel CLI..."
cd "$REPO_ROOT"
if go build -o bin/sentinel cmd/sentinel/main/main.go; then
  echo "✅ Backend build successful"
  
  # Start the backend API server
  echo "Starting the backend API server..."
  cd "$REPO_ROOT"
  bin/sentinel api -p 8080 --cors &
  BACKEND_PID=$!

  # Check if backend started successfully
  sleep 2
  if ! ps -p $BACKEND_PID > /dev/null; then
    echo "⚠️ Failed to start backend server!"
    echo "Falling back to mock mode..."
    USE_MOCK=true
  else
    echo "✅ Backend API server running on http://localhost:8080"
  fi
else
  echo "⚠️ Backend build failed!"
  echo "Falling back to mock mode..."
  USE_MOCK=true
fi

# If we need to use mock mode, set up the mock environment variable
if [ "$USE_MOCK" = true ]; then
  echo "Setting up mock data mode..."
  
  # Create .env.local file to enable mock mode
  cd "$REPO_ROOT/web-ui"
  echo "VITE_USE_MOCK_API=true" > .env.local
  echo "✅ Created .env.local with mock mode enabled"
  
  # Check if we need to create the mock API handler file
  if [ ! -f "src/mocks/handlers.ts" ]; then
    echo "Creating mock API handlers..."
    
    # Create mocks directory if it doesn't exist
    mkdir -p src/mocks
    
    # Create handlers.ts file
    cat > src/mocks/handlers.ts << 'EOF'
import { http, HttpResponse } from 'msw'

// Types
export interface Agent {
  id: string;
  name: string;
  model: string;
  status: 'running' | 'stopped' | 'error';
  created_at: string;
  updated_at: string;
  description?: string;
  image?: string;
}

export interface Conversation {
  id: string;
  agent_id: string;
  title: string;
  created_at: string;
  updated_at: string;
  messages: Message[];
}

export interface Message {
  id: string;
  conversation_id: string;
  role: 'user' | 'assistant';
  content: string;
  created_at: string;
}

// Mock data
const mockAgents: Agent[] = [
  {
    id: '1',
    name: 'General Assistant',
    model: 'claude-3-opus',
    status: 'running',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    description: 'General purpose AI assistant',
    image: 'https://source.unsplash.com/random/300x300/?robot',
  },
  {
    id: '2',
    name: 'Code Assistant',
    model: 'gpt-4',
    status: 'running',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    description: 'Specialized in code assistance and debugging',
    image: 'https://source.unsplash.com/random/300x300/?code',
  },
  {
    id: '3',
    name: 'Image Analyzer',
    model: 'claude-3-sonnet',
    status: 'stopped',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    description: 'Visual analysis and image understanding',
    image: 'https://source.unsplash.com/random/300x300/?camera',
  },
];

const mockConversations: Record<string, Conversation[]> = {
  '1': [
    {
      id: '101',
      agent_id: '1',
      title: 'Project Planning Discussion',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      messages: [
        {
          id: '1001',
          conversation_id: '101',
          role: 'user',
          content: 'Can you help me plan a new software project?',
          created_at: new Date().toISOString(),
        },
        {
          id: '1002',
          conversation_id: '101',
          role: 'assistant',
          content: 'Of course! Let\'s start by defining the project scope and objectives. What kind of software are you looking to build?',
          created_at: new Date().toISOString(),
        },
      ],
    },
  ],
  '2': [
    {
      id: '201',
      agent_id: '2',
      title: 'Debugging React Component',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      messages: [
        {
          id: '2001',
          conversation_id: '201',
          role: 'user',
          content: 'I have a React component that\'s not rendering correctly. Here\'s the code...',
          created_at: new Date().toISOString(),
        },
        {
          id: '2002',
          conversation_id: '201',
          role: 'assistant',
          content: 'Looking at your code, I see a few issues. First, you\'re not properly managing state with useEffect...',
          created_at: new Date().toISOString(),
        },
      ],
    },
  ],
};

// API handlers
export const handlers = [
  // Authentication
  http.post('/api/auth/login', () => {
    return HttpResponse.json({
      token: 'mock-jwt-token',
      user: {
        id: '1',
        username: 'demo_user',
        email: 'demo@example.com',
      },
    })
  }),

  // Agents
  http.get('/api/agents', () => {
    return HttpResponse.json(mockAgents)
  }),

  http.get('/api/agents/:id', ({ params }) => {
    const { id } = params
    const agent = mockAgents.find(a => a.id === id) || mockAgents[0]
    return HttpResponse.json(agent)
  }),

  http.post('/api/agents', async ({ request }) => {
    const data = await request.json()
    const newAgent: Agent = {
      id: Math.random().toString(36).substring(2, 9),
      name: data.name,
      model: data.model,
      status: 'running',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      description: data.description || '',
    }
    mockAgents.push(newAgent)
    return HttpResponse.json(newAgent)
  }),

  // Conversations
  http.get('/api/agents/:agentId/conversations', ({ params }) => {
    const { agentId } = params
    return HttpResponse.json(mockConversations[agentId as string] || [])
  }),

  http.get('/api/agents/:agentId/conversations/:conversationId', ({ params }) => {
    const { agentId, conversationId } = params
    const conversations = mockConversations[agentId as string] || []
    const conversation = conversations.find(c => c.id === conversationId) || conversations[0]
    return HttpResponse.json(conversation)
  }),

  http.post('/api/agents/:agentId/conversations/:conversationId/messages', async ({ request, params }) => {
    const data = await request.json()
    const { agentId, conversationId } = params
    
    // Create user message
    const userMsg: Message = {
      id: Math.random().toString(36).substring(2, 9),
      conversation_id: conversationId as string,
      role: 'user',
      content: data.content,
      created_at: new Date().toISOString(),
    }
    
    // Create assistant response
    const assistantMsg: Message = {
      id: Math.random().toString(36).substring(2, 9),
      conversation_id: conversationId as string,
      role: 'assistant',
      content: 'This is a mock response from the assistant. In a real app, this would be generated by the AI model.',
      created_at: new Date(Date.now() + 1000).toISOString(),
    }
    
    // Add messages to conversation
    const conversations = mockConversations[agentId as string] || []
    const conversation = conversations.find(c => c.id === conversationId)
    
    if (conversation) {
      conversation.messages.push(userMsg)
      conversation.messages.push(assistantMsg)
    }
    
    return HttpResponse.json(userMsg)
  }),
]
EOF

    # Create browser.ts file
    cat > src/mocks/browser.ts << 'EOF'
import { setupWorker } from 'msw/browser'
import { handlers } from './handlers'

export const worker = setupWorker(...handlers)
EOF

    # Create node.ts file for testing environments
    cat > src/mocks/node.ts << 'EOF'
import { setupServer } from 'msw/node'
import { handlers } from './handlers'

export const server = setupServer(...handlers)
EOF

    # Update main.tsx to include MSW in development
    if [ -f "src/main.tsx" ]; then
      cat > src/main.tsx.new << 'EOF'
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'

async function bootstrap() {
  // Setup mock server in development
  if (import.meta.env.VITE_USE_MOCK_API === 'true') {
    console.log('🔶 Using mock API in development mode')
    const { worker } = await import('./mocks/browser')
    await worker.start({ 
      onUnhandledRequest: 'bypass' 
    })
  }

  ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  )
}

bootstrap()
EOF
      mv src/main.tsx.new src/main.tsx
    fi

    # Add MSW dependencies
    echo "Installing MSW (Mock Service Worker) for API mocking..."
    npm install msw --save-dev
    
    echo "✅ Mock API setup complete"
  else
    echo "✅ Mock API handlers already exist"
  fi
fi

# Start the frontend
echo "Starting the frontend development server..."
cd "$REPO_ROOT/web-ui"
npm run dev &
FRONTEND_PID=$!

# Check if frontend started successfully
sleep 3
if ! ps -p $FRONTEND_PID > /dev/null; then
  echo "❌ Failed to start frontend server!"
  exit 1
fi
echo "✅ Frontend server running on http://localhost:5173"

echo "==================================================="
echo "Both services are now running!"
if [ "$USE_MOCK" = true ]; then
  echo "⚠️ RUNNING IN MOCK MODE - Using mocked backend data"
  echo "- Frontend UI: http://localhost:5173"
else 
  echo "- Backend API: http://localhost:8080"
  echo "- Frontend UI: http://localhost:5173"
fi
echo "==================================================="
echo "Press Ctrl+C to stop all services."

# Wait for user to press Ctrl+C
wait 