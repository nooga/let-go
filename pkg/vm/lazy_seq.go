/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "sync"

// LazySeq delays computation of a sequence until first/next is called.
// This is the foundation for lazy operations like map, filter, etc.
type LazySeq struct {
	fn Fn       // thunk that produces the seq when called
	s  Seq      // cached realized seq
	sv Value    // intermediate value from fn
	mu sync.Mutex
}

func NewLazySeq(fn Fn) *LazySeq {
	return &LazySeq{fn: fn}
}

// seq realizes the lazy seq if not already done
func (l *LazySeq) seq() Seq {
	l.mu.Lock()
	defer l.mu.Unlock()

	// If we have the thunk, call it
	if l.fn != nil {
		sv, err := l.fn.Invoke(nil)
		if err != nil {
			return nil // TODO: handle error better
		}
		l.sv = sv
		l.fn = nil
	}

	// If we have an intermediate value, convert to seq
	if l.sv != nil {
		if l.sv == NIL {
			l.s = nil
		} else if seq, ok := l.sv.(Seq); ok {
			l.s = seq
		} else if seqable, ok := l.sv.(Sequable); ok {
			l.s = seqable.Seq()
		}
		l.sv = nil
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
	s := l.seq()
	if s == nil {
		return NIL
	}
	return s.First()
}

func (l *LazySeq) More() Seq {
	s := l.seq()
	if s == nil {
		return EmptyList
	}
	return s.More()
}

func (l *LazySeq) Next() Seq {
	s := l.seq()
	if s == nil {
		return nil
	}
	return s.Next()
}

func (l *LazySeq) Cons(val Value) Seq {
	return NewCons(val, l)
}

func (l *LazySeq) Seq() Seq {
	return l.seq()
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

