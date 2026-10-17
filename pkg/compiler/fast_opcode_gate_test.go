/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// compileCode compiles src in the core ns and returns the top-level chunk's
// code words.
func compileCode(t *testing.T, src string) []int32 {
	t.Helper()
	c := NewTransientCompiler(consts, rt.NS(rt.NameCoreNS))
	chunk, _, err := c.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("CompileMultiple(%q): %v", src, err)
	}
	return chunk.Code()
}

// hasOpcode reports whether op appears in an opcode position. Operands can
// collide with small opcode numbers, so walk the stream by opcode width.
func hasOpcode(code []int32, op int32) bool {
	for ip := 0; ip < len(code); {
		cur := code[ip] & 0xff
		if cur == op {
			return true
		}
		ip += 1 + opcodeOperandWords(cur)
	}
	return false
}

func opcodeOperandWords(op int32) int {
	switch op {
	case vm.OP_LOAD_CONST, vm.OP_LOAD_ARG, vm.OP_INVOKE, vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE, vm.OP_JUMP,
		vm.OP_POP_N, vm.OP_DUP_NTH, vm.OP_LOAD_VAR, vm.OP_SET_VAR, vm.OP_MAKE_CLOSURE,
		vm.OP_LOAD_CLOSEDOVER, vm.OP_PUSH_CLOSEDOVER, vm.OP_TAIL_CALL, vm.OP_FINALLY_END, vm.OP_MAKE_MULTI_ARITY, vm.OP_RECUR_FN:
		return 1
	case vm.OP_TRY_PUSH:
		return 2
	case vm.OP_RECUR:
		return 3
	}
	return 0
}

// The fast-opcode gate accepts a head only when it resolves to the core var
// itself: unqualified, qualified to clojure.core (what syntax-quote emits,
// #1045), or through an alias.
func TestFastOpcodeGateAcceptsHeadsThatResolveToTheCoreVar(t *testing.T) {
	for _, tc := range []struct {
		src string
		op  int32
	}{
		{`(let [a (double-array 1)] (aget a 0))`, vm.OP_AGET},
		{`(let [a (double-array 1)] (aset a 0 1.5))`, vm.OP_ASET},
		{`(let [a (double-array 1)] (clojure.core/aget a 0))`, vm.OP_AGET},
		{`(let [x 1] (clojure.core/+ x 2))`, vm.OP_ADD},
		{`(let [x 1] (clojure.core/< x 2))`, vm.OP_LT},
	} {
		if !hasOpcode(compileCode(t, tc.src), tc.op) {
			t.Errorf("%s: opcode %d not emitted", tc.src, tc.op)
		}
	}
}

// A local that shadows the core name, a nested aget, and a redefinition in
// the current ns all stay on the call path, and the program sees their
// semantics rather than the opcode's.
func TestFastOpcodeGateRejectsShadowsAndRedefinitions(t *testing.T) {
	shadow := `(let [aget (fn [a i] :mine) a (double-array 1)] (aget a 0))`
	if hasOpcode(compileCode(t, shadow), vm.OP_AGET) {
		t.Error("shadowing local: OP_AGET emitted")
	}
	if got := mustEval(t, shadow); got != vm.Keyword("mine") {
		t.Errorf("shadowing local: got %v, want :mine", got)
	}

	nested := `(let [a (object-array 1)] (aset a 0 (long-array 1)) (aget a 0 0))`
	if hasOpcode(compileCode(t, nested), vm.OP_AGET) {
		t.Error("3-arg aget: OP_AGET emitted")
	}

	redef := `(ns fastop.redef) (defn aget [a i] :redefined) (aget (double-array 1) 0)`
	if hasOpcode(compileCode(t, redef), vm.OP_AGET) {
		t.Error("ns-level redefinition: OP_AGET emitted")
	}
	if got := mustEval(t, redef); got != vm.Keyword("redefined") {
		t.Errorf("ns-level redefinition: got %v, want :redefined", got)
	}
}

// The opcode and the core fn agree on results and on error text.
func TestFastAgetAsetMatchTheCoreFns(t *testing.T) {
	for _, tc := range []struct{ fast, slow string }{
		{`(let [a (double-array 3)] (aset a 1 2) (aget a 1))`,
			`(let [a (double-array 3) g aget s aset] (s a 1 2) (g a 1))`},
		{`(let [a (long-array 2)] (aset a 0 300) [(aget a 0) (aget a 1)])`,
			`(let [a (long-array 2) g aget s aset] (s a 0 300) [(g a 0) (g a 1)])`},
		{`(let [a (double-array 1)] (try (aget a 1) (catch Exception e (ex-message e))))`,
			`(let [a (double-array 1) g aget] (try (g a 1) (catch Exception e (ex-message e))))`},
		{`(try (aget "s" 0) (catch Exception e (ex-message e)))`,
			`(let [g aget] (try (g "s" 0) (catch Exception e (ex-message e))))`},
		{`(let [a (long-array 1)] (try (aset a 0 "x") (catch Exception e (ex-message e))))`,
			`(let [a (long-array 1) s aset] (try (s a 0 "x") (catch Exception e (ex-message e))))`},
	} {
		fast, slow := mustEval(t, tc.fast), mustEval(t, tc.slow)
		if !vm.ValueEquals(fast, slow) {
			t.Errorf("%s\n  opcode: %v\n  core fn: %v", tc.fast, fast, slow)
		}
	}
}
