import React, { useState } from 'react';
import { FiDatabase, FiCpu, FiBriefcase } from 'react-icons/fi';

// Import the MemoryManager component
const MemoryManager = React.lazy(() => import('../components/agents/MemoryManager'));

interface Agent {
  id: string;
  name: string;
  description: string;
}

const Memory: React.FC = () => {
  const [selectedAgent, setSelectedAgent] = useState<string | null>(null);
  
  // Mock data for development
  const mockAgents: Agent[] = [
    { id: '1', name: 'Research Assistant', description: 'Helps with academic research' },
    { id: '2', name: 'Image Analyzer', description: 'Analyzes images and provides descriptions' },
    { id: '3', name: 'Code Assistant', description: 'Helps with programming tasks' },
    { id: '4', name: 'Data Analyst', description: 'Assists with data analysis' },
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center mb-8">
        <FiDatabase className="text-blue-500 mr-3 text-3xl" />
        <h1 className="text-3xl font-bold text-white">Memory Management</h1>
      </div>

      <div className="bg-gray-800 rounded-lg p-6 mb-8">
        <h2 className="text-xl font-semibold text-white mb-4 flex items-center">
          <FiBriefcase className="mr-2" /> Select Agent
        </h2>
        <p className="text-gray-400 mb-6">
          Choose an agent to view and manage its memory store. Each agent has its own isolated memory stores.
        </p>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {mockAgents.map((agent) => (
            <div 
              key={agent.id}
              className={`bg-gray-900 rounded-lg p-4 cursor-pointer transition-all hover:shadow-lg hover:bg-gray-850 border-2 ${
                selectedAgent === agent.id ? 'border-blue-500' : 'border-transparent'
              }`}
              onClick={() => setSelectedAgent(agent.id)}
            >
              <div className="flex items-center mb-2">
                <FiCpu className="text-blue-500 mr-2" />
                <h3 className="text-lg font-medium text-white">{agent.name}</h3>
              </div>
              <p className="text-gray-400 text-sm">{agent.description}</p>
            </div>
          ))}
        </div>
      </div>

      {selectedAgent ? (
        <React.Suspense fallback={
          <div className="bg-gray-900 rounded-lg p-8 flex justify-center items-center">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
          </div>
        }>
          <MemoryManager />
        </React.Suspense>
      ) : (
        <div className="bg-gray-900 rounded-lg p-8 text-center">
          <FiDatabase className="text-gray-600 text-5xl mx-auto mb-4" />
          <h3 className="text-xl font-medium text-white mb-2">No Agent Selected</h3>
          <p className="text-gray-400">
            Select an agent from the list above to view and manage its memory.
          </p>
        </div>
      )}

      <div className="bg-gray-800 rounded-lg p-6 mt-8">
        <h2 className="text-xl font-semibold text-white mb-4">About Agent Memory</h2>
        <div className="text-gray-400 space-y-4">
          <p>
            SentinelStacks provides sophisticated memory management capabilities for AI agents, allowing them to
            store and retrieve information across conversations and sessions.
          </p>
          <p>
            <strong className="text-white">Key-Value Store:</strong> Simple storage for discrete pieces of information
            that can be retrieved by exact key matching.
          </p>
          <p>
            <strong className="text-white">Vector Store:</strong> Semantic storage that allows for similarity search
            and retrieval of information based on meaning rather than exact matching.
          </p>
          <p>
            Memory can be configured to persist between sessions or be ephemeral, depending on your agent's
            needs and privacy requirements.
          </p>
        </div>
      </div>
    </div>
  );
};

export default Memory; 