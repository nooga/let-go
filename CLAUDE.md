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

## Common Gotchas

- `and`/`or` are macros (short-circuiting), not functions
- `Cons.RawCount()` is O(n) — don't call on potentially infinite seqs
- `LazySeq.RawCount()` forces full realization — never call on infinite seqs
- The `seq` builtin must skip `RawCount()` for Cons/LazySeq types
- Vector destructuring uses `nth` (works on any seq), not `get` (vectors only)
- Type objects (IntType, StringType, etc.) are singleton pointers — comparable with `==`
- `concat` is eager (to avoid issues with quasiquote at compile time)
