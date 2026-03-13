# Scripts

Collection of small scripts to automate stuff.

## Structure

### Scripts

- [branchname](./branchname/) - Generates a git branch name from type, issue key, and description
- [taskbranch](./taskbranch/) - Selects an assigned Jira task and turns it into a branch name
- [msgit](./msgit/) - Generates a git commit message from staged changes using Ollama

### Utils

Reusable libraries that can be imported by scripts.

- [Jira](./internal/jira/README.md) - Jira REST API v2 client with PAT authentication
- [ollama](./internal/ollama/README.md) - AI stuff using [ollama.com](https://ollama.com/)

## Development

### Requirements

- Go 1.24+
- [Ollama](https://ollama.com) running locally (for `msgit`)

```sh
# Install binaries to $GOPATH/bin (make available system-wide)
go install ./branchname
go install ./taskbranch
go install ./msgit

# Run a script
go run ./branchname --type feat --issue PROJ-123 --description "add login page"
go run ./taskbranch
go run ./msgit

# Pipe branchname directly to git
git checkout -b $(branchname --type feat --issue PROJ-123 --description "add login page")

# Run tests
go test ./...

# Format
go fmt ./...

# Lint
go vet ./...
```

Required environment variables (see individual script/library READMEs for details):
- `JIRA_DOMAIN` - Jira instance domain (e.g. `your-company.atlassian.net`)
- `JIRA_TOKEN` - Personal Access Token
- `OLLAMA_HOST` - Ollama API base URL (default: `http://localhost:11434`)

## Why Go?

Nice and easy to run, no need to install dependencies, just run.

Also, I have been meaning to learn Go for a while, so this is a good excuse to do so.
