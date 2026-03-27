package ai

import (
	"fmt"

	ollama "scripts/internal/ai/providers/ollama"
)

type Config struct {
	Provider string
}

func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case "ollama":
		return ollama.New()
	default:
		return nil, fmt.Errorf("ai: unsupported provider %q", cfg.Provider)
	}
}
