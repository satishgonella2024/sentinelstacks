# NLP-to-Agent Generator

This example demonstrates how to create SentinelStacks agents on-the-fly from natural language descriptions.

## Overview

The NLP-to-Agent Generator is a tool that allows users to describe agents in plain English. The tool processes the natural language description using an LLM, generates a valid Sentinelfile, and builds a working agent automatically.

This approach makes agent creation more accessible to non-technical users and enables rapid prototyping of agents.

## How It Works

The NLP-to-Agent Generator follows these steps:

1. Takes a natural language description of an agent's desired functionality
2. Sends the description to an LLM with a specialized prompt
3. The LLM converts the description into a structured Sentinelfile YAML
4. The generator creates the necessary files and builds the agent
5. The agent is ready to run

## Example Usage

```bash
# Run the NLP-to-Agent generator
go run examples/nlp-generator/nlp_to_agent.go
```

Then enter your natural language description, for example:

```
I want an agent that can help students with math homework.
It should be patient, able to break down complex problems
into simpler steps, provide explanations for each step,
and give examples when needed. It should also be able to
recognize when the student is struggling and offer
extra guidance.
```

After processing, the tool will generate a Sentinelfile and build the agent.

## Integration with SentinelStacks CLI

In a full implementation, this functionality could be integrated into the SentinelStacks CLI:

```bash
# Create an agent from natural language
sentinel create --from-nlp "Create a customer support agent that can..."

# Or interactively
sentinel create --from-nlp-file description.txt

# Or using a dialog
sentinel create --interactive
```

## Technical Details

The generator uses:

1. **Language Detection**: Identifies important entities, capabilities, and tools from the natural language
2. **Prompt Engineering**: Crafts a specific prompt to guide the LLM in generating YAML
3. **Schema Validation**: Ensures the generated YAML conforms to the Sentinelfile specification
4. **Post-Processing**: Augments the generated definition with additional metadata and improvements

## Extending the Generator

The generator can be extended to:

- Allow iterative refinement of generated agents
- Support specific templates or starting points
- Enable fine-tuning for domain-specific agents
- Include custom tool configuration based on natural language descriptions

## Benefits

- **Accessibility**: Makes agent creation accessible to non-technical users
- **Rapid Prototyping**: Quickly create agents for testing and exploration
- **Iteration**: Easily refine and improve agents through natural language feedback
- **Exploration**: Discover capabilities you might not have considered in manual creation 