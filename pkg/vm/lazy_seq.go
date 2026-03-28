/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "sync"

// LazySeq delays computation of a sequence until first/next is called.
// This is the foundation for lazy operations like map, filter, etc.
type LazySeq struct {
	fn  Fn       // thunk that produces the seq when called
	s   Seq      // cached realized seq
	sv  Value    // intermediate value from fn
	err error    // error from thunk realization (propagated on access)
	mu  sync.Mutex
}

func NewLazySeq(fn Fn) *LazySeq {
	return &LazySeq{fn: fn}
}

// sval realizes the thunk and returns the raw value without converting to seq.
// Used for unwrapping nested LazySeqs without locking issues.
func (l *LazySeq) sval() Value {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.fn != nil {
		sv, err := l.fn.Invoke(nil)
		if err != nil {
			l.err = err
			l.fn = nil
			panic(&thrownPanic{err: err})
		}
		l.sv = sv
		l.fn = nil
	}
	if l.sv != nil {
		return l.sv
	}
	// Already fully resolved to l.s
	if l.s != nil {
		return l.s
	}
	return nil
}

// seq realizes the lazy seq if not already done
func (l *LazySeq) seq() Seq {
	l.sval() // ensure thunk is called (may panic with thrownPanic)

	l.mu.Lock()
	defer l.mu.Unlock()

	// Re-raise stored error on repeated access
	if l.err != nil {
		panic(&thrownPanic{err: l.err})
	}

	if l.s != nil {
		return l.s
	}

	if l.sv != nil {
		sv := l.sv
		l.sv = nil
		if sv == nil || sv == NIL {
			l.s = nil
		} else if seq, ok := sv.(Seq); ok {
			l.s = seq
		} else if seqable, ok := sv.(Sequable); ok {
			l.s = seqable.Seq()
		}
	}

	return l.s
}

func (l *LazySeq) String() string {
	s := l.seq()
	if s == nil {
		return "()"
	}
	return s.String()
}

func (l *LazySeq) Type() ValueType    { return ListType }
func (l *LazySeq) Unbox() interface{} { return l.seq() }

func (l *LazySeq) First() Value {
	s := l.Resolve()
	if s == nil {
		return NIL
	}
	return s.First()
}

func (l *LazySeq) More() Seq {
	s := l.Resolve()
	if s == nil {
		return EmptyList
	}
	return s.More()
}

func (l *LazySeq) Next() Seq {
	s := l.Resolve()
	if s == nil {
		return nil
	}
	return s.Next()
}

// Resolve returns the fully realized non-LazySeq inner seq.
// Iteratively unwraps nested LazySeqs without accumulating Go stack.
func (l *LazySeq) Resolve() Seq {
	s := l.seq()
	for s != nil {
		inner, ok := s.(*LazySeq)
		if !ok {
			return s
		}
		s = inner.seq()
	}
	return nil
}

func (l *LazySeq) Cons(val Value) Seq {
	return NewCons(val, l)
}

func (l *LazySeq) Seq() Seq {
	return l.Resolve()
}

func (l *LazySeq) Count() Value {
	n := 0
	for s := l.seq(); s != nil && s != EmptyList; s = s.Next() {
		n++
	}
	return Int(n)
}

func (l *LazySeq) RawCount() int {
	n := 0
	for s := l.seq(); s != nil && s != EmptyList; s = s.Next() {
		n++
	}
	return n
}

func (l *LazySeq) Empty() Collection { return EmptyList }

func (l *LazySeq) Conj(val Value) Collection {
	return NewCons(val, l)
}

func (l *LazySeq) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l *LazySeq) ValueAtOr(key Value, notFound Value) Value {
	idx, ok := key.(Int)
	if !ok {
		return notFound
	}
	i := int(idx)
	s := l.seq()
	for j := 0; s != nil && s != EmptyList; j++ {
		if j == i {
			return s.First()
		}
		s = s.Next()
	}
	return notFound
}

