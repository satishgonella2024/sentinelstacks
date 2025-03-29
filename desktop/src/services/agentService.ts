// During development, we'll use the mock service 
// In production, we would use the real service implementation

// Reexport everything from the mock service for development
export * from './mockAgentService';

import { Agent } from '../types/Agent';
import { invoke } from '@tauri-apps/api/tauri';
import * as mockService from './mockAgentService';

// Environment check for development mode
const isDev = process.env.NODE_ENV === 'development';
const useMockService = isDev && process.env.REACT_APP_USE_MOCK === 'true';

/**
 * Real implementation of the agent service that communicates with the Tauri backend
 */
class RealAgentService {
  async getAgents(): Promise<Agent[]> {
    try {
      return await invoke<Agent[]>('get_agents');
    } catch (error) {
      console.error('Error fetching agents:', error);
      throw error;
    }
  }

  async getAgent(id: string): Promise<Agent> {
    try {
      return await invoke<Agent>('get_agent', { id });
    } catch (error) {
      console.error('Error fetching agent:', error);
      throw error;
    }
  }

  async createAgent(agent: Omit<Agent, 'id'>): Promise<Agent> {
    try {
      return await invoke<Agent>('create_agent', { agent });
    } catch (error) {
      console.error('Error creating agent:', error);
      throw error;
    }
  }

  async updateAgent(id: string, agent: Partial<Agent>): Promise<Agent> {
    try {
      return await invoke<Agent>('update_agent', { id, agent });
    } catch (error) {
      console.error('Error updating agent:', error);
      throw error;
    }
  }

  async deleteAgent(id: string): Promise<void> {
    try {
      await invoke('delete_agent', { id });
    } catch (error) {
      console.error('Error deleting agent:', error);
      throw error;
    }
  }

  async startAgent(id: string): Promise<void> {
    try {
      await invoke('start_agent', { id });
    } catch (error) {
      console.error('Error starting agent:', error);
      throw error;
    }
  }

  async stopAgent(id: string): Promise<void> {
    try {
      await invoke('stop_agent', { id });
    } catch (error) {
      console.error('Error stopping agent:', error);
      throw error;
    }
  }
}

// Create service instances
const realService = new RealAgentService();
const mockAgentService = mockService as typeof realService;

// Export the appropriate service based on environment
export const {
  getAgents,
  getAgent,
  createAgent,
  updateAgent,
  deleteAgent,
  startAgent,
  stopAgent,
} = useMockService ? mockAgentService : realService;

// Also export the service class for direct usage if needed
export const AgentService = useMockService ? mockAgentService.constructor : RealAgentService;