package network

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// NewNetworkCmd creates the network command
func NewNetworkCmd(dataDir string) *cobra.Command {
	networkManager, err := NewNetworkManager(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing network manager: %v\n", err)
		os.Exit(1)
	}
	
	// Initialize messaging system
	messagingSystem := NewMessagingSystem(networkManager)
	
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage agent networks",
		Long:  `Create and manage networks for agent-to-agent communication`,
	}

	// Add network management subcommands
	cmd.AddCommand(newNetworkCreateCmd(networkManager))
	cmd.AddCommand(newNetworkListCmd(networkManager))
	cmd.AddCommand(newNetworkConnectCmd(networkManager))
	cmd.AddCommand(newNetworkDisconnectCmd(networkManager))
	cmd.AddCommand(newNetworkRemoveCmd(networkManager))
	cmd.AddCommand(newNetworkInspectCmd(networkManager))
	
	// Add multimodal messaging subcommands
	cmd.AddCommand(newNetworkMessageCmd(networkManager, messagingSystem))
	cmd.AddCommand(newNetworkConfigCmd(networkManager))

	return cmd
}

// newNetworkCreateCmd creates the network create command
func newNetworkCreateCmd(networkManager *NetworkManager) *cobra.Command {
	var (
		driver  string
		formats []string
		configStr string
	)

	cmd := &cobra.Command{
		Use:   "create [network_name]",
		Short: "Create a new agent network",
		Long:  `Create a new network for agents to communicate with each other`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			
			// Parse config if provided
			config := make(map[string]interface{})
			if configStr != "" {
				if err := json.Unmarshal([]byte(configStr), &config); err != nil {
					return fmt.Errorf("invalid config JSON: %w", err)
				}
			}
			
			// Add formats to config if provided
			if len(formats) > 0 {
				config["supported_formats"] = formats
			}
			
			fmt.Printf("Creating network '%s' with driver '%s'\n", networkName, driver)
			
			// Create the network
			network, err := networkManager.CreateNetwork(networkName, driver, config)
			if err != nil {
				return fmt.Errorf("failed to create network: %w", err)
			}
			
			fmt.Printf("Network '%s' created successfully with ID %s\n", network.Name, network.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&driver, "driver", "default", "Network driver to use")
	cmd.Flags().StringSliceVar(&formats, "formats", []string{"text"}, "Supported message formats (text, image, audio, video, binary, json)")
	cmd.Flags().StringVar(&configStr, "config", "", "Network configuration as JSON string")
	return cmd
}

// newNetworkListCmd creates the network list command
func newNetworkListCmd(networkManager *NetworkManager) *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List networks",
		Long:    `List all networks available for agent communication`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Listing all networks:")
			
			// List networks
			networks, err := networkManager.ListNetworks()
			if err != nil {
				return fmt.Errorf("failed to list networks: %w", err)
			}
			
			if len(networks) == 0 {
				fmt.Println("No networks found")
				return nil
			}
			
			// Print header
			fmt.Printf("%-20s %-12s %-10s %-15s %-10s %s\n", "NAME", "ID", "DRIVER", "FORMATS", "STATUS", "AGENTS")
			fmt.Printf("%-20s %-12s %-10s %-15s %-10s %s\n", "----", "--", "------", "-------", "------", "------")
			
			// Print networks
			for _, network := range networks {
				idDisplay := network.ID
				if len(idDisplay) > 8 {
					idDisplay = idDisplay[:8]
				}
				
				// Format supported formats
				formats := "text"
				if len(network.SupportedFormats) > 0 {
					if len(network.SupportedFormats) <= 2 {
						formats = strings.Join(network.SupportedFormats, ",")
					} else {
						formats = fmt.Sprintf("%s,+%d", 
							strings.Join(network.SupportedFormats[:2], ","), 
							len(network.SupportedFormats)-2)
					}
				}
				
				fmt.Printf("%-20s %-12s %-10s %-15s %-10s %d\n", 
					network.Name,
					idDisplay,
					network.Driver,
					formats,
					network.Status,
					len(network.Agents))
			}
			
			return nil
		},
	}
}

