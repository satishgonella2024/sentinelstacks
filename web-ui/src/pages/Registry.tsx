import React, { useState } from 'react';

interface RegistryItem {
  id: string;
  name: string;
  provider: string;
  url: string;
  status: 'active' | 'inactive';
  modelsCount: number;
  lastSync: string;
}

const Registry: React.FC = () => {
  // Mock registry data
  const [registries] = useState<RegistryItem[]>([
    {
      id: '1',
      name: 'OpenAI Models',
      provider: 'OpenAI',
      url: 'https://api.openai.com/v1/models',
      status: 'active',
      modelsCount: 12,
      lastSync: '2023-07-15T14:30:00Z'
    },
    {
      id: '2',
      name: 'Anthropic Models',
      provider: 'Anthropic',
      url: 'https://api.anthropic.com/v1/models',
      status: 'active',
      modelsCount: 5,
      lastSync: '2023-07-18T09:15:00Z'
    },
    {
      id: '3',
      name: 'Hugging Face Hub',
      provider: 'Hugging Face',
      url: 'https://huggingface.co/api/models',
      status: 'active',
      modelsCount: 38,
      lastSync: '2023-07-10T11:45:00Z'
    },
    {
      id: '4',
      name: 'Internal Model Registry',
      provider: 'Self-hosted',
      url: 'http://model-registry.internal:8080',
      status: 'inactive',
      modelsCount: 7,
      lastSync: '2023-06-30T16:20:00Z'
    }
  ]);

  // Format date
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  };

  return (
    <div className="container px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold">Model Registry</h1>
        <button className="px-4 py-2 bg-primary-600 text-white rounded">
          Connect New Registry
        </button>
      </div>

      <div className="bg-gray-800 rounded-lg overflow-hidden shadow-lg mb-8">
        <table className="w-full text-left">
          <thead>
            <tr className="border-b border-gray-700">
              <th className="px-6 py-4 font-medium">Name</th>
              <th className="px-6 py-4 font-medium">Provider</th>
              <th className="px-6 py-4 font-medium hidden md:table-cell">Models</th>
              <th className="px-6 py-4 font-medium hidden lg:table-cell">Last Sync</th>
              <th className="px-6 py-4 font-medium">Status</th>
              <th className="px-6 py-4 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody>
            {registries.map(registry => (
              <tr key={registry.id} className="border-b border-gray-700 hover:bg-gray-700">
                <td className="px-6 py-4">{registry.name}</td>
                <td className="px-6 py-4">{registry.provider}</td>
                <td className="px-6 py-4 hidden md:table-cell">{registry.modelsCount}</td>
                <td className="px-6 py-4 hidden lg:table-cell">
                  {formatDate(registry.lastSync)}
                </td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded text-xs ${
                    registry.status === 'active' 
                      ? 'bg-green-900 text-green-300' 
                      : 'bg-red-900 text-red-300'
                  }`}>
                    {registry.status}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <div className="flex space-x-2">
                    <button className="text-primary-400 hover:text-primary-300">
                      Sync
                    </button>
                    <button className="text-primary-400 hover:text-primary-300">
                      Edit
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="bg-gray-800 rounded-lg p-6 shadow-lg">
        <h2 className="text-xl font-bold mb-4">Registry Status</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-gray-700 rounded-lg p-4">
            <h3 className="text-sm font-medium text-gray-400 mb-1">Total Models</h3>
            <p className="text-2xl font-bold">{registries.reduce((sum, reg) => sum + reg.modelsCount, 0)}</p>
          </div>
          <div className="bg-gray-700 rounded-lg p-4">
            <h3 className="text-sm font-medium text-gray-400 mb-1">Active Registries</h3>
            <p className="text-2xl font-bold">{registries.filter(r => r.status === 'active').length}</p>
          </div>
          <div className="bg-gray-700 rounded-lg p-4">
            <h3 className="text-sm font-medium text-gray-400 mb-1">Last Sync</h3>
            <p className="text-2xl font-bold">
              {formatDate(
                registries
                  .map(r => new Date(r.lastSync).getTime())
                  .reduce((max, date) => Math.max(max, date), 0)
                  .toString()
              )}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Registry; 