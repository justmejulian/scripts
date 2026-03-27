# Ollama Provider

Ollama implementation for `scripts/internal/ai`.

This package is intended to be used via `scripts/internal/ai`, not imported directly by app code.

## Configuration

| Env var | Default | Description |
|---|---|---|
| `OLLAMA_HOST` | `http://localhost:11434` | Base URL of the Ollama server |

## Notes

- Maps `ai.Request` to Ollama's `/api/chat` endpoint
- Supports the shared `Think` option
- Strips `<think>...</think>` output before returning text
