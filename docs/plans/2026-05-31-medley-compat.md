# medley Compatibility Implementation Plan

> **For agentic workers:** Use executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `weavejester/medley` 1.10.0 load under let-go so its namespace
compiles and the pure data/map/seq functions work, by adding targeted
Clojure-compat aliases for the JVM references medley reaches on its `:clj` and
`:default` reader-conditional branches.

**Tech Stack:** Go (let-go runtime), Clojure-dialect `.cljc` source (medley),
let-go's `installClojureCompatAliases` compat layer.

---

## Design

### Background & root cause

medley's `core.cljc` reaches JVM-only references on its reader-conditional
branches. let-go selects a branch and then must *resolve every symbol at
compile time* — because loading a namespace compiles every `defn`, even ones
the caller never invokes. Any unresolved `clojure.lang.*` / `java.*` symbol
aborts the whole namespace load, so even `m/boolean?` becomes unavailable.

Branch selection in let-go is priority-based: `:lg` > `:clj` (when matched) >
`:default`. lgx now spawns `lg` with `LG_READ_CLJ=1`, which turns on `:clj`
matching. This **exposes more JVM interop** than the older `:default`-only path,
because medley's `:clj` branches use `java.util.ArrayList`, `Throwable`,
`java.util.UUID` statics, and `java.util.regex.Pattern`.

### Scope (decided)

- **Target outcome:** medley namespace **loads**; the ~50 pure fns work.
  The interop-dependent fns must *resolve* (so the ns loads) but may throw at
  runtime. This is "load + core fns work", not "fully functional medley".
- **Handle `:clj` refs by resolving them in let-go** (compat aliases/stubs),
  preserving lgx's `LG_READ_CLJ=1` behavior. No reader-priority changes.
- **PR boundary:** separate commit(s) on branch `read-clj-minimal`, distinct
  from the existing `.clj`-resolution + `LG_READ_CLJ` commit, so this can
  become its own upstream PR.

### Already done (recovered from stash — present in working tree, uncommitted)

These two fixes are already in `pkg/rt/hierarchy.go` and `pkg/rt/lang.go`:

1. **`clojure.lang.IEditableCollection`** (medley:41–42). `hierarchy.go`: added
   `cljIEditableCollection` marker, wired into `directTypeParents` for editable
   collections only — vectors, `MapType`, `SetType` (the merged
   `MapType,SortedMapType` and `SetType,SortedSetType` cases were split because
   sorted variants have no transient form, so they must NOT report
   `IEditableCollection`). `lang.go`: added it to the `installClojureCompatAliases`
   marker list and extended the `instance?` builtin to answer marker symbols via
   the type→interface ancestry (`directTypeAncestors`).
2. **`clojure.lang.MapEntry.`** (medley:105–106). `lang.go`: aliased
   `->clojure.lang.MapEntry` and `clojure.lang.MapEntry.` to the existing
   `MapEntry` `create` fn (mirrors the `Integer.`/`->Integer` precedent).

This plan **builds on** those; it does not redo them. They get committed as
part of Task 1.

### The architecture

All new changes are **additive entries in the existing Clojure-compat layer** —
no changes to reader, resolver, or compiler core. Two files:

- `pkg/rt/lang.go` → `installClojureCompatAliases(ns)` (~line 7341): the single
  function registering JVM class aliases, constructor sugar, and static
  namespaces. Each blocker fix is one or two `ns.Def(...)` / `DefNSBare(...)`
  lines, mirroring existing entries (`java.util.UUID`, `Long/MAX_VALUE`,
  `clojure.lang.MapEntry.`).
- `pkg/rt/hierarchy.go` → already modified for `IEditableCollection`; no further
  change needed (queue/Pattern/Throwable are leaf markers, not ancestry nodes).

The one piece of *mechanism* (already recovered) is the `instance?` extension
that lets bare marker symbols answer `instance?` via type ancestry. Everything
else is lookup-table data.

