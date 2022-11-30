/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
)

// Opcodes
const (
	OP_NOOP int32 = iota // do nothing

	OP_LOAD_CONST // load constant LDC (index int32)
	OP_LOAD_ARG   // load argument LDA (index int32)

	OP_INVOKE // invoke function
	OP_RETURN // return from function

	OP_BRANCH_TRUE  // branch if truthy BRT (offset int32)
	OP_BRANCH_FALSE // branch if falsy BRF (offset int32)
	OP_JUMP         // jump by offset JMP (offset int32)

	OP_POP     // pop value from the stack and discard it
	OP_POP_N   // save top and pop n elements from the stack PON (n int32)
	OP_DUP_NTH // duplicate nth value from the stack OPN (n int32)

	OP_SET_VAR        // set var
	OP_LOAD_VAR       // push var root
	OP_LOAD_CONST_VAR // push derefed const

	OP_MAKE_CLOSURE    // make a closure out of fn
	OP_LOAD_CLOSEDOVER // load closed over LDK (index int32)
	OP_PUSH_CLOSEDOVER // push closed over value to a closure

	OP_RECUR    // loop recurse REC (offset int32, argc int32)
	OP_RECUR_FN // function recurse REF (argc int32)

	OP_TRACE_ENABLE  // enable tracing
	OP_TRACE_DISABLE // disable tracing

	OP_MAKE_MULTI_ARITY // make multi-arity function (n int32)
)

func OpcodeToString(op int32) string {
	inst := op & 0xff
	sp := (op >> 16) & 0xffff
	ops := []string{
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
		"LOAD_CONST_VAR",
		"MAKE_CLOSURE",
		"LOAD_CLOSEDOVER",
		"PUSH_CLOSEDOVER",
		"RECUR",
		"RECUR_FN",
		"TRACE_ENABLE",
		"TRACE_DISABLE",
		"MAKE_MULTI_ARITY",
	}
	if int(inst) < len(ops) {
		return fmt.Sprintf("%d/%-16s", sp, ops[inst])
	}
	return "???"
}

// CodeChunk holds bytecode and provides facilities for reading and writing it
type CodeChunk struct {
	maxStack int
	consts   *Consts
	code     []int32
	length   int
}

func NewCodeChunk(consts *Consts) *CodeChunk {
	return &CodeChunk{
		consts: consts,
		code:   []int32{},
		length: 0,
	}
}

