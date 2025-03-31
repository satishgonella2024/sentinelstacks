import { configureStore, Middleware } from '@reduxjs/toolkit'
import { setupListeners } from '@reduxjs/toolkit/query'
import { apiSlice } from '../services/api'

// Reducers
import authReducer from './slices/authSlice'
import uiReducer from './slices/uiSlice'
import agentsReducer from './slices/agentsSlice'

// API Logger middleware for debugging
const apiLogger: Middleware = () => (next) => (action: any) => {
  // Log actions related to API requests
  if (action.type && typeof action.type === 'string' && action.type.endsWith('/executeQuery')) {
    console.log('API Request:', action);
  }
  
  // Log actions related to API responses
  if (action.type && typeof action.type === 'string' && 
      (action.type.endsWith('/executeQuery/fulfilled') || 
       action.type.endsWith('/executeQuery/rejected'))) {
    console.log('API Response:', action);
  }
  
  return next(action);
};

export const store = configureStore({
  reducer: {
    auth: authReducer,
    ui: uiReducer,
    agents: agentsReducer,
    [apiSlice.reducerPath]: apiSlice.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(apiSlice.middleware, apiLogger),
  devTools: true,
})

// Enable refetchOnFocus and refetchOnReconnect
setupListeners(store.dispatch)

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch 