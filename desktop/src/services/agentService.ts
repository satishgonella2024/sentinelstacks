import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig } from 'axios';
import { Agent } from '../types/Agent';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const API_VERSION = import.meta.env.VITE_API_VERSION;

const api: AxiosInstance = axios.create({
  baseURL: `${API_BASE_URL}/api/${API_VERSION}`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor for authentication
api.interceptors.request.use((config: AxiosRequestConfig) => {
  const token = localStorage.getItem('auth_token');
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error?.response?.status === 401) {
      // Handle unauthorized access
      localStorage.removeItem('auth_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

/**
 * Fetches all agents from the backend
 */
export const getAgents = async (): Promise<Agent[]> => {
  try {
    const response = await api.get<Agent[]>('/agents');
    return response.data;
  } catch (error) {
    console.error('Error fetching agents:', error);
    throw error;
  }
};

/**
 * Fetches a single agent by ID
 */
export const getAgent = async (id: string): Promise<Agent> => {
  try {
    const response = await api.get<Agent>(`/agents/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching agent ${id}:`, error);
    throw error;
  }
};

/**
 * Creates a new agent
 */
export const createAgent = async (agent: Omit<Agent, 'id'>): Promise<Agent> => {
  try {
    const response = await api.post<Agent>('/agents', agent);
    return response.data;
  } catch (error) {
    console.error('Error creating agent:', error);
    throw error;
  }
};

/**
 * Updates an existing agent
 */
export const updateAgent = async (id: string, agent: Partial<Agent>): Promise<Agent> => {
  try {
    const response = await api.put<Agent>(`/agents/${id}`, agent);
    return response.data;
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
    await api.delete(`/agents/${id}`);
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
    await api.post(`/agents/${id}/start`);
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
    await api.post(`/agents/${id}/stop`);
  } catch (error) {
    console.error(`Error stopping agent ${id}:`, error);
    throw error;
  }
};

interface AgentMemory {
  messages: Array<{
    role: string;
    content: string;
    timestamp: string;
  }>;
  vectorStore: {
    size: number;
    lastUpdated: string;
  };
}

interface AgentConversation {
  messages: Array<{
    id: string;
    role: string;
    content: string;
    timestamp: string;
  }>;
  metadata: {
    totalMessages: number;
    lastMessageAt: string;
  };
}

export const getAgentMemory = async (id: string): Promise<AgentMemory> => {
  try {
    const response = await api.get<AgentMemory>(`/agents/${id}/memory`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching agent ${id} memory:`, error);
    throw error;
  }
};

export const getAgentConversation = async (id: string): Promise<AgentConversation> => {
  try {
    const response = await api.get<AgentConversation>(`/agents/${id}/conversation`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching agent ${id} conversation:`, error);
    throw error;
  }
};