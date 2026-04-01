package models

import (
	"scripts/internal/ai/providers/zen"
	"scripts/internal/ai/spec/model"
)

type MinimaxM25FreeConfig struct {
	Default string
}

var MinimaxM25Free = model.Model[MinimaxM25FreeConfig]{
	Info: model.Info{
		Name:     "minimax-m2.5-free",
		Provider: zen.Name,
		Endpoint: zen.EndpointChat,
	},
	Config: MinimaxM25FreeConfig{
		Default: "",
	},
}
