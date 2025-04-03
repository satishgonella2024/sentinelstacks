# SentinelStacks Web UI

A modern, enterprise-grade web interface for the SentinelStacks AI Agent Management Platform.

## Vision

SentinelStacks Web UI provides a sophisticated yet intuitive interface for enterprise users to create, deploy, manage, and monitor AI agents at scale. The interface combines professional design with futuristic elements to create an engaging user experience without sacrificing enterprise functionality.

## Comparison with Legacy Web Interface

| Feature | Modern Web UI | Legacy Web Interface |
|---------|--------------|---------------------|
| **Technology** | React, TypeScript, Tailwind CSS | Vanilla HTML, CSS, JavaScript |
| **Architecture** | Component-based with state management | Monolithic single-page application |
| **API Integration** | RTK Query with mock data support | Direct fetch calls |
| **UI Design** | Enterprise-grade with tailored components | Bootstrap-based general styling |
| **Features** | Rich dashboard, analytics, agent building | Basic agent management and chat |
| **Code Organization** | Structured folders and TypeScript types | Single HTML file with embedded JS |
| **Maintainability** | High - modular, typed components | Low - tightly coupled code |
| **Extensibility** | High - component system | Limited - requires HTML modifications |

## Key Features

- **Landing Page**: Captivating introduction to SentinelStacks capabilities
- **Dashboard**: Enterprise command center for all AI agent operations
- **Agent Builder**: Intuitive interface for creating and configuring agents
- **Agent Management**: Monitoring, metrics, and lifecycle management
- **Enterprise Controls**: Role-based access, audit logs, and governance

## Technology Stack

- **Framework**: React with TypeScript
- **Styling**: Tailwind CSS with custom enterprise extensions
- **State Management**: Redux Toolkit
- **API Communication**: Axios
- **Visualization**: D3.js for enterprise-grade charts
- **Animation**: Framer Motion for meaningful transitions
- **3D Effects**: Three.js for data visualization
- **Build Tool**: Vite

## Development

### Prerequisites

- Node.js 16+
- npm or yarn

### Getting Started

1. Clone the repository
2. Navigate to the `web-ui` directory
3. Install dependencies: `npm install` or `yarn`
4. Start the development server: `npm run dev` or `yarn dev`
5. Open your browser to `http://localhost:5173`

### Running with Mock API

For development without a backend:

```bash
VITE_USE_MOCK_API=true npm run dev
```

Or use the provided script:

```bash
./run-fixed-app.sh
```

### Project Structure

```
web-ui/
├── public/          # Static assets
├── src/
│   ├── assets/      # Images, fonts, etc.
│   ├── components/  # Reusable UI components
│   ├── context/     # React context providers
│   ├── hooks/       # Custom React hooks
│   ├── pages/       # Top-level page components
│   ├── services/    # API and service integrations
│   ├── styles/      # Global styles and Tailwind config
│   └── utils/       # Utility functions
├── README.md        # Project documentation
└── package.json     # Project dependencies and scripts
```

## Design Principles

1. **Professional Sophistication**: Clean lines, thoughtful spacing, consistent color scheme
2. **Enterprise Reliability**: Visual cues that communicate stability, security, scalability
3. **Intuitive Innovation**: Futuristic elements that enhance understanding, not distract
4. **Narrative Visualization**: Interactive elements that tell the SentinelStacks story
5. **Actionable Excitement**: Converting interest into clear next steps

## Integration

The web UI connects to the SentinelStacks API server for all data operations. See the API documentation for details on available endpoints and authentication requirements.
