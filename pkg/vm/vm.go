/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// Opcodes
const (
	OP_NOOP int32 = iota // do nothing

	OP_LOAD_CONST // load constant LDC (index int32)
	OP_LOAD_ARG   // load argument LDA (index int32)

	OP_INVOKE // invoke function
	OP_RETURN // return from function

	OP_BRANCH_TRUE  // branch if truthy BRANCH_TRUE (offset int32)
	OP_BRANCH_FALSE // branch if falsy BRF (offset int32)
	OP_JUMP         // jump by offset JMP (offset int32)

	OP_POP     // pop value from the stack and discard it
	OP_POP_N   // save top and pop n elements from the stack POP_N (n int32)
	OP_DUP_NTH // duplicate nth value from the stack OPN (n int32)

	OP_SET_VAR  // set var
	OP_LOAD_VAR // push var root

	OP_MAKE_CLOSURE    // make a closure out of fn
	OP_LOAD_CLOSEDOVER // load closed over LDK (index int32)
	OP_PUSH_CLOSEDOVER // push closed over value to a closure

	OP_RECUR    // loop recurse RECUR (offset int32, argc int32)
	OP_RECUR_FN // function recurse REF (argc int32)

	OP_MAKE_MULTI_ARITY // make multi-arity function (n int32)
	OP_TAIL_CALL        // like OP_INVOKE but re-uses the frame

	OP_TRY_PUSH // push exception handler (catchOffset int32, finallyOffset int32)
	OP_TRY_POP  // pop exception handler (normal completion)
	OP_THROW    // throw top-of-stack value

	// Specialized arithmetic/comparison opcodes (binary, stack: [a b] -> [result])
	// These inline the common int64 fast path and fall back to generic NumXxx for other types.
	OP_ADD // + (2 args)
	OP_SUB // - (2 args)
	OP_MUL // * (2 args)
	OP_BIT_AND
	OP_BIT_OR
	OP_BIT_XOR
	OP_BIT_AND_NOT
	OP_BIT_SHIFT_LEFT
	OP_BIT_SHIFT_RIGHT
	OP_UNSIGNED_BIT_SHIFT_RIGHT
	OP_LT  // < (2 args)
	OP_LTE // <= (2 args)
	OP_GT  // > (2 args)
	OP_GTE // >= (2 args)
	OP_EQ  // = (2 args)
	OP_INC // inc (1 arg)
	OP_DEC // dec (1 arg)
	OP_BIT_NOT
	OP_QUOT // quot — integer quotient, truncated toward zero (2 args)
	OP_DIV  // / — true division; int/int yields a Ratio (or Int when exact), any float yields Float (2 args)

	OP_FINALLY_END // end of a finally block (finallyOffset int32, negative): rethrow the pending error after an abnormal entry

	OP_COUNT // sentinel — keep last; must equal len(opcodeNames) (enforced at init)
)

// opcodeNames maps opcode values to their mnemonics, in enum order. It is the
// canonical inventory of the opcode set: OpcodeToString indexes it and
// OpcodeSetSignature hashes it, so it must stay in sync with the const block
// above (same count, same order).
var opcodeNames = []string{
	"NOOP",
	"LOAD_CONST",
	"LOAD_ARG",
	"INVOKE",
	"RETURN",
	"BRANCH_T",
	"BRANCH_F",
	"JUMP",
	"POP",
	"POP_N",
	"DUP_NTH",
	"SET_VAR",
	"LOAD_VAR",
	"MAKE_CLOSURE",
	"LOAD_CLOSEDOVER",
	"PUSH_CLOSEDOVER",
	"RECUR",
	"RECUR_FN",
	"MAKE_MULTI_ARITY",
	"TAIL_CALL",
	"TRY_PUSH",
	"TRY_POP",
	"THROW",
	"ADD",
	"SUB",
	"MUL",
	"BIT_AND",
	"BIT_OR",
	"BIT_XOR",
	"BIT_AND_NOT",
	"BIT_SHIFT_LEFT",
	"BIT_SHIFT_RIGHT",
	"UNSIGNED_BIT_SHIFT_RIGHT",
	"LT",
	"LTE",
	"GT",
	"GTE",
	"EQ",
	"INC",
	"DEC",
	"BIT_NOT",
	"QUOT",
	"DIV",
	"FINALLY_END",
}

// A new opcode must land in both the const block and opcodeNames; the
// signature (and the disassembler) read only opcodeNames, so a missed update
// there would let bytecode change under an unchanged signature.
func init() {
	if len(opcodeNames) != int(OP_COUNT) {
		panic(fmt.Sprintf("opcodeNames out of sync with the opcode enum: %d names, %d opcodes", len(opcodeNames), OP_COUNT))
	}
}

func OpcodeToString(op int32) string {
	inst := op & 0xff
	sp := (op >> 16) & 0xffff
	if int(inst) < len(opcodeNames) {
		return fmt.Sprintf("%d/%-16s", sp, opcodeNames[inst])
	}
	return "???"
}

// ComputeSignatureForNames computes an opcode-set signature for a given list
// of opcode names. Used by the bytecode migration system to compute signatures
// for older opcode sets.
func ComputeSignatureForNames(names []string) (count int, hash uint64) {
	h := fnv.New64a()
	for _, name := range names {
		h.Write([]byte(name))
		h.Write([]byte{0})
	}
	return len(names), h.Sum64()
}

// OpcodeSetSignature identifies this VM's opcode set: the opcode count and an
// FNV-64a hash of the mnemonics joined in enum order. Bytecode bundles embed
// it at encode time (bytecode.CapOpcodeSet) so a decoder on a VM with a
// different opcode enum — opcodes inserted, removed, or reordered — can reject
// the bundle instead of executing shifted opcodes.
func OpcodeSetSignature() (count int, hash uint64) {
	return ComputeSignatureForNames(opcodeNames)
}

// CodeChunk holds bytecode and provides facilities for reading and writing it
// LocalVar is a debug-info entry mapping a stack slot to the original
// source name of a local binding (let/loop binding, fn parameter, catch
// binding). Carried so crash traces from bundle-loaded code (e.g. under
// WASM) can name local variables, not just slots.
type LocalVar struct {
	Slot int
	Name string
}

type CodeChunk struct {
	maxStack  int
	consts    *Consts
	code      []int32
	length    int
	sourceMap *SourceMap
	localVars []LocalVar
}

func NewCodeChunk(consts *Consts) *CodeChunk {
	return &CodeChunk{
		consts: consts,
		code:   []int32{},
		length: 0,
	}
}

// NewCodeChunkWithCapacity creates an empty chunk with preallocated code
// storage for n instructions.
func NewCodeChunkWithCapacity(consts *Consts, n int) *CodeChunk {
	return &CodeChunk{
		consts: consts,
		code:   make([]int32, 0, n),
		length: 0,
	}
}

