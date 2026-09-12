---
status: planning
last-verified: 2026-09-10
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
17. [Executable Evidence](#17-executable-evidence)
18. [Appendix A: Deviations from the Oracle](#appendix-a-deviations-from-the-oracle)
19. [Appendix B: Oracle Transcripts](#appendix-b-oracle-transcripts)
20. [Appendix C: Expansion Examples](#appendix-c-expansion-examples)

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

**One reporting seam.** Every line produced by a test run flows through the `report` multimethod. Default and TAP reporters print inside `with-test-out`; the Go bridge replaces printing with `testing.T` calls. No assertion, runner, or library function prints directly. Rebinding `report` silences, replaces, or forwards all output.

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
| 1. Traces | 3.5, 4.8, 4.10, 5 | Structured traces from the existing unwind chain for every thrown exception, on both backends, plus `Throwable->map`, `clojure.stacktrace`, and better uncaught-error output for every user. |
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

`use-fixtures` stores fixtures on namespace metadata under auto-resolved keywords. The storage behavior matches Clojure; the canonical keyword namespace differs as documented below:

| Key | Value | Written by | Read by |
|---|---|---|---|
| `:test/each-fixtures` | seq of fixture fns | `(use-fixtures :each ...)` | `test-vars`, kaocha |
| `:test/once-fixtures` | seq of fixture fns | `(use-fixtures :once ...)` | `test-vars`, kaocha |

The canonical namespace is `test`, so the keyword `::each-fixtures` inside `test.lg` reads as `:test/each-fixtures`. The reader resolves `::t/each-fixtures` through the namespace object named by alias `t`; because `clojure.test` aliases that same canonical `test` object, external harness code also produces `:test/each-fixtures`. The implementation therefore stores the canonical `:test/*` keys. Literal JVM keys such as `:clojure.test/each-fixtures` are a documented alias-model deviation (Appendix A), not a second source of fixture state.

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
| `*test-out*` | writer handle | the registered host stdout, or the current root of `*out*` before any host replacement | See Section 10.2. |
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

**The trace is the unwound error chain.** No capture happens at the throw site and no live frame stack is walked. A `throw` returns `ThrownError{value}` as a Go `error`; every call site the error passes through on the way out wraps it in an `ExecutionError` carrying that site's function name and `SourceInfo`, which is what the bytecode VM does today at `OP_INVOKE`. The chain that arrives at a `catch` therefore already lists every frame between the throw and the catch, innermost first. The design adds three things: the lowered-Go backend wraps the same way (Section 4.10), the catch keeps the chain instead of discarding it, and `ex-trace` renders it.

```
RECORD ExInfo:                          -- existing pointer struct, two fields added
    message : String
    data    : Map
    cause   : Error | None              -- a Go error: ThrownError, ExecutionError chain, or a host error
    meta    : Value
    chain   : Error | None              -- set once at first catch; never overwritten (CAS)

FUNCTION catch_bind(err : Error) -> Value:
    -- The single seam both backends use (vm.ErrorToValue today).
    value = thrown_value(err)                       -- ThrownError.Value, or a boxed exception for a Go error
    IF value IS *ExInfo AND value.chain IS None:
        compare_and_swap(value.chain, None, err)    -- first catch wins; a reused exception keeps its first trace
    RETURN value

FUNCTION box_go_error(err : Error) -> *ExInfo:
    -- A Go error that reaches let-go is boxed with the Go error retained as cause.
    RETURN ExInfo(innermost_message(err), {}, cause = err)

FUNCTION ex_trace(v) -> Vector<Frame> | None:
    IF v IS *ExInfo AND v.chain IS NOT None: RETURN frames_of(v.chain)
    RETURN None

FUNCTION frames_of(err) -> Vector<Frame>:
    -- Walk ExecutionError links outward. A "calling <fn>" link with SourceInfo is an LG or NATIVE frame.
    -- A GoPanicError contributes its captured Go frames as GO frames. Resolution is lazy: the chain is
    -- stored, frames are built on read.
```

`chain` is a computed field. It participates in neither `=` nor `hash`, so `(= e (ex-info "x" {}))` is unaffected by whether `e` was ever thrown, and `WithMeta` copies it like any other field. Setting it once matches the JVM, where a `Throwable` captures its stack at construction and rethrowing the same object keeps that trace.

Scalars have no trace. A thrown string, keyword, number, boolean, or `nil` is a bare Go value with no field to hold a chain, so `ex-trace` on it returns `nil` and the test reporter falls back to the assertion's own source position (Section 10.1). After the Section 13 migration the only scalar throws in the repository are test fixtures.

---

## 4. Runtime Extensions

These are the Go-side capabilities the ported namespaces require and the runtime lacks (September 2026), verified against the working tree. They are general-purpose rather than test-specific. Most are small; the provenance-bearing datum path required for correct arbitrary-value traces is deliberately cross-cutting and lands separately in Slice 1. Trace primitives are in Section 5.

### 4.1 Namespace Metadata

`(meta ns-obj)` returns the namespace's metadata map. `(alter-meta! ns-obj f & args)` updates it. Currently `alter-meta!` raises "expected Atom or Var" for a namespace. Since `*ns*` inside a file is the namespace object, `(alter-meta! *ns* assoc k v)` works without further lookup.

### 4.2 Form Source Positions

The reader already records a `SourceInfo` for every identity-bearing list or cons form it produces, nested forms included, in the `vm.FormSource` side table (`readList` in `pkg/compiler/reader.go`). Calls attempting to record non-hashable vector and map values are intentionally ignored by the current table and are not part of this slice. The compiler and macroexpander copy entries onto rewritten and expanded list forms, and `compileForm` emits them per instruction. What is missing is the Clojure-facing surface: `(meta form)` does not consult the table and `&form` inside a macro is nil.

Three readers are added, none of which changes the reader or the bundle format:

1. `meta` on a `*List` or `*Cons` merges `{:line L :column C}` from the form's `FormSource` entry when one exists.
2. Macro invocation passes the call form so `&form` sees that metadata.
3. `def` copies `:file`, `:line`, and `:column` from the def form's entry into var metadata (Section 4.4).

Granularity is per supported list/cons form. An `ERROR` resolves to the throwing subform because that instruction carries its own entry; a `FAIL` resolves to the `is` form because expansion inherits the `is` form's entry. `do-report` (Section 10.1) reads the stack first and var metadata second; neither path needs `&form`.

### 4.3 `*file*`

A dynamic var `*file*` in `clojure.core` is bound to the path being loaded during `load`, `require`, and the CLI file runner, and to `"NO_SOURCE_PATH"` at the REPL, matching the Clojure 1.12.5 REPL oracle. The compiler already knows the path via `SetSource`; this exposes it.

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

### 4.9 `clojure.string/split` Conformance

Clojure's `clojure.string/split` is a thin wrapper over `java.util.regex.Pattern.split`, and `clojure.test.tap` calls `String.split` on the same path. Both are the one-argument form, which the Javadoc defines as the two-argument form with a limit of zero:

> If the limit *n* is zero then the pattern will be applied as many times as possible, the array can have any length, and trailing empty strings will be discarded.
>
> — `java.lang.String#split(String, int)`, https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/String.html#split(java.lang.String,int)

The current let-go `split` in `pkg/rt/lang.go` wraps Go's `strings.Split`/`Regexp.Split`, which keep trailing empty strings, returns a list rather than a vector, and has no limit arity. It is brought to Clojure's contract:

```
FUNCTION string_split(s : String, re : Regex, limit : Integer = 0) -> Vector<String>:
    IF s == "":            RETURN [""]              -- no match: the input itself
    IF limit > 0:          RETURN at most `limit` fields; the last holds the rest of the input
    fields = every field between matches, interior empty fields included
    IF limit == 0:         drop trailing empty fields from `fields`
    RETURN vector(fields)                            -- limit < 0 keeps every field
```

A leading empty field is kept whenever the first match is at index zero and the match is non-empty, exactly as Java does. Only the pattern form is specified; whether a string delimiter continues to be accepted is unchanged by this document.

`clojure.string/split` returns a vector and follows the limit-zero rule. [R-string-split-limit-zero]

<!-- evidence: @R-string-split-limit-zero -->
| form | value | out |
|---|---|---|
| `(str/split "a\n" #"\n")` | `["a"]` | |
| `(str/split "a\n\n" #"\n")` | `["a"]` | |
| `(str/split "\n" #"\n")` | `[]` | |
| `(str/split "" #"\n")` | `[""]` | |
| `(str/split "a\n\nb" #"\n")` | `["a" "" "b"]` | |
| `(str/split "\na" #"\n")` | `["" "a"]` | |
| `(str/split "a:b:c" #":" 2)` | `["a" "b:c"]` | |
| `(str/split "a::" #":" -1)` | `["a" "" ""]` | |
| `(vector? (str/split "a" #"\n"))` | `true` | |

### 4.10 Error Boundaries: Lowered Go and Natives

Two existing gaps make the chain incomplete, and both are closed in Slice 1.

**Lowered Go wraps at every call site.** Generated code today propagates errors with a bare `return nil, err`, so a throw inside a lowered function unwinds through it without recording a frame. The IR already attaches the full `SourceInfo` of the originating form to every instruction, calls included, but the IR bridge exposes only `source-info-symbol` to `.lg`. The bridge gains `source-info-file`, `source-info-line`, and `source-info-column`, and the call emitter in `lower_go.lg` wraps on the error path with the same `calling <fn>` shape the VM uses:

```
-- emitted for every lowered call whose result may be an error
IF err != nil:
    RETURN nil, vm.NewExecutionError("calling <fn>").WithSource(<file>, <line>, <column>).Wrap(err)
```

The happy path is unchanged, so this does not affect the bench ratchet. Positions come from the same `SourceInfo` the VM resolves through `LookupSource`, so the two backends produce the same frames for the same program. A call whose form has no `FormSource` entry records `<unknown>` as the VM does.

**Natives preserve the chain.** A native that re-wraps an error with `fmt.Errorf("...: %v", err)` flattens the chain to a string and loses the `ThrownError` inside it, so a user's `(throw (ex-info ...))` passing through that native arrives at the catch as a generic exception. Every such site becomes `%w`. A test asserts that a value thrown from a callback survives, with its trace, through `map`, `reduce`, `sort`, `apply`, and every native that invokes let-go code.

**Go errors are boxed with their cause.** `box_go_error` (Section 3.5) keeps the Go error as `cause`, so `ex-cause` on a runtime error returns a boxed Go error value whose class is the Go type name and whose message is `Error()`, and `Throwable->map` walks it. Today the box is built with a `nil` cause and the Go error is unrecoverable from let-go.

---

## 5. Stack Traces

### 5.1 Existing Machinery

The runtime records source positions per instruction (`CodeChunk.LookupSource` over `SourceInfo` in `pkg/vm/source.go`), chains `ExecutionError` values with a source per call, converts that chain into a `:trace` list of strings under `ex-data` for runtime errors (`errorToValue` in `pkg/vm/errors.go`), captures the Go stack when a Go panic is recovered (`GoPanicError`), and prints a `stack trace:` block for uncaught errors at the top level. This section unifies those facilities and exposes them to Lisp without adding a second trace mechanism.

### 5.2 Capture Rules

1. A trace is the `ExecutionError` chain that unwinds from the throw to the catch, read at the catch and stored on the caught `ExInfo` (Section 3.5). `(throw (ex-info ...))`, a runtime error raised in a native, and a recovered Go panic all produce a trace the same way. Nothing is captured at the throw site.
2. Each `OP_INVOKE` in the VM and each lowered call site (Section 4.10) contributes one `LG` frame. A native that raises an error contributes one `NATIVE` frame whose position is the let-go call site, matching `errorToValue` today.
3. A recovered Go panic contributes `GO` frames from the captured Go stack, filtered to frames outside the Go runtime and the VM dispatch loop, followed by the `NATIVE` frame of the primitive and the `LG` frames above it.
4. A trace crossing the Go boundary more than once (Lisp calls a native, which invokes Lisp, which throws) is one vector in call order, because the native returns the callback's error unchanged and its own call site wraps it. Each crossing contributes a `NATIVE` frame named for the native (`apply`, `map`, `reduce`, multimethod and protocol dispatch). These are the frames `internal?` filters in Section 10.1.
5. A tail call (`OP_TAIL_CALL`) reuses the frame, so the tail-called function does not appear between its caller and callee. A self-recursive function in tail position shows one frame, not its recursion depth. This matches `recur` on the JVM.
6. A catch stores the chain on the `ExInfo` only if it has none. Rethrowing a caught exception, `(throw e)`, starts a new `ThrownError` around the same object; the outer frames wrap it, and `ex-trace` returns the stored chain followed by the new one. Throwing a new `ex-info` from a catch gives it a new trace, and when the caught value is its cause, `ex-cause` reaches the old trace.
7. A caught exception may be stored in a collection, an atom, or a var and inspected later; the trace is part of the object.
8. A thrown scalar has no trace (Section 3.5). `ex-trace` returns `nil` and reporters fall back to form position.

### 5.3 Primitives

```
FUNCTION ex_trace(v : Any) -> Vector<Frame> | None:
    -- The frames unwound when v was thrown. None if v was never thrown or is a scalar.

FUNCTION Throwable_to_map(v : Any) -> Map:
    -- Clojure's Throwable->map key/omission rules, over let-go Frame maps.
    root = root_cause(v)
    via = []
    FOR EACH x IN cause_chain(v):
        entry = {type: class_symbol(x)}
        IF ex_message(x) IS NOT None: entry.message = ex_message(x)
        IF ex_data(x) IS NOT None:    entry.data = ex_data(x)
        IF first(ex_trace(x)) EXISTS: entry.at = first(ex_trace(x))
        append(via, entry)
    result = {via: via, trace: ex_trace(root) OR []}
    IF ex_message(root) IS NOT None: result.cause = ex_message(root)
    IF ex_data(root) IS NOT None:    result.data = ex_data(root)
    IF get(ex_data(v), :clojure.error/phase) IS NOT None:
        result.phase = get(ex_data(v), :clojure.error/phase)
    RETURN result
```

`ex-trace` and `Throwable->map` are interned in `clojure.core`. There is no `current-stack-trace`: the runtime has no live frame chain to walk (calls recurse on the Go stack, and `ExecContext` holds bindings and scope only), and nothing in `clojure.test` needs one once positions come from form source (Section 4.2). `file-position` is therefore absent (Appendix A).

`cause_chain` follows `ex-cause`, which returns the stored `cause` whether it is an `ExInfo`, a boxed Go error, or a thrown scalar.

`:trace` under `ex-data` for runtime errors is retained for compatibility but is **derived lazily from the chain** on read, as a vector of `Frame` maps rather than strings. Because it is computed from `chain`, it is excluded from equality exactly as `chain` is, and it cannot drift from `ex-trace` on rethrow. The string form was never documented and has no known consumer.

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
    IF ex_data(v) IS NOT None:
        newline
        print pr_str(ex_data(v))

FUNCTION exception_value?(v) -> Boolean:
    -- True for ExInfo and normalized runtime/Go exception values carrying an
    -- ExceptionClass. Deliberately false for arbitrary values that only
    -- catch-matches? Throwable because they were thrown.

FUNCTION print_trace_frames(v, n : Integer | None = None):
    frames = ex_trace(v) OR []
    print " at "
    IF first(frames) EXISTS:
        print_trace_element(first(frames))
    ELSE:
        print "[empty stack trace]"
    newline
    rest_frames = rest(frames)
    IF n IS NOT None: rest_frames = take(max(n - 1, 0), rest_frames)
    FOR EACH f IN rest_frames:
        print "    "; print_trace_element(f); newline

FUNCTION print_stack_trace(v, n : Integer | None = None):
    print_throwable(v)
    newline
    print_trace_frames(v, n)

FUNCTION print_cause_trace(v, n : Integer | None = None):
    print_stack_trace(v, n)
    c = ex_cause(v)
    WHILE c IS NOT None:
        print "Caused by: "
        print_stack_trace(c, n)
        c = ex_cause(c)

-- Behavior:
--   - Output goes to *out*. Callers wrap in with-test-out.
--   - A value with no trace prints the header and ` at [empty stack trace]`.
--   - n = nil prints all frames.
--   - As in Clojure 1.12.5, the first frame prints even when n <= 0; n limits
--     only how many frames after the first are printed.
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
| `:default` | `assert_predicate` when `sequential?(form) AND function?(first(form))`, else `assert_any`. The `sequential?` guard is what lets `(is true)`, `(is x)`, and `(is [1 2])` reach `assert_any` without calling `first` on a non-sequence. |
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
    RAISE IllegalArgumentException("The number of args doesn't match are's argv.")
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
    alter_meta!(*ns*, assoc, :test/each-fixtures, fixtures)

METHOD use_fixtures :once(fixture_type, fixtures):
    alter_meta!(*ns*, assoc, :test/once-fixtures, fixtures)

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
        once = join_fixtures(meta(ns)[:test/once-fixtures])
        each = join_fixtures(meta(ns)[:test/each-fixtures])
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
--   - run_all_tests with no matching namespace calls (apply run-tests ()) which
--     dispatches to the zero-arity and therefore runs the CURRENT namespace, as
--     Clojure 1.12.5 does (verified: from a namespace with one deftest,
--     (run-all-tests #"no-such") returns {:test 1 ...}). Callers must not invoke
--     it from inside a test of the current namespace. successful? treats
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
    -- Step 1: the `is` form's own source position (Section 4.2), threaded into
    -- the expansion by `try-expr` as a literal {file line} map
    pos = *assertion-position*
    IF pos IS NOT None: RETURN pos
    -- Step 2: find the var under test
    v = first(*testing-vars*)
    IF v IS NOT None AND meta(v).line IS NOT None:
        RETURN {file: meta(v).file, line: meta(v).line}
    -- Step 3: fall back to unknown
    RETURN {file: None, line: None}

FUNCTION error_position(thrown) -> {file, line}:
    -- Innermost user frame of the thrown exception; for a scalar or an
    -- exception with no trace, the assertion's own position.
    frames = drop_while(fn(f): internal?(f), ex_trace(thrown) OR [])
    IF frames IS NOT empty:
        RETURN {file: first(frames).file, line: first(frames).line}
    RETURN fail_position()

FUNCTION internal?(f : Frame) -> Boolean:
    -- Frames the test machinery and dispatch natives add between the user's
    -- code and the throw (Section 5.2 rule 4).
    RETURN f.fn starts with "test/" OR f.fn starts with "clojure.test/"
        OR f.fn IN {"apply", "partial", "map", "reduce", "sort", "sort-by"}
        OR f.fn starts with "core/-dispatch" OR f.fn starts with "clojure.core/-dispatch"
```

`*assertion-position*` is bound by `try-expr` around the assertion body to the `is` form's `{:file :line}` from `FormSource`, or `nil` when the form has no entry. It is the only way a `:fail` learns its line; no stack is consulted.

The observable result matches Clojure: `FAIL` and `ERROR` events carry `file` and `line`; `PASS` events do not. The fallback chain terminates in `nil`, which `testing-vars-str` renders as `(:)`.

### 10.2 `*test-out*` and `with-test-out`

```
registered_host_stdout = None

FUNCTION install_host_output_roots(stdout, stderr):
    -- One serialized host-lifecycle operation, completed before guest execution.
    registered_host_stdout = stdout
    set_root(*out*, stdout)
    set_root(*err*, stderr)
    IF test/*test-out* EXISTS: set_root(test/*test-out*, stdout)

FUNCTION register_test_out(var):
    -- Called immediately after test/*test-out* is interned, so both load orders work.
    set_root(var, registered_host_stdout OR root_binding(*out*))

MACRO with_test_out(body) -> Form:
    RETURN a form that:
        BINDING *out* = *test-out*:
            evaluates body
```

Behavior:
- Every default `report` method and every `tap-report` method wraps its printing in `with-test-out`.
- Because `*test-out*` is captured once and not rebound by `*out*` bindings, `(with-out-str (run-tests 'x))` returns `""`. Test output goes to the original stdout, matching the oracle (Appendix B.3). To capture, bind `*test-out*` to an `io/buffer`, run, then read `io/buffer-str`. Section 4.5 makes that binding effective.
- `installIOBuiltins` registers the native stdout/stderr pair. The test namespace calls `register-test-out` when it interns `*test-out*`. This makes initialization correct whether the test namespace loads before or after host setup.
- Every Go host that replaces root output uses `install-host-output-roots`; direct host `SetRoot` calls on `*out*` are prohibited. The WASM template replaces its two direct `SetRoot` calls with one operation before namespace or main chunks execute, so `*out*` and `*test-out*` select the same `HostWriter` lifecycle-wise.
- `api.WithStdout` is a per-runtime binding rather than a process-root replacement. It pushes the same handle for `*out*` and `*test-out*` on the run's `ExecContext`, then pops both on every exit. Temporary captures used by `with-out-str`, nREPL, and individual WASM requests still bind only `*out*`; they do not silently capture test reports.
- User Lisp changes to the `*out*` root do not implicitly change `*test-out*`, preserving them as distinct vars. If root replacement after guest execution begins is added later, the paired update and readers require the same synchronization rather than two unrelated `SetRoot` calls.
- `(identical? *test-out* *out*)` at a fresh REPL remains `true`.

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

```

`file-position` (deprecated in Clojure 1.2) is not provided: it reads the live JVM stack and let-go has no equivalent (Section 5.3, Appendix A).

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
    IF ex_trace(v) IS NOT None AND exception_value?(v):
        print_cause_trace(v, *stack-trace-depth*)   -- Section 5.4
    ELSE IF ex_trace(v) IS NOT None:                -- let-go permits non-exception throws
        prn(v)
        print_trace_frames(v, *stack-trace-depth*)
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
    -- Mirrors tap.clj: (.split data "\\n"), i.e. Java limit-zero split (Section 4.9).
    FOR EACH line IN string_split(data, #"\n"):
        println("# " + line)                        -- printing "" yields "# "

FUNCTION print_tap_pass(msg : String):
    println("ok " + msg)

FUNCTION print_tap_fail(msg : String):
    println("not ok " + msg)
```

All four return `nil`. `print-tap-diagnostic` is implemented with `clojure.string/split` and inherits its Section 4.9 contract rather than stripping lines itself, so it matches Clojure's direct Java call. [R-tap-diagnostic-split] This matters for `:error` diagnostics, whose trace text ends in a newline. The negative-plan check is a let-go addition (Appendix A); Clojure prints `1..-1`.

```clj-repl @R-tap-diagnostic-split
(require '[clojure.test.tap :refer [print-tap-diagnostic print-tap-plan print-tap-pass print-tap-fail]])
(print-tap-diagnostic "a\nb")
;; out: "# a\n# b\n"
;=> nil
(print-tap-diagnostic "")
;; out: "# \n"
(print-tap-diagnostic "a\n")
;; out: "# a\n"
(print-tap-diagnostic "a\n\n")
;; out: "# a\n"
(print-tap-diagnostic "\n")
;; out: ""
(print-tap-diagnostic "a\n\nb")
;; out: "# a\n# \n# b\n"
(print-tap-diagnostic " a ")
;; out: "#  a \n"
(print-tap-pass "m")
;; out: "ok m\n"
;=> nil
(print-tap-fail "m")
;; out: "not ok m\n"
;=> nil
(print-tap-plan 0)
;; out: "1..0\n"
;=> nil
```

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
    -- The loader reports namespace declarations from the forms it already reads;
    -- the harness does not reparse source or infer this from CurrentNS.
    loaded, load_error = load_file_with_result(path)
    IF load_error IS NOT None: FAIL file with copied error text
    IF NOT loaded.saw_ns_form:
        -- Load-only file: its assertions ran at load time. A deftest defined
        -- here would silently never run, so that is the failure, not the
        -- missing ns form.
        IF loaded.final_ns HAS any var with :test metadata: FAIL file with "deftest outside a namespace"
        RETURN true

    -- Step 2: run the current namespace through the public API
    ns      = loaded.final_ns
    summary, run_error = invoke(test, "run-tests", ns) -- rt.LookupVar + rt.InvokeValue
    IF run_error IS NOT None: FAIL file with copied error text
    success, result_error = invoke(test, "successful?", summary)
    IF result_error IS NOT None: FAIL file with copied error text
    RETURN success == true

-- Behavior:
--   - The harness compiles no source strings.
--   - The harness holds no test state. Discovery, counters, and results are the
--     let-go vars of Section 3.
--   - The compiler/file loader returns `saw_ns_form` and `final_ns` from its
--     existing top-level form walk. Changing CurrentNS through `in-ns` is not
--     mistaken for a declaration.
--   - A file with no ns form is load-only. The corpus has 18 today: the 16
--     `test/gold-aot/*.lg` fixtures (driven by TestGold, not deftests),
--     `test/top_level_do_test.lg`, and `test/in_ns_auto_refer_test.lg`, which
--     deliberately uses only `in-ns`. The walk additionally skips `gold-aot`
--     as it skips `gogen`, since those files are lowering fixtures.
--   - Per-file dynamic-binding isolation (vm.RunWithBindings around the run) is unchanged.
--   - Output goes through *test-out*, which is stdout, so `go test -v` shows the same text
--     a user sees.
```

### 12.3 Fast Follow: The Report Bridge

The bridge binds `report` to a context-aware native fn for the whole load-and-run operation, so each `ReportEvent` is frozen on the runner goroutine before any Go harness code sees it. The let-go side is unchanged: the bridge is just another reporter, like `tap-report`.

`testing.T.Run` runs its body on a new goroutine and blocks the caller until that body returns. The let-go runner therefore cannot call `t.Run` itself, or it would block inside its own `BEGIN_TEST_VAR` and never emit `END_TEST_VAR`. The bridge inverts control: `run-tests` executes on its own goroutine and hands every event across an unbuffered channel to the harness goroutine, which owns `t` and calls `t.Run` from there. The unbuffered send is a rendezvous, so let-go does not advance until Go has consumed the event, and ordering and attribution hold.

```
RECORD BridgeEnvelope:                  -- immutable, Go-owned
    sequence        : UInt64
    kind            : ReportType
    scope_id        : UInt64 | None
    parent_scope_id : UInt64 | None
    namespace_text  : String
    qualified_name  : String
    failure_text    : String
    default_text    : String
    summary         : copied integer fields | None

-- No field is a vm.Value, Var, Namespace, persistent collection, dynamic-var
-- reference, writer, or lazily rendered object.

ENUM TerminalKind: NORMAL, INVOKE_ERROR, PANIC, ABORTED, CANCELLED

RECORD Terminal:
    kind        : TerminalKind
    successful  : Boolean
    summary     : copied integer fields | None
    error_text  : String
    panic_text  : String
    go_stack    : String

FUNCTION bridge_report(ec, events, cancel) -> NativeFn:
    -- State below is confined to the runner goroutine.
    next_sequence = 0
    next_scope_id = 0
    scope_stack = []
    RETURN context_aware_fn(event):
        CASE event.type:
            PASS:  inc_report_counter(:pass)
            FAIL:  inc_report_counter(:fail)
            ERROR: inc_report_counter(:error)

        envelope = freeze_event_on_runner(ec, event)
        -- freeze_event reads *testing-vars*, *testing-contexts*, and trace depth;
        -- renders names, expected, actual, traces, default text, and complete
        -- failure text; and copies summary integers before bindings can change.
        envelope.sequence = next_sequence++

        IF event.type == BEGIN_TEST_VAR:
            envelope.scope_id = ++next_scope_id
            envelope.parent_scope_id = last(scope_stack) OR None
            push(scope_stack, envelope.scope_id)
        ELSE IF event.type == END_TEST_VAR:
            REQUIRE scope_stack IS NOT empty
            envelope.scope_id = last(scope_stack)
        ELSE:
            envelope.scope_id = last(scope_stack) OR None

        SELECT:
            send(events, envelope)                 -- unbuffered rendezvous
            receive(cancel): RAISE bridge_cancelled

        IF event.type == END_TEST_VAR: pop(scope_stack)

RECORD ConsumerState:
    protocol_error : ProtocolError | None
    expected_sequence : UInt64

FUNCTION protocol_fault(state, events, text):
    IF state.protocol_error IS None: state.protocol_error = ProtocolError(text)
    drain(events)                              -- read until the sole producer closes

FUNCTION consume_scope(sub : testing.T, begin : BridgeEnvelope, events, state):
    FOR EACH envelope IN events:
        IF envelope.sequence != state.expected_sequence++:
            protocol_fault(state, events, "non-contiguous event sequence")
            RETURN
        CASE envelope.kind:
            BEGIN_TEST_VAR:
                IF envelope.parent_scope_id != begin.scope_id:
                    protocol_fault(state, events, "invalid nested scope")
                    RETURN
                sub.Run(envelope.qualified_name,
                        fn(child): consume_scope(child, envelope, events, state))
            FAIL, ERROR:
                IF envelope.scope_id != begin.scope_id:
                    protocol_fault(state, events, "event attributed to wrong scope")
                    RETURN
                sub.Error(envelope.failure_text)   -- text is never a printf format
            END_TEST_VAR:
                IF envelope.scope_id != begin.scope_id:
                    protocol_fault(state, events, "scope ended out of order")
                    RETURN
                RETURN
            PASS:
                IF envelope.scope_id != begin.scope_id:
                    protocol_fault(state, events, "pass attributed to wrong scope")
                    RETURN
            ELSE:
                IF envelope.scope_id != begin.scope_id:
                    protocol_fault(state, events, "event attributed to wrong scope")
                    RETURN
                sub.Log(envelope.default_text)
    IF state.protocol_error IS None:
        state.protocol_error = ProtocolError("event stream closed before end-test-var")

FUNCTION consume_events(t : testing.T, events) -> ProtocolError | None:
    -- A protocol mismatch is recorded once and switches to drain mode. The
    -- consumer never returns early and strands an unbuffered producer.
    state = ConsumerState(expected_sequence = 0)
    FOR EACH envelope IN events:
        IF envelope.sequence != state.expected_sequence++:
            protocol_fault(state, events, "non-contiguous event sequence")
            RETURN state.protocol_error
        CASE envelope.kind:
            BEGIN_TEST_NS:
                IF envelope.scope_id IS NOT None:
                    protocol_fault(state, events, "namespace began inside a var scope")
                    RETURN state.protocol_error
                t.Log(envelope.default_text)
            BEGIN_TEST_VAR:
                IF envelope.parent_scope_id IS NOT None:
                    protocol_fault(state, events, "top-level var has a parent scope")
                    RETURN state.protocol_error
                t.Run(envelope.qualified_name,
                      fn(sub): consume_scope(sub, envelope, events, state))
            FAIL, ERROR:
                IF envelope.scope_id IS NOT None:
                    protocol_fault(state, events, "assertion attributed outside active scope")
                    RETURN state.protocol_error
                t.Error(envelope.failure_text)  -- assertion emitted directly by test-ns-hook
            PASS:
                IF envelope.scope_id IS NOT None:
                    protocol_fault(state, events, "pass attributed outside active scope")
                    RETURN state.protocol_error
            END_TEST_NS, SUMMARY:
                IF envelope.scope_id IS NOT None:
                    protocol_fault(state, events, "top-level event has a var scope")
                    RETURN state.protocol_error
                t.Log(envelope.default_text)
            END_TEST_VAR:
                protocol_fault(state, events, "end-test-var without matching begin")
                RETURN state.protocol_error
            ELSE:
                IF envelope.scope_id IS NOT None:
                    protocol_fault(state, events, "top-level event has a var scope")
                    RETURN state.protocol_error
                t.Log(envelope.default_text)
        IF state.protocol_error IS NOT None: RETURN state.protocol_error
    RETURN state.protocol_error

FUNCTION run_file_tests_bridged(t, path) -> Boolean:
    events = new unbuffered Channel<BridgeEnvelope>
    done   = new Channel<Terminal>(capacity = 1)
    cancel = new close-only Channel

    -- Step 1: all let-go loading, state access, formatting, and invocation stay
    -- on this runner goroutine.
    GO:
        terminal = None
        DEFER:
            IF panic was recovered:
                terminal = Terminal(PANIC, panic_text = copied_text,
                                    go_stack = copied_go_stack)
            ELSE IF terminal IS None:
                terminal = Terminal(ABORTED)
            send(done, terminal)                   -- capacity one: never waits
            close(events)                          -- runner is the sole closer

        TRY:
            WITH rt.WithBinding(ec, report_var, bridge_report(ec, events, cancel)):
                loaded, load_error = load_file_with_result(path)
                IF load_error IS NOT None:
                    terminal = Terminal(INVOKE_ERROR, error_text = copy_text(load_error))
                    RETURN
                IF NOT loaded.saw_ns_form:
                    terminal = Terminal(INVOKE_ERROR, error_text = "no test namespace")
                    RETURN
                summary, run_error = invoke(test, "run-tests", loaded.final_ns)
                IF run_error IS NOT None:
                    terminal = Terminal(INVOKE_ERROR, error_text = copy_text(run_error))
                    RETURN
                success, result_error = invoke(test, "successful?", summary)
                IF result_error IS NOT None:
                    terminal = Terminal(INVOKE_ERROR, error_text = copy_text(result_error))
                    RETURN
                -- Copy and evaluate while still on the runner ExecContext.
                terminal = Terminal(NORMAL,
                                    summary = copy_summary(summary),
                                    successful = success == true)
        CATCH bridge_cancelled:
            terminal = Terminal(CANCELLED)

    -- Step 2: this goroutine performs only testing.T calls and string handling.
    DEFER close(cancel)
    protocol_error = consume_events(t, events)     -- always drains to close
    terminal = receive(done)                       -- exactly one terminal
    protocol_failed = protocol_error IS NOT None
    IF protocol_failed: t.Error(protocol_error.text)
    CASE terminal.kind:
        NORMAL:       RETURN terminal.successful AND NOT protocol_failed
        INVOKE_ERROR: t.Error(terminal.error_text); RETURN false
        PANIC:        t.Error(terminal.panic_text + "\n" + terminal.go_stack); RETURN false
        ABORTED:      t.Error("let-go runner exited without a result"); RETURN false
        CANCELLED:    t.Error("let-go runner cancelled"); RETURN false

-- Behavior:
--   - Each deftest is a Go subtest named <ns>/<var>, selectable with -run.
--   - The bridge increments counters itself, exactly as default methods and
--     tap-report do, since binding `report` replaces the methods that would
--     otherwise count. successful? on the returned summary remains the authoritative
--     pass/fail result; the Error calls agree by construction.
--   - The bridge does not wrap printing in with-test-out. testing.T owns the output.
--   - Nested BEGIN/END pairs receive explicit scope IDs. A composed deftest opens
--     a nested Go subtest; its END returns only from that matching recursive scope,
--     so later outer events stay attributed to the outer subtest.
--   - The let-go goroutine holds the ExecContext for the whole load and run. The
--     harness goroutine never touches let-go state or formats a live value.
--   - One outer defer converts every normal return, invocation error, panic, or
--     runtime.Goexit-like abandonment into exactly one buffered Terminal and closes
--     events. Every producer send also selects on cancel, so consumer abandonment
--     cannot leave the producer blocked. Protocol faults drain before returning.
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
| `:trace` under `ex-data` as strings | Vector of `Frame` maps, derived from the chain on read, excluded from equality | Readers of the string form must switch. |

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

**Trace capture for values never thrown.** `(ex-trace (ex-info "x" {}))` on a value constructed but not thrown returns `nil`. The trace is the unwind chain, which exists only once the value is thrown and caught. Extension point: an `ex-info` arity that takes an explicit chain.

**`current-stack-trace`.** A trace of the running code without a throw. The runtime has no live frame chain (Section 5.3); adding one costs a push and pop per call on both backends for a primitive nothing in `clojure.test` needs. Extension point: a linked frame stack on `ExecContext`, which would also enable parallel bridged subtests.

**Specific exception classes for Go runtime errors.** `(/ 1 0)` and `(nth [] 5)` raise `java.lang.Exception`, not `ArithmeticException` or `IndexOutOfBoundsException`, so `(is (thrown? ArithmeticException (/ 1 0)))` reports `ERROR`. Excluded by the scope of #472, which tags wrapped Go errors as `Exception` deliberately, and by the preference not to model more of the JVM inside Go. `(is (thrown? Exception ...))` covers them. Extension point: the single `class:` tag in `errorToValue` (`pkg/vm/errors.go`), which could consult the error kind.

---

## 15. Design Decision Rationale

**Why discover tests by `:test` metadata instead of keeping the registry?** kaocha and cognitect test-runner read `(:test (meta var))` via `ns-interns` and `ns-publics`. A registry is invisible to them. Metadata also makes `(math-test)` composable and enables `run-test` to check it.

**Why is the TAP plan printed last rather than first?** Clojure's `tap-report :summary` prints it after the run for both `run-tests` and `run-all-tests`, since the plan counts assertions and no runner knows that count in advance. Perl's Test::More prints first only when the script declares `plan tests => N`; with `done_testing` it prints last. TAP 14 allows either. Plan-first would require a dry run and would diverge from the oracle byte-for-byte.

**Why capture `*test-out*` once instead of following `*out*` dynamically?** The oracle does, and `with-out-str` around `run-tests` returning `""` is documented Clojure behavior that harnesses rely on when binding `*test-out*` to their own writer. Following `*out*` dynamically would make `(binding [*test-out* w] (with-out-str ...))` ambiguous.

**Why coordinate host output installation with `test.lg` registration?** `test.lg` can load from the bundle before an embedder replaces the root binding of `*out*`, but another host can install output first. `install-host-output-roots` and `register-test-out` cover both orders and choose the same live host handle once. Capturing only at bundle load would freeze a handle the embedder later abandons. The observable semantics are the same as Clojure's; only the initialization handshake differs.

**Why an atom for `*report-counters*` instead of a plain map with `set!`?** `inc-report-counter` derefs and `commute`s a reference; kaocha's `with-report-counters` binds `*report-counters*` to a ref and derefs it. A plain map breaks both. let-go's `ref`, `commute`, and `dosync` aliases let the Clojure source run as written.

**Why take `:file` and `:line` from the stack instead of only from `&form` metadata?** The runtime already tracks a source position per instruction and builds traces from it. The stack is the source Clojure uses, works for `ERROR` events whose throw site is inside a helper the `is` form never saw, and is what `print-stack-trace` needs anyway. Var metadata remains the second fallback step, not the primary.

**Why expose `FormSource` through `meta` rather than attaching metadata in the reader?** The side table already holds a position for every supported identity-bearing list/cons form and is propagated through macroexpansion; attaching metadata in the reader would duplicate it, touch every list allocation, and require the bundle to round-trip form metadata. Exposing the table is three small readers and no format change.

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
- [ ] `use-fixtures :each` in a namespace puts a seq under `:test/each-fixtures` in `(meta *ns*)`

### 16.2 Data Model (Section 3)

- [ ] Every var in the Section 3.3 table exists in `test`, is dynamic, and has the stated default
- [ ] `*test-result*`, `*registered-tests*`, `*each-fixtures*`, `*once-fixtures*`, `register-test!`, `clear-registered-tests!` are gone
- [ ] `(meta (var some-deftest))` contains `:test`, `:ns`, `:name`, `:file`, `:line`
- [ ] Every `:fail` and `:error` event carries at least the Section 3.1 keys; a `report` method receiving an extra key ignores it
- [ ] A frame from `ex-trace` is a map with `:fn`, `:kind`, `:file`, `:line`, and `:column` for `:lg` frames

### 16.3 Runtime Extensions (Section 4)

- [ ] `(alter-meta! *ns* assoc :k 1)` succeeds and `(:k (meta *ns*))` is `1`
- [ ] `(meta (read-string "(a b)"))` contains `:line` and `:column`, and so does the nested `(+ 2 2)` inside `(read-string "(is (= 4 (+ 2 2)))")`
- [ ] Inside a macro, `(meta &form)` carries the call site's `:line`
- [ ] `*file*` is bound to the path during file load and to `"NO_SOURCE_PATH"` at the REPL
- [ ] `(binding [*out* (io/buffer)] (println "x"))` leaves stdout untouched and the buffer holding `"x\n"`
- [ ] `(try (throw "s") (catch Throwable e :caught))` returns `:caught` and `(try (throw "s") (catch Exception e :caught) (catch Throwable e :thr))` returns `:thr` (regression guard for #476)
- [ ] `test/catch_dispatch_test.lg` still passes
- [ ] `(ex-message "s")` returns `"s"` when `"s"` was the thrown value
- [ ] `rt.WithBinding` pushes a binding for the call and pops it on every exit path, including a panic
- [ ] Evidence `@R-string-split-limit-zero` (Section 4.9) passes

### 16.4 Stack Traces (Section 5)

- [ ] `(try (throw (ex-info "x" {})) (catch e (ex-trace e)))` is a non-empty vector whose first frame has the throw site's `:file` and `:line`
- [ ] `(try (throw "s") (catch e (ex-trace e)))` is `nil`, and `(is (throw "s"))` reports `ERROR` at the `is` form's file and line
- [ ] A trace read after the throwing frames have been reused by later calls still reports the original positions
- [ ] A runtime error inside `nth` yields a first frame of kind `:native` whose position is the let-go call site, and `(ex-cause e)` is a boxed Go error whose message is the Go `Error()` text
- [ ] A recovered Go panic yields at least one `:go` frame with a Go file and line, followed by the `:native` frame
- [ ] Lisp calling native calling Lisp that throws yields one vector with the frames in call order and one `:native` frame per crossing; this holds through `map`, `reduce`, `sort`, and `apply`
- [ ] A value thrown from a callback survives every native that invokes let-go code with its class, message, data, and trace intact (no native flattens a `ThrownError`)
- [ ] A rethrown exception keeps its trace and gains the outer frames; a new `ex-info` thrown from a `catch` has a new trace and `ex-cause` reaches the old value and its trace
- [ ] `(def boom (ex-info "x" {}))` thrown from two sites reports the first site from both catches, and `(ex-trace boom)` after the first catch is non-nil
- [ ] `(= e (ex-info "x" {}))` and `(hash e)` are unchanged by whether `e` was thrown; `(:trace (ex-data e))` on a runtime error equals `(ex-trace e)`
- [ ] A self-recursive function throwing at depth 100 in tail position shows one frame for itself; the same function in non-tail position shows 100
- [ ] Interpreter and lowered-Go execution produce the same frame sequence for every case above, including a throw from inside the lowered standard library
- [ ] A chained `(Throwable->map e)` has outer-to-root `:via`, original-value `:phase`, root `:cause`/`:data`/non-empty `:trace`, and `:at` on every via entry whose value has a trace
- [ ] An implementation-level exception with absent message/data and no trace omits top-level `:cause`/`:data` and per-via `:message`/`:data`/`:at`, while retaining `:via` and an empty `:trace`
- [ ] `(print-stack-trace e 2)` prints ex-data on its own line, then a first frame prefixed ` at ` and one later frame prefixed by four spaces
- [ ] With an empty trace, `(print-stack-trace e)` prints ` at [empty stack trace]`; with `n <= 0`, a non-empty trace still prints its first frame
- [ ] `(print-cause-trace e)` prints a `Caused by:` section per cause
- [ ] `(ex-trace (ex-info "x" {}))` on an unthrown value is `nil`

### 16.5 Assertions (Section 6)

- [ ] `(is (= 4 (+ 2 2)))` emits one `:pass` event with `:actual (= 4 4)` and returns `true`
- [ ] `(is (= 10 (+ 5 4)) "m")` emits one `:fail` with `:actual (not (= 10 9))` and `:message "m"`, returns `false`
- [ ] `(is nil)` emits `:fail` via `:always-fail`
- [ ] `(is true)`, `(is :ok)`, `(is x)` with `x` bound to `1`, and `(is [1 2])` each go through `assert-any`, emit one `:pass`, and return the value; `(is false)` emits one `:fail` with `:actual false`
- [ ] `(is (thrown? Exception (throw (ex-info "x" {}))))` emits `:pass` and returns the exception
- [ ] `(is (thrown? Exception 1))` emits `:fail` with `:actual nil`
- [ ] `(is (thrown-with-msg? Exception #"oo" (throw (ex-info "boom" {}))))` emits `:pass`
- [ ] `(is (throw (ex-info "x" {})))` emits exactly one `:error` and nothing escapes
- [ ] `(is (some-macro ...))` goes through `assert-any`
- [ ] `(are [x y] (= x y) 1 1 2 2)` emits two `:pass` events; mismatched macroexpansion fails with root cause `java.lang.IllegalArgumentException` and exact message `"The number of args doesn't match are's argv."`
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
- [ ] Installing a host writer before or after `test` registration makes root `*out*` and `*test-out*` identical to that writer; the WASM entry uses the same hook
- [ ] `api.WithStdout` dynamically binds both `*out*` and `*test-out*`; `with-out-str` and other temporary captures bind only `*out*`
- [ ] Default output for the reference example matches Appendix B.1 modulo file path
- [ ] A zero-failure run prints only the `Testing` and `Ran` lines
- [ ] An `:error` for `(throw "s")` prints `  actual: "s"`, then a first frame prefixed ` at ` and later frames prefixed by four spaces
- [ ] An `:error` for an `ex-info` throw prints a header and `at` lines, no more than `*stack-trace-depth*` when bound

### 16.10 clojure.test.tap (Section 11)

- [ ] Evidence `@R-tap-diagnostic-split` (Section 11.1) passes: diagnostics, pass, fail, and plan output, split semantics, and `nil` returns
- [ ] `(print-tap-plan -1)` raises `ex-info` (let-go extension, Appendix A)
- [ ] `(with-tap-output (run-all-tests #"my\.test.*"))` prints Appendix B.2 modulo namespace print form and file path, then returns the summary
- [ ] The plan line is the last line for `run-tests`, `run-all-tests`, and `run-test-var`
- [ ] Plan count equals the number of `ok`/`not ok` lines
- [ ] An `:error` under TAP renders its trace as consecutive `#`-prefixed lines
- [ ] After `with-tap-output` exits, `(run-tests 'x)` prints default output

### 16.11 Go Harness (Section 12)

- [ ] `test/language_test.go` compiles no source strings; it uses `rt.LookupVar` and `rt.InvokeValue`
- [ ] After loading a file, the harness runs exactly the namespace current at the end of that file
- [ ] A file with no `ns` form that defines no `deftest` passes as load-only; one that defines a `deftest` fails with "deftest outside a namespace"
- [ ] `test/in_ns_auto_refer_test.lg` and `test/top_level_do_test.lg` pass unmodified; `test/gold-aot/` is skipped by the walk
- [ ] A loader unit test with a counting one-shot reader proves `saw_ns_form` and `final_ns` come from one existing top-level form walk; `(in-ns ...)` alone changes `final_ns` without setting `saw_ns_form`
- [ ] `LG_TEST_REPORTER=summary` selects the summary-based run; unset selects the bridge once it lands
- [ ] Summary run: a file with failing assertions fails its Go subtest with the summary in the message
- [ ] Bridge: `run-tests` executes on a goroutine other than the one that owns `testing.T`, and the harness never deadlocks on a file with one or more deftests
- [ ] Bridge: every `deftest` is a Go subtest named `<ns>/<var>` selectable with `-run`
- [ ] Bridge: a frozen envelope contains only copied Go strings, integers, booleans, and IDs—no VM value, var, namespace, collection, writer, or lazy rendering
- [ ] Bridge: a `:fail` event produces exactly one `t.Error` on the matching subtest, with the Section 10.4 text
- [ ] Bridge: an `:error` event produces exactly one `t.Error` with the rendered trace; a `:go` frame shows its Go file and line
- [ ] Bridge: `:pass` produces no Go output; counters still reach the summary
- [ ] Bridge: a deftest invoked from another deftest becomes a nested subtest, and events after its matching end remain on the outer subtest
- [ ] Bridge: harness stdout is empty; all text flows through `t.Log`/`t.Error`, and percent signs in messages remain literal
- [ ] Bridge: load/invoke errors, panics, cancellation, premature channel close, and mismatched scope IDs each produce one terminal result without panic, timeout, deadlock, or a stranded producer
- [ ] `go test ./test/...` passes on the migrated corpus

### 16.12 Migration (Section 13)

- [ ] `test/tap/tap-example.lg` exists and produces Appendix B.1 and B.2 under the CLI
- [ ] No non-asserting demonstration script remains under `test/`
- [ ] `make generate` and `make check-generated` are clean

### 16.13 Integration Smoke Test

The reference example from Section 1.7, run through both reporters in one process. Output is captured by binding `*test-out*` to `*out*` inside `with-out-str`, which works on both engines. [R-smoke-run-and-tap]

```clj-test @R-smoke-run-and-tap
(ns my.test.tap-example
  (:require [clojure.test :refer [deftest is run-tests run-all-tests testing successful?]]
            [clojure.test.tap :refer [with-tap-output]]
            [clojure.string :as str]))

(deftest math-test
  (testing "simple addition"
    (is (= 4 (+ 2 2)))
    (is (= 5 (+ 2 3))))
  (testing "deliberate-failure-test"
    (is (= 10 (+ 5 4)) "This assertion will fail intentionally")))

(ns my.test.smoke
  (:require [clojure.test :refer [deftest is run-tests run-all-tests testing
                                  *test-out* *testing-vars* *testing-contexts*]]
            [clojure.test.tap :refer [with-tap-output]]
            [clojure.string :as str]))

(defn- capture
  "Run f as if at top level: *testing-vars* and *testing-contexts* are reset so
   the inner run's headers do not include this deftest's own var."
  [f]
  (let [result (volatile! nil)
        out (with-out-str
              (binding [*test-out* *out* *testing-vars* (list) *testing-contexts* (list)]
                (vreset! result (f))))]
    [out @result]))

(deftest plain-reporter
  (let [[out summary] (capture #(run-all-tests #"my\.test\.tap-example"))]
    (is (str/includes? out "\nTesting my.test.tap-example\n"))
    (is (str/includes? out "\nFAIL in (math-test) ("))
    (is (str/includes? out "deliberate-failure-test\nThis assertion will fail intentionally\n"))
    (is (str/includes? out "expected: (= 10 (+ 5 4))\n  actual: (not (= 10 9))\n"))
    (is (str/ends-with? out "\nRan 1 tests containing 3 assertions.\n1 failures, 0 errors.\n"))
    (is (= {:test 1 :pass 2 :fail 1 :error 0 :type :summary} summary))))

(deftest tap-reporter
  (let [[out summary] (capture #(with-tap-output (run-all-tests #"my\.test\.tap-example")))
        lines (str/split-lines out)]
    (is (= 2 (count (filter #(str/starts-with? % "ok ") lines))))
    (is (= 1 (count (filter #(str/starts-with? % "not ok ") lines))))
    (is (= "1..3" (last (remove str/blank? lines))))
    (is (some #(= % "# {:type :begin-test-var, :var #'my.test.tap-example/math-test}") lines))
    (is (= {:test 1 :pass 2 :fail 1 :error 0 :type :summary} summary))))

(ns my.test.error-example (:require [clojure.test :refer [deftest is]]))
(deftest boom (is (= 1 (nth [] 5))))

(ns my.test.smoke)
(deftest error-through-native
  (let [[out summary] (capture #(run-tests 'my.test.error-example))]
    (is (str/includes? out "\nERROR in (boom) ("))
    (is (str/includes? out "expected: (= 1 (nth [] 5))\n  actual: "))
    (is (not (str/includes? out "ERROR in test:")))
    (is (= {:test 1 :pass 0 :fail 0 :error 1 :type :summary} summary))))
```

The rendered class and message of the `:error` are engine-specific (the JVM prints `java.lang.IndexOutOfBoundsException`), so they are a separate let-go-only block: [R-smoke-error-text]

```clj-test @R-smoke-error-text oracle=none
(ns my.test.smoke-error-text
  (:require [clojure.test :refer [deftest is run-tests *test-out* *testing-vars* *testing-contexts*]]
            [clojure.string :as str]))
(deftest error-text
  (let [out (with-out-str
              (binding [*test-out* *out* *testing-vars* (list) *testing-contexts* (list)]
                (run-tests 'my.test.error-example)))]
    (is (str/includes? out "  actual: java.lang.Exception: nth index out of bounds"))
    (is (str/includes? out " at nth ("))))
```

Prose-only cases that remain:

- `go test ./test/ -run 'TestRunner/tap-example.lg' -v` reports `--- FAIL: TestRunner/tap-example.lg/my.test.tap-example/math-test` with `(math-test) (test/tap/tap-example.lg:9)` and `actual: (not (= 10 9))` (bridge, Section 12.3)
- The failure line names `test/tap/tap-example.lg:9` when the example is loaded from that file

---

## 17. Executable Evidence

The Definition of Done in Section 16 is prose. Nothing runs it, so a wrong expectation (the trailing-empty split rule in the first draft) is caught only by a reviewer who happens to check. This section makes the document's own assertions executable, in the document, against any Clojure engine, with the expected output inline so no oracle pass is needed to validate a run.

### 17.1 Shape

The spec is literate. A provable statement carries a marker, `[R-<slug>]`, and its evidence is a fenced block or a table tagged with the same slug, placed beside it:

````
```clj-repl @R-tap-diagnostic-split
(print-tap-diagnostic "a\n")
;; out: "# a\n"
;=> nil
```
````

Pairing is checked both ways: a marker without evidence and evidence without a marker are both lint errors. The fence grammar is the one `rfc-tangle` already reads (` ```<type> @R-<slug> `, tables preceded by `<!-- evidence: @R-<slug> -->`), so the block is extracted verbatim into `<spec>.<slug>.<type>` by a tool that never interprets content.

### 17.2 Vocabularies Are Data, Runners Are Generated

An evidence *type* is a **vocabulary file**: rule lines mapping a statement pattern to a Clojure template, plus a `boot` preamble. No per-type program interprets the block. One generic generator reads the vocabulary, walks the tangled block statement by statement, instantiates the matching template, and writes a runner `.cljc`. The runner is then executed by whichever engine is selected. This is a literate untangle with generated code around the vocabulary's own harness, and it mirrors the shell `rfc-flow` runner with Clojure forms in the templates instead of commands.

```
# clj-repl.vocab — top-level forms with inline expectations
boot #?(:clj (import 'clojure.lang.ExceptionInfo))
boot (require '[clojure.test :refer :all] '[clojure.string :as str])
boot <spec-evidence preamble: case, expect-value, expect-out, expect-throws, report, exit>
{form*}                  => (case! {id} (fn [] {form*}))
;=> {v*}                 => (expect-value! {id} '{v*})
;; out: {s*}             => (expect-out! {id} {s*})
;; throws: {c} {m*}      => (expect-throws! {id} {c} {m*})
```

| Type | Statement shapes | Verdict |
|---|---|---|
| `clj-repl` | A top-level form, then any of `;=> <edn>` (value, compared with `=` after reading), `;; out: "<string>"` (captured `*out*` and `*test-out*`, compared byte-exact), `;; throws: <Class> "<message>"` (class by `instance?`, message by `=`). A form with no expectation must evaluate without throwing. | All cases hold. |
| `clj-table` | Markdown rows `\| form \| value \| out \|`. Each row compiles to the same three statements as `clj-repl`, which is the SLIM shape: a table is an instruction list. Empty cell means no expectation for that column. | All rows hold. |
| `clj-test` | Ordinary `deftest` forms. The runner wraps them in a namespace and calls `run-tests`. | `successful?` |

Template variables: `{id}` is the 1-based case number within the block, `{form*}` the rest of the line (multi-line forms are joined until the reader has a complete form). Everything in a block is plain Clojure plus `clojure.test` and `clojure.string`. Dialect setup, such as an import a JVM needs, lives in the vocabulary's `boot` lines under reader conditionals, never in evidence.

### 17.3 Engines

The generated runner is portable Clojure. The only engine-specific fact is the launch command, taken from `CLJ_ENGINE` and defaulting to `./lg`:

| Engine | Command | Role |
|---|---|---|
| let-go | `./lg <runner.cljc>` | The gate. Runs in `TestSpecEvidence` under `go test ./test/...`, so it rides the pre-push hook and CI. |
| Clojure 1.12.5 | `clojure -M <runner.cljc>` | Validates the expectations themselves. `make spec-oracle`; needs a JDK, so it is a separate CI job. |
| Others (babashka, jank) | their launcher | Same file, no changes. |

The runner prints one line per case, `ok <id>` or `FAIL <id> want <edn> got <edn>`, and exits 1 if any case failed. Cases tagged `oracle=none` in the fence info string are let-go extensions with no JVM behavior (Appendix A); they run under every engine but only the let-go verdict counts.

### 17.4 Promotion

A mismatch is not only a red run. `spec-evidence --accept <spec.md>` reruns the runner, parses the `FAIL` lines, and rewrites the `;=>`, `;; out:`, or `;; throws:` line (or the table cell) for each failed case with the engine's actual result, then sets `last-verified` in the masthead to today. The result is a diff on the spec, reviewed in git like any other change. Two engines separate the two kinds of mismatch:

| JVM run | let-go run | Meaning | Action |
|---|---|---|---|
| fails | any | The document's expectation is wrong | `--accept` from the JVM run |
| passes | fails | A conformance gap in let-go | Fix the runtime; promote nothing |
| passes | passes | Conformant | none |

Promotion from the let-go run alone is permitted only for `oracle=none` cases.

### 17.5 Tooling and Wiring

- `scripts/spec-evidence.lg`: tangle, generate, run, and `--accept`. Written in `.lg`, so the reader that parses evidence forms is the engine under test; the JVM pass is the independent check on that reader.
- `scripts/spec-vocabs/clj-repl.vocab`, `clj-table.vocab`, `clj-test.vocab`: the three vocabularies, reviewed as evidence machinery.
- `test/spec_evidence_test.go`: `TestSpecEvidence` runs every `docs/specs/*.md` containing tagged evidence under `./lg`; one Go subtest per block.
- `make spec-oracle`: the same under `CLJ_ENGINE="clojure -M"`.
- `make spec-lint`: marker and evidence pairing, fence grammar, and that every `clj-*` block tangles and generates without running.

### 17.6 Coverage Rule

A Definition of Done line that can be stated as a form and an expected result must be a fence, and Section 16 cites the slug instead of restating the case. Lines that remain prose are those with no in-process oracle: the bridge deadlock and cancellation cases, the `go test -run` transcripts, and ratchet results. A prose line is visibly unverified, which is the point.

The first three fences convert cases that were verified by hand during review: the split contract (Section 4.9), the TAP diagnostic rules (Section 11.1), and the integration smoke test (Section 16.13).

---

## Appendix A: Deviations from the Oracle

| Behavior | Clojure | let-go | Reason |
|---|---|---|---|
| `print-tap-plan` with negative n | prints `1..-1` | raises `ex-info` | A negative plan is invalid TAP. |
| `ex-message` on a non-exception | not applicable; cannot throw one | returns `(str v)` | let-go permits throwing any value. |
| `catch Throwable` | Throwables only | any thrown value, strings included | Pre-existing let-go rule (#476); `Exception` stays typed. |
| Stack capture for an `ex-info` value | captured when the JVM Throwable is constructed | captured on its first `throw`; an unthrown value has no trace | Avoids a stack walk for every constructed value and supports arbitrary let-go throw values. |
| Fixture metadata keys through the canonical `test` namespace | `:clojure.test/each-fixtures`, `:clojure.test/once-fixtures` | `:test/each-fixtures`, `:test/once-fixtures` | `clojure.test` aliases the canonical namespace object; auto-resolved keywords use its canonical name. |
| Class of a Go runtime error | `ArithmeticException`, `IndexOutOfBoundsException`, ... | `java.lang.Exception` | Scope of #472; no JVM classification of Go errors (Section 14). |
| Metadata on a list returned by `read-string` | `nil` | `:line` and `:column` merged from `FormSource` | Reuses the source table needed for precise test diagnostics; no reader or bundle change. |
| `&form` metadata | reader-attached | merged from `FormSource` on macro invocation | Same keys at the macro boundary. |
| `do-template` location | `clojure.template` | `test` | No other consumer; add alias later if needed. |
| `file-position` | reads JVM stack | absent | No live frame chain in let-go; deprecated since Clojure 1.2 and unused by `clojure.test`. |
| Trace of a thrown scalar | not applicable; cannot throw one | `nil`; reporters use the assertion's position | A bare Go string or keyword has no field to hold a chain. |
| Tail-call frames | `recur` shows one frame | `OP_TAIL_CALL` shows one frame | Same semantics; stated so recursion depth is not expected in a trace. |
| Frame shape | `[class method file line]` vector | `Frame` map with `:kind` | `clojure.main/ex-triage` destructures frames positionally; a port of it needs a `[fn kind file line]` projection. |
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
| `(keys (meta *ns*))` after `(use-fixtures :each f)` in Clojure | `(:clojure.test/each-fixtures)`; let-go uses `(:test/each-fixtures)` (Appendix A) |
| `*initial-report-counters*` | `{:test 0, :pass 0, :fail 0, :error 0}` |
| Nested `testing` contexts in a diagnostic | `# Arithmetic with positive integers` |
| An `ex-info` `:error` under TAP | `#   actual:clojure.lang.ExceptionInfo: boom`, then `# {:x 1}`, then frame lines beginning `#  at ` / `#     ` |
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
