/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

import (
	"errors"
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// ErrUnsupportedOp is returned when Build encounters a bytecode opcode
// it doesn't yet handle. Callers should fall back to the original
// bytecode in this case.
var ErrUnsupportedOp = errors.New("ir.Build: unsupported opcode")

// Build constructs an SSA IR Function from a bytecode chunk.
//
// Algorithm (indexed RPN + block-arg SSA, stack-shape pre-pass):
//
//	Phase 1: scan for basic block leaders (branch targets + insts
//	         following branches). Each leader starts a new block.
//	         Compute reachability via CFG BFS.
//
//	Phase 1.5 (stack-shape pre-pass): abstract-interpret the bytecode
//	         to compute the value-stack depth at each block's leader IP.
//	         Iterative worklist algorithm to fixed point — handles loop
//	         back-edges (RECUR) naturally. Each reachable block ends up
//	         with a definite entryHeight ≥ 0; unreachable blocks stay
//	         at the sentinel -1.
//
//	Phase 2: for each block in CFG order, walk the bytecode with a
//	         parallel "value stack". The stack starts pre-populated
//	         with `entryHeight` OpBlockArg nodes (slot 0 = deepest).
//	         Each opcode is translated to an IR node, popping/pushing
//	         NodeIDs from the stack. Stack underflow is now impossible
//	         by construction; any underflow indicates a builder bug.
//
//	Phase 3: at block joins, each predecessor's terminator passes its
//	         exit-stack as branch args, sliced to exactly the
//	         successor's entryHeight (= len(Block.Params)). The stack-
//	         shape pre-pass guarantees every reachable predecessor has
//	         enough on its exit-stack — no fallback required.
//
// On the first unsupported opcode, returns ErrUnsupportedOp.
func Build(chunk *vm.CodeChunk, name string, arity int, variadic bool) (*Function, error) {
	if chunk == nil {
		return nil, errors.New("ir.Build: nil chunk")
	}
	b := &builder{
		chunk:     chunk,
		f:         NewFunction(name, arity, variadic, chunk.Consts()),
		ipToBlock: map[int]BlockID{},
	}

	// Phase 1: discover basic blocks and compute reachability.
	if err := b.discoverBlocks(); err != nil {
		return nil, err
	}

	// Phase 1.5: compute entry-stack-height per block by abstract interpretation.
	if err := b.computeStackHeights(); err != nil {
		return nil, err
	}

	// Phase 2: walk each block, build nodes via RPN simulation. The
	// value stack starts pre-populated with entry-height BlockArg
	// nodes per block.
	if err := b.buildBlocks(); err != nil {
		return nil, err
	}

	// Phase 3: wire up each predecessor's branch args to its
	// successors' pre-declared Params.
	if err := b.insertBlockArgs(); err != nil {
		return nil, err
	}

	// Phase 4: prune unreachable blocks (e.g., the dead fall-through
	// after OP_RECUR). Compacts Blocks and remaps BlockID references in
	// branch targets, predecessor lists, and Entry.
	b.pruneUnreachable()

	return b.f, nil
}

type builder struct {
	chunk        *vm.CodeChunk
	f            *Function
	ipToBlock    map[int]BlockID  // bytecode IP → block (for leaders)
	blockIPs     []int            // sorted leader IPs; index = blockID
	exitStack    [][]NodeID       // exit-stack per block (set during buildBlocks)
	reachable    map[BlockID]bool // reachable[bid] = true iff bid is forward-reachable from Entry
	entryHeights []int            // entryHeights[bid] = value-stack depth at block entry; -1 = unreachable
}

// discoverBlocks finds all basic block leaders. A leader is:
//   - the first instruction (IP 0)
//   - any branch target
//   - any instruction immediately following a terminator
//
// Each leader gets a BlockID. Sets b.ipToBlock and b.blockIPs.
// Also computes b.reachable via CFG BFS from the entry block.
func (b *builder) discoverBlocks() error {
	code := b.chunk.Code()
	if len(code) == 0 {
		return errors.New("ir.Build: empty chunk")
	}

	// IPs that are block leaders.
	leaders := map[int]bool{0: true}

	// Walk once, recording branch targets + fall-throughs after terminators.
	i := 0
	for i < len(code) {
		op := code[i] & 0xff
		stride := instStride(op)
		if stride == 0 {
			return fmt.Errorf("%w: opcode 0x%x at ip %d", ErrUnsupportedOp, op, i)
		}
		switch op {
		case vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE:
			offset := int(code[i+1])
			leaders[i+offset] = true
			leaders[i+stride] = true // fall-through
		case vm.OP_JUMP:
			offset := int(code[i+1])
			leaders[i+offset] = true
			// no fall-through — but next ip is a leader if reachable
			// (we leave it to leader-set discovery to find it via other branches)
			if i+stride < len(code) {
				leaders[i+stride] = true
			}
		case vm.OP_RECUR:
			// Backward jump to loop header. offset is positive; runtime
			// computes f.ip -= offset.
			offset := int(code[i+1])
			leaders[i-offset] = true // loop header
			leaders[i+stride] = true // fall-through is unreachable, but harmless to mark
		case vm.OP_RETURN, vm.OP_TAIL_CALL, vm.OP_THROW:
			// terminators with no successor in-function (TAIL_CALL ends
			// this function's execution; RETURN/THROW unwind).
			if i+stride < len(code) {
				leaders[i+stride] = true
			}
		}
		i += stride
	}

	// Allocate one BlockID per leader, sorted by IP.
	// (The Function already has Entry = block 0 from NewFunction; we
	// reuse it for IP 0.)
	sortedIPs := sortedKeys(leaders)
	b.blockIPs = sortedIPs
	for i, ip := range sortedIPs {
		var id BlockID
		if i == 0 {
			id = b.f.Entry // IP 0 is the entry block
		} else {
			id = b.f.AddBlock()
		}
		b.ipToBlock[ip] = id
	}

	// Reachability pass: BFS from the entry block, using bytecode-level
	// successor information. We need this so insertBlockArgs can skip
	// dead predecessors (e.g. the unreachable fall-through after OP_RECUR)
	// without reverting to the unsound t.Args != nil check.
	b.reachable = make(map[BlockID]bool, len(sortedIPs))
	b.reachable[b.f.Entry] = true
	worklist := []BlockID{b.f.Entry}
	for len(worklist) > 0 {
		cur := worklist[0]
		worklist = worklist[1:]
		for _, succ := range b.blockSuccessors(cur) {
			if !b.reachable[succ] {
				b.reachable[succ] = true
				worklist = append(worklist, succ)
			}
		}
	}

	return nil
}

// computeStackHeights performs an abstract-interpretation pre-pass:
// for each block, compute the value-stack depth on entry. Iterative
// worklist algorithm to fixed point; required for loops where the
// header's entry-height depends on its back-edge predecessors.
//
// Domain: a non-negative integer (stack depth) per block, with -1 as
// the "unreached" sentinel.
//
// For each block, walking its bytecode and summing per-opcode stack
// effects yields its exit-height. Each successor's entry-height must
// equal the predecessor's exit-height (or be the max over all preds —
// the abstract interpretation tolerates a single value per block,
// requiring all predecessors to agree). For our supported bytecode
// shapes (one frame of CLJS-style compilation), this consistency holds.
func (b *builder) computeStackHeights() error {
	code := b.chunk.Code()
	n := len(b.f.Blocks)
	b.entryHeights = make([]int, n)
	for i := range b.entryHeights {
		b.entryHeights[i] = -1
	}
	b.entryHeights[b.f.Entry] = 0

	worklist := []BlockID{b.f.Entry}
	for len(worklist) > 0 {
		bid := worklist[0]
		worklist = worklist[1:]

		// Walk this block's bytecode, summing stack effects.
		bi := int(bid)
		leaderIP := b.blockIPs[bi]
		var endIP int
		if bi+1 < len(b.blockIPs) {
			endIP = b.blockIPs[bi+1]
		} else {
			endIP = len(code)
		}
		h := b.entryHeights[bid]
		i := leaderIP
		for i < endIP {
			op := code[i] & 0xff
			stride := instStride(op)
			if stride == 0 {
				return fmt.Errorf("%w: opcode 0x%x at ip %d", ErrUnsupportedOp, op, i)
			}
			h += stackEffect(op, code, i)
			if h < 0 {
				return fmt.Errorf("ir.Build: stack-shape: block %d underflowed at ip %d (height=%d after op 0x%x)", bid, i, h, op)
			}
			i += stride
		}
		exitHeight := h

		// Propagate to successors. Each successor's entry-height should
		// equal the predecessor's exit-height. If we discover a successor
		// for the first time, set it and enqueue. If we discover an
		// inconsistency, fail loudly — well-formed bytecode shouldn't
		// hit this.
		for _, succ := range b.blockSuccessors(bid) {
			if b.entryHeights[succ] == -1 {
				b.entryHeights[succ] = exitHeight
				worklist = append(worklist, succ)
			} else if b.entryHeights[succ] != exitHeight {
				return fmt.Errorf(
					"ir.Build: stack-shape inconsistency at block %d: pred %d arrives with height %d, was %d",
					succ, bid, exitHeight, b.entryHeights[succ])
			}
		}
	}

	return nil
}

// stackEffect returns the net change in value-stack depth caused by the
// opcode at code[ip]. For terminators that branch with the stack still
// populated (JUMP, RECUR, BRANCH_*), the effect reflects what's left on
// the stack at the moment of branching — i.e., the branch's pops are
// already accounted for. Successors' entry-heights are then equal to
// the predecessor's exit-height.
//
// Coverage parallels instStride: any opcode instStride handles, this
// function must handle.
func stackEffect(op int32, code []int32, ip int) int {
	switch op & 0xff {
	case vm.OP_NOOP:
		return 0
	case vm.OP_LOAD_CONST, vm.OP_LOAD_VAR, vm.OP_LOAD_ARG:
		return 1
	case vm.OP_SET_VAR:
		return 0 // pops 1, pushes 1
	case vm.OP_ADD, vm.OP_SUB, vm.OP_MUL,
		vm.OP_LT, vm.OP_LTE, vm.OP_GT, vm.OP_GTE, vm.OP_EQ:
		return -1 // pops 2, pushes 1
	case vm.OP_INC, vm.OP_DEC:
		return 0 // pops 1, pushes 1
	case vm.OP_POP:
		return -1
	case vm.OP_POP_N:
		// POP_N saves the top, drops N below, pushes top back. Net = -N.
		return -int(code[ip+1])
	case vm.OP_DUP_NTH:
		return 1
	case vm.OP_INVOKE:
		// Pops fn + arity args, pushes 1. Net = -arity.
		return -int(code[ip+1])
	case vm.OP_TAIL_CALL:
		// Pops fn + arity args, leaves none on stack (terminator).
		// Net (effective at terminator) = -(arity+1).
		return -int(code[ip+1]) - 1
	case vm.OP_RETURN:
		return -1
	case vm.OP_JUMP:
		return 0
	case vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE:
		return -1 // pops cond
	case vm.OP_RECUR:
		// RECUR(offset, argc, ignore): copies top argc to temp, drops
		// argc*2 + ignore from the top, pushes argc back. Net change in
		// stack depth = -(argc*2 + ignore) + argc = -(argc + ignore).
		//
		// At fixed-point, the target loop-header's entry-height equals
		// the height just before RECUR minus (argc + ignore). For the
		// common case argc==targetParams, ignore==0, the target entry-
		// height equals (heightBeforeRecur - argc) — i.e., RECUR drops
		// the old loop bindings AND positions the new ones; the
		// abstract domain just sees -(argc + ignore) net.
		argc := int(code[ip+2])
		ignore := int(code[ip+3])
		return -(argc + ignore)
	case vm.OP_THROW:
		return -1
	}
	// Any opcode instStride accepts must be handled above. If we hit
	// this, instStride has been extended without updating stackEffect.
	return 0
}

// buildBlocks walks each block's bytecode range and emits IR nodes.
// Each reachable block starts with entryHeight pre-declared OpBlockArg
// nodes on the value stack (slot 0 = deepest). Unreachable blocks are
// still translated but with no Params (their content is dead anyway).
//
// Records exit-stack per block for the block-arg phase.
func (b *builder) buildBlocks() error {
	code := b.chunk.Code()
	b.exitStack = make([][]NodeID, len(b.f.Blocks))

	for bi, leaderIP := range b.blockIPs {
		blockID := b.ipToBlock[leaderIP]

		// End IP: next leader's IP, or end of chunk.
		var endIP int
		if bi+1 < len(b.blockIPs) {
			endIP = b.blockIPs[bi+1]
		} else {
			endIP = len(code)
		}

		// Initialize the value stack with entryHeight OpBlockArg nodes.
		// Slot 0 (deepest) gets Aux=0; top of pre-populated stack is
		// Aux=entryHeight-1. These also become Block.Params in the same
		// order (Params[0] = deepest = stack[0]).
		//
		// Unreachable blocks (entryHeights[bid] == -1) skip this and
		// rely on the tolerant pop in translateOp to synthesize dummy
		// BlockArgs on underflow — their content is dead, but we still
		// translate them so all blocks have a terminator.
		entryH := b.entryHeights[blockID]
		var stack []NodeID
		if entryH > 0 {
			stack = make([]NodeID, 0, entryH)
			params := make([]NodeID, 0, entryH)
			for i := 0; i < entryH; i++ {
				pid := b.f.AddNode(Node{Op: OpBlockArg, Aux: i, Block: blockID})
				params = append(params, pid)
				stack = append(stack, pid)
			}
			b.f.Blocks[blockID].Params = params
		}

		// Walk this block's instructions.
		i := leaderIP
		for i < endIP {
			op := code[i] & 0xff
			stride := instStride(op)
			next := i + stride
			if err := b.translateOp(blockID, op, code, i, &stack); err != nil {
				return err
			}
			i = next
		}

		// If this block didn't emit a terminator (block ends without
		// an explicit JUMP/BRANCH/RETURN/TAIL_CALL), it falls through
		// to the next block. Synthesize an OpBranch terminator.
		if b.f.Blocks[blockID].Term == 0 && bi+1 < len(b.blockIPs) {
			// Special case: NodeID 0 is the entry block's first node;
			// any other block's Term==0 means "no terminator set".
			nextLeaderIP := b.blockIPs[bi+1]
			nextBlock := b.ipToBlock[nextLeaderIP]
			// Carry the current block's exit-stack as branch args so the
			// successor's block-params can be wired up by insertBlockArgs.
			// Use stackSnapshot to ensure a non-nil slice (even for an
			// empty stack) so branchArgsAlreadySet can detect it.
			bt := &BranchTarget{Target: nextBlock, Args: stackSnapshot(stack)}
			termID := b.f.AddNode(Node{Op: OpBranch, Aux: bt, Block: blockID})
			b.f.SetTerminator(blockID, termID)
			b.f.AddPred(nextBlock, blockID)
		}

		b.exitStack[blockID] = stack
	}
	return nil
}

// translateOp emits the IR node(s) for a single bytecode op, updating
// the value stack. Stack underflow is impossible by construction (the
// stack-shape pre-pass guarantees every block enters with the right
// number of slots); any underflow here is a builder bug and surfaces
// as a hard error.
func (b *builder) translateOp(blockID BlockID, op int32, code []int32, ip int, stack *[]NodeID) error {
	push := func(n Node) NodeID {
		n.Block = blockID
		if si := b.chunk.LookupSource(ip); si != nil && len(n.SourceInfos) == 0 {
			n.SourceInfos = []vm.SourceInfo{*si}
		}
		id := b.f.AddNode(n)
		b.f.AppendToBlock(blockID, id)
		*stack = append(*stack, id)
		return id
	}
	pushNoStack := func(n Node) NodeID {
		// For terminators (no stack output).
		n.Block = blockID
		if si := b.chunk.LookupSource(ip); si != nil && len(n.SourceInfos) == 0 {
			n.SourceInfos = []vm.SourceInfo{*si}
		}
		id := b.f.AddNode(n)
		return id
	}
	pop := func() (NodeID, error) {
		if len(*stack) == 0 {
			// Stack underflow. For unreachable blocks the abstract
			// interpretation didn't compute an entry-height, so we
			// synthesize a placeholder BlockArg so translation
			// completes (the block's content is dead — never
			// executed; insertBlockArgs skips unreachable preds).
			if !b.reachable[blockID] {
				idx := len(b.f.Blocks[blockID].Params)
				id := b.f.AddNode(Node{Op: OpBlockArg, Aux: idx, Block: blockID})
				b.f.Blocks[blockID].Params = append(b.f.Blocks[blockID].Params, id)
				return id, nil
			}
			return 0, fmt.Errorf("ir.Build: stack underflow in block %d at ip %d (abstract interpretation bug)", blockID, ip)
		}
		v := (*stack)[len(*stack)-1]
		*stack = (*stack)[:len(*stack)-1]
		return v, nil
	}
	popN := func(n int) ([]NodeID, error) {
		out := make([]NodeID, n)
		for i := n - 1; i >= 0; i-- {
			v, err := pop()
			if err != nil {
				return nil, err
			}
			out[i] = v
		}
		return out, nil
	}

	switch op {
	case vm.OP_NOOP:
		// nothing
		return nil

	case vm.OP_LOAD_CONST:
		idx := int(code[ip+1])
		val := b.chunk.Consts().AllValues()[idx]
		push(Node{Op: OpConst, Aux: val})
		return nil

	case vm.OP_LOAD_ARG:
		idx := int(code[ip+1])
		push(Node{Op: OpLoadArg, Aux: idx})
		return nil

	case vm.OP_LOAD_VAR:
		idx := int(code[ip+1])
		val := b.chunk.Consts().AllValues()[idx]
		push(Node{Op: OpLoadVar, Aux: val})
		return nil

	case vm.OP_SET_VAR:
		v, err := pop()
		if err != nil {
			return err
		}
		// SET_VAR's var-name is on the stack just below; original VM
		// pops the var ref + value. Here we model the var-ref as the
		// value already loaded; this is approximate. For the spike,
		// the support is limited.
		push(Node{Op: OpSetVar, Refs: []NodeID{v}})
		return nil

	case vm.OP_ADD, vm.OP_SUB, vm.OP_MUL,
		vm.OP_LT, vm.OP_LTE, vm.OP_GT, vm.OP_GTE, vm.OP_EQ:
		args, err := popN(2)
		if err != nil {
			return err
		}
		push(Node{Op: bytecodeToIROp(op), Refs: args})
		return nil

	case vm.OP_INC, vm.OP_DEC:
		v, err := pop()
		if err != nil {
			return err
		}
		push(Node{Op: bytecodeToIROp(op), Refs: []NodeID{v}})
		return nil

	case vm.OP_INVOKE:
		arity := int(code[ip+1])
		args, err := popN(arity)
		if err != nil {
			return err
		}
		fn, err := pop()
		if err != nil {
			return err
		}
		refs := make([]NodeID, 0, arity+1)
		refs = append(refs, fn)
		refs = append(refs, args...)
		push(Node{Op: OpCall, Refs: refs, Aux: arity})
		return nil

	case vm.OP_TAIL_CALL:
		arity := int(code[ip+1])
		args, err := popN(arity)
		if err != nil {
			return err
		}
		fn, err := pop()
		if err != nil {
			return err
		}
		refs := make([]NodeID, 0, arity+1)
		refs = append(refs, fn)
		refs = append(refs, args...)
		// TailCall is a terminator.
		termID := pushNoStack(Node{Op: OpTailCall, Refs: refs, Aux: arity})
		b.f.SetTerminator(blockID, termID)
		return nil

	case vm.OP_POP:
		_, err := pop()
		if err != nil {
			return err
		}
		// Don't emit anything — POP just discards the top stack value,
		// which in IR means the previously-pushed node is now unused.
		// DCE can prune it if pure.
		return nil

	case vm.OP_POP_N:
		// Pop n elements; save top, then drop n below, then push top.
		// VM semantics: keeps the topmost value, discards N just below.
		n := int(code[ip+1])
		top, err := pop()
		if err != nil {
			return err
		}
		for k := 0; k < n; k++ {
			if _, err := pop(); err != nil {
				return err
			}
		}
		*stack = append(*stack, top)
		return nil

	case vm.OP_DUP_NTH:
		// Duplicate the Nth-from-top value back onto the stack. In SSA,
		// this is just reusing an existing NodeID — no new node. The
		// stack-shape pre-pass guarantees enough depth for reachable
		// blocks; unreachable blocks get synthesized placeholder
		// BlockArgs (dead content).
		n := int(code[ip+1])
		if len(*stack) < n+1 {
			if !b.reachable[blockID] {
				// Synthesize placeholder BlockArgs at the bottom of the
				// stack to satisfy the DUP.
				for len(*stack) < n+1 {
					idx := len(b.f.Blocks[blockID].Params)
					id := b.f.AddNode(Node{Op: OpBlockArg, Aux: idx, Block: blockID})
					b.f.Blocks[blockID].Params = append(b.f.Blocks[blockID].Params, id)
					*stack = append([]NodeID{id}, (*stack)...)
				}
			} else {
				return fmt.Errorf("ir.Build: DUP_NTH %d in block %d at ip %d would underflow (stack depth %d)", n, blockID, ip, len(*stack))
			}
		}
		idx := len(*stack) - 1 - n
		*stack = append(*stack, (*stack)[idx])
		return nil

	case vm.OP_RETURN:
		v, err := pop()
		if err != nil {
			return err
		}
		termID := pushNoStack(Node{Op: OpReturn, Refs: []NodeID{v}})
		b.f.SetTerminator(blockID, termID)
		return nil

	case vm.OP_JUMP:
		offset := int(code[ip+1])
		target := b.ipToBlock[ip+offset]
		// Like fall-through, an explicit JUMP terminator carries its
		// exit-stack as branch args so the successor's block-params can
		// be wired up by insertBlockArgs. branchArgsAlreadySet will see
		// len(Args) == n and skip this predecessor on the next phase.
		bt := &BranchTarget{Target: target, Args: stackSnapshot(*stack)}
		termID := pushNoStack(Node{Op: OpBranch, Aux: bt})
		b.f.SetTerminator(blockID, termID)
		b.f.AddPred(target, blockID)
		return nil

	case vm.OP_RECUR:
		// Backward jump to loop header with new loop-binding values.
		// Stack at this point: [...old-stuff..., new-arg-1, ..., new-arg-N].
		// We pop the N new args; the "drop old locals" is implicit in
		// SSA (those values just become dead — DCE will remove if pure).
		offset := int(code[ip+1])
		argc := int(code[ip+2])
		// ignore (code[ip+3]) is for VM stack discipline; irrelevant in SSA.

		args, err := popN(argc)
		if err != nil {
			return err
		}
		target := b.ipToBlock[ip-offset]
		bt := &BranchTarget{Target: target, Args: args}
		termID := pushNoStack(Node{Op: OpBranch, Aux: bt})
		b.f.SetTerminator(blockID, termID)
		b.f.AddPred(target, blockID)
		return nil

	case vm.OP_BRANCH_FALSE, vm.OP_BRANCH_TRUE:
		cond, err := pop()
		if err != nil {
			return err
		}
		offset := int(code[ip+1])
		var trueIP, falseIP int
		stride := instStride(op)
		if op == vm.OP_BRANCH_FALSE {
			// Branch taken if cond is false.
			falseIP = ip + offset
			trueIP = ip + stride
		} else {
			trueIP = ip + offset
			falseIP = ip + stride
		}
		trueT := &BranchTarget{Target: b.ipToBlock[trueIP]}
		falseT := &BranchTarget{Target: b.ipToBlock[falseIP]}
		ct := &CondTarget{True: trueT, False: falseT}
		termID := pushNoStack(Node{Op: OpBranchIf, Refs: []NodeID{cond}, Aux: ct})
		b.f.SetTerminator(blockID, termID)
		b.f.AddPred(trueT.Target, blockID)
		b.f.AddPred(falseT.Target, blockID)
		return nil

	default:
		return fmt.Errorf("%w: opcode 0x%x at ip %d", ErrUnsupportedOp, op, ip)
	}
}

// insertBlockArgs wires each predecessor's branch args to its
// successors' pre-declared Params.
//
// With the stack-shape pre-pass, every block's Params are already
// populated in buildBlocks (length = entryHeight). Predecessors'
// exit-stacks are guaranteed (by the abstract interpretation's fixed
// point) to be at least entryHeight deep. The predecessor passes its
// deepest entryHeight values as the branch args.
//
// Branch terminators that were emitted with explicit Args (RECUR,
// JUMP, fall-through) are left alone (branchArgsAlreadySet skips them).
func (b *builder) insertBlockArgs() error {
	for blockID := range b.f.Blocks {
		bid := BlockID(blockID)
		if bid == b.f.Entry {
			continue // entry block has no incoming branches
		}
		required := len(b.f.Blocks[bid].Params)
		if required == 0 {
			continue
		}
		for _, p := range b.f.Blocks[bid].Preds {
			if !b.reachable[p] {
				continue
			}
			termID := b.f.Blocks[p].Term
			term := b.f.Node(termID)

			if branchArgsAlreadySet(term, bid, required) {
				continue
			}

			predStack := b.exitStack[p]
			if len(predStack) < required {
				return fmt.Errorf(
					"ir.Build: block %d expects %d args but pred %d only has %d on exit-stack",
					bid, required, p, len(predStack))
			}
			start := len(predStack) - required
			passed := make([]NodeID, required)
			copy(passed, predStack[start:])
			setBranchArgs(term, bid, passed)
		}
	}
	return nil
}

// pruneUnreachable removes blocks not reachable from Entry. It
// compacts f.Blocks (no gaps) and remaps all BlockID references —
// f.Entry, each block's Preds, each terminator's BranchTarget/CondTarget
// targets, and each Node's Block field — through the old→new ID map.
//
// Reachable blocks retain their relative order. Branch targets pointing
// at a pruned block would represent dynamically unreachable code (the
// only way to get one is the dead post-RECUR JUMP) — those edges go
// away with the pruned predecessor.
func (b *builder) pruneUnreachable() {
	// Build old→new BlockID remap. Reachable blocks get sequential new IDs.
	remap := make(map[BlockID]BlockID, len(b.f.Blocks))
	newBlocks := make([]Block, 0, len(b.f.Blocks))
	for i := range b.f.Blocks {
		oldID := BlockID(i)
		if !b.reachable[oldID] {
			continue
		}
		newID := BlockID(len(newBlocks))
		remap[oldID] = newID
		blk := b.f.Blocks[i]
		blk.ID = newID
		newBlocks = append(newBlocks, blk)
	}
	if len(newBlocks) == len(b.f.Blocks) {
		return // nothing to prune
	}

	// Rewrite Preds (drop any pred that was pruned) and Block field on
	// every node (only nodes belonging to retained blocks survive
	// meaningfully — nodes on pruned blocks become dangling, which is
	// fine since nothing references them).
	for i := range newBlocks {
		blk := &newBlocks[i]
		filtered := blk.Preds[:0]
		for _, p := range blk.Preds {
			if np, ok := remap[p]; ok {
				filtered = append(filtered, np)
			}
		}
		blk.Preds = filtered
	}

	// Rewrite each surviving node's Block + branch-target IDs.
	for i := range b.f.Nodes {
		n := &b.f.Nodes[i]
		if newID, ok := remap[n.Block]; ok {
			n.Block = newID
		}
		switch t := n.Aux.(type) {
		case *BranchTarget:
			if newID, ok := remap[t.Target]; ok {
				t.Target = newID
			}
		case *CondTarget:
			if t.True != nil {
				if newID, ok := remap[t.True.Target]; ok {
					t.True.Target = newID
				}
			}
			if t.False != nil {
				if newID, ok := remap[t.False.Target]; ok {
					t.False.Target = newID
				}
			}
		}
	}

	b.f.Entry = remap[b.f.Entry]
	b.f.Blocks = newBlocks
}

// setBranchArgs replaces term's branch args for the edge to bid.
func setBranchArgs(term *Node, bid BlockID, args []NodeID) {
	switch t := term.Aux.(type) {
	case *BranchTarget:
		if t.Target == bid {
			t.Args = args
		}
	case *CondTarget:
		if t.True != nil && t.True.Target == bid {
			t.True.Args = args
		}
		if t.False != nil && t.False.Target == bid {
			t.False.Args = args
		}
	}
}

// branchArgsAlreadySet reports whether term's branch to bid already has
// exactly n args wired. Used to skip predecessors (like RECUR, JUMP,
// and fall-through) that carry args set at translation time.
//
// We use len(t.Args) == n (not t.Args != nil) so that a stale or
// incorrectly-populated Args slice does not silently pass. The strong
// invariant: "wired" means "wired with the exact right count".
func branchArgsAlreadySet(term *Node, bid BlockID, n int) bool {
	switch t := term.Aux.(type) {
	case *BranchTarget:
		return t.Target == bid && len(t.Args) == n
	case *CondTarget:
		if t.True != nil && t.True.Target == bid && len(t.True.Args) == n {
			return true
		}
		if t.False != nil && t.False.Target == bid && len(t.False.Args) == n {
			return true
		}
	}
	return false
}

// --- helpers ---

// instStride returns the number of int32 words an opcode occupies,
// or 0 if the opcode is unsupported by Build.
func instStride(op int32) int {
	switch op & 0xff {
	case vm.OP_NOOP, vm.OP_POP, vm.OP_RETURN,
		vm.OP_ADD, vm.OP_SUB, vm.OP_MUL,
		vm.OP_LT, vm.OP_LTE, vm.OP_GT, vm.OP_GTE, vm.OP_EQ,
		vm.OP_INC, vm.OP_DEC, vm.OP_SET_VAR:
		return 1
	case vm.OP_LOAD_CONST, vm.OP_LOAD_VAR, vm.OP_LOAD_ARG,
		vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE, vm.OP_JUMP,
		vm.OP_INVOKE, vm.OP_TAIL_CALL,
		vm.OP_POP_N, vm.OP_DUP_NTH:
		return 2
	case vm.OP_RECUR:
		return 4 // offset, argc, ignore
	}
	return 0
}

// stackSnapshot returns a non-nil copy of s (empty slice for nil/empty s).
// Used when storing exit-stack values in BranchTarget.Args so that
// branchArgsAlreadySet can distinguish "not yet set" (nil) from "empty
// block had nothing to pass" (non-nil empty slice).
func stackSnapshot(s []NodeID) []NodeID {
	out := make([]NodeID, len(s))
	copy(out, s)
	return out
}

// blockSuccessors returns the BlockIDs that are direct successors of bid
// based on the bytecode of that block. Called during discoverBlocks after
// b.ipToBlock and b.blockIPs are populated.
//
// Rules:
//   - OP_BRANCH_TRUE / OP_BRANCH_FALSE → two successors (target + fall-through)
//   - OP_JUMP → one successor (target)
//   - OP_RECUR → one successor (loop header at ip - offset)
//   - OP_RETURN / OP_TAIL_CALL / OP_THROW → no in-function successors
//   - No terminator in range (fall-through block) → next block in blockIPs
func (b *builder) blockSuccessors(bid BlockID) []BlockID {
	code := b.chunk.Code()

	// Find this block's index in blockIPs to get its IP range.
	// bid == b.ipToBlock[blockIPs[bi]] by construction.
	bi := int(bid) // blockIPs index == BlockID index (both built in the same sorted loop)
	leaderIP := b.blockIPs[bi]
	var endIP int
	if bi+1 < len(b.blockIPs) {
		endIP = b.blockIPs[bi+1]
	} else {
		endIP = len(code)
	}

	// Walk the block's instructions to find the last (terminating) one.
	i := leaderIP
	for i < endIP {
		op := code[i] & 0xff
		stride := instStride(op)
		if stride == 0 {
			// Unknown opcode — treat as fall-through to be safe.
			break
		}
		next := i + stride
		if next >= endIP {
			// This is the last instruction in the block.
			switch op {
			case vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE:
				offset := int(code[i+1])
				target := b.ipToBlock[i+offset]
				fallThrough := b.ipToBlock[i+stride]
				return []BlockID{target, fallThrough}
			case vm.OP_JUMP:
				offset := int(code[i+1])
				target := b.ipToBlock[i+offset]
				return []BlockID{target}
			case vm.OP_RECUR:
				offset := int(code[i+1])
				target := b.ipToBlock[i-offset]
				return []BlockID{target}
			case vm.OP_RETURN, vm.OP_TAIL_CALL, vm.OP_THROW:
				return nil // no in-function successors
			default:
				// Fall-through to next block.
				if bi+1 < len(b.blockIPs) {
					return []BlockID{b.ipToBlock[b.blockIPs[bi+1]]}
				}
				return nil
			}
		}
		i = next
	}

	// Block was empty or ended without a recognized terminator: fall-through.
	if bi+1 < len(b.blockIPs) {
		return []BlockID{b.ipToBlock[b.blockIPs[bi+1]]}
	}
	return nil
}

func sortedKeys(m map[int]bool) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	// simple sort, small n
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}
