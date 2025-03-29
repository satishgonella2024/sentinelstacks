import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import MainLayout from './components/layout/MainLayout';
// ToastProvider is now handled directly with a Toaster component
import AppRoutes from './routes';
import './styles.css';

const App: React.FC = () => {
  return (
    <Router>
      <MainLayout>
        <AppRoutes />
      </MainLayout>
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 4000,
          style: {
            background: 'var(--toast-bg, #ffffff)',
            color: 'var(--toast-color, #1f2937)',
          },
          success: {
            className: '!bg-green-50 !text-green-800 dark:!bg-green-900 dark:!text-green-100',
            iconTheme: {
              primary: '#22c55e',
              secondary: '#ffffff',
            },
          },
          error: {
            className: '!bg-red-50 !text-red-800 dark:!bg-red-900 dark:!text-red-100',
            iconTheme: {
              primary: '#ef4444',
              secondary: '#ffffff',
            },
          },
        }}
      />
    </Router>
  );
};

export default App;