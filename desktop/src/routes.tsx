import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Agents from './pages/Agents';
import AgentCreate from './pages/AgentCreate';
import AgentDetail from './pages/AgentDetail';
import AgentEdit from './pages/AgentEdit';
import Monitoring from './pages/Monitoring';
import Settings from './pages/Settings';

const AppRoutes: React.FC = () => {
  return (
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
  );
};

export default AppRoutes; 