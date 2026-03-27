package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"scripts/internal/ai/spec"
)

func (c *Client) Generate(ctx context.Context, req spec.Request) (spec.Response, error) {
	body, err := json.Marshal(map[string]any{
		"model": req.Model,
		"messages": []map[string]string{
			{"role": "user", "content": req.Prompt},
		},
		"stream": false,
		"think":  req.Think,
	})
	if err != nil {
		return spec.Response{}, err
	}

	resp, err := c.sendRequest(ctx, "POST", "/api/chat", bytes.NewReader(body))
	if err != nil {
		return spec.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return spec.Response{}, fmt.Errorf("ollama: unexpected status %s", resp.Status)
	}

	var raw struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return spec.Response{}, err
	}

	return spec.Response{Text: normalize(raw.Message.Content)}, nil
}
