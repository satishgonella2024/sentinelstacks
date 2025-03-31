import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { store } from './context/store'
import App from './App.tsx'
import './styles/index.css'

async function bootstrap() {
  // Setup mock server in development
  if (import.meta.env.VITE_USE_MOCK_API === 'true') {
    console.log('🔶 Using mock API in development mode')
    const { worker } = await import('./mocks/browser')
    await worker.start({ 
      onUnhandledRequest: 'bypass' 
    })
  }

  ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
      <Provider store={store}>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </Provider>
    </React.StrictMode>,
  )
}

bootstrap()
