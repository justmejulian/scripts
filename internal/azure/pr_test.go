package azure

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePR_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req CreatePRRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Title != "My PR" {
			t.Errorf("unexpected title: %s", req.Title)
		}
		if req.SourceRefName != "refs/heads/feature" {
			t.Errorf("unexpected sourceRefName: %s", req.SourceRefName)
		}
		if req.TargetRefName != "refs/heads/main" {
			t.Errorf("unexpected targetRefName: %s", req.TargetRefName)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(PullRequest{
			PullRequestID: 42,
			Title:         req.Title,
			SourceRefName: req.SourceRefName,
			TargetRefName: req.TargetRefName,
			URL:           "https://example.com/pr/42",
		})
	}))
	defer srv.Close()

	c := &Client{baseURL: srv.URL, authHeader: "Basic dGVzdA==", http: &http.Client{}}
	pr, err := c.CreatePR(context.Background(), "myproject", "myrepo", CreatePRRequest{
		Title:         "My PR",
		SourceRefName: "refs/heads/feature",
		TargetRefName: "refs/heads/main",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.PullRequestID != 42 {
		t.Errorf("unexpected PR ID: %d", pr.PullRequestID)
	}
	if pr.URL != "https://example.com/pr/42" {
		t.Errorf("unexpected URL: %s", pr.URL)
	}
}

func TestCreatePR_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))
	defer srv.Close()

	c := &Client{baseURL: srv.URL, authHeader: "Basic dGVzdA==", http: &http.Client{}}
	_, err := c.CreatePR(context.Background(), "myproject", "myrepo", CreatePRRequest{
		Title:         "My PR",
		SourceRefName: "refs/heads/feature",
		TargetRefName: "refs/heads/main",
	})
	if err == nil {
		t.Fatal("expected error for non-201 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusConflict {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestCreatePR_OmitsEmptyDescription(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body["description"]; ok {
			t.Error("description should be omitted when empty")
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(PullRequest{PullRequestID: 1})
	}))
	defer srv.Close()

	c := &Client{baseURL: srv.URL, authHeader: "Basic dGVzdA==", http: &http.Client{}}
	c.CreatePR(context.Background(), "myproject", "myrepo", CreatePRRequest{
		Title:         "My PR",
		SourceRefName: "refs/heads/feature",
		TargetRefName: "refs/heads/main",
	})
}
