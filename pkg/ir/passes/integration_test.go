/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package passes

import (
	"errors"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/ir"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestIntegration_MandelbrotHelpers runs Build → ConstFold → CSE → DCE
// → Lower on the complex-number helpers from examples/mandelbrot.lg.
// If a helper fails to Build because of an unsupported opcode, that's a
// SKIP with the opcode logged (not a failure — Build coverage is owed
// separate work). If Build succeeds but Lower fails, that IS a failure
// (Lower is what Tasks 1-7 just made robust).
//
// The fixture uses flat args (cadd a b c d) rather than destructured
// pairs ([[a b] [c d]]) so the test exercises arithmetic and CSE without
// also requiring destructuring support in Build.
func TestIntegration_MandelbrotHelpers(t *testing.T) {
	src := `
(def scale 1000)
(defn cadd [a b c d] (vector (+ a c) (+ b d)))
(defn cmul [a b c d]
  (vector (/ (- (* a c) (* b d)) scale)
          (/ (+ (* a d) (* b c)) scale)))
(defn cmagsq [a b]
  (/ (+ (* a a) (* b b)) scale))
`
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	if _, _, err := ctx.CompileMultiple(strings.NewReader(src)); err != nil {
		t.Fatalf("compile: %v", err)
	}

	cases := []struct {
		name  string
		arity int
	}{
		{"cadd", 4},
		{"cmul", 4},
		{"cmagsq", 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := ns.Lookup(vm.Symbol(c.name))
			if v == nil {
				t.Fatalf("symbol %s not found after compile", c.name)
			}
			fnVal := v.(*vm.Var).Deref()
			fn, ok := fnVal.(*vm.Func)
			if !ok {
				t.Fatalf("%s is not a Func: %T", c.name, fnVal)
			}

			irFn, err := ir.Build(fn.Chunk(), c.name, c.arity, false)
			if err != nil {
				if errors.Is(err, ir.ErrUnsupportedOp) {
					t.Skipf("Build unsupported for %s: %v", c.name, err)
					return
				}
				t.Fatalf("Build failed for %s: %v", c.name, err)
			}

			// Run the optimization passes.
			ConstFold(irFn)
			CSE(irFn)
			DCE(irFn)

			if _, err := ir.Lower(irFn); err != nil {
				t.Errorf("Lower failed for %s after passes: %v\nIR:\n%s", c.name, err, ir.Dump(irFn))
				return
			}
			t.Logf("%s: Build + ConstFold + CSE + DCE + Lower round-tripped", c.name)
		})
	}
}
