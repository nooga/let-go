/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// mustEval evaluates src for its value, failing the test if it doesn't compile.
func mustEval(t *testing.T, src string) vm.Value {
	t.Helper()
	v, err := Eval(src)
	if err != nil {
		t.Fatalf("Eval(%q): %v", src, err)
	}
	return v
}

// compileFormStackEffect compiles v as a form on a fresh compiler and reports
// how many values the emitted code leaves on the stack.
func compileFormStackEffect(v vm.Value) (int, error) {
	c := NewTransientCompiler(consts, rt.NS(rt.NameCoreNS))
	c.chunk = vm.NewCodeChunk(c.consts)
	before := c.sp
	if err := c.compileForm(v); err != nil {
		return 0, err
	}
	return c.sp - before, nil
}

// compileForm's contract is that a form either fails to compile or leaves
// exactly one value on the stack. The type switch had no default, so a value
// whose ValueType was not enumerated emitted no bytecode at all and returned
// no error: the compile succeeded and the call blew up much later, indexing a
// stack one slot short (issue #943). Assert the contract over every value type
// that can reach a form, so the next unhandled type is a test failure here
// rather than a runtime panic in someone's macro.
func TestCompileFormAlwaysPushesOneValueOrErrors(t *testing.T) {
	// One representative per distinct ValueType a macroexpansion can splice in.
	// Kept as source rather than Go constructors so the list reads as the
	// language sees it.
	sources := []string{
		// scalars and tagged literals
		"1", "1.5", `"s"`, "nil", "true", ":k", `\c`, "'conj", "1N", "1/2", "1.0M",
		`#uuid "00000000-0000-0000-0000-000000000000"`, `#inst "2026-01-01T00:00:00Z"`, `#"re"`,
		"(fn [] 1)", "conj", "(type 1)",
		// vectors, both representations
		"[]", "[1 2 3]", "(vec (range 40))", "(reduce conj [] (range 33))",
		// maps, all three representations
		"{}", "{:a 1}", "(into {} (map (fn [i] [i i]) (range 20)))", "(sorted-map :b 2 :a 1)",
		// sets, all three representations
		"#{}", "#{1 2}", "(into #{} (range 20))", "(sorted-set 3 1 2)",
		// lists and the other seq representations
		"()", "'(1 2)", "(range 3)", "(repeat 2 1)", "(map inc [1 2])", "(seq [1 2])",
		"(seq (reduce conj [] (range 33)))", "(cons 1 '(2))",
		// opaque runtime values: no form can carry these, so they must error
		"(atom 1)", "(volatile! 1)", "(delay 1)", "(promise)", "(transient [])",
		"(transient {})", "(transient #{})", "(conj clojure.lang.PersistentQueue/EMPTY 1)",
		"(ex-info \"x\" {})",
	}

	for _, src := range sources {
		t.Run(src, func(t *testing.T) {
			v := mustEval(t, src)
			delta, err := compileFormStackEffect(v)
			if err != nil {
				// A refusal is a fine outcome — it just has to name the type.
				if !strings.Contains(err.Error(), v.Type().Name()) {
					t.Fatalf("compileForm(%s) refused without naming the type: %v", src, err)
				}
				return
			}
			if delta != 1 {
				t.Fatalf("compileForm(%s) (%s) left %d values on the stack, want 1",
					src, v.Type().Name(), delta)
			}
		})
	}
}

// unknownType stands in for a ValueType added to pkg/vm later. The point of
// the default arm is that such a type is refused by name instead of compiling
// to nothing, so the guard keeps working for types that don't exist yet.
type unknownType struct{}

func (t *unknownType) String() string     { return t.Name() }
func (t *unknownType) Type() vm.ValueType { return vm.TypeType }
func (t *unknownType) Unbox() any         { return reflect.TypeFor[*unknownType]() }
func (t *unknownType) Name() string       { return "let-go.test.Unknown" }
func (t *unknownType) Box(any) (vm.Value, error) {
	return vm.NIL, fmt.Errorf("can't box %s", t.Name())
}

type unknownValue struct{}

func (u unknownValue) String() string     { return "#<unknown>" }
func (u unknownValue) Type() vm.ValueType { return &unknownType{} }
func (u unknownValue) Unbox() any         { return nil }

func TestCompileFormRefusesUnhandledValueType(t *testing.T) {
	_, err := compileFormStackEffect(unknownValue{})
	if err == nil {
		t.Fatal("compileForm accepted a value of an unhandled type; it must refuse by name")
	}
	if !strings.Contains(err.Error(), "let-go.test.Unknown") {
		t.Fatalf("refusal does not name the type: %v", err)
	}
}

// The bug as reported: a macro folding a conj-built vector into its expansion
// worked up to 32 elements and crashed at 33, because conj promotes
// ArrayVector to PersistentVector there and only ArrayVector had a case.
func TestMacroFoldedVectorCompilesAtAnySize(t *testing.T) {
	for _, n := range []int{0, 1, 31, 32, 33, 64, 200} {
		src := fmt.Sprintf(`(do (defmacro folded-%d [] (reduce conj [] (range %d)))
                                (defn use-folded-%d [] (folded-%d))
                                (use-folded-%d))`, n, n, n, n, n)
		v, err := Eval(src)
		if err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}
		got, ok := v.(vm.Collection)
		if !ok {
			t.Fatalf("n=%d: got %s, want a collection", n, v.Type().Name())
		}
		if got.RawCount() != n {
			t.Fatalf("n=%d: folded vector has %d elements", n, got.RawCount())
		}
	}
}

