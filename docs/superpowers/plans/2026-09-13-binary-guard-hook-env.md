---
status: active
last-verified: 2026-09-13
---

# Binary Guard Hook Environment Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the binary-guard temporary Git fixtures pass inside real Git hooks without bypassing a gate.

**Architecture:** The existing `binaryGuardNoRefEnv` boundary removes inherited hook refs and repository-local `GIT_*` variables. `binaryGuardGit` gives that environment to each child Git process; guard-script subprocesses already use the helper.

**Tech Stack:** Go test, Git, prek pre-push.

---

## Chunk 1: Test-only isolation

### Task 1: Reproduce and fix foreign-repository Git resolution

**Files:**
- Modify: `test/reject_compiled_binaries_test.go`

- [ ] Add a test with two temporary repos and hook-style `GIT_DIR`/`GIT_WORK_TREE`; assert `binaryGuardGit` targets the foreign repo and `binaryGuardNoRefEnv` omits inherited `GIT_*`.
- [ ] Run the focused Go test with pinned Go 1.26.5; require the expected failure before implementation.
- [ ] Filter repository-local Git variables in `binaryGuardNoRefEnv` and assign its result to `binaryGuardGit`'s `cmd.Env`.
- [ ] Rerun focused and full short suites, vet/build, generated-manifest, and diff checks.
- [ ] Commit the test-only fix on #870, rebase #871 onto it with `rebase.updateRefs=false`, and verify #871 independently.
- [ ] Retry the actual pre-push hook without bypassing it.
