import React, { useState } from 'react';
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline';
import { Link } from 'react-router-dom';
import { Agent, ModelProvider, Tool } from '../../../types/Agent';
import ModelConfig from './ModelConfig';
import ToolSelection from './ToolSelection';
import CapabilitySelection from './CapabilitySelection';

// Initial form state for new agents
const getInitialFormState = () => ({
  name: '',
  description: '',
  model: {
    provider: 'openai' as ModelProvider,
    name: 'gpt-4',
    options: {
      temperature: 0.7,
      max_tokens: 2000
    }
  },
  tools: [] as Tool[],
  capabilities: ['conversation'] as string[],
  memory: {
    persistence: true,
    vectorStorage: true,
    messageCount: 0,
    lastUpdated: new Date().toISOString()
  }
});

interface AgentFormProps {
  agent?: Agent; // If provided, we're editing an existing agent
  isSubmitting: boolean;
  error: string | null;
  onSubmit: (formData: any) => void;
  cancelPath: string;
}

const AgentForm: React.FC<AgentFormProps> = ({
  agent,
  isSubmitting,
  error,
  onSubmit,
  cancelPath
}) => {
  // Initialize with existing agent data or defaults
  const [formData, setFormData] = useState(agent || getInitialFormState());

  // Handlers for basic field changes
  const handleTextChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  // Handlers for model configuration
  const handleProviderChange = (provider: ModelProvider) => {
    // Set default model for the selected provider
    let defaultModel: string;
    switch (provider) {
      case 'openai':
        defaultModel = 'gpt-4';
        break;
      case 'claude':
        defaultModel = 'claude-3-opus';
        break;
      case 'ollama':
        defaultModel = 'llama2';
        break;
      default:
        defaultModel = 'gpt-4';
    }
    
    setFormData(prev => ({
      ...prev,
      model: {
        ...prev.model,
        provider,
        name: defaultModel
      }
    }));
  };

  const handleModelChange = (modelName: string) => {
    setFormData(prev => ({
      ...prev,
      model: {
        ...prev.model,
        name: modelName
      }
    }));
  };

  const handleTemperatureChange = (temperature: number) => {
    setFormData(prev => ({
      ...prev,
      model: {
        ...prev.model,
        options: {
          ...prev.model.options,
          temperature
        }
      }
    }));
  };

  const handleMaxTokensChange = (maxTokens: number) => {
    setFormData(prev => ({
      ...prev,
      model: {
        ...prev.model,
        options: {
          ...prev.model.options,
          max_tokens: maxTokens
        }
      }
    }));
  };

  // Handler for tool selection
  const handleToggleTool = (tool: Tool) => {
    setFormData(prev => {
      const toolExists = prev.tools.some(t => t.id === tool.id);
      const tools = toolExists
        ? prev.tools.filter(t => t.id !== tool.id)
        : [...prev.tools, tool];

      return { ...prev, tools };
    });
  };

  // Handler for capability selection
  const handleToggleCapability = (capabilityId: string) => {
    setFormData(prev => {
      const capabilities = prev.capabilities.includes(capabilityId)
        ? prev.capabilities.filter(c => c !== capabilityId)
        : [...prev.capabilities, capabilityId];

      return { ...prev, capabilities };
    });
  };

  // Handler for memory options
  const handleToggleMemoryOption = (option: 'persistence' | 'vectorStorage') => {
    setFormData(prev => ({
      ...prev,
      memory: {
        ...prev.memory,
        [option]: !prev.memory[option]
      }
    }));
  };

  // Form submission handler
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      {/* Error display */}
      {error && (
        <div className="bg-red-50 dark:bg-red-900 p-4 rounded-md">
          <div className="flex">
            <div className="flex-shrink-0">
              <ExclamationTriangleIcon className="h-5 w-5 text-red-400" aria-hidden="true" />
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800 dark:text-red-200">Error</h3>
              <div className="mt-2 text-sm text-red-700 dark:text-red-300">
                <p>{error}</p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Basic Information */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">Basic Information</h2>
        <div className="grid grid-cols-1 gap-6">
          <div>
            <label htmlFor="name" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Name <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              name="name"
              id="name"
              value={formData.name}
              onChange={handleTextChange}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600"
              placeholder="My Agent"
              required
            />
          </div>

          <div>
            <label htmlFor="description" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Description
            </label>
            <textarea
              name="description"
              id="description"
              rows={3}
              value={formData.description}
              onChange={handleTextChange}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600"
              placeholder="What does this agent do?"
            />
          </div>
        </div>
      </div>

      {/* Model Configuration */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <ModelConfig
          provider={formData.model.provider}
          modelName={formData.model.name}
          temperature={formData.model.options.temperature}
          maxTokens={formData.model.options.max_tokens}
          onProviderChange={handleProviderChange}
          onModelChange={handleModelChange}
          onTemperatureChange={handleTemperatureChange}
          onMaxTokensChange={handleMaxTokensChange}
        />
      </div>

      {/* Capabilities */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <CapabilitySelection
          selectedCapabilities={formData.capabilities}
          onToggleCapability={handleToggleCapability}
        />
      </div>

      {/* Tools */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <ToolSelection
          selectedTools={formData.tools}
          onToggleTool={handleToggleTool}
        />
      </div>

      {/* Memory Options */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">Memory Options</h2>
        <div className="space-y-4">
          <div className="relative flex items-start">
            <div className="flex items-center h-5">
              <input
                id="persistence"
                name="persistence"
                type="checkbox"
                checked={formData.memory.persistence}
                onChange={() => handleToggleMemoryOption('persistence')}
                className="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500 dark:border-gray-600"
              />
            </div>
            <div className="ml-3 text-sm">
              <label htmlFor="persistence" className="font-medium text-gray-700 dark:text-gray-300">
                Persistence
              </label>
              <p className="text-gray-500 dark:text-gray-400">Save agent state between sessions</p>
            </div>
          </div>

          <div className="relative flex items-start">
            <div className="flex items-center h-5">
              <input
                id="vectorStorage"
                name="vectorStorage"
                type="checkbox"
                checked={formData.memory.vectorStorage}
                onChange={() => handleToggleMemoryOption('vectorStorage')}
                className="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500 dark:border-gray-600"
              />
            </div>
            <div className="ml-3 text-sm">
              <label htmlFor="vectorStorage" className="font-medium text-gray-700 dark:text-gray-300">
                Vector Storage
              </label>
              <p className="text-gray-500 dark:text-gray-400">Enable semantic search in the agent's memory</p>
            </div>
          </div>
        </div>
      </div>

      {/* Form Actions */}
      <div className="flex justify-end">
        <Link
          to={cancelPath}
          className="mr-3 inline-flex justify-center py-2 px-4 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600 dark:hover:bg-gray-600"
        >
          Cancel
        </Link>
        <button
          type="submit"
          className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Saving...' : agent ? 'Update Agent' : 'Create Agent'}
        </button>
      </div>
    </form>
  );
};

export default AgentForm;