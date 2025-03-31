import React from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'

const NotFound: React.FC = () => {
  return (
    <div className="min-h-[calc(100vh-4rem)] flex flex-col items-center justify-center px-4 text-center">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
      >
        <h1 className="text-9xl font-display font-bold text-primary-500 mb-2">404</h1>
        <h2 className="text-3xl font-display text-white mb-6">Page Not Found</h2>
        <p className="text-xl text-gray-400 mb-8 max-w-md mx-auto">
          The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.
        </p>
        <Link 
          to="/"
          className="py-3 px-6 bg-primary-600 hover:bg-primary-500 transition-colors rounded-lg text-white font-medium shadow-glow inline-block"
        >
          Return Home
        </Link>
      </motion.div>
      
      {/* Decorative elements */}
      <motion.div 
        className="absolute top-1/3 left-1/4 w-12 h-12 rounded-full border-4 border-primary-600"
        animate={{ 
          y: [0, -50, 0],
          opacity: [0.7, 0.2, 0.7],
          scale: [1, 1.2, 1]
        }}
        transition={{ 
          repeat: Infinity, 
          duration: 5,
          ease: "easeInOut"
        }}
      />
      
      <motion.div 
        className="absolute bottom-1/3 right-1/4 w-8 h-8 rounded-full border-4 border-secondary-600"
        animate={{ 
          y: [0, 40, 0],
          opacity: [0.5, 0.2, 0.5],
          scale: [1, 1.1, 1]
        }}
        transition={{ 
          repeat: Infinity, 
          duration: 4,
          ease: "easeInOut",
          delay: 1
        }}
      />
      
      <motion.div 
        className="absolute bottom-1/4 left-1/3 w-6 h-6 rounded-full border-4 border-accent-600"
        animate={{ 
          y: [0, 30, 0],
          opacity: [0.6, 0.3, 0.6],
          scale: [1, 1.3, 1]
        }}
        transition={{ 
          repeat: Infinity, 
          duration: 6,
          ease: "easeInOut",
          delay: 2
        }}
      />
    </div>
  )
}

export default NotFound 