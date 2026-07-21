/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// A trailing comment (or #_ discard) must not clobber the last value: the
// reader surfaces those as the VOID sentinel and CompileMultiple used to
// compile-and-run it like a real form, so `42 ;; done` evaluated to VOID.
func TestTrailingCommentKeepsLastValue(t *testing.T) {
	for _, src := range []string{
		"42 ;; done",
		"42\n;; done\n",
		"42 #_(discarded)",
		"42 #?(:cljr 1)",
	} {
		v, err := Eval(src)
		if err != nil {
			t.Fatalf("Eval(%q): %v", src, err)
		}
		if n, ok := v.(vm.Int); !ok || int(n) != 42 {
			t.Fatalf("Eval(%q) = %v (%s), want 42", src, v, v.Type().Name())
		}
	}
}

// Input holding only no-value forms still evaluates to VOID, not nil — the
// REPL relies on this to echo nothing for a comment-only line.
func TestCommentOnlyInputEvaluatesToVoid(t *testing.T) {
	for _, src := range []string{
		";; just a comment",
		"#_(all discarded)",
	} {
		v, err := Eval(src)
		if err != nil {
			t.Fatalf("Eval(%q): %v", src, err)
		}
		if v != vm.VOID {
			t.Fatalf("Eval(%q) = %v (%s), want VOID", src, v, v.Type().Name())
		}
	}
}

// Skipped forms must not leave dead LOAD_CONST/POP pairs in the chunk: a
// comment between two forms compiles identically to no comment at all.
func TestVoidFormsEmitNoCode(t *testing.T) {
	compile := func(src string) []int32 {
		c := NewTransientCompiler(consts, rt.NS(rt.NameCoreNS))
		chunk, _, err := c.CompileMultiple(strings.NewReader(src))
		if err != nil {
			t.Fatalf("CompileMultiple(%q): %v", src, err)
		}
		return chunk.Code()
	}
	plain := compile("1 2")
	commented := compile("1 ;; between\n#_(dead) 2")
	if len(plain) != len(commented) {
		t.Fatalf("comments changed emitted code size: %d words vs %d", len(plain), len(commented))
	}
}
