import { useState } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import MainLayout from "./components/layout/MainLayout";
import Dashboard from "./pages/Dashboard";
import Agents from "./pages/Agents";
import Monitoring from "./pages/Monitoring";
import Settings from "./pages/Settings";

function App() {
  return (
    <div className="h-screen bg-gray-50 dark:bg-gray-900">
      <Router>
        <Routes>
          <Route path="/" element={<MainLayout />}>
            <Route index element={<Dashboard />} />
            <Route path="agents" element={<Agents />} />
            <Route path="monitoring" element={<Monitoring />} />
            <Route path="settings" element={<Settings />} />
          </Route>
        </Routes>
      </Router>
    </div>
  );
}

export default App;
