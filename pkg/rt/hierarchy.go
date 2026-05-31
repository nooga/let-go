/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"sync"

	"github.com/nooga/let-go/pkg/vm"
)

var (
	hierarchyMu           sync.Mutex
	typeParentsMu         sync.RWMutex
	globalHierarchy       = emptyHierarchy()
	registeredTypeParents = map[vm.ValueType]*vm.PersistentSet{}
)

var (
	hKeyParents     = vm.Keyword("parents")
	hKeyAncestors   = vm.Keyword("ancestors")
	hKeyDescendants = vm.Keyword("descendants")
)

var (
	cljAssociative           = vm.Symbol("clojure.lang.Associative")
	cljCounted               = vm.Symbol("clojure.lang.Counted")
	cljIndexed               = vm.Symbol("clojure.lang.Indexed")
	cljSeqable               = vm.Symbol("clojure.lang.Seqable")
	cljIPersistentCollection = vm.Symbol("clojure.lang.IPersistentCollection")
	cljIReduce               = vm.Symbol("clojure.lang.IReduce")
	cljIEditableCollection   = vm.Symbol("clojure.lang.IEditableCollection")
	// Throwable is registered bare (not under clojure.lang) by
	// installClojureCompatAliases; ExInfo reports it as an ancestor so medley's
	// ex-message/ex-cause (guarded by (instance? Throwable ex)) work on ex-info.
	cljThrowable = vm.Symbol("Throwable")
)

func emptyHierarchy() *vm.PersistentMap {
	return vm.EmptyPersistentMap.
		Assoc(hKeyParents, vm.EmptyPersistentMap).(*vm.PersistentMap).
		Assoc(hKeyAncestors, vm.EmptyPersistentMap).(*vm.PersistentMap).
		Assoc(hKeyDescendants, vm.EmptyPersistentMap).(*vm.PersistentMap)
}

func installHierarchyBuiltins(ns *vm.Namespace) {
	ns.Def("make-hierarchy", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return emptyHierarchy(), nil
	}))

	ns.Def("register-type-parent!", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		child, ok := vs[0].(vm.ValueType)
		if !ok {
			return vm.NIL, fmt.Errorf("register-type-parent! expected child type")
		}
		if vs[1] == vm.NIL {
			return vm.NIL, nil
		}
		typeParentsMu.Lock()
		registeredTypeParents[child] = setConj(registeredTypeParents[child], vs[1])
		typeParentsMu.Unlock()
		return vm.NIL, nil
	}))

	ns.Def("derive", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		switch len(vs) {
		case 2:
			if err := validateGlobalDerive(vs[0], vs[1]); err != nil {
				return vm.NIL, err
			}
			hierarchyMu.Lock()
			defer hierarchyMu.Unlock()
			next, err := deriveIn(globalHierarchy, vs[0], vs[1])
			if err != nil {
				return vm.NIL, err
			}
			globalHierarchy = next
			return vm.NIL, nil
		case 3:
			h, err := asHierarchy(vs[0])
			if err != nil {
				return vm.NIL, err
			}
			if vs[1] == vm.NIL || vs[2] == vm.NIL {
				return vm.NIL, fmt.Errorf("derive requires non-nil tag and parent")
			}
			return deriveIn(h, vs[1], vs[2])
		default:
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
	}))

	ns.Def("underive", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		switch len(vs) {
		case 2:
			hierarchyMu.Lock()
			defer hierarchyMu.Unlock()
			if next, err := underiveIn(globalHierarchy, vs[0], vs[1]); err == nil {
				globalHierarchy = next
			}
			return vm.NIL, nil
		case 3:
			h, err := asHierarchy(vs[0])
			if err != nil {
				return vm.NIL, err
			}
			return underiveIn(h, vs[1], vs[2])
		default:
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
	}))

	ns.Def("parents", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		switch len(vs) {
		case 1:
			hierarchyMu.Lock()
			h := globalHierarchy
			hierarchyMu.Unlock()
			return parentsFor(h, vs[0], true, false), nil
		case 2:
			h, err := asHierarchy(vs[0])
			if err != nil {
				return vm.NIL, nil
			}
			return parentsFor(h, vs[1], true, true), nil
		default:
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
	}))

	ns.Def("ancestors", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		switch len(vs) {
		case 1:
			hierarchyMu.Lock()
			h := globalHierarchy
			hierarchyMu.Unlock()
			return ancestorsFor(h, vs[0], true, false), nil
		case 2:
			h, err := asHierarchy(vs[0])
			if err != nil {
				return vm.NIL, nil
			}
			return ancestorsFor(h, vs[1], true, true), nil
		default:
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
	}))

	ns.Def("descendants", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		switch len(vs) {
		case 1:
			if shouldThrowDescendants(vs[0]) {
				return vm.NIL, fmt.Errorf("descendants does not support root type %s", vs[0])
			}
			hierarchyMu.Lock()
			h := globalHierarchy
			hierarchyMu.Unlock()
			return descendantsFor(h, vs[0], false), nil
		case 2:
			if shouldThrowDescendants(vs[1]) {
				return vm.NIL, fmt.Errorf("descendants does not support root type %s", vs[1])
			}
			h, err := asHierarchy(vs[0])
			if err != nil {
				return vm.NIL, nil
			}
			return descendantsFor(h, vs[1], false), nil
		default:
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
	}))

	ns.Def("isa?", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		switch len(vs) {
		case 2:
			hierarchyMu.Lock()
			h := globalHierarchy
			hierarchyMu.Unlock()
			return vm.Boolean(isaIn(h, vs[0], vs[1])), nil
		case 3:
			h, err := asHierarchy(vs[0])
			if err != nil {
				return vm.FALSE, nil
			}
			return vm.Boolean(isaIn(h, vs[1], vs[2])), nil
		default:
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
	}))
}

