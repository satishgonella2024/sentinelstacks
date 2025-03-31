import React, { useState } from 'react'
import { useGetAgentsQuery } from '@services/api'
import { motion } from 'framer-motion'

const Agents: React.FC = () => {
  const { data: agents, isLoading, error } = useGetAgentsQuery()
  const [filterStatus, setFilterStatus] = useState<string>('all')
  
  console.log('Agents page data:', agents)
  
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
    show: { y: 0, opacity: 1 }
  }
  
  if (isLoading) return <div className="p-4">Loading agents...</div>
  
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
          <h1 className="text-3xl font-display text-white mb-2">Agents</h1>
          <p className="text-gray-400">Manage your AI agents</p>
        </div>
        
        <div className="glass p-8 rounded-lg text-center">
          <h2 className="text-xl text-white mb-4">No agents found</h2>
          <p className="text-gray-400 mb-6">You don't have any agents yet. Create your first agent to get started.</p>
          <button className="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg transition-colors">
            Create Agent
          </button>
        </div>
      </div>
    )
  }
  
  // Filter agents by status if needed
  const filteredAgents = filterStatus === 'all' 
    ? agents 
    : agents.filter(agent => agent.status === filterStatus)
  
  return (
    <div className="p-4 max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-display text-white mb-2">Agents</h1>
        <p className="text-gray-400">Manage your AI agents</p>
      </div>
      
      <div className="mb-6 flex justify-between items-center">
        <div className="flex space-x-2">
          <button 
            className={`px-3 py-1 text-sm rounded-full ${filterStatus === 'all' ? 'bg-primary-600 text-white' : 'bg-gray-800 text-gray-400'}`}
            onClick={() => setFilterStatus('all')}
          >
            All
          </button>
          <button 
            className={`px-3 py-1 text-sm rounded-full ${filterStatus === 'active' ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-400'}`}
            onClick={() => setFilterStatus('active')}
          >
            Active
          </button>
          <button 
            className={`px-3 py-1 text-sm rounded-full ${filterStatus === 'idle' ? 'bg-gray-600 text-white' : 'bg-gray-800 text-gray-400'}`}
            onClick={() => setFilterStatus('idle')}
          >
            Idle
          </button>
        </div>
        
        <button className="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg transition-colors">
          Create Agent
        </button>
      </div>
      
      {filteredAgents.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          No agents match your current filter
        </div>
      ) : (
        <motion.div 
          className="grid grid-cols-1 md:grid-cols-3 gap-6"
          variants={container}
          initial="hidden"
          animate="show"
        >
          {filteredAgents.map((agent) => (
            <motion.div key={agent.id} variants={item}>
              <div className="glass p-6 rounded-lg hover:shadow-lg transition-all duration-300">
                <div className="flex justify-between items-start mb-4">
                  <h3 className="text-xl font-semibold text-white">{agent.name}</h3>
                  <div className={`w-3 h-3 rounded-full ${
                    agent.status === 'active' ? 'bg-green-500' : 
                    agent.status === 'error' ? 'bg-red-500' : 'bg-gray-500'
                  }`}></div>
                </div>
                
                <p className="text-gray-400 mb-3 line-clamp-2">{agent.description}</p>
                
                <div className="grid grid-cols-2 gap-2 mb-4">
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">Model</span>
                    {agent.model}
                  </div>
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">Created</span>
                    {new Date(agent.created).toLocaleDateString()}
                  </div>
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">Status</span>
                    {agent.status}
                  </div>
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">Last Active</span>
                    {new Date(agent.lastActive).toLocaleDateString()}
                  </div>
                </div>
                
                <div className="flex space-x-2">
                  <button className="flex-1 px-3 py-2 bg-primary-600 hover:bg-primary-500 text-white text-sm rounded transition-colors">
                    Chat
                  </button>
                  <button className="flex-1 px-3 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded transition-colors">
                    Details
                  </button>
                  <button className="px-3 py-2 bg-red-600 hover:bg-red-500 text-white text-sm rounded transition-colors">
                    Stop
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