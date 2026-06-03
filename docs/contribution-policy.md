---
status: active
last-verified: 2026-06-01
authoritative-for:
  - design-contracts
  - regression-checkpoints
  - callback-error-contract
  - go-deps-schema
  - host-interop-syntax-map
  - wasm-constraints
supersedes:
  - master-plan.md (on direction — self-host as committed direction, gogen_ir as deployment path)
---

# Contribution policy

This document captures the design contracts, regression checkpoints, and
interop conventions that govern changes to `let-go`. It is the
human-readable counterpart to the CI gates and `make` targets that
enforce them. Source: [issue #97 — Comprehensive Go interop redesign
roadmap](https://github.com/nooga/let-go/issues/97).

Changes that conflict with anything below need explicit project-owner
sign-off in the PR thread, not just a green CI run.

---

## 1. Design contracts

These are load-bearing. New code is expected to fit inside them, not
ask them to move.

1. **Standalone interpreter stays portable.** REPL, nREPL, builtins,
   and bytecode loading must work from just the interpreter binary —
   no Go toolchain on the host. Authoring features (`lg deps`,
   generator, walker) live behind opt-in surfaces and must not pull
   into the minimal runtime distribution.
2. **`let-go` is embeddable.** It must remain usable as a scripting
   engine inside a host Go program. Public Go-side API surface (`pkg/api`,
   `pkg/vm`, `pkg/rt`) changes are reviewed against this constraint.
3. **Bytecode + VM are first-class.** Native/AOT lowering (`-tags
   gogen_ir`, `lower-go`, future `lg build`) is additive. The VM still
   powers `eval`, dynamic features, REPL workflows, embedded use, and
   runtime compilation. No change may make any of those depend on the
   AOT path.
4. **Host interop is a Clojure-compatible superset.** Same spelling
   and same semantics where they exist; Go-specific extensions are
   purely additive. The mapping table is in §6.
5. **Authored binaries go through generated Go.** The target for
   `lg build` is a Go build that includes only the necessary `let-go`
   modules and Go dependencies — not a bytecode blob appended to the
   stock `lg` binary. The same machinery should eventually drive
   minimal / network / full distributions.
6. **Performance and footprint are tracked across modes.** Four
   surfaces: interpreter, embedded, wasm, native. Per surface, per
   program in the perf corpus we track `cold-start`, `warm-run`,
   `peak-rss`, and `artifact-size`. The numbers live in `test/perf/`
   so deltas are visible in PR diffs.

## 2. Direction: self-host

The longer arc is `let-go` implemented in `let-go` — compiler, reader,
runtime ops, eventually the VM dispatch loop — written in `.lg` and
lowered to Go for the shipped binary. This is the committed direction
(issue #97, owner reply).

Two consequences for sequencing PRs:

- **Smart-wrapper coverage is scoped by what the compiler touches.**
  Not "as much as we feel like." The packages `pkg/compiler/` and
  `pkg/vm/` use directly (`reflect`, `bytes`, `unicode`, …) are the
  reference user. Coverage gaps the self-host path doesn't need are
  not blockers.
- **`gogen_ir` build is the deployment path for the self-hosted
  compiler.** Bootstrap parity (untagged vs `-tags gogen_ir`) is not
  optional polish — it is the gate that says self-host works.

## 3. Regression checkpoints

A PR that trips any of these requires an explicit review checkpoint in
the thread — not a silent override.

| Checkpoint | Threshold | Enforcement |
| --- | --- | --- |
| Cold-start regression | > 10 % on any surface | perf corpus, `test/perf/` (in flight) |
| Artifact-size regression | > 10 % on any surface | perf corpus, `test/perf/` (in flight) |
| Bootstrap parity divergence | any count/bucket delta | `make parity-full` — local target, CI wiring deferred (see note below) |
| Stale `core_compiled.lgb` | committed bytes differ from regeneration | `make check-generated` — CI gate |
| Stale `core_go_lowered/` | committed tree fails to compile under `-tags gogen_ir` | `make check-generated` — CI gate |
| Existing test suite | any regression | `make test` — CI gate (already in place) |

**Why `parity-full` isn't yet a CI gate.** Wiring it requires
`lgbgen --target=go` to produce `pkg/rt/core_go_lowered/` cleanly from
current `main` — and today that gen path errors on a `core.lg` reader
EOF. The fix lands in a follow-up; until then the policy describes the
target end-state but CI only enforces the staleness gates that protect
the artifacts `parity-full` depends on.

The two staleness gates are siblings. `core_compiled.lgb` is the bytecode
bundle loaded by the untagged build; `core_go_lowered/` is the generated
Go linked into `-tags gogen_ir` builds. Both derive from
`pkg/rt/core/**/*.lg`. If either falls behind the sources, the two engines
quietly run different versions of the IR pipeline — `parity-full` then
diverges on bucket hashes even when pass/fail counts match. Caught the
hard way 2026-05-28.

`make parity-full` runs `scripts/gogen-parity.sh --full`, covering:
1. clojure-test-suite (jank) — end-user-observable `clojure.core`
2. ir-stress `lower-go` mode — AOT lowering correctness
3. ir-stress `ir-compile` mode — eval-mode IR compile correctness

Each suite runs once under the untagged build and once under
`-tags gogen_ir`; bucket-by-symbol diffs must match.

## 4. Callback error contract

Signature-driven, not ambient. Decided in issue #97 and binding here.

- If the Go callback type returns `error`, let-go failures route
  through that slot and other return values are zero-valued.
- If there is no error slot, a private panic sentinel is raised
  inside the proxy and recovered at the outer let-go/native boundary.
- Async escaped callbacks without an error slot need API-specific
  policy and documented process-level behavior — they do not get a
  default.

Out of scope (rejected): sticky "next VM op" errors; ambient `*err*`
as the primary contract.

## 5. Interop schema

### `:go-deps`

A single ns-form key for stdlib and external Go packages.

```clojure
(ns my.app
  (:require ...)
  (:go-deps
    "net/http"
    "github.com/foo/bar@v1.2.3"
    {:pkg "golang.org/x/sys/unix" :platform :linux}))
```

- Entries are strings (`"path"` or `"path@version"`) or maps for
  alias / platform / explicit version metadata.
- No `:git/url` / `:git/sha` initially — Go modules carry that via
  `go.mod` and `replace`. Add only if a real gap appears.

### Generator manifest

When the generator emits binding manifests, the rule is:

- **Keywords** for `let-go`/compiler-owned tags: `:kind :func`,
  `:kind :method`, `:kind :error`, `:variadic? true`.
- **Strings** for Go-owned identifiers: `"net/http"`,
  `"github.com/foo/bar"`, `"ListenAndServe"`.

Each binding entry carries import path, package name, generated
alias, Go identifier, params, and results as distinct fields — not
embedded in the let-go symbol.

## 6. Host-interop syntax map

The target is full Clojure compatibility on the JVM-interop surface,
with Go-specific extensions where there is no Clojure analogue.

| Clojure (JVM)              | ClojureScript      | let-go on Go                              | Status      |
| ---                        | ---                | ---                                       | ---         |
| `(.method obj args)`       | same               | `Boxed.InvokeMethod`                      | shipped     |
| `(.-field obj)`            | same               | explicit field access                     | follow-up   |
| `(set! (.-field obj) v)`   | same               | field setter                              | follow-up   |
| `Class/staticMethod`       | `js/foo.bar`       | `net:http/Get`                            | shipped     |
| `Class/STATIC_FIELD`       | `js/CONSTANT`      | `net:http/StatusOK`                       | shipped     |
| `(Class. args)`            | `(js/Foo. args)`   | `(Foo. args)`                             | target      |
| `:import` in `ns`          | `:require :refer`  | `:go-deps` ns key                         | this policy |
| `proxy` / `reify`          | `reify`            | reify → Go structural interfaces          | follow-up   |
| `bean`                     | —                  | `(go/bean struct)`                        | partial     |
| `clj->js` / `js->clj`      | both               | `clj->go` / `go->clj`                     | follow-up   |
| `clojure.java.io`          | —                  | `clojure.go.io`                           | follow-up   |

Field-access semantics: `(.X obj)` falls method→field (method
preferred when both exist); `(.-X obj)` is explicit field access.
This matches Clojure on the JVM.

Constructor syntax: prefer `(Foo. args)` (Clojure-style). `(Foo/new args)`
is not introduced — one way to spell a constructor.

## 7. WASM-specific constraints

The wasm surface surfaces footprint tradeoffs earliest, so it gets
explicit rules:

- `.wasm` size and browser cold-start are first-class checkpoint
  numbers, not derived from the artifact-size column.
- `lg deps` carries a `js/wasm` compatibility filter with clear
  diagnostics for packages that cannot work in wasm — fail fast at
  resolve time, not deep in `go build`.
- If the default binary grows to carry deps/generator machinery,
  verify the wasm build does not pay for authoring-only code unless
  the user opted into it.

## 8. Distribution UX target

The end-user authoring entry point is `lg deps`. Generator machinery
ships in the default binary only if it does not noticeably grow size
or startup. Otherwise it lives behind an authoring-only build path.

The default `lg` binary stays lean.

---

## How CI enforces this today

`make check-generated` and `make parity-full` run on every PR via
`.github/workflows/go.yml`. The generated-artifacts check regenerates
the bundle and fails if the committed bytes differ, and verifies the
committed lowered tree compiles under `-tags gogen_ir`; it fails the
build with an explicit remediation message. The parity
check is the longer gate (~5 min) but is what guarantees the
self-host substrate keeps working.

`make test` (the existing path) still runs as before — it does not
go away. Parity is layered on top.

## How to update this policy

This document is owner-reviewed. Changes need a PR thread referencing
the rationale, and ideally a link to a related issue discussion.
Tightening a threshold (e.g. moving the 10 % gate to 5 %) does not
need full review; loosening one does.