func asHierarchy(v vm.Value) (*vm.PersistentMap, error) {
	h, ok := v.(*vm.PersistentMap)
	if !ok {
		return nil, fmt.Errorf("invalid hierarchy")
	}
	for _, key := range []vm.Keyword{hKeyParents, hKeyAncestors, hKeyDescendants} {
		if _, ok := h.ValueAt(key).(*vm.PersistentMap); !ok {
			return nil, fmt.Errorf("invalid hierarchy")
		}
	}
	return h, nil
}

func validateGlobalDerive(tag, parent vm.Value) error {
	if tag == vm.NIL || parent == vm.NIL {
		return fmt.Errorf("derive requires non-nil tag and parent")
	}
	if !isGlobalDeriveTag(tag) {
		return fmt.Errorf("derive tag must be a namespaced symbol/keyword or type")
	}
	if !isNamespacedSymbolic(parent) {
		return fmt.Errorf("derive parent must be a namespaced symbol/keyword")
	}
	return nil
}

func isGlobalDeriveTag(v vm.Value) bool {
	if _, ok := v.(vm.ValueType); ok {
		return true
	}
	return isNamespacedSymbolic(v)
}

func isNamespacedSymbolic(v vm.Value) bool {
	switch x := v.(type) {
	case vm.Keyword:
		return x.Namespace() != vm.NIL
	case vm.Symbol:
		return x.Namespace() != vm.NIL
	default:
		return false
	}
}

func deriveIn(h *vm.PersistentMap, tag, parent vm.Value) (*vm.PersistentMap, error) {
	if valueEqual(tag, parent) || isaIn(h, parent, tag) {
		return nil, fmt.Errorf("cyclic derivation")
	}

	parents := hierarchySection(h, hKeyParents)
	parents = assocSet(parents, tag, parent)
	return rebuildHierarchy(parents), nil
}

func underiveIn(h *vm.PersistentMap, tag, parent vm.Value) (*vm.PersistentMap, error) {
	parents := hierarchySection(h, hKeyParents)
	current := setAt(parents, tag)
	if current == nil || current.Contains(parent) == vm.FALSE {
		return h, nil
	}
	current = current.Disj(parent)
	if current.RawCount() == 0 {
		parents = parents.Dissoc(tag).(*vm.PersistentMap)
	} else {
		parents = parents.Assoc(tag, current).(*vm.PersistentMap)
	}
	return rebuildHierarchy(parents), nil
}

func rebuildHierarchy(parents *vm.PersistentMap) *vm.PersistentMap {
	ancestors := vm.EmptyPersistentMap
	descendants := vm.EmptyPersistentMap

	for _, tag := range mapKeys(parents) {
		ancs := collectAncestors(parents, tag, map[vm.Value]bool{})
		if ancs != nil && ancs.RawCount() > 0 {
			ancestors = ancestors.Assoc(tag, ancs).(*vm.PersistentMap)
			for _, anc := range setValues(ancs) {
				descendants = assocSet(descendants, anc, tag)
			}
		}
	}

	return vm.EmptyPersistentMap.
		Assoc(hKeyParents, parents).(*vm.PersistentMap).
		Assoc(hKeyAncestors, ancestors).(*vm.PersistentMap).
		Assoc(hKeyDescendants, descendants).(*vm.PersistentMap)
}

