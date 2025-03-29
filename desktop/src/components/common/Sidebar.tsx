import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
  RectangleGroupIcon,
  CpuChipIcon,
  Cog6ToothIcon,
  ChartBarIcon,
  ServerIcon,
  DocumentTextIcon
} from '@heroicons/react/24/outline';

interface NavItem {
  name: string;
  href: string;
  icon: (props: React.ComponentProps<'svg'>) => JSX.Element;
}

const navigation: NavItem[] = [
  { name: 'Dashboard', href: '/', icon: RectangleGroupIcon },
  { name: 'Agents', href: '/agents', icon: CpuChipIcon },
  { name: 'Monitoring', href: '/monitoring', icon: ChartBarIcon },
  { name: 'Registry', href: '/registry', icon: ServerIcon },
  { name: 'Documentation', href: '/docs', icon: DocumentTextIcon },
  { name: 'Settings', href: '/settings', icon: Cog6ToothIcon },
];

const Sidebar: React.FC = () => {
  const location = useLocation();

  return (
    <div className="hidden md:flex md:w-64 md:flex-col md:fixed md:inset-y-0">
      <div className="flex-1 flex flex-col min-h-0 bg-gray-800">
        <div className="flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
          <div className="flex items-center flex-shrink-0 px-4">
            <img
              className="h-8 w-auto"
              src="/sentinel-logo.svg"
              alt="SentinelStacks"
            />
            <span className="ml-2 text-white text-lg font-semibold">SentinelStacks</span>
          </div>
          <nav className="mt-5 flex-1 px-2 space-y-1">
            {navigation.map((item) => {
              const current = location.pathname === item.href || 
                              (item.href !== '/' && location.pathname.startsWith(item.href));
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`
                    group flex items-center px-2 py-2 text-sm font-medium rounded-md
                    ${current
                      ? 'bg-gray-900 text-white'
                      : 'text-gray-300 hover:bg-gray-700 hover:text-white'}
                  `}
                >
                  <item.icon
                    className={`
                      mr-3 flex-shrink-0 h-6 w-6
                      ${current
                        ? 'text-white'
                        : 'text-gray-400 group-hover:text-gray-300'}
                    `}
                    aria-hidden="true"
                  />
                  {item.name}
                </Link>
              );
            })}
          </nav>
        </div>
        <div className="flex-shrink-0 flex border-t border-gray-700 p-4">
          <div className="flex-shrink-0 w-full group block">
            <div className="flex items-center">
              <div>
                <p className="text-sm font-medium text-white">v0.1.0-alpha</p>
                <p className="text-xs font-medium text-gray-400 group-hover:text-gray-300">
                  Development Build
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;