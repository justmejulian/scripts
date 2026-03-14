package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func (c *Client) AddComment(ctx context.Context, key, body string) error {
	payload, err := json.Marshal(map[string]any{"body": body})
	if err != nil {
		return err
	}

	resp, err := c.sendRequest(ctx, "POST", c.baseURL+"/issue/"+key+"/comment", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}
	return nil
}
