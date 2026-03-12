package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Issue struct {
	Key    string
	Title  string
	Status string
}

func (c *Client) GetIssue(ctx context.Context, key string) (Issue, error) {
	resp, err := c.sendRequest(ctx, "GET", c.baseURL+"/issue/"+key, nil)
	if err != nil {
		return Issue{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Issue{}, &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var raw struct {
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
			Status  struct {
				Name string `json:"name"`
			} `json:"status"`
		} `json:"fields"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return Issue{}, err
	}

	return Issue{
		Key:    raw.Key,
		Title:  raw.Fields.Summary,
		Status: raw.Fields.Status.Name,
	}, nil
}

func (c *Client) CreateIssue(ctx context.Context, project, summary, issueType string) (Issue, error) {
	body, err := json.Marshal(map[string]any{
		"fields": map[string]any{
			"project":   map[string]string{"key": project},
			"summary":   summary,
			"issuetype": map[string]string{"name": issueType},
		},
	})
	if err != nil {
		return Issue{}, err
	}

	resp, err := c.sendRequest(ctx, "POST", c.baseURL+"/issue", bytes.NewReader(body))
	if err != nil {
		return Issue{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return Issue{}, &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var raw struct {
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
			Status  struct {
				Name string `json:"name"`
			} `json:"status"`
		} `json:"fields"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return Issue{}, err
	}

	return Issue{
		Key:    raw.Key,
		Title:  raw.Fields.Summary,
		Status: raw.Fields.Status.Name,
	}, nil
}

func (c *Client) UpdateIssue(ctx context.Context, key string, fields map[string]any) error {
	body, err := json.Marshal(map[string]any{"fields": fields})
	if err != nil {
		return err
	}

	resp, err := c.sendRequest(ctx, "PUT", c.baseURL+"/issue/"+key, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}
	return nil
}
