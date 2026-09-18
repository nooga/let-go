---
status: active
last-verified: 2026-09-14
---

# Linux Quality Temp-Path CI Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make #871's Linux quality tests and docs-frontmatter CI jobs pass without changing runtime temp-dir semantics.

**Architecture:** Add an explicit separator at each identified temp-directory-plus-filename construction in the quality script and tests. Add the required YAML header to one existing plan. The already failing Linux-style tests are the behavioral regression, with no new API.

**Tech Stack:** let-go `.lg`, Go 1.26.5 test runner, GitHub Actions, Python frontmatter checker.

---

## Chunk 1: Reproduce, fix, verify

### Task 1: Portable paths and missing header

**Files:**
- Modify: `scripts/quality.lg` (`extract-rev` Git archive tarball)
- Modify: `test/quality_cli_test.lg` (nine EDN temp paths)
- Modify: `test/quality_reap_test.lg` (probe script and pid paths)
- Modify: `docs/superpowers/plans/2026-09-13-quality-scorer-review-fixes.md` (frontmatter only)
- Spec: `docs/superpowers/specs/2026-09-14-quality-temp-path-ci-design.md`

- [ ] Confirm the two red gates before any implementation edit:

```sh
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" TMPDIR=/tmp go test ./test -run '^TestRunner/(quality_cli_test.lg|quality_delta_test.lg)$' -count=1
rtk proxy python3 scripts/docs_frontmatter_hook.py --check docs/superpowers/plans/2026-09-13-quality-scorer-review-fixes.md
```

  Expected: Go test fails with `/tmpquality-cli-*.edn` and `/tmpquality-rev-*.tar`; checker fails with `no frontmatter`. These failures were already reproduced on the #871 parent branch. `TMPDIR=/tmp` is diagnostic only; do not change the developer's global environment.

- [ ] In `scripts/quality.lg` line near `extract-rev`, change only the tarball expression:

```clojure
;; before
(let [tarball (str (os/temp-dir) name ".tar")]
;; after
(let [tarball (str (os/temp-dir) "/" name ".tar")]
```

  Keep `revision-dir`, extraction, cleanup, Git/jj selection, and error data unchanged.

- [ ] In `test/quality_cli_test.lg`, change all nine expressions `(str (os/temp-dir) "quality-...edn")` to `(str (os/temp-dir) "/quality-...edn")`. The exact basenames are `quality-cli-test.edn`, `quality-pr-delta-test.edn`, `quality-direction-test.edn`, `quality-order-test.edn`, `quality-raw-test.edn`, `quality-terms-test.edn`, `quality-dead-test.edn`, `quality-nolint.edn`, and `quality-classes.edn`. Preserve each test's existing cleanup and assertions.

- [ ] In `test/quality_reap_test.lg`, change the three bare concatenations to `(str (os/temp-dir) "/quality-reap-probe.lg")`, `(str (os/temp-dir) "/quality-reap.pid")` in the shell command, and the same pid expression in `slurp`. This is a portability sweep, not a change to the test's expected-fail orphan-process semantics. Its runtime assertion may skip when `../bin/lg` is absent; therefore run `rtk proxy rg -n '\(os/temp-dir\) (name|"quality-)' scripts/quality.lg test/quality_cli_test.lg test/quality_reap_test.lg` and require **no matches** (rg exit 1) as a static sweep of all 13 changed constructions, including the bare pid expression inside the shell command.

- [ ] Prepend this exact YAML block to `docs/superpowers/plans/2026-09-13-quality-scorer-review-fixes.md`; leave all prior text intact:

```yaml
---
status: active
last-verified: 2026-09-14
---

```

- [ ] Re-run the two formerly red commands above; expect both exit 0. Run the default focused quality CLI/delta tests, Go short suite with package parallelism disabled, build/vet, and diff check:

```sh
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" go test ./test -run '^TestRunner/(quality_cli_test.lg|quality_delta_test.lg)$' -count=1
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" go test -p 1 -short ./...
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" go build ./...
rtk proxy env -u GOROOT PATH="/Users/ndn/.local/share/mise/installs/go/1.26.5/bin:$PATH" go vet ./...
rtk git diff --check
```

  Each must exit 0; `diff --check` prints nothing. Inspect the diff to ensure no other path/temp behavior changed.

- [ ] Commit only the four modified files on the #871 follow-up branch. The controller will request spec/code review, integrate the isolated implementation commit, and push with a newly queried exact lease. Only `bench-ratchet` may be skipped on that push, as previously approved. Verify both failed CI jobs turn green; if not, diagnose the new failure before claiming completion.
