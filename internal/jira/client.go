// Package jira provides a client for the Jira REST API v2.
package jira

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	baseURL    string
	authHeader string
	http       *http.Client
}

func NewClientFromEnv() (Client, error) {
	domain := os.Getenv("JIRA_DOMAIN")
	token := os.Getenv("JIRA_TOKEN")

	if domain == "" || token == "" {
		return Client{}, fmt.Errorf("missing env vars: JIRA_DOMAIN and JIRA_TOKEN are required")
	}

	return NewClientPAT(domain, token), nil
}

func NewClientPAT(domain, token string) Client {
	return Client{
		baseURL:    "https://" + domain + "/rest/api/2",
		authHeader: "Bearer " + token,
		http:       &http.Client{},
	}
}

func (c Client) sendRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return c.http.Do(req)
}
