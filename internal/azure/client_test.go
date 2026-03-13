package azure

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClientFromEnv_MissingVars(t *testing.T) {
	t.Setenv("AZURE_DEVOPS_ORG", "")
	t.Setenv("AZURE_DEVOPS_TOKEN", "")

	_, err := NewClientFromEnv()
	if err == nil {
		t.Fatal("expected error when env vars are missing")
	}
}

func TestNewClientFromEnv_MissingOrg(t *testing.T) {
	t.Setenv("AZURE_DEVOPS_ORG", "")
	t.Setenv("AZURE_DEVOPS_TOKEN", "mytoken")

	_, err := NewClientFromEnv()
	if err == nil {
		t.Fatal("expected error when AZURE_DEVOPS_ORG is missing")
	}
}

func TestNewClientFromEnv_MissingToken(t *testing.T) {
	t.Setenv("AZURE_DEVOPS_ORG", "myorg")
	t.Setenv("AZURE_DEVOPS_TOKEN", "")

	_, err := NewClientFromEnv()
	if err == nil {
		t.Fatal("expected error when AZURE_DEVOPS_TOKEN is missing")
	}
}

func TestNewClientFromEnv_Success(t *testing.T) {
	t.Setenv("AZURE_DEVOPS_ORG", "myorg")
	t.Setenv("AZURE_DEVOPS_TOKEN", "mytoken")

	c, err := NewClientFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.baseURL != "https://dev.azure.com/myorg" {
		t.Errorf("unexpected baseURL: %s", c.baseURL)
	}
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte(":mytoken"))
	if c.authHeader != expected {
		t.Errorf("unexpected authHeader: %s", c.authHeader)
	}
}

func TestNewClientPAT(t *testing.T) {
	c := NewClientPAT("myorg", "mytoken")
	if c.baseURL != "https://dev.azure.com/myorg" {
		t.Errorf("unexpected baseURL: %s", c.baseURL)
	}
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte(":mytoken"))
	if c.authHeader != expected {
		t.Errorf("unexpected authHeader: %s", c.authHeader)
	}
	if c.http == nil {
		t.Error("expected http client to be set")
	}
}

func TestSendRequest_SetsHeaders(t *testing.T) {
	var gotAuth, gotAccept, gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotAccept = r.Header.Get("Accept")
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		authHeader: "Basic dGVzdA==",
		http:       &http.Client{},
	}

	resp, err := c.sendRequest(context.Background(), http.MethodGet, srv.URL+"/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if gotAuth != "Basic dGVzdA==" {
		t.Errorf("unexpected Authorization header: %s", gotAuth)
	}
	if gotAccept != "application/json" {
		t.Errorf("unexpected Accept header: %s", gotAccept)
	}
	if gotContentType != "application/json" {
		t.Errorf("unexpected Content-Type header: %s", gotContentType)
	}
}

func TestSendRequest_WithBody(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		authHeader: "Basic dGVzdA==",
		http:       &http.Client{},
	}

	body := strings.NewReader(`{"key":"value"}`)
	resp, err := c.sendRequest(context.Background(), http.MethodPost, srv.URL+"/test", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if gotBody != `{"key":"value"}` {
		t.Errorf("unexpected body: %s", gotBody)
	}
}

func TestSendRequest_InvalidURL(t *testing.T) {
	c := &Client{
		baseURL:    "",
		authHeader: "Basic dGVzdA==",
		http:       &http.Client{},
	}

	_, err := c.sendRequest(context.Background(), http.MethodGet, "://invalid-url", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestURL(t *testing.T) {
	c := NewClientPAT("myorg", "mytoken")
	got := c.url("myproject", "wit/workitems/42")
	want := "https://dev.azure.com/myorg/myproject/_apis/wit/workitems/42?api-version=7.2"
	if got != want {
		t.Errorf("url() = %s, want %s", got, want)
	}
}
