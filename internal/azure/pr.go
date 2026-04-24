package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

type IdentityRef struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type PRComment struct {
	ID              int         `json:"id"`
	ParentCommentID int         `json:"parentCommentId"`
	Author          IdentityRef `json:"author"`
	Content         string      `json:"content"`
	PublishedDate   time.Time   `json:"publishedDate"`
	CommentType     string      `json:"commentType"`
}

type PRThread struct {
	ID            int         `json:"id"`
	Comments      []PRComment `json:"comments"`
	Status        string      `json:"status"`
	PublishedDate time.Time   `json:"publishedDate"`
}

type prThreadsResponse struct {
	Value []PRThread `json:"value"`
}

func (c *Client) GetPRThreads(ctx context.Context, project, repo string, prID int) ([]PRThread, error) {
	url := c.urlPreview(project, fmt.Sprintf("git/repositories/%s/pullrequests/%d/threads", repo, prID))
	resp, err := c.sendRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	var result prThreadsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("azure: decode response: %w", err)
	}

	return result.Value, nil
}

func (c *Client) CreatePR(ctx context.Context, project, repo string, req CreatePRRequest) (*PullRequest, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("azure: marshal request: %w", err)
	}

	url := c.urlPreview(project, fmt.Sprintf("git/repositories/%s/pullrequests", repo))
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
