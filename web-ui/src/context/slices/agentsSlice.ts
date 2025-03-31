import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Agent {
  id: string
  name: string
  description: string
  model: string
  image: string
  status: 'idle' | 'active' | 'error'
  created: string
  lastActive: string
  systemPrompt: string
  isMultimodal: boolean
  capabilities?: string[]
  tags?: string[]
}

export interface Message {
  id: string
  role: 'system' | 'user' | 'assistant'
  content: string | Array<{
    type: 'text' | 'image'
    content: string
  }>
  timestamp: string
  agentId: string
}

interface Conversation {
  id: string
  agentId: string
  title: string
  messages: Message[]
  created: string
  updated: string
}

interface AgentsState {
  agents: Agent[]
  selectedAgentId: string | null
  conversations: Conversation[]
  selectedConversationId: string | null
  isLoading: boolean
  error: string | null
}

const initialState: AgentsState = {
  agents: [],
  selectedAgentId: null,
  conversations: [],
  selectedConversationId: null,
  isLoading: false,
  error: null
}

const agentsSlice = createSlice({
  name: 'agents',
  initialState,
  reducers: {
    setAgents: (state, action: PayloadAction<Agent[]>) => {
      state.agents = action.payload
    },
    selectAgent: (state, action: PayloadAction<string>) => {
      state.selectedAgentId = action.payload
    },
    setConversations: (state, action: PayloadAction<Conversation[]>) => {
      state.conversations = action.payload
    },
    selectConversation: (state, action: PayloadAction<string>) => {
      state.selectedConversationId = action.payload
    },
    addMessage: (state, action: PayloadAction<Message>) => {
      const conversation = state.conversations.find(
        c => c.id === state.selectedConversationId
      )
      if (conversation) {
        conversation.messages.push(action.payload)
        conversation.updated = new Date().toISOString()
      }
    },
    createConversation: (state, action: PayloadAction<Omit<Conversation, 'messages' | 'created' | 'updated'>>) => {
      const newConversation: Conversation = {
        ...action.payload,
        messages: [],
        created: new Date().toISOString(),
        updated: new Date().toISOString()
      }
      state.conversations.push(newConversation)
      state.selectedConversationId = newConversation.id
    }
  }
})

export const {
  setAgents,
  selectAgent,
  setConversations,
  selectConversation,
  addMessage,
  createConversation
} = agentsSlice.actions

export default agentsSlice.reducer