func (c *CodeChunk) Debug() {
	consts := c.consts
	fmt.Println("code:")
	i := 0
	for i < len(c.code) {
		op, _ := c.Get(i)
		switch op & 0xff {
		case OP_RECUR:
			arg, _ := c.Get32(i + 1)
			arg2, _ := c.Get32(i + 2)
			arg3, _ := c.Get32(i + 3)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg, arg2, arg3)
			i += 4
		case OP_LOAD_ARG, OP_BRANCH_TRUE, OP_BRANCH_FALSE, OP_JUMP, OP_POP_N, OP_DUP_NTH, OP_INVOKE, OP_LOAD_CLOSEDOVER, OP_RECUR_FN, OP_MAKE_MULTI_ARITY:
			arg, _ := c.Get32(i + 1)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg)
			i += 2
		case OP_LOAD_CONST, OP_LOAD_CONST_VAR:
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
	c.code = append(c.code, o.code...)
	c.length += len(o.code)
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

// Frame is a single interpreter context
type Frame struct {
	stack       []Value
	args        []Value
	closedOvers []Value
	argc        int
	consts      *Consts
	constsc     int
	code        *CodeChunk
	ip          int
	sp          int
	debug       bool
}

func NewFrame(code *CodeChunk, args []Value) *Frame {
	return &Frame{
		stack:   make([]Value, code.maxStack),
		args:    args,
		argc:    len(args),
		consts:  code.consts,
		constsc: code.consts.count(),
		code:    code,
		ip:      0,
		sp:      0,
		debug:   false,
	}
}

func NewDebugFrame(code *CodeChunk, args []Value) *Frame {
	f := NewFrame(code, args)
	f.debug = true
	return f
}

func (f *Frame) push(v Value) error {
	if f.sp >= f.code.maxStack {
		f.stackDbg()
		return NewExecutionError("stack overflow")
	}
	f.stack[f.sp] = v
	f.sp++
	return nil
}

func (f *Frame) pushMult(v []Value) error {
	l := len(v)
	if f.sp >= f.code.maxStack-l {
		f.stackDbg()
		return NewExecutionError("stack overflow")
	}

	for i := 0; i < l; i++ {
		f.stack[f.sp] = v[i]
		f.sp++
	}
	return nil
}

func (f *Frame) pop() (Value, error) {
	if f.sp == 0 {
		f.stackDbg()
		return NIL, NewExecutionError("stack underflow")
	}
	f.sp--
	v := f.stack[f.sp]
	//f.stack[f.sp] = nil
	return v, nil
}

func (f *Frame) nth(n int) (Value, error) {
	i := f.sp - 1 - n
	if i < 0 {
		f.stackDbg()
		return NIL, NewExecutionError("nth: stack underflow")
	}
	return f.stack[i], nil
}

func (f *Frame) mult(start int, count int) ([]Value, error) {
	if count < 0 {
		f.stackDbg()
		return nil, NewExecutionError("mult: count 0 or negative")
	}
	i := f.sp - start
	if i-count < 0 {
		f.stackDbg()
		return nil, NewExecutionError("mult: stack underflow")
	}
	return f.stack[i-count : i], nil
}

func (f *Frame) drop(n int) error {
	top := f.sp - 1
	if top < 0 {
		f.stackDbg()
		f.code.Debug()
		return NewExecutionError("drop: stack underflow")
	}
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

func (f *Frame) Run() (Value, error) {
	if f.debug {
		fmt.Print("run", f.args, "\n")
		f.code.Debug()
	}
	for {
		inst := f.code.code[f.ip]
		if f.debug {
			f.stackDbg()
			fmt.Println("#", f.ip, OpcodeToString(inst))
		}
		switch inst & 0xff {
		case OP_NOOP:
			f.ip++

		case OP_TRACE_ENABLE:
			fmt.Print("# tracing frame, args: ", f.args, "\n")
			f.code.Debug()
			f.debug = true
			f.ip += 1

		case OP_TRACE_DISABLE:
			f.debug = false
			f.ip += 1

		case OP_LOAD_CONST:
			idx := f.code.code[f.ip+1]
			if int(idx) >= f.constsc {
				return NIL, NewExecutionError("const lookup out of bounds")
			}
			err := f.push(f.consts.get(int(idx)))
			if err != nil {
				return NIL, NewExecutionError("const push failed").Wrap(err)
			}
			f.ip += 2

		case OP_LOAD_ARG:
			idx := f.code.code[f.ip+1]
			if int(idx) >= f.argc {
				return NIL, NewExecutionError("argument lookup out of bounds")
			}
			err := f.push(f.args[idx])
			if err != nil {
				return NIL, NewExecutionError("argument push failed").Wrap(err)
			}
			f.ip += 2

		case OP_RETURN:
			v, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("return failed").Wrap(err)
			}
			return v, nil

		case OP_INVOKE:
			arity := f.code.code[f.ip+1]
			var out Value
			if arity > 0 {
				fraw, err := f.nth(int(arity))
				if err != nil {
					return NIL, NewExecutionError("invoke instruction failed").Wrap(err)
				}
				fn, ok := fraw.(Fn)
				if !ok {
					return NIL, NewTypeError(fraw, "is not a function", nil)
				}
				a, err := f.mult(0, int(arity))
				if err != nil {
					return NIL, NewExecutionError("popping arguments failed").Wrap(err)
				}
				out, err = fn.Invoke(a)
				if err != nil {
					// FIXME this is an exception, we should handle it
					return NIL, NewExecutionError(fmt.Sprintf("calling %s", fn.String())).Wrap(err)
				}
				err = f.drop(int(arity) + 1)
				if err != nil {
					return NIL, NewExecutionError("cleaning stack after call").Wrap(err)
				}
			} else {
				fraw, err := f.pop()
				if err != nil {
					return NIL, NewExecutionError("invoke instruction failed").Wrap(err)
				}
				fn, ok := fraw.(Fn)
				if !ok {
					return NIL, NewTypeError(fraw, "is not a function", nil)
				}
				out, err = fn.Invoke(nil)
				if err != nil {
					// FIXME this is an exception, we should handle it
					return NIL, NewExecutionError(fmt.Sprintf("calling %s", fn.String())).Wrap(err)
				}
			}
			err := f.push(out)
			if err != nil {
				return NIL, NewExecutionError("pushing return value failed").Wrap(err)
			}
			f.ip += 2

		case OP_BRANCH_TRUE:
			offset := f.code.code[f.ip+1]
			v, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("BRT pop condition").Wrap(err)
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
				return NIL, NewExecutionError("BRT pop condition").Wrap(err)
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
				return NIL, NewExecutionError("POP failed").Wrap(err)
			}
			f.ip++

		case OP_POP_N:
			v, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("PON top value").Wrap(err)
			}
			num := f.code.code[f.ip+1]
			err = f.drop(int(num))
			if err != nil {
				return NIL, NewExecutionError("PON drop").Wrap(err)
			}
			err = f.push(v)
			if err != nil {
				return NIL, NewExecutionError("PON push").Wrap(err)
			}
			f.ip += 2

		case OP_DUP_NTH:
			num := f.code.code[f.ip+1]
			val, err := f.nth(int(num))
			if err != nil {
				return NIL, NewExecutionError("DPN get nth").Wrap(err)
			}
			err = f.push(val)
			if err != nil {
				return NIL, NewExecutionError("DPN push").Wrap(err)
			}
			f.ip += 2

		case OP_SET_VAR:
			val, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("STV pop value failed").Wrap(err)
			}
			varr, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("STV pop var failed").Wrap(err)
			}
			varrd, ok := varr.(*Var)
			if !ok {
				return NIL, NewExecutionError("STV invalid Var").Wrap(err)
			}
			varrd.SetRoot(val)
			err = f.push(varr)
			if err != nil {
				return NIL, NewExecutionError("STV push var failed").Wrap(err)
			}
			f.ip++

		case OP_LOAD_VAR:
			// note this avoids pop-push dance
			idx := f.sp - 1
			if idx < 0 {
				return NIL, NewExecutionError("LDV stack underflow")
			}
			varr, ok := f.stack[idx].(*Var)
			if !ok {
				return NIL, NewExecutionError("LDV invalid var on stack")
			}
			f.stack[idx] = varr.Deref()
			f.ip++

		case OP_LOAD_CONST_VAR:
			idx := f.code.code[f.ip+1]
			if int(idx) >= f.constsc {
				return NIL, NewExecutionError("const lookup out of bounds")
			}
			err := f.push(f.consts.get(int(idx)).(*Var).Deref())
			if err != nil {
				return NIL, NewExecutionError("const push failed").Wrap(err)
			}
			f.ip += 2

		case OP_MAKE_CLOSURE:
			idx := f.sp - 1
			if idx < 0 {
				return NIL, NewExecutionError("MKC stack underflow")
			}
			fn, ok := f.stack[idx].(*Func)
			if !ok {
				return NIL, NewExecutionError("MKC invalid func on stack")
			}
			f.stack[idx] = fn.MakeClosure()
			f.ip++

		case OP_LOAD_CLOSEDOVER:
			idx := f.code.code[f.ip+1]
			// FIXME cache closedOvers count
			if int(idx) >= len(f.closedOvers) {
				return NIL, NewExecutionError("closed over lookup out of bounds")
			}
			err := f.push(f.closedOvers[idx])
			if err != nil {
				return NIL, NewExecutionError("closed over push failed").Wrap(err)
			}
			f.ip += 2

		case OP_PUSH_CLOSEDOVER:
			val, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("popping closed over value failed").Wrap(err)
			}
			idx := f.sp - 1
			if idx < 0 {
				return NIL, NewExecutionError("PAK stack overflow").Wrap(err)
			}
			cls := f.stack[idx]
			if cls.Type() != FuncType {
				return NIL, NewExecutionError("PAK expected a Fn")
			}
			fun, ok := cls.(*Closure)
			if !ok {
				return NIL, NewExecutionError("PAK invalid closure on stack")
			}
			fun.closedOvers = append(fun.closedOvers, val)
			f.ip++

		case OP_RECUR_FN:
			arity := f.code.code[f.ip+1]
			a, err := f.mult(0, int(arity))
			if err != nil {
				return NIL, NewExecutionError("popping arguments failed").Wrap(err)
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
				return NIL, NewExecutionError("REC popping arguments failed").Wrap(err)
			}
			err = f.drop(int(argc)*2 + int(ignore))
			if err != nil {
				return NIL, NewExecutionError("REC popping old locals").Wrap(err)
			}
			err = f.pushMult(a)
			if err != nil {
				return NIL, NewExecutionError("REC pushing new locals").Wrap(err)
			}

			f.ip -= int(offset)

		case OP_MAKE_MULTI_ARITY:
			nfns := f.code.code[f.ip+1]
			n := int(nfns)
			fns, err := f.mult(0, n)
			if err != nil {
				return NIL, NewExecutionError("MAKE_MULTI_ARITY popping arguments failed").Wrap(err)
			}
			f.sp -= n

			fn, err := makeMultiArity(fns)
			if err != nil {
				return NIL, NewExecutionError("MAKE_MULTI_ARITY failed").Wrap(err)
			}

			err = f.push(fn)
			if err != nil {
				return NIL, NewExecutionError("MAKE_MULTI_ARITY push failed").Wrap(err)
			}

			f.ip += 2

		default:
			return NIL, NewExecutionError("unknown instruction")
		}
	}
}
