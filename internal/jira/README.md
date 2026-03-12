# jira

Go client for the Jira REST API v2, supporting issue CRUD, JQL search, and user lookups.

## Authentication

Requires a Personal Access Token (PAT). Set the following environment variables:

```sh
export JIRA_DOMAIN=your-company.atlassian.net
export JIRA_TOKEN=your-pat-token
```

Or initialize the client directly:

```go
client := jira.NewClientPAT("your-company.atlassian.net", "your-token")
```

## Usage

```go
import "scripts/internal/jira"

client, err := jira.NewClientFromEnv()
if err != nil {
    log.Fatal(err)
}

// Get a single issue
issue, err := client.GetIssue(ctx, "PROJ-123")

// Create an issue
issue, err := client.CreateIssue(ctx, "PROJ", "Summary text", "Task")

// Update issue fields
err = client.UpdateIssue(ctx, "PROJ-123", map[string]any{
    "summary": "Updated summary",
})

// Search with JQL
issues, err := client.SearchIssues(ctx, "assignee = currentUser() AND status = 3", []string{"key", "summary", "status"})

// Get authenticated user
user, err := client.GetCurrentUser(ctx)

// Find users
users, err := client.FindUser(ctx, "jane")
```

## Types

```go
type Issue struct {
    Key    string
    Title  string
    Status string
}

type User struct {
    AccountID    string
    DisplayName  string
    EmailAddress string
}
```

## Error Handling

API errors are returned as `*APIError` with an HTTP status code:

```go
var apiErr *jira.APIError
if errors.As(err, &apiErr) {
    fmt.Println(apiErr.StatusCode) // e.g. 404
}
```

## Tests

```sh
go test ./internal/jira/...
```
