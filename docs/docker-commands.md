# SentinelStacks Docker-Inspired Command Reference

SentinelStacks provides a Docker-inspired command interface for managing AI agents, networks, volumes, and multi-agent systems. This guide documents all the available commands and their usage.

## Setup

Before using the commands, make sure you have properly installed SentinelStacks:

```bash
# Build the Sentinel binary
go build -o sentinel main.go

# Verify installation
./sentinel version
```

## Network Commands

Networks enable communication between agents, allowing them to exchange information and collaborate.

### Creating Networks

```bash
# Create a simple network
./sentinel network create my-network

# Create a network with a specific driver
./sentinel network create my-network --driver advanced
```

### Listing Networks

```bash
# List all networks
./sentinel network ls
```

### Connecting Agents to Networks

```bash
# Connect an agent to a network
./sentinel network connect my-network agent-id
```

### Disconnecting Agents from Networks

```bash
# Disconnect an agent from a network
./sentinel network disconnect my-network agent-id
```

### Inspecting Networks

```bash
# View detailed information about a network
./sentinel network inspect my-network
```

### Removing Networks

```bash
# Remove a network
./sentinel network rm my-network

# Force remove a network even if it has connected agents
./sentinel network rm my-network --force
```

## Volume Commands

Volumes provide persistent memory for agents, allowing them to store and retrieve information across sessions.

### Creating Volumes

```bash
# Create a simple volume
./sentinel volume create my-volume

# Create a volume with a specific size
./sentinel volume create my-volume --size 2GB

# Create an encrypted volume
./sentinel volume create secure-volume --size 1GB --encrypted
```

### Listing Volumes

```bash
# List all volumes
./sentinel volume ls
```

### Mounting Volumes

```bash
# Mount a volume to an agent
./sentinel volume mount my-volume agent-id

# Mount a volume with a specific path
./sentinel volume mount my-volume agent-id --path /custom/path
```

### Unmounting Volumes

```bash
# Unmount a volume from an agent
./sentinel volume unmount my-volume agent-id
```

### Inspecting Volumes

```bash
# View detailed information about a volume
./sentinel volume inspect my-volume
```

### Removing Volumes

```bash
# Remove a volume
./sentinel volume rm my-volume

# Force remove a volume even if it is mounted
./sentinel volume rm my-volume --force
```

## Compose Commands

Compose commands enable managing multi-agent systems defined in YAML configuration files.

### Creating Multi-Agent Systems

```bash
# Create and start a multi-agent system defined in a YAML file
./sentinel compose up -f my-compose.yaml

# Create in detached mode
./sentinel compose up -f my-compose.yaml -d
```

### Listing Multi-Agent Systems

```bash
# List all running multi-agent systems
./sentinel compose ls
```

### Pausing Multi-Agent Systems

```bash
# Pause all agents in a multi-agent system
./sentinel compose pause system-id
```

### Resuming Multi-Agent Systems

```bash
# Resume all agents in a multi-agent system
./sentinel compose resume system-id
```

### Viewing Logs

```bash
# View logs from all agents in a multi-agent system
./sentinel compose logs system-id

# Follow log output
./sentinel compose logs system-id -f

# Show only the last N lines
./sentinel compose logs system-id --tail 100
```

### Stopping and Removing Multi-Agent Systems

```bash
# Stop and remove a multi-agent system
./sentinel compose down system-id

# Remove volumes as well
./sentinel compose down system-id --volumes
```

## System Commands

System commands provide monitoring and maintenance capabilities for your SentinelStacks installation.

### System Information

```bash
# Display system information
./sentinel system info
```

### Disk Usage

```bash
# Show disk usage by different components
./sentinel system df
```

### Removing Unused Data

```bash
# Remove unused data
./sentinel system prune

# Remove all unused data, including volumes
./sentinel system prune --all --volumes
```

### System Events

```bash
# View system events
./sentinel system events

# Filter events by type
./sentinel system events --filter volume

# Show events from the last hour
./sentinel system events --since 1h

# Show verbose event details
./sentinel system events -v
```

## Agent Interaction Commands

SentinelStacks also provides commands for interacting directly with agents.

### Executing Commands

```bash
# Execute a one-time command without creating an agent
./sentinel exec "What is the capital of France?"
```

### Interactive Shell

```bash
# Start an interactive shell with a running agent
./sentinel shell agent-id
```

## Registry Commands

Registry commands allow managing agent images in remote registries.

### Pulling Images

```bash
# Pull an agent image from a registry
./sentinel pull registry.example.com/research-assistant:latest
```

### Pushing Images

```bash
# Push an agent image to a registry
./sentinel push research-assistant:latest registry.example.com/research-assistant:latest
```

### Searching Images

```bash
# Search for agent images in a registry
./sentinel search "research assistant"
```

### Authentication

```bash
# Log in to a registry
./sentinel login registry.example.com

# Log out from a registry
./sentinel logout registry.example.com
```

## Compose File Format

Multi-agent systems are defined using YAML files similar to Docker Compose. Here's an example:

```yaml
name: research-team

networks:
  brain-net:
    driver: default
  data-net:
    driver: default

volumes:
  research-memory:
    size: 2GB
  output-memory:
    size: 1GB
    encrypted: true

agents:
  coordinator:
    image: sentinelstacks/agent:coordinator
    networks:
      - brain-net
      - data-net
    volumes:
      - research-memory:/memory
    environment:
      ROLE: coordinator
      TASK: research_coordination
    resources:
      memory: 1GB
      
  researcher:
    image: sentinelstacks/agent:researcher
    networks:
      - brain-net
    volumes:
      - research-memory:/memory/read-only
    environment:
      ROLE: researcher
      TOPIC: ai_safety
    resources:
      memory: 2GB
      
  writer:
    image: sentinelstacks/agent:writer
    networks:
      - brain-net
      - data-net
    volumes:
      - output-memory:/memory
    environment:
      ROLE: writer
      FORMAT: academic_paper
    resources:
      memory: 1GB
```

This configuration defines a multi-agent system with three agents (coordinator, researcher, and writer) connected via networks and sharing volumes.

## Troubleshooting

If you encounter issues with the SentinelStacks commands, try the following:

1. Check the system information:
   ```bash
   ./sentinel system info
   ```

2. View recent system events:
   ```bash
   ./sentinel system events --since 15m
   ```

3. Check if there are any network or volume conflicts:
   ```bash
   ./sentinel network ls
   ./sentinel volume ls
   ```

4. Make sure you have sufficient disk space:
   ```bash
   ./sentinel system df
   ```

## Configuration

SentinelStacks configuration is stored in `~/.sentinel/config.yaml`. You can modify this file to customize the behavior of the system.

## Where to Go Next

For more information on SentinelStacks:

- [Main Documentation](README.md)
- [Agent Development Guide](agent-development.md)
- [API Reference](api-reference.md)
- [Advanced Configuration](advanced-configuration.md)
