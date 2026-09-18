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
# evidence: skip until clojure.test, clojure.test.tap, and the trace primitives this spec describes are implemented
evidence: skip
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
    fn      : String                -- qualified function name, e.g. "my.test.tap-example/math-test",
                                     -- or the callee's name at a call-site link, e.g. "throw"
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
| `NATIVE` | Rendered as `at <fn> (<file>:<line>:<column>)` using the let-go call site — the same shape the runtime's existing top-level printer already uses for a native call. Counted. |
| `GO` | Rendered as `at <fn> (<file>:<line>) [go]`. Counted. Omitted entirely when a reporter requests `lg_only`. |

A trace is a vector of `Frame`, innermost first.

**The trace is the unwound error chain.** No live frame stack is walked. `throw` is a native; `(throw ...)` is an ordinary call that returns `ThrownError{value}` as a Go `error`, and that call's own call-site link wraps it in one `ExecutionError` — a `CALL_NATIVE` link named `throw`, carrying the throw form's position (Section 4.2) — before the enclosing function's handler runs. The bytecode VM already wraps every error leaving a call at `OP_INVOKE`; Section 4.2 is what makes the position it resolves there the call form's own, not its last argument's. So even a throw caught in the same function has its throw form as the innermost frame. A runtime error raised directly by a native, or a recovered Go panic, has no such call to unwind through at its own origin; that error's chain starts instead at the native's own call-site link (Section 5.2 rule 2) or at the panic's captured Go frames (Section 5.2 rule 3). A specialized VM opcode that fails — the arithmetic, comparison, and bit-operation fast paths the compiler emits for a known binary call to a core native (Section 4.10) — has no `OP_INVOKE` to wrap it either, since it never makes the call it specializes; its chain starts the same way a call to that native would have wrapped it: a `CALL_NATIVE` link named for the operator (`+`, `-`, `*`, `inc`, `<`, `bit-and`, ...) at the failing form's own `SourceInfo`, on both backends. From there outward every call site the error passes through on its way out wraps it in another `ExecutionError` link carrying that site's callee name and `SourceInfo`. The chain that arrives at a `catch` therefore already lists every frame between the error's origin and the catch, innermost first. The design adds three things: the lowered-Go backend records a link at every call site the same way the VM does (Section 4.10), the catch keeps the chain instead of discarding it, and `ex-trace` renders it.

```
ENUM LinkKind:                          -- added
    CALL_LG        -- wraps an error leaving a call to a function compiled from .lg source
    CALL_NATIVE    -- wraps an error leaving a call to a Go-implemented primitive

RECORD ExecutionError:                  -- existing struct, two fields added
    message : String
    source  : SourceInfo | None
    cause   : Error | None
    kind    : LinkKind | None           -- added; None for a link that is not a frame ("integer overflow")
    fn      : String | None             -- added; the callee's name for a call-site link, None for a link with kind None

RECORD ExInfo:                          -- existing pointer struct, one field added
    message : String
    data    : Map
    cause   : Error | None              -- a Go error: either another *ExInfo given to ex-info,
                                         -- or a host error reached through host_error
    meta    : Value
    class   : ExceptionClass            -- existing; ExceptionInfo for ex-info, Exception for a boxed runtime error
    chain   : Error | None              -- added: set once at first catch; never overwritten (CAS)

FUNCTION step(err : Error) -> Error | None:
    -- One raw step down the Go chain, with no frame-link skipping: err.GetCause() when err
    -- implements it (pkg/errors.Error, as ExecutionError and TypeError do), else the standard
    -- library's errors.Unwrap(err) (for a native's fmt.Errorf("...: %w", err) wrap). None at
    -- the end of the chain. thrown_value, frames_of, host_error, and next_link all walk with
    -- this same step, so an ExecutionError link and a %w wrap are each just one more kind of
    -- link to step through, on the same footing.

FUNCTION thrown_value(err : Error) -> Value:
    -- Walks err with step looking for a ThrownError anywhere in the chain and returns its
    -- Value when found; a native's fmt.Errorf("...: %w", err) wrap around a ThrownError is
    -- one more link stepped through, not a barrier, so the wrap is transparent to catch
    -- dispatch. Returns box_caught(err) when no ThrownError is found.

FUNCTION catch_bind(err : Error) -> Value:
    -- The single seam both backends use (vm.ErrorToValue today).
    value = thrown_value(err)                       -- ThrownError.Value, or a boxed exception for a Go error
    IF value IS *ExInfo AND value.chain IS None:
        compare_and_swap(value.chain, None, err)    -- first catch wins; a reused exception keeps its first trace
    RETURN value

FUNCTION host_error(err : Error) -> Error | None:
    -- The first link in err's own chain that is not an ExecutionError frame link: an
    -- ExecutionError whose kind is None is a host error in its own right (e.g. the one
    -- NewExecutionError("integer overflow") constructs), not a link to skip past. Steps through
    -- every ExecutionError whose kind is CALL_LG or CALL_NATIVE. None when the chain is
    -- only ExecutionError frame links. box_caught calls host_error only when thrown_value found
    -- no ThrownError in err's chain, so neither host_error nor next_link ever steps onto one.

FUNCTION next_link(err : Error) -> Error | None:
    -- One step down the Go chain: step(err), then skip ExecutionError frame links the same way
    -- host_error does. None at the end of the chain.

FUNCTION box_caught(err : Error) -> *ExInfo:
    -- The box a catch binds when err carries no ThrownError: a runtime error or recovered panic.
    -- data stays {} rather than None so ex-data can carry the compatibility :trace key (Section 5.3);
    -- a cause box further down the chain (ex_cause, below) carries no data at all. The box represents
    -- h, the head host error, directly -- it does not duplicate h as its own cause.
    h = host_error(err)
    RETURN ExInfo(h.Error(), {}, cause = next_link(h), class = Exception)

FUNCTION ex_cause(v) -> Value | None:
    IF v IS NOT *ExInfo OR v.cause IS None: RETURN None
    c = v.cause
    IF c IS *ExInfo: RETURN c                                              -- the value given to ex-info
    RETURN ExInfo(c.Error(), None, cause = next_link(c), class = Exception)  -- a cause box; no data, so
                                                                              -- ex-data returns nil for it,
                                                                              -- as Clojure's does for a
                                                                              -- non-IExceptionInfo throwable.
                                                                              -- class is Exception, not a class
                                                                              -- derived from c's Go type: finer
                                                                              -- host classes are out of scope
                                                                              -- (Section 14).

FUNCTION ex_trace(v) -> Vector<Frame> | None:
    IF v IS *ExInfo AND v.chain IS NOT None: RETURN frames_of(v.chain)
    RETURN None

FUNCTION frames_of(err) -> Vector<Frame>:
    -- Walks err from the head with step, reading kind, fn, and source off each
    -- ExecutionError link; never parsing message. A link whose kind is CALL_LG or
    -- CALL_NATIVE is an LG or NATIVE frame naming the callee at the call site — the
    -- throw native's own call-site link included. A link whose kind is None contributes
    -- no frame. A GoPanicError contributes its captured Go frames as GO frames. A link
    -- of any other type (a %w wrap, for instance) is stepped through and contributes no
    -- frame. step descends from the head (the outermost wrap) toward the error's origin,
    -- collecting one frame-link per step; only these links are reversed before returning.
    -- A GoPanicError's captured Go frames are already innermost first (the Go runtime's own
    -- stack capture at the recover site), so the result is those GO frames in their captured
    -- order, followed by the reversed frame links — innermost link first — starting with the
    -- NATIVE frame of the panicking primitive. Resolution is lazy: the chain is stored,
    -- frames are built on read.
```

