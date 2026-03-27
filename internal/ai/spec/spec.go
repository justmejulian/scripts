package spec

import "context"

type Provider interface {
	Generate(ctx context.Context, req Request) (Response, error)
}

type Model[C any] struct {
	Name     string
	Provider string
	Config   C
}

type Request struct {
	Prompt string
	Model  string
	Config string
}

type Response struct {
	Text string
}
