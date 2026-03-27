package models

import (
	"scripts/internal/ai/providers/zen"
	"scripts/internal/ai/spec"
)

type BigPickleConfig struct {
	Default string
}

var BigPickle = spec.Model[BigPickleConfig]{
	Name:     "big-pickle",
	Provider: zen.Name,
	Config: BigPickleConfig{
		Default: "",
	},
}
