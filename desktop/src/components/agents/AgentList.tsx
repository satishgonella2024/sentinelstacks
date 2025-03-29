import React, { useState } from 'react';
import { 
  PlusIcon, 
  MagnifyingGlassIcon,
  Squares2X2Icon,
  ListBulletIcon
} from '@heroicons/react/24/outline';
import AgentCard from './AgentCard';
import AgentListItem from './AgentListItem';
import LoadingSpinner from '../common/LoadingSpinner';
import { useAgentList } from '../../hooks/useAgentList';
import { Link } from 'react-router-dom';

type ViewMode = 'grid' | 'list';

const AgentList: React.FC = () => {
  const { 
    filteredAgents, 
    isLoading, 
    error, 
    searchTerm, 
    setSearchTerm, 
    statusFilter, 
    setStatusFilter, 
    fetchAgents 
  } = useAgentList();
  
  const [viewMode, setViewMode] = useState<ViewMode>('grid');

  const handleRefresh = () => {
    fetchAgents();
  };

  const renderAgents = () => {
    if (isLoading) {
      return (
        <div className="flex items-center justify-center h-64">
          <LoadingSpinner size="md" />
        </div>
      );
    }

    if (error) {
      return (
        <div className="flex flex-col items-center justify-center h-64 text-gray-500 dark:text-gray-400">
          <p className="mb-4 text-red-500">{error}</p>
          <button 
            onClick={handleRefresh}
            className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Retry
          </button>
        </div>
      );
    }

    if (filteredAgents.length === 0) {
      return (
        <div className="flex flex-col items-center justify-center h-64 text-gray-500 dark:text-gray-400">
          <p className="mb-4">No agents found</p>
          <Link 
            to="/agents/create"
            className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Create Agent
          </Link>
        </div>
      );
    }

    return (
      <div className={viewMode === 'grid' ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6" : "space-y-4"}>
        {filteredAgents.map(agent => (
          viewMode === 'grid' ? (
            <AgentCard key={agent.id} agent={agent} onStatusChange={handleRefresh} />
          ) : (
            <AgentListItem key={agent.id} agent={agent} onStatusChange={handleRefresh} />
          )
        ))}
      </div>
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Agents</h1>
        <div className="flex items-center gap-4">
          <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <MagnifyingGlassIcon className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="text"
              className="block w-full pl-10 pr-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md leading-5 bg-white dark:bg-gray-800 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
              placeholder="Search agents..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
          <div className="flex items-center space-x-2 border border-gray-300 dark:border-gray-700 rounded-md">
            <button
              className={`p-2 ${viewMode === 'grid' ? 'bg-gray-100 dark:bg-gray-700' : 'bg-white dark:bg-gray-800'} rounded-l-md`}
              onClick={() => setViewMode('grid')}
              aria-label="Grid view"
            >
              <Squares2X2Icon className="h-5 w-5 text-gray-500 dark:text-gray-400" />
            </button>
            <button
              className={`p-2 ${viewMode === 'list' ? 'bg-gray-100 dark:bg-gray-700' : 'bg-white dark:bg-gray-800'} rounded-r-md`}
              onClick={() => setViewMode('list')}
              aria-label="List view"
            >
              <ListBulletIcon className="h-5 w-5 text-gray-500 dark:text-gray-400" />
            </button>
          </div>
          <select
            className="block w-full py-2 px-3 border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
          >
            <option value="all">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
            <option value="error">Error</option>
          </select>
          <Link
            to="/agents/create"
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            <PlusIcon className="h-5 w-5 mr-2" />
            New Agent
          </Link>
        </div>
      </div>
      
      {renderAgents()}
    </div>
  );
};

export default AgentList;