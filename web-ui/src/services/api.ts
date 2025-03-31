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

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api',
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
    }),

    getAgents: builder.query<Agent[], void>({
      query: () => '/agents',
      providesTags: ['Agents'],
    }),

    getAgent: builder.query<Agent, string>({
      query: (id) => `/agents/${id}`,
      providesTags: (result, error, id) => [{ type: 'Agents', id }],
    }),

    createAgent: builder.mutation<Agent, CreateAgentRequest>({
      query: (agent) => ({
        url: '/agents',
        method: 'POST',
        body: agent,
      }),
      invalidatesTags: ['Agents'],
    }),

    getConversations: builder.query<
      { id: string; agentId: string; title: string; created: string; updated: string }[],
      string
    >({
      query: (agentId) => `/agents/${agentId}/conversations`,
      providesTags: ['Conversations'],
    }),

    getConversation: builder.query<
      { id: string; agentId: string; title: string; messages: Message[]; created: string; updated: string },
      { agentId: string; conversationId: string }
    >({
      query: ({ agentId, conversationId }) => `/agents/${agentId}/conversations/${conversationId}`,
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