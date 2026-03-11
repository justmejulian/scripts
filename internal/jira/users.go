package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type User struct {
	AccountID    string
	DisplayName  string
	EmailAddress string
}

func (c Client) GetCurrentUser(ctx context.Context) (User, error) {
	resp, err := c.sendRequest(ctx, "GET", c.baseURL+"/myself", nil)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var raw struct {
		AccountID    string `json:"accountId"`
		DisplayName  string `json:"displayName"`
		EmailAddress string `json:"emailAddress"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return User{}, err
	}

	return User{
		AccountID:    raw.AccountID,
		DisplayName:  raw.DisplayName,
		EmailAddress: raw.EmailAddress,
	}, nil
}

func (c Client) FindUser(ctx context.Context, query string) ([]User, error) {
	resp, err := c.sendRequest(ctx, "GET", c.baseURL+"/user/search?query="+url.QueryEscape(query), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var raw []struct {
		AccountID    string `json:"accountId"`
		DisplayName  string `json:"displayName"`
		EmailAddress string `json:"emailAddress"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	users := make([]User, len(raw))
	for i, r := range raw {
		users[i] = User{
			AccountID:    r.AccountID,
			DisplayName:  r.DisplayName,
			EmailAddress: r.EmailAddress,
		}
	}
	return users, nil
}
