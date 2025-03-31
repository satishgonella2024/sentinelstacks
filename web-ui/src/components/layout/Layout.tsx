import React from 'react'
import { Outlet } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { RootState } from '@context/store'

// Components
import Sidebar from './Sidebar'
import Header from './Header'
import ThinkBubble from '../common/ThinkBubble'
import Notification from '../common/Notification'

const Layout: React.FC = () => {
  const { sidebarOpen, notification, thinkBubbles } = useSelector((state: RootState) => state.ui)
  
  return (
    <div className="flex h-screen bg-background-900 text-white">
      <Sidebar isOpen={sidebarOpen} />
      
      <div className="flex flex-col flex-1 h-full overflow-hidden">
        <Header />
        
        <main className="flex-1 overflow-auto p-4">
          <Outlet />
        </main>
        
        {notification && notification.show && (
          <Notification 
            message={notification.message} 
            type={notification.type} 
          />
        )}
        
        {thinkBubbles.show && (
          <ThinkBubble position={thinkBubbles.position} />
        )}
      </div>
    </div>
  )
}

export default Layout 