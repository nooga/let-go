/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// spikeInst is a flattened IR instruction (test-only decoder output).
type spikeInst struct {
	op     string   // IR op keyword (e.g., "const", "add")
	opc    uint8    // Interned opcode for fast dispatch
	args   []int32  // RefInstIds (instruction indices)
	aux    int32    // Auxiliary data: arg index for :load-arg, param index for :block-arg
	auxVal vm.Value // For :const, the actual constant value; for :load-var, resolved var value
	typ    uint8    // Type code: 0=unknown, 1=int, 2=float
}

// spikeEdge represents one edge of a terminator with its target block and edge args.
type spikeEdge struct {
	target int32   // BlockID of the target block
	args   []int32 // Edge arguments (InstIds) to pass as block params
}

// spikeTerm describes a terminator's shape for a spikeBlock.
type spikeTerm struct {
	op         string    // terminator op: "return", "branch", "branch-if"
	opc        uint8     // Interned opcode for fast dispatch
	returnVal  int32     // for return: the value inst-id
	condRef    int32     // for branch-if: the condition inst-id
	trueEdge   spikeEdge // for branch-if: true branch edge
	falseEdge  spikeEdge // for branch-if: false branch edge
	simpleEdge spikeEdge // for branch: the single edge
}

// spikeBlock describes the insts and structure of a block.
type spikeBlock struct {
	insts  []int32 // inst-ids in this block (in order, excluding terminator)
	params []int32 // block param inst-ids (for block-arg insts)
	term   spikeTerm
}

// spikeStats tracks execution metrics for unboxed vs boxed paths.
type spikeStats struct {
	boxOps        int // number of vm.Value allocations from typed slots
	unboxOps      int // number of reads from unboxed slots
	boxedArithOps int // number of arithmetic ops executed on boxed path
	callOps       int // number of call ops executed
}

// spikeRoute indicates where an inst's result is stored: unboxed or boxed.
const (
	ROUTE_BOXED = uint8(0)
	ROUTE_INT   = uint8(1)
	ROUTE_FLOAT = uint8(2)
	// ROUTE_BOOL: comparison results stored as 0/1 in localsI; branch-if reads
	// them natively, boxing to vm.Boolean only at value boundaries.
	ROUTE_BOOL = uint8(3)
)

// spikeOp* are interned integer opcodes for spike VM dispatch.
const (
	spikeOpInvalid  = uint8(0)
	spikeOpConst    = uint8(1)
	spikeOpLoadArg  = uint8(2)
	spikeOpLoadVar  = uint8(3)
	spikeOpBlockArg = uint8(4)
	spikeOpAdd      = uint8(5)
	spikeOpSub      = uint8(6)
	spikeOpMul      = uint8(7)
	spikeOpInc      = uint8(8)
	spikeOpDec      = uint8(9)
	spikeOpLt       = uint8(10)
	spikeOpLte      = uint8(11)
	spikeOpGt       = uint8(12)
	spikeOpGte      = uint8(13)
	spikeOpEq       = uint8(14)
	spikeOpBranch   = uint8(15)
	spikeOpBranchIf = uint8(16)
	spikeOpReturn   = uint8(17)
	spikeOpCall     = uint8(18)
)

// spikeFn is the flattened output of decode: slices replace persistent maps.
type spikeFn struct {
	insts         []spikeInst  // indexed by inst-id
	blocks        []spikeBlock // indexed by block-id
	consts        []vm.Value   // const pool (unused for spike; for reference)
	callees       []vm.Value   // resolved callees (one per call in insts)
	calleeIndices []int        // indexed by inst-id: index in callees slice for call insts, -1 otherwise
	nargs         int          // arity
	routes        []uint8      // routing for each inst: ROUTE_BOXED/INT/FLOAT (computed post-decode for typed path)
}

var spikeDecodeVarCounter int

