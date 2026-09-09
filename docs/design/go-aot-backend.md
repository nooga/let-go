---
status: active
last-verified: 2026-06-05
authoritative-for:
  - go-aot-backend-design
  - two-tier-aot-approach
  - aot-compile-driver
human-verified:
---

## Go AOT backend — compiling let-go to Go while keeping the runtime

This document proposes a second backend that compiles let-go code to Go, preserving the runtime (`pkg/vm` and `pkg/rt`) semantics and interop. Dynamic `eval` continues to use the VM and targets the same `Var`s, enabling mixed compiled+interpreted systems and ultimate AOT via the Go toolchain.

### Goals

- Preserve let-go semantics (vars, dynamic redefinition, seq/coercion, errors) by keeping the same runtime types and APIs.
- Allow compiled functions to call interpreted functions and vice versa via shared `Var` indirection.
- Reduce interpreter overhead by lowering bytecode to direct Go code where profitable.
- Keep the developer workflow familiar: `lg build -target go` produces a Go module/binary with the same behavior.

### Two-tier approach (staged)

1. Embedding tier (quick win): emit Go that embeds const pools and bytecode as Go data and registers `*Func`/`*Closure` into `Var`s at `init()`. This removes image I/O and lets the Go linker/tinyalloc help, but still uses the interpreter loop. Comparable to a built-in runtime image.

2. Native lowering tier (performance): compile let-go functions to Go functions that implement the same logic without the bytecode interpreter loop. Use the same `Value`/`Fn` interfaces so the VM and runtime can call them normally.

### Interop model (critical)

- Vars as the boundary: compiled code resolves `Var`s by symbol and either:
  - Reads the root on each invocation (full dynamism, slower), or
  - Captures the `*Var` pointer and `Deref()`s on each call (fast indirection), or
  - Inlines a constant value for `^:const` or explicitly annotated `^:inline` vars (AOT-only optimization).
- Register compiled functions as `NativeFn` implementations (or a dedicated `CompiledFn`) so interpreted code can call them via normal `invoke`.
- Dynamic `eval` continues to compile to bytecode and install roots in the same `Var`s. Compiled code that dereferences vars at call-time observes redefinitions.

### Code generation strategy

- Use `go/ast` + `go/printer` to build safe Go code with imports and namespacing.
- For each let-go namespace, generate a Go package (or a single `lggen_app` package) with:
  - Constants/const-pool reified as Go values (scalars/collections using existing constructors).
  - A function per let-go function:
    - Signature: `func f(env *vm.Env, args []vm.Value) (vm.Value, error)` or a small fixed-arity wrapper when known.
    - Implement logic by translating bytecode ops to structured Go control flow and direct calls to runtime helpers.
    - Self tail-calls lowered to loops to maintain TCO semantics.
  - `init()` that registers these as `NativeFn`s in their `Var`s and sets metadata.

### Lowering details

- Stack elimination: represent VM stack slots as local Go variables; allocate `[]vm.Value` only when necessary (variadics/rest).
- Control flow: map bytecode basic blocks to Go `for`/`switch`/`goto` (prefer structured loops; avoid `goto` unless necessary for performance clarity).
- Tail-call optimization:
  - Self-tail recursion: transform to `for { ... update args ... continue }`.
  - Mutual tail recursion: optional trampoline helper (lower priority, can remain dynamic via VM path initially).
- Calls:
  - Direct calls to known compiled functions via Go calls when target is statically known and annotations allow devirtualization.
  - Otherwise, call through `Var` indirection: `v := ns.Resolve("sym"); fn := v.Deref().(vm.Fn); return fn.Invoke(args)`.
- Numerics: continue to use `vm.Int` fast paths per `value-representation-and-numeric-performance.md`.

### Error handling and stack traces

- Wrap errors with source locations where possible; generate a lightweight mapping from function/pc to source span for better error messages.
- Include function names and ns in Go function names for recognizability in profiles (e.g., `lg_ns_core_inc`).

### Build and tooling

- CLI: `lg build -target go [-o outdir]` emits Go files and optionally a `go.mod` if building a standalone binary.
- Integration modes:
  - Library: emit a Go package that can be `import`ed and `init()` registers vars.
  - Binary: emit `main.go` that installs stdlib (image or compiled) and application ns, then runs `-m your.main/ns -f your-fn`.

