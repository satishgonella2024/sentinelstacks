import React from 'react';
import { Tool } from '../../types/Agent';

interface AgentToolBadgeProps {
  tool: Tool;
  className?: string;
}

const AgentToolBadge: React.FC<AgentToolBadgeProps> = ({ tool, className = '' }) => {
  return (
    <span 
      className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 ${className}`}
      title={`${tool.description} (v${tool.version})`}
    >
      {tool.name}
    </span>
  );
};

export default AgentToolBadge;