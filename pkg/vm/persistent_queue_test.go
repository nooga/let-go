/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

func mkQueue(vs ...int) *PersistentQueue {
	var q Collection = EmptyPersistentQueue
	for _, v := range vs {
		q = q.Conj(Int(v))
	}
	return q.(*PersistentQueue)
}

func queueElems(q *PersistentQueue) []int {
	var out []int
	for s := q.Seq(); s != nil && s != EmptyList; s = s.Next() {
		out = append(out, int(s.First().(Int)))
	}
	return out
}

func TestQueueEmpty(t *testing.T) {
	q := EmptyPersistentQueue
	if q.RawCount() != 0 {
		t.Fatalf("empty RawCount = %d, want 0", q.RawCount())
	}
	if q.Peek() != NIL {
		t.Fatalf("empty Peek = %v, want NIL", q.Peek())
	}
	if q.Pop() != q {
		t.Fatalf("Pop of empty should return the empty queue unchanged")
	}
	if q.Seq() != EmptyList {
		t.Fatalf("empty Seq should be EmptyList")
	}
}

func TestQueueFIFOPeekPop(t *testing.T) {
	q := mkQueue(1, 2, 3)
	if got := q.Peek(); got != Int(1) {
		t.Fatalf("Peek = %v, want 1", got)
	}
	if got := q.Pop().Peek(); got != Int(2) {
		t.Fatalf("Pop().Peek() = %v, want 2", got)
	}
	if got := q.Pop().Pop().Peek(); got != Int(3) {
		t.Fatalf("Pop().Pop().Peek() = %v, want 3", got)
	}
	if got := q.RawCount(); got != 3 {
		t.Fatalf("RawCount = %d, want 3", got)
	}
}

func TestQueueImmutable(t *testing.T) {
	q := mkQueue(1, 2, 3)
	_ = q.Conj(Int(4))
	_ = q.Pop()
	if got := queueElems(q); len(got) != 3 || got[0] != 1 || got[2] != 3 {
		t.Fatalf("original queue mutated: %v", got)
	}
}

func TestQueueConjAppendsRear(t *testing.T) {
	q := mkQueue(1, 2, 3).Conj(Int(4)).(*PersistentQueue)
	if got := q.Peek(); got != Int(1) {
		t.Fatalf("conj should not change the front: Peek = %v, want 1", got)
	}
	if got := queueElems(q); len(got) != 4 || got[3] != 4 {
		t.Fatalf("conj should append to the rear: %v", got)
	}
}

// Exercises front/rear rebalancing: pop past the seeded front so the rear is
// promoted, repeatedly, over a large run.
func TestQueueFIFOLarge(t *testing.T) {
	const n = 500
	q := EmptyPersistentQueue
	for i := 0; i < n; i++ {
		q = q.Conj(Int(i)).(*PersistentQueue)
	}
	if q.RawCount() != n {
		t.Fatalf("RawCount = %d, want %d", q.RawCount(), n)
	}
	for i := 0; i < n; i++ {
		if got := q.Peek(); got != Int(i) {
			t.Fatalf("drain step %d: Peek = %v, want %d", i, got, i)
		}
		q = q.Pop()
	}
	if q.RawCount() != 0 {
		t.Fatalf("drained queue RawCount = %d, want 0", q.RawCount())
	}
	if q.Peek() != NIL {
		t.Fatalf("drained queue Peek = %v, want NIL", q.Peek())
	}
}

func TestQueueSeqOrder(t *testing.T) {
	got := queueElems(mkQueue(10, 20, 30))
	if len(got) != 3 || got[0] != 10 || got[1] != 20 || got[2] != 30 {
		t.Fatalf("Seq order = %v, want [10 20 30]", got)
	}
}

func TestQueueEquals(t *testing.T) {
	a := mkQueue(1, 2, 3)
	// Same elements, but built via a conj-then-drain path so the internal
	// front/rear split differs — Equals must still hold.
	b := mkQueue(9, 1, 2, 3).Pop()
	if !a.Equals(b) {
		t.Fatalf("queues with equal FIFO contents should be Equals")
	}
	if a.Equals(mkQueue(1, 2)) {
		t.Fatalf("queues of different length must not be Equals")
	}
	if a.Equals(mkQueue(3, 2, 1)) {
		t.Fatalf("order-different queues must not be Equals")
	}
	if a.Equals(NewList([]Value{Int(1), Int(2), Int(3)})) {
		t.Fatalf("a queue must not equal a non-queue value")
	}
}

func TestQueueType(t *testing.T) {
	if mkQueue(1).Type() != QueueType {
		t.Fatalf("queue Type should be QueueType")
	}
	if QueueType.Name() != "clojure.lang.PersistentQueue" {
		t.Fatalf("QueueType.Name = %q", QueueType.Name())
	}
}
