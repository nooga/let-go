/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Package passes — IR optimization passes.
package passes

import (
	"github.com/nooga/let-go/pkg/ir"
	"github.com/nooga/let-go/pkg/vm"
)

// ConstFold finds pure arithmetic ops whose operands are all constants
// and replaces them with the folded constant. Repeats to fixed point.
//
// Operates in-place. Doesn't remove the old nodes (DCE does that);
// just rewires uses to the new const node.
func ConstFold(f *ir.Function) (changed bool) {
	for changed = false; ; {
		anyChange := false
		for id, n := range f.Nodes {
			if !n.Op.IsPure() {
				continue
			}
			if n.Op == ir.OpConst {
				continue
			}
			folded, ok := tryFold(f, &n)
			if !ok {
				continue
			}
			// Collect operand spans before clearing Refs.
			var operandSpans []vm.SourceInfo
			for _, ref := range n.Refs {
				operandSpans = ir.MergeSourceInfo(operandSpans, f.Nodes[ref].SourceInfos...)
			}
			// Apply the fold.
			f.Nodes[id].Op = ir.OpConst
			f.Nodes[id].Refs = nil
			f.Nodes[id].Aux = folded
			// Union operand spans into the rewritten node's existing spans.
			f.Nodes[id].SourceInfos = ir.MergeSourceInfo(f.Nodes[id].SourceInfos, operandSpans...)
			anyChange = true
			changed = true
		}
		if !anyChange {
			break
		}
	}
	return changed
}

func tryFold(f *ir.Function, n *ir.Node) (vm.Value, bool) {
	switch n.Op {
	case ir.OpAdd, ir.OpSub, ir.OpMul,
		ir.OpLt, ir.OpLte, ir.OpGt, ir.OpGte, ir.OpEq:
		if len(n.Refs) != 2 {
			return nil, false
		}
		a := f.Node(n.Refs[0])
		b := f.Node(n.Refs[1])
		if a.Op != ir.OpConst || b.Op != ir.OpConst {
			return nil, false
		}
		av, aok := a.Aux.(vm.Value)
		bv, bok := b.Aux.(vm.Value)
		if !aok || !bok {
			return nil, false
		}
		return foldBinary(n.Op, av, bv)
	case ir.OpInc, ir.OpDec:
		if len(n.Refs) != 1 {
			return nil, false
		}
		a := f.Node(n.Refs[0])
		if a.Op != ir.OpConst {
			return nil, false
		}
		av, ok := a.Aux.(vm.Value)
		if !ok {
			return nil, false
		}
		return foldUnary(n.Op, av)
	}
	return nil, false
}

func foldBinary(op ir.Op, a, b vm.Value) (vm.Value, bool) {
	ai, aOk := a.(vm.Int)
	bi, bOk := b.(vm.Int)
	if !aOk || !bOk {
		return nil, false
	}
	switch op {
	case ir.OpAdd:
		return vm.Int(int(ai) + int(bi)), true
	case ir.OpSub:
		return vm.Int(int(ai) - int(bi)), true
	case ir.OpMul:
		return vm.Int(int(ai) * int(bi)), true
	case ir.OpLt:
		return vm.Boolean(int(ai) < int(bi)), true
	case ir.OpLte:
		return vm.Boolean(int(ai) <= int(bi)), true
	case ir.OpGt:
		return vm.Boolean(int(ai) > int(bi)), true
	case ir.OpGte:
		return vm.Boolean(int(ai) >= int(bi)), true
	case ir.OpEq:
		return vm.Boolean(int(ai) == int(bi)), true
	}
	return nil, false
}

func foldUnary(op ir.Op, a vm.Value) (vm.Value, bool) {
	ai, ok := a.(vm.Int)
	if !ok {
		return nil, false
	}
	switch op {
	case ir.OpInc:
		return vm.Int(int(ai) + 1), true
	case ir.OpDec:
		return vm.Int(int(ai) - 1), true
	}
	return nil, false
}