// ReserveLocalVars ensures the local-var debug table can hold n entries
// without growing.
func (c *CodeChunk) ReserveLocalVars(n int) {
	if cap(c.localVars) >= n {
		return
	}
	next := make([]LocalVar, 0, n)
	next = append(next, c.localVars...)
	c.localVars = next
}

// AddLocalVar records the source name bound to a stack slot (debug info).
func (c *CodeChunk) AddLocalVar(slot int, name string) {
	c.localVars = append(c.localVars, LocalVar{Slot: slot, Name: name})
}

// LocalVars returns the chunk's local-variable debug table (may be nil).
func (c *CodeChunk) LocalVars() []LocalVar { return c.localVars }

func (c *CodeChunk) Debug() {
	consts := c.consts
	fmt.Println("code:")
	i := 0
	for i < len(c.code) {
		op, _ := c.Get(i)
		switch op & 0xff {
		case OP_TRY_PUSH:
			arg, _ := c.Get32(i + 1)
			arg2, _ := c.Get32(i + 2)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg, arg2)
			i += 3
		case OP_RECUR:
			arg, _ := c.Get32(i + 1)
			arg2, _ := c.Get32(i + 2)
			arg3, _ := c.Get32(i + 3)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg, arg2, arg3)
			i += 4
		case OP_LOAD_ARG, OP_BRANCH_TRUE, OP_BRANCH_FALSE, OP_JUMP, OP_POP_N, OP_DUP_NTH, OP_INVOKE, OP_LOAD_CLOSEDOVER, OP_RECUR_FN, OP_MAKE_MULTI_ARITY, OP_TAIL_CALL, OP_FINALLY_END:
			arg, _ := c.Get32(i + 1)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg)
			i += 2
		case OP_LOAD_CONST, OP_LOAD_VAR:
			arg, _ := c.Get32(i + 1)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg, "<-", consts.get(arg))
			i += 2
		default:
			fmt.Println("  ", i, ":", OpcodeToString(op))
			i++
		}
	}
}

func (c *CodeChunk) Length() int {
	return c.length
}

func (c *CodeChunk) Append(insts ...int32) {
	c.code = append(c.code, insts...)
	c.length = len(c.code)
}

func (c *CodeChunk) Append32(val int) {
	c.Append(int32(val))
}

func (c *CodeChunk) AppendChunk(o *CodeChunk) {
	if o.maxStack > c.maxStack {
		c.maxStack = o.maxStack
	}
	// Merge source maps with IP offset
	if o.sourceMap != nil {
		o.sourceMap.materialize() // realize a lazy map before reading its entries
		if c.sourceMap == nil {
			c.sourceMap = NewSourceMap()
		}
		base := c.length
		for _, e := range o.sourceMap.entries {
			c.sourceMap.Add(base+e.startIP, e.info)
		}
	}
	// Carry local-variable debug entries. Slots are frame-relative (not IP),
	// so they are appended as-is.
	c.localVars = append(c.localVars, o.localVars...)
	c.code = append(c.code, o.code...)
	c.length += len(o.code)
}

// AddSourceInfo records the source location for the current bytecode offset.
func (c *CodeChunk) AddSourceInfo(info SourceInfo) {
	c.AddSourceInfoAt(c.length, info)
}

// AddSourceInfoAt records the source location for the provided bytecode
// offset.
func (c *CodeChunk) AddSourceInfoAt(ip int, info SourceInfo) {
	if c.sourceMap == nil {
		c.sourceMap = NewSourceMap()
	}
	c.sourceMap.Add(ip, info)
}

// ReserveSourceMap ensures the source-map backing storage can hold n entries
// without growing.
func (c *CodeChunk) ReserveSourceMap(n int) {
	if c.sourceMap == nil {
		c.sourceMap = NewSourceMapWithCapacity(n)
		return
	}
	if cap(c.sourceMap.entries) >= n {
		return
	}
	next := NewSourceMapWithCapacity(n)
	next.entries = append(next.entries, c.sourceMap.entries...)
	c.sourceMap = next
}

// SetSourceMap replaces the chunk's source map. Used by the decoder to inject a
// lazily-materialized map (whose entries are decoded on first Lookup) so bundle
// load skips per-chunk source-map allocation for the common no-error path.
func (c *CodeChunk) SetSourceMap(sm *SourceMap) {
	c.sourceMap = sm
}

// LookupSource finds the source location for a given instruction pointer.
func (c *CodeChunk) LookupSource(ip int) *SourceInfo {
	if c.sourceMap == nil {
		return nil
	}
	return c.sourceMap.Lookup(ip)
}

func (c *CodeChunk) Get(idx int) (int32, error) {
	if idx >= c.length {
		return 0, NewExecutionError("bytecode fetch out of bounds")
	}
	return c.code[idx], nil
}

func (c *CodeChunk) Get32(idx int) (int, error) {
	n, err := c.Get(idx)
	return int(n), err
}

func (c *CodeChunk) Update32(address int, value int32) {
	c.code[address] = value
}

func (c *CodeChunk) SetMaxStack(max int) {
	c.maxStack = max
}

// Code returns the bytecode slice.
func (c *CodeChunk) Code() []int32 { return c.code }

// MaxStack returns the max stack depth.
func (c *CodeChunk) MaxStack() int { return c.maxStack }

// GetSourceMap returns the source map, may be nil.
func (c *CodeChunk) GetSourceMap() *SourceMap { return c.sourceMap }

// Consts returns the const pool reference.
func (c *CodeChunk) Consts() *Consts { return c.consts }

type exHandler struct {
	catchIP   int   // absolute IP of catch block (-1 if none, or already consumed)
	finallyIP int   // absolute IP of the shared finally block, entered by the abnormal path (-1 if none)
	savedSP   int   // stack depth to restore
	pending   error // in-flight error while the abnormal finally runs (nil otherwise)
}

// Frame is a single interpreter context
type Frame struct {
	stack       []Value
	args        []Value
	argbuf      []Value // frame-owned arguments used by same-frame tail replacement
	closedOvers []Value
	argc        int
	consts      *Consts
	constsc     int
	code        *CodeChunk
	ip          int
	sp          int
	debug       bool
	handlers    []exHandler  // exception handler stack (nil when unused)
	ec          *ExecContext // per-execution context (dynamic bindings); nil = none installed
	parent      *Frame       // suspended caller while the dispatch loop runs this frame
	prevOp      uint8        // opcode profiler state, preserved while this frame is suspended
	profileOn   bool         // profiler gate sampled at frame entry
}

