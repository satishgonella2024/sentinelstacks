import React from 'react';

interface PageContainerProps {
  children: React.ReactNode;
  title?: string;
  actions?: React.ReactNode;
}

const PageContainer: React.FC<PageContainerProps> = ({ 
  children, 
  title, 
  actions 
}) => {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      {(title || actions) && (
        <div className="flex flex-col md:flex-row md:items-center justify-between pb-5 border-b border-gray-200 dark:border-gray-700 mb-5">
          {title && <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-4 md:mb-0">{title}</h1>}
          {actions && <div className="flex space-x-3">{actions}</div>}
        </div>
      )}
      {children}
    </div>
  );
};

export default PageContainer;