**Reuse of existing let-go infrastructure** — several blockers get *real*
implementations for free, not just load-only stubs:
- let-go already has `vm.RegexType` → `java.util.regex.Pattern` aliases to it.
- let-go already has `vm.UUIDType`, `ParseUUID`, and `random-uuid`/`parse-uuid`
  builtins → `java.util.UUID/fromString` and `/randomUUID` become functional.

### Per-blocker design

`MapType`/`SetType` etc. refer to medley line numbers in
`<medley>/src/medley/core.cljc` (medley 1.10.0).

| # | medley ref (line) | Registration in `installClojureCompatAliases` | Kind |
|---|---|---|---|
| 1a | `clojure.lang.PersistentQueue` bare in `instance?` (195) | `ns.Def("clojure.lang.PersistentQueue", vm.Symbol("clojure.lang.PersistentQueue"))` — leaf marker | load-only |
| 1b | `clojure.lang.PersistentQueue/EMPTY` (188) | `DefNSBare("clojure.lang.PersistentQueue").Def("EMPTY", vm.EmptyList)` | load-only stub |
| 2 | `(java.util.ArrayList.)` / `(java.util.ArrayList. n)` (458, 525) | `ns.Def("java.util.ArrayList", marker)` + ctor sugar `ns.Def("java.util.ArrayList.", ctorStub)` + `ns.Def("->java.util.ArrayList", ctorStub)` | load-only stub |
| 3 | `Throwable` bare in `instance?` (671, 680) | `ns.Def("Throwable", vm.Symbol("Throwable"))` — leaf marker | semantically correct |
| 4a | `java.util.UUID/fromString` (695) | `DefNSBare("java.util.UUID").Def("fromString", <fn over ParseUUID>)` | **real** |
| 4b | `java.util.UUID/randomUUID` (703) | same NS `.Def("randomUUID", <existing random-uuid fn>)` | **real** |
| 5 | `java.util.regex.Pattern` bare in `instance?` (711) | `ns.Def("java.util.regex.Pattern", vm.RegexType)` | **real** |

**Key decisions:**

- **#1b EMPTY** binds to `vm.EmptyList` so `(into (queue) coll)` partially
  evaluates, but `queue?` is always `false` (nothing has the marker as an
  ancestor) — load-only semantics, documented as a stub.
- **#2 ArrayList ctor stub** returns a value sufficient for the ns to compile.
  `partition-between` / `sliding` then call `.add`/`.toArray`/`.clear`/`.size`
  on it; these throw cleanly at runtime if those fns are ever called. Acceptable
  under the chosen scope. The `.method` instance-call forms already compile in
  let-go (runtime dispatch via `vm.Receiver`), so no compiler change is needed —
  only the constructor must resolve.
- **#3 Throwable marker** → `instance?` marker path returns `false` (let-go has
  no JVM exceptions). medley's `ex-message`/`ex-cause` guard with
  `(instance? Throwable ex)`, so they return `nil` — which exactly matches
  medley's own documented "for all other types returns nil" fallback. This is
  correct behavior, not a degradation.
- **#4 UUID** reuses `vm.ParseUUID` and the existing `random-uuid` builtin →
  `m/uuid` and `m/random-uuid` genuinely work.
- **#5 Pattern → RegexType** → `m/regexp?` genuinely works against let-go
  regexes.

**Net functional outcome (better than load-only):** `uuid`, `random-uuid`,
`regexp?`, `ex-message`, `ex-cause` all work correctly. Only `queue`/`queue?`
and `partition-between`/`sliding` are degraded (load but may throw).

### Error handling

- Unresolved-symbol compile errors are the failure mode being eliminated; each
  alias removes one. The ArrayList/queue stubs deliberately let runtime calls
  fail with let-go's normal method-dispatch / arity errors rather than silently
  returning wrong data.
- No new panics introduced; all registrations are `ns.Def` of valid values.

### Testing strategy

