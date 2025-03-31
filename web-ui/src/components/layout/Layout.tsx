import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const location = useLocation();
  
  const isActive = (path: string) => location.pathname === path;
  
  const navItems = [
    { path: '/dashboard', label: 'Dashboard' },
    { path: '/agents', label: 'Agents' },
    { path: '/registry', label: 'Registry' },
    { path: '/documentation', label: 'Documentation' },
    { path: '/settings', label: 'Settings' },
  ];
  
  return (
    <div className="flex h-screen bg-gray-900 text-white">
      {/* Sidebar for desktop */}
      <div className="hidden md:flex md:w-64 bg-gray-800 flex-col">
        <div className="p-4 border-b border-gray-700">
          <Link to="/" className="text-xl font-bold">SentinelStacks</Link>
        </div>
        
        <nav className="flex-1 p-4">
          <ul className="space-y-2">
            {navItems.map(item => (
              <li key={item.path}>
                <Link 
                  to={item.path}
                  className={`block px-4 py-2 rounded transition ${
                    isActive(item.path) 
                      ? 'bg-primary-600 text-white' 
                      : 'text-gray-300 hover:bg-gray-700'
                  }`}
                >
                  {item.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
        
        <div className="p-4 border-t border-gray-700">
          <div className="flex items-center">
            <div className="w-8 h-8 rounded-full bg-gray-600 mr-3"></div>
            <div>
              <div className="font-medium">User</div>
              <div className="text-xs text-gray-400">Admin</div>
            </div>
          </div>
        </div>
      </div>
      
      {/* Mobile header */}
      <div className="flex flex-col flex-1">
        <header className="md:hidden bg-gray-800 p-4 flex items-center justify-between">
          <Link to="/" className="text-xl font-bold">SentinelStacks</Link>
          
          <button 
            className="text-gray-400 hover:text-white"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? 'Close' : 'Menu'}
          </button>
        </header>
        
        {/* Mobile menu */}
        {isMobileMenuOpen && (
          <div className="md:hidden bg-gray-800 border-b border-gray-700">
            <nav className="p-4">
              <ul className="space-y-2">
                {navItems.map(item => (
                  <li key={item.path}>
                    <Link 
                      to={item.path}
                      className={`block px-4 py-2 rounded ${
                        isActive(item.path) 
                          ? 'bg-primary-600 text-white' 
                          : 'text-gray-300 hover:bg-gray-700'
                      }`}
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      {item.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </nav>
          </div>
        )}
        
        {/* Main content */}
        <main className="flex-1 overflow-auto">
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;
