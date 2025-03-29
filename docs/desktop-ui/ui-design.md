# SentinelStacks Desktop UI Design

This document showcases the UI design and mockups for the SentinelStacks desktop application.

## Main Dashboard

The main dashboard provides an overview of available agents and quick access to common actions.

### Agents Overview

![Agents Overview](../images/desktop-ui/agents-overview.png)

The agents view displays all available agents as cards, showing:
- Agent name and description
- Capabilities/tags
- Quick action buttons (Run, Edit)
- A "Create New Agent" card for adding new agents

### Agent Details

![Agent Details](../images/desktop-ui/agent-details.png)

The agent detail view provides comprehensive information and interaction with a specific agent:

**Chat Tab:**
- Interactive chat interface
- Message history
- Input area for sending commands

**Logs Tab:**
- Real-time execution logs
- Log filtering and search
- Timestamp and log level indicators

**Memory Tab:**
- View agent's memory entries
- Search through memory
- Memory usage statistics

**Configuration Tab:**
- Edit agent properties
- Configure capabilities and tools
- Advanced YAML configuration

## Registry

![Registry](../images/desktop-ui/registry.png)

The registry view allows browsing and managing agents from the central repository:
- Search and filtering
- Sort by popularity, date, or name
- Install agents with a single click
- View agent details before installation

## Settings

![Settings](../images/desktop-ui/settings.png)

The settings view provides application configuration:
- Theme selection (light/dark)
- Model provider settings
- API keys management
- Application preferences

## Design System

### Colors

Primary palette:
- Primary: #00AFC8 (Cyan)
- Secondary: #6366F1 (Indigo)
- Accent: #F59E0B (Amber)
- Success: #10B981 (Green)
- Warning: #F59E0B (Amber)
- Error: #EF4444 (Red)

Neutral palette (Light mode):
- Background: #F9FAFB
- Surface: #FFFFFF
- Border: #E5E7EB
- Text: #111827

Neutral palette (Dark mode):
- Background: #111827
- Surface: #1F2937
- Border: #374151
- Text: #F9FAFB

### Typography

- Headings: Inter, sans-serif
- Body: Inter, sans-serif
- Monospace: JetBrains Mono, monospace

Font sizes:
- xs: 0.75rem
- sm: 0.875rem
- base: 1rem
- lg: 1.125rem
- xl: 1.25rem
- 2xl: 1.5rem
- 3xl: 1.875rem
- 4xl: 2.25rem

### Components

#### Buttons

Primary button:
- Background: #00AFC8
- Text: white
- Hover: #008BA0
- Active: #006878

Secondary button:
- Background: #E5E7EB (light) / #374151 (dark)
- Text: #111827 (light) / #F9FAFB (dark)
- Hover: #D1D5DB (light) / #4B5563 (dark)
- Active: #9CA3AF (light) / #6B7280 (dark)

#### Cards

- Background: white (light) / #1F2937 (dark)
- Border: #E5E7EB (light) / #374151 (dark)
- Border radius: 0.5rem
- Shadow: 0 1px 3px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06)

#### Inputs

- Background: white (light) / #374151 (dark)
- Border: #D1D5DB (light) / #4B5563 (dark)
- Focus border: #00AFC8
- Text: #111827 (light) / #F9FAFB (dark)
- Placeholder: #9CA3AF

#### Tags

- Background: #E5E7EB (light) / #374151 (dark)
- Text: #4B5563 (light) / #D1D5DB (dark)
- Border radius: 0.25rem
- Padding: 0.25rem 0.5rem

### Responsive Breakpoints

- sm: 640px
- md: 768px
- lg: 1024px
- xl: 1280px
- 2xl: 1536px

## Interaction Patterns

### Agent Creation

1. Click "Create New Agent" card
2. Choose creation method (wizard or natural language)
3. Enter agent details (name, description)
4. Select model provider and model
5. Configure capabilities and tools
6. Review and create

### Agent Execution

1. Click "Run" on an agent card
2. Agent detail view opens with chat interface
3. Enter commands in the chat input
4. View responses in real-time
5. Access logs and memory during execution

### Registry Navigation

1. Click "Registry" in sidebar
2. Browse available agents
3. Use search and filters to find specific agents
4. Click "Install" to add an agent to your local collection
5. Installed agents appear in the "Agents" view

## Accessibility Considerations

- All colors meet WCAG 2.1 AA contrast requirements
- All interactive elements are keyboard accessible
- Focus states are clearly visible
- Screen reader support for all UI elements
- Responsive design for different screen sizes
- Support for system color schemes and reduced motion preferences
