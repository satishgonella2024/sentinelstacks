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
  http.post('/v1/auth/login', () => {
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
  http.get('/v1/agents', () => {
    return HttpResponse.json({
      agents: mockAgents
    })
  }),

  http.get('/v1/agents/:id', ({ params }) => {
    const { id } = params
    const agent = mockAgents.find(a => a.id === id) || mockAgents[0]
    return HttpResponse.json(agent)
  }),

  http.post('/v1/agents', async ({ request }) => {
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
  http.get('/v1/agents/:agentId/conversations', ({ params }) => {
    const { agentId } = params
    return HttpResponse.json(mockConversations[agentId as string] || [])
  }),

  http.get('/v1/agents/:agentId/conversations/:conversationId', ({ params }) => {
    const { agentId, conversationId } = params
    const conversations = mockConversations[agentId as string] || []
    const conversation = conversations.find(c => c.id === conversationId) || conversations[0]
    return HttpResponse.json(conversation)
  }),

  http.post('/v1/agents/:agentId/conversations/:conversationId/messages', async ({ request, params }) => {
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
