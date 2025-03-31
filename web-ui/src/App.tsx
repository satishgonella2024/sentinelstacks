import { Routes, Route } from 'react-router-dom'
import { lazy, Suspense } from 'react'

// Lazy load pages
const Landing = lazy(() => import('@pages/Landing'))
const Dashboard = lazy(() => import('@pages/Dashboard'))
const NotFound = lazy(() => import('@pages/NotFound'))

// Layout component
import Layout from '@components/layout/Layout'

// Loading component
import LoadingScreen from '@components/common/LoadingScreen'

function App() {
  return (
    <Suspense fallback={<LoadingScreen />}>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Landing />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </Suspense>
  )
}

export default App 