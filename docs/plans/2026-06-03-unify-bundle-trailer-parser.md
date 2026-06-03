# Unify Bundle Trailer Parser Implementation Plan

> **Status: ✅ COMPLETED (2026-06-03).** See the Implementation Summary at the end.

> **For agentic workers:** Use executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Collapse the two standalone-binary trailer parsers in `lg.go` into one validated helper, restoring dropped error checks and guarding against a corrupt trailer that can panic at `make([]byte, lgbSize)`.

**Tech Stack:** Go (`lg.go`, root `package main`).

---

## Design

`lg.go` parses the appended trailer of a standalone binary in two places that duplicate the `LGB2`/`LGBX` magic discrimination and layout:

- `readBundledLGB` (extracts the lgb + resource archive at startup) checks every `Seek`/`ReadFull`.
- `getBaseBinarySize` (strips an existing payload when re-bundling) re-implements the same logic but drops the `Seek` error checks and does unchecked `total - int64(lgbSize) - …` arithmetic.

Two problems follow from the duplication:

1. **Drift** — the two parsers can diverge (already have, on error handling).
2. **Panic on a corrupt trailer** — `lgbSize` is an untrusted `uint64`. A value with the high bit set makes `int64(lgbSize)` negative, so the guarding `Seek(-20-…)` lands past EOF (a valid seek) instead of erroring, and `make([]byte, lgbSize)` then allocates a garbage-huge slice and panics. `getBaseBinarySize` similarly can return a bogus/negative base size.

### Approach

Extract one parser used by both call sites:

```go
type bundleKind int // bundleNone, bundleLegacy, bundleV2

// parseBundleTrailer reads and validates f's appended trailer.
func parseBundleTrailer(f *os.File) (kind bundleKind, lgbSize, resSize uint64, err error)
```

