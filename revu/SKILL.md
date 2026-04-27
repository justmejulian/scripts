---
name: revu
description: "Understand and write Azure DevOps PR review comments using the REVU marker format. Use when reading existing review threads or adding inline review comments to source files."
---

# revu

Inline PR review comments are embedded in source files using `REVU` markers.

## Reading existing threads

Threads appear as code comments above the reviewed line:

```go
// REVU[100038] @dave: we call this function three times. why.
// REVU[100038] @sarah: the first two are optimistic. the third one actually works.
// REVU[100038] @dave: i need to lie down
chargeCard(user, amount)
```

- `REVU[<id>]` — thread ID, same across all comments in a thread
- `@author` — who wrote the comment
- Status is either `active` (unresolved) or `resolved`

## Adding a new comment

Write `REVU[NEW]` above the line being commented on:

```go
// REVU[NEW] this will panic if user is nil
chargeCard(user, amount)
```

## Replying to an existing thread

Place `REVU[NEW]` on the line immediately below the last comment of the thread:

```go
// REVU[100038] @dave: we call this function three times. why.
// REVU[NEW] agreed, extracting to chargeCardOnce() in next PR
chargeCard(user, amount)
```

The adjacent thread ID signals a reply rather than a new thread.

## Comment prefix by file type

| Extensions | Prefix |
|---|---|
| `.go` `.js` `.ts` `.java` `.c` `.cpp` `.cs` `.swift` `.kt` `.rs` | `//` |
| `.py` `.rb` `.sh` `.yaml` `.yml` `.toml` `.conf` | `#` |
| `.sql` | `--` |
| default | `//` |
