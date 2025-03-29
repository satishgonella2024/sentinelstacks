import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { 
  ArrowLeftIcon,
  PlayIcon, 
  StopIcon,
  PencilIcon,
  TrashIcon,
  ExclamationTriangleIcon,
  ClipboardDocumentIcon,
  ChevronDownIcon,
  ArrowPathIcon,
  CodeBracketIcon,
  ChatBubbleBottomCenterTextIcon,
  ServerIcon,
  KeyIcon
} from '@heroicons/react/24/outline';
import { Agent, AgentStatus } from '../types/Agent';
import { getAgentById, startAgent, stopAgent, deleteAgent } from '../services/agentService';

const AgentDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [agent, setAgent] = useState<Agent | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'conversation' | 'memory' | 'settings'>('conversation');
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [runningAgentOperation, setRunningAgentOperation] = useState(false);

  useEffect(() => {
    if (id) {
      fetchAgent(id);
    }
  }, [id]);

  const fetchAgent = async (agentId: string) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const data = await getAgentById(agentId);
      setAgent(data);
    } catch (err) {
      setError(`Failed to load agent: ${err instanceof Error ? err.message : String(err)}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleStartAgent = async () => {
    if (!agent) return;
    
    setRunningAgentOperation(true);
    try {
      await startAgent(agent.id);
      // Update agent status
      setAgent(prev => prev ? { ...prev, status: 'active' } : null);
    } catch (err) {
      setError(`Failed to start agent: ${err instanceof Error ? err.message : String(err)}`);
    } finally {
      setRunningAgentOperation(false);
    }
  };

  const handleStopAgent = async () => {
    if (!agent) return;
    
    setRunningAgentOperation(true);
    try {
      await stopAgent(agent.id);
      // Update agent status
      setAgent(prev => prev ? { ...prev, status: 'inactive' } : null);
    } catch (err) {
      setError(`Failed to stop agent: ${err instanceof Error ? err.message : String(err)}`);
    } finally {
      setRunningAgentOperation(false);
    }
  };

  const handleDeleteAgent = async () => {
    if (!agent) return;
    
    try {
      await deleteAgent(agent.id);
      navigate('/agents');
    } catch (err) {
      setError(`Failed to delete agent: ${err instanceof Error ? err.message : String(err)}`);
      setShowDeleteConfirm(false);
    }
  };

  const getStatusBadge = (status: AgentStatus) => {
    switch (status) {
      case 'active':
        return (
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
            <span className="h-2 w-2 rounded-full bg-green-400 mr-1.5" />
            Active
          </span>
        );
      case 'inactive':
        return (
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200">
            <span className="h-2 w-2 rounded-full bg-gray-400 mr-1.5" />
            Inactive
          </span>
        );
      case 'error':
        return (
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200">
            <ExclamationTriangleIcon className="h-3 w-3 mr-1" />
            Error
          </span>
        );
      default:
        return null;
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <ArrowPathIcon className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-screen">
        <div className="text-red-500 mb-4">
          <ExclamationTriangleIcon className="w-12 h-12" />
        </div>
        <h2 className="text-xl font-semibold mb-2">Error</h2>
        <p className="text-gray-600 dark:text-gray-400 mb-4">{error}</p>
        <Link
          to="/agents"
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
        >
          Back to Agents
        </Link>
      </div>
    );
  }

  if (!agent) {
    return (
      <div className="flex flex-col items-center justify-center h-screen">
        <h2 className="text-xl font-semibold mb-4">Agent not found</h2>
        <Link
          to="/agents"
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
        >
          Back to Agents
        </Link>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-6">
      {/* Header with back button, title and actions */}
      <div className="flex flex-col md:flex-row md:items-center justify-between mb-6">
        <div className="flex items-center mb-4 md:mb-0">
          <Link to="/agents" className="mr-4 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300">
            <ArrowLeftIcon className="w-5 h-5" />
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{agent.name}</h1>
            <div className="flex items-center mt-1">
              {getStatusBadge(agent.status)}
              <span className="ml-2 text-sm text-gray-500 dark:text-gray-400">
                {agent.model.provider} / {agent.model.name}
              </span>
            </div>
          </div>
        </div>
        
        <div className="flex items-center space-x-2">
          {agent.status === 'active' ? (
            <button
              onClick={handleStopAgent}
              disabled={runningAgentOperation}
              className="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
            >
              {runningAgentOperation ? (
                <ArrowPathIcon className="animate-spin -ml-0.5 mr-2 h-4 w-4" />
              ) : (
                <StopIcon className="-ml-0.5 mr-2 h-4 w-4" />
              )}
              Stop Agent
            </button>
          ) : (
            <button
              onClick={handleStartAgent}
              disabled={runningAgentOperation}
              className="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50"
            >
              {runningAgentOperation ? (
                <ArrowPathIcon className="animate-spin -ml-0.5 mr-2 h-4 w-4" />
              ) : (
                <PlayIcon className="-ml-0.5 mr-2 h-4 w-4" />
              )}
              Start Agent
            </button>
          )}
          
          <Link
            to={`/agents/${agent.id}/edit`}
            className="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            <PencilIcon className="-ml-0.5 mr-2 h-4 w-4" />
            Edit
          </Link>
          
          <button
            onClick={() => setShowDeleteConfirm(true)}
            className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 dark:bg-gray-800 dark:text-gray-200 dark:border-gray-600 dark:hover:bg-gray-700"
          >
            <TrashIcon className="-ml-0.5 mr-2 h-4 w-4" />
            Delete
          </button>
        </div>
      </div>
      
      {/* Agent description */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6 mb-6">
        <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Description</h2>
        <p className="text-gray-600 dark:text-gray-400">{agent.description}</p>
      </div>
      
      {/* Tabs */}
      <div className="border-b border-gray-200 dark:border-gray-700 mb-6">
        <nav className="-mb-px flex space-x-8">
          <button
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'conversation'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
            onClick={() => setActiveTab('conversation')}
          >
            <ChatBubbleBottomCenterTextIcon className="w-5 h-5 inline mr-2" />
            Conversation
          </button>
          <button
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'memory'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
            onClick={() => setActiveTab('memory')}
          >
            <ServerIcon className="w-5 h-5 inline mr-2" />
            Memory
          </button>
          <button
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'settings'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
            onClick={() => setActiveTab('settings')}
          >
            <KeyIcon className="w-5 h-5 inline mr-2" />
            Settings
          </button>
        </nav>
      </div>
      
      {/* Tab content */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
        {activeTab === 'conversation' && (
          <div className="p-6">
            <div className="flex flex-col h-96 border rounded-lg">
              <div className="flex-grow p-4 overflow-y-auto">
                {/* Conversation history would go here */}
                <div className="text-center text-gray-500 dark:text-gray-400 my-8">
                  No conversation history yet.
                </div>
              </div>
              <div className="border-t p-4">
                <div className="flex">
                  <input
                    type="text"
                    className="flex-grow block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600"
                    placeholder="Type a message..."
                  />
                  <button
                    className="ml-2 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                  >
                    Send
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
        
        {activeTab === 'memory' && (
          <div className="p-6">
            <div className="mb-6">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Memory Configuration</h3>
              <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-md">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Persistence</p>
                    <p className="mt-1">{agent.memory.persistence ? 'Enabled' : 'Disabled'}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Vector Storage</p>
                    <p className="mt-1">{agent.memory.vectorStorage ? 'Enabled' : 'Disabled'}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Message Count</p>
                    <p className="mt-1">{agent.memory.messageCount}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Last Updated</p>
                    <p className="mt-1">{new Date(agent.memory.lastUpdated).toLocaleString()}</p>
                  </div>
                </div>
              </div>
            </div>
            
            <div>
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Memory Explorer</h3>
              <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-md">
                <p className="text-gray-500 dark:text-gray-400 text-center my-8">
                  Memory explorer is under development.
                </p>
              </div>
            </div>
          </div>
        )}
        
        {activeTab === 'settings' && (
          <div className="p-6">
            <div className="mb-6">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Model Configuration</h3>
              <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-md">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Provider</p>
                    <p className="mt-1">{agent.model.provider}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Model</p>
                    <p className="mt-1">{agent.model.name}</p>
                  </div>
                  {agent.model.endpoint && (
                    <div className="col-span-2">
                      <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Endpoint</p>
                      <p className="mt-1">{agent.model.endpoint}</p>
                    </div>
                  )}
                  <div className="col-span-2">
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Options</p>
                    <div className="mt-1 bg-gray-100 dark:bg-gray-800 p-2 rounded">
                      <pre className="text-xs overflow-x-auto">
                        {JSON.stringify(agent.model.options, null, 2)}
                      </pre>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="mb-6">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Capabilities</h3>
              <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-md">
                <div className="flex flex-wrap gap-2">
                  {agent.capabilities.map((capability, index) => (
                    <span 
                      key={index}
                      className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200"
                    >
                      {capability}
                    </span>
                  ))}
                </div>
              </div>
            </div>
            
            <div>
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Tools</h3>
              <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-md">
                {agent.tools.length === 0 ? (
                  <p className="text-gray-500 dark:text-gray-400 text-center my-2">
                    No tools configured
                  </p>
                ) : (
                  <ul className="divide-y divide-gray-200 dark:divide-gray-700">
                    {agent.tools.map((tool) => (
                      <li key={tool.id} className="py-3">
                        <div className="flex items-center justify-between">
                          <div>
                            <h4 className="text-sm font-medium">{tool.name}</h4>
                            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                              {tool.description}
                            </p>
                          </div>
                          <div className="text-xs text-gray-500 dark:text-gray-400">
                            v{tool.version}
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
      
      {/* Delete confirmation modal */}
      {showDeleteConfirm && (
        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 transition-opacity" aria-hidden="true">
              <div className="absolute inset-0 bg-gray-500 opacity-75 dark:bg-gray-900 dark:opacity-90"></div>
            </div>
            
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
            
            <div className="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full p-6">
              <div>
                <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 dark:bg-red-900">
                  <ExclamationTriangleIcon className="h-6 w-6 text-red-600 dark:text-red-200" aria-hidden="true" />
                </div>
                <div className="mt-3 text-center sm:mt-5">
                  <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
                    Delete Agent
                  </h3>
                  <div className="mt-2">
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      Are you sure you want to delete this agent? This action cannot be undone.
                      All data associated with this agent will be permanently removed.
                    </p>
                  </div>
                </div>
              </div>
              <div className="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense">
                <button
                  type="button"
                  className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:col-start-2 sm:text-sm"
                  onClick={handleDeleteAgent}
                >
                  Delete
                </button>
                <button
                  type="button"
                  className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:col-start-1 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-600"
                  onClick={() => setShowDeleteConfirm(false)}
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AgentDetail;