- **Go unit tests** (`test/medley_compat_test.go`, new) in **`package test`** —
  NOT `package rt`. `pkg/compiler` imports `pkg/rt`, so a `package rt` test that
  drives the compiler would create an import cycle and fail to build; the
  existing `pkg/rt/*_test.go` files are all internal `package rt` and none of
  them run lisp. The `test/` package already imports both `pkg/compiler` and
  `pkg/rt` (see `test/language_test.go`, `test/zz_compat_test.go`). Evaluate
  snippets via:
  `compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS)).CompileMultiple(strings.NewReader(expr))`,
  asserting with testify (`github.com/stretchr/testify/assert`). Tests cover:
  marker `instance?` semantics, UUID statics round-trip, Pattern→RegexType, and
  compile-resolution smoke for every previously-blocking form. Run the medley
  end-to-end load from this package too if convenient, or via Task 5's manual
  invocation.
- **End-to-end:** `LG_READ_CLJ=1 lg -source-paths <medley>/src main.lg` prints
  `"Bool?" false`; a `verify.lg` exercises pure fns + real interop fns.
- **Regression guard:** `go test ./pkg/... ./test/` green; the
  `test/zz_compat_test.go` clojure-test-suite failing-set compared to a clean
  baseline to prove zero regressions; `go vet` + `golangci-lint`.
- **Doc-drift** (per let-go AGENTS.md "drift > silence"): verify the
  `reader.go` env-var comment names `LG_READ_CLJ` (not the dead
  `LETGO_READ_CLJ`); update lgx `docs/issues/clojure-lib-compat.md`.

### Out of scope

Real `PersistentQueue` collection type; real mutable `ArrayList` backing.
`queue`/`queue?`/`partition-between`/`sliding` remain load-only.

---

## File Structure

- **Modify** `pkg/rt/lang.go` — add blocker aliases #1–#5 inside
  `installClojureCompatAliases`; the already-recovered IEditableCollection /
  MapEntry / `instance?` changes also live here.
- **Modify** `pkg/rt/hierarchy.go` — already-recovered IEditableCollection
  wiring; no new changes expected.
- **Create** `test/medley_compat_test.go` (`package test`) — Go unit tests for
  the new aliases and marker `instance?` semantics. Placed in `test/` (not
  `pkg/rt/`) to avoid the `pkg/compiler`→`pkg/rt` import cycle; mirrors
  `test/language_test.go` / `test/zz_compat_test.go`.
- **Modify** `pkg/compiler/reader.go` — fix stale `LETGO_READ_CLJ` comment to
  `LG_READ_CLJ` if mismatched (doc-drift).
- **Modify (lgx repo)** `/Users/andrew/Projects/lgx/docs/issues/clojure-lib-compat.md`
  — record medley now loading.

Medley source path used in commands:
`/home/agent.guest/.lgx/gitlibs/github.com/weavejester/medley/1.10.0`

---

## Tasks

### Task 1: Commit the recovered IEditableCollection + MapEntry fixes

Establishes the baseline commit so subsequent blocker work is reviewable
separately. These changes are already in the working tree (stash-recovered).

**Files:**
- Modify: `pkg/rt/hierarchy.go` (already changed)
- Modify: `pkg/rt/lang.go` (already changed)

- [ ] **Step 1: Confirm the recovered changes are present and build**
  Run: `cd /Users/andrew/Projects/let-go && grep -c IEditableCollection pkg/rt/hierarchy.go pkg/rt/lang.go && go build -o /tmp/lg-medley ./`
  Expected: non-zero counts in both files; build exits 0.

- [ ] **Step 2: Confirm no regressions from the recovered changes**
  Run: `go test ./pkg/rt/... ./pkg/compiler/...`
  Expected: PASS (ok for both packages).

- [ ] **Step 3: Stage only the two source files (not mise.toml)**
  Run: `git add pkg/rt/hierarchy.go pkg/rt/lang.go`
  Note: leave the unrelated `mise.toml` working-tree change unstaged.

