import React, { useState, useRef, useEffect } from 'react'
import { motion } from 'framer-motion'
import { useCreateAgentMutation } from '@/services/api'
import { useNavigate } from 'react-router-dom'
import { FiInfo, FiSend, FiCpu, FiFeather, FiCode, FiBookOpen, FiChevronRight } from 'react-icons/fi'

import AgentTemplatesShowcase from '@/components/agents/AgentTemplatesShowcase'

// Example template prompts for quick selection
const QUICK_TEMPLATES = [
  {
    name: 'Customer Support Assistant',
    description: 'An agent that helps with customer inquiries and troubleshooting.',
    prompt: 'You are a customer support agent for SentinelStacks. Help users with their questions about our platform, troubleshoot issues, and provide precise, friendly guidance. If you don\'t know the answer, acknowledge that and offer to escalate the issue to a human support team.'
  },
  {
    name: 'Code Assistant',
    description: 'Specialized in software development assistance and code review.',
    prompt: 'You are a senior software developer. Help users write clean, efficient code. Suggest improvements to their code, provide debugging help, explain programming concepts, and answer technical questions. Focus on providing actionable, practical advice.'
  },
  {
    name: 'Content Writer',
    description: 'Creates and refines various types of written content.',
    prompt: 'You are a versatile content writer. Help users create and refine various types of written content including blog posts, marketing copy, emails, and social media posts. Focus on clarity, engagement, and appropriate tone for the target audience.'
  }
]

