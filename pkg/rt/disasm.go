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
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
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
		vm.OP_LOAD_CONST, vm.OP_LOAD_VAR, vm.OP_FINALLY_END:
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

// projectConst maps a const-pool value to EDN-safe, DETERMINISTIC data that
// preserves identifiers. Readable scalars pass through unchanged; a Var becomes
// its qualified symbol (the identifier it carries); a Func becomes [:fn name];
// anything else (atoms, boxed Go objects, data structures) becomes a type tag.
// This is what keeps the dump free of non-deterministic 0x heap pointers (which
// is all you get from String() on live boxed objects) while keeping the names
// that matter. (Vector/map consts are tagged, not recursed — a follow-up.)
func projectConst(v vm.Value) vm.Value {
	switch x := v.(type) {
	case vm.Symbol, vm.Keyword, vm.String, vm.Int, vm.Float, vm.Boolean,
		vm.Char, *vm.Nil, *vm.BigInt, *vm.Ratio, *vm.BigDecimal:
		return v
	case *vm.Var:
		// "#'ns/name" -> symbol ns/name (the identifier is right there)
		return vm.Symbol(strings.TrimPrefix(x.String(), "#'"))
	case *vm.Func:
		if name := x.FuncName(); name != "" {
			return vm.NewArrayVector([]vm.Value{vm.Keyword("fn"), vm.Symbol(name)})
		}
		return vm.NewArrayVector([]vm.Value{vm.Keyword("fn")})
	default:
		return vm.NewArrayVector([]vm.Value{vm.Keyword("opaque"), vm.String(v.Type().Name())})
	}
}

// disassembleChunkResolved is disassembleChunk plus inline operand resolution:
// LOAD_CONST / LOAD_VAR rows get the referenced const's identifier appended
// (via projectConst), so e.g. [:LOAD_VAR 51 clojure.core/seq] instead of a bare
// index. Index is the global const index; AllValues() is the layered pool.
func disassembleChunkResolved(chunk *vm.CodeChunk) vm.Value {
	code := chunk.Code()
	consts := chunk.Consts().AllValues()
	var rows []vm.Value
	i := 0
	for i < len(code) {
		op := code[i]
		stride := opcodeStride(op)
		row := make([]vm.Value, 0, stride+1)
		row = append(row, vm.Keyword(opcodeMnemonic(op)))
		for j := 1; j < stride && i+j < len(code); j++ {
			row = append(row, vm.Int(int(code[i+j])))
		}
		switch op & 0xff {
		case vm.OP_LOAD_CONST, vm.OP_LOAD_VAR:
			if i+1 < len(code) {
				if idx := int(code[i+1]); idx >= 0 && idx < len(consts) {
					row = append(row, projectConst(consts[idx]))
				}
			}
		}
		rows = append(rows, vm.NewArrayVector(row))
		i += stride
	}
	return vm.NewArrayVector(rows)
}

// installDisasmNS installs the `disasm` Clojure namespace with a single
// fn: disassemble.
func init() { RegisterInstaller(installDisasmNS) }

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

	// decode-bundle — open a .lgb bundle FILE and return its chunks as data,
	// so a Lisp layer can lazily disassemble each via `disassemble`/`constants`.
	// Reuses the SAME module-import path the runtime uses to load core
	// (bytecode.DecodeToExecUnit + the DefNSBare/DefStub resolver from
	// loadPrecompiledBundle), so chunks come back fully wired (nested funcs
	// resolved) without executing anything. Returns:
	//   {:main <boxed-CodeChunk>
	//    :namespaces [{:ns "name" :chunk-index i :chunk <boxed-CodeChunk>} ...]}
	// Each boxed chunk is accepted directly by `disassemble`/`constants`
	// (chunkOf already unwraps *Boxed of *CodeChunk).
	decodeBundle, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("decode-bundle: expected 1 arg (path), got %d", len(vs))
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("decode-bundle: expected String path, got %s", vs[0].Type().Name())
		}
		data, err := os.ReadFile(string(path))
		if err != nil {
			return vm.NIL, fmt.Errorf("decode-bundle: %w", err)
		}
		// Same resolver loadPrecompiledBundle uses: create bare namespaces and
		// stub vars so decoding never triggers the loader or warns on shadow.
		resolve := func(nsName, name string) *vm.Var {
			n := DefNSBare(nsName)
			if v := n.LookupLocal(vm.Symbol(name)); v != nil {
				return v
			}
			return n.DefStub(name)
		}
		unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(data), resolve)
		if err != nil {
			return vm.NIL, fmt.Errorf("decode-bundle: %w", err)
		}
		nsEntries := make([]vm.Value, 0, len(unit.NSOrder))
		for idx, name := range unit.NSOrder {
			ch := unit.NSChunks[name]
			if ch == nil {
				continue
			}
			nsEntries = append(nsEntries, vm.NewArrayMap([]vm.Value{
				vm.Keyword("ns"), vm.String(name),
				vm.Keyword("chunk-index"), vm.Int(idx),
				vm.Keyword("chunk"), vm.NewBoxed(ch),
			}))
		}
		// Bundle-level const pool, projected ONCE (all chunks share it, so
		// per-chunk consts would just repeat it 30x). AllValues = layered pool.
		var poolVals []vm.Value
		if unit.MainChunk != nil {
			poolVals = unit.MainChunk.Consts().AllValues()
		}
		consts := make([]vm.Value, len(poolVals))
		for i, c := range poolVals {
			consts[i] = projectConst(c)
		}
		out := vm.NewArrayMap([]vm.Value{
			vm.Keyword("main"), vm.NewBoxed(unit.MainChunk),
			vm.Keyword("namespaces"), vm.NewArrayVector(nsEntries),
			vm.Keyword("consts"), vm.NewArrayVector(consts),
		})
		return out, nil
	})

	// disassemble-resolved — like disassemble, but each LOAD_CONST/LOAD_VAR row
	// carries the referenced identifier inline (via projectConst), so the
	// disassembly reads without a separate const-pool lookup.
	disasmResolved, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("disassemble-resolved: expected 1 arg, got %d", len(vs))
		}
		chunk, err := chunkOf(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return disassembleChunkResolved(chunk), nil
	})

	ns := vm.NewNamespace("disasm")
	ns.Def("disassemble", disasm)
	ns.Def("disassemble-resolved", disasmResolved)
	ns.Def("constants", constants)
	ns.Def("decode-bundle", decodeBundle)
	RegisterNS(ns)
}
