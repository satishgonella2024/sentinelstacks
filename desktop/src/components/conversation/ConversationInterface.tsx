import React, { useState, useEffect, useRef } from 'react';
import { PaperAirplaneIcon, ArrowPathIcon, DocumentTextIcon } from '@heroicons/react/24/outline';
import { getAgentConversation } from '../../services/agentService';

interface Message {
  id: string;
  role: string;
  content: string;
  timestamp: string;
}

interface ConversationData {
  messages: Message[];
  metadata: {
    totalMessages: number;
    lastMessageAt: string;
  };
}

interface ConversationInterfaceProps {
  agentId: string;
}

const ConversationInterface: React.FC<ConversationInterfaceProps> = ({ agentId }) => {
  const [conversation, setConversation] = useState<ConversationData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState('');
  const [isSending, setIsSending] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetchConversation();
  }, [agentId]);

  useEffect(() => {
    scrollToBottom();
  }, [conversation]);

  const fetchConversation = async () => {
    setIsLoading(true);
    try {
      const data = await getAgentConversation(agentId);
      setConversation(data);
    } catch (err) {
      setError(`Failed to load conversation: ${err instanceof Error ? err.message : String(err)}`);
    } finally {
      setIsLoading(false);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = async () => {
    if (!message.trim()) return;
    
    setIsSending(true);
    
    // Simulate adding the user message
    const userMessage: Message = {
      id: `temp-${Date.now()}`,
      role: 'user',
      content: message,
      timestamp: new Date().toISOString()
    };
    
    // Update the UI optimistically
    setConversation(prev => {
      if (!prev) return {
        messages: [userMessage],
        metadata: {
          totalMessages: 1,
          lastMessageAt: userMessage.timestamp
        }
      };
      
      return {
        messages: [...prev.messages, userMessage],
        metadata: {
          totalMessages: prev.metadata.totalMessages + 1,
          lastMessageAt: userMessage.timestamp
        }
      };
    });
    
    setMessage('');
    
    // Simulate agent response (in a real app, you would call your API)
    setTimeout(() => {
      const agentMessage: Message = {
        id: `temp-${Date.now()}`,
        role: 'assistant',
        content: `This is a simulated response. In a real implementation, this would be handled by the agent's API.\n\nYou said: "${userMessage.content}"`,
        timestamp: new Date().toISOString()
      };
      
      setConversation(prev => {
        if (!prev) return {
          messages: [userMessage, agentMessage],
          metadata: {
            totalMessages: 2,
            lastMessageAt: agentMessage.timestamp
          }
        };
        
        return {
          messages: [...prev.messages, agentMessage],
          metadata: {
            totalMessages: prev.metadata.totalMessages + 1,
            lastMessageAt: agentMessage.timestamp
          }
        };
      });
      
      setIsSending(false);
    }, 1500);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 dark:bg-red-900 p-4 rounded-md">
        <p className="text-red-600 dark:text-red-200">{error}</p>
        <button
          onClick={fetchConversation}
          className="mt-2 inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          <ArrowPathIcon className="h-4 w-4 mr-1" />
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Conversation header */}
      <div className="bg-white dark:bg-gray-800 p-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex justify-between items-center">
          <h3 className="text-lg font-medium">Conversation</h3>
          {conversation && (
            <span className="text-sm text-gray-500 dark:text-gray-400">
              {conversation.metadata.totalMessages} messages
            </span>
          )}
        </div>
      </div>
      
      {/* Message history */}
      <div className="flex-grow overflow-y-auto p-4 bg-gray-50 dark:bg-gray-900">
        {!conversation || conversation.messages.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center p-4">
            <DocumentTextIcon className="h-12 w-12 text-gray-400 dark:text-gray-500 mb-2" />
            <p className="text-gray-600 dark:text-gray-400 mb-1">No messages yet</p>
            <p className="text-sm text-gray-500 dark:text-gray-500">
              Start a conversation with this agent by typing a message below.
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {conversation.messages.map((msg) => (
              <div
                key={msg.id}
                className={`p-4 rounded-lg max-w-3xl ${
                  msg.role === 'user'
                    ? 'ml-auto bg-primary-100 dark:bg-primary-900'
                    : 'mr-auto bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700'
                }`}
              >
                <div className="flex justify-between items-start mb-1">
                  <span className="font-medium capitalize text-sm">
                    {msg.role === 'user' ? 'You' : 'Agent'}
                  </span>
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    {new Date(msg.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                  </span>
                </div>
                <p className="whitespace-pre-wrap">{msg.content}</p>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>
      
      {/* Message input */}
      <div className="bg-white dark:bg-gray-800 p-4 border-t border-gray-200 dark:border-gray-700">
        <div className="flex items-end">
          <textarea
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Type your message..."
            className="flex-grow block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            rows={2}
            disabled={isSending}
          />
          <button
            onClick={handleSendMessage}
            disabled={isSending || !message.trim()}
            className="ml-3 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
          >
            {isSending ? (
              <ArrowPathIcon className="h-5 w-5 animate-spin" />
            ) : (
              <PaperAirplaneIcon className="h-5 w-5" />
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ConversationInterface;