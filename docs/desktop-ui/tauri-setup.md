# Setting Up the SentinelStacks Desktop App with Tauri

This guide walks through the initial setup process for the SentinelStacks desktop application using Tauri, React, and TypeScript.

## Prerequisites

Before getting started, ensure you have the following installed:

- Node.js (v14 or later)
- Rust (latest stable version)
- Tauri CLI

### Installing Rust

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

### Installing Tauri CLI

```bash
npm install -g @tauri-apps/cli
```

## Project Setup

### 1. Create a new Tauri + React + TypeScript project

```bash
# Navigate to the project root
cd /Users/subrahmanyagonella/SentinelStacks

# Create the desktop UI directory
mkdir -p desktop

# Create a new Tauri project with React and TypeScript template
npm create tauri-app@latest desktop -- --template react-ts
```

### 2. Navigate to the project and install dependencies

```bash
cd desktop
npm install
```

### 3. Install required dependencies

```bash
# React Router for navigation
npm install react-router-dom

# Tailwind CSS for styling
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p

# UI component libraries
npm install @headlessui/react @heroicons/react

# State management
npm install zustand
```

### 4. Configure Tailwind CSS

Update the `tailwind.config.js` file:

```js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#E6F7F9',
          100: '#CCEFF4',
          200: '#99DFE9',
          300: '#66CFDE',
          400: '#33BFD3',
          500: '#00AFC8',
          600: '#008BA0',
          700: '#006878',
          800: '#004450',
          900: '#002228',
        },
      },
    },
  },
  plugins: [],
}
```

Update the `src/styles.css` file:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

### 5. Configure Tauri

Update the `src-tauri/tauri.conf.json` file:

```json
{
  "build": {
    "beforeDevCommand": "npm run dev",
    "beforeBuildCommand": "npm run build",
    "devPath": "http://localhost:1420",
    "distDir": "../dist",
    "withGlobalTauri": false
  },
  "package": {
    "productName": "SentinelStacks",
    "version": "1.0.0"
  },
  "tauri": {
    "allowlist": {
      "all": false,
      "shell": {
        "all": false,
        "open": true,
        "execute": true,
        "sidecar": true,
        "scope": [
          {
            "name": "sentinel",
            "cmd": "sentinel",
            "args": true
          }
        ]
      },
      "dialog": {
        "all": true
      },
      "fs": {
        "all": false,
        "readFile": true,
        "writeFile": true,
        "readDir": true,
        "copyFile": true,
        "createDir": true,
        "removeDir": true,
        "removeFile": true,
        "scope": ["$HOME/.sentinel/**"]
      },
      "path": {
        "all": true
      },
      "window": {
        "all": true
      }
    },
    "bundle": {
      "active": true,
      "icon": [
        "icons/32x32.png",
        "icons/128x128.png",
        "icons/128x128@2x.png",
        "icons/icon.icns",
        "icons/icon.ico"
      ],
      "identifier": "io.sentinelstacks.desktop",
      "targets": "all",
      "windows": {
        "certificateThumbprint": null,
        "digestAlgorithm": "sha256",
        "timestampUrl": ""
      }
    },
    "security": {
      "csp": null
    },
    "updater": {
      "active": false
    },
    "windows": [
      {
        "fullscreen": false,
        "height": 768,
        "resizable": true,
        "title": "SentinelStacks",
        "width": 1200,
        "minWidth": 800,
        "minHeight": 600
      }
    ]
  }
}
```

### 6. Set Up Project Structure

Create the following directory structure:

```
src/
├── assets/          # Static assets
├── components/      # Reusable UI components
│   ├── layout/      # Layout components
│   ├── agents/      # Agent-related components
│   ├── registry/    # Registry-related components
│   └── ui/          # Generic UI components
├── hooks/           # Custom hooks
├── pages/           # Main pages
├── services/        # API integration
├── store/           # State management
├── types/           # TypeScript types
└── utils/           # Utility functions
```

Run the following commands to create this structure:

