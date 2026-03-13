# branchname

Generates a git branch name from a task type, issue key, and description.

## Usage

```sh
branchname --type feat --issue PROJ-123 --description "add login page"
# feat/PROJ-123-add-login-page
```

All flags are optional; missing parts are simply omitted:

```sh
branchname --type feat --issue PROJ-123
# feat/PROJ-123

branchname --issue PROJ-123 --description "add login page"
# PROJ-123-add-login-page
```

## Pipe to git

```sh
# Create and switch to the new branch
git checkout -b $(branchname --type feat --issue PROJ-123 --description "add login page")

# Or with xargs
branchname --type feat --issue PROJ-123 --description "add login page" | xargs git checkout -b

# Copy to clipboard (macOS)
branchname --type feat --issue PROJ-123 --description "add login page" | pbcopy
```
