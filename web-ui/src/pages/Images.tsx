import React from 'react'
import { motion } from 'framer-motion'

// Define types based on the API response
interface ImageInfo {
  id: string;
  name: string;
  tag: string;
  created_at: string;
  size: number;
  llm: string;
  parameters?: Record<string, any>;
}

interface ImagesResponse {
  images: ImageInfo[];
}

const Images: React.FC = () => {
  // Mock images data since we don't have the actual useGetImagesQuery hook
  const mockImages: ImageInfo[] = [
    {
      id: "sha256:abcdef1234567890",
      name: "user/chatbot",
      tag: "latest",
      created_at: new Date().toISOString(),
      size: 5242880,
      llm: "claude-3-haiku-20240307"
    },
    {
      id: "sha256:9876543210abcdef",
      name: "user/research-assistant",
      tag: "v1.0",
      created_at: new Date(Date.now() - 86400000).toISOString(),
      size: 8388608,
      llm: "claude-3-opus-20240229"
    }
  ];
  
  const isLoading = false;
  const error: Error | null = null;
  const images = mockImages;
  
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
  
  if (isLoading) return <div className="p-4">Loading images...</div>
  
  if (error) {
    return (
      <div className="p-4">
        <div className="text-red-500 mb-2">Error loading images</div>
        <div className="text-sm bg-gray-800 p-4 rounded">
          {error.message || "Unknown error occurred"}
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
        <h1 className="text-3xl font-display text-white mb-2">Agent Images</h1>
        <p className="text-gray-400">Browse and manage available agent images</p>
      </div>
      
      <div className="mb-6 flex justify-between items-center">
        <div>
          <span className="text-gray-400">{images.length} images available</span>
        </div>
        
        <button className="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg transition-colors">
          Import Image
        </button>
      </div>
      
      {images.length === 0 ? (
        <div className="glass p-8 rounded-lg text-center">
          <h2 className="text-xl text-white mb-4">No images found</h2>
          <p className="text-gray-400 mb-6">No agent images are currently available.</p>
        </div>
      ) : (
        <motion.div 
          className="grid grid-cols-1 md:grid-cols-2 gap-6"
          variants={container}
          initial="hidden"
          animate="show"
        >
          {images.map((image: ImageInfo) => (
            <motion.div key={image.id} variants={item}>
              <div className="glass p-6 rounded-lg hover:shadow-lg transition-all duration-300">
                <div className="flex justify-between items-start mb-4">
                  <h3 className="text-xl font-semibold text-white">{image.name}</h3>
                  <span className="px-2 py-1 bg-gray-700 rounded-full text-xs">{image.tag}</span>
                </div>
                
                <div className="grid grid-cols-2 gap-2 mb-4">
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">LLM</span>
                    {image.llm}
                  </div>
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">Created</span>
                    {new Date(image.created_at).toLocaleDateString()}
                  </div>
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">Size</span>
                    {(image.size / (1024 * 1024)).toFixed(2)} MB
                  </div>
                  <div className="text-xs text-gray-500">
                    <span className="block text-gray-400">ID</span>
                    <span className="truncate">{image.id.substring(0, 12)}</span>
                  </div>
                </div>
                
                <div className="flex space-x-2">
                  <button className="flex-1 px-3 py-2 bg-primary-600 hover:bg-primary-500 text-white text-sm rounded transition-colors">
                    Run Agent
                  </button>
                  <button className="flex-1 px-3 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded transition-colors">
                    Details
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

export default Images 