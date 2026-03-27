package models

import (
	"scripts/internal/ai/providers/zen"
	"scripts/internal/ai/spec"
)

type GLM5FreeConfig struct {
	Default          string
	ThinkingEnabled  string
	ThinkingDisabled string
}

var GLM5Free = spec.Model[GLM5FreeConfig]{
	Name:     "glm-5",
	Provider: zen.Name,
	Config: GLM5FreeConfig{
		Default:          `{"thinking":{"type":"enabled"}}`,
		ThinkingEnabled:  `{"thinking":{"type":"enabled"}}`,
		ThinkingDisabled: `{"thinking":{"type":"disabled"}}`,
	},
}
