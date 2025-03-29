# SentinelStacks Desktop UI

This is the desktop UI implementation for SentinelStacks.

## Implementation Status

The following components have been successfully implemented:

1. **Core Structure**
   * Layout System with MainLayout and PageContainer components
   * Navigation with Header and Sidebar components
   * Routing for all major pages

2. **Agent Management**
   * Agent Listing with grid and list views, filtering
   * Agent Creation with model, tool, and capability selection
   * Agent Editing with form reuse for editing
   * Agent Detail View with a tabbed interface

3. **Reusable Components**
   * Status Badge for displaying agent statuses
   * Loading Spinner for consistent loading indication
   * Model Configuration form section
   * Tool Selection component
   * Capability Selection component

4. **State Management**
   * Custom useAgent and useAgentList hooks
   * Service layer for API interaction
   * TypeScript type definitions

5. **User Experience**
   * Dark Mode support
   * Responsive Design
   * Loading States
   * Consistent Styling with Tailwind CSS

6. **Recent Enhancements**
   * Memory Visualization component
   * Conversation Interface
   * Toast Notification system
   * Error Handling

## Getting Started

### Setup with the Fix Script

To ensure all dependencies are properly installed, run the fix script first:

```bash
# Make the script executable
chmod +x ./fix-dependencies.sh

# Run the script
./fix-dependencies.sh
```

### Manual Setup

If you prefer to set up manually:

1. First, install the dependencies:
   ```bash
   npm install
   ```

2. If you encounter PostCSS errors, make sure the PostCSS configuration is correct:
   ```bash
   # Check that postcss.config.cjs contains:
   module.exports = {
     plugins: {
       tailwindcss: {},
       autoprefixer: {},
     },
   }
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

4. Open [http://localhost:5173](http://localhost:5173) in your browser to see the application.

## Troubleshooting

If you encounter build errors:

1. PostCSS Configuration: Make sure both project root and desktop directory have proper postcss.config.cjs files.
2. Tailwind Configuration: Ensure tailwind.config.js is properly configured.
3. Missing Dependencies: Install any missing dependencies with `npm install [package-name]`.

## Next Steps

The following features are still in progress:

1. **Registry Integration**: Implement registry browsing and publishing
2. **Settings Management**: Complete the settings interface
3. **Accessibility**: Test and improve keyboard navigation and screen reader support
4. **Desktop Packaging**: Package the application for desktop using Tauri

## Tech Stack

- **Frontend**: React, TypeScript, Tailwind CSS
- **State Management**: React Hooks
- **Routing**: React Router
- **UI Components**: Custom components with Tailwind
- **API Client**: Axios
- **Toast Notifications**: react-hot-toast
- **Charts**: Recharts
- **Desktop Framework**: Tauri
