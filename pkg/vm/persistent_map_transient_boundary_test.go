/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"sync/atomic"
	"testing"
)

// The in-place transient insert reuses an owned node's spare capacity, so
// its correctness hangs on two boundaries: the shift must be right for any
// insert position (begin/middle/end), and spare capacity on a persisted
// node must never be reused by a later transient (ensureEditable must
// clone). These tests pin both directly.

// nodeSlots reads back the (slot, key) pairs of a bitmap node in slot order.
func nodeSlots(t *testing.T, n *hmapBitmapNode) []string {
	t.Helper()
	var out []string
	for slot := uint32(0); slot < 32; slot++ {
		bit := uint32(1) << slot
		if n.bitmap&bit == 0 {
			continue
		}
		idx := hmapIndex(n.bitmap, bit)
		k, ok := n.array[2*idx].(Value)
		if !ok {
			t.Fatalf("slot %d: key is %T, want Value", slot, n.array[2*idx])
		}
		out = append(out, fmt.Sprintf("%d=%s", slot, k.String()))
	}
	return out
}

// Driving assocTransient with crafted hashes (slot = hash & 0x1f at shift 0)
// gives deterministic control of the insert position within the node array.
func TestAssocTransientInsertPositions(t *testing.T) {
	edit := new(atomic.Bool)
	edit.Store(true)
	n := &hmapBitmapNode{edit: edit, array: make([]any, 0, 64)}

	var added bool
	insert := func(slot uint32) {
		nn := n.assocTransient(edit, 0, slot, String(fmt.Sprintf("k%d", slot)), Int(int(slot)), &added)
		if nn != n {
			t.Fatalf("slot %d: owned insert returned a new node (in-place path not taken)", slot)
		}
	}

	// middle anchor, then end, begin, and interior positions — each order
	// exercises a different shift span in the in-place branch.
	for _, slot := range []uint32{16, 31, 0, 8, 24, 4} {
		insert(slot)
	}
	if cap(n.array) != 64 {
		t.Fatalf("cap = %d, want 64 (in-place growth must not reallocate)", cap(n.array))
	}
	want := []string{`0="k0"`, `4="k4"`, `8="k8"`, `16="k16"`, `24="k24"`, `31="k31"`}
	got := nodeSlots(t, n)
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("slot layout after in-place inserts:\n got %v\nwant %v", got, want)
	}
}

// Capacity-exhaustion boundary: an owned node with len==cap must take the
// realloc branch, preserve order across it, and then grow in place again
// within the new headroom.
func TestAssocTransientReallocBoundary(t *testing.T) {
	edit := new(atomic.Bool)
	edit.Store(true)
	n := &hmapBitmapNode{
		edit:   edit,
		bitmap: 1 << 5,
		array:  []any{String("k5"), Int(5)}, // len == cap == 2, no reserve
	}

	var added bool
	// Descending-position insert across the realloc: k2 lands before k5.
	n = n.assocTransient(edit, 0, 2, String("k2"), Int(2), &added)
	if len(n.array) != 4 {
		t.Fatalf("len = %d, want 4", len(n.array))
	}
	if cap(n.array) != 8 { // min(2*need, 64) with need=4
		t.Fatalf("cap = %d, want 8 (realloc headroom)", cap(n.array))
	}
	want := []string{`2="k2"`, `5="k5"`}
	if got := nodeSlots(t, n); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("layout after realloc insert:\n got %v\nwant %v", got, want)
	}

	// Next inserts fit the new headroom: node and backing array must be
	// reused (in-place), including another begin-position shift.
	before := n
	n = n.assocTransient(edit, 0, 9, String("k9"), Int(9), &added)
	n = n.assocTransient(edit, 0, 0, String("k0"), Int(0), &added)
	if n != before {
		t.Fatal("in-headroom inserts reallocated the node")
	}
	if cap(n.array) != 8 || len(n.array) != 8 {
		t.Fatalf("len/cap = %d/%d, want 8/8", len(n.array), cap(n.array))
	}
	want = []string{`0="k0"`, `2="k2"`, `5="k5"`, `9="k9"`}
	if got := nodeSlots(t, n); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("layout after in-headroom inserts:\n got %v\nwant %v", got, want)
	}
}

// A persisted map built transiently carries nodes with spare capacity. A
// second transient over that map must clone before mutating — never write
// into the frozen map's spare slots.
func TestTransientPersistedNodeNotReused(t *testing.T) {
	t1 := NewTransientMap(EmptyPersistentMap)
	const n = 20
	for i := 0; i < n; i++ {
		if _, err := t1.Assoc(Int(i), Int(i*10)); err != nil {
			t.Fatal(err)
		}
	}
	m1, err := t1.Persistent()
	if err != nil {
		t.Fatal(err)
	}

	// Mutate a second transient over the persisted result: new keys plus
	// overwrites of existing ones.
	t2 := NewTransientMap(m1)
	for i := 0; i < n; i++ {
		if _, err := t2.Assoc(Int(100+i), Int(i)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := t2.Assoc(Int(3), String("overwritten")); err != nil {
		t.Fatal(err)
	}

	// m1 must be untouched: same count, original values, no new keys.
	if m1.RawCount() != n {
		t.Fatalf("persisted map count changed: %d, want %d", m1.RawCount(), n)
	}
	for i := 0; i < n; i++ {
		if got := m1.ValueAtOr(Int(i), NIL); got != Int(i*10) {
			t.Fatalf("persisted map key %d = %v, want %d", i, got, i*10)
		}
	}
	if got := m1.ValueAtOr(Int(105), NIL); got != NIL {
		t.Fatalf("persisted map leaked a second-transient key: %v", got)
	}

	m2, err := t2.Persistent()
	if err != nil {
		t.Fatal(err)
	}
	if m2.RawCount() != 2*n {
		t.Fatalf("second map count = %d, want %d", m2.RawCount(), 2*n)
	}
	if got := m2.ValueAtOr(Int(3), NIL); got != String("overwritten") {
		t.Fatalf("second map overwrite lost: %v", got)
	}
}
