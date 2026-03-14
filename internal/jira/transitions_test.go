package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTransitionIssue_Success(t *testing.T) {
	var postCalled bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/rest/api/2/issue/PROJ-1/transitions" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"transitions": []map[string]any{
					{"id": "10", "name": "In Review"},
					{"id": "20", "name": "Done"},
				},
			})
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == "/rest/api/2/issue/PROJ-1/transitions" {
			postCalled = true
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode body: %v", err)
			}
			trans := body["transition"].(map[string]any)
			if trans["id"] != "10" {
				t.Errorf("unexpected transition id: %v", trans["id"])
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.TransitionIssue(context.Background(), "PROJ-1", "In Review")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !postCalled {
		t.Error("expected POST to transitions endpoint")
	}
}

func TestTransitionIssue_CaseInsensitive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]any{
				"transitions": []map[string]any{
					{"id": "10", "name": "In Review"},
				},
			})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.TransitionIssue(context.Background(), "PROJ-1", "in review")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransitionIssue_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"transitions": []map[string]any{
				{"id": "20", "name": "Done"},
				{"id": "30", "name": "To Do"},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.TransitionIssue(context.Background(), "PROJ-1", "In Review")
	if err == nil {
		t.Fatal("expected error for missing transition")
	}
	if !strings.Contains(err.Error(), "In Review") {
		t.Errorf("error should mention missing transition name, got: %v", err)
	}
	if !strings.Contains(err.Error(), "Done") {
		t.Errorf("error should list available transitions, got: %v", err)
	}
}

func TestTransitionIssue_GetAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.TransitionIssue(context.Background(), "PROJ-1", "In Review")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestTransitionIssue_PostAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]any{
				"transitions": []map[string]any{
					{"id": "10", "name": "In Review"},
				},
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.TransitionIssue(context.Background(), "PROJ-1", "In Review")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}
