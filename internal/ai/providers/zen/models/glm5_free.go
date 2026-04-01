package models

import (
	"scripts/internal/ai/providers/zen"
	"scripts/internal/ai/spec/model"
)

type GLM5FreeConfig struct {
	Default          string
	ThinkingEnabled  string
	ThinkingDisabled string
}

var GLM5Free = model.Model[GLM5FreeConfig]{
	Info: model.Info{
		Name:     "glm-5",
		Provider: zen.Name,
		Endpoint: zen.EndpointChat,
	},
	Config: GLM5FreeConfig{
		Default:          `{"thinking":{"type":"enabled"}}`,
		ThinkingEnabled:  `{"thinking":{"type":"enabled"}}`,
		ThinkingDisabled: `{"thinking":{"type":"disabled"}}`,
	},
}