A call-site link's frame is named for the **callee** and positioned at the **call site**. `errorToValue`'s top-level printer keeps using the `calling <fn>` message text for a call-site link; `frames_of` does not depend on it.

`chain` is a computed field. It participates in neither `=` nor `hash`, so `(= e (ex-info "x" {}))` is unaffected by whether `e` was ever thrown, and `WithMeta` copies it like any other field. Setting it once matches the JVM, where a `Throwable` captures its stack at construction and rethrowing the same object keeps that trace.

**Advance.** A cause box's own `cause` is `next_link` of the error it boxes, never that error again, so each `ex-cause` step moves one link further down a finite Go chain and `cause_chain` (Section 5.3) terminates at the chain's end with `None`.

**Terminal.** A Go error with neither `GetCause` nor `Unwrap` boxes with `cause = None`.

**Frame links are not causes.** An `ExecutionError` link whose `kind` is `CALL_LG` or `CALL_NATIVE` is a frame, not a cause; `host_error` and `next_link` skip it, so a trace frame never appears as a `:via` entry (Section 5.3). An `ExecutionError` whose `kind` is `None` is not a frame link — it is a host error in its own right and is never skipped.

**Cause boxes have no trace.** `chain` is `None` on a cause box, so `ex-trace` returns `nil` for it and `Throwable_to_map` omits `:at` for it (Section 5.3).

**Identity.** A cause box is built fresh on every `ex-cause` read. Exceptions compare by identity in let-go, as on the JVM, so two calls on the same Go error return boxes that are neither `=` nor identical, even though each carries the same message, data, and class.

Scalars have no trace. A thrown string, keyword, number, boolean, or `nil` is a bare Go value with no field to hold a chain, so `ex-trace` on it returns `nil` and the test reporter falls back to the assertion's own source position (Section 10.1). After the Section 13 migration the only scalar throws in the repository are test fixtures.

---

## 4. Runtime Extensions

These are the Go-side capabilities the ported namespaces require and the runtime lacks (September 2026), verified against the working tree. They are general-purpose rather than test-specific. Most are small. Trace primitives are in Section 5.

### 4.1 Namespace Metadata

`(meta ns-obj)` returns the namespace's metadata map. `(alter-meta! ns-obj f & args)` updates it. Currently `alter-meta!` raises "expected Atom or Var" for a namespace. Since `*ns*` inside a file is the namespace object, `(alter-meta! *ns* assoc k v)` works without further lookup.

### 4.2 Form Source Positions

The reader already records a `SourceInfo` for every identity-bearing list or cons form it produces, nested forms included, in the `vm.FormSource` side table (`readList` in `pkg/compiler/reader.go`). Calls attempting to record non-hashable vector and map values are intentionally ignored by the current table and are not part of this slice. The compiler and macroexpander copy entries onto rewritten and expanded list forms, and `compileForm` emits them per instruction. What is missing is the Clojure-facing surface: `(meta form)` does not consult the table and `&form` inside a macro is nil.

**A call's position at the invoke instruction is the call form's, not its last argument's.** `compileForm` records a form's `SourceInfo` once, at the instruction offset where that form's code starts (`c.chunk.AddSourceInfo`, `pkg/compiler/compiler.go`), and `SourceMap.Lookup` returns the entry with the greatest `startIP` not after the instruction (`pkg/vm/source.go`). A call compiles its callee and each argument before emitting `OP_INVOKE` or `OP_TAIL_CALL`, and only a list or cons argument records its own entry (non-hashable vector and map values are intentionally ignored by the table, above, and a symbol or scalar argument records nothing), so today the position resolved at the invoke instruction is the most recently recorded entry before it — the last argument's own entry when that argument is itself a list or cons form, or the call form's own entry, still standing from when the call itself began compiling, when the last argument records nothing: a multi-line `(throw (ex-info ...))` reports the `ex-info` form's line, not the `throw` form's, but `(throw x)` with `x` a bound symbol already reports the `throw` form's own position, since `x` adds no later entry to displace it. The compiler gains one more source-map entry: it records the call form's `SourceInfo` again immediately before emitting `OP_INVOKE` and `OP_TAIL_CALL`, once the callee and arguments are compiled, so `LookupSource` at the invoke instruction resolves to the call form. This is one extra source-map entry per call with arguments; it costs compile time and source-map size only, nothing on the execution path. The IR builder does not share this gap: `build-form` sets the current form's `SourceInfo` before dispatching into it and restores the caller's on return, so by the time a `:call` instruction is added the context's `SourceInfo` is back to the call form's own, not the last argument's, on both backends alike.

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

**Lowered Go wraps at every call site.** Generated code today propagates errors with a bare `return nil, err`, so a throw inside a lowered function unwinds through it without recording a frame. The IR already attaches the full `SourceInfo` of the originating form to every instruction, calls included, but the IR bridge exposes only `source-info-symbol` to `.lg`. The bridge gains `source-info-file`, `source-info-line`, and `source-info-column`, and the call emitter in `lower_go.lg` wraps on the error path with the same `calling <fn>` shape the VM uses, setting `kind` and `fn` on the link:

```
-- emitted for every lowered call whose result may be an error
IF err != nil:
    RETURN nil, vm.NewExecutionError("calling <fn>").
        WithSource(<file>, <line>, <column>).
        WithFrame(<kind>, "<fn>").
        Wrap(err)
```

`<fn>` is the callee's name and `<kind>` is `CALL_LG` for a direct Go call to a lowered function, known statically at emission, or, for a dynamic call through `rt.InvokeValue*`, the callee's origin read at the wrap (LG-origin or native-origin, below).

