# `*command-line-args*` Implementation Plan

> **For agentic workers:** Use executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `core/*command-line-args*` — a Clojure/Babashka-compatible var holding exactly the user's CLI arguments — so scripts and tools stop reverse-engineering them from `os/args`.

**Tech Stack:** Go (the `lg` host: `lg.go`, `pkg/rt/lang.go`, `pkg/vm`), let-go for the e2e test fixtures.

---

## Design

### Problem

`os/args` is the entire process argv, captured once at runtime init (`pkg/rt/os.go:26`, `vm.ToLetGo(os.Args)`). Its shape depends on how the program was launched:

| Invocation | `os/args` contents |
|---|---|
| `lg app.lg a b` | `["lg" "app.lg" "a" "b"]` |
| `lg -d app.lg a b` | `["lg" "-d" "app.lg" "a" "b"]` |
| `./bundled a b` | `["./bundled" "a" "b"]` |

Because the layout shifts, every consumer hand-rolls a heuristic to recover the user's args: `(drop 2 os/args)` (breaks when a flag precedes the script, or for a bundle), `.lg`-suffix sniffing (breaks if a real arg ends in `.lg`), or an `LGX_RUN` env flag plus a `--` marker (breaks on env inheritance and on a literal `--` in the user's command). All of these are guesses layered on top of a value that has already lost the needed information.

### Key insight

`lg` is the only layer that knows *authoritatively* where the script ends and the user's args begin — it is the one that ran `flag.Parse()`. After parsing, `flag.Args()[0]` is the script and `flag.Args()[1:]` are the user args; in a bundle, `os.Args[1:]` are all user args (the bundle skips flag parsing via the `checkBundledLGB` early return). Compute that slice once in `lg` and publish it as a var, and every downstream heuristic disappears.

### Decisions (locked)

1. **Name & location:** `core/*command-line-args*`, matching Clojure and Babashka so portable code reads it unmodified. `os/args` is left untouched (still the full argv).
2. **Empty value:** `nil` when the user passed no args; a seq of strings otherwise. Clojure-exact, so `(if *command-line-args* ...)` and `(when-let [a (first *command-line-args*)] ...)` behave as in Clojure.
3. **`--` handling:** none. `lg` reports the tokens after the script verbatim; a literal `--` (e.g. `app git checkout -- file`) is preserved. `lg` carries zero arg heuristics. (Migrating lgx/wtr to drop their `--`/`LGX_RUN` machinery is a separate follow-up, out of scope here.)

### The rule, stated plainly

`*command-line-args*` = the command-line positionals **after the first** (the script), verbatim.

| Invocation | `*command-line-args*` |
|---|---|
| `lg app.lg a b` | `("a" "b")` |
| `lg -source-paths . app.lg a b` | `("a" "b")` — flags before the script don't shift it |
| `lg app.lg git checkout -- file` | `("git" "checkout" "--" "file")` |
| `lg app.lg` / `./bundle` / REPL / `-e`-only | `nil` |
| `./bundle a b` | `("a" "b")` |

### Components & change points

**1. Define the var** in `pkg/rt/lang.go`, alongside `*compiling-aot*` / `*in-wasm*` (~line 6005):

```go
ns.Def("*command-line-args*", vm.NIL)
```

`nil` is the correct default for early core load, AOT compilation, and the REPL.

**2. A helper in `lg.go`** that converts a `[]string` into the var's value:

```go
func commandLineArgsValue(args []string) vm.Value {
    if len(args) == 0 {
        return vm.NIL // Clojure-exact: nil when empty
    }
    vs := make([]vm.Value, len(args))
    for i, a := range args {
        vs[i] = vm.String(a)
    }
    return vm.NewList(vs) // a seq of strings otherwise
}
```

`vm.NewList` already exists (`pkg/vm/list.go:200`).

**3. Two set-points**, both using the established pattern
`rt.CoreNS.Lookup("*command-line-args*").(*vm.Var).SetRoot(...)` (identical to how `*compiling-aot*` is flipped at `lg.go:781`):

- **Normal path:** immediately after `context := initCompiler(debug)` (`lg.go:767`), set unconditionally so it covers script / `-e` / compile / bundle / wasm modes uniformly. Value = `files[1:]` when `len(files) >= 1`, else nil. `files` is already in scope from `files := flag.Args()` (`lg.go:762`). Placing it here guarantees `CoreNS` is ready and the var is set before `runFile` / `compileLG` / `bundleBinary` / `buildWasm` execute user code.
- **Bundled path:** inside the `checkBundledLGB` early-return branch, after `ctx := initCompiler(false)` (`lg.go:702`) and before the namespace-chunk loop (`lg.go:732`). Value = `os.Args[1:]`. A bundle skips `flag.Parse`, and its top-level forms read the var at startup, so it must be set before any chunk runs.

