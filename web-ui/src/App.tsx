import React, { lazy, Suspense } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import NotFound from '@pages/NotFound'
import Landing from '@pages/Landing'
import Layout from '@components/layout/Layout'

// Lazy load pages to reduce initial bundle size
const Dashboard = lazy(() => import('@pages/Dashboard'))
const Agents = lazy(() => import('@pages/Agents'))
const Builder = lazy(() => import('@pages/Builder'))
const Images = lazy(() => import('@pages/Images'))
const Analytics = lazy(() => import('@pages/Analytics'))
const Settings = lazy(() => import('@pages/Settings'))

const App: React.FC = () => {
  return (
    <Suspense fallback={<div className="loading-fallback">Loading...</div>}>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/dashboard" element={<Layout><Dashboard /></Layout>} />
        <Route path="/agents" element={<Layout><Agents /></Layout>} />
        <Route path="/builder" element={<Layout><Builder /></Layout>} />
        <Route path="/images" element={<Layout><Images /></Layout>} />
        <Route path="/analytics" element={<Layout><Analytics /></Layout>} />
        <Route path="/settings" element={<Layout><Settings /></Layout>} />
        <Route path="/404" element={<NotFound />} />
        <Route path="*" element={<Navigate to="/404" replace />} />
      </Routes>
    </Suspense>
  )
}

export default App 