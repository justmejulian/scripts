# Scripts

Collection of small scripts to automate stuff.

## Structure

### Scripts

- [taskfinder](./taskfinder/) - Lists Jira issues assigned to the current user

### Utils

Reusable libraries that can be imported by scripts.

- [Jira](./internal/jira/README.md) - Jira REST API v2 client with PAT authentication
- [ollama](./internal/ollama/README.md) - AI stuff using [ollama.com](https://ollama.com/)

## Development

```sh
# Run a script
go run ./taskfinder

# Run tests
go test ./...

# Format
go fmt ./...
```

Required environment variables (see individual script/library READMEs for details):
- `JIRA_BASE_URL` - Jira instance base URL
- `JIRA_TOKEN` - Personal Access Token

## Why Go?

Nice and easy to run, no need to install dependencies, just run.

Also, I have been meaning to learn Go for a while, so this is a good excuse to do so.
