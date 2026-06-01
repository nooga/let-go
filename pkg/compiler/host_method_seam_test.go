/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func runHM(t *testing.T, src string) vm.Value {
	t.Helper()
	cp := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := NewCompiler(cp, ns)
	ch, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v, err := vm.NewFrame(ch, nil).Run()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	return v
}

// The host-method seam lets a (type, method) pair be registered so that a
// `.method` dot-form on a value of that type dispatches to the registered fn
// (receiver first). This backs Java-method calls on let-go values without a
// general reflective bridge.
func TestHostMethodSeam(t *testing.T) {
	// A method registered on a deftype, dispatched via a dot-form.
	if v := runHM(t, `(do
                        (deftype Pair [a b])
                        (register-host-method! Pair 'sum (fn [p] (+ (.a p) (.b p))))
                        (.sum (->Pair 3 4)))`); v.String() != "7" {
		t.Fatalf("deftype host method: want 7, got %v", v)
	}

	// A registered method taking arguments.
	if v := runHM(t, `(do
                        (deftype Acc [n])
                        (register-host-method! Acc 'plus (fn [a x] (+ (.n a) x)))
                        (.plus (->Acc 5) 3))`); v.String() != "8" {
		t.Fatalf("host method with arg: want 8, got %v", v)
	}

	// A method registered on a built-in collection type.
	if v := runHM(t, `(do
                        (register-host-method! (type []) 'firstish (fn [v] (first v)))
                        (.firstish [10 20]))`); v.String() != "10" {
		t.Fatalf("native-type host method: want 10, got %v", v)
	}

	// A registered method whose name collides with a core fn (pop) must win
	// over the generic name→fn fallback — otherwise `.pop` on a custom type
	// would be hijacked by core `pop`. (malli needs .nth/.cons/.pop etc.)
	if v := runHM(t, `(do
                        (deftype Stack [top])
                        (register-host-method! Stack 'pop (fn [s] (.top s)))
                        (.pop (->Stack 99)))`); v.String() != "99" {
		t.Fatalf("registered .pop must win over core pop: want 99, got %v", v)
	}
}
