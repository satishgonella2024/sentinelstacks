# Team Collaboration Example

This example demonstrates how to build a team of collaborative agents using SentinelStacks, where multiple specialized agents work together under the coordination of a manager agent.

## Overview

The Team Collaboration example showcases a system of four interacting agents:

1. **Project Manager** - Coordinates the team, assigns tasks, and monitors progress
2. **Research Specialist** - Gathers and analyzes information
3. **Content Writer** - Creates and edits content based on research
4. **Quality Reviewer** - Evaluates content and provides improvement feedback

This multi-agent system demonstrates how SentinelStacks can be used to create complex workflows with specialized agents that communicate and collaborate to achieve a common goal.

## Sentinelfiles

### Project Manager Agent

```yaml
name: project-manager
description: A project manager agent that coordinates a team of specialized agents to complete complex tasks.
capabilities:
  - Task decomposition and delegation
  - Progress tracking and coordination
  - Decision making based on team input
  - Problem resolution and resource allocation
  - Quality assessment and feedback
model:
  base: claude3
  parameters:
    temperature: 0.2
    top_p: 0.9
# Additional configuration omitted for brevity
communication:
  channels:
    - agent_message_bus:
        access: read_write
        purpose: Primary communication channel with other agents
    - human_interface:
        access: read_write
        purpose: Communication with human supervisors
team:
  members:
    - researcher:
        role: Information gathering and analysis
        capabilities: ["research", "summarization", "fact-checking"]
    - writer:
        role: Content creation and editing
        capabilities: ["writing", "editing", "formatting"]
    - reviewer:
        role: Quality control and feedback
        capabilities: ["evaluation", "critique", "improvement suggestions"]
```

### Research Specialist Agent

```yaml
name: research-specialist
description: A specialized research agent that gathers and analyzes information for the project team.
capabilities:
  - Web research and information retrieval
  - Data analysis and pattern recognition
  - Fact-checking and source verification
  - Summarization and knowledge extraction
  - Trend identification and market analysis
# Additional configuration omitted for brevity
team:
  role: researcher
  reports_to: project-manager
  collaborates_with:
    - writer
    - reviewer
```

### Content Writer Agent

```yaml
name: content-writer
description: A specialized writing agent that creates, edits, and formats content based on research findings.
capabilities:
  - Professional content creation
  - Editing and proofreading
  - Formatting and document structure
  - Style adaptation for different audiences
  - Visual content suggestions
# Additional configuration omitted for brevity
team:
  role: writer
  reports_to: project-manager
  collaborates_with:
    - researcher
    - reviewer
```

### Quality Reviewer Agent

```yaml
name: quality-reviewer
description: A specialized quality control agent that evaluates content, provides feedback, and ensures high standards.
capabilities:
  - Critical evaluation and analysis
  - Constructive feedback formulation
  - Quality standards enforcement
  - Consistency verification
  - Improvement recommendations
# Additional configuration omitted for brevity
team:
  role: reviewer
  reports_to: project-manager
  collaborates_with:
    - researcher
    - writer
```

## Building and Running the Team

Build each agent separately:

```bash
# Build the manager agent
./bin/sentinel build -t demo/project-manager:v1 -f examples/team-collaboration/manager.Sentinelfile --llm anthropic --llm-model claude-3-sonnet

# Build the researcher agent
./bin/sentinel build -t demo/research-specialist:v1 -f examples/team-collaboration/researcher.Sentinelfile --llm anthropic --llm-model claude-3-sonnet

# Build the writer agent
./bin/sentinel build -t demo/content-writer:v1 -f examples/team-collaboration/writer.Sentinelfile --llm anthropic --llm-model claude-3-sonnet

# Build the reviewer agent
./bin/sentinel build -t demo/quality-reviewer:v1 -f examples/team-collaboration/reviewer.Sentinelfile --llm anthropic --llm-model claude-3-sonnet
```

Run the team with:

```bash
# Start the agent network with all four agents
./bin/sentinel team run --agents demo/project-manager:v1,demo/research-specialist:v1,demo/content-writer:v1,demo/quality-reviewer:v1
```

## Communication System

The agents communicate through a message bus system:

1. The **Project Manager** receives tasks from the human user
2. The manager breaks down the task and assigns subtasks to specialists
3. The **Research Specialist** gathers information and reports back
4. The **Content Writer** creates content based on the research
5. The **Quality Reviewer** evaluates the content and provides feedback
6. The writer revises based on feedback
7. The manager consolidates the final output and delivers it to the user

## Example Use Cases

- Creating research reports on complex topics
- Developing content marketing campaigns
- Producing technical documentation
- Analyzing and summarizing large volumes of information
- Generating multi-perspective analyses of issues

## Advanced Features

This example demonstrates several advanced features:

- **Multi-agent collaboration**: Agents working together in a coordinated workflow
- **Hierarchical organization**: Manager-subordinate relationship structure
- **Specialized roles**: Each agent has a specific function in the overall process
- **Communication channels**: Structured message passing between agents
- **Workflow orchestration**: Sequential and parallel task execution

## Customization

You can customize this team by:

- Adding more specialized agents
- Changing the team structure and reporting lines
- Adjusting the communication patterns
- Modifying the workflow processes
- Specializing the team for particular domains or tasks 