# Scripts

Collection of small scripts to automate stuff.

## Structure

### Scripts

- [branchname](./branchname/) - Generates a git branch name from type, issue key, and description
- [taskbranch](./taskbranch/) - Selects an assigned Jira task and turns it into a branch name
- [msgit](./msgit/) - Generates a git commit message from staged changes using Ollama
- [createpr](./createpr/) - Creates a pull request in Azure DevOps from the current branch

### Utils

Reusable libraries that can be imported by scripts.

- [azure](./internal/azure/) - Azure DevOps REST API v7.2 client (pull requests)
- [branchname](./internal/branchname/) - Branch name building and slugification helpers
- [fzf](./internal/fzf/) - Wrapper around `fzf` for interactive selection
- [git](./internal/git/) - Git helpers (e.g. current branch)
- [jira](./internal/jira/README.md) - Jira REST API v2 client with PAT authentication
- [ollama](./internal/ollama/README.md) - AI stuff using [ollama.com](https://ollama.com/)

## Development

### Requirements

- Go 1.24+
- [Ollama](https://ollama.com) running locally (for `msgit`)

```sh
# Install all binaries to $GOPATH/bin (make available system-wide)
go install ./...

# Or install individually
go install ./branchname
go install ./taskbranch
go install ./msgit
go install ./createpr

# Check installed binary version (git commit hash) and compare with current repo
go version -m $(which branchname)  # look for vcs.revision in the output
git rev-parse HEAD

# Run tests
go test ./...

# Format
go fmt ./...

# Lint
go vet ./...
```

Required environment variables (see individual script/library READMEs for details):
- `AZURE_DEVOPS_ORG` - Azure DevOps organisation name
- `AZURE_DEVOPS_TOKEN` - Azure DevOps Personal Access Token
- `JIRA_DOMAIN` - Jira instance domain (e.g. `your-company.atlassian.net`)
- `JIRA_TOKEN` - Jira Personal Access Token
- `OLLAMA_HOST` - Ollama API base URL (default: `http://localhost:11434`)

## Why Go?

Nice and easy to run, no need to install dependencies, just run.

Also, I have been meaning to learn Go for a while, so this is a good excuse to do so.
