/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"reflect"
	"strings"
)

// PersistentQueue is an immutable FIFO queue with the same shape as Clojure's
// clojure.lang.PersistentQueue: a `front` seq holding the elements ready to be
// read/popped, and a `rear` slice collecting newly conj'd elements. When the
// front drains, the rear is promoted (in insertion order) to become the new
// front. peek reads the front head, pop drops it, conj appends to the rear —
// each O(1) amortized, and every operation returns a new queue (the receiver
// is never mutated).

type theQueueType struct{}

func (t *theQueueType) String() string  { return t.Name() }
func (t *theQueueType) Type() ValueType { return TypeType }
func (t *theQueueType) Unbox() any      { return reflect.TypeFor[*theQueueType]() }
func (t *theQueueType) Name() string    { return "clojure.lang.PersistentQueue" }

func (t *theQueueType) Box(bare any) (Value, error) {
	arr, ok := bare.([]Value)
	if !ok {
		return EmptyPersistentQueue, NewTypeError(bare, "can't be boxed as", t)
	}
	var q Collection = EmptyPersistentQueue
	for _, v := range arr {
		q = q.Conj(v)
	}
	return q, nil
}

// QueueType is the type of PersistentQueues.
var QueueType *theQueueType = &theQueueType{}

// PersistentQueue is an immutable FIFO queue. The empty queue has front == nil.
type PersistentQueue struct {
	front Seq // elements ready to read; nil only when the whole queue is empty
	rear  []Value
	count int
	meta  Value
}

// EmptyPersistentQueue is the canonical empty queue (clojure.lang.PersistentQueue/EMPTY).
var EmptyPersistentQueue *PersistentQueue = &PersistentQueue{}

// Type implements Value.
func (q *PersistentQueue) Type() ValueType { return QueueType }

// Unbox implements Value — a flat []Value in FIFO order.
func (q *PersistentQueue) Unbox() any {
	var out []Value // grown by append; see Seq() on the omitted make() capacity
	for s := q.Seq(); s != nil && s != EmptyList; s = s.Next() {
		out = append(out, s.First())
	}
	return out
}

// Meta / WithMeta implement IMeta.
func (q *PersistentQueue) Meta() Value {
	if q.meta == nil {
		return NIL
	}
	return q.meta
}

func (q *PersistentQueue) WithMeta(m Value) Value {
	cp := *q
	cp.meta = m
	return &cp
}

// RawCount / Count implement Counted.
func (q *PersistentQueue) RawCount() int { return q.count }
func (q *PersistentQueue) Count() Value  { return Int(q.count) }

// Empty implements Collection.
func (q *PersistentQueue) Empty() Collection {
	if q.meta != nil {
		return &PersistentQueue{meta: q.meta}
	}
	return EmptyPersistentQueue
}

// Conj implements Collection — appends to the rear (Clojure enqueue semantics).
// The very first element seeds the front instead, so peek always has somewhere
// to read from.
func (q *PersistentQueue) Conj(value Value) Collection {
	if q.count == 0 {
		return &PersistentQueue{front: EmptyList.Cons(value), count: 1, meta: q.meta}
	}
	nr := make([]Value, len(q.rear)+1)
	copy(nr, q.rear)
	nr[len(q.rear)] = value
	return &PersistentQueue{front: q.front, rear: nr, count: q.count + 1, meta: q.meta}
}

// Seq implements Sequable — front elements first, then the rear in FIFO order.
// Returns EmptyList for the empty queue (the seq builtin maps that to nil).
func (q *PersistentQueue) Seq() Seq {
	if q.count == 0 {
		return EmptyList
	}
	// Grown by append rather than pre-sized from q.count: q.count is the
	// already-materialized element count (front seq + rear), not an untrusted
	// size, but leaving the make() capacity off keeps CodeQL's allocation-size
	// check quiet and matches Record.Seq().
	var elems []Value
	for s := q.front; s != nil && s != EmptyList; s = s.Next() {
		elems = append(elems, s.First())
	}
	elems = append(elems, q.rear...)
	l, _ := ListType.Box(elems)
	return l.(*List)
}

// Peek returns the front element, or NIL when empty.
func (q *PersistentQueue) Peek() Value {
	if q.count == 0 {
		return NIL
	}
	return q.front.First()
}

// Pop returns the queue without its front element. Popping the empty queue
// returns it unchanged (matching Clojure, which does not throw here). When the
// front seq drains, the rear is promoted to become the new front in FIFO order.
func (q *PersistentQueue) Pop() *PersistentQueue {
	if q.count == 0 {
		return q
	}
	f1 := q.front.Next() // nil once the front held a single element
	if f1 == nil {
		// Front drained. count-1 == len(rear); promote the rear (in order),
		// or fall back to the empty queue when the rear is also empty.
		if len(q.rear) == 0 {
			return q.Empty().(*PersistentQueue)
		}
		nf, _ := ListType.Box(q.rear)
		return &PersistentQueue{front: nf.(*List), count: q.count - 1, meta: q.meta}
	}
	return &PersistentQueue{front: f1, rear: q.rear, count: q.count - 1, meta: q.meta}
}

// Hash implements Hashable — order-sensitive, over the FIFO sequence.
func (q *PersistentQueue) Hash() uint32 {
	if q.count == 0 {
		return 0
	}
	return hashOrdered(q.Seq())
}

// Equals implements value equality — another queue with the same elements in
// the same FIFO order.
func (q *PersistentQueue) Equals(other Value) bool {
	o, ok := other.(*PersistentQueue)
	if !ok || q.count != o.count {
		return false
	}
	a, b := q.Seq(), o.Seq()
	for a != nil && a != EmptyList {
		if !valueEquiv(a.First(), b.First()) {
			return false
		}
		a, b = a.Next(), b.Next()
	}
	return true
}

func (q *PersistentQueue) String() string {
	b := &strings.Builder{}
	b.WriteString("#queue (")
	first := true
	for s := q.Seq(); s != nil && s != EmptyList; s = s.Next() {
		if !first {
			b.WriteRune(' ')
		}
		b.WriteString(s.First().String())
		first = false
	}
	b.WriteString(")")
	return b.String()
}
