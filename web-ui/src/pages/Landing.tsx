import React from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'

const Landing: React.FC = () => {
  return (
    <div className="min-h-[calc(100vh-4rem)] flex flex-col">
      {/* Hero Section */}
      <section className="relative flex-1 flex flex-col items-center justify-center px-4 text-center py-20">
        <div className="absolute inset-0 bg-gradient-to-b from-background-900 via-background-900 to-primary-900 opacity-20"></div>
        
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="relative z-10 max-w-4xl mx-auto"
        >
          <h1 className="text-4xl md:text-6xl font-display font-bold text-white leading-tight mb-4">
            Enterprise AI <span className="text-primary-400">Agent Management</span> Platform
          </h1>
          
          <p className="text-xl text-gray-300 mb-8 max-w-2xl mx-auto">
            Build, deploy, and manage intelligent agents with multimodal capabilities for your enterprise applications.
          </p>
          
          <div className="flex flex-col md:flex-row justify-center gap-4">
            <Link to="/dashboard" className="py-3 px-6 bg-primary-600 hover:bg-primary-500 transition-colors rounded-lg text-white font-medium shadow-glow">
              Get Started
            </Link>
            <a href="https://docs.sentinelstacks.com" target="_blank" rel="noopener noreferrer" className="py-3 px-6 bg-background-800 hover:bg-background-700 transition-colors rounded-lg text-white font-medium">
              Documentation
            </a>
          </div>
        </motion.div>
      </section>
      
      {/* Features Section */}
      <section className="relative py-20 px-4 bg-background-800 bg-opacity-70">
        <div className="max-w-7xl mx-auto">
          <motion.div 
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.8, delay: 0.2 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl md:text-4xl font-display font-bold text-white mb-4">Advanced Features</h2>
            <p className="text-lg text-gray-300 max-w-2xl mx-auto">SentinelStacks provides a powerful platform for managing your AI agents with all the capabilities you need.</p>
          </motion.div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[
              {
                title: 'Multimodal Support',
                description: 'Process text, images, and other data types with powerful multimodal AI models.',
                icon: '🖼️',
                delay: 0.1
              },
              {
                title: 'Agent Orchestration',
                description: 'Coordinate multiple agents to solve complex tasks and workflows.',
                icon: '🔄',
                delay: 0.2
              },
              {
                title: 'Secure Deployment',
                description: 'Deploy and monitor agents with enterprise-grade security and compliance.',
                icon: '🔒',
                delay: 0.3
              },
              {
                title: 'Model Integration',
                description: 'Connect to various AI models like OpenAI, Claude, and open source alternatives.',
                icon: '🧠',
                delay: 0.4
              },
              {
                title: 'Analytics & Insights',
                description: 'Track performance, usage, and get actionable insights about your agents.',
                icon: '📊',
                delay: 0.5
              },
              {
                title: 'Custom Workflows',
                description: 'Build automated workflows and integrate with your existing tools.',
                icon: '⚙️',
                delay: 0.6
              }
            ].map((feature, index) => (
              <motion.div 
                key={index}
                className="glass p-6 rounded-lg border-t border-gray-700"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: feature.delay }}
                viewport={{ once: true }}
                whileHover={{ y: -5, boxShadow: '0 10px 20px rgba(0,0,0,0.2)' }}
              >
                <div className="text-4xl mb-4">{feature.icon}</div>
                <h3 className="text-xl font-display font-semibold text-white mb-2">{feature.title}</h3>
                <p className="text-gray-400">{feature.description}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>
      
      {/* Call to Action */}
      <section className="py-20 px-4 text-center relative overflow-hidden">
        <div className="max-w-4xl mx-auto relative z-10">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            whileInView={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.7 }}
            viewport={{ once: true }}
          >
            <h2 className="text-3xl md:text-4xl font-display font-bold text-white mb-6">
              Ready to Supercharge Your AI Workflows?
            </h2>
            <p className="text-xl text-gray-300 mb-8">
              Get started with SentinelStacks today and take control of your AI agents.
            </p>
            <Link to="/dashboard" className="py-3 px-8 bg-primary-600 hover:bg-primary-500 transition-colors rounded-lg text-white font-medium shadow-glow text-lg">
              Launch Dashboard
            </Link>
          </motion.div>
        </div>
        
        {/* Background decoration */}
        <div className="absolute -bottom-10 -right-10 w-64 h-64 bg-primary-500 rounded-full filter blur-[100px] opacity-20"></div>
        <div className="absolute -top-10 -left-10 w-64 h-64 bg-secondary-500 rounded-full filter blur-[100px] opacity-20"></div>
      </section>
    </div>
  )
}

export default Landing 