import React, { useState } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';
import { Bars3Icon, XMarkIcon, ChartBarIcon, ComputerDesktopIcon, ChartPieIcon, Cog6ToothIcon } from '@heroicons/react/24/outline';
import Logo from '../Logo';

const MainLayout: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const location = useLocation();

  const navigation = [
    { name: 'Dashboard', href: '/', icon: ChartPieIcon },
    { name: 'Agents', href: '/agents', icon: ComputerDesktopIcon },
    { name: 'Monitoring', href: '/monitoring', icon: ChartBarIcon },
    { name: 'Settings', href: '/settings', icon: Cog6ToothIcon },
  ];

  const isActive = (path: string) => {
    return location.pathname === path;
  };

  return (
    <div className="flex h-screen bg-gray-50 dark:bg-gray-900">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-20 bg-gray-900/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div
        className={`fixed inset-y-0 left-0 z-30 w-64 transform bg-white dark:bg-gray-800 shadow-lg transition-transform duration-200 lg:translate-x-0 lg:static lg:z-auto ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        }`}
      >
        {/* Sidebar header */}
        <div className="flex items-center justify-between h-16 px-4 border-b dark:border-gray-700">
          <Link to="/" className="flex items-center space-x-3">
            <Logo className="h-8 w-8" />
            <span className="text-lg font-semibold text-primary-600 dark:text-primary-500">
              SentinelStacks
            </span>
          </Link>
          <button
            className="p-2 text-gray-500 rounded-md lg:hidden hover:text-gray-900 hover:bg-gray-100 dark:hover:text-white dark:hover:bg-gray-700"
            onClick={() => setSidebarOpen(false)}
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        {/* Navigation */}
        <nav className="p-4 space-y-1">
          {navigation.map((item) => {
            const active = isActive(item.href);
            return (
              <Link
                key={item.name}
                to={item.href}
                className={`flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors duration-150 ${
                  active
                    ? 'bg-primary-50 text-primary-600 dark:bg-primary-900/50 dark:text-primary-500'
                    : 'text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700'
                }`}
              >
                <item.icon className={`w-5 h-5 mr-3 ${
                  active ? 'text-primary-600 dark:text-primary-500' : 'text-gray-400 dark:text-gray-500'
                }`} />
                {item.name}
              </Link>
            );
          })}
        </nav>
      </div>

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
        {/* Header */}
        <header className="sticky top-0 z-10 flex-shrink-0 h-16 bg-white dark:bg-gray-800 shadow-sm">
          <div className="flex items-center justify-between h-full px-4">
            <button
              className="p-2 text-gray-500 rounded-md lg:hidden hover:text-gray-900 hover:bg-gray-100 dark:hover:text-white dark:hover:bg-gray-700"
              onClick={() => setSidebarOpen(true)}
            >
              <Bars3Icon className="w-6 h-6" />
            </button>
            <div className="flex items-center space-x-4">
              {/* Add header content here (e.g., user profile, notifications) */}
            </div>
          </div>
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default MainLayout; 