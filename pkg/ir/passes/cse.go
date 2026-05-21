/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package passes

import (
	"fmt"
	"strings"

	"github.com/nooga/let-go/pkg/ir"
	"github.com/nooga/let-go/pkg/vm"
)

// CSE — common subexpression elimination.
//
// Block-local for the spike: walks each block's body and merges
// duplicate pure expressions. A node is a duplicate of an earlier
// node if both have:
//   - The same Op
//   - The same Refs (after redirecting via previously-CSE'd nodes)
//   - The same Aux (for Const, that's the vm.Value equality)
//
// When a duplicate is detected:
//   - The duplicate node is marked dead (Op = OpInvalid)
//   - A redirect entry maps the dead node's ID to the original's ID
//
// After the per-block walk, every Refs[] in the function is rewritten
// to follow redirects. This canonicalizes the IR.
//
// Cross-block CSE would need dominance — defer to a later pass.
func CSE(f *ir.Function) (changed bool) {
	redirect := make(map[ir.NodeID]ir.NodeID, len(f.Nodes))

	for bid := range f.Blocks {
		blk := &f.Blocks[bid]

		// key → first NodeID that computed this value (in this block).
		// Key is a string built from Op + redirected Refs + Aux.
		seen := map[string]ir.NodeID{}

		kept := blk.Nodes[:0]
		for _, nid := range blk.Nodes {
			n := f.Node(nid)
			if !n.Op.IsPure() {
				kept = append(kept, nid)
				continue
			}
			// Build the canonical key. Follow redirects on operand refs
			// so two nodes with operands that became equivalent compare equal.
			key := nodeKey(n, redirect)
			if prev, ok := seen[key]; ok {
				// Duplicate. Inherit its spans into the survivor before marking dead.
				f.Nodes[prev].SourceInfos = ir.MergeSourceInfo(f.Nodes[prev].SourceInfos, f.Nodes[nid].SourceInfos...)
				redirect[nid] = prev
				f.Nodes[nid].Op = ir.OpInvalid
				changed = true
				continue
			}
			seen[key] = nid
			kept = append(kept, nid)
		}
		blk.Nodes = kept
	}

	if !changed {
		return false
	}

	// Final pass: rewrite all Refs to follow redirects. Includes
	// terminator Refs and branch-target Args.
	for i := range f.Nodes {
		n := &f.Nodes[i]
		for j, r := range n.Refs {
			if final, ok := follow(redirect, r); ok {
				n.Refs[j] = final
			}
		}
		switch t := n.Aux.(type) {
		case *ir.BranchTarget:
			if t != nil {
				rewriteArgs(t.Args, redirect)
			}
		case *ir.CondTarget:
			if t != nil {
				if t.True != nil {
					rewriteArgs(t.True.Args, redirect)
				}
				if t.False != nil {
					rewriteArgs(t.False.Args, redirect)
				}
			}
		}
	}
	return true
}

// nodeKey computes a canonical key for n. Two nodes with equal keys
// compute the same value. Refs are normalized via the redirect table.
func nodeKey(n *ir.Node, redirect map[ir.NodeID]ir.NodeID) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d|", n.Op)
	for _, r := range n.Refs {
		final, _ := follow(redirect, r)
		fmt.Fprintf(&sb, "%d,", final)
	}
	sb.WriteByte('|')
	// Aux serialization. For Const, hash by the vm.Value's string repr;
	// for LoadArg, by arg index; for LoadVar, by var identity.
	switch v := n.Aux.(type) {
	case vm.Value:
		// Two Const nodes are equal iff their boxed values' String() match.
		// For Int/Float/Bool/String/Keyword this is reliable; for compound
		// values, may over-collapse if values stringify identically but
		// aren't `=`. Conservative refinement (compare actual values via
		// vm.ValueEquals) deferred.
		sb.WriteString(v.String())
	case int:
		fmt.Fprintf(&sb, "%d", v)
	default:
		// For terminators and other ops, Aux often references mutable
		// state (BranchTarget pointers). Don't hash it — terminators
		// aren't CSE candidates anyway (impure).
		fmt.Fprintf(&sb, "%p", v)
	}
	return sb.String()
}

// follow walks the redirect chain to its terminus. Returns the final
// NodeID and whether any redirect was followed.
func follow(redirect map[ir.NodeID]ir.NodeID, id ir.NodeID) (ir.NodeID, bool) {
	final := id
	moved := false
	for {
		next, ok := redirect[final]
		if !ok {
			return final, moved
		}
		final = next
		moved = true
	}
}

func rewriteArgs(args []ir.NodeID, redirect map[ir.NodeID]ir.NodeID) {
	for i, a := range args {
		if final, ok := follow(redirect, a); ok {
			args[i] = final
		}
	}
}
