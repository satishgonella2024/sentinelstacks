import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface UIState {
  theme: 'light' | 'dark'
  sidebarOpen: boolean
  notification: {
    show: boolean
    message: string
    type: 'success' | 'error' | 'info' | 'warning'
  } | null
  thinkBubbles: {
    show: boolean
    position: 'dashboard' | 'chat' | 'builder' | 'explorer'
  }
}

const initialState: UIState = {
  theme: (localStorage.getItem('theme') as 'light' | 'dark') || 'dark',
  sidebarOpen: true,
  notification: null,
  thinkBubbles: {
    show: true,
    position: 'dashboard'
  }
}

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    toggleTheme: (state) => {
      const newTheme = state.theme === 'light' ? 'dark' : 'light'
      state.theme = newTheme
      localStorage.setItem('theme', newTheme)
    },
    toggleSidebar: (state) => {
      state.sidebarOpen = !state.sidebarOpen
    },
    showNotification: (state, action: PayloadAction<Omit<NonNullable<UIState['notification']>, 'show'>>) => {
      state.notification = {
        show: true,
        ...action.payload
      }
    },
    hideNotification: (state) => {
      if (state.notification) {
        state.notification.show = false
      }
    },
    toggleThinkBubbles: (state) => {
      state.thinkBubbles.show = !state.thinkBubbles.show
    },
    setThinkBubblePosition: (state, action: PayloadAction<UIState['thinkBubbles']['position']>) => {
      state.thinkBubbles.position = action.payload
    }
  }
})

export const { 
  toggleTheme, 
  toggleSidebar, 
  showNotification, 
  hideNotification,
  toggleThinkBubbles,
  setThinkBubblePosition
} = uiSlice.actions

export default uiSlice.reducer 