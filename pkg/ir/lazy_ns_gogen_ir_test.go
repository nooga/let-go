//go:build gogen_ir

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

const lazyNSChildEnv = "LG_TEST_LAZY_NS_REQUIRE_CHILD"

// Requiring a lazily loaded bundled namespace whose lowered Go overrides are
// linked in must not redefine its own names over its clojure.core refers: the
// overrides have to wait for the namespace's load. Each require runs in a
// fresh child process, so no other test can have loaded the namespace (and
// printed its warnings) first.
func TestRequiringALazyNamespacePrintsNoShadowWarnings(t *testing.T) {
	if os.Getenv(lazyNSChildEnv) != "" {
		requireInThisProcess(t, os.Getenv(lazyNSChildEnv))
		return
	}
	for _, ns := range []string{"edn", "walk", "set"} {
		t.Run(ns, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestRequiringALazyNamespacePrintsNoShadowWarnings$")
			cmd.Env = append(os.Environ(), lazyNSChildEnv+"="+ns)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("child requiring %s failed: %v\n%s", ns, err, out)
			}
			if strings.Contains(string(out), "already refers to") {
				t.Fatalf("(require '%s) printed shadow warnings:\n%s", ns, out)
			}
		})
	}
}

func requireInThisProcess(t *testing.T, ns string) {
	consts := vm.NewConsts()
	ctx := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	rt.SetNSLoader(resolver.NewNSResolver(ctx, []string{"."}))
	c := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	c.SetSource("lazy-ns-require-child")
	if _, _, err := c.CompileMultiple(strings.NewReader("(require '" + ns + ")")); err != nil {
		t.Fatalf("(require '%s): %v", ns, err)
	}
}
