import React from 'react';
import { NavLink } from 'react-router-dom';
import {
  HomeIcon,
  CpuChipIcon,
  ServerIcon,
  CogIcon
} from '@heroicons/react/24/outline';

const Sidebar: React.FC = () => {
  const navigation = [
    { name: 'Dashboard', href: '/', icon: HomeIcon },
    { name: 'Agents', href: '/agents', icon: CpuChipIcon },
    { name: 'Registry', href: '/registry', icon: ServerIcon },
    { name: 'Settings', href: '/settings', icon: CogIcon },
  ];

  return (
    <div className="flex flex-col w-64 bg-gray-800 text-white">
      <div className="flex items-center justify-center h-16 border-b border-gray-700">
        <span className="text-xl font-bold">SentinelStacks</span>
      </div>
      <nav className="flex-1 px-2 py-4 space-y-1">
        {navigation.map((item) => (
          <NavLink
            key={item.name}
            to={item.href}
            className={({ isActive }) =>
              `flex items-center px-4 py-2 rounded-md ${
                isActive
                  ? 'bg-gray-900 text-white'
                  : 'text-gray-300 hover:bg-gray-700'
              }`
            }
          >
            <item.icon className="w-5 h-5 mr-3" />
            {item.name}
          </NavLink>
        ))}
      </nav>
      <div className="px-4 py-2 border-t border-gray-700">
        <div className="flex items-center">
          <div className="w-8 h-8 bg-gray-600 rounded-full flex items-center justify-center">
            <CpuChipIcon className="w-4 h-4 text-white" />
          </div>
          <div className="ml-2">
            <div className="text-sm font-medium">SentinelStacks</div>
            <div className="text-xs text-gray-400">v0.1.0</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;