// Frame reuse via a mutex-guarded LIFO.
//
// We previously used sync.Pool, but profiling fib(30) showed ~25% of
// CPU went to pool overhead (sync.(*Pool).Get/Put/pin) because the
// per-P bookkeeping (runtime.procPin etc.) is expensive relative to
// our tiny per-call work. A plain mutex-guarded stack is faster per
// Get/Put for the common single-goroutine case, and still safe under
// multi-goroutine workloads (per the async/go macro).
//
// The stack is capped at framePoolCap entries — beyond that, frames
// returned to the pool are dropped and left for GC. This bounds the
// memory let-go keeps live and matches sync.Pool's eventual-GC behavior.

const framePoolCap = 256

var (
	framePoolMu    sync.Mutex
	framePoolStack []*Frame // LIFO; pop from end
)

func acquireFrame() *Frame {
	framePoolMu.Lock()
	n := len(framePoolStack)
	if n == 0 {
		framePoolMu.Unlock()
		return &Frame{}
	}
	f := framePoolStack[n-1]
	framePoolStack[n-1] = nil // help GC see we no longer hold this slot
	framePoolStack = framePoolStack[:n-1]
	framePoolMu.Unlock()
	return f
}

func releaseFrame(f *Frame) {
	framePoolMu.Lock()
	if len(framePoolStack) < framePoolCap {
		framePoolStack = append(framePoolStack, f)
	}
	// else: drop on the floor; GC will reclaim
	framePoolMu.Unlock()
}

func NewFrame(code *CodeChunk, args []Value) *Frame {
	f := acquireFrame()
	needed := max(code.maxStack, 4)
	if cap(f.stack) >= needed {
		f.stack = f.stack[:needed]
	} else {
		f.stack = make([]Value, needed)
	}
	f.args = args
	f.argbuf = f.argbuf[:0]
	f.argc = len(args)
	f.closedOvers = nil
	f.consts = code.consts
	f.constsc = code.consts.count()
	f.code = code
	f.ip = 0
	f.sp = 0
	f.debug = false
	f.ec = nil
	f.parent = nil
	f.prevOp = 0
	f.profileOn = false
	if f.handlers != nil {
		f.handlers = f.handlers[:0]
	}
	return f
}

// ReleaseFrame returns a frame to the pool.
// We don't clear stack slots — they'll be overwritten on reuse.
// We only nil out the large reference fields to avoid pinning code/const objects.
func ReleaseFrame(f *Frame) {
	f.args = nil
	clear(f.argbuf)
	f.argbuf = nil
	f.closedOvers = nil
	f.consts = nil
	f.code = nil
	f.handlers = nil
	f.ec = nil
	f.parent = nil
	releaseFrame(f)
}

// newFrameForBytecodeCall is the non-tail transition for a prepared bytecode
// target. Fixed-arity target.args may borrow the suspended parent's operand
// stack; that is safe for a child frame because the parent does not resume
// until the child has completed.
func newFrameForBytecodeCall(target bytecodeCallTarget, ec *ExecContext) *Frame {
	child := NewFrame(target.fn.chunk, target.args)
	child.closedOvers = target.closedOvers
	child.ec = ec
	return child
}

// installBytecodeCall is the same-frame transition used by the existing
// direct-*Func tail-call path. target.args may be a window into f.stack, so
// copy it into the frame-owned buffer before resetting or reusing that stack.
func installBytecodeCall(f *Frame, target bytecodeCallTarget) {
	argc := len(target.args)
	if cap(f.argbuf) < argc {
		f.argbuf = make([]Value, argc)
	} else {
		clear(f.argbuf)
		f.argbuf = f.argbuf[:argc]
	}
	copy(f.argbuf, target.args)

	f.args = f.argbuf
	f.argc = argc
	f.closedOvers = target.closedOvers
	f.code = target.fn.chunk
	f.consts = f.code.consts
	f.constsc = f.code.consts.count()
	f.ip = 0
	f.sp = 0
	needed := max(f.code.maxStack, 4)
	if cap(f.stack) < needed {
		f.stack = make([]Value, needed)
	} else {
		f.stack = f.stack[:needed]
	}
}

func NewDebugFrame(code *CodeChunk, args []Value) *Frame {
	f := NewFrame(code, args)
	f.debug = true
	return f
}

// Fast-path stack operations. The compiler guarantees correct stack depth,
// so bounds checks are skipped for performance. Debug mode uses safe variants.

func (f *Frame) push(v Value) error {
	f.stack[f.sp] = v
	f.sp++
	return nil
}

func (f *Frame) pushMult(v []Value) error {
	copy(f.stack[f.sp:], v)
	f.sp += len(v)
	return nil
}

func (f *Frame) pop() (Value, error) {
	f.sp--
	return f.stack[f.sp], nil
}

func (f *Frame) nth(n int) (Value, error) {
	return f.stack[f.sp-1-n], nil
}

func (f *Frame) mult(start int, count int) ([]Value, error) {
	i := f.sp - start
	return f.stack[i-count : i], nil
}

func (f *Frame) drop(n int) error {
	f.sp -= n
	return nil
}

// Safe variants for debug mode
func (f *Frame) pushSafe(v Value) error {
	if f.sp >= f.code.maxStack {
		f.stackDbg()
		return NewExecutionError("stack overflow")
	}
	f.stack[f.sp] = v
	f.sp++
	return nil
}

func (f *Frame) popSafe() (Value, error) {
	if f.sp == 0 {
		f.stackDbg()
		return NIL, NewExecutionError("stack underflow")
	}
	f.sp--
	return f.stack[f.sp], nil
}

func (f *Frame) nthSafe(n int) (Value, error) {
	i := f.sp - 1 - n
	if i < 0 {
		f.stackDbg()
		return NIL, NewExecutionError("nth: stack underflow")
	}
	return f.stack[i], nil
}

func (f *Frame) dropSafe(n int) error {
	f.sp -= n
	if f.sp < 0 {
		f.stackDbg()
		f.code.Debug()
		return NewExecutionError("drop: stack underflow")
	}
	return nil
}

func (f *Frame) stackDbg() {
	fmt.Printf(";   stack [%d/%d]:\n", f.sp, f.code.maxStack)
	for i := 0; i < f.sp; i++ {
		fmt.Printf(";   %4d: %s\n", i, f.stack[i].String())
	}
	fmt.Println()
}

