import React from 'react';
import { 
  PlayIcon, 
  StopIcon,
  PencilIcon,
  ClockIcon
} from '@heroicons/react/24/outline';
import { Agent } from '../../types/Agent';
import { Link } from 'react-router-dom';
import { startAgent, stopAgent } from '../../services/agentService';
import StatusBadge from '../common/StatusBadge';
import AgentToolBadge from './AgentToolBadge';

interface AgentCardProps {
  agent: Agent;
  onStatusChange?: () => void;
}

const AgentCard: React.FC<AgentCardProps> = ({ agent, onStatusChange }) => {
  const handleStartAgent = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    try {
      await startAgent(agent.id);
      if (onStatusChange) onStatusChange();
    } catch (error) {
      console.error(`Error starting agent ${agent.id}:`, error);
      // Show error notification
    }
  };

  const handleStopAgent = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    try {
      await stopAgent(agent.id);
      if (onStatusChange) onStatusChange();
    } catch (error) {
      console.error(`Error stopping agent ${agent.id}:`, error);
      // Show error notification
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Never';
    
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffSecs = Math.floor(diffMs / 1000);
    const diffMins = Math.floor(diffSecs / 60);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);
    
    if (diffSecs < 60) return `${diffSecs} seconds ago`;
    if (diffMins < 60) return `${diffMins} minutes ago`;
    if (diffHours < 24) return `${diffHours} hours ago`;
    if (diffDays < 30) return `${diffDays} days ago`;
    
    return date.toLocaleDateString();
  };

  return (
    <Link to={`/agents/${agent.id}`} className="block">
      <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg hover:shadow-md transition-shadow duration-300">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-start justify-between">
            <div className="flex-1 min-w-0">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white truncate">
                {agent.name}
              </h3>
              <div className="mt-1 flex items-center">
                <StatusBadge status={agent.status} />
                <span className="ml-2 text-sm text-gray-500 dark:text-gray-400">
                  {agent.model.provider} / {agent.model.name}
                </span>
              </div>
            </div>
            <div className="flex shrink-0 space-x-1">
              {agent.status === 'active' ? (
                <button
                  onClick={handleStopAgent}
                  className="inline-flex items-center p-1 border border-transparent rounded-full shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                >
                  <StopIcon className="h-4 w-4" aria-hidden="true" />
                </button>
              ) : (
                <button
                  onClick={handleStartAgent}
                  className="inline-flex items-center p-1 border border-transparent rounded-full shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                >
                  <PlayIcon className="h-4 w-4" aria-hidden="true" />
                </button>
              )}
              <Link
                to={`/agents/${agent.id}/edit`}
                className="inline-flex items-center p-1 border border-transparent rounded-full shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                onClick={(e) => e.stopPropagation()}
              >
                <PencilIcon className="h-4 w-4" aria-hidden="true" />
              </Link>
            </div>
          </div>
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400 line-clamp-2">
            {agent.description}
          </p>
          
          <div className="mt-4 flex flex-wrap gap-2">
            {agent.tools.slice(0, 2).map(tool => (
              <AgentToolBadge key={tool.id} tool={tool} />
            ))}
            {agent.tools.length > 2 && (
              <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200">
                +{agent.tools.length - 2} more
              </span>
            )}
          </div>
          
          <div className="mt-4 flex items-center text-sm text-gray-500 dark:text-gray-400">
            <ClockIcon className="mr-1.5 h-4 w-4 text-gray-400 dark:text-gray-500" />
            <span>Last active: {formatDate(agent.lastActiveAt)}</span>
          </div>
        </div>
      </div>
    </Link>
  );
};

export default AgentCard;