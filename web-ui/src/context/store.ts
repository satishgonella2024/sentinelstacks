import { configureStore } from '@reduxjs/toolkit'
import { setupListeners } from '@reduxjs/toolkit/query'
import { apiSlice } from '@services/api'

// Reducers
import authReducer from './slices/authSlice'
import uiReducer from './slices/uiSlice'
import agentsReducer from './slices/agentsSlice'

export const store = configureStore({
  reducer: {
    auth: authReducer,
    ui: uiReducer,
    agents: agentsReducer,
    [apiSlice.reducerPath]: apiSlice.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(apiSlice.middleware),
  devTools: process.env.NODE_ENV !== 'production',
})

// Enable refetchOnFocus and refetchOnReconnect
setupListeners(store.dispatch)

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch 