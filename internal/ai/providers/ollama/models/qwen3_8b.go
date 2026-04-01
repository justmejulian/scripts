package models

import (
	"scripts/internal/ai/providers/ollama"
	"scripts/internal/ai/spec/model"
)

type Qwen3_8BConfig struct {
	Default       string
	ThinkEnabled  string
	ThinkDisabled string
}

var Qwen3_8B = model.Model[Qwen3_8BConfig]{
	Info: model.Info{
		Name:     "qwen3:8b",
		Provider: ollama.Name,
	},
	Config: Qwen3_8BConfig{
		Default:       `{"think":true}`,
		ThinkEnabled:  `{"think":true}`,
		ThinkDisabled: `{"think":false}`,
	},
}
