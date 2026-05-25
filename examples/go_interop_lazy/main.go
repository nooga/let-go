// Worked example: a Go program embeds let-go, then consumes
// potentially-infinite let-go sequences from Go using the new vm.Iter
// (Go 1.23+ iter.Seq[Value]) and vm.SeqToSlice helpers.
//
// Demonstrates:
//   - Early-break over an infinite sequence (primes)
//   - Bounded realization with SeqToSlice (Fibonacci)
//   - Bidirectional fns: let-go calling a Go fn inside a lazy pipeline
//   - Passing a (take N infinite-seq) to a Go fn taking []int
//
// Run: go run ./examples/go_interop_lazy
package main

import (
	"fmt"
	"log"

	"github.com/nooga/let-go/pkg/api"
	"github.com/nooga/let-go/pkg/vm"
)

func main() {
	c, err := api.NewLetGo("demo")
	if err != nil {
		log.Fatal(err)
	}

	demo1_InfiniteSeqEarlyBreak(c)
	demo2_BoundedRealization(c)
	demo3_GoFnInsideLazyPipeline(c)
	demo4_TakeFromInfiniteIntoGoSlice(c)
	demo5_FilterByGoPredicate(c)
}

// Demo 1: Iterate an infinite let-go sequence from Go, breaking out
// when a condition holds. vm.Iter preserves laziness — the seq isn't
// realized past what the consumer pulls.
func demo1_InfiniteSeqEarlyBreak(c *api.LetGo) {
	fmt.Println("=== Demo 1: infinite seq, early break ===")

	// A let-go expression that defines an infinite lazy seq of primes.
	// api.Run compiles a single form, so wrap multiple defs in (do …).
	_, err := c.Run(`(do
		(defn divisible? [n d] (zero? (mod n d)))
		(defn prime? [n]
		  (and (> n 1)
		       (not-any? #(divisible? n %) (range 2 n))))
		(def primes (filter prime? (iterate inc 2))))`)
	if err != nil {
		log.Fatal(err)
	}

	v, err := c.Run("primes")
	if err != nil {
		log.Fatal(err)
	}

	// Pull primes from Go until we find one above 1000. vm.Iter yields
	// values one at a time; breaking stops further realization.
	var found []vm.Value
	for p := range vm.Iter(v) {
		found = append(found, p)
		if int64(p.(vm.Int)) > 1000 {
			break
		}
	}
	fmt.Printf("  pulled %d primes, last = %v\n\n", len(found), found[len(found)-1])
}

// Demo 2: Use SeqToSlice when you know the seq is bounded and want a
// concrete []Value to work with in Go.
func demo2_BoundedRealization(c *api.LetGo) {
	fmt.Println("=== Demo 2: bounded realization with SeqToSlice ===")

	_, err := c.Run(`(defn fibs []
		(map first (iterate (fn [[a b]] [b (+ a b)]) [0 1])))`)
	if err != nil {
		log.Fatal(err)
	}

	v, err := c.Run("(take 10 (fibs))")
	if err != nil {
		log.Fatal(err)
	}

	first10, err := vm.SeqToSlice(v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  first 10 Fibonacci: %v\n\n", first10)
}

// Demo 3: Register a Go function in let-go, then call it inside a lazy
// pipeline. Demonstrates that the Go fn is invoked once per element as
// the consumer pulls — not eagerly at composition time.
func demo3_GoFnInsideLazyPipeline(c *api.LetGo) {
	fmt.Println("=== Demo 3: Go fn called from let-go in a lazy pipeline ===")

	calls := 0
	if err := c.Def("go-square", func(n int) int {
		calls++
		return n * n
	}); err != nil {
		log.Fatal(err)
	}

	// Build the lazy pipeline. Go-square is referenced inside `map` but
	// not invoked yet — laziness defers the call.
	v, err := c.Run("(map go-square (iterate inc 0))")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  go-square invocations before pulling: %d\n", calls)

	// Pull 5 elements via vm.Iter. The Go fn fires exactly 5 times.
	var squares []vm.Value
	count := 0
	for x := range vm.Iter(v) {
		squares = append(squares, x)
		count++
		if count >= 5 {
			break
		}
	}
	fmt.Printf("  first 5 squares: %v\n", squares)
	fmt.Printf("  go-square invocations after pulling 5: %d\n\n", calls)
}

// Demo 4: Pass a (take N infinite-seq) to a Go function expecting []int.
// The reflect bridge in pkg/vm/native_func.go auto-realizes the seq
// into a []vm.Value, then does per-element conversion via the existing
// struct_mapping machinery — so Go-side gets a real []int, not a seq.
func demo4_TakeFromInfiniteIntoGoSlice(c *api.LetGo) {
	fmt.Println("=== Demo 4: (take N infinite) → Go fn([]int) ===")

	if err := c.Def("sum-ints", func(xs []int) int {
		total := 0
		for _, x := range xs {
			total += x
		}
		return total
	}); err != nil {
		log.Fatal(err)
	}

	v, err := c.Run("(sum-ints (take 100 (iterate inc 1)))")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  sum of 1..100 = %v   (expected 5050)\n\n", v.Unbox())
}

// Demo 5: A Go-defined predicate threaded through let-go's lazy filter,
// consumed back in Go via vm.Iter. The predicate is called once per
// element as the consumer pulls.
func demo5_FilterByGoPredicate(c *api.LetGo) {
	fmt.Println("=== Demo 5: Go predicate inside let-go filter, consumed from Go ===")

	if err := c.Def("go-prime?", isPrime); err != nil {
		log.Fatal(err)
	}

	v, err := c.Run("(filter go-prime? (iterate inc 2))")
	if err != nil {
		log.Fatal(err)
	}

	// Pull the first 5 primes found by the Go predicate.
	var primes []vm.Value
	count := 0
	for p := range vm.Iter(v) {
		primes = append(primes, p)
		count++
		if count >= 5 {
			break
		}
	}
	fmt.Printf("  first 5 primes via Go predicate: %v\n", primes)
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}
