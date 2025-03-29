import React from 'react';
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline';
import { AgentStatus } from '../../types/Agent';

interface StatusBadgeProps {
  status: AgentStatus;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ status }) => {
  switch (status) {
    case 'active':
      return (
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
          <span className="h-2 w-2 rounded-full bg-green-400 mr-1.5" />
          Active
        </span>
      );
    case 'inactive':
      return (
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200">
          <span className="h-2 w-2 rounded-full bg-gray-400 mr-1.5" />
          Inactive
        </span>
      );
    case 'error':
      return (
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200">
          <ExclamationTriangleIcon className="h-3 w-3 mr-1" />
          Error
        </span>
      );
    default:
      return null;
  }
};

export default StatusBadge;