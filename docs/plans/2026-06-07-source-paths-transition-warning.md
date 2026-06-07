# Source-Paths Transition Warning Implementation Plan

> **For agentic workers:** Use executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Warn users, for a release or two, that the current directory (`.`) is no longer added to the namespace search path implicitly when `-source-paths`/`LG_SOURCE_PATHS` is set — so the behavior change from commit `5c75260` does not break them silently.

**Tech Stack:** Go 1.26 (`package main` in `lg.go`), stdlib `slices`, existing e2e test harness in `source_paths_e2e_test.go`.

---

## Design

### Background

Commit `5c75260` changed `buildSearchPaths` so that an explicit `-source-paths` flag or `LG_SOURCE_PATHS` env var *fully* defines the namespace search path. Previously `.` (the current directory) was always prepended; now it is only searched if listed explicitly. A user who upgrades and runs `lg -source-paths src app.lg` will silently stop finding namespaces that live in the current directory.

The review comment asks for a transition warning "for a release or two" so this change is surfaced rather than silent.

### Approach

Two pieces:

1. **A generic warning helper** `warn(format, args...)` in `lg.go`. It prints `warning: <message>` to stderr unless `LG_SUPPRESS_WARNINGS` is set to any non-empty value. This is a reusable switch for *intentional* transition/deprecation notices — source-paths is its first caller, and the same env var will silence future ones. The non-empty-value convention matches existing boolean env vars in the codebase (`LG_BOXARGS_DEBUG`, `LG_PANIC_STACK`, `LG_READ_CLJ`, all gated with `os.Getenv(...) != ""`). The `warning:` prefix mirrors the existing `error:` prefix used throughout `lg.go`.

2. **The trigger**, inside `buildSearchPaths`. In the explicit branch (`explicitSet || envSet`), after computing the parsed `paths`, emit the warning when `"."` is not among them. This deliberately **includes the empty case** (`-source-paths ""` / `LG_SOURCE_PATHS=`, which parse to no paths): per the design decision, any explicit setting that lacks `.` warns. Placing the check inside `buildSearchPaths` means a single decision point covers both call sites (bundled exec at `lg.go:690`, normal run at `lg.go:755`), and it fires at most once per process because `buildSearchPaths` runs once per process (the bundled branch `return`s before the normal path).

The default path (no flag and no env var → keeps `.` plus any `deps.edn` `:paths`) never warns. Only an exact `"."` entry suppresses the warning; alternate spellings like `"./"` or an absolute path to the cwd are out of scope (documented behavior: "add `.` to the list").

### Warning message

```
warning: the current directory (".") is no longer added to the namespace search path automatically when -source-paths or LG_SOURCE_PATHS is set; add "." to the list to keep searching it. This notice will be removed in a future release (set LG_SUPPRESS_WARNINGS=1 to silence).
```

### Error handling

None beyond the gate: the helper is best-effort output to stderr. No new error paths.

### Testing strategy

E2e subtests in `source_paths_e2e_test.go` (build the real binary, assert on combined output), mirroring the existing `TestSourcePathsControlSearchPath` style and reusing `buildLG(t)` plus the `cleanEnv`/`run` pattern. Cases:

- explicit flag without `.` → output **contains** the warning
- explicit flag **with** `.` → output does **not** contain the warning
- empty env (`LG_SOURCE_PATHS=`) → output **contains** the warning
- explicit env without `.` → output **contains** the warning
- `LG_SUPPRESS_WARNINGS=1` set alongside an explicit path without `.` → output does **not** contain the warning
- plain default run (no flag, no env) → output does **not** contain the warning

A short, stable sentinel substring (e.g. `is no longer added to the namespace search path`) is used for assertions so minor wording tweaks do not break tests.

## File Structure

- **Modify `lg.go`** — add the `warn` helper; add the trigger inside `buildSearchPaths`; add `"slices"` to imports. Update the `-source-paths` flag help text only if needed (optional; the README is the primary doc surface).
- **Modify `source_paths_e2e_test.go`** — add a new test function for the warning behavior (keep it separate from `TestSourcePathsControlSearchPath` for clarity).
- **Modify `README.md`** — extend the existing "Source paths" section with a sentence on the transition warning and the `LG_SUPPRESS_WARNINGS=1` opt-out.

---

## Task 1: Generic `warn` helper + trigger in `buildSearchPaths`

**Files:**
- Modify: `lg.go`
- Test: `source_paths_e2e_test.go`

- [x] **Step 1: Write the failing e2e test**
  In `source_paths_e2e_test.go`, add `TestSourcePathsTransitionWarning`. Reuse the existing pattern (`buildLG(t)`, a `cleanEnv` that strips `LG_SOURCE_PATHS` *and* `LG_SUPPRESS_WARNINGS`, and a `run` helper that sets `cmd.Dir` to a temp dir and returns combined output). Define `const warnMarker = "is no longer added to the namespace search path"`. Add subtests:
  - `explicit flag without dot warns` → `-source-paths <libDir> -e "(println 1)"`, assert output **contains** `warnMarker`.
  - `explicit flag with dot does not warn` → `-source-paths .:<libDir> -e "(println 1)"`, assert output does **not** contain `warnMarker`.
  - `empty env warns` → env `LG_SOURCE_PATHS=`, assert **contains**.
  - `explicit env without dot warns` → env `LG_SOURCE_PATHS=<libDir>`, assert **contains**.
  - `suppress env silences` → `-source-paths <libDir>` plus env `LG_SUPPRESS_WARNINGS=1`, assert does **not** contain.
  - `default run does not warn` → no flag, no env, assert does **not** contain.

