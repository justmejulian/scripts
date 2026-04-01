package zen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"scripts/internal/ai/spec"
	"scripts/internal/ai/utils/requestconfig"
)

func (c *Client) generateResponses(ctx context.Context, req spec.Request) (spec.Response, error) {
	bodyMap := map[string]any{
		"model": req.Model.Name,
		"input": req.Prompt,
	}

	bodyMap, err := requestconfig.Apply("zen", bodyMap, req.Config, "model", "input")
	if err != nil {
		return spec.Response{}, err
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return spec.Response{}, err
	}

	resp, err := c.sendRequest(ctx, "POST", "/zen/v1/responses", bytes.NewReader(body))
	if err != nil {
		return spec.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return spec.Response{}, fmt.Errorf("zen: unexpected status %s", resp.Status)
	}

	var raw struct {
		Output []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return spec.Response{}, err
	}

	for _, out := range raw.Output {
		if out.Type != "message" {
			continue
		}
		for _, c := range out.Content {
			if c.Type == "output_text" {
				return spec.Response{Text: strings.TrimSpace(c.Text)}, nil
			}
		}
	}

	return spec.Response{}, fmt.Errorf("zen: no output_text in response")
}
