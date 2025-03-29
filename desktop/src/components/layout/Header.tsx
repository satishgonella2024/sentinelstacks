import React, { useState } from 'react';
import { MoonIcon, SunIcon, BellIcon } from '@heroicons/react/24/outline';

const Header: React.FC = () => {
  const [darkMode, setDarkMode] = useState(
    window.matchMedia('(prefers-color-scheme: dark)').matches
  );

  const toggleDarkMode = () => {
    const newMode = !darkMode;
    setDarkMode(newMode);
    
    if (newMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  };

  return (
    <header className="flex justify-between items-center h-16 px-6 bg-white dark:bg-gray-800 shadow">
      <div>
        <h1 className="text-xl font-bold text-gray-800 dark:text-white">
          SentinelStacks
        </h1>
      </div>
      <div className="flex items-center space-x-4">
        <button
          onClick={toggleDarkMode}
          className="p-2 rounded-full text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
          aria-label="Toggle dark mode"
        >
          {darkMode ? (
            <SunIcon className="h-5 w-5" />
          ) : (
            <MoonIcon className="h-5 w-5" />
          )}
        </button>
        <button
          className="p-2 rounded-full text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 relative"
          aria-label="Notifications"
        >
          <BellIcon className="h-5 w-5" />
          <span className="absolute top-0 right-0 h-2 w-2 rounded-full bg-red-500"></span>
        </button>
        <div className="flex items-center">
          <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white font-bold">
            A
          </div>
          <div className="ml-2 text-gray-800 dark:text-white">Admin</div>
        </div>
      </div>
    </header>
  );
};

export default Header;