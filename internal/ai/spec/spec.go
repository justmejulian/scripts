package spec

import (
	"context"

	"scripts/internal/ai/spec/model"
)

type Provider interface {
	Generate(ctx context.Context, req Request) (Response, error)
}

type Request struct {
	Prompt string
	Model  model.Info
	Config string
}

type Response struct {
	Text string
}