```bash
mkdir -p src/{assets,components,hooks,pages,services,store,types,utils}
mkdir -p src/components/{layout,agents,registry,ui}
```

### 7. Create Rust Command Interfaces

Update the Rust code in `src-tauri/src/main.rs` to interface with the SentinelStacks CLI:

```rust
#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use std::process::Command;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct Agent {
    name: String,
    version: String,
    description: Option<String>,
    model: Option<String>,
    memory_type: Option<String>,
}

#[tauri::command]
fn list_agents() -> Result<Vec<Agent>, String> {
    let output = Command::new("sentinel")
        .args(["registry", "list", "--json"])
        .output()
        .map_err(|e| format!("Failed to execute command: {}", e))?;

    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr);
        return Err(format!("Command exited with error: {}", stderr));
    }

    let stdout = String::from_utf8_lossy(&output.stdout);
    let agents: Vec<Agent> = serde_json::from_str(&stdout)
        .map_err(|e| format!("Failed to parse output: {}", e))?;

    Ok(agents)
}

#[tauri::command]
fn create_agent(name: String, description: String, model: String, memory_type: String) -> Result<(), String> {
    let output = Command::new("sentinel")
        .args([
            "agent", 
            "create", 
            "--name", &name,
            "--description", &description,
            "--model", &model,
            "--memory", &memory_type
        ])
        .output()
        .map_err(|e| format!("Failed to execute command: {}", e))?;

    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr);
        return Err(format!("Command exited with error: {}", stderr));
    }

    Ok(())
}

#[tauri::command]
fn run_agent(name: String, version: Option<String>) -> Result<(), String> {
    let mut args = vec!["agent", "run", "--name", &name];
    let version_arg: String;
    
    if let Some(ver) = version {
        version_arg = format!("--version={}", ver);
        args.push(&version_arg);
    }

    let output = Command::new("sentinel")
        .args(args)
        .output()
        .map_err(|e| format!("Failed to execute command: {}", e))?;

    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr);
        return Err(format!("Command exited with error: {}", stderr));
    }

    Ok(())
}

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            list_agents,
            create_agent,
            run_agent,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

### 8. Setting Up the Main Layout

Create a main layout component:

```tsx
// src/components/layout/MainLayout.tsx
import React, { useState, useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';

const MainLayout: React.FC = () => {
  const [darkMode, setDarkMode] = useState<boolean>(false);
  
  useEffect(() => {
    // Check system preference
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    setDarkMode(prefersDark);
    
    // Listen for changes
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = (e: MediaQueryListEvent) => {
      setDarkMode(e.matches);
    };
    
    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);
  
  return (
    <div className={`h-screen flex flex-col ${darkMode ? 'dark' : ''}`}>
      <Header toggleDarkMode={() => setDarkMode(!darkMode)} />
      <div className="flex-1 flex overflow-hidden">
        <Sidebar />
        <main className="flex-1 overflow-auto bg-slate-100 dark:bg-slate-900 text-slate-900 dark:text-white">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default MainLayout;
```

### 9. Setting Up Routing

Create a routes configuration:

```tsx
// src/App.tsx
import { useState } from 'react';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import MainLayout from './components/layout/MainLayout';
import AgentsPage from './pages/AgentsPage';
import RegistryPage from './pages/RegistryPage';
import HistoryPage from './pages/HistoryPage';
import SettingsPage from './pages/SettingsPage';
import AgentDetailPage from './pages/AgentDetailPage';

const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <AgentsPage />,
      },
      {
        path: 'agents/:agentId',
        element: <AgentDetailPage />,
      },
      {
        path: 'registry',
        element: <RegistryPage />,
      },
      {
        path: 'history',
        element: <HistoryPage />,
      },
      {
        path: 'settings',
        element: <SettingsPage />,
      },
    ],
  },
]);

function App() {
  return <RouterProvider router={router} />;
}

