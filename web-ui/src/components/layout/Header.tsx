import React from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { toggleTheme, toggleSidebar } from '@context/slices/uiSlice'
import { RootState } from '@context/store'
import { useLocation } from 'react-router-dom'

const Header: React.FC = () => {
  const dispatch = useDispatch()
  const { theme } = useSelector((state: RootState) => state.ui)
  const { user } = useSelector((state: RootState) => state.auth)
  const location = useLocation()
  
  // Get page title based on current route
  const getPageTitle = () => {
    const path = location.pathname
    
    if (path === '/') return 'Home'
    if (path === '/dashboard') return 'Dashboard'
    if (path.includes('/agents')) return 'Agent Management'
    if (path.includes('/builder')) return 'Agent Builder'
    if (path.includes('/analytics')) return 'Analytics'
    if (path.includes('/settings')) return 'Settings'
    
    return 'SentinelStacks'
  }
  
  return (
    <header className="h-16 flex items-center justify-between border-b border-gray-800 px-4">
      <div className="flex items-center">
        <button
          onClick={() => dispatch(toggleSidebar())}
          className="mr-4 p-2 rounded-full hover:bg-gray-800 transition-colors"
          aria-label="Toggle sidebar"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
        </button>
        
        <h1 className="text-xl font-display text-white">{getPageTitle()}</h1>
      </div>
      
      <div className="flex items-center space-x-3">
        <button
          onClick={() => dispatch(toggleTheme())}
          className="p-2 rounded-full hover:bg-gray-800 transition-colors"
          aria-label={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
        >
          {theme === 'dark' ? (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
              />
            </svg>
          ) : (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
              />
            </svg>
          )}
        </button>
        
        <button className="p-2 rounded-full hover:bg-gray-800 transition-colors" aria-label="Notifications">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
            />
          </svg>
        </button>
        
        <div className="flex items-center ml-4">
          <div className="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white font-medium">
            {user ? user.username.charAt(0).toUpperCase() : 'G'}
          </div>
          <span className="ml-2 text-sm font-medium text-gray-300">
            {user ? user.username : 'Guest'}
          </span>
        </div>
      </div>
    </header>
  )
}

export default Header 