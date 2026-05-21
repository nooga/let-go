/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	crand "crypto/rand"
	_ "embed"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

var nsRegistry map[string]*vm.Namespace

var (
	tapsMu sync.Mutex
	taps   []vm.Fn
)

// nsAliases maps alternative namespace names to canonical names.
// e.g. "clojure.core" → "core", "clojure.test" → "test", "clojure.string" → "string"
// Both names resolve to the same *Namespace object.
var nsAliases = map[string]string{
	"clojure.core":   "core",
	"clojure.test":   "test",
	"clojure.string": "string",
	"clojure.set":    "set",
	"clojure.walk":   "walk",
	"clojure.edn":    "edn",
	"clojure.zip":    "zip",
	"clojure.data":   "data",
	"clojure.pprint": "pprint",
}

// resolveNSAlias returns the canonical name for a namespace.
func resolveNSAlias(name string) string {
	if canonical, ok := nsAliases[name]; ok {
		return canonical
	}
	return name
}

func externalNSName(name string) string {
	for alias, canonical := range nsAliases {
		if canonical == name {
			return alias
		}
	}
	return name
}

type NSLoader interface {
	Load(string) *vm.Namespace
}

var nsLoader NSLoader

func SetNSLoader(loader NSLoader) {
	nsLoader = loader
}

func GetNSLoader() NSLoader {
	return nsLoader
}

func init() {
	nsRegistry = make(map[string]*vm.Namespace)

	// Register global namespace lookup so qualified symbols (foo/x) work
	vm.SetNSLookup(func(name string) *vm.Namespace {
		return nsRegistry[resolveNSAlias(name)]
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
	installMathNS()
	installTermNS()
	installSyscallNS()
	installUnixNS()
	installSystemNS()
	installGogenNS()
	installDisasmNS()
	installProfileNS()
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
	return LookupOrRegisterNS(resolveNSAlias(name))
}

func RegisterNS(namespace *vm.Namespace) *vm.Namespace {
	// Auto-refer CoreNS so user code can use clojure.core symbols (defn, def,
	// require, etc.) after switching into this namespace via (in-ns 'foo) —
	// matching JVM Clojure semantics. Skip when registering CoreNS itself, and
	// when CoreNS isn't installed yet (very early init).
	if CoreNS != nil && namespace != CoreNS {
		namespace.Refer(CoreNS, "", true)
	}
	nsRegistry[resolveNSAlias(namespace.Name())] = namespace
	return namespace
}

// nsNeedsLoad tracks namespaces that were pre-registered during bytecode
// decoding but still need their precompiled chunks executed.
var nsNeedsLoad = map[string]bool{}

// MarkNSNeedsLoad flags a namespace as needing on-demand loading even though
// it already exists in the registry (created during bytecode decoding).
func MarkNSNeedsLoad(name string) {
	nsNeedsLoad[name] = true
}

// LookupNS returns a namespace if it exists, nil otherwise. Does not create.
func LookupNS(name string) *vm.Namespace {
	return nsRegistry[resolveNSAlias(name)]
}

// DefNSBare creates and registers a minimal namespace (with CoreNS refer)
// without triggering the loader. Used during bytecode decoding to create
// var homes for not-yet-loaded namespaces.
func DefNSBare(name string) *vm.Namespace {
	name = resolveNSAlias(name)
	if e := nsRegistry[name]; e != nil {
		return e
	}
	ns := vm.NewNamespace(name)
	if CoreNS != nil {
		ns.Refer(CoreNS, "", true)
	}
	nsRegistry[name] = ns
	return ns
}

func LookupOrRegisterNS(name string) *vm.Namespace {
	e := nsRegistry[name]
	if e != nil && !nsNeedsLoad[name] {
		return e
	}
	if nsLoader != nil {
		// Clear the flag before loading to prevent re-entrancy loops
		delete(nsNeedsLoad, name)
		n := nsLoader.Load(name)
		if n != nil {
			nsRegistry[name] = n
			return n
		}
	}
	// Check if loading side-effected the registry (in-ns during load creates the ns)
	if e := nsRegistry[name]; e != nil {
		delete(nsNeedsLoad, name)
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

// mapLike is any map type with Lookup + Counted + Sequable
type mapLike interface {
	vm.Lookup
	vm.Counted
	vm.Sequable
}

// setLike is any set type with Contains + Counted + Sequable
type setLike interface {
	vm.Keyed
	vm.Counted
	vm.Sequable
}

// sentinel is a unique Value used to detect missing keys.
var sentinel vm.Value = vm.Symbol("__sentinel__")

func mapEquals(a, b mapLike) bool {
	if a.RawCount() != b.RawCount() {
		return false
	}
	seq := a.Seq()
	for seq != nil && seq != vm.EmptyList {
		entry := seq.First()
		if k, v, ok := mapEntryKV(entry); ok {
			bv := b.ValueAtOr(k, sentinel)
			if bv == sentinel {
				return false
			}
			if !valueEquals(v, bv) {
				return false
			}
		}
		seq = seq.Next()
	}
	return true
}

func setEquals(a, b setLike) bool {
	if a.RawCount() != b.RawCount() {
		return false
	}
	seq := a.Seq()
	for seq != nil && seq != vm.EmptyList {
		if b.Contains(seq.First()) == vm.FALSE {
			return false
		}
		seq = seq.Next()
	}
	return true
}

func isMapType(v vm.Value) bool {
	switch v.(type) {
	case vm.Map, *vm.PersistentMap, *vm.SortedMap:
		return true
	}
	return false
}

func canConjMapEntry(v vm.Value) bool {
	if isMapType(v) {
		return true
	}
	if _, _, ok := vm.MapEntryKV(v); ok {
		return true
	}
	switch pv := v.(type) {
	case vm.PersistentVector:
		return pv.RawCount() == 2
	}
	return false
}

func isSetType(v vm.Value) bool {
	switch v.(type) {
	case vm.Set, *vm.PersistentSet, *vm.SortedSet:
		return true
	}
	return false
}

// crossMapEquals compares two map-like values of potentially different types.
func crossMapEquals(a, b vm.Value) bool {
	// Get count and lookup for both sides
	ac, al := mapCountAndLookup(a)
	bc, bl := mapCountAndLookup(b)
	if ac != bc || al == nil || bl == nil {
		return false
	}
	// Iterate a's entries and check b
	iterateMap(a, func(k, v vm.Value) bool {
		bv := bl.ValueAtOr(k, sentinel)
		if bv == sentinel || !valueEquals(v, bv) {
			ac = -1 // signal mismatch
			return false
		}
		return true
	})
	return ac >= 0
}

func mapCountAndLookup(v vm.Value) (int, vm.Lookup) {
	switch m := v.(type) {
	case vm.Map:
		return len(m), m
	case *vm.PersistentMap:
		return m.RawCount(), m
	case *vm.SortedMap:
		return m.RawCount(), m
	}
	return -1, nil
}

func iterateMap(v vm.Value, f func(k, v vm.Value) bool) {
	switch m := v.(type) {
	case vm.Map:
		for k, val := range m {
			if !f(k, val) {
				return
			}
		}
	case vm.Sequable:
		seq := m.Seq()
		for seq != nil && seq != vm.EmptyList {
			entry := seq.First()
			if k, v, ok := mapEntryKV(entry); ok {
				if !f(k, v) {
					return
				}
			}
			seq = seq.Next()
		}
	}
}

func mapEntryKV(entry vm.Value) (vm.Value, vm.Value, bool) {
	return vm.MapEntryKV(entry)
}

// crossSetEquals compares two set-like values of potentially different types.
func crossSetEquals(a, b vm.Value) bool {
	ac := setCount(a)
	bc := setCount(b)
	if ac != bc {
		return false
	}
	bk := asKeyed(b)
	if bk == nil {
		return false
	}
	match := true
	iterateSet(a, func(v vm.Value) bool {
		if bk.Contains(v) == vm.FALSE {
			match = false
			return false
		}
		return true
	})
	return match
}

func setCount(v vm.Value) int {
	switch s := v.(type) {
	case vm.Set:
		return len(s)
	case *vm.PersistentSet:
		return s.RawCount()
	case *vm.SortedSet:
		return s.RawCount()
	}
	return -1
}

func asKeyed(v vm.Value) vm.Keyed {
	if k, ok := v.(vm.Keyed); ok {
		return k
	}
	return nil
}

func iterateSet(v vm.Value, f func(vm.Value) bool) {
	switch s := v.(type) {
	case vm.Set:
		for k := range s {
			if !f(k) {
				return
			}
		}
	case vm.Sequable:
		seq := s.Seq()
		for seq != nil && seq != vm.EmptyList {
			if !f(seq.First()) {
				return
			}
			seq = seq.Next()
		}
	}
}

// valueEquals performs deep equality comparison for Clojure semantics
func valueEquals(a, b vm.Value) bool {
	// Handle nil
	if isNilValue(a) && isNilValue(b) {
		return true
	}
	if isNilValue(a) || isNilValue(b) {
		return false
	}
	if vm.IsNumber(a) && vm.IsNumber(b) {
		return vm.NumEq(a, b)
	}
	switch av := a.(type) {
	case *vm.Range:
		if av == b {
			return true
		}
	case *vm.InfiniteRange:
		if av == b {
			return true
		}
	}

	// Allow cross-type comparison for numbers and vectors
	if a.Type() != b.Type() {
		if vm.IsNumber(a) && vm.IsNumber(b) {
			return vm.NumEq(a, b)
		}
		// Cross-type vector equality (ArrayVector vs PersistentVector)
		if eq, ok := a.(interface{ Equals(vm.Value) bool }); ok {
			if eq.Equals(b) {
				return true
			}
		}
		if eq, ok := b.(interface{ Equals(vm.Value) bool }); ok {
			if eq.Equals(a) {
				return true
			}
		}
		// Cross-type sequential equality: compare any two seq-like
		// values element-by-element (lists, vectors, cons, lazy seqs,
		// ranges — but not maps or sets).
		if isSequentialType(a) && isSequentialType(b) {
			as := toSeq(a)
			bs := toSeq(b)
			for as != nil && bs != nil {
				if !valueEquals(as.First(), bs.First()) {
					return false
				}
				as, bs = as.Next(), bs.Next()
			}
			return as == nil && bs == nil
		}
		// Cross-type map equality (SortedMap vs Map vs PersistentMap)
		if isMapType(a) && isMapType(b) {
			return crossMapEquals(a, b)
		}
		// Cross-type set equality (SortedSet vs Set vs PersistentSet)
		if isSetType(a) && isSetType(b) {
			return crossSetEquals(a, b)
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
			if !nilListEquivalent(av[i], bv[i]) && !valueEquals(av[i], bv[i]) {
				return false
			}
		}
		return true
	case *vm.List:
		// b could be any Seq-like type (List, Cons, ArrayVectorSeq, etc.)
		if av == vm.EmptyList {
			if bl, ok := b.(*vm.List); ok {
				return bl == vm.EmptyList || bl.RawCount() == 0
			}
			if isSequentialType(b) {
				return toSeq(b) == nil
			}
			return false
		}
		if bl, ok := b.(*vm.List); ok && bl == vm.EmptyList {
			return false
		}
		bs, ok := b.(vm.Seq)
		if !ok {
			return false
		}
		as := vm.Seq(av)
		for as != nil && bs != nil {
			if !listElementEquals(as.First(), bs.First()) {
				return false
			}
			as, bs = as.Next(), bs.Next()
		}
		return as == nil && bs == nil
	case vm.Map:
		if bm, ok := b.(vm.Map); ok {
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
		}
		// Cross-type: vm.Map vs SortedMap/PersistentMap
		if bm, ok := b.(mapLike); ok {
			if len(av) != bm.RawCount() {
				return false
			}
			for k, v := range av {
				bv := bm.ValueAtOr(k, sentinel)
				if bv == sentinel || !valueEquals(v, bv) {
					return false
				}
			}
			return true
		}
		return false
	case *vm.PersistentMap:
		if bm, ok := b.(*vm.PersistentMap); ok {
			return av.Equals(bm)
		}
		if bm, ok := b.(*vm.SortedMap); ok {
			return mapEquals(av, bm)
		}
		return false
	case *vm.SortedMap:
		if bm, ok := b.(*vm.SortedMap); ok {
			return mapEquals(av, bm)
		}
		if bm, ok := b.(*vm.PersistentMap); ok {
			return mapEquals(av, bm)
		}
		if bm, ok := b.(vm.Map); ok {
			if av.RawCount() != len(bm) {
				return false
			}
			for k, v := range bm {
				sv := av.ValueAtOr(k, sentinel)
				if sv == sentinel || !valueEquals(v, sv) {
					return false
				}
			}
			return true
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
		if bs, ok := b.(*vm.PersistentSet); ok {
			return setEquals(av, bs)
		}
		if bs, ok := b.(*vm.SortedSet); ok {
			return setEquals(av, bs)
		}
		return false
	case *vm.SortedSet:
		if bs, ok := b.(*vm.SortedSet); ok {
			return setEquals(av, bs)
		}
		if bs, ok := b.(*vm.PersistentSet); ok {
			return setEquals(av, bs)
		}
		return false
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
		// Same-type seq comparison (Range, Cons, LazySeq, etc.)
		if isSequentialType(a) && isSequentialType(b) {
			as := toSeq(a)
			bs := toSeq(b)
			for as != nil && bs != nil {
				if !valueEquals(as.First(), bs.First()) {
					return false
				}
				as, bs = as.Next(), bs.Next()
			}
			return as == nil && bs == nil
		}
		return a == b
	}
}

func isNilValue(v vm.Value) bool {
	return v == nil || v == vm.NIL
}

func isNaNValue(v vm.Value) bool {
	if f, ok := v.(vm.Float); ok {
		return math.IsNaN(float64(f))
	}
	if f, ok := v.(vm.Float32); ok {
		return math.IsNaN(float64(f))
	}
	return false
}

func listElementEquals(a, b vm.Value) bool {
	if isNaNValue(a) && isNaNValue(b) {
		return true
	}
	return valueEquals(a, b)
}

func nilListEquivalent(a, b vm.Value) bool {
	return (a == vm.NIL && b == vm.EmptyList) || (a == vm.EmptyList && b == vm.NIL)
}

// isSequentialType returns true for types that participate in cross-type
// sequential equality (lists, vectors, cons, lazy seqs, ranges — not maps/sets).
func isSequentialType(v vm.Value) bool {
	if _, ok := v.(vm.String); ok {
		return false
	}
	switch v.(type) {
	case vm.ArrayVector, vm.PersistentVector, *vm.PersistentVector, *vm.List, *vm.Cons, *vm.LazySeq:
		return true
	case *vm.PersistentMap, vm.Map, *vm.PersistentSet, vm.Set, *vm.SortedMap, *vm.SortedSet:
		return false
	}
	// Range and other Seq implementations
	_, isSeq := v.(vm.Seq)
	return isSeq
}

// toSeq converts a sequential value to a Seq for element-by-element comparison.
func toSeq(v vm.Value) vm.Seq {
	if v == vm.NIL || v == vm.EmptyList {
		return nil
	}
	if ls, ok := v.(*vm.LazySeq); ok {
		s := ls.Resolve()
		if s == vm.EmptyList {
			return nil
		}
		return s
	}
	if c, ok := v.(vm.Counted); ok && c.RawCount() == 0 {
		return nil
	}
	if sq, ok := v.(vm.Sequable); ok {
		s := sq.Seq()
		if s == vm.EmptyList {
			return nil
		}
		return s
	}
	if s, ok := v.(vm.Seq); ok {
		return s
	}
	return nil
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

// forceSeq fully resolves a LazySeq chain to a concrete seq (or nil if empty).
// Non-LazySeq inputs pass through unchanged.
func forceSeq(s vm.Seq) vm.Seq {
	if s == nil {
		return nil
	}
	if ls, ok := s.(*vm.LazySeq); ok {
		return ls.Resolve()
	}
	return s
}

func mapLazy1(f vm.Fn, s vm.Seq) vm.Seq {
	if s == nil {
		return nil
	}
	captured := s
	thunk, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		// Force the captured seq before reading First(): a LazySeq that
		// resolves to empty must produce EmptyList, not f(nil).
		head := forceSeq(captured)
		if head == nil {
			return vm.EmptyList, nil
		}
		v, err := f.Invoke([]vm.Value{head.First()})
		if err != nil {
			return vm.NIL, err
		}
		rest := head.Next()
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
		// Force every captured seq before reading First(); if any is empty,
		// the whole multi-coll map ends.
		heads := make([]vm.Seq, len(captured))
		for i, s := range captured {
			h := forceSeq(s)
			if h == nil {
				return vm.EmptyList, nil
			}
			heads[i] = h
		}
		fargs := make([]vm.Value, len(heads))
		nexts := make([]vm.Seq, len(heads))
		for i, s := range heads {
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

func fnComparator(comp vm.Fn) vm.Comparator {
	return func(a, b vm.Value) (int, error) {
		r, err := comp.Invoke([]vm.Value{a, b})
		if err != nil {
			return 0, err
		}
		if br, ok := r.(vm.Boolean); ok {
			if bool(br) {
				return -1, nil
			}
			r, err = comp.Invoke([]vm.Value{b, a})
			if err != nil {
				return 0, err
			}
			if reverse, ok := r.(vm.Boolean); ok && bool(reverse) {
				return 1, nil
			}
			return 0, nil
		}
		n, ok := vm.ToInt(r)
		if !ok {
			return 0, fmt.Errorf("comparator returned non-numeric value %s", r.Type().Name())
		}
		switch {
		case n < 0:
			return -1, nil
		case n > 0:
			return 1, nil
		default:
			return 0, nil
		}
	}
}

func invokeMethodFallback(rec vm.Value, name vm.Symbol, args []vm.Value, originalErr error) (vm.Value, error) {
	if isCompatChecker(rec) && len(args) == 1 {
		switch name {
		case "isLong":
			if _, boxed := args[0].(*vm.DTypeInstance); boxed {
				return vm.FALSE, nil
			}
			_, ok := args[0].(vm.Int)
			return vm.Boolean(ok), nil
		case "isDouble":
			if _, boxed := args[0].(*vm.DTypeInstance); boxed {
				return vm.FALSE, nil
			}
			_, ok := args[0].(vm.Float)
			return vm.Boolean(ok), nil
		}
	}
	if name == "reduce" && len(args) == 1 {
		if v := CoreNS.Lookup(vm.Symbol("reduce")); v != vm.NIL {
			if methodVar, ok := v.(*vm.Var); ok {
				if fn, ok := methodVar.Deref().(vm.Fn); ok {
					return fn.Invoke([]vm.Value{args[0], rec})
				}
			}
		}
	}
	if name == "reduce" && len(args) == 2 {
		if v := CoreNS.Lookup(vm.Symbol("reduce")); v != vm.NIL {
			if methodVar, ok := v.(*vm.Var); ok {
				if fn, ok := methodVar.Deref().(vm.Fn); ok {
					return fn.Invoke([]vm.Value{args[0], args[1], rec})
				}
			}
		}
	}
	if v := CurrentNS.Deref().(*vm.Namespace).Lookup(name); v != vm.NIL {
		if methodVar, ok := v.(*vm.Var); ok {
			if fn, ok := methodVar.Deref().(vm.Fn); ok {
				return fn.Invoke(append([]vm.Value{rec}, args...))
			}
		}
	}
	for _, ns := range nsRegistry {
		if v := ns.LookupLocal(name); v != nil {
			if fn, ok := v.Deref().(vm.Fn); ok {
				return fn.Invoke(append([]vm.Value{rec}, args...))
			}
		}
	}
	return vm.NIL, originalErr
}

func isCompatChecker(v vm.Value) bool {
	inst, ok := v.(*vm.DTypeInstance)
	if !ok {
		return false
	}
	return strings.HasSuffix(inst.DType().Name(), "Checker")
}

func roundBigDecimalValue(v *vm.BigDecimal, precision int, modeName string) (vm.Value, error) {
	if precision <= 0 {
		return vm.NIL, fmt.Errorf("precision must be positive")
	}
	f, _ := v.Val().Float64()
	if f == 0 {
		return v, nil
	}
	sign := 1
	if f < 0 {
		sign = -1
		f = -f
	}
	exp := int(math.Floor(math.Log10(f)))
	scaleExp := precision - 1 - exp
	rat, ok := new(big.Rat).SetString(v.Val().Text('f', -1))
	if !ok {
		return vm.NIL, fmt.Errorf("cannot round BigDecimal %s", v)
	}
	if sign < 0 {
		rat.Abs(rat)
	}
	scaled := new(big.Rat).Mul(rat, pow10Rat(scaleExp))
	rounded, err := roundRatToInt(scaled, sign, modeName)
	if err != nil {
		return vm.NIL, err
	}
	result := new(big.Rat).SetInt(rounded)
	result.Quo(result, pow10Rat(scaleExp))
	if sign < 0 {
		result.Neg(result)
	}
	return vm.NewBigDecimal(new(big.Float).SetPrec(vm.BigDecimalPrecConst).SetRat(result)), nil
}

func pow10Rat(exp int) *big.Rat {
	if exp >= 0 {
		return new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exp)), nil))
	}
	return new(big.Rat).Inv(new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-exp)), nil)))
}

