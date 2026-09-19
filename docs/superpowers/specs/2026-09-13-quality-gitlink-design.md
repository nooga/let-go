---
status: active
last-verified: 2026-09-13
---

# Gitlink Changes in Quality Deltas

## Problem

`quality.diff/git-changes` reads changed text with `git show <commit>:<path>`.
Gitlinks (mode `160000`) are commit pointers, not file blobs; `git show` for
their paths fails and aborts an otherwise valid quality delta. A dirty
submodule can expose this even when the outer repository has no commit change.

## Design

Classify a Gitlink from its Git patch mode, alongside the existing binary and
symlink classification. When Git omits the mode for a dirty submodule, recognize
its paired `Subproject commit` hunk only if no regular-file mode is present.
Preserve its path and status, but omit old/new text so
`quality.touched/build` reports a path-only entry with no invented form metrics.
Use the same behavior for commit-to-commit and commit-to-working-tree deltas.
No jj dependency, workspace path assumption, or extra Git process is added.

## Verification

A plain-Git temporary repository will contain a real mode-`160000` Gitlink.
The test first proves commit-to-commit and dirty working-tree deltas fail on
the current loader, then verifies both return path-only entries. Focused and
short Go suites will run with the pinned Go binary.
