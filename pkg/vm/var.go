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

import "fmt"

type Var struct {
	root  Value
	nsref *Namespace
	ns    string
	name  string
}

func (v *Var) Invoke(values []Value) Value {
	f, ok := v.root.(Fn)
	if !ok {
		return NIL // FIXME this should be an error
	}
	return f.Invoke(values)
}

func (v *Var) Arity() int {
	f, ok := v.root.(Fn)
	if !ok {
		return 0 // FIXME this should be an error
	}
	return f.Arity()
}

func NewVar(nsref *Namespace, ns string, name string) *Var {
	return &Var{
		nsref: nsref,
		ns:    ns,
		name:  name,
	}
}

func (v *Var) SetRoot(val Value) *Var {
	v.root = val
	return v
}

func (v *Var) Deref() Value {
	return v.root
}

func (v *Var) Type() ValueType {
	return v.Deref().Type()
}

func (v *Var) Unbox() interface{} {
	return v.Deref().Unbox()
}

func (v *Var) String() string {
	return fmt.Sprintf("#'%s/%s", v.ns, v.name)
}
