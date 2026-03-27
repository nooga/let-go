# let-go Development Guide

## Build & Test

```bash
go build ./...          # build all packages
go test ./... -count=1  # run all tests
go test ./test/ -count=1 -v  # run .lg integration tests with output
go test ./pkg/vm/ -run TestPersistentMap -v  # run specific Go tests
go test ./pkg/compiler/ -bench=. -benchmem   # run benchmarks
```

## Testing Policy

**Always add tests for new features.** Do NOT rely on `go run . -e '...'` for verification — write proper test files instead.

- **Go-level unit tests**: `pkg/vm/*_test.go` for VM types, data structures, opcodes
- **Integration tests**: `test/*_test.lg` for language-level features (macros, builtins, seq ops, etc.)
- **Benchmarks**: `pkg/vm/bench_test.go` and `pkg/compiler/bench_test.go`

### Writing .lg tests

Tests use the `test` namespace with `deftest`, `testing`, and `is`:

```clojure
(ns test.my-feature-test
  (:require [test :refer :all]))

(deftest my-feature
  (testing "basic behavior"
    (is (= expected (my-fn input)))))
```

The test runner (`test/language_test.go`) picks up all `*.lg` files in `test/` automatically.

### When to add which kind of test

- New VM type or data structure → Go unit test in `pkg/vm/`
- New builtin function → .lg integration test in `test/`
- New macro or language feature → .lg integration test in `test/`
- Performance-sensitive change → benchmark in `bench_test.go`

## Project Structure

- `pkg/vm/` — VM types, opcodes, frame execution, persistent data structures
- `pkg/compiler/` — reader (lexer/parser), compiler (bytecode generation), eval
- `pkg/rt/` — runtime: builtins (lang.go), core library (core/core.lg), HTTP, JSON, OS
- `pkg/errors/` — error base types
- `pkg/resolver/` — namespace file loader
- `pkg/api/` — public Go API
- `test/` — .lg integration tests
- `examples/` — example programs

## Architecture Notes

- Values implement `vm.Value` interface (Type, Unbox, String)
- Seq protocol: `Next()` returns Go `nil` at end of sequence, `More()` returns `EmptyList`
- `seqOf()` in lang.go coerces any value to a Seq (prefers Sequable.Seq() for non-lazy types)
- Frame pooling via `sync.Pool` — call `ReleaseFrame(f)` after `f.Run()`
- Small int cache (-128..255) via `MakeInt()` to avoid interface boxing
- HAMT persistent maps/sets — `PersistentMap` and `PersistentSet` backed by hash-array mapped tries
- `Hashable` interface for cached/efficient hashing — all primitive types implement it
- Protocols dispatch on `ValueType` — `extend-type` registers in a map from ValueType to impl
- Multimethods dispatch on arbitrary function return value
- try/catch uses handler stack on Frame, `ThrownError` propagates through error returns
- `thrownPanic` for errors crossing native function boundaries (lazy seq realization)
- Error reporting: `FormSource` maps forms to source locations, `SourceMap` on CodeChunk

## Known Limitations

### Recursive lazy seq stack overflow
Lazy sequences built with `(cons val (lazy-seq ...))` patterns overflow the Go stack when deeply chained. This affects `repeatedly`, `iterate` + `take` + `vec`, and deep `filter`/`map` chains on infinite sequences.

**Root cause:** `Cons.Next()` realizes `LazySeq` tails on the Go call stack. Each level adds ~500 bytes. Go's 1GB stack limit means ~2M levels max, but in practice the overhead per level (thunk invocation via `Func.Invoke → Frame.Run`) is much higher.

**What works:**
- `(first (iterate inc 0))`, `(first (next (iterate inc 0)))` — manual traversal
- `(vec (take 5 (range 1000)))` — `range` is not lazy-seq based
- `(vec (filter even? (range 100)))` — finite range sources
- `(take 5 (for [x (range 1000)] (* x x)))` — for with finite range

**What overflows:**
- `(vec (take 5 (filter even? (iterate inc 0))))` — lazy chain on infinite iterate
- `(vec (repeatedly 100 f))` — repeatedly uses recursive lazy-seq
- `(count (repeatedly 5 f))` — forces full realization via RawCount

**Workaround:** Use `range` instead of `iterate` for finite sequences. Use `loop`/`recur` for accumulation instead of `vec` on lazy seqs.

**Proper fix:** Implement trampoline-style `LazySeq` that doesn't accumulate Go stack frames during realization. This requires restructuring `Cons.Next()` to not eagerly realize LazySeq tails.

## Common Gotchas

- `and`/`or` are macros (short-circuiting), not functions
- `Cons.RawCount()` is O(n) — don't call on potentially infinite seqs
- `LazySeq.RawCount()` forces full realization — never call on infinite seqs
- The `seq` and `map` builtins must skip `RawCount()` for Cons/LazySeq types
- Vector destructuring uses `nth` (works on any seq), not `get` (vectors only)
- Type objects (IntType, StringType, etc.) are singleton pointers — comparable with `==`
- `concat` is eager (to avoid issues with quasiquote at compile time)
