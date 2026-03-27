package ollama

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"scripts/internal/ai/registry"
	"scripts/internal/ai/spec"
)

const Name = "ollama"

func init() {
	registry.Register(Name, New)
}

type Client struct {
	baseURL string
	http    *http.Client
}

func New() (spec.Provider, error) {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		return nil, fmt.Errorf("ollama: OLLAMA_HOST is required")
	}

	return &Client{
		baseURL: host,
		http:    &http.Client{},
	}, nil
}

func (c *Client) sendRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.http.Do(req)
}
