/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// A native-baked multimethod is frozen once at namespace-load completion. The
// frozen flag is the currency signal the generated native type-switch consults:
// while the var still holds the frozen *MultiFn, the baked arms are current.
// AddMethod / RemoveMethod return a NEW, unfrozen *MultiFn, so any late
// defmethod replaces the var's value with an unfrozen one and the generated
// guard falls back to runtime dispatch.
func TestMultiFnFreezeNativeInvalidatedByAddMethod(t *testing.T) {
	m := NewMultiFn("mm", nil, Keyword("default"))

	if m.IsNativeFrozen() {
		t.Fatal("a fresh MultiFn must not be native-frozen")
	}

	m.FreezeNative()
	if !m.IsNativeFrozen() {
		t.Fatal("FreezeNative must mark the MultiFn frozen")
	}

	// A late defmethod produces a new MultiFn — it must NOT be frozen, so the
	// generated guard sees the change and stops trusting the baked arms.
	extended := m.AddMethod(Keyword("x"), nil)
	if extended.IsNativeFrozen() {
		t.Fatal("AddMethod result must not inherit the frozen flag")
	}
	// The original frozen pointer is untouched (the var no longer points at it,
	// but freezing in place must not have been clobbered by AddMethod).
	if !m.IsNativeFrozen() {
		t.Fatal("AddMethod must not unfreeze the original MultiFn")
	}

	removed := m.RemoveMethod(Keyword("x"))
	if removed.IsNativeFrozen() {
		t.Fatal("RemoveMethod result must not inherit the frozen flag")
	}
}
