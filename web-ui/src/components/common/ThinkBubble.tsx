import React, { useState } from 'react'
import { AnimatePresence, motion } from 'framer-motion'

interface ThinkBubbleProps {
  position: 'dashboard' | 'chat' | 'builder' | 'explorer'
}

const ThinkBubble: React.FC<ThinkBubbleProps> = ({ position }) => {
  const [isOpen, setIsOpen] = useState(true)
  
  // Content based on position
  const getBubbleContent = () => {
    switch (position) {
      case 'dashboard':
        return {
          title: 'Dashboard Insights',
          text: 'Your agents have processed 125 conversations this week, a 15% increase from last week.',
          icon: '📊'
        }
      case 'chat':
        return {
          title: 'Chat Tips',
          text: 'Try asking your agent to summarize the previous conversation or explain its reasoning.',
          icon: '💬'
        }
      case 'builder':
        return {
          title: 'Agent Builder',
          text: 'Multimodal agents can process both text and images for more versatile conversations.',
          icon: '🔧'
        }
      case 'explorer':
        return {
          title: 'Explorer Guide',
          text: 'You can filter agents by model type or status using the controls above.',
          icon: '🔍'
        }
      default:
        return {
          title: 'SentinelStacks',
          text: 'Welcome to your enterprise AI agent management platform.',
          icon: '✨'
        }
    }
  }
  
  const content = getBubbleContent()
  
  if (!isOpen) return (
    <motion.button
      className="fixed bottom-4 right-4 w-12 h-12 rounded-full bg-primary-500 text-white flex items-center justify-center shadow-glow-sm hover:bg-primary-400 transition-colors"
      onClick={() => setIsOpen(true)}
      initial={{ scale: 0.8, opacity: 0 }}
      animate={{ scale: 1, opacity: 1 }}
      whileHover={{ scale: 1.1 }}
    >
      {content.icon}
    </motion.button>
  )
  
  return (
    <AnimatePresence>
      <motion.div
        className="fixed bottom-4 right-4 w-80 glass p-4 rounded-lg shadow-glow-sm"
        initial={{ y: 50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        exit={{ y: 50, opacity: 0 }}
        transition={{ type: 'spring', stiffness: 500, damping: 30 }}
      >
        <div className="flex justify-between items-start mb-2">
          <h3 className="font-display text-lg text-primary-300">{content.title}</h3>
          <button 
            onClick={() => setIsOpen(false)}
            className="text-gray-400 hover:text-white transition-colors"
          >
            ×
          </button>
        </div>
        
        <p className="text-sm text-gray-300 mb-3">{content.text}</p>
        
        <div className="flex justify-end">
          <button 
            className="text-xs text-primary-300 hover:text-primary-200 transition-colors"
            onClick={() => setIsOpen(false)}
          >
            Dismiss
          </button>
        </div>
      </motion.div>
    </AnimatePresence>
  )
}

export default ThinkBubble 