func collectAncestors(parents *vm.PersistentMap, tag vm.Value, seen map[vm.Value]bool) *vm.PersistentSet {
	result := vm.EmptyPersistentSet
	ps := setAt(parents, tag)
	for _, p := range setValues(ps) {
		if seen[p] {
			continue
		}
		seen[p] = true
		result = setConj(result, p)
		for _, gp := range setValues(collectAncestors(parents, p, seen)) {
			result = setConj(result, gp)
		}
	}
	return result
}

func parentsFor(h *vm.PersistentMap, tag vm.Value, includeType bool, emptySet bool) vm.Value {
	if !validLookupTag(tag) {
		return vm.NIL
	}
	result := setAt(hierarchySection(h, hKeyParents), tag)
	if includeType {
		result = setUnion(result, directTypeParents(tag))
	}
	return nilOrSet(result, emptySet)
}

func ancestorsFor(h *vm.PersistentMap, tag vm.Value, includeType bool, emptySet bool) vm.Value {
	if !validLookupTag(tag) {
		return vm.NIL
	}
	result := setAt(hierarchySection(h, hKeyAncestors), tag)
	if includeType {
		if vt, ok := tag.(vm.ValueType); ok {
			result = setUnion(result, directTypeAncestors(vt))
		}
		for _, p := range setValues(directTypeParents(tag)) {
			result = setConj(result, p)
			if pt, ok := p.(vm.ValueType); ok {
				result = setUnion(result, directTypeAncestors(pt))
			}
		}
	}
	return nilOrSet(result, emptySet)
}

func descendantsFor(h *vm.PersistentMap, tag vm.Value, emptySet bool) vm.Value {
	if !validLookupTag(tag) {
		return vm.NIL
	}
	return nilOrSet(setAt(hierarchySection(h, hKeyDescendants), tag), emptySet)
}

func isaIn(h *vm.PersistentMap, tag, parent vm.Value) bool {
	if valueEqual(tag, parent) {
		return true
	}
	if parents := parentsFor(h, tag, true, false); parents != vm.NIL {
		for _, p := range setValues(parents.(*vm.PersistentSet)) {
			if valueEqual(p, parent) || isaIn(h, p, parent) {
				return true
			}
		}
	}
	return false
}

func validLookupTag(v vm.Value) bool {
	switch v.(type) {
	case vm.Keyword, vm.Symbol, vm.ValueType, *vm.Protocol:
		return v != vm.NIL
	default:
		return false
	}
}

func shouldThrowDescendants(v vm.Value) bool {
	return v == vm.AnyType
}

func hierarchySection(h *vm.PersistentMap, key vm.Keyword) *vm.PersistentMap {
	m, _ := h.ValueAt(key).(*vm.PersistentMap)
	if m == nil {
		return vm.EmptyPersistentMap
	}
	return m
}

func assocSet(m *vm.PersistentMap, key, value vm.Value) *vm.PersistentMap {
	return m.Assoc(key, setConj(setAt(m, key), value)).(*vm.PersistentMap)
}

func setAt(m *vm.PersistentMap, key vm.Value) *vm.PersistentSet {
	if m == nil {
		return nil
	}
	s, _ := m.ValueAt(key).(*vm.PersistentSet)
	return s
}

func setConj(s *vm.PersistentSet, value vm.Value) *vm.PersistentSet {
	if s == nil {
		s = vm.EmptyPersistentSet
	}
	return s.Conj(value).(*vm.PersistentSet)
}

func setUnion(a, b *vm.PersistentSet) *vm.PersistentSet {
	result := a
	if result == nil {
		result = vm.EmptyPersistentSet
	}
	for _, v := range setValues(b) {
		result = setConj(result, v)
	}
	return result
}

func nilOrSet(s *vm.PersistentSet, emptySet bool) vm.Value {
	if s == nil || s.RawCount() == 0 {
		if emptySet {
			return vm.EmptyPersistentSet
		}
		return vm.NIL
	}
	return s
}

