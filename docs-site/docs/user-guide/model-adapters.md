# Model Adapters

SentinelStacks supports multiple AI model providers through adapters. This page documents the available adapters and how to configure them.

## Supported Providers

| Provider | Models | Local/Remote | Function Calling |
|----------|--------|--------------|-----------------|
| Ollama | llama3, mistral, phi-2, etc. | Local | Basic |
| OpenAI | gpt-3.5-turbo, gpt-4, etc. | Remote | Advanced |
| Claude | claude-3-opus, claude-3-sonnet, etc. | Remote | Advanced |

## Ollama Adapter

Ollama allows you to run various open-source models locally on your machine.

### Configuration

```yaml
model:
  provider: "ollama"
  name: "llama3"
  options:
    temperature: 0.7
    top_p: 0.9
    max_tokens: 2000
```

### Available Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `temperature` | float | 0.7 | Controls randomness (0-1) |
| `top_p` | float | 0.9 | Nucleus sampling parameter |
| `top_k` | integer | 40 | Number of tokens to consider |
| `max_tokens` | integer | 2048 | Maximum generation length |

### Setup

1. Install Ollama from [ollama.ai](https://ollama.ai)
2. Pull the model you want to use:
   ```bash
   ollama pull llama3
   ```
3. Ensure Ollama is running in the background

## OpenAI Adapter

The OpenAI adapter connects to OpenAI's API for models like GPT-3.5 and GPT-4.

### Configuration

```yaml
model:
  provider: "openai"
  name: "gpt-4"
  options:
    temperature: 0.3
    max_tokens: 1000
```

### Available Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `temperature` | float | 0.7 | Controls randomness (0-1) |
| `max_tokens` | integer | 1000 | Maximum generation length |
| `top_p` | float | 1.0 | Nucleus sampling parameter |
| `presence_penalty` | float | 0.0 | Penalizes repeated tokens |
| `frequency_penalty` | float | 0.0 | Penalizes frequent tokens |

### Setup

1. Get an API key from [OpenAI](https://platform.openai.com)
2. Configure the key:
   ```bash
   sentinel config set openai.api_key=<your-api-key>
   ```

## Claude Adapter

The Claude adapter connects to Anthropic's Claude models.

### Configuration

```yaml
model:
  provider: "claude"
  name: "claude-3-opus"
  options:
    temperature: 0.5
    max_tokens: 1500
```

### Available Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `temperature` | float | 0.7 | Controls randomness (0-1) |
| `max_tokens` | integer | 1000 | Maximum generation length |
| `top_p` | float | 0.9 | Nucleus sampling parameter |
| `top_k` | integer | 40 | Number of tokens to consider |

### Setup

1. Get an API key from [Anthropic](https://anthropic.com)
2. Configure the key:
   ```bash
   sentinel config set claude.api_key=<your-api-key>
   ```

## Custom Adapter Configuration

You can override adapter settings when running an agent:

```bash
sentinel agent run my-assistant --model llama3 --endpoint http://custom-ollama-server:11434
```

## Adding New Adapters

SentinelStacks is designed to be extensible. To add support for new model providers, you can create custom adapters that implement the ModelAdapter interface.

See the [Developer Guide](../developer-guide/model-adapters.md) for information on creating custom adapters.
