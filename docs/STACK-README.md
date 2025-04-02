# SentinelStacks - Stack Engine

SentinelStacks' Stack Engine enables you to define and execute multi-agent workflows where agents work together to accomplish complex tasks. This document provides guidance on using the stack functionality.

## What are Stacks?

Stacks are declarative definitions of multi-agent workflows. They define:

1. Which agents to use
2. How agents are connected (the data flow)
3. Parameters for each agent
4. Execution order dependencies

## Key Components

The Stack Engine consists of the following components:

- **StackSpec** - The core data structure defining a stack and its agents
- **DAG Runner** - Executes agents in topological order based on dependencies
- **Stack State Manager** - Manages execution state and data passing between agents
- **Agent Runtime** - Executes individual agents using the sentinel runtime

## Quick Start

### 1. Initialize a Stack

Create a new stack using the `stack init` command:

```bash
sentinel stack init my-analysis-stack
```

This will create a Stackfile.yaml in the current directory. You can also create a stack from a template:

```bash
sentinel stack init my-analysis-stack --template=analyzer
```

Or using natural language:

```bash
sentinel stack init nlp-stack --nl="Create a pipeline with a text extractor, a sentiment analyzer, and a summarizer agent"
```

### 2. Run a Stack

Execute a stack using the `stack run` command:

```bash
sentinel stack run -f Stackfile.yaml
```

You can also provide additional input data:

```bash
sentinel stack run -f Stackfile.yaml --input='{"text": "Sample text to analyze"}'
```

### 3. List Available Stacks

List stacks in the current directory and the global stacks directory:

```bash
sentinel stack list --all
```

### 4. Inspect a Stack

Examine the structure of a stack:

```bash
sentinel stack inspect Stackfile.yaml
```

Output the execution plan as a graph:

```bash
sentinel stack inspect Stackfile.yaml --format=dot > stack.dot
```

## Stack File Format

A Stackfile.yaml follows this structure:

```yaml
name: my-analysis-stack
description: A stack for analyzing text data
version: 1.0.0

agents:
  - id: text-extractor
    uses: text-extractor:latest
    params:
      format: "json"

  - id: sentiment-analyzer
    uses: sentiment-analyzer:latest
    inputFrom:
      - text-extractor
    params:
      model: "default"

  - id: summarizer
    uses: text-summarizer:latest
    inputFrom:
      - text-extractor
      - sentiment-analyzer
    params:
      style: "concise"
```

### Key Fields

- **name**: Unique identifier for the stack
- **description**: Human-readable description
- **version**: Semantic version of the stack
- **agents**: List of agent specifications

For each agent:
- **id**: Unique identifier for the agent within the stack
- **uses**: Reference to the agent image (name:tag)
- **inputFrom**: List of agent IDs to take input from
- **inputKey**: (Optional) Specific key to extract from source agent's output
- **outputKey**: (Optional) Key to store this agent's output under
- **params**: Custom parameters to pass to the agent
- **depends**: (Optional) Additional dependencies that don't involve data passing

## Examples

See the `examples/stacks/` directory for sample stacks that demonstrate different patterns:

- `simple_analysis.yaml` - Basic text analysis pipeline
- `data_pipeline.yaml` - Complex data processing pipeline
- `research_assistant.yaml` - Research and summarization workflow
- `chat_enhancement.yaml` - Conversational enhancement stack

## Best Practices

1. **Atomicity**: Design agents to do one thing well
2. **Clear Dependencies**: Make data dependencies explicit
3. **Error Handling**: Include validation agents for critical points
4. **Testing**: Test stacks with representative data
5. **Versioning**: Version your stacks and reference specific agent versions

## Advanced Features

### Context Propagation

Agents can access outputs from any of their dependencies using the `inputFrom` field. By default, all outputs from the dependency are provided, but you can use `inputKey` to select specific data.

### Execution Models

The Stack Engine supports different execution models:

- **Sequential**: Agents run in topological order (default)
- **Parallel**: Independent agents run concurrently
- **Conditional**: Some agents may be skipped based on conditions

### Custom Runtime Configuration

You can configure execution parameters using flags:

```bash
sentinel stack run -f Stackfile.yaml --timeout=30 --verbose
```

## Troubleshooting

### Common Issues

1. **Agent Not Found**: Ensure the agent is available locally or in the registry
2. **Dependency Cycles**: Check for circular dependencies between agents
3. **Execution Timeout**: Increase the timeout for complex stacks
4. **Data Format Mismatch**: Ensure agents can process the data format they receive

### Debugging

Use the `--verbose` flag to see detailed execution information:

```bash
sentinel stack run -f Stackfile.yaml --verbose
```
