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

// Opcodes
const (
	OPNOP uint8 = iota // do nothing

	OPLDC // load constant LDC (index uint8)
	OPLDA // load argument LDA (index uint8)

	OPINV // invoke function
	OPRET // return from function

	OPFIT // jump forward if truthy JIT (offset uint8)
	OPBIT // jump backward if truthy JIT (offset uint8)
	OPJUF // jump forward (offset uint8)
	OPJUB // jump backward (offset uint8)
)

// CodeChunk holds bytecode and provides facilities for reading and writing it
type CodeChunk struct {
	consts []Value
	code   []uint8
	length int
}

func NewCodeChunk() *CodeChunk {
	return &CodeChunk{
		consts: []Value{},
		code:   []uint8{},
		length: 0,
	}
}

func (c *CodeChunk) Append(insts ...uint8) {
	c.code = append(c.code, insts...)
	c.length = len(c.code)
}

func (c *CodeChunk) Get(idx int) uint8 {
	return c.code[idx]
}

const defaultStackSize = 32

// Frame is a single interpreter context
type Frame struct {
	stack   []Value
	args    []Value
	argc    int
	consts  []Value
	constsc int
	code    *CodeChunk
	ip      int
	sp      int
}

func NewFrame(code *CodeChunk, args []Value) *Frame {
	return &Frame{
		stack:   make([]Value, defaultStackSize),
		args:    args,
		argc:    len(args),
		consts:  code.consts,
		constsc: len(code.consts),
		code:    code,
		ip:      0,
		sp:      0,
	}
}

func (f *Frame) Push(v Value) error {
	if f.sp >= defaultStackSize-1 {
		return NewExecutionError("stack overflow")
	}
	f.stack[f.sp] = v
	f.sp++
	return nil
}

func (f *Frame) Pop() (Value, error) {
	if f.sp == 0 {
		return NIL, NewExecutionError("stack underflow")
	}
	f.sp--
	v := f.stack[f.sp]
	//f.stack[f.sp] = nil
	return v, nil
}

func (f *Frame) Nth(n int) (Value, error) {
	i := f.sp - 1 - n
	if i < 0 {
		return NIL, NewExecutionError("Nth: stack underflow")
	}
	return f.stack[i], nil
}

func (f *Frame) Mult(start int, count int) ([]Value, error) {
	if count < 0 {
		return nil, NewExecutionError("Mult: count 0 or negative")
	}
	i := f.sp - start
	if i-count < 0 {
		return nil, NewExecutionError("Mult: stack underflow")
	}
	return f.stack[i-count : i], nil
}

func (f *Frame) Drop(n int) error {
	top := f.sp - 1
	if top < 0 {
		return NewExecutionError("Drop: stack underflow")
	}
	f.sp -= n
	if f.sp < 0 {
		return NewExecutionError("Drop: stack underflow")
	}
	// for i := top; i >= f.sp; i-- {
	// 	f.stack[i] = nil
	// }
	return nil
}

func (f *Frame) Run() (Value, error) {
	for {
		inst := f.code.Get(f.ip)
		switch inst {
		case OPNOP:
			f.ip++
		case OPLDC:
			idx := int(f.code.Get(f.ip + 1))
			if idx >= f.constsc {
				return NIL, NewExecutionError("const lookup out of bounds")
			}
			err := f.Push(f.consts[idx])
			if err != nil {
				return NIL, NewExecutionError("const push failed").Wrap(err)
			}
			f.ip += 2
		case OPLDA:
			idx := int(f.code.Get(f.ip + 1))
			if idx >= f.argc {
				return NIL, NewExecutionError("argument lookup out of bounds")
			}
			err := f.Push(f.args[idx])
			if err != nil {
				return NIL, NewExecutionError("argument push failed").Wrap(err)
			}
			f.ip += 2
		case OPRET:
			v, err := f.Pop()
			if err != nil {
				return NIL, NewExecutionError("return failed").Wrap(err)
			}
			return v, nil
		case OPINV:
			fraw, err := f.Nth(0)
			if err != nil {
				return NIL, NewExecutionError("invoke instruction failed").Wrap(err)
			}
			fn, ok := fraw.(Fn)
			if !ok {
				return NIL, NewTypeError(fraw, "is not a function", nil)
			}
			arity := fn.Arity()
			a, err := f.Mult(1, arity)
			if err != nil {
				return NIL, NewExecutionError("popping arguments failed").Wrap(err)
			}
			out := fn.Invoke(a)
			err = f.Drop(arity + 1)
			if err != nil {
				return NIL, NewExecutionError("cleaning stack after call").Wrap(err)
			}
			err = f.Push(out)
			if err != nil {
				return NIL, NewExecutionError("pushing return value failed").Wrap(err)
			}
			f.ip++
		case OPFIT:
			offset := int(f.code.Get(f.ip + 1))
			v, err := f.Pop()
			if err != nil {
				return NIL, NewExecutionError("FIT pop condition").Wrap(err)
			}
			if v == NIL || v == TRUE {
				f.ip += 2
				continue
			}
			f.ip += offset
		case OPBIT:
			offset := int(f.code.Get(f.ip + 1))
			v, err := f.Pop()
			if err != nil {
				return NIL, NewExecutionError("BIT pop condition").Wrap(err)
			}
			if v == NIL || v == TRUE {
				f.ip += 2
				continue
			}
			f.ip -= offset

		default:
			return NIL, NewExecutionError("unknown instruction")
		}

	}
}
