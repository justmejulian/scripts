# AI

Provider-agnostic AI abstraction for scripts in this repo.

App code should use this package to construct a provider and submit requests.

## Usage

```go
import "scripts/internal/ai"

provider, err := ai.NewProvider(ai.Config{Provider: "ollama"})
if err != nil {
	panic(err)
}

resp, err := provider.Generate(ctx, ai.Request{
	Prompt: "Say hello world.",
	Model:  "qwen3:8b",
	Think:  false,
})
```

## Providers

- `ollama` - local Ollama chat provider
  - requires:
    - [Ollama](https://ollama.com) running locally
    - `OLLAMA_HOST` env var
- `zen` - [opencode.ai](https://opencode.ai) cloud provider
  - requires:
    - `ZEN_API_KEY` env var

## Adding a provider

Create `internal/ai/providers/<name>/` with two files:

- `client.go` — implement `spec.Provider`, export `const Name = "<name>"`, a `New() (spec.Provider, error)` constructor, and an `init()` that calls `registry.Register(Name, New)`
- `models.go` — string constants for model names (optional but conventional)
