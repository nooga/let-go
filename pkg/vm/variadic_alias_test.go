/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// TestVariadicInvokeDoesNotAliasCallerArgs pins the fix for a slice-aliasing
// bug in variadic invocation.
//
// invokeIn used to pack the rest-list with
//
//	args = append(args[0:arity-1], restlist)
//
// The reslice keeps the caller's backing array, so append WROTE the packed
// rest-list over the caller's element arity-1 — for a plain (fn [& as]),
// args[0]. A Go caller reusing one []Value across invocations got the right
// answer on the first call and garbage on the second, since its arguments had
// been replaced by the previous call's rest-list. Lisp call sites never saw it
// because they build a fresh argument vector per call.
func TestVariadicInvokeDoesNotAliasCallerArgs(t *testing.T) {
	// (fn [& as] as) — arity 1, variadic: the shape that clobbers args[0].
	chunk := NewCodeChunk(NewConsts())
	chunk.Append(OP_LOAD_ARG, 0)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(2)
	fn := MakeFunc(1, true, chunk)

	args := []Value{Int(7), Int(11)}
	first, err := fn.Invoke(args)
	if err != nil {
		t.Fatalf("first invoke: %v", err)
	}
	if args[0] != Value(Int(7)) || args[1] != Value(Int(11)) {
		t.Fatalf("caller args mutated: got [%v %v], want [7 11]", args[0], args[1])
	}

	second, err := fn.Invoke(args)
	if err != nil {
		t.Fatalf("second invoke: %v", err)
	}
	if first.String() != second.String() {
		t.Errorf("repeated invoke differs: first=%v second=%v", first, second)
	}
}
