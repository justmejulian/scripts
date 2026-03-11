package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCurrentUser_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/rest/api/2/myself" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"accountId":    "abc123",
			"displayName":  "Alice",
			"emailAddress": "alice@example.com",
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	user, err := c.GetCurrentUser(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.AccountID != "abc123" {
		t.Errorf("unexpected account ID: %s", user.AccountID)
	}
	if user.DisplayName != "Alice" {
		t.Errorf("unexpected display name: %s", user.DisplayName)
	}
	if user.EmailAddress != "alice@example.com" {
		t.Errorf("unexpected email: %s", user.EmailAddress)
	}
}

func TestGetCurrentUser_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetCurrentUser(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestGetCurrentUser_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetCurrentUser(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFindUser_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/rest/api/2/user/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "alice" {
			t.Errorf("unexpected query param: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"accountId": "abc123", "displayName": "Alice", "emailAddress": "alice@example.com"},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	users, err := c.FindUser(context.Background(), "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].AccountID != "abc123" || users[0].DisplayName != "Alice" || users[0].EmailAddress != "alice@example.com" {
		t.Errorf("unexpected user: %+v", users[0])
	}
}

func TestFindUser_QueryEncoding(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("query")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.FindUser(context.Background(), "alice smith")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery != "alice smith" {
		t.Errorf("unexpected decoded query: %s", gotQuery)
	}
}

func TestFindUser_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.FindUser(context.Background(), "alice")
	if err == nil {
		t.Fatal("expected error for 403 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestFindUser_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.FindUser(context.Background(), "alice")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFindUser_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	users, err := c.FindUser(context.Background(), "nobody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}
