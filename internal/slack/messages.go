package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func (c *Client) PostMessage(ctx context.Context, channel, text string) error {
	body, err := json.Marshal(map[string]string{
		"channel": channel,
		"text":    text,
	})
	if err != nil {
		return err
	}

	resp, err := c.sendRequest(ctx, http.MethodPost, baseURL+"/chat.postMessage", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var raw struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return err
	}
	if !raw.OK {
		return &APIError{Code: raw.Error}
	}
	return nil
}
