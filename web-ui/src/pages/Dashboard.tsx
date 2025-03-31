import React, { useState, useRef, useEffect } from 'react'
import { useSelector } from 'react-redux'
import { useGetAgentsQuery } from '@services/api'
import { RootState } from '@context/store'
import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'

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
  const { data: agentsData, isLoading, error } = useGetAgentsQuery()
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [chatMessage, setChatMessage] = useState('')
  const [chatHistory, setChatHistory] = useState<{role: 'user' | 'assistant', content: string, timestamp: Date}[]>([])
  const [chatAgent, setChatAgent] = useState<string | null>(null)
  const [isTyping, setIsTyping] = useState(false)
  const chatEndRef = useRef<HTMLDivElement>(null)
  
  // Scroll chat to bottom when messages change
  useEffect(() => {
    if (chatEndRef.current) {
      chatEndRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [chatHistory])
  
  // Extract agents array from response
  const agents = agentsData?.agents || [];
  
  console.log('Dashboard agents data:', agents)
  console.log('Dashboard error:', error)
  
  // Use mock data if there's an error (especially 500 errors)
  const displayAgents = error ? 
    [
      {
        id: '1',
        name: 'Assistant Bot',
        description: 'General purpose assistant for everyday tasks',
        model: 'gpt-4',
        image: 'openai/gpt-4:latest',
        status: 'active',
        created: new Date().toISOString(),
        lastActive: new Date().toISOString(),
        systemPrompt: 'You are a helpful assistant.',
        isMultimodal: false
      },
      {
        id: '2',
        name: 'Image Analyzer',
        description: 'Specialized in analyzing images and providing descriptions',
        model: 'claude-3-opus-20240229',
        image: 'anthropic/claude-3-opus:latest',
        status: 'idle',
        created: new Date(Date.now() - 86400000).toISOString(),
        lastActive: new Date(Date.now() - 3600000).toISOString(),
        systemPrompt: 'You analyze images and provide detailed descriptions.',
        isMultimodal: true
      }
    ] : agents || []
  
  const filteredAgents = displayAgents.filter(agent => 
    filterStatus === 'all' || agent.status === filterStatus
  )
  
  const totalAgents = displayAgents.length
  const activeAgents = displayAgents.filter(a => a.status === 'active').length
  const multimodalAgents = displayAgents.filter(a => a.isMultimodal).length
  
  // Handle sending chat messages
  const handleSendMessage = () => {
    if (!chatMessage.trim() || !chatAgent) return
    
    // Add user message
    const userMessage = {
      role: 'user' as const,
      content: chatMessage,
      timestamp: new Date()
    }
    
    setChatHistory(prev => [...prev, userMessage])
    setChatMessage('')
    setIsTyping(true)
    
    // Simulate AI response after a delay
    setTimeout(() => {
      const botMessage = {
        role: 'assistant' as const,
        content: getAIResponse(chatMessage, chatAgent),
        timestamp: new Date()
      }
      
      setChatHistory(prev => [...prev, botMessage])
      setIsTyping(false)
    }, 1000 + Math.random() * 2000) // Random delay between 1-3 seconds
  }
  
  // Generate mock AI responses
  const getAIResponse = (message: string, agentId: string) => {
    const agent = displayAgents.find(a => a.id === agentId)
    
    if (!agent) return "I'm sorry, I can't process your request right now."
    
    const lowercaseMessage = message.toLowerCase()
    
    // Some basic response patterns
    if (lowercaseMessage.includes('hello') || lowercaseMessage.includes('hi')) {
      return `Hello! I'm ${agent.name}. How can I assist you today?`
    }
    
    if (lowercaseMessage.includes('help') || lowercaseMessage.includes('what can you do')) {
      return `I'm ${agent.name}, an AI assistant built with ${agent.model}. I can help answer questions, provide information, and assist with various tasks.`
    }
    
    if (lowercaseMessage.includes('thank')) {
      return "You're welcome! Is there anything else I can help with?"
    }
    
    if (lowercaseMessage.includes('weather')) {
      return "I don't have real-time access to weather data, but I can help you find a weather service or discuss weather-related topics."
    }
    
    if (lowercaseMessage.includes('image') || lowercaseMessage.includes('picture')) {
      return agent.isMultimodal 
        ? "I can process and discuss images. You can upload an image and I'll analyze it for you." 
        : "I don't have image processing capabilities, but I can discuss images conceptually."
    }
    
    // Generic responses
    const genericResponses = [
      "That's an interesting question. Let me think about it...",
      "I understand your query. Based on my knowledge, I would suggest...",
      "Thank you for sharing that information. Would you like me to provide more details or suggestions?",
      "I can definitely help with that. Here's what I know about this topic...",
      "That's a great point. I'd add that there are several perspectives to consider..."
    ]
    
    return genericResponses[Math.floor(Math.random() * genericResponses.length)]
  }
  
  // Function to get agent name from ID
  const getAgentName = (id: string) => {
    const agent = displayAgents.find(a => a.id === id)
    return agent ? agent.name : 'Unknown Agent'
  }
  
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
  
  // In case of server error with status 500, still render the dashboard with mock data
  if (error && !('status' in error && error.status === 500)) {
    console.error('Dashboard error details:', error)
    return (
      <div className="p-4">
        <div className="text-red-500 mb-2">Error loading dashboard data</div>
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
  
  return (
    <div className="p-4 max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-display text-white mb-2">Dashboard</h1>
        <p className="text-gray-400">Overview of your SentinelStacks agents and activity</p>
      </div>
      
      {/* Main Dashboard Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        {/* Left Column - Stats and Agents */}
        <div className="lg:col-span-2 space-y-6">
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
                className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-4"
                variants={container}
                initial="hidden"
                animate="show"
              >
                {filteredAgents.slice(0, 6).map((agent) => (
                  <motion.div key={agent.id} variants={item}>
                    <div className="glass p-4 rounded-lg cursor-pointer hover:border hover:border-primary-500 transition-all"
                      onClick={() => setChatAgent(agent.id)}
                    >
                      <div className="flex justify-between items-start mb-3">
                        <h3 className="text-white font-semibold">{agent.name}</h3>
                        <div className={`w-2 h-2 rounded-full ${agent.status === 'active' ? 'bg-green-500' : agent.status === 'error' ? 'bg-red-500' : 'bg-gray-500'}`}></div>
                      </div>
                      <div className="text-sm text-gray-400 mb-2">Model: {agent.model}</div>
                      <div className="text-xs text-gray-500 mb-3">Last active: {new Date(agent.lastActive).toLocaleString()}</div>
                      <div className="flex space-x-2">
                        <button className="px-3 py-1 text-xs bg-primary-600 hover:bg-primary-500 text-white rounded transition-colors"
                          onClick={(e) => {
                            e.stopPropagation();
                            setChatAgent(agent.id);
                          }}
                        >
                          Chat
                        </button>
                        <Link to="/agents" className="px-3 py-1 text-xs bg-gray-700 hover:bg-gray-600 text-white rounded transition-colors">
                          Details
                        </Link>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </motion.div>
            )}
            
            <div className="mt-4 text-center">
              <Link to="/agents" className="text-primary-400 hover:text-primary-300 text-sm">
                View all agents →
              </Link>
            </div>
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
        
        {/* Right Column - Chat Interface */}
        <div className="lg:col-span-1">
          <div className="glass rounded-lg flex flex-col h-[600px]">
            <div className="p-4 border-b border-gray-800 flex justify-between items-center">
              <h2 className="text-xl font-display text-white">
                {chatAgent ? `Chat with ${getAgentName(chatAgent)}` : 'Quick Chat'}
              </h2>
              {chatAgent && (
                <button 
                  onClick={() => setChatAgent(null)}
                  className="text-gray-400 hover:text-white text-sm bg-gray-800 px-2 py-1 rounded"
                >
                  Change Agent
                </button>
              )}
            </div>
            
            {!chatAgent ? (
              <div className="flex-1 flex flex-col items-center justify-center p-6 text-center">
                <div className="text-5xl mb-4">👋</div>
                <h3 className="text-white text-lg mb-2">Select an Agent to Chat</h3>
                <p className="text-gray-400 text-sm mb-6">Choose an agent from the list to start a conversation</p>
                <div className="grid grid-cols-2 gap-2 w-full max-w-xs">
                  {displayAgents.slice(0, 4).map(agent => (
                    <button
                      key={agent.id}
                      onClick={() => setChatAgent(agent.id)}
                      className="px-3 py-2 bg-primary-600 hover:bg-primary-500 text-white text-sm rounded transition-colors"
                    >
                      {agent.name}
                    </button>
                  ))}
                </div>
              </div>
            ) : (
              <>
                <div className="flex-1 overflow-y-auto p-4">
                  {chatHistory.length === 0 ? (
                    <div className="text-center text-gray-500 py-8">
                      <p>Start a conversation with {getAgentName(chatAgent)}</p>
                    </div>
                  ) : (
                    <div className="space-y-4">
                      {chatHistory.map((msg, index) => (
                        <div key={index} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                          <div className={`max-w-[70%] rounded-lg p-3 ${msg.role === 'user' ? 'bg-primary-600 text-white' : 'bg-gray-800 text-white'}`}>
                            <p>{msg.content}</p>
                            <p className="text-xs opacity-70 mt-1">{msg.timestamp.toLocaleTimeString()}</p>
                          </div>
                        </div>
                      ))}
                      {isTyping && (
                        <div className="flex justify-start">
                          <div className="bg-gray-800 text-white rounded-lg p-3 max-w-[70%]">
                            <div className="flex space-x-2">
                              <div className="w-2 h-2 rounded-full bg-gray-500 animate-bounce"></div>
                              <div className="w-2 h-2 rounded-full bg-gray-500 animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                              <div className="w-2 h-2 rounded-full bg-gray-500 animate-bounce" style={{ animationDelay: '0.4s' }}></div>
                            </div>
                          </div>
                        </div>
                      )}
                      <div ref={chatEndRef} />
                    </div>
                  )}
                </div>
                <div className="p-4 border-t border-gray-800">
                  <form 
                    onSubmit={(e) => {
                      e.preventDefault();
                      handleSendMessage();
                    }}
                    className="flex space-x-2"
                  >
                    <input
                      type="text"
                      value={chatMessage}
                      onChange={(e) => setChatMessage(e.target.value)}
                      placeholder="Type your message..."
                      className="flex-1 bg-gray-800 border border-gray-700 rounded-lg p-2 text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                      disabled={isTyping}
                    />
                    <button
                      type="submit"
                      className={`px-4 py-2 ${isTyping ? 'bg-gray-700 cursor-not-allowed' : 'bg-primary-600 hover:bg-primary-500'} text-white rounded-lg transition-colors`}
                      disabled={isTyping || !chatMessage.trim()}
                    >
                      Send
                    </button>
                  </form>
                </div>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default Dashboard