// handleError checks if there's an active try/catch handler and dispatches to it.
// Returns true if the error was handled (caller should continue the dispatch loop).
//
// Handler protocol (see tryCompiler for the emitted layout):
//   - An armed catch (catchIP >= 0) receives the error value. Entering it
//     consumes the catch; the handler stays on the stack when a finally
//     exists so a throw from the catch body still routes through it.
//   - With no armed catch but a finally, control enters the shared finally
//     block with the error parked in pending; OP_FINALLY_END at its end
//     re-dispatches the parked error.
//   - A throw escaping an abnormal finally (pending != nil) replaces the
//     in-flight error: the handler is discarded and unwinding continues
//     with the new error.
func (f *Frame) handleError(err error) bool {
	for len(f.handlers) > 0 {
		h := &f.handlers[len(f.handlers)-1]
		if h.pending != nil {
			f.handlers = f.handlers[:len(f.handlers)-1]
			continue
		}
		if h.catchIP >= 0 {
			catchIP := h.catchIP
			savedSP := h.savedSP
			h.catchIP = -1
			if h.finallyIP < 0 {
				f.handlers = f.handlers[:len(f.handlers)-1]
			}
			f.sp = savedSP
			f.push(errorToValue(err))
			f.ip = catchIP
			return true
		}
		if h.finallyIP >= 0 {
			h.pending = err
			f.sp = h.savedSP
			// The finally block is shared with the normal path, which enters
			// it with the try/catch result on the stack. Push a placeholder
			// so locals inside the finally resolve to the same slots.
			f.push(NIL)
			f.ip = h.finallyIP
			return true
		}
		f.handlers = f.handlers[:len(f.handlers)-1]
	}
	return false
}

// RunProtected runs the frame with panic recovery for thrownPanic.
// Use at top-level entry points (REPL, file eval). Internal calls use Run() directly.
func (f *Frame) RunProtected() (result Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			if tp, ok := r.(*thrownPanic); ok {
				err = tp.err
			} else {
				// Convert Go panic to let-go error with source info
				srcInfo := f.code.LookupSource(f.ip)
				msg := fmt.Sprintf("%v", r)
				err = NewExecutionError(msg).WithSource(srcInfo)
			}
		}
	}()
	return f.Run()
}

// suspendedCallArity validates the continuation stored in a suspended parent.
// Returning an error here is preferable to either panicking or trusting corrupt
// bytecode and silently dropping an arbitrary number of operand-stack values.
func suspendedCallArity(f *Frame) (int, *ExecutionError) {
	if f == nil || f.code == nil {
		return 0, NewExecutionError("invalid suspended call: missing frame code")
	}
	if f.ip < 0 || f.ip >= len(f.code.code) {
		return 0, NewExecutionError(fmt.Sprintf("invalid suspended call: instruction pointer %d is out of bounds", f.ip))
	}
	if f.ip+1 >= len(f.code.code) {
		return 0, NewExecutionError(fmt.Sprintf("invalid suspended call at %d: missing arity operand", f.ip))
	}
	if f.sp < 0 || f.sp > len(f.stack) {
		return 0, NewExecutionError(fmt.Sprintf("invalid suspended call at %d: stack depth %d is out of bounds", f.ip, f.sp))
	}
	op := f.code.code[f.ip] & 0xff
	if op != OP_INVOKE && op != OP_TAIL_CALL {
		return 0, NewExecutionError(fmt.Sprintf("invalid suspended call at %d: found %s", f.ip, OpcodeToString(f.code.code[f.ip])))
	}
	arity := int(f.code.code[f.ip+1])
	if arity < 0 || arity+1 > f.sp {
		return 0, NewExecutionError(fmt.Sprintf("invalid suspended call at %d: arity %d exceeds stack depth %d", f.ip, arity, f.sp))
	}
	return arity, nil
}

// wrapCallSite attributes err to the call currently under f.ip. The callee is
// still on f's operand stack (the call opcodes peek rather than pop), so the
// name is recoverable without extra bookkeeping.
func wrapCallSite(f *Frame, err error) error {
	arity, siteErr := suspendedCallArity(f)
	if siteErr != nil {
		return NewExecutionError(siteErr.message).Wrap(err)
	}
	srcInfo := f.code.LookupSource(f.ip)
	name := "fn"
	calleeIndex := f.sp - 1 - arity
	if fn, ok := AsFn(f.stack[calleeIndex]); ok {
		name = fnName(fn)
	}
	return NewExecutionError(fmt.Sprintf("calling %s", name)).WithSource(srcInfo).Wrap(err)
}

type frameRunState struct {
	root    *Frame
	current *Frame
}

func enterFrame(f *Frame) {
	if allocAttrEnabled {
		attrPushFrame(f)
	}
	// Dynamically-scoped tracing (*lg-trace*). Coarse gate first: TraceArmed is
	// false until *lg-trace* is first set truthy, so the precise per-frame Deref
	// only runs once tracing has been used. When *lg-trace* resolves truthy in
	// this frame's context, trace this frame — and, because the binding
	// propagates down the shared stack, every frame it calls.
	if TraceArmed.Load() {
		if tv := TraceVar.Load(); tv != nil && IsTruthy(f.ec.deref(tv)) {
			f.debug = true
		}
	}
	if f.debug {
		fmt.Print("run", f.args, "\n")
		f.code.Debug()
	}
	// The newly entered frame starts a fresh opcode-pair sequence at zero.
	// Suspended parents retain prevOp on their own Frame, so resuming records
	// the next parent-local pair; cross-frame pairs are intentionally omitted.
	f.prevOp = 0
	f.profileOn = ProfilingEnabled.Load()
}

// leaveFrame takes no frame on purpose: attribution keeps its own shadow
// stack, and a parameter here invites passing a frame that may already be
// back in the pool (OP_RETURN releases the child before state.current is
// rebound). No parameter, no use-after-release to reason about.
func leaveFrame() {
	if allocAttrEnabled {
		attrPopFrame()
	}
}

// releaseFailedFrames drops the current failed frame and each unhandled
// suspended parent. The root is owned by Run's caller and is never pooled here.
// It returns the first parent whose handler accepts err, or nil if none does.
func releaseFailedFrames(state *frameRunState, err error) (*Frame, error) {
	failed := state.current
	parent := failed.parent
	failed.parent = nil
	if failed != state.root {
		ReleaseFrame(failed)
	}
	for parent != nil {
		next := parent.parent
		err = wrapCallSite(parent, err)
		if parent.handleError(err) {
			state.current = parent
			return parent, err
		}
		leaveFrame()
		parent.parent = nil
		if parent != state.root {
			ReleaseFrame(parent)
		}
		parent = next
	}
	return nil, err
}

func releasePanickedFrames(state *frameRunState) {
	failed := state.current
	parent := failed.parent
	failed.parent = nil
	if failed != state.root {
		ReleaseFrame(failed)
	}
	for parent != nil {
		next := parent.parent
		leaveFrame()
		parent.parent = nil
		if parent != state.root {
			ReleaseFrame(parent)
		}
		parent = next
	}
}

// Run executes this frame to completion. The dispatch loop descends into bytecode
// callees within one dispatch loop, keeping the Go stack flat. A parent link
// on each live frame carries the suspended call chain without a retained slice.
func (f *Frame) Run() (result Value, resultErr error) {
	f.parent = nil
	state := frameRunState{root: f, current: f}
	// A defer statement in this function would cost deferreturn scaffolding on
	// every host->VM entry even when never registered (the compiler emits the
	// return-path walk unconditionally), so the panic-cleanup defer lives in a
	// separate wrapper taken only when allocation attribution is on. Without
	// attribution, a Go panic leaks the chain's pooled frames to GC — exactly
	// what pre-#645 did when a panic skipped runBytecodeCallTarget's
	// ReleaseFrame; the panic itself is preserved either way.
	if allocAttrEnabled {
		return runChainProtected(&state)
	}
	return runChain(&state)
}

