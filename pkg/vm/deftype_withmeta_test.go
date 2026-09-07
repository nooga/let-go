/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"sync/atomic"
	"testing"
)

// WithMeta must not copy the rendering guard. The guard is 1 while the
// instance's toString override runs, and a whole-struct copy taken at that
// moment — (with-meta this ...) inside the override is the natural way —
// would never run its own override again. White-box: raise the guard and
// check the copy starts clean. The end-to-end shape lives in
// test/reify_object_test.lg (reify-object-with-meta-inside-tostring).
func TestDTypeInstance_WithMeta_FreshRenderingGuard(t *testing.T) {
	dt := NewDType("Point", []Symbol{"x"})
	inst := NewDTypeInstance(dt, []Value{Int(1)})
	atomic.StoreUint32(&inst.rendering, 1)

	cp, ok := inst.WithMeta(NIL).(*DTypeInstance)
	if !ok {
		t.Fatalf("WithMeta returned %T, want *DTypeInstance", inst.WithMeta(NIL))
	}
	if atomic.LoadUint32(&cp.rendering) != 0 {
		t.Fatal("WithMeta copied the rendering guard; a copy must start at 0")
	}
	if cp.dtype != dt {
		t.Fatal("WithMeta lost the dtype")
	}
	if &cp.fields[0] != &inst.fields[0] {
		t.Fatal("WithMeta should share field storage with the original")
	}
}