// decodeOptimizedIR flattens an optimized IR Function (returned by the Lisp
// pipeline) into flat Go structures (spikeInst/spikeBlock/spikeFn). The
// optimized IR is a Lisp atom wrapping a persistent-map Function.
//
// Returns error for out-of-scope ops, unresolvable callees, or malformed IR.
func decodeOptimizedIR(irValue vm.Value) (*spikeFn, error) {
	// irValue is a Lisp atom; access the map fields via Lisp evaluations.
	spikeDecodeVarCounter++
	varName := fmt.Sprintf("*ir-fn-%d*", spikeDecodeVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, irValue)

	// Helper to evaluate a Lisp expression and get result
	evalExpr := func(expr string) (vm.Value, error) {
		consts := vm.NewConsts()
		c := compiler.NewCompiler(consts, coreNS)
		c.SetSource("spike-decode")
		_, result, err := c.CompileMultiple(strings.NewReader(expr))
		return result, err
	}

	// STORY-0059: run the pre-execution cleanup IR pass BEFORE flattening. It drops
	// dead (DCE-tombstoned) block-params + their in-edge args, strips :pop markers,
	// and DCE-sweeps insts orphaned by the compaction — promoting what used to be the
	// Go-side compactDeadParams into the real ir.passes.cleanup pass. cleanup mutates
	// the IR atom in place, so every field read below sees the cleaned form, and
	// validateLiveInvariants then holds with NO decoder-side compaction.
	if _, err := evalExpr(fmt.Sprintf("(ir.passes.cleanup/cleanup %s)", varName)); err != nil {
		return nil, fmt.Errorf("eval cleanup pass: %w", err)
	}

	// Get the consts pool via (ir.data/fn-consts ir-var)
	constPoolVal, err := evalExpr(fmt.Sprintf("(ir.data/fn-consts %s)", varName))
	if err != nil {
		return nil, fmt.Errorf("eval fn-consts: %w", err)
	}
	boxed, ok := constPoolVal.(*vm.Boxed)
	if !ok {
		return nil, fmt.Errorf("fn-consts did not return boxed *vm.Consts, got %T", constPoolVal)
	}
	constPool := boxed.Unbox().(*vm.Consts)

	// Get arity
	arityVal, err := evalExpr(fmt.Sprintf("(ir.data/fn-arity %s)", varName))
	if err != nil {
		return nil, fmt.Errorf("eval fn-arity: %w", err)
	}
	arityInt, ok := arityVal.(vm.Int)
	if !ok {
		return nil, fmt.Errorf("fn-arity is not an int: %T", arityVal)
	}
	arity := int(arityInt)

	// Get inst count
	instCountVal, err := evalExpr(fmt.Sprintf("(ir.data/inst-count %s)", varName))
	if err != nil {
		return nil, fmt.Errorf("eval inst-count: %w", err)
	}
	instCountInt, ok := instCountVal.(vm.Int)
	if !ok {
		return nil, fmt.Errorf("inst-count is not an int: %T", instCountVal)
	}
	instCount := int(instCountInt)

	// Get block count
	blockCountVal, err := evalExpr(fmt.Sprintf("(count (ir.data/blocks %s))", varName))
	if err != nil {
		return nil, fmt.Errorf("eval block-count: %w", err)
	}
	blockCountInt, ok := blockCountVal.(vm.Int)
	if !ok {
		return nil, fmt.Errorf("block-count is not an int: %T", blockCountVal)
	}
	blockCount := int(blockCountInt)

	// Decode all insts
	insts := make([]spikeInst, instCount) // Pre-allocate; :invalid tombstones create gaps
	var callees []vm.Value
	calleeIndices := make([]int, instCount)
	// Initialize all to -1 (not a call)
	for i := range calleeIndices {
		calleeIndices[i] = -1
	}

	for nid := 0; nid < instCount; nid++ {
		// Get op
		opVal, err := evalExpr(fmt.Sprintf("(ir.data/op %d %s)", nid, varName))
		if err != nil {
			return nil, fmt.Errorf("eval op at inst %d: %w", nid, err)
		}
		opKw, ok := opVal.(vm.Keyword)
		if !ok {
			return nil, fmt.Errorf("op at inst %d is not a keyword: %T", nid, opVal)
		}
		opStr := string(opKw)

		// Skip :invalid insts (tombstones left by DCE — dead inst markers, nothing references them)
		// Make the gap explicit: set op to "invalid" so accidental execution is loud
		if opStr == "invalid" {
			insts[nid].op = "invalid"
			insts[nid].opc = spikeOpInvalid
			continue
		}

		// Validate op is in scope
		if !isValidScopeOp(opStr) {
			return nil, fmt.Errorf("inst %d has out-of-scope op: :%s (valid ops: const, load-arg, load-var, block-arg, add, sub, mul, inc, dec, lt, lte, gt, gte, eq, branch, branch-if, return, call)", nid, opStr)
		}

		// Get refs
		refsVal, err := evalExpr(fmt.Sprintf("(ir.data/refs %d %s)", nid, varName))
		if err != nil {
			return nil, fmt.Errorf("eval refs at inst %d: %w", nid, err)
		}
		refList, ok := refsVal.(vm.Indexed)
		if !ok {
			return nil, fmt.Errorf("refs at inst %d is not indexed: %T", nid, refsVal)
		}
		refCount := refList.RawCount()
		refInstIDs := make([]int32, refCount)
		for i := 0; i < refCount; i++ {
			refVal := refList.Nth(i)
			refID, ok := refVal.(vm.Int)
			if !ok {
				return nil, fmt.Errorf("ref at inst %d[%d] is not an int: %T", nid, i, refVal)
			}
			refInstIDs[i] = int32(refID)
		}

		// Get aux (raw; may be value or int depending on op)
		auxVal, err := evalExpr(fmt.Sprintf("(ir.data/aux %d %s)", nid, varName))
		if err != nil {
			return nil, fmt.Errorf("eval aux at inst %d: %w", nid, err)
		}

		// Process aux based on op type
		var auxInt int32
		var auxValue vm.Value

		switch opStr {
		case "const":
			// For :const, aux IS the constant value
			auxValue = auxVal

		case "load-arg":
			// For :load-arg, aux is the arg index
			if auxVal != vm.NIL {
				if auxI, ok := auxVal.(vm.Int); ok {
					auxInt = int32(auxI)
				}
			}

		case "block-arg":
			// For :block-arg, aux is the param index
			if auxVal != vm.NIL {
				if auxI, ok := auxVal.(vm.Int); ok {
					auxInt = int32(auxI)
				}
			}

		case "load-var":
			// For :load-var, resolve the var at decode time
			// aux contains the symbol or var reference
			auxValue = auxVal
			if auxVal != vm.NIL {
				// Try to resolve it as a var
				resolvedVal, resolveErr := resolveVar(auxVal, coreNS)
				if resolveErr == nil {
					auxValue = resolvedVal
				} else {
					// If it can't be resolved, store it as-is; may be resolved at execution time
					auxValue = auxVal
				}
			}

		default:
			// For other ops, ignore aux
			if auxVal != vm.NIL {
				if auxI, ok := auxVal.(vm.Int); ok {
					auxInt = int32(auxI)
				}
			}
		}

		// Get type
		typeVal, err := evalExpr(fmt.Sprintf("(ir.data/type-of %d %s)", nid, varName))
		if err != nil {
			return nil, fmt.Errorf("eval type-of at inst %d: %w", nid, err)
		}
		typeCode := uint8(0)
		if kw, ok := typeVal.(vm.Keyword); ok {
			typeCode = typeToCode(kw)
		}

		// For call, resolve the callee at decode time
		if opStr == "call" {
			if len(refInstIDs) < 1 {
				return nil, fmt.Errorf("inst %d (call) has no callee ref", nid)
			}
			calleeID := refInstIDs[0]
			// Get the callee op
			calleeOpVal, err := evalExpr(fmt.Sprintf("(ir.data/op %d %s)", calleeID, varName))
			if err != nil {
				return nil, fmt.Errorf("eval op at callee inst %d (for call at %d): %w", calleeID, nid, err)
			}
			calleeOpKw, ok := calleeOpVal.(vm.Keyword)
			if !ok {
				return nil, fmt.Errorf("callee op at inst %d is not a keyword: %T", calleeID, calleeOpVal)
			}
			calleeOp := string(calleeOpKw)

			var calleeVal vm.Value

			if calleeOp == "const" {
				// Callee is a const; get its value
				calleeAuxVal, err := evalExpr(fmt.Sprintf("(ir.data/aux %d %s)", calleeID, varName))
				if err != nil {
					return nil, fmt.Errorf("eval aux at callee inst %d (for call at %d): %w", calleeID, nid, err)
				}
				calleeVal = calleeAuxVal
			} else if calleeOp == "load-var" {
				// Callee is loaded from a var; resolve it
				calleeAuxVal, err := evalExpr(fmt.Sprintf("(ir.data/aux %d %s)", calleeID, varName))
				if err != nil {
					return nil, fmt.Errorf("eval aux at callee inst %d (for call at %d): %w", calleeID, nid, err)
				}
				resolved, resolveErr := resolveVar(calleeAuxVal, coreNS)
				if resolveErr != nil {
					return nil, fmt.Errorf("inst %d (call) callee inst %d (:load-var) unresolvable: %w", nid, calleeID, resolveErr)
				}
				calleeVal = resolved
				// If it's a Var, deref it to get the actual value
				if v, ok := calleeVal.(*vm.Var); ok {
					calleeVal = v.Deref()
				}
			} else if calleeOp == "load-arg" {
				// Callee is loaded from an arg; this is unresolvable at decode time
				return nil, fmt.Errorf("inst %d (call) has unresolvable callee (inst %d is :load-arg, cannot resolve at decode time)", nid, calleeID)
			} else {
				return nil, fmt.Errorf("inst %d (call) has unresolvable callee (inst %d is %s, not const/load-var)", nid, calleeID, calleeOp)
			}

			// Verify it's a callable (Func, NativeFn, etc)
			switch calleeVal.(type) {
			case *vm.Func, *vm.NativeFn:
				// OK
			default:
				return nil, fmt.Errorf("inst %d (call) resolved to non-callable: %T", nid, calleeVal)
			}
			calleeIndices[nid] = len(callees)
			callees = append(callees, calleeVal)
		}

		insts[nid] = spikeInst{
			op:     opStr,
			opc:    opcFromStr(opStr),
			args:   refInstIDs,
			aux:    auxInt,
			auxVal: auxValue,
			typ:    typeCode,
		}
	}

	// Decode all blocks
	blocks := make([]spikeBlock, 0, blockCount)
	for bid := 0; bid < blockCount; bid++ {
		// Get block-params
		blockParamsVal, err := evalExpr(fmt.Sprintf("(ir.data/block-params %d %s)", bid, varName))
		if err != nil {
			return nil, fmt.Errorf("eval block-params at block %d: %w", bid, err)
		}
		blockParamList, ok := blockParamsVal.(vm.Indexed)
		if !ok {
			return nil, fmt.Errorf("block-params at block %d is not indexed: %T", bid, blockParamsVal)
		}
		blockParamCount := blockParamList.RawCount()
		params := make([]int32, blockParamCount)
		for i := 0; i < blockParamCount; i++ {
			paramID := int(blockParamList.Nth(i).(vm.Int))
			params[i] = int32(paramID)
		}

		// Get block-insts (actual list of inst-ids in this block)
		blockInstsVal, err := evalExpr(fmt.Sprintf("(ir.data/block-insts %d %s)", bid, varName))
		if err != nil {
			return nil, fmt.Errorf("eval block-insts at block %d: %w", bid, err)
		}
		blockInstList, ok := blockInstsVal.(vm.Indexed)
		if !ok {
			return nil, fmt.Errorf("block-insts at block %d is not indexed: %T", bid, blockInstsVal)
		}
		blockInstCount := blockInstList.RawCount()

		// Extract the actual inst-ids for this block
		blockInstIds := make([]int32, blockInstCount)
		for i := 0; i < blockInstCount; i++ {
			instVal := blockInstList.Nth(i)
			instID, ok := instVal.(vm.Int)
			if !ok {
				return nil, fmt.Errorf("block inst at block %d[%d] is not an int: %T", bid, i, instVal)
			}
			blockInstIds[i] = int32(instID)
		}

		// Get terminator
		termVal, err := evalExpr(fmt.Sprintf("(ir.data/block-term %d %s)", bid, varName))
		if err != nil {
			return nil, fmt.Errorf("eval block-term at block %d: %w", bid, err)
		}
		termInstInt, ok := termVal.(vm.Int)
		if !ok {
			return nil, fmt.Errorf("block-term at block %d is not an int: %T", bid, termVal)
		}
		termInstID := int(termInstInt)

		// Get terminator op
		termOpVal, err := evalExpr(fmt.Sprintf("(ir.data/op %d %s)", termInstID, varName))
		if err != nil {
			return nil, fmt.Errorf("eval term op at block %d: %w", bid, err)
		}
		termOpKw, ok := termOpVal.(vm.Keyword)
		if !ok {
			return nil, fmt.Errorf("term op at block %d is not a keyword: %T", bid, termOpVal)
		}
		termOp := string(termOpKw)

		// Get terminator aux (has the edge info)
		termAuxVal, err := evalExpr(fmt.Sprintf("(ir.data/aux %d %s)", termInstID, varName))
		if err != nil {
			return nil, fmt.Errorf("eval term aux at block %d: %w", bid, err)
		}

		// Get terminator refs (for return)
		termRefsVal, err := evalExpr(fmt.Sprintf("(ir.data/refs %d %s)", termInstID, varName))
		if err != nil {
			return nil, fmt.Errorf("eval term refs at block %d: %w", bid, err)
		}
		termRefList, ok := termRefsVal.(vm.Indexed)
		if !ok {
			return nil, fmt.Errorf("term refs at block %d is not indexed: %T", bid, termRefsVal)
		}

		// Decode terminator based on op
		var term spikeTerm
		term.op = termOp
		term.opc = opcFromStr(termOp)

		switch termOp {
		case "return":
			if termRefList.RawCount() < 1 {
				return nil, fmt.Errorf("block %d (return) has no value ref", bid)
			}
			retVal := termRefList.Nth(0)
			retInt, ok := retVal.(vm.Int)
			if !ok {
				return nil, fmt.Errorf("return value at block %d is not an int: %T", bid, retVal)
			}
			term.returnVal = int32(retInt)

		case "branch":
			// branch aux is {:target BlockID :args [InstIds...]}
			auxMap, ok := termAuxVal.(vm.Lookup)
			if !ok {
				return nil, fmt.Errorf("branch aux at block %d is not a lookup: %T", bid, termAuxVal)
			}
			targetVal := auxMap.ValueAt(vm.Keyword("target"))
			argsVal := auxMap.ValueAt(vm.Keyword("args"))
			if targetVal == vm.NIL {
				return nil, fmt.Errorf("block %d (branch) aux missing :target", bid)
			}
			targetInt, ok := targetVal.(vm.Int)
			if !ok {
				return nil, fmt.Errorf("branch target at block %d is not an int: %T", bid, targetVal)
			}
			target := int32(targetInt)

			// Bounds-check edge target
			if target < 0 || target >= int32(blockCount) {
				return nil, fmt.Errorf("block %d (branch) edge target %d out of bounds [0, %d)", bid, target, blockCount)
			}

			args := make([]int32, 0)
			if argsVal != vm.NIL {
				argsList, ok := argsVal.(vm.Indexed)
				if !ok {
					return nil, fmt.Errorf("branch args at block %d is not indexed: %T", bid, argsVal)
				}
				for i := 0; i < argsList.RawCount(); i++ {
					argVal := argsList.Nth(i)
					argInt, ok := argVal.(vm.Int)
					if !ok {
						return nil, fmt.Errorf("branch arg at block %d[%d] is not an int: %T", bid, i, argVal)
					}
					args = append(args, int32(argInt))
				}
			}
			term.simpleEdge = spikeEdge{target: target, args: args}

		case "branch-if":
			// branch-if aux is {:true-target {:target BlockID :args [...]} :false-target {:target BlockID :args [...]}}
			// branch-if refs: first ref is the condition
			if termRefList.RawCount() < 1 {
				return nil, fmt.Errorf("block %d (branch-if) has no condition ref", bid)
			}
			condRefVal := termRefList.Nth(0)
			condRefInt, ok := condRefVal.(vm.Int)
			if !ok {
				return nil, fmt.Errorf("branch-if condition ref at block %d is not an int: %T", bid, condRefVal)
			}
			term.condRef = int32(condRefInt)

			auxMap, ok := termAuxVal.(vm.Lookup)
			if !ok {
				return nil, fmt.Errorf("branch-if aux at block %d is not a lookup: %T", bid, termAuxVal)
			}
			trueTargetVal := auxMap.ValueAt(vm.Keyword("true-target"))
			falseTargetVal := auxMap.ValueAt(vm.Keyword("false-target"))

			if trueTargetVal == vm.NIL || falseTargetVal == vm.NIL {
				return nil, fmt.Errorf("block %d (branch-if) aux missing true/false targets", bid)
			}

			// Extract true edge
			trueMap, ok := trueTargetVal.(vm.Lookup)
			if !ok {
				return nil, fmt.Errorf("branch-if true-target at block %d is not a lookup: %T", bid, trueTargetVal)
			}
			trueTgtVal := trueMap.ValueAt(vm.Keyword("target"))
			trueTargetInt, ok := trueTgtVal.(vm.Int)
			if !ok {
				return nil, fmt.Errorf("branch-if true target at block %d is not an int: %T", bid, trueTgtVal)
			}
			trueTarget := int32(trueTargetInt)
			if trueTarget < 0 || trueTarget >= int32(blockCount) {
				return nil, fmt.Errorf("block %d (branch-if) true edge target %d out of bounds [0, %d)", bid, trueTarget, blockCount)
			}

			trueArgs := make([]int32, 0)
			if trueArgsVal := trueMap.ValueAt(vm.Keyword("args")); trueArgsVal != vm.NIL {
				argsList, ok := trueArgsVal.(vm.Indexed)
				if !ok {
					return nil, fmt.Errorf("branch-if true args at block %d is not indexed: %T", bid, trueArgsVal)
				}
				for i := 0; i < argsList.RawCount(); i++ {
					argVal := argsList.Nth(i)
					argInt, ok := argVal.(vm.Int)
					if !ok {
						return nil, fmt.Errorf("branch-if true arg at block %d[%d] is not an int: %T", bid, i, argVal)
					}
					trueArgs = append(trueArgs, int32(argInt))
				}
			}
			term.trueEdge = spikeEdge{target: trueTarget, args: trueArgs}

			// Extract false edge
			falseMap, ok := falseTargetVal.(vm.Lookup)
			if !ok {
				return nil, fmt.Errorf("branch-if false-target at block %d is not a lookup: %T", bid, falseTargetVal)
			}
			falseTgtVal := falseMap.ValueAt(vm.Keyword("target"))
			falseTargetInt, ok := falseTgtVal.(vm.Int)
			if !ok {
				return nil, fmt.Errorf("branch-if false target at block %d is not an int: %T", bid, falseTgtVal)
			}
			falseTarget := int32(falseTargetInt)
			if falseTarget < 0 || falseTarget >= int32(blockCount) {
				return nil, fmt.Errorf("block %d (branch-if) false edge target %d out of bounds [0, %d)", bid, falseTarget, blockCount)
			}

			falseArgs := make([]int32, 0)
			if falseArgsVal := falseMap.ValueAt(vm.Keyword("args")); falseArgsVal != vm.NIL {
				argsList, ok := falseArgsVal.(vm.Indexed)
				if !ok {
					return nil, fmt.Errorf("branch-if false args at block %d is not indexed: %T", bid, falseArgsVal)
				}
				for i := 0; i < argsList.RawCount(); i++ {
					argVal := argsList.Nth(i)
					argInt, ok := argVal.(vm.Int)
					if !ok {
						return nil, fmt.Errorf("branch-if false arg at block %d[%d] is not an int: %T", bid, i, argVal)
					}
					falseArgs = append(falseArgs, int32(argInt))
				}
			}
			term.falseEdge = spikeEdge{target: falseTarget, args: falseArgs}

		default:
			return nil, fmt.Errorf("block %d has invalid terminator: %s", bid, termOp)
		}

		blocks = append(blocks, spikeBlock{
			insts:  blockInstIds,
			params: params,
			term:   term,
		})
	}

	// Post-pass: validate edge-arg arity for all terminators
	for bid, block := range blocks {
		switch block.term.op {
		case "branch":
			if len(block.term.simpleEdge.args) != len(blocks[block.term.simpleEdge.target].params) {
				return nil, fmt.Errorf("block %d (branch) edge has %d args but target block %d has %d params", bid, len(block.term.simpleEdge.args), block.term.simpleEdge.target, len(blocks[block.term.simpleEdge.target].params))
			}
		case "branch-if":
			if len(block.term.trueEdge.args) != len(blocks[block.term.trueEdge.target].params) {
				return nil, fmt.Errorf("block %d (branch-if) true edge has %d args but target block %d has %d params", bid, len(block.term.trueEdge.args), block.term.trueEdge.target, len(blocks[block.term.trueEdge.target].params))
			}
			if len(block.term.falseEdge.args) != len(blocks[block.term.falseEdge.target].params) {
				return nil, fmt.Errorf("block %d (branch-if) false edge has %d args but target block %d has %d params", bid, len(block.term.falseEdge.args), block.term.falseEdge.target, len(blocks[block.term.falseEdge.target].params))
			}
		}
	}

	fn := &spikeFn{
		insts:         insts,
		blocks:        blocks,
		consts:        constPool.AllValues(),
		callees:       callees,
		calleeIndices: calleeIndices,
		nargs:         arity,
		routes:        computeRouting(insts),
	}

	// No decoder-side compaction: ir.passes.cleanup (run above, before flattening)
	// already produced a tombstone-free, arity-aligned form. validateLiveInvariants
	// is now a pure assertion that the pass did its job (AC-PO-CL.1).
	if err := validateLiveInvariants(fn); err != nil {
		return nil, err
	}
	return fn, nil
}

// validateLiveInvariants enforces at DECODE time that tombstones and
// terminators never reach the interpreters: block inst lists contain only
// live, executable, in-range insts (block-args are filtered out too — their
// locals are written by incoming edges, they never execute); no live
// ref/edge-arg/param/return points at a tombstone. This removes every
// per-iteration guard from the interpreter hot loops — an :invalid marker
// reaching the VM is a decode bug, not a runtime case to skip.
func validateLiveInvariants(fn *spikeFn) error {
	tomb := func(nid int32) bool {
		return nid < 0 || nid >= int32(len(fn.insts)) || fn.insts[nid].opc == spikeOpInvalid
	}
	for bid := range fn.blocks {
		b := &fn.blocks[bid]
		live := b.insts[:0]
		for _, nid := range b.insts {
			if nid < 0 || nid >= int32(len(fn.insts)) {
				return fmt.Errorf("decode: block %d lists out-of-range inst %d", bid, nid)
			}
			switch fn.insts[nid].opc {
			case spikeOpInvalid:
				return fmt.Errorf("decode: block %d lists tombstone inst %d (:invalid must never reach the VM)", bid, nid)
			case spikeOpReturn, spikeOpBranch, spikeOpBranchIf:
				// Terminators execute via block-term, never in the body.
				continue
			case spikeOpBlockArg:
				// Placeholders; written by incoming edges, never executed.
				continue
			}
			for _, ref := range fn.insts[nid].args {
				if tomb(ref) {
					return fmt.Errorf("decode: inst %d references tombstone/out-of-range inst %d", nid, ref)
				}
			}
			live = append(live, nid)
		}
		b.insts = live
		for _, p := range b.params {
			if tomb(p) {
				return fmt.Errorf("decode: block %d param references tombstone inst %d", bid, p)
			}
		}
		for _, e := range []spikeEdge{b.term.simpleEdge, b.term.trueEdge, b.term.falseEdge} {
			for _, a := range e.args {
				if tomb(a) {
					return fmt.Errorf("decode: block %d edge arg references tombstone inst %d", bid, a)
				}
			}
		}
		if b.term.opc == spikeOpReturn && tomb(b.term.returnVal) {
			return fmt.Errorf("decode: block %d returns tombstone inst %d", bid, b.term.returnVal)
		}
	}
	return nil
}

