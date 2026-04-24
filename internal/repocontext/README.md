# repocontext

Derives the Azure DevOps project and repository from the current working directory path.

## Directory structure

Repositories must be cloned into a two-level structure where the parent directory is the Azure DevOps project name and the directory itself is the repository name:

```
<project>/
└── <repo>/
    └── ... (your code)
```

**Example:**

```
my-project/
└── my-service/
    └── src/
```

Running any tool from within `my-service/` will resolve to `project=my-project`, `repo=my-service`.

## Why

Azure DevOps scopes repositories under a project. By mirroring that structure on disk, no extra configuration or flags are needed — the context is implicit from where you run the command.
