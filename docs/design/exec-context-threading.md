---
status: active
last-verified: 2026-06-10
human-verified:
---

# Threading an execution context (off goid)

**Status:** design note / gating decision for the dynamic-binding refactor
**Workspace:** `.workspaces/dynvar-binding`
**Decision leaning:** Option 3 — unify Scope + dynamic-var bindings into one
explicitly-threaded execution context; delete goroutine-ID keying.

## 1. Problem

The runtime has **two** ambient, goroutine-ID-keyed mechanisms for "what
context is the calling goroutine running under":

1. **Scope GLS** (`pkg/vm/scope_gls.go`) — `scopeByGID sync.Map` keyed by
   `goID()`; `CurrentScope()` / `CurrentContext()` read it. Carries structured
   concurrency (cancel/await tree) and the cancellation `context.Context` that
   blocking native ops select on.
2. **Dynamic-var binding store** (`pkg/vm/binding_frames.go`) —
   `goroutineLocalBindingStore.frames map[int64]*frame` keyed by `goID()`;
   `Var.Deref` / `PushBinding` / `PopBinding` read it.

Both resolve "current" through `goID()` (parse `runtime.Stack` output;
`goIDFast` is an intentionally-disabled no-op). This costs us:

- **The goid footgun.** A goroutine that exits with a live binding (panic
  without recover, or a binding meant to outlive it) leaks its frame; Go
  **reuses goids**, so a later goroutine on that id sees the dead one's state.
  The map also grows unbounded under leaks.
- **Cost.** `goID()` parses a stack on the fast-miss path.
- **Inconsistency / duplication.** Two parallel goid maps for the same
  underlying need.

It is **load-bearing**, which is why it can't simply be reverted to a
process-global stack: `pmap` / `parallel` (`pkg/rt/parallel.go:50`) and
**lgbgen's concurrent lowering** (N workers, each binding
`*lowered-registry*` / `*pooled-consts*` / `*ti-counters*` via the pipeline)
depend on **per-goroutine binding isolation**. A shared global stack would let
8 lowering workers corrupt each other's registries → nondeterministic lowering
(the exact failure class lg-cwqr §4 fixed).

The wall: **`Fn.Invoke([]Value) (Value, error)` is context-free**
(`pkg/vm/value.go:63`). Hand-written builtins therefore reach for *ambient*
context. ~17 `Var.Deref()` sites across 7 dynamic vars (`*out*`, `*err*`,
`*in*`, `*ns*`, `*compiling-aot*`, …) do this, and `*out*`/`*err*` — exactly
what PR #207 routes print through — are hand-written NativeFns, **not** lowered
`.lg`. So the ambient lookup is not an IR artifact; it's a VM calling
convention, and it must be addressed, not deferred.

**#206 is the lever here** (now merged to upstream as `d6d0128`, the
`main@upstream` tip). It routes `print`/`pr`/`prn`/`println` through `*out*`
and error sites through `*err*` by funnelling them into two helpers —
`WriteToOut(s)` / `WriteToErr(s)` (`pkg/rt/iort.go:364`) — each of which calls
`resolveIOHandleVar("*out*"|"*err*")` to deref the dynamic var (with an
early-boot `os.Stdout`/`os.Stderr` fallback). So the IO half of the ambient
reads — the bulk of the ~17 sites, and *all* of #207's concern — is already
**consolidated behind ~3 helpers** (`WriteToOut`, `WriteToErr`,
`resolveIOHandleVar`) rather than scattered. Incorporating #206 therefore
*shrinks* the dig-in. (#206 also exports `LookupCoreVar` so `nrepl` can
`PushBinding`/`PopBinding` the IO vars for scoped capture — another binding
consumer the ExecContext must keep working.)