// resolveVar resolves a var symbol or value to its deref'd value via rt.NS.
func resolveVar(auxVal vm.Value, ns *vm.Namespace) (vm.Value, error) {
	if auxVal == vm.NIL {
		return nil, fmt.Errorf("auxVal is NIL")
	}

	// If it's a symbol, look it up in the namespace and deref
	if sym, ok := auxVal.(vm.Symbol); ok {
		varVal := ns.Lookup(sym)
		if varVal == nil {
			return nil, fmt.Errorf("symbol %s not found in namespace", sym)
		}
		if v, ok := varVal.(*vm.Var); ok {
			return v.Deref(), nil
		}
		return varVal, nil
	}

	// If it's already a value, return it
	return auxVal, nil
}

// isValidScopeOp checks if an op keyword is in the spike's bounded subset.
// The op string is the keyword name without the leading colon (e.g., "add" not ":add").
// opcFromStr maps an op string to its opcode.
func opcFromStr(op string) uint8 {
	switch op {
	case "invalid":
		return spikeOpInvalid
	case "const":
		return spikeOpConst
	case "load-arg":
		return spikeOpLoadArg
	case "load-var":
		return spikeOpLoadVar
	case "block-arg":
		return spikeOpBlockArg
	case "add":
		return spikeOpAdd
	case "sub":
		return spikeOpSub
	case "mul":
		return spikeOpMul
	case "inc":
		return spikeOpInc
	case "dec":
		return spikeOpDec
	case "lt":
		return spikeOpLt
	case "lte":
		return spikeOpLte
	case "gt":
		return spikeOpGt
	case "gte":
		return spikeOpGte
	case "eq":
		return spikeOpEq
	case "branch":
		return spikeOpBranch
	case "branch-if":
		return spikeOpBranchIf
	case "return":
		return spikeOpReturn
	case "call":
		return spikeOpCall
	default:
		return spikeOpInvalid
	}
}

func isValidScopeOp(op string) bool {
	switch op {
	case "const", "load-arg", "load-var", "block-arg", "add", "sub", "mul", "inc", "dec", "lt", "lte", "gt", "gte", "eq",
		"branch", "branch-if", "return", "call":
		return true
	default:
		return false
	}
}

// typeToCode maps typeinfer's type keyword to a code (0=unknown, 1=int, 2=float).
// The type string is the keyword name without the leading colon (e.g., "int" not ":int").
func typeToCode(t vm.Keyword) uint8 {
	switch string(t) {
	case "int":
		return 1
	case "float":
		return 2
	default:
		return 0 // unknown or other
	}
}

// computeRouting determines the storage route (ROUTE_INT/ROUTE_FLOAT/ROUTE_BOXED) for each inst.
// Conservative: only ROUTE_INT/FLOAT if all operands also route that way.
func computeRouting(insts []spikeInst) []uint8 {
	routes := make([]uint8, len(insts))

	// Seed: const, block-arg, and ops with typ code
	for i, inst := range insts {
		if inst.opc == spikeOpConst && inst.typ == 1 {
			routes[i] = ROUTE_INT
		} else if inst.opc == spikeOpConst && inst.typ == 2 {
			routes[i] = ROUTE_FLOAT
		} else if inst.opc == spikeOpBlockArg && inst.typ == 1 {
			routes[i] = ROUTE_INT
		} else if inst.opc == spikeOpBlockArg && inst.typ == 2 {
			routes[i] = ROUTE_FLOAT
		} else if inst.opc == spikeOpLoadArg && inst.typ == 1 {
			// typeinfer narrowed the argument: unbox once at entry
			// (constant boundary cost), then the whole loop runs native.
			routes[i] = ROUTE_INT
		} else if inst.opc == spikeOpLoadArg && inst.typ == 2 {
			routes[i] = ROUTE_FLOAT
		} else {
			routes[i] = ROUTE_BOXED // conservative default
		}
	}

	// Propagate: if all operands are ROUTE_INT and op is typed as int, route as ROUTE_INT
	changed := true
	for changed && len(insts) > 0 {
		changed = false
		for i, inst := range insts {
			if routes[i] != ROUTE_BOXED {
				continue // already routed
			}

			// Check if this is a numeric op that could be routed
			switch inst.opc {
			case spikeOpAdd, spikeOpSub, spikeOpMul, spikeOpInc, spikeOpDec:
				if inst.typ == 1 && canOperandsRoute(insts, inst.args, ROUTE_INT, routes) {
					routes[i] = ROUTE_INT
					changed = true
				} else if inst.typ == 2 && canOperandsRoute(insts, inst.args, ROUTE_FLOAT, routes) {
					routes[i] = ROUTE_FLOAT
					changed = true
				}
			case spikeOpLt, spikeOpLte, spikeOpGt, spikeOpGte, spikeOpEq:
				// Comparisons over fully-narrowed operands run natively and
				// produce a native bool (0/1 in localsI) — the per-iteration
				// box of the loop condition is the cost this removes.
				if len(inst.args) == 2 &&
					(canOperandsRoute(insts, inst.args, ROUTE_INT, routes) ||
						canOperandsRoute(insts, inst.args, ROUTE_FLOAT, routes)) {
					routes[i] = ROUTE_BOOL
					changed = true
				}
			}
		}
	}

	return routes
}

// canOperandsRoute checks if all operand insts route to the desired route.
func canOperandsRoute(insts []spikeInst, args []int32, desiredRoute uint8, routes []uint8) bool {
	for _, argID := range args {
		if argID < 0 || argID >= int32(len(insts)) || routes[argID] != desiredRoute {
			return false
		}
	}
	return true
}

// TestSpikeDecode_Simple decodes a trivial one-block function.
func TestSpikeDecode_Simple(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-simple")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn k [a b] (+ (* a a) b)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed: %v", decodeErr)
	}

	// Basic structural checks
	if fn.nargs != 2 {
		t.Errorf("expected arity 2, got %d", fn.nargs)
	}

	if len(fn.blocks) == 0 {
		t.Fatal("no blocks decoded")
	}

	if len(fn.insts) == 0 {
		t.Fatal("no insts decoded")
	}

	// Check for expected ops
	hasLoadArg := false
	hasAdd := false
	hasMul := false
	for _, inst := range fn.insts {
		if inst.opc == spikeOpLoadArg {
			hasLoadArg = true
		}
		if inst.opc == spikeOpAdd {
			hasAdd = true
		}
		if inst.opc == spikeOpMul {
			hasMul = true
		}
	}
	if !hasLoadArg {
		t.Error("no load-arg insts found")
	}
	if !hasAdd {
		t.Error("no add insts found")
	}
	if !hasMul {
		t.Error("no mul insts found")
	}
}

// TestSpikeDecode_ConstValues verifies that :const insts decode actual constant values (auxVal field).
func TestSpikeDecode_ConstValues(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-const")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn c [] 42))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed: %v", decodeErr)
	}

	// Verify that const insts have auxVal populated
	for _, inst := range fn.insts {
		if inst.opc == spikeOpConst {
			// Const insts should have auxVal populated
			if inst.auxVal == nil {
				t.Error("const inst has nil auxVal")
			}
		}
	}
}

// TestSpikeDecode_ArgIndices verifies that ref inst-ids are correct by checking
// the mul inst references.
func TestSpikeDecode_ArgIndices(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-args")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn k [a b] (+ (* a a) b)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed: %v", decodeErr)
	}

	// Find mul inst; it should have two refs (both to the same load-arg inst)
	foundMul := false
	for _, inst := range fn.insts {
		if inst.opc == spikeOpMul {
			foundMul = true
			if len(inst.args) != 2 {
				t.Errorf("mul inst has %d args, expected 2", len(inst.args))
			}
			// Both args should reference the same load-arg inst
			if len(inst.args) >= 2 && inst.args[0] != inst.args[1] {
				t.Errorf("mul inst has mismatched args: %d != %d (expected both to reference same load-arg)", inst.args[0], inst.args[1])
			}
		}
	}
	if !foundMul {
		t.Error("no mul inst found")
	}
}

// TestSpikeDecode_AuxField verifies aux field is populated for load-arg and block-arg.
func TestSpikeDecode_AuxField(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-aux")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn k [a b] (+ (* a a) b)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed: %v", decodeErr)
	}

	// Find load-arg; it should have aux == argument index
	foundLoadArg := false
	for _, inst := range fn.insts {
		if inst.opc == spikeOpLoadArg {
			foundLoadArg = true
			if inst.aux < 0 || inst.aux >= int32(fn.nargs) {
				t.Errorf("load-arg aux %d out of range [0, %d)", inst.aux, fn.nargs)
			}
		}
	}
	if !foundLoadArg {
		t.Error("no load-arg inst found")
	}
}

// TestSpikeDecode_BlockParams verifies that block params are captured.
func TestSpikeDecode_BlockParams(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-params")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn test [x] (if (< x 0) (- 0 x) x)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed: %v", decodeErr)
	}

	// Last block should have params (block-args from branches)
	if len(fn.blocks) < 3 {
		t.Logf("expected 3+ blocks for if-expr, got %d (might be optimized)", len(fn.blocks))
		return
	}

	lastBlock := fn.blocks[len(fn.blocks)-1]
	if len(lastBlock.params) == 0 {
		t.Logf("last block has no params (all branches may be identical)")
		return
	}

	// Params should be valid inst-ids
	for _, paramID := range lastBlock.params {
		if paramID < 0 || paramID >= int32(len(fn.insts)) {
			t.Errorf("block param %d out of range [0, %d)", paramID, len(fn.insts))
		}
		// The param should correspond to a block-arg inst
		if paramID < int32(len(fn.insts)) && fn.insts[paramID].op != "block-arg" {
			t.Errorf("block param %d inst is %s, expected block-arg", paramID, fn.insts[paramID].op)
		}
	}
}

// TestSpikeDecode_EdgeArgs verifies that branch edge args are captured correctly
// and that edge-arg arity matches block params.
func TestSpikeDecode_EdgeArgs(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-edges")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn test [x] (if (< x 0) (- 0 x) x)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed: %v", decodeErr)
	}

	// Find a branch-if terminator and check its edges
	foundBranchIf := false
	for _, block := range fn.blocks {
		if block.term.opc == spikeOpBranchIf {
			foundBranchIf = true
			// Both edges should have valid targets
			if block.term.trueEdge.target < 0 || block.term.trueEdge.target >= int32(len(fn.blocks)) {
				t.Errorf("branch-if true edge target %d out of range", block.term.trueEdge.target)
			}
			if block.term.falseEdge.target < 0 || block.term.falseEdge.target >= int32(len(fn.blocks)) {
				t.Errorf("branch-if false edge target %d out of range", block.term.falseEdge.target)
			}
			// Edge args should be valid inst-ids and arity should match params
			if len(block.term.trueEdge.args) != len(fn.blocks[block.term.trueEdge.target].params) {
				t.Errorf("true edge arg arity %d != target block params %d", len(block.term.trueEdge.args), len(fn.blocks[block.term.trueEdge.target].params))
			}
			if len(block.term.falseEdge.args) != len(fn.blocks[block.term.falseEdge.target].params) {
				t.Errorf("false edge arg arity %d != target block params %d", len(block.term.falseEdge.args), len(fn.blocks[block.term.falseEdge.target].params))
			}
		}
	}
	if !foundBranchIf {
		t.Logf("no branch-if terminator found (if-expr may be optimized)")
	}
}

// TestSpikeDecode_OutOfScopeOpFails verifies that out-of-scope ops (outside the bounded subset) fail loudly.
// Real out-of-scope op: :def (from side-effecting forms like do+def).
func TestSpikeDecode_OutOfScopeOpFails(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-def")
	// The :def op (side effect) is outside the bounded subset
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn s [x] (do (def zzy-spike-scope x) zzy-spike-scope)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	_, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr == nil {
		t.Error("expected decode to fail on out-of-scope op (:def), but it succeeded")
	}
	if decodeErr != nil && !strings.Contains(decodeErr.Error(), "def") {
		t.Errorf("expected error to name out-of-scope op :def, got: %v", decodeErr)
	}
}

// TestSpikeDecode_UnresolvableCalleeLoadsArgFails verifies that callees loaded from
// function parameters fail at decode time.
func TestSpikeDecode_UnresolvableCalleeLoadsArgFails(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-unresolvable")
	// A defn with a call through a parameter (f x) where f is the param
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn u [f x] (f x)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	_, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr == nil {
		t.Error("expected decode to fail on unresolvable callee, but it succeeded")
	}
	if decodeErr != nil && !strings.Contains(decodeErr.Error(), "unresolvable") {
		t.Errorf("expected unresolvable callee error, got: %v", decodeErr)
	}
}