// newNetworkConnectCmd creates the network connect command
func newNetworkConnectCmd(networkManager *NetworkManager) *cobra.Command {
	return &cobra.Command{
		Use:   "connect [network_name] [agent_id]",
		Short: "Connect an agent to a network",
		Long:  `Connect an existing agent to a specified network`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			agentID := args[1]
			
			fmt.Printf("Connecting agent '%s' to network '%s'\n", agentID, networkName)
			
			// Connect agent to network
			if err := networkManager.ConnectAgent(networkName, agentID); err != nil {
				return fmt.Errorf("failed to connect agent to network: %w", err)
			}
			
			fmt.Printf("Agent '%s' successfully connected to network '%s'\n", agentID, networkName)
			return nil
		},
	}
}

// newNetworkDisconnectCmd creates the network disconnect command
func newNetworkDisconnectCmd(networkManager *NetworkManager) *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect [network_name] [agent_id]",
		Short: "Disconnect an agent from a network",
		Long:  `Disconnect an agent from a specified network`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			agentID := args[1]
			
			fmt.Printf("Disconnecting agent '%s' from network '%s'\n", agentID, networkName)
			
			// Disconnect agent from network
			if err := networkManager.DisconnectAgent(networkName, agentID); err != nil {
				return fmt.Errorf("failed to disconnect agent from network: %w", err)
			}
			
			fmt.Printf("Agent '%s' successfully disconnected from network '%s'\n", agentID, networkName)
			return nil
		},
	}
}

// newNetworkRemoveCmd creates the network remove command
func newNetworkRemoveCmd(networkManager *NetworkManager) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "rm [network_name]",
		Aliases: []string{"remove"},
		Short:   "Remove a network",
		Long:    `Remove a specified network`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			
			fmt.Printf("Removing network '%s'\n", networkName)
			
			// Get network by name
			network, err := networkManager.GetNetworkByName(networkName)
			if err != nil {
				return fmt.Errorf("failed to find network: %w", err)
			}
			
			// Remove network
			if err := networkManager.DeleteNetwork(network.ID, force); err != nil {
				return fmt.Errorf("failed to remove network: %w", err)
			}
			
			fmt.Printf("Network '%s' successfully removed\n", networkName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force removal even if network has connected agents")
	return cmd
}

// newNetworkInspectCmd creates the network inspect command
func newNetworkInspectCmd(networkManager *NetworkManager) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect [network_name]",
		Short: "Display detailed information on a network",
		Long:  `Display detailed information about a network, including connected agents`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			
			fmt.Printf("Inspecting network '%s'\n\n", networkName)
			
			// Get network details
			network, err := networkManager.GetNetworkByName(networkName)
			if err != nil {
				return fmt.Errorf("failed to inspect network: %w", err)
			}
			
			// Display network details
			fmt.Printf("Network: %s\n", network.Name)
			fmt.Printf("  ID: %s\n", network.ID)
			fmt.Printf("  Created: %s\n", network.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("  Status: %s\n", network.Status)
			fmt.Printf("  Driver: %s\n", network.Driver)
			
			// Display supported formats
			fmt.Printf("  Supported Formats: %s\n", strings.Join(network.SupportedFormats, ", "))
			
			// Display configuration if present
			if network.Config != nil && len(network.Config) > 0 {
				fmt.Println("  Configuration:")
				configData, _ := json.MarshalIndent(network.Config, "    ", "  ")
				fmt.Println(string(configData))
			}
			
			if len(network.Agents) > 0 {
				fmt.Printf("  Connected Agents (%d):\n", len(network.Agents))
				for _, agentID := range network.Agents {
					fmt.Printf("    - %s\n", agentID)
				}
			} else {
				fmt.Println("  No connected agents")
			}
			
			return nil
		},
	}
}

