/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	_ "embed"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

var nsRegistry map[string]*vm.Namespace

type NSLoader interface {
	Load(string) *vm.Namespace
}

var nsLoader NSLoader

func SetNSLoader(loader NSLoader) {
	nsLoader = loader
}

func init() {
	nsRegistry = make(map[string]*vm.Namespace)

	// Register global namespace lookup so qualified symbols (foo/x) work
	vm.SetNSLookup(func(name string) *vm.Namespace {
		return nsRegistry[name]
	})

	// Wire up ValueEquals for OP_EQ fast path in the VM
	vm.SetValueEquals(func(a, b vm.Value) bool {
		return valueEquals(a, b)
	})

	initTypeMappings()
	installLangNS()
	installHttpNS()
	installOsNS()
	installJSONNS()
	installIoNS()
	installAsyncNS()
	installTransitNS()
	installPodsNS()
	// walk namespace is embedded via WalkSrc and will be loaded on demand
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
	return vm.FuzzySymbolLookup(ns, vm.Symbol(name.(vm.String)), true)
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
	if nsLoader != nil {
		n := nsLoader.Load(name)
		if n != nil {
			nsRegistry[name] = n
			return n
		}
	}
	// Check if loading side-effected the registry (in-ns during load creates the ns)
	if e := nsRegistry[name]; e != nil {
		return e
	}
	nsRegistry[name] = vm.NewNamespace(name)
	nsRegistry[name].Refer(CoreNS, "", true)
	return nsRegistry[name]
}

func LookupOrRegisterNSNoLoad(name string) *vm.Namespace {
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

// valueEquals performs deep equality comparison for Clojure semantics
func valueEquals(a, b vm.Value) bool {
	// Handle nil
	if a == vm.NIL && b == vm.NIL {
		return true
	}
	if a == vm.NIL || b == vm.NIL {
		return false
	}

	// Allow cross-type comparison for numbers and vectors
	if a.Type() != b.Type() {
		if vm.IsNumber(a) && vm.IsNumber(b) {
			return vm.NumEq(a, b)
		}
		// Cross-type vector equality (ArrayVector vs PersistentVector)
		if eq, ok := a.(interface{ Equals(vm.Value) bool }); ok {
			return eq.Equals(b)
		}
		if eq, ok := b.(interface{ Equals(vm.Value) bool }); ok {
			return eq.Equals(a)
		}
		return false
	}

	// Handle collections specially
	switch av := a.(type) {
	case vm.ArrayVector:
		bv, ok := b.(vm.ArrayVector)
		if !ok {
			// Could be PersistentVector — use Equals
			if eq, ok2 := a.(interface{ Equals(vm.Value) bool }); ok2 {
				return eq.Equals(b)
			}
			return false
		}
		if len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !valueEquals(av[i], bv[i]) {
				return false
			}
		}
		return true
	case *vm.List:
		// b could be any Seq-like type (List, Cons, ArrayVectorSeq, etc.)
		bs, ok := b.(vm.Seq)
		if !ok {
			return false
		}
		as := vm.Seq(av)
		for as != nil && bs != nil {
			if !valueEquals(as.First(), bs.First()) {
				return false
			}
			as, bs = as.Next(), bs.Next()
		}
		return as == nil && bs == nil
	case vm.Map:
		bm := b.(vm.Map)
		if len(av) != len(bm) {
			return false
		}
		for k, v := range av {
			bv, ok := bm[k]
			if !ok || !valueEquals(v, bv) {
				return false
			}
		}
		return true
	case *vm.PersistentMap:
		if bm, ok := b.(*vm.PersistentMap); ok {
			return av.Equals(bm)
		}
		return false
	case vm.Set:
		bs := b.(vm.Set)
		if len(av) != len(bs) {
			return false
		}
		for k := range av {
			if _, ok := bs[k]; !ok {
				return false
			}
		}
		return true
	case *vm.PersistentSet:
		bs, ok := b.(*vm.PersistentSet)
		if !ok {
			return false
		}
		if av.RawCount() != bs.RawCount() {
			return false
		}
		// Check every element of av is in bs
		seq := av.Seq()
		for seq != nil && seq != vm.EmptyList {
			if bs.Contains(seq.First()) == vm.FALSE {
				return false
			}
			seq = seq.Next()
		}
		return true
	case *vm.BigInt:
		if bv, ok := b.(*vm.BigInt); ok {
			return av.Equals(bv)
		}
		return false
	default:
		// Try Equals interface for types that implement it (PersistentVector, etc.)
		if eq, ok := a.(interface{ Equals(vm.Value) bool }); ok {
			return eq.Equals(b)
		}
		return a == b
	}
}

func seqOf(v vm.Value) (vm.Seq, error) {
	if v == vm.NIL {
		return nil, nil
	}
	if v == vm.EmptyList {
		return nil, nil
	}
	// For concrete collections (not LazySeq), prefer Sequable.Seq() which
	// produces a stable seq view (e.g. MapSeq with cached entries).
	// For LazySeq, return it directly — callers that iterate must
	// handle empty lazy seqs by checking First()/Next() properly.
	if _, isLazy := v.(*vm.LazySeq); !isLazy {
		if sq, ok := v.(vm.Sequable); ok {
			s := sq.Seq()
			if s == nil || s == vm.EmptyList {
				return nil, nil
			}
			return s, nil
		}
	}
	if s, ok := v.(vm.Seq); ok {
		return s, nil
	}
	return nil, fmt.Errorf("don't know how to create ISeq from %s", v.Type())
}

func mapLazy1(f vm.Fn, s vm.Seq) vm.Seq {
	if s == nil {
		return nil
	}
	captured := s
	thunk, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		v, err := f.Invoke([]vm.Value{captured.First()})
		if err != nil {
			return vm.NIL, err
		}
		rest := captured.Next()
		tail := mapLazy1(f, rest)
		if tail == nil {
			return vm.EmptyList.Cons(v), nil
		}
		return vm.NewCons(v, tail), nil
	})
	return vm.NewLazySeq(thunk.(vm.Fn))
}

func mapLazyN(f vm.Fn, seqs []vm.Seq) vm.Seq {
	for _, s := range seqs {
		if s == nil {
			return nil
		}
	}
	captured := make([]vm.Seq, len(seqs))
	copy(captured, seqs)
	thunk, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		fargs := make([]vm.Value, len(captured))
		nexts := make([]vm.Seq, len(captured))
		for i, s := range captured {
			fargs[i] = s.First()
			nexts[i] = s.Next()
		}
		v, err := f.Invoke(fargs)
		if err != nil {
			return vm.NIL, err
		}
		tail := mapLazyN(f, nexts)
		if tail == nil {
			return vm.EmptyList.Cons(v), nil
		}
		return vm.NewCons(v, tail), nil
	})
	return vm.NewLazySeq(thunk.(vm.Fn))
}

