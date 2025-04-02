# SentinelStacks Visualizations

This directory contains Mermaid diagrams for visualizing the SentinelStacks architecture and roadmap. 

## Available Diagrams

1. **Phase Roadmap** (`phase-roadmap.mmd`) - A Gantt chart showing the development timeline and key milestones
2. **Enhanced Architecture** (`enhanced-architecture.mmd`) - A flowchart showing the overall system architecture with planned enhancements
3. **Stack Engine Detail** (`stack-engine-detail.mmd`) - A class diagram detailing the internal structure of the Stack Engine component

## Rendering the Diagrams

These diagrams are written in Mermaid markdown syntax and can be rendered in multiple ways:

### Using the Mermaid CLI

```bash
npx @mermaid-js/mermaid-cli -i phase-roadmap.mmd -o phase-roadmap.png
```

### Using a Mermaid-compatible Markdown Renderer

Many documentation systems support Mermaid natively, including:
- GitHub Markdown
- GitLab Markdown
- MkDocs with the mermaid plugin
- Docusaurus

### Using the Mermaid Live Editor

You can copy the content of any `.mmd` file into the [Mermaid Live Editor](https://mermaid.live/) to visualize and export the diagrams.

## Updating Diagrams

When making significant architectural changes, please update these diagrams to reflect the current design. The diagrams should be kept in sync with the implementation.

For more information about Mermaid syntax, see the [Mermaid documentation](https://mermaid-js.github.io/mermaid/).
