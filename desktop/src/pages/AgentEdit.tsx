import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import MainLayout from '../components/layout/MainLayout';
import PageContainer from '../components/layout/PageContainer';
import AgentForm from '../components/agents/forms/AgentForm';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { useAgent } from '../hooks/useAgent';
import { updateAgent } from '../services/agentService';

const AgentEdit: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  const { agent, isLoading, error: fetchError } = useAgent(id);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Combine fetch error and submission error
  useEffect(() => {
    if (fetchError) {
      setError(fetchError);
    }
  }, [fetchError]);

  const handleSubmit = async (formData: any) => {
    if (!id) return;
    
    setIsSubmitting(true);
    setError(null);

    try {
      // Validate form
      if (!formData.name.trim()) {
        throw new Error('Agent name is required');
      }

      // Update the agent
      await updateAgent(id, {
        name: formData.name,
        description: formData.description,
        model: formData.model,
        tools: formData.tools,
        capabilities: formData.capabilities,
        memory: formData.memory,
      });

      // Navigate back to the agent's detail page
      navigate(`/agents/${id}`);
    } catch (err) {
      setError(`Failed to update agent: ${err instanceof Error ? err.message : String(err)}`);
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return (
      <MainLayout>
        <div className="flex items-center justify-center h-screen">
          <LoadingSpinner size="lg" />
        </div>
      </MainLayout>
    );
  }

  if (!agent && !isLoading) {
    return (
      <MainLayout>
        <div className="flex flex-col items-center justify-center h-screen">
          <h2 className="text-xl font-semibold mb-4">Agent not found</h2>
          <Link
            to="/agents"
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
          >
            Back to Agents
          </Link>
        </div>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <PageContainer
        title={`Edit Agent: ${agent?.name}`}
        actions={
          <Link
            to={`/agents/${id}`}
            className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-200 dark:hover:bg-gray-700"
          >
            <ArrowLeftIcon className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
            Back to Agent
          </Link>
        }
      >
        {agent && (
          <AgentForm
            agent={agent}
            isSubmitting={isSubmitting}
            error={error}
            onSubmit={handleSubmit}
            cancelPath={`/agents/${id}`}
          />
        )}
      </PageContainer>
    </MainLayout>
  );
};

export default AgentEdit;