const Builder: React.FC = () => {
  const [createAgent, { isLoading }] = useCreateAgentMutation()
  const navigate = useNavigate()
  const inputRef = useRef<HTMLTextAreaElement>(null)
  
  // Builder mode (Natural Language or Templates)
  const [builderMode, setBuilderMode] = useState<'nlp' | 'templates'>('nlp')
  
  // For NLP interface
  const [prompt, setPrompt] = useState('')
  const [submittedPrompt, setSubmittedPrompt] = useState(false)
  const [isAnalyzing, setIsAnalyzing] = useState(false)
  const [selectedTemplate, setSelectedTemplate] = useState(-1)
  
  // Agent data (will be populated from NLP analysis or template selection)
  const [agentData, setAgentData] = useState({
    name: '',
    description: '',
    model: 'gpt-4',
    systemPrompt: '',
    isMultimodal: false
  })
  
  // Step state (for stepper display)
  const [currentStep, setCurrentStep] = useState(1)
  
  // Process natural language input to extract agent details
  const processNaturalLanguage = (input: string) => {
    setIsAnalyzing(true)
    
    // Simulate NLP processing with a delay
    setTimeout(() => {
      let name = ''
      let description = ''
      let model = 'gpt-4'
      let systemPrompt = input
      let isMultimodal = false
      
      // Extract name (look for phrases like "create a", "build a", "named", "called")
      const namePatterns = [
        /create (?:a|an) ([a-z0-9 -]+) agent/i,
        /build (?:a|an) ([a-z0-9 -]+) agent/i,
        /(?:agent|assistant) (?:named|called) "?([a-z0-9 -]+)"?/i,
        /([a-z0-9 -]+) (?:agent|assistant) that/i
      ]
      
      for (const pattern of namePatterns) {
        const match = input.match(pattern)
        if (match && match[1]) {
          name = match[1].trim()
          // Capitalize first letter of each word
          name = name.split(' ').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
          if (name) break
        }
      }
      
      // If no name was found, create a generic one
      if (!name) {
        if (input.toLowerCase().includes('code') || input.toLowerCase().includes('programming')) {
          name = 'Code Assistant'
        } else if (input.toLowerCase().includes('customer') || input.toLowerCase().includes('support')) {
          name = 'Customer Support'
        } else if (input.toLowerCase().includes('content') || input.toLowerCase().includes('writing')) {
          name = 'Content Writer'
        } else {
          name = 'AI Assistant'
        }
      }
      
      // Extract description
      if (input.length > 30) {
        description = input.substring(0, 120).trim()
        if (description.endsWith('.')) {
          description = description.substring(0, description.lastIndexOf('.') + 1)
        } else {
          description += '...'
        }
      } else {
        description = `A specialized agent for ${name.toLowerCase()} tasks.`
      }
      
      // Detect if model should be multimodal
      if (
        input.toLowerCase().includes('image') || 
        input.toLowerCase().includes('visual') || 
        input.toLowerCase().includes('photo') ||
        input.toLowerCase().includes('picture')
      ) {
        isMultimodal = true
        
        // If multimodal, suggest a capable model
        if (input.toLowerCase().includes('high quality') || input.toLowerCase().includes('accurate')) {
          model = 'claude-3-opus-20240229' // Best quality multimodal
        } else {
          model = 'claude-3-sonnet-20240229' // Good balance multimodal
        }
      } else if (input.toLowerCase().includes('code') || input.toLowerCase().includes('programming')) {
        model = 'gpt-4' // Good for code
      } else if (input.toLowerCase().includes('reasoning') || input.toLowerCase().includes('thinking')) {
        model = 'claude-3-opus-20240229' // Good for reasoning
      } else if (input.toLowerCase().includes('fast') || input.toLowerCase().includes('quick')) {
        model = 'llama-3-70b-instruct' // Faster response times
      }
      
      setAgentData({
        name,
        description,
        model,
        systemPrompt: input,
        isMultimodal
      })
      
      setIsAnalyzing(false)
      setSubmittedPrompt(true)
      setCurrentStep(2)
    }, 1500)
  }
  
  const handleNaturalLanguageSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (prompt.trim()) {
      processNaturalLanguage(prompt.trim())
    }
  }
  
  const handleQuickTemplateSelect = (index: number) => {
    setSelectedTemplate(index)
    const template = QUICK_TEMPLATES[index]
    setPrompt(template.prompt)
    setAgentData(prev => ({
      ...prev,
      name: template.name,
      description: template.description,
      systemPrompt: template.prompt
    }))
  }
  
  const handleFullTemplateSelect = (template: any) => {
    setAgentData({
      name: template.title,
      description: template.description,
      model: template.model,
      systemPrompt: template.prompt,
      isMultimodal: template.isMultimodal
    })
    setCurrentStep(2)
  }
  
  const handleEditField = (field: keyof typeof agentData, value: string | boolean) => {
    setAgentData(prev => ({
      ...prev,
      [field]: value
    }))
  }
  
  const handleSubmit = async () => {
    try {
      setCurrentStep(3) // Move to final step
      const result = await createAgent(agentData).unwrap()
      console.log('Agent created:', result)
      
      // Navigate to agents page after successful creation
      setTimeout(() => {
        navigate('/agents')
      }, 1500)
    } catch (error) {
      console.error('Failed to create agent:', error)
      alert('Failed to create agent. Please try again.')
      setCurrentStep(2) // Go back to edit step
    }
  }
  
  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1,
      transition: { 
        when: "beforeChildren",
        staggerChildren: 0.1
      }
    }
  }
  
  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: { y: 0, opacity: 1 }
  }
  
  // Focus the input when component mounts if in NLP mode
  useEffect(() => {
    if (builderMode === 'nlp' && inputRef.current) {
      inputRef.current.focus()
    }
  }, [builderMode])
  
  return (
    <div className="p-4 max-w-5xl mx-auto">
      <motion.div 
        className="mb-8"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <h1 className="text-4xl font-display text-white mb-3">Agent Builder</h1>
        <p className="text-xl text-primary-400">Create a specialized AI agent for your organization</p>
      </motion.div>
      
      {/* Builder Mode Selector */}
      {currentStep === 1 && (
        <div className="flex mb-8 rounded-xl overflow-hidden">
          <button 
            className={`flex-1 py-3 px-4 ${builderMode === 'nlp' 
              ? 'bg-primary-600 text-white' 
              : 'bg-gray-800 text-gray-400 hover:bg-gray-700'}`}
            onClick={() => setBuilderMode('nlp')}
          >
            Natural Language
          </button>
          <button 
            className={`flex-1 py-3 px-4 ${builderMode === 'templates' 
              ? 'bg-primary-600 text-white' 
              : 'bg-gray-800 text-gray-400 hover:bg-gray-700'}`}
            onClick={() => setBuilderMode('templates')}
          >
            Templates
          </button>
        </div>
      )}
      
      {/* Stepper */}
      <div className="mb-10">
        <div className="flex items-center justify-between relative">
          <div className="absolute h-1 bg-gray-700 w-full top-1/2 -translate-y-1/2 -z-10"></div>
          
          {[1, 2, 3].map((step) => (
            <div key={step} className="flex flex-col items-center">
              <div 
                className={`w-10 h-10 rounded-full flex items-center justify-center mb-2 
                  ${currentStep > step ? 'bg-green-500 text-white'
                  : currentStep === step ? 'bg-primary-500 ring-4 ring-primary-500/30 text-white' 
                  : 'bg-gray-800 text-gray-400'}`}
              >
                {currentStep > step ? '✓' : step}
              </div>
              <div className={`text-sm ${currentStep >= step ? 'text-white' : 'text-gray-500'}`}>
                {step === 1 ? 'Select' : step === 2 ? 'Configure' : 'Create'}
              </div>
            </div>
          ))}
        </div>
      </div>
      
      <div className="glass backdrop-blur-glass p-8 rounded-xl border border-gray-700/50 shadow-glass">
        {/* Step 1: Select Mode (NLP or Templates) */}
        {currentStep === 1 && builderMode === 'nlp' && (
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
          >
            <motion.h2 
              className="text-2xl font-display text-white mb-6"
              variants={itemVariants}
            >
              Describe Your Agent
            </motion.h2>
            
            <motion.p 
              className="text-gray-300 mb-8"
              variants={itemVariants}
            >
              Tell me what kind of AI agent you want to create. Describe its purpose, behavior, and capabilities in natural language.
            </motion.p>
            
            <motion.form onSubmit={handleNaturalLanguageSubmit} variants={itemVariants}>
              <div className="mb-6 relative">
                <textarea
                  ref={inputRef}
                  value={prompt}
                  onChange={(e) => setPrompt(e.target.value)}
                  placeholder="e.g., Create a customer support agent that helps users with product issues, can escalate to human agents, and maintains a friendly, professional tone..."
                  className="w-full bg-gray-800/80 border border-gray-700 rounded-lg p-4 text-white focus:outline-none focus:ring-2 focus:ring-primary-500 min-h-32 resize-none"
                ></textarea>
                
                <button 
                  type="submit"
                  disabled={!prompt.trim() || isAnalyzing}
                  className={`absolute bottom-4 right-4 rounded-full p-2 
                    ${prompt.trim() && !isAnalyzing 
                      ? 'bg-primary-500 hover:bg-primary-400 text-white' 
                      : 'bg-gray-700 text-gray-500 cursor-not-allowed'}`}
                  aria-label="Submit"
                >
                  {isAnalyzing ? (
                    <div className="w-6 h-6 rounded-full border-t-2 border-primary-300 animate-spin"></div>
                  ) : (
                    <FiSend size={18} />
                  )}
                </button>
              </div>
            </motion.form>
            
            <motion.div className="mt-8" variants={itemVariants}>
              <h3 className="text-lg font-medium text-white mb-4">Or start with a quick template:</h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {QUICK_TEMPLATES.map((template, index) => (
                  <div 
                    key={index}
                    onClick={() => handleQuickTemplateSelect(index)}
                    className={`p-4 rounded-lg border cursor-pointer transition-all
                      ${selectedTemplate === index 
                        ? 'border-primary-500 bg-primary-500/10 shadow-glow' 
                        : 'border-gray-700 bg-gray-800/60 hover:bg-gray-800'}`}
                  >
                    <h4 className="text-base font-medium text-white mb-1">{template.name}</h4>
                    <p className="text-sm text-gray-400">{template.description}</p>
                  </div>
                ))}
              </div>
            </motion.div>
          </motion.div>
        )}
        
        {/* Step 1: Template Selection */}
        {currentStep === 1 && builderMode === 'templates' && (
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
          >
            <AgentTemplatesShowcase onSelectTemplate={handleFullTemplateSelect} />
          </motion.div>
        )}
        
        {/* Step 2: Review and Edit */}
        {currentStep === 2 && (
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
          >
            <motion.h2 
              className="text-2xl font-display text-white mb-6"
              variants={itemVariants}
            >
              Configure Your Agent
            </motion.h2>
            
            <motion.p 
              className="text-gray-300 mb-6"
              variants={itemVariants}
            >
              Review and customize your agent's details before creation.
            </motion.p>
            
            <motion.div className="space-y-6" variants={itemVariants}>
              <div>
                <label className="block text-gray-400 mb-2">Agent Name</label>
                <input 
                  type="text" 
                  value={agentData.name} 
                  onChange={(e) => handleEditField('name', e.target.value)}
                  className="w-full bg-gray-800/80 border border-gray-700 rounded-lg p-3 text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
              
              <div>
                <label className="block text-gray-400 mb-2">Description</label>
                <textarea 
                  value={agentData.description} 
                  onChange={(e) => handleEditField('description', e.target.value)}
                  className="w-full bg-gray-800/80 border border-gray-700 rounded-lg p-3 text-white focus:outline-none focus:ring-2 focus:ring-primary-500 min-h-20 resize-none"
                ></textarea>
              </div>
              
              <div>
                <label className="block text-gray-400 mb-2">Model</label>
                <select 
                  value={agentData.model} 
                  onChange={(e) => handleEditField('model', e.target.value)}
                  className="w-full bg-gray-800/80 border border-gray-700 rounded-lg p-3 text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                >
                  <option value="gpt-4">GPT-4</option>
                  <option value="gpt-4-turbo">GPT-4 Turbo</option>
                  <option value="claude-3-opus-20240229">Claude 3 Opus</option>
                  <option value="claude-3-sonnet-20240229">Claude 3 Sonnet</option>
                  <option value="llama-3-70b-instruct">Llama 3 70B</option>
                </select>
              </div>
              
              <div>
                <div className="flex justify-between items-center mb-2">
                  <label className="text-gray-400">System Prompt</label>
                  <div className="flex items-center text-sm text-gray-500">
                    <FiInfo className="mr-1" size={14} />
                    <span>Defines your agent's core behavior</span>
                  </div>
                </div>
                <textarea 
                  value={agentData.systemPrompt} 
                  onChange={(e) => handleEditField('systemPrompt', e.target.value)}
                  className="w-full bg-gray-800/80 border border-gray-700 rounded-lg p-3 text-white focus:outline-none focus:ring-2 focus:ring-primary-500 min-h-32 resize-none"
                ></textarea>
              </div>
              
              <div>
                <label className="flex items-center space-x-2 cursor-pointer">
                  <input 
                    type="checkbox" 
                    checked={agentData.isMultimodal} 
                    onChange={(e) => handleEditField('isMultimodal', e.target.checked)}
                    className="w-5 h-5 bg-gray-800 border border-gray-700 rounded focus:ring-2 focus:ring-primary-500 text-primary-500"
                  />
                  <span className="text-gray-400">Enable image capabilities (multimodal)</span>
                </label>
              </div>
            </motion.div>
            
            <motion.div className="mt-8 flex justify-between" variants={itemVariants}>
              <button 
                onClick={() => setCurrentStep(1)}
                className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
              >
                Back
              </button>
              
              <button 
                onClick={handleSubmit}
                disabled={!agentData.name || !agentData.systemPrompt || isLoading}
                className={`px-6 py-2 ${
                  !agentData.name || !agentData.systemPrompt || isLoading
                    ? 'bg-gray-700 cursor-not-allowed'
                    : 'bg-primary-600 hover:bg-primary-500'
                } text-white rounded-lg transition-colors`}
              >
                {isLoading ? 'Creating...' : 'Create Agent'}
              </button>
            </motion.div>
          </motion.div>
        )}
        
        {/* Step 3: Creation Confirmation */}
        {currentStep === 3 && (
          <motion.div
            className="text-center py-6"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.5 }}
          >
            <div className="flex justify-center mb-6">
              <div className="w-16 h-16 rounded-full bg-green-500/20 flex items-center justify-center text-green-500 animate-pulse">
                <FiCpu size={32} />
              </div>
            </div>
            
            <h2 className="text-2xl font-display text-white mb-4">Creating Your Agent</h2>
            <p className="text-gray-400 max-w-md mx-auto">
              {agentData.name} is being created and configured. You'll be redirected to the Agents page when complete.
            </p>
            
            <div className="mt-6 flex justify-center">
              <div className="w-8 h-8 border-t-2 border-primary-500 rounded-full animate-spin"></div>
            </div>
          </motion.div>
        )}
      </div>
      
      {currentStep === 1 && builderMode === 'nlp' && (
        <motion.div 
          className="mt-8 glass backdrop-blur-glass p-6 rounded-xl border border-gray-700/50"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3, duration: 0.5 }}
        >
          <h3 className="text-lg font-medium text-white mb-3 flex items-center">
            <FiFeather className="mr-2" size={18} />
            Tips for Describing Your Agent
          </h3>
          
          <ul className="space-y-2 text-gray-300">
            <li className="flex items-start">
              <span className="text-primary-400 mr-2">•</span>
              <span>Mention the agent's <strong>purpose</strong> (e.g., "customer support", "content writer")</span>
            </li>
            <li className="flex items-start">
              <span className="text-primary-400 mr-2">•</span>
              <span>Specify <strong>tone and style</strong> (e.g., "professional", "casual", "technical")</span>
            </li>
            <li className="flex items-start">
              <span className="text-primary-400 mr-2">•</span>
              <span>Include <strong>domain knowledge</strong> requirements (e.g., "expertise in JavaScript")</span>
            </li>
            <li className="flex items-start">
              <span className="text-primary-400 mr-2">•</span>
              <span>Indicate if <strong>image capabilities</strong> are needed (for analyzing visual content)</span>
            </li>
          </ul>
        </motion.div>
      )}
    </div>
  )
}

export default Builder