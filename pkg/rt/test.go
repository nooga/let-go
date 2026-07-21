/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"github.com/nooga/let-go/pkg/vm"
)

// nolint
func installTestNS() {
	// no-op; test namespace is embedded and loaded by resolver on demand
	_ = vm.NIL
}