- [ ] **Step 4: Commit**
  `git commit -m "feat(rt): IEditableCollection + MapEntry compat for medley"`

### Task 2: PersistentQueue marker + EMPTY stub (blocker #1)

**Files:**
- Modify: `pkg/rt/lang.go`
- Test: `test/medley_compat_test.go` (new, `package test`)

- [ ] **Step 1: Write the failing test**
  Create `test/medley_compat_test.go` in `package test`. Add a small helper that
  evaluates a string and returns `(vm.Value, error)` via
  `compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS)).CompileMultiple(strings.NewReader(expr))`
  (model on `test/language_test.go`). Add a table-driven test (`TestMedleyCompat`)
  asserting `(instance? clojure.lang.PersistentQueue [1 2])` → `vm.FALSE` and that
  `clojure.lang.PersistentQueue/EMPTY` evaluates without a compile error. Use
  `github.com/stretchr/testify/assert`.

- [ ] **Step 2: Run test to verify it fails**
  Run: `go test ./test/ -run TestMedleyCompat`
  Expected: FAIL with `Can't resolve clojure.lang.PersistentQueue`.

- [ ] **Step 3: Implement**
  In `installClojureCompatAliases`, add the `clojure.lang.PersistentQueue` leaf
  marker to the marker list (alongside `IEditableCollection`), and a
  `DefNSBare("clojure.lang.PersistentQueue").Def("EMPTY", vm.EmptyList)`.

- [ ] **Step 4: Run test to verify it passes**
  Run: `go test ./test/ -run TestMedleyCompat`
  Expected: PASS.

- [ ] **Step 5: Commit**
  `git commit -am "feat(rt): PersistentQueue compat stub for medley"`

### Task 3: ArrayList constructor stub (blocker #2)

**Files:**
- Modify: `pkg/rt/lang.go`
- Test: `test/medley_compat_test.go`

- [ ] **Step 1: Write the failing test**
  Assert that `(defn f [] (java.util.ArrayList.))` and
  `(defn g [n] (java.util.ArrayList. n))` compile without error (eval returns no
  error). Do not assert runtime behavior of `.add`/`.toArray` (out of scope).

- [ ] **Step 2: Run test to verify it fails**
  Run: `go test ./test/ -run TestMedleyCompat`
  Expected: FAIL with `Can't resolve ->java.util.ArrayList`.

- [ ] **Step 3: Implement**
  Add `ns.Def("java.util.ArrayList", <marker>)`, plus constructor sugar
  `ns.Def("java.util.ArrayList.", ctorStub)` and
  `ns.Def("->java.util.ArrayList", ctorStub)`, where `ctorStub` is a NativeFn
  accepting 0 or 1 args and returning a load-only placeholder value. Document in
  a comment that this is load-only (partition-between/sliding throw at runtime).

- [ ] **Step 4: Run test to verify it passes**
  Run: `go test ./test/ -run TestMedleyCompat`
  Expected: PASS.

- [ ] **Step 5: Commit**
  `git commit -am "feat(rt): java.util.ArrayList ctor stub for medley"`

### Task 4: Throwable marker, UUID statics, Pattern alias (blockers #3, #4, #5)

**Files:**
- Modify: `pkg/rt/lang.go`
- Test: `test/medley_compat_test.go`

- [ ] **Step 1: Write the failing tests**
  Assert: `(instance? Throwable nil)` → `vm.FALSE` (compiles + correct);
  `(java.util.UUID/fromString "00000000-0000-0000-0000-000000000000")` →
  a `UUIDType` value; `(java.util.UUID/randomUUID)` → a `UUIDType` value;
  `(instance? java.util.regex.Pattern #"x")` → `vm.TRUE`.

- [ ] **Step 2: Run tests to verify they fail**
  Run: `go test ./test/ -run TestMedleyCompat`
  Expected: FAIL with `Can't resolve Throwable` / `java.util.UUID/fromString` /
  `java.util.regex.Pattern`.

