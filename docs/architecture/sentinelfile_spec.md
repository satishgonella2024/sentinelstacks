# Sentinelfile Specification

## Overview

A Sentinelfile is a natural language definition file that describes an AI agent's capabilities, behavior, and requirements. Unlike traditional configuration files that use structured formats like YAML or JSON, Sentinelfiles use plain language that is parsed by an LLM to extract structured information.

## File Format

A Sentinelfile is a plain text file named `Sentinelfile` (no extension). While the format is flexible, the following structure is recommended for clarity:

```
# Sentinelfile for [Agent Name]

[High-level description of the agent's purpose]

The agent should have access to [resources/knowledge/tools].

It should be able to:
- [Capability 1]
- [Capability 2]
- [Capability 3]
...

The agent should use [LLM model] as its base model.

It should maintain state about [state information].

When the conversation starts, the agent should [initialization behavior].

When the conversation ends, the agent should [termination behavior].

Allow the agent to access the following tools:
- [Tool 1]
- [Tool 2]
...

Set [parameter name] to [value].
```

## Example Sentinelfile

```
# Sentinelfile for ResearchAssistant

Create an agent that helps with academic research tasks, finding relevant papers and summarizing their content.

The agent should have access to academic databases and search engines.

It should be able to:
- Search for academic papers on a given topic
- Summarize key findings from papers
- Extract methods and results sections
- Compare multiple papers
- Generate literature review outlines
- Create proper citations in multiple formats

The agent should use claude-3.7-sonnet as its base model.

It should maintain state about the current research topic, papers that have been reviewed, and key findings.

When the conversation starts, the agent should introduce itself as a research assistant and ask about the user's research area.

When the conversation ends, the agent should summarize the research findings and suggest next steps.

Allow the agent to access the following tools:
- Academic search
- PDF parser
- Citation generator
- Web browser
- Note taking

Set search_depth to 15 papers maximum.
Set citation_format to APA by default.
```

## Parsing Process

When a Sentinelfile is processed by the `sentinel build` command:

1. The natural language is sent to an LLM (default: Claude)
2. The LLM extracts structured information into a standardized JSON format
3. The structured definition is validated for completeness and consistency
4. The definition is packaged into a Sentinel Image

## Structured Output

While users write in natural language, the parser converts this to a structured format:

```json
{
  "name": "research-assistant",
  "description": "An agent that helps with academic research tasks",
  "baseModel": "claude-3.7-sonnet",
  "capabilities": [
    "academic_search",
    "summarization",
    "extraction",
    "comparison",
    "outline_generation",
    "citation_generation"
  ],
  "tools": [
    "academic_search",
    "pdf_parser",
    "citation_generator",
    "web_browser",
    "note_taking"
  ],
  "stateSchema": {
    "research_topic": "string",
    "reviewed_papers": "array",
    "key_findings": "map"
  },
  "lifecycle": {
    "initialization": "Introduce as research assistant and ask about research area",
    "termination": "Summarize findings and suggest next steps"
  },
  "parameters": {
    "search_depth": 15,
    "citation_format": "APA"
  }
}
```

## Best Practices

1. **Be Specific**: Clearly define capabilities and behavior
2. **Use Plain Language**: Avoid technical jargon unless necessary
3. **List Capabilities**: Bullet points help the parser identify discrete functions
4. **Specify Tools**: Explicitly list required tool access
5. **Define Initialization**: Describe how the agent should start conversations
6. **Set Parameters**: Include any configuration parameters with clear values
7. **Maintain Consistency**: Ensure descriptions don't contradict themselves

## Limitations

1. Highly technical or domain-specific terminology may be misinterpreted
2. Very complex agent behaviors might need to be simplified
3. Custom tool configurations may require additional specification
4. The parser may occasionally miss nuanced requirements

## Validation

SentinelStacks includes a validation step that:

1. Checks for internal consistency
2. Verifies that all required fields are present
3. Confirms that specified tools are available
4. Validates that the base model is supported
5. Ensures that parameters have valid values

You can manually validate a Sentinelfile with:

```bash
sentinel validate
```
