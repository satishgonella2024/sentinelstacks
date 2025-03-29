import React from 'react';

const Monitoring: React.FC = () => {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Monitoring</h1>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* System Metrics */}
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm">
          <h2 className="text-lg font-semibold mb-4">System Metrics</h2>
          <div className="space-y-6">
            {/* CPU Usage */}
            <div>
              <div className="flex justify-between mb-2">
                <span>CPU Usage</span>
                <span className="text-primary-500">45%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div className="bg-primary-500 h-2 rounded-full" style={{ width: '45%' }}></div>
              </div>
            </div>

            {/* Memory Usage */}
            <div>
              <div className="flex justify-between mb-2">
                <span>Memory Usage</span>
                <span className="text-primary-500">2.4GB / 8GB</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div className="bg-primary-500 h-2 rounded-full" style={{ width: '30%' }}></div>
              </div>
            </div>

            {/* Disk Usage */}
            <div>
              <div className="flex justify-between mb-2">
                <span>Disk Usage</span>
                <span className="text-primary-500">120GB / 500GB</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div className="bg-primary-500 h-2 rounded-full" style={{ width: '24%' }}></div>
              </div>
            </div>
          </div>
        </div>

        {/* Network Traffic */}
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm">
          <h2 className="text-lg font-semibold mb-4">Network Traffic</h2>
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <span>Inbound</span>
              <span className="text-green-500">2.5 MB/s</span>
            </div>
            <div className="flex justify-between items-center">
              <span>Outbound</span>
              <span className="text-blue-500">1.8 MB/s</span>
            </div>
            <div className="flex justify-between items-center">
              <span>Active Connections</span>
              <span className="text-primary-500">245</span>
            </div>
          </div>
        </div>

        {/* Logs */}
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm lg:col-span-2">
          <h2 className="text-lg font-semibold mb-4">Recent Logs</h2>
          <div className="space-y-3">
            <div className="text-sm">
              <span className="text-gray-500 dark:text-gray-400">[2024-03-29 13:45:22]</span>
              <span className="ml-2 text-yellow-500">[WARNING]</span>
              <span className="ml-2">High memory usage detected on agent web-01</span>
            </div>
            <div className="text-sm">
              <span className="text-gray-500 dark:text-gray-400">[2024-03-29 13:44:15]</span>
              <span className="ml-2 text-green-500">[INFO]</span>
              <span className="ml-2">Agent db-01 backup completed successfully</span>
            </div>
            <div className="text-sm">
              <span className="text-gray-500 dark:text-gray-400">[2024-03-29 13:43:01]</span>
              <span className="ml-2 text-red-500">[ERROR]</span>
              <span className="ml-2">Connection timeout on agent cache-01</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Monitoring; 