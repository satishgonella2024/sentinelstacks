import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { Agent, Message } from '../context/slices/agentsSlice'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    id: string
    username: string
    email: string
    role: 'admin' | 'user'
  }
}

export interface CreateAgentRequest {
  name: string
  description: string
  model: string
  systemPrompt: string
  isMultimodal: boolean
}

// Mock data for development
const mockAgents: Agent[] = [
  {
    id: '1',
    name: 'Assistant Bot',
    description: 'General purpose assistant for everyday tasks',
    model: 'gpt-4',
    image: 'openai/gpt-4:latest',
    status: 'active',
    created: new Date().toISOString(),
    lastActive: new Date().toISOString(),
    systemPrompt: 'You are a helpful assistant.',
    isMultimodal: false
  },
  {
    id: '2',
    name: 'Image Analyzer',
    description: 'Specialized in analyzing images and providing descriptions',
    model: 'claude-3-opus-20240229',
    image: 'anthropic/claude-3-opus:latest',
    status: 'idle',
    created: new Date(Date.now() - 86400000).toISOString(),
    lastActive: new Date(Date.now() - 3600000).toISOString(),
    systemPrompt: 'You analyze images and provide detailed descriptions.',
    isMultimodal: true
  },
  {
    id: '3',
    name: 'Code Assistant',
    description: 'Specialized in helping with programming and code review',
    model: 'llama-3-70b-instruct',
    image: 'meta/llama3:latest',
    status: 'idle',
    created: new Date(Date.now() - 172800000).toISOString(),
    lastActive: new Date(Date.now() - 86400000).toISOString(),
    systemPrompt: 'You are a code assistant. Help with programming tasks and code review.',
    isMultimodal: false
  }
];

const mockConversations = [
  {
    id: '1',
    agentId: '1',
    title: 'Project Planning Discussion',
    created: new Date(Date.now() - 86400000).toISOString(),
    updated: new Date(Date.now() - 3600000).toISOString(),
    messages: [
      {
        id: '1',
        role: 'user' as const,
        content: 'I need help planning my project timeline',
        timestamp: new Date(Date.now() - 86400000).toISOString(),
        agentId: '1'
      },
      {
        id: '2',
        role: 'assistant' as const,
        content: 'I can help with that. What kind of project are you working on?',
        timestamp: new Date(Date.now() - 86390000).toISOString(),
        agentId: '1'
      }
    ]
  }
];

