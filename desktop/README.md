# SentinelStacks Desktop UI

This is the desktop application for SentinelStacks, built with Tauri, React, and TypeScript.

## Architecture

The desktop UI follows a clean, component-based architecture:

```
src/
├── components/              # Reusable UI components
│   ├── agents/              # Agent-related components
│   │   ├── forms/           # Form components for agents
│   │   ├── AgentCard.tsx    # Card display for agents
│   │   ├── AgentList.tsx    # List display for agents
│   │   └── AgentListItem.tsx # List item display for agents
│   ├── common/              # Common UI components
│   │   ├── Header.tsx       # Application header
│   │   ├── Sidebar.tsx      # Application sidebar
│   │   ├── StatusBadge.tsx  # Status indicator badge
│   │   └── LoadingSpinner.tsx # Loading indicator
│   ├── layout/              # Layout components
│   │   ├── MainLayout.tsx   # Main application layout
│   │   └── PageContainer.tsx # Container for page content
│   └── memory/              # Memory visualization components
├── pages/                   # Application pages
│   ├── Dashboard.tsx        # Dashboard page
│   ├── Agents.tsx           # Agents listing page
│   ├── AgentDetail.tsx      # Agent detail page
│   ├── AgentCreate.tsx      # Agent creation page
│   ├── AgentEdit.tsx        # Agent editing page
│   ├── Monitoring.tsx       # Monitoring page
│   └── Settings.tsx         # Settings page
├── services/                # API services
│   ├── agentService.ts      # Service for agent operations
│   └── memoryService.ts     # Service for memory operations
├── hooks/                   # Custom React hooks
│   ├── useAgent.ts          # Hook for agent operations
│   └── useAgentList.ts      # Hook for agent listing
└── types/                   # TypeScript type definitions
    ├── Agent.ts             # Agent-related types
    └── Tool.ts              # Tool-related types
```

## Development

### Prerequisites

- Node.js 18+ (LTS recommended)
- Rust (for Tauri)
- Go 1.21+ (for backend)

### Setup

1. Install dependencies:
```bash
npm install
```

2. Run in development mode:
```bash
npm run tauri dev
```

### Building

Build the desktop application for production:
```bash
npm run tauri build
```

This will create platform-specific installers in the `src-tauri/target/release/bundle` directory.

## Features

### Completed (✅)

- Agent management interface
  - List view and grid view
  - Status indicators
  - Agent filtering and search
  - Action buttons (start/stop/edit)
- Agent detail view
  - Information display
  - Status control
  - Tabbed interface for different sections
- Agent creation and editing
  - Form validation
  - Model selection
  - Tool selection
  - Capability configuration
  - Memory options

### In Progress (🚧)

- Memory visualization
- Conversation interface
- Registry integration
- Settings management

## Component Guide

### Agent Components

- **AgentCard**: Displays an agent as a card in grid view
- **AgentListItem**: Displays an agent as a row in list view
- **AgentList**: Container for displaying multiple agents with filtering
- **AgentForm**: Form for creating or editing agents

### Form Components

- **ModelConfig**: Form section for configuring model provider and parameters
- **ToolSelection**: Form section for selecting agent tools
- **CapabilitySelection**: Form section for selecting agent capabilities

### Layout Components

- **MainLayout**: Main application layout with sidebar and header
- **PageContainer**: Container for page content with consistent styling

## Styling

The application uses Tailwind CSS for styling. The theme is configured to support both light and dark modes.

### Color Scheme

- Primary: Blue (#3b82f6)
- Success: Green (#10b981)
- Warning: Orange (#f59e0b)
- Danger: Red (#ef4444)
- Neutral: Gray (#6b7280)

### Dark Mode

Dark mode is implemented using Tailwind's dark mode feature. The application respects the user's system preference and allows manual toggling.