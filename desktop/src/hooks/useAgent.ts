import { useState, useEffect } from 'react';
import toast from '../utils/toast';
import { Agent } from '../types/Agent';
import { getAgent, startAgent, stopAgent, deleteAgent } from '../services/agentService';

export function useAgent(id: string | undefined) {
  const [agent, setAgent] = useState<Agent | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isBusy, setIsBusy] = useState(false);

  useEffect(() => {
    if (id) {
      fetchAgent(id);
    }
  }, [id]);

  const fetchAgent = async (agentId: string) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const data = await getAgent(agentId);
      setAgent(data);
    } catch (err) {
      const errorMsg = `Failed to load agent: ${err instanceof Error ? err.message : String(err)}`;
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setIsLoading(false);
    }
  };

  const handleStartAgent = async () => {
    if (!agent) return;
    
    setIsBusy(true);
    try {
      await startAgent(agent.id);
      // Update agent status
      setAgent(prev => prev ? { ...prev, status: 'active' } : null);
      toast.success(`Agent ${agent.name} started successfully`);
      return true;
    } catch (err) {
      const errorMsg = `Failed to start agent: ${err instanceof Error ? err.message : String(err)}`;
      setError(errorMsg);
      toast.error(errorMsg);
      return false;
    } finally {
      setIsBusy(false);
    }
  };

  const handleStopAgent = async () => {
    if (!agent) return;
    
    setIsBusy(true);
    try {
      await stopAgent(agent.id);
      // Update agent status
      setAgent(prev => prev ? { ...prev, status: 'inactive' } : null);
      toast.success(`Agent ${agent.name} stopped successfully`);
      return true;
    } catch (err) {
      const errorMsg = `Failed to stop agent: ${err instanceof Error ? err.message : String(err)}`;
      setError(errorMsg);
      toast.error(errorMsg);
      return false;
    } finally {
      setIsBusy(false);
    }
  };

  const handleDeleteAgent = async () => {
    if (!agent) return false;
    
    setIsBusy(true);
    try {
      await deleteAgent(agent.id);
      toast.success(`Agent ${agent.name} deleted successfully`);
      return true;
    } catch (err) {
      const errorMsg = `Failed to delete agent: ${err instanceof Error ? err.message : String(err)}`;
      setError(errorMsg);
      toast.error(errorMsg);
      return false;
    } finally {
      setIsBusy(false);
    }
  };

  return {
    agent,
    isLoading,
    error,
    isBusy,
    fetchAgent,
    startAgent: handleStartAgent,
    stopAgent: handleStopAgent,
    deleteAgent: handleDeleteAgent
  };
}