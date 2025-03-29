export type AgentStatus = 'active' | 'inactive' | 'error';

export type ModelProvider = 'openai' | 'claude' | 'ollama';

export interface Tool {
  id: string;
  name: string;
  description: string;
  version: string;
}

export interface ModelConfig {
  provider: ModelProvider;
  name: string;
  endpoint?: string;
  options: {
    temperature?: number;
    top_p?: number;
    max_tokens?: number;
    [key: string]: any;
  };
}

export interface AgentMemory {
  persistence: boolean;
  vectorStorage: boolean;
  messageCount: number;
  lastUpdated: string;
}

export interface Agent {
  id: string;
  name: string;
  description: string;
  status: AgentStatus;
  model: ModelConfig;
  tools: Tool[];
  capabilities: string[];
  memory: AgentMemory;
  createdAt: string;
  updatedAt: string;
  lastActiveAt?: string;
}