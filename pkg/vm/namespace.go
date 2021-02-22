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

type Namespace struct {
	name     string
	registry map[Symbol]*Var
}

func NewNamespace(name string) *Namespace {
	return &Namespace{
		name:     name,
		registry: map[Symbol]*Var{},
	}
}

func (n *Namespace) Def(name string, val Value) *Var {
	s := Symbol(name)
	va := NewVar(n, n.name, name)
	va.SetRoot(val)
	n.registry[s] = va
	return va
}

func (n *Namespace) LookupOrAdd(symbol Symbol) Value {
	val, ok := n.registry[symbol]
	if !ok {
		return n.Def(string(symbol), NIL)
	}
	return val
}
