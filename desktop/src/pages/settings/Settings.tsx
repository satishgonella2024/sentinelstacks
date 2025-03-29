import React from 'react';

const Settings: React.FC = () => {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Settings</h1>
      <p className="text-gray-600 dark:text-gray-400">
        Configure your SentinelStacks environment settings here.
      </p>
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6 flex items-center justify-center h-64">
        <p className="text-gray-500 dark:text-gray-400 text-center">
          Settings panel coming soon.
        </p>
      </div>
    </div>
  );
};

export default Settings;