**Methodology: isolate on a clean branch off `main@upstream`** (`d6d0128`,
which already carries #206), and *duplicate* the relevant binding code onto it —
rather than building on our 60-commit branch and merging #206 in. The 60-commit
branch predates #206 (its `println` writes straight to `os.Stdout`) and is
entangled with the IR/lowering work; starting clean means #206 is present by
construction (no merge) and the binding work stays reviewable in isolation. See
§7 step 0.

## 2. Goal

Eliminate goid. Achieve per-execution isolation via a single **`ExecContext`**
that carries *both* the Scope and the dynamic-var binding stack, **explicitly
threaded** through execution rather than looked up by goroutine id. Goroutine
boundaries propagate it explicitly — and already call the right primitives
(`SnapshotBindings` / `RunWithBindings` at every spawn site:
`pkg/rt/lang.go:5118,5720`, `pkg/rt/parallel.go:50`, lgbgen workers).

## 3. Design — one `ExecContext`, threaded

```go
type ExecContext struct {
    scope    *Scope          // structured concurrency + cancellation ctx
    bindings *bindingStack   // dynamic-var stack (plain, NOT goid-keyed)
}
```

- **Home:** the VM eval `Frame` (`pkg/vm/vm.go:270`) — already threaded through
  the interpreter loop — gains an `*ExecContext`. Child frames inherit the
  pointer; a new dynamic extent (a `binding` form, `with-scope`) pushes onto
  the carried stacks.
- **Compiled deref** of a dynamic var reads `frame.ec.bindings`. The VM holds
  the frame when it executes a deref op, so no ambient lookup is needed.
- **Goroutine spawn** captures `ec.snapshot()` and seeds the child frame's
  `ExecContext`. The child's ec dies with its goroutine — **no map entry to
  forget, no goid reuse hazard, no leak.** `RunWithBindings(snap, fn)` becomes
  `RunWithContext(snap, fn)` (same shape; every call site already exists).
- **Delete:** `scopeByGID`, `scopedLive`, `goID`, `goid_fast_go124.go`,
  `goid_fast_stub.go`, and the `goroutineLocalBindingStore`.

## 4. The NativeFn convention — the bounded "dig-in"

`Invoke([]Value)` is context-free, and **only ~17 builtins + the Scope ops**
need the context; the other hundreds are pure (`args → value`). So we do **not**
break `Invoke`. Add an *additive*, context-aware path:

```go
// Fn (optional, additive): default delegates to Invoke and ignores ec.
type CtxFn interface { InvokeCtx(ec *ExecContext, args []Value) (Value, error) }
```

- The VM, at every NativeFn call site, calls `InvokeCtx(currentEC, args)` when
  the fn implements `CtxFn`, else falls back to `Invoke(args)`. Pure builtins
  need **zero** changes.
- The migration is smaller than "~17 sites" thanks to #206's consolidation:
  - **IO family** (`*out*`/`*err*`, all of #207): thread `ec` through the
    three chokepoints — `WriteToOut(ec, s)`, `WriteToErr(ec, s)`,
    `resolveIOHandleVar(ec, name)` (read `ec.bindings` instead of ambient
    `Deref`). The ~5 print builtins that call them become `CtxFn` and forward
    `ec`. That's the whole IO surface.
  - **Non-IO** ambient readers (`*ns*`, `*compiling-aot*`, `*in*`, the
    `CurrentScope`/`CurrentContext` consumers) migrate individually — a small
    remaining set.

This keeps the change additive and the blast radius equal to the actual set of
context-dependent builtins, not the whole `Fn` surface.

## 5. The third surface — lowered (gogen_ir) code *is* an IR concern

`core_go_lowered/**` functions are native Go. When lowered let-go code reads a
dynamic var or calls a context-needing builtin, it too needs `ec`. So the
**lowering must thread `ec` through lowered function signatures**
(`lower_go.lg` emits `func f(ec *vm.ExecContext, …)` and forwards it at call
sites). This is the part that *is* an artifact of the IR system — and it's
where the lowering work and the VM work meet:

| Surface | Mechanism | IR artifact? |
|---|---|---|
| Interpreter eval loop | `ec` on the `Frame` | no (VM) |
| Hand-written NativeFns (~17) | `InvokeCtx` convention | no (VM) |
| Lowered gogen_ir functions | `ec` threaded through emitted signatures | **yes (lowering)** |

So the honest answer to "can the NativeFn path be deferred as an IR artifact?"
is **partly**: the lowered-code half is an IR-lowering change; the hand-written
NativeFn half is not, and gates the design.

## 6. #207 end-to-end under this design

1. `api.New(WithStdout(w))` → `Run` creates an `ExecContext` and pushes
   `*out* → w` onto `ec.bindings` for the Run's dynamic extent.
2. `println` (now a `CtxFn`) → `WriteToOut(ec, …)` → `resolveIOHandleVar(ec,
   "*out*")` reads `ec.bindings[*out*]` → `w`. (This is #206's chokepoint, now
   `ec`-threaded — no new resolution path is invented.)
3. Two runtimes Run concurrently on different goroutines → **distinct `ec`** →
   isolated, with no goid, no global, no leak.
4. mparrett's caveat disappears, and the `binding_frames.go:33` / `var.go` line
   refs in our #207 notes resolve once this lands (they're local-only today).

This requires #206 to be present (its `WriteToOut`/`*out*` routing is what
`WithStdout` binds against), which is why incorporating upstream is step 0.

## 7. Migration sequence

0. **Branch off clean `main@upstream` + duplicate the relevant code.** Base the
   workspace on `d6d0128` (already done) — so #206's print/`*out*` routing
   (`WriteToOut`/`WriteToErr`/`resolveIOHandleVar`, `LookupCoreVar`) is present
   by construction, no merge. Then duplicate the binding mechanism we want to
   refactor onto it — `binding_frames.go`, the `var.go` binding hooks, and the
   binding tests — **explicitly dropping `goid_fast_go124.go`/`goid_fast_stub.go`
   (the rejected hack)** and reconciling against main's process-global `var.go`.
   Verify it builds + `go test ./pkg/vm/...` green: that proves the feature is
   self-contained on clean main and gives a known-green baseline to refactor
   toward `ExecContext`.
1. Introduce `ExecContext` + thread it through the eval `Frame`. No behavior
   change yet: `ec.bindings` mirrors the existing store.
2. Add `CtxFn`/`InvokeCtx` with an ignoring default; wire the VM call sites.
3. Migrate the ~17 context-needing builtins to read `ec`.
4. Thread `ec` through lowered-fn signatures in `lower_go.lg`; `make generate`.
5. Make spawn sites propagate `ec` explicitly (`RunWithContext`).
6. Fold Scope into `ec`; migrate `CurrentScope`/`CurrentContext` readers.
7. Delete goid: `scopeByGID`, `scopedLive`, `goID`, `goid_fast_*.go`, the
   `goroutineLocalBindingStore`.
8. Verify: `pmap`/`parallel`/lgbgen concurrency + determinism preserved
   (`make check-generated`); `go test ./...`; the `binding-unwind` /
   non-local-exit cases (a workspace already exists for that).

## 8. Risks & open questions

- **Eval-loop coverage.** Every place a dynamic-var deref or NativeFn call
  happens must carry the `Frame`/`ec`. Audit the fast-path opcodes in
  `vm.go` (`execPrecompiled`, the dispatch around line 946) that bypass
  `NativeFn.Invoke`.
- **Lowered-signature churn.** Threading `ec` through `core_go_lowered/**`
  touches every lowered fn signature → large generated diff; must stay
  deterministic (check-generated).
- **Structured concurrency.** Scope's cancel/await tree is independent of how
  "current" is found, so folding it into `ec` is mechanical — but `with-scope`
  and detached goroutines must seed the child `ec.scope` correctly.
- **Perf.** A threaded pointer is cheaper than a `goID()` parse; expect neutral
  or better. Confirm on `var_deref_bench_test.go`.
- **Scope of v1.** Could land the binding half first (steps 1–5, 7-for-bindings)
  and fold Scope in later — but that leaves goid alive for Scope, so the
  "delete goid" win is only realized when both move. Recommend doing both.
- **Duplication-reconciliation surface (step 0).** We are *not* merging the
  60-commit branch; we duplicate the binding files onto clean main. The only
  reconciliation is our `var.go` binding hooks vs main's process-global
  `var.go` — a focused port, not a 60-commit merge. `test/io_binding_test.lg`
  (added by #206, already on the base) is the regression anchor for the
  eventual `ec`-threaded `WriteToOut`.
