import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { Agent, Message } from '../context/slices/agentsSlice'

import { API_CONFIG } from '../api-config'

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

// Mock data for development with enriched sample agents
const mockAgents: Agent[] = [
  {
    id: '1',
    name: 'Research Assistant',
    description: 'Helps with academic research, citation management, and research methodology.',
    model: 'claude-3-opus-20240229',
    image: 'anthropic/claude-3-opus:latest',
    status: 'active',
    created: new Date().toISOString(),
    lastActive: new Date().toISOString(),
    systemPrompt: 'You are a research assistant specializing in academic research. Help users find relevant papers, cite properly in different formats, and design research methodologies. Provide thoughtful analysis of academic papers and research findings.',
    isMultimodal: false,
    capabilities: ['Deep research', 'Citation formatting', 'Literature review', 'Methodology design'],
    tags: ['academic', 'research', 'papers']
  },
  {
    id: '2',
    name: 'Image Analyzer',
    description: 'Specialized in analyzing images and providing detailed visual descriptions',
    model: 'claude-3-opus-20240229',
    image: 'anthropic/claude-3-opus:latest',
    status: 'idle',
    created: new Date(Date.now() - 86400000).toISOString(),
    lastActive: new Date(Date.now() - 3600000).toISOString(),
    systemPrompt: 'You analyze images and provide detailed descriptions. Identify objects, people, settings, actions, text content, style elements, and other visual information. Maintain accuracy and avoid hallucinations.',
    isMultimodal: true,
    capabilities: ['Object detection', 'Scene description', 'Text recognition', 'Visual analysis'],
    tags: ['vision', 'images', 'multimodal']
  },
  {
    id: '3',
    name: 'Code Assistant',
    description: 'Specialized in helping with programming and code review',
    model: 'gpt-4',
    image: 'openai/gpt-4:latest',
    status: 'active',
    created: new Date(Date.now() - 172800000).toISOString(),
    lastActive: new Date(Date.now() - 86400000).toISOString(),
    systemPrompt: 'You are a code assistant. Help with programming tasks and code review. Provide working, well-documented code examples. Explain concepts clearly and help debug issues efficiently.',
    isMultimodal: false,
    capabilities: ['Code generation', 'Bug fixing', 'Code review', 'Algorithm design'],
    tags: ['coding', 'programming', 'development']
  },
  {
    id: '4',
    name: 'Data Analyst',
    description: 'Helps with data analysis, visualization, and statistical interpretation',
    model: 'claude-3-sonnet-20240229',
    image: 'anthropic/claude-3-sonnet:latest',
    status: 'idle',
    created: new Date(Date.now() - 259200000).toISOString(),
    lastActive: new Date(Date.now() - 172800000).toISOString(),
    systemPrompt: 'You are a data analysis specialist. Help users interpret data, create visualizations, perform statistical analysis, and generate insights from datasets. Explain statistical concepts clearly and provide actionable recommendations.',
    isMultimodal: false,
    capabilities: ['Statistical analysis', 'Data visualization', 'Trend detection', 'Insight generation'],
    tags: ['data', 'statistics', 'analysis']
  },
  {
    id: '5',
    name: 'Content Strategist',
    description: 'Creates and refines marketing content and communication strategies',
    model: 'gpt-4-turbo',
    image: 'openai/gpt-4-turbo:latest',
    status: 'error',
    created: new Date(Date.now() - 345600000).toISOString(),
    lastActive: new Date(Date.now() - 259200000).toISOString(),
    systemPrompt: 'You are a content strategist and copywriter. Help users create compelling marketing copy, social media content, email campaigns, and content strategies. Focus on clarity, engagement, and appropriate tone for the target audience.',
    isMultimodal: false,
    capabilities: ['Copywriting', 'Social media content', 'Email campaigns', 'Content planning'],
    tags: ['marketing', 'content', 'copywriting']
  },
  {
    id: '6',
    name: 'Legal Assistant',
    description: 'Provides information on legal concepts and contract analysis',
    model: 'llama-3-70b-instruct',
    image: 'meta/llama3:latest',
    status: 'active',
    created: new Date(Date.now() - 432000000).toISOString(),
    lastActive: new Date(Date.now() - 345600000).toISOString(),
    systemPrompt: 'You are a legal information assistant. Help users understand legal concepts, analyze contracts, and navigate legal documentation. Always clarify you are not providing legal advice and users should consult qualified attorneys for specific situations.',
    isMultimodal: false,
    capabilities: ['Contract review', 'Legal research', 'Document analysis', 'Legal education'],
    tags: ['legal', 'contracts', 'law']
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
  },
  {
    id: '2',
    agentId: '3',
    title: 'Debugging React Components',
    created: new Date(Date.now() - 172800000).toISOString(),
    updated: new Date(Date.now() - 86400000).toISOString(),
    messages: [
      {
        id: '3',
        role: 'user' as const,
        content: 'I have a React component that\'s not rendering correctly',
        timestamp: new Date(Date.now() - 172800000).toISOString(),
        agentId: '3'
      },
      {
        id: '4',
        role: 'assistant' as const,
        content: 'Let\'s troubleshoot that. Can you share the component code that\'s causing issues?',
        timestamp: new Date(Date.now() - 172790000).toISOString(),
        agentId: '3'
      }
    ]
  }
];

// Flag to control whether to use mock data or real API data
const USE_REAL_DATA = API_CONFIG.USE_REAL_DATA;

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

    getAgents: builder.query<{ agents: Agent[] }, void>({
      query: () => '/agents',
      transformResponse: (_response, meta) => {
        console.log('Agents API response:', _response, meta);
        
        // Handle server errors
        if (meta?.response && meta.response.status >= 500) {
          console.warn('Server error occurred in agents endpoint', meta.response);
          
          // Always fall back to mock data on server error
          console.log('Falling back to mock data due to server error');
          return { agents: mockAgents };
        }
        
        // Process successful response
        if (meta?.response?.ok && _response) {
          try {
            const responseData = _response as any;
            console.log('Processing successful response:', responseData);
            
            if (responseData.agents && Array.isArray(responseData.agents)) {
              return responseData as { agents: Agent[] };
            } else if (Array.isArray(responseData)) {
              return { agents: responseData as Agent[] };
            } else {
              console.warn('Unexpected response format:', responseData);
              return { agents: [] };
            }
          } catch (e) {
            console.error('Error parsing API response:', e);
            return { agents: [] };
          }
        }
        
        console.log('Unhandled response case - using mock data');
        return { agents: mockAgents };
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
            isMultimodal: arg.isMultimodal,
            capabilities: [],
            tags: []
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