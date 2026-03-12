# msgit

**msgit** — _message it_ — generates a commit message from your staged changes using Ollama.

The name is a mashup of "message" and "git": you stage your changes, run `msgit`, and it messages you what to commit.

## Requirements

- [Ollama](https://ollama.com) running locally with `qwen3:8b` pulled:
  ```sh
  ollama pull qwen3:8b
  ```

## Usage

```sh
# Preview the generated message
git add <files>
msgit

# Pipe directly into git commit
msgit | git commit -F -

# Open editor pre-filled with the generated message
go run ./msgit | git commit -e -F -

# Copy to clipboard (macOS)
msgit | pbcopy
```

## How it works

`msgit` collects three pieces of context from git:

1. **Staged diff** (`git diff --cached`) — the actual changes
2. **Branch name** — hints at the feature or ticket
3. **Recent commits** — matches the style and scope of the project

It sends these to `qwen3:8b` via the local Ollama API and prints the result to stdout. Thinking tokens (`<think>…</think>`) that qwen3 emits are stripped before output.

The generated message follows [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <summary>

<optional body explaining why>
```

## Configuration

| Variable      | Default                  | Description              |
|---------------|--------------------------|--------------------------|
| `OLLAMA_HOST` | `http://localhost:11434` | Ollama API base URL      |
