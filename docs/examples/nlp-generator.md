# NLP-to-Agent Generator

## Overview

The NLP-to-Agent Generator allows users to create SentinelStacks agents from natural language descriptions. This powerful tool bridges the gap between human language and agent definitions, making agent creation accessible to users without technical expertise in YAML or agent configuration.

## How It Works

The NLP-to-Agent process works as follows:

1. **Natural Language Input**: Users describe the agent they want to create in plain English.
2. **LLM Processing**: The description is sent to a Large Language Model (like Claude).
3. **YAML Generation**: The LLM generates a structured Sentinelfile YAML based on the description.
4. **Validation**: The generated YAML is validated against the Sentinelfile schema.
5. **File Creation**: The necessary files are created in the output directory.
6. **Agent Building**: The agent is built using the SentinelStacks CLI.

## Example Interaction

**User description:**
```
Create a travel planning assistant that can help users find flights, hotels, and attractions.
It should be friendly and conversational, with knowledge of popular destinations. The agent
should ask about budget, dates, and preferences, then provide personalized recommendations.
```

**Generated Sentinelfile:**
```yaml
name: travel-assistant
description: Travel planning assistant that helps with flights, hotels, and attractions
capabilities:
  - Find and recommend flights based on user preferences
  - Suggest hotels within user's budget
  - Recommend attractions and activities at destinations
  - Create personalized travel itineraries
model:
  base: llama3
  parameters:
    temperature: 0.7
    top_p: 0.95
state:
  - conversation_history
  - user_preferences:
      travel_dates: null
      destination: null
      budget: null
      interests: []
tools:
  - web_search:
      purpose: Looking up travel information
  - calculator:
      purpose: Calculating costs and budgets
  - date_time:
      purpose: Checking travel dates and availability
initialization:
  introduction: "Hello! I'm your travel planning assistant. I can help you find flights, hotels, and attractions for your next trip. Where are you thinking of going?"
termination:
  farewell: "Thank you for planning your trip with me. Have a wonderful journey!"
personality:
  tone: friendly
  style: conversational
  traits:
    - knowledgeable
    - helpful
    - enthusiastic
```

## Implementation

The NLP-to-Agent Generator is implemented in Go. Here's a simplified version of how it processes natural language input:

```go
func (g *Generator) ProcessNaturalLanguage(input string) (*SentinelfileResponse, error) {
    // Extract agent name and tag from input if possible
    g.extractNameAndTag(input)
    
    // Prepare prompt for LLM
    prompt := g.buildPrompt(input)
    
    // Send to LLM API
    yamlContent, metadata, err := g.callLLM(prompt)
    if err != nil {
        return nil, fmt.Errorf("LLM processing error: %v", err)
    }
    
    response := &SentinelfileResponse{
        Sentinelfile: yamlContent,
        Metadata:     metadata,
    }
    
    // Generate files if enabled
    if g.GenerateFiles {
        err = g.generateFiles(response)
        if err != nil {
            return response, fmt.Errorf("file generation error: %v", err)
        }
    }
    
    return response, nil
}
```

## CLI Integration

The NLP-to-Agent Generator is integrated with the SentinelStacks CLI, allowing users to create agents directly from natural language descriptions. The following commands are available:

```bash
# Create an agent from a natural language description
sentinel create --from-nlp "Create a customer service agent that helps with product inquiries"

# Create an agent from a file containing a natural language description
sentinel create --from-nlp-file agent-description.txt

# Create an agent interactively
sentinel create --interactive

# Apply a template to the natural language description
sentinel create --from-nlp "Answer math questions" --template tutor

# Edit the generated Sentinelfile before building
sentinel create --from-nlp "Create a travel planning assistant" --edit

# Specify LLM provider and model
sentinel create --from-nlp "Create a code review assistant" --llm anthropic --llm-model claude-3-opus
```

The CLI integration provides several options for customizing the agent creation process:

- **Templates**: Use pre-defined templates for common agent types
- **Interactive Mode**: Guided agent creation with prompts and suggestions
- **Edit Before Building**: Review and modify the generated Sentinelfile
- **LLM Selection**: Choose which LLM provider and model to use

## Benefits

The NLP-to-Agent Generator offers several advantages:

1. **Accessibility**: Non-technical users can create agents without learning YAML or understanding the full agent configuration.
2. **Rapid Prototyping**: Quickly test different agent concepts without writing detailed configuration files.
3. **Exploration**: Discover the capabilities of the agent system through natural language descriptions.

## Advanced Usage

For more advanced use cases, you can:

- **Use Templates**: Start with a template and customize it with natural language
- **Hybrid Approach**: Generate a base Sentinelfile with NLP and then refine it manually
- **Iterative Refinement**: Generate an agent, test it, and then regenerate with more specific instructions 