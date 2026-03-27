package zen

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"scripts/internal/ai/registry"
	"scripts/internal/ai/spec"
)

const Name = "zen"

func init() {
	registry.Register(Name, New)
}

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func New() (spec.Provider, error) {
	key := os.Getenv("ZEN_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("zen: ZEN_API_KEY is required")
	}

	return &Client{
		baseURL: "https://opencode.ai",
		apiKey:  key,
		http:    &http.Client{},
	}, nil
}

func (c *Client) sendRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	return c.http.Do(req)
}
