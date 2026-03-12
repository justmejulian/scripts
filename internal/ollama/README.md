# Ollama

Go client for a locally-running [Ollama](https://ollama.com) instance.

## Usage

```go
import "scripts/internal/ollama"

c := ollama.NewClient("qwen3:8b")
reply, err := c.Chat(ctx, "Say hello world.")
```

## Configuration

| Env var | Default | Description |
|---|---|---|
| `OLLAMA_HOST` | `http://localhost:11434` | Base URL of the Ollama server |

## API

### `NewClient(model string) *Client`

Creates a new client targeting the given model.

### `(*Client) Chat(ctx context.Context, prompt string) (string, error)`

Sends a single user message and returns the assistant's reply. Streaming is disabled; the full response is returned at once.

## Example

See [`ollama/main.go`](../../ollama/main.go) for a runnable example.
