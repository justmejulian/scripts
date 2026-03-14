// Package slack provides a client for the Slack Web API.
package slack

import (
	"context"
	"io"
	"net/http"
	"os"

	"fmt"
)

const baseURL = "https://slack.com/api"

type Client struct {
	authHeader string
	http       *http.Client
}

func NewClientFromEnv() (*Client, error) {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("missing env var: SLACK_TOKEN is required")
	}
	return NewClient(token), nil
}

func NewClient(token string) *Client {
	return &Client{
		authHeader: "Bearer " + token,
		http:       &http.Client{},
	}
}

func (c *Client) sendRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return c.http.Do(req)
}
