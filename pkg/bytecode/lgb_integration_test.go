package bytecode_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// captureStdout runs fn and returns everything written to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

// compileSource compiles src and returns the compiled chunk + consts.
func compileSource(t *testing.T, src string) (*vm.CodeChunk, *vm.Consts) {
	t.Helper()
	consts := vm.NewConsts()
	ctx := compiler.NewCompiler(consts, rt.NS("user"))
	rt.SetNSLoader(resolver.NewNSResolver(ctx, []string{"."}))
	ctx.SetSource("<test>")
	chunk, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	return chunk, consts
}

// runSource compiles+runs src and returns stdout.
func runSource(t *testing.T, src string) string {
	t.Helper()
	consts := vm.NewConsts()
	ctx := compiler.NewCompiler(consts, rt.NS("user"))
	rt.SetNSLoader(resolver.NewNSResolver(ctx, []string{"."}))
	ctx.SetSource("<test>")
	return captureStdout(t, func() {
		_, _, err := ctx.CompileMultiple(strings.NewReader(src))
		if err != nil {
			t.Fatalf("source run failed: %v", err)
		}
	})
}

// compileThenRunLGB compiles src to LGB, decodes+runs the LGB, returns stdout.
func compileThenRunLGB(t *testing.T, src string) string {
	t.Helper()
	chunk, consts := compileSource(t, src)

	// Encode to LGB
	var buf bytes.Buffer
	if err := bytecode.EncodeCompilation(&buf, consts, chunk); err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	lgbData := buf.Bytes()

	// Set up resolver for the decode+run phase
	freshConsts := vm.NewConsts()
	ctx := compiler.NewCompiler(freshConsts, rt.NS("user"))
	rt.SetNSLoader(resolver.NewNSResolver(ctx, []string{"."}))

	// Decode
	resolve := func(nsName, name string) *vm.Var {
		n := rt.DefNSBare(nsName)
		v := n.LookupLocal(vm.Symbol(name))
		if v == nil {
			return n.Def(name, vm.NIL)
		}
		return v
	}
	unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(lgbData), resolve)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	// Run
	return captureStdout(t, func() {
		f := vm.NewFrame(unit.MainChunk, nil)
		_, err := f.RunProtected()
		vm.ReleaseFrame(f)
		if err != nil {
			t.Fatalf("lgb run failed: %v", err)
		}
	})
}

func assertSourceMatchesLGB(t *testing.T, name, src string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		srcOut := runSource(t, src)
		lgbOut := compileThenRunLGB(t, src)
		if srcOut != lgbOut {
			t.Errorf("output mismatch:\nsource: %q\nlgb:    %q", srcOut, lgbOut)
		}
	})
}

func TestLGBParity(t *testing.T) {
	assertSourceMatchesLGB(t, "hello-world", `
		(println "hello world")
	`)

	assertSourceMatchesLGB(t, "arithmetic", `
		(println (+ 1 2 3))
		(println (* 4 5))
		(println (- 10 3))
	`)

	assertSourceMatchesLGB(t, "def-and-fn", `
		(def x 42)
		(defn add [a b] (+ a b))
		(println x)
		(println (add 10 20))
	`)

	assertSourceMatchesLGB(t, "multi-arity-fn", `
		(defn greet
			([name] (str "Hello, " name))
			([greeting name] (str greeting ", " name)))
		(println (greet "world"))
		(println (greet "Howdy" "partner"))
	`)

	assertSourceMatchesLGB(t, "closures", `
		(defn make-counter []
			(let [n (atom 0)]
				(fn [] (swap! n inc) @n)))
		(let [c (make-counter)]
			(println (c))
			(println (c))
			(println (c)))
	`)

	assertSourceMatchesLGB(t, "loop-recur", `
		(println (loop [i 0 sum 0]
			(if (>= i 10) sum (recur (inc i) (+ sum i)))))
	`)

	assertSourceMatchesLGB(t, "higher-order", `
		(println (vec (map inc [1 2 3])))
		(println (vec (filter even? (range 10))))
		(println (reduce + 0 (range 5)))
	`)

	assertSourceMatchesLGB(t, "lazy-seqs", `
		(println (vec (take 5 (iterate #(* 2 %) 1))))
		(println (vec (take 3 (drop 5 (range)))))
	`)

	assertSourceMatchesLGB(t, "destructuring", `
		(let [{:keys [a b] :or {b 99}} {:a 1}]
			(println a b))
		(let [[x y & rest] [1 2 3 4 5]]
			(println x y (vec rest)))
	`)

	assertSourceMatchesLGB(t, "macros", `
		(defmacro unless [test & body]
			` + "`" + `(if (not ~test) (do ~@body)))
		(unless false (println "macro works"))
		(unless true (println "should not print"))
	`)

	assertSourceMatchesLGB(t, "try-catch", `
		(println
			(try
				(throw (ex-info "boom" {:code 42}))
				(catch e
					(str "caught: " (ex-message e)))))
	`)

	assertSourceMatchesLGB(t, "protocols", `
		(defprotocol Greetable
			(greet [this]))
		(defrecord Person [name])
		(extend-type Person
			Greetable
			(greet [this] (str "Hi, " (:name this))))
		(println (greet (->Person "Alice")))
	`)

	assertSourceMatchesLGB(t, "transducers", `
		(println (into [] (comp (map inc) (filter even?) (take 3)) (range 20)))
		(println (transduce (map inc) + 0 [1 2 3 4 5]))
	`)

	assertSourceMatchesLGB(t, "atoms-and-volatile", `
		(let [a (atom 0)]
			(dotimes [_ 5] (swap! a inc))
			(println @a))
		(let [v (volatile! 0)]
			(dotimes [_ 3] (vswap! v inc))
			(println @v))
	`)

	assertSourceMatchesLGB(t, "persistent-collections", `
		(println (assoc {:a 1} :b 2))
		(println (conj #{1 2} 3))
		(println (conj [1 2] 3))
		(println (into {} [[:a 1] [:b 2]]))
	`)

	assertSourceMatchesLGB(t, "string-operations", `
		(println (str "hello" " " "world"))
		(println (subs "abcdef" 2 4))
		(println (upper-case "hello"))
		(println (count "four"))
	`)

	assertSourceMatchesLGB(t, "ns-require-string", `
		(ns user (:require [string :as s]))
		(println (s/join ", " ["a" "b" "c"]))
		(println (s/upper-case "hello"))
	`)

	assertSourceMatchesLGB(t, "ns-require-set", `
		(ns user (:require [set]))
		(println (set/union #{1 2} #{2 3}))
		(println (set/intersection #{1 2 3} #{2 3 4}))
	`)

	assertSourceMatchesLGB(t, "multimethods", `
		(defmulti area :shape)
		(defmethod area :circle [m] (* 3 (:r m) (:r m)))
		(defmethod area :rect [m] (* (:w m) (:h m)))
		(println (area {:shape :circle :r 5}))
		(println (area {:shape :rect :w 3 :h 4}))
	`)

	assertSourceMatchesLGB(t, "dynamic-vars", `
		(def ^:dynamic *x* 10)
		(println *x*)
		(binding [*x* 42]
			(println *x*))
		(println *x*)
	`)

	assertSourceMatchesLGB(t, "delay-force", `
		(def d (delay (do (println "computing") 42)))
		(println (force d))
		(println (force d))
	`)
}
