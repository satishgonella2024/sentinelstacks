import React, { useState, useEffect } from 'react';
import { ArchiveBoxIcon, ClockIcon, ServerIcon, DocumentTextIcon } from '@heroicons/react/24/outline';
import { getAgentMemory } from '../../services/agentService';

interface MemoryItem {
  role: string;
  content: string;
  timestamp: string;
}

interface VectorStoreInfo {
  size: number;
  lastUpdated: string;
}

interface MemoryData {
  messages: MemoryItem[];
  vectorStore: VectorStoreInfo;
}

interface MemoryVisualizerProps {
  agentId: string;
}

const MemoryVisualizer: React.FC<MemoryVisualizerProps> = ({ agentId }) => {
  const [memoryData, setMemoryData] = useState<MemoryData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeView, setActiveView] = useState<'messages' | 'vectors'>('messages');

  useEffect(() => {
    const fetchMemoryData = async () => {
      setIsLoading(true);
      try {
        const data = await getAgentMemory(agentId);
        setMemoryData(data);
      } catch (err) {
        setError(`Failed to load memory data: ${err instanceof Error ? err.message : String(err)}`);
      } finally {
        setIsLoading(false);
      }
    };

    fetchMemoryData();
  }, [agentId]);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 dark:bg-red-900 p-4 rounded-md">
        <p className="text-red-600 dark:text-red-200">{error}</p>
      </div>
    );
  }

  if (!memoryData) {
    return (
      <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-md">
        <p className="text-gray-500 dark:text-gray-400 text-center">No memory data available.</p>
      </div>
    );
  }

  return (
    <div>
      {/* Memory Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
        <div className="bg-white dark:bg-gray-800 shadow-sm rounded-lg p-4">
          <div className="flex items-center">
            <DocumentTextIcon className="h-8 w-8 text-primary-500" />
            <div className="ml-3">
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Total Messages</p>
              <p className="text-xl font-semibold">{memoryData.messages.length}</p>
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 shadow-sm rounded-lg p-4">
          <div className="flex items-center">
            <ServerIcon className="h-8 w-8 text-primary-500" />
            <div className="ml-3">
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Vector Store Size</p>
              <p className="text-xl font-semibold">{memoryData.vectorStore.size} items</p>
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 shadow-sm rounded-lg p-4">
          <div className="flex items-center">
            <ClockIcon className="h-8 w-8 text-primary-500" />
            <div className="ml-3">
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Last Updated</p>
              <p className="text-xl font-semibold">
                {new Date(memoryData.vectorStore.lastUpdated).toLocaleString()}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Memory Content */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
        {/* Tab Navigation */}
        <div className="border-b border-gray-200 dark:border-gray-700">
          <nav className="-mb-px flex">
            <button
              className={`py-4 px-6 border-b-2 font-medium text-sm ${
                activeView === 'messages'
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
              }`}
              onClick={() => setActiveView('messages')}
            >
              <DocumentTextIcon className="w-5 h-5 inline mr-2" />
              Messages
            </button>
            <button
              className={`py-4 px-6 border-b-2 font-medium text-sm ${
                activeView === 'vectors'
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
              }`}
              onClick={() => setActiveView('vectors')}
            >
              <ServerIcon className="w-5 h-5 inline mr-2" />
              Vector Store
            </button>
          </nav>
        </div>

        {/* Content Area */}
        <div className="p-6">
          {activeView === 'messages' && (
            <div className="space-y-4 max-h-96 overflow-y-auto">
              {memoryData.messages.length === 0 ? (
                <p className="text-gray-500 dark:text-gray-400 text-center py-4">
                  No message history available.
                </p>
              ) : (
                memoryData.messages.map((message, index) => (
                  <div
                    key={index}
                    className={`p-4 rounded-lg ${
                      message.role === 'user'
                        ? 'bg-gray-100 dark:bg-gray-700 ml-12'
                        : 'bg-primary-50 dark:bg-primary-900 mr-12'
                    }`}
                  >
                    <div className="flex justify-between items-start">
                      <span className="font-semibold capitalize">{message.role}</span>
                      <span className="text-xs text-gray-500 dark:text-gray-400">
                        {new Date(message.timestamp).toLocaleString()}
                      </span>
                    </div>
                    <p className="mt-2 whitespace-pre-wrap">{message.content}</p>
                  </div>
                ))
              )}
            </div>
          )}

          {activeView === 'vectors' && (
            <div>
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-medium">Vector Store</h3>
                <span className="text-sm text-gray-500 dark:text-gray-400">
                  {memoryData.vectorStore.size} items
                </span>
              </div>
              
              {memoryData.vectorStore.size === 0 ? (
                <div className="bg-gray-50 dark:bg-gray-900 p-8 rounded-lg flex flex-col items-center justify-center">
                  <ArchiveBoxIcon className="h-12 w-12 text-gray-400 dark:text-gray-500 mb-2" />
                  <p className="text-gray-500 dark:text-gray-400 text-center">
                    No vector data available.
                  </p>
                </div>
              ) : (
                <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-lg">
                  <p className="text-gray-600 dark:text-gray-300">
                    Vector embeddings visualization is under development. Current storage includes {memoryData.vectorStore.size} indexed items.
                  </p>
                  
                  <div className="mt-4 h-64 border border-gray-200 dark:border-gray-700 rounded-lg flex items-center justify-center">
                    <div className="text-center p-4">
                      <ServerIcon className="h-10 w-10 mx-auto text-gray-400 dark:text-gray-500 mb-2" />
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        Vector visualization will be available in a future update.
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MemoryVisualizer;
