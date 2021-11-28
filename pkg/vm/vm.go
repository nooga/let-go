/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
 * persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package vm

import (
	"encoding/binary"
	"fmt"
)

// Opcodes
const (
	OP_NOOP uint8 = iota // do nothing

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
)

func OpcodeToString(op uint8) string {
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
	}
	if int(op) < len(ops) {
		return ops[op]
	}
	return "???"
}

// CodeChunk holds bytecode and provides facilities for reading and writing it
type CodeChunk struct {
	maxStack int
	consts   *[]Value
	code     []uint8
	length   int
}

func NewCodeChunk(consts *[]Value) *CodeChunk {
	return &CodeChunk{
		consts: consts,
		code:   []uint8{},
		length: 0,
	}
}

func (c *CodeChunk) Debug() {
	//fmt.Println("consts:")
	consts := *c.consts
	//for i := range consts {
	//	fmt.Println("  [", i, "] =", consts[i])
	//}
	fmt.Println("code:")
	i := 0
	for i < len(c.code) {
		op, _ := c.Get(i)
		switch op {
		case OP_RECUR:
			arg, _ := c.Get32(i + 1)
			arg2, _ := c.Get32(i + 5)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg, arg2)
			i += 9
		case OP_LOAD_ARG, OP_BRANCH_TRUE, OP_BRANCH_FALSE, OP_JUMP, OP_POP_N, OP_DUP_NTH, OP_INVOKE, OP_LOAD_CLOSEDOVER, OP_RECUR_FN, OP_MAKE_CLOSURE:
			arg, _ := c.Get32(i + 1)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg)
			i += 5
		case OP_LOAD_CONST:
			arg, _ := c.Get32(i + 1)
			fmt.Println("  ", i, ":", OpcodeToString(op), arg, "<-", consts[arg])
			i += 5
		default:
			fmt.Println("  ", i, ":", OpcodeToString(op))
			i++
		}
	}
}

func (c *CodeChunk) Length() int {
	return c.length
}

func (c *CodeChunk) Append(insts ...uint8) {
	c.code = append(c.code, insts...)
	c.length = len(c.code)
}

func (c *CodeChunk) Append32(val int) {
	n := make([]uint8, 4)
	binary.LittleEndian.PutUint32(n, uint32(val))
	c.code = append(c.code, n...)
	c.length = len(c.code)
}

func (c *CodeChunk) AppendChunk(o *CodeChunk) {
	if o.maxStack > c.maxStack {
		c.maxStack = o.maxStack
	}
	c.code = append(c.code, o.code...)
	c.length += len(o.code)
}

func (c *CodeChunk) Get(idx int) (uint8, error) {
	if idx >= c.length {
		return 0, NewExecutionError("bytecode fetch out of bounds")
	}
	return c.code[idx], nil
}

func (c *CodeChunk) Get32(idx int) (int, error) {
	if idx >= c.length || idx+4 > c.length {
		return 0, NewExecutionError("bytecode wide fetch out of bounds")
	}
	return int(binary.LittleEndian.Uint32(c.code[idx:])), nil
}

func (c *CodeChunk) Update32(address int, value int) {
	binary.LittleEndian.PutUint32(c.code[address:address+4], uint32(value))
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
	consts      []Value
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
		consts:  *code.consts,
		constsc: len(*code.consts),
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
		return NewExecutionError("drop: stack underflow")
	}
	f.sp -= n
	if f.sp < 0 {
		f.stackDbg()
		return NewExecutionError("drop: stack underflow")
	}
	// for i := top; i >= f.sp; i-- {
	// 	f.stack[i] = nil
	// }
	return nil
}

func trunc(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[0:n] + " ..."
}

func (f *Frame) stackDbg() {
	//fmt.Println("IP = ", f.ip)
	//f.code.Debug()
	fmt.Printf("VM stack [%d/%d]:", f.sp, f.code.maxStack)
	for i := 0; i < f.sp; i++ {
		fmt.Print("   ", trunc(f.stack[i].String(), 32))

	}
	fmt.Println()
}

