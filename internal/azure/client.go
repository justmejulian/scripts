// Package azure provides a client for the Azure DevOps REST API v7.2-preview.
package azure

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiVersion = "7.2-preview"

type Client struct {
	baseURL    string
	authHeader string
	http       *http.Client
}

func NewClientFromEnv() (*Client, error) {
	org := os.Getenv("AZURE_DEVOPS_ORG")
	token := os.Getenv("AZURE_DEVOPS_TOKEN")

	if org == "" || token == "" {
		return nil, fmt.Errorf("missing env vars: AZURE_DEVOPS_ORG and AZURE_DEVOPS_TOKEN are required")
	}

	return NewClientPAT(org, token), nil
}

func NewClientPAT(org, token string) *Client {
	encoded := base64.StdEncoding.EncodeToString([]byte(":" + token))
	return &Client{
		baseURL:    "https://dev.azure.com/" + org,
		authHeader: "Basic " + encoded,
		http:       &http.Client{},
	}
}

func (c *Client) url(project, path string) string {
	return fmt.Sprintf("%s/%s/_apis/%s?api-version=%s", c.baseURL, project, path, apiVersion)
}

func (c *Client) sendRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	return c.sendRequestWithContentType(ctx, method, url, "application/json", body)
}

func (c *Client) sendRequestWithContentType(ctx context.Context, method, url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", contentType)

	return c.http.Do(req)
}
