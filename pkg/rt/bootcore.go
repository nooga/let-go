/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

// This file carries no bootstrap build constraint on purpose: BootCore must
// exist in every build so a downstream caller compiles, and under -tags
// bootstrap the embedded core (CoreCompiledLGB) is empty, making LoadCoreBundle
// return the documented error instead of an undefined-symbol failure.
package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// BootCore decodes and runs the embedded precompiled clojure.core bundle so a
// program that links only pkg/rt + pkg/vm — with no compiler — can resolve and
// invoke the full clojure.core surface plus the lg baseline and hybrid
// namespaces. It returns a fresh ExecContext ready to Invoke.
//
// The motivating case is an AOT-lowered binary (scripts/lg-compile output): its
// generated funcs call core via ec.Invoke(rt.CachedVarFn(..., "clojure.core",
// name), ...), which needs those vars bound. Programs that touch only NATIVE
// core fns (str, +, vector, …) already work from a bare vm.NewExecContext(),
// because installers register those at package init; BootCore is what
// additionally binds the .lg-DEFINED core (and let-go.core / let-go.types).
//
// Call once at startup. It shares the compiler-free boot spine with LoadCore
// and compiler.loadPrecompiledBundle via LoadCoreBundle; the eager hybrid
// replay is what BootCore needs over LoadCore, since it returns straight to the
// caller with no namespace loader installed.
func BootCore() (*vm.ExecContext, error) {
	if _, err := LoadCoreBundle(CoreLoadOptions{EagerHybrids: true}); err != nil {
		return nil, fmt.Errorf("BootCore: %w", err)
	}
	return vm.NewExecContext(), nil
}
