/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	_ "embed"

	"github.com/nooga/let-go/pkg/vm"
)

//go:embed core/test.lg
var TestSrc string

// nolint
func installTestNS() {
	ns := vm.NewNamespace("test")

	// vars
	CurrentNS = ns.Def("*ns*", ns)

	RegisterNS(ns)
}