// newNetworkMessageCmd creates the network message command
func newNetworkMessageCmd(networkManager *NetworkManager, messagingSystem *MessagingSystem) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "message",
		Short: "Manage agent messages",
		Long:  `Send and receive messages between agents in a network`,
	}
	
	// Add message subcommands
	cmd.AddCommand(newMessageSendCmd(networkManager, messagingSystem))
	cmd.AddCommand(newMessageListCmd(networkManager, messagingSystem))
	cmd.AddCommand(newMessageGetCmd(networkManager, messagingSystem))
	
	return cmd
}

// newMessageSendCmd creates the message send command
func newMessageSendCmd(networkManager *NetworkManager, messagingSystem *MessagingSystem) *cobra.Command {
	var (
		format       string
		content      string
		attachments  []string
		metadataStr  string
	)

	cmd := &cobra.Command{
		Use:   "send [network_name] [sender_id]",
		Short: "Send a message to a network",
		Long:  `Send a message from an agent to a network`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			senderID := args[1]
			
			// Verify format
			if format == "" {
				format = "text"
			}
			
			// Create message
			message := Message{
				SenderID:    senderID,
				Format:      MessageFormat(format),
				Content:     content,
				Timestamp:   time.Now(),
				Attachments: []Attachment{},
			}
			
			// Parse metadata if provided
			if metadataStr != "" {
				metadata := make(map[string]interface{})
				if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
					return fmt.Errorf("invalid metadata JSON: %w", err)
				}
				message.Metadata = metadata
			}
			
			// Process attachments if provided
			if len(attachments) > 0 {
				for _, attachment := range attachments {
					parts := strings.SplitN(attachment, ":", 2)
					if len(parts) != 2 {
						return fmt.Errorf("invalid attachment format: %s (should be 'format:path')", attachment)
					}
					
					info, err := os.Stat(parts[1])
					if err != nil {
						return fmt.Errorf("attachment file not found: %s", parts[1])
					}
					
					message.Attachments = append(message.Attachments, Attachment{
						Format: MessageFormat(parts[0]),
						Path:   parts[1],
						Name:   filepath.Base(parts[1]),
						Size:   info.Size(),
					})
				}
			}
			
			fmt.Printf("Sending %s message from '%s' to network '%s'\n", format, senderID, networkName)
			
			// Send message
			sentMessage, err := messagingSystem.SendMessage(networkName, message)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			
			fmt.Printf("Message sent successfully with ID %s\n", sentMessage.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "text", "Message format (text, image, audio, video, binary, json)")
	cmd.Flags().StringVar(&content, "content", "", "Message content (text or file path)")
	cmd.Flags().StringSliceVar(&attachments, "attach", []string{}, "Attachments in format 'format:path'")
	cmd.Flags().StringVar(&metadataStr, "metadata", "", "Message metadata as JSON string")
	cmd.MarkFlagRequired("content")
	
	return cmd
}

