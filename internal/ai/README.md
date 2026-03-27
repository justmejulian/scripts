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
    - [Ollama](https://ollama.com) running locally
