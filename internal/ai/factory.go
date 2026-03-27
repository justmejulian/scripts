package ai

import "scripts/internal/ai/registry"

type Config struct {
	Provider string
}

func NewProvider(cfg Config) (Provider, error) {
	return registry.New(cfg.Provider)
}