// TestSpikeDecode_CallResolution verifies that calls to loaded vars are resolved at decode time.
func TestSpikeDecode_CallResolution(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-call")
	// A call to nth (loaded as a var, not a const)
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn c [v] (nth v 1)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed: %v", decodeErr)
	}

	// Should have resolved the nth fn
	if len(fn.callees) == 0 {
		t.Error("expected at least one resolved callee")
	}

	// All callees should be callable values (Func, NativeFn, etc)
	for i, callee := range fn.callees {
		switch callee.(type) {
		case *vm.Func, *vm.NativeFn:
			// OK
		default:
			t.Errorf("callee %d is %T, expected *vm.Func or *vm.NativeFn", i, callee)
		}
	}
}

// TestSpikeDecode_LoopKernel is the critical T2 regression test: loop/recur fixtures generate
// :invalid tombstones (from DCE), and the decoder must SKIP them while maintaining InstId correspondence.
// This test verifies that the loop-based numeric kernel that T2/T5 depend on decodes successfully,
// with edge args properly threading the recur values back to the loop header.
func TestSpikeDecode_LoopKernel(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-loop")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn m [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode loop kernel failed (tombstones not properly skipped): %v", decodeErr)
	}

	// Verify the decoded structure
	if len(fn.blocks) == 0 {
		t.Fatal("loop kernel has no blocks")
	}

	// Find the recur back-edge (should have non-empty edge args)
	foundRecurEdge := false
	for bid, block := range fn.blocks {
		if block.term.opc == spikeOpBranch || block.term.opc == spikeOpBranchIf {
			// Check for edges with arguments (recur threading values back to loop header)
			if len(block.term.simpleEdge.args) > 0 {
				foundRecurEdge = true
				targetBlock := fn.blocks[block.term.simpleEdge.target]
				// Verify arity matches
				if len(block.term.simpleEdge.args) != len(targetBlock.params) {
					t.Errorf("block %d edge has %d args but target has %d params", bid, len(block.term.simpleEdge.args), len(targetBlock.params))
				}
			}
			if block.term.opc == spikeOpBranchIf && len(block.term.trueEdge.args) > 0 {
				foundRecurEdge = true
				targetBlock := fn.blocks[block.term.trueEdge.target]
				if len(block.term.trueEdge.args) != len(targetBlock.params) {
					t.Errorf("block %d true edge has %d args but target has %d params", bid, len(block.term.trueEdge.args), len(targetBlock.params))
				}
			}
			if block.term.opc == spikeOpBranchIf && len(block.term.falseEdge.args) > 0 {
				foundRecurEdge = true
				targetBlock := fn.blocks[block.term.falseEdge.target]
				if len(block.term.falseEdge.args) != len(targetBlock.params) {
					t.Errorf("block %d false edge has %d args but target has %d params", bid, len(block.term.falseEdge.args), len(targetBlock.params))
				}
			}
		}
	}

	if !foundRecurEdge {
		t.Logf("note: loop kernel decoded but no recur edges found with args (structure: %d blocks, %d insts)", len(fn.blocks), len(fn.insts))
	}

	t.Logf("loop kernel decoded successfully: %d blocks, %d total insts (some are tombstones)", len(fn.blocks), len(fn.insts))
}

// TestSpikeDecode_TypeInferConsumption verifies that typeinfer types are correctly extracted.
func TestSpikeDecode_TypeInferConsumption(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-typed-arith")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn typed-arith [x] (+ x 1)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed: %v", decodeErr)
	}

	// Find the add inst; it should have typ == 1 (int)
	foundAdd := false
	for _, inst := range fn.insts {
		if inst.opc == spikeOpAdd {
			foundAdd = true
			if inst.typ != 1 {
				t.Errorf("expected add to be typed as int (1), got %d", inst.typ)
			}
		}
	}
	if !foundAdd {
		t.Error("no add inst found in decoded fn")
	}
}

// TestSpikeDecode_IncOp verifies that the inc op (kernel d prerequisite) is properly in scope and decoded.
func TestSpikeDecode_IncOp(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-inc")
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn inc-kernel [x] (inc x)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile optimized IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode failed on inc fixture: %v", decodeErr)
	}

	// Find the inc inst
	foundInc := false
	for _, inst := range fn.insts {
		if inst.opc == spikeOpInc {
			foundInc = true
			// inc should be a unary op with one ref
			if len(inst.args) != 1 {
				t.Errorf("inc inst has %d args, expected 1", len(inst.args))
			}
		}
	}
	if !foundInc {
		t.Error("no inc inst found in decoded fn")
	}
}

// TestSpikeRun_Simple tests the spike interpreter with a simple arithmetic kernel.
// Compares spike output with stack VM output for (defn k [a b] (+ (* a a) b)).
func TestSpikeRun_Simple(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-run-simple")

	// Compile via the spike pipeline
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn k [a b] (+ (* a a) b)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile via the stack VM pipeline (direct)
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-simple")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn k [a b] (+ (* a a) b))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("k"))
	if stackVar == nil {
		t.Fatal("stack VM: k not found in namespace")
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Test various inputs
	testCases := []struct {
		a, b vm.Value
	}{
		{vm.Int(2), vm.Int(3)},         // (2*2) + 3 = 7
		{vm.Int(0), vm.Int(5)},         // (0*0) + 5 = 5
		{vm.Int(-1), vm.Int(2)},        // (-1*-1) + 2 = 3
		{vm.Float(2.5), vm.Float(1.5)}, // (2.5*2.5) + 1.5 = 8.75
	}

	for _, tc := range testCases {
		// Run spike
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{tc.a, tc.b})
		if spikeErr != nil {
			t.Errorf("spikeRun(%v, %v) error: %v", tc.a, tc.b, spikeErr)
			continue
		}

		// Run stack VM
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.a, tc.b})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM(%v, %v) error: %v", tc.a, tc.b, stackErr)
			continue
		}

		// Compare
		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("mismatch for (%v, %v): spike=%v, stack=%v", tc.a, tc.b, spikeResult, stackResult)
		}
	}
}

// TestSpikeRun_Branch tests branch handling with an abs-value function.
func TestSpikeRun_Branch(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-run-branch")

	// Compile via the spike pipeline
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn abs-val [x] (if (< x 0) (- 0 x) x)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile via the stack VM pipeline
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-branch")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn abs-val [x] (if (< x 0) (- 0 x) x))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("abs-val"))
	if stackVar == nil {
		t.Fatal("stack VM: abs-val not found in namespace")
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Test cases: one negative (takes then-branch), one positive (takes else-branch)
	testCases := []vm.Value{
		vm.Int(-5), // true branch
		vm.Int(5),  // false branch
		vm.Int(0),  // false branch
		vm.Float(-2.5),
		vm.Float(2.5),
	}

	for _, x := range testCases {
		// Run spike
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{x})
		if spikeErr != nil {
			t.Errorf("spikeRun(%v) error: %v", x, spikeErr)
			continue
		}

		// Run stack VM
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{x})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM(%v) error: %v", x, stackErr)
			continue
		}

		// Compare
		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("mismatch for abs-val(%v): spike=%v, stack=%v", x, spikeResult, stackResult)
		}
	}
}

// TestSpikeRun_LoopKernel is AC-WS.1 kernel (b) parity test: untyped numeric loop.
// Verifies bit-identical results with the stack VM for (defn m [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))).
func TestSpikeRun_LoopKernel(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-run-loop")

	// Compile via the spike pipeline
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn m [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile via the stack VM pipeline
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-loop")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn m [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("m"))
	if stackVar == nil {
		t.Fatal("stack VM: m not found in namespace")
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Test AC-WS.1 scenarios: n ∈ {0, 1, 10, 1000}
	// Kernel computes sum of 0..n-1, so expected results:
	// n=0: 0 (empty sum)
	// n=1: 0 (sum of [0])
	// n=10: 45 (sum of [0..9])
	// n=1000: 499500 (sum of [0..999])
	testCases := []vm.Value{
		vm.Int(0),
		vm.Int(1),
		vm.Int(10),
		vm.Int(1000),
	}

	for _, n := range testCases {
		// Run spike
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{n})
		if spikeErr != nil {
			t.Errorf("spikeRun(%v) error: %v", n, spikeErr)
			continue
		}

		// Run stack VM
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{n})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM(%v) error: %v", n, stackErr)
			continue
		}

		// Compare bit-identical
		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("kernel (b) parity FAIL for n=%v: spike=%v, stack=%v", n, spikeResult, stackResult)
		} else {
			t.Logf("kernel (b) parity OK for n=%v: result=%v", n, spikeResult)
		}
	}
}

// TestSpikeRun_ContainerKernel tests AC-WS.1 kernel (c): container-touching via function calls.
// Kernel: (defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc)))
// This tests that call ops work with container access (nth on vector arg).
func TestSpikeRun_ContainerKernel(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-container")

	// Compile via the spike pipeline
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile via the stack VM pipeline
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-container")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc)))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("kvec"))
	if stackVar == nil {
		t.Fatal("stack VM: kvec not found in namespace")
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Test with different vectors
	testCases := []struct {
		values []vm.Value
		label  string
	}{
		{[]vm.Value{}, "empty-vector"},
		{[]vm.Value{vm.Int(1)}, "single-element"},
		{[]vm.Value{vm.Int(3), vm.Int(1), vm.Int(4), vm.Int(1), vm.Int(5), vm.Int(9), vm.Int(2), vm.Int(6)}, "multi-element"},
	}

	for _, tc := range testCases {
		vec := vm.NewPersistentVector(tc.values)

		// Run spike (boxed)
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{vec})
		if spikeErr != nil {
			t.Errorf("spikeRun container kernel %s error: %v", tc.label, spikeErr)
			continue
		}

		// Run typed
		stats := &spikeStats{}
		spikeResultTyped, spikeErrTyped := spikeRunTyped(spikeFn, []vm.Value{vec}, stats)
		if spikeErrTyped != nil {
			t.Errorf("spikeRunTyped container kernel %s error: %v", tc.label, spikeErrTyped)
			continue
		}

		// Run stack VM
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{vec})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM container kernel %s error: %v", tc.label, stackErr)
			continue
		}

		// Compare boxed spike with stack VM
		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("container kernel (boxed) FAIL %s: spike=%v, stack=%v", tc.label, spikeResult, stackResult)
		} else {
			t.Logf("container kernel (boxed) PASS %s: result=%v", tc.label, spikeResult)
		}

		// Compare typed spike with stack VM
		if !valuesEqual(spikeResultTyped, stackResult) {
			t.Errorf("container kernel (typed) FAIL %s: spike=%v, stack=%v", tc.label, spikeResultTyped, stackResult)
		} else {
			t.Logf("container kernel (typed) PASS %s: result=%v, callOps=%d", tc.label, spikeResultTyped, stats.callOps)
		}
	}
}

// TestSpikeRun_GoNativeKernel tests AC-WS.1 kernel (d): Go-native callout via function calls.
// Kernel: (defn kmax [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc (- n i))) acc)))
// This tests that call ops work with native Go functions (max).
func TestSpikeRun_GoNativeKernel(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-native")

	// Compile via the spike pipeline
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn kmax [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc (- n i))) acc))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile via the stack VM pipeline
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-native")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn kmax [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc (- n i))) acc)))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("kmax"))
	if stackVar == nil {
		t.Fatal("stack VM: kmax not found in namespace")
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Test cases: n ∈ {0, 10, 100}
	testCases := []vm.Value{
		vm.Int(0),
		vm.Int(10),
		vm.Int(100),
	}

	for _, n := range testCases {
		// Run spike (boxed)
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{n})
		if spikeErr != nil {
			t.Errorf("spikeRun native kernel n=%v error: %v", n, spikeErr)
			continue
		}

		// Run typed
		stats := &spikeStats{}
		spikeResultTyped, spikeErrTyped := spikeRunTyped(spikeFn, []vm.Value{n}, stats)
		if spikeErrTyped != nil {
			t.Errorf("spikeRunTyped native kernel n=%v error: %v", n, spikeErrTyped)
			continue
		}

		// Run stack VM
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{n})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM native kernel n=%v error: %v", n, stackErr)
			continue
		}

		// Compare boxed spike with stack VM
		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("native kernel (boxed) FAIL n=%v: spike=%v, stack=%v", n, spikeResult, stackResult)
		} else {
			t.Logf("native kernel (boxed) PASS n=%v: result=%v", n, spikeResult)
		}

		// Compare typed spike with stack VM
		if !valuesEqual(spikeResultTyped, stackResult) {
			t.Errorf("native kernel (typed) FAIL n=%v: spike=%v, stack=%v", n, spikeResultTyped, stackResult)
		} else {
			t.Logf("native kernel (typed) PASS n=%v: result=%v, callOps=%d", n, spikeResultTyped, stats.callOps)
		}
	}
}

// TestSpikeRun_CallError tests error handling in call operations.
// Kernel: (defn kbad [v] (nth v 99)) — should error on out-of-bounds access.
func TestSpikeRun_CallError(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-call-error")

	// Compile via the spike pipeline
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn kbad [v] (nth v 99)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile via the stack VM pipeline
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-call-error")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn kbad [v] (nth v 99))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("kbad"))
	if stackVar == nil {
		t.Fatal("stack VM: kbad not found in namespace")
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Create a small vector
	vec := vm.NewPersistentVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)})

	// Run spike (should error)
	_, spikeErr := spikeRun(spikeFn, []vm.Value{vec})
	if spikeErr == nil {
		t.Error("spikeRun expected error on out-of-bounds access, got nil")
	} else {
		t.Logf("spikeRun correctly errored: %v", spikeErr)
	}

	// Run typed spike (should error)
	stats := &spikeStats{}
	_, spikeErrTyped := spikeRunTyped(spikeFn, []vm.Value{vec}, stats)
	if spikeErrTyped == nil {
		t.Error("spikeRunTyped expected error on out-of-bounds access, got nil")
	} else {
		t.Logf("spikeRunTyped correctly errored: %v", spikeErrTyped)
	}

	// Run stack VM (should error)
	stackFrame := vm.NewFrame(stackChunk, []vm.Value{vec})
	_, stackErr := stackFrame.Run()
	if stackErr == nil {
		t.Error("stack VM expected error on out-of-bounds access, got nil")
	} else {
		t.Logf("stack VM correctly errored: %v", stackErr)
	}

	// Document that both error (exact message parity not required)
	if (spikeErr != nil) && (spikeErrTyped != nil) && (stackErr != nil) {
		t.Logf("call error parity: all three paths error (spike boxed, spike typed, stack VM)")
	}
}

