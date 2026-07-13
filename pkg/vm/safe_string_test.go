/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"errors"
	"testing"
)

func TestSafeStringConvertsLazyThrowToError(t *testing.T) {
	boom := errors.New("boom while printing")
	thunk, err := NativeFnType.Wrap(func(_ []Value) (Value, error) {
		return NIL, boom
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := SafeString(NewLazySeq(thunk.(Fn)))
	if !errors.Is(err, boom) {
		t.Fatalf("SafeString error = %v, want %v", err, boom)
	}
	if got != "" {
		t.Fatalf("SafeString value = %q, want empty string on error", got)
	}
}
