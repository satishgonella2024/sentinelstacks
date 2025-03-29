import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import MainLayout from './components/layout/MainLayout';
import ToastProvider from './components/common/ToastProvider';
import AppRoutes from './routes';
import './styles.css';

const App: React.FC = () => {
  return (
    <Router>
      <ToastProvider>
        <MainLayout>
          <AppRoutes />
        </MainLayout>
      </ToastProvider>
    </Router>
  );
};

export default App;