// runChainProtected owns the panic-path cleanup that keeps allocation
// attribution balanced and pooled frames released. Kept out of Run and
// runChain so their fast paths carry no defer scaffolding.
func runChainProtected(state *frameRunState) (result Value, resultErr error) {
	defer func() {
		if r := recover(); r != nil {
			// runLoopAttr already left the current frame. Balance and
			// release every suspended child-owned frame before preserving the panic.
			releasePanickedFrames(state)
			panic(r)
		}
	}()
	return runChain(state)
}

// runChain drives the dispatch loop and the error-resume protocol for the
// frame chain rooted at state.root. The attribution gate lives here, hoisted
// out of the loop: a dispatcher function between runChain and the loop would
// be too big to inline (runLoopInner far exceeds the inlining budget), so it
// would cost a real call on every host->VM entry and every error-resume
// iteration. A defer statement inside the bare loop itself would cost
// deferreturn scaffolding on every VM entry, so the deferring variant stays
// a separate function.
func runChain(state *frameRunState) (result Value, resultErr error) {
	entering := true
	attr := allocAttrEnabled
	for {
		if attr {
			result, resultErr = state.current.runLoopAttr(state, entering)
		} else {
			result, resultErr = state.current.runLoopInner(state, entering)
		}
		entering = false
		if resultErr == nil {
			return result, nil
		}
		var resumed *Frame
		resumed, resultErr = releaseFailedFrames(state, resultErr)
		if resumed == nil {
			return NIL, resultErr
		}
	}
}

// runLoopAttr balances attrPopFrame for whichever logical frame is current
// when the loop unwinds; suspended parents remain active. Attribution keeps
// its own frame stack, so the wrapper can leave on behalf of the loop without
// tracking its final frame.
func (f *Frame) runLoopAttr(state *frameRunState, entering bool) (Value, error) {
	defer leaveFrame()
	return f.runLoopInner(state, entering)
}

// Cold-path error constructors, kept out of runLoopInner. Error construction
// only runs on failure, but when the constructor bodies (fmt varargs setup,
// struct literals, source-map lookups) are inlined into the ~25KB dispatch
// function they enlarge it and degrade its register allocation (the #706/#719
// FrameDispatch fragility). //go:noinline keeps each body a plain call so the
// hot function stays small.

//go:noinline
func execErr(msg string) error {
	return NewExecutionError(msg)
}

//go:noinline
func execWrap(msg string, err error) error {
	return NewExecutionError(msg).Wrap(err)
}

//go:noinline
func errNotAFunction(v Value) error {
	return NewTypeError(v, "is not a function", nil)
}

// wrapCallErr attributes err to the call at f.ip, naming the callee. Same
// message shape as wrapCallSite, but for the in-arm sites where the callee fn
// is already in hand.
//
//go:noinline
func wrapCallErr(f *Frame, fn Fn, err error) error {
	srcInfo := f.code.LookupSource(f.ip)
	return NewExecutionError(fmt.Sprintf("calling %s", fnName(fn))).WithSource(srcInfo).Wrap(err)
}

//go:noinline
func errIntOverflow() error {
	return NewExecutionError("integer overflow")
}

//go:noinline
func errDivByZero() error {
	return NewExecutionError("integer division by zero")
}

//go:noinline
func errBitOpType(name string) error {
	return fmt.Errorf("%s expected Int", name)
}

// debugTraceInst prints the per-instruction trace when frame tracing is on.
// Kept out of line so the fmt varargs boxing does not sit at the top of the
// dispatch loop's hottest block.
//
//go:noinline
func debugTraceInst(f *Frame, inst int32) {
	f.stackDbg()
	fmt.Println("#", f.ip, OpcodeToString(inst))
}