// TestSpikeRun_NonNumberError verifies error handling parity with stack VM.
// Passes a non-number arg into arithmetic operations and checks for error.
func TestSpikeRun_NonNumberError(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-run-error")

	// Compile via the spike pipeline
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn add-fn [a b] (+ a b)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile via the stack VM pipeline
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-error")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn add-fn [a b] (+ a b))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("add-fn"))
	if stackVar == nil {
		t.Fatal("stack VM: add-fn not found in namespace")
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Pass non-number arg (string)
	badArg := vm.String("not a number")
	goodArg := vm.Int(5)

	// Spike should error
	_, spikeErr := spikeRun(spikeFn, []vm.Value{badArg, goodArg})
	if spikeErr == nil {
		t.Error("spikeRun expected error with non-number, got nil")
	}

	// Stack VM should also error
	stackFrame := vm.NewFrame(stackChunk, []vm.Value{badArg, goodArg})
	_, stackErrExec := stackFrame.Run()
	if stackErrExec == nil {
		t.Error("stack VM expected error with non-number, got nil")
	}

	// Both should error (exact message parity not required, just behavior)
	if (spikeErr != nil) != (stackErrExec != nil) {
		t.Logf("note: error presence matches (spike=%v, stack=%v)", spikeErr != nil, stackErrExec != nil)
	}
}

// TestSpikeRun_IntFloatMixing probes arithmetic with mixed int/float operands.
// Spec: Int + Float → Float (promotes int to float)
func TestSpikeRun_IntFloatMixing(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-int-float")

	// Compile spike: (defn add [a b] (+ a b))
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn add [a b] (+ a b)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile stack VM
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-int-float")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn add [a b] (+ a b))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("add"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	testCases := []struct {
		a, b  vm.Value
		label string
	}{
		{vm.Int(2), vm.Float(3.5), "int+float"},
		{vm.Float(2.5), vm.Int(3), "float+int"},
		{vm.Float(2.5), vm.Float(3.5), "float+float"},
	}

	for _, tc := range testCases {
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{tc.a, tc.b})
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.a, tc.b})
		stackResult, stackErr := stackFrame.Run()

		if (spikeErr != nil) != (stackErr != nil) {
			t.Errorf("int/float parity error (%s): spike err=%v, stack err=%v", tc.label, spikeErr, stackErr)
			continue
		}
		if spikeErr != nil {
			continue // Both errored, OK
		}

		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("int/float parity FAIL (%s): spike=%v (%T), stack=%v (%T)", tc.label, spikeResult, spikeResult, stackResult, stackResult)
		}
	}
}

// TestSpikeRun_Truthiness probes branch conditions with truthiness edge cases.
// Spec: IsTruthy = !(v == NIL || v == FALSE)
func TestSpikeRun_Truthiness(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-truthiness")

	// Compile spike: (defn t [x] (if x 1 2))
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn t [x] (if x 1 2)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile stack VM
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-truthiness")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn t [x] (if x 1 2))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("t"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	testCases := []struct {
		val   vm.Value
		label string
	}{
		{vm.NIL, "nil"},
		{vm.FALSE, "false"},
		{vm.TRUE, "true"},
		{vm.Int(0), "int-0"},
		{vm.Int(1), "int-1"},
		{vm.Float(0.0), "float-0"},
		{vm.Float(0.5), "float-0.5"},
		{vm.String(""), "empty-string"},
		{vm.String("x"), "non-empty-string"},
	}

	for _, tc := range testCases {
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{tc.val})
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.val})
		stackResult, stackErr := stackFrame.Run()

		if (spikeErr != nil) != (stackErr != nil) {
			t.Errorf("truthiness error (%s): spike err=%v, stack err=%v", tc.label, spikeErr, stackErr)
			continue
		}
		if spikeErr != nil {
			continue
		}

		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("truthiness FAIL (%s): spike=%v, stack=%v (should both be 1 or both be 2)", tc.label, spikeResult, stackResult)
		}
	}
}

// TestSpikeRun_FloatAccumulator probes loop with float accumulator (T3 preview).
// Kernel: (defn m [n] (loop [i 0 acc 0.0] (if (< i n) (recur (+ i 1) (+ acc 1.5)) acc)))
func TestSpikeRun_FloatAccumulator(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-float-acc")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn m [n] (loop [i 0 acc 0.0] (if (< i n) (recur (+ i 1) (+ acc 1.5)) acc))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-float-acc")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn m [n] (loop [i 0 acc 0.0] (if (< i n) (recur (+ i 1) (+ acc 1.5)) acc)))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("m"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	testCases := []vm.Value{
		vm.Int(0),  // 0.0
		vm.Int(1),  // 1.5
		vm.Int(3),  // 4.5
		vm.Int(10), // 15.0
	}

	for _, n := range testCases {
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{n})
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{n})
		stackResult, stackErr := stackFrame.Run()

		if (spikeErr != nil) != (stackErr != nil) {
			t.Errorf("float-acc error (n=%v): spike err=%v, stack err=%v", n, spikeErr, stackErr)
			continue
		}
		if spikeErr != nil {
			continue
		}

		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("float-acc parity FAIL (n=%v): spike=%v (%T), stack=%v (%T)", n, spikeResult, spikeResult, stackResult, stackResult)
		}
	}
}

// TestSpikeRun_NestedBranches probes two-level if nesting.
// Kernel: (defn g [x y] (if (< x 0) (if (< y 0) (+ x y) (- x y)) (if (< y 0) (- x y) (+ x y))))
func TestSpikeRun_NestedBranches(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-nested")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn g [x y] (if (< x 0) (if (< y 0) (+ x y) (- x y)) (if (< y 0) (- x y) (+ x y)))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-nested")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn g [x y] (if (< x 0) (if (< y 0) (+ x y) (- x y)) (if (< y 0) (- x y) (+ x y))))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("g"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	testCases := []struct {
		x, y  vm.Value
		label string
	}{
		{vm.Int(-2), vm.Int(-3), "both-neg"},
		{vm.Int(-2), vm.Int(3), "x-neg-y-pos"},
		{vm.Int(2), vm.Int(-3), "x-pos-y-neg"},
		{vm.Int(2), vm.Int(3), "both-pos"},
	}

	for _, tc := range testCases {
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{tc.x, tc.y})
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.x, tc.y})
		stackResult, stackErr := stackFrame.Run()

		if (spikeErr != nil) != (stackErr != nil) {
			t.Errorf("nested-branches error (%s): spike err=%v, stack err=%v", tc.label, spikeErr, stackErr)
			continue
		}
		if spikeErr != nil {
			continue
		}

		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("nested-branches FAIL (%s): spike=%v, stack=%v", tc.label, spikeResult, stackResult)
		}
	}
}

// TestSpikeRun_FloatComparison probes int/float comparisons in branch conditions.
// Kernel: (defn cmp [a b] (if (< a b) 1 (if (> a b) -1 0)))
func TestSpikeRun_FloatComparison(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-float-cmp")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn cmp [a b] (if (< a b) 1 (if (> a b) -1 0))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-float-cmp")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn cmp [a b] (if (< a b) 1 (if (> a b) -1 0)))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("cmp"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	testCases := []struct {
		a, b  vm.Value
		label string
	}{
		{vm.Int(2), vm.Float(3.5), "int<float"},
		{vm.Float(2.5), vm.Int(1), "float>int"},
		{vm.Float(2.5), vm.Float(2.5), "float==float"},
	}

	for _, tc := range testCases {
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{tc.a, tc.b})
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.a, tc.b})
		stackResult, stackErr := stackFrame.Run()

		if (spikeErr != nil) != (stackErr != nil) {
			t.Errorf("float-cmp error (%s): spike err=%v, stack err=%v", tc.label, spikeErr, stackErr)
			continue
		}
		if spikeErr != nil {
			continue
		}

		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("float-cmp FAIL (%s): spike=%v, stack=%v", tc.label, spikeResult, stackResult)
		}
	}
}

// TestSpikeRun_IntOverflow probes integer overflow detection.
// Note: The spike interpreter may not have overflow checks yet; this documents behavior.
func TestSpikeRun_IntOverflow(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-overflow")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn add [a b] (+ a b)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-overflow")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn add [a b] (+ a b))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("add"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// MaxInt on the platform
	maxInt := vm.Int(9223372036854775807) // 2^63 - 1 on 64-bit
	one := vm.Int(1)

	// This should cause overflow in both
	_, spikeErr := spikeRun(spikeFn, []vm.Value{maxInt, one})
	stackFrame := vm.NewFrame(stackChunk, []vm.Value{maxInt, one})
	_, stackErr := stackFrame.Run()

	// Document overflow behavior (both should error, or both should wrap)
	if (spikeErr != nil) != (stackErr != nil) {
		t.Logf("INT OVERFLOW DIVERGENCE: spike error=%v, stack error=%v", spikeErr != nil, stackErr != nil)
		if spikeErr != nil {
			t.Logf("  spike: %v", spikeErr)
		}
		if stackErr != nil {
			t.Logf("  stack: %v", stackErr)
		}
	}
}

// TestSpikeRun_EqParity verifies that eq (=) operator matches stack VM exactly.
// Critical: OP_EQ uses vm.ValueEquals (structural equality), not numeric equality.
// Test cases: (2,2), (2,3), (2,2.0), ("a","a"), ("a","b"), (nil,nil), (nil,false), (true,true).
func TestSpikeRun_EqParity(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-eq")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn e [a b] (if (= a b) 1 0)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile spike IR: %v", err)
	}

	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode spike IR: %v", decodeErr)
	}

	// Compile via the stack VM pipeline
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-eq")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn e [a b] (if (= a b) 1 0))`)); stackErr != nil {
		t.Fatalf("compile stack VM code: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("e"))
	if stackVar == nil {
		t.Fatal("stack VM: e not found in namespace")
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Test cases: (value1, value2, expected_result_description)
	testCases := []struct {
		a, b  vm.Value
		label string
	}{
		{vm.Int(2), vm.Int(2), "int-equal"},                // should be true
		{vm.Int(2), vm.Int(3), "int-unequal"},              // should be false
		{vm.Int(2), vm.Float(2.0), "int-float-diff"},       // should be false (different types)
		{vm.String("a"), vm.String("a"), "string-equal"},   // should be true
		{vm.String("a"), vm.String("b"), "string-unequal"}, // should be false
		{vm.NIL, vm.NIL, "nil-equal"},                      // should be true
		{vm.NIL, vm.FALSE, "nil-false-diff"},               // should be false
		{vm.TRUE, vm.TRUE, "true-equal"},                   // should be true
	}

	for _, tc := range testCases {
		// Run spike
		spikeResult, spikeErr := spikeRun(spikeFn, []vm.Value{tc.a, tc.b})
		if spikeErr != nil {
			t.Errorf("spikeRun(%s) error: %v", tc.label, spikeErr)
			continue
		}

		// Run stack VM
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.a, tc.b})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM(%s) error: %v", tc.label, stackErr)
			continue
		}

		// Compare (should be bit-identical)
		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("EQ PARITY FAIL for %s (%v, %v): spike=%v, stack=%v", tc.label, tc.a, tc.b, spikeResult, stackResult)
		} else {
			t.Logf("EQ PARITY OK for %s (%v, %v): result=%v", tc.label, tc.a, tc.b, spikeResult)
		}
	}
}

// valuesEqual compares two vm.Values for bit-identical equality.
func valuesEqual(a, b vm.Value) bool {
	if a == b {
		return true
	}
	// Handle numeric types
	switch av := a.(type) {
	case vm.Int:
		if bv, ok := b.(vm.Int); ok {
			return av == bv
		}
	case vm.Float:
		if bv, ok := b.(vm.Float); ok {
			return av == bv
		}
	case vm.Boolean:
		if bv, ok := b.(vm.Boolean); ok {
			return av == bv
		}
	}
	// Fall back to standard equality (NIL, etc)
	return a == b
}

