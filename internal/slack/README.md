# slack

Go client for the Slack Web API, supporting message posting.

## Authentication

Requires a Bot token with the `chat:write` scope. Set the following environment variable:

```sh
export SLACK_TOKEN=xoxb-your-bot-token
```

Or initialize the client directly:

```go
client := slack.NewClient("xoxb-your-bot-token")
```

## Usage

```go
import "scripts/internal/slack"

client, err := slack.NewClientFromEnv()
if err != nil {
    log.Fatal(err)
}

// Post a message to a channel
err = client.PostMessage(ctx, "#general", "Hello!")
```

## Error Handling

API errors are returned as `*APIError` with Slack's error code:

```go
var apiErr *slack.APIError
if errors.As(err, &apiErr) {
    fmt.Println(apiErr.Code) // e.g. "channel_not_found"
}
```
