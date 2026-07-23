/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Command lgprimgen generates pkg/rt/zz_primitives_generated.go from the
// //lg:native directives in the runtime's Go sources.
//
// It is split out from cmd/lginterop precisely so it imports NOTHING that boots
// the let-go runtime (no pkg/compiler, no pkg/rt). lginterop pulls in
// pkg/compiler for its interop-EDN mode, whose package init compiles core.lg —
// which fails mid-migration, when primitives have been moved out of
// installLangNS but the registrar that replaces them hasn't been generated yet.
// The generator must run in exactly that inconsistent state (it is what makes it
// consistent again), so it stays a pure go/ast → Go codegen tool.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nooga/let-go/internal/primgen"
)

func main() {
	srcDir := flag.String("primitives", "pkg/rt", "directory containing //lg:-annotated Go sources")
	out := flag.String("primitives-out", "pkg/rt/zz_primitives_generated.go", "output file for the generated registrar")
	goPkg := flag.String("go-pkg", "github.com/nooga/let-go/pkg/rt", "Go import path of the scanned sources")
	flag.Parse()

	if err := primgen.Generate(*srcDir, *out, *goPkg); err != nil {
		fmt.Fprintf(os.Stderr, "lgprimgen: %v\n", err)
		os.Exit(1)
	}
}
