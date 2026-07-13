/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestFormatErrorIncludesEveryCompileChainSource(t *testing.T) {
	const sourceName = "diagnostic-source.lg"
	const source = "(def broken\n  (fn []\n    (let [:tag 1]\n      1)))"
	compiler := NewCompiler(vm.NewConsts(), rt.CoreNS).SetSource(sourceName)

	_, err := compiler.Compile(source)
	if err == nil {
		t.Fatal("Compile unexpectedly succeeded")
	}
	rendered := vm.FormatError(err)
	for _, location := range []string{
		sourceName + ":1:",
		sourceName + ":2:",
		sourceName + ":3:",
	} {
		if !strings.Contains(rendered, location) {
			t.Errorf("formatted error missing %s:\n%s", location, rendered)
		}
	}
	if !strings.Contains(rendered, "compiling def value") ||
		!strings.Contains(rendered, "compiling fn body") ||
		!strings.Contains(rendered, "let binding name must be a symbol") {
		t.Fatalf("formatted error lost compile-chain messages:\n%s", rendered)
	}
}
