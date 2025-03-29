import React from 'react';
import { CheckIcon, PlusIcon } from '@heroicons/react/24/outline';

// Sample capabilities
const availableCapabilities = [
  { id: 'conversation', name: 'Conversation', description: 'Natural language conversation with memory' },
  { id: 'code', name: 'Code Generation', description: 'Generate and analyze code' },
  { id: 'infrastructure', name: 'Infrastructure', description: 'Manage cloud infrastructure' },
  { id: 'security', name: 'Security', description: 'Security analysis and recommendations' },
  { id: 'documentation', name: 'Documentation', description: 'Generate and maintain documentation' }
];

interface CapabilitySelectionProps {
  selectedCapabilities: string[];
  onToggleCapability: (capabilityId: string) => void;
}

const CapabilitySelection: React.FC<CapabilitySelectionProps> = ({ 
  selectedCapabilities, 
  onToggleCapability 
}) => {
  return (
    <div>
      <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">Capabilities</h2>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        {availableCapabilities.map(capability => (
          <div
            key={capability.id}
            className={`relative flex items-start p-4 border rounded-lg cursor-pointer transition-colors ${
              selectedCapabilities.includes(capability.id) 
                ? 'bg-primary-50 border-primary-200 dark:bg-primary-900 dark:border-primary-700' 
                : 'bg-white border-gray-200 hover:bg-gray-50 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700'
            }`}
            onClick={() => onToggleCapability(capability.id)}
          >
            <div className="min-w-0 flex-1 text-sm">
              <label className="font-medium text-gray-700 dark:text-gray-300 select-none">
                {capability.name}
              </label>
              <p className="text-gray-500 dark:text-gray-400">{capability.description}</p>
            </div>
            <div className="ml-3 flex items-center h-5">
              {selectedCapabilities.includes(capability.id) ? (
                <CheckIcon className="h-5 w-5 text-primary-600" />
              ) : (
                <PlusIcon className="h-5 w-5 text-gray-400" />
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default CapabilitySelection;