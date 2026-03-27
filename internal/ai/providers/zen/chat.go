package zen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"scripts/internal/ai/spec"
)

func (c *Client) Generate(ctx context.Context, req spec.Request) (spec.Response, error) {
	body, err := json.Marshal(map[string]any{
		"model": req.Model,
		"messages": []map[string]string{
			{"role": "user", "content": req.Prompt},
		},
	})
	if err != nil {
		return spec.Response{}, err
	}

	resp, err := c.sendRequest(ctx, "POST", "/zen/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return spec.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return spec.Response{}, fmt.Errorf("zen: unexpected status %s", resp.Status)
	}

	var raw struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return spec.Response{}, err
	}

	if len(raw.Choices) == 0 {
		return spec.Response{}, fmt.Errorf("zen: no choices in response")
	}

	return spec.Response{Text: strings.TrimSpace(raw.Choices[0].Message.Content)}, nil
}
