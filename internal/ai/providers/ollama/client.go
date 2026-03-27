package ollama

import (
	"context"
	"io"
	"net/http"
	"os"

	"scripts/internal/ai/spec"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New() spec.Provider {
	host := os.Getenv("OLLAMA_HOST")
	// todo fail if not found
	if host == "" {
		host = "http://localhost:11434"
	}

	return &Client{
		baseURL: host,
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
