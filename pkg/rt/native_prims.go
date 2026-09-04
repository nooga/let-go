/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"math"
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
	// Resolve a bytecode reducer once and reuse one frame for the whole
	// fold; per-element ec.Invoke pays resolution, frame-pool traffic, and
	// frame init on every element. Non-bytecode/variadic reducers fall back
	// to the generic invoke below.
	pc := ec.PrepareCall(mfn, 2)
	if pc != nil {
		defer pc.Release()
	}
	// Reused two-arg buffer, shared by every path below: ec.Invoke does not
	// retain its argument slice, so one buffer per reduce beats one per
	// element — the single largest allocation site in a reduce loop.
	fargs := []vm.Value{nil, nil}
	// ArrayVector fast path: flat slice with O(1) indexing — reduce by
	// direct index with zero seq/chunk allocation instead of via seqOf.
	if av, ok := coll.(vm.ArrayVector); ok {
		var acc vm.Value
		i := 0
		if hasInit {
			acc = initVal
		} else {
			acc = av[0]
			i = 1
		}
		for ; i < len(av); i++ {
			var res vm.Value
			var err error
			if pc != nil {
				res, err = pc.Call2(acc, av[i])
			} else {
				fargs[0] = acc
				fargs[1] = av[i]
				res, err = ec.Invoke(mfn, fargs)
			}
			if err != nil {
				return vm.NIL, err
			}
			if r, ok := res.(*vm.Reduced); ok {
				return r.Deref(), nil
			}
			acc = res
		}
		return acc, nil
	}
	// Range fast path: a range is (start, end, step) arithmetic — reduce
	// by direct iteration with zero seq/chunk allocation. Mirrors the
	// ArrayVector fast path above.
	if rng, ok := coll.(*vm.Range); ok {
		start, end, step := rng.Bounds()
		var acc vm.Value
		i := start
		if hasInit {
			acc = initVal
		} else {
			acc = vm.Int(i)
			i += step
		}
		for (step > 0 && i < end) || (step < 0 && i > end) {
			var res vm.Value
			var err error
			if pc != nil {
				res, err = pc.Call2(acc, vm.Int(i))
			} else {
				fargs[0] = acc
				fargs[1] = vm.Int(i)
				res, err = ec.Invoke(mfn, fargs)
			}
			if err != nil {
				return vm.NIL, err
			}
			if r, ok := res.(*vm.Reduced); ok {
				return r.Deref(), nil
			}
			acc = res
			i += step
		}
		return acc, nil
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
				if pc != nil {
					acc, err = pc.Call2(acc, c.Nth(i))
				} else {
					fargs[0] = acc
					fargs[1] = c.Nth(i)
					acc, err = ec.Invoke(mfn, fargs)
				}
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
		if pc != nil {
			acc, err = pc.Call2(acc, seq.First())
		} else {
			fargs[0] = acc
			fargs[1] = seq.First()
			acc, err = ec.Invoke(mfn, fargs)
		}
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

// Namespace implements (namespace x) for Named values.
//
//lg:native
//lg:name namespace
func Namespace(v vm.Value) (vm.Value, error) {
	named, ok := v.(vm.Named)
	if !ok {
		return vm.NIL, fmt.Errorf("namespace expected Named")
	}
	return named.Namespace(), nil
}

// PushBinding implements (push-binding! v val) — the dynamic-binding
// lifecycle, so it needs the ExecContext holding the binding stack.
//
//lg:native
//lg:name push-binding!
func PushBinding(ec *vm.ExecContext, v vm.Value, val vm.Value) (vm.Value, error) {
	vr, ok := v.(*vm.Var)
	if !ok {
		return vm.NIL, fmt.Errorf("push-binding expected Var")
	}
	ec.PushBinding(vr, val)
	return vm.NIL, nil
}

// PopBinding implements (pop-binding! v).
//
//lg:native
//lg:name pop-binding!
func PopBinding(ec *vm.ExecContext, v vm.Value) (vm.Value, error) {
	vr, ok := v.(*vm.Var)
	if !ok {
		return vm.NIL, fmt.Errorf("pop-binding expected Var")
	}
	ec.PopBinding(vr)
	return vm.NIL, nil
}

// regexSubmatchVector renders every capture group, with nil for groups that
// did not participate in the match. Promoted from a lang.go closure alongside
// regexSubmatchValue; both original call sites are same-package and unaffected.
func regexSubmatchVector(s string, indices []int) vm.ArrayVector {
	result := make(vm.ArrayVector, len(indices)/2)
	for i := range result {
		start, end := indices[2*i], indices[2*i+1]
		if start < 0 {
			result[i] = vm.NIL
			continue
		}
		result[i] = vm.String(s[start:end])
	}
	return result
}

// regexSubmatchValue shapes a submatch-index result: a bare match returns the
// String, a grouped match returns the vector. Promoted from a closure inside
// lang.go's install fn so the direct-call natives here can share it — the two
// original call sites are in the same package and are unaffected.
func regexSubmatchValue(s string, indices []int) vm.Value {
	if len(indices) == 2 {
		return vm.String(s[indices[0]:indices[1]])
	}
	return regexSubmatchVector(s, indices)
}

// Some implements (some pred coll). Body extracted verbatim from lang.go's
// closure (including the chunked fast path) rather than retyped, so the
// direct-call path cannot drift from the interpreter path. Needs the
// ExecContext to invoke the predicate.
//
//lg:native
//lg:name some
func Some(ec *vm.ExecContext, pred vm.Value, coll vm.Value) (vm.Value, error) {
	f, ok := vm.AsFn(pred)
	if !ok {
		return vm.NIL, fmt.Errorf("some expected Fn")
	}
	seq, err := seqOf(coll)
	if err != nil {
		return vm.NIL, fmt.Errorf("some expected Seq")
	}
	// seqOf returns LazySeq directly; resolve here so an unresolved-empty
	// LazySeq becomes nil instead of invoking the predicate once with nil.
	if ls, ok := seq.(*vm.LazySeq); ok {
		seq = ls.Resolve()
	}
	// Resolve a bytecode predicate once and reuse one frame for the whole
	// walk: per-element ec.Invoke would pay resolution, frame-pool traffic,
	// and frame init on every element. Non-bytecode/variadic predicates fall
	// back to the generic invoke below.
	pc := ec.PrepareCall(f, 1)
	if pc != nil {
		defer pc.Release()
	}
	fargs := []vm.Value{nil}
	for seq != nil {
		// Chunked fast path: walk whole chunks via Nth to avoid per-element
		// seq.Next()/First() churn on chunked sources such as vectors/ranges.
		if cs, ok := vm.AsChunkedSeq(seq); ok {
			c := cs.ChunkedFirst()
			n := c.ChunkCount()
			for i := 0; i < n; i++ {
				var v vm.Value
				var err error
				if pc != nil {
					v, err = pc.Call1(c.Nth(i))
				} else {
					fargs[0] = c.Nth(i)
					v, err = ec.Invoke(f, fargs)
				}
				if err != nil {
					return vm.NIL, err
				}
				if vm.IsTruthy(v) {
					return v, nil
				}
			}
			seq = cs.ChunkedNext()
			continue
		}
		var v vm.Value
		var err error
		if pc != nil {
			v, err = pc.Call1(seq.First())
		} else {
			fargs[0] = seq.First()
			v, err = ec.Invoke(f, fargs)
		}
		if err != nil {
			return vm.NIL, err
		}
		if vm.IsTruthy(v) {
			return v, nil
		}
		seq = seq.Next()
	}

	return vm.NIL, nil
}

// Int implements (int x). Body extracted verbatim from lang.go's closure —
// including the int32 range guard and the per-type coercion table — rather
// than retyped, so the direct-call path cannot drift from the interpreter.
//
//lg:native
//lg:name int
func Int(v vm.Value) (vm.Value, error) {
	const minInt32 = -2147483648
	const maxInt32 = 2147483647
	coerce := func(f float64) (vm.Value, error) {
		if f < minInt32 || f > maxInt32 {
			return vm.NIL, fmt.Errorf("%s can't be coerced to int", v)
		}
		return vm.MakeInt(int(math.Trunc(f))), nil
	}
	switch v := v.(type) {
	case vm.Int:
		if v < minInt32 || v > maxInt32 {
			return vm.NIL, fmt.Errorf("%s can't be coerced to int", v)
		}
		return v, nil
	case vm.Float:
		return coerce(float64(v))
	case vm.Char:
		return vm.Int(int(v)), nil
	case *vm.BigInt:
		if !v.Val().IsInt64() {
			return vm.NIL, fmt.Errorf("%s can't be coerced to int", v)
		}
		i := v.Val().Int64()
		if i < minInt32 || i > maxInt32 {
			return vm.NIL, fmt.Errorf("%s can't be coerced to int", v)
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
		return vm.NIL, fmt.Errorf("%s can't be coerced to int", v)
	}
}

// rangeInt coerces a range bound: Int passes through; anything else is a
// graceful error (a raw .(vm.Int) assertion panics). Promoted from a lang.go
// closure so the direct-call range arities can share it.
func rangeInt(v vm.Value, what string) (vm.Int, error) {
	if n, ok := v.(vm.Int); ok {
		return n, nil
	}
	return 0, fmt.Errorf("range %s must be an integer, got %s",
		what, v.Type().Name())
}
