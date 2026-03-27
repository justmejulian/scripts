package models

import (
	"scripts/internal/ai/providers/ollama"
	"scripts/internal/ai/spec"
)

type Qwen3_5_4BConfig struct {
	Default       string
	ThinkEnabled  string
	ThinkDisabled string
}

var Qwen3_5_4B = spec.Model[Qwen3_5_4BConfig]{
	Name:     "qwen3.5:4b",
	Provider: ollama.Name,
	Config: Qwen3_5_4BConfig{
		Default:       `{"think":true}`,
		ThinkEnabled:  `{"think":true}`,
		ThinkDisabled: `{"think":false}`,
	},
}
