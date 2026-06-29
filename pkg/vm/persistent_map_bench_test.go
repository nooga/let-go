/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

func buildPersistentMapForBench(size int) *PersistentMap {
	m := EmptyPersistentMap
	for i := range size {
		m = m.Assoc(Int(i), Int(i)).(*PersistentMap)
	}
	return m
}

func benchmarkPersistentMapAssocOverwrite(b *testing.B, size int) {
	base := buildPersistentMapForBench(size)
	keys := make([]Value, size)
	for i := range size {
		keys[i] = Int(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%size]
		_ = base.Assoc(key, Int(i)).(*PersistentMap)
	}
}

func benchmarkTransientMapAssocOverwrite(b *testing.B, size int) {
	base := buildPersistentMapForBench(size)
	keys := make([]Value, size)
	for i := range size {
		keys[i] = Int(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm := NewTransientMap(base)
		if _, err := tm.Assoc(keys[i%size], Int(i)); err != nil {
			b.Fatal(err)
		}
		if _, err := tm.Persistent(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkPersistentMapAssocOverwriteBatch(b *testing.B, size int) {
	base := buildPersistentMapForBench(size)
	keys := make([]Value, size)
	for i := range size {
		keys[i] = Int(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := base
		for j := range size {
			key := keys[j]
			m = m.Assoc(key, Int(i+j)).(*PersistentMap)
		}
	}
}

func benchmarkTransientMapAssocOverwriteBatch(b *testing.B, size int) {
	base := buildPersistentMapForBench(size)
	keys := make([]Value, size)
	for i := range size {
		keys[i] = Int(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm := NewTransientMap(base)
		for j := range size {
			if _, err := tm.Assoc(keys[j], Int(i+j)); err != nil {
				b.Fatal(err)
			}
		}
		if _, err := tm.Persistent(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkPersistentMapAssocInsertBatch(b *testing.B, batch int) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := EmptyPersistentMap
		for j := range batch {
			k := Int(i*batch + j)
			m = m.Assoc(k, k).(*PersistentMap)
		}
	}
}

func benchmarkTransientMapAssocInsertBatch(b *testing.B, batch int) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm := NewTransientMap(EmptyPersistentMap)
		for j := range batch {
			k := Int(i*batch + j)
			if _, err := tm.Assoc(k, k); err != nil {
				b.Fatal(err)
			}
		}
		if _, err := tm.Persistent(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPersistentMapAssocOverwrite1K(b *testing.B) {
	benchmarkPersistentMapAssocOverwrite(b, 1024)
}

func BenchmarkTransientMapAssocOverwrite1K(b *testing.B) {
	benchmarkTransientMapAssocOverwrite(b, 1024)
}

func BenchmarkPersistentMapAssocInsertBatch1K(b *testing.B) {
	benchmarkPersistentMapAssocInsertBatch(b, 1024)
}

func BenchmarkTransientMapAssocInsertBatch1K(b *testing.B) {
	benchmarkTransientMapAssocInsertBatch(b, 1024)
}

func BenchmarkPersistentMapAssocOverwriteBatch1K(b *testing.B) {
	benchmarkPersistentMapAssocOverwriteBatch(b, 1024)
}

func BenchmarkTransientMapAssocOverwriteBatch1K(b *testing.B) {
	benchmarkTransientMapAssocOverwriteBatch(b, 1024)
}

func benchmarkDepsTransientMapPersistentSet(b *testing.B, keys, valsPerKey int) {
	keyVals := make([]Value, keys)
	for i := range keys {
		keyVals[i] = Int(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm := NewTransientMap(EmptyPersistentMap)
		for k := range keys {
			key := keyVals[k]
			cur := EmptyPersistentSet
			if v := tm.ValueAt(key); v != NIL {
				cur = v.(*PersistentSet)
			}
			for j := range valsPerKey {
				cur = cur.Conj(Int(i*valsPerKey + j)).(*PersistentSet)
			}
			if _, err := tm.Assoc(key, cur); err != nil {
				b.Fatal(err)
			}
		}
		if _, err := tm.Persistent(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkDepsTransientMapTransientSet(b *testing.B, keys, valsPerKey int) {
	keyVals := make([]Value, keys)
	for i := range keys {
		keyVals[i] = Int(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm := NewTransientMap(EmptyPersistentMap)
		for k := range keys {
			key := keyVals[k]
			ts := NewTransientSet(EmptyPersistentSet)
			if v := tm.ValueAt(key); v != NIL {
				existing, ok := v.(*TransientSet)
				if ok {
					ts = existing
				}
			}
			for j := range valsPerKey {
				if _, err := ts.Conj(Int(i*valsPerKey + j)); err != nil {
					b.Fatal(err)
				}
			}
			if _, err := tm.Assoc(key, ts); err != nil {
				b.Fatal(err)
			}
		}
		if _, err := tm.Persistent(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkDepsTransientMapVector(b *testing.B, keys, valsPerKey int) {
	keyVals := make([]Value, keys)
	for i := range keys {
		keyVals[i] = Int(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm := NewTransientMap(EmptyPersistentMap)
		for k := range keys {
			key := keyVals[k]
			tv := NewTransientVector(nil)
			if v := tm.ValueAt(key); v != NIL {
				existing, ok := v.(*TransientVector)
				if ok {
					tv = existing
				}
			}
			for j := range valsPerKey {
				if _, err := tv.Conj(Int(i*valsPerKey + j)); err != nil {
					b.Fatal(err)
				}
			}
			if _, err := tm.Assoc(key, tv); err != nil {
				b.Fatal(err)
			}
		}
		if _, err := tm.Persistent(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDepsTransientMapPersistentSet64x16(b *testing.B) {
	benchmarkDepsTransientMapPersistentSet(b, 64, 16)
}

func BenchmarkDepsTransientMapTransientSet64x16(b *testing.B) {
	benchmarkDepsTransientMapTransientSet(b, 64, 16)
}

func BenchmarkDepsTransientMapVector64x16(b *testing.B) {
	benchmarkDepsTransientMapVector(b, 64, 16)
}
