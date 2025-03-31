import React, { useState } from 'react'
import { useGetAgentsQuery } from '@/services/api'
import { motion } from 'framer-motion'
import { Link, useNavigate } from 'react-router-dom'
import { FiCpu, FiSearch, FiPlus, FiFilter, FiTag, FiMessageCircle, FiSettings, FiTrash2, FiPause, FiPlay } from 'react-icons/fi'

const Agents: React.FC = () => {
  const { data: agentsData, isLoading, error } = useGetAgentsQuery()
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const navigate = useNavigate()
  
  // Extract agents array from response
  const agents = agentsData?.agents || [];
  
  console.log('Agents page data:', agents)
  
  // Get unique categories from agent tags
  const allCategories = ['all']
  agents.forEach(agent => {
    if (agent.tags) {
      agent.tags.forEach(tag => {
        if (!allCategories.includes(tag)) {
          allCategories.push(tag)
        }
      })
    }
  })
  
  // Filter agents by status, search query, and category
  const filteredAgents = agents.filter(agent => {
    // Filter by status
    if (filterStatus !== 'all' && agent.status !== filterStatus) {
      return false
    }
    
    // Filter by search query
    if (searchQuery && !agent.name.toLowerCase().includes(searchQuery.toLowerCase()) && 
        !agent.description.toLowerCase().includes(searchQuery.toLowerCase())) {
      return false
    }
    
    // Filter by category
    if (selectedCategory !== 'all' && (!agent.tags || !agent.tags.includes(selectedCategory))) {
      return false
    }
    
    return true
  })
  
  // Animation variants
  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  }
  
  const item = {
    hidden: { y: 20, opacity: 0 },
    show: { y: 0, opacity: 1, transition: { duration: 0.4 } }
  }
  
  // Loading state
  if (isLoading) return (
    <div className="flex justify-center items-center h-64">
      <div className="flex flex-col items-center">
        <div className="w-12 h-12 border-t-2 border-b-2 border-primary-500 rounded-full animate-spin mb-4"></div>
        <p className="text-gray-400">Loading agents...</p>
      </div>
    </div>
  )
  
  // Handle error state
  if (error) {
    console.error('Agents error details:', error)
    return (
      <div className="p-4">
        <div className="text-red-500 mb-2">Error loading agents</div>
        <div className="text-sm bg-gray-800 p-4 rounded">
          {error instanceof Error 
            ? error.message 
            : JSON.stringify(error, null, 2)}
        </div>
        <button 
          className="mt-4 px-4 py-2 bg-primary-600 text-white rounded"
          onClick={() => window.location.reload()}
        >
          Retry
        </button>
      </div>
    )
  }
  
  // If no agents available
  if (!agents || agents.length === 0) {
    return (
      <div className="p-4 max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-4xl font-display text-white mb-2">Agents</h1>
          <p className="text-gray-400">Manage your AI agents</p>
        </div>
        
        <div className="glass p-12 rounded-xl text-center">
          <div className="w-20 h-20 mx-auto mb-6 rounded-full bg-primary-500/20 flex items-center justify-center">
            <FiCpu size={40} className="text-primary-400" />
          </div>
          <h2 className="text-2xl text-white mb-4">No agents found</h2>
          <p className="text-gray-400 mb-6 max-w-md mx-auto">
            You don't have any agents yet. Create your first agent to get started with powerful AI assistance.
          </p>
          <button 
            className="px-6 py-3 bg-primary-600 hover:bg-primary-500 text-white rounded-lg transition-colors inline-flex items-center"
            onClick={() => navigate('/builder')}
          >
            <FiPlus className="mr-2" />
            Create Agent
          </button>
        </div>
      </div>
    )
  }
  
  return (
    <div className="p-4 max-w-7xl mx-auto">
      <motion.div 
        className="mb-8"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <h1 className="text-4xl font-display text-white mb-3">Agents</h1>
        <p className="text-xl text-primary-400">Manage your AI assistants</p>
      </motion.div>
      
      {/* Search and filter bar */}
      <motion.div 
        className="mb-8 glass p-4 rounded-xl"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1, duration: 0.5 }}
      >
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="relative w-full md:w-1/3">
            <FiSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search agents..."
              className="w-full bg-gray-800/50 border border-gray-700 rounded-lg py-2 pl-10 pr-4 text-white focus:outline-none focus:ring-1 focus:ring-primary-500"
            />
          </div>
          
          <div className="flex items-center space-x-2 overflow-x-auto py-1 md:py-0 md:justify-end">
            <span className="text-gray-400 whitespace-nowrap flex items-center">
              <FiFilter className="mr-1" size={14} />
              Status:
            </span>
            <button 
              className={`px-3 py-1 text-sm rounded-full ${filterStatus === 'all' ? 'bg-primary-600 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700'}`}
              onClick={() => setFilterStatus('all')}
            >
              All
            </button>
            <button 
              className={`px-3 py-1 text-sm rounded-full ${filterStatus === 'active' ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700'}`}
              onClick={() => setFilterStatus('active')}
            >
              Active
            </button>
            <button 
              className={`px-3 py-1 text-sm rounded-full ${filterStatus === 'idle' ? 'bg-gray-600 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700'}`}
              onClick={() => setFilterStatus('idle')}
            >
              Idle
            </button>
            <button 
              className={`px-3 py-1 text-sm rounded-full ${filterStatus === 'error' ? 'bg-red-600 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700'}`}
              onClick={() => setFilterStatus('error')}
            >
              Error
            </button>
          </div>
        </div>
        
        {/* Categories filter */}
        <div className="mt-4 flex flex-wrap gap-2 items-center">
          <span className="text-gray-400 flex items-center mr-1">
            <FiTag className="mr-1" size={14} />
            Categories:
          </span>
          {allCategories.map(category => (
            <button
              key={category}
              className={`px-3 py-1 text-sm rounded-full ${
                selectedCategory === category
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
              }`}
              onClick={() => setSelectedCategory(category)}
            >
              {category === 'all' ? 'All' : category}
            </button>
          ))}
        </div>
      </motion.div>
      
      <div className="mb-6 flex justify-between items-center">
        <div className="text-gray-400">
          {filteredAgents.length} {filteredAgents.length === 1 ? 'agent' : 'agents'} found
        </div>
        
        <button 
          className="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg transition-colors inline-flex items-center"
          onClick={() => navigate('/builder')}
        >
          <FiPlus className="mr-2" />
          Create Agent
        </button>
      </div>
      
      {filteredAgents.length === 0 ? (
        <div className="glass p-8 rounded-xl text-center">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-800 flex items-center justify-center">
            <FiSearch size={24} className="text-gray-500" />
          </div>
          <h3 className="text-xl text-white mb-2">No matching agents</h3>
          <p className="text-gray-400">Try adjusting your filters or search query</p>
        </div>
      ) : (
        <motion.div 
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
          variants={container}
          initial="hidden"
          animate="show"
        >
          {filteredAgents.map((agent) => (
            <motion.div key={agent.id} variants={item}>
              <div className="glass card-hover p-6 rounded-xl border border-gray-700/50 h-full flex flex-col">
                <div className="flex justify-between items-start mb-4">
                  <h3 className="text-xl font-display text-white">{agent.name}</h3>
                  <div className={`w-3 h-3 rounded-full ${
                    agent.status === 'active' ? 'bg-green-500 status-active' : 
                    agent.status === 'error' ? 'bg-red-500 status-error' : 'bg-yellow-500 status-idle'
                  }`}></div>
                </div>
                
                <p className="text-gray-300 mb-4 line-clamp-2 flex-grow">{agent.description}</p>
                
                {/* Tags */}
                {agent.tags && agent.tags.length > 0 && (
                  <div className="mb-4 flex flex-wrap gap-2">
                    {agent.tags.map((tag, i) => (
                      <span key={i} className="text-xs px-2 py-1 rounded-full bg-gray-800 text-gray-400">
                        {tag}
                      </span>
                    ))}
                  </div>
                )}
                
                {/* Capabilities */}
                {agent.capabilities && agent.capabilities.length > 0 && (
                  <div className="mb-4">
                    <h4 className="text-sm text-gray-400 mb-2">Capabilities</h4>
                    <div className="grid grid-cols-2 gap-2 text-xs">
                      {agent.capabilities.map((capability, i) => (
                        <div key={i} className="bg-gray-800/50 rounded px-2 py-1 text-gray-300">
                          {capability}
                        </div>
                      ))}
                    </div>
                  </div>
                )}
                
                <div className="grid grid-cols-2 gap-2 mb-4 mt-auto">
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">Model</span>
                    {agent.model.split('-').slice(0, 2).join(' ')}
                  </div>
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">Created</span>
                    {new Date(agent.created).toLocaleDateString()}
                  </div>
                </div>
                
                <div className="flex space-x-2">
                  <button className="flex-1 px-3 py-2 bg-primary-600 hover:bg-primary-500 text-white text-sm rounded transition-colors flex items-center justify-center">
                    <FiMessageCircle className="mr-1" size={14} />
                    Chat
                  </button>
                  <button className="px-2 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded transition-colors">
                    <FiSettings size={14} />
                  </button>
                  {agent.status === 'active' ? (
                    <button className="px-2 py-2 bg-yellow-600 hover:bg-yellow-500 text-white text-sm rounded transition-colors">
                      <FiPause size={14} />
                    </button>
                  ) : (
                    <button className="px-2 py-2 bg-green-600 hover:bg-green-500 text-white text-sm rounded transition-colors">
                      <FiPlay size={14} />
                    </button>
                  )}
                  <button className="px-2 py-2 bg-red-600 hover:bg-red-500 text-white text-sm rounded transition-colors">
                    <FiTrash2 size={14} />
                  </button>
                </div>
              </div>
            </motion.div>
          ))}
        </motion.div>
      )}
    </div>
  )
}

export default Agents