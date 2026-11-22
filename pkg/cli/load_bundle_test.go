/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type nsSource struct{ name, src string }

// deferrableBundle compiles each namespace into one shared const pool and
// encodes them with their pool ranges, as lgbgen does, so a namespace whose
// functions only its own code reaches can be left undecoded. The first one is
// never deferred (its chunk is chunk 0) and the last one is the main chunk.
func deferrableBundle(t *testing.T, nss []nsSource) []byte {
	t.Helper()
	consts := vm.NewConsts()
	chunks := map[string]*vm.CodeChunk{}
	ranges := map[string]bytecode.ConstRange{}
	order := make([]string, 0, len(nss))
	for _, n := range nss {
		c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
		lo := len(consts.Values())
		chunk, _, err := c.CompileMultiple(strings.NewReader(n.src))
		if err != nil {
			t.Fatalf("compiling %s: %v", n.name, err)
		}
		chunks[n.name] = chunk
		ranges[n.name] = bytecode.ConstRange{Lo: lo, Hi: len(consts.Values())}
		order = append(order, n.name)
	}
	var buf bytes.Buffer
	if err := bytecode.EncodeBundleDeferrable(&buf, consts, chunks, order, ranges, false); err != nil {
		t.Fatalf("encoding: %v", err)
	}
	// Compiling ran each namespace; drop them so the load starts as a fresh
	// process would, with nothing loaded.
	for _, n := range nss {
		rt.RemoveNS(n.name)
	}
	return buf.Bytes()
}

// cleanupNS removes test namespaces from the global registry when the test ends.
func cleanupNS(t *testing.T, names ...string) {
	t.Helper()
	t.Cleanup(func() {
		for _, n := range names {
			rt.RemoveNS(n)
		}
	})
}

// withoutNSLoader runs the rest of the test with no namespace loader installed.
func withoutNSLoader(t *testing.T) {
	t.Helper()
	prev := rt.GetNSLoader()
	rt.SetNSLoader(nil)
	t.Cleanup(func() { rt.SetNSLoader(prev) })
}

func resetRoot(t *testing.T, ns, name string) {
	t.Helper()
	v := rt.NS(ns).LookupLocal(vm.Symbol(name))
	if v == nil {
		t.Fatalf("%s/%s not defined", ns, name)
	}
	v.SetRoot(vm.NIL)
}

func varValue(t *testing.T, ns, name string) vm.Value {
	t.Helper()
	n := rt.LookupNS(ns)
	if n == nil {
		t.Fatalf("namespace %s not registered", ns)
	}
	v := n.LookupLocal(vm.Symbol(name))
	if v == nil {
		t.Fatalf("%s/%s not defined", ns, name)
	}
	return v.Deref()
}

func decodeDeferringAll(t *testing.T, data []byte) *bytecode.ExecUnit {
	t.Helper()
	unit, err := bytecode.DecodeBundle(data, rt.LGBVarResolver, func(string) bool { return true })
	if err != nil {
		t.Fatalf("decoding: %v", err)
	}
	return unit
}

func TestBundleTableRequireWithoutLoader(t *testing.T) {
	cleanupNS(t, "lbt.first", "lbt.b", "lbt.a")
	data := deferrableBundle(t, []nsSource{
		{"lbt.first", "(ns lbt.first)"},
		{"lbt.b", "(ns lbt.b) (def x 42)"},
		{"lbt.a", "(ns lbt.a) (def y 1)"},
	})
	unit := decodeDeferringAll(t, data)
	if !unit.DeferredNS("lbt.b") {
		t.Fatal("lbt.b was decoded eagerly")
	}
	rt.RegisterBundleNamespaces(unit, []string{"lbt.b"})
	withoutNSLoader(t)

	if _, err := rt.RequireNS("lbt.b"); err != nil {
		t.Fatalf("RequireNS(lbt.b): %v", err)
	}
	if got := varValue(t, "lbt.b", "x"); got != vm.Int(42) {
		t.Fatalf("lbt.b/x = %v, want 42", got)
	}
	if unit.DeferredNS("lbt.b") {
		t.Fatal("lbt.b still deferred after require")
	}
	if rt.NSNeedsLoad("lbt.b") {
		t.Fatal("lbt.b still needs-load after require")
	}
}

func TestBundleTableFailureRetries(t *testing.T) {
	cleanupNS(t, "lbt.ctl", "lbt.first", "lbt.f", "lbt.g")
	// Compiling runs top-level forms, so the chunk throws only once the
	// switch is on.
	fail := rt.NS("lbt.ctl").Def("fail", vm.FALSE)
	data := deferrableBundle(t, []nsSource{
		{"lbt.first", "(ns lbt.first)"},
		{"lbt.f", `(ns lbt.f) (def z (when lbt.ctl/fail (throw (ex-info "boom" {}))))`},
		{"lbt.g", "(ns lbt.g) (def y 1)"},
	})
	unit := decodeDeferringAll(t, data)
	if !unit.DeferredNS("lbt.f") {
		t.Fatal("lbt.f was decoded eagerly")
	}
	rt.RegisterBundleNamespaces(unit, []string{"lbt.f"})
	withoutNSLoader(t)
	fail.SetRoot(vm.TRUE)

	for i := 1; i <= 2; i++ {
		_, err := rt.RequireNS("lbt.f")
		if err == nil || !strings.Contains(err.Error(), "boom") {
			t.Fatalf("require %d: err = %v, want boom", i, err)
		}
		if !rt.NSNeedsLoad("lbt.f") {
			t.Fatalf("require %d: lbt.f no longer needs-load after a failed load", i)
		}
	}
}
