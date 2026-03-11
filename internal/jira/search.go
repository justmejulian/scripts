package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func (c Client) SearchIssues(ctx context.Context, jql string, fields []string) ([]Issue, error) {
	body, err := json.Marshal(map[string]any{
		"jql":    jql,
		"fields": fields,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.sendRequest(ctx, "POST", c.baseURL+"/search", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var result struct {
		Issues []struct {
			Key    string `json:"key"`
			Fields struct {
				Summary string `json:"summary"`
				Status  struct {
					Name string `json:"name"`
				} `json:"status"`
			} `json:"fields"`
		} `json:"issues"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	issues := make([]Issue, len(result.Issues))
	for i, r := range result.Issues {
		issues[i] = Issue{
			Key:    r.Key,
			Title:  r.Fields.Summary,
			Status: r.Fields.Status.Name,
		}
	}
	return issues, nil
}
