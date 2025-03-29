import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import MainLayout from '../components/layout/MainLayout';
import PageContainer from '../components/layout/PageContainer';
import AgentForm from '../components/agents/forms/AgentForm';
import { createAgent } from '../services/agentService';

const AgentCreate: React.FC = () => {
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (formData: any) => {
    setIsSubmitting(true);
    setError(null);

    try {
      // Validate form
      if (!formData.name.trim()) {
        throw new Error('Agent name is required');
      }

      // Create the agent
      const newAgent = await createAgent({
        name: formData.name,
        description: formData.description,
        model: formData.model,
        tools: formData.tools,
        capabilities: formData.capabilities,
        memory: formData.memory,
        status: 'inactive',
      });

      // Navigate to the new agent's detail page
      navigate(`/agents/${newAgent.id}`);
    } catch (err) {
      setError(`Failed to create agent: ${err instanceof Error ? err.message : String(err)}`);
      setIsSubmitting(false);
    }
  };

  return (
    <MainLayout>
      <PageContainer
        title="Create New Agent"
        actions={
          <Link
            to="/agents"
            className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-200 dark:hover:bg-gray-700"
          >
            <ArrowLeftIcon className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
            Back to Agents
          </Link>
        }
      >
        <AgentForm
          isSubmitting={isSubmitting}
          error={error}
          onSubmit={handleSubmit}
          cancelPath="/agents"
        />
      </PageContainer>
    </MainLayout>
  );
};

export default AgentCreate;