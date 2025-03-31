import React from 'react'

const Settings: React.FC = () => {
  return (
    <div className="p-4 max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-display text-white mb-2">Settings</h1>
        <p className="text-gray-400">Configure your Sentinel Stacks environment</p>
      </div>
      
      <div className="glass p-8 rounded-lg">
        <div className="mb-8">
          <h2 className="text-xl font-display text-white mb-4">System Settings</h2>
          <p className="text-gray-400 mb-6">
            Control panel for system-wide configurations is under development.
          </p>
        </div>
        
        <div className="mb-8">
          <h2 className="text-xl font-display text-white mb-4">API Keys</h2>
          <p className="text-gray-400 mb-6">
            Manage API keys and integrations for your agents.
          </p>
        </div>
        
        <div className="mb-8">
          <h2 className="text-xl font-display text-white mb-4">User Preferences</h2>
          <p className="text-gray-400 mb-6">
            Customize your experience with personal preferences.
          </p>
        </div>
      </div>
    </div>
  )
}

export default Settings 