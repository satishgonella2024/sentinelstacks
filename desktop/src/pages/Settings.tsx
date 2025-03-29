import React from 'react';

const Settings: React.FC = () => {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Settings</h1>
      <div className="max-w-3xl space-y-6">
        {/* General Settings */}
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm">
          <h2 className="text-lg font-semibold mb-4">General Settings</h2>
          <div className="space-y-4">
            <div>
              <label htmlFor="theme" className="block text-sm font-medium mb-1">
                Theme
              </label>
              <select
                id="theme"
                className="input w-full"
                defaultValue="system"
              >
                <option value="light">Light</option>
                <option value="dark">Dark</option>
                <option value="system">System</option>
              </select>
            </div>
            <div>
              <label htmlFor="language" className="block text-sm font-medium mb-1">
                Language
              </label>
              <select
                id="language"
                className="input w-full"
                defaultValue="en"
              >
                <option value="en">English</option>
                <option value="es">Spanish</option>
                <option value="fr">French</option>
              </select>
            </div>
          </div>
        </div>

        {/* Agent Defaults */}
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm">
          <h2 className="text-lg font-semibold mb-4">Agent Defaults</h2>
          <div className="space-y-4">
            <div>
              <label htmlFor="check-interval" className="block text-sm font-medium mb-1">
                Health Check Interval (seconds)
              </label>
              <input
                type="number"
                id="check-interval"
                className="input w-full"
                defaultValue={30}
                min={5}
              />
            </div>
            <div>
              <label htmlFor="timeout" className="block text-sm font-medium mb-1">
                Connection Timeout (seconds)
              </label>
              <input
                type="number"
                id="timeout"
                className="input w-full"
                defaultValue={10}
                min={1}
              />
            </div>
            <div>
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded border-gray-300 text-primary-600 shadow-sm focus:border-primary-300 focus:ring focus:ring-primary-200 focus:ring-opacity-50"
                  defaultChecked
                />
                <span className="text-sm font-medium">Enable Auto-Recovery</span>
              </label>
            </div>
          </div>
        </div>

        {/* Notifications */}
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm">
          <h2 className="text-lg font-semibold mb-4">Notifications</h2>
          <div className="space-y-4">
            <div>
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded border-gray-300 text-primary-600 shadow-sm focus:border-primary-300 focus:ring focus:ring-primary-200 focus:ring-opacity-50"
                  defaultChecked
                />
                <span className="text-sm font-medium">Email Notifications</span>
              </label>
            </div>
            <div>
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded border-gray-300 text-primary-600 shadow-sm focus:border-primary-300 focus:ring focus:ring-primary-200 focus:ring-opacity-50"
                  defaultChecked
                />
                <span className="text-sm font-medium">Desktop Notifications</span>
              </label>
            </div>
          </div>
        </div>

        {/* Save Button */}
        <div className="flex justify-end">
          <button className="btn btn-primary">Save Changes</button>
        </div>
      </div>
    </div>
  );
};

export default Settings; 