**Specialized opcodes wrap on the error path too.** `OP_ADD`, `OP_SUB`, `OP_MUL`, `OP_INC`, the bit operations (`OP_BIT_AND`, `OP_BIT_OR`, `OP_BIT_XOR`, `OP_BIT_AND_NOT`, `OP_BIT_SHIFT_LEFT`, `OP_BIT_SHIFT_RIGHT`, `OP_BIT_NOT`), and the comparisons (`OP_LT`, `OP_LTE`, `OP_GT`, `OP_GTE`) are fast paths the compiler emits for a known binary call to a core arithmetic/comparison native, avoiding `NativeFn.Invoke`, interface boxing, and `RecoverPanic` (`pkg/vm/vm.go`). Today, on failure, each calls `f.handleError` directly with the bare error from `checkedAddInt`/`checkedSubInt`/`checkedMulInt`, `errIntOverflow`, `errBitOpType`, or the shared `NumAdd`/`NumSub`/`NumMul`/`NumLt`/`NumLe`/`NumGt`/`NumGe` helpers, without going through `wrapCallErr` first, so `frames_of` sees a kind-less host error and produces no frame for that link. The design adds the same wrap `wrapCallErr` already gives a native call: before `f.handleError` runs, the opcode's error path wraps the error in a `CALL_NATIVE` link named for the operator it specializes (`+`, `-`, `*`, `inc`, `<`, `bit-and`, ...) at the form's own `SourceInfo` (Section 4.2), exactly as a call to that operator would have. The lowered-Go backend's emitter wraps its numeric-op helper's error path the same way, so a specialized-operation failure produces the identical first frame on both backends.

**`throw` is an ordinary call on both backends.** `(throw ...)` has no dedicated instruction or IR op on either backend: the reader and the IR builder resolve `throw` as an ordinary symbol, and the compiler and `lower_go.lg` each lower a call to it exactly like a call to any other native. The wrap above already covers it. Because `throw` is a hand-written native (`CoreThrowf`, `pkg/rt/lang.go`) with no `SourceInfo`, its call-site link is `CALL_NATIVE`, `fn = "throw"`, positioned at the throw form — the same rule Section 5.2(2) gives for `nth`. A throw caught in the same function therefore still has its throw form as the innermost frame, with no separate mechanism needed: the bytecode VM already wraps the native's error at `OP_INVOKE` before the enclosing frame's handler runs (`wrapCallErr`), and Section 4.2 is what makes the position resolved there the throw form's own rather than its last argument's.

**A callee's origin, not its Go type, sets `kind`.**

```
ENUM Origin:                            -- added
    LG       -- a function the lowering emitter boxes from .lg source
    NATIVE   -- every other native
```

`*Func`, `*Closure`, and `*MultiArityFn` are LG-origin by construction, and never carry any other origin: a `*Func` holds a `*CodeChunk` with a `sourceMap` (`pkg/vm/func.go`, `pkg/vm/vm.go`), `*Closure` reaches one through its `*Func`, and `*MultiArityFn` has no `*Func` field of its own, so it reads the origin from any of its per-arity `Fn` branches — hand-written or lowered, every branch of one function shares the same origin (`pkg/vm/func.go`). A hand-written `NativeFn` holds no source and no origin field today (`pkg/vm/native_func.go`). Lowered Go is the one addition: a function lowered from `.lg` source is boxed as a `NativeFn` through `rt.BoxNativeFn`, which today sets neither a name nor any source on the box. `NativeFn` gains an `origin : Origin` field — the addition — that the lowering emitter sets to `LG` when it boxes a lowered function; a hand-written native leaves it `NATIVE`. `source : SourceInfo | None` stays what it already is, optional location data set when known, rather than the discriminator: a stripped or sourceless lowered function can lose its `source` without becoming native, because `origin` is set independently at the point the emitter boxes it. The emitter also populates the existing `name` field, to the same string `fnName` (`pkg/vm/vm.go`) reports for that function on the VM, including its anonymous-function spelling, `"anonymous fn"`, so a dynamic call to a lowered function names its `CALL_LG` link identically on both backends. One criterion, both backends: a callee whose origin is `LG` is `CALL_LG`. A direct Go call between lowered functions is statically `CALL_LG`, known at emission. In the VM, `wrapCallErr` and `wrapCallSite` set `kind` from the callee in hand and `fn` from `fnName`. Neither fallback path silently drops a call site: when `wrapCallSite` cannot read the callee as an `Fn` at all, the link is `CALL_NATIVE` named `"fn"`, since the callee is not an `.lg` function; when the call site itself cannot be located (`pkg/vm/vm.go`, the arity-error path), the link's `kind` is `None` — a host error, not a frame.

**A frame's `fn` names the callee's binding, and `fnName` always yields one.** A native bound to a var reports that var's name: `Namespace.Def` already does this for a native it binds, calling `SetName` on any `*NativeFn` value passed to it (`pkg/vm/namespace.go`), so `ns.Def("nth", nthf)` (`pkg/rt/lang.go`) names that native `nth` at the moment it is defined. The name is then lost on one path: the generated-primitive registrar rebinds the same, already-interned var through `setPrimitiveRoot` (`pkg/rt/native_prims_lifecycle.go`), which calls `existing.SetRoot` rather than `Namespace.Def` and never renames the fresh adapter — built unnamed by `vm.NativeFnType.Wrap` — that it installs as the new root, so `fnName` reports `"native fn"` for `nth` today (`pkg/vm/vm.go`). The design closes this one path: `setPrimitiveRoot` names the adapter it installs the same way `Namespace.Def` does, from the var it is rebinding. A function lowered from `.lg` source reports the name `fnName` gives it on the VM, `"anonymous fn"` included for an anonymous one, by the naming mechanism the paragraph above describes. `fnName`'s other fallbacks are not defects to close: `"native fn"` is the legitimate name for a native bound to no var — a Go closure produced at run time rather than interned — and `"fn"` is `wrapCallSite`'s name for a callee that is not a function at all (above).

The happy path is unchanged, so this does not affect the bench ratchet. Positions come from the same `SourceInfo` the VM resolves through `LookupSource`, so the two backends produce the same frames for the same program. A call whose form has no `FormSource` entry records `<unknown>` as the VM does.

**Natives preserve the chain.** A native that re-wraps an error with `fmt.Errorf("...: %v", err)` flattens the chain to a string and loses the `ThrownError` inside it, so a user's `(throw (ex-info ...))` passing through that native arrives at the catch as a generic exception. Every such site becomes `%w`. `unwrapThrown` (`pkg/vm/errors.go`) walks only `ExecutionError` links today, so a `ThrownError` under a `%w` wrap is invisible to it; `unwrapThrown` gains `step` (Section 3.5) as its walk, the same one `thrown_value`, `frames_of`, `host_error`, and `next_link` use, so a `%w` wrap is an ordinary link it steps through rather than a barrier. A test asserts that a value thrown from a callback survives, with its trace, through `map` (once its lazy sequence is realized — Section 5.2 rule 9), `reduce`, `sort`, `apply`, and every native that invokes let-go code. The survival requirement — class, message, data, and trace intact — holds through `map` regardless of the frames a crossing there produces; `map` never contributes a `NATIVE` crossing frame of its own (Section 5.2 rule 9).