func (f *Frame) runLoopInner(state *frameRunState, entering bool) (Value, error) {
	if entering {
		enterFrame(f)
	}
	for {
		inst := f.code.code[f.ip]
		if f.debug {
			debugTraceInst(f, inst)
		}
		if f.profileOn {
			currOp := uint8(inst & 0xff)
			RecordOpcode(f.prevOp, currOp)
			f.prevOp = currOp
		}
		switch inst & 0xff {
		case OP_NOOP:
			f.ip++

		case OP_LOAD_CONST:
			idx := f.code.code[f.ip+1]
			if int(idx) >= f.constsc {
				return NIL, execErr("const lookup out of bounds")
			}
			err := f.push(f.consts.get(int(idx)))
			if err != nil {
				return NIL, execWrap("const push failed", err)
			}
			f.ip += 2

		case OP_LOAD_ARG:
			idx := f.code.code[f.ip+1]
			if int(idx) >= f.argc {
				return NIL, execErr("argument lookup out of bounds")
			}
			err := f.push(f.args[idx])
			if err != nil {
				return NIL, execWrap("argument push failed", err)
			}
			f.ip += 2

		case OP_RETURN:
			v, err := f.pop()
			if err != nil {
				return NIL, execWrap("return failed", err)
			}
			if f.parent == nil {
				return v, nil
			}
			child := f
			parent := child.parent
			child.parent = nil
			leaveFrame()
			ReleaseFrame(child)
			f = parent
			state.current = f
			callArity, callErr := suspendedCallArity(f)
			if callErr != nil {
				if f.handleError(callErr) {
					continue
				}
				return NIL, callErr
			}
			if e := f.drop(callArity + 1); e != nil {
				return NIL, execWrap("cleaning stack after call", e)
			}
			if e := f.push(v); e != nil {
				return NIL, execWrap("pushing return value failed", e)
			}
			f.ip += 2

		case OP_INVOKE:
			arity := f.code.code[f.ip+1]
			var out Value
			if arity > 0 {
				fraw, err := f.nth(int(arity))
				if err != nil {
					return NIL, execWrap("invoke instruction failed", err)
				}
				fn, ok := AsFn(fraw)
				if !ok {
					return NIL, errNotAFunction(fraw)
				}
				a, err := f.mult(0, int(arity))
				if err != nil {
					return NIL, execWrap("popping arguments failed", err)
				}
				target, direct, cerr := resolveBytecodeCall(fn, a)
				if cerr != nil {
					wrapped := wrapCallErr(f, fn, cerr)
					if f.handleError(wrapped) {
						continue
					}
					return NIL, wrapped
				}
				if direct {
					child := newFrameForBytecodeCall(target, f.ec)
					child.parent = f
					f = child
					state.current = f
					enterFrame(f)
					continue
				}
				out, err = f.ec.Invoke(fn, a)
				if err != nil {
					wrapped := wrapCallErr(f, fn, err)
					if f.handleError(wrapped) {
						continue
					}
					return NIL, wrapped
				}
				err = f.drop(int(arity) + 1)
				if err != nil {
					return NIL, execWrap("cleaning stack after call", err)
				}
			} else {
				// Peek rather than pop, so the callee slot remains for the
				// uniform drop(arity+1) when a descended child returns.
				fraw, err := f.nth(0)
				if err != nil {
					return NIL, execWrap("invoke instruction failed", err)
				}
				fn, ok := AsFn(fraw)
				if !ok {
					return NIL, errNotAFunction(fraw)
				}
				target, direct, cerr := resolveBytecodeCall(fn, nil)
				if cerr != nil {
					wrapped := wrapCallErr(f, fn, cerr)
					if f.handleError(wrapped) {
						continue
					}
					return NIL, wrapped
				}
				if direct {
					child := newFrameForBytecodeCall(target, f.ec)
					child.parent = f
					f = child
					state.current = f
					enterFrame(f)
					continue
				}
				if _, perr := f.pop(); perr != nil {
					return NIL, execWrap("invoke instruction failed", perr)
				}
				out, err = f.ec.Invoke(fn, nil)
				if err != nil {
					wrapped := wrapCallErr(f, fn, err)
					if f.handleError(wrapped) {
						continue
					}
					return NIL, wrapped
				}
			}
			err := f.push(out)
			if err != nil {
				return NIL, execWrap("pushing return value failed", err)
			}
			f.ip += 2

		case OP_TAIL_CALL:
			arity := f.code.code[f.ip+1]
			// Keep the callee on the operand stack when descending so the
			// return path can use one drop(arity+1) protocol for every arity.
			fraw, err := f.nth(int(arity))
			if err != nil {
				return NIL, execWrap("invoke instruction failed", err)
			}
			fn, ok := AsFn(fraw)
			if !ok {
				return NIL, errNotAFunction(fraw)
			}
			var a []Value
			if arity > 0 {
				a, err = f.mult(0, int(arity))
				if err != nil {
					return NIL, execWrap("popping arguments failed", err)
				}
			}

			target, direct, resolveErr := resolveBytecodeCall(fn, a)
			if resolveErr != nil {
				wrapped := wrapCallErr(f, fn, resolveErr)
				if f.handleError(wrapped) {
					continue
				}
				return NIL, wrapped
			}
			if direct {
				if _, reuse := fn.(*Func); reuse {
					installBytecodeCall(f, target)
					continue
				}
				// Closure, multi-arity, and metadata-wrapped bytecode
				// callees descend without re-entering Frame.Run. #620 will
				// move these targets onto the same-frame transition.
				child := newFrameForBytecodeCall(target, f.ec)
				child.parent = f
				f = child
				state.current = f
				enterFrame(f)
				continue
			}

			out, err := f.ec.Invoke(fn, a)
			if err != nil {
				wrapped := wrapCallErr(f, fn, err)
				if f.handleError(wrapped) {
					continue
				}
				return NIL, wrapped
			}
			// A native/non-bytecode tail call is terminal. Land its result
			// on the compiler-emitted RETURN so frame-chain unwind remains
			// centralized in OP_RETURN.
			if e := f.drop(int(arity) + 1); e != nil {
				return NIL, execWrap("cleaning stack after call", e)
			}
			if e := f.push(out); e != nil {
				return NIL, execWrap("pushing return value failed", e)
			}
			f.ip += 2
			continue

		case OP_BRANCH_TRUE:
			offset := f.code.code[f.ip+1]
			v, err := f.pop()
			if err != nil {
				return NIL, execWrap("BRANCH_TRUE pop condition", err)
			}
			if !IsTruthy(v) {
				f.ip += 2
				continue
			}
			f.ip += int(offset)

		case OP_BRANCH_FALSE:
			offset := f.code.code[f.ip+1]
			v, err := f.pop()
			if err != nil {
				return NIL, execWrap("BRANCH_FALSE pop condition", err)
			}
			if IsTruthy(v) {
				f.ip += 2
				continue
			}
			f.ip += int(offset)

		case OP_JUMP:
			offset := f.code.code[f.ip+1]
			f.ip += int(offset)

		case OP_POP:
			_, err := f.pop()
			if err != nil {
				return NIL, execWrap("POP failed", err)
			}
			f.ip++

		case OP_POP_N:
			v, err := f.pop()
			if err != nil {
				return NIL, execWrap("POP_N top value", err)
			}
			num := f.code.code[f.ip+1]
			err = f.drop(int(num))
			if err != nil {
				return NIL, execWrap("POP_N drop", err)
			}
			err = f.push(v)
			if err != nil {
				return NIL, execWrap("POP_N push", err)
			}
			f.ip += 2

		case OP_DUP_NTH:
			num := f.code.code[f.ip+1]
			val, err := f.nth(int(num))
			if err != nil {
				return NIL, execWrap("DUP_NTH get nth", err)
			}
			err = f.push(val)
			if err != nil {
				return NIL, execWrap("DUP_NTH push", err)
			}
			f.ip += 2

		case OP_SET_VAR:
			val, err := f.pop()
			if err != nil {
				return NIL, execWrap("SET_VAR pop value failed", err)
			}
			varr, err := f.pop()
			if err != nil {
				return NIL, execWrap("SET_VAR pop var failed", err)
			}
			varrd, ok := varr.(*Var)
			if !ok {
				return NIL, execWrap("SET_VAR invalid Var", err)
			}
			// (set! *v* val) mutates the current execution's top dynamic
			// binding (thread-local, matching Clojure) when one is active.
			// PERMISSIVE DEVIATION: with no binding in scope let-go falls
			// through to mutating the root, whereas Clojure throws
			// "Can't set!: ... from non-binding thread". This leniency is
			// load-bearing — core test.lg's run-tests set!s *report-counters*
			// etc. at the root with no surrounding binding.
			if !f.ec.setBinding(varrd, val) {
				varrd.SetRoot(val)
			}
			// Arm frame tracing when *lg-trace* is set truthy (pointer-compared;
			// no-op for every other var).
			armTraceIfTruthy(varrd, val)
			err = f.push(varr)
			if err != nil {
				return NIL, execWrap("SET_VAR push var failed", err)
			}
			f.ip++

		case OP_LOAD_VAR:
			idx := f.code.code[f.ip+1]
			if int(idx) >= f.constsc {
				return NIL, execErr("const lookup out of bounds")
			}
			err := f.push(f.ec.deref(f.consts.get(int(idx)).(*Var)))
			if err != nil {
				return NIL, execWrap("const push failed", err)
			}
			f.ip += 2

		case OP_MAKE_CLOSURE:
			idx := f.sp - 1
			if idx < 0 {
				return NIL, execErr("MAKE_CLOSURE stack underflow")
			}
			type closureCreator interface {
				MakeClosure() Fn
			}
			fn, ok := f.stack[idx].(closureCreator)
			if !ok {
				return NIL, execErr("MAKE_CLOSURE invalid func on stack")
			}
			f.stack[idx] = fn.MakeClosure()
			f.ip++

		case OP_LOAD_CLOSEDOVER:
			idx := f.code.code[f.ip+1]
			if int(idx) >= len(f.closedOvers) {
				return NIL, execErr("closed over lookup out of bounds")
			}
			err := f.push(f.closedOvers[idx])
			if err != nil {
				return NIL, execWrap("closed over push failed", err)
			}
			f.ip += 2

		case OP_PUSH_CLOSEDOVER:
			val, err := f.pop()
			if err != nil {
				return NIL, execWrap("popping closed over value failed", err)
			}
			idx := f.sp - 1
			if idx < 0 {
				return NIL, execWrap("PUSH_CLOSEDOVER stack overflow", err)
			}
			cls := f.stack[idx]
			if cls.Type() != FuncType {
				return NIL, execErr("PUSH_CLOSEDOVER expected a Fn")
			}
			fun, ok := cls.(*Closure)
			if !ok {
				return NIL, execErr("PUSH_CLOSEDOVER invalid closure on stack")
			}
			fun.closedOvers = append(fun.closedOvers, val)
			f.ip++

		case OP_RECUR_FN:
			arity := f.code.code[f.ip+1]
			a, err := f.mult(0, int(arity))
			if err != nil {
				return NIL, execWrap("popping arguments failed", err)
			}
			copy(f.args, a)
			f.argc = int(arity)
			f.sp = 0
			f.ip = 0

		case OP_RECUR:
			offset := f.code.code[f.ip+1]
			argc := f.code.code[f.ip+2]
			ignore := f.code.code[f.ip+3]
			a, err := f.mult(0, int(argc))
			if err != nil {
				return NIL, execWrap("RECUR popping arguments failed", err)
			}
			err = f.drop(int(argc)*2 + int(ignore))
			if err != nil {
				return NIL, execWrap("RECUR popping old locals", err)
			}
			err = f.pushMult(a)
			if err != nil {
				return NIL, execWrap("RECUR pushing new locals", err)
			}

			f.ip -= int(offset)

		case OP_MAKE_MULTI_ARITY:
			nfns := f.code.code[f.ip+1]
			n := int(nfns)
			fns, err := f.mult(0, n)
			if err != nil {
				return NIL, execWrap("MAKE_MULTI_ARITY popping arguments failed", err)
			}
			f.sp -= n

			fn, err := MakeMultiArity(fns)
			if err != nil {
				return NIL, execWrap("MAKE_MULTI_ARITY failed", err)
			}

			err = f.push(fn)
			if err != nil {
				return NIL, execWrap("MAKE_MULTI_ARITY push failed", err)
			}

			f.ip += 2

		case OP_TRY_PUSH:
			catchOffset := f.code.code[f.ip+1]
			finallyOffset := f.code.code[f.ip+2]
			// Offset 0 is a sentinel for "absent": a finally-only try has no
			// catch block, and a throw must route to the finally, then rethrow.
			catchIP := -1
			if catchOffset != 0 {
				catchIP = f.ip + int(catchOffset)
			}
			finallyIP := -1
			if finallyOffset != 0 {
				finallyIP = f.ip + int(finallyOffset)
			}
			f.handlers = append(f.handlers, exHandler{
				catchIP:   catchIP,
				finallyIP: finallyIP,
				savedSP:   f.sp,
			})
			f.ip += 3

		case OP_TRY_POP:
			if len(f.handlers) > 0 {
				f.handlers = f.handlers[:len(f.handlers)-1]
			}
			f.ip++

		case OP_THROW:
			v, _ := f.pop()
			thrown := NewThrownError(v)
			if f.handleError(thrown) {
				continue
			}
			return NIL, thrown

		case OP_FINALLY_END:
			// End of a finally block. On a normal entry this is a no-op; the
			// handler was already retired by TRY_POP. After an abnormal entry
			// the try's own handler is still on top with the error parked in
			// pending (handleError keeps it), so resume unwinding with it.
			// Ownership is checked by finally address: an OUTER handler may
			// also carry a pending error while this try runs normally inside
			// that handler's finally, and it must not be resumed from here.
			if n := len(f.handlers); n > 0 {
				h := f.handlers[n-1]
				if h.pending != nil && h.finallyIP == f.ip+int(f.code.code[f.ip+1]) {
					f.handlers = f.handlers[:n-1]
					if f.handleError(h.pending) {
						continue
					}
					return NIL, h.pending
				}
			}
			f.ip += 2

		// --- Specialized fast-path opcodes ---
		// These avoid NativeFn.Invoke, interface boxing, and RecoverPanic.
		// The compiler emits these for known binary calls to core arithmetic/comparison.

		case OP_ADD:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					r, ok := checkedAddInt(ai, bi)
					if !ok {
						if f.handleError(errIntOverflow()) {
							continue
						}
						return NIL, errIntOverflow()
					}
					f.stack[f.sp-2] = r
					f.sp--
					f.ip++
					continue
				}
			}
			r, err := NumAdd(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = r
			f.sp--
			f.ip++

		case OP_SUB:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					r, ok := checkedSubInt(ai, bi)
					if !ok {
						if f.handleError(errIntOverflow()) {
							continue
						}
						return NIL, errIntOverflow()
					}
					f.stack[f.sp-2] = r
					f.sp--
					f.ip++
					continue
				}
			}
			r, err := NumSub(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = r
			f.sp--
			f.ip++

		case OP_MUL:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					r, ok := checkedMulInt(ai, bi)
					if !ok {
						if f.handleError(errIntOverflow()) {
							continue
						}
						return NIL, errIntOverflow()
					}
					f.stack[f.sp-2] = r
					f.sp--
					f.ip++
					continue
				}
			}
			r, err := NumMul(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = r
			f.sp--
			f.ip++

		case OP_BIT_AND:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			ai, ok := a.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-and")) {
					continue
				}
				return NIL, errBitOpType("bit-and")
			}
			bi, ok := b.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-and")) {
					continue
				}
				return NIL, errBitOpType("bit-and")
			}
			f.stack[f.sp-2] = MakeInt(int(ai) & int(bi))
			f.sp--
			f.ip++

		case OP_BIT_OR:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			ai, ok := a.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-or")) {
					continue
				}
				return NIL, errBitOpType("bit-or")
			}
			bi, ok := b.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-or")) {
					continue
				}
				return NIL, errBitOpType("bit-or")
			}
			f.stack[f.sp-2] = MakeInt(int(ai) | int(bi))
			f.sp--
			f.ip++

		case OP_BIT_XOR:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			ai, ok := a.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-xor")) {
					continue
				}
				return NIL, errBitOpType("bit-xor")
			}
			bi, ok := b.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-xor")) {
					continue
				}
				return NIL, errBitOpType("bit-xor")
			}
			f.stack[f.sp-2] = MakeInt(int(ai) ^ int(bi))
			f.sp--
			f.ip++

		case OP_BIT_AND_NOT:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			ai, ok := a.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-and-not")) {
					continue
				}
				return NIL, errBitOpType("bit-and-not")
			}
			bi, ok := b.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-and-not")) {
					continue
				}
				return NIL, errBitOpType("bit-and-not")
			}
			f.stack[f.sp-2] = MakeInt(int(ai) &^ int(bi))
			f.sp--
			f.ip++

		case OP_BIT_SHIFT_LEFT:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			ai, ok := a.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-shift-left")) {
					continue
				}
				return NIL, errBitOpType("bit-shift-left")
			}
			bi, ok := b.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-shift-left")) {
					continue
				}
				return NIL, errBitOpType("bit-shift-left")
			}
			f.stack[f.sp-2] = MakeInt(int(ai) << uint(bi))
			f.sp--
			f.ip++

		case OP_BIT_SHIFT_RIGHT:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			ai, ok := a.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-shift-right")) {
					continue
				}
				return NIL, errBitOpType("bit-shift-right")
			}
			bi, ok := b.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-shift-right")) {
					continue
				}
				return NIL, errBitOpType("bit-shift-right")
			}
			f.stack[f.sp-2] = MakeInt(int(ai) >> uint(bi))
			f.sp--
			f.ip++

		case OP_UNSIGNED_BIT_SHIFT_RIGHT:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			ai, ok := a.(Int)
			if !ok {
				if f.handleError(errBitOpType("unsigned-bit-shift-right")) {
					continue
				}
				return NIL, errBitOpType("unsigned-bit-shift-right")
			}
			bi, ok := b.(Int)
			if !ok {
				if f.handleError(errBitOpType("unsigned-bit-shift-right")) {
					continue
				}
				return NIL, errBitOpType("unsigned-bit-shift-right")
			}
			f.stack[f.sp-2] = MakeInt(int(uint(ai) >> uint(bi)))
			f.sp--
			f.ip++

		case OP_QUOT:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					if bi == 0 {
						if f.handleError(errDivByZero()) {
							continue
						}
						return NIL, errDivByZero()
					}
					// Int quot is truncated toward zero — Go's / on
					// signed ints matches Clojure's quot semantics.
					f.stack[f.sp-2] = Int(int64(ai) / int64(bi))
					f.sp--
					f.ip++
					continue
				}
			}
			r, err := NumQuot(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = r
			f.sp--
			f.ip++

		case OP_DIV:
			// True division (clojure.core//). NumDiv produces the exact
			// numeric-tower result: Int/Int -> Ratio (or Int when exact),
			// any Float -> Float. Identical to the `/` core fn so the IR
			// path (OP_DIV / rt.DivValue) never diverges from the interpreter.
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			r, err := NumDiv(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = r
			f.sp--
			f.ip++

		case OP_BIT_NOT:
			a := f.stack[f.sp-1]
			ai, ok := a.(Int)
			if !ok {
				if f.handleError(errBitOpType("bit-not")) {
					continue
				}
				return NIL, errBitOpType("bit-not")
			}
			f.stack[f.sp-1] = MakeInt(^int(ai))
			f.ip++

		case OP_LT:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					f.stack[f.sp-2] = Boolean(int64(ai) < int64(bi))
					f.sp--
					f.ip++
					continue
				}
			}
			r, err := NumLt(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = Boolean(r)
			f.sp--
			f.ip++

		case OP_LTE:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					f.stack[f.sp-2] = Boolean(int64(ai) <= int64(bi))
					f.sp--
					f.ip++
					continue
				}
			}
			r, err := NumLe(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = Boolean(r)
			f.sp--
			f.ip++

		case OP_GT:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					f.stack[f.sp-2] = Boolean(int64(ai) > int64(bi))
					f.sp--
					f.ip++
					continue
				}
			}
			r, err := NumGt(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = Boolean(r)
			f.sp--
			f.ip++

		case OP_GTE:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					f.stack[f.sp-2] = Boolean(int64(ai) >= int64(bi))
					f.sp--
					f.ip++
					continue
				}
			}
			r, err := NumGe(a, b)
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-2] = Boolean(r)
			f.sp--
			f.ip++

		case OP_EQ:
			b := f.stack[f.sp-1]
			a := f.stack[f.sp-2]
			// Int fast path (most common in arithmetic code)
			if ai, ok := a.(Int); ok {
				if bi, ok := b.(Int); ok {
					f.stack[f.sp-2] = Boolean(ai == bi)
					f.sp--
					f.ip++
					continue
				}
			}
			// Keyword fast path
			if ak, ok := a.(Keyword); ok {
				if bk, ok := b.(Keyword); ok {
					f.stack[f.sp-2] = Boolean(ak == bk)
					f.sp--
					f.ip++
					continue
				}
			}
			f.stack[f.sp-2] = Boolean(ValueEquals(a, b))
			f.sp--
			f.ip++

		case OP_INC:
			a := f.stack[f.sp-1]
			if ai, ok := a.(Int); ok {
				r, ok := checkedAddInt(ai, 1)
				if !ok {
					if f.handleError(errIntOverflow()) {
						continue
					}
					return NIL, errIntOverflow()
				}
				f.stack[f.sp-1] = r
				f.ip++
				continue
			}
			r, err := NumAdd(a, Int(1))
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-1] = r
			f.ip++

		case OP_DEC:
			a := f.stack[f.sp-1]
			if ai, ok := a.(Int); ok {
				r, ok := checkedSubInt(ai, 1)
				if !ok {
					if f.handleError(errIntOverflow()) {
						continue
					}
					return NIL, errIntOverflow()
				}
				f.stack[f.sp-1] = r
				f.ip++
				continue
			}
			r, err := NumSub(a, Int(1))
			if err != nil {
				if f.handleError(err) {
					continue
				}
				return NIL, err
			}
			f.stack[f.sp-1] = r
			f.ip++

		default:
			return NIL, execErr("unknown instruction")
		}
	}
}

// fnName returns a human-readable name for a function value.
func fnName(fn Fn) string {
	switch f := fn.(type) {
	case *Func:
		if f.name != "" {
			return f.name
		}
		return "anonymous fn"
	case *Closure:
		return fnName(f.fn)
	case *MultiArityFn:
		if f.name != "" {
			return f.name
		}
		return "anonymous fn"
	case *NativeFn:
		if f.name != "" {
			return f.name
		}
		return "native fn"
	default:
		return "fn"
	}
}
