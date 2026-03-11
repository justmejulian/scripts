package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(srv *httptest.Server) Client {
	return Client{
		baseURL:    srv.URL + "/rest/api/2",
		authHeader: "Bearer testtoken",
		http:       &http.Client{},
	}
}

func TestGetIssue_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/rest/api/2/issue/PROJ-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"key": "PROJ-1",
			"fields": map[string]any{
				"summary": "Test issue",
				"status":  map[string]string{"name": "In Progress"},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	issue, err := c.GetIssue(context.Background(), "PROJ-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.Key != "PROJ-1" {
		t.Errorf("unexpected key: %s", issue.Key)
	}
	if issue.Title != "Test issue" {
		t.Errorf("unexpected title: %s", issue.Title)
	}
	if issue.Status != "In Progress" {
		t.Errorf("unexpected status: %s", issue.Status)
	}
}

func TestGetIssue_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetIssue(context.Background(), "PROJ-999")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestGetIssue_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetIssue(context.Background(), "PROJ-1")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCreateIssue_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/rest/api/2/issue" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"key": "PROJ-2",
			"fields": map[string]any{
				"summary": "New issue",
				"status":  map[string]string{"name": "To Do"},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	issue, err := c.CreateIssue(context.Background(), "PROJ", "New issue", "Task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.Key != "PROJ-2" {
		t.Errorf("unexpected key: %s", issue.Key)
	}
	if issue.Title != "New issue" {
		t.Errorf("unexpected title: %s", issue.Title)
	}
}

func TestCreateIssue_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.CreateIssue(context.Background(), "PROJ", "Bad issue", "Task")
	if err == nil {
		t.Fatal("expected error for 400 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestCreateIssue_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.CreateIssue(context.Background(), "PROJ", "New issue", "Task")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestUpdateIssue_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/rest/api/2/issue/PROJ-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.UpdateIssue(context.Background(), "PROJ-1", map[string]any{"summary": "Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateIssue_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.UpdateIssue(context.Background(), "PROJ-1", map[string]any{"summary": "Updated"})
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
