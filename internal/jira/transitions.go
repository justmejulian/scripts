package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) TransitionIssue(ctx context.Context, key, transitionName string) error {
	resp, err := c.sendRequest(ctx, "GET", c.baseURL+"/issue/"+key+"/transitions", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var raw struct {
		Transitions []transition `json:"transitions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return err
	}

	var target *transition
	available := make([]string, 0, len(raw.Transitions))
	for i, t := range raw.Transitions {
		available = append(available, t.Name)
		if strings.EqualFold(t.Name, transitionName) {
			target = &raw.Transitions[i]
		}
	}

	if target == nil {
		return fmt.Errorf("jira: transition %q not found for %s; available: %s",
			transitionName, key, strings.Join(available, ", "))
	}

	body, err := json.Marshal(map[string]any{
		"transition": map[string]string{"id": target.ID},
	})
	if err != nil {
		return err
	}

	resp2, err := c.sendRequest(ctx, "POST", c.baseURL+"/issue/"+key+"/transitions", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusNoContent {
		return &APIError{StatusCode: resp2.StatusCode, Status: resp2.Status}
	}
	return nil
}