// nolint
func installLangNS() {
	plus, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.MakeInt(0), nil
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		acc := vs[0]
		for i := 1; i < len(vs); i++ {
			var err error
			acc, err = vm.NumAdd(acc, vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return acc, nil
	})

	mul, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.MakeInt(1), nil
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		acc := vs[0]
		for i := 1; i < len(vs); i++ {
			var err error
			acc, err = vm.NumMul(acc, vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return acc, nil
	})

	sub, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vm.NumNeg(vs[0])
		}
		acc := vs[0]
		for i := 1; i < len(vs); i++ {
			var err error
			acc, err = vm.NumSub(acc, vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return acc, nil
	})

	div, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		acc := vs[0]
		for i := 1; i < len(vs); i++ {
			var err error
			acc, err = vm.NumDiv(acc, vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return acc, nil
	})

	equals, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		length := len(vs)
		if length < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}

		for i := 1; i < length; i++ {
			if !valueEquals(vs[0], vs[i]) {
				return vm.FALSE, nil
			}
		}
		return vm.TRUE, nil
	})

	gt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		for i := 0; i < len(vs)-1; i++ {
			r, err := vm.NumGt(vs[i], vs[i+1])
			if err != nil {
				return vm.NIL, err
			}
			if !r {
				return vm.FALSE, nil
			}
		}
		return vm.TRUE, nil
	})

	lt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		for i := 0; i < len(vs)-1; i++ {
			r, err := vm.NumLt(vs[i], vs[i+1])
			if err != nil {
				return vm.NIL, err
			}
			if !r {
				return vm.FALSE, nil
			}
		}
		return vm.TRUE, nil
	})

	ge, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		for i := 0; i < len(vs)-1; i++ {
			r, err := vm.NumGe(vs[i], vs[i+1])
			if err != nil {
				return vm.NIL, err
			}
			if !r {
				return vm.FALSE, nil
			}
		}
		return vm.TRUE, nil
	})

	le, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		for i := 0; i < len(vs)-1; i++ {
			r, err := vm.NumLe(vs[i], vs[i+1])
			if err != nil {
				return vm.NIL, err
			}
			if !r {
				return vm.FALSE, nil
			}
		}
		return vm.TRUE, nil
	})

	mod, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.NumMod(vs[0], vs[1])
	})

	abs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.NumAbs(vs[0])
	})

	// and/or are now short-circuiting macros defined in core.lg

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
		seq, err := seqOf(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		// Realize lazy seqs to check emptiness
		if ls, ok := seq.(*vm.LazySeq); ok {
			seq = ls.Seq()
		}
		if seq == nil || seq == vm.EmptyList {
			return vm.ArrayVector{}, nil
		}
		ret := []vm.Value{}
		for seq != nil {
			ret = append(ret, seq.First())
			seq = seq.Next()
		}
		return vm.NewArrayVector(ret), nil
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
		return vm.NIL, fmt.Errorf("keyword expects keyword, symbol, or string, got %s", vs[0].Type())
	})

	// symbol(name) or symbol(ns, name)
	symbolf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		toStr := func(v vm.Value) (string, bool) {
			switch s := v.(type) {
			case vm.String:
				return string(s), true
			case vm.Symbol:
				return string(s), true
			default:
				return "", false
			}
		}
		if len(vs) == 1 {
			if s, ok := toStr(vs[0]); ok {
				return vm.Symbol(s), nil
			}
			return vm.NIL, fmt.Errorf("symbol expected String or Symbol")
		}
		nsStr, ok1 := toStr(vs[0])
		nameStr, ok2 := toStr(vs[1])
		if !ok1 || !ok2 {
			return vm.NIL, fmt.Errorf("symbol expected String or Symbol")
		}
		return vm.Symbol(nsStr + "/" + nameStr), nil
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
			return vm.NIL, err
		}
		return colla.Assoc(key, v), nil
	})

	cons, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		elem := vs[0]
		if vs[1] == vm.NIL {
			return vm.EmptyList.Cons(elem), nil
		}
		seq, err := seqOf(vs[1])
		if err != nil {
			return vm.NIL, fmt.Errorf("cons expected Seq")
		}
		if seq == nil {
			return vm.EmptyList.Cons(elem), nil
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
		switch s := vs[0].(type) {
		case *vm.PersistentSet:
			result := s
			for _, v := range vs[1:] {
				result = result.Disj(v)
			}
			return result, nil
		case vm.Set:
			for _, v := range vs[1:] {
				s = s.Disj(v)
			}
			return s, nil
		default:
			return vm.NIL, fmt.Errorf("disj expected Set")
		}
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
		if vs[0] == vm.NIL {
			return vm.NIL, nil
		}
		if seq, ok := vs[0].(vm.Seq); ok {
			return seq.First(), nil
		}
		if sq, ok := vs[0].(vm.Sequable); ok {
			s := sq.Seq()
			if s == nil || s == vm.EmptyList {
				return vm.NIL, nil
			}
			return s.First(), nil
		}
		return vm.NIL, fmt.Errorf("first expected Seq")
	})

	second, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
		}
		seq, err := seqOf(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("second expected Seq")
		}
		if seq == nil {
			return vm.NIL, nil
		}
		n := seq.Next()
		if n == nil {
			return vm.NIL, nil
		}
		return n.First(), nil
	})

	next, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
		}
		seq, err := seqOf(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("next expected Seq")
		}
		if seq == nil {
			return vm.NIL, nil
		}
		n := seq.Next()
		if n == nil {
			return vm.NIL, nil
		}
		return n, nil
	})

	rest, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.EmptyList, nil
		}
		s, err := seqOf(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("rest expected Seq")
		}
		if s == nil {
			return vm.EmptyList, nil
		}
		return s.More(), nil
	})

	seq, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
		}
		// Check for empty collection before calling Seq (Clojure semantics: seq of empty returns nil)
		// Skip for types that may be infinite or expensive to count (Cons, LazySeq)
		switch vs[0].(type) {
		case *vm.Cons, *vm.LazySeq:
			// Don't count — could be infinite
		default:
			if coll, ok := vs[0].(vm.Collection); ok {
				if coll.RawCount() == 0 {
					return vm.NIL, nil
				}
			}
		}
		sqbl, ok := vs[0].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("seq expected Seqauble")
		}
		n := sqbl.Seq()
		// Return nil for empty sequences (Clojure semantics)
		if n == nil || n == vm.EmptyList {
			return vm.NIL, nil
		}
		return n, nil
	})

	isSeq, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		_, ok := vs[0].(vm.Seq)
		return vm.Boolean(ok), nil
	})

	isColl, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		_, ok := vs[0].(vm.Collection)
		return vm.Boolean(ok), nil
	})

	empty, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		coll, ok := vs[0].(vm.Collection)
		if !ok {
			return vm.NIL, fmt.Errorf("empty expected Collection")
		}
		return coll.Empty(), nil
	})

	get, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		vl := len(vs)
		if vl < 2 || vl > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		key := vs[1]
		as, ok := vs[0].(vm.Lookup)
		if !ok {
			// Not a collection - return default value if provided, otherwise nil
			if vl == 3 {
				return vs[2], nil
			}
			return vm.NIL, nil
		}
		if vl == 2 {
			return as.ValueAt(key), nil
		}
		return as.ValueAtOr(key, vs[2]), nil
	})

	// nth: indexed access that works on any sequential type.
	// Fast path for vectors (O(1)), linear walk for seqs (O(n)).
	nthf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		vl := len(vs)
		if vl < 2 || vl > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		idx, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("nth index must be an integer")
		}
		n := int(idx)
		var notFound vm.Value = vm.NIL
		if vl == 3 {
			notFound = vs[2]
		}
		if vs[0] == vm.NIL {
			return notFound, nil
		}
		// Fast path: Lookup types (ArrayVector, PersistentVector, String)
		if l, ok := vs[0].(vm.Lookup); ok {
			v := l.ValueAtOr(vm.Int(n), notFound)
			return v, nil
		}
		// Seq path: linear walk
		s, err := seqOf(vs[0])
		if err != nil {
			return notFound, nil
		}
		for i := 0; s != nil; i++ {
			if i == n {
				return s.First(), nil
			}
			s = s.Next()
		}
		return notFound, nil
	})

	count, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		seq, ok := vs[0].(vm.Counted)
		if !ok {
			return vm.NIL, fmt.Errorf("count expected Counted")
		}
		return seq.Count(), nil
	})

	// map builtin: eager for small counted collections, lazy otherwise
	mapf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		mfn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("map expected Fn")
		}
		// single collection path
		if len(vs) == 2 {
			s, err := seqOf(vs[1])
			if err != nil {
				return vm.NIL, fmt.Errorf("map expected Sequable")
			}
			if s == nil || s == vm.EmptyList {
				return vm.EmptyList, nil
			}
			// Eager fast path: small counted collections (≤32 elements)
			// Skip RawCount for LazySeq/Cons — could be infinite
			length := 0
			switch vs[1].(type) {
			case *vm.LazySeq, *vm.Cons:
				// Don't count — use lazy path
			default:
				if col, ok := vs[1].(vm.Counted); ok {
					length = col.RawCount()
				}
			}
			if length > 0 && length <= 32 {
				newseq := make([]vm.Value, length)
				i := 0
				for s != nil {
					newseq[i], err = mfn.Invoke([]vm.Value{s.First()})
					if err != nil {
						return vm.NIL, err
					}
					s = s.Next()
					i++
				}
				ret, _ := vm.ListType.Box(newseq[:i])
				return ret, nil
			}
			// lazy path
			return mapLazy1(mfn, s), nil
		}
		// multi-collection path
		colls := vs[1:]
		seqs := make([]vm.Seq, len(colls))
		for i := range colls {
			s, err := seqOf(colls[i])
			if err != nil {
				return vm.NIL, fmt.Errorf("map expected Sequable collection")
			}
			if s == nil || s == vm.EmptyList {
				return vm.EmptyList, nil
			}
			seqs[i] = s
		}
		// Check if all collections are small and counted (skip LazySeq/Cons)
		minlen := math.MaxInt
		allCounted := true
		for i := range colls {
			switch colls[i].(type) {
			case *vm.LazySeq, *vm.Cons:
				allCounted = false
				continue
			}
			if coll, ok := colls[i].(vm.Counted); ok {
				c := coll.RawCount()
				if c < minlen {
					minlen = c
				}
			} else {
				allCounted = false
				break
			}
		}
		if allCounted && minlen > 0 && minlen <= 32 {
			newseq := make([]vm.Value, minlen)
			for i := 0; i < minlen; i++ {
				fargs := make([]vm.Value, len(seqs))
				for j := range seqs {
					fargs[j] = seqs[j].First()
					seqs[j] = seqs[j].Next()
				}
				newseq[i], err = mfn.Invoke(fargs)
				if err != nil {
					return vm.NIL, err
				}
			}
			return vm.ListType.Box(newseq)
		}
		// lazy path for multi-collection
		return mapLazyN(mfn, seqs), nil
	})

	mapv, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		v, err := mapf.(vm.Fn).Invoke(vs)
		if err != nil {
			return vm.NIL, err
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
		// Handle nil and empty collections
		if vs[sidx] == vm.NIL {
			if len(vs) == 3 {
				return vs[1], nil
			}
			return mfn.Invoke(nil)
		}
		// Check for empty collection first (skip for lazy/cons — RawCount forces realization)
		switch vs[sidx].(type) {
		case *vm.LazySeq, *vm.Cons:
			// don't call RawCount — could be infinite
		default:
			if coll, ok := vs[sidx].(vm.Collection); ok {
				if coll.RawCount() == 0 {
					if len(vs) == 3 {
						return vs[1], nil
					}
					return mfn.Invoke(nil)
				}
			}
		}
		seq, err := seqOf(vs[sidx])
		if err != nil {
			return vm.NIL, fmt.Errorf("reduce expected Seq")
		}
		if seq == nil {
			if len(vs) == 3 {
				return vs[1], nil
			}
			return mfn.Invoke(nil)
		}
		var acc vm.Value
		if len(vs) == 3 {
			acc = vs[1]
		} else {
			acc = seq.First()
			seq = seq.Next()
		}
		for seq != nil {
			acc, err = mfn.Invoke([]vm.Value{acc, seq.First()})
			if err != nil {
				return vm.NIL, err
			}
			if r, ok := acc.(*vm.Reduced); ok {
				return r.Deref(), nil
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
		seq, err := seqOf(vs[1])
		if err != nil {
			return vm.NIL, fmt.Errorf("some expected Seq")
		}
		for seq != nil {
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
			if vs[i] == vm.NIL {
				continue
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
		if vs[1] == vm.NIL {
			return f.Invoke(nil)
		}
		if av, ok := vs[1].(vm.ArrayVector); ok {
			return f.Invoke(av)
		}
		seq, err := seqOf(vs[1])
		if err != nil {
			return vm.NIL, fmt.Errorf("apply expected Seq")
		}
		if seq == nil {
			return f.Invoke(nil)
		}
		var args []vm.Value
		for seq != nil {
			args = append(args, seq.First())
			seq = seq.Next()
		}
		return f.Invoke(args)
	})

	inNs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		sym := vs[0]
		if sym.Type() != vm.SymbolType {
			return vm.NIL, fmt.Errorf("in-ns expected Symbol")
		}
		nns := LookupOrRegisterNSNoLoad(string(sym.(vm.Symbol)))
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

	aliasf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		al, ok := vs[0].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("alias expected Symbol")
		}
		nsSym, ok := vs[1].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("alias expected Symbol")
		}
		cns := CurrentNS.Deref().(*vm.Namespace)
		target := NS(string(nsSym))
		cns.Alias(al, target)
		return vm.NIL, nil
	})

	referList, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		nsSym, ok := vs[0].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("refer-list expected ns Symbol")
		}
		arr, ok := vs[1].(vm.ArrayVector)
		if !ok {
			return vm.NIL, fmt.Errorf("refer-list expected vector of Symbols")
		}
		syms := make([]vm.Symbol, 0, len(arr))
		for i := range arr {
			if s, ok := arr[i].(vm.Symbol); ok {
				syms = append(syms, s)
			}
		}
		cns := CurrentNS.Deref().(*vm.Namespace)
		target := NS(string(nsSym))
		// Convert []vm.Symbol to []vm.Symbol type alias in vm
		vmSyms := make([]vm.Symbol, len(syms))
		copy(vmSyms, syms)
		cns.ReferList(target, vmSyms)
		return vm.NIL, nil
	})

	// removed resolve-var helper (prefer compile-time resolution)

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
			if vs[i] == vm.NIL {
				continue
			}
			vseq, err := seqOf(vs[i])
			if err != nil {
				return vm.NIL, fmt.Errorf("concat expected Seq")
			}
			for vseq != nil {
				ret = append(ret, vseq.First())
				vseq = vseq.Next()
			}
		}
		r, err := vm.ListType.Box(ret)
		if err != nil {
			return vm.NIL, fmt.Errorf("concat failed: %w", err)
		}
		return r, nil
	})

	// slurp (reintroduced)
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

	// swap-vals!: like swap! but returns [old new]
	swapVals, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		at, ok := vs[0].(*vm.Atom)
		if !ok {
			return vm.NIL, fmt.Errorf("swap-vals! expected Atom")
		}
		fn, ok := vs[1].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("swap-vals! expected Fn")
		}
		old := at.Deref()
		newVal, err := at.Swap(fn, vs[2:])
		if err != nil {
			return vm.NIL, err
		}
		return vm.ArrayVector{old, newVal}, nil
	})

	// reset-vals!: like reset! but returns [old new]
	resetVals, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		at, ok := vs[0].(*vm.Atom)
		if !ok {
			return vm.NIL, fmt.Errorf("reset-vals! expected Atom")
		}
		old := at.Deref()
		at.Reset(vs[1])
		return vm.ArrayVector{old, vs[1]}, nil
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
				fmt.Println(err)
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
		return vm.MakeInt(i), err
	})

	max, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		m := vs[0]
		for i := 1; i < len(vs); i++ {
			gt, err := vm.NumGt(vs[i], m)
			if err != nil {
				return vm.NIL, err
			}
			if gt {
				m = vs[i]
			}
		}
		return m, nil
	})

	min, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		m := vs[0]
		for i := 1; i < len(vs); i++ {
			lt, err := vm.NumLt(vs[i], m)
			if err != nil {
				return vm.NIL, err
			}
			if lt {
				m = vs[i]
			}
		}
		return m, nil
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
			if err != nil {
				return false // abort: previous comparison failed
			}
			var b vm.Value
			b, err = comp.Invoke([]vm.Value{temp[i], temp[j]})
			if err != nil {
				return false
			}
			return vm.IsTruthy(b)
		})
		if err != nil {
			return vm.NIL, err
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

	strReplace, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("str-replace expected String")
		}
		r, ok := vs[2].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("str-replace expected String")
		}
		switch vs[1].(type) {
		case vm.String:
			return vm.String(strings.ReplaceAll(string(s), string(vs[1].(vm.String)), string(r))), nil
		case *vm.Regex:
			return vm.String(vs[1].(*vm.Regex).ReplaceAll(string(s), string(r))), nil
		default:
			return vm.NIL, fmt.Errorf("str-replace expected String or Regex")
		}
	})

	intf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if i, ok := vs[0].(vm.Int); ok {
			return i, nil
		}
		if f, ok := vs[0].(vm.Float); ok {
			return vm.MakeInt(int(f)), nil
		}
		if i, ok := vs[0].(vm.Char); ok {
			return vm.Int(int(i)), nil
		}
		return vm.NIL, fmt.Errorf("%s can't be coerced to int", vs[0])
	})

	floatf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		f, ok := vm.ToFloat(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("%s can't be coerced to float", vs[0])
		}
		return vm.Float(f), nil
	})

	isNumber, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.Boolean(vm.IsNumber(vs[0])), nil
	})

	isFloat, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		_, ok := vs[0].(vm.Float)
		return vm.Boolean(ok), nil
	})

	isInt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		_, ok := vs[0].(vm.Int)
		return vm.Boolean(ok), nil
	})

	char, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		if i, ok := vs[0].(vm.Int); ok {
			return vm.Char(rune(i)), nil
		}
		return vm.NIL, fmt.Errorf("%s can't be coerced to char", vs[0])
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

	peek, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.ArrayVector:
			if len(v) == 0 {
				return vm.NIL, nil
			}
			return v[len(v)-1], nil
		case vm.PersistentVector:
			if v.RawCount() == 0 {
				return vm.NIL, nil
			}
			return v.ValueAt(vm.Int(v.RawCount() - 1)), nil
		case vm.Seq:
			return v.First(), nil
		default:
			if sq, ok := vs[0].(vm.Sequable); ok {
				s := sq.Seq()
				if s == nil || s == vm.EmptyList {
					return vm.NIL, nil
				}
				return s.First(), nil
			}
			return vm.NIL, fmt.Errorf("peek expected Seq or Vec")
		}
	})

	pop, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch vs[0].(type) {
		case vm.PersistentVector:
			v := vs[0].(vm.PersistentVector)
			if v.RawCount() < 1 {
				return vm.NIL, fmt.Errorf("can't pop empty vector")
			}
			// Rebuild without last element
			vals := v.Unbox().([]vm.Value)
			return vm.NewPersistentVector(vals[:len(vals)-1]), nil
		case vm.ArrayVector:
			v := vs[0].(vm.ArrayVector)
			if v.RawCount() < 1 {
				return vm.NIL, fmt.Errorf("can't pop empty vector")
			}
			return vm.ArrayVector(v[0 : len(v)-1]), nil
		case vm.Seq:
			r := vs[0].(vm.Seq).Next()
			if r == nil {
				return vm.NIL, fmt.Errorf("can't pop empty seq")
			}
			return r, nil
		default:
			return vm.NIL, fmt.Errorf("pop expected Seq or Vec")
		}
	})

	iterate, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		f, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("iterate expected a function")
		}
		return vm.NewIterate(f, vs[1]), nil
	})

	repeat, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vm.NewRepeat(vs[0], -1), nil
		}
		n, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("repeat expected an Int")
		}
		if int(n) <= 0 {
			return vm.EmptyList, nil
		}
		return vm.NewRepeat(vs[1], int(n)), nil
	})

	refer, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		cns := CurrentNS.Deref().(*vm.Namespace)
		s, ok := vs[0].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("refer expected Symbol")
		}
		alias := ""
		if len(vs) > 1 {
			if str, ok := vs[1].(vm.String); ok {
				alias = string(str)
			}
		}
		all := true
		if len(vs) > 2 {
			if b, ok := vs[2].(vm.Boolean); ok {
				all = bool(b)
			}
		}
		cns.Refer(NS(string(s)), alias, all)
		return vm.NIL, nil
	})

	// String utility builtins (for string namespace)
	trimf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("trim expected String")
		}
		return vm.String(strings.TrimSpace(string(s))), nil
	})

	trimlf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("triml expected String")
		}
		return vm.String(strings.TrimLeft(string(s), " \t\n\r")), nil
	})

	trimrf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("trimr expected String")
		}
		return vm.String(strings.TrimRight(string(s), " \t\n\r")), nil
	})

	upperCase, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("upper-case expected String")
		}
		return vm.String(strings.ToUpper(string(s))), nil
	})

	lowerCase, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("lower-case expected String")
		}
		return vm.String(strings.ToLower(string(s))), nil
	})

	startsWith, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("starts-with? expected String")
		}
		p, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("starts-with? expected String prefix")
		}
		return vm.Boolean(strings.HasPrefix(string(s), string(p))), nil
	})

	endsWith, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("ends-with? expected String")
		}
		p, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("ends-with? expected String suffix")
		}
		return vm.Boolean(strings.HasSuffix(string(s), string(p))), nil
	})

	includesStr, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("includes? expected String")
		}
		p, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("includes? expected String substr")
		}
		return vm.Boolean(strings.Contains(string(s), string(p))), nil
	})

	// subs: substring
	subs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("subs expected String")
		}
		start, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("subs expected Int start")
		}
		str := string(s)
		si := int(start)
		if si < 0 || si > len(str) {
			return vm.NIL, fmt.Errorf("string index out of range")
		}
		if len(vs) == 3 {
			end, ok := vs[2].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("subs expected Int end")
			}
			ei := int(end)
			if ei < si || ei > len(str) {
				return vm.NIL, fmt.Errorf("string index out of range")
			}
			return vm.String(str[si:ei]), nil
		}
		return vm.String(str[si:]), nil
	})

	// format: sprintf-style string formatting
	formatf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		fmtStr, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("format expected String")
		}
		fmts := string(fmtStr)
		args := make([]interface{}, len(vs)-1)
		// Scan format string to determine which args need float promotion
		vi := 0
		for fi := 0; fi < len(fmts) && vi < len(args); fi++ {
			if fmts[fi] != '%' {
				continue
			}
			fi++ // skip %
			if fi >= len(fmts) {
				break
			}
			if fmts[fi] == '%' {
				continue // %% literal
			}
			// Skip flags, width, precision
			for fi < len(fmts) && (fmts[fi] == '-' || fmts[fi] == '+' || fmts[fi] == ' ' || fmts[fi] == '0' || fmts[fi] == '#' || (fmts[fi] >= '0' && fmts[fi] <= '9') || fmts[fi] == '.') {
				fi++
			}
			if fi >= len(fmts) {
				break
			}
			verb := fmts[fi]
			switch v := vs[vi+1].(type) {
			case vm.Int:
				if verb == 'f' || verb == 'e' || verb == 'g' || verb == 'E' || verb == 'G' {
					args[vi] = float64(v)
				} else {
					args[vi] = int(v)
				}
			case vm.Float:
				args[vi] = float64(v)
			case vm.String:
				args[vi] = string(v)
			case vm.Boolean:
				args[vi] = bool(v)
			default:
				args[vi] = vs[vi+1].Unbox()
			}
			vi++
		}
		return vm.String(fmt.Sprintf(string(fmtStr), args...)), nil
	})

	// rand: returns a random float between 0 (inclusive) and 1 (exclusive)
	// or between 0 and n
	randf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.Float(rand.Float64()), nil
		}
		if len(vs) == 1 {
			if n, ok := vs[0].(vm.Int); ok {
				return vm.Float(rand.Float64() * float64(n)), nil
			}
			if n, ok := vs[0].(vm.Float); ok {
				return vm.Float(rand.Float64() * float64(n)), nil
			}
			return vm.NIL, fmt.Errorf("rand expected number")
		}
		return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
	})

	// rand-int: returns a random integer between 0 (inclusive) and n (exclusive)
	randInt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		n, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("rand-int expected Int")
		}
		if int(n) <= 0 {
			return vm.MakeInt(0), nil
		}
		return vm.MakeInt(rand.Intn(int(n))), nil
	})

	// rand-nth: returns a random element from a collection
	randNth, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		coll, ok := vs[0].(vm.Collection)
		if !ok {
			return vm.NIL, fmt.Errorf("rand-nth expected Collection")
		}
		n := coll.RawCount()
		if n == 0 {
			return vm.NIL, fmt.Errorf("rand-nth called on empty collection")
		}
		idx := rand.Intn(n)
		if l, ok := vs[0].(vm.Lookup); ok {
			return l.ValueAt(vm.Int(idx)), nil
		}
		// Fallback: iterate
		s, _ := seqOf(vs[0])
		for i := 0; i < idx; i++ {
			s = s.Next()
		}
		return s.First(), nil
	})

	// shuffle: returns a random permutation of a collection as a vector
	shuffle, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		// Collect into slice
		s, err := seqOf(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		var vals []vm.Value
		for s != nil {
			vals = append(vals, s.First())
			s = s.Next()
		}
		// Fisher-Yates shuffle
		rand.Shuffle(len(vals), func(i, j int) {
			vals[i], vals[j] = vals[j], vals[i]
		})
		return vm.NewArrayVector(vals), nil
	})

	// transient: create a transient (mutable) version of a persistent collection
	transientf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case *vm.PersistentMap:
			return vm.NewTransientMap(v), nil
		case vm.ArrayVector:
			return vm.NewTransientVector([]vm.Value(v)), nil
		case vm.PersistentVector:
			vals := v.Unbox().([]vm.Value)
			return vm.NewTransientVector(vals), nil
		default:
			return vm.NIL, fmt.Errorf("transient not supported on %s", vs[0].Type().Name())
		}
	})

	// persistent!: freeze a transient back to a persistent collection
	persistentf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case *vm.TransientMap:
			return v.Persistent(), nil
		case *vm.TransientVector:
			return v.Persistent(), nil
		default:
			return vm.NIL, fmt.Errorf("persistent! not supported on %s", vs[0].Type().Name())
		}
	})

	// conj!: mutating conj on a transient
	conjBang, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch t := vs[0].(type) {
		case *vm.TransientMap:
			var err error
			for i := 1; i < len(vs); i++ {
				t, err = t.Conj(vs[i])
				if err != nil {
					return vm.NIL, err
				}
			}
			return t, nil
		case *vm.TransientVector:
			var err error
			for i := 1; i < len(vs); i++ {
				t, err = t.Conj(vs[i])
				if err != nil {
					return vm.NIL, err
				}
			}
			return t, nil
		default:
			return vm.NIL, fmt.Errorf("conj! not supported on %s", vs[0].Type().Name())
		}
	})

	// assoc!: mutating assoc on a transient
	assocBang, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 3 || len(vs)%2 == 0 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch t := vs[0].(type) {
		case *vm.TransientMap:
			var err error
			for i := 1; i < len(vs); i += 2 {
				t, err = t.Assoc(vs[i], vs[i+1])
				if err != nil {
					return vm.NIL, err
				}
			}
			return t, nil
		case *vm.TransientVector:
			var err error
			for i := 1; i < len(vs); i += 2 {
				t, err = t.Assoc(vs[i], vs[i+1])
				if err != nil {
					return vm.NIL, err
				}
			}
			return t, nil
		default:
			return vm.NIL, fmt.Errorf("assoc! not supported on %s", vs[0].Type().Name())
		}
	})

	// dissoc!: mutating dissoc on a transient map
	dissocBang, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		t, ok := vs[0].(*vm.TransientMap)
		if !ok {
			return vm.NIL, fmt.Errorf("dissoc! expected TransientMap")
		}
		var err error
		for i := 1; i < len(vs); i++ {
			t, err = t.Dissoc(vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return t, nil
	})

	// make-record-type: create a RecordType with name and field keywords
	makeRecordType, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		name, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("make-record-type expected String name")
		}
		fields := make([]vm.Keyword, len(vs)-1)
		for i := 1; i < len(vs); i++ {
			kw, ok := vs[i].(vm.Keyword)
			if !ok {
				return vm.NIL, fmt.Errorf("make-record-type expected Keyword fields")
			}
			fields[i-1] = kw
		}
		return vm.NewRecordType(string(name), fields), nil
	})

	// make-record: create a Record from a RecordType and a map
	makeRecord, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		rt, ok := vs[0].(*vm.RecordType)
		if !ok {
			return vm.NIL, fmt.Errorf("make-record expected RecordType")
		}
		m, ok := vs[1].(*vm.PersistentMap)
		if !ok {
			return vm.NIL, fmt.Errorf("make-record expected Map")
		}
		return vm.NewRecord(rt, m), nil
	})

	// record?: check if a value is a Record
	isRecord, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		_, ok := vs[0].(*vm.Record)
		return vm.Boolean(ok), nil
	})

	// defprotocol*: create a protocol (called by defprotocol macro)
	defProtocol, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		name, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("defprotocol* expected String name")
		}
		methods := make([]vm.Symbol, len(vs)-1)
		for i := 1; i < len(vs); i++ {
			s, ok := vs[i].(vm.Symbol)
			if !ok {
				return vm.NIL, fmt.Errorf("defprotocol* expected Symbol method names")
			}
			methods[i-1] = s
		}
		return vm.NewProtocol(string(name), methods), nil
	})

	// extend-type*: extend a protocol for a type (called by extend-type macro)
	extendType, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		protocol, ok := vs[0].(*vm.Protocol)
		if !ok {
			return vm.NIL, fmt.Errorf("extend-type* expected Protocol")
		}
		implMap, ok := vs[2].(*vm.PersistentMap)
		if !ok {
			return vm.NIL, fmt.Errorf("extend-type* expected map of implementations")
		}
		// vs[1] is the type to extend — either a ValueType or nil
		if vs[1] == vm.NIL {
			protocol.ExtendNil(implMap)
		} else {
			vt, ok := vs[1].(vm.ValueType)
			if !ok {
				return vm.NIL, fmt.Errorf("extend-type* expected a type, got %s", vs[1].Type().Name())
			}
			protocol.Extend(vt, implMap)
		}
		return vm.NIL, nil
	})

	// make-protocol-fn: create a ProtocolFn for dispatch
	makeProtocolFn, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		protocol, ok := vs[0].(*vm.Protocol)
		if !ok {
			return vm.NIL, fmt.Errorf("make-protocol-fn expected Protocol")
		}
		methodName, ok := vs[1].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("make-protocol-fn expected Symbol")
		}
		return vm.NewProtocolFn(protocol, methodName), nil
	})

	// satisfies?: check if a value's type implements a protocol
	satisfies, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		protocol, ok := vs[0].(*vm.Protocol)
		if !ok {
			return vm.NIL, fmt.Errorf("satisfies? expected Protocol")
		}
		return vm.Boolean(protocol.Satisfies(vs[1])), nil
	})

	// defmulti*: create a multimethod (called by defmulti macro)
	defMulti, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		name, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("defmulti* expected String name")
		}
		dispatchFn, ok := vs[1].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("defmulti* expected Fn")
		}
		var defaultVal vm.Value = vm.Keyword("default")
		if len(vs) == 3 {
			defaultVal = vs[2]
		}
		return vm.NewMultiFn(string(name), dispatchFn, defaultVal), nil
	})

	// defmethod*: add a method to a multimethod (called by defmethod macro)
	defMethod, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		mf, ok := vs[0].(*vm.MultiFn)
		if !ok {
			return vm.NIL, fmt.Errorf("defmethod* expected MultiFn")
		}
		dispatchVal := vs[1]
		method, ok := vs[2].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("defmethod* expected Fn")
		}
		return mf.AddMethod(dispatchVal, method), nil
	})

	// methods: return the method map of a multimethod
	methods, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		mf, ok := vs[0].(*vm.MultiFn)
		if !ok {
			return vm.NIL, fmt.Errorf("methods expected MultiFn")
		}
		return mf.Methods(), nil
	})

	// pr-str: print readably to string (with quotes on strings)
	prStr, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		b := &strings.Builder{}
		for i := range vs {
			if i > 0 {
				b.WriteRune(' ')
			}
			b.WriteString(vs[i].String())
		}
		return vm.String(b.String()), nil
	})

	// prn: print readably + newline
	prn, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		b := &strings.Builder{}
		for i := range vs {
			if i > 0 {
				b.WriteRune(' ')
			}
			b.WriteString(vs[i].String())
		}
		fmt.Println(b)
		return vm.NIL, nil
	})

	// re-find: find first match of regex in string
	reFind, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		re, ok := vs[0].(*vm.Regex)
		if !ok {
			return vm.NIL, fmt.Errorf("re-find expected Regex")
		}
		s, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("re-find expected String")
		}
		matches := re.FindStringSubmatch(string(s))
		if matches == nil {
			return vm.NIL, nil
		}
		if len(matches) == 1 {
			return vm.String(matches[0]), nil
		}
		result := make(vm.ArrayVector, len(matches))
		for i, m := range matches {
			result[i] = vm.String(m)
		}
		return result, nil
	})

	// re-matches: match entire string against regex
	reMatches, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		re, ok := vs[0].(*vm.Regex)
		if !ok {
			return vm.NIL, fmt.Errorf("re-matches expected Regex")
		}
		s, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("re-matches expected String")
		}
		matches := re.FindStringSubmatch(string(s))
		if matches == nil || matches[0] != string(s) {
			return vm.NIL, nil
		}
		if len(matches) == 1 {
			return vm.String(matches[0]), nil
		}
		result := make(vm.ArrayVector, len(matches))
		for i, m := range matches {
			result[i] = vm.String(m)
		}
		return result, nil
	})

	// re-seq: return lazy seq of all matches
	reSeq, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		re, ok := vs[0].(*vm.Regex)
		if !ok {
			return vm.NIL, fmt.Errorf("re-seq expected Regex")
		}
		s, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("re-seq expected String")
		}
		all := re.FindAllString(string(s), -1)
		if all == nil {
			return vm.EmptyList, nil
		}
		vals := make([]vm.Value, len(all))
		for i, m := range all {
			vals[i] = vm.String(m)
		}
		return vm.ListType.Box(vals)
	})

	// require loads a namespace by name (like Clojure's require function for REPL use)
	// Supports: (require 'foo), (require '[foo :as f]), (require '[foo :refer [a b]])
	requiref, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		cns := CurrentNS.Deref().(*vm.Namespace)
		for _, v := range vs {
			switch arg := v.(type) {
			case vm.Symbol:
				NS(string(arg)) // triggers autoloading
			case *vm.ArrayVector:
				// Vector form: [ns-name :as alias] or [ns-name :refer [syms...]]
				if arg.RawCount() < 1 {
					return vm.NIL, fmt.Errorf("require: empty vector")
				}
				nsName, ok := arg.ValueAt(vm.Int(0)).(vm.Symbol)
				if !ok {
					return vm.NIL, fmt.Errorf("require: first element must be a symbol")
				}
				target := NS(string(nsName))
				// Parse options
				for i := 1; i < arg.RawCount()-1; i += 2 {
					opt := arg.ValueAt(vm.Int(int64(i)))
					val := arg.ValueAt(vm.Int(int64(i + 1)))
					switch opt {
					case vm.Keyword("as"):
						if alias, ok := val.(vm.Symbol); ok {
							cns.Alias(alias, target)
						}
					case vm.Keyword("refer"):
						if val == vm.Keyword("all") {
							cns.Refer(target, "", true)
						} else if vec, ok := val.(*vm.ArrayVector); ok {
							syms := make([]vm.Symbol, vec.RawCount())
							for j := 0; j < vec.RawCount(); j++ {
								syms[j] = vec.ValueAt(vm.Int(int64(j))).(vm.Symbol)
							}
							cns.ReferList(target, syms)
						}
					}
				}
			default:
				return vm.NIL, fmt.Errorf("require expected Symbol or Vector, got %s", v.Type().Name())
			}
		}
		return vm.NIL, nil
	})

	// find-ns returns the namespace with the given name, or nil
	findNs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("find-ns expected Symbol")
		}
		ns := nsRegistry[string(s)]
		if ns == nil {
			return vm.NIL, nil
		}
		return ns, nil
	})

	// all-ns returns a list of all loaded namespaces
	allNs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		var nss []vm.Value
		for _, ns := range nsRegistry {
			nss = append(nss, ns)
		}
		return vm.NewList(nss), nil
	})

	// the-ns returns the namespace for a symbol, throwing if not found
	theNs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.Symbol)
		if !ok {
			// If already a namespace, return it
			if ns, ok := vs[0].(*vm.Namespace); ok {
				return ns, nil
			}
			return vm.NIL, fmt.Errorf("the-ns expected Symbol or Namespace")
		}
		ns := nsRegistry[string(s)]
		if ns == nil {
			return vm.NIL, fmt.Errorf("no namespace: %s found", s)
		}
		return ns, nil
	})

	// ns-name returns the name of a namespace as a symbol
	nsName, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		ns, ok := vs[0].(*vm.Namespace)
		if !ok {
			return vm.NIL, fmt.Errorf("ns-name expected Namespace")
		}
		return vm.Symbol(ns.Name()), nil
	})

	// lazy-seq* creates a LazySeq from a thunk function
	lazySeq, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("lazy-seq* expected 1 argument, got %d", len(vs))
		}
		fn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("lazy-seq* expected a function")
		}
		return vm.NewLazySeq(fn), nil
	})

	pushBinding, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		v, ok := vs[0].(*vm.Var)
		if !ok {
			return vm.NIL, fmt.Errorf("push-binding expected Var")
		}
		v.PushBinding(vs[1])
		return vm.NIL, nil
	})

	popBinding, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		v, ok := vs[0].(*vm.Var)
		if !ok {
			return vm.NIL, fmt.Errorf("pop-binding expected Var")
		}
		v.PopBinding()
		return vm.NIL, nil
	})

	withMeta, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
		}
		m, ok := vs[0].(vm.IMeta)
		if !ok {
			return vm.NIL, fmt.Errorf("with-meta not supported on %s", vs[0].Type().Name())
		}
		return m.WithMeta(vs[1]), nil
	})

	metaf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
		}
		m, ok := vs[0].(vm.IMeta)
		if !ok {
			return vm.NIL, nil
		}
		return m.Meta(), nil
	})

	// throw
	throwf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.NIL, vm.NewThrownError(vs[0])
	})

	// ex-info
	exInfo, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		msg, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("ex-info expected String message")
		}
		data, ok := vs[1].(*vm.PersistentMap)
		if !ok {
			return vm.NIL, fmt.Errorf("ex-info expected Map data")
		}
		var cause error
		if len(vs) == 3 {
			if ei, ok := vs[2].(*vm.ExInfo); ok {
				cause = ei
			}
		}
		return vm.NewExInfo(string(msg), data, cause), nil
	})

	// ex-message
	exMessage, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		if ei, ok := vs[0].(*vm.ExInfo); ok {
			return vm.String(ei.Message()), nil
		}
		return vm.NIL, nil
	})

	// ex-data
	exData, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		if ei, ok := vs[0].(*vm.ExInfo); ok {
			return ei.Data(), nil
		}
		return vm.NIL, nil
	})

	// ex-cause
	exCause, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		if ei, ok := vs[0].(*vm.ExInfo); ok {
			if c := ei.Cause(); c != nil {
				if cev, ok := c.(*vm.ExInfo); ok {
					return cev, nil
				}
			}
		}
		return vm.NIL, nil
	})

	// transformer-seq* — (transformer-seq* xform coll) → lazy seq
	// Lazily pulls elements from coll through the transducer xform.
	// Uses a buffer-based approach: each source element may produce 0, 1, or many outputs.
	transformerSeq, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("transformer-seq* expects 2 args")
		}
		xformFn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("transformer-seq* expected xform Fn")
		}

		src, err := seqOf(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		if src == nil {
			return vm.EmptyList, nil
		}

		// Shared mutable state
		type tstate struct {
			src      vm.Seq
			buf      []vm.Value // output items waiting to be yielded
			xf       vm.Fn      // the xform'd reducing fn
			done     bool       // completion called
			stopped  bool       // early termination via reduced
		}

		// The base reducing fn just appends to the buffer.
		// We pass a pointer to the state's buf so the rf can append.
		st := &tstate{src: src}

		bufRf, _ := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
			switch len(args) {
			case 0:
				return vm.NIL, nil // init
			case 1:
				return args[0], nil // completion — identity
			case 2:
				// step: accumulate the output item
				st.buf = append(st.buf, args[1])
				return args[0], nil // return accumulator unchanged
			}
			return vm.NIL, nil
		})

		xfResult, err := xformFn.Invoke([]vm.Value{bufRf})
		if err != nil {
			return vm.NIL, err
		}
		xf, ok := xfResult.(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("xform did not return a function")
		}
		st.xf = xf

		// Pull from source through xform until buf has items or source is exhausted
		fillBuf := func() {
			for len(st.buf) == 0 && st.src != nil && !st.stopped {
				item := st.src.First()
				st.src = st.src.Next()

				result, err := st.xf.Invoke([]vm.Value{vm.NIL, item})
				if err != nil {
					st.stopped = true
					return
				}

				if vm.IsReduced(result) {
					st.stopped = true
					st.src = nil
				}
			}

			// Source exhausted or stopped — call completion to flush
			if (st.src == nil || st.stopped) && !st.done {
				st.done = true
				st.xf.Invoke([]vm.Value{vm.NIL}) // completion arity
			}
		}

		// Build lazy seq that drains the buffer
		var buildSeq func() *vm.LazySeq
		buildSeq = func() *vm.LazySeq {
			thunk, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
				// Drain buffer first
				if len(st.buf) > 0 {
					item := st.buf[0]
					st.buf = st.buf[1:]
					return vm.NewCons(item, buildSeq()), nil
				}
				// Buffer empty — pull more from source
				fillBuf()
				if len(st.buf) == 0 {
					return nil, nil // done
				}
				item := st.buf[0]
				st.buf = st.buf[1:]
				return vm.NewCons(item, buildSeq()), nil
			})
			return vm.NewLazySeq(thunk.(vm.Fn))
		}

		return buildSeq(), nil
	})

	// delay — (delay body) is a macro in core.lg, but we need delay* as the constructor
	delayStar, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("delay* expects 1 arg (thunk fn)")
		}
		fn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("delay* expected Fn")
		}
		return vm.NewDelay(fn), nil
	})

	// force — deref a delay (or return value if not a delay)
	force, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("force expects 1 arg")
		}
		if d, ok := vs[0].(*vm.Delay); ok {
			return d.Force()
		}
		return vs[0], nil
	})

	// delay? — test if value is a Delay
	isDelay, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(*vm.Delay)
		return vm.Boolean(ok), nil
	})

	// realized? — test if a Delay, Promise, or Future has been realized
	isRealized, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		if d, ok := vs[0].(*vm.Delay); ok {
			return vm.Boolean(d.IsRealized()), nil
		}
		if p, ok := vs[0].(*vm.Promise); ok {
			return vm.Boolean(p.IsRealized()), nil
		}
		return vm.FALSE, nil
	})

	// volatile! — create a volatile mutable box
	volatilef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("volatile! expects 1 arg")
		}
		return vm.NewVolatile(vs[0]), nil
	})

	// vreset! — set volatile value
	vreset, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("vreset! expects 2 args")
		}
		v, ok := vs[0].(*vm.Volatile)
		if !ok {
			return vm.NIL, fmt.Errorf("vreset! expected Volatile")
		}
		return v.Reset(vs[1]), nil
	})

	// vswap! — apply fn to volatile value: (vswap! vol f args...)
	vswap, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("vswap! expects at least 2 args")
		}
		v, ok := vs[0].(*vm.Volatile)
		if !ok {
			return vm.NIL, fmt.Errorf("vswap! expected Volatile")
		}
		fn, ok := vs[1].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("vswap! expected Fn")
		}
		args := make([]vm.Value, 1+len(vs)-2)
		args[0] = v.Deref()
		copy(args[1:], vs[2:])
		result, err := fn.Invoke(args)
		if err != nil {
			return vm.NIL, err
		}
		return v.Reset(result), nil
	})

	// reduced — wrap a value to signal early termination
	reducedf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("reduced expects 1 arg")
		}
		return vm.NewReduced(vs[0]), nil
	})

	// reduced? — test if value is Reduced
	isReducedf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		return vm.Boolean(vm.IsReduced(vs[0])), nil
	})

	// compare — generic comparison: -1, 0, 1
	comparef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("compare expects 2 args")
		}
		a, b := vs[0], vs[1]
		// nil sorts before everything
		if a == vm.NIL && b == vm.NIL {
			return vm.MakeInt(0), nil
		}
		if a == vm.NIL {
			return vm.MakeInt(-1), nil
		}
		if b == vm.NIL {
			return vm.MakeInt(1), nil
		}
		// Numbers (including BigInt)
		if vm.IsNumber(a) && vm.IsNumber(b) {
			lt, err := vm.NumLt(a, b)
			if err != nil {
				return vm.NIL, err
			}
			if lt {
				return vm.MakeInt(-1), nil
			}
			gt, err := vm.NumGt(a, b)
			if err != nil {
				return vm.NIL, err
			}
			if gt {
				return vm.MakeInt(1), nil
			}
			return vm.MakeInt(0), nil
		}
		// Strings
		if sa, ok := a.(vm.String); ok {
			if sb, ok := b.(vm.String); ok {
				as, bs := string(sa), string(sb)
				if as < bs {
					return vm.MakeInt(-1), nil
				}
				if as > bs {
					return vm.MakeInt(1), nil
				}
				return vm.MakeInt(0), nil
			}
		}
		// Keywords
		if ka, ok := a.(vm.Keyword); ok {
			if kb, ok := b.(vm.Keyword); ok {
				as, bs := string(ka), string(kb)
				if as < bs {
					return vm.MakeInt(-1), nil
				}
				if as > bs {
					return vm.MakeInt(1), nil
				}
				return vm.MakeInt(0), nil
			}
		}
		// Booleans (false < true)
		if ba, ok := a.(vm.Boolean); ok {
			if bb, ok := b.(vm.Boolean); ok {
				if ba == bb {
					return vm.MakeInt(0), nil
				}
				if !bool(ba) {
					return vm.MakeInt(-1), nil
				}
				return vm.MakeInt(1), nil
			}
		}
		return vm.NIL, fmt.Errorf("compare: cannot compare %s and %s", a.Type().Name(), b.Type().Name())
	})

	// print — like println but no newline, space-separated
	printf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		for i, v := range vs {
			if i > 0 {
				fmt.Print(" ")
			}
			if s, ok := v.(vm.String); ok {
				fmt.Print(string(s))
			} else {
				fmt.Print(v.String())
			}
		}
		return vm.NIL, nil
	})

	// pr — print readably (like prn without newline)
	prf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		for i, v := range vs {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(v.String())
		}
		return vm.NIL, nil
	})

	// --- Bitwise ops ---

	bitAnd, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-and expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-and expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-and expected Int")
		}
		return vm.MakeInt(int(a) & int(b)), nil
	})

	bitOr, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-or expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-or expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-or expected Int")
		}
		return vm.MakeInt(int(a) | int(b)), nil
	})

	bitXor, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-xor expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-xor expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-xor expected Int")
		}
		return vm.MakeInt(int(a) ^ int(b)), nil
	})

	bitNot, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("bit-not expects 1 arg")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-not expected Int")
		}
		return vm.MakeInt(^int(a)), nil
	})

	bitShiftLeft, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-shift-left expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-shift-left expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-shift-left expected Int")
		}
		return vm.MakeInt(int(a) << uint(b)), nil
	})

	bitShiftRight, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-shift-right expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-shift-right expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-shift-right expected Int")
		}
		return vm.MakeInt(int(a) >> uint(b)), nil
	})

	unsignedBitShiftRight, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("unsigned-bit-shift-right expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("unsigned-bit-shift-right expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("unsigned-bit-shift-right expected Int")
		}
		return vm.MakeInt(int(uint(a) >> uint(b))), nil
	})

	bitTest, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-test expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-test expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-test expected Int")
		}
		return vm.Boolean(int(a)&(1<<uint(b)) != 0), nil
	})

	bitSet, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-set expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-set expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-set expected Int")
		}
		return vm.MakeInt(int(a) | (1 << uint(b))), nil
	})

	bitClear, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-clear expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-clear expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-clear expected Int")
		}
		return vm.MakeInt(int(a) &^ (1 << uint(b))), nil
	})

	// re-groups — find all submatch groups: (re-groups regex str) → vector of [match group1 group2 ...]
	reGroups, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("re-groups expects 2 args")
		}
		re, ok := vs[0].(*vm.Regex)
		if !ok {
			return vm.NIL, fmt.Errorf("re-groups expected Regex")
		}
		s, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("re-groups expected String")
		}
		all := re.FindAllStringSubmatch(string(s), -1)
		if all == nil {
			return vm.NIL, nil
		}
		result := make([]vm.Value, len(all))
		for i, match := range all {
			group := make(vm.ArrayVector, len(match))
			for j, m := range match {
				group[j] = vm.String(m)
			}
			result[i] = group
		}
		return vm.NewArrayVector(result), nil
	})

	// promise — create a promise
	promisef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NewPromise(), nil
	})

	// deliver — deliver a value to a promise
	deliver, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("deliver expects 2 args")
		}
		p, ok := vs[0].(*vm.Promise)
		if !ok {
			return vm.NIL, fmt.Errorf("deliver expected Promise")
		}
		return p.Deliver(vs[1]), nil
	})

	// future — run body in a goroutine, return a promise that delivers the result
	// (future* thunk) — internal, macro wraps body
	futureStar, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("future* expects 1 arg (thunk fn)")
		}
		fn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("future* expected Fn")
		}
		p := vm.NewPromise()
		go func() {
			v, err := fn.Invoke(nil)
			if err != nil {
				p.Deliver(vm.NIL)
			} else {
				p.Deliver(v)
			}
		}()
		return p, nil
	})

	// add-watch — (add-watch atom key fn)
	addWatch, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("add-watch expects 3 args")
		}
		a, ok := vs[0].(*vm.Atom)
		if !ok {
			return vm.NIL, fmt.Errorf("add-watch expected Atom")
		}
		fn, ok := vs[2].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("add-watch expected Fn")
		}
		a.AddWatch(vs[1], fn)
		return vs[0], nil
	})

	// remove-watch — (remove-watch atom key)
	removeWatch, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("remove-watch expects 2 args")
		}
		a, ok := vs[0].(*vm.Atom)
		if !ok {
			return vm.NIL, fmt.Errorf("remove-watch expected Atom")
		}
		a.RemoveWatch(vs[1])
		return vs[0], nil
	})

	// alter-meta! — (alter-meta! ref f & args)
	alterMeta, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("alter-meta! expects at least 2 args")
		}
		a, ok := vs[0].(*vm.Atom)
		if !ok {
			return vm.NIL, fmt.Errorf("alter-meta! expected Atom")
		}
		fn, ok := vs[1].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("alter-meta! expected Fn")
		}
		return a.AlterMeta(fn, vs[2:])
	})

	// subvec — (subvec v start) or (subvec v start end)
	subvecf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("subvec expects 2-3 args")
		}
		start, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("subvec expected Int start")
		}
		s := int(start)

		switch v := vs[0].(type) {
		case vm.ArrayVector:
			end := len(v)
			if len(vs) == 3 {
				e, ok := vs[2].(vm.Int)
				if !ok {
					return vm.NIL, fmt.Errorf("subvec expected Int end")
				}
				end = int(e)
			}
			if s < 0 || end > len(v) || s > end {
				return vm.NIL, fmt.Errorf("subvec: index out of bounds")
			}
			result := make([]vm.Value, end-s)
			copy(result, v[s:end])
			return vm.NewArrayVector(result), nil
		case vm.PersistentVector:
			end := v.Count().(vm.Int)
			if len(vs) == 3 {
				e, ok := vs[2].(vm.Int)
				if !ok {
					return vm.NIL, fmt.Errorf("subvec expected Int end")
				}
				end = e
			}
			if s < 0 || int(end) > int(v.Count().(vm.Int)) || s > int(end) {
				return vm.NIL, fmt.Errorf("subvec: index out of bounds")
			}
			result := make([]vm.Value, int(end)-s)
			for i := s; i < int(end); i++ {
				result[i-s] = v.ValueAt(vm.Int(i))
			}
			return vm.NewArrayVector(result), nil
		default:
			return vm.NIL, fmt.Errorf("subvec expected vector")
		}
	})

	// fn? — test if value is callable
	isFn, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(vm.Fn)
		return vm.Boolean(ok), nil
	})

	// unreduced — unwrap Reduced, or return value as-is
	unreduced, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("unreduced expects 1 arg")
		}
		if r, ok := vs[0].(*vm.Reduced); ok {
			return r.Deref(), nil
		}
		return vs[0], nil
	})

	// ensure-reduced — if already Reduced, return as-is; otherwise wrap
	ensureReduced, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("ensure-reduced expects 1 arg")
		}
		if _, ok := vs[0].(*vm.Reduced); ok {
			return vs[0], nil
		}
		return vm.NewReduced(vs[0]), nil
	})

	if err != nil {
		panic("lang NS init failed")
	}

	ns := vm.NewNamespace(NameCoreNS)

	// vars
	CurrentNS = ns.Def("*ns*", ns)

	// Bootstrap no-op ns macro so source files can declare namespaces before core macro is loaded.
	// Expands (ns name ...) to (in-ns 'name), ignoring options.
	nsMacro, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, nil
		}
		nameSym, ok := vs[0].(vm.Symbol)
		if !ok {
			return vm.NIL, nil
		}
		quoteSym := vm.Symbol("quote")
		inNsSym := vm.Symbol("in-ns")
		quoted := vm.EmptyList.Cons(nameSym).Cons(quoteSym)
		form := vm.EmptyList.Cons(quoted).Cons(inNsSym)
		return form, nil
	})
	// Mark as macro
	_ = ns.Def("ns", nsMacro)
	(ns.Lookup("ns").(*vm.Var)).SetMacro()

	// primitive fns
	ns.Def("+", plus)
	ns.Def("*", mul)
	ns.Def("-", sub)
	ns.Def("/", div)

	ns.Def("=", equals)
	ns.Def("gt", gt)
	ns.Def("lt", lt)
	ns.Def("ge", ge)
	ns.Def("le", le)
	ns.Def("mod", mod)
	ns.Def("abs", abs)

	// and/or are now macros in core.lg (short-circuiting)
	// ns.Def("and", and)
	// ns.Def("or", or)
	ns.Def("not", not)

	ns.Def("set-macro!", setMacro)
	ns.Def("gensym", gensym)
	ns.Def("in-ns", inNs)
	ns.Def("use", use)
	ns.Def("alias", aliasf)
	ns.Def("name", name)
	ns.Def("namespace", namespace)

	ns.Def("vector", vector)
	ns.Def("vec", vec)
	ns.Def("hash-map", hashMap)
	ns.Def("list", list)
	ns.Def("range", rangef)
	ns.Def("keyword", keyword)
	ns.Def("symbol", symbolf)
	ns.Def("hash-set", hashSet)

	ns.Def("seq", seq)
	ns.Def("seq?", isSeq)

	// basic predicates needed during early core bootstrap
	ns.Def("coll?", isColl)

	ns.Def("empty", empty)

	ns.Def("assoc", assoc)
	ns.Def("dissoc", dissoc)
	ns.Def("update", update)
	ns.Def("cons", cons)
	ns.Def("conj", conj)
	ns.Def("disj", disj)
	ns.Def("first", first)
	ns.Def("second", second)
	ns.Def("next", next)
	ns.Def("rest", rest)
	ns.Def("get", get)
	ns.Def("nth", nthf)
	ns.Def("count", count)
	ns.Def("contains?", contains)

	ns.Def("map*", mapf)
	ns.Def("mapv", mapv)
	ns.Def("reduce", reduce)
	ns.Def("concat*", concat)
	ns.Def("some", some)

	ns.Def("println", printlnf)

	ns.Def("type", typef)

	ns.Def("apply*", apply)
	ns.Def("deref", deref)

	ns.Def("atom", atom)
	ns.Def("reset!", reset)
	ns.Def("swap!", swap)
	ns.Def("swap-vals!", swapVals)
	ns.Def("reset-vals!", resetVals)

	ns.Def("now", now)

	ns.Def("slurp", slurp)
	ns.Def("spit", spit)
	ns.Def("lines", lines)

	ns.Def("parse-int", parseInt)
	ns.Def("max", max)
	ns.Def("min", min)

	ns.Def("sort", sort)

	ns.Def(".", methodInvoke)

	// async
	ns.Def("go*", gof)
	ns.Def("chan", chanf)
	ns.Def(">!", chanput)
	ns.Def("<!", changet)
	ns.Def(">!!", chanput)
	ns.Def("<!!", changet)

	ns.Def("int", intf)
	ns.Def("float", floatf)
	ns.Def("double", floatf)
	ns.Def("number?", isNumber)
	ns.Def("float?", isFloat)
	ns.Def("int?", isInt)
	ns.Def("char", char)

	ns.Def("str", str)
	ns.Def("split", split)
	ns.Def("str-replace", strReplace)
	ns.Def("re-pattern", regex)
	// namespace utilities
	ns.Def("refer-list", referList)
	ns.Def("refer", refer)
	ns.Def("trim", trimf)
	ns.Def("triml", trimlf)
	ns.Def("trimr", trimrf)
	ns.Def("upper-case", upperCase)
	ns.Def("lower-case", lowerCase)
	ns.Def("starts-with?", startsWith)
	ns.Def("ends-with?", endsWith)
	ns.Def("includes?", includesStr)
	ns.Def("subs", subs)
	ns.Def("format", formatf)
	ns.Def("rand", randf)
	ns.Def("rand-int", randInt)
	ns.Def("rand-nth", randNth)
	ns.Def("shuffle", shuffle)
	ns.Def("transient", transientf)
	ns.Def("persistent!", persistentf)
	ns.Def("conj!", conjBang)
	ns.Def("assoc!", assocBang)
	ns.Def("dissoc!", dissocBang)
	ns.Def("make-record-type", makeRecordType)
	ns.Def("make-record", makeRecord)
	ns.Def("record?", isRecord)
	ns.Def("defprotocol*", defProtocol)
	ns.Def("make-protocol-fn", makeProtocolFn)
	ns.Def("extend-type*", extendType)
	ns.Def("satisfies?", satisfies)
	ns.Def("defmulti*", defMulti)
	ns.Def("defmethod*", defMethod)
	ns.Def("methods", methods)
	ns.Def("pr-str", prStr)
	ns.Def("prn", prn)
	ns.Def("re-find", reFind)
	ns.Def("re-matches", reMatches)
	ns.Def("re-seq", reSeq)
	ns.Def("require", requiref)
	ns.Def("find-ns", findNs)
	ns.Def("all-ns", allNs)
	ns.Def("the-ns", theNs)
	ns.Def("ns-name", nsName)

	ns.Def("peek", peek)
	ns.Def("pop", pop)

	ns.Def("iterate", iterate)
	ns.Def("repeat", repeat)
	ns.Def("lazy-seq*", lazySeq)

	ns.Def("with-meta", withMeta)
	ns.Def("meta", metaf)
	ns.Def("push-binding!", pushBinding)
	ns.Def("pop-binding!", popBinding)

	ns.Def("throw", throwf)
	ns.Def("ex-info", exInfo)
	ns.Def("ex-message", exMessage)
	ns.Def("ex-data", exData)
	ns.Def("ex-cause", exCause)

	ns.Def("delay*", delayStar)
	ns.Def("force", force)
	ns.Def("delay?", isDelay)
	ns.Def("realized?", isRealized)
	ns.Def("volatile!", volatilef)
	ns.Def("vreset!", vreset)
	ns.Def("vswap!", vswap)
	// bigint — coerce to BigInt
	bigintf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("bigint expects 1 arg")
		}
		switch v := vs[0].(type) {
		case vm.Int:
			return vm.NewBigIntFromInt64(int64(v)), nil
		case vm.Float:
			return vm.NewBigIntFromInt64(int64(v)), nil
		case *vm.BigInt:
			return v, nil
		case vm.String:
			bi, ok := vm.NewBigIntFromString(string(v))
			if !ok {
				return vm.NIL, fmt.Errorf("cannot parse bigint: %s", v)
			}
			return bi, nil
		}
		return vm.NIL, fmt.Errorf("cannot coerce %s to bigint", vs[0].Type().Name())
	})

	// bigint? — test if value is BigInt
	isBigInt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		return vm.Boolean(vm.IsBigInt(vs[0])), nil
	})

	ns.Def("bigint", bigintf)
	ns.Def("bigint?", isBigInt)
	ns.Def("transformer-seq*", transformerSeq)
	ns.Def("compare", comparef)
	ns.Def("fn?", isFn)
	ns.Def("bit-and", bitAnd)
	ns.Def("bit-or", bitOr)
	ns.Def("bit-xor", bitXor)
	ns.Def("bit-not", bitNot)
	ns.Def("bit-shift-left", bitShiftLeft)
	ns.Def("bit-shift-right", bitShiftRight)
	ns.Def("unsigned-bit-shift-right", unsignedBitShiftRight)
	ns.Def("bit-test", bitTest)
	ns.Def("bit-set", bitSet)
	ns.Def("bit-clear", bitClear)
	ns.Def("re-groups", reGroups)
	ns.Def("promise", promisef)
	ns.Def("deliver", deliver)
	ns.Def("future*", futureStar)
	ns.Def("add-watch", addWatch)
	ns.Def("remove-watch", removeWatch)
	ns.Def("alter-meta!", alterMeta)
	ns.Def("subvec", subvecf)
	ns.Def("print", printf)
	ns.Def("pr", prf)
	ns.Def("reduced", reducedf)
	ns.Def("reduced?", isReducedf)
	ns.Def("unreduced", unreduced)
	ns.Def("ensure-reduced", ensureReduced)

	// IO builtins (open, close!, read-line, write!, etc.)
	installIOBuiltins(ns)

	CoreNS = ns

	RegisterNS(ns)
}
