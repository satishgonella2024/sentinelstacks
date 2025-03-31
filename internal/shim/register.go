package shim

import (
	"github.com/sentinelstacks/sentinel/internal/shim/claude"
	"github.com/sentinelstacks/sentinel/internal/shim/mock"
	"github.com/sentinelstacks/sentinel/internal/shim/ollama"
	"github.com/sentinelstacks/sentinel/internal/shim/openai"
)

// init registers known providers
func init() {
	// Register Claude provider
	RegisterProviderFactory("claude", func() Provider {
		return claude.NewProvider().(Provider)
	})

	// Register OpenAI provider
	RegisterProviderFactory("openai", func() Provider {
		return openai.NewProvider().(Provider)
	})

	// Register Ollama provider
	RegisterProviderFactory("ollama", func() Provider {
		return ollama.NewProvider().(Provider)
	})

	// Register Mock provider for testing
	RegisterProviderFactory("mock", func() Provider {
		return mock.NewMockProvider().(Provider)
	})
}
