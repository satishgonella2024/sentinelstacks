import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { 
  FiCode, FiBookOpen, FiEdit3, FiBarChart2, 
  FiShield, FiImage, FiDollarSign, FiMessageSquare,
  FiHeadphones, FiServer, FiGlobe, FiZap
} from 'react-icons/fi';

import AgentTemplateCard from './AgentTemplateCard';

// Template categories
const categories = [
  { id: 'all', name: 'All' },
  { id: 'productivity', name: 'Productivity' },
  { id: 'creative', name: 'Creative' },
  { id: 'coding', name: 'Coding' },
  { id: 'research', name: 'Research' },
  { id: 'business', name: 'Business' },
  { id: 'specialized', name: 'Specialized' }
];

// Agent templates
const templates = [
  {
    id: 'research-assistant',
    title: 'Research Assistant',
    description: 'An agent that helps with academic research, citations, and finding relevant sources.',
    icon: <FiBookOpen size={18} />,
    categories: ['research', 'productivity'],
    prompt: 'You are a research assistant specializing in academic research. Help users find relevant papers, cite properly in different formats, and design research methodologies. Provide thoughtful analysis of academic papers and research findings.',
    model: 'claude-3-opus-20240229',
    isMultimodal: false,
    isPopular: true
  },
  {
    id: 'code-assistant',
    title: 'Code Assistant',
    description: 'Helps with programming tasks, code reviews, and fixing bugs across multiple languages.',
    icon: <FiCode size={18} />,
    categories: ['coding', 'productivity'],
    prompt: 'You are a code assistant. Help with programming tasks and code review. Provide working, well-documented code examples. Explain concepts clearly and help debug issues efficiently.',
    model: 'gpt-4',
    isMultimodal: false,
    isPopular: true
  },
  {
    id: 'content-writer',
    title: 'Content Writer',
    description: 'Creates and refines various types of written content like blogs, marketing copy, and emails.',
    icon: <FiEdit3 size={18} />,
    categories: ['creative', 'business'],
    prompt: 'You are a versatile content writer. Help users create and refine various types of written content including blog posts, marketing copy, emails, and social media posts. Focus on clarity, engagement, and appropriate tone for the target audience.',
    model: 'claude-3-sonnet-20240229',
    isMultimodal: false,
    isPopular: false
  },
  {
    id: 'data-analyst',
    title: 'Data Analyst',
    description: 'Analyzes data, creates visualizations, and provides statistical insights from datasets.',
    icon: <FiBarChart2 size={18} />,
    categories: ['research', 'business'],
    prompt: 'You are a data analysis specialist. Help users interpret data, create visualizations, perform statistical analysis, and generate insights from datasets. Explain statistical concepts clearly and provide actionable recommendations.',
    model: 'claude-3-sonnet-20240229',
    isMultimodal: false,
    isPopular: false
  },
  {
    id: 'cybersecurity-advisor',
    title: 'Cybersecurity Advisor',
    description: 'Provides guidance on security practices, risk assessments, and threat mitigation.',
    icon: <FiShield size={18} />,
    categories: ['specialized', 'coding'],
    prompt: 'You are a cybersecurity advisor. Provide expert guidance on security best practices, risk assessments, and threat mitigation strategies. Explain complex security concepts in accessible terms and help users make their systems more secure.',
    model: 'gpt-4',
    isMultimodal: false,
    isPopular: false
  },
  {
    id: 'image-analyst',
    title: 'Image Analyst',
    description: 'Analyzes and describes images, identifies objects, and extracts visual information.',
    icon: <FiImage size={18} />,
    categories: ['specialized'],
    prompt: 'You analyze images and provide detailed descriptions. Identify objects, people, settings, actions, text content, style elements, and other visual information. Maintain accuracy and avoid hallucinations.',
    model: 'claude-3-opus-20240229',
    isMultimodal: true,
    isPopular: false
  },
  {
    id: 'financial-advisor',
    title: 'Financial Advisor',
    description: 'Provides financial analysis, investment guidance, and budgeting assistance.',
    icon: <FiDollarSign size={18} />,
    categories: ['business', 'specialized'],
    prompt: 'You are a financial information assistant. Help users understand financial concepts, analyze investment options, and create budgeting plans. Always clarify you are not providing professional financial advice, and users should consult qualified financial advisors for specific situations.',
    model: 'claude-3-opus-20240229',
    isMultimodal: false,
    isPopular: false
  },
  {
    id: 'customer-support',
    title: 'Customer Support',
    description: 'Helps with product inquiries, troubleshooting, and customer service tasks.',
    icon: <FiHeadphones size={18} />,
    categories: ['business', 'productivity'],
    prompt: 'You are a customer support agent. Help users with their product inquiries, troubleshoot common issues, and provide friendly, efficient service. Prioritize user satisfaction while following company policies. Escalate complex issues appropriately.',
    model: 'claude-3-sonnet-20240229',
    isMultimodal: false,
    isPopular: true
  },
  {
    id: 'language-tutor',
    title: 'Language Tutor',
    description: 'Teaches languages, helps with grammar, and provides conversation practice.',
    icon: <FiMessageSquare size={18} />,
    categories: ['creative', 'specialized'],
    prompt: 'You are a language tutor. Help users learn new languages by explaining grammar concepts, providing vocabulary practice, and engaging in conversational practice. Adapt to the learner\'s level and provide supportive feedback to encourage progress.',
    model: 'gpt-4',
    isMultimodal: false,
    isPopular: false
  },
  {
    id: 'system-admin',
    title: 'System Administrator',
    description: 'Provides guidance on server management, networking, and infrastructure issues.',
    icon: <FiServer size={18} />,
    categories: ['coding', 'specialized'],
    prompt: 'You are a system administration assistant. Help users with server management, networking configuration, and infrastructure troubleshooting. Provide clear step-by-step instructions and explain the rationale behind recommended solutions.',
    model: 'gpt-4',
    isMultimodal: false,
    isPopular: false
  },
  {
    id: 'travel-planner',
    title: 'Travel Planner',
    description: 'Helps plan trips, recommend destinations, and create customized itineraries.',
    icon: <FiGlobe size={18} />,
    categories: ['productivity', 'creative'],
    prompt: 'You are a travel planning assistant. Help users plan trips by recommending destinations, creating itineraries, suggesting activities, and providing travel tips. Consider the user\'s preferences, budget, and time constraints when making recommendations.',
    model: 'claude-3-sonnet-20240229',
    isMultimodal: false,
    isPopular: false
  },
  {
    id: 'brainstorming',
    title: 'Brainstorming Partner',
    description: 'Helps generate ideas, solve problems, and think creatively across domains.',
    icon: <FiZap size={18} />,
    categories: ['creative', 'productivity'],
    prompt: 'You are a brainstorming partner. Help users generate ideas, think creatively about problems, and develop innovative solutions. Ask thought-provoking questions, suggest unexpected perspectives, and build on the user\'s ideas to expand their thinking.',
    model: 'claude-3-opus-20240229',
    isMultimodal: false,
    isPopular: false
  }
];

