## Value representation and numeric performance

This note documents how values (especially numbers) are represented today, performance implications in Go, and actionable optimizations. It complements the VM and collections plans.

### Current representation

- Core interfaces in `pkg/vm/value.go`:

  - `Value` with methods `Type() ValueType`, `Unbox() interface{}`, and `String()`.
  - Numeric type `Int` is defined in `pkg/vm/int.go` as `type Int int` and implements `Value` via value receivers.

- Boxing/unboxing paths:
  - `Int.Unbox()` returns `int(l)`; `Int.Type()` returns shared `IntType`.
  - `BoxValue` (`value.go`) uses reflection to convert unknown Go values into `Value` instances; it routes `reflect.Int` to `IntType.Box`.
  - Runtime numeric ops in `pkg/rt/lang.go` often use `v.Unbox().(int)` per operand.

### Performance characteristics in Go

- Storing `Int` in an interface does not allocate; it is stored directly in the interface data word.
- Method calls on value receivers (`Type`, `Unbox`, `String`) are cheap; however, repeated `Unbox`+type-assert in tight loops adds avoidable overhead.
- `Int.String()` currently uses `fmt.Sprintf` which is slower than `strconv.Itoa`.
- Reflection in `BoxValue` is expensive and should be avoided on hot paths; OK at interop boundaries.

### Weak spots and improvements

- Arithmetic via `Unbox` in `rt/lang.go`:

  - Today: `n += vs[i].Unbox().(int)` for many core functions (`+`, `-`, `*`, `/`, `mod`, comparisons).
  - Improvement: directly assert to `Int` and convert: `n += int(vs[i].(vm.Int))`. This removes an interface method call and a type assertion on the result.

- Integer to string

  - Today: `fmt.Sprintf("%d", int(l))` in `Int.String()`.
  - Improvement: `strconv.Itoa(int(l))`.

- Native numeric functions

  - Ensure all hot natives (in `pkg/rt/lang.go`) use direct `Int` asserts and avoid `Unbox`. For mixed-type inputs, keep clear errors, but keep the happy path fast.

- Boxing at boundaries

  - Prefer `NativeFnType.Wrap`/`WrapNoErr` to register natives to avoid reflection per call. Reserve `NativeFnType.Box` for rare interop where a plain Go func must be adapted.
  - Avoid calling `BoxValue` in loops; convert inputs once before iteration if needed.

- Future: structural equality/hash
  - When adding `Equiv`/`Hash`, keep fast-path numeric equality (`Int` vs `Int`) and hash (`int(l)`) simple and inlined.

### Relevant files and pointers

- `pkg/vm/value.go`: core `Value` interfaces and `BoxValue` reflection-based converter.
- `pkg/vm/int.go`: `Int` implementation; switch `String` to `strconv.Itoa`.
- `pkg/rt/lang.go`: numeric primitives (`+`, `-`, `*`, `/`, `mod`, `gt`, `lt`, `max`, `min`), replace `Unbox` paths with direct `Int` assertions.
- `pkg/vm/native_func.go`: prefer `Wrap`/`WrapNoErr` over `Box` in hot paths.

### Optimization checklist

- [ ] Change `Int.String()` to use `strconv.Itoa`.
- [ ] Replace `Unbox().(int)` with direct `Int` assertions in all numeric natives in `rt/lang.go`.
- [ ] Audit other numeric consumers (e.g., range, comparisons) for the same pattern.
- [ ] Ensure all runtime natives are registered via `Wrap`/`WrapNoErr` (no `Box` on hot paths).
- [ ] Add microbenchmarks: arithmetic loops over vectors; compare before/after `Unbox` removal and `String` change.

### Benchmarks to add

- `BenchmarkAddInts` operating over `[]Value{Int(...)}` using current and optimized paths.
- String conversion benchmark for `Int.String()` before vs after `strconv.Itoa`.
