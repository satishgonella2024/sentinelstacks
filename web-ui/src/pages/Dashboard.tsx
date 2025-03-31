import React, { useState } from 'react'
import { useSelector } from 'react-redux'
import { useGetAgentsQuery } from '@services/api'
import { RootState } from '@context/store'
import { motion } from 'framer-motion'

// Card components
const StatCard = ({ title, value, icon, color }: { title: string; value: string; icon: string; color: string }) => (
  <motion.div 
    className={`glass p-6 rounded-lg ${color}`}
    whileHover={{ y: -5, boxShadow: '0 10px 20px rgba(0,0,0,0.2)' }}
  >
    <div className="flex justify-between items-center">
      <div>
        <h3 className="text-gray-400 text-sm">{title}</h3>
        <p className="text-white text-2xl font-display mt-1">{value}</p>
      </div>
      <div className="text-3xl opacity-80">{icon}</div>
    </div>
  </motion.div>
)

const AgentCard = ({ name, model, status, lastActive }: { name: string; model: string; status: string; lastActive: string }) => {
  const statusColors = {
    idle: 'bg-gray-500',
    active: 'bg-green-500',
    error: 'bg-red-500',
  }
  
  const statusColor = statusColors[status as keyof typeof statusColors] || statusColors.idle
  
  return (
    <motion.div 
      className="glass p-4 rounded-lg"
      whileHover={{ y: -5, boxShadow: '0 10px 20px rgba(0,0,0,0.2)' }}
    >
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-white font-semibold">{name}</h3>
        <div className={`w-2 h-2 rounded-full ${statusColor}`}></div>
      </div>
      <div className="text-sm text-gray-400 mb-2">Model: {model}</div>
      <div className="text-xs text-gray-500">Last active: {new Date(lastActive).toLocaleString()}</div>
    </motion.div>
  )
}

const Dashboard: React.FC = () => {
  const { data: agents, isLoading, error } = useGetAgentsQuery()
  const [filterStatus, setFilterStatus] = useState<string>('all')
  
  const filteredAgents = agents?.filter(agent => 
    filterStatus === 'all' || agent.status === filterStatus
  ) || []
  
  const totalAgents = agents?.length || 0
  const activeAgents = agents?.filter(a => a.status === 'active').length || 0
  const multimodalAgents = agents?.filter(a => a.isMultimodal).length || 0
  
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
  
  if (isLoading) return <div className="p-4">Loading dashboard data...</div>
  if (error) return <div className="p-4 text-red-500">Error loading dashboard data</div>
  
  return (
    <div className="p-4 max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-display text-white mb-2">Dashboard</h1>
        <p className="text-gray-400">Overview of your SentinelStacks agents and activity</p>
      </div>
      
      {/* Stats Cards */}
      <motion.div 
        className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8"
        variants={container}
        initial="hidden"
        animate="show"
      >
        <motion.div variants={item}>
          <StatCard title="Total Agents" value={totalAgents.toString()} icon="🤖" color="border-l-4 border-primary-500" />
        </motion.div>
        <motion.div variants={item}>
          <StatCard title="Active Agents" value={activeAgents.toString()} icon="✅" color="border-l-4 border-green-500" />
        </motion.div>
        <motion.div variants={item}>
          <StatCard title="Multimodal Agents" value={multimodalAgents.toString()} icon="🖼️" color="border-l-4 border-secondary-500" />
        </motion.div>
      </motion.div>
      
      {/* Agent List */}
      <div className="mb-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-display text-white">Recent Agents</h2>
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
        </div>
        
        {filteredAgents.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            No agents match your current filter
          </div>
        ) : (
          <motion.div 
            className="grid grid-cols-1 md:grid-cols-3 gap-4"
            variants={container}
            initial="hidden"
            animate="show"
          >
            {filteredAgents.slice(0, 6).map((agent) => (
              <motion.div key={agent.id} variants={item}>
                <AgentCard
                  name={agent.name}
                  model={agent.model}
                  status={agent.status}
                  lastActive={agent.lastActive}
                />
              </motion.div>
            ))}
          </motion.div>
        )}
      </div>
      
      {/* Activity Feed */}
      <div>
        <h2 className="text-xl font-display text-white mb-4">Recent Activity</h2>
        <div className="glass rounded-lg p-4">
          <div className="space-y-4">
            {[1, 2, 3].map((item) => (
              <div key={item} className="flex items-start space-x-3 pb-4 border-b border-gray-800">
                <div className="w-10 h-10 rounded-full bg-background-800 flex-shrink-0 flex items-center justify-center">
                  <span className="text-sm">🔄</span>
                </div>
                <div>
                  <p className="text-white">Agent ChatAssistant processed a request</p>
                  <p className="text-sm text-gray-500">{new Date(Date.now() - item * 3600000).toLocaleString()}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

export default Dashboard 