- [ ] **Step 3: Implement**
  Add: `Throwable` leaf marker; `DefNSBare("java.util.UUID")` with
  `fromString` (a NativeFn over `vm.ParseUUID`, erroring on invalid input) and
  `randomUUID`. For `randomUUID`, reuse the existing `random-uuid` builtin by
  value: `ns.Lookup("random-uuid").(*vm.Var).Deref()` — do NOT reach for the
  local `randomUUID` Go var (it's out of scope inside
  `installClojureCompatAliases`); this mirrors the `Integer.`/`int` alias two
  lines below the target site. `random-uuid` is registered before
  `installClojureCompatAliases` runs, so the lookup resolves. Finally
  `ns.Def("java.util.regex.Pattern", vm.RegexType)`.
  Note: `java.util.UUID` is already `ns.Def`'d as a bare class alias to
  `vm.UUIDType`; adding a same-named `DefNSBare` namespace for the statics is
  fine and has direct precedent in `Boolean` (both `ns.Def("Boolean", ...)` and
  `DefNSBare("Boolean")` coexist).

- [ ] **Step 4: Run tests to verify they pass**
  Run: `go test ./test/ -run TestMedleyCompat`
  Expected: PASS.

- [ ] **Step 5: Commit**
  `git commit -am "feat(rt): Throwable/UUID/Pattern compat for medley"`

### Task 5: End-to-end medley load verification

**Files:**
- (none committed; uses external medley source)

- [ ] **Step 1: Rebuild lg**
  Run: `go build -o /tmp/lg-medley ./`
  Expected: exit 0.

- [ ] **Step 2: Run the with-medley example**
  Run: `cd /Users/andrew/Projects/lgx/examples/clojure-libs/with-medley && LG_READ_CLJ=1 /tmp/lg-medley -source-paths /home/agent.guest/.lgx/gitlibs/github.com/weavejester/medley/1.10.0/src main.lg`
  Expected: prints `"Bool?" false`, no load error.

- [ ] **Step 3: Run a spectrum verify script**
  Create a temp `verify.lg` requiring `[medley.core :as m]` and exercising:
  `m/map-vals`, `m/index-by`, `m/assoc-some`, `m/filter-vals`,
  `(m/uuid "…")`, `(m/random-uuid)`, `(m/regexp? #"x")`, and the marker
  `instance?` cases. Run with the same `LG_READ_CLJ=1 -source-paths` invocation.
  Expected: pure fns print correct values; uuid/random-uuid/regexp? work.

- [ ] **Step 4: Regression — full test suite**
  Run: `go test ./pkg/... ./test/`
  Expected: PASS. If `test/zz_compat_test.go` reports failures, compare against a
  clean-`HEAD` baseline (worktree) to confirm zero NEW failures.

- [ ] **Step 5: vet + lint**
  Run: `go vet ./pkg/rt/... && golangci-lint run pkg/rt/... 2>/dev/null || true`
  Expected: no new issues in changed files.

### Task 6: Doc-drift + issue doc updates

**Files:**
- Modify: `pkg/compiler/reader.go` (if comment stale)
- Modify (lgx): `/Users/andrew/Projects/lgx/docs/issues/clojure-lib-compat.md`

- [ ] **Step 1: Fix the env-var comment in reader.go**
  Verify the `matchCljConditional` doc comment names `LG_READ_CLJ`. If it still
  says `LETGO_READ_CLJ`, update it to match the implemented variable.
  Run: `grep -n "READ_CLJ" pkg/compiler/reader.go`

- [ ] **Step 2: Update the lgx issue doc**
  In lgx `docs/issues/clojure-lib-compat.md`, mark the medley
  `:default`/`instance?` items resolved and note the let-go compat aliases that
  fixed them (reference this plan).

- [ ] **Step 3: Commit (let-go)**
  `git commit -am "docs(rt): fix LG_READ_CLJ comment; note medley compat"`

- [ ] **Step 4: Commit (lgx)**
  In the lgx repo, commit the issue-doc update with a clear one-line message.
