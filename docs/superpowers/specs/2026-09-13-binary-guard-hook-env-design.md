---
status: active
last-verified: 2026-09-13
---

# Keep binary-guard test repositories isolated from Git hooks

## Problem

Git exports repository-local variables such as `GIT_DIR` and `GIT_WORK_TREE`
to hooks. The binary-guard tests create temporary Git repositories, but their
child commands inherit those variables. During pre-push, `git` can therefore
operate on the source worktree instead of the temporary repository and fail
the Go test gate.

## Design

Clean `GIT_*` variables from the environment used by this test file's
temporary-repository Git commands and guard-script subprocesses. Retain other
environment variables (including `PATH`) and the guard's deliberately supplied
`PRE_COMMIT_FROM_REF`/`PRE_COMMIT_TO_REF` values. The existing helper that strips
inherited `PRE_COMMIT_*` refs is the single boundary for this cleanup.

Changing the hook configuration would affect every pre-push job and requires
separate CI approval; this fix remains test-local. No hook is skipped.

## Evidence

A regression uses two disposable Git repositories, sets hook-style `GIT_DIR`
and `GIT_WORK_TREE` to the first, and asserts that the helper initializes and
uses the second. The test must fail before the change. Then run focused tests,
the full short suite, and the actual pre-push gate on the #870 layer; cascade
the commit to #871 and verify its layer independently.
