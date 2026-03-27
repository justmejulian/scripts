package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"scripts/internal/ai/spec"
	"scripts/internal/ai/utils/requestconfig"
)

func (c *Client) Generate(ctx context.Context, req spec.Request) (spec.Response, error) {
	bodyMap := map[string]any{
		"model": req.Model,
		"messages": []map[string]string{
			{"role": "user", "content": req.Prompt},
		},
		"stream": false,
	}

	bodyMap, err := requestconfig.Apply("ollama", bodyMap, req.Config, "model", "messages", "stream")
	if err != nil {
		return spec.Response{}, err
	}

	body, err := json.Marshal(bodyMap)
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
