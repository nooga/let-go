## VM performance and calling convention — Audit and Optimization Plan

This document summarizes performance bottlenecks in the VM interpreter and function calling convention, with concrete optimization steps, file pointers, and a rollout plan.

### Overview

- The interpreter loop lives in `pkg/vm/vm.go` (`Frame.Run`). Calls are performed via `OP_INVOKE` and `OP_TAIL_CALL` opcodes.
- Bytecode functions are represented by `*Func`/`*Closure` in `pkg/vm/func.go`; native functions by `*NativeFn` in `pkg/vm/native_func.go`.
- Current calling convention:
  - Caller slices out arguments from its stack (`mult`) and passes them to callee `Invoke`.
  - `Invoke` constructs a new `Frame` for bytecode functions; variadics package the rest into a `List`.
  - TCO exists only for `*Func` under `OP_TAIL_CALL`, not for `*Closure`.

### Hotspots and weak spots

- Frame and stack allocation per call

  - Each bytecode call (`Func.Invoke`/`Closure.Invoke`) allocates a new `Frame` and a `[]Value` stack sized to the callee’s `maxStack`.
  - Impact: GC pressure under call-heavy workloads.

- Argument slice retains caller’s stack backing array

  - In `OP_INVOKE`, `a := f.mult(...)` returns a slice into the caller’s `f.stack`. Passing `a` into `Func.Invoke`/`Closure.Invoke` causes the callee frame to reference the caller’s stack array via `args`, retaining it for the callee’s lifetime.
  - Impact: unwanted memory retention, increased heap live set.

- TCO limited to raw `*Func`

  - `OP_TAIL_CALL` only reuses the frame when the callee is `*Func`. For `*Closure`, it falls back to a normal call, missing TCO opportunities, especially in recursive closures.

- TCO growth branch still retains old stack

  - When reusing the frame in `TAIL_CALL` and growing the stack (`len(f.stack) < f.code.maxStack`), `f.args = a` assigns a slice backed by the previous frame’s stack array, retaining it.

- Generic invocation path overhead

  - All calls use a generic path (`nth`/`mult`, slice creation, dynamic dispatch). No specialized opcodes for tiny arities (0–3), which are very common.

- Reflection in `NativeFnType.Box`

  - Boxing arbitrary Go functions uses `reflect` to build args and make calls on every invocation. This is slow compared to direct function pointers.
  - Core runtime already prefers `Wrap`/`WrapNoErr`, which avoid reflect; ensure this consistently.

- Variadic rest packing
  - Variadic bytecode functions always allocate a `List` for rest args. This is semantically correct but adds allocation on hot variadics.

### Concrete optimization plan

Phase A — Correctness and retention fixes

- Copy args before entering bytecode callee frames

  - In `OP_INVOKE` and `OP_TAIL_CALL`, create a fresh small `[]Value` (or pooled) and copy the args slice `a` before passing to bytecode `Invoke` or assigning to `f.args`.
  - In `TAIL_CALL` growth branch, always copy `a` into `f.args` (do not assign the slice directly).

- Add TCO for closures
  - In `OP_TAIL_CALL`, handle `*Closure` similarly to `*Func`: reuse the frame while swapping `code`, `consts`, `closedOvers`, resetting `ip/sp`, and setting `args`.

Phase B — Allocation and dispatch overhead

- Pool frames and stacks

  - Introduce `sync.Pool` for `Frame` objects and `[]Value` stacks sized to common `maxStack` buckets. Provide `reset()` methods to clear state.
  - Ensure stacks are returned to pool at function return and not retained by arg slices (addressed by Phase A).

- Small-arity invoke opcodes

  - Add `OP_INVOKE_0/1/2/3` and `OP_TAIL_CALL_0/1/2/3` to avoid slicing/`nth`/`mult` overhead and repeated bounds checks.
  - Update compiler to emit specialized opcodes for calls up to arity 3.

- Ensure native functions use direct wrappers
  - Audit `rt/lang.go` and other registrations to confirm `NativeFnType.Wrap`/`WrapNoErr` are used instead of `Box` for hot paths.
  - Reserve `Box` for rare interop cases.

Phase C — Throughput improvements in iteration

- Reducer/transducer fast path
  - Add `Reducible` interface and prefer it in `reduce` for vectors/ranges/maps to avoid seq allocation.
  - Implement chunked seqs for vectors/ranges to reduce per-element overhead in `map`/`reduce` pipelines.

### Rollout and risk mitigation

- Start with Phase A: no API changes, reduces memory retention and enables more TCO.
- Phase B introduces new opcodes and pools; guarded by benchmarks to validate gains; maintain a fallback path.
- Phase C depends on collection refactors; schedule after `PersistentVector` correctness (see `docs/clojurelike-refactor-plan.md`).

### Benchmarks and validation

- Microbenchmarks:
  - Bytecode call overhead: arities 0–3 vs generic; with/without pools.
  - Tail-recursive function: with/without closure TCO.
  - Deep call chains: arg copy impact and GC.
- End-to-end:
  - Hot functional pipelines (`map`/`filter`/`reduce`) on vectors/ranges.
  - Native-heavy workloads to ensure `Wrap`/`WrapNoErr` paths are fast.

### Relevant files and pointers

- `pkg/vm/vm.go`

  - `Frame.Run`: `OP_INVOKE` and `OP_TAIL_CALL` arg handling, frame reuse, and memory retention points.
  - `Frame.mult`, `Frame.nth`, `Frame.push/pop/drop`: generic helpers used in hot paths.

- `pkg/vm/func.go`

  - `Func.Invoke` / `Closure.Invoke`: `NewFrame` allocation, variadic rest packing, and callee `args` lifetime.
  - `Func.MakeClosure`, `Closure` internals: info needed for closure TCO.

- `pkg/vm/native_func.go`

  - `NativeFnType.Box`: reflection-based boxing (slow, avoid in hot paths).
  - `NativeFnType.Wrap`/`WrapNoErr`: direct wrappers (fast path for runtime natives).

- `pkg/rt/lang.go`
  - Ensure all registered native functions use `Wrap`/`WrapNoErr`.

### Implementation checklist

- [ ] Copy args before bytecode calls and in tail-call growth branch.
- [ ] Extend `OP_TAIL_CALL` to reuse frame for `*Closure`.
- [ ] Introduce frame/stack pools with `sync.Pool` and integrate lifecycle.
- [ ] Add `INVOKE_0/1/2/3` and `TAIL_CALL_0/1/2/3`, update compiler emission and interpreter.
- [ ] Audit runtime natives to avoid `Box` in hot paths.
- [ ] Add microbenchmarks for calls/TCO; monitor allocations (pprof) and throughput.
- [ ] Later: add `Reducible` and chunked seqs when collection refactor lands.
