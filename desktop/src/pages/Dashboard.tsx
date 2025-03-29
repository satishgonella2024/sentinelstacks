import React from 'react';
import { ChartBarIcon, UsersIcon, ServerIcon } from '@heroicons/react/24/outline';

const Dashboard: React.FC = () => {
  return (
    <div className="max-w-7xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Dashboard</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
          <div className="flex items-center">
            <ServerIcon className="w-8 h-8 text-primary-600 dark:text-primary-500" />
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Active Agents</p>
              <p className="text-2xl font-semibold text-gray-900 dark:text-white">12</p>
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
          <div className="flex items-center">
            <ChartBarIcon className="w-8 h-8 text-primary-600 dark:text-primary-500" />
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">System Load</p>
              <p className="text-2xl font-semibold text-gray-900 dark:text-white">45%</p>
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
          <div className="flex items-center">
            <UsersIcon className="w-8 h-8 text-primary-600 dark:text-primary-500" />
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Connected Users</p>
              <p className="text-2xl font-semibold text-gray-900 dark:text-white">24</p>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
        <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">System Overview</h2>
        <p className="text-gray-600 dark:text-gray-300">
          Welcome to SentinelStacks Dashboard. Monitor your system's performance and manage your agents from this central location.
        </p>
      </div>
    </div>
  );
};

export default Dashboard; 