- [x] **Step 2: Run the test to verify it fails**
  Run: `go test -run TestSourcePathsTransitionWarning .`
  Expected: FAIL — the warning is never emitted, so the "contains" subtests fail.

- [x] **Step 3: Implement the `warn` helper and trigger**
  In `lg.go`:
  - Add `"slices"` to the import block.
  - Add the `warn(format string, args ...any)` helper near `buildSearchPaths`: return early if `os.Getenv("LG_SUPPRESS_WARNINGS") != ""`, otherwise `fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)`. Include a doc comment describing it as the switch for intentional transition warnings.
  - In `buildSearchPaths`, change the explicit branch to bind `paths := resolver.PathsFromInputs(...)`, then `if !slices.Contains(paths, ".") { warn(<message>) }`, then `return paths`. Use the message from the Design section.

- [x] **Step 4: Run the test to verify it passes**
  Run: `go test -run TestSourcePathsTransitionWarning .`
  Expected: PASS.

- [x] **Step 5: Run the existing source-paths suite to confirm no regressions**
  Run: `go test -run TestSourcePaths . ./pkg/resolver/...`
  Expected: PASS (existing `TestSourcePathsControlSearchPath` and resolver tests unaffected — they assert on resolution behavior, not warning text).

- [x] **Step 6: Commit**
  `git commit -am "feat(source-paths): warn when explicit search path drops implicit '.'"`

## Task 2: Document the warning and opt-out

**Files:**
- Modify: `README.md`

- [x] **Step 1: Update the "Source paths" section**
  In the existing "### Source paths" section, after the sentence explaining that the current directory is not searched implicitly, add a short note: setting the search path without `.` prints a transition warning (to be removed in a future release) to help spot reliance on the old implicit-cwd behavior, and `LG_SUPPRESS_WARNINGS=1` silences it (and any future transition warnings). Use the /writing-clearly skill for the wording.

- [x] **Step 2: Verify the docs build/render**
  Run: `git diff README.md`
  Expected: the new sentence reads cleanly and matches surrounding style; no broken markdown.

- [x] **Step 3: Commit**
  `git commit -am "docs(source-paths): note transition warning and LG_SUPPRESS_WARNINGS opt-out"`

## Task 3: Full verification

- [x] **Step 1: Build and vet**
  Run: `go build ./... && go vet ./...`
  Expected: clean.

- [x] **Step 2: Run the full test suite**
  Run: `go test ./...`
  Expected: PASS.

---

## Status: Completed (2026-06-07)

All three tasks implemented and verified.

### What was implemented

- **`lg.go`** — added a generic `warn(format, args...)` helper that prints
  `warning: …` to stderr unless `LG_SUPPRESS_WARNINGS` is set to any non-empty
  value (matching the codebase's `os.Getenv(...) != ""` boolean-env idiom). In
  `buildSearchPaths`, the explicit branch now warns via `slices.Contains` when
  the parsed path lacks `"."` — covering the empty case (`-source-paths ""` /
  `LG_SOURCE_PATHS=`) per the design decision. Added the `slices` import.
- **`source_paths_e2e_test.go`** — new `TestSourcePathsTransitionWarning` with
  six subtests (explicit flag with/without `.`, empty env, explicit env,
  `LG_SUPPRESS_WARNINGS` suppression, default run), asserting on a stable
  `warnMarker` substring so wording tweaks don't break tests.
- **`README.md`** — extended the "Source paths" section with a note on the
  transition warning and the `LG_SUPPRESS_WARNINGS=1` opt-out.

### Verification

- `go test -run TestSourcePathsTransitionWarning .` — PASS
- `go test -run TestSourcePaths . ./pkg/resolver/...` — PASS (no regressions)
- `go build ./... && go vet ./...` — clean
- `go test ./...` — PASS
- Second-opinion review via `review-with-codex` (uncommitted scope) — no
  correctness issues found.

### Notes / known nuances

- In the bundled-binary branch, `buildSearchPaths()` runs at `lg.go:690`
  *before* `flag.Parse()`, so `flag.Visit` finds nothing and only the
  `LG_SOURCE_PATHS` env channel can trigger the warning there. Harmless
  (bundles are pre-resolved; the env channel still works) and pre-existing;
  not changed by this work.
- Only an exact `"."` entry suppresses the warning; alternate spellings
  (`"./"`, an absolute cwd path) still warn, as documented.

### Commits

- `48de96b` feat(source-paths): warn when explicit search path drops implicit '.'
- `9260167` docs(source-paths): note transition warning and LG_SUPPRESS_WARNINGS opt-out
