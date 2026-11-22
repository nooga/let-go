/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
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

// writeBundle writes a deferrable bundle of nss to a temp file and returns its path.
func writeBundle(t *testing.T, nss []nsSource) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bundle.lgb")
	if err := os.WriteFile(path, deferrableBundle(t, nss), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// threeNS is a bundle whose middle namespace, prefix.b, is left deferred:
// prefix.first is chunk 0 and prefix.a is the main chunk.
func threeNS(prefix string, bSrc string) []nsSource {
	return []nsSource{
		{prefix + ".first", "(ns " + prefix + ".first)"},
		{prefix + ".b", "(ns " + prefix + ".b) " + bSrc},
		{prefix + ".a", "(ns " + prefix + ".a) (def y 1)"},
	}
}

func loadBundleForm(t *testing.T, path string) (vm.Value, error) {
	t.Helper()
	return runForm(initCompiler(false), fmt.Sprintf("(load-bundle %q)", path))
}

func TestLoadBundleCompiledByLgC(t *testing.T) {
	cleanupNS(t, "lbc.lib")
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "lib.lg"), filepath.Join(dir, "lib.lgb")
	if err := os.WriteFile(src, []byte("(ns lbc.lib) (def answer 42) (defn twice [n] (* 2 n))"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := initCompiler(false)
	if err := compileLG(ctx, resolver.NewNSResolver(ctx, nil), src, dst, ""); err != nil {
		t.Fatalf("compileLG: %v", err)
	}
	rt.RemoveNS("lbc.lib")

	if _, err := loadBundleForm(t, dst); err != nil {
		t.Fatalf("load-bundle: %v", err)
	}
	if got := varValue(t, "lbc.lib", "answer"); got != vm.Int(42) {
		t.Fatalf("lbc.lib/answer = %v, want 42", got)
	}
	got, err := runForm(initCompiler(false), "(lbc.lib/twice 21)")
	if err != nil || got != vm.Int(42) {
		t.Fatalf("(lbc.lib/twice 21) = %v, %v; want 42", got, err)
	}
}

func TestLoadBundleDefers(t *testing.T) {
	cleanupNS(t, "lbd.first", "lbd.b", "lbd.a")
	path := writeBundle(t, threeNS("lbd", "(def x 7)"))

	if _, err := loadBundleForm(t, path); err != nil {
		t.Fatalf("load-bundle: %v", err)
	}
	if got := varValue(t, "lbd.a", "y"); got != vm.Int(1) {
		t.Fatalf("lbd.a/y = %v, want 1", got)
	}
	if !rt.NSNeedsLoad("lbd.b") {
		t.Fatal("lbd.b is not needs-load after load-bundle")
	}
	if _, err := runForm(initCompiler(false), "(require 'lbd.b)"); err != nil {
		t.Fatalf("require: %v", err)
	}
	if got := varValue(t, "lbd.b", "x"); got != vm.Int(7) {
		t.Fatalf("lbd.b/x = %v, want 7", got)
	}
}

func TestLoadBundleTwice(t *testing.T) {
	cleanupNS(t, "lbt2.first", "lbt2.b", "lbt2.a")
	path := writeBundle(t, threeNS("lbt2", "(def x 7)"))
	if _, err := rt.LoadBundle(path); err != nil {
		t.Fatalf("first load: %v", err)
	}
	if _, err := rt.RequireNS("lbt2.b"); err != nil {
		t.Fatalf("require: %v", err)
	}
	rt.LookupNS("lbt2.a").LookupLocal("y").SetRoot(vm.NIL)

	if _, err := rt.LoadBundle(path); err != nil {
		t.Fatalf("second load: %v", err)
	}
	if got := varValue(t, "lbt2.a", "y"); got != vm.Int(1) {
		t.Fatalf("lbt2.a/y = %v after the second load, want 1", got)
	}
	if rt.NSNeedsLoad("lbt2.b") {
		t.Fatal("the second load marked lbt2.b, which was already loaded")
	}
}

func TestLoadBundleLeavesLoadedNamespace(t *testing.T) {
	cleanupNS(t, "lbl.first", "lbl.b", "lbl.a")
	path := writeBundle(t, threeNS("lbl", "(def x :bundle)"))
	if _, err := runForm(initCompiler(false), "(ns lbl.b) (def x :source)"); err != nil {
		t.Fatal(err)
	}

	if _, err := rt.LoadBundle(path); err != nil {
		t.Fatalf("load: %v", err)
	}
	if rt.NSNeedsLoad("lbl.b") {
		t.Fatal("load-bundle marked lbl.b, which source had loaded")
	}
	if _, err := rt.RequireNS("lbl.b"); err != nil {
		t.Fatalf("require: %v", err)
	}
	if got := varValue(t, "lbl.b", "x"); got != vm.Keyword("source") {
		t.Fatalf("lbl.b/x = %v, want :source", got)
	}
}

func TestLoadBundleRestoresNS(t *testing.T) {
	cleanupNS(t, "lbn.first", "lbn.b", "lbn.a")
	path := writeBundle(t, threeNS("lbn", "(def x 7)"))
	user := rt.NS("user")
	prev := rt.CurrentNS.Deref()
	rt.CurrentNS.SetRoot(user)
	t.Cleanup(func() { rt.CurrentNS.SetRoot(prev) })

	if _, err := rt.LoadBundle(path); err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := rt.CurrentNS.Deref(); got != user {
		t.Fatalf("*ns* = %v after load-bundle, want user", got)
	}
	if _, err := rt.RequireNS("lbn.b"); err != nil {
		t.Fatalf("require: %v", err)
	}
	if got := rt.CurrentNS.Deref(); got != user {
		t.Fatalf("*ns* = %v after require, want user", got)
	}
}

func TestLoadBundleBadPath(t *testing.T) {
	dir := t.TempDir()
	empty := filepath.Join(dir, "empty.lgb")
	junk := filepath.Join(dir, "junk.lgb")
	for p, body := range map[string]string{empty: "", junk: "not an lgb"} {
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for name, p := range map[string]string{
		"missing":   filepath.Join(dir, "missing.lgb"),
		"directory": dir,
		"empty":     empty,
		"not lgb":   junk,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := loadBundleForm(t, p)
			if err == nil || !strings.Contains(err.Error(), "load-bundle") || !strings.Contains(err.Error(), p) {
				t.Fatalf("err = %v, want one naming load-bundle and %s", err, p)
			}
		})
	}
}

func TestLoadBundleMainRequiresDeferred(t *testing.T) {
	cleanupNS(t, "lbr.first", "lbr.b", "lbr.a")
	withoutNSLoader(t)
	path := writeBundle(t, []nsSource{
		{"lbr.first", "(ns lbr.first)"},
		{"lbr.b", "(ns lbr.b) (def x 7)"},
		{"lbr.a", "(ns lbr.a (:require [lbr.b])) (def y lbr.b/x)"},
	})
	if _, err := rt.LoadBundle(path); err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := varValue(t, "lbr.a", "y"); got != vm.Int(7) {
		t.Fatalf("lbr.a/y = %v, want 7", got)
	}
}

func TestLoadBundleReturnsMainValue(t *testing.T) {
	cleanupNS(t, "lbv.first", "lbv.b", "lbv.a")
	path := writeBundle(t, []nsSource{
		{"lbv.first", "(ns lbv.first)"},
		{"lbv.b", "(ns lbv.b) (def x 7)"},
		{"lbv.a", "(ns lbv.a) (def y 1) :main-value"},
	})
	got, err := rt.LoadBundle(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got != vm.Keyword("main-value") {
		t.Fatalf("load-bundle returned %v, want :main-value", got)
	}
}
