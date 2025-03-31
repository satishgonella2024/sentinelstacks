import React from 'react';
import { motion } from 'framer-motion';
import { FiCpu, FiArrowRight } from 'react-icons/fi';

interface AgentTemplateCardProps {
  title: string;
  description: string;
  icon: React.ReactNode;
  categories: string[];
  onClick: () => void;
  isPopular?: boolean;
}

const AgentTemplateCard: React.FC<AgentTemplateCardProps> = ({
  title,
  description,
  icon,
  categories,
  onClick,
  isPopular = false
}) => {
  return (
    <motion.div 
      className="glass rounded-xl border border-gray-700/50 h-full flex flex-col cursor-pointer card-hover overflow-hidden"
      whileHover={{ scale: 1.02 }}
      onClick={onClick}
    >
      {isPopular && (
        <div className="bg-accent-purple text-white px-3 py-1 text-xs uppercase font-semibold tracking-wider absolute right-0 top-4 -mr-10 rotate-45">
          Popular
        </div>
      )}
      
      <div className="p-6 flex flex-col h-full">
        <div className="flex items-center mb-4">
          <div className="w-10 h-10 rounded-full bg-primary-600/20 flex items-center justify-center mr-3 text-primary-400">
            {icon}
          </div>
          <h3 className="text-lg font-display text-white">{title}</h3>
        </div>
        
        <p className="text-gray-300 mb-4 flex-grow">{description}</p>
        
        <div className="mt-auto">
          <div className="flex flex-wrap gap-2 mb-4">
            {categories.map((category, index) => (
              <span key={index} className="text-xs px-2 py-1 rounded-full bg-gray-800 text-gray-400">
                {category}
              </span>
            ))}
          </div>
          
          <button className="w-full px-4 py-2 bg-primary-600/30 hover:bg-primary-600/50 text-primary-300 rounded-lg transition-colors flex items-center justify-center group">
            <span>Use This Template</span>
            <FiArrowRight className="ml-2 group-hover:translate-x-1 transition-transform" />
          </button>
        </div>
      </div>
    </motion.div>
  );
};

export default AgentTemplateCard;