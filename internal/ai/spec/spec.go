package spec

import "context"

type Provider interface {
	Generate(ctx context.Context, req Request) (Response, error)
}

type Request struct {
	Prompt string
	Model  string
	Think  bool
}

type Response struct {
	Text string
}
