/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"runtime"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func stringFn(t *testing.T, name string) vm.Fn {
	t.Helper()
	v := NS("string").Lookup(vm.Symbol(name))
	if v == nil || v == vm.NIL {
		t.Fatalf("string/%s not found", name)
	}
	if vr, ok := v.(*vm.Var); ok {
		v = vr.Deref()
	}
	fn, ok := v.(vm.Fn)
	if !ok {
		t.Fatalf("string/%s is not an Fn, got %T", name, v)
	}
	return fn
}

// bytesPerCall reports the bytes allocated per call of f, averaged over n.
func bytesPerCall(n int, f func()) uint64 {
	f()
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	for i := 0; i < n; i++ {
		f()
	}
	runtime.ReadMemStats(&after)
	return (after.TotalAlloc - before.TotalAlloc) / uint64(n)
}

// index-of and last-index-of used to convert the whole string to
// []rune per call (4 bytes per rune), so a tokenizer's (index-of s c from)
// loop allocated O(n) per step. On a 30k-rune string that was ~120KB a call;
// walking the UTF-8 in place allocates only the boxed result.
func TestStringIndexFnsDoNotCopyTheString(t *testing.T) {
	s := vm.String(strings.Repeat("ab,", 10000) + "é")
	for _, name := range []string{"index-of", "last-index-of"} {
		fn := stringFn(t, name)
		args := []vm.Value{s, vm.String("é"), vm.MakeInt(3)}
		if name == "last-index-of" {
			args[2] = vm.MakeInt(30000)
		}
		if b := bytesPerCall(50, func() {
			if _, err := fn.Invoke(args); err != nil {
				t.Fatal(err)
			}
		}); b > 1024 {
			t.Errorf("string/%s on a 30k-rune string: %d bytes/call", name, b)
		}
	}
}