func roundRatToInt(r *big.Rat, sign int, modeName string) (*big.Int, error) {
	q := new(big.Int).Quo(r.Num(), r.Denom())
	rem := new(big.Int).Rem(r.Num(), r.Denom())
	if rem.Sign() == 0 {
		return q, nil
	}
	ceil := new(big.Int).Add(q, big.NewInt(1))
	switch modeName {
	case "up":
		return ceil, nil
	case "down":
		return q, nil
	case "ceiling":
		if sign > 0 {
			return ceil, nil
		}
		return q, nil
	case "floor":
		if sign < 0 {
			return ceil, nil
		}
		return q, nil
	case "unnecessary":
		return nil, fmt.Errorf("rounding necessary")
	case "half-up", "half-down", "half-even":
		twiceRem := new(big.Int).Mul(rem, big.NewInt(2))
		cmp := twiceRem.Cmp(r.Denom())
		if cmp > 0 || (cmp == 0 && modeName == "half-up") {
			return ceil, nil
		}
		if cmp == 0 && modeName == "half-even" && q.Bit(0) == 1 {
			return ceil, nil
		}
		return q, nil
	default:
		return nil, fmt.Errorf("unknown rounding mode %s", modeName)
	}
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
		if len(vs) == 1 {
			return vm.NumDiv(vm.MakeInt(1), vs[0])
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

	// Apostrophe arithmetic: identical to + - * but promotes to BigInt on
	// int64 overflow instead of wrapping silently.
	plusP, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.MakeInt(0), nil
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		acc := vs[0]
		for i := 1; i < len(vs); i++ {
			var err error
			acc, err = vm.NumAddP(acc, vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return acc, nil
	})

	mulP, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.MakeInt(1), nil
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		acc := vs[0]
		for i := 1; i < len(vs); i++ {
			var err error
			acc, err = vm.NumMulP(acc, vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return acc, nil
	})

	subP, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vm.NumNegP(vs[0])
		}
		acc := vs[0]
		for i := 1; i < len(vs); i++ {
			var err error
			acc, err = vm.NumSubP(acc, vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return acc, nil
	})

	// unchecked-* family — silent wrap on int64 overflow.
	//
	// Mirrors Clojure's clojure.core/unchecked-add, unchecked-subtract,
	// unchecked-multiply, unchecked-negate, unchecked-divide-int. Each is
	// strict-arity (2-ary for binary ops, 1-ary for unary) to match Clojure's
	// inliner signatures.
	//
	// Required for porting hash functions, splittable RNGs, and any other
	// code that needs modular int64 / u64 semantics. Pair with the existing
	// unsigned-bit-shift-right to build xxh3-64, splitmix64, SipHash, etc.
	//
	// Float and BigInt inputs are coerced to int64 (matching Clojure's
	// long-arithmetic semantics). Use `+`/`-`/`*` for checked arithmetic,
	// or `+'`/`-'`/`*'` for BigInt-promoting arithmetic.

	uncheckedAdd, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, ok := vm.ToInt(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-add expected integer, got %s", vs[0].Type().Name())
		}
		b, ok := vm.ToInt(vs[1])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-add expected integer, got %s", vs[1].Type().Name())
		}
		return vm.MakeInt(int(int64(a) + int64(b))), nil
	})
	if err != nil {
		panic(err)
	}

	uncheckedSubtract, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, ok := vm.ToInt(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-subtract expected integer, got %s", vs[0].Type().Name())
		}
		b, ok := vm.ToInt(vs[1])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-subtract expected integer, got %s", vs[1].Type().Name())
		}
		return vm.MakeInt(int(int64(a) - int64(b))), nil
	})
	if err != nil {
		panic(err)
	}

	uncheckedMultiply, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, ok := vm.ToInt(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-multiply expected integer, got %s", vs[0].Type().Name())
		}
		b, ok := vm.ToInt(vs[1])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-multiply expected integer, got %s", vs[1].Type().Name())
		}
		return vm.MakeInt(int(int64(a) * int64(b))), nil
	})
	if err != nil {
		panic(err)
	}

	uncheckedNegate, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, ok := vm.ToInt(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-negate expected integer, got %s", vs[0].Type().Name())
		}
		// Note: -Long/MIN_VALUE wraps to Long/MIN_VALUE in 2's-complement int64.
		// This matches Clojure's unchecked-negate behavior.
		return vm.MakeInt(int(-int64(a))), nil
	})
	if err != nil {
		panic(err)
	}

	uncheckedDivideInt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, ok := vm.ToInt(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-divide-int expected integer, got %s", vs[0].Type().Name())
		}
		b, ok := vm.ToInt(vs[1])
		if !ok {
			return vm.NIL, fmt.Errorf("unchecked-divide-int expected integer, got %s", vs[1].Type().Name())
		}
		if b == 0 {
			return vm.NIL, fmt.Errorf("divide by zero")
		}
		// Note: Long/MIN_VALUE / -1 overflows in 2's-complement and throws.
		// This matches Clojure's unchecked-divide-int (and JVM IDIV) behavior.
		if int64(a) == math.MinInt64 && int64(b) == -1 {
			return vm.NIL, fmt.Errorf("integer overflow")
		}
		return vm.MakeInt(int(int64(a) / int64(b))), nil
	})
	if err != nil {
		panic(err)
	}

	uncheckedLong, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.Int:
			return v, nil
		case *vm.BigInt:
			// Low-64 of two's-complement representation, matching Clojure JVM.
			// big.Int doesn't expose two's-complement directly; mask with 2^64,
			// then reinterpret the low-64 bits as signed int64.
			mask := new(big.Int).Lsh(big.NewInt(1), 64)
			lo := new(big.Int).Mod(v.Val(), mask)
			return vm.MakeInt(int(int64(lo.Uint64()))), nil
		case vm.Float:
			return vm.MakeInt(int(int64(float64(v)))), nil
		default:
			return vm.NIL, fmt.Errorf("unchecked-long expected integer or float, got %s", vs[0].Type().Name())
		}
	})
	if err != nil {
		panic(err)
	}

	uncheckedInt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.Int:
			return vm.MakeInt(int(int32(int64(v)))), nil
		case *vm.BigInt:
			// Low-32 of two's-complement representation, matching Clojure JVM.
			// Mask with 2^64, take low-32 bits as int32, then sign-extend.
			mask := new(big.Int).Lsh(big.NewInt(1), 64)
			lo := new(big.Int).Mod(v.Val(), mask)
			return vm.MakeInt(int(int32(lo.Uint64()))), nil
		case vm.Float:
			return vm.MakeInt(int(int32(float64(v)))), nil
		default:
			return vm.NIL, fmt.Errorf("unchecked-int expected integer or float, got %s", vs[0].Type().Name())
		}
	})
	if err != nil {
		panic(err)
	}

	uncheckedShort, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.Int:
			return vm.MakeInt(int(int16(int64(v)))), nil
		case *vm.BigInt:
			// Low-16 of two's-complement representation, matching Clojure JVM.
			mask := new(big.Int).Lsh(big.NewInt(1), 64)
			lo := new(big.Int).Mod(v.Val(), mask)
			return vm.MakeInt(int(int16(lo.Uint64()))), nil
		case vm.Float:
			return vm.MakeInt(int(int16(float64(v)))), nil
		default:
			return vm.NIL, fmt.Errorf("unchecked-short expected integer or float, got %s", vs[0].Type().Name())
		}
	})
	if err != nil {
		panic(err)
	}

	uncheckedByte, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.Int:
			return vm.MakeInt(int(int8(int64(v)))), nil
		case *vm.BigInt:
			// Low-8 of two's-complement representation, matching Clojure JVM.
			mask := new(big.Int).Lsh(big.NewInt(1), 64)
			lo := new(big.Int).Mod(v.Val(), mask)
			return vm.MakeInt(int(int8(lo.Uint64()))), nil
		case vm.Float:
			return vm.MakeInt(int(int8(float64(v)))), nil
		default:
			return vm.NIL, fmt.Errorf("unchecked-byte expected integer or float, got %s", vs[0].Type().Name())
		}
	})
	if err != nil {
		panic(err)
	}

	uncheckedChar, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.Int:
			return vm.Char(rune(uint16(int64(v)))), nil
		case *vm.BigInt:
			// Low-16 of two's-complement representation, then uint16 → Char (rune).
			mask := new(big.Int).Lsh(big.NewInt(1), 64)
			lo := new(big.Int).Mod(v.Val(), mask)
			return vm.Char(rune(uint16(lo.Uint64()))), nil
		case vm.Float:
			return vm.Char(rune(uint16(int64(float64(v))))), nil
		default:
			return vm.NIL, fmt.Errorf("unchecked-char expected integer or float, got %s", vs[0].Type().Name())
		}
	})
	if err != nil {
		panic(err)
	}

	uncheckedDouble, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.Int:
			return vm.Float(float64(int64(v))), nil
		case *vm.BigInt:
			// Convert BigInt to float64 via big.Float (handles overflow → ±Inf).
			f, _ := new(big.Float).SetInt(v.Val()).Float64()
			return vm.Float(f), nil
		case vm.Float:
			return v, nil
		default:
			return vm.NIL, fmt.Errorf("unchecked-double expected numeric, got %s", vs[0].Type().Name())
		}
	})
	if err != nil {
		panic(err)
	}

	uncheckedFloat, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.Int:
			// Cast through float32 to introduce float32-precision loss + overflow → ±Inf.
			return vm.Float(float64(float32(int64(v)))), nil
		case *vm.BigInt:
			f, _ := new(big.Float).SetInt(v.Val()).Float64()
			return vm.Float(float64(float32(f))), nil
		case vm.Float:
			return vm.Float(float64(float32(float64(v)))), nil
		default:
			return vm.NIL, fmt.Errorf("unchecked-float expected numeric, got %s", vs[0].Type().Name())
		}
	})
	if err != nil {
		panic(err)
	}

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

	notEq, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.FALSE, nil
		}
		if len(vs) == 2 && isNaNValue(vs[0]) && isNaNValue(vs[1]) {
			return vm.FALSE, nil
		}
		eq, err := equals.(vm.Fn).Invoke(vs)
		if err != nil {
			return vm.NIL, err
		}
		return vm.Boolean(!vm.IsTruthy(eq)), nil
	})

	gt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vm.TRUE, nil
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
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vm.TRUE, nil
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
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vm.TRUE, nil
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
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 1 {
			return vm.TRUE, nil
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

	complement, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		f, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("complement expected Fn")
		}
		wrapped, wrapErr := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
			v, err := f.Invoke(args)
			if err != nil {
				return vm.NIL, err
			}
			return vm.Boolean(!vm.IsTruthy(v)), nil
		})
		if wrapErr != nil {
			return vm.NIL, wrapErr
		}
		return wrapped, nil
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
	hashMap, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs)%2 != 0 {
			return vm.NIL, fmt.Errorf("hash-map requires an even number of arguments, got %d", len(vs))
		}
		return vm.NewMap(vs), nil
	})
	arrayMap, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs)%2 != 0 {
			return vm.NIL, fmt.Errorf("array-map requires an even number of arguments, got %d", len(vs))
		}
		return vm.NewArrayMap(vs), nil
	})
	hashSet, err := vm.NativeFnType.WrapNoErr(vm.NewSet)

	sortedMap, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs)%2 != 0 {
			return vm.NIL, fmt.Errorf("sorted-map requires even number of arguments, got %d", len(vs))
		}
		return vm.NewSortedMap(nil, vs), nil
	})

	sortedSet, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NewSortedSet(nil, vs), nil
	})

	sortedMapBy, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("sorted-map-by requires a comparator")
		}
		comp, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("sorted-map-by first arg must be a function")
		}
		kvs := vs[1:]
		if len(kvs)%2 != 0 {
			return vm.NIL, fmt.Errorf("sorted-map-by requires even number of key-value arguments")
		}
		return vm.NewSortedMap(fnComparator(comp), kvs), nil
	})

	sortedSetBy, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("sorted-set-by requires a comparator")
		}
		comp, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("sorted-set-by first arg must be a function")
		}
		return vm.NewSortedSet(fnComparator(comp), vs[1:]), nil
	})

	vec, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL || vs[0] == vm.EmptyList {
			return vm.ArrayVector{}, nil
		}
		// Empty string → empty vector
		if s, ok := vs[0].(vm.String); ok && len(string(s)) == 0 {
			return vm.ArrayVector{}, nil
		}
		if v, ok := vs[0].(vm.ArrayVector); ok {
			return v, nil
		}
		if a, ok := vs[0].(*vm.TypedArray); ok && a.Kind() == vm.ArrayObject {
			return vm.ArrayVector(a.Unbox().([]vm.Value)), nil
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
			// Infinite range: (range) -> lazy seq 0, 1, 2, ...
			return vm.NewInfiniteRange(0, 1), nil
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
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if len(vs) == 2 {
			// (keyword ns name) — both must be strings (or nil ns)
			var nsStr, nameStr string
			if vs[0] != vm.NIL {
				switch n := vs[0].(type) {
				case vm.String:
					nsStr = string(n)
				default:
					return vm.NIL, fmt.Errorf("keyword namespace must be a string, got %s", vs[0].Type())
				}
			}
			switch n := vs[1].(type) {
			case vm.String:
				nameStr = string(n)
			default:
				return vm.NIL, fmt.Errorf("keyword name must be a string, got %s", vs[1].Type())
			}
			if nsStr == "" && vs[0] == vm.NIL {
				return vm.Keyword(nameStr), nil
			}
			return vm.Keyword(nsStr + "/" + nameStr), nil
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
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
			case vm.Keyword:
				return string(s), true
			default:
				return "", false
			}
		}
		if len(vs) == 1 {
			if vs[0] == vm.NIL {
				return vm.NIL, fmt.Errorf("symbol expected String, Symbol, Keyword, or Var")
			}
			if v, ok := vs[0].(*vm.Var); ok {
				return vm.Symbol(externalNSName(v.NS()) + "/" + v.VarName()), nil
			}
			if s, ok := toStr(vs[0]); ok {
				return vm.Symbol(s), nil
			}
			return vm.NIL, fmt.Errorf("symbol expected String or Symbol")
		}
		nsStr := ""
		if vs[0] != vm.NIL {
			ns, ok := vs[0].(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf("symbol expected String namespace")
			}
			nsStr = string(ns)
		}
		name, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("symbol expected String name")
		}
		nameStr := string(name)
		if nsStr == "" && vs[0] == vm.NIL {
			return vm.Symbol(nameStr), nil
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
			if ret == vm.NIL {
				return vm.NIL, fmt.Errorf("assoc failed for key %s", vs[i].String())
			}
		}
		return ret, nil
	})

	dissoc, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.NIL, fmt.Errorf("wrong number of arguments 0")
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		coll, ok := vs[0].(vm.Associative)
		if !ok {
			return vm.NIL, fmt.Errorf("dissoc expected Associative")
		}
		ret := coll
		for i := 1; i < len(vs); i++ {
			ret = ret.Dissoc(vs[i])
			if vs[0] != vm.NIL && ret == vm.NIL {
				return vm.NIL, fmt.Errorf("dissoc failed for key %s", vs[i].String())
			}
		}
		return ret, nil
	})

	update, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		// Treat nil as empty map
		if vs[0] == vm.NIL {
			vs[0] = vm.EmptyPersistentMap
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
		fn, ok := vs[2].(vm.Fn)
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
			return vm.NewCons(elem, nil), nil
		}
		seq, err := seqOf(vs[1])
		if err != nil {
			return vm.NIL, fmt.Errorf("cons expected Seq")
		}
		if seq == nil {
			return vm.NewCons(elem, nil), nil
		}
		return vm.NewCons(elem, seq), nil
	})

	conj, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.ArrayVector{}, nil
		}
		if len(vs) == 1 {
			return vs[0], nil
		}
		var seq vm.Collection
		if vs[0] == vm.NIL {
			seq = vm.EmptyList
		} else {
			if _, ok := vs[0].(vm.String); ok {
				return vm.NIL, fmt.Errorf("conj expected Collection")
			}
			var ok bool
			seq, ok = vs[0].(vm.Collection)
			if !ok {
				if s, ok := vs[0].(vm.Seq); ok {
					for i := 1; i < len(vs); i++ {
						if seq == nil {
							seq = vm.NewCons(vs[i], s)
						} else {
							seq = seq.Conj(vs[i])
						}
					}
					return seq, nil
				}
				return vm.NIL, fmt.Errorf("conj expected Collection")
			}
		}
		for i := 1; i < len(vs); i++ {
			if isMapType(seq) && !canConjMapEntry(vs[i]) {
				return vm.NIL, fmt.Errorf("conj expected map entry")
			}
			seq = seq.Conj(vs[i])
		}
		return seq, nil
	})

	disj, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
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
		case *vm.SortedSet:
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
		if vs[0] == vm.NIL {
			return vm.FALSE, nil
		}
		// Keyed types: maps and sets
		if s, ok := vs[0].(vm.Keyed); ok {
			return s.Contains(vs[1]), nil
		}
		// Vectors: contains? checks if index exists
		if idx, ok := vs[1].(vm.Int); ok {
			i := int(idx)
			if c, ok := vs[0].(vm.Counted); ok {
				return vm.Boolean(i >= 0 && i < c.RawCount()), nil
			}
		}
		// Strings: contains? checks if index exists (must be integer key)
		if s, ok := vs[0].(vm.String); ok {
			if idx, ok := vs[1].(vm.Int); ok {
				return vm.Boolean(int(idx) >= 0 && int(idx) < len([]rune(string(s)))), nil
			}
			return vm.NIL, fmt.Errorf("contains? on a string requires an integer key, got %s", vs[1].Type().Name())
		}
		return vm.FALSE, nil
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
		var n vm.Seq
		if sqbl, ok := vs[0].(vm.Sequable); ok {
			n = sqbl.Seq()
		} else if s, ok := vs[0].(vm.Seq); ok {
			n = s
		} else {
			return vm.NIL, fmt.Errorf("seq expected Seqable")
		}
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
		switch vs[0].(type) {
		case *vm.Nil, vm.String, *vm.TypedArray,
			vm.ArrayVector, *vm.PersistentVector,
			*vm.PersistentMap, *vm.PersistentSet,
			*vm.SortedSet, *vm.SortedMap:
			return vm.FALSE, nil
		}
		_, ok := vs[0].(vm.Seq)
		return vm.Boolean(ok), nil
	})

	isList, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		_, ok := vs[0].(*vm.List)
		return vm.Boolean(ok), nil
	})

	isColl, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch vs[0].(type) {
		case *vm.Nil, vm.String, *vm.TypedArray:
			return vm.FALSE, nil
		}
		_, ok := vs[0].(vm.Collection)
		return vm.Boolean(ok), nil
	})

	empty, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL || vs[0].Type() == vm.StringType {
			return vm.NIL, nil
		}
		if _, ok := vs[0].(*vm.InfiniteRange); ok {
			return vm.EmptyList, nil
		}
		if _, ok := vs[0].(*vm.Record); ok {
			return vm.NIL, fmt.Errorf("empty is not supported on records")
		}
		coll, ok := vs[0].(vm.Collection)
		if !ok {
			return vm.NIL, nil
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

	keyf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if e, ok := vs[0].(vm.MapEntry); ok {
			return e.Key, nil
		}
		return vm.NIL, fmt.Errorf("key expects map entry")
	})

	valf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if e, ok := vs[0].(vm.MapEntry); ok {
			return e.Value, nil
		}
		return vm.NIL, fmt.Errorf("val expects map entry")
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
		hasDefault := vl == 3
		var notFound vm.Value = vm.NIL
		if hasDefault {
			notFound = vs[2]
		}
		if vs[0] == vm.NIL {
			if hasDefault {
				return notFound, nil
			}
			return vm.NIL, nil
		}
		// Fast path: Lookup types (ArrayVector, PersistentVector, String)
		if l, ok := vs[0].(vm.Lookup); ok {
			if hasDefault {
				return l.ValueAtOr(vm.Int(n), notFound), nil
			}
			if c, ok := vs[0].(vm.Counted); ok && (n < 0 || n >= c.RawCount()) {
				return vm.NIL, fmt.Errorf("nth index out of bounds")
			}
			return l.ValueAt(vm.Int(n)), nil
		}
		// Seq path: linear walk
		if n < 0 {
			if hasDefault {
				return notFound, nil
			}
			return vm.NIL, fmt.Errorf("nth index out of bounds")
		}
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
		if !hasDefault {
			return vm.NIL, fmt.Errorf("nth index out of bounds")
		}
		return notFound, nil
	})

	count, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.MakeInt(0), nil
		}
		if s, ok := vs[0].(vm.String); ok {
			return vm.MakeInt(len([]rune(string(s)))), nil
		}
		seq, ok := vs[0].(vm.Counted)
		if !ok {
			return vm.NIL, fmt.Errorf("count expected Counted")
		}
		return seq.Count(), nil
	})

	// map builtin: always lazy. Clojure semantics require laziness on all
	// inputs — small counted collections must still defer realization so
	// consumers (rose trees, short-circuit take, side-effecting f) work.
	// Use mapv for eager realization to a vector.
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
		// 3-arg form with nil coll: return init regardless of fn
		if len(vs) == 3 && vs[2] == vm.NIL {
			return vs[1], nil
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

		return vm.NIL, nil
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
			b.WriteString(strValue(vs[i]))
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

	excludeInCurrentNs, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		cns := CurrentNS.Deref().(*vm.Namespace)
		for _, v := range vs {
			sym, ok := v.(vm.Symbol)
			if !ok {
				return vm.NIL, fmt.Errorf("exclude-in-current-ns expected Symbol, got %s", v.Type().Name())
			}
			cns.Exclude(string(sym))
		}
		return vm.NIL, nil
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
		name, ok := vs[1].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("method-invoke expected Symbol")
		}
		rec, ok := vs[0].(vm.Receiver)
		if !ok {
			return invokeMethodFallback(vs[0], name, vs[2:], fmt.Errorf("method-invoke expected Receiver"))
		}
		result, err := rec.InvokeMethod(name, vs[2:])
		if err == nil {
			return result, nil
		}
		return invokeMethodFallback(rec, name, vs[2:], err)
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
		if vs[0] == vm.NIL {
			return vm.NIL, fmt.Errorf("name called on nil")
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
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if (len(vs)-1)%2 != 0 {
			return vm.NIL, fmt.Errorf("atom options must be key/value pairs")
		}
		var meta vm.Value
		var validator vm.Fn
		for i := 1; i < len(vs); i += 2 {
			switch vs[i] {
			case vm.Keyword("meta"):
				if vs[i+1] != vm.NIL && !isMapType(vs[i+1]) {
					return vm.NIL, fmt.Errorf("atom :meta must be nil or map")
				}
				meta = vs[i+1]
			case vm.Keyword("validator"):
				if vs[i+1] == vm.NIL {
					validator = nil
					continue
				}
				fn, ok := vs[i+1].(vm.Fn)
				if !ok {
					return vm.NIL, fmt.Errorf("atom :validator must be nil or function")
				}
				validator = fn
			}
		}
		return vm.NewAtomWithMetaValidator(vs[0], meta, validator)
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
		return at.Reset(vs[1])
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
		if _, err := at.Reset(vs[1]); err != nil {
			return vm.NIL, err
		}
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
		if err != nil {
			return vm.NIL, nil // Clojure returns nil for unparseable
		}
		return vm.MakeInt(i), nil
	})

	max, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		m := vs[0]
		if isNaNValue(m) {
			return m, nil
		}
		for i := 1; i < len(vs); i++ {
			if isNaNValue(vs[i]) {
				return vs[i], nil
			}
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
		if isNaNValue(m) {
			return m, nil
		}
		for i := 1; i < len(vs); i++ {
			if isNaNValue(vs[i]) {
				return vs[i], nil
			}
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

	// compareValues delegates to the vm package's DefaultCompare
	compareValues := vm.DefaultCompare

	comparef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		c, err := compareValues(vs[0], vs[1])
		if err != nil {
			return vm.NIL, err
		}
		return vm.MakeInt(c), nil
	})

	sort, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		var comp vm.Comparator
		var coll vm.Collection
		var ok bool
		if len(vs) == 2 {
			compFn, ok := vs[0].(vm.Fn)
			if !ok {
				return vm.NIL, fmt.Errorf("sort expected a comparator function")
			}
			comp = fnComparator(compFn)
			coll, ok = vs[1].(vm.Collection)
			if !ok {
				return vm.NIL, fmt.Errorf("sort expected a Collection")
			}
		} else {
			comp = nil // use default compare
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
				return false
			}
			if comp != nil {
				var c int
				c, err = comp(temp[i], temp[j])
				if err != nil {
					return false
				}
				return c < 0
			}
			// Default: general compare
			var c int
			c, err = compareValues(temp[i], temp[j])
			if err != nil {
				return false
			}
			return c < 0
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
		var frags []string
		if len(vs) == 2 {
			switch delim := vs[1].(type) {
			case vm.String:
				frags = strings.Split(string(s), string(delim))
			case *vm.Regex:
				frags = delim.Split(string(s), -1)
			default:
				return vm.NIL, fmt.Errorf("split expected String or Regex")
			}
		} else {
			frags = strings.Split(string(s), "")
		}
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
		const minInt32 = -2147483648
		const maxInt32 = 2147483647
		coerce := func(f float64) (vm.Value, error) {
			if f < minInt32 || f > maxInt32 {
				return vm.NIL, fmt.Errorf("%s can't be coerced to int", vs[0])
			}
			return vm.MakeInt(int(math.Trunc(f))), nil
		}
		switch v := vs[0].(type) {
		case vm.Int:
			if v < minInt32 || v > maxInt32 {
				return vm.NIL, fmt.Errorf("%s can't be coerced to int", vs[0])
			}
			return v, nil
		case vm.Float:
			return coerce(float64(v))
		case vm.Char:
			return vm.Int(int(v)), nil
		case *vm.BigInt:
			if !v.Val().IsInt64() {
				return vm.NIL, fmt.Errorf("%s can't be coerced to int", vs[0])
			}
			i := v.Val().Int64()
			if i < minInt32 || i > maxInt32 {
				return vm.NIL, fmt.Errorf("%s can't be coerced to int", vs[0])
			}
			return vm.MakeInt(int(i)), nil
		case *vm.BigDecimal:
			f, _ := v.Val().Float64()
			return coerce(f)
		case *vm.Ratio:
			f, _ := v.Val().Float64()
			return coerce(f)
		case vm.Boolean:
			if bool(v) {
				return vm.MakeInt(1), nil
			}
			return vm.MakeInt(0), nil
		default:
			return vm.NIL, fmt.Errorf("%s can't be coerced to int", vs[0])
		}
	})

	longf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		const minInt64 = float64(-9223372036854775808)
		const maxInt64 = float64(9223372036854775807)
		coerce := func(f float64) (vm.Value, error) {
			if f < minInt64 || f > maxInt64 {
				return vm.NIL, fmt.Errorf("%s can't be coerced to long", vs[0])
			}
			return vm.Int(int64(math.Trunc(f))), nil
		}
		switch v := vs[0].(type) {
		case vm.Int:
			return v, nil
		case vm.Float:
			return coerce(float64(v))
		case vm.Char:
			return vm.Int(int(v)), nil
		case *vm.BigInt:
			if !v.Val().IsInt64() {
				return vm.NIL, fmt.Errorf("%s can't be coerced to long", vs[0])
			}
			return vm.Int(v.Val().Int64()), nil
		case *vm.BigDecimal:
			f, _ := v.Val().Float64()
			return coerce(f)
		case *vm.Ratio:
			f, _ := v.Val().Float64()
			return coerce(f)
		case vm.Boolean:
			if bool(v) {
				return vm.MakeInt(1), nil
			}
			return vm.MakeInt(0), nil
		default:
			return vm.NIL, fmt.Errorf("%s can't be coerced to long", vs[0])
		}
	})

	floatf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		f, ok := vm.ToFloat(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("%s can't be coerced to float", vs[0])
		}
		if math.IsInf(f, 0) {
			return vm.NIL, fmt.Errorf("%s can't be coerced to float", vs[0])
		}
		f32 := float32(f)
		if math.IsInf(float64(f32), 0) {
			return vm.NIL, fmt.Errorf("%s can't be coerced to float", vs[0])
		}
		return vm.Float32(float64(f32)), nil
	})

	doublef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		f, ok := vm.ToFloat(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("%s can't be coerced to double", vs[0])
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
		switch vs[0].(type) {
		case vm.Float, vm.Float32:
			return vm.TRUE, nil
		}
		ok := false
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
		switch v := vs[0].(type) {
		case vm.Int:
			if int(v) < 0 || int(v) > 0x10FFFF {
				return vm.NIL, fmt.Errorf("value out of range for char: %d", v)
			}
			return vm.Char(rune(v)), nil
		case vm.Char:
			return v, nil
		case vm.String:
			runes := []rune(string(v))
			if len(runes) == 1 {
				return vm.Char(runes[0]), nil
			}
			return vm.NIL, fmt.Errorf("%s can't be coerced to char", vs[0])
		case *vm.BigInt:
			n := v.Unbox().(*big.Int).Int64()
			if n < 0 || n > 0x10FFFF {
				return vm.NIL, fmt.Errorf("value out of range for char: %d", n)
			}
			return vm.Char(rune(n)), nil
		default:
			return vm.NIL, fmt.Errorf("%s can't be coerced to char", vs[0])
		}
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
		if vs[0] == vm.NIL {
			return vm.NIL, nil
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
		case *vm.List:
			if v == vm.EmptyList {
				return vm.NIL, nil
			}
			return v.First(), nil
		default:
			return vm.NIL, fmt.Errorf("peek not supported on %s", vs[0].Type())
		}
	})

	pop, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
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
		if _, ok := vs[0].(vm.Boolean); ok {
			return vm.NIL, fmt.Errorf("repeat expected an Int")
		}
		ni, ok := vm.ToInt(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("repeat expected an Int")
		}
		n := vm.Int(ni)
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

	// subs: substring (character-indexed, not byte-indexed)
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
		runes := []rune(string(s))
		si := int(start)
		if si < 0 || si > len(runes) {
			return vm.NIL, fmt.Errorf("string index out of range")
		}
		if len(vs) == 3 {
			end, ok := vs[2].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("subs expected Int end")
			}
			ei := int(end)
			if ei < si || ei > len(runes) {
				return vm.NIL, fmt.Errorf("string index out of range")
			}
			return vm.String(string(runes[si:ei])), nil
		}
		return vm.String(string(runes[si:])), nil
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

	// random-uuid: generate a random UUID v4
	randomUUID, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		var buf [16]byte
		_, err := crand.Read(buf[:])
		if err != nil {
			return vm.NIL, fmt.Errorf("random-uuid: %w", err)
		}
		buf[6] = (buf[6] & 0x0f) | 0x40 // version 4
		buf[8] = (buf[8] & 0x3f) | 0x80 // variant 10
		s := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
		return vm.NewUUID(s), nil
	})

	// rand-nth: returns a random element from a collection
	randNth, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
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
		if vs[0] == vm.NIL {
			return vm.NIL, fmt.Errorf("shuffle not supported on nil")
		}
		switch vs[0].(type) {
		case vm.String:
			return vm.NIL, fmt.Errorf("shuffle not supported on string")
		case *vm.PersistentMap, vm.Map, *vm.SortedMap:
			return vm.NIL, fmt.Errorf("shuffle not supported on map")
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
		case *vm.PersistentSet:
			return vm.NewTransientSet(v), nil
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
			return v.Persistent()
		case *vm.TransientVector:
			return v.Persistent()
		case *vm.TransientSet:
			return v.Persistent()
		default:
			return vm.NIL, fmt.Errorf("persistent! not supported on %s", vs[0].Type().Name())
		}
	})

	// conj!: mutating conj on a transient
	conjBang, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return vm.NewTransientVector(nil), nil
		}
		if len(vs) == 1 {
			return vs[0], nil
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
		case *vm.TransientSet:
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
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch t := vs[0].(type) {
		case *vm.TransientMap:
			var err error
			for i := 1; i < len(vs); i += 2 {
				val := vm.Value(vm.NIL)
				if i+1 < len(vs) {
					val = vs[i+1]
				}
				t, err = t.Assoc(vs[i], val)
				if err != nil {
					return vm.NIL, err
				}
			}
			return t, nil
		case *vm.TransientVector:
			var err error
			for i := 1; i < len(vs); i += 2 {
				val := vm.Value(vm.NIL)
				if i+1 < len(vs) {
					val = vs[i+1]
				}
				t, err = t.Assoc(vs[i], val)
				if err != nil {
					return vm.NIL, err
				}
			}
			return t, nil
		default:
			return vm.NIL, fmt.Errorf("assoc! not supported on %s", vs[0].Type().Name())
		}
	})

	// disj!: mutating disj on a transient set
	disjBang, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		t, ok := vs[0].(*vm.TransientSet)
		if !ok {
			return vm.NIL, fmt.Errorf("disj! expected TransientSet")
		}
		var err error
		for i := 1; i < len(vs); i++ {
			t, err = t.Disj(vs[i])
			if err != nil {
				return vm.NIL, err
			}
		}
		return t, nil
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

	// make-deftype: create a DType (deftype class) with a name and field symbols.
	makeDType, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		name, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("make-deftype expected String name")
		}
		fields := make([]vm.Symbol, len(vs)-1)
		for i := 1; i < len(vs); i++ {
			sym, ok := vs[i].(vm.Symbol)
			if !ok {
				return vm.NIL, fmt.Errorf("make-deftype expected Symbol fields")
			}
			fields[i-1] = sym
		}
		return vm.NewDType(string(name), fields), nil
	})

	// make-deftype-instance: construct an instance of a DType from positional field values.
	makeDTypeInstance, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		dt, ok := vs[0].(*vm.DType)
		if !ok {
			return vm.NIL, fmt.Errorf("make-deftype-instance expected DType, got %s", vs[0].Type().Name())
		}
		// Copy to avoid aliasing the VM's args slice (which may be reused).
		fields := make([]vm.Value, len(vs)-1)
		copy(fields, vs[1:])
		return vm.NewDTypeInstance(dt, fields), nil
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

	// prn-str: print readably + newline to string
	prnStr, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		b := &strings.Builder{}
		for i := range vs {
			if i > 0 {
				b.WriteRune(' ')
			}
			b.WriteString(vs[i].String())
		}
		b.WriteRune('\n')
		return vm.String(b.String()), nil
	})

	// print-str: print human-readably to string (no quotes on strings)
	printStr, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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
		return vm.String(b.String()), nil
	})

	// println-str: print human-readably + newline to string
	printlnStr, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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
		b.WriteRune('\n')
		return vm.String(b.String()), nil
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
			case vm.ArrayVector:
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
						} else if vec, ok := val.(vm.ArrayVector); ok {
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

	resolvef, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		sym, ok := vs[0].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("resolve expected Symbol")
		}
		cns := CurrentNS.Deref().(*vm.Namespace)
		if v := cns.Lookup(sym); v != vm.NIL {
			return v, nil
		}
		return vm.NIL, nil
	})
	if err != nil {
		panic(err)
	}

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

	boundFnStar, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("bound-fn* expects 1 arg")
		}
		fn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("bound-fn* expected Fn")
		}
		snap := vm.SnapshotBindings()
		if len(snap) == 0 {
			return fn, nil
		}
		wrapped, _ := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
			return vm.RunWithBindings(snap, func() (vm.Value, error) {
				return fn.Invoke(args)
			})
		})
		return wrapped, nil
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
			src     vm.Seq
			buf     []vm.Value // output items waiting to be yielded
			xf      vm.Fn      // the xform'd reducing fn
			done    bool       // completion called
			stopped bool       // early termination via reduced
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

	// realized? — test if a Delay, Promise, Future, or LazySeq has been realized
	isRealized, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if d, ok := vs[0].(*vm.Delay); ok {
			return vm.Boolean(d.IsRealized()), nil
		}
		if p, ok := vs[0].(*vm.Promise); ok {
			return vm.Boolean(p.IsRealized()), nil
		}
		if s, ok := vs[0].(*vm.LazySeq); ok {
			return vm.Boolean(s.IsRealized()), nil
		}
		return vm.NIL, fmt.Errorf("realized? expected delay, promise, future, or lazy seq")
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
	comparef, err = vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("compare expects 2 args")
		}
		c, err := vm.DefaultCompare(vs[0], vs[1])
		if err != nil {
			return vm.NIL, err
		}
		return vm.MakeInt(c), nil
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

	bitAndNot, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-and-not expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-and-not expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-and-not expected Int")
		}
		return vm.MakeInt(int(a) &^ int(b)), nil
	})

	bitFlip, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("bit-flip expects 2 args")
		}
		a, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-flip expected Int")
		}
		b, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("bit-flip expected Int")
		}
		return vm.MakeInt(int(a) ^ (1 << uint(b))), nil
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
		snap := vm.SnapshotBindings()
		p := vm.NewPromise()
		go func() {
			v, err := vm.RunWithBindings(snap, func() (vm.Value, error) {
				return fn.Invoke(nil)
			})
			if err != nil {
				p.Deliver(vm.NIL)
			} else {
				p.Deliver(v)
			}
		}()
		return p, nil
	})

	// add-tap / remove-tap / tap> — debug tap queue (synchronous)
	addTap, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("add-tap expects 1 arg")
		}
		fn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("add-tap expected Fn")
		}
		tapsMu.Lock()
		taps = append(taps, fn)
		tapsMu.Unlock()
		return vm.NIL, nil
	})

	removeTap, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("remove-tap expects 1 arg")
		}
		tapsMu.Lock()
		for i, t := range taps {
			if t == vs[0] {
				taps = append(taps[:i], taps[i+1:]...)
				break
			}
		}
		tapsMu.Unlock()
		return vm.NIL, nil
	})

	tapBang, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("tap> expects 1 arg")
		}
		tapsMu.Lock()
		snap := make([]vm.Fn, len(taps))
		copy(snap, taps)
		tapsMu.Unlock()
		for _, t := range snap {
			_, _ = t.Invoke([]vm.Value{vs[0]})
		}
		return vm.TRUE, nil
	})

	// add-watch — (add-watch atom-or-var key fn)
	addWatch, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("add-watch expects 3 args")
		}
		fn, ok := vs[2].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("add-watch expected Fn")
		}
		switch ref := vs[0].(type) {
		case *vm.Atom:
			ref.AddWatch(vs[1], fn)
		case *vm.Var:
			ref.AddWatch(vs[1], fn)
		default:
			return vm.NIL, fmt.Errorf("add-watch expected Atom or Var")
		}
		return vs[0], nil
	})

	// remove-watch — (remove-watch atom-or-var key)
	removeWatch, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("remove-watch expects 2 args")
		}
		switch ref := vs[0].(type) {
		case *vm.Atom:
			ref.RemoveWatch(vs[1])
		case *vm.Var:
			ref.RemoveWatch(vs[1])
		default:
			return vm.NIL, fmt.Errorf("remove-watch expected Atom or Var")
		}
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

	getValidator, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("get-validator expects 1 arg")
		}
		a, ok := vs[0].(*vm.Atom)
		if !ok {
			return vm.NIL, fmt.Errorf("get-validator expected Atom")
		}
		return a.Validator(), nil
	})

	// subvec — (subvec v start) or (subvec v start end)
	subvecf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("subvec expects 2-3 args")
		}
		s, ok := vm.ToInt(vs[1])
		if !ok {
			return vm.NIL, fmt.Errorf("subvec expected Int start")
		}

		switch v := vs[0].(type) {
		case vm.ArrayVector:
			end := len(v)
			if len(vs) == 3 {
				e, ok := vm.ToInt(vs[2])
				if !ok {
					return vm.NIL, fmt.Errorf("subvec expected Int end")
				}
				end = e
			}
			if s < 0 || end > len(v) || s > end {
				return vm.NIL, fmt.Errorf("subvec: index out of bounds")
			}
			result := make([]vm.Value, end-s)
			copy(result, v[s:end])
			return vm.NewArrayVector(result), nil
		case vm.PersistentVector:
			end := int(v.Count().(vm.Int))
			if len(vs) == 3 {
				e, ok := vm.ToInt(vs[2])
				if !ok {
					return vm.NIL, fmt.Errorf("subvec expected Int end")
				}
				end = e
			}
			if s < 0 || end > int(v.Count().(vm.Int)) || s > end {
				return vm.NIL, fmt.Errorf("subvec: index out of bounds")
			}
			result := make([]vm.Value, end-s)
			for i := s; i < end; i++ {
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
		// fn? returns true only for actual functions, not all IFn implementors
		// (keywords, maps, sets, vectors are invokable but are not fn?)
		switch vs[0].(type) {
		case *vm.NativeFn, *vm.Func, *vm.Closure, *vm.MultiArityFn, *vm.MultiFn:
			return vm.TRUE, nil
		default:
			return vm.FALSE, nil
		}
	})

	// double? — true only for float64 values; float? accepts float32 and float64.
	isDouble, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(vm.Float)
		return vm.Boolean(ok), nil
	})

	// instance? — type check (simplified: checks if type name matches)
	instancep, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		// We accept type objects (e.g. IntType) and check if the value's type matches
		if t, ok := vs[0].(vm.ValueType); ok {
			return vm.Boolean(vs[1].Type() == t), nil
		}
		return vm.FALSE, nil
	})

	// ifn? — true if value implements Fn (invokable: functions, keywords, maps, sets, vectors)
	isIFn, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		switch vs[0].(type) {
		case vm.Fn, vm.Keyword, vm.Symbol, *vm.PersistentMap, *vm.PersistentSet,
			vm.ArrayVector, vm.PersistentVector, *vm.SortedMap, *vm.SortedSet, *vm.Promise:
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})

	// identical? — reference/value identity
	identical, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.Boolean(vs[0] == vs[1]), nil
	})

	// any? — returns true for everything (every value satisfies any?)
	anyp, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.TRUE, nil
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
	ns.Def("*compiling-aot*", vm.FALSE)
	ns.Def("*in-wasm*", vm.FALSE)
	ns.Def("Object", vm.AnyType)
	// True if the host terminal renders ANSI escape sequences. Defaults to
	// true; flipped to false on platforms that don't (e.g. plan9 / rio —
	// see term_plan9.go).
	ns.Def("*ansi?*", vm.TRUE)

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
	ns.Def("+'", plusP)
	ns.Def("*'", mulP)
	ns.Def("-'", subP)
	ns.Def("unchecked-add", uncheckedAdd)
	ns.Def("unchecked-subtract", uncheckedSubtract)
	ns.Def("unchecked-multiply", uncheckedMultiply)
	ns.Def("unchecked-negate", uncheckedNegate)
	ns.Def("unchecked-divide-int", uncheckedDivideInt)
	ns.Def("unchecked-long", uncheckedLong)
	ns.Def("unchecked-int", uncheckedInt)
	ns.Def("unchecked-short", uncheckedShort)
	ns.Def("unchecked-byte", uncheckedByte)
	ns.Def("unchecked-char", uncheckedChar)
	ns.Def("unchecked-double", uncheckedDouble)
	ns.Def("unchecked-float", uncheckedFloat)

	ns.Def("=", equals)
	ns.Def("not=", notEq)
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
	ns.Def("complement", complement)

	ns.Def("set-macro!", setMacro)
	ns.Def("gensym", gensym)
	ns.Def("in-ns", inNs)
	ns.Def("exclude-in-current-ns", excludeInCurrentNs)
	ns.Def("use", use)
	ns.Def("alias", aliasf)
	ns.Def("name", name)
	ns.Def("namespace", namespace)

	ns.Def("vector", vector)
	ns.Def("vec", vec)
	ns.Def("hash-map", hashMap)
	ns.Def("array-map", arrayMap)
	ns.Def("list", list)
	ns.Def("range", rangef)
	ns.Def("keyword", keyword)
	ns.Def("symbol", symbolf)
	ns.Def("hash-set", hashSet)
	ns.Def("sorted-map", sortedMap)
	ns.Def("sorted-set", sortedSet)
	ns.Def("sorted-map-by", sortedMapBy)
	ns.Def("sorted-set-by", sortedSetBy)

	ns.Def("seq", seq)
	ns.Def("seq?", isSeq)
	ns.Def("list?", isList)

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
	ns.Def("key", keyf)
	ns.Def("val", valf)
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
	ns.Def("parse-long", parseInt)
	ns.Def("max", max)
	ns.Def("min", min)

	ns.Def("compare", comparef)
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
	ns.Def("long", longf)
	ns.Def("byte", intf)
	ns.Def("short", intf)
	ns.Def("float", floatf)
	ns.Def("double", doublef)
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
	ns.Def("random-uuid", randomUUID)
	ns.Def("rand-int", randInt)
	ns.Def("rand-nth", randNth)
	ns.Def("shuffle", shuffle)
	ns.Def("transient", transientf)
	ns.Def("persistent!", persistentf)
	ns.Def("conj!", conjBang)
	ns.Def("assoc!", assocBang)
	ns.Def("disj!", disjBang)
	ns.Def("dissoc!", dissocBang)
	ns.Def("make-record-type", makeRecordType)
	ns.Def("make-record", makeRecord)
	ns.Def("record?", isRecord)
	ns.Def("make-deftype", makeDType)
	ns.Def("make-deftype-instance", makeDTypeInstance)
	ns.Def("defprotocol*", defProtocol)
	installHierarchyBuiltins(ns)
	ns.Def("make-protocol-fn", makeProtocolFn)
	ns.Def("extend-type*", extendType)
	ns.Def("satisfies?", satisfies)
	ns.Def("defmulti*", defMulti)
	ns.Def("defmethod*", defMethod)
	ns.Def("methods", methods)
	ns.Def("pr-str", prStr)
	ns.Def("prn", prn)
	ns.Def("prn-str", prnStr)
	ns.Def("print-str", printStr)
	ns.Def("println-str", printlnStr)
	ns.Def("re-find", reFind)
	ns.Def("re-matches", reMatches)
	ns.Def("re-seq", reSeq)
	ns.Def("require", requiref)
	ns.Def("find-ns", findNs)
	ns.Def("resolve", resolvef)
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
	ns.Def("bound-fn*", boundFnStar)

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
			f := float64(v)
			if math.IsNaN(f) || math.IsInf(f, 0) {
				return vm.NIL, fmt.Errorf("cannot coerce %s to bigint", vs[0])
			}
			decimal := strconv.FormatFloat(f, 'g', -1, 64)
			i, _, err := new(big.Float).SetPrec(4096).SetMode(big.ToZero).Parse(decimal, 10)
			if err != nil || i == nil {
				return vm.NIL, fmt.Errorf("cannot coerce %s to bigint", vs[0])
			}
			bi, _ := i.Int(nil)
			return vm.NewBigInt(bi), nil
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

	// bigint?/big-int? — test if value is BigInt. Clojure core does not expose
	// a BigInt-specific predicate; big-int? is useful for compatibility suites.
	isBigInt, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		return vm.Boolean(vm.IsBigInt(vs[0])), nil
	})

	ns.Def("bigint", bigintf)
	ns.Def("bigint?", isBigInt)
	ns.Def("big-int?", isBigInt)

	// ratio? — test if value is Ratio
	isRatio, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		return vm.Boolean(vm.IsRatio(vs[0])), nil
	})
	ns.Def("ratio?", isRatio)

	// decimal? — test if value is BigDecimal
	isDecimal, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		return vm.Boolean(vm.IsBigDecimal(vs[0])), nil
	})
	ns.Def("decimal?", isDecimal)

	// sorted? — test if value is a sorted collection
	isSorted, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		switch vs[0].(type) {
		case *vm.SortedMap, *vm.SortedSet:
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("sorted?", isSorted)

	// map? — test if value is a map (hash or sorted)
	isMap, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		switch vs[0].(type) {
		case *vm.PersistentMap, *vm.SortedMap:
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("map?", isMap)

	// set? — test if value is a set (hash or sorted)
	isSet, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		switch vs[0].(type) {
		case *vm.PersistentSet, *vm.SortedSet:
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("set?", isSet)

	// map-entry? — test if value is a key/value pair from a map
	isMapEntry, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		if _, ok := vs[0].(vm.MapEntry); ok {
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("map-entry?", isMapEntry)

	// lazy-seq? — test if value is a thunk-backed lazy seq. LazySeq.Type()
	// reports ListType for print/seq compatibility, so user-side type
	// checks can't reach the underlying Go type; this predicate does.
	isLazySeq, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		if _, ok := vs[0].(*vm.LazySeq); ok {
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("lazy-seq?", isLazySeq)

	// reversible? — test if value supports rseq
	isReversible, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		switch vs[0].(type) {
		case vm.ArrayVector, vm.PersistentVector, *vm.SortedMap, *vm.SortedSet:
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("reversible?", isReversible)

	// rseq — reverse seq for sorted collections and vectors
	rseqf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		switch v := vs[0].(type) {
		case *vm.SortedMap:
			s := v.RSeq()
			if s == vm.EmptyList {
				return vm.NIL, nil
			}
			return s, nil
		case *vm.SortedSet:
			s := v.RSeq()
			if s == vm.EmptyList {
				return vm.NIL, nil
			}
			return s, nil
		case vm.ArrayVector:
			n := len(v)
			if n == 0 {
				return vm.NIL, nil
			}
			var s vm.Seq = vm.EmptyList
			for i := 0; i < n; i++ {
				s = vm.NewCons(v[i], s)
			}
			return s, nil
		case vm.PersistentVector:
			n := v.RawCount()
			if n == 0 {
				return vm.NIL, nil
			}
			var s vm.Seq = vm.EmptyList
			for i := 0; i < n; i++ {
				s = vm.NewCons(v.ValueAt(vm.MakeInt(i)), s)
			}
			return s, nil
		}
		return vm.NIL, fmt.Errorf("rseq not supported on: %s", vs[0].Type())
	})
	ns.Def("rseq", rseqf)

	// numerator — return numerator of a Ratio
	numeratorf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		if r, ok := vs[0].(*vm.Ratio); ok {
			num := r.Val().Num()
			return vm.MaybeDowngrade(new(big.Int).Set(num)), nil
		}
		return vm.NIL, fmt.Errorf("numerator expects a Ratio")
	})
	ns.Def("numerator", numeratorf)

	// denominator — return denominator of a Ratio
	denominatorf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		if r, ok := vs[0].(*vm.Ratio); ok {
			den := r.Val().Denom()
			return vm.MaybeDowngrade(new(big.Int).Set(den)), nil
		}
		return vm.NIL, fmt.Errorf("denominator expects a Ratio")
	})
	ns.Def("denominator", denominatorf)

	// bigdec — coerce to BigDecimal
	bigdecf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		switch v := vs[0].(type) {
		case *vm.BigDecimal:
			return v, nil
		case vm.Int:
			return vm.NewBigDecimalFromInt64(int64(v)), nil
		case vm.Float:
			return vm.NewBigDecimalFromFloat64(float64(v)), nil
		case *vm.BigInt:
			f, _ := new(big.Float).SetPrec(vm.BigDecimalPrecConst).SetInt(v.Val()).Float64()
			return vm.NewBigDecimalFromFloat64(f), nil
		case *vm.Ratio:
			f, _ := v.Val().Float64()
			return vm.NewBigDecimalFromFloat64(f), nil
		case vm.String:
			bd, ok := vm.NewBigDecimalFromString(string(v))
			if !ok {
				return vm.NIL, fmt.Errorf("cannot parse bigdec: %s", v)
			}
			return bd, nil
		}
		return vm.NIL, fmt.Errorf("cannot coerce %s to bigdec", vs[0].Type().Name())
	})
	ns.Def("bigdec", bigdecf)

	roundBigdec, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("round-bigdec expects precision, rounding mode, and value")
		}
		precision, ok := vm.ToInt(vs[0])
		if !ok {
			return vm.NIL, fmt.Errorf("round-bigdec expected integer precision")
		}
		var modeName string
		switch m := vs[1].(type) {
		case vm.Keyword:
			modeName = string(m)
		case vm.Symbol:
			modeName = string(m)
		case vm.String:
			modeName = string(m)
		default:
			return vm.NIL, fmt.Errorf("round-bigdec expected rounding mode")
		}
		modeName = strings.TrimPrefix(strings.ToLower(modeName), ":")
		bd, ok := vs[2].(*vm.BigDecimal)
		if !ok {
			return vs[2], nil
		}
		return roundBigDecimalValue(bd, precision, modeName)
	})
	ns.Def("round-bigdec", roundBigdec)

	// rationalize — convert to Ratio
	rationalizef, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments")
		}
		switch v := vs[0].(type) {
		case *vm.Ratio:
			return v, nil
		case vm.Int:
			return vm.MaybeSimplifyRatio(big.NewRat(int64(v), 1)), nil
		case *vm.BigInt:
			return vm.MaybeSimplifyRatio(new(big.Rat).SetInt(v.Val())), nil
		case vm.Float:
			f := float64(v)
			if math.IsInf(f, 0) || math.IsNaN(f) {
				return vm.NIL, fmt.Errorf("cannot rationalize %s", vs[0].Type().Name())
			}
			if r, ok := new(big.Rat).SetString(strconv.FormatFloat(f, 'f', -1, 64)); ok {
				return vm.MaybeSimplifyRatio(r), nil
			}
			return vm.NIL, fmt.Errorf("cannot rationalize %s", vs[0].Type().Name())
		case *vm.BigDecimal:
			if r, ok := new(big.Rat).SetString(v.Val().Text('f', -1)); ok {
				return vm.MaybeSimplifyRatio(r), nil
			}
			return vm.NIL, fmt.Errorf("cannot rationalize %s", vs[0].Type().Name())
		}
		return vm.NIL, fmt.Errorf("cannot rationalize %s", vs[0].Type().Name())
	})
	ns.Def("rationalize", rationalizef)

	ns.Def("transformer-seq*", transformerSeq)
	ns.Def("compare", comparef)
	ns.Def("fn?", isFn)
	ns.Def("ifn?", isIFn)
	ns.Def("double?", isDouble)
	ns.Def("instance?", instancep)
	ns.Def("identical?", identical)
	ns.Def("any?", anyp)
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
	ns.Def("bit-and-not", bitAndNot)
	ns.Def("bit-flip", bitFlip)
	ns.Def("re-groups", reGroups)
	ns.Def("promise", promisef)
	ns.Def("deliver", deliver)
	ns.Def("future*", futureStar)
	ns.Def("add-tap", addTap)
	ns.Def("remove-tap", removeTap)
	ns.Def("tap>", tapBang)
	ns.Def("add-watch", addWatch)
	ns.Def("remove-watch", removeWatch)
	ns.Def("alter-meta!", alterMeta)
	ns.Def("get-validator", getValidator)
	ns.Def("subvec", subvecf)
	ns.Def("print", printf)
	ns.Def("pr", prf)
	ns.Def("reduced", reducedf)
	ns.Def("reduced?", isReducedf)
	ns.Def("unreduced", unreduced)
	ns.Def("ensure-reduced", ensureReduced)

	// quot — integer division (truncated toward zero)
	quotf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.NumQuot(vs[0], vs[1])
	})
	ns.Def("quot", quotf)

	// rem — remainder of truncated division (sign follows dividend)
	remf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.NumRem(vs[0], vs[1])
	})
	ns.Def("rem", remf)

	// hash — returns the hash code of a value
	hashf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.MakeInt(int(vm.HashValue(vs[0]))), nil
	})
	ns.Def("hash", hashf)

	// parse-double — parse string to float
	parseDouble, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("parse-double expected String")
		}
		f, err := strconv.ParseFloat(string(s), 64)
		if err != nil {
			return vm.NIL, nil // Clojure returns nil for unparseable
		}
		return vm.Float(f), nil
	})
	ns.Def("parse-double", parseDouble)

	// parse-boolean — parse "true"/"false" to boolean
	parseBool, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("parse-boolean expected String")
		}
		switch string(s) {
		case "true":
			return vm.TRUE, nil
		case "false":
			return vm.FALSE, nil
		}
		return vm.NIL, nil
	})
	ns.Def("parse-boolean", parseBool)

	// NaN? — test if value is NaN
	isNaN, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		switch v := vs[0].(type) {
		case vm.Float:
			return vm.Boolean(math.IsNaN(float64(v))), nil
		case vm.Float32:
			return vm.Boolean(math.IsNaN(float64(v))), nil
		case vm.Int, *vm.BigInt, *vm.Ratio, *vm.BigDecimal:
			return vm.FALSE, nil
		default:
			return vm.NIL, fmt.Errorf("NaN? requires a number, got %s", vs[0].Type().Name())
		}
	})
	ns.Def("NaN?", isNaN)

	// infinite? — test if value is +/-Inf
	isInfinite, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		if f, ok := vs[0].(vm.Float); ok {
			return vm.Boolean(math.IsInf(float64(f), 0)), nil
		}
		return vm.FALSE, nil
	})
	ns.Def("infinite?", isInfinite)

	// boolean? — test if value is boolean
	isBool, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(vm.Boolean)
		return vm.Boolean(ok), nil
	})
	ns.Def("boolean?", isBool)

	// char? — test if value is char
	isChar, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(vm.Char)
		return vm.Boolean(ok), nil
	})
	ns.Def("char?", isChar)

	// var? — test if value is a var
	isVar, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(*vm.Var)
		return vm.Boolean(ok), nil
	})
	ns.Def("var?", isVar)

	// string/index-of — find first index of substring/char
	indexOf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("index-of expected String as first arg")
		}
		var needle string
		switch v := vs[1].(type) {
		case vm.String:
			needle = string(v)
		case vm.Char:
			needle = string(rune(v))
		default:
			return vm.NIL, fmt.Errorf("index-of expected String or Char as second arg")
		}
		str := string(s)
		if len(vs) == 3 {
			from, ok := vs[2].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("index-of expected Int as third arg")
			}
			idx := strings.Index(str[int(from):], needle)
			if idx == -1 {
				return vm.NIL, nil
			}
			return vm.MakeInt(idx + int(from)), nil
		}
		idx := strings.Index(str, needle)
		if idx == -1 {
			return vm.NIL, nil
		}
		return vm.MakeInt(idx), nil
	})
	ns.Def("index-of", indexOf)

	// string/last-index-of — find last index of substring/char
	lastIndexOf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("last-index-of expected String as first arg")
		}
		var needle string
		switch v := vs[1].(type) {
		case vm.String:
			needle = string(v)
		case vm.Char:
			needle = string(rune(v))
		default:
			return vm.NIL, fmt.Errorf("last-index-of expected String or Char as second arg")
		}
		str := string(s)
		if len(vs) == 3 {
			from, ok := vs[2].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("last-index-of expected Int as third arg")
			}
			idx := strings.LastIndex(str[:int(from)+1], needle)
			if idx == -1 {
				return vm.NIL, nil
			}
			return vm.MakeInt(idx), nil
		}
		idx := strings.LastIndex(str, needle)
		if idx == -1 {
			return vm.NIL, nil
		}
		return vm.MakeInt(idx), nil
	})
	ns.Def("last-index-of", lastIndexOf)

	// --- Array operations ---

	// Helper: build a typed array from size or seq
	buildArray := func(kind vm.ArrayKind, vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		// (x-array n) or (x-array n init)
		if n, ok := vs[0].(vm.Int); ok {
			size := int(n)
			if size < 0 {
				return vm.NIL, fmt.Errorf("negative array size: %d", size)
			}
			var arr *vm.TypedArray
			switch kind {
			case vm.ArrayByte:
				arr = vm.NewByteArray(size)
			case vm.ArrayInt:
				arr = vm.NewIntArray(size)
			case vm.ArrayFloat:
				arr = vm.NewFloatArray(size)
			case vm.ArrayObject:
				arr = vm.NewObjectArray(size)
			}
			if len(vs) == 2 {
				for i := 0; i < size; i++ {
					if err := arr.Set(i, vs[1]); err != nil {
						return vm.NIL, err
					}
				}
			}
			return arr, nil
		}
		// (x-array coll)
		s, serr := seqOf(vs[0])
		if serr != nil {
			return vm.NIL, serr
		}
		if s == nil {
			return vm.NIL, fmt.Errorf("cannot create array from %s", vs[0].Type().Name())
		}
		var vals []vm.Value
		for ; s != nil; s = s.Next() {
			vals = append(vals, s.First())
		}
		var arr *vm.TypedArray
		switch kind {
		case vm.ArrayByte:
			data := make([]byte, len(vals))
			for i, v := range vals {
				n, ok := v.(vm.Int)
				if !ok {
					return vm.NIL, fmt.Errorf("byte-array element must be Int, got %s", v.Type().Name())
				}
				data[i] = byte(n)
			}
			arr = vm.NewByteArrayFrom(data)
		case vm.ArrayInt:
			data := make([]int64, len(vals))
			for i, v := range vals {
				switch n := v.(type) {
				case vm.Int:
					data[i] = int64(n)
				default:
					return vm.NIL, fmt.Errorf("int-array element must be Int, got %s", v.Type().Name())
				}
			}
			arr = vm.NewIntArrayFrom(data)
		case vm.ArrayFloat:
			data := make([]float64, len(vals))
			for i, v := range vals {
				f, ok := vm.ToFloat(v)
				if !ok {
					return vm.NIL, fmt.Errorf("double-array element must be numeric, got %s", v.Type().Name())
				}
				data[i] = f
			}
			arr = vm.NewFloatArrayFrom(data)
		case vm.ArrayObject:
			arr = vm.NewObjectArrayFrom(vals)
		}
		return arr, nil
	}

	byteArrayf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return buildArray(vm.ArrayByte, vs)
	})
	ns.Def("byte-array", byteArrayf)

	intArrayf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return buildArray(vm.ArrayInt, vs)
	})
	ns.Def("int-array", intArrayf)
	ns.Def("long-array", intArrayf)

	doubleArrayf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return buildArray(vm.ArrayFloat, vs)
	})
	ns.Def("double-array", doubleArrayf)
	ns.Def("float-array", doubleArrayf)

	objectArrayf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return buildArray(vm.ArrayObject, vs)
	})
	ns.Def("object-array", objectArrayf)

	// make-array — (make-array type size)
	makeArrayf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("make-array expects 2 args (type, size)")
		}
		size, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("make-array size must be Int")
		}
		n := int(size)
		if n < 0 {
			return vm.NIL, fmt.Errorf("negative array size: %d", n)
		}
		switch t := vs[0].(type) {
		case vm.Keyword:
			switch string(t) {
			case "byte":
				return vm.NewByteArray(n), nil
			case "int", "long":
				return vm.NewIntArray(n), nil
			case "double", "float":
				return vm.NewFloatArray(n), nil
			case "object":
				return vm.NewObjectArray(n), nil
			}
			return vm.NIL, fmt.Errorf("unknown array type: %s", t)
		default:
			// Default to object-array
			return vm.NewObjectArray(n), nil
		}
	})
	ns.Def("make-array", makeArrayf)

	// aget — (aget arr idx) or (aget arr idx idx2 ...)
	agetf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		arr, ok := vs[0].(*vm.TypedArray)
		if !ok {
			return vm.NIL, fmt.Errorf("aget expects array, got %s", vs[0].Type().Name())
		}
		idx, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("aget index must be Int")
		}
		i := int(idx)
		if i < 0 || i >= arr.Len() {
			return vm.NIL, fmt.Errorf("array index %d out of bounds for length %d", i, arr.Len())
		}
		val := arr.Get(i)
		// Support nested access: (aget arr i j k ...)
		for _, extra := range vs[2:] {
			inner, ok := val.(*vm.TypedArray)
			if !ok {
				return vm.NIL, fmt.Errorf("aget nested: expected array, got %s", val.Type().Name())
			}
			idx, ok := extra.(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("aget index must be Int")
			}
			j := int(idx)
			if j < 0 || j >= inner.Len() {
				return vm.NIL, fmt.Errorf("array index %d out of bounds for length %d", j, inner.Len())
			}
			val = inner.Get(j)
		}
		return val, nil
	})
	ns.Def("aget", agetf)

	// aset — (aset arr idx val)
	asetf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		arr, ok := vs[0].(*vm.TypedArray)
		if !ok {
			return vm.NIL, fmt.Errorf("aset expects array, got %s", vs[0].Type().Name())
		}
		idx, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("aset index must be Int")
		}
		i := int(idx)
		if i < 0 || i >= arr.Len() {
			return vm.NIL, fmt.Errorf("array index %d out of bounds for length %d", i, arr.Len())
		}
		if err := arr.Set(i, vs[2]); err != nil {
			return vm.NIL, err
		}
		return vs[2], nil
	})
	ns.Def("aset", asetf)

	// alength — (alength arr)
	alengthf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		arr, ok := vs[0].(*vm.TypedArray)
		if !ok {
			return vm.NIL, fmt.Errorf("alength expects array, got %s", vs[0].Type().Name())
		}
		return vm.MakeInt(arr.Len()), nil
	})
	ns.Def("alength", alengthf)

	// aclone — (aclone arr)
	aclonef, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		arr, ok := vs[0].(*vm.TypedArray)
		if !ok {
			return vm.NIL, fmt.Errorf("aclone expects array, got %s", vs[0].Type().Name())
		}
		return arr.Clone(), nil
	})
	ns.Def("aclone", aclonef)

	// to-array — (to-array coll) creates object-array from seq
	toArrayf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		s, serr := seqOf(vs[0])
		if serr != nil {
			return vm.NIL, serr
		}
		if s == nil {
			return vm.NewObjectArray(0), nil
		}
		var vals []vm.Value
		for ; s != nil; s = s.Next() {
			vals = append(vals, s.First())
		}
		return vm.NewObjectArrayFrom(vals), nil
	})
	ns.Def("to-array", toArrayf)

	// into-array — (into-array seq) or (into-array type seq)
	intoArrayf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 1 {
			return toArrayf.(*vm.NativeFn).Invoke(vs)
		}
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		kw, ok := vs[0].(vm.Keyword)
		if !ok {
			// If not a keyword, just make object-array from second arg
			return buildArray(vm.ArrayObject, vs[1:])
		}
		switch string(kw) {
		case "byte":
			return buildArray(vm.ArrayByte, vs[1:])
		case "int", "long":
			return buildArray(vm.ArrayInt, vs[1:])
		case "double", "float":
			return buildArray(vm.ArrayFloat, vs[1:])
		case "object":
			return buildArray(vm.ArrayObject, vs[1:])
		}
		return vm.NIL, fmt.Errorf("unknown array type: %s", kw)
	})
	ns.Def("into-array", intoArrayf)

	// bytes — coerce string to byte-array
	bytesf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch v := vs[0].(type) {
		case vm.String:
			data := []byte(string(v))
			return vm.NewByteArrayFrom(data), nil
		case *vm.TypedArray:
			if v.Kind() == vm.ArrayByte {
				return v, nil
			}
		}
		return vm.NIL, fmt.Errorf("bytes expects String or byte-array, got %s", vs[0].Type().Name())
	})
	ns.Def("bytes", bytesf)

	// ints — coerce seq to int-array
	intsf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if a, ok := vs[0].(*vm.TypedArray); ok && a.Kind() == vm.ArrayInt {
			return a, nil
		}
		return buildArray(vm.ArrayInt, vs)
	})
	ns.Def("ints", intsf)
	ns.Def("longs", intsf)

	// doubles — coerce seq to double-array
	doublesf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		if a, ok := vs[0].(*vm.TypedArray); ok && a.Kind() == vm.ArrayFloat {
			return a, nil
		}
		return buildArray(vm.ArrayFloat, vs)
	})
	ns.Def("doubles", doublesf)

	// array? — type predicate
	isArrayf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(*vm.TypedArray)
		return vm.Boolean(ok), nil
	})
	ns.Def("array?", isArrayf)

	// alter-var-root — alter a var's root binding via (f root & args).
	// Reads/writes the root, bypassing any current dynamic binding.
	// TODO: the read/apply/write is not atomic. let-go evaluates synchronously
	// today, so contention does not arise. If concurrent evaluation is added,
	// centralize the read/apply/write under Var-level synchronization.
	alterVarRoot, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("alter-var-root expects at least 2 args")
		}
		v, ok := vs[0].(*vm.Var)
		if !ok {
			return vm.NIL, fmt.Errorf("alter-var-root expects a Var")
		}
		fn, ok := vs[1].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("alter-var-root expects a function")
		}
		result, err := v.AlterRootArgs(fn, vs[2:])
		if err != nil {
			return vm.NIL, err
		}
		return result, nil
	})
	ns.Def("alter-var-root", alterVarRoot)

	// var-get — get the value of a Var
	varGet, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("var-get expects 1 arg")
		}
		v, ok := vs[0].(*vm.Var)
		if !ok {
			return vm.NIL, fmt.Errorf("var-get expects a Var")
		}
		return v.Deref(), nil
	})
	ns.Def("var-get", varGet)

	// macroexpand — expand a macro form once
	macroexpandf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("macroexpand expects 1 arg")
		}
		form := vs[0]
		lst, ok := form.(*vm.List)
		if !ok {
			return form, nil
		}
		if lst == vm.EmptyList || lst.First() == nil {
			return form, nil
		}
		sym, ok := lst.First().(vm.Symbol)
		if !ok {
			return form, nil
		}
		// Resolve the symbol to a var
		symStr := string(sym)
		var resolved *vm.Var
		if strings.Contains(symStr, "/") {
			parts := strings.SplitN(symStr, "/", 2)
			rns := nsRegistry[resolveNSAlias(parts[0])]
			if rns != nil {
				if v := rns.Lookup(vm.Symbol(parts[1])); v != vm.NIL {
					resolved, _ = v.(*vm.Var)
				}
			}
		} else {
			if CoreNS != nil {
				if v := CoreNS.Lookup(vm.Symbol(symStr)); v != vm.NIL {
					resolved, _ = v.(*vm.Var)
				}
			}
		}
		if resolved == nil || !resolved.IsMacro() {
			return form, nil
		}
		// Call the macro with the form's args
		args := make([]vm.Value, 0)
		for s := lst.Next(); s != nil; s = s.Next() {
			args = append(args, s.First())
		}
		macroFn, ok := resolved.Deref().(vm.Fn)
		if !ok {
			return form, nil
		}
		return macroFn.Invoke(args)
	})
	ns.Def("macroexpand", macroexpandf)

	// macroexpand-1 — same as macroexpand for now
	ns.Def("macroexpand-1", macroexpandf)

	// sleep — sleep for n milliseconds
	sleepf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("sleep expects 1 arg")
		}
		switch v := vs[0].(type) {
		case vm.Int:
			time.Sleep(time.Duration(int64(v)) * time.Millisecond)
		case vm.Float:
			time.Sleep(time.Duration(int64(v)) * time.Millisecond)
		default:
			return vm.NIL, fmt.Errorf("sleep expects a number")
		}
		return vm.NIL, nil
	})
	ns.Def("sleep", sleepf)

	// intern — intern a var in a namespace
	internf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("intern expects 2 or 3 args")
		}
		var targetNS *vm.Namespace
		switch n := vs[0].(type) {
		case vm.Symbol:
			targetNS = nsRegistry[resolveNSAlias(string(n))]
		case *vm.Namespace:
			targetNS = n
		default:
			return vm.NIL, fmt.Errorf("intern expects a namespace or symbol")
		}
		if targetNS == nil {
			return vm.NIL, fmt.Errorf("namespace not found")
		}
		sym, ok := vs[1].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("intern expects a symbol name")
		}
		if len(vs) == 2 {
			if existing := targetNS.LookupLocal(sym); existing != nil {
				return existing, nil
			}
			return targetNS.Def(string(sym), vm.NIL), nil
		}
		v := targetNS.Def(string(sym), vs[2])
		return v, nil
	})
	ns.Def("intern", internf)

	// create-ns — create or return existing namespace
	createNsf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("create-ns expects 1 arg")
		}
		sym, ok := vs[0].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("create-ns expects a symbol")
		}
		return NS(string(sym)), nil
	})
	ns.Def("create-ns", createNsf)

	// pop! — pop from a transient vector
	popBang, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("pop! expects 1 arg")
		}
		tv, ok := vs[0].(*vm.TransientVector)
		if !ok {
			return vm.NIL, fmt.Errorf("pop! expects a transient vector")
		}
		return tv.Pop()
	})
	ns.Def("pop!", popBang)

	// with-out-str — capture stdout as string (macro helper)
	withOutStrf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("with-out-str* expects 1 arg (a thunk)")
		}
		fn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("with-out-str* expects a function")
		}
		// Capture stdout
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			return vm.NIL, err
		}
		os.Stdout = w
		_, callErr := fn.Invoke(nil)
		w.Close()
		os.Stdout = old
		buf := make([]byte, 64*1024)
		n, _ := r.Read(buf)
		r.Close()
		if callErr != nil {
			return vm.NIL, callErr
		}
		return vm.String(buf[:n]), nil
	})
	ns.Def("with-out-str*", withOutStrf)

	// uuid? — type predicate for UUID
	isUUID, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(*vm.UUID)
		return vm.Boolean(ok), nil
	})
	ns.Def("uuid?", isUUID)

	// inst? — type predicate for Instant
	isInst, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(*vm.Instant)
		return vm.Boolean(ok), nil
	})
	ns.Def("inst?", isInst)

	// parse-uuid — parse a string into a UUID
	parseUUID, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("parse-uuid expects 1 arg")
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("parse-uuid expects a string, got %s", vs[0].Type().Name())
		}
		u := vm.ParseUUID(string(s))
		if u == nil {
			return vm.NIL, nil
		}
		return u, nil
	})
	ns.Def("parse-uuid", parseUUID)

	// == — numeric equality across numeric categories.
	numericEq, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.TRUE, nil
		}
		for i := 1; i < len(vs); i++ {
			if !vm.NumEquivalent(vs[0], vs[i]) {
				return vm.FALSE, nil
			}
		}
		return vm.TRUE, nil
	})
	ns.Def("==", numericEq)

	// IO builtins (open, close!, read-line, write!, etc.)
	installIOBuiltins(ns)

	CoreNS = ns
	vm.SetCoreNamespace(ns)

	RegisterNS(ns)
	installClojureCompatAliases(ns)
}

