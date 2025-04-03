import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { FiSearch, FiSave, FiTrash2, FiCpu } from 'react-icons/fi';

interface MemoryItem {
  key: string;
  value: string;
  timestamp: string;
  type: 'key-value' | 'vector';
}

const MemoryManager: React.FC = () => {
  const { agentId } = useParams<{ agentId: string }>();
  const [searchQuery, setSearchQuery] = useState('');
  const [memoryItems, setMemoryItems] = useState<MemoryItem[]>([]);
  const [selectedMemoryType, setSelectedMemoryType] = useState<'key-value' | 'vector' | 'all'>('all');
  const [newKey, setNewKey] = useState('');
  const [newValue, setNewValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Mock data for development
  const mockMemoryItems: MemoryItem[] = [
    { key: 'user_preferences', value: '{"theme":"dark","notifications":true}', timestamp: new Date().toISOString(), type: 'key-value' },
    { key: 'conversation_history', value: 'Last conversation was about AI ethics and safety measures.', timestamp: new Date(Date.now() - 86400000).toISOString(), type: 'key-value' },
    { key: 'research_topics', value: 'Machine learning, neural networks, reinforcement learning', timestamp: new Date(Date.now() - 172800000).toISOString(), type: 'vector' },
    { key: 'project_deadlines', value: 'AI safety report due on 2023-12-15', timestamp: new Date(Date.now() - 259200000).toISOString(), type: 'key-value' },
    { key: 'technical_concepts', value: 'Transformers architecture, attention mechanisms, and embeddings.', timestamp: new Date(Date.now() - 345600000).toISOString(), type: 'vector' },
  ];

  useEffect(() => {
    // In a real implementation, this would fetch from the API
    setIsLoading(true);
    // Simulate API call
    setTimeout(() => {
      setMemoryItems(mockMemoryItems);
      setIsLoading(false);
    }, 500);
  }, [agentId]);

  const handleSearch = () => {
    setIsLoading(true);
    // Simulate search API call
    setTimeout(() => {
      if (searchQuery.trim() === '') {
        setMemoryItems(mockMemoryItems.filter(item => 
          selectedMemoryType === 'all' ? true : item.type === selectedMemoryType
        ));
      } else {
        const filtered = mockMemoryItems.filter(item => 
          (item.key.toLowerCase().includes(searchQuery.toLowerCase()) || 
           item.value.toLowerCase().includes(searchQuery.toLowerCase())) &&
          (selectedMemoryType === 'all' ? true : item.type === selectedMemoryType)
        );
        setMemoryItems(filtered);
      }
      setIsLoading(false);
    }, 300);
  };

  const handleStoreMemory = () => {
    if (!newKey.trim() || !newValue.trim()) {
      setError('Both key and value are required');
      return;
    }

    setIsLoading(true);
    setError(null);

    // Simulate API call to store memory
    setTimeout(() => {
      const newItem: MemoryItem = {
        key: newKey,
        value: newValue,
        timestamp: new Date().toISOString(),
        type: 'key-value' // Default to key-value for new entries
      };
      
      setMemoryItems([newItem, ...memoryItems]);
      setNewKey('');
      setNewValue('');
      setIsLoading(false);
    }, 500);
  };

  const handleDeleteMemory = (key: string) => {
    setIsLoading(true);
    
    // Simulate API call to delete memory
    setTimeout(() => {
      setMemoryItems(memoryItems.filter(item => item.key !== key));
      setIsLoading(false);
    }, 300);
  };

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleString();
  };

  return (
    <div className="bg-gray-900 rounded-lg p-6 my-4">
      <h2 className="text-2xl font-semibold text-white mb-6 flex items-center">
        <FiCpu className="mr-2" /> Memory Management
      </h2>
      
      {/* Memory Type Filter */}
      <div className="mb-6">
        <div className="text-sm text-gray-400 mb-2">Filter by Type</div>
        <div className="flex space-x-3">
          <button
            className={`px-4 py-2 rounded-md ${selectedMemoryType === 'all' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-300'}`}
            onClick={() => setSelectedMemoryType('all')}
          >
            All
          </button>
          <button
            className={`px-4 py-2 rounded-md ${selectedMemoryType === 'key-value' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-300'}`}
            onClick={() => setSelectedMemoryType('key-value')}
          >
            Key-Value
          </button>
          <button
            className={`px-4 py-2 rounded-md ${selectedMemoryType === 'vector' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-300'}`}
            onClick={() => setSelectedMemoryType('vector')}
          >
            Vector
          </button>
        </div>
      </div>

      {/* Search Memory */}
      <div className="mb-6">
        <div className="text-sm text-gray-400 mb-2">Search Memory</div>
        <div className="flex">
          <input
            type="text"
            placeholder="Search by key or value..."
            className="bg-gray-800 text-white px-4 py-2 rounded-l-md w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
          />
          <button 
            className="bg-blue-600 text-white px-4 py-2 rounded-r-md hover:bg-blue-700 focus:outline-none"
            onClick={handleSearch}
          >
            <FiSearch />
          </button>
        </div>
      </div>

      {/* Store New Memory */}
      <div className="bg-gray-800 p-4 rounded-md mb-6">
        <div className="text-sm text-gray-400 mb-2">Store New Memory</div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-3">
          <input
            type="text"
            placeholder="Key"
            className="bg-gray-700 text-white px-4 py-2 rounded-md w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={newKey}
            onChange={(e) => setNewKey(e.target.value)}
          />
          <input
            type="text"
            placeholder="Value"
            className="bg-gray-700 text-white px-4 py-2 rounded-md w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={newValue}
            onChange={(e) => setNewValue(e.target.value)}
          />
        </div>
        {error && <div className="text-red-500 text-sm mb-3">{error}</div>}
        <button 
          className="bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700 focus:outline-none flex items-center justify-center"
          onClick={handleStoreMemory}
          disabled={isLoading}
        >
          <FiSave className="mr-2" /> Store Memory
        </button>
      </div>

      {/* Memory Items */}
      <div className="bg-gray-800 rounded-md">
        <div className="text-sm text-gray-400 p-4 border-b border-gray-700">Memory Items</div>
        {isLoading ? (
          <div className="p-8 flex justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-blue-500"></div>
          </div>
        ) : memoryItems.length === 0 ? (
          <div className="p-8 text-center text-gray-500">No memory items found</div>
        ) : (
          <div className="divide-y divide-gray-700">
            {memoryItems.map((item) => (
              <div key={item.key} className="p-4 hover:bg-gray-750">
                <div className="flex justify-between items-start mb-2">
                  <div className="flex-1">
                    <div className="flex items-center">
                      <span className="text-white font-medium">{item.key}</span>
                      <span className={`ml-2 px-2 py-1 text-xs rounded-full ${
                        item.type === 'key-value' ? 'bg-purple-900 text-purple-300' : 'bg-teal-900 text-teal-300'
                      }`}>
                        {item.type}
                      </span>
                    </div>
                    <div className="text-sm text-gray-400 mt-1">
                      {formatTimestamp(item.timestamp)}
                    </div>
                  </div>
                  <button 
                    className="text-gray-400 hover:text-red-500 transition-colors"
                    onClick={() => handleDeleteMemory(item.key)}
                  >
                    <FiTrash2 />
                  </button>
                </div>
                <div className="bg-gray-700 p-3 rounded-md">
                  <pre className="text-sm text-gray-300 whitespace-pre-wrap">
                    {item.value}
                  </pre>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default MemoryManager; 