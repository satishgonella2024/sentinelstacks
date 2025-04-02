# Network Module for SentinelStacks

This module implements Docker-inspired network management commands for SentinelStacks AI agents.

## Features

- Create, list, inspect, and remove networks
- Connect and disconnect agents to/from networks
- Full multimodal messaging support (text, images, audio, video, binary, JSON)
- Advanced configuration options for networks
- Messaging system with attachments

## Usage

```bash
# Create a network
sentinel network create <network_name> [--driver <driver>] [--formats text,image,audio] [--config '{"key": "value"}']

# List networks
sentinel network ls

# Inspect a network
sentinel network inspect <network_name>

# Connect an agent to a network
sentinel network connect <network_name> <agent_id>

# Disconnect an agent from a network
sentinel network disconnect <network_name> <agent_id>

# Remove a network
sentinel network rm <network_name> [--force]

# Send a message
sentinel network message send <network_name> <sender_id> --content "Hello" --format text [--attach "format:path"]

# List messages
sentinel network message ls <network_name>

# Get message details
sentinel network message get <network_name> <message_id>

# Update network configuration
sentinel network config <network_name> --config '{"key": "value"}'
```

## Supported Message Formats

- `text`: Plain text messages
- `image`: Image files
- `audio`: Audio files
- `video`: Video files
- `binary`: Binary data
- `json`: Structured JSON data

## Multimodal Support

This module provides a full implementation of multimodal messaging between agents, allowing them to exchange various types of content beyond just text.
