/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// A collection literal whose elements are all constants is one value, shared
// by every evaluation, as Clojure compiles it. A literal with a non-constant
// element or with metadata is built fresh each time, as in Clojure.
func TestConstantCollectionLiteralsAreShared(t *testing.T) {
	cases := []identityCase{
		{"vector", "(do (defn mk [] [1 2 3]) (identical? (mk) (mk)))", true},
		{"map", "(do (defn mk [] {:a 1 :b 2}) (identical? (mk) (mk)))", true},
		{"set", "(do (defn mk [] #{:a :b}) (identical? (mk) (mk)))", true},
		{"nested", "(do (defn mk [] [1 {:k [2 #{3}]}]) (identical? (mk) (mk)))", true},
		{"empty vector", "(do (defn mk [] []) (identical? (mk) (mk)))", true},
		{"non-constant element", "(do (defn mk [x] [x 1]) (identical? (mk 1) (mk 1)))", false},
		{"nested non-constant element", "(do (defn mk [x] [1 {:k x}]) (identical? (mk 1) (mk 1)))", false},
		{"metadata", "(do (defn mk [] ^{:t 1} [1 2]) (identical? (mk) (mk)))", false},
		{"char", `(do (defn mk [] [\a 1]) (identical? (mk) (mk)))`, true},
		{"bignum, ratio, decimal", "(do (defn mk [] [1N 1/2 1.5M]) (identical? (mk) (mk)))", true},
		{"uuid", `(do (defn mk [] [#uuid "00000000-0000-0000-0000-000000000001"]) (identical? (mk) (mk)))`, true},
		{"regex", `(do (defn mk [] [#"x"]) (identical? (mk) (mk)))`, true},
		{"map over eight entries", "(do (defn mk [] {:a 1 :b 2 :c 3 :d 4 :e 5 :f 6 :g 7 :h 8 :i 9 :j 10}) (identical? (mk) (mk)))", true},
		{"vector over 32 elements", "(do (defn mk [] [" + strings.Repeat("1 ", 40) + "]) (identical? (mk) (mk)))", true},
	}
	checkIdentity(t, cases)
}

// Sharing a literal is unobservable except through identical?: every
// operation on the shared value copies, the map literal keeps its insertion
// order, and a literal's metadata still arrives.
func TestSharedLiteralsStayIntact(t *testing.T) {
	src := `(do
	  (defn v [] [1 2 3])
	  (defn m [] {:b 1 :a 2 :c 3})
	  (defn s [] #{:a})
	  (defn w [] ^{:t 1} [1 2])
	  (conj (v) 4) (assoc (v) 0 9) (with-meta (v) {:x 1}) (persistent! (conj! (transient (v)) 5))
	  (assoc (m) :z 9) (dissoc (m) :a) (persistent! (assoc! (transient (m)) :q 1))
	  (conj (s) :b) (disj (s) :a)
	  [(v) (seq (m)) (s) (meta (w)) (meta (v)) (= (v) [1 2 3])])`
	out, err := Eval(src)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), "[[1 2 3] ([:b 1] [:a 2] [:c 3]) #{:a} {:t 1} nil true]"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

// identityCase evaluates src, an (identical? …) of two evaluations of a form,
// and expects shared.
type identityCase struct {
	name   string
	src    string
	shared bool
}

func checkIdentity(t *testing.T, cases []identityCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Eval(tc.src)
			if err != nil {
				t.Fatal(err)
			}
			if got := out == vm.TRUE; got != tc.shared {
				t.Fatalf("identical? = %v, want %v", got, tc.shared)
			}
		})
	}
}

// Each collection literal is its own constant: a map literal keeps its own key
// order even beside an equal one written in another order, and two literal
// sites are two values, as in Clojure; one site stays one value.
func TestMapLiteralsKeepTheirOwnKeyOrder(t *testing.T) {
	out, err := Eval(`(do
	  (defn ab [] {:a 1 :b 2})
	  (defn ba [] {:b 2 :a 1})
	  [(vec (keys (ab))) (vec (keys (ba))) (identical? [1 2] [1 2]) (identical? (ab) (ab))])`)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), "[[:a :b] [:b :a] false true]"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
