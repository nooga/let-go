---
status: active
last-verified: 2026-09-14
---

# Test-Process Git Environment Isolation Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make #871's quality Git fixtures pass under the real pre-push hook without changing production shell behavior.

**Architecture:** `TestMain` clears inherited `GIT_*` variables before building `lg` or running any tests. A focused Go test checks the process boundary; the existing `.lg` Git fixture proves the downstream interpreter receives a clean environment.

**Tech Stack:** Go 1.26.5, Git, `go test`, prek.

---

## Chunk 1: Test-only environment isolation

### Task 1: Guard the test-process boundary

**Files:**
- Modify: `test/fanout_ratchet_test.go` (`TestMain` and nearby helper)
- Test: `test/fanout_ratchet_test.go` (focused regression)
- Spec: `docs/superpowers/specs/2026-09-14-test-process-git-env-isolation-design.md`

- [ ] Add this regression test. It exercises the actual `TestMain` boundary because
  `TestMain` runs before every selected test. `PATH` is a non-Git sentinel that
  must survive, and the optional hook sentinel gives a second explicit check.

```go
func TestHookGitEnvIsolated(t *testing.T) {
    for _, entry := range os.Environ() {
        name, _, _ := strings.Cut(entry, "=")
        if strings.HasPrefix(name, "GIT_") {
            t.Fatalf("inherited Git variable %s", name)
        }
    }
    if os.Getenv("PATH") == "" {
        t.Fatal("non-Git PATH was removed")
    }
    if got, ok := os.LookupEnv("LG_TEST_HOOK_SENTINEL"); ok && got != "kept" {
        t.Fatalf("non-Git sentinel = %q, want kept", got)
    }
}
```

- [ ] Run the red tests with a disposable repo; `GOFLAGS` is diagnostic-only
  because Go's VCS stamping otherwise rejects the deliberately foreign
  `GIT_DIR` before the test binary starts. Do not use the source checkout as
  the disposable target.

```sh
hook_repo=$(rtk proxy mktemp -d /private/tmp/let-go-hook-env.XXXXXX)
rtk proxy git init -q "$hook_repo"
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" GOFLAGS=-buildvcs=false GIT_DIR="$hook_repo/.git" GIT_WORK_TREE="$hook_repo" LG_TEST_HOOK_SENTINEL=kept go test ./test -run '^TestHookGitEnvIsolated$' -count=1
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" GOFLAGS=-buildvcs=false GIT_DIR="$hook_repo/.git" GIT_WORK_TREE="$hook_repo" go test ./test -run '^TestRunner/quality_diff_test.lg$' -count=1
```

  Expected before fix: first command fails with `inherited Git variable GIT_DIR`;
  second fails with `fixture git failed` (or equivalent wrong-repository error).

- [ ] At the first line of `TestMain`, clear only `GIT_*` entries. `strings` and
  `os` are already imported; leave production `os/sh` and the existing
  per-command filter unchanged.

```go
for _, entry := range os.Environ() {
    name, _, _ := strings.Cut(entry, "=")
    if strings.HasPrefix(name, "GIT_") {
        if err := os.Unsetenv(name); err != nil {
            panic(err)
        }
    }
}
```

- [ ] Rerun both focused commands above; expect `ok github.com/nooga/let-go/test`.
- [ ] Run the following checks; expect exit 0 and no `diff --check` output.

```sh
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" go test -short ./...
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" go build ./...
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" go vet ./...
rtk git diff --check
```

- [ ] Commit only `test/fanout_ratchet_test.go` on #871 (this plan is committed before implementation).
- [ ] Query the current remote tip, substitute that exact SHA in the lease,
  then push with only the previously approved hook skip. Do not use a cached
  SHA or skip any other hook.

```sh
rtk proxy git ls-remote upstream refs/heads/quality/complexity-scorer
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" PREK_SKIP=bench-ratchet GIT_SSH_COMMAND="ssh -o ServerAliveInterval=30 -o ServerAliveCountMax=20" git push --force-with-lease=refs/heads/quality/complexity-scorer:<FRESH_SHA> upstream HEAD:refs/heads/quality/complexity-scorer
rtk proxy git ls-remote upstream refs/heads/quality/complexity-scorer
rtk proxy git rev-parse HEAD
```

  Expected: the last two SHAs match.
