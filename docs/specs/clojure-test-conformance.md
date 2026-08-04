---
status: planning
last-verified: 2026-09-03
authoritative-for:
  - clojure-test-api-design
  - clojure-test-tap
  - stack-trace-primitives
supersedes:
  - testing-and-conformance.md (on clojure-test-api-design — the clojure.test layer is now specified here; the older doc remains authoritative for conformance strategy)
human-verified:
---

# clojure.test Conformance and clojure.test.tap

This document specifies a `clojure.test` implementation for let-go that is faithful enough for unmodified external Clojure test harnesses to drive it, together with the `clojure.test.tap` reporter and the stack-trace primitives both depend on. It is written for a coding agent implementing the work in let-go's Clojure dialect, with Go changes limited to runtime capabilities the dialect cannot express.

## Table of Contents

1. [Overview and Goals](#1-overview-and-goals)
2. [Architecture](#2-architecture)
3. [Data Model](#3-data-model)
4. [Runtime Extensions](#4-runtime-extensions)
5. [Stack Traces](#5-stack-traces)
6. [Assertions](#6-assertions)
7. [Defining Tests](#7-defining-tests)
8. [Fixtures](#8-fixtures)
9. [Running Tests](#9-running-tests)
10. [Reporting](#10-reporting)
11. [clojure.test.tap](#11-clojuretesttap)
12. [Go Harness Integration](#12-go-harness-integration)
13. [Migration of Existing Tests](#13-migration-of-existing-tests)
14. [Out of Scope](#14-out-of-scope)
15. [Design Decision Rationale](#15-design-decision-rationale)
16. [Definition of Done](#16-definition-of-done)
17. [Appendix A: Deviations from the Oracle](#appendix-a-deviations-from-the-oracle)
18. [Appendix B: Oracle Transcripts](#appendix-b-oracle-transcripts)
19. [Appendix C: Expansion Examples](#appendix-c-expansion-examples)

---

## 1. Overview and Goals

### 1.1 What This Is

This spec rewrites let-go's `test` namespace (aliased as `clojure.test`) to match Clojure's public contract, ports `clojure/test/tap.clj` as a new `test.tap` namespace (aliased as `clojure.test.tap`), and defines a `stacktrace` namespace (aliased as `clojure.stacktrace`) backed by new trace primitives. It is written for the coding agent implementing the work and the maintainers reviewing it.

### 1.2 Problem Statement

**Status quo.** let-go's `test.lg` provides `deftest`, `is`, `testing`, `are`, fixtures, `run-tests`, `run-test-var`, and `run-test`. There is no `run-all-tests`. Assertions print `PASS`/`FAIL` lines directly. Tests are discovered via a registry populated by `register-test!` rather than through var metadata. The namespace lacks a `report` multimethod, `*test-out*`, `with-test-out`, `do-report`, `test-var`, and summary map returns. Fixtures are process-global rather than per-namespace. Runtime errors carry stack traces, but `ex-info` throws do not, and no function exposes the current stack.

**Pain.** External harnesses fail. kaocha discovers tests via `(filter (comp :test meta val) (ns-interns ns))`, reads fixtures from namespace metadata, and rebinds `do-report`. cognitect test-runner (upstream issue #738) calls `(apply run-tests nses)` and extracts `:fail` and `:error` from the returned map. `clojure.test.tap` relies on `(binding [report tap-report] ...)`, which requires a `report` var. Each harness breaks against the current namespace. The console also prints `ERROR in test:` lines with traces for conditions that are ordinary assertion failures, because library code throws bare strings and assertion machinery prints outside any reporting seam. The Go harness sees only a boolean and cannot attribute failures to specific deftests.

**Solution.** Port `clojure.test` semantically. Add the runtime capabilities the port requires. Port `clojure.test.tap`. Bridge the Go harness to the public `clojure.test` API through interop, using the `report` multimethod as the connection to Go's `testing` package.

### 1.3 Design Principles

**Oracle fidelity.** Public function names, arities, return values, var names, metadata keys, and printed output match Clojure 1.12.5 on the JVM. Every deviation from the oracle is listed in Appendix A with justification.

**One reporting seam.** Every line of test output flows through the `report` multimethod inside `with-test-out`. No assertion, runner, or library function prints directly. Rebinding `report` silences, replaces, or forwards all output.

**State lives in let-go.** Test discovery, counters, fixtures, and results are Clojure-visible vars and metadata. The Go harness reads and invokes them through interop and holds no parallel state.

**Prefer the dialect.** Behavior is implemented in `.lg` files. Go changes are limited to the runtime capabilities in Sections 4 and 5 that the dialect cannot express, plus the harness.

**Throw data, not strings.** Library code raises `ex-info` values so `thrown?` and `catch` behave as in Clojure. No assertion failure is mistaken for an uncaught exception.

**Traces are structured.** A stack frame is a map, not a formatted string. Rendering is a separate step that every consumer, including the Go harness, can replace.

### 1.4 Reference Projects

These are exemplars and oracles, not dependencies. The implementer may consult any combination.

| Name | Language | Relevance |
|---|---|---|
| [clojure/test.clj](https://github.com/clojure/clojure/blob/master/src/clj/clojure/test.clj) | Clojure | The semantic oracle. Sections 6–10 restate its definitions. |
| [clojure/test/tap.clj](https://github.com/clojure/clojure/blob/master/src/clj/clojure/test/tap.clj) | Clojure | Ported in Section 11. |
| [clojure/stacktrace.clj](https://github.com/clojure/clojure/blob/master/src/clj/clojure/stacktrace.clj) | Clojure | Ported in Section 5.4. |
| [lambdaisland/kaocha](https://github.com/lambdaisland/kaocha) | Clojure | `kaocha.type.ns` and `kaocha.type.var` show which vars and metadata keys a harness reads. |
| [cognitect-labs/test-runner](https://github.com/cognitect-labs/test-runner) | Clojure | The harness named in upstream issue #738; uses `ns-publics`, `alter-meta!` on vars, `run-tests`, and the summary map. |
| [TAP version 14](https://testanything.org/tap-version-14-specification.html) | Spec | Governs plan placement and test point syntax. |

### 1.5 Scope

This spec defines the `clojure.test` public API and its printed output, the `clojure.test.tap` public API and its printed output, the trace primitives and `clojure.stacktrace`, the runtime extensions those require, the Go harness contract including the report bridge, and the migration of existing let-go tests.

This spec excludes `clojure.test.junit`, namespace discovery from the filesystem (`clojure.tools.namespace`), and running kaocha or cognitect test-runner end to end. See Section 14.

**Relationship to the current `test.lg`.** Three recent changes invested in the registry-based design this spec replaces: #673 (register each test var once), #671 (`run-test-var` and `run-test` over the registry), and #754 (`thrown?` and `thrown-with-msg?` by hand-written expansion in `is`). Their behavior survives: a re-evaluated `deftest` still runs once (Section 7), `run-test-var` and `run-test` keep their contracts (Section 9.4), and `thrown?` keeps its pass, fail, and error semantics (Section 6.4). Their mechanisms do not: metadata discovery replaces the registry, and the `assert-expr` multimethod replaces the special-cased expansion. The typed-catch dispatch from #472 and #476 is kept as-is and relied on (Section 4.6).

### 1.6 Notation

Pseudocode follows NLSpec conventions: UPPER CASE keywords, `snake_case` names, `PascalCase` types, and `--` comments. Three domain keywords appear throughout:

- `MACRO name(params) -> Form:` defines an expansion-time function whose body builds and returns a form. Appendix C shows the resulting Clojure forms as illustrations; the pseudocode is normative.
- `METHOD multimethod dispatch_value(params):` defines one case of a multimethod.
- `BINDING var = value:` introduces a dynamic scope for the indented block, equivalent to Clojure's `binding`.

Clojure symbols are written as-is inside pseudocode identifiers when they are the public name being specified, for example `run_tests` for `run-tests`. Map access is written `m.key` for `(:key m)`, and `meta(v).test` for `(:test (meta v))`.

### 1.7 Reference Example

Examples and transcripts throughout this document use the test file `test/tap/tap-example.lg`:

```clojure
(ns my.test.tap-example
  (:require [clojure.test :refer [deftest is run-tests run-all-tests testing]]
            [clojure.test.tap :refer [with-tap-output]]))

(deftest math-test
  (testing "simple addition"
    (is (= 4 (+ 2 2)))
    (is (= 5 (+ 2 3))))
  (testing "deliberate-failure-test"
    (is (= 10 (+ 5 4)) "This assertion will fail intentionally")))
```

### 1.8 Delivery Order

The work lands as three pull requests, each independently valuable and each leaving the corpus green:

| Slice | Sections | Delivers on its own |
|---|---|---|
| 1. Traces | 3.5, 5, 4.8 | Structured traces for every thrown value, `Throwable->map`, `clojure.stacktrace`, better uncaught-error output for every user. |
| 2. Test port | 4.1 through 4.7, 6 through 11, 13 | `clojure.test` and `clojure.test.tap` at Clojure fidelity; harness in summary mode (Section 12.2). |
| 3. Bridge | 12.3 | Per-deftest Go subtests with attributed failures. |

Slice 2 depends on slice 1 for `:file` and `:line`. Slice 3 depends on slice 2 for the `report` seam. The Definition of Done subsections map onto the slices by the section numbers above.

---

## 2. Architecture

### 2.1 Layers

```text
+---------------------------------------------------------------+
|  External harnesses: kaocha, cognitect test-runner, user code |
+---------------------------------------------------------------+
|  clojure.test.tap   (pkg/rt/core/test/tap.lg)                 |
+---------------------------------------------------------------+
|  clojure.test       (pkg/rt/core/test.lg)                     |
|  clojure.stacktrace (pkg/rt/core/stacktrace.lg)               |
+---------------------------------------------------------------+
|  Runtime extensions: trace primitives, ns metadata, form      |
|  positions via meta, *file*, *out* handle coercion,           |
|  def metadata, binding helper                                 |
+---------------------------------------------------------------+
|  Go harness (test/language_test.go): invokes clojure.test     |
|  through rt.LookupVar / rt.InvokeValue, and binds `report`    |
|  to a native fn that forwards events into testing.T           |
+---------------------------------------------------------------+
```

Layers depend on those below. `clojure.test` knows nothing of TAP or Go. The Go harness does not know how tests are stored; it consumes the same `report` events any reporter does.

### 2.2 Namespaces and Aliases

| Canonical namespace | Alias | File | Notes |
|---|---|---|---|
| `test` | `clojure.test` | `pkg/rt/core/test.lg` | Alias exists in the `nsAliases` table in `pkg/rt/lang.go`. |
| `test.tap` | `clojure.test.tap` | `pkg/rt/core/test/tap.lg` | Add alias entry to `nsAliases`. Delete any existing shim namespace named `clojure.test.tap`. |
| `stacktrace` | `clojure.stacktrace` | `pkg/rt/core/stacktrace.lg` | New. |

An alias makes `(require '[clojure.test :as t])`, `(:require [clojure.test :refer :all])`, `(find-ns 'clojure.test)`, and `clojure.test/report` resolve to the canonical namespace. Vars interned in `test` are visible under both names.

All three namespaces are included in the compiled bundle and none is marked `lgbgen:skip`, since nothing in them depends on source-loading.

### 2.3 Namespace Metadata Keys

`use-fixtures` stores fixtures on namespace metadata under namespaced keywords, exactly as Clojure does:

| Key | Value | Written by | Read by |
|---|---|---|---|
| `:clojure.test/each-fixtures` | seq of fixture fns | `(use-fixtures :each ...)` | `test-vars`, kaocha |
| `:clojure.test/once-fixtures` | seq of fixture fns | `(use-fixtures :once ...)` | `test-vars`, kaocha |

The canonical namespace is `test`, so the keyword `::each-fixtures` inside `test.lg` reads as `:test/each-fixtures`. The implementation writes the keys as `:clojure.test/each-fixtures` and `:clojure.test/once-fixtures` literally, enabling kaocha's `(::t/each-fixtures ns-meta)` (with `t` aliased to `clojure.test`) to find them.

---

## 3. Data Model

### 3.1 Report Events

```
RECORD ReportEvent:
    type      : ReportType                -- dispatch value for `report`
    message   : String | None             -- the optional msg argument to `is`
    expected  : Form | None               -- the unevaluated assertion form
    actual    : Any                       -- see classification table below
    file      : String | None             -- added by do-report for FAIL and ERROR
    line      : Integer | None            -- added by do-report for FAIL and ERROR
    ns        : Namespace | None          -- BEGIN_TEST_NS / END_TEST_NS only
    var       : Var | None                -- BEGIN_TEST_VAR / END_TEST_VAR only

    -- SUMMARY events carry the counter fields directly
    test      : Integer | None
    pass      : Integer | None
    fail      : Integer | None
    error     : Integer | None
```

```
ENUM ReportType:
    PASS            -- :pass, an assertion held
    FAIL            -- :fail, an assertion did not hold
    ERROR           -- :error, an assertion or test body threw
    SUMMARY         -- :summary, end of a run-tests / run-test-var call
    BEGIN_TEST_NS   -- :begin-test-ns, test-ns is about to run a namespace
    END_TEST_NS     -- :end-test-ns, test-ns finished a namespace
    BEGIN_TEST_VAR  -- :begin-test-var, test-var is about to run a var
    END_TEST_VAR    -- :end-test-var, test-var finished a var
```

| Type | `actual` holds | Default `report` behavior |
|---|---|---|
| `PASS` | The form with subforms evaluated, e.g. `(= 4 4)`; for `assert-any`, the value | Increment `pass`. Print nothing. |
| `FAIL` | `(not <evaluated form>)`; for `assert-any`, the falsy value; for `thrown?`, `nil` | Increment `fail`. Print the FAIL block (Section 10.4). |
| `ERROR` | The thrown value | Increment `error`. Print the ERROR block (Section 10.4). |
| `SUMMARY` | n/a | Print "Ran N tests containing M assertions." and "F failures, E errors." |
| `BEGIN_TEST_NS` | n/a | Print a blank line and "Testing <ns-name>". |
| `END_TEST_NS` | n/a | Print nothing. |
| `BEGIN_TEST_VAR` | n/a | Print nothing. |
| `END_TEST_VAR` | n/a | Print nothing. |
| any other keyword | n/a | Print the event map via `prn`. |

**Key floor.** The keys above are Clojure's and are the minimum every event carries with Clojure's meaning. An implementation may add keys, for example `column` on `FAIL` and `ERROR`, and never removes or renames one. Every `report` method ignores keys it does not know, which is what lets a harness such as kaocha merge its own keys through a replacement `do-report` without breaking the default methods.

### 3.2 Summary

```
RECORD Summary:
    test   : Integer     -- deftest vars invoked
    pass   : Integer     -- passing assertions
    fail   : Integer     -- failing assertions
    error  : Integer     -- assertions or test bodies that threw
    type   : Keyword     -- always :summary
```

`(successful? summary)` is true when `fail` and `error` are both zero. A missing key counts as zero.

### 3.3 Dynamic Vars

Every var below is `^:dynamic` and interned in `test`. The Default column is the root value.

| Var | Type | Default | Description |
|---|---|---|---|
| `*load-tests*` | Boolean | `true` | When false, `deftest`, `deftest-`, `with-test`, `set-test` expand to nothing (or to the bare definition). |
| `*stack-trace-depth*` | Integer or nil | `nil` | Maximum trace frames rendered for an ERROR. `nil` renders everything. |
| `*report-counters*` | ref of map, or nil | `nil` | Bound by `test-ns` and `run-test-var`. `inc-report-counter` is a no-op while nil. |
| `*initial-report-counters*` | map | `{:test 0 :pass 0 :fail 0 :error 0}` | Seed for `*report-counters*`. |
| `*testing-vars*` | list of Var | `()` | Innermost var first. `conj` prepends. |
| `*testing-contexts*` | list of String | `()` | Innermost context first. `conj` prepends. |
| `*test-out*` | writer handle | the root binding of `*out*` when runtime bootstrap completes | See Section 10.2. |
| `report` | multimethod | the default methods | Rebound by reporters such as `with-tap-output` and the Go bridge. |
| `test-var` | fn | the default | Rebindable so harnesses can wrap test execution. |

`*report-counters*` is a ref. let-go defines `ref` as `atom`, `alter` and `commute` as `swap!`, and `dosync` as `do`; the Clojure source therefore runs unchanged. The vars `*test-result*`, `*registered-tests*`, `*each-fixtures*`, and `*once-fixtures*` are removed.

### 3.4 Var Metadata

| Key | Set by | Meaning |
|---|---|---|
| `:test` | `deftest`, `deftest-`, `with-test`, `set-test` | A zero-arity fn holding the test body. Its presence is what makes a var a test. |
| `:ns` | `def` | The namespace object the var is interned in. |
| `:name` | `def` | The var's symbol. |
| `:file`, `:line`, `:column` | `def` | Source position of the defining form. |
| `:private` | `deftest-` | Excludes the var from `ns-publics`. |

Harnesses group by `(comp :ns meta)` and read `:file`/`:line` from var metadata, so `def` attaches these keys (Section 4.4).

### 3.5 Stack Frames

```
RECORD Frame:
    fn      : String                -- qualified function name, e.g. "my.test.tap-example/math-test"
    kind    : FrameKind
    file    : String | None         -- source path for LG frames; Go file for GO frames
    line    : Integer | None
    column  : Integer | None        -- LG frames only

ENUM FrameKind:
    LG       -- :lg, a let-go function compiled from .lg source
    NATIVE   -- :native, a Go-implemented primitive called from let-go; file/line are the let-go call site
    GO       -- :go, a Go function inside a native primitive, present only when a Go panic was recovered
```

| Kind | Meaning |
|---|---|
| `LG` | Rendered as `at <fn> (<file>:<line>:<column>)`. Counted toward `*stack-trace-depth*`. |
| `NATIVE` | Rendered as `at <fn> (<file>:<line>:<column>)` using the let-go call site, matching the runtime's behavior for `native fn`. Counted. |
| `GO` | Rendered as `at <fn> (<file>:<line>) [go]`. Counted. Omitted entirely when a reporter requests `lg_only`. |

A trace is a vector of `Frame`, innermost first.

---

## 4. Runtime Extensions

These are the Go-side capabilities the ported namespaces require and the runtime lacks (September 2026), verified against the working tree. Each is a small, general-purpose feature, not test-specific. Trace primitives are in Section 5.

### 4.1 Namespace Metadata

`(meta ns-obj)` returns the namespace's metadata map. `(alter-meta! ns-obj f & args)` updates it. Currently `alter-meta!` raises "expected Atom or Var" for a namespace. Since `*ns*` inside a file is the namespace object, `(alter-meta! *ns* assoc k v)` works without further lookup.

### 4.2 Form Source Positions

The reader already records a `SourceInfo` for every list, vector, map, and `#()` literal it produces, nested forms included, in the `vm.FormSource` side table keyed by form identity (`readList` in `pkg/compiler/reader.go`). The compiler and macroexpander copy entries onto rewritten and expanded forms, and `compileForm` emits them per instruction. What is missing is the Clojure-facing surface: `(meta form)` does not consult the table and `&form` inside a macro is nil.

Three readers are added, none of which changes the reader or the bundle format:

1. `meta` on a `*List` or `*Cons` merges `{:line L :column C}` from the form's `FormSource` entry when one exists.
2. Macro invocation passes the call form so `&form` sees that metadata.
3. `def` copies `:file`, `:line`, and `:column` from the def form's entry into var metadata (Section 4.4).

Granularity is per form. An `ERROR` resolves to the throwing subform because that instruction carries its own entry; a `FAIL` resolves to the `is` form because expansion inherits the `is` form's entry. `do-report` (Section 10.1) reads the stack first and var metadata second; neither path needs `&form`.

### 4.3 `*file*`

A dynamic var `*file*` in `clojure.core` is bound to the path being loaded during `load`, `require`, and the CLI file runner, and to `"NO_SOURCE_FILE"` at the REPL. The compiler already knows the path via `SetSource`; this exposes it.

### 4.4 `def` Metadata

`def` attaches `:ns`, `:name`, and, when the def form has a `FormSource` entry, `:file`, `:line`, and `:column` to the var.

### 4.5 `*out*` Handle Coercion

The `*out*` resolver in `pkg/rt/iort.go` accepts an `IOHandle` or a raw `os.File` and falls back to stdout for anything else. It also accepts the values `io/buffer` and `io/writer` produce, so `(binding [*out* (io/buffer)] (println "x"))` writes into the buffer. Currently the text goes to the terminal and the buffer stays empty. This seam lets a test bind `*test-out*` to a buffer and read it with `io/buffer-str`.

### 4.6 Typed Catch Dispatch (existing)

No change. `catch-matches?` in `pkg/rt/exceptions.go` (#476) already dispatches `catch` clauses through the `ExceptionClass` hierarchy in `pkg/vm/exception_class.go` (#472): `Throwable` is the bottom and matches every thrown value, strings included, while every other class matches by identity or registered ancestry. `try-expr` and `test-var` therefore use `catch Throwable`, exactly as Clojure does, and no thrown value escapes an `is`. `(catch Exception e ...)` keeps its current meaning and does not see non-exception throws. Wrapped Go runtime errors carry the class `java.lang.Exception` and no finer class (Section 14, Appendix A).

### 4.7 `ex-message` on Non-Exceptions

`(ex-message v)` returns `(str v)` for a thrown non-exception value, so `thrown-with-msg?` can match a message regardless of how the value was thrown. This is a let-go-only extension listed in Appendix A.

### 4.8 Dynamic Binding from Go

Go code can push a binding of a let-go var for the duration of a call. `with-out-str*` already does this via `PushBinding` and a deferred `PopBinding` on the caller's `ExecContext`. Expose the pair as `rt.WithBinding(ec, v, value, fn)` so the harness bridge (Section 12.3) avoids reaching into `ExecContext` internals. The helper pops the binding on every exit path, including panics.

---

## 5. Stack Traces

### 5.1 Existing Machinery

The runtime records source positions per instruction (`CodeChunk.LookupSource` over `SourceInfo` in `pkg/vm/source.go`), chains `ExecutionError` values with a source per call, converts that chain into a `:trace` list of strings under `ex-data` for runtime errors (`errorToValue` in `pkg/vm/errors.go`), captures the Go stack when a Go panic is recovered (`GoPanicError`), and prints a `stack trace:` block for uncaught errors at the top level. This section unifies those facilities and exposes them to Lisp without adding a second trace mechanism.

### 5.2 Capture Rules

1. Every `throw`, of any value, captures the let-go stack at the throw site and associates it with the thrown value. Currently only runtime errors raised from Go carry a trace; `(throw (ex-info ...))` and `(throw "s")` carry none. Capture walks the live frame chain before unwinding and records an immutable `(chunk, ip, fn-name)` triple per frame; `SourceInfo` is resolved through `LookupSource` only when the trace is first read. `Frame` objects are pooled and reused after unwind, so a capture never retains a `Frame` pointer.
2. A runtime error raised inside a native primitive contributes a `NATIVE` frame whose `file`, `line`, and `column` are the let-go call site, matching `errorToValue` today.
3. A recovered Go panic contributes `GO` frames from the captured Go stack (filtered to frames outside the Go runtime and VM dispatch loop), followed by the `NATIVE` frame of the primitive and the `LG` frames above it.
4. A trace crossing the Go boundary multiple times (Lisp calls native, which invokes Lisp, which throws) forms a single vector with frames in call order, innermost first. Each boundary crossing contributes a `NATIVE` frame.
5. A thrown value that is rethrown by a `catch` keeps its original trace.
6. A `catch` that throws a new value gives the new value a new trace; the caught value is reachable via `ex-cause` when the new value is an `ex-info` constructed with it as cause.

### 5.3 Primitives

```
FUNCTION current_stack_trace() -> Vector<Frame>:
    -- Frames of the caller, innermost first. The frame for
    -- current_stack_trace itself is excluded.

FUNCTION ex_trace(v : Any) -> Vector<Frame> | None:
    -- The trace captured when v was thrown. None if v was never thrown.
    -- Works for ex-info values, runtime errors, and non-exception values.

FUNCTION Throwable_to_map(v : Any) -> Map:
    -- Clojure's Throwable->map shape.
    RETURN {
        cause : ex_message(root_cause(v)),
        data  : ex_data(root_cause(v)),          -- omitted when nil
        via   : [ {type: class_name(x), message: ex_message(x),
                   data: ex_data(x), at: first(ex_trace(x))}
                  FOR EACH x IN cause_chain(v) ],
        trace : ex_trace(v)
    }
```

`current-stack-trace`, `ex-trace`, and `Throwable->map` are interned in `clojure.core`. `ex-data` continues to expose `:trace` for runtime errors as today, now as a vector of `Frame` maps rather than strings; the string form was never documented and has no known consumer.

### 5.4 `clojure.stacktrace`

```
FUNCTION root_cause(v) -> Any:
    WHILE ex_cause(v) IS NOT None:
        v = ex_cause(v)
    RETURN v

FUNCTION print_trace_element(frame : Frame):
    CASE frame.kind:
        LG, NATIVE: print "<fn> (<file>:<line>:<column>)"
        GO:         print "<fn> (<file>:<line>) [go]"

FUNCTION print_throwable(v):
    print class_name(v) ": " ex_message(v)
    IF ex_data(v) IS NOT None: print " " pr_str(ex_data(v))

FUNCTION print_stack_trace(v, n : Integer | None = None):
    -- Step 1: print the header
    print_throwable(v)
    -- Step 2: print frames, depth-limited
    frames = ex_trace(v) OR []
    IF n IS NOT None: frames = take(n, frames)
    FOR EACH f IN frames:
        print " at "; print_trace_element(f); newline

FUNCTION print_cause_trace(v, n : Integer | None = None):
    print_stack_trace(v, n)
    c = ex_cause(v)
    WHILE c IS NOT None:
        print "Caused by: "
        print_stack_trace(c, n)
        c = ex_cause(c)

-- Behavior:
--   - Output goes to *out*. Callers wrap in with-test-out.
--   - A value with no trace prints only the header line.
--   - n = nil prints all frames.
```

`e` (print the root cause of `*e`) is omitted; let-go's REPL does not bind `*e`.

---

## 6. Assertions

### 6.1 `is`

```
MACRO is(form) -> Form:
    RETURN is_with_message(form, None)

MACRO is(form, msg) -> Form:
    RETURN try_expr(msg, form)

MACRO try_expr(msg, form) -> Form:
    RETURN a try form that:
        evaluates assert_expr(msg, form)
        CATCH Throwable t:
            do_report({type: ERROR, message: msg, expected: quote(form), actual: t})
```

`is` returns the value of the inner expression on success and `nil` or the thrown value in special forms, exactly as the expansions prescribe. `is` never prints. See Appendix C.1.

### 6.2 `assert-expr` Dispatch

`assert-expr` is a public multimethod:

```
FUNCTION assert_expr_dispatch(msg, form) -> DispatchValue:
    IF form IS None:            RETURN :always-fail
    IF form IS a seq:           RETURN first(form)        -- a symbol
    RETURN :default
```

| Dispatch value | Expansion |
|---|---|
| `:always-fail` | `do_report({type: FAIL, message: msg})` |
| `:default` | `assert_predicate` when `function?(first(form))`, else `assert_any` |
| `'instance?` | Evaluate class and object; `actual` is `(class object)` |
| `'thrown?` | Section 6.4 |
| `'thrown-with-msg?` | Section 6.4 |

Users may add methods.

### 6.3 `assert-predicate` and `assert-any`

```
FUNCTION assert_predicate(msg, form) -> Form:
    pred = first(form)
    args = rest(form)
    RETURN a form that, when evaluated:
        values = evaluate each of args, in order, into a list
        result = apply(pred, values)
        IF result:
            do_report({type: PASS, message: msg, expected: quote(form),
                       actual: cons(quote(pred), values)})
        ELSE:
            do_report({type: FAIL, message: msg, expected: quote(form),
                       actual: list(quote(not), cons(quote(pred), values))})
        RETURN result

FUNCTION assert_any(msg, form) -> Form:
    RETURN a form that, when evaluated:
        value = evaluate form
        IF value:
            do_report({type: PASS, message: msg, expected: quote(form), actual: value})
        ELSE:
            do_report({type: FAIL, message: msg, expected: quote(form), actual: value})
        RETURN value

FUNCTION function?(x) -> Boolean:
    IF x IS a symbol:
        v = resolve(x)
        IF v IS None: RETURN false
        value = get_possibly_unbound_var(v)
        RETURN value IS NOT None AND fn?(value) AND NOT meta(v).macro
    RETURN fn?(x)
```

`(= 10 (+ 5 4))` reports `actual (not (= 10 9))`; `(is (some-macro ...))` reports the raw value. See Appendix C.2.

### 6.4 `thrown?` and `thrown-with-msg?`

```
FUNCTION assert_thrown(msg, form) -> Form:
    klass = second(form)
    body  = nthnext(form, 2)
    RETURN a try form that:
        evaluates body
        do_report({type: FAIL, message: msg, expected: quote(form), actual: None})
        CATCH klass e:
            do_report({type: PASS, message: msg, expected: quote(form), actual: e})
            RETURN e

FUNCTION assert_thrown_with_msg(msg, form) -> Form:
    klass = nth(form, 1)
    re    = nth(form, 2)
    body  = nthnext(form, 3)
    RETURN a try form that:
        evaluates body
        do_report({type: FAIL, message: msg, expected: quote(form), actual: None})
        CATCH klass e:
            IF re_find(re, ex_message(e)):
                do_report({type: PASS, message: msg, expected: quote(form), actual: e})
            ELSE:
                do_report({type: FAIL, message: msg, expected: quote(form), actual: e})
            RETURN e

-- Behavior:
--   - Body returns normally: FAIL with actual nil.
--   - Body throws a matching class: PASS. The thrown value is returned.
--   - Body throws a non-matching class: the value propagates to try_expr and becomes ERROR.
--   - Throwable matches every thrown value; Exception matches exception values only (Section 4.6).
```

### 6.5 `are`

```
MACRO are(argv, expr, rest_args) -> Form:
    n_argv = count(argv)
    n_args = count(rest_args)
    IF n_argv == 0 AND n_args == 0:
        RETURN do_template(argv, is_form(expr))
    IF n_argv > 0 AND n_args > 0 AND n_args MOD n_argv == 0:
        RETURN do_template(argv, is_form(expr), rest_args)
    RAISE ex_info("The number of args doesn't match are's argv.", {})
```

`do-template` lives in `test` (Clojure keeps it in `clojure.template`; see Appendix A).

### 6.6 `testing`

```
MACRO testing(string, body) -> Form:
    RETURN a form that:
        BINDING *testing-contexts* = conj(*testing-contexts*, string):
            evaluates body
```

`testing` prints nothing. `testing-contexts-str` joins the contexts, outermost first, with a single space: `"Arithmetic with positive integers"`.

---

## 7. Defining Tests

```
MACRO deftest(name, body) -> Form:
    IF NOT *load-tests*: RETURN None
    test_fn = fn_form([], body)
    RETURN def_form(with_meta(name, {test: test_fn}),
                    fn_form([], call(test_var, var_form(name))))

MACRO deftest_private(name, body) -> Form:          -- deftest-
    same as deftest with {test: test_fn, private: true}

MACRO with_test(definition, body) -> Form:
    IF NOT *load-tests*: RETURN definition
    RETURN a form that:
        v = evaluate definition                    -- a var
        alter_meta!(v, assoc, :test, fn_form([], body))
        RETURN v

MACRO set_test(name, body) -> Form:
    IF NOT *load-tests*: RETURN None
    RETURN call(alter_meta!, var_form(name), assoc, :test, fn_form([], body))

-- Behavior:
--   - A test is any var whose metadata contains a fn under :test. No registry exists;
--     register-test! and clear-registered-tests! are removed.
--   - Calling the var's value, (math-test), runs the test through test-var with
--     reporting. Tests compose.
--   - Re-evaluating a deftest replaces the metadata on the same var. Nothing is counted twice.
```

See Appendix C.3.

---

## 8. Fixtures

```
METHOD use_fixtures :each(fixture_type, fixtures):
    alter_meta!(*ns*, assoc, :clojure.test/each-fixtures, fixtures)

METHOD use_fixtures :once(fixture_type, fixtures):
    alter_meta!(*ns*, assoc, :clojure.test/once-fixtures, fixtures)

FUNCTION default_fixture(f):
    f()

FUNCTION compose_fixtures(f1, f2) -> Fixture:
    RETURN fn(g): f1(fn(): f2(g))

FUNCTION join_fixtures(fixtures) -> Fixture:
    RETURN reduce(compose_fixtures, default_fixture, fixtures)

-- Behavior:
--   - Fixtures are per namespace, not per process. A second use-fixtures :each call
--     in the same namespace replaces the first.
--   - join_fixtures of an empty or None collection returns default_fixture.
--   - test-ns-hook, when defined in a namespace, bypasses fixtures entirely (Section 9.3).
```

---

## 9. Running Tests

### 9.1 `test-var`

```
FUNCTION test_var(v):                                   -- ^:dynamic
    t = meta(v).test
    IF t IS None: RETURN None
    BINDING *testing-vars* = conj(*testing-vars*, v):
        do_report({type: BEGIN_TEST_VAR, var: v})
        inc_report_counter(:test)
        TRY:
            t()
        CATCH Throwable e:
            do_report({type: ERROR,
                       message: "Uncaught exception, not in assertion.",
                       expected: None, actual: e})
        do_report({type: END_TEST_VAR, var: v})
```

### 9.2 `test-vars` and `test-all-vars`

```
FUNCTION test_vars(vars):
    FOR EACH (ns, vs) IN group_by(fn(v): meta(v).ns, vars):
        once = join_fixtures(meta(ns)[:clojure.test/once-fixtures])
        each = join_fixtures(meta(ns)[:clojure.test/each-fixtures])
        once(fn():
            FOR EACH v IN vs:
                IF meta(v).test IS NOT None:
                    each(fn(): test_var(v)))

FUNCTION test_all_vars(ns):
    test_vars(vals(ns_interns(ns)))
```

Order within a namespace matches the order `ns-interns` yields, which is unspecified, as in Clojure.

### 9.3 `test-ns`

```
FUNCTION test_ns(ns) -> Map:
    BINDING *report-counters* = ref(*initial-report-counters*):
        ns_obj = the_ns(ns)
        do_report({type: BEGIN_TEST_NS, ns: ns_obj})
        hook = find_var(symbol(str(ns_name(ns_obj)), "test-ns-hook"))
        IF hook IS NOT None:
            var_get(hook)()
        ELSE:
            test_all_vars(ns_obj)
        do_report({type: END_TEST_NS, ns: ns_obj})
        RETURN deref(*report-counters*)
```

### 9.4 `run-tests`, `run-all-tests`, `run-test-var`, `run-test`

```
FUNCTION run_tests() -> Summary:
    RETURN run_tests(*ns*)

FUNCTION run_tests(namespaces...) -> Summary:
    summary = assoc(merge_with(+, map(test_ns, namespaces)...), :type, :summary)
    do_report(summary)
    RETURN summary

FUNCTION run_all_tests() -> Summary:
    RETURN run_tests(all_ns()...)

FUNCTION run_all_tests(re) -> Summary:
    matching = filter(fn(ns): re_matches(re, name(ns_name(ns))), all_ns())
    RETURN run_tests(matching...)

FUNCTION run_test_var(v) -> Summary:
    BINDING *report-counters* = ref(*initial-report-counters*):
        ns_obj = meta(v).ns
        do_report({type: BEGIN_TEST_NS, ns: ns_obj})
        test_vars([v])
        do_report({type: END_TEST_NS, ns: ns_obj})
        summary = assoc(deref(*report-counters*), :type, :summary)
    do_report(summary)
    RETURN summary

MACRO run_test(test_symbol) -> Form:
    v = resolve(test_symbol)
    IF v IS None:
        print to *err* "Unable to resolve <test_symbol> to a test function."
        RETURN None
    IF meta(v).test IS None:
        print to *err* "<test_symbol> is not a test."
        RETURN None
    RETURN call(run_test_var, v)

-- Behavior:
--   - run_tests with no arguments runs the current namespace only.
--   - run_all_tests filters with re-matches (a full match, not re-find).
--     #"my.test.*" matches my.test.tap-example; #"tap-example" does not.
--   - run_all_tests with no matching namespace returns {:type :summary}. The default
--     SUMMARY method prints "Ran nil tests", as Clojure does. successful? treats
--     missing keys as zero.
--   - All runners return the summary map. None sets a global success flag.
```

---

## 10. Reporting

### 10.1 `report` and `do-report`

`report` is a `^:dynamic` multimethod dispatching on `type`. Default methods are the rows of the table in Section 3.1.

```
FUNCTION do_report(m):
    CASE m.type:
        FAIL:   report(merge(fail_position(), m))
        ERROR:  report(merge(error_position(m.actual), m))
        ELSE:   report(m)

FUNCTION fail_position() -> {file, line}:
    -- Step 1: find the innermost frame outside the test machinery
    frames = drop_while(fn(f): internal?(f), current_stack_trace())
    IF frames IS NOT empty:
        RETURN {file: first(frames).file, line: first(frames).line}
    -- Step 2: find the var under test
    v = first(*testing-vars*)
    IF v IS NOT None AND meta(v).line IS NOT None:
        RETURN {file: meta(v).file, line: meta(v).line}
    -- Step 3: fall back to unknown
    RETURN {file: None, line: None}

FUNCTION error_position(thrown) -> {file, line}:
    frames = ex_trace(thrown)
    IF frames IS NOT None AND frames IS NOT empty:
        RETURN {file: first(frames).file, line: first(frames).line}
    RETURN fail_position()

FUNCTION internal?(f : Frame) -> Boolean:
    RETURN f.fn starts with "test/" OR f.fn starts with "clojure.test/"
        OR f.fn starts with "core/" OR f.fn starts with "clojure.core/"
```

The observable result matches Clojure: `FAIL` and `ERROR` events carry `file` and `line`; `PASS` events do not. The fallback chain terminates in `nil`, which `testing-vars-str` renders as `(:)`.

### 10.2 `*test-out*` and `with-test-out`

```
*test-out* = root_binding(*out*)         -- captured when runtime bootstrap completes

MACRO with_test_out(body) -> Form:
    RETURN a form that:
        BINDING *out* = *test-out*:
            evaluates body
```

Behavior:
- Every default `report` method and every `tap-report` method wraps its printing in `with-test-out`.
- Because `*test-out*` is captured once and not rebound by `*out*` bindings, `(with-out-str (run-tests 'x))` returns `""`. Test output goes to the original stdout, matching the oracle (Appendix B.3). To capture, bind `*test-out*` to an `io/buffer`, run, then read `io/buffer-str`. Section 4.5 makes that binding effective.
- The capture happens at the end of runtime bootstrap, after an embedder has installed its host writer as the root binding of `*out*` (the WASM entry in `pkg/rt/wasm/rendermain.go` does this), not when `test.lg` loads from the bundle. Capturing at bundle load would point `*test-out*` at a process stdout the embedder has replaced, and every report would vanish. The runtime's last init hook sets the root; `(identical? *test-out* *out*)` at a fresh REPL is still `true`.

### 10.3 Helpers

```
FUNCTION inc_report_counter(name):
    IF *report-counters* IS NOT None:
        dosync(commute(*report-counters*, update_in, [name], fnil(inc, 0)))

FUNCTION testing_vars_str(m) -> String:
    names = reverse(map(fn(v): meta(v).name, *testing-vars*))
    RETURN str("(", names, ") (", m.file, ":", m.line, ")")
    -- Examples: "(math-test) (test/tap/tap-example.lg:9)" or "(math-test) (:)"

FUNCTION testing_contexts_str() -> String:
    RETURN join(" ", reverse(*testing-contexts*))

FUNCTION get_possibly_unbound_var(v) -> Any | None:
    TRY: RETURN var_get(v)
    CATCH Throwable: RETURN None

FUNCTION file_position(n : Integer) -> [String | None, Integer | None]:
    -- deprecated in Clojure 1.2; kept for source compatibility
    f = nth(current_stack_trace(), n, None)
    IF f IS None: RETURN [None, None]
    RETURN [f.file, f.line]
```

### 10.4 Default Printed Output

Text is exact, including the leading newline that `println "\nFAIL in"` produces.

```text
FAIL:
  <blank line>
  FAIL in (<vars>) (<file>:<line>)
  <testing-contexts-str, only when *testing-contexts* is non-empty>
  <message, only when message is non-nil>
  expected: <pr-str expected>
    actual: <pr-str actual>

ERROR:
  <blank line>
  ERROR in (<vars>) (<file>:<line>)
  <contexts, message as above>
  expected: <pr-str expected>
    actual: <render_actual actual>

SUMMARY:
  <blank line>
  Ran <test> tests containing <pass+fail+error> assertions.
  <fail> failures, <error> errors.

BEGIN-TEST-NS:
  <blank line>
  Testing <ns-name>
```

```
FUNCTION render_actual(v):
    IF ex_trace(v) IS NOT None:                     -- any thrown value with a trace
        print_cause_trace(v, *stack-trace-depth*)   -- Section 5.4
    ELSE:
        prn(v)
```

**Console hygiene.** No function in `test`, `test.tap`, or `stacktrace` prints except through `report` or when user code calls it directly. Library validation failures raise `ex-info` and are reported as `ERROR` only when genuinely uncaught. A run with zero failures prints only the `Testing` and `Ran` lines. The current `ERROR in test: <e>` line from the runner is removed; the `ERROR` event is the sole path.

---

## 11. clojure.test.tap

### 11.1 Public Functions

```
FUNCTION print_tap_plan(n : Integer):
    IF n < 0: RAISE ex_info("TAP plan count cannot be negative", {n: n})
    println("1.." + n)

FUNCTION print_tap_diagnostic(data : String):
    FOR EACH line IN split(data, "\n", keep_trailing_empty = true):
        println("# " + line)                        -- printing "" yields "# "

FUNCTION print_tap_pass(msg : String):
    println("ok " + msg)

FUNCTION print_tap_fail(msg : String):
    println("not ok " + msg)
```

All four return `nil`. `print-tap-diagnostic` with `""` prints the line `"# "`. The negative-plan check is a let-go addition (Appendix A); Clojure prints `1..-1`.

### 11.2 `tap-report`

`tap-report` is a `^:dynamic` multimethod on `type`:

```
FUNCTION print_diagnostics(data):
    IF *testing-contexts* IS NOT empty:  print_tap_diagnostic(testing_contexts_str())
    IF data.message IS NOT None:         print_tap_diagnostic(data.message)
    print_tap_diagnostic("expected:" + pr_str(data.expected))
    IF data.type == PASS:
        print_tap_diagnostic("  actual:" + pr_str(data.actual))
    ELSE:
        print_tap_diagnostic("  actual:" + with_out_str(render_actual(data.actual)))

METHOD tap_report :default(data):
    with_test_out: print_tap_diagnostic(pr_str(data))

METHOD tap_report :pass(data):
    with_test_out:
        inc_report_counter(:pass)
        print_tap_pass(testing_vars_str(data))
        print_diagnostics(data)

METHOD tap_report :fail(data):
    with_test_out:
        inc_report_counter(:fail)
        print_tap_fail(testing_vars_str(data))
        print_diagnostics(data)

METHOD tap_report :error(data):
    with_test_out:
        inc_report_counter(:error)
        print_tap_fail(testing_vars_str(data))
        print_diagnostics(data)

METHOD tap_report :summary(data):
    with_test_out: print_tap_plan(data.pass + data.fail + data.error)
```

Because `BEGIN_TEST_NS`, `END_TEST_NS`, `BEGIN_TEST_VAR`, and `END_TEST_VAR` have no TAP method, they fall to `:default` and print as `# {:type :begin-test-ns, :ns <ns print form>}` lines, exactly as the oracle shows. Multi-line `actual` values (error traces) become several `#`-prefixed lines because `print-tap-diagnostic` splits on newlines.

### 11.3 `with-tap-output`

```
MACRO with_tap_output(body) -> Form:
    RETURN a form that:
        BINDING report = tap_report:
            evaluates body
```

Behavior:
- The `:summary` method prints the plan after all test points for `run-tests`, `run-all-tests`, and `run-test-var`. TAP 14 permits a trailing plan when no test point follows, and `SUMMARY` is the last event any runner emits.
- Plan count equals the number of assertions (`pass + fail + error`), not deftests.
- Test points are unnumbered, matching Clojure. A TAP harness maintains its own counter.
- Outside `with-tap-output`, `report` reverts to the default methods. Nothing is installed globally.

---

## 12. Go Harness Integration

The Go harness in `test/language_test.go` runs every `.lg` file under `test/` in one shared runtime. It currently compiles `(clear-registered-tests!)` and `(run-tests)` as strings and reads `*test-result*`. Those vars no longer exist. The harness is delivered in two steps: a summary-based run first, then the report bridge. Both share the same discovery and invocation shape; only the `report` binding differs.

Discovery follows Clojure's own `(run-tests)` with no arguments: after a file loads, the harness runs the namespace that is current at the end of the file, which is what `clojure.main -i file.clj -e "(run-tests)"` does. A file that defines several test namespaces gets its last one tested; the others are reachable only through an explicit `(run-tests 'a 'b)` in the file itself. The harness keeps no namespace-to-file map and diffs no namespace sets.

### 12.1 Configuration

| Key | Type | Default | Description |
|---|---|---|---|
| `LG_TEST_REPORTER` | `bridge` or `summary` | `bridge` once Section 12.3 lands, `summary` before | Selects how the harness consumes results. `summary` exists for bisecting bridge bugs. |

**Resolution precedence** (highest first):

1. `LG_TEST_REPORTER` environment variable
2. Default from the table above

### 12.2 Step One: Summary-Based Run

```
FUNCTION run_file_tests(path) -> Boolean:
    -- Step 1: load the file (defines deftest vars, leaves its last ns current)
    load_file(path)

    -- Step 2: run the current namespace through the public API
    ns      = deref(CurrentNS)
    summary = invoke(test, "run-tests", ns)            -- rt.LookupVar + rt.InvokeValue
    RETURN invoke(test, "successful?", summary) == true

-- Behavior:
--   - The harness compiles no source strings.
--   - The harness holds no test state. Discovery, counters, and results are the
--     let-go vars of Section 3.
--   - A file with no ns form leaves core current; core has no :test vars, so the
--     run is empty and the harness fails the file with "no test namespace".
--   - Per-file dynamic-binding isolation (vm.RunWithBindings around the run) is unchanged.
--   - Output goes through *test-out*, which is stdout, so `go test -v` shows the same text
--     a user sees.
```

### 12.3 Fast Follow: The Report Bridge

The bridge binds `report` to a native fn for the duration of Step 2, so each `ReportEvent` becomes a call on Go's `testing.T`. The let-go side is unchanged: the bridge is just another reporter, like `tap-report`.

`testing.T.Run` runs its body on a new goroutine and blocks the caller until that body returns. The let-go runner therefore cannot call `t.Run` itself, or it would block inside its own `BEGIN_TEST_VAR` and never emit `END_TEST_VAR`. The bridge inverts control: `run-tests` executes on its own goroutine and hands every event across an unbuffered channel to the harness goroutine, which owns `t` and calls `t.Run` from there. The unbuffered send is a rendezvous, so let-go does not advance until Go has consumed the event, and ordering and attribution hold.

```
FUNCTION bridge_report(events : Channel<ReportEvent>) -> NativeFn:
    RETURN fn(event):
        CASE event.type:
            PASS:  inc_report_counter(:pass)      -- counters stay in let-go
            FAIL:  inc_report_counter(:fail)
            ERROR: inc_report_counter(:error)
        send(events, event)                       -- unbuffered: blocks until consumed

FUNCTION consume_events(t : testing.T, events : Channel<ReportEvent>):
    FOR EACH event IN events:                     -- until the channel closes
        CASE event.type:
            BEGIN_TEST_NS:
                t.Logf("Testing %s", ns_name(event.ns))
            BEGIN_TEST_VAR:
                name = qualified_name(event.var)
                t.Run(name, fn(sub): consume_var(sub, events))   -- blocks here until END_TEST_VAR
            FAIL, ERROR:                          -- an `is` outside any deftest
                t.Errorf(event_text(event))
            END_TEST_NS, SUMMARY:
                t.Logf(default_text(event))       -- same text Section 10.4 prints
            ELSE:
                t.Logf("%s", pr_str(event))

FUNCTION consume_var(sub : testing.T, events : Channel<ReportEvent>):
    FOR EACH event IN events:
        CASE event.type:
            FAIL:          sub.Errorf(fail_text(event))
            ERROR:         sub.Errorf(error_text(event))
            END_TEST_VAR:  RETURN                 -- subtest body returns; t.Run unblocks
            PASS:          nothing
            ELSE:          sub.Logf("%s", pr_str(event))

FUNCTION event_text(event) -> String:
    IF event.type == ERROR: RETURN error_text(event)
    RETURN fail_text(event)

FUNCTION fail_text(event) -> String:
    RETURN testing_vars_str(event) + "\n"
         + (testing_contexts_str() + "\n" IF *testing-contexts* IS NOT empty ELSE "")
         + (event.message + "\n" IF event.message IS NOT None ELSE "")
         + "expected: " + pr_str(event.expected) + "\n"
         + "  actual: " + pr_str(event.actual)

FUNCTION error_text(event) -> String:
    RETURN testing_vars_str(event) + "\n"
         + "expected: " + pr_str(event.expected) + "\n"
         + "  actual: " + with_out_str(render_actual(event.actual))

FUNCTION run_file_tests_bridged(t, path) -> Boolean:
    load_file(path)
    ns     = deref(CurrentNS)
    events = new unbuffered Channel<ReportEvent>
    result = new Channel<Summary>

    -- Step 1: let-go runs on its own goroutine
    GO:
        WITH rt.WithBinding(report_var, bridge_report(events)):
            summary = invoke(test, "run-tests", ns)
        close(events)
        send(result, summary)

    -- Step 2: the harness goroutine owns t and drives subtests
    consume_events(t, events)
    summary = receive(result)
    RETURN invoke(test, "successful?", summary) == true

-- Behavior:
--   - Each deftest is a Go subtest named <ns>/<var>, selectable with -run.
--   - The bridge increments counters itself, exactly as default methods and
--     tap-report do, since binding `report` replaces the methods that would
--     otherwise count. successful? on the returned summary remains the authoritative
--     pass/fail result; the Errorf calls agree by construction.
--   - The bridge does not wrap printing in with-test-out. testing.T owns the output.
--   - The let-go goroutine holds the ExecContext for the whole run; the harness
--     goroutine never touches let-go state. A panic on the let-go goroutine closes
--     the channel and is re-raised on the harness goroutine after consume_events returns.
--   - GO frames in an error trace show their Go file and line, making a crash inside
--     a native primitive attributable from the Go test output.
```

### 12.4 Security Considerations

The harness reads one environment variable, `LG_TEST_REPORTER`, only inside the Go test binary, and accepts two literal values. Error output includes source paths, which the top-level error printer already discloses. No new input crosses a trust boundary and no path from the environment is executed or opened.

---

## 13. Migration of Existing Tests

Existing `test/*.lg` files use `deftest`, `is`, `testing`, `are`, and `use-fixtures`. They continue to work. The following changes are observable:

| Before | After | Action |
|---|---|---|
| `PASS <form>` printed per assertion | Nothing printed on pass | None. |
| `FAIL <form> - msg` on one line | Multi-line FAIL block (Section 10.4) | None. |
| `Testing: a > b` printed by `testing` | Nothing printed. Contexts appear only in failure blocks, joined by spaces. | Update tests asserting on `>` separators. |
| `(run-tests)` runs everything registered | Runs the current namespace | Harness change (Section 12). Scripts wanting everything call `(run-all-tests)`. |
| `*test-result*` boolean | `(successful? summary)` | `.lg` or Go readers of `*test-result*` must switch. |
| `use-fixtures` global | Per namespace | Update files relying on fixtures leaking across namespaces. |
| `(throw (str ...))` in core library code | `(throw (ex-info ...))` | Convert the 40 core sites. Behavior under `catch Throwable` and bare `catch` is unchanged. |
| `:trace` under `ex-data` as strings | Vector of `Frame` maps | Readers of the string form must switch. |

Add the reference example as `test/tap/tap-example.lg`. Its failure line reads `(math-test) (test/tap/tap-example.lg:9)`. Remove any demonstration scripts that print TAP output without asserting from `test/`, since the harness runs every `.lg` file there.

**Patterns to avoid in new tests.** Tests capturing reporter output bind `*test-out*` to an `io/buffer` and read it with `io/buffer-str`. `with-out-str` around a runner returns `""`. Tests needing counters bind `*report-counters*` to `(ref *initial-report-counters*)` or call `run-test-var`, and increment via `inc-report-counter`. Tests define probes with `deftest` rather than attaching `:test` metadata by hand. Tests needing an integer from a regex group use `parse-long`.

Bundle regeneration follows the repository rule: after editing any `pkg/rt/core/**/*.lg`, run `make generate` and verify with `make check-generated`.

---

## 14. Out of Scope

**clojure.test.junit.** XML output for JUnit-compatible tooling. Excluded because no let-go consumer needs it yet. Extension point: a `junit.lg` that binds `report` to its own multimethod, exactly as Section 11.3 does.

**Filesystem namespace discovery.** `clojure.tools.namespace.find` (used by cognitect test-runner). Excluded because it is a separate library, not part of `clojure.test`. Extension point: a port whose `find-namespaces-in-dir` returns symbols that `run-tests` accepts unchanged.

**Running kaocha or cognitect test-runner end to end.** Both depend on `clojure.spec`, `clojure.tools.cli`, and JVM classpath machinery. This spec makes `clojure.test` present the contract they consume. Running them is the acceptance test of a later effort, tracked in upstream issue #738.

**Numbered TAP test points and YAML blocks.** Clojure emits neither. Extension point: `print-tap-pass` and `print-tap-fail` can prepend a counter from a dynamic var without changing callers.

**Refs with transactional semantics.** `ref` remains an alias of `atom`. `inc-report-counter` does a single `commute`, so an atom is observably equivalent.

**Parallel test execution under the bridge.** `t.Parallel()` inside bridged subtests. Excluded because the shared `ExecContext` and dynamic bindings are not goroutine-partitioned. Extension point: one `ExecContext` per subtest, once the runtime supports forking contexts with inherited bindings.

**Trace capture for values never thrown.** `(ex-trace (ex-info "x" {}))` on a value constructed but not thrown returns `nil`. Excluded because capturing at construction would cost every `ex-info` a stack walk. Extension point: an `ex-info` arity that takes an explicit trace.

**Specific exception classes for Go runtime errors.** `(/ 1 0)` and `(nth [] 5)` raise `java.lang.Exception`, not `ArithmeticException` or `IndexOutOfBoundsException`, so `(is (thrown? ArithmeticException (/ 1 0)))` reports `ERROR`. Excluded by the scope of #472, which tags wrapped Go errors as `Exception` deliberately, and by the preference not to model more of the JVM inside Go. `(is (thrown? Exception ...))` covers them. Extension point: the single `class:` tag in `errorToValue` (`pkg/vm/errors.go`), which could consult the error kind.

---

## 15. Design Decision Rationale

**Why discover tests by `:test` metadata instead of keeping the registry?** kaocha and cognitect test-runner read `(:test (meta var))` via `ns-interns` and `ns-publics`. A registry is invisible to them. Metadata also makes `(math-test)` composable and enables `run-test` to check it.

**Why is the TAP plan printed last rather than first?** Clojure's `tap-report :summary` prints it after the run for both `run-tests` and `run-all-tests`, since the plan counts assertions and no runner knows that count in advance. Perl's Test::More prints first only when the script declares `plan tests => N`; with `done_testing` it prints last. TAP 14 allows either. Plan-first would require a dry run and would diverge from the oracle byte-for-byte.

**Why capture `*test-out*` once instead of following `*out*` dynamically?** The oracle does, and `with-out-str` around `run-tests` returning `""` is documented Clojure behavior that harnesses rely on when binding `*test-out*` to their own writer. Following `*out*` dynamically would make `(binding [*test-out* w] (with-out-str ...))` ambiguous.

**Why capture at the end of bootstrap rather than when `test.lg` loads?** `test.lg` loads from the bundle before any embedder code runs, and embedders replace the root binding of `*out*` afterward. Capturing at bundle load would freeze a handle the embedder has already abandoned. The observable semantics are the same as Clojure's; only the moment moves.

**Why an atom for `*report-counters*` instead of a plain map with `set!`?** `inc-report-counter` derefs and `commute`s a reference; kaocha's `with-report-counters` binds `*report-counters*` to a ref and derefs it. A plain map breaks both. let-go's `ref`, `commute`, and `dosync` aliases let the Clojure source run as written.

**Why take `:file` and `:line` from the stack instead of only from `&form` metadata?** The runtime already tracks a source position per instruction and builds traces from it. The stack is the source Clojure uses, works for `ERROR` events whose throw site is inside a helper the `is` form never saw, and is what `print-stack-trace` needs anyway. Var metadata remains the second fallback step, not the primary.

**Why expose `FormSource` through `meta` rather than attaching metadata in the reader?** The side table already holds a position for every form and is propagated through macroexpansion; attaching metadata in the reader would duplicate it, touch every list allocation, and require the bundle to round-trip form metadata. Exposing the table is three small readers and no format change.

**Why structured `Frame` maps instead of the existing `"fn (file:line:col)"` strings?** Reporters need the fields separately: `do-report` reads `file` and `line`, the Go bridge distinguishes `GO` frames, `*stack-trace-depth*` counts frames. Parsing strings back apart is fragile and the string form had no documented consumer.

**Why rely on `catch Throwable` instead of adding a new catch-all rule?** The typed-catch dispatch from #476 already makes `Throwable` the bottom of the class hierarchy, matching every thrown value including strings, while keeping `Exception` typed. `try-expr` using `catch Throwable` is therefore both Clojure-exact and sufficient for "every `is` yields exactly one report event". Widening `Exception` as well would have changed the meaning of 36 existing typed-catch sites for no gain.

**Why capture trace triples rather than frame pointers or resolved positions?** Frames are pooled and reused after unwind, so pointers go stale; resolving positions eagerly pays a source lookup on every throw, including the caught-and-discarded ones. Immutable `(chunk, ip, fn)` triples are cheap to record and resolve correctly whenever they are read.

**Why bridge the Go harness through `report` instead of parsing printed output or reading counters?** Parsing output is brittle and loses the var boundary. Counters give a total but no attribution. `report` is the sole seam every event flows through, it is what `clojure.test.tap` itself uses, and binding it from Go needs no let-go change. The summary-based Step One exists only so the migration lands in two reviewable pieces.

**Why does the bridge run let-go on its own goroutine and drive `t.Run` from the harness goroutine?** `t.Run` blocks its caller until the subtest body returns, so the let-go runner cannot call it without deadlocking on its own `END_TEST_VAR`. Driving execution from Go instead, by calling `run-test-var` per var inside `t.Run`, would run `:once` fixtures per var and skip `test-ns-hook`. Moving the runner to a goroutine and handing events across a rendezvous channel keeps let-go in control of execution and Go in control of `testing.T`.

**Why delete a `clojure.test.tap` shim namespace instead of keeping it?** A shim that re-defs functions is not an alias: `(var clojure.test.tap/print-tap-pass)` and `(var test.tap/print-tap-pass)` would be different vars, and macros are not re-exported. The `nsAliases` table already solves this for `clojure.test`.

**Why keep `do-template` in `test` rather than a `template` namespace?** Nothing else in let-go needs `clojure.template`, and `are` is its sole caller. A `clojure.template` alias can be added later without moving code.

---

## 16. Definition of Done

### 16.1 Architecture (Section 2)

- [ ] `(find-ns 'clojure.test.tap)` and `(find-ns 'test.tap)` return the same namespace object
- [ ] `(find-ns 'clojure.stacktrace)` and `(find-ns 'stacktrace)` return the same namespace object
- [ ] `(:require [clojure.test.tap :refer [with-tap-output]])` resolves the macro
- [ ] No namespace named `clojure.test.tap` exists other than through the alias
- [ ] `use-fixtures :each` in a namespace puts a seq under `:clojure.test/each-fixtures` in `(meta *ns*)`

### 16.2 Data Model (Section 3)

- [ ] Every var in the Section 3.3 table exists in `test`, is dynamic, and has the stated default
- [ ] `*test-result*`, `*registered-tests*`, `*each-fixtures*`, `*once-fixtures*`, `register-test!`, `clear-registered-tests!` are gone
- [ ] `(meta (var some-deftest))` contains `:test`, `:ns`, `:name`, `:file`, `:line`
- [ ] Every `:fail` and `:error` event carries at least the Section 3.1 keys; a `report` method receiving an extra key ignores it
- [ ] A frame from `current-stack-trace` is a map with `:fn`, `:kind`, `:file`, `:line`, and `:column` for `:lg` frames

### 16.3 Runtime Extensions (Section 4)

- [ ] `(alter-meta! *ns* assoc :k 1)` succeeds and `(:k (meta *ns*))` is `1`
- [ ] `(meta (read-string "(a b)"))` contains `:line` and `:column`, and so does the nested `(+ 2 2)` inside `(read-string "(is (= 4 (+ 2 2)))")`
- [ ] Inside a macro, `(meta &form)` carries the call site's `:line`
- [ ] `*file*` is bound to the path during file load and to `"NO_SOURCE_FILE"` at the REPL
- [ ] `(binding [*out* (io/buffer)] (println "x"))` leaves stdout untouched and the buffer holding `"x\n"`
- [ ] `(try (throw "s") (catch Throwable e :caught))` returns `:caught` and `(try (throw "s") (catch Exception e :caught) (catch Throwable e :thr))` returns `:thr` (regression guard for #476)
- [ ] `test/catch_dispatch_test.lg` still passes
- [ ] `(ex-message "s")` returns `"s"` when `"s"` was the thrown value
- [ ] `rt.WithBinding` pushes a binding for the call and pops it on every exit path, including a panic

### 16.4 Stack Traces (Section 5)

- [ ] `(try (throw (ex-info "x" {})) (catch e (ex-trace e)))` is a non-empty vector whose first frame has the throw site's `:file` and `:line`
- [ ] `(try (throw "s") (catch e (ex-trace e)))` is a non-empty vector
- [ ] A trace read after the throwing frames have been reused by later calls still reports the original positions
- [ ] A runtime error inside `nth` yields a first frame of kind `:native` whose position is the let-go call site
- [ ] A recovered Go panic yields at least one `:go` frame with a Go file and line, followed by the `:native` frame
- [ ] Lisp calling native calling Lisp that throws yields one vector with the frames in call order and one `:native` frame per crossing
- [ ] A rethrown value keeps its trace; a new `ex-info` thrown from a `catch` has a new trace and `ex-cause` reaches the old value
- [ ] `(Throwable->map e)` has `:cause`, `:via`, and `:trace` keys with the Section 5.3 shapes
- [ ] `(print-stack-trace e 2)` prints the header and exactly two `at` lines
- [ ] `(print-cause-trace e)` prints a `Caused by:` section per cause
- [ ] `(ex-trace (ex-info "x" {}))` on an unthrown value is `nil`

### 16.5 Assertions (Section 6)

- [ ] `(is (= 4 (+ 2 2)))` emits one `:pass` event with `:actual (= 4 4)` and returns `true`
- [ ] `(is (= 10 (+ 5 4)) "m")` emits one `:fail` with `:actual (not (= 10 9))` and `:message "m"`, returns `false`
- [ ] `(is nil)` emits `:fail` via `:always-fail`
- [ ] `(is (thrown? Exception (throw (ex-info "x" {}))))` emits `:pass` and returns the exception
- [ ] `(is (thrown? Exception 1))` emits `:fail` with `:actual nil`
- [ ] `(is (thrown-with-msg? Exception #"oo" (throw (ex-info "boom" {}))))` emits `:pass`
- [ ] `(is (throw (ex-info "x" {})))` emits exactly one `:error` and nothing escapes
- [ ] `(is (some-macro ...))` goes through `assert-any`
- [ ] `(are [x y] (= x y) 1 1 2 2)` emits two `:pass` events; mismatched arity raises `ex-info`
- [ ] `testing` prints nothing; `(testing "A" (testing "B" (testing-contexts-str)))` is `"A B"`

### 16.6 Defining Tests (Section 7)

- [ ] `deftest` produces a var with `:test` metadata whose value fn calls `test-var` on itself
- [ ] `deftest-` marks the var `:private`
- [ ] `with-test` and `set-test` attach `:test` without changing the var's value
- [ ] With `*load-tests*` bound false, `deftest` expands to nil and defines nothing
- [ ] Re-evaluating a `deftest` does not cause the test to run twice under `run-tests`

### 16.7 Fixtures (Section 8)

- [ ] `:each` fixture runs once per test var; `:once` fixture runs once per namespace
- [ ] Fixtures in namespace A do not run for tests in namespace B
- [ ] `(join-fixtures nil)` returns a fn that just calls its argument
- [ ] A namespace defining `test-ns-hook` runs the hook and no fixtures

### 16.8 Running Tests (Section 9)

- [ ] `(run-tests)` with no args runs only `*ns*`
- [ ] `(run-tests 'a 'b)` returns `{:test .. :pass .. :fail .. :error .. :type :summary}` summed over both
- [ ] `(run-all-tests #"my\.test.*")` uses `re-matches`; `#"tap-example"` matches nothing
- [ ] `(run-test-var (var t))` binds fresh counters and returns a summary
- [ ] `(run-test t)` on a non-test var prints to `*err*` and returns nil
- [ ] `test-var` is rebindable with `binding` and the runner honors the rebinding
- [ ] `(successful? {:type :summary})` is `true`
- [ ] Event order per run: `:begin-test-ns`, then per var `:begin-test-var`, assertions, `:end-test-var`, then `:end-test-ns`, then `:summary`

### 16.9 Reporting (Section 10)

- [ ] `(binding [report (fn [m] nil)] (run-tests 'x))` prints nothing
- [ ] `:fail` and `:error` events carry `:file` and `:line`; `:pass` events do not
- [ ] A `:fail` inside a helper fn called from a deftest reports the `is` call site, not the helper's definition
- [ ] An `:error` from a throw inside a helper reports the throw site
- [ ] `(with-out-str (run-tests 'x))` returns `""` and output reaches stdout
- [ ] `(binding [*test-out* b] (run-tests 'x))` with `b` from `io/buffer` leaves stdout untouched and `(io/buffer-str b)` holds the output
- [ ] Under the WASM entry, test output reaches the host writer without any binding of `*test-out*`
- [ ] Default output for the reference example matches Appendix B.1 modulo file path
- [ ] A zero-failure run prints only the `Testing` and `Ran` lines
- [ ] An `:error` for `(throw "s")` prints `actual: "s"` with no trace lines
- [ ] An `:error` for an `ex-info` throw prints a header and `at` lines, no more than `*stack-trace-depth*` when bound

### 16.10 clojure.test.tap (Section 11)

- [ ] `(with-out-str (print-tap-diagnostic "a\nb"))` is `"# a\n# b\n"`; `""` gives `"# \n"`
- [ ] `(with-out-str (print-tap-pass "m"))` is `"ok m\n"`; fail is `"not ok m\n"`
- [ ] `(with-out-str (print-tap-plan 0))` is `"1..0\n"`; `-1` raises `ex-info`
- [ ] All four return `nil`
- [ ] `(with-tap-output (run-all-tests #"my\.test.*"))` prints Appendix B.2 modulo namespace print form and file path, then returns the summary
- [ ] The plan line is the last line for `run-tests`, `run-all-tests`, and `run-test-var`
- [ ] Plan count equals the number of `ok`/`not ok` lines
- [ ] An `:error` under TAP renders its trace as consecutive `#`-prefixed lines
- [ ] After `with-tap-output` exits, `(run-tests 'x)` prints default output

### 16.11 Go Harness (Section 12)

- [ ] `test/language_test.go` compiles no source strings; it uses `rt.LookupVar` and `rt.InvokeValue`
- [ ] After loading a file, the harness runs exactly the namespace current at the end of that file
- [ ] A file with no `ns` form fails with "no test namespace"
- [ ] `LG_TEST_REPORTER=summary` selects the summary-based run; unset selects the bridge once it lands
- [ ] Summary run: a file with failing assertions fails its Go subtest with the summary in the message
- [ ] Bridge: `run-tests` executes on a goroutine other than the one that owns `testing.T`, and the harness never deadlocks on a file with one or more deftests
- [ ] Bridge: every `deftest` is a Go subtest named `<ns>/<var>` selectable with `-run`
- [ ] Bridge: a `:fail` event produces exactly one `t.Errorf` on the matching subtest, with the Section 10.4 text
- [ ] Bridge: an `:error` event produces exactly one `t.Errorf` with the rendered trace; a `:go` frame shows its Go file and line
- [ ] Bridge: `:pass` produces no Go output; counters still reach the summary
- [ ] Bridge: harness stdout is empty; all text flows through `t.Log`/`t.Errorf`
- [ ] `go test ./test/...` passes on the migrated corpus

### 16.12 Migration (Section 13)

- [ ] `test/tap/tap-example.lg` exists and produces Appendix B.1 and B.2 under the CLI
- [ ] No non-asserting demonstration script remains under `test/`
- [ ] `make generate` and `make check-generated` are clean

### 16.13 Integration Smoke Test

```
-- Runs in the lg CLI against test/tap/tap-example.lg
load_file("test/tap/tap-example.lg")

plain = capture_test_out(fn(): run_all_tests(#"my\.test.*"))
ASSERT plain CONTAINS "\nTesting my.test.tap-example\n"
ASSERT plain CONTAINS "\nFAIL in (math-test) (test/tap/tap-example.lg:9)\n"
ASSERT plain CONTAINS "deliberate-failure-test\nThis assertion will fail intentionally\n"
ASSERT plain CONTAINS "expected: (= 10 (+ 5 4))\n  actual: (not (= 10 9))\n"
ASSERT plain ENDS WITH "\nRan 1 tests containing 3 assertions.\n1 failures, 0 errors.\n"
ASSERT return value == {:test 1 :pass 2 :fail 1 :error 0 :type :summary}

tap = capture_test_out(fn(): with_tap_output(run_all_tests(#"my\.test.*")))
lines = split(tap, "\n")
ASSERT count(lines matching /^ok /) == 2
ASSERT count(lines matching /^not ok /) == 1
ASSERT lines[first ok] == "ok (math-test) (:)"
ASSERT the not-ok line == "not ok (math-test) (test/tap/tap-example.lg:9)"
ASSERT last non-empty line == "1..3"
ASSERT lines CONTAIN "# {:type :begin-test-ns, :ns " + print_form(ns) + "}"
ASSERT lines CONTAIN "# {:type :begin-test-var, :var #'my.test.tap-example/math-test}"

-- Error path: a deftest whose helper throws through a native primitive
load_file("test/tap/tap-error-example.lg")     -- (deftest boom (is (= 1 (nth [] 5))))
err = capture_test_out(fn(): run_tests('my.test.tap-error-example))
ASSERT err CONTAINS "\nERROR in (boom) (test/tap/tap-error-example.lg:"
ASSERT err CONTAINS "  actual: java.lang.Exception: nth index out of bounds"
ASSERT err CONTAINS " at nth (test/tap/tap-error-example.lg:"
ASSERT err DOES NOT CONTAIN "ERROR in test:"

-- capture_test_out binds *test-out* to an io/buffer and returns io/buffer-str

-- Bridge (Go side)
result = go test ./test/ -run 'TestRunner/tap-example.lg' -v
ASSERT result CONTAINS "--- FAIL: TestRunner/tap-example.lg/my.test.tap-example/math-test"
ASSERT result CONTAINS "(math-test) (test/tap/tap-example.lg:9)"
ASSERT result CONTAINS "actual: (not (= 10 9))"
```

---

## Appendix A: Deviations from the Oracle

| Behavior | Clojure | let-go | Reason |
|---|---|---|---|
| `print-tap-plan` with negative n | prints `1..-1` | raises `ex-info` | A negative plan is invalid TAP. |
| `ex-message` on a non-exception | not applicable; cannot throw one | returns `(str v)` | let-go permits throwing any value. |
| `catch Throwable` | Throwables only | any thrown value, strings included | Pre-existing let-go rule (#476); `Exception` stays typed. |
| Class of a Go runtime error | `ArithmeticException`, `IndexOutOfBoundsException`, ... | `java.lang.Exception` | Scope of #472; no JVM classification of Go errors (Section 14). |
| `&form` metadata | reader-attached | merged from `FormSource` on read | Same keys; no reader or bundle change. |
| `do-template` location | `clojure.template` | `test` | No other consumer; add alias later if needed. |
| `file-position` | reads JVM stack | reads `current-stack-trace` | Same contract over let-go frames. |
| Namespace print form in `:default` diagnostics | `#object[clojure.lang.Namespace 0x.. "name"]` | let-go's namespace print form | Pointer identity is meaningless; harnesses read `:type` and name. |
| Frame shape | `[class method file line]` vector | `Frame` map with `:kind` | let-go frames have no class; the Go boundary requires a kind. |
| `clojure.stacktrace/e` | prints root cause of `*e` | absent | No `*e` in the let-go REPL. |
| `:trace` under `ex-data` | not present | vector of `Frame` | Pre-existing let-go convention made structured. |

---

## Appendix B: Oracle Transcripts

Produced with Clojure CLI 1.12.5.1654 on 2026-09-03 from the reference example (Section 1.7) saved as `my/test/tap_example.clj`. The JVM prints `tap_example.clj`. let-go prints the `.lg` path.

### B.1 `(run-all-tests #"my.test.*")`

```text

Testing my.test.tap-example

FAIL in (math-test) (tap_example.clj:9)
deliberate-failure-test
This assertion will fail intentionally
expected: (= 10 (+ 5 4))
  actual: (not (= 10 9))

Ran 1 tests containing 3 assertions.
1 failures, 0 errors.
```

Return value: `{:test 1, :pass 2, :fail 1, :error 0, :type :summary}`

### B.2 `(with-tap-output (run-all-tests #"my.test.*"))`

```text
# {:type :begin-test-ns, :ns #object[clojure.lang.Namespace 0x408613cc "my.test.tap-example"]}
# {:type :begin-test-var, :var #'my.test.tap-example/math-test}
ok (math-test) (:)
# simple addition
# expected:(= 4 (+ 2 2))
#   actual:(= 4 4)
ok (math-test) (:)
# simple addition
# expected:(= 5 (+ 2 3))
#   actual:(= 5 5)
not ok (math-test) (tap_example.clj:9)
# deliberate-failure-test
# This assertion will fail intentionally
# expected:(= 10 (+ 5 4))
#   actual:(not (= 10 9))
# {:type :end-test-var, :var #'my.test.tap-example/math-test}
# {:type :end-test-ns, :ns #object[clojure.lang.Namespace 0x408613cc "my.test.tap-example"]}
1..3
```

Return value: `{:test 1, :pass 2, :fail 1, :error 0, :type :summary}`. `(with-tap-output (run-tests 'my.test.tap-example))` produces identical output.

### B.3 Semantic probes

| Probe | Result |
|---|---|
| `(identical? *test-out* *out*)` at the REPL | `true` |
| `(with-out-str (run-tests 'my.test.tap-example))` | `""` (output went to the terminal) |
| `*report-counters*` at the root | `nil` |
| `(keys (meta *ns*))` after `(use-fixtures :each f)` | `(:clojure.test/each-fixtures)` |
| `*initial-report-counters*` | `{:test 0, :pass 0, :fail 0, :error 0}` |
| Nested `testing` contexts in a diagnostic | `# Arithmetic with positive integers` |
| An `:error` under TAP | `#   actual:#error {` followed by `:cause`, `:via`, `:trace` lines, each `#`-prefixed |
| Perl `plan tests => 2` | `1..2` first |
| Perl `done_testing` | `1..2` last |

---

## Appendix C: Expansion Examples

Illustrative only; the pseudocode in Section 6 and Section 7 is normative.

### C.1 `is` and `try-expr`

```clojure
(is (= 4 (+ 2 2)) "msg")
;; expands to
(try-expr "msg" (= 4 (+ 2 2)))
;; expands to
(try (let [values (list 4 (+ 2 2))
           result (apply = values)]
       (if result
         (do-report {:type :pass :message "msg" :expected '(= 4 (+ 2 2)) :actual (cons '= values)})
         (do-report {:type :fail :message "msg" :expected '(= 4 (+ 2 2)) :actual (list 'not (cons '= values))}))
       result)
     (catch Throwable t
       (do-report {:type :error :message "msg" :expected '(= 4 (+ 2 2)) :actual t})))
```

### C.2 `assert-any`

```clojure
(is (some-macro x))
;; expands to
(try (let [value (some-macro x)]
       (if value
         (do-report {:type :pass :message nil :expected '(some-macro x) :actual value})
         (do-report {:type :fail :message nil :expected '(some-macro x) :actual value}))
       value)
     (catch Throwable t ...))
```

### C.3 `deftest`

```clojure
(deftest math-test (is (= 4 (+ 2 2))))
;; expands to
(def ^{:test (fn [] (is (= 4 (+ 2 2))))} math-test
  (fn [] (test-var (var math-test))))
```

### C.4 `with-tap-output`

```clojure
(with-tap-output (run-tests 'my.test.tap-example))
;; expands to
(binding [report tap-report] (run-tests 'my.test.tap-example))
```
