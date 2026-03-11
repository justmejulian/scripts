package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchIssues_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/rest/api/2/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"issues": []map[string]any{
				{
					"key": "PROJ-1",
					"fields": map[string]any{
						"summary": "First issue",
						"status":  map[string]string{"name": "To Do"},
					},
				},
				{
					"key": "PROJ-2",
					"fields": map[string]any{
						"summary": "Second issue",
						"status":  map[string]string{"name": "Done"},
					},
				},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	issues, err := c.SearchIssues(context.Background(), "project = PROJ", []string{"summary", "status"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key != "PROJ-1" || issues[0].Title != "First issue" || issues[0].Status != "To Do" {
		t.Errorf("unexpected first issue: %+v", issues[0])
	}
	if issues[1].Key != "PROJ-2" || issues[1].Title != "Second issue" || issues[1].Status != "Done" {
		t.Errorf("unexpected second issue: %+v", issues[1])
	}
}

func TestSearchIssues_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"issues": []any{}})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	issues, err := c.SearchIssues(context.Background(), "project = EMPTY", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}

func TestSearchIssues_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.SearchIssues(context.Background(), "project = PROJ", nil)
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

func TestSearchIssues_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.SearchIssues(context.Background(), "project = PROJ", nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
