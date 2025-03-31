package sentinel

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	versionCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/version"
)

// TestVersionCommand tests the version command separately with direct output capturing
func TestVersionCommand(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Save the original outputWriter
	origWriter := versionCmd.GetOutputWriter()

	// Set outputWriter to our buffer
	versionCmd.SetOutputWriter(&buf)

	// Restore outputWriter after test
	defer versionCmd.SetOutputWriter(origWriter)

	// Create a new version command
	cmd := versionCmd.NewVersionCmd()

	// Execute command with default options
	err := cmd.Execute()
	require.NoError(t, err)

	// Check output
	output := buf.String()
	assert.Contains(t, output, "SentinelStacks AI Agent Management System")
	assert.Contains(t, output, "Version:")
	assert.Contains(t, output, versionCmd.Version)

	// Reset buffer and test short version
	buf.Reset()
	cmd.SetArgs([]string{"--short"})
	err = cmd.Execute()
	require.NoError(t, err)

	// Check output
	output = buf.String()
	assert.Equal(t, versionCmd.Version+"\n", output)
}

// TestImagesCommand tests the images command separately with direct output capturing
func TestImagesCommand(t *testing.T) {
	// Create a new root command
	rootCmd := newTestRootCmd()

	// Test images command
	cmd, _, err := rootCmd.Find([]string{"images"})
	require.NotNil(t, cmd, "Images command should exist")
	require.NoError(t, err, "Finding images command should not error")

	// Execute command with quiet flag
	cmd.SetArgs([]string{"--quiet"})
	err = cmd.Execute()
	assert.NoError(t, err, "Images command should execute without error in quiet mode")

	// Execute command with default settings
	cmd.SetArgs([]string{})
	err = cmd.Execute()
	assert.NoError(t, err, "Images command should execute without error with default settings")

	// Execute command with format option
	cmd.SetArgs([]string{"--format", "{{.ID}}"})
	err = cmd.Execute()
	assert.NoError(t, err, "Images command should execute without error with format option")

	// Execute command with filter option
	cmd.SetArgs([]string{"--filter", "name=translator"})
	err = cmd.Execute()
	assert.NoError(t, err, "Images command should execute without error with filter option")
}

// Integration tests for the CLI commands
func TestCLIHelpCommand(t *testing.T) {
	// Create a new test command
	cmd := newTestRootCmd()

	// Create buffers for stdout and stderr
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Set args for help command
	cmd.SetArgs([]string{"--help"})

	// Execute command
	err := cmd.Execute()
	require.NoError(t, err)

	// Get output
	output := stdout.String() + stderr.String()

	// Check expected content
	assert.Contains(t, output, "SentinelStacks")
	assert.Contains(t, output, "Available Commands:")
	assert.Contains(t, output, "build")
	assert.Contains(t, output, "run")
	assert.Contains(t, output, "ps")
	assert.Contains(t, output, "images")
}

// Create a test version of rootCmd for integration testing
func newTestRootCmd() *cobra.Command {
	// Create a copy of the rootCmd
	testCmd := &cobra.Command{
		Use:   rootCmd.Use,
		Short: rootCmd.Short,
		Long:  rootCmd.Long,
	}

	// Copy PersistentFlags
	testCmd.PersistentFlags().AddFlagSet(rootCmd.PersistentFlags())

	// Add commands
	for _, cmd := range rootCmd.Commands() {
		testCmd.AddCommand(cmd)
	}

	return testCmd
}

// This is a more advanced integration test that executes multiple commands in sequence
func TestWorkflow(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping workflow test in short mode")
	}

	// This would simulate a full workflow:
	// 1. Initialize an agent
	// 2. Build the agent
	// 3. Run the agent
	// 4. List running agents
	// 5. Stop the agent

	// For simplicity, this test is not fully implemented
	// as it would require a more complex test environment
	t.Skip("workflow test not fully implemented")
}

func TestMultimodalCommand(t *testing.T) {
	rootCmd := newTestRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Test with the --help flag
	rootCmd.SetArgs([]string{"multimodal", "--help"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
	output := buf.String()

	// Check the output
	assert.Contains(t, output, "analyze-image")
	assert.Contains(t, output, "multimodal")
}

func TestChatCommand(t *testing.T) {
	rootCmd := newTestRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Test with the --help flag
	rootCmd.SetArgs([]string{"chat", "--help"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
	output := buf.String()

	// Check the output
	assert.Contains(t, output, "Start an interactive chat session")
	assert.Contains(t, output, "--provider")
	assert.Contains(t, output, "--model")
	assert.Contains(t, output, "--images")
}