// A PersistentVector form is a vector form, not a constant: its elements are
// subforms and have to be evaluated, exactly as in the ArrayVector case.
func TestPersistentVectorFormEvaluatesItsElements(t *testing.T) {
	// 33 elements so conj promotes past the ArrayVector threshold, with the
	// last two being forms rather than literals.
	v, err := Eval(`(do (def big (conj (reduce conj [] (range 31)) '(+ 1 2) 'marker))
                        (def marker 99)
                        ((eval (list 'fn [] big))))`)
	if err != nil {
		t.Fatalf("Eval: %v", err)
	}
	coll, ok := v.(vm.Collection)
	if !ok || coll.RawCount() != 33 {
		t.Fatalf("got %v (%s), want a 33-element vector", v, v.Type().Name())
	}
	nth, ok := v.(interface{ Nth(int) vm.Value })
	if !ok {
		t.Fatalf("result %s is not indexable", v.Type().Name())
	}
	if got := nth.Nth(31); got != vm.Int(3) {
		t.Fatalf("element 31 = %v, want 3 — the (+ 1 2) subform was not evaluated", got)
	}
	if got := nth.Nth(32); got != vm.Int(99) {
		t.Fatalf("element 32 = %v, want 99 — the symbol subform was not resolved", got)
	}
}

// Sorted collections keep both their type and their ordering across a compile:
// rebuilding them as array-map/hash-set would round-trip the contents while
// silently dropping the sort, which is the same class of quiet wrong answer.
func TestSortedCollectionFormsKeepTypeAndOrder(t *testing.T) {
	for _, tc := range []struct{ src, wantType, want string }{
		{`(sorted-map :c 3 :a 1 :b 2)`, "let-go.lang.PersistentTreeMap", "{:a 1, :b 2, :c 3}"},
		{`(sorted-set 3 1 2)`, "let-go.lang.PersistentTreeSet", "#{1 2 3}"},
	} {
		v, err := Eval(fmt.Sprintf("((eval (list 'fn [] %s)))", tc.src))
		if err != nil {
			t.Fatalf("Eval(%s): %v", tc.src, err)
		}
		if v.Type().Name() != tc.wantType {
			t.Fatalf("%s compiled to %s, want %s", tc.src, v.Type().Name(), tc.wantType)
		}
		if v.String() != tc.want {
			t.Fatalf("%s compiled to %s, want %s", tc.src, v, tc.want)
		}
	}
}

// A custom comparator is not expressible as a form, so the compiler has to say
// so rather than emit a default-ordered rebuild that quietly reorders the keys.
func TestSortedCollectionWithCustomComparatorIsRefused(t *testing.T) {
	for _, src := range []string{`(sorted-map-by > :a 1 :b 2)`, `(sorted-set-by > 1 2)`} {
		v := mustEval(t, src)
		_, err := compileFormStackEffect(v)
		if err == nil {
			t.Fatalf("compileForm(%s) accepted a custom comparator it cannot carry", src)
		}
		if !strings.Contains(err.Error(), "comparator") {
			t.Fatalf("refusal for %s does not mention the comparator: %v", src, err)
		}
	}
}

// Every seq representation is a call form, the way a list is. LazySeq already
// behaved this way; Range, Repeat and the vector seqs report distinct
// ValueTypes and used to fall through the switch instead.
func TestSeqFormsCompileAsCalls(t *testing.T) {
	// A lazy seq already reported ListType, so it was never affected; the
	// distinct-ValueType seqs below are the ones that used to fall through.
	for _, tc := range []struct{ src, want string }{
		{`(map inc [1 2])`, "let-go.lang.PersistentList"},
		{`(range 3)`, "let-go.lang.Range"},
		{`(repeat 2 1)`, "let-go.lang.Repeat"},
		{`(seq (reduce conj [] (range 33)))`, "let-go.lang.Sequence"},
		{`(seq [1 2])`, "let-go.lang.PersistentList"},
	} {
		v := mustEval(t, tc.src)
		if got := v.Type().Name(); got != tc.want {
			t.Fatalf("%s is a %s, not the %s this case is meant to cover", tc.src, got, tc.want)
		}
		delta, err := compileFormStackEffect(v)
		if err != nil {
			t.Fatalf("compileForm(%s): %v", tc.src, err)
		}
		if delta != 1 {
			t.Fatalf("compileForm(%s) left %d values on the stack, want 1", tc.src, delta)
		}
	}
	// And the call is a real one: the head of the seq is invoked.
	v, err := Eval(`((eval (list 'fn [] (map (fn [i] (if (zero? i) 'inc 41)) (range 2))))) `)
	if err != nil {
		t.Fatalf("Eval: %v", err)
	}
	if v != vm.Int(42) {
		t.Fatalf("seq form evaluated to %v, want 42", v)
	}
}
