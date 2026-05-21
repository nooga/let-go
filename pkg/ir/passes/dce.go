/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package passes

import "github.com/nooga/let-go/pkg/ir"

// DCE removes pure nodes whose results are never used.
//
// Algorithm:
//  1. Compute use-list (ComputeUses).
//  2. Find pure nodes with empty use-list.
//  3. Remove them from their block's body list and mark them dead.
//     (We don't physically remove from f.Nodes — that would reshuffle
//     indices. Instead we set Op to OpInvalid and skip them at lowering.)
//  4. Repeat to fixed point (removing a node may make its operands dead).
func DCE(f *ir.Function) (changed bool) {
	for {
		uses := ir.ComputeUses(f)
		anyChange := false
		for bid := range f.Blocks {
			blk := &f.Blocks[bid]
			kept := blk.Nodes[:0]
			for _, nid := range blk.Nodes {
				n := f.Node(nid)
				// Keep if used or impure (terminators don't appear in blk.Nodes
				// so we only worry about pure body values).
				if !n.Op.IsPure() {
					kept = append(kept, nid)
					continue
				}
				if len(uses[nid]) > 0 {
					kept = append(kept, nid)
					continue
				}
				// Dead. Mark and drop.
				f.Nodes[nid].Op = ir.OpInvalid
				anyChange = true
				changed = true
			}
			blk.Nodes = kept
		}
		if !anyChange {
			break
		}
	}
	return changed
}
