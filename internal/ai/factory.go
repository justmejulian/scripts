package ai

import (
	"fmt"

	ollama "scripts/internal/ai/providers/ollama"
	zen "scripts/internal/ai/providers/zen"
)

type Config struct {
	Provider string
}

func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case "ollama":
		return ollama.New()
	case "zen":
		return zen.New()
	default:
		return nil, fmt.Errorf("ai: unsupported provider %q", cfg.Provider)
	}
}