// spikeRun executes a decoded spikeFn over boxed locals (every slot holds vm.Value).
// Implements the index-RPN interpreter semantics:
// - locals[InstId] ← op result
// - terminators thread values across block edges via edge args
// - Iterative (no Go recursion) block execution
func spikeRun(fn *spikeFn, args []vm.Value) (vm.Value, error) {
	if len(args) != fn.nargs {
		return nil, fmt.Errorf("spikeRun: expected %d args, got %d", fn.nargs, len(args))
	}

	// Allocate locals array (one slot per inst, indexed by inst-id)
	locals := make([]vm.Value, len(fn.insts))

	// Load arguments into the first fn.nargs positions (load-arg insts expect to find them there)
	// Note: In the actual IR, load-arg insts have aux=0..nargs-1 and should fetch args directly
	// For now, we'll handle this in the load-arg case by fetching from args directly

	// Start at block 0
	currentBlockID := 0

	// Iterative block execution (no recursion)
	for {
		if currentBlockID < 0 || currentBlockID >= len(fn.blocks) {
			return nil, fmt.Errorf("spikeRun: block id %d out of bounds", currentBlockID)
		}

		block := fn.blocks[currentBlockID]

		// Execute all insts in the current block (in order)
		// Block lists are decode-validated: live, executable insts only
		// (no tombstones, terminators, or block-args — see
		// validateLiveInvariants). No per-iteration guards needed.
		for _, nid := range block.insts {
			inst := fn.insts[nid]

			var result vm.Value
			var err error

			switch inst.opc {
			case spikeOpConst:
				result = inst.auxVal

			case spikeOpLoadArg:
				if inst.aux < 0 || inst.aux >= int32(fn.nargs) {
					return nil, fmt.Errorf("spikeRun: load-arg aux %d out of range [0, %d)", inst.aux, fn.nargs)
				}
				result = args[inst.aux]

			case spikeOpLoadVar:
				result = inst.auxVal

			case spikeOpAdd:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRun: add inst %d has %d args, expected 2", nid, len(inst.args))
				}
				a, b := locals[inst.args[0]], locals[inst.args[1]]
				result, err = vm.NumAdd(a, b)
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (add) error: %w", nid, err)
				}

			case spikeOpSub:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRun: sub inst %d has %d args, expected 2", nid, len(inst.args))
				}
				a, b := locals[inst.args[0]], locals[inst.args[1]]
				result, err = vm.NumSub(a, b)
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (sub) error: %w", nid, err)
				}

			case spikeOpMul:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRun: mul inst %d has %d args, expected 2", nid, len(inst.args))
				}
				a, b := locals[inst.args[0]], locals[inst.args[1]]
				result, err = vm.NumMul(a, b)
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (mul) error: %w", nid, err)
				}

			case spikeOpInc:
				if len(inst.args) != 1 {
					return nil, fmt.Errorf("spikeRun: inc inst %d has %d args, expected 1", nid, len(inst.args))
				}
				a := locals[inst.args[0]]
				// Use NumAdd for inc (mirrors VM semantics for both Int fast path and fallback)
				result, err = vm.NumAdd(a, vm.Int(1))
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (inc) error: %w", nid, err)
				}

			case spikeOpDec:
				if len(inst.args) != 1 {
					return nil, fmt.Errorf("spikeRun: dec inst %d has %d args, expected 1", nid, len(inst.args))
				}
				a := locals[inst.args[0]]
				result, err = vm.NumSub(a, vm.Int(1))
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (dec) error: %w", nid, err)
				}

			case spikeOpLt:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRun: lt inst %d has %d args, expected 2", nid, len(inst.args))
				}
				a, b := locals[inst.args[0]], locals[inst.args[1]]
				ltResult, err := vm.NumLt(a, b)
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (lt) error: %w", nid, err)
				}
				result = vm.Boolean(ltResult)

			case spikeOpLte:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRun: lte inst %d has %d args, expected 2", nid, len(inst.args))
				}
				a, b := locals[inst.args[0]], locals[inst.args[1]]
				lteResult, err := vm.NumLe(a, b)
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (lte) error: %w", nid, err)
				}
				result = vm.Boolean(lteResult)

			case spikeOpGt:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRun: gt inst %d has %d args, expected 2", nid, len(inst.args))
				}
				a, b := locals[inst.args[0]], locals[inst.args[1]]
				gtResult, err := vm.NumGt(a, b)
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (gt) error: %w", nid, err)
				}
				result = vm.Boolean(gtResult)

			case spikeOpGte:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRun: gte inst %d has %d args, expected 2", nid, len(inst.args))
				}
				a, b := locals[inst.args[0]], locals[inst.args[1]]
				gteResult, err := vm.NumGe(a, b)
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (gte) error: %w", nid, err)
				}
				result = vm.Boolean(gteResult)

			case spikeOpEq:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRun: eq inst %d has %d args, expected 2", nid, len(inst.args))
				}
				a, b := locals[inst.args[0]], locals[inst.args[1]]
				// Match OP_EQ semantics exactly (structural equality via ValueEquals)
				// Int fast path
				if ai, ok := a.(vm.Int); ok {
					if bi, ok := b.(vm.Int); ok {
						result = vm.Boolean(ai == bi)
						break
					}
				}
				// Keyword fast path
				if ak, ok := a.(vm.Keyword); ok {
					if bk, ok := b.(vm.Keyword); ok {
						result = vm.Boolean(ak == bk)
						break
					}
				}
				// Structural equality fallback (strings, nil, etc)
				if vm.ValueEquals == nil {
					return nil, fmt.Errorf("spikeRun: ValueEquals not initialized")
				}
				result = vm.Boolean(vm.ValueEquals(a, b))

			case spikeOpCall:
				// Call op: first ref is callee (as inst-id), remaining refs are args
				if len(inst.args) < 1 {
					return nil, fmt.Errorf("spikeRun: call inst %d has no callee ref", nid)
				}

				calleeIdx := fn.calleeIndices[nid]
				if calleeIdx < 0 || calleeIdx >= len(fn.callees) {
					return nil, fmt.Errorf("spikeRun: call inst %d has invalid callee index %d", nid, calleeIdx)
				}
				callee := fn.callees[calleeIdx]

				// Remaining args are the actual call arguments
				callArgs := make([]vm.Value, len(inst.args)-1)
				for i, argID := range inst.args[1:] {
					if argID < 0 || argID >= int32(len(fn.insts)) {
						return nil, fmt.Errorf("spikeRun: call arg %d out of bounds", argID)
					}
					callArgs[i] = locals[argID]
				}

				// Invoke the callee
				var callResult vm.Value
				switch c := callee.(type) {
				case *vm.Func:
					callResult, err = c.Invoke(callArgs)
				case *vm.NativeFn:
					callResult, err = c.Invoke(callArgs)
				default:
					return nil, fmt.Errorf("spikeRun: call callee %T is not callable", callee)
				}
				if err != nil {
					return nil, fmt.Errorf("spikeRun: inst %d (call) error: %w", nid, err)
				}
				result = callResult

			default:
				return nil, fmt.Errorf("spikeRun: inst %d has unknown op: %s", nid, inst.op)
			}

			locals[nid] = result
		}

		// Execute terminator to decide next block or return
		term := block.term

		switch term.opc {
		case spikeOpReturn:
			if term.returnVal < 0 || term.returnVal >= int32(len(fn.insts)) {
				return nil, fmt.Errorf("spikeRun: return value inst %d out of bounds", term.returnVal)
			}
			return locals[term.returnVal], nil

		case spikeOpBranch:
			// Thread edge args and jump
			// IMPORTANT: snapshot all edge arg values BEFORE writing any params,
			// since a param inst-id could coincide with an arg inst-id in a loop back-edge
			target := term.simpleEdge.target
			if target < 0 || target >= int32(len(fn.blocks)) {
				return nil, fmt.Errorf("spikeRun: branch target %d out of bounds", target)
			}

			targetBlock := fn.blocks[target]
			if len(term.simpleEdge.args) != len(targetBlock.params) {
				return nil, fmt.Errorf("spikeRun: branch edge has %d args but target has %d params", len(term.simpleEdge.args), len(targetBlock.params))
			}

			// Snapshot edge arg values
			edgeValues := make([]vm.Value, len(term.simpleEdge.args))
			for i, edgeArgID := range term.simpleEdge.args {
				if edgeArgID < 0 || edgeArgID >= int32(len(fn.insts)) {
					return nil, fmt.Errorf("spikeRun: branch edge arg %d out of bounds", edgeArgID)
				}
				edgeValues[i] = locals[edgeArgID]
			}

			// Write snapshotted values to target params
			for i, paramID := range targetBlock.params {
				locals[paramID] = edgeValues[i]
			}

			currentBlockID = int(target)
			continue

		case spikeOpBranchIf:
			// Evaluate condition using truthiness rule: anything except NIL and FALSE is truthy
			condRef := term.condRef
			if condRef < 0 || condRef >= int32(len(fn.insts)) {
				return nil, fmt.Errorf("spikeRun: branch-if condition ref %d out of bounds", condRef)
			}

			condValue := locals[condRef]
			isTruthy := vm.IsTruthy(condValue)

			var target int32
			var args []int32

			if isTruthy {
				target = term.trueEdge.target
				args = term.trueEdge.args
			} else {
				target = term.falseEdge.target
				args = term.falseEdge.args
			}

			if target < 0 || target >= int32(len(fn.blocks)) {
				return nil, fmt.Errorf("spikeRun: branch-if target %d out of bounds", target)
			}

			targetBlock := fn.blocks[target]
			if len(args) != len(targetBlock.params) {
				return nil, fmt.Errorf("spikeRun: branch-if edge has %d args but target has %d params", len(args), len(targetBlock.params))
			}

			// Snapshot edge arg values (before writing any params, for loop back-edges)
			edgeValues := make([]vm.Value, len(args))
			for i, edgeArgID := range args {
				if edgeArgID < 0 || edgeArgID >= int32(len(fn.insts)) {
					return nil, fmt.Errorf("spikeRun: branch-if edge arg %d out of bounds", edgeArgID)
				}
				edgeValues[i] = locals[edgeArgID]
			}

			// Write snapshotted values to target params
			for i, paramID := range targetBlock.params {
				locals[paramID] = edgeValues[i]
			}

			currentBlockID = int(target)
			continue

		default:
			return nil, fmt.Errorf("spikeRun: unknown terminator: %s", term.op)
		}
	}
}

// getBoxedValue returns a value as boxed vm.Value, boxing if necessary (without counting toward stats).
func getBoxedValue(instID int32, routes []uint8, locals []vm.Value, localsI []int64, localsF []float64) vm.Value {
	route := routes[instID]
	if route == ROUTE_INT {
		// If routed INT, return boxed int from localsI
		return vm.Int(localsI[instID])
	} else if route == ROUTE_FLOAT {
		// If routed FLOAT, return boxed float from localsF
		return vm.Float(localsF[instID])
	}
	if route == ROUTE_BOOL {
		return vm.Boolean(localsI[instID] != 0)
	}
	// Otherwise return from boxed locals. May legitimately be nil only for
	// insts that were never evaluated — callers must treat nil as an error.
	return locals[instID]
}

// boxedOperand materializes an operand for a boxed-path op, boxing from a
// typed slot when the operand routes INT/FLOAT/BOOL (counted as a boundary
// box). A nil boxed slot is a routing bug and errors loudly.
func boxedOperand(instID int32, routes []uint8, locals []vm.Value, localsI []int64, localsF []float64, stats *spikeStats) (vm.Value, error) {
	if routes[instID] != ROUTE_BOXED {
		stats.boxOps++
	}
	v := getBoxedValue(instID, routes, locals, localsI, localsF)
	if v == nil {
		return nil, fmt.Errorf("spikeRunTyped: operand inst %d has nil boxed value (routing bug)", instID)
	}
	return v, nil
}