// Flag to control whether to use mock data or real API data
const USE_REAL_DATA = true;

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: '/v1',
    prepareHeaders: (headers, { getState }) => {
      // @ts-ignore - we'll properly type this later
      const token = getState().auth.token
      if (token) {
        headers.set('authorization', `Bearer ${token}`)
      }
      return headers
    },
  }),
  tagTypes: ['Agents', 'Conversations'],
  endpoints: (builder) => ({
    login: builder.mutation<LoginResponse, LoginRequest>({
      query: (credentials) => ({
        url: '/auth/login',
        method: 'POST',
        body: credentials,
      }),
      transformResponse: (_response, meta) => {
        console.log('Login API response:', _response, meta);
        
        if (USE_REAL_DATA && meta?.response?.ok && _response) {
          try {
            return _response as LoginResponse;
          } catch (e) {
            console.error('Error parsing login response:', e);
          }
        }
        
        // Fallback to mock data if needed
        return {
          token: 'mock-jwt-token',
          user: {
            id: 'user-1',
            username: 'admin',
            email: 'admin@sentinelstacks.com',
            role: 'admin'
          }
        };
      }
    }),

    getAgents: builder.query<Agent[], void>({
      query: () => '/agents',
      transformResponse: (_response, meta) => {
        console.log('Agents API response:', _response, meta);
        
        // Handle server errors
        if (meta?.response && meta.response.status >= 500) {
          console.warn('Server error occurred in agents endpoint', meta.response);
          
          // Always fall back to mock data on server error
          console.log('Falling back to mock data due to server error');
          return mockAgents;
        }
        
        // Process successful response
        if (meta?.response?.ok && _response) {
          try {
            const responseData = _response as any;
            console.log('Processing successful response:', responseData);
            
            if (responseData.agents && Array.isArray(responseData.agents)) {
              return responseData.agents as Agent[];
            } else if (Array.isArray(responseData)) {
              return responseData as Agent[];
            } else {
              console.warn('Unexpected response format:', responseData);
              return [];
            }
          } catch (e) {
            console.error('Error parsing API response:', e);
            return [];
          }
        }
        
        console.log('Unhandled response case - using mock data');
        return mockAgents;
      },
      // Add retry logic for the agents endpoint
      extraOptions: {
        maxRetries: 3
      },
      providesTags: ['Agents'],
    }),

    getAgent: builder.query<Agent, string>({
      query: (id) => `/agents/${id}`,
      transformResponse: (_response, meta, arg) => {
        console.log(`Agent ${arg} API response:`, _response, meta);
        
        if (USE_REAL_DATA && meta?.response?.ok && _response) {
          try {
            return _response as Agent;
          } catch (e) {
            console.error('Error parsing API response:', e);
          }
        }
        
        // Fallback to mock data if needed
        if (!USE_REAL_DATA) {
          console.log('Using mock agent data');
          return mockAgents.find(a => a.id === arg) || mockAgents[0];
        }
        
        throw new Error(`Agent with ID ${arg} not found`);
      },
      providesTags: (result, error, id) => [{ type: 'Agents', id }],
    }),

    createAgent: builder.mutation<Agent, CreateAgentRequest>({
      query: (agent) => ({
        url: '/agents',
        method: 'POST',
        body: agent,
      }),
      transformResponse: (_response, meta, arg) => {
        console.log('Create agent API response:', _response, meta);
        
        if (USE_REAL_DATA && meta?.response?.ok && _response) {
          try {
            return _response as Agent;
          } catch (e) {
            console.error('Error parsing API response:', e);
          }
        }
        
        // Fallback to mock data if needed
        if (!USE_REAL_DATA) {
          console.log('Using mock create agent response');
          return {
            id: Math.random().toString(36).substring(2, 9),
            name: arg.name,
            description: arg.description,
            model: arg.model,
            image: `default/${arg.model}:latest`,
            status: 'idle',
            created: new Date().toISOString(),
            lastActive: new Date().toISOString(),
            systemPrompt: arg.systemPrompt,
            isMultimodal: arg.isMultimodal
          };
        }
        
        throw new Error('Failed to create agent');
      },
      invalidatesTags: ['Agents'],
    }),

    getConversations: builder.query<
      { id: string; agentId: string; title: string; created: string; updated: string }[],
      string
    >({
      query: (agentId) => `/agents/${agentId}/conversations`,
      transformResponse: (_response, meta, arg) => {
        console.log('Conversations API response:', _response, meta);
        
        if (USE_REAL_DATA && meta?.response?.ok && _response) {
          try {
            return _response as { id: string; agentId: string; title: string; created: string; updated: string }[];
          } catch (e) {
            console.error('Error parsing API response:', e);
          }
        }
        
        // Fallback to mock data if needed
        if (!USE_REAL_DATA) {
          console.log('Using mock conversations data');
          return mockConversations
            .filter(c => c.agentId === arg)
            .map(({ id, agentId, title, created, updated }) => ({ id, agentId, title, created, updated }));
        }
        
        return [];
      },
      providesTags: ['Conversations'],
    }),

    getConversation: builder.query<
      { id: string; agentId: string; title: string; messages: Message[]; created: string; updated: string },
      { agentId: string; conversationId: string }
    >({
      query: ({ agentId, conversationId }) => `/agents/${agentId}/conversations/${conversationId}`,
      transformResponse: (_response, meta, arg) => {
        console.log('Conversation API response:', _response, meta);
        
        if (USE_REAL_DATA && meta?.response?.ok && _response) {
          try {
            return _response as { id: string; agentId: string; title: string; messages: Message[]; created: string; updated: string };
          } catch (e) {
            console.error('Error parsing API response:', e);
          }
        }
        
        // Fallback to mock data if needed
        if (!USE_REAL_DATA) {
          console.log('Using mock conversation data');
          const convo = mockConversations.find(c => c.agentId === arg.agentId && c.id === arg.conversationId);
          return convo || mockConversations[0];
        }
        
        throw new Error(`Conversation not found for agent ${arg.agentId}`);
      },
      providesTags: (result, error, arg) => [{ type: 'Conversations', id: arg.conversationId }],
    }),

    sendMessage: builder.mutation<
      Message,
      { agentId: string; conversationId: string; message: Omit<Message, 'id' | 'timestamp'> }
    >({
      query: ({ agentId, conversationId, message }) => ({
        url: `/agents/${agentId}/conversations/${conversationId}/messages`,
        method: 'POST',
        body: message,
      }),
      transformResponse: (_response, meta, arg) => {
        console.log('Send message API response:', _response, meta);
        
        if (USE_REAL_DATA && meta?.response?.ok && _response) {
          try {
            return _response as Message;
          } catch (e) {
            console.error('Error parsing API response:', e);
          }
        }
        
        // Fallback to mock data if needed
        if (!USE_REAL_DATA) {
          console.log('Using mock send message response');
          const userMsg: Message = {
            id: Math.random().toString(36).substring(2, 9),
            role: arg.message.role,
            content: arg.message.content,
            timestamp: new Date().toISOString(),
            agentId: arg.agentId
          };
          
          // Simulate response after a short delay in mock mode
          setTimeout(() => {
            const botMsg: Message = {
              id: Math.random().toString(36).substring(2, 9),
              role: 'assistant',
              content: `This is a simulated response to: "${typeof arg.message.content === 'string' 
                ? arg.message.content 
                : 'your message'}"`,
              timestamp: new Date().toISOString(),
              agentId: arg.agentId
            };
          }, 500);
          
          return userMsg;
        }
        
        throw new Error('Failed to send message');
      },
      invalidatesTags: (result, error, arg) => [{ type: 'Conversations', id: arg.conversationId }],
    }),
  }),
})

export const {
  useLoginMutation,
  useGetAgentsQuery,
  useGetAgentQuery,
  useCreateAgentMutation,
  useGetConversationsQuery,
  useGetConversationQuery,
  useSendMessageMutation,
} = apiSlice

export const api = apiSlice; 