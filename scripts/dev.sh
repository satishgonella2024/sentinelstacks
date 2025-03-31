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
  bin/sentinel api -p 8080 -e true &
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

# If we need to use mock mode, ensure the API service has mock data
if [ "$USE_MOCK" = true ]; then
  echo "Setting up mock data mode..."
  
  # Check if the mock data is already in place
  cd "$REPO_ROOT/web-ui/src/services"
  if ! grep -q "mockAgents" api.ts; then
    echo "Adding mock data to api.ts..."
    
    cat > api.ts.new << 'EOF'
// API service for SentinelStacks
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

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

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: string;
    username: string;
    email: string;
  };
}

export interface CreateAgentRequest {
  name: string;
  model: string;
  description?: string;
  system_prompt?: string;
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

// API definition
export const api = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({ 
    baseUrl: '/api',
    prepareHeaders: (headers, { getState }) => {
      // Get token from state
      const token = (getState() as any).auth?.token;
      
      // Add auth header if token exists
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
      
      return headers;
    },
  }),
  tagTypes: ['Agent', 'Conversation'],
  endpoints: (builder) => ({
    login: builder.mutation<LoginResponse, LoginRequest>({
      query: (credentials) => ({
        url: '/auth/login',
        method: 'POST',
        body: credentials,
      }),
      transformResponse: (response: LoginResponse) => {
        // In mock mode, return fake data
        return {
          token: 'mock-jwt-token',
          user: {
            id: '1',
            username: 'demo_user',
            email: 'demo@example.com',
          },
        };
      },
    }),
    
    getAgents: builder.query<Agent[], void>({
      query: () => '/agents',
      transformResponse: (response: Agent[]) => {
        // In mock mode, return fake data
        return mockAgents;
      },
      providesTags: ['Agent'],
    }),
    
    getAgent: builder.query<Agent, string>({
      query: (id) => `/agents/${id}`,
      transformResponse: (response: Agent, _meta, arg) => {
        // In mock mode, return fake data
        return mockAgents.find(a => a.id === arg) || mockAgents[0];
      },
      providesTags: (_result, _error, id) => [{ type: 'Agent', id }],
    }),
    
    createAgent: builder.mutation<Agent, CreateAgentRequest>({
      query: (agent) => ({
        url: '/agents',
        method: 'POST',
        body: agent,
      }),
      transformResponse: (response: Agent, _meta, arg) => {
        // In mock mode, create a new agent
        return {
          id: Math.random().toString(36).substring(2, 9),
          name: arg.name,
          model: arg.model,
          status: 'running',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          description: arg.description,
        };
      },
      invalidatesTags: ['Agent'],
    }),
    
    getConversations: builder.query<Conversation[], string>({
      query: (agentId) => `/agents/${agentId}/conversations`,
      transformResponse: (response: Conversation[], _meta, arg) => {
        // In mock mode, return fake data
        return mockConversations[arg] || [];
      },
      providesTags: (_result, _error, id) => [{ type: 'Conversation', id }],
    }),
    
    getConversation: builder.query<Conversation, { agentId: string, conversationId: string }>({
      query: ({ agentId, conversationId }) => `/agents/${agentId}/conversations/${conversationId}`,
      transformResponse: (response: Conversation, _meta, arg) => {
        // In mock mode, return fake data
        const conversations = mockConversations[arg.agentId] || [];
        return conversations.find(c => c.id === arg.conversationId) || conversations[0];
      },
      providesTags: (_result, _error, arg) => [{ type: 'Conversation', id: arg.conversationId }],
    }),
    
    sendMessage: builder.mutation<Message, { agentId: string, conversationId: string, content: string }>({
      query: ({ agentId, conversationId, content }) => ({
        url: `/agents/${agentId}/conversations/${conversationId}/messages`,
        method: 'POST',
        body: { content },
      }),
      transformResponse: (response: Message, _meta, arg) => {
        // In mock mode, create a fake message response
        const userMsg: Message = {
          id: Math.random().toString(36).substring(2, 9),
          conversation_id: arg.conversationId,
          role: 'user',
          content: arg.content,
          created_at: new Date().toISOString(),
        };
        
        // Mock assistant response after a delay
        setTimeout(() => {
          // This would add the assistant message to the state in a real app
          console.log('Assistant would respond here in a real app');
        }, 1000);
        
        return userMsg;
      },
      invalidatesTags: (_result, _error, arg) => [{ type: 'Conversation', id: arg.conversationId }],
    }),
  }),
});

// Export hooks
export const {
  useLoginMutation,
  useGetAgentsQuery,
  useGetAgentQuery,
  useCreateAgentMutation,
  useGetConversationsQuery,
  useGetConversationQuery,
  useSendMessageMutation,
} = api;
EOF

    mv api.ts.new api.ts
    echo "✅ Mock data added to API service"
  else
    echo "✅ Mock data already present in API service"
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