# Test Agent

This is a simple test agent for SentinelStacks. It demonstrates the basic functionality of the agent system.

## Configuration

The agent is defined in two files:

- `agentfile.natural.txt`: Natural language description of the agent
- `agentfile.yaml`: Structured YAML configuration (generated from the natural language description)

## Capabilities

This agent has the following capabilities:

- Basic conversation
- Remembering previous interactions
- Simple question answering

## Running the Agent

To run this agent, use the following command from the root of the SentinelStacks repository:

```bash
./sentinel agent run test-agent
```

## Testing Prompts

Here are some example prompts to try with the agent:

```
Hello, who are you?
```

```
What can you do?
```

```
Remember that my favorite color is blue.
```

```
What is my favorite color?
```

## State Management

The agent's state is stored in the `agent.state.json` file. This includes conversation history and any variables that the agent needs to remember between sessions.
