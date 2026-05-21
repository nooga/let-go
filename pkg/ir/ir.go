/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

import "github.com/nooga/let-go/pkg/vm"

// NodeID is an index into Function.Nodes — the identity of an SSA value.
// Zero is a valid id; the first emitted node is NodeID(0).
type NodeID int32

// BlockID is an index into Function.Blocks.
type BlockID int32

// Node is one IR value or terminator. The flat-array layout means every
// node has a stable index referenced by other nodes via NodeID. There
// are no Go pointers between nodes; just integer references. This
// keeps the IR cache-friendly and easy to serialize / mutate / dump.
type Node struct {
	Op    Op       // kind
	Refs  []NodeID // operands; semantics depend on Op
	Aux   any      // op-specific payload (const value, var ref, branch target)
	Block BlockID  // which block this node lives in
	Type  Type     // result type; Unknown until type inference runs

	// SourceInfos accumulates source spans associated with this node.
	// Build seeds with exactly one entry per node (from chunk.LookupSource).
	// Optimization passes that merge or rewrite nodes union the sets so
	// runtime error messages can attribute the IP to any of the originating
	// source locations. nil/empty slice = unknown.
	SourceInfos []vm.SourceInfo
}

// Type is a coarse classification used by type inference / specialization.
// Unknown is the initial state; passes refine it.
type Type uint8

const (
	TypeUnknown Type = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeNil
	TypeString
	TypeAny // catch-all for things we know are vm.Value but not which concrete kind
)

func (t Type) String() string {
	switch t {
	case TypeUnknown:
		return "?"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeBool:
		return "bool"
	case TypeNil:
		return "nil"
	case TypeString:
		return "string"
	case TypeAny:
		return "any"
	}
	return "??"
}

// BranchTarget describes the destination of a branch — the target
// block plus the values passed to that block's parameters. Carbon-style
// block-arg semantics; replaces phi nodes.
type BranchTarget struct {
	Target BlockID
	Args   []NodeID // one per Block.Params; positional
}

// CondTarget is the payload for OpBranchIf: two BranchTargets, one
// for the true case and one for false.
type CondTarget struct {
	True, False *BranchTarget
}

// Block is a sequence of straight-line nodes ending in a terminator.
type Block struct {
	ID     BlockID
	Params []NodeID // OpBlockArg nodes, declared at block entry
	Nodes  []NodeID // body nodes, in execution order (excluding terminator)
	Term   NodeID   // terminator (Op.IsTerminator())
	Preds  []BlockID
}

// Function is the unit of optimization — one let-go fn.
type Function struct {
	Name       string
	Arity      int
	IsVariadic bool

	Nodes  []Node  // all nodes, indexed by NodeID
	Blocks []Block // basic blocks, Block.ID == index
	Entry  BlockID

	// SourceConsts is the constant pool from the source bytecode. We
	// reference constants by index (saved in Const Aux) so we don't
	// have to clone the pool. The lowering step writes back into the
	// same pool.
	SourceConsts *vm.Consts
}

// NewFunction creates an empty Function with one entry block.
//
// The entry block starts with no Params. Function arguments are
// referenced via OpLoadArg nodes wherever they're consumed — they live
// in f.args[] at runtime, not on the value stack, so they don't fit the
// BlockArg model (which represents values supplied by a predecessor's
// branch-with-args). Callers that need a fn-arg reference should
// f.AddNode(Node{Op: OpLoadArg, Aux: i, Block: someBlock}) explicitly.
func NewFunction(name string, arity int, variadic bool, consts *vm.Consts) *Function {
	f := &Function{
		Name:         name,
		Arity:        arity,
		IsVariadic:   variadic,
		SourceConsts: consts,
	}
	f.Entry = f.AddBlock()
	return f
}

// AddBlock allocates a new block, returns its ID.
func (f *Function) AddBlock() BlockID {
	id := BlockID(len(f.Blocks))
	f.Blocks = append(f.Blocks, Block{ID: id})
	return id
}

// AddNode appends a node, returns its ID.
func (f *Function) AddNode(n Node) NodeID {
	id := NodeID(len(f.Nodes))
	f.Nodes = append(f.Nodes, n)
	return id
}

// AppendToBlock adds an existing node ID to the block's body.
func (f *Function) AppendToBlock(b BlockID, id NodeID) {
	f.Blocks[b].Nodes = append(f.Blocks[b].Nodes, id)
}

// SetTerminator sets the block's terminator. Called once per block.
func (f *Function) SetTerminator(b BlockID, id NodeID) {
	f.Blocks[b].Term = id
}

// Node returns a pointer into the flat array. Stable across Function
// lifetime as long as no node is added (which would grow the slice
// and invalidate pointers — callers must re-fetch after AddNode).
func (f *Function) Node(id NodeID) *Node {
	return &f.Nodes[id]
}

// AddPred records that `pred` branches into `b`.
func (f *Function) AddPred(b, pred BlockID) {
	for _, p := range f.Blocks[b].Preds {
		if p == pred {
			return
		}
	}
	f.Blocks[b].Preds = append(f.Blocks[b].Preds, pred)
}

// MergeSourceInfo unions src into slice with dedup by (File, Line, Column).
// Returns the possibly-extended slice. Variadic so callers can union one or
// many spans at once (e.g. ir.MergeSourceInfo(dst, operand.SourceInfos...)).
func MergeSourceInfo(slice []vm.SourceInfo, src ...vm.SourceInfo) []vm.SourceInfo {
	for _, si := range src {
		dup := false
		for _, existing := range slice {
			if existing.File == si.File && existing.Line == si.Line && existing.Column == si.Column {
				dup = true
				break
			}
		}
		if !dup {
			slice = append(slice, si)
		}
	}
	return slice
}
