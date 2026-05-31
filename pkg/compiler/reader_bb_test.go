/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func evalStr(t *testing.T, s string) vm.Value {
	t.Helper()
	cp := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := NewCompiler(cp, ns)
	ch, _, err := ctx.CompileMultiple(strings.NewReader(s))
	if err != nil {
		t.Fatalf("compile %q: %v", s, err)
	}
	v, err := vm.NewFrame(ch, nil).Run()
	if err != nil {
		t.Fatalf("run %q: %v", s, err)
	}
	return v
}

// The :bb reader-conditional feature is opt-in (off by default, like :clj), so
// libraries shipping babashka-compatible fallbacks can be read on let-go.
func TestReaderConditional_BB(t *testing.T) {
	defer SetMatchBbConditional(false)

	// Off: :bb is skipped, :default wins.
	SetMatchBbConditional(false)
	if v := evalStr(t, `#?(:bb 1 :default 2)`); v != vm.Int(2) {
		t.Fatalf(":bb off → want 2 (:default), got %v", v)
	}

	// On: :bb matches.
	SetMatchBbConditional(true)
	if v := evalStr(t, `#?(:bb 1 :default 2)`); v != vm.Int(1) {
		t.Fatalf(":bb on → want 1 (:bb), got %v", v)
	}

	// First-match semantics: :bb (listed first) wins even with :clj off.
	if v := evalStr(t, `#?(:bb 10 :clj 20)`); v != vm.Int(10) {
		t.Fatalf(":bb-first → want 10, got %v", v)
	}
}

// The set-read-bb! runtime fn toggles the same flag.
func TestSetReadBb_Wiring(t *testing.T) {
	defer SetMatchBbConditional(false)
	SetMatchBbConditional(false)
	evalStr(t, `(set-read-bb! true)`)
	if !matchBbConditional {
		t.Fatal("(set-read-bb! true) should enable :bb matching")
	}
	evalStr(t, `(set-read-bb! false)`)
	if matchBbConditional {
		t.Fatal("(set-read-bb! false) should disable :bb matching")
	}
}
