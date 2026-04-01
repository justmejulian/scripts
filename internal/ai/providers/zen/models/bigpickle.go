package models

import (
	"scripts/internal/ai/providers/zen"
	"scripts/internal/ai/spec/model"
)

type BigPickleConfig struct {
	Default string
}

var BigPickle = model.Model[BigPickleConfig]{
	Info: model.Info{
		Name:     "big-pickle",
		Provider: zen.Name,
		Endpoint: zen.EndpointChat,
	},
	Config: BigPickleConfig{
		Default: "",
	},
}
