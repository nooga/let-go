/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

// ComputeUses populates Function.Uses[v] = list of NodeIDs that
// reference v as an operand or as a branch-target arg. This is the
// def→use direction; the Node.Refs field is use→def.
//
// Run after Build; rerun whenever a pass mutates the IR.
//
// Uses-list is the foundation for:
//   - DCE: pure nodes with empty Uses (and not a terminator's operand) are dead
//   - CSE: when merging two equivalent values, redirect all uses of one to the other
//   - LICM: a loop-invariant value can be hoisted if none of its uses are loop-defs
//   - inlining decisions: a single-use call site is a strong candidate
type Uses [][]NodeID // indexed by NodeID; Uses[v] = nodes that reference v

// ComputeUses returns the def→use index for f.
func ComputeUses(f *Function) Uses {
	uses := make(Uses, len(f.Nodes))
	for id, n := range f.Nodes {
		_ = id
		// Direct operand refs.
		for _, r := range n.Refs {
			uses[r] = append(uses[r], NodeID(id))
		}
		// Branch-target args (live in Aux for terminators).
		switch t := n.Aux.(type) {
		case *BranchTarget:
			if t != nil {
				for _, a := range t.Args {
					uses[a] = append(uses[a], NodeID(id))
				}
			}
		case *CondTarget:
			if t != nil {
				if t.True != nil {
					for _, a := range t.True.Args {
						uses[a] = append(uses[a], NodeID(id))
					}
				}
				if t.False != nil {
					for _, a := range t.False.Args {
						uses[a] = append(uses[a], NodeID(id))
					}
				}
			}
		}
	}
	return uses
}