- Stat the file for `total` size; read the trailing 4-byte magic.
- Unrecognized magic → `bundleNone`, `err == nil` (a normal, non-bundled binary).
- `LGBX` → read the 12-byte trailer (`lgbSize`, `resSize == 0`); `LGB2` → read the 20-byte trailer (`lgbSize`, `resSize`). Check every `Seek`/`ReadFull`.
- **Size guard:** require `lgbSize + resSize + trailerLen ≤ total`, computed in `uint64` (so a crafted size can't overflow or wrap negative). On violation return a `"corrupt bundle"` error.

Rewire the callers:

- `readBundledLGB` — on `err` or `bundleNone`, return `nil, nil` (graceful; never panics). Otherwise seek to the validated offset and read; sizes are already bounded so `make` is safe.
- `getBaseBinarySize` — on `err`, propagate it (re-bundling fails cleanly); on `bundleNone`, return `total`; otherwise `total - lgbSize - resSize - trailerLen`.

### Behavior

- Valid legacy and v2 bundles parse exactly as before (no format change).
- A corrupt self-bundle makes `readBundledLGB` return `nil`, so the binary falls back to its normal CLI/REPL path instead of crashing.
- A corrupt trailer during `lg -b` re-bundling surfaces a `"corrupt bundle"` error.

### Testing

A focused unit test in `resources_e2e_test.go` (same `package main`, so it can call the unexported helper and `getBaseBinarySize`/`readBundledLGB` directly):

- Crafted file: base bytes + a `LGB2` trailer with a huge `lgbSize` (high bit set) → `readBundledLGB` returns `nil, nil` (no panic) and `getBaseBinarySize` returns a non-nil error.
- Valid `LGBX` and `LGB2` trailers parse with the expected sizes / base size.
- A non-bundled file (random bytes, no magic) → `bundleNone`, `getBaseBinarySize` returns `total`.

Existing `TestResourceBundle` and `TestLegacyBundleStillRuns` must still pass (real bundles round-trip unchanged).

---

## File Structure

- **`lg.go`** *(modify)* — add `bundleKind` + `parseBundleTrailer`; rewrite `readBundledLGB` and `getBaseBinarySize` to call it.
- **`resources_e2e_test.go`** *(modify)* — add `TestParseBundleTrailer` covering corrupt, legacy, v2, and non-bundle cases.

---

## Implementation Steps

### Task 1: Extract and validate `parseBundleTrailer`, rewire both callers

**Files:**
- Modify: `lg.go`
- Test: `resources_e2e_test.go`

- [ ] **Step 1: Write the focused test**
  Add `TestParseBundleTrailer` in `resources_e2e_test.go`. Build temp files with `os.WriteFile`:
  - base + valid `LGBX` trailer (known `lgbSize`) → `getBaseBinarySize` returns the base length; `readBundledLGB` returns the lgb bytes, `nil` res.
  - base + valid `LGB2` trailer (known `lgbSize`/`resSize`) → both sizes recovered; base length correct.
  - base + `LGB2` trailer with `lgbSize = 0xFFFFFFFFFFFFFFFF` → `readBundledLGB` returns `nil, nil` (no panic) and `getBaseBinarySize` returns a non-nil error.
  - random bytes, no magic → `getBaseBinarySize` returns the full file size, no error.

- [ ] **Step 2: Run the test to confirm it fails / panics**
  Run: `go test -timeout 120s -run TestParseBundleTrailer -count=1 .`
  Expected: FAIL (undefined `parseBundleTrailer`, or panic on the crafted-size case before the guard exists).

- [ ] **Step 3: Implement the helper and rewire**
  Add `bundleKind` (`bundleNone`/`bundleLegacy`/`bundleV2`) and `parseBundleTrailer(f *os.File) (bundleKind, lgbSize, resSize uint64, err error)` per the design: stat `total`, read trailing magic, read the 12-/20-byte trailer with all `Seek`/`ReadFull` errors checked, and enforce `lgbSize + resSize + trailerLen ≤ total` (uint64 math) returning a `"corrupt bundle"` error otherwise. Rewrite `readBundledLGB` to call it (return `nil, nil` on `err`/`bundleNone`; else seek to the validated offset and read) and `getBaseBinarySize` to call it (propagate `err`; return `total` on `bundleNone`; else `total - lgbSize - resSize - trailerLen`).

- [ ] **Step 4: Run the focused test + regressions**
  Run: `go test -timeout 200s -run 'TestParseBundleTrailer|TestResourceBundle|TestLegacyBundleStillRuns' -count=1 .`
  Expected: PASS.

- [ ] **Step 5: Verify formatting, vet, and lint**
  Run: `gofmt -l lg.go resources_e2e_test.go && go vet . && go build ./...`
  Expected: no gofmt output, vet/build clean.

- [ ] **Step 6: Commit**
  `git commit -m "refactor(lg): single validated bundle-trailer parser; guard corrupt sizes"`

---

## Implementation Summary (2026-06-03)

Done in commit `5e97654` (`lg.go`, `resources_e2e_test.go`).

- Added `bundleKind` (`bundleNone`/`bundleLegacy`/`bundleV2`) with a `trailerLen()` method, and one `parseBundleTrailer(f)` that discriminates the magic, reads the 12-/20-byte trailer with every `Seek`/`ReadFull` error checked, and validates the claimed payload against the file size.
- Rewired `readBundledLGB` (returns `nil, nil` on a corrupt or absent trailer — no panic) and `getBaseBinarySize` (propagates a `"corrupt bundle"` error, returns `total` for a non-bundle) onto the helper. Both lost their duplicated parsing.
- Size guard lives in `payloadFitsFile`, which subtracts step by step instead of summing, so the fit test can't overflow `uint64` on a huge/sparse file with crafted sizes.

**Issue encountered (caught by codex review):** the first guard summed the sizes (`lgb + res + trailerLen`), which can wrap `uint64` when each field is individually `≤ total` on a near-`MaxInt64` file. Replaced with the subtractive `payloadFitsFile` and added `TestPayloadFitsFile` covering the overflow case directly (deterministic, no huge file needed).

**Verification:** `TestParseBundleTrailer`, `TestPayloadFitsFile`, `TestResourceBundle`, and `TestLegacyBundleStillRuns` pass; `gofmt`/`go vet`/`go build` clean; `golangci-lint` reports 0 issues; final codex review found nothing.