interface AgentTemplatesShowcaseProps {
  onSelectTemplate: (template: any) => void;
}

const AgentTemplatesShowcase: React.FC<AgentTemplatesShowcaseProps> = ({ onSelectTemplate }) => {
  const [selectedCategory, setSelectedCategory] = useState('all');
  
  // Filter templates by selected category
  const filteredTemplates = selectedCategory === 'all'
    ? templates
    : templates.filter(template => template.categories.includes(selectedCategory));
  
  return (
    <div>
      <div className="mb-6">
        <h2 className="text-2xl font-display text-white mb-6">Agent Templates</h2>
        
        {/* Category tabs */}
        <div className="flex overflow-x-auto pb-2 gap-2">
          {categories.map(category => (
            <button
              key={category.id}
              className={`px-4 py-2 rounded-lg whitespace-nowrap ${
                selectedCategory === category.id
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
              }`}
              onClick={() => setSelectedCategory(category.id)}
            >
              {category.name}
            </button>
          ))}
        </div>
      </div>
      
      {/* Templates grid */}
      <motion.div 
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ staggerChildren: 0.1 }}
      >
        {filteredTemplates.map(template => (
          <AgentTemplateCard
            key={template.id}
            title={template.title}
            description={template.description}
            icon={template.icon}
            categories={template.categories}
            isPopular={template.isPopular}
            onClick={() => onSelectTemplate(template)}
          />
        ))}
      </motion.div>
    </div>
  );
};

export default AgentTemplatesShowcase;