import React from 'react'

const LoadingScreen: React.FC = () => {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-background-900 bg-opacity-80 z-50">
      <div className="text-center">
        <div className="relative w-32 h-32 mb-4 mx-auto">
          {/* Outer circle */}
          <div className="absolute inset-0 border-t-4 border-b-4 border-primary-500 rounded-full animate-spin [animation-duration:3s]"></div>
          
          {/* Middle circle */}
          <div className="absolute inset-4 border-t-4 border-b-4 border-secondary-500 rounded-full animate-spin [animation-duration:2s] [animation-direction:reverse]"></div>
          
          {/* Inner circle */}
          <div className="absolute inset-8 border-t-4 border-b-4 border-accent-500 rounded-full animate-spin [animation-duration:1s]"></div>
          
          {/* Center dot */}
          <div className="absolute inset-[40%] bg-white rounded-full animate-pulse"></div>
        </div>
        
        <h2 className="text-2xl font-display text-white mb-1">Loading</h2>
        <p className="text-sm text-gray-300">Initializing SentinelStacks</p>
      </div>
    </div>
  )
}

export default LoadingScreen 