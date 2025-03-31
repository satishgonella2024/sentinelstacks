# LLM Provider Configuration

SentinelStacks supports multiple LLM (Large Language Model) providers for building and running agents. This guide explains how to configure different providers.

## Supported Providers

Currently, SentinelStacks supports the following LLM providers:

1. **Claude** - Anthropic's Claude models
2. **Ollama** - Self-hosted models through Ollama, including Llama and other open-source models
3. **OpenAI** - Support planned but not yet implemented

## Configuration

You can configure LLM providers through the CLI configuration:

```bash
# Set the default provider
sentinel config set llm.provider [provider_name]

# Set the API key (if required)
sentinel config set llm.api_key [your_api_key]

# Set a custom endpoint (for self-hosted models)
sentinel config set [provider_name].endpoint [endpoint_url]

# Set the default model
sentinel config set llm.model [model_name]
```

You can also override these settings on a per-command basis:

```bash
sentinel build -t my-image --llm ollama --llm-endpoint https://example.com --llm-model llama3
```

## Provider-Specific Configuration

### Claude

Claude requires an API key from Anthropic.

```bash
# Set Claude as the default provider
sentinel config set llm.provider claude

# Set your Claude API key
sentinel config set llm.api_key your_api_key_here

# Set the default Claude model
sentinel config set llm.model claude-3.7-sonnet
```

### Ollama

Ollama can be used with a local installation or a remote endpoint.

```bash
# Set Ollama as the default provider
sentinel config set llm.provider ollama

# For local Ollama installation (default)
# No additional configuration needed

# For remote Ollama endpoint
sentinel config set ollama.endpoint https://your-ollama-endpoint.com

# Set the model to use
sentinel config set llm.model llama3
```

#### Available Ollama Models

The models available depend on your Ollama installation. Common models include:

- `llama3`
- `llama3:8b`
- `llama3:70b` 
- `mistral`
- `mixtral`
- `gemma`

Check your Ollama installation for a complete list of available models.

## Using Different Providers for Different Commands

You can use different providers for different stages in your workflow:

```bash
# Build using Claude for parsing
sentinel build -t my-agent:v1 --llm claude

# Run using Ollama for execution
sentinel run my-agent:v1 --llm ollama --llm-model llama3
```

## Environment Variables

You can also use environment variables to set LLM configurations:

- `SENTINEL_LLM_PROVIDER` - The LLM provider to use
- `SENTINEL_LLM_ENDPOINT` - Custom endpoint for the LLM provider
- `SENTINEL_LLM_MODEL` - The model to use
- `SENTINEL_API_KEY` - API key for the LLM provider

Example:

```bash
export SENTINEL_LLM_PROVIDER=ollama
export SENTINEL_LLM_ENDPOINT=https://model.gonella.co.uk
export SENTINEL_LLM_MODEL=llama3
sentinel build -t my-agent:v1
``` 