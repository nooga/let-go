/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestLoweredOpsDevirtualized guards the lgbgen driver half of the deftype
// native lowering (T9): lower-ns-to-go captures defprotocol/deftype forms,
// emits native interface/struct/method decls, and binds the registries that
// devirtualize (->ConstOp) to a native constructor and (op-infer recv ...)
// to a direct method call — but only if the driver FORWARDS those forms.
// cmd/lgbgen's form filter used to keep only defn/defmulti forms, silently
// starving the machinery: ir.ops lowered with every protocol call as a
// boxed rt.CachedVarFn trampoline and no type decls at all.
//
// Content-based over the lowered --target=go artifact, so a filter regression
// fails here without re-running lgbgen. That tree (pkg/rt/core_go_lowered/) is
// gitignored and a clean checkout does not contain it — `go generate ./pkg/rt/`
// only rebuilds core_compiled.lgb, not the --target=go tree — so this SKIPS
// when the artifact is absent (the default -short CI lane and fresh clones)
// rather than failing on ENOENT. It runs wherever the tree has been generated
// (local dev after `make generate`); the same forwarding regression is also
// caught in CI by the -tags gogen_ir wire build (undefined symbols) and by
// check-generated.
func TestLoweredOpsDevirtualized(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	loweredOps := filepath.Join(repoRoot, "pkg", "rt", "core_go_lowered", "ir", "ops", "ops.go")
	b, err := os.ReadFile(loweredOps)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("lowered --target=go tree absent (gitignored; run `make generate`); " +
				"forwarding regressions are covered by the -tags gogen_ir wire build and check-generated")
		}
		t.Fatalf("read lowered ir/ops: %v", err)
	}
	src := string(b)

	for _, marker := range []string{
		// deftype decls emitted natively (go-export-name collapses interior
		// caps: ConstOp -> Constop). #555's ops.lg reifies only :const
		// (op-registry {:const (->ConstOp)}); further op-types (Tryop, …)
		// are added as the catalog reifies them — assert the ones present.
		"type Constop struct",
		// protocol interface emitted
		"type Irop interface",
		// protocol-method call on a known-concrete receiver devirtualized
		// to a direct method call (any receiver ident, so match loosely)
		".OpInfer(ec",
	} {
		if !strings.Contains(src, marker) {
			t.Errorf("lowered ir/ops missing %q — defprotocol/deftype forms not forwarded to lower-ns-to-go (lgbgen form filter) or devirtualization regressed", marker)
		}
	}

	// The devirtualized arms must not fall back to the boxed trampoline for
	// the protocol method itself.
	if strings.Contains(src, `"ir.ops", "op-infer"`) {
		t.Errorf("lowered ir/ops still dispatches op-infer through rt.CachedVarFn — protocol call not devirtualized")
	}
}