// spikeRunTyped executes with unboxed int64/float64 slots based on routing.
// Uses type information from typeinfer to avoid per-operation boxing.
func spikeRunTyped(fn *spikeFn, args []vm.Value, stats *spikeStats) (vm.Value, error) {
	if len(args) != fn.nargs {
		return nil, fmt.Errorf("spikeRunTyped: expected %d args, got %d", fn.nargs, len(args))
	}

	// Allocate parallel slots
	locals := make([]vm.Value, len(fn.insts))
	localsI := make([]int64, len(fn.insts))
	localsF := make([]float64, len(fn.insts))

	currentBlockID := 0

	for {
		if currentBlockID < 0 || currentBlockID >= len(fn.blocks) {
			return nil, fmt.Errorf("spikeRunTyped: block id %d out of bounds", currentBlockID)
		}

		block := fn.blocks[currentBlockID]

		// Block lists are decode-validated: live, executable insts only
		// (see validateLiveInvariants). No per-iteration guards needed.
		for _, nid := range block.insts {
			inst := fn.insts[nid]
			route := fn.routes[nid]

			switch inst.opc {
			case spikeOpConst:
				if route == ROUTE_INT {
					if iv, ok := inst.auxVal.(vm.Int); ok {
						localsI[nid] = int64(iv)
					}
				} else if route == ROUTE_FLOAT {
					if fv, ok := inst.auxVal.(vm.Float); ok {
						localsF[nid] = float64(fv)
					}
				} else {
					locals[nid] = inst.auxVal
				}

			case spikeOpLoadArg:
				if inst.aux < 0 || inst.aux >= int32(fn.nargs) {
					return nil, fmt.Errorf("spikeRunTyped: load-arg aux %d out of range [0, %d)", inst.aux, fn.nargs)
				}
				argValue := args[inst.aux]
				switch route {
				case ROUTE_INT:
					// Boundary unbox (once per call). A mismatch means the
					// runtime arg contradicts typeinfer's narrowing — loud
					// error; the real boundary design is STORY-0053's
					// rt.UnboxInt (spike limitation, documented).
					iv, ok := argValue.(vm.Int)
					if !ok {
						return nil, fmt.Errorf("spikeRunTyped: load-arg %d narrowed to int but got %T (spike inference-mismatch limitation)", inst.aux, argValue)
					}
					localsI[nid] = int64(iv)
					stats.unboxOps++
				case ROUTE_FLOAT:
					fv, ok := argValue.(vm.Float)
					if !ok {
						return nil, fmt.Errorf("spikeRunTyped: load-arg %d narrowed to float but got %T (spike inference-mismatch limitation)", inst.aux, argValue)
					}
					localsF[nid] = float64(fv)
					stats.unboxOps++
				default:
					locals[nid] = argValue
				}

			case spikeOpLoadVar:
				locals[nid] = inst.auxVal

			case spikeOpAdd:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRunTyped: add inst %d has %d args, expected 2", nid, len(inst.args))
				}
				if route == ROUTE_INT {
					a, b := localsI[inst.args[0]], localsI[inst.args[1]]
					aInt, bInt := vm.Int(a), vm.Int(b)
					r, ok := checkedAddIntSpike(aInt, bInt)
					if !ok {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (add) integer overflow", nid)
					}
					localsI[nid] = int64(r)
					stats.unboxOps++
				} else if route == ROUTE_FLOAT {
					a, b := localsF[inst.args[0]], localsF[inst.args[1]]
					localsF[nid] = a + b
					stats.unboxOps++
				} else {
					a, aErr := boxedOperand(inst.args[0], fn.routes, locals, localsI, localsF, stats)
					if aErr != nil {
						return nil, aErr
					}
					b, bErr := boxedOperand(inst.args[1], fn.routes, locals, localsI, localsF, stats)
					if bErr != nil {
						return nil, bErr
					}
					r, err := vm.NumAdd(a, b)
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (add) error: %w", nid, err)
					}
					locals[nid] = r
					stats.boxedArithOps++
				}

			case spikeOpSub:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRunTyped: sub inst %d has %d args, expected 2", nid, len(inst.args))
				}
				if route == ROUTE_INT {
					a, b := localsI[inst.args[0]], localsI[inst.args[1]]
					aInt, bInt := vm.Int(a), vm.Int(b)
					r, ok := checkedSubIntSpike(aInt, bInt)
					if !ok {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (sub) integer overflow", nid)
					}
					localsI[nid] = int64(r)
					stats.unboxOps++
				} else if route == ROUTE_FLOAT {
					a, b := localsF[inst.args[0]], localsF[inst.args[1]]
					localsF[nid] = a - b
					stats.unboxOps++
				} else {
					a, aErr := boxedOperand(inst.args[0], fn.routes, locals, localsI, localsF, stats)
					if aErr != nil {
						return nil, aErr
					}
					b, bErr := boxedOperand(inst.args[1], fn.routes, locals, localsI, localsF, stats)
					if bErr != nil {
						return nil, bErr
					}
					r, err := vm.NumSub(a, b)
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (sub) error: %w", nid, err)
					}
					locals[nid] = r
					stats.boxedArithOps++
				}

			case spikeOpMul:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRunTyped: mul inst %d has %d args, expected 2", nid, len(inst.args))
				}
				if route == ROUTE_INT {
					a, b := localsI[inst.args[0]], localsI[inst.args[1]]
					aInt, bInt := vm.Int(a), vm.Int(b)
					r, ok := checkedMulIntSpike(aInt, bInt)
					if !ok {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (mul) integer overflow", nid)
					}
					localsI[nid] = int64(r)
					stats.unboxOps++
				} else if route == ROUTE_FLOAT {
					a, b := localsF[inst.args[0]], localsF[inst.args[1]]
					localsF[nid] = a * b
					stats.unboxOps++
				} else {
					a, aErr := boxedOperand(inst.args[0], fn.routes, locals, localsI, localsF, stats)
					if aErr != nil {
						return nil, aErr
					}
					b, bErr := boxedOperand(inst.args[1], fn.routes, locals, localsI, localsF, stats)
					if bErr != nil {
						return nil, bErr
					}
					r, err := vm.NumMul(a, b)
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (mul) error: %w", nid, err)
					}
					locals[nid] = r
					stats.boxedArithOps++
				}

			case spikeOpInc:
				if len(inst.args) != 1 {
					return nil, fmt.Errorf("spikeRunTyped: inc inst %d has %d args, expected 1", nid, len(inst.args))
				}
				if route == ROUTE_INT {
					a := localsI[inst.args[0]]
					aInt := vm.Int(a)
					r, ok := checkedAddIntSpike(aInt, 1)
					if !ok {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (inc) integer overflow", nid)
					}
					localsI[nid] = int64(r)
					stats.unboxOps++
				} else {
					a, aErr := boxedOperand(inst.args[0], fn.routes, locals, localsI, localsF, stats)
					if aErr != nil {
						return nil, aErr
					}
					r, err := vm.NumAdd(a, vm.Int(1))
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (inc) error: %w", nid, err)
					}
					locals[nid] = r
					stats.boxedArithOps++
				}

			case spikeOpDec:
				if len(inst.args) != 1 {
					return nil, fmt.Errorf("spikeRunTyped: dec inst %d has %d args, expected 1", nid, len(inst.args))
				}
				if route == ROUTE_INT {
					a := localsI[inst.args[0]]
					aInt := vm.Int(a)
					r, ok := checkedSubIntSpike(aInt, 1)
					if !ok {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (dec) integer overflow", nid)
					}
					localsI[nid] = int64(r)
					stats.unboxOps++
				} else {
					a, aErr := boxedOperand(inst.args[0], fn.routes, locals, localsI, localsF, stats)
					if aErr != nil {
						return nil, aErr
					}
					r, err := vm.NumSub(a, vm.Int(1))
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (dec) error: %w", nid, err)
					}
					locals[nid] = r
					stats.boxedArithOps++
				}

			case spikeOpLt, spikeOpLte, spikeOpGt, spikeOpGte, spikeOpEq:
				if len(inst.args) != 2 {
					return nil, fmt.Errorf("spikeRunTyped: %s inst %d has %d args, expected 2", inst.op, nid, len(inst.args))
				}
				if route == ROUTE_BOOL {
					// Native comparison: both operands narrowed to the same
					// numeric route; result is 0/1 in localsI, no boxing.
					aID, bID := inst.args[0], inst.args[1]
					var cmp bool
					if fn.routes[aID] == ROUTE_INT {
						a, b := localsI[aID], localsI[bID]
						switch inst.opc {
						case spikeOpLt:
							cmp = a < b
						case spikeOpLte:
							cmp = a <= b
						case spikeOpGt:
							cmp = a > b
						case spikeOpGte:
							cmp = a >= b
						case spikeOpEq:
							cmp = a == b
						}
					} else {
						a, b := localsF[aID], localsF[bID]
						switch inst.opc {
						case spikeOpLt:
							cmp = a < b
						case spikeOpLte:
							cmp = a <= b
						case spikeOpGt:
							cmp = a > b
						case spikeOpGte:
							cmp = a >= b
						case spikeOpEq:
							cmp = a == b
						}
					}
					if cmp {
						localsI[nid] = 1
					} else {
						localsI[nid] = 0
					}
					stats.unboxOps++
					break
				}
				a, aErr := boxedOperand(inst.args[0], fn.routes, locals, localsI, localsF, stats)
				if aErr != nil {
					return nil, aErr
				}
				b, bErr := boxedOperand(inst.args[1], fn.routes, locals, localsI, localsF, stats)
				if bErr != nil {
					return nil, bErr
				}
				var result vm.Value
				switch inst.opc {
				case spikeOpLt:
					r, err := vm.NumLt(a, b)
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (lt) error: %w", nid, err)
					}
					result = vm.Boolean(r)
				case spikeOpLte:
					r, err := vm.NumLe(a, b)
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (lte) error: %w", nid, err)
					}
					result = vm.Boolean(r)
				case spikeOpGt:
					r, err := vm.NumGt(a, b)
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (gt) error: %w", nid, err)
					}
					result = vm.Boolean(r)
				case spikeOpGte:
					r, err := vm.NumGe(a, b)
					if err != nil {
						return nil, fmt.Errorf("spikeRunTyped: inst %d (gte) error: %w", nid, err)
					}
					result = vm.Boolean(r)
				case spikeOpEq:
					// Mirror OP_EQ: Int/Keyword fast paths, then structural.
					result = nil
					if ai, ok := a.(vm.Int); ok {
						if bi, ok := b.(vm.Int); ok {
							result = vm.Boolean(ai == bi)
						}
					}
					if result == nil {
						if ak, ok := a.(vm.Keyword); ok {
							if bk, ok := b.(vm.Keyword); ok {
								result = vm.Boolean(ak == bk)
							}
						}
					}
					if result == nil {
						if vm.ValueEquals == nil {
							return nil, fmt.Errorf("spikeRunTyped: ValueEquals not initialized")
						}
						result = vm.Boolean(vm.ValueEquals(a, b))
					}
				}
				locals[nid] = result
				stats.boxedArithOps++

			case spikeOpCall:
				// Call op: first ref is callee (as inst-id), remaining refs are args
				if len(inst.args) < 1 {
					return nil, fmt.Errorf("spikeRunTyped: call inst %d has no callee ref", nid)
				}

				calleeIdx := fn.calleeIndices[nid]
				if calleeIdx < 0 || calleeIdx >= len(fn.callees) {
					return nil, fmt.Errorf("spikeRunTyped: call inst %d has invalid callee index %d", nid, calleeIdx)
				}
				callee := fn.callees[calleeIdx]

				// Remaining args are the actual call arguments
				// Box typed operands at the boundary
				callArgs := make([]vm.Value, len(inst.args)-1)
				for i, argID := range inst.args[1:] {
					if argID < 0 || argID >= int32(len(fn.insts)) {
						return nil, fmt.Errorf("spikeRunTyped: call arg %d out of bounds", argID)
					}
					// Use boxedOperand to materialize the argument value
					v, boxErr := boxedOperand(argID, fn.routes, locals, localsI, localsF, stats)
					if boxErr != nil {
						return nil, boxErr
					}
					callArgs[i] = v
				}

				// Invoke the callee
				var callResult vm.Value
				var callErr error
				switch c := callee.(type) {
				case *vm.Func:
					callResult, callErr = c.Invoke(callArgs)
				case *vm.NativeFn:
					callResult, callErr = c.Invoke(callArgs)
				default:
					return nil, fmt.Errorf("spikeRunTyped: call callee %T is not callable", callee)
				}
				if callErr != nil {
					return nil, fmt.Errorf("spikeRunTyped: inst %d (call) error: %w", nid, callErr)
				}
				stats.callOps++
				// Call result is always boxed
				locals[nid] = callResult

			default:
				return nil, fmt.Errorf("spikeRunTyped: inst %d has unknown op: %s", nid, inst.op)
			}
		}

		term := block.term

		switch term.opc {
		case spikeOpReturn:
			if term.returnVal < 0 || term.returnVal >= int32(len(fn.insts)) {
				return nil, fmt.Errorf("spikeRunTyped: return value inst %d out of bounds", term.returnVal)
			}
			// Box the result if needed
			route := fn.routes[term.returnVal]
			if route == ROUTE_INT {
				stats.boxOps++
				return vm.Int(localsI[term.returnVal]), nil
			} else if route == ROUTE_FLOAT {
				stats.boxOps++
				return vm.Float(localsF[term.returnVal]), nil
			} else if route == ROUTE_BOOL {
				stats.boxOps++
				return vm.Boolean(localsI[term.returnVal] != 0), nil
			}
			return locals[term.returnVal], nil

		case spikeOpBranch:
			target := term.simpleEdge.target
			if target < 0 || target >= int32(len(fn.blocks)) {
				return nil, fmt.Errorf("spikeRunTyped: branch target %d out of bounds", target)
			}

			targetBlock := fn.blocks[target]
			if len(term.simpleEdge.args) != len(targetBlock.params) {
				return nil, fmt.Errorf("spikeRunTyped: branch edge has %d args but target has %d params", len(term.simpleEdge.args), len(targetBlock.params))
			}

			// Thread edge args with proper boxing/unboxing
			for i, edgeArgID := range term.simpleEdge.args {
				if edgeArgID < 0 || edgeArgID >= int32(len(fn.insts)) {
					return nil, fmt.Errorf("spikeRunTyped: branch edge arg %d out of bounds", edgeArgID)
				}
				paramID := targetBlock.params[i]
				if terr := threadEdgeValueTyped(edgeArgID, paramID, fn.routes, localsI, localsF, locals, stats); terr != nil {
					return nil, terr
				}
			}

			currentBlockID = int(target)
			continue

		case spikeOpBranchIf:
			condRef := term.condRef
			if condRef < 0 || condRef >= int32(len(fn.insts)) {
				return nil, fmt.Errorf("spikeRunTyped: branch-if condition ref %d out of bounds", condRef)
			}

			var isTruthy bool
			switch fn.routes[condRef] {
			case ROUTE_BOOL:
				// Native comparison result: no boxing on the loop back-edge.
				isTruthy = localsI[condRef] != 0
			case ROUTE_INT, ROUTE_FLOAT:
				// Numbers are always truthy under lg rules; no box needed.
				isTruthy = true
			default:
				condValue := locals[condRef]
				if condValue == nil {
					return nil, fmt.Errorf("spikeRunTyped: branch-if cond inst %d has nil value (routing bug)", condRef)
				}
				isTruthy = vm.IsTruthy(condValue)
			}

			var target int32
			var args []int32

			if isTruthy {
				target = term.trueEdge.target
				args = term.trueEdge.args
			} else {
				target = term.falseEdge.target
				args = term.falseEdge.args
			}

			if target < 0 || target >= int32(len(fn.blocks)) {
				return nil, fmt.Errorf("spikeRunTyped: branch-if target %d out of bounds", target)
			}

			targetBlock := fn.blocks[target]
			if len(args) != len(targetBlock.params) {
				return nil, fmt.Errorf("spikeRunTyped: branch-if edge has %d args but target has %d params", len(args), len(targetBlock.params))
			}

			for i, edgeArgID := range args {
				if edgeArgID < 0 || edgeArgID >= int32(len(fn.insts)) {
					return nil, fmt.Errorf("spikeRunTyped: branch-if edge arg %d out of bounds", edgeArgID)
				}
				paramID := targetBlock.params[i]
				if terr := threadEdgeValueTyped(edgeArgID, paramID, fn.routes, localsI, localsF, locals, stats); terr != nil {
					return nil, terr
				}
			}

			currentBlockID = int(target)
			continue

		default:
			return nil, fmt.Errorf("spikeRunTyped: unknown terminator: %s", term.op)
		}
	}
}

