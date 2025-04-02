# Registry Sharing Example

This example demonstrates how to share stacks and agents through the SentinelStacks registry.

## Prerequisites

- SentinelStacks CLI installed
- Registry account (sign up at https://registry.sentinelstacks.io)

## Step 1: Create a Simple Stack

We'll create a simple stack for text analysis:

```bash
# Create a new stack with the analyzer template
sentinel stack init text-analyzer --template analyzer

# This creates a Stackfile.yaml with three agents:
# - data-processor
# - analyzer 
# - summarizer
```

## Step 2: Push the Stack to the Registry

```bash
# Log in to the registry
sentinel login

# Push the stack to the registry
sentinel stack push Stackfile.yaml --sign

# The stack is now available in the registry under your username
```

## Step 3: Pull the Stack from Another Machine

```bash
# On another machine, log in to the registry
sentinel login

# Search for the stack
sentinel stack search text-analyzer

# Pull the stack
sentinel stack pull username/text-analyzer:latest --extract-agents

# This will:
# 1. Download the stack definition
# 2. Verify its signature
# 3. Pull all required agents
```

## Step 4: Run the Stack

```bash
# Run the stack with some input
sentinel stack run Stackfile.yaml --input="This is a sample text for analysis"

# The result will be a summary and analysis of the text
```

## Custom Formats

This example uses our standardized file formats:

- Text analyzer stack package: `text-analyzer-1.0.0.stack.sntl`
- Stack definition file: `text-analyzer.stack.yaml`
- Agent package example: `data-processor-1.0.0.agent.sntl`

## Package Structure

Our stack package contains:

```
sentinel.manifest.json             # Package manifest
text-analyzer.stack.yaml          # Main stack definition
README.md                         # Documentation
examples/                         # Example inputs and outputs
```

## Security

All packages are signed with developer keys and can be verified:

```bash
# Verify a package
sentinel verify text-analyzer-1.0.0.stack.sntl

# Import a trusted key
sentinel key import colleague-key.pub --id colleague
```

## Example Registry Commands

The SentinelStacks registry system uses a familiar command interface:

```bash
# List your stacks
sentinel stack list --registry

# Get versions of a stack
sentinel stack versions text-analyzer

# Check for updates
sentinel registry outdated
```