func (f *Frame) Run() (Value, error) {
	if f.debug {
		fmt.Print("run")
		f.code.Debug()
	}
	for {
		inst, _ := f.code.Get(f.ip)
		//if f.debug {
		//	fmt.Println("exec", f.ip, OpcodeToString(inst))
		//	f.stackDbg()
		//}
		switch inst {
		case OP_NOOP:
			f.ip++

		case OP_LOAD_CONST:
			idx, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("const push failed").Wrap(err)
			}
			if idx >= f.constsc {
				return NIL, NewExecutionError("const lookup out of bounds")
			}
			err = f.push(f.consts[idx])
			if err != nil {
				return NIL, NewExecutionError("const push failed").Wrap(err)
			}
			f.ip += 5

		case OP_LOAD_ARG:
			idx, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("get argument index failed").Wrap(err)
			}
			if idx >= f.argc {
				return NIL, NewExecutionError("argument lookup out of bounds")
			}
			err = f.push(f.args[idx])
			if err != nil {
				return NIL, NewExecutionError("argument push failed").Wrap(err)
			}
			f.ip += 5

		case OP_RETURN:
			v, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("return failed").Wrap(err)
			}
			return v, nil

		case OP_INVOKE:
			arity, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("INV arg count").Wrap(err)
			}
			fraw, err := f.nth(arity)
			if err != nil {
				return NIL, NewExecutionError("invoke instruction failed").Wrap(err)
			}
			fn, ok := fraw.(Fn)
			if !ok {
				return NIL, NewTypeError(fraw, "is not a function", nil)
			}
			a, err := f.mult(0, arity)
			if err != nil {
				return NIL, NewExecutionError("popping arguments failed").Wrap(err)
			}
			args := make([]Value, len(a))
			copy(args, a)
			out := fn.Invoke(args)
			err = f.drop(arity + 1)
			if err != nil {
				return NIL, NewExecutionError("cleaning stack after call").Wrap(err)
			}
			err = f.push(out)
			if err != nil {
				return NIL, NewExecutionError("pushing return value failed").Wrap(err)
			}
			f.ip += 5

		case OP_BRANCH_TRUE:
			offset, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("BRT offset").Wrap(err)
			}
			v, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("BRT pop condition").Wrap(err)
			}
			if !IsTruthy(v) {
				f.ip += 5
				continue
			}
			f.ip += offset

		case OP_BRANCH_FALSE:
			offset, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("BRT offset").Wrap(err)
			}
			v, err := f.pop()
			if err != nil {
				return NIL, NewExecutionError("BRT pop condition").Wrap(err)
			}
			if IsTruthy(v) {
				f.ip += 5
				continue
			}
			f.ip += offset

		case OP_JUMP:
			offset, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("JMP offset").Wrap(err)
			}
			f.ip += offset

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
			num, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("PON get argument").Wrap(err)
			}
			err = f.drop(num)
			if err != nil {
				return NIL, NewExecutionError("PON drop").Wrap(err)
			}
			err = f.push(v)
			if err != nil {
				return NIL, NewExecutionError("PON push").Wrap(err)
			}
			f.ip += 5

		case OP_DUP_NTH:
			num, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("DPN get argument").Wrap(err)
			}
			val, err := f.nth(num)
			if err != nil {
				return NIL, NewExecutionError("DPN get nth").Wrap(err)
			}
			err = f.push(val)
			if err != nil {
				return NIL, NewExecutionError("DPN push").Wrap(err)
			}
			f.ip += 5

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
			idx, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("const push failed").Wrap(err)
			}
			if idx >= f.constsc {
				return NIL, NewExecutionError("const lookup out of bounds")
			}
			err = f.push(f.consts[idx].(*Var).Deref())
			if err != nil {
				return NIL, NewExecutionError("const push failed").Wrap(err)
			}
			f.ip += 5

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
			idx, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("get closed over index failed").Wrap(err)
			}
			// FIXME cache closedOvers count
			if idx >= len(f.closedOvers) {
				return NIL, NewExecutionError("closed over lookup out of bounds")
			}
			err = f.push(f.closedOvers[idx])
			if err != nil {
				return NIL, NewExecutionError("closed over push failed").Wrap(err)
			}
			f.ip += 5

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
			arity, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("REF arg count").Wrap(err)
			}
			a, err := f.mult(0, arity)
			if err != nil {
				return NIL, NewExecutionError("popping arguments failed").Wrap(err)
			}
			copy(f.args, a)
			f.argc = arity
			f.sp = 0
			f.ip = 0

		case OP_RECUR:
			offset, err := f.code.Get32(f.ip + 1)
			if err != nil {
				return NIL, NewExecutionError("REC reading offset").Wrap(err)
			}
			argc, err := f.code.Get32(f.ip + 5)
			if err != nil {
				return NIL, NewExecutionError("REC reading argc").Wrap(err)
			}
			a, err := f.mult(0, argc)
			if err != nil {
				return NIL, NewExecutionError("REC popping arguments failed").Wrap(err)
			}
			err = f.drop(argc * 2)
			if err != nil {
				return NIL, NewExecutionError("REC popping old locals").Wrap(err)
			}
			err = f.pushMult(a)
			if err != nil {
				return NIL, NewExecutionError("REC pushing new locals").Wrap(err)
			}

			f.ip -= offset

		default:
			return NIL, NewExecutionError("unknown instruction")
		}
		if f.debug {
			fmt.Println("exec", f.ip, OpcodeToString(inst))
			f.stackDbg()
		}
	}
}
