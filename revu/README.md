# revu

CLI tool for reviewing Azure DevOps pull requests from the terminal.

Inspired by [jez's CLI code review workflow](https://blog.jez.io/cli-code-review/). I use [diffview.nvim](https://github.com/sindrets/diffview.nvim) for diffs — but reviewing in the terminal while posting comments via the browser meant constant context-switching.

Old workaround: drop inline comments in the code, step through them with `git add -p`, then manually post each one via the browser. The friction was intentional — re-reading your own comments before posting catches noise and half-formed thoughts before they land on reviewers.

A byproduct of keeping review threads in source files as plain text: AI can read them too. With `revu sync`, an AI agent sees the full thread context next to the code and can help draft replies or suggest fixes — no copy-pasting from the browser.

## Setup

```sh
export AZURE_DEVOPS_ORG=your-org
export AZURE_DEVOPS_TOKEN=your-pat
```

Run from inside the repo being reviewed. The project and repository are derived from the last two path segments of the working directory.

## Commands

### `revu comments`

Print all PR comments for the current branch.

```
[active] /src/auth/session.go:42
REVU[100042] @dave: this sleep(2000) is doing a lot of heavy lifting

[active] /src/billing/charge.go:7
REVU[100038] @dave: we call this function three times. why.
REVU[100038] @sarah: the first two are optimistic. the third one actually works.
REVU[100038] @dave: i need to lie down
```

### `revu sync`

Inject PR comments into source files as code comments above the reviewed lines. Use `git diff` to review inline.

```sh
revu sync                # inject all thread comments
revu sync --active-only  # skip resolved threads
```

Comments are prefixed with `REVU[<threadID>]` so they can be identified and stripped. Sync is idempotent — existing REVU comments are removed before re-inserting. `REVU[NEW]` lines added by the user are preserved.

To write a multi-line comment, prefix each continuation line with `REVU[NEW]`:

```go
// REVU[NEW] Three issues with current approach:
// REVU[NEW]  1. Status values duplicated across handlers.
// REVU[NEW]  2. Plain string return makes follow-up checks harder.
// REVU[NEW] Fix: return a typed result and share the mapping logic.
func getStatus() string {
```

All consecutive `REVU[NEW]` lines are joined into a single PR comment on upload.

### `revu clean`

Remove all injected REVU comments from source files, including any `REVU[NEW]` lines.

```sh
revu clean
```

**Example diff after sync:**

```diff
+// REVU[100038] @dave: we call this function three times. why.
+// REVU[100038] @sarah: the first two are optimistic. the third one actually works.
+// REVU[100038] @dave: i need to lie down
 chargeCard(user, amount);
```

#### Supported comment prefixes

| Extensions | Prefix |
|---|---|
| `.go` `.js` `.ts` `.java` `.c` `.cpp` `.cs` `.swift` `.kt` `.rs` | `//` |
| `.py` `.rb` `.sh` `.yaml` `.yml` `.toml` `.conf` | `#` |
| `.sql` | `--` |
| default | `//` |

## Typical workflow

```sh
revu sync --active-only   # inject unresolved comments
git diff                  # review inline
revu clean                # restore files when done
```

## Neovim integration

[revu.nvim](https://github.com/justmejulian/.dotfiles/tree/main/.config/nvim/local/revu.nvim) — Neovim plugin for revu. Run commands and navigate injected comments without leaving the editor.

## Claude / AI usage

[SKILL.md](SKILL.md) teaches Claude the REVU marker format — reading threads, adding `REVU[NEW]` comments, and replying to existing threads. It follows the [Agent Skills](https://agentskills.io) open standard and works with Claude Code and GitHub Copilot.

### Install

**Personal** (available across all projects):

```sh
mkdir -p ~/.claude/skills/revu
ln -s /path/to/revu/SKILL.md ~/.claude/skills/revu/SKILL.md
```

**Project** (checked into the repo being reviewed):

```sh
mkdir -p .claude/skills/revu
cp /path/to/revu/SKILL.md .claude/skills/revu/SKILL.md
```

### Usage

```sh
revu sync --active-only   # inject unresolved threads into source files
```

Then in Claude Code:

```
/revu
```

Claude reads the injected REVU threads and can draft replies, suggest fixes, or add new `REVU[NEW]` comments inline. Finish with `revu upload` to post them as PR threads, then `revu clean` to restore the files.

## Build

```sh
go build -o revu .
```
