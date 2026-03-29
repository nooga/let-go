package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// eval compiles and runs a let-go expression, returning the result.
func eval(b *testing.B, src string) vm.Value {
	b.Helper()
	consts := vm.NewConsts()
	ctx := NewCompiler(consts, rt.NS(rt.NameCoreNS))
	chunk, result, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		b.Fatal(err)
	}
	_ = chunk
	return result
}

// benchEval runs src b.N times, measuring compile+execute.
func benchEval(b *testing.B, src string) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		consts := vm.NewConsts()
		ctx := NewCompiler(consts, rt.NS(rt.NameCoreNS))
		_, _, err := ctx.CompileMultiple(strings.NewReader(src))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// benchExec compiles once, then runs b.N times (execution only).
func benchExec(b *testing.B, src string) {
	b.Helper()
	consts := vm.NewConsts()
	ctx := NewCompiler(consts, rt.NS(rt.NameCoreNS))
	chunk, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := vm.NewFrame(chunk, nil)
		f.Run()
	}
}

// ============================================================================
// Init — core.lg loading: source compilation vs precompiled bytecode
// ============================================================================

func BenchmarkInitFromSource(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := vm.NewConsts()
		ctx := NewCompiler(c, rt.NS(rt.NameCoreNS))
		ctx.SetSource("<core>")
		_, _, err := ctx.CompileMultiple(strings.NewReader(rt.CoreSrc))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInitFromLGB(b *testing.B) {
	if len(rt.CoreCompiledLGB) == 0 {
		b.Skip("no precompiled core_compiled.lgb")
	}
	resolve := func(ns, name string) *vm.Var {
		n := rt.NS(ns)
		v := n.Lookup(vm.Symbol(name))
		if v == vm.NIL {
			return n.Def(name, vm.NIL)
		}
		return v.(*vm.Var)
	}
	for i := 0; i < b.N; i++ {
		unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(rt.CoreCompiledLGB), resolve)
		if err != nil {
			b.Fatal(err)
		}
		f := vm.NewFrame(unit.MainChunk, nil)
		_, err = f.RunProtected()
		vm.ReleaseFrame(f)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ============================================================================
// Arithmetic — tight numeric loops
// ============================================================================

func BenchmarkArithmetic(b *testing.B) {
	b.Run("sum-loop-100", func(b *testing.B) {
		benchExec(b, `(loop [i 0 sum 0] (if (< i 100) (recur (+ i 1) (+ sum i)) sum))`)
	})
	b.Run("sum-loop-10000", func(b *testing.B) {
		benchExec(b, `(loop [i 0 sum 0] (if (< i 10000) (recur (+ i 1) (+ sum i)) sum))`)
	})
	b.Run("fib-25", func(b *testing.B) {
		benchExec(b, `
		(do
			(defn fib [n]
				(if (< n 2) n (+ (fib (- n 1)) (fib (- n 2)))))
			(fib 25))`)
	})
}

// ============================================================================
// Function calls — compiled fn, closure, native
// ============================================================================

func BenchmarkFunctionCalls(b *testing.B) {
	b.Run("direct-call-loop", func(b *testing.B) {
		benchExec(b, `
		(do
			(defn add1 [x] (+ x 1))
			(loop [i 0] (if (< i 10000) (do (add1 i) (recur (+ i 1))) i)))`)
	})
	b.Run("closure-call-loop", func(b *testing.B) {
		benchExec(b, `
		(do
			(defn make-adder [n] (fn [x] (+ x n)))
			(let [add5 (make-adder 5)]
				(loop [i 0] (if (< i 10000) (do (add5 i) (recur (+ i 1))) i))))`)
	})
	b.Run("native-call-loop", func(b *testing.B) {
		benchExec(b, `(loop [i 0] (if (< i 10000) (recur (+ i 1)) i))`)
	})
	b.Run("higher-order", func(b *testing.B) {
		benchExec(b, `
		(do
			(defn apply-twice [f x] (f (f x)))
			(loop [i 0 x 0] (if (< i 1000) (recur (+ i 1) (apply-twice inc x)) x)))`)
	})
}

// ============================================================================
// Collection operations
// ============================================================================

func BenchmarkCollections(b *testing.B) {
	b.Run("vector-conj-100", func(b *testing.B) {
		benchExec(b, `(loop [v [] i 0] (if (< i 100) (recur (conj v i) (+ i 1)) v))`)
	})
	b.Run("list-conj-100", func(b *testing.B) {
		benchExec(b, `(loop [l () i 0] (if (< i 100) (recur (conj l i) (+ i 1)) l))`)
	})
	b.Run("map-assoc-100", func(b *testing.B) {
		benchExec(b, `(loop [m {} i 0] (if (< i 100) (recur (assoc m i (* i 10)) (+ i 1)) m))`)
	})
	b.Run("reduce-vector-1000", func(b *testing.B) {
		benchExec(b, `(reduce + 0 (vec (range 1000)))`)
	})
	b.Run("reduce-list-1000", func(b *testing.B) {
		benchExec(b, `(reduce + 0 (range 1000))`)
	})
}

// ============================================================================
// Seq operations — map, filter, take with lazy evaluation
// ============================================================================

func BenchmarkSeqOps(b *testing.B) {
	b.Run("map-small-vec", func(b *testing.B) {
		benchExec(b, `(vec (map inc [1 2 3 4 5 6 7 8 9 10]))`)
	})
	b.Run("map-large-range", func(b *testing.B) {
		benchExec(b, `(reduce + 0 (map inc (range 1000)))`)
	})
	b.Run("filter-range", func(b *testing.B) {
		benchExec(b, `(reduce + 0 (filter even? (range 1000)))`)
	})
	b.Run("take-from-range", func(b *testing.B) {
		benchExec(b, `(vec (take 100 (range 10000)))`)
	})
	b.Run("lazy-chain", func(b *testing.B) {
		benchExec(b, `(reduce + 0 (take 100 (filter even? (map inc (range 1000)))))`)
	})
}

// ============================================================================
// String operations
// ============================================================================

func BenchmarkStrings(b *testing.B) {
	b.Run("str-concat-100", func(b *testing.B) {
		benchExec(b, `(loop [s "" i 0] (if (< i 100) (recur (str s "x") (+ i 1)) s))`)
	})
	b.Run("split-join", func(b *testing.B) {
		benchExec(b, `(apply str (interpose "," (split "a,b,c,d,e,f,g,h" ",")))`)
	})
}

// ============================================================================
// Recursion patterns
// ============================================================================

func BenchmarkRecursion(b *testing.B) {
	b.Run("tail-recur-10000", func(b *testing.B) {
		benchExec(b, `(loop [i 0] (if (< i 10000) (recur (+ i 1)) i))`)
	})
	b.Run("non-tail-fib-20", func(b *testing.B) {
		benchExec(b, `
		(do
			(defn fib [n]
				(if (< n 2) n (+ (fib (- n 1)) (fib (- n 2)))))
			(fib 20))`)
	})
	b.Run("mutual-even-odd-1000", func(b *testing.B) {
		benchExec(b, `
		(do
			(declare my-even? my-odd?)
			(defn my-even? [n] (if (= n 0) true (my-odd? (- n 1))))
			(defn my-odd? [n] (if (= n 0) false (my-even? (- n 1))))
			(my-even? 1000))`)
	})
}
