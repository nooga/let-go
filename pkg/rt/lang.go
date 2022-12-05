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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

var nsRegistry map[string]*vm.Namespace

func init() {
	nsRegistry = make(map[string]*vm.Namespace)

	installLangNS()
	installHttpNS()
}

func AllNSes() map[string]*vm.Namespace {
	return nsRegistry
}

func FuzzyNamespacedSymbolLookup(currentNS *vm.Namespace, s vm.Symbol) []vm.Symbol {
	sns := s.Namespace()
	var ns *vm.Namespace
	if sns != vm.NIL {
		ns = nsRegistry[string(sns.(vm.String))]
	} else {
		ns = currentNS
	}
	name := s.Name()
	return vm.FuzzySymbolLookup(ns, vm.Symbol(name.(vm.String)))
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

// nolint
func installLangNS() {
	plus, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		n := 0
		for i := range vs {
			n += vs[i].Unbox().(int)
		}
		return vm.Int(n), nil
	})

	mul, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		n := 1
		for i := range vs {
			n *= vs[i].Unbox().(int)
		}
		return vm.Int(n), nil
	})

	sub, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		n := vs[0].Unbox().(int)
		if len(vs) == 1 {
			return vm.Int(-n), nil
		}
		for i := 1; i < len(vs); i++ {
			n -= vs[i].Unbox().(int)
		}
		return vm.Int(n), nil
	})

	div, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		n := vs[0].Unbox().(int)
		for _, e := range vs[1:] {
			n /= e.Unbox().(int)
		}
		return vm.Int(n), nil
	})

	equals, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		length := len(vs)
		if length < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}

		for i := 1; i < length; i++ {
			if vs[0] != vs[i] {
				return vm.FALSE, nil
			}
		}
		return vm.TRUE, nil
	})

	gt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		ret, err := vm.BooleanType.Box(vs[0].Unbox().(int) > vs[1].Unbox().(int))
		if err != nil {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return ret, nil
	})

	lt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		ret, err := vm.BooleanType.Box(vs[0].Unbox().(int) < vs[1].Unbox().(int))
		if err != nil {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return ret, nil
	})

	mod, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		ret := vs[0].Unbox().(int) % vs[1].Unbox().(int)
		return vm.Int(ret), nil
	})

	and, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		var ret vm.Value = vm.TRUE
		if len(vs) == 1 {
			return vs[0], nil
		}
		for i := range vs {
			if !vm.IsTruthy(vs[i]) {
				return vs[i], nil
			}
			ret = vs[i]
		}
		return ret, nil
	})

	or, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		var ret vm.Value = vm.NIL
		if len(vs) == 1 {
			return vs[0], nil
		}
		for i := range vs {
			if vm.IsTruthy(vs[i]) {
				return vs[i], nil
			}
			ret = vs[i]
		}
		return ret, nil
	})

	not, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.Boolean(!vm.IsTruthy(vs[0])), nil
	})

	setMacro, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		m := vs[0].(*vm.Var)
		m.SetMacro()
		return m, nil
	})

	gensym, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		prefix := "G__"
		if len(vs) == 1 {
			arg, ok := vs[0].(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf("gensym expected String")
			}
			prefix = string(arg)
		}
		return vm.Symbol(fmt.Sprintf("%s%d", prefix, nextID())), nil
	})

	vector, err := vm.NativeFnType.WrapNoErr(vm.NewArrayVector)
	list, err := vm.NativeFnType.WrapNoErr(vm.NewList)
	hashMap, err := vm.NativeFnType.WrapNoErr(vm.NewMap)
	hashSet, err := vm.NativeFnType.WrapNoErr(vm.NewSet)

	vec, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL || vs[0] == vm.EmptyList {
			return vm.ArrayVector{}, nil
		}

		if v, ok := vs[0].(vm.ArrayVector); ok {
			return v, nil
		}
		if seq, ok := vs[0].(vm.Seq); ok {
			ret := []vm.Value{}
			for seq != vm.EmptyList {
				ret = append(ret, seq.First())
				seq = seq.Next()
			}
			return vm.NewArrayVector(ret), nil
		}
		return vm.NIL, nil // FIXME is this an error?
	})

	rangef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.EmptyList, nil
		}
		if len(vs) == 1 {
			return vm.NewRange(0, vs[0].(vm.Int), 1), nil
		}
		if len(vs) == 2 {
			return vm.NewRange(vs[0].(vm.Int), vs[1].(vm.Int), 1), nil
		}
		if len(vs) == 3 {
			return vm.NewRange(vs[0].(vm.Int), vs[1].(vm.Int), vs[2].(vm.Int)), nil
		}
		return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
	})

	keyword, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if k, ok := vs[0].(vm.Keyword); ok {
			return k, nil
		}
		if k, ok := vs[0].(vm.Symbol); ok {
			return vm.Keyword(k), nil
		}
		if k, ok := vs[0].(vm.String); ok {
			return vm.Keyword(k), nil
		}
		// TODO handle namespaces
		return vm.NIL, nil // TODO is this an error?
	})

	assoc, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 3 || len(vs)%2 == 0 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		coll, ok := vs[0].(vm.Associative)
		if !ok {
			return vm.NIL, fmt.Errorf("assoc expected Associative")
		}
		ret := coll
		for i := 1; i < len(vs); i += 2 {
			ret = ret.Assoc(vs[i], vs[i+1])
		}
		return ret, nil
	})

	dissoc, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		coll, ok := vs[0].(vm.Associative)
		if !ok {
			return vm.NIL, fmt.Errorf("dissoc expected Associative")
		}
		ret := coll
		for i := 1; i < len(vs); i++ {
			ret = ret.Dissoc(vs[i])
		}
		return ret, nil
	})

	update, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		colla, ok := vs[0].(vm.Associative)
		if !ok {
			return vm.NIL, fmt.Errorf("update expected Associative")
		}
		collg, ok := vs[0].(vm.Lookup)
		if !ok {
			return vm.NIL, fmt.Errorf("update expected Lookup")
		}
		key := vs[1]
		fn := vs[2].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("update expected Fn")
		}
		args := []vm.Value{collg.ValueAt(key)}
		if len(vs) > 3 {
			args = append(args, vs[3:]...)
		}
		v, err := fn.Invoke(args)
		if err != nil {
			return vm.NIL, err // FIXME wrap this?
		}
		return colla.Assoc(key, v), nil
	})

	cons, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		elem := vs[0]
		seq, ok := vs[1].(vm.Seq)
		if !ok {
			return vm.NIL, fmt.Errorf("cons expected Seq")
		}
		return seq.Cons(elem), nil
	})

	conj, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		seq, ok := vs[0].(vm.Collection)
		if !ok {
			return vm.NIL, fmt.Errorf("conj expected Collection")
		}
		for i := 1; i < len(vs); i++ {
			seq = seq.Conj(vs[i])
		}
		return seq, nil
	})

	disj, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		s, ok := vs[0].(vm.Set)
		if !ok {
			return vm.NIL, fmt.Errorf("conj expected Set")
		}

		for _, v := range vs[1:] {
			s = s.Disj(v)
		}
		return s, nil
	})

	contains, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		s, ok := vs[0].(vm.Keyed)
		if !ok {
			return vm.NIL, fmt.Errorf("contains? expected Set")
		}

		return s.Contains(vs[1]), nil
	})

	first, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			return vm.NIL, fmt.Errorf("first expected Seq")
		}
		return seq.First(), nil
	})

	second, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			return vm.NIL, fmt.Errorf("second expected Seq")
		}
		return seq.Next().First(), nil
	})

	next, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			return vm.NIL, fmt.Errorf("next expected Seq")
		}

		n := seq.Next()

		// FIXME move that to Seq.Next()
		if n.(vm.Collection).Count().(vm.Int) == 0 {
			return vm.NIL, nil
		}
		return n, nil
	})

	seq, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		seq, ok := vs[0].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("seq expected Seqauble")
		}
		n := seq.Seq()
		return n, nil
	})

	get, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		vl := len(vs)
		if vl < 2 || vl > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		key := vs[1]
		as, ok := vs[0].(vm.Lookup)
		if !ok {
			return vm.NIL, nil // this is not an error
		}
		if vl == 2 {
			return as.ValueAt(key), nil
		}
		return as.ValueAtOr(key, vs[2]), nil
	})

	count, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		seq, ok := vs[0].(vm.Collection)
		if !ok {
			return vm.NIL, fmt.Errorf("count expected Collection")
		}
		return seq.Count(), nil
	})

	// FIXME write real ones later because this is naiiiiive
	mapf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		mfn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("map expected Fn")
		}
		sequable, ok := vs[1].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("map expected Sequable")
		}
		seq := sequable.Seq()
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
					newseq[i], err = mfn.Invoke([]vm.Value{seq.First()})
					if err != nil {
						return vm.NIL, err // FIXME wrap this?
					}
					seq = seq.Next()
					i++
				}
				ret, _ := vm.ListType.Box(newseq)
				return ret, nil
			}
			return vm.EmptyList, nil
		}
		// if we have more colls
		colls := vs[1:]
		seqs := make([]vm.Seq, len(colls))
		minlen := math.MaxInt
		for i := range colls {
			collx, ok := colls[i].(vm.Collection)
			if !ok {
				return vm.NIL, fmt.Errorf("map expected Collection")
			}
			c := collx.RawCount()
			if c < minlen {
				minlen = c
			}
			seqs[i] = colls[i].(vm.Seq)
		}
		if minlen == 0 {
			return vm.EmptyList, nil
		}
		newseq := make([]vm.Value, minlen)
		for i := 0; i < minlen; i++ {
			fargs := make([]vm.Value, len(seqs))
			for j := range seqs {
				fargs[j] = seqs[j].First()
				seqs[j] = seqs[j].Next()
			}
			newseq[i], err = mfn.Invoke(fargs)
			if err != nil {
				return vm.NIL, err // TODO wrap this somehow
			}
		}
		return vm.ListType.Box(newseq)
	})

	mapv, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		v, err := mapf.(vm.Fn).Invoke(vs)
		if err != nil {
			return vm.NIL, err // TODO wrap this
		}
		return vec.(vm.Fn).Invoke([]vm.Value{v})
	})

	reduce, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		mfn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("reduce expected Fn")
		}
		sidx := 1
		if len(vs) == 3 {
			sidx = 2
		}
		seq, ok := vs[sidx].(vm.Seq)
		if !ok {
			return vm.NIL, fmt.Errorf("reduce expected Seq")
		}
		var acc vm.Value
		if len(vs) == 3 {
			acc = vs[1]
		} else {
			acc = seq.First()
			seq = seq.Next()
		}
		for seq != vm.EmptyList {
			acc, err = mfn.Invoke([]vm.Value{acc, seq.First()})
			if err != nil {
				return vm.NIL, err // FIXME wrap this somehow
			}
			seq = seq.Next()
		}

		return acc, nil
	})

	some, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		f, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("some expected Fn")
		}
		seq, ok := vs[1].(vm.Seq)
		if !ok {
			return vm.NIL, fmt.Errorf("some expected Seq")
		}
		for seq != vm.EmptyList {
			v, err := f.Invoke([]vm.Value{seq.First()})
			if err != nil {
				return vm.NIL, err
			}
			if vm.IsTruthy(v) {
				return v, nil
			}
			seq = seq.Next()
		}

		return vm.FALSE, nil
	})

	printlnf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		b := &strings.Builder{}
		for i := range vs {
			if i > 0 {
				b.WriteRune(' ')
			}
			if vs[i].Type() == vm.StringType {
				b.WriteString(string(vs[i].(vm.String)))
				continue
			} else if vs[i].Type() == vm.CharType {
				b.WriteRune(rune(vs[i].(vm.Char)))
				continue
			}
			b.WriteString(vs[i].String())
		}
		fmt.Println(b)
		return vm.NIL, nil
	})

	str, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		b := &strings.Builder{}
		for i := range vs {
			if vs[i].Type() == vm.StringType {
				b.WriteString(string(vs[i].(vm.String)))
				continue
			} else if vs[i].Type() == vm.CharType {
				b.WriteRune(rune(vs[i].(vm.Char)))
				continue
			}
			b.WriteString(vs[i].String())
		}
		return vm.String(b.String()), nil
	})

	typef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		t := vs[0].Type()
		if t == vm.NilType {
			return vm.NIL, nil
		}
		return t, nil
	})

	apply, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		f, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("apply expected Fn")
		}
		switch vs[1].Type() {
		case vm.ArrayVectorType:
			return f.Invoke(vs[1].(vm.ArrayVector))
		case vm.ListType:
			args := vs[1].Unbox().([]vm.Value)
			return f.Invoke(args)
		}
		return vm.NIL, nil // FIXME is this an error?
	})

	inNs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		sym := vs[0]
		if sym.Type() != vm.SymbolType {
			return vm.NIL, fmt.Errorf("in-ns expected Symbol")
		}
		nns := LookupOrRegisterNS(string(sym.(vm.Symbol)))
		CurrentNS.SetRoot(nns)
		return nns, nil
	})

	use, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		cns := CurrentNS.Deref().(*vm.Namespace)
		for i := range vs {
			s, ok := vs[i].(vm.Symbol)
			if !ok {
				return vm.NIL, fmt.Errorf("use expected Symbol")
			}
			cns.Refer(NS(string(s)), "", true)
		}
		return vm.NIL, nil
	})

	now, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NewBoxed(time.Now()), nil
	})

	methodInvoke, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		rec, ok := vs[0].(vm.Receiver)
		if !ok {
			return vm.NIL, fmt.Errorf("method-invoke expected Receiver")
		}
		name, ok := vs[1].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("method-invoke expected Symbol")
		}
		return rec.InvokeMethod(name, vs[2:])
	})

	deref, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		ref, ok := vs[0].(vm.Reference)
		if !ok {
			return vm.NIL, fmt.Errorf("deref expected Reference")
		}
		return ref.Deref(), nil
	})

	concat, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		var ret []vm.Value
		for i := range vs {
			vseq, ok := vs[i].(vm.Seq)
			if !ok {
				return vm.NIL, fmt.Errorf("concat expected Seq")
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
			return vm.NIL, fmt.Errorf("concat failed: %w", err)
		}
		return r, nil
	})

	slurp, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		filename, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("slurp expected String")
		}
		data, err := os.ReadFile(string(filename))
		if err != nil {
			return vm.NIL, fmt.Errorf("slurp failed: %w", err)
		}
		return vm.String(data), nil
	})

	spit, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		filename, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("spit expected String")
		}
		contents, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("spit expected String")
		}
		err := os.WriteFile(string(filename), []byte(contents), 0644)
		if err != nil {
			return vm.NIL, fmt.Errorf("spit failed: %w", err)
		}
		return vm.NIL, nil
	})

	name, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if ok {
			return s, nil
		}
		named, ok := vs[0].(vm.Named)
		if !ok {
			return vm.NIL, fmt.Errorf("name expected Named")
		}
		return named.Name(), nil
	})

	namespace, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		named, ok := vs[0].(vm.Named)
		if !ok {
			return vm.NIL, fmt.Errorf("namespace expected Named")
		}
		return named.Namespace(), nil
	})

	atom, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.NewAtom(vs[0]), nil
	})

	// (swap! a fn)
	swap, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		at, ok := vs[0].(*vm.Atom)
		if !ok {
			return vm.NIL, fmt.Errorf("swap expected Atom")
		}
		fn, ok := vs[1].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("swap expected Fn")
		}
		return at.Swap(fn, vs[2:])
	})

	// (reset! a fn)
	reset, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		at, ok := vs[0].(*vm.Atom)
		if !ok {
			return vm.NIL, fmt.Errorf("reset expected Atom")
		}
		return at.Reset(vs[1]), nil
	})

	gof, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		at, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("go expected Fn")
		}
		ret := make(vm.Chan)
		go func() {
			v, err := at.Invoke(nil)
			if err != nil {
				fmt.Println(err) // FIXME handle this properly
			}
			ret <- v
			close(ret)
		}()
		return ret, nil
	})

	chanf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return make(vm.Chan), nil
	})

	chanput, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		ch, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf(">! expected Chan")
		}
		if vs[1] == vm.NIL {
			return vm.NIL, fmt.Errorf(">! can't put nil on chan")
		}
		ch <- vs[1]
		return vm.TRUE, nil
	})

	changet, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		ch, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("<! expected Chan")
		}
		v, ok := <-ch
		if !ok {
			return vm.NIL, nil // this is not an error
		}
		return v, nil
	})

	lines, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("lines expected String")
		}
		ss := strings.Split(string(s), "\n")
		av := make([]vm.Value, len(ss))
		for i := range ss {
			av[i] = vm.String(ss[i])
		}
		return vm.ArrayVector(av), nil
	})

	parseInt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("parse-int expected String")
		}
		i, err := strconv.Atoi(string(s))
		return vm.Int(i), err
	})

	max, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		fst, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("max expected Int")
		}
		max := vm.Int(fst)
		for i := 1; i < len(vs); i++ {
			if n, ok := vs[i].(vm.Int); ok {
				if n > max {
					max = n
				}
			} else {
				return vm.NIL, fmt.Errorf("max expected Int")
			}
		}
		return max, err
	})

	min, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		fst, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("min expected Int")
		}
		min := vm.Int(fst)
		for i := 1; i < len(vs); i++ {
			if n, ok := vs[i].(vm.Int); ok {
				if n < min {
					min = n
				}
			} else {
				return vm.NIL, fmt.Errorf("min expected Int")
			}
		}
		return min, err
	})

	sort, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		var comp vm.Fn
		var coll vm.Collection
		var ok bool
		if len(vs) == 2 {
			comp, ok = vs[0].(vm.Fn)
			if !ok {
				return vm.NIL, fmt.Errorf("sort expected a comparator function")
			}
			coll, ok = vs[1].(vm.Collection)
			if !ok {
				return vm.NIL, fmt.Errorf("sort expected a Collection")
			}
		} else {
			comp = lt.(vm.Fn)
			coll, ok = vs[0].(vm.Collection)
			if !ok {
				return vm.NIL, fmt.Errorf("sort expected a Collection")
			}
		}
		temp := make([]vm.Value, coll.RawCount())
		seq := coll.(vm.Sequable).Seq()
		for i := range temp {
			temp[i] = seq.First()
			seq = seq.Next()
		}
		var err error
		sort.SliceStable(temp, func(i, j int) bool {
			var b vm.Value
			b, err = comp.Invoke([]vm.Value{temp[i], temp[j]})
			return vm.IsTruthy(b)
		})
		if err != nil {
			return vm.NIL, fmt.Errorf("sort expected a Collection")
		}

		return vm.ListType.Box(temp)
	})

	split, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("split expected String")
		}
		delim := ""
		if len(vs) == 2 {
			delimv := vs[1].(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf("split expected String")
			}
			delim = string(delimv)
		}
		frags := strings.Split(string(s), delim)
		var ret vm.Seq = vm.EmptyList
		l := len(frags)
		for i := range frags {
			ret = ret.Cons(vm.String(frags[l-i-1]))
		}
		return ret, nil
	})

	intf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if i, ok := vs[0].(vm.Int); ok {
			return i, nil
		}
		if i, ok := vs[0].(vm.Char); ok {
			return vm.Int(int(i)), nil
		}
		return vm.NIL, fmt.Errorf("%s can't be coerced to int", vs[0])
	})

	regex, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if s, ok := vs[0].(vm.String); ok {
			return vm.NewRegex(string(s))
		}
		return vm.NIL, fmt.Errorf("regex expected String")
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
	ns.Def("vec", vec)
	ns.Def("hash-map", hashMap)
	ns.Def("list", list)
	ns.Def("range", rangef)
	ns.Def("keyword", keyword)
	ns.Def("hash-set", hashSet)

	ns.Def("seq", seq)

	ns.Def("assoc", assoc)
	ns.Def("dissoc", dissoc)
	ns.Def("update", update)
	ns.Def("cons", cons)
	ns.Def("conj", conj)
	ns.Def("disj", disj)
	ns.Def("first", first)
	ns.Def("second", second)
	ns.Def("next", next)
	ns.Def("get", get)
	ns.Def("count", count)
	ns.Def("contains?", contains)

	ns.Def("map", mapf)
	ns.Def("mapv", mapv)
	ns.Def("reduce", reduce)
	ns.Def("concat", concat)
	ns.Def("some", some)

	ns.Def("println", printlnf)

	ns.Def("type", typef)

	ns.Def("apply", apply)
	ns.Def("deref", deref)

	ns.Def("atom", atom)
	ns.Def("reset!", reset)
	ns.Def("swap!", swap)

	// FIXME move this later outside the core
	ns.Def("now", now)

	ns.Def("slurp", slurp)
	ns.Def("spit", spit)
	ns.Def("lines", lines)

	ns.Def("parse-int", parseInt)
	ns.Def("max", max)
	ns.Def("min", min)

	ns.Def("sort", sort)

	// FIXME move this to VM later
	ns.Def(".", methodInvoke)

	// FIXME move to async
	ns.Def("go*", gof)
	ns.Def("chan", chanf)
	ns.Def(">!", chanput)
	ns.Def("<!", changet)

	ns.Def("int", intf)

	ns.Def("str", str)
	ns.Def("split", split)
	ns.Def("regex", regex)

	CoreNS = ns

	RegisterNS(ns)
}
