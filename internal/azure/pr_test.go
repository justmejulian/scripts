package azure

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPRThreads_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(prThreadsResponse{
			Value: []PRThread{
				{
					ID:     1,
					Status: "active",
					Comments: []PRComment{
						{ID: 1, Content: "looks good", CommentType: "text", Author: IdentityRef{DisplayName: "Alice"}},
					},
				},
				{
					ID:     2,
					Status: "resolved",
					Comments: []PRComment{
						{ID: 1, Content: "nit: rename this", CommentType: "text", Author: IdentityRef{DisplayName: "Bob"}},
						{ID: 2, Content: "done", CommentType: "text", Author: IdentityRef{DisplayName: "Alice"}, ParentCommentID: 1},
					},
				},
			},
		})
	}))
	defer srv.Close()

	c := &Client{baseURL: srv.URL, authHeader: "Basic dGVzdA==", http: &http.Client{}}
	threads, err := c.GetPRThreads(context.Background(), "myproject", "myrepo", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(threads) != 2 {
		t.Fatalf("expected 2 threads, got %d", len(threads))
	}
	if threads[0].Status != "active" {
		t.Errorf("unexpected status: %s", threads[0].Status)
	}
	if threads[0].Comments[0].Content != "looks good" {
		t.Errorf("unexpected comment content: %s", threads[0].Comments[0].Content)
	}
	if len(threads[1].Comments) != 2 {
		t.Errorf("expected 2 comments in thread 2, got %d", len(threads[1].Comments))
	}
}

func TestGetPRThreads_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := &Client{baseURL: srv.URL, authHeader: "Basic dGVzdA==", http: &http.Client{}}
	_, err := c.GetPRThreads(context.Background(), "myproject", "myrepo", 99)
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code: %d", apiErr.StatusCode)
	}
}

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