**Go errors are boxed at two points.** `box_caught` (Section 3.5) is the box a catch binds for a runtime error or recovered panic; it represents `h = host_error(err)`, the first link in the Go chain that is not an `ExecutionError` frame link, directly — its message is `h.Error()` and its `cause` is `next_link(h)`, not `h` itself, so the head host error is not duplicated as its own cause. `ex-cause` builds a second, lazier box on each read, for a Go error reached as a cause: its class is `Exception` (Section 14 puts finer host classes out of scope), its message is `Error()`, and its own `cause` is `next_link` of the error it boxes rather than that error again, so repeated `ex-cause` calls advance one link at a time down the chain, and `Throwable->map` (Section 5.3) walks it to `None`.

---

## 5. Stack Traces

### 5.1 Existing Machinery

The runtime records source positions per instruction (`CodeChunk.LookupSource` over `SourceInfo` in `pkg/vm/source.go`), chains `ExecutionError` values with a source per call, converts that chain into a `:trace` list of strings under `ex-data` for runtime errors (`errorToValue` in `pkg/vm/errors.go`), captures the Go stack when a Go panic is recovered (`GoPanicError`), and prints a `stack trace:` block for uncaught errors at the top level. This section unifies those facilities and exposes them to Lisp without adding a second trace mechanism.

### 5.2 Capture Rules

