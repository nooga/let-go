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

type NSLoader interface {
	Load(string) *vm.Namespace
}

var nsLoader NSLoader

func SetNSLoader(loader NSLoader) {
	nsLoader = loader
}

func init() {
	nsRegistry = make(map[string]*vm.Namespace)

	installLangNS()
	installHttpNS()
	installOsNS()
	installJSONNS()
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

	// Same type check
	if a.Type() != b.Type() {
		return false
	}

	// Handle collections specially
	switch av := a.(type) {
	case vm.ArrayVector:
		bv := b.(vm.ArrayVector)
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
	default:
		// For primitives, use Go equality
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
	// Skip LazySeq to avoid forcing realization prematurely
	// (e.g. cons must not realize its tail).
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
			if !valueEquals(vs[0], vs[i]) {
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

		a, ok := vs[0].Unbox().(int)
		if !ok {
			return vm.NIL, fmt.Errorf("can't cast %s to number", vs[0].Type().Name())
		}

		b, ok := vs[1].Unbox().(int)
		if !ok {
			return vm.NIL, fmt.Errorf("can't cast %s to number", vs[1].Type().Name())
		}

		ret, err := vm.BooleanType.Box(a < b)
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

	abs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		ret := vs[0].Unbox().(int)
		if ret < 0 {
			ret = -ret
		}
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
		seq, err := seqOf(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		if seq == nil {
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
		// TODO handle namespaces
		return vm.NIL, nil // TODO is this an error?
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
			return vm.NIL, err // FIXME wrap this?
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
		if coll, ok := vs[0].(vm.Collection); ok {
			if coll.RawCount() == 0 {
				return vm.NIL, nil
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
			// Check for empty collection that implements Seq but has no elements
			if coll, ok := vs[1].(vm.Collection); ok && coll.RawCount() == 0 {
				return vm.EmptyList, nil
			}
			// eager path for small counted collections
			length := 0
			if col, ok := vs[1].(vm.Collection); ok {
				length = col.RawCount()
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
			if coll, ok := colls[i].(vm.Collection); ok && coll.RawCount() == 0 {
				return vm.EmptyList, nil
			}
			s, err := seqOf(colls[i])
			if err != nil {
				return vm.NIL, fmt.Errorf("map expected Sequable collection")
			}
			if s == nil || s == vm.EmptyList {
				return vm.EmptyList, nil
			}
			seqs[i] = s
		}
		// Check if all collections are small and counted
		minlen := math.MaxInt
		allCounted := true
		for i := range colls {
			if coll, ok := colls[i].(vm.Collection); ok {
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
		// Handle nil and empty collections
		if vs[sidx] == vm.NIL {
			if len(vs) == 3 {
				return vs[1], nil
			}
			return mfn.Invoke(nil)
		}
		// Check for empty collection first
		if coll, ok := vs[sidx].(vm.Collection); ok {
			if coll.RawCount() == 0 {
				if len(vs) == 3 {
					return vs[1], nil
				}
				return mfn.Invoke(nil)
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
		// fmt.Printf("[aliasf] in ns=%s set alias %q -> %s\n", cns.Name(), string(al), target.Name())
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
		if i, ok := vs[0].(vm.Char); ok {
			return vm.Int(int(i)), nil
		}
		return vm.NIL, fmt.Errorf("%s can't be coerced to int", vs[0])
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
	ns.Def("abs", abs)

	ns.Def("and", and)
	ns.Def("or", or)
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
	ns.Def("count", count)
	ns.Def("contains?", contains)

	ns.Def("map", mapf)
	ns.Def("mapv", mapv)
	ns.Def("reduce", reduce)
	ns.Def("concat", concat)
	ns.Def("some", some)

	ns.Def("println", printlnf)

	ns.Def("type", typef)

	ns.Def("apply*", apply)
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
	ns.Def("char", char)

	ns.Def("str", str)
	ns.Def("split", split)
	ns.Def("str-replace", strReplace)
	ns.Def("re-pattern", regex)
	// namespace utilities
	ns.Def("refer-list", referList)
	ns.Def("refer", refer)

	ns.Def("peek", peek)
	ns.Def("pop", pop)

	ns.Def("iterate", iterate)
	ns.Def("repeat", repeat)
	ns.Def("lazy-seq*", lazySeq)

	CoreNS = ns

	RegisterNS(ns)
}
