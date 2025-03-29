import React from 'react';
import { ModelProvider } from '../../../types/Agent';

// Available model options
const modelOptions = {
  openai: [
    { name: 'gpt-3.5-turbo', label: 'GPT-3.5 Turbo' },
    { name: 'gpt-4', label: 'GPT-4' },
    { name: 'gpt-4-turbo', label: 'GPT-4 Turbo' }
  ],
  claude: [
    { name: 'claude-3-opus', label: 'Claude 3 Opus' },
    { name: 'claude-3-sonnet', label: 'Claude 3 Sonnet' },
    { name: 'claude-3-haiku', label: 'Claude 3 Haiku' }
  ],
  ollama: [
    { name: 'llama2', label: 'Llama 2' },
    { name: 'mistral', label: 'Mistral' },
    { name: 'codellama', label: 'CodeLlama' }
  ]
};

interface ModelConfigProps {
  provider: ModelProvider;
  modelName: string;
  temperature: number;
  maxTokens: number;
  onProviderChange: (provider: ModelProvider) => void;
  onModelChange: (modelName: string) => void;
  onTemperatureChange: (temperature: number) => void;
  onMaxTokensChange: (maxTokens: number) => void;
}

const ModelConfig: React.FC<ModelConfigProps> = ({
  provider,
  modelName,
  temperature,
  maxTokens,
  onProviderChange,
  onModelChange,
  onTemperatureChange,
  onMaxTokensChange
}) => {
  return (
    <div>
      <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">Model Configuration</h2>
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <div>
          <label htmlFor="model-provider" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Provider
          </label>
          <select
            id="model-provider"
            value={provider}
            onChange={(e) => onProviderChange(e.target.value as ModelProvider)}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600"
          >
            <option value="openai">OpenAI</option>
            <option value="claude">Claude</option>
            <option value="ollama">Ollama (Local)</option>
          </select>
        </div>

        <div>
          <label htmlFor="model-name" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Model
          </label>
          <select
            id="model-name"
            value={modelName}
            onChange={(e) => onModelChange(e.target.value)}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600"
          >
            {modelOptions[provider].map(model => (
              <option key={model.name} value={model.name}>
                {model.label}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label htmlFor="temperature" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Temperature
          </label>
          <div className="mt-1 flex items-center">
            <input
              type="range"
              id="temperature"
              min="0"
              max="1"
              step="0.1"
              value={temperature}
              onChange={(e) => onTemperatureChange(parseFloat(e.target.value))}
              className="flex-grow mr-2 h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
            />
            <span className="text-sm text-gray-900 dark:text-white min-w-[40px] text-center">{temperature.toFixed(1)}</span>
          </div>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Lower values (0.0) are more deterministic, higher values (1.0) are more creative.
          </p>
        </div>

        <div>
          <label htmlFor="max-tokens" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Max Tokens
          </label>
          <input
            type="number"
            id="max-tokens"
            value={maxTokens}
            onChange={(e) => onMaxTokensChange(parseInt(e.target.value, 10) || 0)}
            step="100"
            min="100"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600"
          />
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Maximum length of the model's output in tokens.
          </p>
        </div>
      </div>
    </div>
  );
};

export default ModelConfig;