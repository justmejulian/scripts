# taskbranch

Fetches Jira tasks assigned to the current user, lets you pick one, and builds a branch name from it.

## Usage

```sh
taskbranch
```

The script will:

1. fetch your assigned Jira tasks
2. let you select a task
3. let you select a branch type
4. print the final branch name

Example output:

```sh
$ taskbranch
Assigned tasks:
1. PROJ-123 - Add login page [In Progress]
2. PROJ-456 - Fix redirect bug [To Do]
Select task [1-2]: 1
Branch types:
1. feat
2. fix
3. chore
4. custom
Select branch type [1-4]: 1
feat/PROJ-123-add-login-page
```

## Flags

```sh
taskbranch --type fix
taskbranch --jql "assignee = currentUser() AND project = PROJ ORDER BY updated DESC"
```

- `--type` skips the branch type prompt
- `--jql` overrides the default Jira query

## Pipe to git

```sh
git checkout -b "$(taskbranch)"
```
