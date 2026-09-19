---
status: active
last-verified: 2026-09-14
---

# Isolate Go tests from Git hook repository variables

## Problem

Git exports `GIT_*` repository bindings to pre-push hooks. The `test` package
starts the `lg` interpreter, whose `.lg` quality tests use `os/sh` to create
and inspect disposable Git repositories. Inherited bindings override `git -C`
repository discovery, so those tests fail only under the hook and can address
the source checkout instead of their fixtures.

## Design

At the start of the package's `TestMain`, remove inherited variables whose
names begin `GIT_` from the Go test process. The boundary is test-only: it
protects TestMain's build, Go fixture commands, and all child interpreter
processes, without changing production `os/sh` or hook configuration. Preserve
non-Git variables. Individual tests remain free to set Git variables after
startup when deliberately testing them. Existing per-command `filterHookEnv`
may remain as defense in depth; no unrelated refactor is required.

## Alternatives considered

- Wrap every `.lg` fixture command with `env -u`: narrow but easy to miss and
  duplicates the same policy across fixtures.
- Add environment overrides to production `os/sh`: useful for other purposes,
  but changes runtime API and scope for a test-harness inheritance bug.

## Evidence

Reproduce with a disposable hook-style `GIT_DIR`/`GIT_WORK_TREE`: the existing
`TestRunner/quality_diff_test.lg` fails against the wrong repository. Add a
focused regression that exercises `TestMain` under a hook-style environment,
asserts `GIT_*` removal and preservation of a non-Git sentinel, and prove it
fails before the fix. Then rerun the same hook-style quality fixture, the short
Go suite, and the
actual #871 pre-push gate. The approved benchmark-ratchet skip applies only to
that pending push; all other hooks remain active.
