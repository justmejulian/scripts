package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreatePRRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	SourceRefName string `json:"sourceRefName"`
	TargetRefName string `json:"targetRefName"`
}

type PullRequest struct {
	PullRequestID int    `json:"pullRequestId"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	SourceRefName string `json:"sourceRefName"`
	TargetRefName string `json:"targetRefName"`
	URL           string `json:"url"`
}

func (c *Client) CreatePR(ctx context.Context, project, repo string, req CreatePRRequest) (*PullRequest, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("azure: marshal request: %w", err)
	}

	url := c.url(project, fmt.Sprintf("git/repositories/%s/pullrequests", repo))
	resp, err := c.sendRequest(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("azure: decode response: %w", err)
	}

	return &pr, nil
}
