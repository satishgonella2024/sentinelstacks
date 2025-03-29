import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Agents from './pages/Agents';
import AgentDetail from './pages/AgentDetail';
import AgentCreate from './pages/AgentCreate';
import AgentEdit from './pages/AgentEdit';
import Monitoring from './pages/Monitoring';
import Settings from './pages/Settings';
import './App.css';

function App() {
  // Set up theme
  useEffect(() => {
    // Check if the user has a theme preference in localStorage
    const storedTheme = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    // Apply dark theme if stored preference is dark or no preference but system is dark
    if (storedTheme === 'dark' || (!storedTheme && prefersDark)) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/agents" element={<Agents />} />
        <Route path="/agents/create" element={<AgentCreate />} />
        <Route path="/agents/:id" element={<AgentDetail />} />
        <Route path="/agents/:id/edit" element={<AgentEdit />} />
        <Route path="/monitoring" element={<Monitoring />} />
        <Route path="/settings" element={<Settings />} />
        {/* Redirect any unknown routes to the dashboard */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;