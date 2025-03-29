import React from 'react';
import MainLayout from '../components/layout/MainLayout';
import PageContainer from '../components/layout/PageContainer';
import AgentList from '../components/agents/AgentList';

const Agents: React.FC = () => {
  return (
    <MainLayout>
      <PageContainer>
        <AgentList />
      </PageContainer>
    </MainLayout>
  );
};

export default Agents;