package models

import (
	"scripts/internal/ai/providers/zen"
	"scripts/internal/ai/spec/model"
)

type GPT5NanoConfig struct {
	Default string
}

var GPT5Nano = model.Model[GPT5NanoConfig]{
	Info: model.Info{
		Name:     "gpt-5-nano",
		Provider: zen.Name,
		Endpoint: zen.EndpointResponses,
	},
	Config: GPT5NanoConfig{
		Default: "",
	},
}
