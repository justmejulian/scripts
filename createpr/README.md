# createpr

Creates a pull request in Azure DevOps from the current branch, then transitions the Jira issue to "In Review" and adds the PR URL as a comment.

## Usage

```sh
# Run from within the target repository
go run ./createpr

# Or if installed
createpr
# https://dev.azure.com/your-org/project/_git/repo/pullrequest/42
```

The PR title is derived automatically from the branch name and the Jira issue title:

```
PROJ-123 feat: Add login page
```

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-target` | `main` | Target branch for the PR |

```sh
createpr -target develop
```

## Requirements

Must be run from inside a git repository whose path follows the convention:

```
.../project/repo/
```

The two trailing path segments are used as the Azure DevOps project and repository names.

The branch must contain a Jira issue key and a type prefix, e.g. `feat/PROJ-123-add-login-page`.

## Environment variables

| Variable | Description |
|----------|-------------|
| `AZURE_DEVOPS_ORG` | Azure DevOps organisation name |
| `AZURE_DEVOPS_TOKEN` | Personal Access Token |
| `JIRA_DOMAIN` | Jira instance domain (e.g. `your-company.atlassian.net`) |
| `JIRA_TOKEN` | Jira Personal Access Token |
