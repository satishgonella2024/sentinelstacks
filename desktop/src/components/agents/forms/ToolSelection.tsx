import React from 'react';
import { CheckIcon, PlusIcon } from '@heroicons/react/24/outline';
import { Tool } from '../../../types/Agent';

// Sample available tools
const availableTools: Tool[] = [
  { 
    id: 'calculator', 
    name: 'Calculator', 
    description: 'Perform mathematical calculations', 
    version: '1.0.0' 
  },
  { 
    id: 'url-fetcher', 
    name: 'URL Fetcher', 
    description: 'Fetch content from URLs', 
    version: '1.0.0' 
  },
  { 
    id: 'terraform', 
    name: 'Terraform', 
    description: 'Manage infrastructure as code', 
    version: '1.0.0' 
  },
  { 
    id: 'weather', 
    name: 'Weather', 
    description: 'Get weather information', 
    version: '1.0.0' 
  }
];

interface ToolSelectionProps {
  selectedTools: Tool[];
  onToggleTool: (tool: Tool) => void;
}

const ToolSelection: React.FC<ToolSelectionProps> = ({ 
  selectedTools, 
  onToggleTool 
}) => {
  const isToolSelected = (toolId: string) => {
    return selectedTools.some(t => t.id === toolId);
  };

  return (
    <div>
      <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">Tools</h2>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        {availableTools.map(tool => (
          <div
            key={tool.id}
            className={`relative flex items-start p-4 border rounded-lg cursor-pointer transition-colors ${
              isToolSelected(tool.id) 
                ? 'bg-primary-50 border-primary-200 dark:bg-primary-900 dark:border-primary-700' 
                : 'bg-white border-gray-200 hover:bg-gray-50 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700'
            }`}
            onClick={() => onToggleTool(tool)}
          >
            <div className="min-w-0 flex-1 text-sm">
              <label className="font-medium text-gray-700 dark:text-gray-300 select-none">
                {tool.name}
              </label>
              <p className="text-gray-500 dark:text-gray-400">{tool.description}</p>
              <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">v{tool.version}</p>
            </div>
            <div className="ml-3 flex items-center h-5">
              {isToolSelected(tool.id) ? (
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

export default ToolSelection;