/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// theHostArrayListType is the ValueType for the mutable java.util.ArrayList
// compat shim — mutable, not a persistent let-go value, so it lives in the
// compat layer rather than in pkg/vm. medley's partition-between / sliding
// and code like (vec (doto (ArrayList. xs) (Collections/shuffle))) build one
// on their :clj branches.
type theHostArrayListType struct{}

func (t *theHostArrayListType) String() string     { return t.Name() }
func (t *theHostArrayListType) Type() vm.ValueType { return vm.TypeType }
func (t *theHostArrayListType) Unbox() any         { return nil }
func (t *theHostArrayListType) Name() string       { return "java.util.ArrayList" }
func (t *theHostArrayListType) Box(any) (vm.Value, error) {
	return vm.NIL, fmt.Errorf("java.util.ArrayList cannot be boxed")
}

// hostArrayListType is the singleton type for hostArrayList values.
var hostArrayListType = &theHostArrayListType{}

// hostArrayList is a mutable, index-addressed list. It is seqable and
// counted, so vec, seq and count read it as they read the JVM's.
type hostArrayList struct{ items []vm.Value }

func (l *hostArrayList) Type() vm.ValueType { return hostArrayListType }
func (l *hostArrayList) Unbox() any         { return l }

// String renders like java.util.AbstractCollection.toString: [a, b, c].
func (l *hostArrayList) String() string {
	parts := make([]string, len(l.items))
	for i, v := range l.items {
		parts[i] = strValue(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func (l *hostArrayList) RawCount() int   { return len(l.items) }
func (l *hostArrayList) Count() vm.Value { return vm.MakeInt(len(l.items)) }

// Seq walks a snapshot, so mutating the list during a walk does not disturb
// it (the JVM would throw ConcurrentModificationException instead).
func (l *hostArrayList) Seq() vm.Seq {
	if len(l.items) == 0 {
		return nil
	}
	return vm.NewList(append([]vm.Value(nil), l.items...)).(vm.Seq)
}

func (l *hostArrayList) index(v vm.Value, size int) (int, error) {
	i, ok := vm.ToInt64(v)
	if !ok {
		return 0, fmt.Errorf("java.util.ArrayList index must be an integer, got %s", v.Type().Name())
	}
	if i < 0 || int(i) >= size {
		return 0, vm.NewThrownError(vm.NewExInfoWithClass(vm.ClassIndexOutOfBounds,
			fmt.Sprintf("Index %d out of bounds for length %d", i, len(l.items)), nil))
	}
	return int(i), nil
}

func (l *hostArrayList) InvokeMethod(name vm.Symbol, args []vm.Value) (vm.Value, error) {
	switch string(name) {
	case "add":
		switch len(args) {
		case 1:
			l.items = append(l.items, args[0])
			return vm.TRUE, nil
		case 2:
			// Inserting at size appends, so the bound is size+1.
			i, err := l.index(args[0], len(l.items)+1)
			if err != nil {
				return vm.NIL, err
			}
			l.items = append(l.items[:i], append([]vm.Value{args[1]}, l.items[i:]...)...)
			return vm.NIL, nil
		}
	case "addAll":
		if len(args) == 1 {
			vals, err := seqValues(args[0])
			if err != nil {
				return vm.NIL, err
			}
			l.items = append(l.items, vals...)
			return vm.Boolean(len(vals) > 0), nil
		}
	case "get":
		if len(args) == 1 {
			i, err := l.index(args[0], len(l.items))
			if err != nil {
				return vm.NIL, err
			}
			return l.items[i], nil
		}
	case "set":
		if len(args) == 2 {
			i, err := l.index(args[0], len(l.items))
			if err != nil {
				return vm.NIL, err
			}
			old := l.items[i]
			l.items[i] = args[1]
			return old, nil
		}
	case "size":
		if len(args) == 0 {
			return vm.MakeInt(len(l.items)), nil
		}
	case "isEmpty":
		if len(args) == 0 {
			return vm.Boolean(len(l.items) == 0), nil
		}
	case "clear":
		if len(args) == 0 {
			l.items = nil
			return vm.NIL, nil
		}
	case "contains":
		if len(args) == 1 {
			return vm.Boolean(l.indexOf(args[0]) >= 0), nil
		}
	case "indexOf":
		if len(args) == 1 {
			return vm.MakeInt(l.indexOf(args[0])), nil
		}
	case "toArray":
		if len(args) == 0 {
			return vm.NewObjectArrayFrom(append([]vm.Value(nil), l.items...)), nil
		}
	case "toString":
		if len(args) == 0 {
			return vm.String(l.String()), nil
		}
	}
	return vm.NIL, fmt.Errorf("java.util.ArrayList has no method .%s/%d", name, len(args))
}

func (l *hostArrayList) indexOf(v vm.Value) int {
	for i, x := range l.items {
		if valueEquals(x, v) {
			return i
		}
	}
	return -1
}

// seqValues collects the elements of any seqable value.
func seqValues(v vm.Value) ([]vm.Value, error) {
	if v == vm.NIL {
		return nil, nil
	}
	s, err := seqOf(v)
	if err != nil {
		return nil, err
	}
	return appendSeqValues(nil, forceSeq(s)), nil
}

// installHostArrayList registers the java.util.ArrayList constructor forms
// and the java.util.Collections statics that act on one. As on the JVM, the
// constructor's single argument is a capacity (an int, ignored) or a
// collection to copy.
func installHostArrayList(ns *vm.Namespace) {
	ctor := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		l := &hostArrayList{}
		switch len(vs) {
		case 0:
		case 1:
			if _, ok := vs[0].(vm.Int); ok {
				break // capacity hint
			}
			vals, err := seqValues(vs[0])
			if err != nil {
				return vm.NIL, fmt.Errorf("java.util.ArrayList. expected an int capacity or a collection: %w", err)
			}
			l.items = vals
		default:
			return vm.NIL, fmt.Errorf("java.util.ArrayList. expects 0 or 1 args, got %d", len(vs))
		}
		return l, nil
	})
	ns.Def("java.util.ArrayList", hostArrayListType)
	for _, n := range []string{"ArrayList.", "->ArrayList", "java.util.ArrayList.", "->java.util.ArrayList"} {
		ns.Def(n, ctor)
	}

	listArg := func(name string, vs []vm.Value) (*hostArrayList, error) {
		if len(vs) != 1 {
			return nil, fmt.Errorf("Collections/%s expects 1 arg", name)
		}
		l, ok := vs[0].(*hostArrayList)
		if !ok {
			// The JVM throws UnsupportedOperationException for an
			// immutable list; every let-go collection is one.
			return nil, vm.NewThrownError(vm.NewExInfoWithClass(vm.ClassUnsupportedOperation,
				fmt.Sprintf("Collections/%s needs a mutable java.util.ArrayList, got %s", name, vs[0].Type().Name()), nil))
		}
		return l, nil
	}
	shuffle := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		l, err := listArg("shuffle", vs)
		if err != nil {
			return vm.NIL, err
		}
		rngShuffle(len(l.items), func(i, j int) {
			l.items[i], l.items[j] = l.items[j], l.items[i]
		})
		return vm.NIL, nil
	})
	reverse := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		l, err := listArg("reverse", vs)
		if err != nil {
			return vm.NIL, err
		}
		for i, j := 0, len(l.items)-1; i < j; i, j = i+1, j-1 {
			l.items[i], l.items[j] = l.items[j], l.items[i]
		}
		return vm.NIL, nil
	})
	for _, nm := range []string{"Collections", "java.util.Collections"} {
		cns := defStaticNS(nm)
		cns.Def("shuffle", shuffle)
		cns.Def("reverse", reverse)
	}
}
