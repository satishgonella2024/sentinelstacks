# Conversation Package

This package provides conversation history management for the Sentinel Stacks runtime.

## Overview

The conversation package implements:

1. A history manager that maintains the conversation between the user and the agent
2. Support for different message types (user, assistant, system)
3. Support for both text and multimodal messages (images, etc.)
4. Persistence of conversation history

## Components

### MessageType

Enum representing the type of message:
- `MessageTypeUser`: Messages from the user
- `MessageTypeAssistant`: Messages from the assistant
- `MessageTypeSystem`: System messages/prompts

### Message

Represents a single message in the conversation, containing:
- ID: Unique identifier
- Role: The sender's role (user, assistant, system)
- Content: Text content (for text messages)
- Contents: Array of multimodal contents (for multimodal messages)
- Timestamp: When the message was created

### History

Manages the conversation history, including:
- ID: The conversation session ID
- AgentID: The ID of the agent handling the conversation
- Messages: The array of messages in the conversation

## Usage

```go
// Create a new conversation history
history := conversation.NewHistory(agentID, sessionID)

// Add messages
history.AddUserMessage("Hello, how can you help me?")
history.AddAssistantMessage("I can assist with various tasks. What do you need help with?")
history.AddSystemMessage("The agent should be helpful and friendly.")

// Convert to multimodal input for sending to LLM
input, err := history.ToMultimodalInput(10) // Last 10 messages

// Save conversation to a file
err := history.SaveToFile("/path/to/conversation.json")

// Load conversation from a file
history, err := conversation.LoadFromFile("/path/to/conversation.json")
```

## Integration

This package integrates with the `multimodal` package to support rich media conversations including text, images, and other content types.