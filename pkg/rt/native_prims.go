/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/vm"
)

// Name returns the name of a value as a string.
// For Symbols and Keywords, returns the name part after the namespace separator.
// For Strings, returns the string itself.
// Mirrors clojure.core/name — `(name x)`.
//
//lg:native
//lg:name name
func Name(v vm.Value) (string, error) {
	if v == vm.NIL {
		return "", fmt.Errorf("name called on nil")
	}
	switch x := v.(type) {
	case vm.String:
		return string(x), nil
	case vm.Symbol:
		name := x.Name()
		if name == vm.NIL {
			return "", fmt.Errorf("name expected Symbol with a name")
		}
		return string(name.(vm.String)), nil
	case vm.Keyword:
		name := x.Name()
		if name == vm.NIL {
			return "", fmt.Errorf("name expected Keyword with a name")
		}
		return string(name.(vm.String)), nil
	default:
		return "", fmt.Errorf("name expected Symbol, Keyword or String")
	}
}

// UpperCase converts a value's string rendering to uppercase.
// Mirrors clojure.string/upper-case — `(upper-case s)`: JVM Clojure calls
// .toString on the argument, so keywords/symbols coerce ((upper-case :a) →
// ":A") and only nil throws. Registered in clojure.string namespace.
//
//lg:native
//lg:ns clojure.string
//lg:name upper-case
func UpperCase(v vm.Value) (string, error) {
	if v == vm.NIL {
		return "", fmt.Errorf("upper-case expected non-nil")
	}
	return strings.ToUpper(strValue(v)), nil
}

// Subs returns a substring from index start to the end.
// Mirrors clojure.core/subs with 2 arguments — `(subs s start)`.
//
//lg:native
//lg:name subs
func Subs(s string, start int) (string, error) {
	if start < 0 {
		return "", fmt.Errorf("string index out of range")
	}

	// subs is character-indexed, but a substring needs byte offsets. Walk
	// the string to map the start rune index to byte offset instead
	// of materializing []rune(s) for the whole string on every call — subs
	// is hot in parsers, and the full conversion dominated allocation.
	str := string(s)
	byteStart := -1
	ri, bo := 0, 0
	for {
		if ri == start {
			byteStart = bo
			break
		}
		if bo >= len(str) {
			break
		}
		_, size := utf8.DecodeRuneInString(str[bo:])
		bo += size
		ri++
	}
	if byteStart < 0 {
		return "", fmt.Errorf("string index out of range") // start past end
	}
	return str[byteStart:], nil
}

// Subs3 returns a substring from index start to end (exclusive).
// Mirrors clojure.core/subs with 3 arguments — `(subs s start end)`.
//
//lg:native
//lg:name subs
func Subs3(s string, start, end int) (string, error) {
	if start < 0 {
		return "", fmt.Errorf("string index out of range")
	}
	if end < start {
		return "", fmt.Errorf("string index out of range")
	}

	// subs is character-indexed, but a substring needs byte offsets. Walk
	// the string to map the start/end rune indices to byte offsets instead
	// of materializing []rune(s) for the whole string on every call — subs
	// is hot in parsers, and the full conversion dominated allocation.
	// Stop as soon as both offsets are known (i.e. at max(start, end) runes).
	str := string(s)
	byteStart, byteEnd := -1, -1
	ri, bo := 0, 0
	for {
		if ri == start {
			byteStart = bo
		}
		if ri == end {
			byteEnd = bo
		}
		if byteStart >= 0 && byteEnd >= 0 {
			break
		}
		if bo >= len(str) {
			break
		}
		_, size := utf8.DecodeRuneInString(str[bo:])
		bo += size
		ri++
	}
	if byteStart < 0 {
		return "", fmt.Errorf("string index out of range") // start past end
	}
	if byteEnd < 0 {
		return "", fmt.Errorf("string index out of range") // end past end
	}
	return str[byteStart:byteEnd], nil
}

// Nth returns the element at index in a collection.
// For Indexed (vectors, strings, arrays, etc.), returns the element at that index.
// For seqs, walks linearly to find the element at that index.
// Mirrors clojure.core/nth with 2 arguments — `(nth coll i)`.
//
//lg:native
//lg:name nth
func Nth(coll vm.Value, i int) (vm.Value, error) {
	if coll == vm.NIL {
		return vm.NIL, nil
	}
	// Fast path: positional collections (vm.Indexed) — those addressed by a
	// 0-based integer index. Maps and sets are key-addressable rather than
	// positional and fall through to the seq walk below, as do lazy seqs.
	// This path is the only correct route for positional-but-not-seqable
	// types (transient vectors, chunks): the seq walk would return nil.
	if ix, ok := coll.(vm.Indexed); ok {
		if i < 0 || i >= ix.RawCount() {
			return vm.NIL, fmt.Errorf("nth index out of bounds")
		}
		return ix.Nth(i), nil
	}
	// Seq path: linear walk
	if i < 0 {
		return vm.NIL, fmt.Errorf("nth index out of bounds")
	}
	s, err := seqOf(coll)
	if err != nil {
		return vm.NIL, err
	}
	if v, ok := nthInSeq(forceSeq(s), i); ok {
		return v, nil
	}
	return vm.NIL, fmt.Errorf("nth index out of bounds")
}

