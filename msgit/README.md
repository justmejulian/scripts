# msgit

**msgit** — _message it_ — generates a commit message from your staged changes using the shared AI package.

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
msg=$(msgit) && git commit -F - <<< "$msg"

# Open editor pre-filled with the generated message
msg=$(msgit) && git commit -e -m "$msg"

# Copy to clipboard (macOS)
msgit | pbcopy
```

> **Why `msg=$(msgit) && git commit ...` instead of `msgit | git commit ...`?**
>
> In a pipeline, both sides run concurrently — git doesn't know or care whether msgit succeeded.
> If msgit fails (e.g. missing API key, nothing staged), git still opens the editor with an empty
> message, and you end up aborting the commit manually. Capturing the output first with `$(...)` and
> using `&&` ensures git only runs when msgit exits cleanly.

## How it works

`msgit` collects three pieces of context from git:

1. **Staged diff** (`git diff --cached`) — the actual changes
2. **Branch name** — hints at the feature or ticket
3. **Recent commits** — matches the style and scope of the project

It sends these to the shared AI package, which currently uses local Ollama with `qwen3:8b`, and prints the result to stdout. Thinking is disabled for `msgit`.

The generated message follows [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <summary>

<optional body explaining why>
```

## Configuration

| Variable      | Default                  | Description              |
|---------------|--------------------------|--------------------------|
| `OLLAMA_HOST` | `http://localhost:11434` | Ollama API base URL      |