export default App;
```

### 10. Creating a Global Store with Zustand

```tsx
// src/store/agentStore.ts
import { create } from 'zustand';
import { invoke } from '@tauri-apps/api/tauri';

export interface Agent {
  name: string;
  version: string;
  description?: string;
  model?: string;
  memory_type?: string;
}

interface AgentState {
  agents: Agent[];
  loading: boolean;
  error: string | null;
  selectedAgent: Agent | null;
  fetchAgents: () => Promise<void>;
  createAgent: (name: string, description: string, model: string, memoryType: string) => Promise<void>;
  runAgent: (name: string, version?: string) => Promise<void>;
  selectAgent: (agent: Agent) => void;
}

export const useAgentStore = create<AgentState>((set, get) => ({
  agents: [],
  loading: false,
  error: null,
  selectedAgent: null,
  
  fetchAgents: async () => {
    set({ loading: true, error: null });
    try {
      const agents = await invoke<Agent[]>('list_agents');
      set({ agents, loading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : String(error), loading: false });
    }
  },
  
  createAgent: async (name, description, model, memoryType) => {
    set({ loading: true, error: null });
    try {
      await invoke('create_agent', { name, description, model, memoryType });
      // Refresh agent list after creation
      await get().fetchAgents();
    } catch (error) {
      set({ error: error instanceof Error ? error.message : String(error), loading: false });
    }
  },
  
  runAgent: async (name, version) => {
    set({ loading: true, error: null });
    try {
      await invoke('run_agent', { name, version });
      set({ loading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : String(error), loading: false });
    }
  },
  
  selectAgent: (agent) => {
    set({ selectedAgent: agent });
  },
}));
```

### 11. Create Basic Pages

Create placeholder pages for each main section:

```tsx
// src/pages/AgentsPage.tsx
import React, { useEffect } from 'react';
import { useAgentStore } from '../store/agentStore';

const AgentsPage: React.FC = () => {
  const { agents, loading, error, fetchAgents } = useAgentStore();
  
  useEffect(() => {
    fetchAgents();
  }, [fetchAgents]);
  
  if (loading) return <div className="p-6">Loading agents...</div>;
  if (error) return <div className="p-6 text-red-500">Error: {error}</div>;
  
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">My Agents</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {agents.map((agent) => (
          <div key={agent.name} className="bg-white dark:bg-slate-800 rounded-lg shadow-md overflow-hidden">
            <div className="bg-gradient-to-r from-cyan-500 to-blue-500 h-2"></div>
            <div className="p-4">
              <div className="font-bold text-lg">{agent.name}</div>
              <div className="text-slate-500 dark:text-slate-400 text-sm mb-3">
                {agent.description || "No description"}
              </div>
              <div className="flex justify-between">
                <button className="bg-cyan-100 dark:bg-cyan-900 text-cyan-700 dark:text-cyan-200 px-3 py-1 text-sm rounded">
                  Run
                </button>
                <button className="bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-200 px-3 py-1 text-sm rounded">
                  Edit
                </button>
              </div>
            </div>
          </div>
        ))}
        
        {/* Add Agent Card */}
        <div className="bg-white dark:bg-slate-800 rounded-lg shadow-md overflow-hidden border-2 border-dashed border-slate-300 dark:border-slate-700 flex items-center justify-center h-48">
          <div className="text-center">
            <div className="w-12 h-12 bg-slate-200 dark:bg-slate-700 rounded-full flex items-center justify-center mx-auto mb-2">
              <svg className="w-6 h-6 text-slate-500 dark:text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <div className="text-slate-500 dark:text-slate-400">Create New Agent</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AgentsPage;
```

### 12. Running the App

Start the development server:

```bash
npm run tauri dev
```

## Building for Production

When you're ready to build the production version:

```bash
npm run tauri build
```

This will create executable files for your platform in the `src-tauri/target/release` directory.

## Next Steps

- Implement the remaining UI components
- Add more Tauri commands for full CLI integration
- Implement the agent chat interface
- Create the agent creation wizard
- Build the registry browsing interface
- Add settings management
