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
          text: 'Your agents are processing user requests efficiently. Consider adding specialized agents for specific domains to improve response quality.',
          icon: '📊',
          type: 'insight'
        }
      case 'chat':
        return {
          title: 'Chat Tips',
          text: 'Try using more specific prompts to get better results. Include details and context in your questions to help the AI understand what you need.',
          icon: '💬',
          type: 'guidance'
        }
      case 'builder':
        return {
          title: 'Agent Building Tips',
          text: 'Creating a well-defined system prompt is crucial. Be specific about the agent\'s role, tone, and limitations to get the best performance.',
          icon: '🔧',
          type: 'suggestion'
        }
      case 'explorer':
        return {
          title: 'Explorer Guide',
          text: 'You can filter agents by model type or status. Try different models for different tasks - Claude excels at reasoning while GPT-4 handles general tasks well.',
          icon: '🔍',
          type: 'guidance'
        }
      default:
        return {
          title: 'SentinelStacks',
          text: 'Welcome to your enterprise AI agent management platform. Need help? Check our documentation or contact support.',
          icon: '✨',
          type: 'achievement'
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
        className={`fixed bottom-4 right-4 w-80 think-bubble p-4 ${content.type} shadow-glow-sm`}
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