### Driving lowering: `lg.compiler` and the `lg-compile` shim

Whole-program lowering is driven by `lg.compiler` (`pkg/rt/core/lg/compiler.lg`): it parses each input into a lowering spec, orders the specs by their requires, evaluates only definitional forms so cross-namespace references resolve, hands the set to `ir.passes.pipeline/lower-all-ns-to-go-result`, and writes the emitted packages. Entry-frame emission is opt-in.

It is an ordinary embedded core namespace, so a released binary drives AOT lowering with no let-go checkout on disk, and edits to it churn `pkg/rt/generated.sums` like any other core source.

**Core-embedded is not the same as bundled.** Build-time tools are embedded as source but excluded from `core_compiled.lgb`, and load on demand when something requires them — `cmd/lgbgen.isBundleSkippedTool` holds the list. Two families qualify today, the `ir.*` pipeline and `lg.compiler.*`, for the same reason: a plain `lg` script never touches them, so decoding them at every process start is pure cost. It is not a small cost. `lg.compiler` is 30 chunks on its own, but its requires pull the whole IR pipeline and `gogen` into the const pool: bundling it takes `core_compiled.lgb` from 311,302 to 1,097,639 bytes (937 to 2,164 chunks) and costs 3.2 to 3.4 times the boot, 1.9 ms to 6.0 ms median on a quiet host and 2.7 ms to 9.2 ms on a loaded one, against an 8 ms budget (interleaved, 25 pairs per condition, 2026-09-08). That fits the budget on a quiet host and flakes the gate on a loaded one, so the skip is a cost decision rather than a hard failure. Either is acceptable to the project (#735); the skip is kept because the cost buys a plain script nothing. The skip is scoped to the bundle alone — the Go-lowering path must still see these namespaces, so `deriveGoLoweringOrder` does not consult it.

**`scripts/lg-compile` is a compatibility surface, not an implementation.** External orchestrators invoke it by path and read its output, so four things are contract: the path, the positional argv, the printed lines, and the exit code. In particular `EMIT-FAIL <path> pkg=<pkg> returned <type>` is scraped from stdout, and an EMIT-FAIL is deliberately not an exit failure — a function that does not lower stays on the trampoline and the program still runs, so partial lowering is a normal outcome the orchestrator judges for itself. Only a fatal error exits non-zero. `lg.compiler/exit-code` owns that policy and `test/e2e` pins both halves.

The intended destination for this layering, including a thin `lg.commands.compile` CLI adapter and the move of `ir.*` under `lg.compiler.ir.*`, is tracked in #786.

### MVP scope (Native lowering tier)

- Support non-variadic functions, fixed arity 0–3, no closures.
- Self-tail recursion elimination.
- Calls via `Var` indirection; no direct-call devirtualization initially.
- Register compiled functions into `Var`s; interpreted code calls them as `NativeFn`s.
- Basic error wrapping with function name and ns.

### Extended scope

- Closures: compile to Go structs capturing closed-overs and implementing `vm.Fn`.
- Variadics/rest args: specialized wrappers to minimize allocations; align with VM small-arity opcodes.
- Direct-call devirtualization when target var is known final (annotated or compiler-proven).
- Source maps and improved stack traces.

### Acceptance criteria

- Mixed execution works: compiled functions can call interpreted functions and vice versa using the same `Var`s.
- Performance: microbenchmarks show clear wins over VM interpreter for hot functions (call overhead, tight loops, tail recursion).
- Semantics preserved: dynamic var redefinition observed by compiled code when not inlined.

### Risks and mitigations

- Code size and compile times: mitigate via per-ns packages and incremental builds.
- Divergence from VM semantics: maintain a single spec and comprehensive tests; reuse runtime helpers for complex ops.
- TCO parity: self-tail elimination first; document limits for mutual recursion; keep VM path as fallback.
- Debuggability: provide source location attachments and recognizable symbols; later add mapping files.

### Relationship to runtime images

- The embedding tier is a Go-native alternative to runtime images: zero I/O and simple distribution.
- Native lowering tier provides maximum AOT performance while keeping the runtime semantics intact.
