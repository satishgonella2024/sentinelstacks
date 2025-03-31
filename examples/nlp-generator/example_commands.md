# NLP-to-Agent Generator Example Commands

This document provides examples of how to use the NLP-to-Agent Generator from the command line.

## Basic Usage

Build and run the example with different demo modes:

```bash
# Build the example
go build

# Run in interactive mode (default)
./nlp-generator

# Run the CLI integration demo
./nlp-generator -mode cli

# Run the package usage demo
./nlp-generator -mode package

# Run both demos
./nlp-generator -mode both
```

## Integration with SentinelStacks CLI

When integrated with the SentinelStacks CLI, you would use commands like:

```bash
# Create an agent from a natural language description
sentinel create --from-nlp "Create a customer service agent that helps users with product inquiries, returns, and order tracking. It should be friendly and efficient."

# Create an agent from a file containing a natural language description
sentinel create --from-nlp-file agent-description.txt

# Create an agent interactively
sentinel create --interactive

# Apply a template to the natural language description
sentinel create --from-nlp "Answer math questions and guide students through problem-solving" --template tutor

# Edit the generated Sentinelfile before building
sentinel create --from-nlp "Create a travel planning assistant" --edit

# Specify LLM provider and model
sentinel create --from-nlp "Create a code review assistant" --llm anthropic --llm-model claude-3-opus
```

## Natural Language Description Examples

Here are some examples of natural language descriptions you can use:

### Research Assistant

```
Create a research assistant that can help with academic research. It should be able to
summarize scientific papers, suggest relevant sources, analyze research methodologies,
and help formulate research questions. The agent should have a formal tone and provide
detailed, well-structured responses with citations when possible.
```

### Customer Support Bot

```
I need a customer support agent called "SupportBot" that can handle product inquiries,
troubleshooting, and returns. It should be friendly but efficient, with detailed knowledge
of our product catalog. The agent should collect relevant information from customers
before suggesting solutions and be able to escalate complex issues to human support.
```

### Educational Tutor

```
Create an educational tutor for high school mathematics. It should be patient and encouraging,
able to explain concepts in multiple ways for different learning styles. The tutor should
provide step-by-step explanations, generate practice problems, check students' work,
and adapt its teaching approach based on student performance and feedback.
```

### Fitness Coach

```
I want a virtual fitness coach that can create personalized workout plans, demonstrate
exercises with form tips, track progress, and provide motivational support. It should
adjust workouts based on user feedback and be knowledgeable about nutrition and recovery.
The coach should have an energetic and supportive personality.
```

## Example Templates

The system supports several templates that can be applied to natural language descriptions:

- `customer-service`: For creating customer service agents
- `research`: For creating research assistant agents
- `tutor`: For creating educational tutoring agents

Templates provide a standardized structure while allowing for customization through the natural language description. 