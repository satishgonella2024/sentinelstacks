import React, { useState } from 'react';
import ReactMarkdown from 'react-markdown';

interface DocSection {
  id: string;
  title: string;
  content: string;
}

const Documentation: React.FC = () => {
  const [activeSection, setActiveSection] = useState('getting-started');
  
  const sections: DocSection[] = [
    {
      id: 'getting-started',
      title: 'Getting Started',
      content: `
# Getting Started with SentinelStacks

SentinelStacks is an AI Agent Management Platform designed for enterprise use. This guide will help you get started with the platform.

## Setting up your first Agent

1. Navigate to the **Agents** page
2. Click the **Create New Agent** button
3. Fill in the required details:
   - Name: A descriptive name for your agent
   - Model: Select the AI model to power your agent
   - Description: A brief description of what your agent does
4. Click **Create Agent** to finish

Once created, your agent will be available for use and will appear in your dashboard.

## Configuring Agent Settings

You can configure various aspects of your agent:

- **Memory**: Adjust how much context your agent retains
- **Permissions**: Control what capabilities your agent has
- **Tools**: Enable specific tools and integrations
- **Fine-tuning**: Apply custom training to your agent
      `
    },
    {
      id: 'model-registry',
      title: 'Model Registry',
      content: `
# Using the Model Registry

The Model Registry allows you to connect to various model providers and manage which models are available to your agents.

## Connecting to a Model Provider

1. Navigate to the **Registry** page
2. Click **Connect New Registry**
3. Select a provider type (OpenAI, Anthropic, Hugging Face, etc.)
4. Enter your API credentials
5. Click **Connect**

Once connected, all available models from that provider will be listed and can be used with your agents.

## Syncing Models

To ensure you have the latest models available:

1. Go to the **Registry** page
2. Find the provider you want to update
3. Click the **Sync** button

This will fetch the latest models and update your registry.
      `
    },
    {
      id: 'conversations',
      title: 'Managing Conversations',
      content: `
# Managing Conversations

SentinelStacks provides tools to manage conversations between users and your AI agents.

## Viewing Conversations

1. Navigate to the **Agents** page
2. Select an agent
3. Click on the **Conversations** tab

Here you'll see all conversations that have occurred with this agent.

## Analyzing Conversation Data

SentinelStacks provides analytics on conversations:

- **Usage patterns**: See when and how often your agents are used
- **Common queries**: Identify frequently asked questions
- **Performance metrics**: Measure response times and user satisfaction
- **Error rates**: Track and troubleshoot agent errors

Use this data to improve your agents over time.
      `
    },
    {
      id: 'security',
      title: 'Security & Compliance',
      content: `
# Security and Compliance

SentinelStacks is built with enterprise security requirements in mind.

## Data Protection

- All data is encrypted in transit and at rest
- Personal Identifiable Information (PII) can be automatically redacted
- Data retention policies can be configured to meet your requirements

## Access Control

SentinelStacks provides role-based access control:

- **Admin**: Full access to all features
- **Manager**: Can create and configure agents
- **User**: Can interact with agents
- **Viewer**: Can view agents and analytics, but cannot make changes

## Compliance Features

- Audit logs for all system actions
- Export tools for compliance reviews
- Integration with enterprise SSO systems
      `
    },
    {
      id: 'api-reference',
      title: 'API Reference',
      content: `
# API Reference

SentinelStacks provides a comprehensive API for programmatic access to all features.

## Authentication

All API requests require authentication using an API key or JWT token:

\`\`\`
Authorization: Bearer YOUR_API_TOKEN
\`\`\`

## Core Endpoints

### Agents

- \`GET /v1/agents\` - List all agents
- \`POST /v1/agents\` - Create a new agent
- \`GET /v1/agents/:id\` - Get a specific agent
- \`PUT /v1/agents/:id\` - Update an agent
- \`DELETE /v1/agents/:id\` - Delete an agent

### Conversations

- \`GET /v1/agents/:id/conversations\` - List conversations for an agent
- \`POST /v1/agents/:id/conversations\` - Start a new conversation
- \`GET /v1/agents/:id/conversations/:convId\` - Get a specific conversation
- \`POST /v1/agents/:id/conversations/:convId/messages\` - Send a message

### Registry

- \`GET /v1/registry\` - List all connected registries
- \`POST /v1/registry\` - Connect a new registry
- \`GET /v1/registry/:id/models\` - List models from a registry
      `
    }
  ];
  
  const activeContent = sections.find(s => s.id === activeSection)?.content || '';
  
  return (
    <div className="container px-4 py-8">
      <h1 className="text-2xl font-bold mb-8">Documentation</h1>
      
      <div className="flex flex-col md:flex-row gap-6">
        {/* Sidebar */}
        <div className="w-full md:w-64 bg-gray-800 rounded-lg p-4">
          <h2 className="text-lg font-medium mb-4">Contents</h2>
          <nav>
            <ul className="space-y-1">
              {sections.map(section => (
                <li key={section.id}>
                  <button
                    className={`w-full text-left px-3 py-2 rounded ${
                      activeSection === section.id 
                        ? 'bg-primary-600 text-white' 
                        : 'text-gray-300 hover:bg-gray-700'
                    }`}
                    onClick={() => setActiveSection(section.id)}
                  >
                    {section.title}
                  </button>
                </li>
              ))}
            </ul>
          </nav>
        </div>
        
        {/* Content */}
        <div className="flex-1 bg-gray-800 rounded-lg p-6">
          <article className="prose prose-invert max-w-none">
            <ReactMarkdown>
              {activeContent}
            </ReactMarkdown>
          </article>
        </div>
      </div>
    </div>
  );
};

export default Documentation; 