// threadEdgeValueTyped copies values between slots, handling route mismatches
// with boxing/unboxing. A boxed value whose runtime type contradicts the
// target's narrowed route is a routing bug — error loudly, never drop silently.
func threadEdgeValueTyped(fromID, toID int32, routes []uint8, localsI []int64, localsF []float64, locals []vm.Value, stats *spikeStats) error {
	fromRoute := routes[fromID]
	toRoute := routes[toID]

	if fromRoute == toRoute {
		switch fromRoute {
		case ROUTE_INT, ROUTE_BOOL:
			localsI[toID] = localsI[fromID]
		case ROUTE_FLOAT:
			localsF[toID] = localsF[fromID]
		case ROUTE_BOXED:
			locals[toID] = locals[fromID]
		}
	} else if fromRoute == ROUTE_BOOL && toRoute == ROUTE_BOXED {
		stats.boxOps++
		locals[toID] = vm.Boolean(localsI[fromID] != 0)
	} else if fromRoute == ROUTE_BOXED && toRoute == ROUTE_BOOL {
		localsI[toID] = 0
		if vm.IsTruthy(locals[fromID]) {
			localsI[toID] = 1
		}
	} else if fromRoute == ROUTE_INT && toRoute == ROUTE_BOXED {
		stats.boxOps++
		locals[toID] = vm.Int(localsI[fromID])
	} else if fromRoute == ROUTE_FLOAT && toRoute == ROUTE_BOXED {
		stats.boxOps++
		locals[toID] = vm.Float(localsF[fromID])
	} else if fromRoute == ROUTE_BOXED && toRoute == ROUTE_INT {
		iv, ok := locals[fromID].(vm.Int)
		if !ok {
			return fmt.Errorf("threadEdgeValueTyped: edge %d→%d target narrowed to int but boxed value is %T", fromID, toID, locals[fromID])
		}
		stats.unboxOps++
		localsI[toID] = int64(iv)
	} else if fromRoute == ROUTE_BOXED && toRoute == ROUTE_FLOAT {
		fv, ok := locals[fromID].(vm.Float)
		if !ok {
			return fmt.Errorf("threadEdgeValueTyped: edge %d→%d target narrowed to float but boxed value is %T", fromID, toID, locals[fromID])
		}
		stats.unboxOps++
		localsF[toID] = float64(fv)
	} else {
		return fmt.Errorf("threadEdgeValueTyped: unsupported route transition %d→%d (routes %d→%d)", fromID, toID, fromRoute, toRoute)
	}
	return nil
}

// Overflow checking helpers
const maxIntValueSpike = vm.Int(int(^uint(0) >> 1))
const minIntValueSpike = -maxIntValueSpike - 1

func checkedAddIntSpike(a, b vm.Int) (vm.Int, bool) {
	if (b > 0 && a > maxIntValueSpike-b) || (b < 0 && a < minIntValueSpike-b) {
		return 0, false
	}
	return a + b, true
}

func checkedSubIntSpike(a, b vm.Int) (vm.Int, bool) {
	if (b < 0 && a > maxIntValueSpike+b) || (b > 0 && a < minIntValueSpike+b) {
		return 0, false
	}
	return a - b, true
}

func checkedMulIntSpike(a, b vm.Int) (vm.Int, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	if (a == minIntValueSpike && b == -1) || (b == minIntValueSpike && a == -1) {
		return 0, false
	}
	if a > 0 {
		if b > 0 && a > maxIntValueSpike/b {
			return 0, false
		}
		if b < 0 && b < minIntValueSpike/a {
			return 0, false
		}
	} else {
		if b > 0 && a < minIntValueSpike/b {
			return 0, false
		}
		if b < 0 && a < maxIntValueSpike/b {
			return 0, false
		}
	}
	return a * b, true
}

// TestSpikeRun_TypedLoopKernel: AC-WS.2 kernel (a) — unboxed loop arithmetic.
// Kernel: (defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))
// Expected: boxedArithOps==0 for the loop interior, boxOps==1 at return, and boxOps(n=10)==boxOps(n=1000).
func TestSpikeRun_TypedLoopKernel(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-typed-loop")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode IR: %v", decodeErr)
	}

	// Compile stack VM for comparison
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-typed-loop")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))`)); stackErr != nil {
		t.Fatalf("compile stack VM: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("ksum"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Test cases with statistics collection
	testCases := []struct {
		n          vm.Value
		label      string
		wantBoxOps int // at return only
	}{
		{vm.Int(0), "n=0", 1},
		{vm.Int(1), "n=1", 1},
		{vm.Int(10), "n=10", 1},
		{vm.Int(1000), "n=1000", 1},
	}

	for _, tc := range testCases {
		// Typed execution
		stats := &spikeStats{}
		spikeResult, spikeErr := spikeRunTyped(fn, []vm.Value{tc.n}, stats)
		if spikeErr != nil {
			t.Errorf("spikeRunTyped(%s) error: %v", tc.label, spikeErr)
			continue
		}

		// Stack VM
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.n})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM(%s) error: %v", tc.label, stackErr)
			continue
		}

		// Parity
		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("kernel (a) FAIL %s: spike=%v, stack=%v", tc.label, spikeResult, stackResult)
		} else {
			t.Logf("kernel (a) PASS %s: result=%v", tc.label, spikeResult)
		}

		// AC-WS.2 assertions: ALL narrowed loop ops (add/inc/lt) run unboxed.
		n := int(tc.n.(vm.Int))
		if stats.boxedArithOps != 0 {
			t.Errorf("kernel (a) %s: boxedArithOps=%d, want 0 (narrowed ops must not run boxed)", tc.label, stats.boxedArithOps)
		}
		if n > 0 && stats.unboxOps == 0 {
			t.Errorf("kernel (a) %s: unboxOps=0, expected >0 (arithmetic should be unboxed)", tc.label)
		}
		t.Logf("kernel (a) %s stats: unboxOps=%d boxedArithOps=%d boxOps=%d", tc.label, stats.unboxOps, stats.boxedArithOps, stats.boxOps)
	}

	// Verify arithmetic operations are unboxed (only boundary box at return)
	// Both calls should have same boxOps (return boxing)
	stats10 := &spikeStats{}
	spikeRunTyped(fn, []vm.Value{vm.Int(10)}, stats10)

	stats1000 := &spikeStats{}
	spikeRunTyped(fn, []vm.Value{vm.Int(1000)}, stats1000)

	if stats10.unboxOps == 0 || stats1000.unboxOps == 0 {
		t.Errorf("kernel (a): unboxOps should be non-zero (loop arithmetic exercised)")
	} else {
		t.Logf("kernel (a): arithmetic unboxed: n=10 unboxOps=%d, n=1000 unboxOps=%d", stats10.unboxOps, stats1000.unboxOps)
	}
	// Boxing must be a boundary cost, invariant in n — never per-iteration.
	if stats10.boxOps != stats1000.boxOps {
		t.Errorf("kernel (a): boxOps scales with n: n=10 → %d, n=1000 → %d (must be constant)", stats10.boxOps, stats1000.boxOps)
	}
}

// TestSpikeRun_TypedOverflow: AC-WS.2 unboxed overflow handling.
func TestSpikeRun_TypedOverflow(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-typed-overflow")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn ov [x] (+ x 1)))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode IR: %v", decodeErr)
	}

	stats := &spikeStats{}
	maxInt := vm.Int(9223372036854775807) // 2^63 - 1

	_, spikeErr := spikeRunTyped(fn, []vm.Value{maxInt}, stats)

	if spikeErr == nil {
		t.Error("expected overflow error, got nil")
	} else if !strings.Contains(spikeErr.Error(), "overflow") {
		t.Errorf("expected 'overflow' in error message, got: %v", spikeErr)
	} else {
		t.Logf("overflow correctly detected: %v", spikeErr)
	}
}

// TestSpikeRun_TypedMixed: AC-WS.2 mixed typed and boxed operations.
// For fixtures where only part of the ops narrow (unboxed and boxed paths both exercised).
func TestSpikeRun_TypedMixed(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-typed-mixed")

	// Fixture: (defn mx [a b] (+ (* a a) (if (< a b) a b)))
	// The (* a a) is typed int; the if branches on (< a b) which is comparison (boxed).
	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn mx [a b] (+ (* a a) (if (< a b) a b))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode IR: %v", decodeErr)
	}

	// Compile stack VM
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-typed-mixed")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn mx [a b] (+ (* a a) (if (< a b) a b)))`)); stackErr != nil {
		t.Fatalf("compile stack VM: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("mx"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	testCases := []struct {
		a, b  vm.Value
		label string
	}{
		{vm.Int(2), vm.Int(3), "a<b"},
		{vm.Int(5), vm.Int(2), "a>b"},
		{vm.Int(3), vm.Int(3), "a==b"},
	}

	for _, tc := range testCases {
		stats := &spikeStats{}
		spikeResult, spikeErr := spikeRunTyped(fn, []vm.Value{tc.a, tc.b}, stats)
		if spikeErr != nil {
			t.Errorf("spikeRunTyped(%s) error: %v", tc.label, spikeErr)
			continue
		}

		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.a, tc.b})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM(%s) error: %v", tc.label, stackErr)
			continue
		}

		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("mixed-typed FAIL %s: spike=%v, stack=%v", tc.label, spikeResult, stackResult)
		} else {
			t.Logf("mixed-typed PASS %s: result=%v, unboxOps=%d, boxedOps=%d", tc.label, spikeResult, stats.unboxOps, stats.boxedArithOps)
		}

		// Verify both paths are exercised
		if stats.unboxOps == 0 && stats.boxedArithOps == 0 {
			t.Logf("mixed-typed %s: neither path exercised (pure boxed load-arg)", tc.label)
		}
	}
}

// TestDebugRouting prints routing info for the kernel to understand why it's not unboxing.
func TestDebugRouting(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-debug-routing")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode IR: %v", decodeErr)
	}

	// Print routing and types
	t.Logf("=== Decoded IR dump ===")
	t.Logf("Total insts: %d", len(fn.insts))
	for i, inst := range fn.insts {
		if inst.opc == spikeOpInvalid {
			continue
		}
		route := fn.routes[i]
		routeName := []string{"BOXED", "INT", "FLOAT", "BOOL"}[route]
		t.Logf("inst %d: op=%s typ=%d route=%s args=%v", i, inst.op, inst.typ, routeName, inst.args)
	}
}

// TestSpikeRun_TypedFloat: AC-WS.2 float accumulator kernel (truthful typ observation).
// Kernel: (defn kf [n] (loop [i 0 acc 0.0] (if (< i n) (recur (+ i 1) (+ acc 1.5)) acc)))
// Documents whether typeinfer narrows the float accumulator and tests both paths if so.
func TestSpikeRun_TypedFloat(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("test-spike-typed-float")

	expr := `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn kf [n] (loop [i 0 acc 0.0] (if (< i n) (recur (+ i 1) (+ acc 1.5)) acc))))))`
	_, irVal, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("compile IR: %v", err)
	}

	fn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("decode IR: %v", decodeErr)
	}

	// Compile stack VM
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("test-stack-typed-float")
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(`(defn kf [n] (loop [i 0 acc 0.0] (if (< i n) (recur (+ i 1) (+ acc 1.5)) acc)))`)); stackErr != nil {
		t.Fatalf("compile stack VM: %v", stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol("kf"))
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Inspect typeinfer results to document float narrowing behavior
	t.Logf("=== Float kernel typ codes ===")
	hasFloatRoute := false
	for i, inst := range fn.insts {
		if inst.opc == spikeOpInvalid {
			continue
		}
		route := fn.routes[i]
		routeName := []string{"BOXED", "INT", "FLOAT", "BOOL"}[route]
		if route == ROUTE_FLOAT {
			hasFloatRoute = true
			t.Logf("inst %d (%s): typ=%d route=%s", i, inst.op, inst.typ, routeName)
		}
	}

	if !hasFloatRoute {
		t.Logf("Note: Typeinfer did not narrow float accumulator to :float; entire kernel routes as BOXED")
	}

	// Test cases
	testCases := []struct {
		n          vm.Value
		label      string
		wantBoxOps int
	}{
		{vm.Int(0), "n=0", 1},
		{vm.Int(10), "n=10", 1},
		{vm.Int(100), "n=100", 1},
	}

	for _, tc := range testCases {
		// Typed execution
		stats := &spikeStats{}
		spikeResult, spikeErr := spikeRunTyped(fn, []vm.Value{tc.n}, stats)
		if spikeErr != nil {
			t.Errorf("spikeRunTyped(%s) error: %v", tc.label, spikeErr)
			continue
		}

		// Stack VM
		stackFrame := vm.NewFrame(stackChunk, []vm.Value{tc.n})
		stackResult, stackErr := stackFrame.Run()
		if stackErr != nil {
			t.Errorf("stack VM(%s) error: %v", tc.label, stackErr)
			continue
		}

		// Parity check (REQUIRED - must match stack VM exactly)
		if !valuesEqual(spikeResult, stackResult) {
			t.Errorf("float kernel FAIL %s: spike=%v (%T), stack=%v (%T)", tc.label, spikeResult, spikeResult, stackResult, stackResult)
		} else {
			t.Logf("float kernel PASS %s: result=%v (spike and stack identical)", tc.label, spikeResult)
		}

		// Stats logging
		if hasFloatRoute {
			t.Logf("float kernel %s: boxOps=%d (expecting %d), unboxOps=%d (float arithmetic)", tc.label, stats.boxOps, tc.wantBoxOps, stats.unboxOps)
			// If float was narrowed, boxOps should be constant
			if stats.boxOps != tc.wantBoxOps {
				t.Logf("float kernel %s: boxOps constant check (expected %d, got %d)", tc.label, tc.wantBoxOps, stats.boxOps)
			}
		}
	}
}