// Nth3 returns the element at index in a collection, or notFound if out of bounds.
// For Indexed (vectors, strings, arrays, etc.), returns the element at that index.
// For seqs, walks linearly to find the element at that index.
// Mirrors clojure.core/nth with 3 arguments — `(nth coll i not-found)`.
//
//lg:native
//lg:name nth
func Nth3(coll vm.Value, i int, notFound vm.Value) (vm.Value, error) {
	if coll == vm.NIL {
		return notFound, nil
	}
	// Fast path: positional collections (vm.Indexed) — those addressed by a
	// 0-based integer index. Maps and sets are key-addressable rather than
	// positional and fall through to the seq walk below, as do lazy seqs.
	// This path is the only correct route for positional-but-not-seqable
	// types (transient vectors, chunks): the seq walk would return nil.
	if ix, ok := coll.(vm.Indexed); ok {
		if i < 0 || i >= ix.RawCount() {
			return notFound, nil
		}
		return ix.Nth(i), nil
	}
	// Seq path: linear walk
	if i < 0 {
		return notFound, nil
	}
	s, err := seqOf(coll)
	if err != nil {
		return notFound, nil
	}
	if v, ok := nthInSeq(forceSeq(s), i); ok {
		return v, nil
	}
	return notFound, nil
}

// Deref dereferences a Reference (atom, var, delay, promise, …).
// Mirrors clojure.core/deref — `(deref ref)` / `@ref`.
//
//lg:native
//lg:name deref
func Deref(ref vm.Value) (vm.Value, error) {
	r, ok := vm.AsRef(ref)
	if !ok {
		return vm.NIL, fmt.Errorf("deref expected Reference")
	}
	return r.Deref(), nil
}

// Deref3 blocks up to ms milliseconds for a promise/future's value, else
// returns timeoutVal. Mirrors clojure.core/deref with 3 args —
// `(deref blocking-ref timeout-ms timeout-val)`.
//
//lg:native
//lg:name deref
func Deref3(bref vm.Value, ms int, timeoutVal vm.Value) (vm.Value, error) {
	b, ok := bref.(vm.BlockingDeref)
	if !ok {
		return vm.NIL, fmt.Errorf("deref with timeout expects a blocking ref (promise or future)")
	}
	return b.DerefTimeout(int64(ms), timeoutVal), nil
}

// Str concatenates the string representation of its arguments.
// Mirrors clojure.core/str — `(str a b c)`.
//
//lg:native
//lg:name str
func Str(args ...vm.Value) (string, error) {
	b := &strings.Builder{}
	for i := range args {
		b.WriteString(strValue(args[i]))
	}
	return b.String(), nil
}

// Get looks up key in an associative collection, returning nil when absent
// or when coll is not a lookup. Mirrors clojure.core/get — `(get coll key)`.
//
//lg:native
//lg:name get
func Get(coll vm.Value, key vm.Value) (vm.Value, error) {
	as, ok := coll.(vm.Lookup)
	if !ok {
		return vm.NIL, nil
	}
	return as.ValueAt(key), nil
}

// Get3 looks up key in coll, returning notFound when absent or when coll is
// not a lookup. Mirrors clojure.core/get — `(get coll key not-found)`.
//
//lg:native
//lg:name get
func Get3(coll vm.Value, key vm.Value, notFound vm.Value) (vm.Value, error) {
	as, ok := coll.(vm.Lookup)
	if !ok {
		return notFound, nil
	}
	return as.ValueAtOr(key, notFound), nil
}

