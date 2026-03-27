package models

import (
	"scripts/internal/ai/providers/zen"
	"scripts/internal/ai/spec"
)

type MinimaxM25FreeConfig struct {
	Default string
}

var MinimaxM25Free = spec.Model[MinimaxM25FreeConfig]{
	Name:     "minimax-m2.5-free",
	Provider: zen.Name,
	Config: MinimaxM25FreeConfig{
		Default: "",
	},
}
