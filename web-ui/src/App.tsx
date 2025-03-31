import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import ErrorBoundary from './components/common/ErrorBoundary';
import Layout from './components/layout/Layout';

// Simple pages (directly imported, no lazy loading)
import Landing from './pages/Landing';
import NotFound from './pages/NotFound';
import Dashboard from './pages/Dashboard';
import Agents from './pages/Agents';
import Settings from './pages/Settings';
import Builder from './pages/Builder';

// Simple fallback for any pages that fail to load
const ErrorFallback = () => (
  <div className="p-6 max-w-lg mx-auto mt-20 bg-red-900 text-white rounded-lg shadow">
    <h2 className="text-xl font-bold mb-4">Failed to load page</h2>
    <p className="mb-4">
      Sorry, we couldn't load this page. You can try refreshing the browser.
    </p>
    <button 
      className="px-4 py-2 bg-white text-red-900 rounded font-medium"
      onClick={() => window.location.href = '/'}
    >
      Return to Home
    </button>
  </div>
);

const App: React.FC = () => {
  return (
    <ErrorBoundary fallback={<ErrorFallback />}>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/dashboard" element={<Layout><Dashboard /></Layout>} />
        <Route path="/agents" element={<Layout><Agents /></Layout>} />
        <Route path="/builder" element={<Layout><Builder /></Layout>} />
        <Route path="/settings" element={<Layout><Settings /></Layout>} />
        <Route path="/404" element={<NotFound />} />
        <Route path="*" element={<Navigate to="/404" replace />} />
      </Routes>
    </ErrorBoundary>
  );
};

export default App;
