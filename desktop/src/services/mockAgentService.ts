import { Agent } from '../types/Agent';
import { v4 as uuidv4 } from 'uuid'; 

// Mock data for development
const mockAgents: Agent[] = [
  {
    id: '1',
    name: 'Data Analysis Assistant',
    description: 'Analyzes data and provides insights, specialized in numerical analysis and visualization.',
    status: 'active',
    model: {
      provider: 'openai',
      name: 'gpt-4',
      options: {
        temperature: 0.7,
        max_tokens: 8000
      }
    },
    tools: [
      { id: 'tool1', name: 'Calculator', description: 'Performs complex calculations', version: '1.0.0' },
      { id: 'tool2', name: 'DataVisualizer', description: 'Creates charts and graphs', version: '1.0.0' }
    ],
    capabilities: ['conversation', 'data-analysis', 'visualization'],
    memory: {
      persistence: true,
      vectorStorage: true,
      messageCount: 24,
      lastUpdated: new Date().toISOString()
    },
    createdAt: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
    lastActiveAt: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString()
  },
  {
    id: '2',
    name: 'Customer Support Agent',
    description: 'Handles customer inquiries and provides support for product questions and issues.',
    status: 'inactive',
    model: {
      provider: 'claude',
      name: 'claude-3-opus',
      options: {
        temperature: 0.5,
        max_tokens: 4000
      }
    },
    tools: [
      { id: 'tool3', name: 'KnowledgeBase', description: 'Searches product documentation', version: '1.1.0' }
    ],
    capabilities: ['conversation', 'customer-support', 'product-knowledge'],
    memory: {
      persistence: true,
      vectorStorage: true,
      messageCount: 120,
      lastUpdated: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
    },
    createdAt: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(),
    lastActiveAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: '3',
    name: 'Code Assistant',
    description: 'Helps with programming tasks, code review, and debugging across multiple languages.',
    status: 'error',
    model: {
      provider: 'openai',
      name: 'gpt-4',
      options: {
        temperature: 0.2,
        max_tokens: 8000
      }
    },
    tools: [
      { id: 'tool4', name: 'GitIntegration', description: 'Interfaces with Git repositories', version: '0.9.0' },
      { id: 'tool5', name: 'CodeAnalyzer', description: 'Static code analysis', version: '1.2.0' }
    ],
    capabilities: ['conversation', 'code-generation', 'debugging', 'code-review'],
    memory: {
      persistence: false,
      vectorStorage: true,
      messageCount: 0,
      lastUpdated: new Date().toISOString()
    },
    createdAt: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    lastActiveAt: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString()
  }
];

// Demo conversation data
const mockConversation = {
  messages: [
    {
      id: '1',
      role: 'user',
      content: 'Hello, can you help me analyze some sales data?',
      timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString()
    },
    {
      id: '2', 
      role: 'assistant',
      content: 'I\'d be happy to help you analyze your sales data. What specific insights are you looking for?',
      timestamp: new Date(Date.now() - 29 * 60 * 1000).toISOString()
    },
    {
      id: '3',
      role: 'user',
      content: 'I want to understand which products performed best last quarter.',
      timestamp: new Date(Date.now() - 28 * 60 * 1000).toISOString()
    },
    {
      id: '4',
      role: 'assistant',
      content: 'To analyze your product performance for last quarter, I\'ll need your sales data. Could you upload a CSV or Excel file with your sales information? Ideally, it should include product names, sales volumes, revenue figures, and dates.',
      timestamp: new Date(Date.now() - 27 * 60 * 1000).toISOString()
    }
  ],
  metadata: {
    totalMessages: 4,
    lastMessageAt: new Date(Date.now() - 27 * 60 * 1000).toISOString()
  }
};

// Demo memory data
const mockMemory = {
  messages: [
    {
      role: 'user',
      content: 'Hello, can you help me analyze some sales data?',
      timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString()
    },
    {
      role: 'assistant',
      content: 'I\'d be happy to help you analyze your sales data. What specific insights are you looking for?',
      timestamp: new Date(Date.now() - 29 * 60 * 1000).toISOString()
    },
    {
      role: 'user',
      content: 'I want to understand which products performed best last quarter.',
      timestamp: new Date(Date.now() - 28 * 60 * 1000).toISOString()
    },
    {
      role: 'assistant',
      content: 'To analyze your product performance for last quarter, I\'ll need your sales data. Could you upload a CSV or Excel file with your sales information? Ideally, it should include product names, sales volumes, revenue figures, and dates.',
      timestamp: new Date(Date.now() - 27 * 60 * 1000).toISOString()
    }
  ],
  vectorStore: {
    size: 24,
    lastUpdated: new Date(Date.now() - 27 * 60 * 1000).toISOString()
  }
};

// Simulate API latency
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

/**
 * Fetches all agents
 */
export const getAgents = async (): Promise<Agent[]> => {
  await delay(800); // Simulate network delay
  return [...mockAgents];
};

/**
 * Fetches a single agent by ID
 */
export const getAgent = async (id: string): Promise<Agent> => {
  await delay(500);
  const agent = mockAgents.find(a => a.id === id);
  if (!agent) {
    throw new Error(`Agent with ID ${id} not found`);
  }
  return {...agent};
};

/**
 * Creates a new agent
 */
export const createAgent = async (agent: Omit<Agent, 'id'>): Promise<Agent> => {
  await delay(1000);
  const newAgent: Agent = {
    ...agent,
    id: uuidv4(),
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    status: 'inactive'
  };
  mockAgents.push(newAgent);
  return {...newAgent};
};

/**
 * Updates an existing agent
 */
export const updateAgent = async (id: string, agent: Partial<Agent>): Promise<Agent> => {
  await delay(800);
  const index = mockAgents.findIndex(a => a.id === id);
  if (index === -1) {
    throw new Error(`Agent with ID ${id} not found`);
  }
  
  mockAgents[index] = {
    ...mockAgents[index],
    ...agent,
    updatedAt: new Date().toISOString()
  };
  
  return {...mockAgents[index]};
};

/**
 * Deletes an agent
 */
export const deleteAgent = async (id: string): Promise<void> => {
  await delay(700);
  const index = mockAgents.findIndex(a => a.id === id);
  if (index === -1) {
    throw new Error(`Agent with ID ${id} not found`);
  }
  
  mockAgents.splice(index, 1);
};

/**
 * Starts an agent
 */
export const startAgent = async (id: string): Promise<void> => {
  await delay(1200);
  const agent = mockAgents.find(a => a.id === id);
  if (!agent) {
    throw new Error(`Agent with ID ${id} not found`);
  }
  
  agent.status = 'active';
  agent.lastActiveAt = new Date().toISOString();
};

/**
 * Stops an agent
 */
export const stopAgent = async (id: string): Promise<void> => {
  await delay(1000);
  const agent = mockAgents.find(a => a.id === id);
  if (!agent) {
    throw new Error(`Agent with ID ${id} not found`);
  }
  
  agent.status = 'inactive';
};

/**
 * Get agent memory
 */
export const getAgentMemory = async (id: string) => {
  await delay(600);
  return {...mockMemory};
};

/**
 * Get agent conversation
 */
export const getAgentConversation = async (id: string) => {
  await delay(600);
  return {...mockConversation};
};