### Why this is safe

- `os/args` is unchanged — no existing script breaks.
- The default `vm.NIL` means the var is always bound, even during core bootstrap and AOT compilation (when top-level forms run under `lg -b`).
- `CoreNS` is populated before `initCompiler` runs in both paths (the existing `*compiling-aot*` write at `lg.go:781` relies on the same guarantee).

### Known non-goal

`lg -e expr a b` already mistreats `a` as a script path (pre-existing behavior in `lg.go:824`). This plan does not change that; `*command-line-args*` simply follows the "positionals after the first" rule. Not worth special-casing.

### Testing strategy

A new `command_line_args_e2e_test.go` builds the real `lg` binary with `buildLG(t)` (defined in `scope_e2e_test.go`) and runs it against a tiny script whose entire body is `(prn *command-line-args*)`, asserting exact stdout. The bundle case mirrors `TestLegacyBundleStillRuns` (`resources_e2e_test.go:403`). This exercises both set-points (normal + bundled) and the verbatim/`nil` semantics end to end, which a unit test on the helper alone could not.

### Docs

Document `*command-line-args*` in `README.md` next to the existing `os/args` description, with the framing: "the user's args after the script; `nil` when none; `os/args` remains the full process argv."

---

## File Structure

- **Modify** `pkg/rt/lang.go` — define `*command-line-args*` var with `nil` default in the core namespace.
- **Modify** `lg.go` — add the `commandLineArgsValue` helper; set the var in the normal path (after `initCompiler`) and in the bundled path (in the `checkBundledLGB` branch).
- **Create** `command_line_args_e2e_test.go` — end-to-end tests for both run modes and all value-semantics cases.
- **Modify** `README.md` — document the new var alongside `os/args`.

---

## Task Structure

### Task 1: Define the `*command-line-args*` var

**Files:**
- Modify: `pkg/rt/lang.go`

- [x] **Step 1: Define the var with a `nil` default**
  In `installLangNS` (or wherever the core namespace `ns` is built — the block at `pkg/rt/lang.go:6002-6010` that defines `*compiling-aot*`, `*in-wasm*`, `*ansi?*`), add `ns.Def("*command-line-args*", vm.NIL)`. Use a comment noting it is set by `lg` at startup to the user's CLI args (the positionals after the script), `nil` when there are none.

- [x] **Step 2: Verify it builds and the var resolves**
  Run: `go build ./... && go run . -e '(prn *command-line-args*)'`
  Expected: builds clean; prints `nil` (no script positional, so the default stands).

- [x] **Step 3: Commit**
  `git commit -am "feat(args): define core/*command-line-args* var (nil default)"`

### Task 2: Populate the var from `lg` (both run paths)

**Files:**
- Modify: `lg.go`

- [x] **Step 1: Add the `commandLineArgsValue` helper**
  Add the helper described in the design (near the other top-level helpers in `lg.go`, e.g. by `buildSearchPaths`): `[]string` → `vm.NIL` when empty, else `vm.NewList` of `vm.String`.

- [x] **Step 2: Set the var in the normal path**
  Immediately after `context := initCompiler(debug)` (`lg.go:767`), compute `userArgs := files[1:]` when `len(files) >= 1`, else a nil slice, and set
  `rt.CoreNS.Lookup("*command-line-args*").(*vm.Var).SetRoot(commandLineArgsValue(userArgs))`.
  `files` comes from `files := flag.Args()` at `lg.go:762`.

- [x] **Step 3: Set the var in the bundled path**
  In the `checkBundledLGB` early-return branch, after `ctx := initCompiler(false)` (`lg.go:702`) and before the namespace-chunk loop (`lg.go:732`), set the same var from `os.Args[1:]`.

- [x] **Step 4: Verify both paths by hand**
  Run: `go build -o /tmp/lg . && printf '(prn *command-line-args*)' > /tmp/cla.lg`
  Then:
  - `/tmp/lg /tmp/cla.lg a b` → expect `("a" "b")`
  - `/tmp/lg /tmp/cla.lg` → expect `nil`
  - `/tmp/lg -source-paths . /tmp/cla.lg a b` → expect `("a" "b")` (a flag before the script does not shift the args)
  - `/tmp/lg /tmp/cla.lg git checkout -- file` → expect `("git" "checkout" "--" "file")`
  - `/tmp/lg -b /tmp/cla /tmp/cla.lg && /tmp/cla a b` → expect `("a" "b")`