func installClojureCompatAliases(ns *vm.Namespace) {
	ns.Def("java.lang.Byte", vm.IntType)
	ns.Def("java.lang.Short", vm.IntType)
	ns.Def("java.lang.Integer", vm.IntType)
	ns.Def("java.lang.Long", vm.IntType)
	ns.Def("java.lang.Float", vm.FloatType)
	ns.Def("java.lang.Double", vm.FloatType)
	ns.Def("java.lang.Boolean", vm.BooleanType)
	ns.Def("java.lang.Object", vm.AnyType)
	ns.Def("java.math.BigDecimal", vm.BigDecimalType)
	ns.Def("java.util.UUID", vm.UUIDType)

	ns.Def("clojure.lang.BigInt", vm.BigIntType)
	ns.Def("clojure.lang.Ratio", vm.RatioType)
	ns.Def("clojure.lang.BigDecimal", vm.BigDecimalType)
	ns.Def("clojure.lang.Atom", vm.AtomType)
	ns.Def("clojure.lang.PersistentHashSet", vm.SetType)
	ns.Def("clojure.lang.IPending", vm.PromiseType)

	for _, name := range []string{
		"clojure.lang.Associative",
		"clojure.lang.Counted",
		"clojure.lang.Indexed",
		"clojure.lang.Seqable",
		"clojure.lang.IPersistentCollection",
		"clojure.lang.IReduce",
	} {
		ns.Def(name, vm.Symbol(name))
	}
	ns.Def("IReduce", vm.Symbol("clojure.lang.IReduce"))
	ns.Def("String", vm.StringType)
	ns.Def("Boolean", vm.BooleanType)
	ns.Def("Integer.", ns.Lookup("int").(*vm.Var).Deref())
	ns.Def("->Integer", ns.Lookup("int").(*vm.Var).Deref())

	longNS := DefNSBare("Long")
	longNS.Def("MAX_VALUE", longCompatValue(9223372036854775807))
	longNS.Def("MIN_VALUE", longCompatValue(-9223372036854775808))
	longNS.Def("TYPE", vm.IntType)

	doubleNS := DefNSBare("Double")
	doubleNS.Def("MAX_VALUE", vm.Float(1.7976931348623157e308))
	doubleNS.Def("MIN_VALUE", vm.Float(4.9e-324))
	doubleNS.Def("TYPE", vm.FloatType)

	integerNS := DefNSBare("Integer")
	integerNS.Def("TYPE", vm.IntType)

	booleanNS := DefNSBare("Boolean")
	booleanNS.Def("TYPE", vm.BooleanType)

	mapEntryNS := DefNSBare("clojure.lang.MapEntry")
	create, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("MapEntry/create expects 2 args")
		}
		return vm.MapEntry{Key: vs[0], Value: vs[1]}, nil
	})
	if err != nil {
		panic(err)
	}
	mapEntryNS.Def("create", create)
}

func longCompatValue(v int64) vm.Value {
	if strconv.IntSize == 64 {
		return vm.Int(v)
	}
	return vm.NewBigIntFromInt64(v)
}

func strValue(v vm.Value) string {
	if v == vm.NIL {
		return ""
	}
	switch x := v.(type) {
	case vm.String:
		return string(x)
	case vm.Char:
		return string(rune(x))
	case vm.Float:
		f := float64(x)
		if math.IsInf(f, 1) {
			return "Infinity"
		}
		if math.IsInf(f, -1) {
			return "-Infinity"
		}
	case *vm.BigInt:
		return x.Val().String()
	case *vm.BigDecimal:
		s := x.String()
		return strings.TrimSuffix(s, "M")
	}
	return v.String()
}
