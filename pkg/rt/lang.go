/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	_ "embed"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

var nsRegistry map[string]*vm.Namespace

func init() {
	nsRegistry = make(map[string]*vm.Namespace)

	installLangNS()
}

func AllNSes() map[string]*vm.Namespace {
	return nsRegistry
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

var gensymID = 0

func nextID() int {
	gensymID++
	return gensymID
}

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

	mod, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			// FIXME error out
			return vm.NIL
		}
		ret := vs[0].Unbox().(int) % vs[1].Unbox().(int)
		return vm.Int(ret)
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

	gensym, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		prefix := "G__"
		if len(vs) == 1 {
			arg, ok := vs[0].(vm.String)
			if !ok {
				// FIXME :P
				return vm.NIL
			}
			prefix = string(arg)
		}
		return vm.Symbol(fmt.Sprintf("%s%d", prefix, nextID()))
	})

	vector, err := vm.NativeFnType.Wrap(vm.NewArrayVector)
	list, err := vm.NativeFnType.Wrap(vm.NewList)
	hashMap, err := vm.NativeFnType.Wrap(vm.NewMap)

	rangef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) == 0 {
			return vm.EmptyList
		}
		if len(vs) == 1 {
			return vm.NewRange(0, vs[0].(vm.Int), 1)
		}
		if len(vs) == 2 {
			return vm.NewRange(vs[0].(vm.Int), vs[1].(vm.Int), 1)
		}
		if len(vs) == 3 {
			return vm.NewRange(vs[0].(vm.Int), vs[1].(vm.Int), vs[2].(vm.Int))
		}
		// FIXME :P
		return vm.NIL
	})

	assoc, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) < 3 || len(vs)%2 == 0 {
			// FIXME error out
			return vm.NIL
		}
		coll, ok := vs[0].(vm.Associative)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		ret := coll
		for i := 1; i < len(vs); i += 2 {
			ret = ret.Assoc(vs[i], vs[i+1])
		}
		return ret
	})

	dissoc, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) < 2 {
			// FIXME error out
			return vm.NIL
		}
		coll, ok := vs[0].(vm.Associative)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		ret := coll
		for i := 1; i < len(vs); i++ {
			ret = ret.Dissoc(vs[i])
		}
		return ret
	})

	update, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) < 3 {
			// FIXME error out
			return vm.NIL
		}
		colla, ok := vs[0].(vm.Associative)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		collg, ok := vs[0].(vm.Lookup)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		key := vs[1]
		fn := vs[2].(vm.Fn)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		args := []vm.Value{collg.ValueAt(key)}
		if len(vs) > 3 {
			args = append(args, vs[3:]...)
		}
		return colla.Assoc(key, fn.Invoke(args))
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

	conj, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			// FIXME error out
			return vm.NIL
		}
		elem := vs[1]
		seq, ok := vs[0].(vm.Collection)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.Conj(elem)
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
		if len(vs) < 2 {
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
		// if we have a single coll
		if len(vs) == 2 {
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
		}
		// if we have more colls
		colls := vs[1:]
		seqs := make([]vm.Seq, len(colls))
		minlen := math.MaxInt
		for i := range colls {
			collx, ok := colls[i].(vm.Collection)
			if !ok {
				// FIXME error
				return vm.NIL
			}
			c := collx.RawCount()
			if c < minlen {
				minlen = c
			}
			seqs[i] = colls[i].(vm.Seq)
		}
		if minlen == 0 {
			return vm.EmptyList
		}
		newseq := make([]vm.Value, minlen)
		for i := 0; i < minlen; i++ {
			fargs := make([]vm.Value, len(seqs))
			for j := range seqs {
				fargs[j] = seqs[j].First()
				seqs[j] = seqs[j].Next()
			}
			newseq[i] = mfn.Invoke(fargs)
		}
		ret, _ := vm.ListType.Box(newseq)
		return ret
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
			acc = mfn.Invoke([]vm.Value{acc, seq.First()})
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

	apply, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			return vm.NIL
		}
		f, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL
		}
		switch vs[1].Type() {
		case vm.ArrayVectorType:
			return f.Invoke(vs[1].(vm.ArrayVector))
		case vm.ListType:
			args := vs[1].Unbox().([]vm.Value)
			return f.Invoke(args)
		}
		return vm.NIL
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

	now, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		return vm.NewBoxed(time.Now())
	})

	methodInvoke, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) < 2 {
			// FIXME handle error
			fmt.Println("not enough args")
			return vm.NIL
		}
		rec, ok := vs[0].(vm.Receiver)
		if !ok {
			// FIXME handle error
			fmt.Println("expected Receiver")
			return vm.NIL
		}
		name, ok := vs[1].(vm.Symbol)
		if !ok {
			// FIXME handle error
			fmt.Println("expected Symbol")
			return vm.NIL
		}
		return rec.InvokeMethod(name, vs[2:])
	})

	deref, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME err
			return vm.NIL
		}
		ref, ok := vs[0].(vm.Reference)
		if !ok {
			return vm.NIL
		}
		return ref.Deref()
	})

	concat, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		var ret []vm.Value
		for i := range vs {
			vseq, ok := vs[i].(vm.Seq)
			if !ok {
				return vm.NIL
			}
			for {
				e := vseq.First()
				ret = append(ret, e)
				vseq = vseq.Next()
				if vseq == vm.EmptyList {
					break
				}
			}
		}
		r, err := vm.ListType.Box(ret)
		if err != nil {
			return vm.NIL
		}
		return r
	})

	slurp, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME err
			return vm.NIL
		}
		filename, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL
		}
		data, err := os.ReadFile(string(filename))
		if err != nil {
			return vm.NIL
		}
		return vm.String(data)
	})

	spit, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 2 {
			// FIXME err
			return vm.NIL
		}
		filename, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL
		}
		contents, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL
		}
		err := os.WriteFile(string(filename), []byte(contents), 0644)
		if err != nil {
			return vm.NIL
		}
		return vm.NIL
	})

	name, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME err
			return vm.NIL
		}
		s, ok := vs[0].(vm.String)
		if ok {
			return s
		}
		named, ok := vs[0].(vm.Named)
		if !ok {
			return vm.NIL
		}
		return named.Name()
	})

	namespace, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		if len(vs) != 1 {
			// FIXME err
			return vm.NIL
		}
		named, ok := vs[0].(vm.Named)
		if !ok {
			return vm.NIL
		}
		return named.Namespace()
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
	ns.Def("mod", mod)

	ns.Def("and", and)
	ns.Def("or", or)
	ns.Def("not", not)

	ns.Def("set-macro!", setMacro)
	ns.Def("gensym", gensym)
	ns.Def("in-ns", inNs)
	ns.Def("use", use)
	ns.Def("name", name)
	ns.Def("namespace", namespace)

	ns.Def("vector", vector)
	ns.Def("hash-map", hashMap)
	ns.Def("list", list)
	ns.Def("range", rangef)

	ns.Def("assoc", assoc)
	ns.Def("dissoc", dissoc)
	ns.Def("update", update)
	ns.Def("cons", cons)
	ns.Def("conj", conj)
	ns.Def("first", first)
	ns.Def("second", second)
	ns.Def("next", next)
	ns.Def("get", get)
	ns.Def("count", count)

	ns.Def("map", mapf)
	ns.Def("reduce", reduce)
	ns.Def("concat", concat)

	ns.Def("println", printlnf)

	ns.Def("type", typef)

	ns.Def("apply", apply)
	ns.Def("deref", deref)

	// FIXME move this later outside the core
	ns.Def("now", now)

	ns.Def("slurp", slurp)
	ns.Def("spit", spit)

	// FIXME move this to VM later
	ns.Def(".", methodInvoke)

	CoreNS = ns

	RegisterNS(ns)
}
