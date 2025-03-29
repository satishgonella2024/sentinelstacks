import React, { useEffect, useState } from 'react';
import { invoke } from '@tauri-apps/api/tauri';
import { 
  CpuChipIcon, 
  ServerIcon, 
  ChatBubbleLeftRightIcon, 
  ClockIcon 
} from '@heroicons/react/24/outline';
import { Link } from 'react-router-dom';

type Agent = {
  name: string;
  description: string;
  model: string;
  memory_type: string;
};

const Dashboard: React.FC = () => {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const result = await invoke<Agent[]>('get_agents');
        setAgents(result);
        setIsLoading(false);
      } catch (err) {
        console.error('Error fetching agents:', err);
        setError('Failed to load agents. Please try again.');
        setIsLoading(false);
      }
    };

    fetchAgents();
  }, []);

  // Stats for dashboard
  const stats = [
    { name: 'Total Agents', value: agents.length, icon: CpuChipIcon, color: 'bg-blue-500' },
    { name: 'Active Agents', value: 0, icon: ServerIcon, color: 'bg-green-500' },
    { name: 'Recent Conversations', value: 0, icon: ChatBubbleLeftRightIcon, color: 'bg-purple-500' },
    { name: 'Uptime', value: '23h 12m', icon: ClockIcon, color: 'bg-yellow-500' },
  ];

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 bg-red-100 text-red-700 rounded-md">
        <p>{error}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Dashboard</h1>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <div
            key={stat.name}
            className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow"
          >
            <div className="flex items-center">
              <div className={`${stat.color} p-3 rounded-md`}>
                <stat.icon className="h-6 w-6 text-white" />
              </div>
              <div className="ml-4">
                <h2 className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {stat.name}
                </h2>
                <p className="text-lg font-semibold text-gray-800 dark:text-white">
                  {stat.value}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Recent Agents */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
        <h2 className="text-lg font-semibold text-gray-800 dark:text-white mb-4">
          Recent Agents
        </h2>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Name
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Description
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Model
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Memory
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {agents.length > 0 ? (
                agents.slice(0, 5).map((agent) => (
                  <tr key={agent.name} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <Link to={`/agents/${agent.name}`} className="text-blue-500 hover:text-blue-700">
                        {agent.name}
                      </Link>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-gray-700 dark:text-gray-300">
                      {agent.description.length > 50
                        ? `${agent.description.substring(0, 50)}...`
                        : agent.description}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-gray-700 dark:text-gray-300">
                      {agent.model}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-gray-700 dark:text-gray-300">
                      {agent.memory_type}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={4} className="px-6 py-4 text-center text-gray-500 dark:text-gray-400">
                    No agents found. <Link to="/agents/create" className="text-blue-500 hover:text-blue-700">Create one</Link>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
        
        {agents.length > 5 && (
          <div className="mt-4 text-right">
            <Link to="/agents" className="text-blue-500 hover:text-blue-700">
              View all agents
            </Link>
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
        <h2 className="text-lg font-semibold text-gray-800 dark:text-white mb-4">
          Quick Actions
        </h2>
        <div className="grid grid-cols-2 gap-4">
          <Link
            to="/agents/create"
            className="flex items-center justify-center p-4 bg-blue-100 dark:bg-blue-900 rounded-lg text-blue-700 dark:text-blue-100 hover:bg-blue-200 dark:hover:bg-blue-800 transition-colors"
          >
            <CpuChipIcon className="h-5 w-5 mr-2" />
            Create Agent
          </Link>
          <Link
            to="/registry"
            className="flex items-center justify-center p-4 bg-green-100 dark:bg-green-900 rounded-lg text-green-700 dark:text-green-100 hover:bg-green-200 dark:hover:bg-green-800 transition-colors"
          >
            <ServerIcon className="h-5 w-5 mr-2" />
            Browse Registry
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;