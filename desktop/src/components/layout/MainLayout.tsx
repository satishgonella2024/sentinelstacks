import React from 'react';
import Sidebar from '../common/Sidebar';
import Header from '../common/Header';

interface MainLayoutProps {
  children: React.ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  return (
    <div className="flex h-screen bg-gray-100 dark:bg-gray-900">
      <Sidebar />
      <div className="flex-1 flex flex-col md:pl-64">
        <Header />
        <main className="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100 dark:bg-gray-900">
          <div className="py-6">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};

export default MainLayout;