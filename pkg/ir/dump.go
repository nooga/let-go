/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

import (
	"fmt"
	"strings"
)

// Dump returns a human-readable text representation of f.
//
// Format (loosely Carbon/SIL-flavoured):
//
//	fn FibFast(arity=1, variadic=false):
//	  entry b0(v0: ?):           ; b0 takes one BlockArg, v0 (the function's arg)
//	    v1 = LoadVar #'<=
//	    v2 = Const 1
//	    v3 = Call v1, v0, v2
//	    BranchIf v3 -> b2() : b1()
//
//	  b1():
//	    v4 = ...
//	    Return v4
//
// Compact, scannable, debuggable.
func Dump(f *Function) string {
	var b strings.Builder
	fmt.Fprintf(&b, "fn %s(arity=%d, variadic=%v):\n",
		f.Name, f.Arity, f.IsVariadic)
	for _, blk := range f.Blocks {
		writeBlock(&b, f, &blk)
	}
	return b.String()
}

func writeBlock(b *strings.Builder, f *Function, blk *Block) {
	tag := ""
	if blk.ID == f.Entry {
		tag = "entry "
	}
	// Block header with parameters.
	fmt.Fprintf(b, "  %sb%d(", tag, blk.ID)
	for i, p := range blk.Params {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(b, "v%d: %s", p, f.Node(p).Type)
	}
	b.WriteString("):")
	if len(blk.Preds) > 0 {
		b.WriteString("    ; preds: ")
		for i, p := range blk.Preds {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(b, "b%d", p)
		}
	}
	b.WriteString("\n")

	// Body nodes.
	for _, id := range blk.Nodes {
		n := f.Node(id)
		writeNode(b, f, id, n)
	}

	// Terminator.
	if blk.Term != 0 || (len(blk.Nodes) > 0) {
		n := f.Node(blk.Term)
		writeNode(b, f, blk.Term, n)
	}
	b.WriteString("\n")
}

func writeNode(b *strings.Builder, f *Function, id NodeID, n *Node) {
	if n.Op.IsTerminator() {
		fmt.Fprintf(b, "    %s", n.Op)
		// Refs (e.g. condition for BranchIf, value for Return).
		for _, r := range n.Refs {
			fmt.Fprintf(b, " v%d", r)
		}
		// Targets.
		switch t := n.Aux.(type) {
		case *BranchTarget:
			fmt.Fprintf(b, " -> b%d(", t.Target)
			for i, a := range t.Args {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(b, "v%d", a)
			}
			b.WriteString(")")
		case *CondTarget:
			fmt.Fprintf(b, " -> b%d(", t.True.Target)
			for i, a := range t.True.Args {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(b, "v%d", a)
			}
			fmt.Fprintf(b, ") : b%d(", t.False.Target)
			for i, a := range t.False.Args {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(b, "v%d", a)
			}
			b.WriteString(")")
		}
		b.WriteString("\n")
		return
	}
	// Regular value.
	fmt.Fprintf(b, "    v%d = %s", id, n.Op)
	for _, r := range n.Refs {
		fmt.Fprintf(b, " v%d", r)
	}
	if n.Aux != nil {
		fmt.Fprintf(b, " ; %v", n.Aux)
	}
	if n.Type != TypeUnknown {
		fmt.Fprintf(b, " : %s", n.Type)
	}
	b.WriteString("\n")
}
