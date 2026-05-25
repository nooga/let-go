/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

// Package rt — disasm.go installs `let-go.lang/disassemble`, a Clojure-
// callable bytecode introspection primitive used by the fusion-analysis
// tooling (and useful for debugging/learning the VM).
//
// `(disassemble f)` returns a vector of [opcode-keyword & args] tuples,
// one per instruction. Arg words (LDC index, branch offset, etc.) appear
// as integers after the opcode keyword. Walks ONLY the direct chunk —
// nested functions are reached by recursing on `(constants f)` and
// disassembling each `*Func` you find there.
//
// Example:
//   (def f (fn [x] (+ x 1)))
//   (disassemble f)
//   ;; => [[:LOAD_ARG 0] [:LOAD_CONST 0] [:ADD] [:RETURN]]

package rt

import (
	"fmt"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// opcodeStride returns the number of int32 words that an opcode + its
// inline args occupy in a code chunk. Must stay in sync with the
// emission logic in pkg/compiler/compiler.go and the dispatch in
// pkg/vm/vm.go's Frame.Run.
func opcodeStride(op int32) int {
	switch op & 0xff {
	case vm.OP_TRY_PUSH:
		return 3 // catchOffset, finallyOffset
	case vm.OP_RECUR:
		return 4 // offset, argc, ignore
	case vm.OP_LOAD_ARG, vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE, vm.OP_JUMP,
		vm.OP_POP_N, vm.OP_DUP_NTH, vm.OP_INVOKE, vm.OP_LOAD_CLOSEDOVER,
		vm.OP_RECUR_FN, vm.OP_MAKE_MULTI_ARITY, vm.OP_TAIL_CALL,
		vm.OP_LOAD_CONST, vm.OP_LOAD_VAR:
		return 2 // one int32 arg
	default:
		return 1
	}
}

// opcodeMnemonic returns the bare mnemonic ("ADD" etc.) without the
// sp prefix that vm.OpcodeToString includes.
func opcodeMnemonic(op int32) string {
	s := vm.OpcodeToString(op)
	if _, after, ok := strings.Cut(s, "/"); ok {
		return strings.TrimSpace(after)
	}
	return s
}

// chunkOf extracts the *CodeChunk from any value that carries one.
// Handles *Func directly and *Var by dereferencing.
func chunkOf(v vm.Value) (*vm.CodeChunk, error) {
	switch x := v.(type) {
	case *vm.Func:
		return x.Chunk(), nil
	case *vm.Var:
		return chunkOf(x.Deref())
	case *vm.Boxed:
		if c, ok := x.Unbox().(*vm.CodeChunk); ok {
			return c, nil
		}
	}
	return nil, fmt.Errorf("disassemble: expected Func, Var-of-Func, or boxed *CodeChunk, got %s",
		v.Type().Name())
}

// disassembleChunk walks chunk's code array and returns a slice of
// instruction vectors. Each inner slice starts with the opcode mnemonic
// (as a Keyword) followed by any inline arg integers.
func disassembleChunk(chunk *vm.CodeChunk) vm.Value {
	code := chunk.Code()
	var rows []vm.Value
	i := 0
	for i < len(code) {
		op := code[i]
		stride := opcodeStride(op)
		row := make([]vm.Value, 0, stride)
		row = append(row, vm.Keyword(opcodeMnemonic(op)))
		for j := 1; j < stride && i+j < len(code); j++ {
			row = append(row, vm.Int(int(code[i+j])))
		}
		rows = append(rows, vm.NewArrayVector(row))
		i += stride
	}
	return vm.NewArrayVector(rows)
}

// installDisasmNS installs the `disasm` Clojure namespace with a single
// fn: disassemble.
func installDisasmNS() {
	disasm, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("disassemble: expected 1 arg, got %d", len(vs))
		}
		chunk, err := chunkOf(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return disassembleChunk(chunk), nil
	})

	// constants — peer at the const pool to find nested functions.
	// Returns a vector of all constant values in order.
	constants, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("constants: expected 1 arg, got %d", len(vs))
		}
		chunk, err := chunkOf(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		consts := chunk.Consts()
		// Values() returns only this layer's consts; AllValues() includes
		// inherited parent consts. For introspection of nested funcs we
		// want all reachable consts.
		all := consts.AllValues()
		out := make([]vm.Value, len(all))
		copy(out, all)
		return vm.NewArrayVector(out), nil
	})

	ns := vm.NewNamespace("disasm")
	ns.Def("disassemble", disasm)
	ns.Def("constants", constants)
	RegisterNS(ns)
}
