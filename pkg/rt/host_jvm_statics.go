/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"strconv"

	"github.com/nooga/let-go/pkg/vm"
)

// defStaticNS returns (creating if needed) a bare namespace that holds host
// static members (e.g. `Util/hash`, `Long/parseLong`). Unlike DefNSBare it does
// NOT auto-refer clojure.core, so defining a member whose name overlaps a core
// var (`Util/hash` vs `clojure.core/hash`) does not print a shadow WARNING at
// every `lg` startup. These namespaces only carry static members; they never
// resolve core names, so the missing refer is harmless.
func defStaticNS(name string) *vm.Namespace {
	name = resolveNSAlias(name)
	nsMu.RLock()
	if e := nsRegistry[name]; e != nil {
		nsMu.RUnlock()
		return e
	}
	nsMu.RUnlock()

	ns := vm.NewNamespace(name)
	nsMu.Lock()
	nsRegistry[name] = ns
	nsMu.Unlock()
	return ns
}

// installJVMStatics registers JVM static methods and (instance? Class x) markers
// that Clojure libraries reach on their :clj branches. Motivated by metosin/malli
// (registry, regex cache, entry parser, string coercion) but each is a general,
// real host-interop resolution backed by a let-go primitive.
func installJVMStatics(ns *vm.Namespace) {
	// --- (instance? Class x) host-class markers ---
	// Map literals evaluate to vm.MapType (let-go.lang.Map), so match that.
	RegisterHostClass("java.util.Map", vm.MapType)
	RegisterHostClass("CharSequence", vm.StringType)
	RegisterHostClass("Pattern", vm.RegexType) // bare java.util.regex.Pattern
	RegisterHostClass("java.util.AbstractList", vm.Symbol("java.util.AbstractList"))
	RegisterHostClass("java.util.Vector", vm.Symbol("java.util.Vector"))

	// --- clojure.lang.LazilyPersistentVector/createOwning(objectArray) -> vector.
	vecFn := ns.Lookup("vec").(*vm.Var).Deref()
	for _, nm := range []string{"LazilyPersistentVector", "clojure.lang.LazilyPersistentVector"} {
		defStaticNS(nm).Def("createOwning", vecFn)
	}

	// --- clojure.lang.PersistentArrayMap/createWithCheck(flatArray) -> map,
	// throwing on a duplicate key. The array is [k0 v0 k1 v1 ...].
	createWithCheck := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("PersistentArrayMap/createWithCheck expects 1 arg")
		}
		arr, ok := vs[0].(*vm.TypedArray)
		if !ok {
			return vm.NIL, fmt.Errorf("PersistentArrayMap/createWithCheck expects an array")
		}
		n := arr.RawCount()
		if n%2 != 0 {
			return vm.NIL, fmt.Errorf("PersistentArrayMap/createWithCheck: odd-length array (%d)", n)
		}
		var m vm.Associative = vm.EmptyPersistentMap
		for i := 0; i+1 < n; i += 2 {
			k := arr.Get(i)
			before := m.(*vm.PersistentMap).RawCount()
			m = m.Assoc(k, arr.Get(i+1))
			if m.(*vm.PersistentMap).RawCount() == before {
				return vm.NIL, fmt.Errorf("duplicate key: %s", k)
			}
		}
		return m, nil
	})
	for _, nm := range []string{"PersistentArrayMap", "clojure.lang.PersistentArrayMap"} {
		defStaticNS(nm).Def("createWithCheck", createWithCheck)
	}

	// --- java.lang.reflect.Array/newInstance(class, n) -> n-element object array
	// (the class argument is ignored — let-go arrays are untyped).
	defStaticNS("Array").Def("newInstance", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 2 {
			if n, ok := vs[1].(vm.Int); ok {
				if n < 0 {
					return vm.NIL, fmt.Errorf("Array/newInstance: negative size %d", n)
				}
				return vm.NewObjectArray(int(n)), nil
			}
		}
		return vm.NIL, fmt.Errorf("Array/newInstance expects (class, int)")
	}))

	// --- clojure.lang.Util/hash + Murmur3/hashLong -> let-go hash; Util/hashCombine
	// -> the standard Clojure combine.
	hashOf := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.Int(int64(vm.HashValue(vs[0]))), nil
	})
	hashCombine := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		a, _ := vs[0].(vm.Int)
		b, _ := vs[1].(vm.Int)
		return vm.Int(int64(a) ^ (int64(b) + 0x9e3779b9 + (int64(a) << 6) + (int64(a) >> 2))), nil
	})
	for _, nm := range []string{"Util", "clojure.lang.Util"} {
		defStaticNS(nm).Def("hash", hashOf)
		defStaticNS(nm).Def("hashCombine", hashCombine)
	}
	for _, nm := range []string{"Murmur3", "clojure.lang.Murmur3"} {
		defStaticNS(nm).Def("hashLong", hashOf)
	}

	// --- number/uuid string parses (Long/parseLong, Float/parseFloat, ...). ---
	parseLong := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("parseLong expects a string")
		}
		// strconv.Atoi (returns int) + MakeInt matches let-go's own parse-long
		// and avoids an int64->int conversion (CodeQL "incorrect conversion").
		n, err := strconv.Atoi(string(s))
		if err != nil {
			return vm.NIL, err
		}
		return vm.MakeInt(n), nil
	})
	parseFloat := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("parseFloat expects a string")
		}
		f, err := strconv.ParseFloat(string(s), 64)
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(f), nil
	})
	defStaticNS("Long").Def("parseLong", parseLong)
	defStaticNS("Integer").Def("parseInt", parseLong)
	defStaticNS("Float").Def("parseFloat", parseFloat)
	defStaticNS("Double").Def("parseDouble", parseFloat)

	// bare UUID/fromString (lang.go registers only the fully-qualified
	// java.util.UUID namespace).
	defStaticNS("UUID").Def("fromString", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("UUID/fromString expects a string")
		}
		u := vm.ParseUUID(string(s))
		if u == nil {
			return vm.NIL, fmt.Errorf("invalid UUID: %q", string(s))
		}
		return u, nil
	}))
}