1. A trace is the `ExecutionError` chain that unwinds from the error's origin to the catch, read at the catch and stored on the caught `ExInfo` (Section 3.5). Every chain starts at a native's call-site link — the `throw` native's own call-site link for a thrown value, or the raising primitive's call-site link for a runtime error raised directly by a native (rule 2) — at a specialized opcode's own call-site link for a runtime error a specialized arithmetic, comparison, or bit-operation opcode raises without ever executing an `OP_INVOKE` (rule 2, Section 4.10) — or at a recovered panic's captured `GO` frames (rule 3). From there outward, every call site the error unwinds through adds one `CALL_LG` or `CALL_NATIVE` link, the same way regardless of how the chain started.
2. Each `OP_INVOKE` in the VM and each lowered call site (Section 4.10) contributes one call-site link, `CALL_LG` when the callee is compiled from `.lg` source and `CALL_NATIVE` otherwise, producing one `LG` or `NATIVE` frame named for the callee at the call site. This is how a native that raises an error directly — `nth` on out-of-range input, for example — surfaces as one `NATIVE` frame whose position is the let-go call site (Section 4.2). A specialized opcode that fails (Section 4.10) contributes the same kind of link on the operator's behalf, without an `OP_INVOKE` ever executing: `CALL_NATIVE`, named for the operator it specializes, at the form's own position — the same `NATIVE` frame a call through `OP_INVOKE` to that operator's native would have produced.
3. A recovered Go panic contributes `GO` frames from the captured Go stack, in their captured order (innermost first) and filtered to frames outside the Go runtime and the VM dispatch loop, followed by the reversed frame links — the `NATIVE` frame of the primitive first, then the `LG` frames above it.
4. A trace crossing the Go boundary more than once (Lisp calls a native, which invokes Lisp, which throws) is one vector in call order, because the native returns the callback's error unchanged and its own call site wraps it. Each crossing contributes one `NATIVE` frame, but only when the function that invokes the callback is itself a Go primitive: `reduce` and `sort` (Go primitives that invoke the callback directly), the native `apply*` that the `.lg` function `apply` delegates to, and multimethod and protocol dispatch, whose crossing frame carries the literal name `fn` rather than the dispatcher's own name (Section 4.10's `wrapCallSite` fallback, since a `MultiFn` or a protocol dispatcher is not a `*Func`). A collection function written in `.lg` (`map`, `mapv`, `filter`, `sort-by`, `apply` itself, ...) contributes an ordinary `LG` frame for its own call site instead of a `NATIVE` crossing frame there; when such a function's own body calls a Go primitive, that primitive's crossing frame is positioned inside `<embedded:core>` rather than at the user's call site — `sort-by` calls `sort` from its own `.lg` body, so the crossing frame is `sort`'s, positioned in `<embedded:core>`, with `sort-by`'s own `LG` frame above it at the user's call site. No link names the *immediate* callback value handed straight to a Go primitive as a callee, since the primitive invokes it with `ec.Invoke` rather than through a compiled call site; but when that callback is itself `.lg` code that makes a further call — its own body, or an intervening `.lg` wrapper such as `mapv`'s `(fn [acc x] (conj! acc (f x)))` — that further call is an ordinary compiled call site and gets an ordinary named frame, in call order, like any other. For example, Lisp `f` calls the native `reduce` with `g` as the reducing function; `g` calls Lisp `h`; `h` throws: the chain reads, innermost first, a `NATIVE` frame named `throw` at the throw form inside `h` (the `throw` native's own call-site link), an `LG` frame named `h` at the call site in `g` (a `CALL_LG` link), and a `NATIVE` frame named `reduce` at the call site in `f` (a `CALL_NATIVE` link) — one `NATIVE` frame per crossing, call order preserved, and no frame names `g`, the callback `reduce` invoked directly. A lazy sequence function (`map`, `filter`, `remove`, `keep`, ...) is not itself a crossing: it returns before its element function ever runs, so an error thrown while realizing its result carries the frames of the realizing call chain (`seq`, `dorun`, `doall`, or whatever consumer forces the sequence), not of the call that built the sequence, and no frame names the `map` call itself (Section 5.2 rule 9, below).
5. A tail call (`OP_TAIL_CALL`) reuses the frame, so the tail-called function does not appear between its caller and callee. A self-recursive function in tail position shows one frame, not its recursion depth. This matches `recur` on the JVM.
6. A catch stores the chain on the `ExInfo` only if it has none. Rethrowing a caught exception, `(throw e)`, starts a new `ThrownError` around the same object, and that call's own call-site link (`CALL_NATIVE`, `throw`) becomes the new chain's innermost frame; the frames the rethrow unwinds through wrap that new chain, not the stored one. `catch_bind` at the outer catch discards this new chain because `chain` is already set (Section 3.5), so `ex-trace` returns the same frames before and after the rethrow. Throwing a new `ex-info` from a catch gives it a new trace, and when the caught value is its cause, `ex-cause` reaches the old trace.
7. A caught exception may be stored in a collection, an atom, or a var and inspected later; the trace is part of the object.
8. A thrown scalar has no trace (Section 3.5). `ex-trace` returns `nil` and reporters fall back to form position.
9. A lazy sequence (`map`, `filter`, `remove`, `keep`, and every other function built on `LazySeq`) returns before its element function ever runs, so it contributes no crossing frame of its own. If it throws at all, it throws only when something realizes it, and the resulting trace carries the frames of that realizing call chain — `seq`, `dorun`, `doall`, or whichever consumer forced the sequence — innermost first, down to the element function itself; no frame names the original `(map ...)` call.

### 5.3 Primitives

```
FUNCTION ex_trace(v : Any) -> Vector<Frame> | None:
    -- The frames unwound when v was thrown. None if v was never thrown or is a scalar.

FUNCTION cause_chain(v) -> Vector<Any>:
    -- v itself, then each successive ex_cause(v) until None, outer-to-root order.
    result = [v]
    c = ex_cause(v)
    WHILE c IS NOT None:
        append(result, c)
        c = ex_cause(c)
    RETURN result

FUNCTION root_cause(v) -> Any:
    RETURN last(cause_chain(v))

FUNCTION trace_source(v) -> Any:
    -- The innermost element of cause_chain(v) that has a stored trace, walking from the
    -- root outward: the root cause itself for an ordinary ex-info chain, where only the
    -- originally-thrown value was ever caught; for a caught runtime error it is v itself
    -- (the caught box), since the cause boxes below it (ex_cause) are built fresh on read
    -- and never caught, so they carry no chain (Section 3.5).
    FOR EACH x IN reverse(cause_chain(v)):
        IF ex_trace(x) IS NOT None: RETURN x
    RETURN root_cause(v)

FUNCTION Throwable_to_map(v : Any) -> Map:
    -- Clojure's Throwable->map key/omission rules, over let-go Frame maps.
    root = root_cause(v)
    via = []
    FOR EACH x IN cause_chain(v):
        entry = {type: class_symbol(x)}
        IF ex_message(x) IS NOT None: entry.message = ex_message(x)
        IF ex_data(x) IS NOT None AND NOT empty(ex_data(x)): entry.data = ex_data(x)
        IF first(ex_trace(x)) EXISTS: entry.at = first(ex_trace(x))
        append(via, entry)
    result = {via: via, trace: ex_trace(trace_source(v)) OR []}
    IF ex_message(root) IS NOT None: result.cause = ex_message(root)
    IF ex_data(root) IS NOT None AND NOT empty(ex_data(root)): result.data = ex_data(root)
    IF get(ex_data(v), :clojure.error/phase) IS NOT None:
        result.phase = get(ex_data(v), :clojure.error/phase)
    RETURN result
```

`ex-trace` and `Throwable->map` are interned in `clojure.core`. There is no `current-stack-trace`: the runtime has no live frame chain to walk (calls recurse on the Go stack, and `ExecContext` holds bindings and scope only), and nothing in `clojure.test` needs one once positions come from form source (Section 4.2). `file-position` is therefore absent (Appendix A).

`cause_chain` follows `ex-cause`, which returns the stored `ExInfo` directly when the cause is one given to `ex-info`, and a fresh box (Section 3.5) for a Go error, one link further down that error's own chain each time. The walk terminates because each box's `cause` is `next_link` of the error it boxes rather than that error again, so it reaches `None` after as many steps as the chain has host-error links. `:via` therefore has exactly one entry per link of the caught Go error's own host-error chain for a caught runtime error, plus one entry per `ex-info` for a chain built from `ex-info` causes — `box_caught` represents the head host error itself rather than adding a duplicate entry for it (Section 3.5), so a terminal Go error produces exactly as many `:via` entries as it has links, never one more.

`:trace` under `ex-data` for runtime errors is retained for compatibility but is **derived lazily from the chain** on read, as a vector of `Frame` maps rather than strings. This is why `box_caught`'s `data` stays `{}` rather than `None` (Section 3.5): `ex-data` on the caught box must still be a map so `:trace` has somewhere to live. Because it is computed from `chain`, it is excluded from equality exactly as `chain` is, and it cannot drift from `ex-trace` on rethrow. The string form was never documented and has no known consumer.

### 5.4 `clojure.stacktrace`

```
-- root_cause is defined in Section 5.3, alongside cause_chain.

FUNCTION print_trace_element(frame : Frame):
    CASE frame.kind:
        LG, NATIVE:
            print str(frame.fn, " (", frame.file, ":", frame.line, ":", frame.column, ")")
        GO:
            print str(frame.fn, " (", frame.file, ":", frame.line, ") [go]")

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
    -- The first frame of the thrown exception's trace, unfiltered, as Clojure's
    -- do-report takes the first element of the Throwable's stack trace. For a
    -- scalar or an exception with no trace, the assertion's own position.
    frames = ex_trace(thrown) OR []
    IF frames IS NOT empty:
        RETURN {file: first(frames).file, line: first(frames).line}
    RETURN fail_position()
```

The first frame is wherever the error originated. For a `(throw ...)` that is the throw form. For an error raised inside embedded core source it is a position in that source, and for a recovered Go panic a Go file and line, as Clojure reports `Numbers.java` for `(/ 1 0)`; the caller's own line is in the trace printed with the `ERROR` block (Section 10.4).

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

Discovery follows which namespaces gained a `deftest` or a `test-ns-hook` while the file loaded: the harness snapshots every var carrying `:test` metadata, plus every namespace's `test-ns-hook` var, immediately before loading a file and, once the load returns, runs every namespace holding a `:test` var that is new or changed, or a `test-ns-hook` var that is new or changed, since that snapshot, in one `(run-tests 'a 'b ...)` call (Section 9.4). A file that interns `deftest`s in several namespaces of its own, whether it reaches each one with `ns` or `in-ns`, gets all of them tested under that one call. The harness keeps no namespace-to-file map; it compares one before/after snapshot of `:test` and `test-ns-hook` vars per file.

### 12.1 Configuration

| Key | Type | Default | Description |
|---|---|---|---|
| `LG_TEST_REPORTER` | `bridge` or `summary` | `bridge` once Section 12.3 lands, `summary` before | Selects how the harness consumes results. `summary` exists for bisecting bridge bugs. |

**Resolution precedence** (highest first):

1. `LG_TEST_REPORTER` environment variable
2. Default from the table above

### 12.2 Step One: Summary-Based Run

```
ENUM FileKind: TESTS, LOAD_ONLY, STRAY_DEFTEST

FUNCTION test_snapshot() -> Map<Var, Value>:
    -- Every var, in every namespace, that carries :test metadata, mapped to
    -- that :test value, plus every namespace's test-ns-hook var (when
    -- interned), mapped to its own value. A namespace whose only test entry
    -- point is test-ns-hook interns no :test var, so without this second
    -- half tested_namespaces would never see it touched (Section 9.3).
    -- Uses existing API only: all-ns, ns-interns, meta (rt.AllNSes from Go).

FUNCTION tested_namespaces(before) -> List<Namespace>:
    -- Namespaces holding at least one var whose :test value is absent from
    -- `before` or not identical to the value recorded there, or whose
    -- test-ns-hook var is absent from `before` or not identical to the value
    -- recorded there: the namespaces in which this load defined or
    -- redefined a deftest, or defined or redefined a test-ns-hook. Ordered
    -- by namespace name, so the run order is deterministic.

RECORD Loaded:
    kind        : FileKind          -- TESTS, LOAD_ONLY, STRAY_DEFTEST
    namespaces  : List<Namespace>   -- non-empty only for TESTS

FUNCTION classify_loaded(initial_ns, before) -> Loaded:
    -- Shared by the summary run (12.2) and the report bridge (12.3), so the
    -- two harness modes cannot classify the same file differently. Mirrors
    -- clojure.test itself: a deftest belongs to the namespace it is interned
    -- in, whatever namespace the file is current in when the load returns,
    -- and however the file reached that namespace (ns or in-ns, which is
    -- in-ns plus refers). clojure.test never asks either question, and
    -- neither does this classifier.
    touched = tested_namespaces(before)
    IF touched CONTAINS initial_ns:
        RETURN Loaded(STRAY_DEFTEST, [])  -- harness policy, not clojure.test: the starting
                                            -- namespace is shared by every file, so a deftest
                                            -- interned there can't be attributed to one. This
                                            -- fires even when the file also went on to touch a
                                            -- namespace of its own: a deftest left behind in the
                                            -- starting namespace is still a stray one.
    IF touched IS EMPTY:
        RETURN Loaded(LOAD_ONLY, [])       -- no namespace gained a test; any assertions ran at load time
    RETURN Loaded(TESTS, touched)

FUNCTION run_file_tests(path) -> Boolean:
    initial_ns = CurrentNS       -- the harness resets *ns* to clojure.core before every file (test/language_test.go)
    before = test_snapshot()
    load_error = load_file(path)
    IF load_error IS NOT None: FAIL file with copied error text
    CASE classify_loaded(initial_ns, before):
        Loaded(STRAY_DEFTEST, _): FAIL file with "deftest outside a namespace"
        Loaded(LOAD_ONLY, _):     RETURN true
        Loaded(TESTS, namespaces):
            -- Step 2: run every touched namespace through the public API, one call
            summary, run_error = invoke(test, "run-tests", namespaces...) -- rt.LookupVar + rt.InvokeValue
            IF run_error IS NOT None: FAIL file with copied error text
            success, result_error = invoke(test, "successful?", summary)
            IF result_error IS NOT None: FAIL file with copied error text
            RETURN success == true

-- Behavior:
--   - The harness compiles no source strings.
--   - The harness holds no test state. Discovery, counters, and results are the
--     let-go vars of Section 3.
--   - `load_file` returns only a load error; the harness reads no namespace back
--     from the loader. `before = test_snapshot()` is taken immediately before the
--     load, on the same goroutine and binding scope as `initial_ns` and the load
--     itself, and `classify_loaded` runs immediately after the load returns.
--   - A file that fails to load after interning some deftests fails with the load
--     error and is never classified; the deftests it already interned stay
--     interned in their namespace. A later file that touches that namespace runs
--     them under its own `run-tests` call, same as any other var an earlier file
--     left behind (below) — accepted, not attributed back to the failed file.
--   - Section 9.4 defines `run-tests` over any number of namespace arguments,
--     merging their summaries with `merge-with +` under one `:type :summary` map;
--     `invoke(test, "run-tests", namespaces...)` on a TESTS classification is
--     exactly `(run-tests 'a 'b)`, so a file that interns deftests in two
--     namespaces produces one summary, not two, matching Clojure.
--   - `run-tests` on a touched namespace runs every `:test` var interned there,
--     including ones an earlier file left behind; that is clojure.test's unit of
--     execution, the same for one namespace or several, not per-var filtering by
--     which file added which var.
--   - `classify_loaded` decides LOAD_ONLY vs. STRAY_DEFTEST vs. TESTS by which
--     namespaces gained tests. The walk skips `test/gold-aot/` (lowering
--     fixtures driven by TestGold, not deftests) before classification, the
--     same as `gogen`; `test/language_test.go`'s skip list does not name
--     `test/gold-aot` today, so this design adds that entry. Of the files the
--     walk does reach, 13 intern no deftest and classify LOAD_ONLY:
--     `test/top_level_do_test.lg` stays in the namespace the harness starts it
--     in, and the rest declare a namespace of their own and define helper
--     functions or fixtures there, but no deftest. A LOAD_ONLY file passes and,
--     in the bridge, sends no envelopes.
--     `test/in_ns_auto_refer_test.lg` enters `json` and `test.in-ns-fresh-target`
--     only to define helper functions and interns both of its deftests in
--     `test.in-ns-auto-refer-test`, so exactly one namespace is touched and it
--     classifies TESTS.
--   - Per-file dynamic-binding isolation (vm.RunWithBindings around the run) is unchanged.
--   - Output goes through *test-out*, which is stdout, so `go test -v` shows the same text
--     a user sees.
--   - Cost: the snapshot walks every namespace's interns once before each file.
--     It runs on the test path only.
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

        rt.CurrentNS.SetRoot(coreNS)   -- reset *ns* to clojure.core, on this runner goroutine,
                                        -- before initial_ns or test_snapshot() read it; CurrentNS
                                        -- is a process-global root, so this cannot run on any
                                        -- other goroutine without racing this one
        TRY:
            WITH rt.WithBinding(ec, report_var, bridge_report(ec, events, cancel)):
                initial_ns = CurrentNS
                before = test_snapshot()
                load_error = load_file(path)
                IF load_error IS NOT None:
                    terminal = Terminal(INVOKE_ERROR, error_text = copy_text(load_error))
                    RETURN
                CASE classify_loaded(initial_ns, before):  -- Section 12.2's classifier; runs on this goroutine
                    Loaded(STRAY_DEFTEST, _):
                        terminal = Terminal(INVOKE_ERROR, error_text = "deftest outside a namespace")
                        RETURN
                    Loaded(LOAD_ONLY, _):
                        terminal = Terminal(NORMAL, successful = true, summary = None)   -- load-only: assertions ran at load time
                        RETURN
                    Loaded(TESTS, namespaces):
                        summary, run_error = invoke(test, "run-tests", namespaces...)
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
--   - Each deftest is a Go subtest named <ns>/<var>, selectable with -run. A file
--     with deftests in two namespaces runs both under one `run-tests` call and
--     one Terminal, and each deftest appears as its own `<ns>/<var>` subtest.
--   - Classification is Section 12.2's `classify_loaded`, run on the runner
--     goroutine, with `before = test_snapshot()` taken on that same goroutine
--     immediately before `load_file`. A LOAD_ONLY file sends no envelopes and a
--     NORMAL terminal with no summary, which the NORMAL arm below accepts.
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

**Why is the trace the unwound error chain instead of a captured frame stack?** A frame carrier has to capture live frames at the throw site and thread them across every Go boundary: interpreter frames, native calls, callbacks into let-go, and lowered Go, which has no VM frames at all. The chain needs none of that. `throw` is a call like any other, and the bytecode VM already wraps every error leaving a call site — including that one — in one `ExecutionError` link carrying `calling <fn>` and that site's `SourceInfo` (`wrapCallSite` and `wrapCallErr` in `pkg/vm/vm.go`, reached only on the error path); a callback's error returns through the same path, so the chain crosses the Go boundary with no frame stack. What loses it today is one seam: `errorToValue` returns a thrown value and discards the chain around it. Keeping that chain on the caught `ExInfo` (Section 3.5) is the whole capture mechanism. `ExInfo` is an existing pointer struct and `WithMeta` copies it whole, so nothing changes in `ArrayVector`, `PersistentMap`, `vm.Var`, or the exported `pkg/vm` surface. The guarantee covers exception values only; a thrown scalar has no field to hold a chain (Appendix A).

The costs are on the error path and in generated code, not the happy path. The lowered core (`pkg/rt/core_go_lowered`, about 4.7 MB of Go) has about 19,000 error-return blocks and no wraps; Section 4.10 adds one wrap call inside each existing `if err != nil` branch, which grows the generated source and the `-tags gogen_ir` compile but not the bytecode binary. The specialized VM opcodes (`OP_ADD`, `OP_SUB`, `OP_MUL`, `OP_INC`, the bit operations, and the comparisons; Section 4.10) gain the same wrap only on their own error branch, an operation that already failed, so the fast path those opcodes exist to keep fast is untouched. Natives that rewrap with `%v` flatten the chain; each becomes `%w`. `ExecutionError` gains two fields, `kind` and `fn`, on the error path only. The bench ratchet should not move, since no happy-path instruction changes; Slice 1 records a before/after `make bench-ratchet` on the same base, and the generated-size and `go build` time of the lowered core before and after the emitter change.

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
- [ ] A frame from `ex-trace` is a map with `:kind`, `:fn`, `:file`, `:line`, and `:column` for `:lg` frames

### 16.3 Runtime Extensions (Section 4)

- [ ] `(alter-meta! *ns* assoc :k 1)` succeeds and `(:k (meta *ns*))` is `1`
- [ ] `(meta (read-string "(a b)"))` contains `:line` and `:column`, and so does the nested `(+ 2 2)` inside `(read-string "(is (= 4 (+ 2 2)))")`
- [ ] A multi-line call's frame reports the line and column of the call form's opening paren, not of its last argument; the same on both backends
- [ ] Inside a macro, `(meta &form)` carries the call site's `:line`
- [ ] `*file*` is bound to the path during file load and to `"NO_SOURCE_PATH"` at the REPL
- [ ] `(binding [*out* (io/buffer)] (println "x"))` leaves stdout untouched and the buffer holding `"x\n"`
- [ ] `(try (throw "s") (catch Throwable e :caught))` returns `:caught` and `(try (throw "s") (catch Exception e :caught) (catch Throwable e :thr))` returns `:thr` (regression guard for #476)
- [ ] `test/catch_dispatch_test.lg` still passes
- [ ] `(ex-message "s")` returns `"s"` when `"s"` was the thrown value
- [ ] `rt.WithBinding` pushes a binding for the call and pops it on every exit path, including a panic
- [ ] Evidence `@R-string-split-limit-zero` (Section 4.9) passes
- [ ] A native bound to a var reports the var's name as `:fn` — `(nth [] 5)` yields `:fn` `nth`, not `native fn`

### 16.4 Stack Traces (Section 5)

- [ ] `(try (throw (ex-info "x" {})) (catch e (ex-trace e)))` is a non-empty vector whose first frame has the throw site's `:file` and `:line`
- [ ] A throw caught in the same function yields exactly one frame: `:native`, `:fn` `throw`, at the throw form
- [ ] `(throw (ex-info "x" {}))` whose argument starts on a later line still reports the `throw` form's line, not the argument's, on both backends
- [ ] A throw inside a helper yields a first frame `:native` `:fn` `throw` at the throw form and a second frame at the call site, naming the helper
- [ ] A call to a function lowered from `.lg` source yields an `:lg` frame, not a `:native` frame, on the lowered backend
- [ ] A dynamic call to a lowered function yields the same `:fn` on both backends
- [ ] `(try (throw "s") (catch e (ex-trace e)))` is `nil`, and `(is (throw "s"))` reports `ERROR` at the `is` form's file and line
- [ ] A trace read after the throwing frames have been reused by later calls still reports the original positions
- [ ] `(try (+ 1 nil) (catch Exception e (ex-trace e)))`, which raises through the `OP_ADD` specialized opcode rather than `OP_INVOKE`, yields a first frame `:native` named `+` at the form's position, on both backends
- [ ] A runtime error inside `nth` yields a first frame of kind `:native` named `nth` at the call site, `(ex-message e)` is the Go error text, and `(ex-cause e)` is `nil` when the Go error wraps nothing, else a box of the wrapped error
- [ ] Repeated `(ex-cause e)` from a caught runtime error reaches `nil` after a number of calls equal to the Go chain's host-error links, and `(Throwable->map e)` on it returns a `:via` with one entry per link
- [ ] `(count (:via (Throwable->map e)))` equals the length of the Go chain for a caught runtime error — `box_caught` represents the head host error itself rather than duplicating it as its own cause (Section 3.5)
- [ ] A Go error with no `GetCause`/`Unwrap` boxes with a `nil` `ex-cause`
- [ ] A recovered Go panic yields at least one `:go` frame with a Go file and line, followed by the `:native` frame
- [ ] Lisp calling native calling Lisp that throws yields one vector with the frames in call order and one `:native` frame per crossing; this holds through `reduce`, `sort`, and `apply`
- [ ] A value thrown while realizing a lazy `map`/`filter` result yields the throw frame first, then the realizing consumer's frames (e.g. `seq`, `dorun`, `doall`), and no frame for the `map`/`filter` call itself
- [ ] A value thrown from a callback survives every native that invokes let-go code with its class, message, data, and trace intact (no native flattens a `ThrownError`)
- [ ] `(= (ex-trace e) trace-before-rethrow)` holds after an outer catch of a rethrown exception, `(throw e)`; the outer frames are not recorded on `e`
- [ ] A new `ex-info` thrown from a `catch` has a new trace and `ex-cause` reaches the old value and its trace
- [ ] `(def boom (ex-info "x" {}))` thrown from two sites reports the first site from both catches, and `(ex-trace boom)` after the first catch is non-nil
- [ ] `(= e (ex-info "x" {}))` and `(hash e)` are unchanged by whether `e` was thrown; `(:trace (ex-data e))` on a runtime error equals `(ex-trace e)`
- [ ] A self-recursive function in tail position shows one frame for itself plus the throw frame, at any depth; the same function in non-tail position shows one frame per activation plus the throw frame — 100 activations (`(deep 99)`) yield 100 `deep` frames plus the throw frame
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
- [ ] An `:error` reports the file and line of its trace's first frame, unfiltered: an error raised inside an embedded core function (e.g. a type error from `compare` reached through `sort-by`) reports that position in `<embedded:core>`, and the caller's own line appears in the printed trace
- [ ] `(with-out-str (run-tests 'x))` returns `""` and output reaches stdout
- [ ] `(binding [*test-out* b] (run-tests 'x))` with `b` from `io/buffer` leaves stdout untouched and `(io/buffer-str b)` holds the output
- [ ] Installing a host writer before or after `test` registration makes root `*out*` and `*test-out*` identical to that writer; the WASM entry uses the same hook
- [ ] `api.WithStdout` dynamically binds both `*out*` and `*test-out*`; `with-out-str` and other temporary captures bind only `*out*`
- [ ] Default output for the reference example matches Appendix B.1 modulo file path
- [ ] A zero-failure run prints only the `Testing` and `Ran` lines
- [ ] An `:error` for `(throw "s")` prints `  actual: "s"` and no frame lines, since a scalar has no trace
- [ ] An `:error` for an `ex-info` throw prints a header, then a first frame prefixed ` at ` and later frames prefixed by four spaces, no more than `*stack-trace-depth*` when bound

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
- [ ] After loading a file, the harness runs exactly the namespaces in which that file defined a `deftest` or a `test-ns-hook`
- [ ] `classify_loaded` covers three observable cases, in both the summary run and the bridge: a file that interns a `deftest` in no namespace passes as load-only; a file that leaves a `deftest` in the namespace the harness started it in fails with "deftest outside a namespace"; a file that interns `deftest`s in one or more namespaces of its own runs all of them under one summary
- [ ] `test/in_ns_auto_refer_test.lg` runs its two `deftest`s, both interned in `test.in-ns-auto-refer-test`, and passes; `test/top_level_do_test.lg` passes as load-only; both unmodified, under both `LG_TEST_REPORTER` modes; `test/gold-aot/` is skipped by the walk
- [ ] A file using `ns` and a file using only `in-ns` to reach the same other namespace classify identically: `classify_loaded` reads only which namespaces gained a `:test` or `test-ns-hook` var since `before`, never how the file reached them
- [ ] A file that interns `deftest`s in two namespaces runs both, in both modes, under one summary, and in the bridge each appears as its own `<ns>/<var>` subtest
- [ ] A file that leaves a `deftest` in the starting namespace and then moves on to intern more in a namespace of its own fails with "deftest outside a namespace"
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
- [ ] A namespace whose only test entry point is a new or changed `test-ns-hook` (no `:test` var) is discovered as touched and run through `run-tests`, in both harness modes

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

`docs/specs/executable-evidence.md` (#861) specifies a general mechanism for a spec to carry its own runnable evidence: the fence grammar, vocabularies, engines, promotion, tooling, and the coverage rule. `specs` there gates a document once it carries at least one `@R-` requirement tag, unless its masthead opts out with `evidence: skip`.

This specification carries no `@R-` tags of its own outside Section 4.9, and it opts out of the gate with `evidence: skip` in its masthead: `clojure.test`, `clojure.test.tap`, and the trace primitives this spec describes are not yet implemented, so a spec-evidence run against this document would exercise runtime capabilities that do not exist. Section 16 is the Definition of Done as written. Once the implementation lands, a Section 16 case that can be stated as a form and an expected result can become an evidence block without changing its meaning, and the opt-out is dropped.

---

## Appendix A: Deviations from the Oracle

| Behavior | Clojure | let-go | Reason |
|---|---|---|---|
| `print-tap-plan` with negative n | prints `1..-1` | raises `ex-info` | A negative plan is invalid TAP. |
| `ex-message` on a non-exception | not applicable; cannot throw one | returns `(str v)` | let-go permits throwing any value. |
| `catch Throwable` | Throwables only | any thrown value, strings included | Pre-existing let-go rule (#476); `Exception` stays typed. |
| Stack capture for an `ex-info` value | captured when the JVM Throwable is constructed | captured at its first `catch`; a value that is thrown but never caught, or never thrown at all, has no trace | Avoids a stack walk for every constructed value and supports arbitrary let-go throw values. |
| Fixture metadata keys through the canonical `test` namespace | `:clojure.test/each-fixtures`, `:clojure.test/once-fixtures` | `:test/each-fixtures`, `:test/once-fixtures` | `clojure.test` aliases the canonical namespace object; auto-resolved keywords use its canonical name. |
| Class of a Go runtime error | `ArithmeticException`, `IndexOutOfBoundsException`, ... | `java.lang.Exception` | Scope of #472; no JVM classification of Go errors (Section 14). |
| Metadata on a list returned by `read-string` | `nil` | `:line` and `:column` merged from `FormSource` | Reuses the source table needed for precise test diagnostics; no reader or bundle change. |
| `&form` metadata | reader-attached | merged from `FormSource` on macro invocation | Same keys at the macro boundary. |
| `do-template` location | `clojure.template` | `test` | No other consumer; add alias later if needed. |
| `file-position` | reads JVM stack | absent | No live frame chain in let-go; deprecated since Clojure 1.2 and unused by `clojure.test`. |
| Trace of a thrown scalar | not applicable; cannot throw one | `nil`; reporters use the assertion's position | A bare Go string or keyword has no field to hold a chain. |
| Tail-call frames | `recur` shows one frame | `OP_TAIL_CALL` shows one frame | Same semantics; stated so recursion depth is not expected in a trace. |
| Frame shape | `[class method file line]` vector | `Frame` map with `:kind` | `clojure.main/ex-triage` destructures frames positionally, so a port needs a `[fn kind file line]` projection; let-go frames also have no class, and the Go boundary requires a kind. |
| Namespace print form in `:default` diagnostics | `#object[clojure.lang.Namespace 0x.. "name"]` | let-go's namespace print form | Pointer identity is meaningless; harnesses read `:type` and name. |
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
| Rethrowing a caught exception object, `(throw e)`, then catching it again | same stack trace as before the rethrow (same frames, same count) |

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
