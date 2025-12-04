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
	// no-op; test namespace is embedded and loaded by resolver on demand
	_ = vm.NIL
}
