import React, { useState, useEffect } from 'react';
import { 
  PlusIcon, 
  MagnifyingGlassIcon,
  ArrowPathIcon,
  Squares2X2Icon,
  ListBulletIcon
} from '@heroicons/react/24/outline';
import AgentCard from './AgentCard';
import AgentListItem from './AgentListItem';
import { Agent } from '../types/Agent';

// Placeholder API service - replace with actual implementation
import { getAgents } from '../services/agentService';

type ViewMode = 'grid' | 'list';

const AgentList: React.FC = () => {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [filteredAgents, setFilteredAgents] = useState<Agent[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [viewMode, setViewMode] = useState<ViewMode>('grid');
  const [statusFilter, setStatusFilter] = useState<string>('all');

  useEffect(() => {
    fetchAgents();
  }, []);

  useEffect(() => {
    filterAgents();
  }, [agents, searchTerm, statusFilter]);

  const fetchAgents = async () => {
    setIsLoading(true);
    try {
      const data = await getAgents();
      setAgents(data);
    } catch (error) {
      console.error('Error fetching agents:', error);
      // Add error handling UI
    } finally {
      setIsLoading(false);
    }
  };

  const filterAgents = () => {
    let filtered = [...agents];
    
    // Apply search filter
    if (searchTerm) {
      filtered = filtered.filter(agent => 
        agent.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        agent.description.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }
    
    // Apply status filter
    if (statusFilter !== 'all') {
      filtered = filtered.filter(agent => agent.status === statusFilter);
    }
    
    setFilteredAgents(filtered);
  };

  const handleCreateAgent = () => {
    // Navigate to agent creation page
    window.location.href = '/agents/create';
  };

  const renderAgents = () => {
    if (isLoading) {
      return (
        <div className="flex items-center justify-center h-64">
          <ArrowPathIcon className="w-8 h-8 animate-spin text-primary-600" />
        </div>
      );
    }

    if (filteredAgents.length === 0) {
      return (
        <div className="flex flex-col items-center justify-center h-64 text-gray-500 dark:text-gray-400">
          <p className="mb-4">No agents found</p>
          <button 
            onClick={handleCreateAgent}
            className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Create Agent
          </button>
        </div>
      );
    }

    return (
      <div className={viewMode === 'grid' ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6" : "space-y-4"}>
        {filteredAgents.map(agent => (
          viewMode === 'grid' ? (
            <AgentCard key={agent.id} agent={agent} />
          ) : (
            <AgentListItem key={agent.id} agent={agent} />
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
          <button
            onClick={handleCreateAgent}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            <PlusIcon className="h-5 w-5 mr-2" />
            New Agent
          </button>
        </div>
      </div>
      
      {renderAgents()}
    </div>
  );
};

export default AgentList;