// Conj adds elements to a collection (front for lists/seqs, back for vectors,
// key/val for maps). Mirrors clojure.core/conj — `(conj coll x y …)`.
//
//lg:native
//lg:name conj
func Conj(args ...vm.Value) (vm.Value, error) {
	if len(args) == 0 {
		return vm.ArrayVector{}, nil
	}
	if len(args) == 1 {
		return args[0], nil
	}
	var seq vm.Collection
	if args[0] == vm.NIL {
		seq = vm.EmptyList
	} else {
		if _, ok := args[0].(vm.String); ok {
			return vm.NIL, fmt.Errorf("conj expected Collection")
		}
		var ok bool
		seq, ok = args[0].(vm.Collection)
		if !ok {
			if s, ok := args[0].(vm.Seq); ok {
				for i := 1; i < len(args); i++ {
					if seq == nil {
						seq = vm.NewCons(args[i], s)
					} else {
						seq = seq.Conj(args[i])
					}
				}
				return seq, nil
			}
			return vm.NIL, fmt.Errorf("conj expected Collection")
		}
	}
	for i := 1; i < len(args); i++ {
		if isMapType(seq) && !canConjMapEntry(args[i]) {
			return vm.NIL, fmt.Errorf("conj expected map entry")
		}
		seq = seq.Conj(args[i])
	}
	return seq, nil
}

// Reduce implements (reduce f coll): the accumulator seeds from the first
// element of coll, and an empty/nil coll invokes f with no arguments.
// Mirrors clojure.core/reduce with 2 arguments.
//
//lg:native
//lg:name reduce
func Reduce(ec *vm.ExecContext, f vm.Value, coll vm.Value) (vm.Value, error) {
	mfn, ok := vm.AsFn(f)
	if !ok {
		return vm.NIL, fmt.Errorf("reduce expected Fn")
	}
	return reduceColl(ec, mfn, coll, false, vm.NIL)
}

// Reduce3 implements (reduce f init coll): the accumulator seeds from init,
// and a nil coll returns init unchanged without calling f. Mirrors
// clojure.core/reduce with 3 arguments.
//
//lg:native
//lg:name reduce
func Reduce3(ec *vm.ExecContext, f vm.Value, init vm.Value, coll vm.Value) (vm.Value, error) {
	// 3-arg form with nil coll: return init regardless of fn
	if coll == vm.NIL {
		return init, nil
	}
	mfn, ok := vm.AsFn(f)
	if !ok {
		return vm.NIL, fmt.Errorf("reduce expected Fn")
	}
	return reduceColl(ec, mfn, coll, true, init)
}

// reduceColl is the shared reduce loop for both arities. When hasInit is
// false the accumulator seeds from the collection's first element (and an
// empty input invokes mfn with no args); when true it seeds from initVal.
func reduceColl(ec *vm.ExecContext, mfn vm.Fn, coll vm.Value, hasInit bool, initVal vm.Value) (vm.Value, error) {
	// Handle nil and empty collections
	if coll == vm.NIL {
		if hasInit {
			return initVal, nil
		}
		return ec.Invoke(mfn, nil)
	}
	// Check for empty collection first (skip for lazy/cons — RawCount forces realization)
	switch coll.(type) {
	case *vm.LazySeq, *vm.Cons:
		// don't call RawCount — could be infinite
	default:
		if cc, ok := coll.(vm.Collection); ok {
			if cc.RawCount() == 0 {
				if hasInit {
					return initVal, nil
				}
				return ec.Invoke(mfn, nil)
			}
		}
	}
	seq, err := seqOf(coll)
	if err != nil {
		return vm.NIL, fmt.Errorf("reduce expected Seq")
	}
	// seqOf returns LazySeq objects without resolving them; an
	// unresolved-empty LazySeq is non-nil but yields First()=NIL.
	// Resolve here so empty inputs hit the early-return path
	// instead of spuriously iterating once with a NIL element.
	if ls, ok := seq.(*vm.LazySeq); ok {
		seq = ls.Resolve()
	}
	if seq == nil {
		if hasInit {
			return initVal, nil
		}
		return ec.Invoke(mfn, nil)
	}
	var acc vm.Value
	if hasInit {
		acc = initVal
	} else {
		acc = seq.First()
		seq = seq.Next()
	}
	for seq != nil {
		// Chunked fast path: when the source exposes a chunk, walk via
		// Nth in a tight inner loop and advance one chunk at a time. This
		// avoids the per-element LazySeq/Cons allocation churn that
		// dominates plain Next()-based reduce on chunked sources.
		if cs, ok := vm.AsChunkedSeq(seq); ok {
			c := cs.ChunkedFirst()
			n := c.ChunkCount()
			for i := 0; i < n; i++ {
				acc, err = ec.Invoke(mfn, []vm.Value{acc, c.Nth(i)})
				if err != nil {
					return vm.NIL, err
				}
				if r, ok := acc.(*vm.Reduced); ok {
					return r.Deref(), nil
				}
			}
			seq = cs.ChunkedNext()
			continue
		}
		acc, err = ec.Invoke(mfn, []vm.Value{acc, seq.First()})
		if err != nil {
			return vm.NIL, err
		}
		if r, ok := acc.(*vm.Reduced); ok {
			return r.Deref(), nil
		}
		seq = seq.Next()
	}

	return acc, nil
}
