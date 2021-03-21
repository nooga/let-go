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

package rt

import (
	_ "embed"
	"fmt"
	"github.com/nooga/let-go/pkg/vm"
	"strings"
)

var nsRegistry map[string]*vm.Namespace

func init() {
	nsRegistry = make(map[string]*vm.Namespace)

	installLangNS()
}

func NS(name string) *vm.Namespace {
	return LookupOrRegisterNS(name)
}

func RegisterNS(namespace *vm.Namespace) *vm.Namespace {
	nsRegistry[namespace.Name()] = namespace
	return namespace
}

func LookupOrRegisterNS(name string) *vm.Namespace {
	e := nsRegistry[name]
	if e != nil {
		return e
	}
	nsRegistry[name] = vm.NewNamespace(name)
	nsRegistry[name].Refer(CoreNS, "", true)
	return nsRegistry[name]
}

//go:embed core/core.lg
var CoreSrc string

const NameCoreNS = "core"

var CoreNS *vm.Namespace
var CurrentNS *vm.Var

//nolint
func installLangNS() {
	plus, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		n := 0
		for i := range vs {
			n += vs[i].Unbox().(int)
		}
		return vm.Int(n)
	})

	mul, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		n := 1
		for i := range vs {
			n *= vs[i].Unbox().(int)
		}
		return vm.Int(n)
	})

	sub, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) < 1 {
			// FIXME error out
			return vm.NIL
		}
		n := vs[0].Unbox().(int)
		if len(vs) == 1 {
			// FIXME error out
			return vm.Int(-n)
		}
		for i := 1; i < len(vs); i++ {
			n -= vs[i].Unbox().(int)
		}
		return vm.Int(n)
	})

	div, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		n := 0
		if len(vs) < 1 {
			// FIXME error out
			return vm.NIL
		}
		for i := range vs {
			n /= vs[i].Unbox().(int)
		}
		return vm.Int(n)
	})

	equals, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		length := len(vs)
		if length < 1 {
			// FIXME error out
			return vm.NIL
		}

		for i := 1; i < length; i++ {
			if vs[0] != vs[i] {
				return vm.FALSE
			}
		}
		return vm.TRUE
	})

	gt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			// FIXME error out
			return vm.NIL
		}
		ret, err := vm.BooleanType.Box(vs[0].Unbox().(int) > vs[1].Unbox().(int))
		if err != nil {
			// FIXME error out
			return vm.NIL
		}
		return ret
	})

	lt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			// FIXME error out
			return vm.NIL
		}
		ret, err := vm.BooleanType.Box(vs[0].Unbox().(int) < vs[1].Unbox().(int))
		if err != nil {
			// FIXME error out
			return vm.NIL
		}
		return ret
	})

	and, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		var ret vm.Value = vm.TRUE
		if len(vs) == 1 {
			return vs[0]
		}
		for i := range vs {
			if !vm.IsTruthy(vs[i]) {
				return vs[i]
			}
			ret = vs[i]
		}
		return ret
	})

	or, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		var ret vm.Value = vm.NIL
		if len(vs) == 1 {
			return vs[0]
		}
		for i := range vs {
			if vm.IsTruthy(vs[i]) {
				return vs[i]
			}
			ret = vs[i]
		}
		return ret
	})

	not, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME this is an error
			return vm.NIL
		}
		return vm.Boolean(!vm.IsTruthy(vs[0]))
	})

	setMacro, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME error out
			return vm.NIL
		}
		m := vs[0].(*vm.Var)
		m.SetMacro()
		return m
	})

	vector, err := vm.NativeFnType.Wrap(vm.NewArrayVector)
	list, err := vm.NativeFnType.Wrap(vm.NewList)
	hashMap, err := vm.NativeFnType.Wrap(vm.NewMap)

	assoc, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 3 {
			// FIXME error out
			return vm.NIL
		}
		seq, ok := vs[0].(vm.Associative)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.Assoc(vs[1], vs[2])
	})

	dissoc, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			// FIXME error out
			return vm.NIL
		}
		seq, ok := vs[0].(vm.Associative)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		key := vs[1]
		return seq.Dissoc(key)
	})

	cons, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			// FIXME error out
			return vm.NIL
		}
		elem := vs[0]
		seq, ok := vs[1].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.Cons(elem)
	})

	first, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME error out
			return vm.NIL
		}
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.First()
	})

	second, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME error out
			return vm.NIL
		}
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.Next().First()
	})

	next, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME error out
			return vm.NIL
		}
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}

		n := seq.Next()

		// FIXME move that to Seq.Next()
		if n.(vm.Collection).Count().(vm.Int) == 0 {
			return vm.NIL
		}
		return n
	})

	get, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		vl := len(vs)
		if vl < 2 || vl > 3 {
			// FIXME return error
			return vm.NIL
		}
		key := vs[1]
		as, ok := vs[0].(vm.Lookup)
		if !ok {
			// FIXME return error
			return vm.NIL
		}
		if vl == 2 {
			return as.ValueAt(key)
		}
		return as.ValueAtOr(key, vs[2])
	})

	count, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME error out
			return vm.NIL
		}
		seq, ok := vs[0].(vm.Collection)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.Count()
	})

	// FIXME write real ones later because this is naiiiiive
	mapf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			// FIXME error out
			return vm.NIL
		}
		mfn, ok := vs[0].(vm.Fn)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		seq, ok := vs[1].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		length := 0
		col, ok := vs[1].(vm.Collection)
		if ok {
			length = col.RawCount()
		}
		if length > 0 {
			newseq := make([]vm.Value, length)
			i := 0
			for seq != vm.EmptyList {
				newseq[i] = mfn.Invoke([]vm.Value{seq.First()})
				seq = seq.Next()
				i++
			}
			ret, _ := vm.ListType.Box(newseq)
			return ret
		}
		return vm.EmptyList
	})

	reduce, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) < 2 || len(vs) > 3 {
			// FIXME error out
			return vm.NIL
		}
		mfn, ok := vs[0].(vm.Fn)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		sidx := 1
		if len(vs) == 3 {
			sidx = 2
		}
		seq, ok := vs[sidx].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		var acc vm.Value
		if len(vs) == 3 {
			acc = vs[1]
		} else {
			acc = seq.First()
			seq = seq.Next()
		}
		for seq != vm.EmptyList {
			acc = mfn.Invoke([]vm.Value{seq.First(), acc})
			seq = seq.Next()
		}

		return acc
	})

	printlnf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		b := &strings.Builder{}
		for i := range vs {
			if i > 0 {
				b.WriteRune(' ')
			}
			if vs[i].Type() == vm.StringType {
				b.WriteString(string(vs[i].(vm.String)))
				continue
			}
			b.WriteString(vs[i].String())
		}
		fmt.Println(b)
		return vm.NIL
	})

	typef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			//FIXME this is an error
			return vm.NIL
		}
		t := vs[0].Type()
		if t == vm.NilType {
			return vm.NIL
		}
		return t
	})

	inNs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME handle error
			return vm.NIL
		}
		sym := vs[0]
		if sym.Type() != vm.SymbolType {
			// FIXME handle error
			return vm.NIL
		}
		nns := LookupOrRegisterNS(string(sym.(vm.Symbol)))
		CurrentNS.SetRoot(nns)
		return nns
	})

	use, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) < 1 {
			// FIXME handle error
			return vm.NIL
		}
		cns := CurrentNS.Deref().(*vm.Namespace)
		for i := range vs {
			s, ok := vs[i].(vm.Symbol)
			if !ok {
				return vm.NIL
			}
			cns.Refer(NS(string(s)), "", true)
		}
		return vm.NIL
	})

	if err != nil {
		panic("lang NS init failed")
	}

	ns := vm.NewNamespace(NameCoreNS)

	// vars
	CurrentNS = ns.Def("*ns*", ns)

	// FIXME implement the primitives in let-go later on and clean up this mess
	// primitive fns
	ns.Def("+", plus)
	ns.Def("*", mul)
	ns.Def("-", sub)
	ns.Def("/", div)

	ns.Def("=", equals)
	ns.Def("gt", gt)
	ns.Def("lt", lt)

	ns.Def("and", and)
	ns.Def("or", or)
	ns.Def("not", not)

	ns.Def("set-macro!", setMacro)
	ns.Def("in-ns", inNs)
	ns.Def("use", use)

	ns.Def("vector", vector)
	ns.Def("hash-map", hashMap)
	ns.Def("list", list)

	ns.Def("assoc", assoc)
	ns.Def("dissoc", dissoc)
	ns.Def("cons", cons)
	ns.Def("first", first)
	ns.Def("second", second)
	ns.Def("next", next)
	ns.Def("get", get)
	ns.Def("count", count)

	ns.Def("map", mapf)
	ns.Def("reduce", reduce)

	ns.Def("println", printlnf)

	ns.Def("type", typef)

	CoreNS = ns

	RegisterNS(ns)
}
