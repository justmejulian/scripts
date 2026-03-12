// Package ollama provides a client for a locally-running Ollama instance.
package ollama

import (
	"context"
	"io"
	"net/http"
	"os"
)

type Client struct {
	baseURL string
	model   string
	http    *http.Client
}

func NewClient(model string) *Client {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}

	return &Client{
		baseURL: host,
		model:   model,
		http:    &http.Client{},
	}
}

func (c *Client) sendRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.http.Do(req)
}
