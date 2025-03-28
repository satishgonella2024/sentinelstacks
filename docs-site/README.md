# SentinelStacks Documentation

This directory contains the SentinelStacks documentation website, built with [MkDocs](https://www.mkdocs.org/) and the [Material theme](https://squidfunk.github.io/mkdocs-material/).

## Local Development

### Prerequisites

- Python 3.7 or higher
- pip

### Setup

1. Install MkDocs and the Material theme:

```bash
pip install mkdocs-material
```

2. Run the local development server:

```bash
cd docs-site
mkdocs serve
```

3. Open your browser to [http://localhost:8000](http://localhost:8000)

The documentation will automatically reload when you save changes to any Markdown file.

## File Structure

- `mkdocs.yml`: Configuration file for the documentation site
- `docs/`: Contains all the Markdown files for the documentation
  - `index.md`: Homepage
  - `getting-started/`: Getting started guides
  - `user-guide/`: User documentation
  - `developer-guide/`: Documentation for contributors
  - `roadmap.md`: Development roadmap

## Adding Content

1. Create a new Markdown file in the appropriate directory
2. Add it to the navigation in `mkdocs.yml`
3. Write your content using Markdown

## Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch. The GitHub Action in `.github/workflows/docs.yml` handles the deployment process.

## Style Guide

- Use ATX-style headers (`#` for h1, `##` for h2, etc.)
- Code blocks should specify the language for syntax highlighting
- Use relative links for internal documentation references
- Include screenshots or diagrams for complex concepts
