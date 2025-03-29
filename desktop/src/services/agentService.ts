import { Agent } from '../types/Agent';
import { invoke } from '@tauri-apps/api/tauri';

/**
 * Fetches all agents from the backend
 */
export const getAgents = async (): Promise<Agent[]> => {
  try {
    // Call the Tauri command to get agents
    const agents = await invoke<Agent[]>('get_agents');
    return agents;
  } catch (error) {
    console.error('Error fetching agents:', error);
    // For development, return mock data
    if (process.env.NODE_ENV === 'development') {
      return getMockAgents();
    }
    throw error;
  }
};

/**
 * Fetches a single agent by ID
 */
export const getAgentById = async (id: string): Promise<Agent> => {
  try {
    // Call the Tauri command to get the agent
    const agent = await invoke<Agent>('get_agent_by_id', { id });
    return agent;
  } catch (error) {
    console.error(`Error fetching agent ${id}:`, error);
    // For development, return mock data
    if (process.env.NODE_ENV === 'development') {
      const mockAgents = getMockAgents();
      const agent = mockAgents.find(a => a.id === id);
      if (agent) {
        return agent;
      }
    }
    throw error;
  }
};

/**
 * Creates a new agent
 */
export const createAgent = async (agent: Omit<Agent, 'id' | 'createdAt' | 'updatedAt'>): Promise<Agent> => {
  try {
    // Call the Tauri command to create the agent
    const newAgent = await invoke<Agent>('create_agent', { agent });
    return newAgent;
  } catch (error) {
    console.error('Error creating agent:', error);
    throw error;
  }
};

/**
 * Updates an existing agent
 */
export const updateAgent = async (id: string, updates: Partial<Agent>): Promise<Agent> => {
  try {
    // Call the Tauri command to update the agent
    const updatedAgent = await invoke<Agent>('update_agent', { id, updates });
    return updatedAgent;
  } catch (error) {
    console.error(`Error updating agent ${id}:`, error);
    throw error;
  }
};

/**
 * Deletes an agent
 */
export const deleteAgent = async (id: string): Promise<void> => {
  try {
    // Call the Tauri command to delete the agent
    await invoke<void>('delete_agent', { id });
  } catch (error) {
    console.error(`Error deleting agent ${id}:`, error);
    throw error;
  }
};

/**
 * Starts an agent
 */
export const startAgent = async (id: string): Promise<void> => {
  try {
    // Call the Tauri command to start the agent
    await invoke<void>('start_agent', { id });
  } catch (error) {
    console.error(`Error starting agent ${id}:`, error);
    throw error;
  }
};

/**
 * Stops an agent
 */
export const stopAgent = async (id: string): Promise<void> => {
  try {
    // Call the Tauri command to stop the agent
    await invoke<void>('stop_agent', { id });
  } catch (error) {
    console.error(`Error stopping agent ${id}:`, error);
    throw error;
  }
};

/**
 * Get mock agents data for development
 */
const getMockAgents = (): Agent[] => {
  return [
    {
      id: '1',
      name: 'Infrastructure Monitor',
      description: 'Monitors cloud infrastructure and reports issues',
      status: 'active',
      model: {
        provider: 'openai',
        name: 'gpt-4',
        options: {
          temperature: 0.7,
          max_tokens: 2000
        }
      },
      tools: [
        {
          id: 'terraform',
          name: 'Terraform',
          description: 'Manages infrastructure as code',
          version: '1.0.0'
        },
        {
          id: 'aws',
          name: 'AWS',
          description: 'AWS cloud operations',
          version: '1.0.0'
        }
      ],
      capabilities: ['conversation', 'code', 'infrastructure'],
      memory: {
        persistence: true,
        vectorStorage: true,
        messageCount: 150,
        lastUpdated: '2025-03-29T10:15:30Z'
      },
      createdAt: '2025-02-15T08:30:00Z',
      updatedAt: '2025-03-29T10:15:30Z',
      lastActiveAt: '2025-03-29T10:15:30Z'
    },
    {
      id: '2',
      name: 'Security Auditor',
      description: 'Performs security checks and compliance audits',
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
        {
          id: 'security-scanner',
          name: 'Security Scanner',
          description: 'Scans for vulnerabilities',
          version: '1.2.0'
        }
      ],
      capabilities: ['conversation', 'security', 'compliance'],
      memory: {
        persistence: true,
        vectorStorage: true,
        messageCount: 75,
        lastUpdated: '2025-03-28T16:45:20Z'
      },
      createdAt: '2025-02-20T11:20:00Z',
      updatedAt: '2025-03-28T16:45:20Z',
      lastActiveAt: '2025-03-28T16:45:20Z'
    },
    {
      id: '3',
      name: 'Code Helper',
      description: 'Assists with coding tasks and code reviews',
      status: 'error',
      model: {
        provider: 'ollama',
        name: 'codellama:13b',
        options: {
          temperature: 0.2,
          max_tokens: 2048
        }
      },
      tools: [
        {
          id: 'code-analyzer',
          name: 'Code Analyzer',
          description: 'Analyzes code quality and suggests improvements',
          version: '0.9.5'
        },
        {
          id: 'git',
          name: 'Git',
          description: 'Git operations',
          version: '1.0.0'
        }
      ],
      capabilities: ['conversation', 'code', 'documentation'],
      memory: {
        persistence: false,
        vectorStorage: false,
        messageCount: 42,
        lastUpdated: '2025-03-29T09:12:15Z'
      },
      createdAt: '2025-03-01T14:10:00Z',
      updatedAt: '2025-03-29T09:12:15Z',
      lastActiveAt: '2025-03-29T09:12:15Z'
    }
  ];
};