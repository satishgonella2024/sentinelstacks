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

interface AgentListItemProps {
  agent: Agent;
  onStatusChange?: () => void;
}

const AgentListItem: React.FC<AgentListItemProps> = ({ agent, onStatusChange }) => {
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
    <Link to={`/agents/${agent.id}`}>
      <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg hover:shadow-md transition-shadow duration-300">
        <div className="px-4 py-4 sm:px-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <h3 className="text-md font-medium text-gray-900 dark:text-white">
                {agent.name}
              </h3>
              <span className="ml-2 text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">
                {agent.model.name}
              </span>
            </div>
            <div className="flex items-center space-x-2">
              <StatusBadge status={agent.status} />
              <div className="flex space-x-1">
                {agent.status === 'active' ? (
                  <button
                    onClick={handleStopAgent}
                    className="p-1 rounded-full text-red-600 hover:bg-red-100 dark:hover:bg-red-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                  >
                    <StopIcon className="h-4 w-4" />
                  </button>
                ) : (
                  <button
                    onClick={handleStartAgent}
                    className="p-1 rounded-full text-green-600 hover:bg-green-100 dark:hover:bg-green-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                  >
                    <PlayIcon className="h-4 w-4" />
                  </button>
                )}
                <Link
                  to={`/agents/${agent.id}/edit`}
                  className="p-1 rounded-full text-blue-600 hover:bg-blue-100 dark:hover:bg-blue-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  onClick={(e) => e.stopPropagation()}
                >
                  <PencilIcon className="h-4 w-4" />
                </Link>
              </div>
            </div>
          </div>
          <div className="mt-2 flex justify-between">
            <div className="text-sm text-gray-500 dark:text-gray-400 truncate">
              {agent.description}
            </div>
            <div className="flex items-center text-xs text-gray-500 dark:text-gray-400">
              <ClockIcon className="mr-1 h-3 w-3" />
              {formatDate(agent.lastActiveAt)}
            </div>
          </div>
          
          <div className="mt-2">
            <div className="flex flex-wrap gap-1">
              {agent.tools.slice(0, 3).map(tool => (
                <AgentToolBadge key={tool.id} tool={tool} />
              ))}
              {agent.tools.length > 3 && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200">
                  +{agent.tools.length - 3} more
                </span>
              )}
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
};

export default AgentListItem;