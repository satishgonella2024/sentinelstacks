# Basic Chatbot Example

This example demonstrates how to build a simple chatbot agent using SentinelStacks with the Llama 3 model via Ollama.

## Sentinelfile

```yaml
name: basicchatbot
description: Create an agent that provides helpful, friendly chat responses with a unique personality.
capabilities:
  - Engage in casual conversation
  - Answer general knowledge questions
  - Remember context from earlier in the conversation
  - Provide thoughtful and nuanced responses
  - Maintain a consistent personality
  - Gracefully handle inappropriate requests
model: 
  base: llama3
state:
  - conversation_history
  - user_preferences
initialization:
  introduction: "Hello! I'm ready to assist you."
termination:
  farewell: "Thank you for chatting with me. Have a great day!"
tools:
  - web_search:
      purpose: For looking up factual information
  - calculator:
      purpose: For performing mathematical operations
  - datetime:
      purpose: For answering time-related queries
personality:
  tone: friendly
  response_length: medium
  memory_depth: 10
```

## Building the Agent

Build the chatbot agent using the following command:

```bash
./bin/sentinel build -t demo/chatbot:v1 -f examples/chatbot/Sentinelfile \
  --llm ollama --llm-endpoint https://your-ollama-endpoint --llm-model llama3
```

## Running the Agent

Run the chatbot with:

```bash
./bin/sentinel run demo/chatbot:v1
```

This will start an interactive session where you can chat with your agent.

## Features

- **Personality**: The agent maintains a friendly tone and consistent personality
- **Context Awareness**: Remembers previous exchanges in the conversation
- **Tool Usage**: Can access web search, calculator, and date/time tools
- **Response Quality**: Provides thoughtful and nuanced responses
- **Safety**: Gracefully handles inappropriate requests

## Customization

You can customize this chatbot by modifying the Sentinelfile:

- Change the `model.base` to use different LLM models
- Adjust the `personality` settings to change tone and response style
- Add or remove `capabilities` to match your use case
- Modify `tools` to grant access to different functions

## Next Steps

- Try building similar agents with different personalities
- Experiment with different LLM providers like Claude
- Add custom tools specific to your application
- Implement specialized knowledge bases for domain-specific chatbots 