func setValues(s *vm.PersistentSet) []vm.Value {
	if s == nil || s.RawCount() == 0 {
		return nil
	}
	var out []vm.Value
	for seq := s.Seq(); seq != nil && seq != vm.EmptyList; seq = seq.Next() {
		out = append(out, seq.First())
	}
	return out
}

func mapKeys(m *vm.PersistentMap) []vm.Value {
	if m == nil || m.RawCount() == 0 {
		return nil
	}
	var out []vm.Value
	for seq := m.Seq(); seq != nil && seq != vm.EmptyList; seq = seq.Next() {
		k, _, ok := vm.MapEntryKV(seq.First())
		if ok {
			out = append(out, k)
		}
	}
	return out
}

func directTypeParents(tag vm.Value) *vm.PersistentSet {
	vt, ok := tag.(vm.ValueType)
	if !ok {
		return nil
	}
	typeParentsMu.RLock()
	result := registeredTypeParents[vt]
	typeParentsMu.RUnlock()
	switch vt {
	case vm.StringType:
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljSeqable)
		result = setConj(result, cljCounted)
		result = setConj(result, cljIndexed)
	case vm.ArrayVectorType, vm.PersistentVectorType:
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljAssociative)
		result = setConj(result, cljCounted)
		result = setConj(result, cljIndexed)
		result = setConj(result, cljSeqable)
		result = setConj(result, cljIPersistentCollection)
		result = setConj(result, cljIReduce)
		result = setConj(result, cljIEditableCollection)
	case vm.MapType:
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljAssociative)
		result = setConj(result, cljCounted)
		result = setConj(result, cljSeqable)
		result = setConj(result, cljIPersistentCollection)
		result = setConj(result, cljIReduce)
		result = setConj(result, cljIEditableCollection)
	case vm.SortedMapType:
		// Sorted maps have no transient form, so no IEditableCollection.
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljAssociative)
		result = setConj(result, cljCounted)
		result = setConj(result, cljSeqable)
		result = setConj(result, cljIPersistentCollection)
		result = setConj(result, cljIReduce)
	case vm.SetType:
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljCounted)
		result = setConj(result, cljSeqable)
		result = setConj(result, cljIPersistentCollection)
		result = setConj(result, cljIReduce)
		result = setConj(result, cljIEditableCollection)
	case vm.SortedSetType:
		// Sorted sets have no transient form, so no IEditableCollection.
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljCounted)
		result = setConj(result, cljSeqable)
		result = setConj(result, cljIPersistentCollection)
		result = setConj(result, cljIReduce)
	case vm.ListType, vm.SequenceType:
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljCounted)
		result = setConj(result, cljSeqable)
		result = setConj(result, cljIPersistentCollection)
		result = setConj(result, cljIReduce)
	case vm.RangeType, vm.RepeatType, vm.IterateType, vm.TypedArrayType:
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljSeqable)
		result = setConj(result, cljIReduce)
	case vm.AtomType:
		result = setConj(result, vm.AnyType)
	case vm.PromiseType:
		result = setConj(result, vm.AnyType)
	case vm.ExInfoType:
		// ex-info values are let-go's equivalent of Clojure's ExceptionInfo,
		// which IS-A Throwable. Reporting the Throwable marker lets medley's
		// (instance? Throwable ex) guard pass so ex-message/ex-cause work.
		result = setConj(result, vm.AnyType)
		result = setConj(result, cljThrowable)
	default:
		if _, ok := vt.(*vm.RecordType); ok {
			result = setConj(result, vm.AnyType)
			result = setConj(result, cljAssociative)
			result = setConj(result, cljCounted)
			result = setConj(result, cljSeqable)
			result = setConj(result, cljIPersistentCollection)
			result = setConj(result, cljIReduce)
		} else if _, ok := vt.(*vm.DType); ok {
			result = setConj(result, vm.AnyType)
		}
	}
	return result
}

func directTypeAncestors(vt vm.ValueType) *vm.PersistentSet {
	result := vm.EmptyPersistentSet
	for _, p := range setValues(directTypeParents(vt)) {
		result = setConj(result, p)
		if pvt, ok := p.(vm.ValueType); ok {
			result = setUnion(result, directTypeAncestors(pvt))
		}
	}
	if result.RawCount() == 0 {
		return nil
	}
	return result
}

func valueEqual(a, b vm.Value) bool {
	if vm.ValueEquals != nil {
		return vm.ValueEquals(a, b)
	}
	return a == b
}
