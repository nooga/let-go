---
status: active
last-verified: 2026-09-14
---

# Quality scorer review fixes — 2026-09-13

Scope: prepare local fixes for PR #871 without moving its bookmark or updating the PR.

1. Add a collision test for a qualified call to `a/foo` from `b/caller`, with `b/foo` present. Require one edge, then preserve only the qualified reference during graph resolution. A bare call still resolves in its own namespace; an aliased call resolves to the alias target.
2. Add failing delta tests for invalid head and measurement errors, and verify that no revision checkout remains. Put each acquired checkout under `try`/`finally` so head extraction and either measure error release all acquired resources. Derive paths from a configurable revision root; never assume `.workspaces`.
3. Make the tracked-file assertion accept the backend actually selected by `tracked-files`; test both Git and jj behavior where available.
4. Verify the shallow-clone SHA concern and selected Go/let-go tests in a Git-only checkout; record any limits. Do not edit CI workflow without approval; report the trigger decision separately.

Use the published `quality/complexity-scorer@upstream` head as the parent. Preserve the divergent local bookmark and unrelated worktree changes.