- [x] **Step 5: Commit**
  `git commit -am "feat(args): populate *command-line-args* from lg in script and bundle modes"`

### Task 3: End-to-end tests

**Files:**
- Create: `command_line_args_e2e_test.go`
- Test: `command_line_args_e2e_test.go`

- [x] **Step 1: Write the failing tests**
  Create `command_line_args_e2e_test.go` (package `main`), modeled on `source_paths_e2e_test.go` for structure and `TestLegacyBundleStillRuns` (`resources_e2e_test.go:403`) for the bundle case. Build the binary once with `bin := buildLG(t)`. Write a fixture script `app.lg` containing exactly `(prn *command-line-args*)` into a `t.TempDir()`. Use `CombinedOutput()` and assert exact trimmed output, per subtest:
  - `app.lg a b` → `("a" "b")`
  - `-source-paths . app.lg a b` → `("a" "b")` (flag before script does not shift the args; because `.` is listed, the source-paths transition warning is *not* emitted — `CombinedOutput()` is clean, no stream separation needed)
  - `app.lg` → `nil`
  - `app.lg git checkout -- file` → `("git" "checkout" "--" "file")`
  - bundle: `-b <out> app.lg`, then run `<out> a b` → `("a" "b")`

- [x] **Step 1b: Lock in the "`os/args` unchanged" claim**
  Add one subtest with a second fixture `both.lg` printing both values (e.g. `(prn *command-line-args*) (prn os/args)`); run `both.lg a b` and assert `os/args` still ends with the full argv (`... "both.lg" "a" "b"`) while `*command-line-args*` is `("a" "b")`. Cheap guard for the design's back-compat promise.

- [x] **Step 2: Run tests to verify they fail (if implemented before Tasks 1-2) or pass**
  Run: `go test -run TestCommandLineArgs ./...`
  Expected: with Tasks 1-2 done, PASS. (If authored first, FAIL on output mismatch / `nil` printed where args expected.)

- [x] **Step 3: Make tests pass**
  If any subtest fails, reconcile the implementation from Task 2 with the expected output. No new production code should be needed beyond Tasks 1-2.

- [x] **Step 4: Run the full e2e suite**
  Run: `go test ./...`
  Expected: PASS (no regressions in existing e2e tests).

- [x] **Step 5: Commit**
  `git commit -am "test(args): e2e coverage for *command-line-args* (script, flag-before-script, nil, --, bundle)"`

### Task 4: Document the var

**Files:**
- Modify: `README.md`

- [x] **Step 1: Document `*command-line-args*`**
  In `README.md`, next to where `os/args` / the `os` surface is described, add `*command-line-args*` with the framing: it holds the user's args (the positionals after the script), is `nil` when there are none, and matches Clojure/Babashka; `os/args` remains the full process argv. Include a one-line example showing the args-after-script and bundle behaviors. Apply the /writing-clearly skill.

- [x] **Step 2: Commit**
  `git commit -am "docs(args): document *command-line-args*"`

---

## Follow-ups (out of scope)

- **lgx:** stop injecting the `--` marker and setting `LGX_RUN` in `cmd-run`; pass user args through directly now that `lg` reports them.
- **wtr:** delete `wtr.args/strip-runner-args` and read `*command-line-args*` instead.
- **let-go's own `scripts/*.lg`:** replace `(drop 2 os/args)` with `*command-line-args*` as dogfood.

---

## Implementation summary — completed 2026-06-08

All four tasks landed as planned, no deviations. Branch `cli-args`, commits:

- `f832932` — define `core/*command-line-args*` (`nil` default) in `pkg/rt/lang.go`.
- `e61d07d` — `commandLineArgsValue` helper + `setCommandLineArgs` in `lg.go`; set in the normal path (`files[1:]` after `initCompiler`) and the bundled path (`os.Args[1:]` in the `checkBundledLGB` branch).
- `72b5372` — `command_line_args_e2e_test.go`: 7 subtests (script, flag-before-script, nil, literal `--`, bundle, plus an `os/args`-unchanged guard).
- `a07845e` — README docs in the Usage section.

**Notes:**
- A single `setCommandLineArgs` wraps the lookup + `SetRoot` so both paths share one code path (small DRY improvement over the plan's inline sketch).
- Verification: full `go test ./...` green; `gofmt` clean.
- Second-opinion review (codex, branch vs `main`): no correctness issues — confirmed both paths, nil-for-empty, and verbatim semantics. Nothing to fix.
- Behavior confirmed by hand: `("a" "b")` for script/flag-before-script/bundle, `nil` for no-args, `("git" "checkout" "--" "file")` for a literal `--`.

Follow-ups above (lgx, wtr, let-go's own scripts) remain out of scope.