// newMessageListCmd creates the message list command
func newMessageListCmd(networkManager *NetworkManager, messagingSystem *MessagingSystem) *cobra.Command {
	var (
		limit  int
		offset int
	)

	cmd := &cobra.Command{
		Use:   "ls [network_name]",
		Short: "List messages in a network",
		Long:  `List messages sent in a specified network`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			
			fmt.Printf("Listing messages in network '%s'\n", networkName)
			
			// Get messages
			messages, err := messagingSystem.GetMessages(networkName, limit, offset)
			if err != nil {
				return fmt.Errorf("failed to list messages: %w", err)
			}
			
			if len(messages) == 0 {
				fmt.Println("No messages found")
				return nil
			}
			
			// Print header
			fmt.Printf("%-12s %-12s %-10s %-19s %s\n", "ID", "SENDER", "FORMAT", "TIMESTAMP", "CONTENT")
			fmt.Printf("%-12s %-12s %-10s %-19s %s\n", "--", "------", "------", "---------", "-------")
			
			// Print messages
			for _, msg := range messages {
				idDisplay := msg.ID
				if len(idDisplay) > 8 {
					idDisplay = idDisplay[:8]
				}
				
				content := msg.Content
				if len(content) > 50 {
					content = content[:47] + "..."
				}
				
				fmt.Printf("%-12s %-12s %-10s %-19s %s\n", 
					idDisplay,
					msg.SenderID,
					msg.Format,
					msg.Timestamp.Format("2006-01-02 15:04:05"),
					content)
			}
			
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 20, "Maximum number of messages to show")
	cmd.Flags().IntVar(&offset, "offset", 0, "Number of messages to skip")
	
	return cmd
}

// newMessageGetCmd creates the message get command
func newMessageGetCmd(networkManager *NetworkManager, messagingSystem *MessagingSystem) *cobra.Command {
	return &cobra.Command{
		Use:   "get [network_name] [message_id]",
		Short: "Get a specific message",
		Long:  `Display detailed information about a message`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			messageID := args[1]
			
			// Get network
			network, err := networkManager.GetNetworkByName(networkName)
			if err != nil {
				return fmt.Errorf("failed to find network: %w", err)
			}
			
			// Get message
			message, err := messagingSystem.GetMessageByID(network.ID, messageID)
			if err != nil {
				return fmt.Errorf("failed to get message: %w", err)
			}
			
			// Display message details
			fmt.Printf("Message ID: %s\n", message.ID)
			fmt.Printf("  Network: %s\n", networkName)
			fmt.Printf("  Sender: %s\n", message.SenderID)
			fmt.Printf("  Format: %s\n", message.Format)
			fmt.Printf("  Timestamp: %s\n", message.Timestamp.Format("2006-01-02 15:04:05"))
			
			if message.Metadata != nil && len(message.Metadata) > 0 {
				fmt.Println("  Metadata:")
				metadataData, _ := json.MarshalIndent(message.Metadata, "    ", "  ")
				fmt.Println(string(metadataData))
			}
			
			fmt.Printf("  Content: %s\n", message.Content)
			
			if len(message.Attachments) > 0 {
				fmt.Printf("  Attachments (%d):\n", len(message.Attachments))
				for i, attachment := range message.Attachments {
					fmt.Printf("    %d. %s (%s, %d bytes)\n", i+1, attachment.Name, attachment.Format, attachment.Size)
					if attachment.Metadata != nil && len(attachment.Metadata) > 0 {
						metadataData, _ := json.MarshalIndent(attachment.Metadata, "        ", "  ")
						fmt.Println("       Metadata:", string(metadataData))
					}
				}
			}
			
			return nil
		},
	}
}

// newNetworkConfigCmd creates the network config command
func newNetworkConfigCmd(networkManager *NetworkManager) *cobra.Command {
	var configStr string

	cmd := &cobra.Command{
		Use:   "config [network_name]",
		Short: "Configure a network",
		Long:  `Update configuration for an existing network`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkName := args[0]
			
			// Parse config
			config := make(map[string]interface{})
			if err := json.Unmarshal([]byte(configStr), &config); err != nil {
				return fmt.Errorf("invalid config JSON: %w", err)
			}
			
			fmt.Printf("Updating configuration for network '%s'\n", networkName)
			
			// Update network
			if err := networkManager.UpdateNetwork(networkName, config); err != nil {
				return fmt.Errorf("failed to update network configuration: %w", err)
			}
			
			fmt.Printf("Network '%s' configuration updated successfully\n", networkName)
			return nil
		},
	}

	cmd.Flags().StringVar(&configStr, "config", "{}", "Network configuration as JSON string")
	cmd.MarkFlagRequired("config")
	
	return cmd
}
