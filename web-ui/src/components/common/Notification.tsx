import React, { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useDispatch } from 'react-redux'
import { hideNotification } from '@context/slices/uiSlice'

interface NotificationProps {
  type: 'success' | 'error' | 'info' | 'warning'
  message: string
}

const Notification: React.FC<NotificationProps> = ({ type, message }) => {
  const [visible, setVisible] = useState(true)
  const dispatch = useDispatch()
  
  // Automatically hide notification after 5 seconds
  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(false)
      setTimeout(() => dispatch(hideNotification()), 300) // Allow animation to complete
    }, 5000)
    
    return () => clearTimeout(timer)
  }, [dispatch])
  
  // Get icon and color based on notification type
  const getNotificationStyles = () => {
    switch (type) {
      case 'success':
        return {
          icon: '✓',
          bgColor: 'bg-green-500',
          borderColor: 'border-green-400',
          iconBg: 'bg-green-600',
        }
      case 'error':
        return {
          icon: '✕',
          bgColor: 'bg-red-500',
          borderColor: 'border-red-400',
          iconBg: 'bg-red-600',
        }
      case 'warning':
        return {
          icon: '!',
          bgColor: 'bg-yellow-500',
          borderColor: 'border-yellow-400',
          iconBg: 'bg-yellow-600',
        }
      case 'info':
      default:
        return {
          icon: 'i',
          bgColor: 'bg-primary-500',
          borderColor: 'border-primary-400',
          iconBg: 'bg-primary-600',
        }
    }
  }
  
  const styles = getNotificationStyles()
  
  return (
    <AnimatePresence>
      {visible && (
        <motion.div
          className={`fixed top-4 right-4 max-w-sm ${styles.bgColor} bg-opacity-90 border-l-4 ${styles.borderColor} rounded shadow-glow-sm`}
          initial={{ x: 100, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          exit={{ x: 100, opacity: 0 }}
          transition={{ type: 'spring', stiffness: 500, damping: 30 }}
        >
          <div className="flex items-center p-4">
            <div className={`flex-shrink-0 w-8 h-8 ${styles.iconBg} rounded-full flex items-center justify-center mr-3`}>
              <span className="text-white font-bold">{styles.icon}</span>
            </div>
            <div className="flex-grow text-white text-sm">{message}</div>
            <button 
              onClick={() => {
                setVisible(false)
                setTimeout(() => dispatch(hideNotification()), 300)
              }}
              className="ml-4 text-white opacity-70 hover:opacity-100"
            >
              ×
            </button>
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}

export default Notification 