/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package gofragments_test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	_ "github.com/nooga/let-go/pkg/gofragments"
	"github.com/nooga/let-go/pkg/vm"
)

// Importing the package is what makes #go read: the entry lands in
// *data-readers*, so the reader returns the fragment body.
func TestImportInstallsTheGoReader(t *testing.T) {
	got, err := compiler.NewLispReader(strings.NewReader("#go{return 1}"), "probe.lg").Read()
	if err != nil {
		t.Fatal(err)
	}
	if got != vm.String("return 1") {
		t.Fatalf("#go{return 1} read as %v", got)
	}
}
