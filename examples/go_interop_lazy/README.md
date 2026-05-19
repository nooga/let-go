# Lazy sequence interop between Go and let-go

A worked example of consuming let-go sequences from Go code using
`vm.Iter` (Go 1.23+ `iter.Seq[Value]`) and `vm.SeqToSlice`, and of
threading Go functions through let-go lazy pipelines.

Run:

```sh
go run ./examples/go_interop_lazy
```

Expected output:

```
=== Demo 1: infinite seq, early break ===
  pulled 169 primes, last = 1009

=== Demo 2: bounded realization with SeqToSlice ===
  first 10 Fibonacci: [0 1 1 2 3 5 8 13 21 34]

=== Demo 3: Go fn called from let-go in a lazy pipeline ===
  go-square invocations before pulling: 0
  first 5 squares: [0 1 4 9 16]
  go-square invocations after pulling 5: 5

=== Demo 4: (take N infinite) → Go fn([]int) ===
  sum of 1..100 = 5050   (expected 5050)

=== Demo 5: Go predicate inside let-go filter, consumed from Go ===
  first 5 primes via Go predicate: [2 3 5 7 11]
```

## Why it works

- `vm.Iter(v)` returns `iter.Seq[Value]`. It's pull-based: each step
  realizes exactly one element of the underlying `LazySeq`. Breaking
  out of the loop stops further realization, so infinite sequences are
  safe to consume.

- `vm.SeqToSlice(v)` realizes a known-bounded sequence into a
  `[]vm.Value`. Use it when you've already taken a finite slice
  (e.g. `(take 10 …)`) and want concrete Go data.

- When a Go function declared with `c.Def` is referenced inside a
  let-go lazy pipeline like `(map go-fn (iterate inc 0))`, the Go fn
  is invoked once per element as the consumer pulls — never eagerly
  at composition time.

- When let-go passes a lazy sequence to a Go function taking `[]T`,
  the runtime auto-realizes the sequence and converts elements via
  the existing `unboxSliceInto` machinery. The Go fn receives a real
  `[]int` (or whatever element type), not a sequence.
