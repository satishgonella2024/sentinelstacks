package logs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/runtime"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelDebug LogLevel = "DEBUG"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Time    time.Time
	Level   LogLevel
	Message string
}

// NewLogsCmd creates a new logs command
func NewLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [agent_id]",
		Short: "View logs from an agent",
		Long:  `View logs from a running or stopped agent.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]
			follow, _ := cmd.Flags().GetBool("follow")
			tail, _ := cmd.Flags().GetInt("tail")
			timestamps, _ := cmd.Flags().GetBool("timestamps")
			level, _ := cmd.Flags().GetString("level")

			return runLogs(agentID, follow, tail, timestamps, level)
		},
	}

	// Add flags
	cmd.Flags().BoolP("follow", "f", false, "Follow log output")
	cmd.Flags().IntP("tail", "n", 10, "Number of lines to show from the end of the logs")
	cmd.Flags().BoolP("timestamps", "t", false, "Show timestamps")
	cmd.Flags().String("level", "info", "Minimum log level to display (debug, info, warn, error)")

	return cmd
}

// runLogs executes the logs command
func runLogs(agentID string, follow bool, tail int, timestamps bool, level string) error {
	// Get the runtime
	rt, err := runtime.GetRuntime()
	if err != nil {
		return fmt.Errorf("failed to get runtime: %w", err)
	}

	// Get agent information first
	agent, err := rt.GetAgent(agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %s", err)
	}

	// Set minimum log level
	minLevel := parseLogLevel(level)

	// Create colorized output
	infoColor := color.New(color.FgGreen).SprintFunc()
	warnColor := color.New(color.FgYellow).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()
	debugColor := color.New(color.FgCyan).SprintFunc()

	// Find the agent log file path
	// In a real implementation, this would come from the runtime
	// For now, we'll construct a path based on the agent ID

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %w", err)
	}

	// Construct log file path
	sentinelDir := filepath.Join(homeDir, ".sentinel")
	logDir := filepath.Join(sentinelDir, "agents", agentID, "logs")
	logFile := filepath.Join(logDir, "agent.log")

	// Check if log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		// For demo, create a dummy log file with sample entries
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("could not create log directory: %w", err)
		}
		if err := createSampleLogFile(logFile, agent.Name); err != nil {
			return fmt.Errorf("could not create sample log file: %w", err)
		}
	}

	// Open the log file
	file, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("could not open log file: %w", err)
	}
	defer file.Close()

	// Print header
	fmt.Printf("Logs for agent %s (%s)\n", agent.Name, agentID)
	fmt.Println(strings.Repeat("-", 80))

	// Read and process log entries
	scanner := bufio.NewScanner(file)
	var entries []LogEntry

	// Read all entries
	for scanner.Scan() {
		line := scanner.Text()
		entry := parseLine(line)
		if entry != nil && isLevelIncluded(entry.Level, minLevel) {
			entries = append(entries, *entry)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %w", err)
	}

	// If tail is specified, only show the last N entries
	if tail > 0 && tail < len(entries) {
		entries = entries[len(entries)-tail:]
	}

	// Print entries
	for _, entry := range entries {
		var levelStr string
		switch entry.Level {
		case LogLevelInfo:
			levelStr = infoColor(string(entry.Level))
		case LogLevelWarn:
			levelStr = warnColor(string(entry.Level))
		case LogLevelError:
			levelStr = errorColor(string(entry.Level))
		case LogLevelDebug:
			levelStr = debugColor(string(entry.Level))
		default:
			levelStr = string(entry.Level)
		}

		if timestamps {
			fmt.Printf("%s [%s] %s\n", entry.Time.Format(time.RFC3339), levelStr, entry.Message)
		} else {
			fmt.Printf("[%s] %s\n", levelStr, entry.Message)
		}
	}

	// If follow flag is set, continue watching for new log entries
	if follow {
		fmt.Println("\nWatching for new log entries... (press Ctrl+C to stop)")

		// Seek to end of file
		if _, err := file.Seek(0, io.SeekEnd); err != nil {
			return fmt.Errorf("could not seek to end of file: %w", err)
		}

		// Create a new scanner
		scanner = bufio.NewScanner(file)

		// Watch for new entries
		for {
			if scanner.Scan() {
				line := scanner.Text()
				entry := parseLine(line)
				if entry != nil && isLevelIncluded(entry.Level, minLevel) {
					var levelStr string
					switch entry.Level {
					case LogLevelInfo:
						levelStr = infoColor(string(entry.Level))
					case LogLevelWarn:
						levelStr = warnColor(string(entry.Level))
					case LogLevelError:
						levelStr = errorColor(string(entry.Level))
					case LogLevelDebug:
						levelStr = debugColor(string(entry.Level))
					default:
						levelStr = string(entry.Level)
					}

					if timestamps {
						fmt.Printf("%s [%s] %s\n", entry.Time.Format(time.RFC3339), levelStr, entry.Message)
					} else {
						fmt.Printf("[%s] %s\n", levelStr, entry.Message)
					}
				}
			} else {
				// Check for scanner errors
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("error reading log file: %w", err)
				}

				// Sleep briefly before trying again
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	return nil
}

// parseLine parses a log line into a LogEntry
func parseLine(line string) *LogEntry {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return nil
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil
	}

	// Parse level
	level := LogLevel(strings.Trim(parts[1], "[]"))

	// Get message
	message := parts[2]

	return &LogEntry{
		Time:    timestamp,
		Level:   level,
		Message: message,
	}
}

// parseLogLevel converts a string level to LogLevel
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		return LogLevelInfo
	}
}

// isLevelIncluded checks if a log level should be included based on the minimum level
func isLevelIncluded(level, minLevel LogLevel) bool {
	levelOrder := map[LogLevel]int{
		LogLevelDebug: 0,
		LogLevelInfo:  1,
		LogLevelWarn:  2,
		LogLevelError: 3,
	}

	return levelOrder[level] >= levelOrder[minLevel]
}

// createSampleLogFile creates a sample log file for demonstration purposes
func createSampleLogFile(path, agentName string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create sample log entries
	now := time.Now()
	entries := []LogEntry{
		{Time: now.Add(-10 * time.Minute), Level: LogLevelInfo, Message: fmt.Sprintf("Agent %s initialization started", agentName)},
		{Time: now.Add(-9 * time.Minute), Level: LogLevelDebug, Message: "Loading configuration"},
		{Time: now.Add(-9 * time.Minute), Level: LogLevelInfo, Message: "Configuration loaded successfully"},
		{Time: now.Add(-8 * time.Minute), Level: LogLevelDebug, Message: "Connecting to LLM provider"},
		{Time: now.Add(-8 * time.Minute), Level: LogLevelWarn, Message: "Connection latency higher than expected"},
		{Time: now.Add(-7 * time.Minute), Level: LogLevelInfo, Message: "Connected to LLM provider"},
		{Time: now.Add(-6 * time.Minute), Level: LogLevelDebug, Message: "Initializing agent state"},
		{Time: now.Add(-5 * time.Minute), Level: LogLevelInfo, Message: "Agent state initialized"},
		{Time: now.Add(-4 * time.Minute), Level: LogLevelError, Message: "Failed to load external tool: web_search"},
		{Time: now.Add(-3 * time.Minute), Level: LogLevelInfo, Message: "Retrying tool initialization"},
		{Time: now.Add(-2 * time.Minute), Level: LogLevelInfo, Message: "External tool loaded successfully"},
		{Time: now.Add(-1 * time.Minute), Level: LogLevelInfo, Message: fmt.Sprintf("Agent %s ready", agentName)},
	}

	// Write entries to file
	for _, entry := range entries {
		line := fmt.Sprintf("%s [%s] %s\n", entry.Time.Format(time.RFC3339), entry.Level, entry.Message)
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}
